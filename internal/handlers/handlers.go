package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sudeeya/metrics-harvester/internal/router"
)

func CreateMetricHandler(router *router.Router) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			metricType  = chi.URLParam(r, "metricType")
			metricName  = chi.URLParam(r, "metricName")
			metricValue = chi.URLParam(r, "metricValue")
		)
		switch metricType {
		case "gauge":
			metric, err := strconv.ParseFloat(metricValue, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			router.PutGauge(metricName, metric)
			w.Header().Set("content-type", "text/plain")
			w.WriteHeader(http.StatusOK)
		case "counter":
			metric, err := strconv.ParseInt(metricValue, 0, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			router.PutCounter(metricName, metric)
			w.Header().Set("content-type", "text/plain")
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

func BadRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
}
