package repository

import (
	"context"

	"github.com/sudeeya/metrics-harvester/internal/metric"
)

type Repository interface {
	PutMetric(ctx context.Context, m metric.Metric) error
	PutBatch(ctx context.Context, metrics []metric.Metric) error
	GetMetric(ctx context.Context, mName string) (metric.Metric, error)
	GetAllMetrics(ctx context.Context) ([]metric.Metric, error)
	Close() error
}
