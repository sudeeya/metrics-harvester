package handlers

/*
func TestMetricHandler(t *testing.T) {
	type result struct {
		code        int
		contentType string
	}
	tests := []struct {
		request *http.Request
		result  result
	}{
		{
			request: httptest.NewRequest(http.MethodPost, "/update/gauge/gauge/12.12", nil),
			result: result{
				code:        200,
				contentType: "text/plain",
			},
		},
		{
			request: httptest.NewRequest(http.MethodPost, "/update/gauge/", nil),
			result: result{
				code:        404,
				contentType: "",
			},
		},
		{
			request: httptest.NewRequest(http.MethodPost, "/update/gauge/gauge/12.12/dummy", nil),
			result: result{
				code:        404,
				contentType: "",
			},
		},
		{
			request: httptest.NewRequest(http.MethodPost, "/update/gauge/gauge/dummy", nil),
			result: result{
				code:        400,
				contentType: "",
			},
		},
		{
			request: httptest.NewRequest(http.MethodGet, "/update/gauge/gauge/12.12", nil),
			result: result{
				code:        404,
				contentType: "",
			},
		},
		{
			request: httptest.NewRequest(http.MethodPost, "/update/counter/counter/12", nil),
			result: result{
				code:        200,
				contentType: "text/plain",
			},
		},
		{
			request: httptest.NewRequest(http.MethodPost, "/update/counter/", nil),
			result: result{
				code:        404,
				contentType: "",
			},
		},
		{
			request: httptest.NewRequest(http.MethodPost, "/update/counter/counter/12/dummy", nil),
			result: result{
				code:        404,
				contentType: "",
			},
		},
		{
			request: httptest.NewRequest(http.MethodPost, "/update/counter/counter/dummy", nil),
			result: result{
				code:        400,
				contentType: "",
			},
		},
		{
			request: httptest.NewRequest(http.MethodGet, "/update/counter/counter/12", nil),
			result: result{
				code:        404,
				contentType: "",
			},
		},
	}
	for _, test := range tests {
		var (
			w       = httptest.NewRecorder()
			storage = storage.NewMemStorage()
			router  = router.NewRouter(storage)
		)
		metricHandler := CreateMetricHandler(router)
		metricHandler(w, test.request)
		response := w.Result()
		assert.Equal(t, response.StatusCode, test.result.code)
		defer response.Body.Close()
		_, err := io.ReadAll(response.Body)
		require.NoError(t, err)
		assert.Equal(t, test.result.contentType, response.Header.Get("content-type"))
	}
}
*/
