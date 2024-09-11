package handlers

import (
	"fmt"
	"net/http/httptest"

	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/sudeeya/metrics-harvester/internal/metric"
	"github.com/sudeeya/metrics-harvester/internal/repository/storage"
	"go.uber.org/zap"
)

func Example() {
	var (
		logger     = zap.NewNop() // without logging
		memStorage = storage.NewMemStorage()
		router     = chi.NewRouter()
	)
	defer memStorage.Close()

	// Adding endpoints
	router.Post("/update/{metricType}/{metricName}/{metricValue}", NewUpdateHandler(logger, memStorage))
	router.Post("/update/", NewJSONUpdateHandler(logger, memStorage))
	router.Post("/updates/", NewBatchHandler(logger, memStorage))
	router.Post("/value/", NewJSONValueHandler(logger, memStorage))
	router.Get("/value/{metricType}/{metricName}", NewValueHandler(logger, memStorage))

	// Creating server instance
	serverExample := httptest.NewServer(router)
	defer serverExample.Close()

	// Creating HTTP client
	client := resty.New().SetBaseURL(serverExample.URL)

	// Endpoint /update/{metricType}/{metricName}/{metricValue}
	client.R().
		Post("/update/gauge/g1/-1")

	// Endpoint /update/
	var (
		c1Delta int64 = 12
		c1            = metric.Metric{ID: "c1", MType: metric.Counter, Delta: &c1Delta}
	)
	client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(c1).
		Post("/update/")

	// Endpoint /updates/
	var (
		c2Delta int64   = 42
		g1Value float64 = 12.12
		batch           = []metric.Metric{
			{ID: "c2", MType: metric.Counter, Delta: &c2Delta},
			{ID: "g1", MType: metric.Gauge, Value: &g1Value},
		}
	)
	client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(batch).
		Post("/updates/")

	// Endpoint /value/{metricType}/{metricName}/
	response, _ := client.R().
		Get("/value/counter/c1")
	fmt.Println(string(response.Body()))

	// Endpoint /value/
	var c2Result metric.Metric
	client.R().
		SetBody(metric.Metric{ID: "c2", MType: "counter"}).
		SetResult(&c2Result).
		Post("/value/")
	fmt.Println(c2Result.GetValue())

	var g1Result metric.Metric
	client.R().
		SetBody(metric.Metric{ID: "g1", MType: "gauge"}).
		SetResult(&g1Result).
		Post("/value/")
	fmt.Println(g1Result.GetValue())

	// Output:
	// 12
	// 42
	// 12.12
}
