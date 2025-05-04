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

	intervalStatus := 3 * time.Second
	intervalFlaky := 10 * time.Second
	timeout := 5 * time.Second

	// create HTTP client with timeout
	client := &http.Client{
		Timeout: timeout,
	}

	flaky_endpoint := fmt.Sprintf("http://%s:%s/flaky", os.Getenv("HOST"), os.Getenv("PORT"))

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("Press Ctrl+C to exit")

	// setup tickers
	tickerStatus := time.NewTicker(intervalStatus)
	defer tickerStatus.Stop()
	tickerFlaky := time.NewTicker(intervalFlaky)
	defer tickerFlaky.Stop()

	status_endpoint := fmt.Sprintf("http://%s:%s/status", os.Getenv("HOST"), os.Getenv("PORT"))

	// main loop
	for {
		select {
		case <-tickerStatus.C:
			getStatus(client, status_endpoint)
		case <-tickerFlaky.C:
			postFlaky(client, flaky_endpoint)
		case sig := <-sigChan:
			fmt.Printf("\nReceived signal %v, shutting down...\n", sig)
			return
		}
	}
}

func getStatus(client *http.Client, endpoint string) {
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

func postFlaky(client *http.Client, endpoint string) {
	fmt.Printf("switch flaky on/off\n")

	resp, err := client.Post(endpoint, "application/json", nil)
	if err != nil {
		fmt.Printf("Request Error: %v", err)
		return
	}
	defer resp.Body.Close()

	// check post response
	if resp.StatusCode != http.StatusAccepted {
		fmt.Printf("StatusCode: %d", resp.StatusCode)
		return
	}
	body, _ := io.ReadAll(resp.Body)

	fmt.Println(string(body))
}
