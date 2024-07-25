package repository

type Repository interface {
	PutGauge(name string, value float64)
	PutCounter(name string, value int64)
}
