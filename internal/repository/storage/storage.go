package storage

import (
	"fmt"
	"slices"
	"strings"

	"github.com/sudeeya/metrics-harvester/internal/metric"
)

type MemStorage struct {
	metrics map[string]metric.Metric
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		metrics: make(map[string]metric.Metric),
	}
}

func (ms *MemStorage) PutGauge(name string, value float64) {
	if _, ok := ms.metrics[name]; !ok {
		ms.metrics[name] = metric.NewGauge(name, value)
	} else {
		ms.metrics[name].(*metric.Gauge).ChangeValue(value)
	}
}

func (ms *MemStorage) PutCounter(name string, value int64) {
	if _, ok := ms.metrics[name]; !ok {
		ms.metrics[name] = metric.NewCounter(name, value)
	} else {
		ms.metrics[name].(*metric.Counter).IncreaseValue(value)
	}
}

func (ms *MemStorage) GetMetric(name string) (metric.Metric, error) {
	if _, ok := ms.metrics[name]; !ok {
		return nil, fmt.Errorf("metric %s is missing", name)
	}
	return ms.metrics[name], nil
}

func (ms *MemStorage) GetAllMetrics() []metric.Metric {
	allMetrics := make([]metric.Metric, len(ms.metrics))
	i := 0
	for _, value := range ms.metrics {
		allMetrics[i] = value
		i++
	}
	slices.SortFunc(allMetrics, func(a, b metric.Metric) int {
		return strings.Compare(a.GetName(), b.GetName())
	})
	return allMetrics
}
