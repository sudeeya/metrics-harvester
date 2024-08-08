package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/sudeeya/metrics-harvester/internal/metric"
	repo "github.com/sudeeya/metrics-harvester/internal/repository"
	"go.uber.org/zap"
)

func NewGetAllMetricsHandler(logger *zap.Logger, repository repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allMetrics := repository.GetAllMetrics()
		response := make([]string, len(allMetrics))
		for i, m := range allMetrics {
			response[i] = fmt.Sprintf("%s: %s", m.ID, m.GetValue())
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("content-type", "text/plain")
		if _, err := w.Write([]byte(strings.Join(response, "\n"))); err != nil {
			logger.Error(err.Error())
		}
	}
}

func NewGetMetricHandler(logger *zap.Logger, repository repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			metricType = chi.URLParam(r, "metricType")
			metricName = chi.URLParam(r, "metricName")
		)
		switch metricType {
		case metric.Gauge, metric.Counter:
			m, ok := repository.GetMetric(metricName)
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Header().Set("content-type", "text/plain")
			if _, err := w.Write([]byte(m.GetValue())); err != nil {
				logger.Error(err.Error())
			}
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

func NewPostMetricHandler(logger *zap.Logger, repository repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			metricType  = chi.URLParam(r, "metricType")
			metricName  = chi.URLParam(r, "metricName")
			metricValue = chi.URLParam(r, "metricValue")
			jsonEncoder = json.NewEncoder(w)
		)
		switch metricType {
		case metric.Gauge:
			value, err := strconv.ParseFloat(metricValue, 64)
			if err != nil {
				logger.Error(err.Error())
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			repository.PutMetric(metric.Metric{ID: metricName, MType: metricType, Value: &value})
			w.WriteHeader(http.StatusOK)
			w.Header().Set("content-type", "application/json")
			m, _ := repository.GetMetric(metricName)
			if err := jsonEncoder.Encode(m); err != nil {
				logger.Error(err.Error())
			}
		case metric.Counter:
			delta, err := strconv.ParseInt(metricValue, 0, 64)
			if err != nil {
				logger.Error(err.Error())
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			repository.PutMetric(metric.Metric{ID: metricName, MType: metricType, Delta: &delta})
			w.WriteHeader(http.StatusOK)
			w.Header().Set("content-type", "application/json")
			m, _ := repository.GetMetric(metricName)
			if err := jsonEncoder.Encode(m); err != nil {
				logger.Error(err.Error())
			}
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

func NewPostJSONMetricHandler(logger *zap.Logger, repository repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var m metric.Metric
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			logger.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("content-type", "application/json")
			return
		}
		m, ok := repository.GetMetric(m.ID)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			w.Header().Set("content-type", "application/json")
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("content-type", "application/json")
		if err := json.NewEncoder(w).Encode(m); err != nil {
			logger.Error(err.Error())
		}
	}
}

func BadRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
}
