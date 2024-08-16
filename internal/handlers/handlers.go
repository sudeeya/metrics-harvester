package handlers

import (
	"compress/gzip"
	"context"
	"database/sql"
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/sudeeya/metrics-harvester/internal/metric"
	repo "github.com/sudeeya/metrics-harvester/internal/repository"
	"go.uber.org/zap"
)

const htmlTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
	<style>
        body {
            background-color: #000000;
            color: #ffffff;
        }
	</style>
</head>
<body>
    <ul>
        {{range .Metrics}}
        <li>{{.ID}}: {{.Value}}</li>
        {{end}}
    </ul>
</body>
`

type htmlMetric struct {
	ID    string
	Value string
}

func NewAllMetricsHandler(logger *zap.Logger, repository repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allMetrics, err := repository.GetAllMetrics()
		if err != nil {
			logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		metrics := make([]htmlMetric, len(allMetrics))
		for i, m := range allMetrics {
			metrics[i].ID = m.ID
			metrics[i].Value = m.GetValue()
		}
		data := struct {
			Metrics []htmlMetric
		}{
			Metrics: metrics,
		}
		t, _ := template.New("page").Parse(htmlTemplate)
		w.Header().Set("content-type", "text/html")
		w.WriteHeader(http.StatusOK)
		if err = t.Execute(w, data); err != nil {
			logger.Error(err.Error())
		}
	}
}

func NewValueHandler(logger *zap.Logger, repository repo.Repository) http.HandlerFunc {
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
			w.Header().Set("content-type", "text/plain")
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte(m.GetValue())); err != nil {
				logger.Error(err.Error())
			}
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

func NewPingHandler(logger *zap.Logger, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		if err := db.PingContext(ctx); err != nil {
			logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
	}
}

func NewUpdateHandler(logger *zap.Logger, repository repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			metricType  = chi.URLParam(r, "metricType")
			metricName  = chi.URLParam(r, "metricName")
			metricValue = chi.URLParam(r, "metricValue")
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
			w.Header().Set("content-type", "text/plain")
			w.WriteHeader(http.StatusOK)
		case metric.Counter:
			delta, err := strconv.ParseInt(metricValue, 0, 64)
			if err != nil {
				logger.Error(err.Error())
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			repository.PutMetric(metric.Metric{ID: metricName, MType: metricType, Delta: &delta})
			w.Header().Set("content-type", "text/plain")
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

func NewJSONUpdateHandler(logger *zap.Logger, repository repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var m metric.Metric
		body, err := decompressIfNeeded(r)
		if err != nil {
			logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := json.NewDecoder(body).Decode(&m); err != nil {
			logger.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		repository.PutMetric(m)
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		m, _ = repository.GetMetric(m.ID)
		if err := json.NewEncoder(w).Encode(m); err != nil {
			logger.Error(err.Error())
		}
	}
}

func decompressIfNeeded(r *http.Request) (io.Reader, error) {
	if strings.Contains(r.Header.Get("content-encoding"), "gzip") {
		return gzip.NewReader(r.Body)
	}
	return r.Body, nil
}

func NewJSONValueHandler(logger *zap.Logger, repository repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestedMetric metric.Metric
		if err := json.NewDecoder(r.Body).Decode(&requestedMetric); err != nil {
			logger.Error(err.Error())
			w.Header().Set("content-type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		m, ok := repository.GetMetric(requestedMetric.ID)
		if !ok {
			w.Header().Set("content-type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(m); err != nil {
			logger.Error(err.Error())
		}
	}
}

func BadRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
}
