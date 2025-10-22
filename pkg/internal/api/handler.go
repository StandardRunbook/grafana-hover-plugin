package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"grafana-plugin-api/internal/analyzer"
	"grafana-plugin-api/internal/config"
)

type Handler struct {
	analyzer      *analyzer.LogAnalyzer
	analyzerError error
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
	logAnalyzer, err := analyzer.NewLogAnalyzer(&cfg.ClickHouse)
	if err != nil {
		log.Printf("Warning: Failed to create log analyzer: %v", err)
		log.Printf("Handler will return mock data for all requests")
		return &Handler{
			analyzer:      nil,
			analyzerError: err,
		}
	}

	return &Handler{
		analyzer:      logAnalyzer,
		analyzerError: nil,
	}
}

func (h *Handler) VerifyTables() error {
	if h.analyzer == nil {
		return h.analyzerError
	}
	return h.analyzer.VerifyTables()
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

	var logGroups []analyzer.LogGroup
	var err error

	// Check if analyzer is available
	if h.analyzer == nil {
		err = h.analyzerError
	} else {
		// Analyze logs using KL divergence
		logGroups, err = h.analyzer.AnalyzeLogs(
			r.Context(),
			req.Org,
			req.Dashboard,
			req.PanelTitle,
			req.MetricName,
			req.StartTime,
			req.EndTime,
		)
	}

	if err != nil {
		log.Printf("Error analyzing logs: %v", err)

		// Return mock data when ClickHouse is not available
		if containsStr(err.Error(), "Connection refused") ||
			containsStr(err.Error(), "connect") ||
			containsStr(err.Error(), "timeout") {

			// Return mock logs that clearly indicate they are examples
			logGroups = []analyzer.LogGroup{
				{
					RepresentativeLogs: []string{
						"⚠️  MOCK DATA: ClickHouse database is not connected",
						"📝 Example anomaly: ERROR: Out of memory on node-3",
						"📝 Example anomaly: WARNING: High CPU usage detected (95%)",
						"📝 Example anomaly: CRITICAL: Disk space below 5%",
					},
					RelativeChange: 2.5,
					KLContribution: 0.8,
					TemplateID:     "mock_error_template",
				},
				{
					RepresentativeLogs: []string{
						"📝 Example pattern: Connection timeout after 30s",
						"📝 Example pattern: Retrying connection attempt 3/5",
					},
					RelativeChange: 1.2,
					KLContribution: 0.4,
					TemplateID:     "mock_warning_template",
				},
				{
					RepresentativeLogs: []string{
						"📝 Example info: Service started successfully",
						"📝 Example info: Health check passed",
					},
					RelativeChange: 0.3,
					KLContribution: 0.1,
					TemplateID:     "mock_info_template",
				},
			}
		} else if containsStr(err.Error(), "does not exist") ||
			containsStr(err.Error(), "UNKNOWN_TABLE") {
			logGroups = []analyzer.LogGroup{
				{
					RepresentativeLogs: []string{
						"⚠️  Required tables missing. Please restart the service to auto-create tables.",
						"📝 MOCK DATA: These are example logs shown because tables don't exist yet",
					},
					RelativeChange: 0.0,
					KLContribution: 0.0,
					TemplateID:     "error",
				},
			}
		} else {
			// Generic error
			logGroups = []analyzer.LogGroup{
				{
					RepresentativeLogs: []string{
						"⚠️  Hover log database encountered an error. Please contact support.",
						"Error details: " + err.Error(),
					},
					RelativeChange: 0.0,
					KLContribution: 0.0,
					TemplateID:     "error",
				},
			}
		}
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

func containsStr(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
