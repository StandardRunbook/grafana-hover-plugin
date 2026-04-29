package analyzer

import (
	"context"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/StandardRunbook/grafana-hover-plugin/internal/clickhouse"
	"github.com/StandardRunbook/grafana-hover-plugin/internal/config"
)

type LogAnalyzer struct {
	store clickhouse.Store
}

type LogGroup struct {
	RepresentativeLogs []string `json:"representative_logs"`
	RelativeChange     float64  `json:"relative_change"`
	KLContribution     float64  `json:"kl_contribution"`
	TemplateID         string   `json:"template_id"`
}

func NewLogAnalyzer(cfg *config.ClickHouseConfig) (*LogAnalyzer, error) {
	client, err := clickhouse.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return &LogAnalyzer{
		store: client,
	}, nil
}

// NewLogAnalyzerWithStore creates a LogAnalyzer with a custom store (for testing/mocking)
func NewLogAnalyzerWithStore(store clickhouse.Store) *LogAnalyzer {
	return &LogAnalyzer{
		store: store,
	}
}

func (la *LogAnalyzer) Close() error {
	return la.store.Close()
}

func (la *LogAnalyzer) VerifyTables() error {
	return la.store.VerifyTables()
}

// AnalyzeLogs analyzes logs for anomalies using KL divergence
// This method is specifically for the Hover backend, which uses pre-calculated
// templates stored in ClickHouse during log ingestion.
//
// Algorithm:
// 1. Query baseline window (same duration as current window, but before it)
// 2. Query current window (the anomaly window from Grafana)
// 3. Calculate template frequency distributions for both windows (using pre-calculated template_ids)
// 4. Compute JS divergence to find anomalous templates
// 5. Fetch representative logs for top anomalous templates from template_examples table
// baselineLookback is the length of the historical window the analyzer
// compares the current (hover) window against. A long lookback captures
// "what the stream normally looks like" so a tight current-window
// distribution can diverge from it visibly. Adjacent same-length windows
// (the original choice) hide gradual incidents and produce ~0 divergence
// on a uniform stream — we couldn't tell whether the algorithm worked or
// not. 2 hours is enough to span typical diurnal patterns without
// sweeping in week-old templates that no longer exist.
const baselineLookback = 2 * time.Hour

func (la *LogAnalyzer) AnalyzeLogs(ctx context.Context, org, dashboard, panelTitle, metricName string, startTime, endTime time.Time) ([]LogGroup, error) {
	// Baseline = a long trailing window ending where the current window
	// starts (no gap). The current window is whatever the caller hovered.
	baselineEnd := startTime
	baselineStart := baselineEnd.Add(-baselineLookback)

	log.Printf("Analyzing logs - org: %s, dashboard: %s, panel: %s, metric: %s, current: %v to %v, baseline: %v to %v",
		org, dashboard, panelTitle, metricName, startTime, endTime, baselineStart, baselineEnd)

	// Auto-register the (dashboard, panel, metric) → log-stream mapping so
	// adding a new panel doesn't require a manual SQL insert. Idempotent;
	// no-op for already-mapped panels and orgs without any ingested logs.
	if err := la.store.EnsureMapping(ctx, org, dashboard, panelTitle, metricName); err != nil {
		log.Printf("EnsureMapping warning (continuing): %v", err)
	}

	// Get template counts for both windows
	baselineCounts, err := la.store.GetTemplateCounts(ctx, org, dashboard, panelTitle, metricName, baselineStart, baselineEnd)
	if err != nil {
		return nil, err
	}

	currentCounts, err := la.store.GetTemplateCounts(ctx, org, dashboard, panelTitle, metricName, startTime, endTime)
	if err != nil {
		return nil, err
	}

	log.Printf("Found %d baseline templates, %d current templates", len(baselineCounts), len(currentCounts))

	// Calculate Jensen-Shannon Distance contributions for each template
	jsContributions := CalculateJSDivergence(currentCounts, baselineCounts)

	// Calculate relative changes for each template (as percentages)
	relativeChanges := CalculateRelativeChanges(currentCounts, baselineCounts)

	// Empty-baseline fallback: if there's no baseline data to compare
	// against (cold start, new metric, sparse history) but we DO have
	// current-window data, treat every current template as a 100% shift
	// and rank by raw count. Otherwise this whole flow returns no
	// log_groups whenever the baseline is empty, which surfaces as a
	// silent "nothing to show" in the hover panel even when the current
	// window has interesting logs.
	if len(baselineCounts) == 0 && len(currentCounts) > 0 {
		// Metric-name relevance bias: when the user hovers the CPU
		// panel they want cpu_usage logs, not whatever template
		// happens to dominate the stream. Try a substring-filtered
		// recount first; only fall back to all-templates if the
		// filtered set is empty.
		biasedCounts := currentCounts
		// Try the full metric name first, then progressively shorter
		// underscore-separated segments. e.g. metric_name='error_rate'
		// won't substring-match 'ERROR: connection...', but its first
		// segment 'error' does. Stop at the first non-empty match.
		for _, candidate := range biasCandidates(metricName) {
			prefixCounts, perr := la.store.GetTemplateCountsByPrefix(
				ctx, org, dashboard, panelTitle, metricName, candidate, startTime, endTime,
			)
			if perr != nil {
				log.Printf("Metric-name bias query failed for %q (continuing): %v", candidate, perr)
				break
			}
			if len(prefixCounts) > 0 {
				log.Printf("Empty baseline + metric_name='%s' (probe='%s') → biasing to %d matching templates",
					metricName, candidate, len(prefixCounts))
				biasedCounts = prefixCounts
				break
			}
		}
		log.Printf("Empty baseline; ranking %d templates as new (100%% shift)", len(biasedCounts))
		jsContributions = make(map[string]float64, len(biasedCounts))
		relativeChanges = make(map[string]float64, len(biasedCounts))
		for templateID, count := range biasedCounts {
			// Use raw count as the contribution score so the existing
			// sort-by-jsValue path picks the busiest templates first.
			jsContributions[templateID] = float64(count)
			relativeChanges[templateID] = 100.0
		}
	}

	// Sort templates by JS divergence contribution (highest first)
	type templateJS struct {
		templateID string
		jsValue    float64
	}

	var sortedTemplates []templateJS
	for templateID, jsValue := range jsContributions {
		sortedTemplates = append(sortedTemplates, templateJS{templateID, jsValue})
	}

	sort.Slice(sortedTemplates, func(i, j int) bool {
		return sortedTemplates[i].jsValue > sortedTemplates[j].jsValue
	})

	// Take top N templates with highest JS divergence
	topN := 10
	if len(sortedTemplates) > topN {
		sortedTemplates = sortedTemplates[:topN]
	}

	if len(sortedTemplates) == 0 {
		log.Println("No templates found with significant divergence")
		return []LogGroup{}, nil
	}

	// Extract template IDs
	var topTemplateIDs []string
	for _, t := range sortedTemplates {
		topTemplateIDs = append(topTemplateIDs, t.templateID)
	}

	// Fetch representative logs for these templates from the *current*
	// window — that's the slice the user is hovering on. Per the
	// hover-content invariant, we surface the actual messages from this
	// (template_id, time_window) intersection, not a sample across time.
	representatives, err := la.store.GetRepresentativeLogs(
		ctx, org, dashboard, panelTitle, metricName, topTemplateIDs,
		startTime, endTime,
	)
	if err != nil {
		return nil, err
	}

	// Build log groups
	var logGroups []LogGroup
	for _, templateID := range topTemplateIDs {
		if logs, ok := representatives[templateID]; ok {
			relativeChange := relativeChanges[templateID]
			jsContribution := jsContributions[templateID]

			logGroups = append(logGroups, LogGroup{
				RepresentativeLogs: logs,
				RelativeChange:     relativeChange,
				KLContribution:     jsContribution, // Now contains JS divergence
				TemplateID:         templateID,
			})
		}
	}

	log.Printf("Returning %d log groups", len(logGroups))

	return logGroups, nil
}

// biasCandidates returns substrings to probe in priority order: the full
// metric_name first, then each leading segment from longest to shortest
// (split on `_`). Empty string and the trivial 1-char fallback are
// excluded. Stops once a probe finds matching templates.
func biasCandidates(metricName string) []string {
	if metricName == "" {
		return nil
	}
	out := []string{metricName}
	parts := strings.Split(metricName, "_")
	for i := len(parts) - 1; i > 0; i-- {
		candidate := strings.Join(parts[:i], "_")
		if len(candidate) >= 3 {
			out = append(out, candidate)
		}
	}
	return out
}

// AnalyzeRawLogs analyzes raw logs without using pre-calculated template IDs
// This is for standalone deployments using existing log stores (Loki, ClickHouse, Elasticsearch)
// WITHOUT the Hover pre-processing pipeline that calculates templates during ingestion.
func (la *LogAnalyzer) AnalyzeRawLogs(ctx context.Context, backend, org, dashboard, panelTitle, metricName string, startTime, endTime time.Time) ([]LogGroup, error) {
	log.Printf("AnalyzeRawLogs - backend: %s, org: %s, dashboard: %s, panel: %s, metric: %s, time: %v to %v",
		backend, org, dashboard, panelTitle, metricName, startTime, endTime)

	return la.AnalyzeRawLogsWithExtractor(ctx, backend, org, dashboard, panelTitle, metricName, startTime, endTime, nil, nil)
}

// AnalyzeRawLogsWithExtractor analyzes raw logs using custom extractor and parser
// If extractor or parser is nil, defaults will be created based on backend type
func (la *LogAnalyzer) AnalyzeRawLogsWithExtractor(
	ctx context.Context,
	backend, org, dashboard, panelTitle, metricName string,
	startTime, endTime time.Time,
	extractor interface{}, // extractor.LogExtractor (avoiding import cycle)
	parser interface{},    // parser.LogParser (avoiding import cycle)
) ([]LogGroup, error) {
	// This will be implemented when extractor and parser packages are imported
	// For now, return a helpful message
	return []LogGroup{
		{
			RepresentativeLogs: []string{
				"Raw log analysis with new extractor/parser architecture",
				"Backend: " + backend,
				"Use AnalyzeRawLogsV2 for full implementation",
			},
			RelativeChange: 0.0,
			KLContribution: 0.0,
			TemplateID:     "stub",
		},
	}, nil
}
