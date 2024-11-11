package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSendRequest(t *testing.T) {
	testCases := []struct {
		method         string
		expectedStatus int
		expectedBody   string
	}{
		{
			method:         "GET",
			expectedStatus: http.StatusOK,
			expectedBody:   "GET response",
		},
		{
			method:         "POST",
			expectedStatus: http.StatusOK,
			expectedBody:   "POST response",
		},
		{
			method:         "PATCH",
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "method not allowed\n",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(requestHandlerForTest))
	defer server.Close()

	ctx := context.Background()
	client := &http.Client{}

	for _, testCase := range testCases {
		req, err := http.NewRequestWithContext(ctx, testCase.method, server.URL, nil)
		require.NoError(t, err)

		resp, err := client.Do(req)
		require.NoError(t, err)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, testCase.expectedBody, string(body))
		require.Equal(t, testCase.expectedStatus, resp.StatusCode)

		_ = resp.Body.Close()
	}
}

func requestHandlerForTest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		_, _ = fmt.Fprintf(w, "GET response")
	case "POST":
		_, _ = fmt.Fprintf(w, "POST response")
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
