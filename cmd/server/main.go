package main

import (
	"net/http"
	"strconv"
	"strings"

	repo "github.com/sudeeya/metrics-harvester/internal/repository"
	"github.com/sudeeya/metrics-harvester/internal/repository/storage"
)

func main() {
	memStorage := storage.NewMemStorage()
	mux := http.NewServeMux()
	mux.HandleFunc("/update/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			path := r.URL.Path[len("/update/"):]
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
				memStorage.PutGauge(splitPath[1], repo.Gauge(metric))
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
				memStorage.PutCounter(splitPath[1], repo.Counter(metric))
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
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
