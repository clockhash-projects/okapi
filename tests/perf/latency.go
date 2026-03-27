package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

const (
	baseURL = "http://localhost:8080"
	timeout = 5 * time.Second
)

var endpoints = []string{
	"/health/all",
	"/services",
	"/incidents",
	"/history",
	"/maintenance",
}

type result struct {
	endpoint string
	latency  string
	err      error
}

func main() {
	fmt.Printf("Starting Go Speed Test against %s...\n\n", baseURL)
	fmt.Printf("%-20s | %-10s\n", "Endpoint", "Latency")
	fmt.Println("-------------------------------------")

	var wg sync.WaitGroup
	results := make(chan result, len(endpoints))
	client := &http.Client{
		Timeout: timeout,
	}

	for _, ep := range endpoints {
		wg.Add(1)
		go func(endpoint string) {
			defer wg.Done()
			start := time.Now()
			resp, err := client.Get(baseURL + endpoint)
			duration := time.Since(start)

			res := result{endpoint: endpoint}
			if err != nil {
				res.err = err
			} else {
				defer resp.Body.Close()
				if resp.StatusCode >= 400 {
					res.err = fmt.Errorf("HTTP %d", resp.StatusCode)
				} else {
					res.latency = fmt.Sprintf("%.4fs", duration.Seconds())
				}
			}
			results <- res
		}(ep)
	}

	// Close channel once all probes finish
	go func() {
		wg.Wait()
		close(results)
	}()

	for res := range results {
		if res.err != nil {
			fmt.Printf("%-20s | Error: %v\n", res.endpoint, res.err)
		} else {
			fmt.Printf("%-20s | %-10s\n", res.endpoint, res.latency)
		}
	}
}
