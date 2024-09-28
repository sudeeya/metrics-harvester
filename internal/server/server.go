package server

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "net/http/pprof"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/sudeeya/metrics-harvester/internal/handlers"
	"github.com/sudeeya/metrics-harvester/internal/metric"
	"github.com/sudeeya/metrics-harvester/internal/middleware"
	repo "github.com/sudeeya/metrics-harvester/internal/repository"
	"github.com/sudeeya/metrics-harvester/internal/repository/database"
	"github.com/sudeeya/metrics-harvester/internal/repository/storage"
)

const limitInSeconds = 10

type Server struct {
	cfg              *Config
	logger           *zap.Logger
	repository       repo.Repository
	handler          http.Handler
	privateKey       *rsa.PrivateKey
	symmetricKeyChan chan []byte
	symmetricKey     *[]byte
}

func NewServer(logger *zap.Logger, cfg *Config, repository repo.Repository) *Server {
	logger.Info("Initializing storage file")
	initializeStorageFile(logger, cfg)

	logger.Info("Initializing repository")
	initializeRepository(logger, cfg, repository)

	logger.Info("Extracting private key from file")
	privateKey := extractPrivateKey(logger, cfg.CryptoKey)

	symmetricKeyChan := make(chan []byte, 1)
	symmetricKey := []byte("")

	router := chi.NewRouter()
	logger.Info("Initializing routes")
	addRoutes(logger, repository, router, privateKey, symmetricKeyChan)

	logger.Info("Initializing middleware")
	handler := middleware.WithCompressing(router)
	handler = middleware.WithDecryption(&symmetricKey, handler)
	handler = middleware.WithSigning([]byte(cfg.Key), handler)
	handler = middleware.WithLogging(logger, handler)
	return &Server{
		cfg:              cfg,
		logger:           logger,
		repository:       repository,
		handler:          handler,
		privateKey:       privateKey,
		symmetricKeyChan: symmetricKeyChan,
		symmetricKey:     &symmetricKey,
	}
}

func initializeStorageFile(logger *zap.Logger, cfg *Config) {
	if file, err := os.Open(cfg.FileStoragePath); !os.IsNotExist(err) {
		file.Close()
		return
	}
	newFile, err := os.Create(cfg.FileStoragePath)
	if err != nil {
		logger.Fatal(err.Error())
	}
	newFile.Close()
	if err := os.WriteFile(cfg.FileStoragePath, []byte("[]"), 0666); err != nil {
		logger.Fatal(err.Error())
	}
}

func initializeRepository(logger *zap.Logger, cfg *Config, repository repo.Repository) {
	switch v := repository.(type) {
	case *database.Database:
		logger.Info("Initializing database")
		ctx, cancel := context.WithTimeout(context.Background(), limitInSeconds*time.Second)
		defer cancel()
		if _, err := v.DB.ExecContext(ctx, database.CreateMetricsTable); err != nil {
			logger.Fatal(err.Error())
		}
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			logger.Error(ctx.Err().Error())
		}
	case *storage.MemStorage:
		if cfg.Restore {
			logger.Info("Initializing metrics with saved values from a file")
			savedData, err := os.ReadFile(cfg.FileStoragePath)
			if err != nil {
				logger.Fatal(err.Error())
			}
			var savedMetrics []metric.Metric
			if err := json.Unmarshal(savedData, &savedMetrics); err != nil {
				logger.Fatal(err.Error())
			}
			ctx := context.Background()
			for _, m := range savedMetrics {
				repository.PutMetric(ctx, m)
			}
		}
	}
}

func addRoutes(logger *zap.Logger, repository repo.Repository, router chi.Router, privateKey *rsa.PrivateKey, symmetricKey chan []byte) {
	router.Get("/value/{metricType}/{metricName}", handlers.NewValueHandler(logger, repository))
	router.Get("/ping", handlers.NewPingHandler(logger, repository))
	router.Get("/", handlers.NewAllMetricsHandler(logger, repository))
	router.Post("/update/{metricType}/{metricName}/{metricValue}", handlers.NewUpdateHandler(logger, repository))
	router.Post("/update/{metricType}/", http.NotFound)
	router.Post("/update/", handlers.NewJSONUpdateHandler(logger, repository))
	router.Post("/updates/", handlers.NewBatchHandler(logger, repository))
	router.Post("/value/", handlers.NewJSONValueHandler(logger, repository))
	router.Post("/key/", handlers.NewKeyHandler(logger, privateKey, symmetricKey))
}

func extractPrivateKey(logger *zap.Logger, file string) *rsa.PrivateKey {
	if file == "" {
		logger.Info("Private key was not given. The server is running without decryption")
		return nil
	}
	pemData, err := os.ReadFile(file)
	if err != nil {
		logger.Fatal(err.Error())
	}
	block, _ := pem.Decode(pemData)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		logger.Sugar().Fatalf("PEM file contains %s, not RSA PRIVATE KEY", block.Type)
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		logger.Fatal(err.Error())
	}
	return privateKey
}

func (s *Server) Run() {
	s.logger.Info("Server is running")
	storeTicker := time.NewTicker(time.Duration(s.cfg.StoreInterval) * time.Second)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		symmetricKey := <-s.symmetricKeyChan
		*s.symmetricKey = symmetricKey
	}()
	go func() {
		if err := http.ListenAndServe(s.cfg.Address, s.handler); err != nil {
			s.logger.Fatal(err.Error())
		}
	}()
	go func() {
		for range storeTicker.C {
			s.logger.Info("Storing all metrics to file")
			s.StoreMetricsToFile()
		}
	}()
	go func() {
		<-sigChan
		s.logger.Info("Server is shutting down")
		s.Shutdown()
	}()
	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", s.cfg.ProfilerPort), nil); err != nil {
			s.logger.Fatal(err.Error())
		}
	}()
	select {}
}

func (s *Server) StoreMetricsToFile() {
	ctx, cancel := context.WithTimeout(context.Background(), limitInSeconds*time.Second)
	defer cancel()
	metrics, err := s.repository.GetAllMetrics(ctx)
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		s.logger.Error(ctx.Err().Error())
	}
	if err != nil {
		s.logger.Fatal(err.Error())
	}
	data, err := json.MarshalIndent(metrics, "", "\t")
	if err != nil {
		s.logger.Fatal(err.Error())
	}
	if err := os.WriteFile(s.cfg.FileStoragePath, data, 0666); err != nil {
		s.logger.Fatal(err.Error())
	}
}

func (s *Server) Shutdown() {
	s.StoreMetricsToFile()
	if err := s.repository.Close(); err != nil {
		s.logger.Fatal(err.Error())
	}
	os.Exit(0)
}
