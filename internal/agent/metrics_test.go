package agent

import "testing"

func BenchmarkMetricsUpdate(b *testing.B) {
	b.StopTimer()
	metrics := NewMetrics()
	b.StartTimer()
	b.Run("Update", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			metrics.Update()
		}
	})
	b.Run("UpdatePSUtil", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if err := metrics.UpdatePSUtil(); err != nil {
				b.Error(err)
			}
		}
	})
}
