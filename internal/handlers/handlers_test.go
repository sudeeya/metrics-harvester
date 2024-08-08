package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"github.com/sudeeya/metrics-harvester/internal/metric"
	"github.com/sudeeya/metrics-harvester/internal/repository/storage"
	"github.com/sudeeya/metrics-harvester/internal/utils"
	"go.uber.org/zap"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)
	response, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	return response, string(body)
}

func TestGetAllMetricsHandler(t *testing.T) {
	var (
		ms = storage.NewMemStorage()
		l  = zap.NewNop()
		ts = httptest.NewServer(CreateGetAllMetricsHandler(l, ms))
	)
	defer ts.Close()
	ms.PutMetric(metric.Metric{ID: "gauge", MType: metric.Gauge, Value: utils.Float64Ptr(12.12)})
	ms.PutMetric(metric.Metric{ID: "counter", MType: metric.Counter, Delta: utils.Int64Ptr(12)})
	type result struct {
		code int
		body string
	}
	tests := []struct {
		result result
	}{
		{
			result: result{
				code: http.StatusOK,
				body: "counter: 12\ngauge: 12.12",
			},
		},
	}
	for _, test := range tests {
		response, body := testRequest(t, ts, "GET", "/")
		defer response.Body.Close()
		require.Equal(t, response.StatusCode, test.result.code)
		require.Equal(t, body, test.result.body)
	}
}

func TestGetMetricHandler(t *testing.T) {
	var (
		ms     = storage.NewMemStorage()
		l      = zap.NewNop()
		router = chi.NewRouter()
		ts     = httptest.NewServer(router)
	)
	defer ts.Close()
	ms.PutMetric(metric.Metric{ID: "gauge", MType: metric.Gauge, Value: utils.Float64Ptr(12.12)})
	ms.PutMetric(metric.Metric{ID: "counter", MType: metric.Counter, Delta: utils.Int64Ptr(12)})
	router.Get("/value/{metricType}/{metricName}", CreateGetMetricHandler(l, ms))
	type result struct {
		code int
		body string
	}
	tests := []struct {
		path   string
		result result
	}{
		{
			path: "/value/counter/counter",
			result: result{
				code: http.StatusOK,
				body: "{\"id\":\"counter\",\"type\":\"counter\",\"delta\":12}\n",
			},
		},
		{
			path: "/value/gauge/dummy",
			result: result{
				code: http.StatusNotFound,
				body: "",
			},
		},
		{
			path: "/value/dummy/dummy",
			result: result{
				code: http.StatusBadRequest,
				body: "",
			},
		},
	}
	for _, test := range tests {
		response, body := testRequest(t, ts, "GET", test.path)
		defer response.Body.Close()
		require.Equal(t, test.result.code, response.StatusCode)
		require.Equal(t, test.result.body, body)
	}
}

func TestPostMetricHandler(t *testing.T) {
	var (
		ms     = storage.NewMemStorage()
		l      = zap.NewNop()
		router = chi.NewRouter()
		ts     = httptest.NewServer(router)
	)
	defer ts.Close()
	router.Post("/update/{metricType}/{metricName}/{metricValue}", CreatePostMetricHandler(l, ms))
	type result struct {
		code int
	}
	tests := []struct {
		path   string
		result result
	}{
		{
			path: "/update/counter/counter/12",
			result: result{
				code: http.StatusOK,
			},
		},
		{
			path: "/update/gauge/gauge/12.12",
			result: result{
				code: http.StatusOK,
			},
		},
		{
			path: "/update/dummy/dummy/12",
			result: result{
				code: http.StatusBadRequest,
			},
		},
		{
			path: "/update/counter/counter/12.12",
			result: result{
				code: http.StatusBadRequest,
			},
		},
		{
			path: "/update/gauge/gauge/dummy",
			result: result{
				code: http.StatusBadRequest,
			},
		},
	}
	for _, test := range tests {
		response, _ := testRequest(t, ts, "POST", test.path)
		defer response.Body.Close()
		require.Equal(t, test.result.code, response.StatusCode)
	}
}
