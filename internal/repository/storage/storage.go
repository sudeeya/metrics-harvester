package storage

import "github.com/sudeeya/metrics-harvester/internal/repository/metrics"

type MemStorage struct {
	metrics map[string]any
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		metrics: make(map[string]any),
	}
}

func (ms *MemStorage) PutGauge(name string, value float64) {
	if _, ok := ms.metrics[name]; !ok {
		ms.metrics[name] = metrics.NewGauge(name, value)
	} else {
		ms.metrics[name].(*metrics.Gauge).ChangeValue(value)
	}
}

func (ms *MemStorage) PutCounter(name string, value int64) {
	if _, ok := ms.metrics[name]; !ok {
		ms.metrics[name] = metrics.NewCounter(name, value)
	} else {
		ms.metrics[name].(*metrics.Counter).IncreaseValue(value)
	}
}
