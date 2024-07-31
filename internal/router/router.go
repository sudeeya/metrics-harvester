package router

import (
	"github.com/go-chi/chi/v5"
	repo "github.com/sudeeya/metrics-harvester/internal/repository"
)

type Router struct {
	chi.Router
	repo.Repository
}

func NewRouter(repository repo.Repository) *Router {
	return &Router{
		Router:     chi.NewRouter(),
		Repository: repository,
	}
}
