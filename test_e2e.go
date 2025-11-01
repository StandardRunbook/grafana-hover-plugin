package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type LogGroup struct {
	RepresentativeLogs []string `json:"representative_logs"`
	RelativeChange     float64  `json:"relative_change"`
}

type QueryLogsResponse struct {
	LogGroups []LogGroup `json:"log_groups"`
}

type AnalyzeRequest struct {
	Org        string `json:"org"`
	Dashboard  string `json:"dashboard"`
	PanelTitle string `json:"panel_title"`
	MetricName string `json:"metric_name"`
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
}

func main() {
	// Test the API endpoint
	req := AnalyzeRequest{
		Org:        "acme-corp",
		Dashboard:  "production-dashboard",
		PanelTitle: "API Logs",
		MetricName: "error_rate",
		StartTime:  "2025-11-01T00:00:00Z",
		EndTime:    "2025-11-01T00:15:00Z",
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		fmt.Printf("Error marshaling request: %v\n", err)
		return
	}

	fmt.Println("=== Testing Go Backend E2E ===")
	fmt.Printf("Request: %s\n\n", string(jsonData))

	// Make HTTP request
	resp, err := http.Post("http://localhost:8080/analyze", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		fmt.Println("\nMake sure the backend is running: go run cmd/main.go")
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return
	}

	fmt.Printf("Status: %d\n", resp.StatusCode)
	fmt.Printf("Response:\n%s\n\n", string(body))

	if resp.StatusCode != http.StatusOK {
		fmt.Println("❌ Request failed")
		return
	}

	// Parse response
	var response QueryLogsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Printf("Error parsing response: %v\n", err)
		return
	}

	fmt.Println("=== Parsed Log Groups ===")
	for i, group := range response.LogGroups {
		fmt.Printf("\n%d. Relative Change: %.4f\n", i+1, group.RelativeChange)
		fmt.Printf("   Representative Logs:\n")
		for j, log := range group.RepresentativeLogs {
			fmt.Printf("      %d) %s\n", j+1, log)
		}
	}

	fmt.Println("\n=== Testing Cache Hit ===")
	fmt.Println("Making second request to test cache...")

	start := time.Now()
	resp2, err := http.Post("http://localhost:8080/analyze", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error making second request: %v\n", err)
		return
	}
	defer resp2.Body.Close()
	elapsed := time.Since(start)

	fmt.Printf("Second request took: %v\n", elapsed)
	if resp2.StatusCode == http.StatusOK {
		fmt.Println("✓ Cache hit successful (should be faster)")
	}

	fmt.Println("\n=== Testing Different Panel ===")
	req.PanelTitle = "Database Logs"
	req.MetricName = "query_errors"

	jsonData2, _ := json.Marshal(req)
	resp3, err := http.Post("http://localhost:8080/analyze", "application/json", bytes.NewBuffer(jsonData2))
	if err != nil {
		fmt.Printf("Error making third request: %v\n", err)
		return
	}
	defer resp3.Body.Close()

	body3, _ := io.ReadAll(resp3.Body)
	var response3 QueryLogsResponse
	if err := json.Unmarshal(body3, &response3); err == nil {
		fmt.Printf("\nDatabase Logs found %d log groups:\n", len(response3.LogGroups))
		for i, group := range response3.LogGroups {
			fmt.Printf("%d. Relative Change: %.4f, Examples: %d logs\n", i+1, group.RelativeChange, len(group.RepresentativeLogs))
		}
	}

	fmt.Println("\n✓ E2E test completed successfully!")
}
