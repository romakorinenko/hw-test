package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func main() {
	host := flag.String("host", "http://localhost:8083", "server baseURL")
	path := flag.String("path", "/", "baseURL path")
	flag.Parse()

	baseURL, err := url.Parse(*host)
	if err != nil {
		return
	}
	baseURL.Path = *path

	sendRequest(baseURL.String(), "GET", "")
	sendRequest(baseURL.String(), "POST", "request body")
}

func sendRequest(url, method, requestBody string) {
	var request *http.Request
	if requestBody != "" {
		req, err := http.NewRequestWithContext(context.Background(), method, url, bytes.NewBufferString(requestBody))
		if err != nil {
			fmt.Printf("%s request creating error: %v\n", method, err)
			return
		}
		req.Header.Set("Content-Type", "text/plain")
		request = req
	} else {
		req, err := http.NewRequestWithContext(context.Background(), method, url, nil)
		if err != nil {
			fmt.Printf("%s request creating error: %v\n", method, err)
			return
		}
		request = req
	}

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Printf("%s request making error: %v\n", method, err)
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("response body reading error: %v\n", err)
		return
	}

	fmt.Printf("Response Status Code: %d\n", resp.StatusCode)
	fmt.Printf("Response Body: %s\n", string(body))
}
