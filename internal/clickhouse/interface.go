package clickhouse

import (
	"context"
	"time"
)

// Store defines the interface for ClickHouse operations
type Store interface {
	GetTemplateCounts(ctx context.Context, org, dashboard, panelTitle, metricName string, startTime, endTime time.Time) (map[string]uint64, error)
	// GetRepresentativeLogs returns actual log messages from the `logs` table
	// for the given template IDs, restricted to the (startTime, endTime) window.
	// This enforces the hover-content invariant: rows shown to the user must
	// come from the exact time slice they're inspecting, not a sample across
	// time. Previously this read from a separate `template_examples` table;
	// that table was removed because it returned representative-rather-than-
	// exact logs.
	GetRepresentativeLogs(ctx context.Context, org, dashboard, panelTitle, metricName string, templateIDs []string, startTime, endTime time.Time) (map[string][]string, error)
	VerifyTables() error
	Close() error
}

// Ensure Client implements Store interface
var _ Store = (*Client)(nil)
