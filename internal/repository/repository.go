package repository

import "github.com/sudeeya/metrics-harvester/internal/metric"

type Repository interface {
	PutMetric(m metric.Metric) error
	PutBatch(metrics []metric.Metric) error
	GetMetric(mName string) (metric.Metric, error)
	GetAllMetrics() ([]metric.Metric, error)
	Close() error
}
