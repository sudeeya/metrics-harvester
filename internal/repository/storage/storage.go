package storage

import repo "github.com/sudeeya/metrics-harvester/internal/repository"

type MemStorage struct {
	gauges   map[string]repo.Gauge
	counters map[string]repo.Counter
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauges:   make(map[string]repo.Gauge),
		counters: make(map[string]repo.Counter),
	}
}

func (ms *MemStorage) PutGauge(metricName string, metricValue repo.Gauge) {
	ms.gauges[metricName] = metricValue
}

func (ms *MemStorage) PutCounter(metricName string, metricValue repo.Counter) {
	if _, ok := ms.counters[metricName]; !ok {
		ms.counters[metricName] = metricValue
	} else {
		ms.counters[metricName] += metricValue
	}
}
