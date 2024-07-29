package metrics

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChangeValue_Basic(t *testing.T) {
	tests := []struct {
		gauge    *Gauge
		newValue float64
		result   *Gauge
	}{
		{
			gauge:    NewGauge("gauge", 0),
			newValue: 12.12,
			result:   NewGauge("gauge", 12.12),
		},
	}
	for _, test := range tests {
		test.gauge.ChangeValue(test.newValue)
		require.EqualValues(t, test.gauge, test.result)
	}
}

func TestChangeValue_Stress(t *testing.T) {
	var (
		n      int64  = 1000000
		gauge  *Gauge = NewGauge("gauge", 0)
		result *Gauge = NewGauge("gauge", float64(n-1))
	)
	var i int64
	for i = 0; i < n; i++ {
		gauge.ChangeValue(float64(i))
	}
	require.EqualValues(t, gauge, result)
}
