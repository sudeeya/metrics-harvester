package router

import (
	"net/http"

	repo "github.com/sudeeya/metrics-harvester/internal/repository"
)

type Router struct {
	repository repo.Repository
	mux        *http.ServeMux
}

func NewRouter(repository repo.Repository) *Router {
	return &Router{
		repository: repository,
		mux:        http.NewServeMux(),
	}
}

func (router *Router) HandleFunc(pattern string, handler func(w http.ResponseWriter, r *http.Request)) {
	router.mux.HandleFunc(pattern, handler)
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.mux.ServeHTTP(w, r)
}
