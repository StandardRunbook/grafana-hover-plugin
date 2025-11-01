package analyzer

import (
	"context"
	"log"
	"sort"
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
//
// Algorithm:
// 1. Query baseline window (same duration as current window, but before it)
// 2. Query current window (the anomaly window from Grafana)
// 3. Calculate template frequency distributions for both windows
// 4. Compute KL divergence to find anomalous templates
// 5. Fetch representative logs for top anomalous templates
func (la *LogAnalyzer) AnalyzeLogs(ctx context.Context, org, dashboard, panelTitle, metricName string, startTime, endTime time.Time) ([]LogGroup, error) {
	// Calculate baseline window (same duration as current window, but before it)
	windowDuration := endTime.Sub(startTime)
	baselineEnd := startTime
	baselineStart := baselineEnd.Add(-windowDuration)

	log.Printf("Analyzing logs - org: %s, dashboard: %s, panel: %s, metric: %s, current: %v to %v, baseline: %v to %v",
		org, dashboard, panelTitle, metricName, startTime, endTime, baselineStart, baselineEnd)

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

	// Fetch representative logs for these templates
	representatives, err := la.store.GetRepresentativeLogs(ctx, org, dashboard, panelTitle, metricName, topTemplateIDs)
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
