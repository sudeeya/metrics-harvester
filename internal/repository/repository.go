// Package repository defines the interaction with an object storing metrics.
package repository

import (
	"context"

	"github.com/sudeeya/metrics-harvester/internal/metric"
)

// Repository describes interaction with an object storing metrics.
type Repository interface {
	// PutMetric inserts a metric into Repository.
	// Returns an error if the metric could not be inserted.
	PutMetric(ctx context.Context, m metric.Metric) error

	// PutBatch inserts a slice of metrics into Repository.
	// The same metric can be present in a slice with different values.
	// Returns an error if any of the metrics could not be inserted.
	PutBatch(ctx context.Context, metrics []metric.Metric) error

	// GetMetric returns a metric by its name (ID).
	// Returns an error if the metric could not be found.
	GetMetric(ctx context.Context, mName string) (metric.Metric, error)

	// GetAllMetrics returns a slice containing all metrics from Repository.
	// Returns an error if the metrics could not be retrieved.
	GetAllMetrics(ctx context.Context) ([]metric.Metric, error)

	// Close closes Repository.
	// Returns an error if the repository could not be closed.
	Close() error
}
