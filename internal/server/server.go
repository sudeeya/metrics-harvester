package server

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sudeeya/metrics-harvester/internal/handlers"
	"github.com/sudeeya/metrics-harvester/internal/metric"
	"github.com/sudeeya/metrics-harvester/internal/middleware"
	repo "github.com/sudeeya/metrics-harvester/internal/repository"
	"go.uber.org/zap"
)

type Server struct {
	cfg        *Config
	db         *sql.DB
	logger     *zap.Logger
	repository repo.Repository
	handler    http.Handler
}

func NewServer(logger *zap.Logger, cfg *Config, repository repo.Repository) *Server {
	logger.Info("Initializing storage file")
	initializeStorageFile(logger, cfg)
	logger.Info("Initializing repository")
	initializeMetrics(logger, cfg, repository)

	var db *sql.DB
	var err error
	if cfg.DatabaseDSN != "" {
		db, err = sql.Open("pgx", cfg.DatabaseDSN)
		if err != nil {
			log.Fatal(err)
		}
	}

	router := chi.NewRouter()
	logger.Info("Initializing routes")
	addRoutes(logger, db, repository, router)
	logger.Info("Initializing middleware")
	handler := middleware.WithCompressing(router)
	handler = middleware.WithLogging(logger, handler)
	return &Server{
		cfg:        cfg,
		db:         db,
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

func initializeMetrics(logger *zap.Logger, cfg *Config, repository repo.Repository) {
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
			repository.PutMetric(m)
		}
	} else {
		logger.Info("Initializing nessessory metrics")
		initializeDefault(repository)
	}
}

func initializeDefault(repository repo.Repository) {
	repository.PutMetric(metric.Metric{ID: "Alloc", MType: metric.Gauge, Value: new(float64)})
	repository.PutMetric(metric.Metric{ID: "BuckHashSys", MType: metric.Gauge, Value: new(float64)})
	repository.PutMetric(metric.Metric{ID: "Frees", MType: metric.Gauge, Value: new(float64)})
	repository.PutMetric(metric.Metric{ID: "GCCPUFraction", MType: metric.Gauge, Value: new(float64)})
	repository.PutMetric(metric.Metric{ID: "GCSys", MType: metric.Gauge, Value: new(float64)})
	repository.PutMetric(metric.Metric{ID: "HeapAlloc", MType: metric.Gauge, Value: new(float64)})
	repository.PutMetric(metric.Metric{ID: "HeapIdle", MType: metric.Gauge, Value: new(float64)})
	repository.PutMetric(metric.Metric{ID: "HeapInuse", MType: metric.Gauge, Value: new(float64)})
	repository.PutMetric(metric.Metric{ID: "HeapObjects", MType: metric.Gauge, Value: new(float64)})
	repository.PutMetric(metric.Metric{ID: "HeapReleased", MType: metric.Gauge, Value: new(float64)})
	repository.PutMetric(metric.Metric{ID: "HeapSys", MType: metric.Gauge, Value: new(float64)})
	repository.PutMetric(metric.Metric{ID: "LastGC", MType: metric.Gauge, Value: new(float64)})
	repository.PutMetric(metric.Metric{ID: "Lookups", MType: metric.Gauge, Value: new(float64)})
	repository.PutMetric(metric.Metric{ID: "MCacheInuse", MType: metric.Gauge, Value: new(float64)})
	repository.PutMetric(metric.Metric{ID: "MCacheSys", MType: metric.Gauge, Value: new(float64)})
	repository.PutMetric(metric.Metric{ID: "MSpanInuse", MType: metric.Gauge, Value: new(float64)})
	repository.PutMetric(metric.Metric{ID: "MSpanSys", MType: metric.Gauge, Value: new(float64)})
	repository.PutMetric(metric.Metric{ID: "Mallocs", MType: metric.Gauge, Value: new(float64)})
	repository.PutMetric(metric.Metric{ID: "NextGC", MType: metric.Gauge, Value: new(float64)})
	repository.PutMetric(metric.Metric{ID: "NumForcedGC", MType: metric.Gauge, Value: new(float64)})
	repository.PutMetric(metric.Metric{ID: "NumGC", MType: metric.Gauge, Value: new(float64)})
	repository.PutMetric(metric.Metric{ID: "OtherSys", MType: metric.Gauge, Value: new(float64)})
	repository.PutMetric(metric.Metric{ID: "PauseTotalNs", MType: metric.Gauge, Value: new(float64)})
	repository.PutMetric(metric.Metric{ID: "PollCount", MType: metric.Counter, Delta: new(int64)})
	repository.PutMetric(metric.Metric{ID: "RandomValue", MType: metric.Gauge, Value: new(float64)})
	repository.PutMetric(metric.Metric{ID: "StackInuse", MType: metric.Gauge, Value: new(float64)})
	repository.PutMetric(metric.Metric{ID: "StackSys", MType: metric.Gauge, Value: new(float64)})
	repository.PutMetric(metric.Metric{ID: "Sys", MType: metric.Gauge, Value: new(float64)})
	repository.PutMetric(metric.Metric{ID: "TotalAlloc", MType: metric.Gauge, Value: new(float64)})
}

func addRoutes(logger *zap.Logger, db *sql.DB, repository repo.Repository, router chi.Router) {
	router.Get("/value/{metricType}/{metricName}", handlers.NewValueHandler(logger, repository))
	router.Get("/ping", handlers.NewPingHandler(logger, db))
	router.Get("/", handlers.NewAllMetricsHandler(logger, repository))
	router.Post("/update/{metricType}/{metricName}/{metricValue}", handlers.NewUpdateHandler(logger, repository))
	router.Post("/update/{metricType}/", http.NotFound)
	router.Post("/update/", handlers.NewJSONUpdateHandler(logger, repository))
	router.Post("/value/", handlers.NewJSONValueHandler(logger, repository))
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
		s.StoreMetricsToFile()
		if err := s.db.Close(); err != nil {
			s.logger.Fatal(err.Error())
		}
		os.Exit(0)
	}()
	select {}
}

func (s *Server) StoreMetricsToFile() {
	metrics, _ := s.repository.GetAllMetrics()
	data, err := json.MarshalIndent(metrics, "", "\t")
	if err != nil {
		s.logger.Fatal(err.Error())
	}
	if err := os.WriteFile(s.cfg.FileStoragePath, data, 0666); err != nil {
		s.logger.Fatal(err.Error())
	}
}
