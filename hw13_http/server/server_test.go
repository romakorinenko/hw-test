package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRequestHandler(t *testing.T) {
	testCases := []struct {
		method         string
		path           string
		expectedStatus int
		expectedBody   string
	}{
		{
			method:         "GET",
			path:           "/",
			expectedStatus: http.StatusOK,
			expectedBody:   "GET response",
		},
		{
			method:         "POST",
			path:           "/",
			expectedStatus: http.StatusOK,
			expectedBody:   "POST response",
		},
		{
			method:         "PATCH",
			path:           "/",
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "method not allowed\n",
		},
	}

	for _, testCase := range testCases {
		req := httptest.NewRequest(testCase.method, testCase.path, nil)
		w := httptest.NewRecorder()
		requestHandler(w, req)

		resp := w.Result()
		_ = resp.Body.Close()
		require.Equal(t, testCase.expectedStatus, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, testCase.expectedBody, string(body))
	}
}
