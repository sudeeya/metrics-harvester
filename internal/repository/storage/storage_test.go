package storage

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/sudeeya/metrics-harvester/internal/repository/metrics"
)

func TestPutGauge(t *testing.T) {
	var (
		ms1 = &MemStorage{metrics: map[string]metrics.Metric{
			"gauge": metrics.NewGauge("gauge", 12),
		}}
		ms2 = &MemStorage{metrics: map[string]metrics.Metric{
			"gauge": metrics.NewGauge("gauge", 12.12),
		}}
		ms3 = &MemStorage{metrics: map[string]metrics.Metric{
			"gauge": metrics.NewGauge("gauge", 12.12),
			"dummy": metrics.NewGauge("dummy", -1),
		}}
	)
	tests := []struct {
		ms     *MemStorage
		name   string
		value  float64
		result *MemStorage
	}{
		{
			ms:     NewMemStorage(),
			name:   "gauge",
			value:  12,
			result: ms1,
		},
		{
			ms:     ms1,
			name:   "gauge",
			value:  12.12,
			result: ms2,
		},
		{
			ms:     ms2,
			name:   "dummy",
			value:  -1,
			result: ms3,
		},
	}
	for _, test := range tests {
		test.ms.PutGauge(test.name, test.value)
		require.EqualValues(t, test.ms, test.result)
	}
}

func TestPutCounter(t *testing.T) {
	var (
		ms1 = &MemStorage{metrics: map[string]metrics.Metric{
			"counter": metrics.NewCounter("counter", 12),
		}}
		ms2 = &MemStorage{metrics: map[string]metrics.Metric{
			"counter": metrics.NewCounter("counter", 24),
		}}
		ms3 = &MemStorage{metrics: map[string]metrics.Metric{
			"counter": metrics.NewCounter("counter", 24),
			"dummy":   metrics.NewCounter("dummy", -1),
		}}
	)
	tests := []struct {
		ms     *MemStorage
		name   string
		value  int64
		result *MemStorage
	}{
		{
			ms:     NewMemStorage(),
			name:   "counter",
			value:  12,
			result: ms1,
		},
		{
			ms:     ms1,
			name:   "counter",
			value:  12,
			result: ms2,
		},
		{
			ms:     ms2,
			name:   "dummy",
			value:  -1,
			result: ms3,
		},
	}
	for _, test := range tests {
		test.ms.PutCounter(test.name, test.value)
		require.EqualValues(t, test.ms, test.result)
	}
}

func TestGetMetric_NoError(t *testing.T) {
	var (
		ms1 = &MemStorage{metrics: map[string]metrics.Metric{
			"gauge": metrics.NewGauge("gauge", 12.12),
		}}
		ms2 = &MemStorage{metrics: map[string]metrics.Metric{
			"counter": metrics.NewCounter("counter", 12),
			"dummy":   metrics.NewCounter("dummy", -1),
		}}
	)
	tests := []struct {
		ms     *MemStorage
		name   string
		result metrics.Metric
	}{
		{
			ms:     ms1,
			name:   "gauge",
			result: metrics.NewGauge("gauge", 12.12),
		},
		{
			ms:     ms2,
			name:   "counter",
			result: metrics.NewCounter("counter", 12),
		},
	}
	for _, test := range tests {
		metric, err := test.ms.GetMetric(test.name)
		require.Nil(t, err)
		require.EqualValues(t, metric, test.result)
	}
}

func TestGetMetric_Error(t *testing.T) {
	var (
		ms1 = &MemStorage{metrics: map[string]metrics.Metric{
			"counter": metrics.NewCounter("counter", 12),
		}}
	)
	tests := []struct {
		ms   *MemStorage
		name string
	}{
		{
			ms:   ms1,
			name: "gauge",
		},
	}
	for _, test := range tests {
		_, err := test.ms.GetMetric(test.name)
		require.NotNil(t, err)
	}
}
