package agent

import (
	"math/rand/v2"
	"runtime"

	"github.com/sudeeya/metrics-harvester/internal/metric"
)

type Metrics struct {
	Values map[string]*metric.Metric
}

func NewMetrics() *Metrics {
	return &Metrics{
		Values: map[string]*metric.Metric{
			"Alloc":         {ID: "Alloc", MType: metric.Gauge, Value: new(float64)},
			"BuckHashSys":   {ID: "BuckHashSys", MType: metric.Gauge, Value: new(float64)},
			"Frees":         {ID: "Frees", MType: metric.Gauge, Value: new(float64)},
			"GCCPUFraction": {ID: "GCCPUFraction", MType: metric.Gauge, Value: new(float64)},
			"GCSys":         {ID: "GCSys", MType: metric.Gauge, Value: new(float64)},
			"HeapAlloc":     {ID: "HeapAlloc", MType: metric.Gauge, Value: new(float64)},
			"HeapIdle":      {ID: "HeapIdle", MType: metric.Gauge, Value: new(float64)},
			"HeapInuse":     {ID: "HeapInuse", MType: metric.Gauge, Value: new(float64)},
			"HeapObjects":   {ID: "HeapObjects", MType: metric.Gauge, Value: new(float64)},
			"HeapReleased":  {ID: "HeapReleased", MType: metric.Gauge, Value: new(float64)},
			"HeapSys":       {ID: "HeapSys", MType: metric.Gauge, Value: new(float64)},
			"LastGC":        {ID: "LastGC", MType: metric.Gauge, Value: new(float64)},
			"Lookups":       {ID: "Lookups", MType: metric.Gauge, Value: new(float64)},
			"MCacheInuse":   {ID: "MCacheInuse", MType: metric.Gauge, Value: new(float64)},
			"MCacheSys":     {ID: "MCacheSys", MType: metric.Gauge, Value: new(float64)},
			"MSpanInuse":    {ID: "MSpanInuse", MType: metric.Gauge, Value: new(float64)},
			"MSpanSys":      {ID: "MSpanSys", MType: metric.Gauge, Value: new(float64)},
			"Mallocs":       {ID: "Mallocs", MType: metric.Gauge, Value: new(float64)},
			"NextGC":        {ID: "NextGC", MType: metric.Gauge, Value: new(float64)},
			"NumForcedGC":   {ID: "NumForcedGC", MType: metric.Gauge, Value: new(float64)},
			"NumGC":         {ID: "NumGC", MType: metric.Gauge, Value: new(float64)},
			"OtherSys":      {ID: "OtherSys", MType: metric.Gauge, Value: new(float64)},
			"PauseTotalNs":  {ID: "PauseTotalNs", MType: metric.Gauge, Value: new(float64)},
			"PollCount":     {ID: "PollCount", MType: metric.Counter, Delta: new(int64)},
			"RandomValue":   {ID: "RandomValue", MType: metric.Gauge, Value: new(float64)},
			"StackInuse":    {ID: "StackInuse", MType: metric.Gauge, Value: new(float64)},
			"StackSys":      {ID: "StackSys", MType: metric.Gauge, Value: new(float64)},
			"Sys":           {ID: "Sys", MType: metric.Gauge, Value: new(float64)},
			"TotalAlloc":    {ID: "TotalAlloc", MType: metric.Gauge, Value: new(float64)},
		},
	}
}

func UpdateMetrics(m *Metrics) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	m.Values["Alloc"].Update(float64(memStats.Alloc))
	m.Values["BuckHashSys"].Update(float64(memStats.BuckHashSys))
	m.Values["Frees"].Update(float64(memStats.Frees))
	m.Values["GCCPUFraction"].Update(float64(memStats.GCCPUFraction))
	m.Values["GCSys"].Update(float64(memStats.GCSys))
	m.Values["HeapAlloc"].Update(float64(memStats.HeapAlloc))
	m.Values["HeapIdle"].Update(float64(memStats.HeapIdle))
	m.Values["HeapInuse"].Update(float64(memStats.HeapInuse))
	m.Values["HeapObjects"].Update(float64(memStats.HeapObjects))
	m.Values["HeapReleased"].Update(float64(memStats.HeapReleased))
	m.Values["HeapSys"].Update(float64(memStats.HeapSys))
	m.Values["LastGC"].Update(float64(memStats.LastGC))
	m.Values["Lookups"].Update(float64(memStats.Lookups))
	m.Values["MCacheInuse"].Update(float64(memStats.MCacheInuse))
	m.Values["MCacheSys"].Update(float64(memStats.MCacheSys))
	m.Values["MSpanInuse"].Update(float64(memStats.MSpanInuse))
	m.Values["MSpanSys"].Update(float64(memStats.MSpanSys))
	m.Values["Mallocs"].Update(float64(memStats.Mallocs))
	m.Values["NextGC"].Update(float64(memStats.NextGC))
	m.Values["NumForcedGC"].Update(float64(memStats.NumForcedGC))
	m.Values["NumGC"].Update(float64(memStats.NumGC))
	m.Values["OtherSys"].Update(float64(memStats.OtherSys))
	m.Values["PauseTotalNs"].Update(float64(memStats.PauseTotalNs))
	m.Values["PollCount"].Update(int64(1))
	m.Values["RandomValue"].Update(rand.Float64())
	m.Values["StackInuse"].Update(float64(memStats.StackInuse))
	m.Values["StackSys"].Update(float64(memStats.StackSys))
	m.Values["Sys"].Update(float64(memStats.Sys))
	m.Values["TotalAlloc"].Update(float64(memStats.TotalAlloc))
}
