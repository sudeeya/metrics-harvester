package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sudeeya/metrics-harvester/internal/handlers"
	repo "github.com/sudeeya/metrics-harvester/internal/repository"
)

type Server struct {
	cfg        *Config
	repository repo.Repository
	router     chi.Router
}

func NewServer(cfg *Config, repository repo.Repository) *Server {
	router := chi.NewRouter()
	addRoutes(repository, router)
	return &Server{
		cfg:        cfg,
		repository: repository,
		router:     router,
	}
}

func addRoutes(repository repo.Repository, router chi.Router) {
	router.Get("/value/{metricType}/{metricName}", handlers.CreateGetMetricHandler(repository))
	router.Get("/", handlers.CreateGetAllMetricsHandler(repository))
	router.Post("/update/{metricType}/{metricName}/{metricValue}", handlers.CreatePostMetricHandler(repository))
	router.Post("/update/{metricType}/", http.NotFound)
	router.Post("/", handlers.BadRequest)
}

func (s *Server) Run() {
	if err := http.ListenAndServe(s.cfg.Address, s.router); err != nil {
		panic(err)
	}
}
