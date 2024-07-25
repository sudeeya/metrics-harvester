package router

import (
	"net/http"
	"strconv"
	"strings"

	repo "github.com/sudeeya/metrics-harvester/internal/repository"
)

type Router struct {
	repository repo.Repository
	mux        *http.ServeMux
}

func NewRouter(repository repo.Repository) *Router {
	mux := http.NewServeMux()
	mux.HandleFunc("/update/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			path := strings.Trim(r.URL.Path[len("/update/"):], "/")
			splitPath := strings.Split(path, "/")
			switch splitPath[0] {
			case "gauge":
				if len(splitPath) == 1 {
					w.WriteHeader(http.StatusNotFound)
					return
				} else if len(splitPath) != 3 {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				metric, err := strconv.ParseFloat(splitPath[2], 64)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				repository.PutGauge(splitPath[1], metric)
			case "counter":
				if len(splitPath) == 1 {
					w.WriteHeader(http.StatusNotFound)
					return
				} else if len(splitPath) != 3 {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				metric, err := strconv.ParseInt(splitPath[2], 0, 64)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				repository.PutCounter(splitPath[1], metric)
			default:
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			w.Header().Set("content-type", "text/plain")
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})
	mux.HandleFunc("/", http.NotFound)
	return &Router{
		repository: repository,
		mux:        mux,
	}
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.mux.ServeHTTP(w, r)
}
