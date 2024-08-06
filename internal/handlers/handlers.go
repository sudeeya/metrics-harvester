package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	repo "github.com/sudeeya/metrics-harvester/internal/repository"
	"go.uber.org/zap"
)

func CreateGetAllMetricsHandler(logger *zap.Logger, repository repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allMetrics := repository.GetAllMetrics()
		response := make([]string, len(allMetrics))
		for i, metric := range allMetrics {
			response[i] = metric.GetName() + ": " + metric.GetValue()
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("content-type", "text/plain")
		if _, err := w.Write([]byte(strings.Join(response, "\n"))); err != nil {
			logger.Error(err.Error())
		}
	}
}

func CreateGetMetricHandler(logger *zap.Logger, repository repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			metricType = chi.URLParam(r, "metricType")
			metricName = chi.URLParam(r, "metricName")
		)
		switch metricType {
		case "gauge", "counter":
			metric, err := repository.GetMetric(metricName)
			if err != nil {
				logger.Error(err.Error())
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Header().Set("content-type", "text/plain")
			if _, err := w.Write([]byte(metric.GetValue())); err != nil {
				logger.Error(err.Error())
			}
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

func CreatePostMetricHandler(logger *zap.Logger, repository repo.Repository) http.HandlerFunc {
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
				logger.Error(err.Error())
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			repository.PutGauge(metricName, metric)
			w.WriteHeader(http.StatusOK)
			w.Header().Set("content-type", "text/plain")
		case "counter":
			metric, err := strconv.ParseInt(metricValue, 0, 64)
			if err != nil {
				logger.Error(err.Error())
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			repository.PutCounter(metricName, metric)
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
