package server

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/sudeeya/metrics-harvester/internal/handlers"
	"github.com/sudeeya/metrics-harvester/internal/metric"
	"github.com/sudeeya/metrics-harvester/internal/middleware"
	repo "github.com/sudeeya/metrics-harvester/internal/repository"
	"github.com/sudeeya/metrics-harvester/internal/repository/database"
	"github.com/sudeeya/metrics-harvester/internal/repository/storage"
	"go.uber.org/zap"
)

type Server struct {
	cfg        *Config
	logger     *zap.Logger
	repository repo.Repository
	handler    http.Handler
}

func NewServer(logger *zap.Logger, cfg *Config, repository repo.Repository) *Server {
	ctx := context.Background()
	logger.Info("Initializing storage file")
	initializeStorageFile(logger, cfg)
	logger.Info("Initializing repository")
	initializeRepository(ctx, logger, cfg, repository)
	router := chi.NewRouter()
	logger.Info("Initializing routes")
	addRoutes(ctx, logger, repository, router)
	logger.Info("Initializing middleware")
	handler := middleware.WithCompressing(router)
	handler = middleware.WithSigning([]byte(cfg.Key), handler)
	handler = middleware.WithLogging(logger, handler)
	return &Server{
		cfg:        cfg,
		logger:     logger,
		repository: repository,
		handler:    handler,
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

func initializeRepository(ctx context.Context, logger *zap.Logger, cfg *Config, repository repo.Repository) {
	switch v := repository.(type) {
	case *database.Database:
		logger.Info("Initializing database")
		if _, err := v.DB.ExecContext(ctx, database.CreateMetricsTable); err != nil {
			logger.Fatal(err.Error())
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
			for _, m := range savedMetrics {
				repository.PutMetric(ctx, m)
			}
		} else {
			logger.Info("Initializing nessessory metrics")
			initializeDefault(ctx, repository)
		}
	}
}

func initializeDefault(ctx context.Context, repository repo.Repository) {
	repository.PutMetric(ctx, metric.Metric{ID: "Alloc", MType: metric.Gauge, Value: new(float64)})
	repository.PutMetric(ctx, metric.Metric{ID: "BuckHashSys", MType: metric.Gauge, Value: new(float64)})
	repository.PutMetric(ctx, metric.Metric{ID: "Frees", MType: metric.Gauge, Value: new(float64)})
	repository.PutMetric(ctx, metric.Metric{ID: "GCCPUFraction", MType: metric.Gauge, Value: new(float64)})
	repository.PutMetric(ctx, metric.Metric{ID: "GCSys", MType: metric.Gauge, Value: new(float64)})
	repository.PutMetric(ctx, metric.Metric{ID: "HeapAlloc", MType: metric.Gauge, Value: new(float64)})
	repository.PutMetric(ctx, metric.Metric{ID: "HeapIdle", MType: metric.Gauge, Value: new(float64)})
	repository.PutMetric(ctx, metric.Metric{ID: "HeapInuse", MType: metric.Gauge, Value: new(float64)})
	repository.PutMetric(ctx, metric.Metric{ID: "HeapObjects", MType: metric.Gauge, Value: new(float64)})
	repository.PutMetric(ctx, metric.Metric{ID: "HeapReleased", MType: metric.Gauge, Value: new(float64)})
	repository.PutMetric(ctx, metric.Metric{ID: "HeapSys", MType: metric.Gauge, Value: new(float64)})
	repository.PutMetric(ctx, metric.Metric{ID: "LastGC", MType: metric.Gauge, Value: new(float64)})
	repository.PutMetric(ctx, metric.Metric{ID: "Lookups", MType: metric.Gauge, Value: new(float64)})
	repository.PutMetric(ctx, metric.Metric{ID: "MCacheInuse", MType: metric.Gauge, Value: new(float64)})
	repository.PutMetric(ctx, metric.Metric{ID: "MCacheSys", MType: metric.Gauge, Value: new(float64)})
	repository.PutMetric(ctx, metric.Metric{ID: "MSpanInuse", MType: metric.Gauge, Value: new(float64)})
	repository.PutMetric(ctx, metric.Metric{ID: "MSpanSys", MType: metric.Gauge, Value: new(float64)})
	repository.PutMetric(ctx, metric.Metric{ID: "Mallocs", MType: metric.Gauge, Value: new(float64)})
	repository.PutMetric(ctx, metric.Metric{ID: "NextGC", MType: metric.Gauge, Value: new(float64)})
	repository.PutMetric(ctx, metric.Metric{ID: "NumForcedGC", MType: metric.Gauge, Value: new(float64)})
	repository.PutMetric(ctx, metric.Metric{ID: "NumGC", MType: metric.Gauge, Value: new(float64)})
	repository.PutMetric(ctx, metric.Metric{ID: "OtherSys", MType: metric.Gauge, Value: new(float64)})
	repository.PutMetric(ctx, metric.Metric{ID: "PauseTotalNs", MType: metric.Gauge, Value: new(float64)})
	repository.PutMetric(ctx, metric.Metric{ID: "PollCount", MType: metric.Counter, Delta: new(int64)})
	repository.PutMetric(ctx, metric.Metric{ID: "RandomValue", MType: metric.Gauge, Value: new(float64)})
	repository.PutMetric(ctx, metric.Metric{ID: "StackInuse", MType: metric.Gauge, Value: new(float64)})
	repository.PutMetric(ctx, metric.Metric{ID: "StackSys", MType: metric.Gauge, Value: new(float64)})
	repository.PutMetric(ctx, metric.Metric{ID: "Sys", MType: metric.Gauge, Value: new(float64)})
	repository.PutMetric(ctx, metric.Metric{ID: "TotalAlloc", MType: metric.Gauge, Value: new(float64)})
}

func addRoutes(ctx context.Context, logger *zap.Logger, repository repo.Repository, router chi.Router) {
	router.Get("/value/{metricType}/{metricName}", handlers.NewValueHandler(ctx, logger, repository))
	router.Get("/ping", handlers.NewPingHandler(ctx, logger, repository))
	router.Get("/", handlers.NewAllMetricsHandler(ctx, logger, repository))
	router.Post("/update/{metricType}/{metricName}/{metricValue}", handlers.NewUpdateHandler(ctx, logger, repository))
	router.Post("/update/{metricType}/", http.NotFound)
	router.Post("/update/", handlers.NewJSONUpdateHandler(ctx, logger, repository))
	router.Post("/updates/", handlers.NewBatchHandler(ctx, logger, repository))
	router.Post("/value/", handlers.NewJSONValueHandler(ctx, logger, repository))
	router.Post("/", handlers.BadRequest)
}

func (s *Server) Run() {
	s.logger.Info("Server is running")
	storeTicker := time.NewTicker(time.Duration(s.cfg.StoreInterval) * time.Second)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
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
	select {}
}

func (s *Server) StoreMetricsToFile() {
	metrics, _ := s.repository.GetAllMetrics(context.Background())
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
