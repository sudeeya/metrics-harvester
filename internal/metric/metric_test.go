package metric

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/sudeeya/metrics-harvester/internal/utils"
)

func TestUpdate(t *testing.T) {
	tests := []struct {
		m      *Metric
		update any
		result *Metric
	}{
		{
			m:      &Metric{ID: "gauge", MType: Gauge, Value: utils.Float64Ptr(12)},
			update: float64(12.12),
			result: &Metric{ID: "gauge", MType: Gauge, Value: utils.Float64Ptr(12.12)},
		},
		{
			m:      &Metric{ID: "counter", MType: Counter, Delta: utils.Int64Ptr(12)},
			update: int64(12),
			result: &Metric{ID: "counter", MType: Counter, Delta: utils.Int64Ptr(24)},
		},
	}
	for _, test := range tests {
		test.m.Update(test.update)
		require.EqualValues(t, test.result.GetValue(), test.m.GetValue())
	}
}
