package analyzer

import (
	"math"
)

const smoothing = 1e-10

// CalculateKLDivergence calculates KL divergence contribution for each template
func CalculateKLDivergence(currentCounts, baselineCounts map[string]uint64) map[string]float64 {
	// Calculate total counts
	var currentTotal, baselineTotal uint64
	for _, count := range currentCounts {
		currentTotal += count
	}
	for _, count := range baselineCounts {
		baselineTotal += count
	}

	// Handle empty cases
	if currentTotal == 0 || baselineTotal == 0 {
		return make(map[string]float64)
	}

	// Get all unique template IDs
	allTemplates := make(map[string]bool)
	for id := range currentCounts {
		allTemplates[id] = true
	}
	for id := range baselineCounts {
		allTemplates[id] = true
	}

	// Calculate KL divergence contribution for each template
	klContributions := make(map[string]float64)

	for templateID := range allTemplates {
		currentCount := currentCounts[templateID]
		baselineCount := baselineCounts[templateID]

		// Calculate probabilities with smoothing
		pCurrent := (float64(currentCount) + smoothing) / (float64(currentTotal) + smoothing*float64(len(allTemplates)))
		pBaseline := (float64(baselineCount) + smoothing) / (float64(baselineTotal) + smoothing*float64(len(allTemplates)))

		// KL divergence contribution: P(current) * log(P(current) / P(baseline))
		klContribution := pCurrent * math.Log(pCurrent/pBaseline)
		klContributions[templateID] = klContribution
	}

	return klContributions
}

// CalculateRelativeChanges calculates relative changes for each template
func CalculateRelativeChanges(currentCounts, baselineCounts map[string]uint64) map[string]float64 {
	// Calculate total counts
	var currentTotal, baselineTotal uint64
	for _, count := range currentCounts {
		currentTotal += count
	}
	for _, count := range baselineCounts {
		baselineTotal += count
	}

	// Handle empty cases
	if currentTotal == 0 || baselineTotal == 0 {
		return make(map[string]float64)
	}

	// Get all unique template IDs
	allTemplates := make(map[string]bool)
	for id := range currentCounts {
		allTemplates[id] = true
	}
	for id := range baselineCounts {
		allTemplates[id] = true
	}

	// Calculate relative changes
	relativeChanges := make(map[string]float64)

	for templateID := range allTemplates {
		currentCount := currentCounts[templateID]
		baselineCount := baselineCounts[templateID]

		// Calculate frequencies
		freqCurrent := float64(currentCount) / float64(currentTotal)
		freqBaseline := float64(baselineCount) / float64(baselineTotal)

		// Relative change: (current - baseline) / baseline
		// Use smoothing to avoid division by zero
		relativeChange := (freqCurrent - freqBaseline) / (freqBaseline + smoothing)
		relativeChanges[templateID] = relativeChange
	}

	return relativeChanges
}
