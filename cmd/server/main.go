package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sudeeya/metrics-harvester/internal/handlers"
	"github.com/sudeeya/metrics-harvester/internal/repository/storage"
	"github.com/sudeeya/metrics-harvester/internal/router"
)

func main() {
	memStorage := storage.NewMemStorage()
	router := router.NewRouter(memStorage)
	metricHandler := handlers.CreateMetricHandler(router)
	router.Post("/", http.NotFound)
	router.Route("/update", func(r chi.Router) {
		r.Post("/", handlers.BadRequest)
		r.Route("/{metricType}", func(r chi.Router) {
			r.Post("/", http.NotFound)
			r.Route("/{metricName}", func(r chi.Router) {
				r.Post("/", handlers.BadRequest)
				r.Route("/{metricValue}", func(r chi.Router) {
					r.Post("/", metricHandler)
				})
			})
		})
	})
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}
}
