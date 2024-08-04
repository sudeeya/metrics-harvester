package metric

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIncreaseValue_Basic(t *testing.T) {
	tests := []struct {
		counter  *Counter
		additive int64
		result   *Counter
	}{
		{
			counter:  NewCounter("count", 0),
			additive: 12,
			result:   NewCounter("count", 12),
		},
	}
	for _, test := range tests {
		test.counter.IncreaseValue(test.additive)
		require.EqualValues(t, test.result, test.counter)
	}
}

func TestIncreaseValue_Stress(t *testing.T) {
	var (
		n       int64 = 1000000
		counter       = NewCounter("count", 0)
		result        = NewCounter("count", n)
	)
	var i int64
	for i = 0; i < n; i++ {
		counter.IncreaseValue(int64(1))
	}
	require.EqualValues(t, result, counter)
}
