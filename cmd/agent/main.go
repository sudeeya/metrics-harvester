package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"time"
)

const serverAddress string = "http://localhost:8080"

var (
	pollInterval time.Duration = 2 * time.Second
)

var metrics = map[string]string{
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

func countCall(f func(*runtime.MemStats)) func(*runtime.MemStats) int {
	count := 0
	return func(ms *runtime.MemStats) int {
		count++
		f(ms)
		return count
	}
}

func formPath(metricType, metricName, metricValue string) string {
	return fmt.Sprintf("%s/update/%s/%s/%s",
		serverAddress, metricType, metricName, metricValue)
}

func main() {
	var (
		memStats          runtime.MemStats
		client            = &http.Client{}
		countReadMemStats = countCall(runtime.ReadMemStats)
		pollCount         int
	)
	for {
		for i := 0; i < 5; i++ {
			time.Sleep(pollInterval)
			pollCount = countReadMemStats(&memStats)
		}
		for metricName, metricType := range metrics {
			var metricValue string
			switch metricName {
			case "PollCount":
				metricValue = fmt.Sprintf("%v", pollCount)
			case "RandomValue":
				metricValue = fmt.Sprintf("%v", rand.Float64())
			default:
				metricValue = reflect.ValueOf(memStats).FieldByName(metricName).String()
			}
			path := formPath(metricType, metricName, metricValue)
			response, err := client.Post(path, "text/plain", nil)
			if err != nil {
				fmt.Println(err)
			}
			_, err = io.Copy(io.Discard, response.Body)
			response.Body.Close()
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
