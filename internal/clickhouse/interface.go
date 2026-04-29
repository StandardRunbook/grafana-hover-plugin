package clickhouse

import (
	"context"
	"time"
)

// Store defines the interface for ClickHouse operations
type Store interface {
	GetTemplateCounts(ctx context.Context, org, dashboard, panelTitle, metricName string, startTime, endTime time.Time) (map[string]uint64, error)
	// GetTemplateCountsByPrefix is GetTemplateCounts narrowed to log
	// messages whose lowercase form contains `prefix`. Used by the
	// empty-baseline relevance bias: when there's no baseline window to
	// diverge from, we still want hovering the CPU panel to surface
	// cpu_usage logs (not whatever template happens to dominate the
	// stream). Substring rather than literal-prefix because metric names
	// like `error_rate` map to log prefixes like `ERROR` — case-folded
	// substring matches both without per-metric mappings.
	GetTemplateCountsByPrefix(ctx context.Context, org, dashboard, panelTitle, metricName, prefix string, startTime, endTime time.Time) (map[string]uint64, error)
	// GetRepresentativeLogs returns actual log messages from the `logs` table
	// for the given template IDs, restricted to the (startTime, endTime) window.
	// This enforces the hover-content invariant: rows shown to the user must
	// come from the exact time slice they're inspecting, not a sample across
	// time. Previously this read from a separate `template_examples` table;
	// that table was removed because it returned representative-rather-than-
	// exact logs.
	GetRepresentativeLogs(ctx context.Context, org, dashboard, panelTitle, metricName string, templateIDs []string, startTime, endTime time.Time) (map[string][]string, error)
	// EnsureMapping registers a (dashboard, panel, metric) → log-stream
	// mapping if one isn't already present, so the analyzer's join
	// against metric_log_hover_mv can resolve. Without it, every new
	// panel needs a manual SQL insert. Picks the org's busiest log
	// stream as the default — fine for single-stream demo orgs.
	EnsureMapping(ctx context.Context, org, dashboard, panelTitle, metricName string) error
	VerifyTables() error
	Close() error
}

// Ensure Client implements Store interface
var _ Store = (*Client)(nil)
