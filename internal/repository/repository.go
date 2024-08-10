package repository

import "github.com/sudeeya/metrics-harvester/internal/metric"

type Repository interface {
	PutMetric(m metric.Metric) error
	GetMetric(mName string) (metric.Metric, bool)
	GetAllMetrics() ([]metric.Metric, error)
}
