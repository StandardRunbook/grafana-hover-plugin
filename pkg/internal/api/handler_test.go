package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"grafana-plugin-api/internal/config"
)

func TestQueryLogsValidation(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		checkError     bool
	}{
		{
			name: "valid request",
			requestBody: QueryLogsRequest{
				Org:        "test-org",
				Dashboard:  "test-dashboard",
				PanelTitle: "test-panel",
				MetricName: "test-metric",
				StartTime:  time.Now().Add(-1 * time.Hour),
				EndTime:    time.Now(),
			},
			expectedStatus: http.StatusOK,
			checkError:     false,
		},
		{
			name: "invalid time range - start after end",
			requestBody: QueryLogsRequest{
				Org:        "test-org",
				Dashboard:  "test-dashboard",
				PanelTitle: "test-panel",
				MetricName: "test-metric",
				StartTime:  time.Now(),
				EndTime:    time.Now().Add(-1 * time.Hour),
			},
			expectedStatus: http.StatusBadRequest,
			checkError:     true,
		},
		{
			name: "missing required field",
			requestBody: map[string]interface{}{
				"org":        "test-org",
				"dashboard":  "test-dashboard",
				"start_time": time.Now().Add(-1 * time.Hour),
				"end_time":   time.Now(),
				// Missing panel_title and metric_name
			},
			expectedStatus: http.StatusBadRequest,
			checkError:     true,
		},
		{
			name:           "invalid JSON",
			requestBody:    "invalid json",
			expectedStatus: http.StatusBadRequest,
			checkError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock handler
			mockHandler := &Handler{
				analyzer:      nil,
				analyzerError: errors.New("mock error"),
			}

			// Marshal request body
			var bodyBytes []byte
			var err error

			switch v := tt.requestBody.(type) {
			case string:
				bodyBytes = []byte(v)
			default:
				bodyBytes, err = json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatalf("Failed to marshal request body: %v", err)
				}
			}

			// Create request
			req := httptest.NewRequest(http.MethodPost, "/query_logs", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			mockHandler.QueryLogs(w, req)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Check for error in response
			if tt.checkError {
				var errResp ErrorResponse
				if err := json.Unmarshal(w.Body.Bytes(), &errResp); err != nil {
					t.Errorf("Failed to unmarshal error response: %v", err)
				}
				if errResp.Error == "" {
					t.Error("Expected error field in response")
				}
			}
		})
	}
}

func TestErrorResponseFormat(t *testing.T) {
	tests := []struct {
		name     string
		response ErrorResponse
	}{
		{
			name: "with code",
			response: ErrorResponse{
				Error:   "TestError",
				Message: "This is a test error",
				Code:    intPtr(400),
			},
		},
		{
			name: "without code",
			response: ErrorResponse{
				Error:   "TestError",
				Message: "This is a test error",
				Code:    nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal to JSON
			data, err := json.Marshal(tt.response)
			if err != nil {
				t.Fatalf("Failed to marshal error response: %v", err)
			}

			// Unmarshal back
			var decoded ErrorResponse
			if err := json.Unmarshal(data, &decoded); err != nil {
				t.Fatalf("Failed to unmarshal error response: %v", err)
			}

			// Check fields
			if decoded.Error != tt.response.Error {
				t.Errorf("Expected error '%s', got '%s'", tt.response.Error, decoded.Error)
			}

			if decoded.Message != tt.response.Message {
				t.Errorf("Expected message '%s', got '%s'", tt.response.Message, decoded.Message)
			}

			// Check code (handling nil)
			if tt.response.Code == nil && decoded.Code != nil {
				t.Error("Expected code to be nil")
			}
			if tt.response.Code != nil {
				if decoded.Code == nil {
					t.Error("Expected code to be non-nil")
				} else if *decoded.Code != *tt.response.Code {
					t.Errorf("Expected code %d, got %d", *tt.response.Code, *decoded.Code)
				}
			}
		})
	}
}

func TestLogGroupResponseFormat(t *testing.T) {
	response := QueryLogsResponse{
		LogGroups: []LogGroup{
			{
				RepresentativeLogs: []string{"Log 1", "Log 2"},
				RelativeChange:     0.5,
			},
			{
				RepresentativeLogs: []string{"Error log"},
				RelativeChange:     1.2,
			},
		},
	}

	// Marshal to JSON
	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}

	// Unmarshal back
	var decoded QueryLogsResponse
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Check log groups count
	if len(decoded.LogGroups) != 2 {
		t.Errorf("Expected 2 log groups, got %d", len(decoded.LogGroups))
	}

	// Check first log group
	if len(decoded.LogGroups[0].RepresentativeLogs) != 2 {
		t.Errorf("Expected 2 representative logs in first group, got %d", len(decoded.LogGroups[0].RepresentativeLogs))
	}

	if decoded.LogGroups[0].RelativeChange != 0.5 {
		t.Errorf("Expected relative change 0.5, got %f", decoded.LogGroups[0].RelativeChange)
	}

	// Check second log group
	if len(decoded.LogGroups[1].RepresentativeLogs) != 1 {
		t.Errorf("Expected 1 representative log in second group, got %d", len(decoded.LogGroups[1].RepresentativeLogs))
	}
}

func TestContainsStr(t *testing.T) {
	tests := []struct {
		s      string
		substr string
		expect bool
	}{
		{"Connection refused", "refused", true},
		{"Connection refused", "timeout", false},
		{"UNKNOWN_TABLE error", "UNKNOWN_TABLE", true},
		{"", "test", false},
		{"test", "", true},
		{"", "", true},
	}

	for _, tt := range tests {
		result := containsStr(tt.s, tt.substr)
		if result != tt.expect {
			t.Errorf("containsStr(%q, %q) = %v, expected %v", tt.s, tt.substr, result, tt.expect)
		}
	}
}

func TestNewHandlerWithoutClickHouse(t *testing.T) {
	// Test that NewHandler doesn't panic when ClickHouse is unavailable
	cfg := &config.Config{
		ClickHouse: config.ClickHouseConfig{
			URL:      "localhost:9999", // Invalid port
			User:     "default",
			Password: "",
			Database: "default",
		},
	}

	handler := NewHandler(cfg)
	if handler == nil {
		t.Fatal("Expected handler to be created even without ClickHouse")
	}

	if handler.analyzer != nil {
		t.Error("Expected analyzer to be nil when ClickHouse is unavailable")
	}

	if handler.analyzerError == nil {
		t.Error("Expected analyzerError to be set when ClickHouse is unavailable")
	}
}

func TestQueryLogsWithoutClickHouse(t *testing.T) {
	// Create handler without ClickHouse
	cfg := &config.Config{
		ClickHouse: config.ClickHouseConfig{
			URL:      "localhost:9999", // Invalid port
			User:     "default",
			Password: "",
			Database: "default",
		},
	}

	handler := NewHandler(cfg)

	// Create valid request
	reqBody := QueryLogsRequest{
		Org:        "test-org",
		Dashboard:  "test-dashboard",
		PanelTitle: "test-panel",
		MetricName: "test-metric",
		StartTime:  time.Now().Add(-1 * time.Hour),
		EndTime:    time.Now(),
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/query_logs", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.QueryLogs(w, req)

	// Should return 200 with mock data
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp QueryLogsResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Should have mock log groups
	if len(resp.LogGroups) == 0 {
		t.Error("Expected mock log groups to be returned")
	}

	// Check that mock data is clearly labeled
	foundMockLabel := false
	for _, group := range resp.LogGroups {
		for _, log := range group.RepresentativeLogs {
			if containsStr(log, "MOCK DATA") || containsStr(log, "Example") {
				foundMockLabel = true
				break
			}
		}
	}

	if !foundMockLabel {
		t.Error("Expected mock data to be clearly labeled")
	}
}

func TestVerifyTablesWithoutClickHouse(t *testing.T) {
	// Create handler without ClickHouse
	cfg := &config.Config{
		ClickHouse: config.ClickHouseConfig{
			URL:      "localhost:9999", // Invalid port
			User:     "default",
			Password: "",
			Database: "default",
		},
	}

	handler := NewHandler(cfg)

	err := handler.VerifyTables()
	if err == nil {
		t.Error("Expected error when verifying tables without ClickHouse")
	}
}

func TestHandlerWithMockAnalyzer(t *testing.T) {
	// Test that handler properly handles nil analyzer
	handler := &Handler{
		analyzer:      nil,
		analyzerError: errors.New("mock connection error"),
	}

	reqBody := QueryLogsRequest{
		Org:        "test-org",
		Dashboard:  "test-dashboard",
		PanelTitle: "test-panel",
		MetricName: "test-metric",
		StartTime:  time.Now().Add(-1 * time.Hour),
		EndTime:    time.Now(),
	}

	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/query_logs", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.QueryLogs(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp QueryLogsResponse
	json.Unmarshal(w.Body.Bytes(), &resp)

	if len(resp.LogGroups) == 0 {
		t.Error("Expected mock data to be returned")
	}
}
