package clickhouse

import (
	"context"
	"time"
)

// MockStore implements the Store interface with mock data
type MockStore struct{}

// NewMockStore creates a new mock ClickHouse store
func NewMockStore() *MockStore {
	return &MockStore{}
}

// GetTemplateCounts returns mock template counts
func (m *MockStore) GetTemplateCounts(ctx context.Context, org, dashboard, panelTitle, metricName string, startTime, endTime time.Time) (map[string]uint64, error) {
	// Return mock data with different templates and counts
	return map[string]uint64{
		"error_template_1":   150, // High count - anomaly
		"error_template_2":   80,
		"warning_template_1": 45,
		"info_template_1":    30,
		"debug_template_1":   20,
	}, nil
}

// GetRepresentativeLogs returns mock representative logs for templates
func (m *MockStore) GetRepresentativeLogs(ctx context.Context, org, dashboard, panelTitle, metricName string, templateIDs []string) (map[string][]string, error) {
	// Mock representative logs for each template
	mockLogs := map[string][]string{
		"error_template_1": {
			"ERROR: Out of memory on node-3 (allocated: 4.2GB, limit: 4GB)",
			"ERROR: Out of memory on node-5 (allocated: 4.5GB, limit: 4GB)",
			"ERROR: Out of memory on node-2 (allocated: 4.1GB, limit: 4GB)",
		},
		"error_template_2": {
			"ERROR: Connection timeout after 30s to database-01",
			"ERROR: Connection timeout after 30s to database-02",
		},
		"warning_template_1": {
			"WARNING: High CPU usage detected (95%) on worker-7",
			"WARNING: High CPU usage detected (96%) on worker-3",
			"WARNING: High CPU usage detected (94%) on worker-1",
		},
		"info_template_1": {
			"INFO: Service started successfully on port 8080",
			"INFO: Service started successfully on port 8081",
		},
		"debug_template_1": {
			"DEBUG: Health check passed in 15ms",
			"DEBUG: Health check passed in 12ms",
		},
	}

	// Filter to only requested template IDs
	result := make(map[string][]string)
	for _, templateID := range templateIDs {
		if logs, exists := mockLogs[templateID]; exists {
			result[templateID] = logs
		}
	}

	return result, nil
}

// VerifyTables always succeeds for mock store
func (m *MockStore) VerifyTables() error {
	return nil
}

// Close is a no-op for mock store
func (m *MockStore) Close() error {
	return nil
}

// Ensure MockStore implements Store interface
var _ Store = (*MockStore)(nil)
