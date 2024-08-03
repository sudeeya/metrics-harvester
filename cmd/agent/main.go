package main

import (
	"flag"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"reflect"
	"runtime"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/sudeeya/metrics-harvester/internal/agent"
)

var cfg agent.Config

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
	flag.StringVar(&cfg.Address, "a", cfg.Address, "Server IP address and port")
	flag.Int64Var(&cfg.PollInterval, "p", cfg.PollInterval, "Polling interval in seconds")
	flag.Int64Var(&cfg.ReportInterval, "r", cfg.ReportInterval, "Report interval in seconds")
	flag.Parse()
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}
	var (
		memStats runtime.MemStats
		client   = &http.Client{}
		metrics  = NewMetrics(&memStats)
	)
	flag.Parse()
	for {
		conductPollCycle(metrics)
		sendMetrics(metrics, client)
	}
}

func conductPollCycle(metrics *Metrics) {
	var i int64
	for i = 0; i < cfg.ReportInterval/cfg.PollInterval; i++ {
		time.Sleep(time.Duration(cfg.PollInterval) * time.Second)
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
			metricValue = fmt.Sprintf("%v", reflect.ValueOf(*metrics.memStats).FieldByName(metricName).Interface())
		}
		url := formURL(metricType, metricName, metricValue)
		response, err := client.Post(url, "text/plain", nil)
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

func formURL(metricType, metricName, metricValue string) string {
	return fmt.Sprintf("http://%s/update/%s/%s/%s",
		cfg.Address, metricType, metricName, metricValue)
}
