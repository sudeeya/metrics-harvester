package handlers

import (
	"net/http"
	"strconv"
	"strings"

	repo "github.com/sudeeya/metrics-harvester/internal/repository"
)

func CreateGaugeHandler(repository repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			path := r.URL.Path[len("/update/gauge/"):]
			splitPath := strings.Split(path, "/")
			if splitPath[0] == "" {
				w.WriteHeader(http.StatusNotFound)
				return
			} else if len(splitPath) != 2 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			metric, err := strconv.ParseFloat(splitPath[1], 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			repository.PutGauge(splitPath[0], metric)
			w.Header().Set("content-type", "text/plain")
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

func CreateCounterHandler(repository repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			path := r.URL.Path[len("/update/counter/"):]
			splitPath := strings.Split(path, "/")
			if splitPath[0] == "" {
				w.WriteHeader(http.StatusNotFound)
				return
			} else if len(splitPath) != 2 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			metric, err := strconv.ParseInt(splitPath[2], 0, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			repository.PutCounter(splitPath[0], metric)
			w.Header().Set("content-type", "text/plain")
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

func BadRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
}
