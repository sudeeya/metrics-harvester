package handlers

import (
	"compress/gzip"
	"context"
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
	"github.com/sudeeya/metrics-harvester/internal/repository/database"
	"github.com/sudeeya/metrics-harvester/internal/repository/storage"
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

func responseOnError(logger *zap.Logger, err error, w http.ResponseWriter, statusCode int) {
	logger.Error(err.Error())
	w.WriteHeader(statusCode)
}

func NewAllMetricsHandler(logger *zap.Logger, repository repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allMetrics, err := repository.GetAllMetrics()
		if err != nil {
			responseOnError(logger, err, w, http.StatusInternalServerError)
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
			m, err := repository.GetMetric(metricName)
			if err != nil {
				responseOnError(logger, err, w, http.StatusNotFound)
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

func NewPingHandler(logger *zap.Logger, repository repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch v := repository.(type) {
		case *database.Database:
			if err := databaseResponse(v, w); err != nil {
				logger.Error(err.Error())
			}
		case *storage.MemStorage:
			http.Error(w, "the database is not in use", http.StatusInternalServerError)
		}
	}
}

func databaseResponse(db *database.Database, w http.ResponseWriter) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := db.DB.PingContext(ctx); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	w.WriteHeader(http.StatusOK)
	return nil
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
				responseOnError(logger, err, w, http.StatusBadRequest)
				return
			}
			if err := repository.PutMetric(metric.Metric{ID: metricName, MType: metricType, Value: &value}); err != nil {
				responseOnError(logger, err, w, http.StatusInternalServerError)
				return
			}
			w.Header().Set("content-type", "text/plain")
			w.WriteHeader(http.StatusOK)
		case metric.Counter:
			delta, err := strconv.ParseInt(metricValue, 0, 64)
			if err != nil {
				responseOnError(logger, err, w, http.StatusBadRequest)
				return
			}
			if err := repository.PutMetric(metric.Metric{ID: metricName, MType: metricType, Delta: &delta}); err != nil {
				responseOnError(logger, err, w, http.StatusInternalServerError)
				return
			}
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
			responseOnError(logger, err, w, http.StatusInternalServerError)
			return
		}
		if err := json.NewDecoder(body).Decode(&m); err != nil {
			responseOnError(logger, err, w, http.StatusBadRequest)
			return
		}
		if err := repository.PutMetric(m); err != nil {
			responseOnError(logger, err, w, http.StatusInternalServerError)
			return
		}
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		m, err = repository.GetMetric(m.ID)
		if err != nil {
			responseOnError(logger, err, w, http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(m); err != nil {
			logger.Error(err.Error())
		}
	}
}

func NewBatchHandler(logger *zap.Logger, repository repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var metrics []metric.Metric
		body, err := decompressIfNeeded(r)
		if err != nil {
			responseOnError(logger, err, w, http.StatusInternalServerError)
			return
		}
		if err := json.NewDecoder(body).Decode(&metrics); err != nil {
			responseOnError(logger, err, w, http.StatusBadRequest)
			return
		}
		if err := repository.PutBatch(metrics); err != nil {
			responseOnError(logger, err, w, http.StatusInternalServerError)
			return
		}
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusOK)
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
			w.Header().Set("content-type", "application/json")
			responseOnError(logger, err, w, http.StatusBadRequest)
			return
		}
		m, err := repository.GetMetric(requestedMetric.ID)
		if err != nil {
			w.Header().Set("content-type", "application/json")
			responseOnError(logger, err, w, http.StatusNotFound)
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
