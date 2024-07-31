package repository

import "github.com/sudeeya/metrics-harvester/internal/repository/metrics"

type Repository interface {
	PutGauge(name string, value float64)
	PutCounter(name string, value int64)
	GetMetric(name string) (metrics.Metric, error)
	GetAllMetrics() []metrics.Metric
}
