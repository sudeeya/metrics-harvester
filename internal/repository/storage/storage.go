package storage

import (
	"fmt"
	"slices"
	"strings"

	"github.com/sudeeya/metrics-harvester/internal/repository/metrics"
)

type MemStorage struct {
	metrics map[string]metrics.Metric
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		metrics: make(map[string]metrics.Metric),
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

func (ms *MemStorage) GetMetric(name string) (metrics.Metric, error) {
	if _, ok := ms.metrics[name]; !ok {
		return nil, fmt.Errorf("metric %s is missing", name)
	}
	return ms.metrics[name], nil
}

func (ms *MemStorage) GetAllMetrics() []metrics.Metric {
	allMetrics := make([]metrics.Metric, len(ms.metrics))
	i := 0
	for _, value := range ms.metrics {
		allMetrics[i] = value
		i++
	}
	slices.SortFunc(allMetrics, func(a, b metrics.Metric) int {
		return strings.Compare(a.GetName(), b.GetName())
	})
	return allMetrics
}
