package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/sudeeya/metrics-harvester/internal/router"
)

func CreateGetAllMetricsHandler(router *router.Router) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allMetrics := router.GetAllMetrics()
		response := make([]string, len(allMetrics))
		for i, metric := range allMetrics {
			response[i] = metric.GetName() + ": " + metric.GetValue()
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("content-type", "text/plain")
		w.Write([]byte(strings.Join(response, "\n")))
	}
}

func CreateGetMetricHandler(router *router.Router) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			metricType = chi.URLParam(r, "metricType")
			metricName = chi.URLParam(r, "metricName")
		)
		switch metricType {
		case "gauge", "counter":
			metric, err := router.GetMetric(metricName)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Header().Set("content-type", "text/plain")
			w.Write([]byte(metric.GetName() + ": " + metric.GetValue()))
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

func CreatePostMetricHandler(router *router.Router) http.HandlerFunc {
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
			w.WriteHeader(http.StatusOK)
			w.Header().Set("content-type", "text/plain")
		case "counter":
			metric, err := strconv.ParseInt(metricValue, 0, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			router.PutCounter(metricName, metric)
			w.WriteHeader(http.StatusOK)
			w.Header().Set("content-type", "text/plain")
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

func BadRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
}
