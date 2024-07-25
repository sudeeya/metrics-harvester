package repository

type Gauge float64

type Counter int64

type Repository interface {
	PutGauge(metricName string, metricValue Gauge)
	PutCounter(metricName string, metricValue Counter)
}
