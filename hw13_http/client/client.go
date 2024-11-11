package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func main() {
	host := flag.String("host", "http://localhost:8083", "server baseUrl")
	path := flag.String("path", "/", "baseUrl path")
	flag.Parse()

	baseUrl, err := url.Parse(*host)
	if err != nil {
		return
	}
	baseUrl.Path = *path

	sendRequest(baseUrl.String(), "GET")
	sendRequest(baseUrl.String(), "POST")
}

func sendRequest(url, method string) {
	req, err := http.NewRequestWithContext(context.Background(), method, url, nil)
	if err != nil {
		fmt.Printf("%s request creating error: %v\n", method, err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
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
