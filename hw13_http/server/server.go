package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	address := flag.String("ip", "localhost", "server ip")
	port := flag.Int("port", 8083, "server port")
	flag.Parse()

	http.HandleFunc("/", requestHandler)

	server := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", *address, *port),
		ReadHeaderTimeout: 1 * time.Second,
	}

	defer func() {
		_ = server.Close()
	}()

	log.Println("server started on", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("starting server error: %v\n", err)
	}
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if _, err := fmt.Fprintf(w, "GET response"); err != nil {
			return
		}
		fmt.Printf("received GET request for %s\n", r.Host+r.URL.Path)
	case "POST":
		body, postErr := io.ReadAll(r.Body)
		if postErr != nil {
			http.Error(w, "cannot read request body", http.StatusBadRequest)
			return
		}
		if _, err := fmt.Fprintf(w, "POST response"); err != nil {
			return
		}
		fmt.Printf("received POST request for %s with body: %s\n", r.Host+r.URL.Path, string(body))
	default:
		fmt.Printf("received %s request is invalid for %s\n", r.Method, r.Host+r.URL.Path)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
