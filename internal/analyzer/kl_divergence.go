package analyzer

import (
	"math"
)

const smoothing = 1e-10

// CalculateJSDivergence calculates Jensen-Shannon Distance contribution for each template
// JSD is symmetric, bounded [0,1], and handles zeros gracefully
func CalculateJSDivergence(currentCounts, baselineCounts map[string]uint64) map[string]float64 {
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

	// Calculate probabilities for all templates
	pCurrent := make(map[string]float64)
	pBaseline := make(map[string]float64)
	pMean := make(map[string]float64)

	vocabSize := float64(len(allTemplates))
	for templateID := range allTemplates {
		// Add-one smoothing (Laplace smoothing)
		pCurrent[templateID] = (float64(currentCounts[templateID]) + 1.0) / (float64(currentTotal) + vocabSize)
		pBaseline[templateID] = (float64(baselineCounts[templateID]) + 1.0) / (float64(baselineTotal) + vocabSize)
		pMean[templateID] = (pCurrent[templateID] + pBaseline[templateID]) / 2.0
	}

	// Calculate JSD contribution for each template
	// JSD = 0.5 * KL(P||M) + 0.5 * KL(Q||M)
	// We calculate the per-template contribution to this divergence
	jsContributions := make(map[string]float64)

	for templateID := range allTemplates {
		p := pCurrent[templateID]
		q := pBaseline[templateID]
		m := pMean[templateID]

		// Per-template JSD contribution
		klPM := 0.0
		klQM := 0.0
		if p > 0 && m > 0 {
			klPM = p * math.Log2(p/m)
		}
		if q > 0 && m > 0 {
			klQM = q * math.Log2(q/m)
		}

		// JSD contribution (normalized to [0,1] by using log2)
		jsContribution := 0.5*klPM + 0.5*klQM
		jsContributions[templateID] = jsContribution
	}

	return jsContributions
}

// CalculateRelativeChanges calculates percentage changes in frequency for each template
// Returns values like 1.5 (150% increase), -0.5 (50% decrease), etc.
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
	if currentTotal == 0 && baselineTotal == 0 {
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

	// Calculate relative changes as percentages
	relativeChanges := make(map[string]float64)
	vocabSize := float64(len(allTemplates))

	for templateID := range allTemplates {
		// Use add-one smoothing for stable frequency estimates
		freqCurrent := (float64(currentCounts[templateID]) + 1.0) / (float64(currentTotal) + vocabSize)
		freqBaseline := (float64(baselineCounts[templateID]) + 1.0) / (float64(baselineTotal) + vocabSize)

		// Percentage change: ((current - baseline) / baseline) * 100
		// This gives values like:
		// - 1.0 = 100% increase (doubled)
		// - 0.5 = 50% increase
		// - -0.5 = 50% decrease
		// - -1.0 = 100% decrease (went to zero)
		percentChange := (freqCurrent - freqBaseline) / freqBaseline
		relativeChanges[templateID] = percentChange
	}

	return relativeChanges
}
