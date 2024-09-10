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
		ms     *MemStorage
		mName  string
		m      metric.Metric
		result *MemStorage
	}{
		{
			ms:     NewMemStorage(),
			mName:  "gauge",
			m:      metric.Metric{ID: "gauge", MType: metric.Gauge, Value: float64Ptr(12)},
			result: ms1,
		},
		{
			ms:     ms1,
			mName:  "gauge",
			m:      metric.Metric{ID: "gauge", MType: metric.Gauge, Value: float64Ptr(12.12)},
			result: ms2,
		},
		{
			ms:     NewMemStorage(),
			mName:  "counter",
			m:      metric.Metric{ID: "counter", MType: metric.Counter, Delta: int64Ptr(12)},
			result: ms3,
		},
		{
			ms:     ms3,
			mName:  "counter",
			m:      metric.Metric{ID: "counter", MType: metric.Counter, Delta: int64Ptr(12)},
			result: ms4,
		},
	}
	for _, test := range tests {
		test.ms.PutMetric(context.Background(), test.m)
		require.EqualValues(t, test.result.metrics[test.mName].GetValue(), test.ms.metrics[test.mName].GetValue())
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
		ms     *MemStorage
		mName  string
		result metric.Metric
	}{
		{
			ms:     ms1,
			mName:  "gauge",
			result: metric.Metric{ID: "gauge", MType: metric.Gauge, Value: float64Ptr(12)},
		},
		{
			ms:     ms2,
			mName:  "counter",
			result: metric.Metric{ID: "counter", MType: metric.Counter, Delta: int64Ptr(12)},
		},
	}
	for _, test := range tests {
		m, err := test.ms.GetMetric(context.Background(), test.mName)
		require.Nil(t, err)
		require.EqualValues(t, test.result.GetValue(), m.GetValue())
	}
}

func TestGetMetric_NotExisting(t *testing.T) {
	var (
		ms1 = &MemStorage{metrics: map[string]metric.Metric{
			"counter": {ID: "counter", MType: metric.Counter, Delta: int64Ptr(12)},
		}}
	)
	tests := []struct {
		ms    *MemStorage
		mName string
	}{
		{
			ms:    ms1,
			mName: "gauge",
		},
	}
	for _, test := range tests {
		_, err := test.ms.GetMetric(context.Background(), test.mName)
		require.NotNil(t, err)
	}
}
