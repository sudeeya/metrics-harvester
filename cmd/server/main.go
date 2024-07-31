package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sudeeya/metrics-harvester/internal/handlers"
	"github.com/sudeeya/metrics-harvester/internal/repository/storage"
	"github.com/sudeeya/metrics-harvester/internal/router"
)

func main() {
	var (
		memStorage           = storage.NewMemStorage()
		router               = router.NewRouter(memStorage)
		getAllMetricsHandler = handlers.CreateGetAllMetricsHandler(router)
		getMetricHandler     = handlers.CreateGetMetricHandler(router)
		postMetricHandler    = handlers.CreatePostMetricHandler(router)
	)
	router.Get("/", getAllMetricsHandler)
	router.Route("/value", func(r chi.Router) {
		r.Get("/", http.NotFound)
		r.Route("/{metricType}", func(r chi.Router) {
			r.Get("/", http.NotFound)
			r.Route("/{metricName}", func(r chi.Router) {
				r.Get("/", getMetricHandler)
			})
		})
	})
	router.Post("/", http.NotFound)
	router.Route("/update", func(r chi.Router) {
		r.Post("/", handlers.BadRequest)
		r.Route("/{metricType}", func(r chi.Router) {
			r.Post("/", http.NotFound)
			r.Route("/{metricName}", func(r chi.Router) {
				r.Post("/", handlers.BadRequest)
				r.Route("/{metricValue}", func(r chi.Router) {
					r.Post("/", postMetricHandler)
				})
			})
		})
	})
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}
}
