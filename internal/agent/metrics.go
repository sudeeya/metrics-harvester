package agent

import (
	"math/rand/v2"
	"runtime"
	"sync"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/sudeeya/metrics-harvester/internal/metric"
)

type Metrics struct {
	mutex  sync.RWMutex
	values map[string]*metric.Metric
}

func NewMetrics() *Metrics {
	return &Metrics{
		values: map[string]*metric.Metric{
			"Alloc":           {ID: "Alloc", MType: metric.Gauge, Value: new(float64)},
			"BuckHashSys":     {ID: "BuckHashSys", MType: metric.Gauge, Value: new(float64)},
			"CPUutilization1": {ID: "CPUutilization1", MType: metric.Gauge, Value: new(float64)},
			"FreeMemory":      {ID: "FreeMemory", MType: metric.Gauge, Value: new(float64)},
			"Frees":           {ID: "Frees", MType: metric.Gauge, Value: new(float64)},
			"GCCPUFraction":   {ID: "GCCPUFraction", MType: metric.Gauge, Value: new(float64)},
			"GCSys":           {ID: "GCSys", MType: metric.Gauge, Value: new(float64)},
			"HeapAlloc":       {ID: "HeapAlloc", MType: metric.Gauge, Value: new(float64)},
			"HeapIdle":        {ID: "HeapIdle", MType: metric.Gauge, Value: new(float64)},
			"HeapInuse":       {ID: "HeapInuse", MType: metric.Gauge, Value: new(float64)},
			"HeapObjects":     {ID: "HeapObjects", MType: metric.Gauge, Value: new(float64)},
			"HeapReleased":    {ID: "HeapReleased", MType: metric.Gauge, Value: new(float64)},
			"HeapSys":         {ID: "HeapSys", MType: metric.Gauge, Value: new(float64)},
			"LastGC":          {ID: "LastGC", MType: metric.Gauge, Value: new(float64)},
			"Lookups":         {ID: "Lookups", MType: metric.Gauge, Value: new(float64)},
			"MCacheInuse":     {ID: "MCacheInuse", MType: metric.Gauge, Value: new(float64)},
			"MCacheSys":       {ID: "MCacheSys", MType: metric.Gauge, Value: new(float64)},
			"MSpanInuse":      {ID: "MSpanInuse", MType: metric.Gauge, Value: new(float64)},
			"MSpanSys":        {ID: "MSpanSys", MType: metric.Gauge, Value: new(float64)},
			"Mallocs":         {ID: "Mallocs", MType: metric.Gauge, Value: new(float64)},
			"NextGC":          {ID: "NextGC", MType: metric.Gauge, Value: new(float64)},
			"NumForcedGC":     {ID: "NumForcedGC", MType: metric.Gauge, Value: new(float64)},
			"NumGC":           {ID: "NumGC", MType: metric.Gauge, Value: new(float64)},
			"OtherSys":        {ID: "OtherSys", MType: metric.Gauge, Value: new(float64)},
			"PauseTotalNs":    {ID: "PauseTotalNs", MType: metric.Gauge, Value: new(float64)},
			"PollCount":       {ID: "PollCount", MType: metric.Counter, Delta: new(int64)},
			"RandomValue":     {ID: "RandomValue", MType: metric.Gauge, Value: new(float64)},
			"StackInuse":      {ID: "StackInuse", MType: metric.Gauge, Value: new(float64)},
			"StackSys":        {ID: "StackSys", MType: metric.Gauge, Value: new(float64)},
			"Sys":             {ID: "Sys", MType: metric.Gauge, Value: new(float64)},
			"TotalAlloc":      {ID: "TotalAlloc", MType: metric.Gauge, Value: new(float64)},
			"TotalMemory":     {ID: "TotalMemory", MType: metric.Gauge, Value: new(float64)},
		},
	}
}

func (m *Metrics) List() []metric.Metric {
	metrics := make([]metric.Metric, 0)
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	for _, m := range m.values {
		metrics = append(metrics, *m)
	}
	return metrics
}

func (m *Metrics) Update() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.values["Alloc"].Update(float64(memStats.Alloc))
	m.values["BuckHashSys"].Update(float64(memStats.BuckHashSys))
	m.values["Frees"].Update(float64(memStats.Frees))
	m.values["GCCPUFraction"].Update(float64(memStats.GCCPUFraction))
	m.values["GCSys"].Update(float64(memStats.GCSys))
	m.values["HeapAlloc"].Update(float64(memStats.HeapAlloc))
	m.values["HeapIdle"].Update(float64(memStats.HeapIdle))
	m.values["HeapInuse"].Update(float64(memStats.HeapInuse))
	m.values["HeapObjects"].Update(float64(memStats.HeapObjects))
	m.values["HeapReleased"].Update(float64(memStats.HeapReleased))
	m.values["HeapSys"].Update(float64(memStats.HeapSys))
	m.values["LastGC"].Update(float64(memStats.LastGC))
	m.values["Lookups"].Update(float64(memStats.Lookups))
	m.values["MCacheInuse"].Update(float64(memStats.MCacheInuse))
	m.values["MCacheSys"].Update(float64(memStats.MCacheSys))
	m.values["MSpanInuse"].Update(float64(memStats.MSpanInuse))
	m.values["MSpanSys"].Update(float64(memStats.MSpanSys))
	m.values["Mallocs"].Update(float64(memStats.Mallocs))
	m.values["NextGC"].Update(float64(memStats.NextGC))
	m.values["NumForcedGC"].Update(float64(memStats.NumForcedGC))
	m.values["NumGC"].Update(float64(memStats.NumGC))
	m.values["OtherSys"].Update(float64(memStats.OtherSys))
	m.values["PauseTotalNs"].Update(float64(memStats.PauseTotalNs))
	m.values["PollCount"].Update(int64(1))
	m.values["RandomValue"].Update(rand.Float64())
	m.values["StackInuse"].Update(float64(memStats.StackInuse))
	m.values["StackSys"].Update(float64(memStats.StackSys))
	m.values["Sys"].Update(float64(memStats.Sys))
	m.values["TotalAlloc"].Update(float64(memStats.TotalAlloc))
}

func (m *Metrics) UpdatePSUtil() error {
	memStats, err := mem.VirtualMemory()
	if err != nil {
		return err
	}
	cpuStats, err := cpu.Percent(0, false)
	if err != nil {
		return err
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.values["CPUutilization1"].Update(cpuStats[0])
	m.values["FreeMemory"].Update(float64(memStats.Free))
	m.values["TotalMemory"].Update(float64(memStats.Total))
	return nil
}
