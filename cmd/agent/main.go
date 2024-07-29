package main

import (
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"reflect"
	"runtime"
	"time"
)

const serverAddress string = "http://localhost:8080"

var pollInterval time.Duration = 2 * time.Second

var typesOfMetrics = map[string]string{
	"Alloc":         "gauge",
	"BuckHashSys":   "gauge",
	"Frees":         "gauge",
	"GCCPUFraction": "gauge",
	"GCSys":         "gauge",
	"HeapAlloc":     "gauge",
	"HeapIdle":      "gauge",
	"HeapInuse":     "gauge",
	"HeapObjects":   "gauge",
	"HeapReleased":  "gauge",
	"HeapSys":       "gauge",
	"LastGC":        "gauge",
	"Lookups":       "gauge",
	"MCacheInuse":   "gauge",
	"MCacheSys":     "gauge",
	"MSpanInuse":    "gauge",
	"MSpanSys":      "gauge",
	"Mallocs":       "gauge",
	"NextGC":        "gauge",
	"NumForcedGC":   "gauge",
	"NumGC":         "gauge",
	"OtherSys":      "gauge",
	"PauseTotalNs":  "gauge",
	"PollCount":     "counter",
	"RandomValue":   "gauge",
	"StackInuse":    "gauge",
	"StackSys":      "gauge",
	"Sys":           "gauge",
	"TotalAlloc":    "gauge",
}

type Metrics struct {
	memStats    *runtime.MemStats
	pollCount   int64
	countFunc   func(*runtime.MemStats) int64
	randomValue float64
}

func NewMetrics(memStats *runtime.MemStats) *Metrics {
	return &Metrics{
		memStats:  memStats,
		countFunc: countCall(runtime.ReadMemStats),
	}
}

func countCall(f func(*runtime.MemStats)) func(*runtime.MemStats) int64 {
	count := int64(0)
	return func(ms *runtime.MemStats) int64 {
		count++
		f(ms)
		return count
	}
}

func (m *Metrics) Update() {
	m.pollCount = m.countFunc(m.memStats)
	m.randomValue = rand.Float64()
}

func main() {
	var (
		memStats runtime.MemStats
		client   = &http.Client{}
		metrics  = NewMetrics(&memStats)
	)
	for {
		conductPollCycle(metrics)
		sendMetrics(metrics, client)
	}
}

func conductPollCycle(metrics *Metrics) {
	for i := 0; i < 5; i++ {
		time.Sleep(pollInterval)
		metrics.Update()
	}
}

func sendMetrics(metrics *Metrics, client *http.Client) {
	for metricName, metricType := range typesOfMetrics {
		var metricValue string
		switch metricName {
		case "PollCount":
			metricValue = fmt.Sprintf("%v", metrics.pollCount)
		case "RandomValue":
			metricValue = fmt.Sprintf("%v", metrics.randomValue)
		default:
			metricValue = reflect.ValueOf(*metrics.memStats).FieldByName(metricName).String()
		}
		path := formPath(metricType, metricName, metricValue)
		response, err := client.Post(path, "text/plain", nil)
		if err != nil {
			panic(err)
		}
		defer response.Body.Close()
		_, err = io.Copy(io.Discard, response.Body)
		if err != nil {
			panic(err)
		}
	}
}

func formPath(metricType, metricName, metricValue string) string {
	return fmt.Sprintf("%s/update/%s/%s/%s",
		serverAddress, metricType, metricName, metricValue)
}
