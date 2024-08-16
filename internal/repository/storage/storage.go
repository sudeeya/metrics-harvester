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

func (ms *MemStorage) PutMetric(m metric.Metric) error {
	value, ok := ms.metrics[m.ID]
	if !ok {
		ms.metrics[m.ID] = m
		return nil
	}
	switch m.MType {
	case metric.Gauge:
		value.Update(*m.Value)
	case metric.Counter:
		value.Update(*m.Delta)
	}
	ms.metrics[m.ID] = value
	return nil
}

func (ms *MemStorage) GetMetric(mName string) (metric.Metric, error) {
	m, ok := ms.metrics[mName]
	if !ok {
		return metric.Metric{}, fmt.Errorf("metric %s is missing", mName)
	}
	return m, nil
}

func (ms *MemStorage) GetAllMetrics() ([]metric.Metric, error) {
	allMetrics := make([]metric.Metric, len(ms.metrics))
	i := 0
	for _, value := range ms.metrics {
		allMetrics[i] = value
		i++
	}
	slices.SortFunc(allMetrics, func(a, b metric.Metric) int {
		return strings.Compare(a.ID, b.ID)
	})
	return allMetrics, nil
}

func (ms *MemStorage) Close() error {
	return nil
}
