package main

import (
	"context"
	"fmt"
	"time"

	"github.com/StandardRunbook/grafana-hover-plugin/internal/clickhouse"
	"github.com/StandardRunbook/grafana-hover-plugin/internal/config"
)

func main() {
	// Try different connection configs
	configs := []*config.ClickHouseConfig{
		{URL: "localhost:9000", Database: "default", User: "default", Password: ""},
		{URL: "localhost:9000", Database: "default", User: "", Password: ""},
		{URL: "127.0.0.1:9000", Database: "default", User: "default", Password: ""},
	}

	var client *clickhouse.Client
	var err error
	var cfg *config.ClickHouseConfig

	for i, testCfg := range configs {
		fmt.Printf("Attempt %d: %s@%s\n", i+1, testCfg.User, testCfg.URL)
		client, err = clickhouse.NewClient(testCfg)
		if err == nil {
			cfg = testCfg
			break
		}
		fmt.Printf("  ❌ %v\n", err)
	}

	if err != nil {
		fmt.Println("\n❌ All connection attempts failed")
		return
	}

	defer client.Close()
	fmt.Printf("\n✓ Connected to ClickHouse with config: %s@%s\n", cfg.User, cfg.URL)

	// Test queries
	fmt.Println("\n=== Testing GetTemplateCounts ===")
	counts, err := client.GetTemplateCounts(
		context.Background(),
		"acme-corp",
		"production-dashboard",
		"API Logs",
		"error_rate",
		time.Date(2025, 11, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 11, 1, 0, 15, 0, 0, time.UTC),
	)
	if err != nil {
		fmt.Printf("❌ Failed to get template counts: %v\n", err)
		return
	}

	fmt.Printf("Found %d templates:\n", len(counts))
	for templateID, count := range counts {
		fmt.Printf("  - %s: %d occurrences\n", templateID, count)
	}

	// Get representative logs
	templateIDs := make([]string, 0, len(counts))
	for id := range counts {
		templateIDs = append(templateIDs, id)
	}

	fmt.Println("\n=== Testing GetRepresentativeLogs ===")
	representatives, err := client.GetRepresentativeLogs(
		context.Background(),
		"acme-corp",
		"production-dashboard",
		"API Logs",
		"error_rate",
		templateIDs,
	)
	if err != nil {
		fmt.Printf("❌ Failed to get representative logs: %v\n", err)
		return
	}

	fmt.Printf("Found representatives for %d templates:\n", len(representatives))
	for templateID, logs := range representatives {
		fmt.Printf("  - %s (%d examples):\n", templateID, len(logs))
		for i, log := range logs {
			if i < 2 { // Show first 2 examples
				fmt.Printf("      %d) %s\n", i+1, log)
			}
		}
	}

	fmt.Println("\n✓ All tests passed!")
}
