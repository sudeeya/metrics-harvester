package metric

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func int64Ptr(i int64) *int64 {
	return &i
}

func float64Ptr(f float64) *float64 {
	return &f
}

func TestUpdate(t *testing.T) {
	tests := []struct {
		m      *Metric
		update any
		result *Metric
	}{
		{
			m:      &Metric{ID: "gauge", MType: Gauge, Value: float64Ptr(12)},
			update: float64(12.12),
			result: &Metric{ID: "gauge", MType: Gauge, Value: float64Ptr(12.12)},
		},
		{
			m:      &Metric{ID: "counter", MType: Counter, Delta: int64Ptr(12)},
			update: int64(12),
			result: &Metric{ID: "counter", MType: Counter, Delta: int64Ptr(24)},
		},
	}
	for _, test := range tests {
		test.m.Update(test.update)
		require.EqualValues(t, test.result.GetValue(), test.m.GetValue())
	}
}
