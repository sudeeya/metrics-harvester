package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sudeeya/metrics-harvester/internal/handlers"
	log "github.com/sudeeya/metrics-harvester/internal/logger"
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

func (s *Server) Run() {
	if err := http.ListenAndServe(s.cfg.Address, s.handler); err != nil {
		s.logger.Fatal(err.Error())
	}
}
