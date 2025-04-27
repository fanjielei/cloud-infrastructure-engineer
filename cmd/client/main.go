package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	endpoint := "http://localhost:8080/status"

	interval := 5 * time.Second
	timeout := 5 * time.Second

	// create HTTP client with timeout
	client := &http.Client{
		Timeout: timeout,
	}

	fmt.Printf("set flasy to be true\n\n")
	client.Post("http://localhost:8080/flaky", "json", nil)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Printf("Starting status monitoring for endpoint: %s\nChecking every %s with %s timeout\n", endpoint, interval, timeout)
	fmt.Println("Press Ctrl+C to exit")

	// setup a ticker
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	checkStatus(client, endpoint)

	// main loop
	for {
		select {
		case <-ticker.C:
			checkStatus(client, endpoint)
		case sig := <-sigChan:
			fmt.Printf("\nReceived signal %v, shutting down...\n", sig)
			return
		}
	}
}

func checkStatus(client *http.Client, endpoint string) {
	fmt.Printf("\n[%s] Checking endpoint: %s\n", time.Now().Format(time.RFC3339), endpoint)

	// measure request duration
	startTime := time.Now()

	// make the request
	resp, err := client.Get(endpoint)
	duration := time.Since(startTime)

	// Report duration regardless of success
	fmt.Printf("  Duration: %v\n", duration)

	// Handle request errors
	if err != nil {
		fmt.Printf("  Success: false\n")
		fmt.Printf("  Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("  Success: false\n")
		fmt.Printf("  Status Code: %d\n", resp.StatusCode)
		fmt.Printf("  Error reading body: %v\n", err)
		return
	}

	// Report status details
	success := resp.StatusCode >= 200 && resp.StatusCode < 300
	fmt.Printf("  Success: %t\n", success)
	fmt.Printf("  Status Code: %d\n", resp.StatusCode)
	fmt.Printf("  Response: %s\n", body)
}
