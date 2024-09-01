package handlers

import (
	"context"
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"
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

func NewAllMetricsHandler(ctx context.Context, logger *zap.Logger, repository repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allMetrics, err := repository.GetAllMetrics(ctx)
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
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		if err = t.Execute(w, data); err != nil {
			logger.Error(err.Error())
		}
	}
}

func NewValueHandler(ctx context.Context, logger *zap.Logger, repository repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			metricType = chi.URLParam(r, "metricType")
			metricName = chi.URLParam(r, "metricName")
		)
		switch metricType {
		case metric.Gauge, metric.Counter:
			m, err := repository.GetMetric(ctx, metricName)
			if err != nil {
				responseOnError(logger, err, w, http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte(m.GetValue())); err != nil {
				logger.Error(err.Error())
			}
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

func NewPingHandler(ctx context.Context, logger *zap.Logger, repository repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch v := repository.(type) {
		case *database.Database:
			if err := databaseResponse(ctx, v, w); err != nil {
				logger.Error(err.Error())
			}
		case *storage.MemStorage:
			http.Error(w, "the database is not in use", http.StatusInternalServerError)
		}
	}
}

func databaseResponse(ctx context.Context, db *database.Database, w http.ResponseWriter) error {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	if err := db.DB.PingContext(ctx); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	w.WriteHeader(http.StatusOK)
	return nil
}

func NewUpdateHandler(ctx context.Context, logger *zap.Logger, repository repo.Repository) http.HandlerFunc {
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
			if err := repository.PutMetric(ctx, metric.Metric{ID: metricName, MType: metricType, Value: &value}); err != nil {
				responseOnError(logger, err, w, http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
		case metric.Counter:
			delta, err := strconv.ParseInt(metricValue, 0, 64)
			if err != nil {
				responseOnError(logger, err, w, http.StatusBadRequest)
				return
			}
			if err := repository.PutMetric(ctx, metric.Metric{ID: metricName, MType: metricType, Delta: &delta}); err != nil {
				responseOnError(logger, err, w, http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

func NewJSONUpdateHandler(ctx context.Context, logger *zap.Logger, repository repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var m metric.Metric
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			responseOnError(logger, err, w, http.StatusBadRequest)
			return
		}
		if err := repository.PutMetric(ctx, m); err != nil {
			responseOnError(logger, err, w, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		m, err := repository.GetMetric(ctx, m.ID)
		if err != nil {
			responseOnError(logger, err, w, http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(m); err != nil {
			logger.Error(err.Error())
		}
	}
}

func NewBatchHandler(ctx context.Context, logger *zap.Logger, repository repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var metrics []metric.Metric
		if err := json.NewDecoder(r.Body).Decode(&metrics); err != nil {
			responseOnError(logger, err, w, http.StatusBadRequest)
			return
		}
		if err := repository.PutBatch(ctx, metrics); err != nil {
			responseOnError(logger, err, w, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
	}
}

func NewJSONValueHandler(ctx context.Context, logger *zap.Logger, repository repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestedMetric metric.Metric
		if err := json.NewDecoder(r.Body).Decode(&requestedMetric); err != nil {
			w.Header().Set("Content-Type", "application/json")
			responseOnError(logger, err, w, http.StatusBadRequest)
			return
		}
		m, err := repository.GetMetric(ctx, requestedMetric.ID)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			responseOnError(logger, err, w, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(m); err != nil {
			logger.Error(err.Error())
		}
	}
}

func BadRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
}
