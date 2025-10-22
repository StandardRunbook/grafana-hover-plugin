package analyzer

import (
	"context"
	"log"
	"sort"
	"time"

	"grafana-plugin-api/internal/clickhouse"
	"grafana-plugin-api/internal/config"
)

type LogAnalyzer struct {
	clickhouse *clickhouse.Client
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
		clickhouse: client,
	}, nil
}

func (la *LogAnalyzer) Close() error {
	return la.clickhouse.Close()
}

func (la *LogAnalyzer) VerifyTables() error {
	return la.clickhouse.VerifyTables()
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
	baselineCounts, err := la.clickhouse.GetTemplateCounts(ctx, org, dashboard, panelTitle, metricName, baselineStart, baselineEnd)
	if err != nil {
		return nil, err
	}

	currentCounts, err := la.clickhouse.GetTemplateCounts(ctx, org, dashboard, panelTitle, metricName, startTime, endTime)
	if err != nil {
		return nil, err
	}

	log.Printf("Found %d baseline templates, %d current templates", len(baselineCounts), len(currentCounts))

	// Calculate KL divergence contributions for each template
	klContributions := CalculateKLDivergence(currentCounts, baselineCounts)

	// Calculate relative changes for each template
	relativeChanges := CalculateRelativeChanges(currentCounts, baselineCounts)

	// Sort templates by KL divergence contribution (highest first)
	type templateKL struct {
		templateID string
		klValue    float64
	}

	var sortedTemplates []templateKL
	for templateID, klValue := range klContributions {
		sortedTemplates = append(sortedTemplates, templateKL{templateID, klValue})
	}

	sort.Slice(sortedTemplates, func(i, j int) bool {
		return sortedTemplates[i].klValue > sortedTemplates[j].klValue
	})

	// Take top N templates with highest KL divergence
	topN := 10
	if len(sortedTemplates) > topN {
		sortedTemplates = sortedTemplates[:topN]
	}

	if len(sortedTemplates) == 0 {
		log.Println("No templates found with significant KL divergence")
		return []LogGroup{}, nil
	}

	// Extract template IDs
	var topTemplateIDs []string
	for _, t := range sortedTemplates {
		topTemplateIDs = append(topTemplateIDs, t.templateID)
	}

	// Fetch representative logs for these templates
	representatives, err := la.clickhouse.GetRepresentativeLogs(ctx, org, dashboard, panelTitle, metricName, topTemplateIDs)
	if err != nil {
		return nil, err
	}

	// Build log groups
	var logGroups []LogGroup
	for _, templateID := range topTemplateIDs {
		if logs, ok := representatives[templateID]; ok {
			relativeChange := relativeChanges[templateID]
			klContribution := klContributions[templateID]

			logGroups = append(logGroups, LogGroup{
				RepresentativeLogs: logs,
				RelativeChange:     relativeChange,
				KLContribution:     klContribution,
				TemplateID:         templateID,
			})
		}
	}

	log.Printf("Returning %d log groups", len(logGroups))

	return logGroups, nil
}
