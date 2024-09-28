// Package handlers provides a collection of HTTP handlers.
package handlers

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/sudeeya/metrics-harvester/internal/metric"
	repo "github.com/sudeeya/metrics-harvester/internal/repository"
	"github.com/sudeeya/metrics-harvester/internal/repository/database"
	"github.com/sudeeya/metrics-harvester/internal/repository/storage"
)

const limitInSeconds = 10

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

// NewAllMetricsHandler returns an http.HandlerFunc that generates
// an HTML list of all metrics from the Repository and writes it to the response.
// If an error occurs while retrieving the metrics, it logs the error and returns an appropriate HTTP status code.
func NewAllMetricsHandler(logger *zap.Logger, repository repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), limitInSeconds*time.Second)
		defer cancel()

		allMetrics, err := repository.GetAllMetrics(ctx)
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			logger.Error(ctx.Err().Error())
		}
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

// NewValueHandler returns an http.HandlerFunc that writes the value of a specified metric to the response.
// The metric type and name are extracted from the URL parameters.
// If the metric type is not supported or an error occurs while retrieving the metric,
// it logs the error and returns an appropriate HTTP status code.
func NewValueHandler(logger *zap.Logger, repository repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			metricType = chi.URLParam(r, "metricType")
			metricName = chi.URLParam(r, "metricName")
		)
		switch metricType {
		case metric.Gauge, metric.Counter:
			m, err := repository.GetMetric(context.Background(), metricName)
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

// NewPingHandler returns an http.HandlerFunc that pings the database to check its availability.
// If the repository is a database, it attempts to ping the database.
// If the repository is an in-memory storage, it returns an error indicating that the database is not in use.
func NewPingHandler(logger *zap.Logger, repository repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch v := repository.(type) {
		case *database.Database:
			if err := databaseResponse(logger, v, w); err != nil {
				logger.Error(err.Error())
			}
		case *storage.MemStorage:
			http.Error(w, "the database is not in use", http.StatusInternalServerError)
		}
	}
}

func databaseResponse(logger *zap.Logger, db *database.Database, w http.ResponseWriter) error {
	ctx, cancel := context.WithTimeout(context.Background(), limitInSeconds*time.Second)
	defer cancel()

	if err := db.DB.PingContext(ctx); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		logger.Error(ctx.Err().Error())
	}

	w.WriteHeader(http.StatusOK)
	return nil
}

// NewUpdateHandler returns an http.HandlerFunc that updates a specified metric.
// The metric type, name and value are extracted from the URL parameters.
// If the metric type is not supported or an error occurs while updating the metric,
// it logs the error and returns an appropriate HTTP status code.
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

			if err := repository.PutMetric(context.Background(), metric.Metric{ID: metricName, MType: metricType, Value: &value}); err != nil {
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

			if err := repository.PutMetric(context.Background(), metric.Metric{ID: metricName, MType: metricType, Delta: &delta}); err != nil {
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

// NewJSONUpdateHandler returns an http.HandlerFunc that updates a specified metric.
// The metric is extracted from the JSON body of the request.
// If the metric type is not supported or an error occurs while updating the metric,
// it logs the error and returns an appropriate HTTP status code.
// After updating the metric, it retrieves the updated metric and returns it in the response.
func NewJSONUpdateHandler(logger *zap.Logger, repository repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var m metric.Metric
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			responseOnError(logger, err, w, http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), limitInSeconds*time.Second)
		defer cancel()

		if err := repository.PutMetric(ctx, m); err != nil {
			responseOnError(logger, err, w, http.StatusInternalServerError)
			return
		}

		m, err := repository.GetMetric(ctx, m.ID)
		if err != nil {
			responseOnError(logger, err, w, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(m); err != nil {
			logger.Error(err.Error())
		}
	}
}

// NewBatchHandler returns an http.HandlerFunc that updates a batch os metrics.
// Metrics are extracted from the JSON body of the request.
// If the metric type is not supported or an error occurs while updating the metric,
// it logs the error and returns an appropriate HTTP status code.
func NewBatchHandler(logger *zap.Logger, repository repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var metrics []metric.Metric
		if err := json.NewDecoder(r.Body).Decode(&metrics); err != nil {
			responseOnError(logger, err, w, http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), limitInSeconds*time.Second)
		defer cancel()

		if err := repository.PutBatch(ctx, metrics); err != nil {
			responseOnError(logger, err, w, http.StatusInternalServerError)
			return
		}
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			logger.Error(ctx.Err().Error())
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
	}
}

// NewJSONValueHandler returns an http.HandlerFunc that writes the JSON of a specified metric to the response.
// The metric type and name are extracted from the JSON body of the request.
// If the metric type is not supported or an error occurs while updating the metric,
// it logs the error and returns an appropriate HTTP status code.
func NewJSONValueHandler(logger *zap.Logger, repository repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestedMetric metric.Metric
		if err := json.NewDecoder(r.Body).Decode(&requestedMetric); err != nil {
			w.Header().Set("Content-Type", "application/json")
			responseOnError(logger, err, w, http.StatusBadRequest)
			return
		}

		m, err := repository.GetMetric(context.Background(), requestedMetric.ID)
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

// NewKeyHandler returns an http.HandlerFunc that writes symmetric key from request to the channel.
func NewKeyHandler(logger *zap.Logger, privateKey *rsa.PrivateKey, symmetricKey chan<- []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		decryptedBody, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, body, nil)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		symmetricKey <- decryptedBody

		w.WriteHeader(http.StatusOK)
	}
}
