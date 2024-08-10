package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sudeeya/metrics-harvester/internal/handlers"
	log "github.com/sudeeya/metrics-harvester/internal/logger"
	"github.com/sudeeya/metrics-harvester/internal/metric"
	repo "github.com/sudeeya/metrics-harvester/internal/repository"
	"go.uber.org/zap"
)

type Server struct {
	cfg        *Config
	logger     *zap.Logger
	repository repo.Repository
	handler    http.Handler
}

func NewServer(cfg *Config, logger *zap.Logger, repository repo.Repository) *Server {
	router := chi.NewRouter()
	addRoutes(logger, repository, router)
	handler := log.WithLogging(logger, router)
	initializeMetrics(repository)
	return &Server{
		cfg:        cfg,
		logger:     logger,
		repository: repository,
		handler:    handler,
	}
}

func addRoutes(logger *zap.Logger, repository repo.Repository, router chi.Router) {
	router.Get("/value/{metricType}/{metricName}", handlers.NewValueHandler(logger, repository))
	router.Get("/", handlers.NewAllMetricsHandler(logger, repository))
	router.Post("/update/{metricType}/{metricName}/{metricValue}", handlers.NewUpdateHandler(logger, repository))
	router.Post("/update/{metricType}/", http.NotFound)
	router.Post("/update/", handlers.NewJSONUpdateHandler(logger, repository))
	router.Post("/value/", handlers.NewJSONValueHandler(logger, repository))
	router.Post("/", handlers.BadRequest)
}

func initializeMetrics(repository repo.Repository) {
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

func (s *Server) Run() {
	if err := http.ListenAndServe(s.cfg.Address, s.handler); err != nil {
		s.logger.Fatal(err.Error())
	}
}
