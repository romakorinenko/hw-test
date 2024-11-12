package main

import (
	"bytes"
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
		body           string
		expectedStatus int
		expectedBody   string
	}{
		{
			method:         "GET",
			body:           "",
			expectedStatus: http.StatusOK,
			expectedBody:   "GET response",
		},
		{
			method:         "POST",
			body:           "request body",
			expectedStatus: http.StatusOK,
			expectedBody:   "POST response",
		},
		{
			method:         "PATCH",
			body:           "",
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "method not allowed\n",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(requestHandlerForTest))
	defer server.Close()

	ctx := context.Background()
	client := &http.Client{}

	for _, testCase := range testCases {
		var request *http.Request
		if testCase.body != "" {
			req, err := http.NewRequestWithContext(ctx, testCase.method, server.URL, bytes.NewBufferString(testCase.body))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "text/plain")
			request = req
		} else {
			req, err := http.NewRequestWithContext(context.Background(), testCase.method, server.URL, nil)
			require.NoError(t, err)
			request = req
		}

		resp, err := client.Do(request)
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
