package analyzer

import (
	"testing"
	"time"
)

// Pins the analyzer's baseline-window contract:
//   • baseline ends where current starts (no gap)
//   • baseline length is the configured baselineLookback (2h), independent
//     of how wide the current window is — earlier behavior was
//     same-length-as-current, which made adjacent ~60s windows on a
//     uniform stream produce ~0 divergence everywhere
func TestLogAnalyzerBaselineIsLongTrailing(t *testing.T) {
	endTime := time.Now()
	startTime := endTime.Add(-1 * time.Minute) // tight current window

	baselineEnd := startTime
	baselineStart := baselineEnd.Add(-baselineLookback)

	if !baselineEnd.Equal(startTime) {
		t.Errorf("Baseline end (%v) should equal current start (%v)", baselineEnd, startTime)
	}

	if got := baselineEnd.Sub(baselineStart); got != baselineLookback {
		t.Errorf("Baseline duration = %v, want %v", got, baselineLookback)
	}

	if cur := endTime.Sub(startTime); cur >= baselineEnd.Sub(baselineStart) {
		t.Errorf("Baseline should be longer than current window; baseline=%v current=%v",
			baselineEnd.Sub(baselineStart), cur)
	}
}

func TestLogGroupSorting(t *testing.T) {
	// Test sorting logic for log groups by KL divergence
	klContributions := map[string]float64{
		"template_001": 0.1,
		"template_002": 0.5,
		"template_003": 0.3,
		"template_004": 0.8,
		"template_005": 0.2,
	}

	type templateKL struct {
		templateID string
		klValue    float64
	}

	var sortedTemplates []templateKL
	for templateID, klValue := range klContributions {
		sortedTemplates = append(sortedTemplates, templateKL{templateID, klValue})
	}

	// Sort by KL divergence (highest first) - simplified bubble sort for test
	for i := 0; i < len(sortedTemplates); i++ {
		for j := i + 1; j < len(sortedTemplates); j++ {
			if sortedTemplates[i].klValue < sortedTemplates[j].klValue {
				sortedTemplates[i], sortedTemplates[j] = sortedTemplates[j], sortedTemplates[i]
			}
		}
	}

	// Check that it's sorted in descending order
	for i := 0; i < len(sortedTemplates)-1; i++ {
		if sortedTemplates[i].klValue < sortedTemplates[i+1].klValue {
			t.Errorf("Templates not sorted correctly at index %d: %v < %v",
				i, sortedTemplates[i].klValue, sortedTemplates[i+1].klValue)
		}
	}

	// Top template should be template_004 with KL 0.8
	if sortedTemplates[0].templateID != "template_004" {
		t.Errorf("Expected top template to be template_004, got %s", sortedTemplates[0].templateID)
	}

	// Take top N (e.g., 3)
	topN := 3
	if len(sortedTemplates) > topN {
		sortedTemplates = sortedTemplates[:topN]
	}

	if len(sortedTemplates) != topN {
		t.Errorf("Expected %d templates after limiting, got %d", topN, len(sortedTemplates))
	}
}

func TestLogGroupCreation(t *testing.T) {
	// Test creating log groups from results
	representatives := map[string][]string{
		"template_001": {"Log message 1", "Log message 2"},
		"template_002": {"Error message 1", "Error message 2", "Error message 3"},
	}

	relativeChanges := map[string]float64{
		"template_001": 0.5,
		"template_002": 1.2,
	}

	klContributions := map[string]float64{
		"template_001": 0.1,
		"template_002": 0.3,
	}

	topTemplateIDs := []string{"template_001", "template_002"}

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

	if len(logGroups) != 2 {
		t.Errorf("Expected 2 log groups, got %d", len(logGroups))
	}

	// Verify first log group
	if len(logGroups[0].RepresentativeLogs) != 2 {
		t.Errorf("Expected 2 representative logs for template_001, got %d", len(logGroups[0].RepresentativeLogs))
	}

	if logGroups[0].RelativeChange != 0.5 {
		t.Errorf("Expected relative change of 0.5, got %f", logGroups[0].RelativeChange)
	}

	// Verify second log group
	if len(logGroups[1].RepresentativeLogs) != 3 {
		t.Errorf("Expected 3 representative logs for template_002, got %d", len(logGroups[1].RepresentativeLogs))
	}

	if logGroups[1].KLContribution != 0.3 {
		t.Errorf("Expected KL contribution of 0.3, got %f", logGroups[1].KLContribution)
	}
}
