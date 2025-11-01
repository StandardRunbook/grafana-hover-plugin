package clickhouse

import (
	"context"
	"time"
)

// Store defines the interface for ClickHouse operations
type Store interface {
	GetTemplateCounts(ctx context.Context, org, dashboard, panelTitle, metricName string, startTime, endTime time.Time) (map[string]uint64, error)
	GetRepresentativeLogs(ctx context.Context, org, dashboard, panelTitle, metricName string, templateIDs []string) (map[string][]string, error)
	VerifyTables() error
	Close() error
}

// Ensure Client implements Store interface
var _ Store = (*Client)(nil)
