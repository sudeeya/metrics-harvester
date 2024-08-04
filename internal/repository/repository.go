package repository

import "github.com/sudeeya/metrics-harvester/internal/metric"

type Repository interface {
	PutGauge(name string, value float64)
	PutCounter(name string, value int64)
	GetMetric(name string) (metric.Metric, error)
	GetAllMetrics() []metric.Metric
}
