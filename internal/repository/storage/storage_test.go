package storage

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/sudeeya/metrics-harvester/internal/metric"
	"github.com/sudeeya/metrics-harvester/internal/utils"
)

func TestPutMetric(t *testing.T) {
	var (
		ms1 = &MemStorage{metrics: map[string]metric.Metric{
			"gauge": {ID: "gauge", MType: metric.Gauge, Value: utils.Float64Ptr(12)},
		}}
		ms2 = &MemStorage{metrics: map[string]metric.Metric{
			"gauge": {ID: "gauge", MType: metric.Gauge, Value: utils.Float64Ptr(12.12)},
		}}
		ms3 = &MemStorage{metrics: map[string]metric.Metric{
			"counter": {ID: "counter", MType: metric.Counter, Delta: utils.Int64Ptr(12)},
		}}
		ms4 = &MemStorage{metrics: map[string]metric.Metric{
			"counter": {ID: "counter", MType: metric.Counter, Delta: utils.Int64Ptr(24)},
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
			m:      metric.Metric{ID: "gauge", MType: metric.Gauge, Value: utils.Float64Ptr(12)},
			result: ms1,
		},
		{
			ms:     ms1,
			mName:  "gauge",
			m:      metric.Metric{ID: "gauge", MType: metric.Gauge, Value: utils.Float64Ptr(12.12)},
			result: ms2,
		},
		{
			ms:     NewMemStorage(),
			mName:  "counter",
			m:      metric.Metric{ID: "counter", MType: metric.Counter, Delta: utils.Int64Ptr(12)},
			result: ms3,
		},
		{
			ms:     ms3,
			mName:  "counter",
			m:      metric.Metric{ID: "counter", MType: metric.Counter, Delta: utils.Int64Ptr(12)},
			result: ms4,
		},
	}
	for _, test := range tests {
		test.ms.PutMetric(test.m)
		require.EqualValues(t, test.result.metrics[test.mName].GetValue(), test.ms.metrics[test.mName].GetValue())
	}
}

func TestGetMetric_Existing(t *testing.T) {
	var (
		ms1 = &MemStorage{metrics: map[string]metric.Metric{
			"gauge": {ID: "gauge", MType: metric.Gauge, Value: utils.Float64Ptr(12)},
		}}
		ms2 = &MemStorage{metrics: map[string]metric.Metric{
			"counter": {ID: "counter", MType: metric.Counter, Delta: utils.Int64Ptr(12)},
			"dummy":   {ID: "dummy", MType: metric.Gauge, Value: utils.Float64Ptr(-1)},
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
			result: metric.Metric{ID: "gauge", MType: metric.Gauge, Value: utils.Float64Ptr(12)},
		},
		{
			ms:     ms2,
			mName:  "counter",
			result: metric.Metric{ID: "counter", MType: metric.Counter, Delta: utils.Int64Ptr(12)},
		},
	}
	for _, test := range tests {
		m, ok := test.ms.GetMetric(test.mName)
		require.Equal(t, true, ok)
		require.EqualValues(t, test.result.GetValue(), m.GetValue())
	}
}

func TestGetMetric_NotExisting(t *testing.T) {
	var (
		ms1 = &MemStorage{metrics: map[string]metric.Metric{
			"counter": {ID: "counter", MType: metric.Counter, Delta: utils.Int64Ptr(12)},
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
		_, ok := test.ms.GetMetric(test.mName)
		require.Equal(t, false, ok)
	}
}
