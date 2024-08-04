package agent

import (
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"runtime"
	"time"

	"github.com/sudeeya/metrics-harvester/internal/metric"
)

type Agent struct {
	cfg    *Config
	client *http.Client
}

func NewAgent(cfg *Config) *Agent {
	return &Agent{
		cfg:    cfg,
		client: &http.Client{},
	}
}

func (a *Agent) Run() {
	metrics := NewMetrics()
	for {
		var i int64
		for i = 0; i < a.cfg.ReportInterval/a.cfg.PollInterval; i++ {
			time.Sleep(time.Duration(a.cfg.PollInterval) * time.Second)
			UpdateMetrics(metrics)
		}
		a.SendMetrics(metrics)
	}
}

func (a *Agent) SendMetrics(metrics *Metrics) {
	for _, metricValue := range metrics.Values {
		url := formURL(a, metricValue)
		response, err := a.client.Post(url, "text/plain", nil)
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

func formURL(a *Agent, m metric.Metric) string {
	var metricType string
	switch m.(type) {
	case *metric.Gauge:
		metricType = "gauge"
	case *metric.Counter:
		metricType = "counter"
	}
	return fmt.Sprintf("%s/update/%s/%s/%s",
		a.cfg.Address, metricType, m.GetName(), m.GetValue())
}

type Metrics struct {
	Values map[string]metric.Metric
}

func NewMetrics() *Metrics {
	return &Metrics{
		Values: map[string]metric.Metric{
			"Alloc":         metric.NewGauge("Alloc", 0),
			"BuckHashSys":   metric.NewGauge("BuckHashSys", 0),
			"Frees":         metric.NewGauge("Frees", 0),
			"GCCPUFraction": metric.NewGauge("GCCPUFraction", 0),
			"GCSys":         metric.NewGauge("GCSys", 0),
			"HeapAlloc":     metric.NewGauge("HeapAlloc", 0),
			"HeapIdle":      metric.NewGauge("HeapIdle", 0),
			"HeapInuse":     metric.NewGauge("HeapInuse", 0),
			"HeapObjects":   metric.NewGauge("HeapObjects", 0),
			"HeapReleased":  metric.NewGauge("HeapReleased", 0),
			"HeapSys":       metric.NewGauge("HeapSys", 0),
			"LastGC":        metric.NewGauge("LastGC", 0),
			"Lookups":       metric.NewGauge("Lookups", 0),
			"MCacheInuse":   metric.NewGauge("MCacheInuse", 0),
			"MCacheSys":     metric.NewGauge("MCacheSys", 0),
			"MSpanInuse":    metric.NewGauge("MSpanInuse", 0),
			"MSpanSys":      metric.NewGauge("MSpanSys", 0),
			"Mallocs":       metric.NewGauge("Mallocs", 0),
			"NextGC":        metric.NewGauge("NextGC", 0),
			"NumForcedGC":   metric.NewGauge("NumForcedGC", 0),
			"NumGC":         metric.NewGauge("NumGC", 0),
			"OtherSys":      metric.NewGauge("OtherSys", 0),
			"PauseTotalNs":  metric.NewGauge("PauseTotalNs", 0),
			"PollCount":     metric.NewCounter("PollCount", 0),
			"RandomValue":   metric.NewGauge("RandomValue", 0),
			"StackInuse":    metric.NewGauge("StackInuse", 0),
			"StackSys":      metric.NewGauge("StackSys", 0),
			"Sys":           metric.NewGauge("Sys", 0),
			"TotalAlloc":    metric.NewGauge("TotalAlloc", 0),
		},
	}
}

func UpdateMetrics(m *Metrics) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	m.Values["Alloc"].(*metric.Gauge).ChangeValue(float64(memStats.Alloc))
	m.Values["BuckHashSys"].(*metric.Gauge).ChangeValue(float64(memStats.BuckHashSys))
	m.Values["Frees"].(*metric.Gauge).ChangeValue(float64(memStats.Frees))
	m.Values["GCCPUFraction"].(*metric.Gauge).ChangeValue(float64(memStats.GCCPUFraction))
	m.Values["GCSys"].(*metric.Gauge).ChangeValue(float64(memStats.GCSys))
	m.Values["HeapAlloc"].(*metric.Gauge).ChangeValue(float64(memStats.HeapAlloc))
	m.Values["HeapIdle"].(*metric.Gauge).ChangeValue(float64(memStats.HeapIdle))
	m.Values["HeapInuse"].(*metric.Gauge).ChangeValue(float64(memStats.HeapInuse))
	m.Values["HeapObjects"].(*metric.Gauge).ChangeValue(float64(memStats.HeapObjects))
	m.Values["HeapReleased"].(*metric.Gauge).ChangeValue(float64(memStats.HeapReleased))
	m.Values["HeapSys"].(*metric.Gauge).ChangeValue(float64(memStats.HeapSys))
	m.Values["LastGC"].(*metric.Gauge).ChangeValue(float64(memStats.LastGC))
	m.Values["Lookups"].(*metric.Gauge).ChangeValue(float64(memStats.Lookups))
	m.Values["MCacheInuse"].(*metric.Gauge).ChangeValue(float64(memStats.MCacheInuse))
	m.Values["MCacheSys"].(*metric.Gauge).ChangeValue(float64(memStats.MCacheSys))
	m.Values["MSpanInuse"].(*metric.Gauge).ChangeValue(float64(memStats.MSpanInuse))
	m.Values["MSpanSys"].(*metric.Gauge).ChangeValue(float64(memStats.MSpanSys))
	m.Values["Mallocs"].(*metric.Gauge).ChangeValue(float64(memStats.Mallocs))
	m.Values["NextGC"].(*metric.Gauge).ChangeValue(float64(memStats.NextGC))
	m.Values["NumForcedGC"].(*metric.Gauge).ChangeValue(float64(memStats.NumForcedGC))
	m.Values["NumGC"].(*metric.Gauge).ChangeValue(float64(memStats.NumGC))
	m.Values["OtherSys"].(*metric.Gauge).ChangeValue(float64(memStats.OtherSys))
	m.Values["PauseTotalNs"].(*metric.Gauge).ChangeValue(float64(memStats.PauseTotalNs))
	m.Values["PollCount"].(*metric.Counter).IncreaseValue(int64(1))
	m.Values["RandomValue"].(*metric.Gauge).ChangeValue(rand.Float64())
	m.Values["StackInuse"].(*metric.Gauge).ChangeValue(float64(memStats.StackInuse))
	m.Values["StackSys"].(*metric.Gauge).ChangeValue(float64(memStats.StackSys))
	m.Values["Sys"].(*metric.Gauge).ChangeValue(float64(memStats.Sys))
	m.Values["TotalAlloc"].(*metric.Gauge).ChangeValue(float64(memStats.TotalAlloc))
}
