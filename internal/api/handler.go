package api

import (
	"container/list"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/StandardRunbook/grafana-hover-plugin/internal/analyzer"
	"github.com/StandardRunbook/grafana-hover-plugin/internal/clickhouse"
	"github.com/StandardRunbook/grafana-hover-plugin/internal/config"
)

type cacheEntry struct {
	key       string
	logGroups []analyzer.LogGroup
	expiresAt time.Time
}

type inFlightRequest struct {
	done       chan struct{}
	result     []analyzer.LogGroup
	err        error
	resultOnce sync.Once
}

type Handler struct {
	analyzer     *analyzer.LogAnalyzer
	cache        map[string]*list.Element // map key to list element
	cacheList    *list.List                // doubly-linked list for LRU order
	cacheMu      sync.Mutex
	cacheTTL     time.Duration
	cacheMaxSize int
	inFlight     map[string]*inFlightRequest
	inFlightMu   sync.Mutex
}

type QueryLogsRequest struct {
	Org        string    `json:"org"`
	Dashboard  string    `json:"dashboard"`
	PanelTitle string    `json:"panel_title"`
	MetricName string    `json:"metric_name"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
}

type LogGroup struct {
	RepresentativeLogs []string `json:"representative_logs"`
	RelativeChange     float64  `json:"relative_change"`
}

type QueryLogsResponse struct {
	LogGroups []LogGroup `json:"log_groups"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    *int   `json:"code,omitempty"`
}

func NewHandler(cfg *config.Config) *Handler {
	var logAnalyzer *analyzer.LogAnalyzer

	// Try to connect to real ClickHouse
	realAnalyzer, err := analyzer.NewLogAnalyzer(&cfg.ClickHouse)
	if err != nil {
		log.Printf("Warning: Failed to connect to ClickHouse: %v", err)
		log.Printf("Using mock ClickHouse with sample data")

		// Use mock ClickHouse store
		mockStore := clickhouse.NewMockStore()
		logAnalyzer = analyzer.NewLogAnalyzerWithStore(mockStore)
	} else {
		logAnalyzer = realAnalyzer
	}

	h := &Handler{
		analyzer:     logAnalyzer,
		cache:        make(map[string]*list.Element),
		cacheList:    list.New(),
		cacheTTL:     10 * time.Second,
		cacheMaxSize: 10,
		inFlight:     make(map[string]*inFlightRequest),
	}

	// Start background cleanup goroutine
	go h.cleanupExpiredCache()

	return h
}

func (h *Handler) VerifyTables() error {
	return h.analyzer.VerifyTables()
}

// generateCacheKey creates a unique cache key from request parameters
func (h *Handler) generateCacheKey(req *QueryLogsRequest) string {
	// Create a deterministic key from all request parameters
	key := fmt.Sprintf("%s|%s|%s|%s|%d|%d",
		req.Org,
		req.Dashboard,
		req.PanelTitle,
		req.MetricName,
		req.StartTime.Unix(),
		req.EndTime.Unix(),
	)

	// Hash the key to keep it compact
	hash := sha256.Sum256([]byte(key))
	return fmt.Sprintf("%x", hash)
}

// getCachedResultOrWait attempts to retrieve a cached result or waits for an in-flight request
func (h *Handler) getCachedResultOrWait(key string) ([]analyzer.LogGroup, error, bool) {
	// First check cache
	h.cacheMu.Lock()
	elem, exists := h.cache[key]
	if exists {
		entry := elem.Value.(*cacheEntry)
		if time.Now().Before(entry.expiresAt) {
			// Move to front (most recently used) - O(1)
			h.cacheList.MoveToFront(elem)
			logGroups := entry.logGroups
			h.cacheMu.Unlock()
			log.Printf("Cache HIT for key: %s", truncateKey(key))
			return logGroups, nil, true
		}
		// Expired, remove it
		h.cacheList.Remove(elem)
		delete(h.cache, key)
	}
	h.cacheMu.Unlock()

	// Check if there's an in-flight request
	h.inFlightMu.Lock()
	inflight, inProgress := h.inFlight[key]
	h.inFlightMu.Unlock()

	if inProgress {
		log.Printf("Waiting for in-flight request: %s", truncateKey(key))
		// Wait for the result from the in-flight request
		<-inflight.done
		if inflight.err != nil {
			log.Printf("Received error from in-flight request: %s", truncateKey(key))
			return nil, inflight.err, true
		}
		log.Printf("Received result from in-flight request: %s", truncateKey(key))
		return inflight.result, nil, true
	}

	return nil, nil, false
}

// startInFlightRequest registers a new in-flight request
func (h *Handler) startInFlightRequest(key string) *inFlightRequest {
	h.inFlightMu.Lock()
	defer h.inFlightMu.Unlock()

	req := &inFlightRequest{
		done: make(chan struct{}),
	}
	h.inFlight[key] = req

	log.Printf("Started in-flight request for key: %s", truncateKey(key))
	return req
}

func truncateKey(key string) string {
	if len(key) > 16 {
		return key[:16]
	}
	return key
}

// evictLRUEntry removes the least recently used cache entry - O(1)
func (h *Handler) evictLRUEntry() {
	// Remove from back of list (least recently used)
	elem := h.cacheList.Back()
	if elem != nil {
		h.cacheList.Remove(elem)
		entry := elem.Value.(*cacheEntry)
		delete(h.cache, entry.key)
		log.Printf("Cache LRU eviction: removed key %s", truncateKey(entry.key))
	}
}

// completeInFlightRequest broadcasts the result to all waiters and stores in cache
func (h *Handler) completeInFlightRequest(key string, logGroups []analyzer.LogGroup, err error) {
	h.inFlightMu.Lock()
	req, exists := h.inFlight[key]
	delete(h.inFlight, key)
	h.inFlightMu.Unlock()

	if !exists {
		return
	}

	// Store result in the request object
	req.resultOnce.Do(func() {
		req.result = logGroups
		req.err = err
	})

	// Store in cache if successful
	if err == nil {
		h.cacheMu.Lock()

		// Evict LRU entry if cache is at max size - O(1)
		if h.cacheList.Len() >= h.cacheMaxSize {
			h.evictLRUEntry()
		}

		// Add to front of list (most recently used) - O(1)
		entry := &cacheEntry{
			key:       key,
			logGroups: logGroups,
			expiresAt: time.Now().Add(h.cacheTTL),
		}
		elem := h.cacheList.PushFront(entry)
		h.cache[key] = elem

		h.cacheMu.Unlock()
		log.Printf("Cache SET for key: %s (TTL: %v, size: %d/%d)", truncateKey(key), h.cacheTTL, h.cacheList.Len(), h.cacheMaxSize)
	}

	// Broadcast to all waiters by closing the done channel
	if err != nil {
		log.Printf("Broadcasted error for key: %s", truncateKey(key))
	} else {
		log.Printf("Broadcasted result for key: %s", truncateKey(key))
	}
	close(req.done)
}

func (h *Handler) QueryLogs(w http.ResponseWriter, r *http.Request) {
	// Only allow POST
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed", "Only POST is allowed")
		return
	}

	var req QueryLogsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// Validate required fields
	if req.Org == "" || req.Dashboard == "" || req.PanelTitle == "" || req.MetricName == "" {
		writeJSONError(w, http.StatusBadRequest, "Invalid request", "Missing required fields")
		return
	}

	// Validate time range
	if !req.StartTime.Before(req.EndTime) {
		writeJSONError(w, http.StatusBadRequest, "Invalid time range", "Start time must be before end time")
		return
	}

	log.Printf("Processing log query - org: %s, dashboard: %s, panel: %s, metric: %s, time range: %v to %v",
		req.Org, req.Dashboard, req.PanelTitle, req.MetricName, req.StartTime, req.EndTime)

	// Generate cache key
	cacheKey := h.generateCacheKey(&req)

	// Check cache or wait for in-flight request
	if cachedLogGroups, cachedErr, found := h.getCachedResultOrWait(cacheKey); found {
		// If we got an error from a waiter, return error response
		if cachedErr != nil {
			log.Printf("Using cached error result: %v", cachedErr)
			writeJSONError(w, http.StatusInternalServerError, "Query failed", cachedErr.Error())
			return
		}

		// Convert to API response format
		apiLogGroups := make([]LogGroup, len(cachedLogGroups))
		for i, group := range cachedLogGroups {
			apiLogGroups[i] = LogGroup{
				RepresentativeLogs: group.RepresentativeLogs,
				RelativeChange:     group.RelativeChange,
			}
		}

		writeJSON(w, http.StatusOK, QueryLogsResponse{
			LogGroups: apiLogGroups,
		})
		return
	}

	// Start in-flight request tracking
	h.startInFlightRequest(cacheKey)

	// Analyze logs using KL divergence
	logGroups, err := h.analyzer.AnalyzeLogs(
		r.Context(),
		req.Org,
		req.Dashboard,
		req.PanelTitle,
		req.MetricName,
		req.StartTime,
		req.EndTime,
	)

	// Complete the in-flight request (broadcasts to waiters and stores in cache)
	h.completeInFlightRequest(cacheKey, logGroups, err)

	// If error occurred, return error response
	if err != nil {
		log.Printf("Error analyzing logs: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "Query failed", err.Error())
		return
	}

	// Convert to API response format
	apiLogGroups := make([]LogGroup, len(logGroups))
	for i, group := range logGroups {
		apiLogGroups[i] = LogGroup{
			RepresentativeLogs: group.RepresentativeLogs,
			RelativeChange:     group.RelativeChange,
		}
	}

	writeJSON(w, http.StatusOK, QueryLogsResponse{
		LogGroups: apiLogGroups,
	})
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeJSONError(w http.ResponseWriter, status int, error, message string) {
	writeJSON(w, status, ErrorResponse{
		Error:   error,
		Message: message,
		Code:    intPtr(status),
	})
}

func intPtr(i int) *int {
	return &i
}

// cleanupExpiredCache runs periodically to remove expired cache entries
func (h *Handler) cleanupExpiredCache() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		h.cacheMu.Lock()
		now := time.Now()
		expired := 0

		// Walk the list and remove expired entries - O(n) but only runs periodically
		var next *list.Element
		for elem := h.cacheList.Front(); elem != nil; elem = next {
			next = elem.Next()
			entry := elem.Value.(*cacheEntry)
			if now.After(entry.expiresAt) {
				h.cacheList.Remove(elem)
				delete(h.cache, entry.key)
				expired++
			}
		}

		if expired > 0 {
			log.Printf("Cache cleanup: removed %d expired entries, %d remaining", expired, h.cacheList.Len())
		}
		h.cacheMu.Unlock()
	}
}

