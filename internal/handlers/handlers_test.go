package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/sudeeya/metrics-harvester/internal/metric"
	"github.com/sudeeya/metrics-harvester/internal/mocks"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	require.NoError(t, err)
	response, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer response.Body.Close()
	responseBody, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	return response, strings.TrimSpace(string(responseBody))
}

func int64Ptr(i int64) *int64 {
	return &i
}

func float64Ptr(f float64) *float64 {
	return &f
}

func TestAllMetricsHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repoMock := mocks.NewMockRepository(ctrl)
	metrics := []metric.Metric{
		{ID: "gauge", MType: metric.Gauge, Value: float64Ptr(12.12)},
		{ID: "counter", MType: metric.Counter, Delta: int64Ptr(12)},
	}
	repoMock.EXPECT().
		GetAllMetrics(gomock.Any()).
		Return(metrics, nil)

	logger := zap.NewNop()
	router := chi.NewRouter()
	router.Get("/", NewAllMetricsHandler(logger, repoMock))
	ts := httptest.NewServer(router)
	defer ts.Close()

	type result struct {
		code int
	}
	tests := []struct {
		name   string
		path   string
		result result
	}{
		{
			name: "simple test",
			path: "/",
			result: result{
				code: http.StatusOK,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response, _ := testRequest(t, ts, "GET", test.path, nil)
			defer response.Body.Close()
			require.Equal(t, test.result.code, response.StatusCode)
		})
	}
}

func TestValueHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repoMock := mocks.NewMockRepository(ctrl)
	gauge := metric.Metric{ID: "gauge", MType: metric.Gauge, Value: float64Ptr(12.12)}
	counter := metric.Metric{ID: "counter", MType: metric.Counter, Delta: int64Ptr(12)}
	repoMock.EXPECT().
		GetMetric(gomock.Any(), "gauge").
		Return(gauge, nil)
	repoMock.EXPECT().
		GetMetric(gomock.Any(), "counter").
		Return(counter, nil)
	repoMock.EXPECT().
		GetMetric(gomock.Any(), "dummy").
		Return(metric.Metric{}, errors.New("dummy"))

	logger := zap.NewNop()
	router := chi.NewRouter()
	router.Get("/value/{metricType}/{metricName}", NewValueHandler(logger, repoMock))
	ts := httptest.NewServer(router)
	defer ts.Close()

	type result struct {
		code int
		body string
	}
	tests := []struct {
		name   string
		path   string
		result result
	}{
		{
			name: "get gauge",
			path: "/value/gauge/gauge",
			result: result{
				code: http.StatusOK,
				body: "12.12",
			},
		},
		{
			name: "get counter",
			path: "/value/counter/counter",
			result: result{
				code: http.StatusOK,
				body: "12",
			},
		},
		{
			name: "try to get non-existent metric",
			path: "/value/gauge/dummy",
			result: result{
				code: http.StatusNotFound,
				body: "",
			},
		},
		{
			name: "try to get metric of non-existent type",
			path: "/value/dummy/dummy",
			result: result{
				code: http.StatusBadRequest,
				body: "",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response, body := testRequest(t, ts, "GET", test.path, nil)
			defer response.Body.Close()
			require.Equal(t, test.result.code, response.StatusCode)
			require.Equal(t, test.result.body, body)
		})
	}
}

func TestUpdateHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repoMock := mocks.NewMockRepository(ctrl)
	gauge := metric.Metric{ID: "gauge", MType: metric.Gauge, Value: float64Ptr(12.12)}
	counter := metric.Metric{ID: "counter", MType: metric.Counter, Delta: int64Ptr(12)}
	repoMock.EXPECT().
		PutMetric(gomock.Any(), gauge).
		Return(nil)
	repoMock.EXPECT().
		PutMetric(gomock.Any(), counter).
		Return(nil)

	logger := zap.NewNop()
	router := chi.NewRouter()
	router.Post("/update/{metricType}/{metricName}/{metricValue}", NewUpdateHandler(logger, repoMock))
	ts := httptest.NewServer(router)
	defer ts.Close()

	type result struct {
		code int
	}
	tests := []struct {
		name   string
		path   string
		result result
	}{
		{
			name: "update counter",
			path: "/update/counter/counter/12",
			result: result{
				code: http.StatusOK,
			},
		},
		{
			name: "update gauge",
			path: "/update/gauge/gauge/12.12",
			result: result{
				code: http.StatusOK,
			},
		},
		{
			name: "try to update metric of non-existent type",
			path: "/update/dummy/dummy/12",
			result: result{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "try to update counter with float",
			path: "/update/counter/counter/12.12",
			result: result{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "try to update gauge with string",
			path: "/update/gauge/gauge/dummy",
			result: result{
				code: http.StatusBadRequest,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response, _ := testRequest(t, ts, "POST", test.path, nil)
			defer response.Body.Close()
			require.Equal(t, test.result.code, response.StatusCode)
		})
	}
}

func TestJSONUpdateHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repoMock := mocks.NewMockRepository(ctrl)
	gauge := metric.Metric{ID: "gauge", MType: metric.Gauge, Value: float64Ptr(12.12)}
	counter := metric.Metric{ID: "counter", MType: metric.Counter, Delta: int64Ptr(12)}
	gomock.InOrder(
		repoMock.EXPECT().
			PutMetric(gomock.Any(), gauge).
			Return(nil),
		repoMock.EXPECT().
			GetMetric(gomock.Any(), "gauge").
			Return(gauge, nil),
	)
	gomock.InOrder(
		repoMock.EXPECT().
			PutMetric(gomock.Any(), counter).
			Return(nil),
		repoMock.EXPECT().
			GetMetric(gomock.Any(), "counter").
			Return(counter, nil),
	)

	logger := zap.NewNop()
	router := chi.NewRouter()
	router.Post("/update/", NewJSONUpdateHandler(logger, repoMock))
	ts := httptest.NewServer(router)
	defer ts.Close()

	jsonCounter, _ := json.Marshal(counter)
	jsonGauge, _ := json.Marshal(gauge)
	type result struct {
		code int
		body string
	}
	tests := []struct {
		name   string
		path   string
		body   io.Reader
		result result
	}{
		{
			name: "update counter",
			path: "/update/",
			body: bytes.NewReader(jsonCounter),
			result: result{
				code: http.StatusOK,
				body: string(jsonCounter),
			},
		},
		{
			name: "update gauge",
			path: "/update/",
			body: bytes.NewReader(jsonGauge),
			result: result{
				code: http.StatusOK,
				body: string(jsonGauge),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response, body := testRequest(t, ts, "POST", test.path, test.body)
			defer response.Body.Close()
			require.Equal(t, test.result.code, response.StatusCode)
			require.Equal(t, test.result.body, body)
		})
	}
}

func TestBatchHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repoMock := mocks.NewMockRepository(ctrl)
	metrics := []metric.Metric{
		{ID: "gauge", MType: metric.Gauge, Value: float64Ptr(12.12)},
		{ID: "counter", MType: metric.Counter, Delta: int64Ptr(12)},
	}
	repoMock.EXPECT().
		PutBatch(gomock.Any(), metrics).
		Return(nil)

	logger := zap.NewNop()
	router := chi.NewRouter()
	router.Post("/updates/", NewBatchHandler(logger, repoMock))
	ts := httptest.NewServer(router)
	defer ts.Close()

	jsonMetrics, _ := json.Marshal(metrics)
	type result struct {
		code int
	}
	tests := []struct {
		name   string
		path   string
		body   io.Reader
		result result
	}{
		{
			name: "simple test",
			path: "/updates/",
			body: bytes.NewReader(jsonMetrics),
			result: result{
				code: http.StatusOK,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response, _ := testRequest(t, ts, "POST", test.path, test.body)
			defer response.Body.Close()
			require.Equal(t, test.result.code, response.StatusCode)
		})
	}
}

func TestJSONValueHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repoMock := mocks.NewMockRepository(ctrl)
	gauge := metric.Metric{ID: "gauge", MType: metric.Gauge, Value: float64Ptr(12.12)}
	counter := metric.Metric{ID: "counter", MType: metric.Counter, Delta: int64Ptr(12)}
	repoMock.EXPECT().
		GetMetric(gomock.Any(), "gauge").
		Return(gauge, nil)
	repoMock.EXPECT().
		GetMetric(gomock.Any(), "counter").
		Return(counter, nil)

	logger := zap.NewNop()
	router := chi.NewRouter()
	router.Post("/value/", NewJSONValueHandler(logger, repoMock))
	ts := httptest.NewServer(router)
	defer ts.Close()

	jsonRawGauge, _ := json.Marshal(metric.Metric{ID: "gauge", MType: metric.Gauge})
	jsonRawCounter, _ := json.Marshal(metric.Metric{ID: "counter", MType: metric.Counter})
	jsonGauge, _ := json.Marshal(gauge)
	jsonCounter, _ := json.Marshal(counter)
	type result struct {
		code int
		body string
	}
	tests := []struct {
		name   string
		path   string
		body   io.Reader
		result result
	}{
		{
			name: "get gauge",
			path: "/value/",
			body: bytes.NewReader(jsonRawGauge),
			result: result{
				code: http.StatusOK,
				body: string(jsonGauge),
			},
		},
		{
			name: "get counter",
			path: "/value/",
			body: bytes.NewReader(jsonRawCounter),
			result: result{
				code: http.StatusOK,
				body: string(jsonCounter),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response, body := testRequest(t, ts, "POST", test.path, test.body)
			defer response.Body.Close()
			require.Equal(t, test.result.code, response.StatusCode)
			require.Equal(t, test.result.body, body)
		})
	}
}
