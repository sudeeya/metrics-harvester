package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sudeeya/metrics-harvester/internal/metric"
)

func int64Ptr(i int64) *int64 {
	return &i
}

func float64Ptr(f float64) *float64 {
	return &f
}

func TestPutMetric(t *testing.T) {
	var (
		ms1 = &MemStorage{metrics: map[string]metric.Metric{
			"gauge": {ID: "gauge", MType: metric.Gauge, Value: float64Ptr(12)},
		}}
		ms2 = &MemStorage{metrics: map[string]metric.Metric{
			"gauge": {ID: "gauge", MType: metric.Gauge, Value: float64Ptr(12.12)},
		}}
		ms3 = &MemStorage{metrics: map[string]metric.Metric{
			"counter": {ID: "counter", MType: metric.Counter, Delta: int64Ptr(12)},
		}}
		ms4 = &MemStorage{metrics: map[string]metric.Metric{
			"counter": {ID: "counter", MType: metric.Counter, Delta: int64Ptr(24)},
		}}
	)
	tests := []struct {
		name   string
		ms     *MemStorage
		mName  string
		m      metric.Metric
		result *MemStorage
	}{
		{
			name:   "put gauge",
			ms:     NewMemStorage(),
			mName:  "gauge",
			m:      metric.Metric{ID: "gauge", MType: metric.Gauge, Value: float64Ptr(12)},
			result: ms1,
		},
		{
			name:   "put the same gauge with a different value",
			ms:     ms1,
			mName:  "gauge",
			m:      metric.Metric{ID: "gauge", MType: metric.Gauge, Value: float64Ptr(12.12)},
			result: ms2,
		},
		{
			name:   "put counter",
			ms:     NewMemStorage(),
			mName:  "counter",
			m:      metric.Metric{ID: "counter", MType: metric.Counter, Delta: int64Ptr(12)},
			result: ms3,
		},
		{
			name:   "put the same counter",
			ms:     ms3,
			mName:  "counter",
			m:      metric.Metric{ID: "counter", MType: metric.Counter, Delta: int64Ptr(12)},
			result: ms4,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.ms.PutMetric(context.Background(), test.m)
			require.Equal(t, test.result.metrics[test.mName].GetValue(), test.ms.metrics[test.mName].GetValue())
		})
	}
}

func TestPutButch(t *testing.T) {
	var (
		ms1 = &MemStorage{metrics: map[string]metric.Metric{
			"gauge":   {ID: "gauge", MType: metric.Gauge, Value: float64Ptr(12.12)},
			"counter": {ID: "counter", MType: metric.Counter, Delta: int64Ptr(12)},
		}}
		ms2 = &MemStorage{metrics: map[string]metric.Metric{
			"gauge":   {ID: "gauge", MType: metric.Gauge, Value: float64Ptr(12.12)},
			"counter": {ID: "counter", MType: metric.Counter, Delta: int64Ptr(12)},
		}}
		ms3 = &MemStorage{metrics: map[string]metric.Metric{
			"gauge":   {ID: "gauge", MType: metric.Gauge, Value: float64Ptr(42)},
			"counter": {ID: "counter", MType: metric.Counter, Delta: int64Ptr(12)},
			"dummy":   {ID: "dummy", MType: metric.Counter, Delta: int64Ptr(-1)},
		}}
	)
	tests := []struct {
		name    string
		ms      *MemStorage
		metrics []metric.Metric
		result  *MemStorage
	}{
		{
			name: "put gauge and counter",
			ms:   NewMemStorage(),
			metrics: []metric.Metric{
				{ID: "gauge", MType: metric.Gauge, Value: float64Ptr(12.12)},
				{ID: "counter", MType: metric.Counter, Delta: int64Ptr(12)}},
			result: ms1,
		},
		{
			name:    "put nothing",
			ms:      ms1,
			metrics: make([]metric.Metric, 0),
			result:  ms2,
		},
		{
			name: "put the same gauge and new counter",
			ms:   ms2,
			metrics: []metric.Metric{
				{ID: "gauge", MType: metric.Gauge, Value: float64Ptr(42)},
				{ID: "dummy", MType: metric.Counter, Delta: int64Ptr(-1)}},
			result: ms3,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.ms.PutBatch(context.Background(), test.metrics)
			require.Equal(t, len(test.result.metrics), len(test.ms.metrics))
			for mName := range test.result.metrics {
				require.Equal(t, test.result.metrics[mName].GetValue(), test.ms.metrics[mName].GetValue())
			}
		})
	}
}

func TestGetMetric_Existing(t *testing.T) {
	var (
		ms1 = &MemStorage{metrics: map[string]metric.Metric{
			"gauge": {ID: "gauge", MType: metric.Gauge, Value: float64Ptr(12)},
		}}
		ms2 = &MemStorage{metrics: map[string]metric.Metric{
			"counter": {ID: "counter", MType: metric.Counter, Delta: int64Ptr(12)},
			"dummy":   {ID: "dummy", MType: metric.Gauge, Value: float64Ptr(-1)},
		}}
	)
	tests := []struct {
		name   string
		ms     *MemStorage
		mName  string
		result metric.Metric
	}{
		{
			name:   "get gauge",
			ms:     ms1,
			mName:  "gauge",
			result: metric.Metric{ID: "gauge", MType: metric.Gauge, Value: float64Ptr(12)},
		},
		{
			name:   "get counter",
			ms:     ms2,
			mName:  "counter",
			result: metric.Metric{ID: "counter", MType: metric.Counter, Delta: int64Ptr(12)},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m, err := test.ms.GetMetric(context.Background(), test.mName)
			require.Nil(t, err)
			require.Equal(t, test.result.GetValue(), m.GetValue())
		})
	}
}

func TestGetMetric_NotExisting(t *testing.T) {
	var (
		ms1 = &MemStorage{metrics: map[string]metric.Metric{
			"counter": {ID: "counter", MType: metric.Counter, Delta: int64Ptr(12)},
		}}
	)
	tests := []struct {
		name  string
		ms    *MemStorage
		mName string
	}{
		{
			name:  "try to get gauge",
			ms:    ms1,
			mName: "gauge",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := test.ms.GetMetric(context.Background(), test.mName)
			require.NotNil(t, err)
		})
	}
}

func TestGetAllMetrics(t *testing.T) {
	var (
		ms1 = &MemStorage{metrics: map[string]metric.Metric{
			"gauge":   {ID: "gauge", MType: metric.Gauge, Value: float64Ptr(42)},
			"counter": {ID: "counter", MType: metric.Counter, Delta: int64Ptr(12)},
			"dummy":   {ID: "dummy", MType: metric.Counter, Delta: int64Ptr(-1)},
		}}
	)
	tests := []struct {
		name   string
		ms     *MemStorage
		result []metric.Metric
	}{
		{
			name:   "empty storage",
			ms:     NewMemStorage(),
			result: make([]metric.Metric, 0),
		},
		{
			name: "not empty storage",
			ms:   ms1,
			result: []metric.Metric{
				{ID: "counter", MType: metric.Counter, Delta: int64Ptr(12)},
				{ID: "dummy", MType: metric.Counter, Delta: int64Ptr(-1)},
				{ID: "gauge", MType: metric.Gauge, Value: float64Ptr(42)},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			metrics, err := test.ms.GetAllMetrics(context.Background())
			require.Nil(t, err)
			require.Equal(t, len(test.result), len(metrics))
			for i := range test.result {
				require.Equal(t, test.result[i].GetValue(), metrics[i].GetValue())
			}
		})
	}
}

func TestClose(t *testing.T) {
	var (
		ms1 = &MemStorage{metrics: map[string]metric.Metric{
			"gauge":   {ID: "gauge", MType: metric.Gauge, Value: float64Ptr(42)},
			"counter": {ID: "counter", MType: metric.Counter, Delta: int64Ptr(12)},
			"dummy":   {ID: "dummy", MType: metric.Counter, Delta: int64Ptr(-1)},
		}}
	)
	tests := []struct {
		name string
		ms   *MemStorage
	}{
		{
			name: "empty storage",
			ms:   NewMemStorage(),
		},
		{
			name: "not empty storage",
			ms:   ms1,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Nil(t, test.ms.Close())
		})
	}
}
