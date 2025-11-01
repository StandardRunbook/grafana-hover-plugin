package api

import (
	"bytes"
	"container/list"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/StandardRunbook/grafana-hover-plugin/internal/analyzer"
	"github.com/StandardRunbook/grafana-hover-plugin/internal/clickhouse"
	"github.com/StandardRunbook/grafana-hover-plugin/internal/config"
)

// Helper function for tests to create mock log groups
func createTestLogGroups() []analyzer.LogGroup {
	return []analyzer.LogGroup{
		{
			RepresentativeLogs: []string{
				"ERROR: Out of memory on node-3",
				"ERROR: Out of memory on node-5",
			},
			RelativeChange: 2.5,
			KLContribution: 0.8,
			TemplateID:     "test_template_1",
		},
		{
			RepresentativeLogs: []string{
				"WARNING: High CPU usage detected",
			},
			RelativeChange: 1.2,
			KLContribution: 0.4,
			TemplateID:     "test_template_2",
		},
	}
}

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
			// Create a mock handler with mock ClickHouse
			mockStore := clickhouse.NewMockStore()
			mockAnalyzer := analyzer.NewLogAnalyzerWithStore(mockStore)
			mockHandler := &Handler{
				analyzer:     mockAnalyzer,
				cache:        make(map[string]*list.Element),
				cacheList:    list.New(),
				cacheTTL:     10 * time.Second,
				cacheMaxSize: 10,
				inFlight:     make(map[string]*inFlightRequest),
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


func TestNewHandlerWithoutClickHouse(t *testing.T) {
	// Test that NewHandler falls back to mock store when ClickHouse is unavailable
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

	if handler.analyzer == nil {
		t.Error("Expected analyzer to be created with mock store when ClickHouse is unavailable")
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

	// Should have log groups from mock store
	if len(resp.LogGroups) == 0 {
		t.Error("Expected log groups to be returned from mock store")
	}
}

func TestVerifyTablesWithoutClickHouse(t *testing.T) {
	// Create handler without ClickHouse (should use mock store)
	cfg := &config.Config{
		ClickHouse: config.ClickHouseConfig{
			URL:      "localhost:9999", // Invalid port
			User:     "default",
			Password: "",
			Database: "default",
		},
	}

	handler := NewHandler(cfg)

	// Mock store's VerifyTables always succeeds
	err := handler.VerifyTables()
	if err != nil {
		t.Errorf("Expected mock store VerifyTables to succeed, got: %v", err)
	}
}

func TestHandlerWithMockAnalyzer(t *testing.T) {
	// Test that handler properly handles mock ClickHouse
	mockStore := clickhouse.NewMockStore()
	mockAnalyzer := analyzer.NewLogAnalyzerWithStore(mockStore)
	handler := &Handler{
		analyzer:     mockAnalyzer,
		cache:        make(map[string]*list.Element),
		cacheList:    list.New(),
		cacheTTL:     10 * time.Second,
		cacheMaxSize: 10,
		inFlight:     make(map[string]*inFlightRequest),
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

func TestCacheKeyGeneration(t *testing.T) {
	cfg := &config.Config{
		ClickHouse: config.ClickHouseConfig{
			URL:      "localhost:9999",
			User:     "default",
			Password: "",
			Database: "default",
		},
	}
	handler := NewHandler(cfg)

	req1 := &QueryLogsRequest{
		Org:        "org1",
		Dashboard:  "dashboard1",
		PanelTitle: "panel1",
		MetricName: "metric1",
		StartTime:  time.Unix(1000, 0),
		EndTime:    time.Unix(2000, 0),
	}

	req2 := &QueryLogsRequest{
		Org:        "org1",
		Dashboard:  "dashboard1",
		PanelTitle: "panel1",
		MetricName: "metric1",
		StartTime:  time.Unix(1000, 0),
		EndTime:    time.Unix(2000, 0),
	}

	req3 := &QueryLogsRequest{
		Org:        "org1",
		Dashboard:  "dashboard1",
		PanelTitle: "panel1",
		MetricName: "metric2", // Different metric
		StartTime:  time.Unix(1000, 0),
		EndTime:    time.Unix(2000, 0),
	}

	key1 := handler.generateCacheKey(req1)
	key2 := handler.generateCacheKey(req2)
	key3 := handler.generateCacheKey(req3)

	// Same requests should generate same key
	if key1 != key2 {
		t.Error("Expected identical requests to generate same cache key")
	}

	// Different requests should generate different keys
	if key1 == key3 {
		t.Error("Expected different requests to generate different cache keys")
	}

	// Keys should be consistent across calls
	key1Again := handler.generateCacheKey(req1)
	if key1 != key1Again {
		t.Error("Expected cache key to be deterministic")
	}
}

func TestCacheBasicHitMiss(t *testing.T) {
	cfg := &config.Config{
		ClickHouse: config.ClickHouseConfig{
			URL:      "localhost:9999",
			User:     "default",
			Password: "",
			Database: "default",
		},
	}
	handler := NewHandler(cfg)

	key := "test-key"

	// Initially should be a miss
	_, _, found := handler.getCachedResultOrWait(key)
	if found {
		t.Error("Expected cache miss for new key")
	}

	// Complete an in-flight request to populate cache
	handler.startInFlightRequest(key)
	testData := createTestLogGroups()
	handler.completeInFlightRequest(key, testData, nil)

	// Now should be a hit
	result, err, found := handler.getCachedResultOrWait(key)
	if !found {
		t.Error("Expected cache hit after insertion")
	}
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(result) != len(testData) {
		t.Errorf("Expected %d log groups, got %d", len(testData), len(result))
	}
}

func TestCacheLRUEviction(t *testing.T) {
	cfg := &config.Config{
		ClickHouse: config.ClickHouseConfig{
			URL:      "localhost:9999",
			User:     "default",
			Password: "",
			Database: "default",
		},
	}
	handler := NewHandler(cfg)

	// Cache max is 10, so insert 11 items
	testData := createTestLogGroups()

	keys := make([]string, 11)
	for i := 0; i < 11; i++ {
		keys[i] = fmt.Sprintf("key-%d", i)
		handler.startInFlightRequest(keys[i])
		handler.completeInFlightRequest(keys[i], testData, nil)
	}

	// First key should be evicted (LRU)
	_, _, found := handler.getCachedResultOrWait(keys[0])
	if found {
		t.Error("Expected first key to be evicted after exceeding cache size")
	}

	// Last key should still be in cache
	_, _, found = handler.getCachedResultOrWait(keys[10])
	if !found {
		t.Error("Expected last key to still be in cache")
	}

	// Cache size should be at max
	handler.cacheMu.Lock()
	size := handler.cacheList.Len()
	handler.cacheMu.Unlock()

	if size != handler.cacheMaxSize {
		t.Errorf("Expected cache size to be %d, got %d", handler.cacheMaxSize, size)
	}
}

func TestCacheLRUAccessOrder(t *testing.T) {
	cfg := &config.Config{
		ClickHouse: config.ClickHouseConfig{
			URL:      "localhost:9999",
			User:     "default",
			Password: "",
			Database: "default",
		},
	}
	handler := NewHandler(cfg)

	// Insert max entries
	testData := createTestLogGroups()
	keys := make([]string, 10)
	for i := 0; i < 10; i++ {
		keys[i] = fmt.Sprintf("key-%d", i)
		handler.startInFlightRequest(keys[i])
		handler.completeInFlightRequest(keys[i], testData, nil)
	}

	// Access the first key to move it to front
	handler.getCachedResultOrWait(keys[0])

	// Insert one more item, should evict key-1 (not key-0)
	newKey := "key-new"
	handler.startInFlightRequest(newKey)
	handler.completeInFlightRequest(newKey, testData, nil)

	// key-0 should still be in cache (recently accessed)
	_, _, found := handler.getCachedResultOrWait(keys[0])
	if !found {
		t.Error("Expected recently accessed key-0 to still be in cache")
	}

	// key-1 should be evicted (LRU)
	_, _, found = handler.getCachedResultOrWait(keys[1])
	if found {
		t.Error("Expected key-1 to be evicted as LRU")
	}
}

func TestRequestCoalescing(t *testing.T) {
	cfg := &config.Config{
		ClickHouse: config.ClickHouseConfig{
			URL:      "localhost:9999",
			User:     "default",
			Password: "",
			Database: "default",
		},
	}
	handler := NewHandler(cfg)

	key := "test-key"
	testData := createTestLogGroups()

	// Start in-flight request
	handler.startInFlightRequest(key)

	var wg sync.WaitGroup
	results := make([][]byte, 5)
	errors := make([]error, 5)

	// Launch 5 concurrent requests waiting on the same key
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			result, err, found := handler.getCachedResultOrWait(key)
			if !found {
				t.Errorf("Goroutine %d: expected to find result from in-flight request", idx)
			}
			errors[idx] = err
			if err == nil && len(result) > 0 {
				results[idx] = []byte(result[0].RepresentativeLogs[0])
			}
		}(i)
	}

	// Give goroutines time to start waiting
	time.Sleep(50 * time.Millisecond)

	// Complete the request - should broadcast to all waiters
	handler.completeInFlightRequest(key, testData, nil)

	// Wait for all goroutines to complete
	wg.Wait()

	// All goroutines should have received the result
	for i, result := range results {
		if len(result) == 0 {
			t.Errorf("Goroutine %d did not receive result", i)
		}
		if errors[i] != nil {
			t.Errorf("Goroutine %d received error: %v", i, errors[i])
		}
	}
}

func TestRequestCoalescingWithError(t *testing.T) {
	cfg := &config.Config{
		ClickHouse: config.ClickHouseConfig{
			URL:      "localhost:9999",
			User:     "default",
			Password: "",
			Database: "default",
		},
	}
	handler := NewHandler(cfg)

	key := "test-key"
	testError := errors.New("test database error")

	// Start in-flight request
	handler.startInFlightRequest(key)

	var wg sync.WaitGroup
	errors := make([]error, 3)

	// Launch 3 concurrent requests waiting on the same key
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			_, err, found := handler.getCachedResultOrWait(key)
			if !found {
				t.Errorf("Goroutine %d: expected to find result from in-flight request", idx)
			}
			errors[idx] = err
		}(i)
	}

	// Give goroutines time to start waiting
	time.Sleep(50 * time.Millisecond)

	// Complete with error - should broadcast error to all waiters
	handler.completeInFlightRequest(key, nil, testError)

	// Wait for all goroutines to complete
	wg.Wait()

	// All goroutines should have received the error
	for i, err := range errors {
		if err == nil {
			t.Errorf("Goroutine %d did not receive error", i)
		}
	}
}

func TestCacheExpiration(t *testing.T) {
	cfg := &config.Config{
		ClickHouse: config.ClickHouseConfig{
			URL:      "localhost:9999",
			User:     "default",
			Password: "",
			Database: "default",
		},
	}
	handler := NewHandler(cfg)

	// Set a very short TTL for testing
	handler.cacheTTL = 100 * time.Millisecond

	key := "test-key"
	testData := createTestLogGroups()

	// Add to cache
	handler.startInFlightRequest(key)
	handler.completeInFlightRequest(key, testData, nil)

	// Should be cached immediately
	_, _, found := handler.getCachedResultOrWait(key)
	if !found {
		t.Error("Expected cache hit immediately after insertion")
	}

	// Wait for expiration
	time.Sleep(150 * time.Millisecond)

	// Should be expired now
	_, _, found = handler.getCachedResultOrWait(key)
	if found {
		t.Error("Expected cache miss after TTL expiration")
	}

	// Cache should have cleaned up the entry
	handler.cacheMu.Lock()
	size := handler.cacheList.Len()
	handler.cacheMu.Unlock()

	if size != 0 {
		t.Errorf("Expected cache to be empty after expiration check, got size %d", size)
	}
}

func TestCacheExpirationCleanup(t *testing.T) {
	cfg := &config.Config{
		ClickHouse: config.ClickHouseConfig{
			URL:      "localhost:9999",
			User:     "default",
			Password: "",
			Database: "default",
		},
	}
	handler := NewHandler(cfg)

	// Set a very short TTL for testing
	handler.cacheTTL = 50 * time.Millisecond

	// Add multiple entries
	testData := createTestLogGroups()
	for i := 0; i < 5; i++ {
		key := fmt.Sprintf("key-%d", i)
		handler.startInFlightRequest(key)
		handler.completeInFlightRequest(key, testData, nil)
	}

	// All should be cached
	handler.cacheMu.Lock()
	initialSize := handler.cacheList.Len()
	handler.cacheMu.Unlock()

	if initialSize != 5 {
		t.Errorf("Expected 5 entries in cache, got %d", initialSize)
	}

	// Wait for expiration
	time.Sleep(100 * time.Millisecond)

	// Manually trigger cleanup (simulating the background goroutine)
	handler.cacheMu.Lock()
	now := time.Now()
	expired := 0
	var next *interface{}
	for elem := handler.cacheList.Front(); elem != nil; {
		entry := elem.Value.(*cacheEntry)
		nextElem := elem.Next()
		if now.After(entry.expiresAt) {
			handler.cacheList.Remove(elem)
			delete(handler.cache, entry.key)
			expired++
		}
		elem = nextElem
		next = nil // Suppress unused warning
		_ = next
	}
	handler.cacheMu.Unlock()

	if expired != 5 {
		t.Errorf("Expected 5 expired entries to be cleaned up, got %d", expired)
	}

	// Cache should be empty
	handler.cacheMu.Lock()
	finalSize := handler.cacheList.Len()
	handler.cacheMu.Unlock()

	if finalSize != 0 {
		t.Errorf("Expected cache to be empty after cleanup, got size %d", finalSize)
	}
}
