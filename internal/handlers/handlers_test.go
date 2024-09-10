package handlers

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/sudeeya/metrics-harvester/internal/metric"
	"github.com/sudeeya/metrics-harvester/internal/mocks"
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

func int64Ptr(i int64) *int64 {
	return &i
}

func float64Ptr(f float64) *float64 {
	return &f
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
			response, body := testRequest(t, ts, "GET", test.path)
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
			response, _ := testRequest(t, ts, "POST", test.path)
			defer response.Body.Close()
			require.Equal(t, test.result.code, response.StatusCode)
		})
	}
}
