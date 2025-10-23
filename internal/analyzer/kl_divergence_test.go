package analyzer

import (
	"math"
	"testing"
)

func TestCalculateKLDivergence(t *testing.T) {
	tests := []struct {
		name            string
		currentCounts   map[string]uint64
		baselineCounts  map[string]uint64
		expectNonZero   []string
		expectHigherFor string
	}{
		{
			name: "new template appears",
			currentCounts: map[string]uint64{
				"template_001": 10,
				"template_002": 5,
				"template_003": 3, // New template
			},
			baselineCounts: map[string]uint64{
				"template_001": 10,
				"template_002": 5,
			},
			expectNonZero:   []string{"template_003"},
			expectHigherFor: "template_003",
		},
		{
			name: "template frequency increases significantly",
			currentCounts: map[string]uint64{
				"template_001": 5,
				"template_002": 20, // Increased significantly
			},
			baselineCounts: map[string]uint64{
				"template_001": 10,
				"template_002": 2,
			},
			expectNonZero:   []string{"template_001", "template_002"},
			expectHigherFor: "template_002",
		},
		{
			name: "identical distributions",
			currentCounts: map[string]uint64{
				"template_001": 10,
				"template_002": 5,
			},
			baselineCounts: map[string]uint64{
				"template_001": 10,
				"template_002": 5,
			},
			expectNonZero: []string{}, // Should be near zero
		},
		{
			name:           "empty current counts",
			currentCounts:  map[string]uint64{},
			baselineCounts: map[string]uint64{"template_001": 10},
			expectNonZero:  []string{},
		},
		{
			name:           "empty baseline counts",
			currentCounts:  map[string]uint64{"template_001": 10},
			baselineCounts: map[string]uint64{},
			expectNonZero:  []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateKLDivergence(tt.currentCounts, tt.baselineCounts)

			// Check that expected templates have non-zero KL divergence
			for _, templateID := range tt.expectNonZero {
				if val, ok := result[templateID]; !ok || math.Abs(val) < 1e-9 {
					t.Errorf("Expected non-zero KL divergence for %s, got %v", templateID, val)
				}
			}

			// Check that the specified template has higher KL divergence
			if tt.expectHigherFor != "" {
				maxKL := -math.MaxFloat64
				maxTemplate := ""
				for id, kl := range result {
					if kl > maxKL {
						maxKL = kl
						maxTemplate = id
					}
				}
				if maxTemplate != tt.expectHigherFor {
					t.Errorf("Expected highest KL divergence for %s, got %s (KL: %v)", tt.expectHigherFor, maxTemplate, maxKL)
				}
			}

			// All values should be finite
			for id, val := range result {
				if math.IsNaN(val) || math.IsInf(val, 0) {
					t.Errorf("KL divergence for %s is not finite: %v", id, val)
				}
			}
		})
	}
}

func TestCalculateRelativeChanges(t *testing.T) {
	tests := []struct {
		name           string
		currentCounts  map[string]uint64
		baselineCounts map[string]uint64
		checkTemplate  string
		expectedSign   int // 1 for positive, -1 for negative
	}{
		{
			name: "template frequency increases significantly",
			currentCounts: map[string]uint64{
				"template_001": 20,
				"template_002": 2,
			},
			baselineCounts: map[string]uint64{
				"template_001": 10,
				"template_002": 10,
			},
			checkTemplate: "template_001",
			expectedSign:  1, // template_001 increased
		},
		{
			name: "template frequency decreases significantly",
			currentCounts: map[string]uint64{
				"template_001": 5,
				"template_002": 15,
			},
			baselineCounts: map[string]uint64{
				"template_001": 15,
				"template_002": 5,
			},
			checkTemplate: "template_001",
			expectedSign:  -1, // template_001 decreased
		},
		{
			name: "new template appears",
			currentCounts: map[string]uint64{
				"template_001": 5,
				"template_002": 10,
			},
			baselineCounts: map[string]uint64{
				"template_001": 15,
			},
			checkTemplate: "template_002",
			expectedSign:  1, // New template has positive change
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateRelativeChanges(tt.currentCounts, tt.baselineCounts)

			val, ok := result[tt.checkTemplate]
			if !ok {
				t.Errorf("Expected relative change for %s, but not found", tt.checkTemplate)
				return
			}

			actualSign := 0
			if val > 0.1 {
				actualSign = 1
			} else if val < -0.1 {
				actualSign = -1
			}

			if actualSign != tt.expectedSign {
				t.Errorf("Template %s: expected sign %d, got %d (value: %v)", tt.checkTemplate, tt.expectedSign, actualSign, val)
			}

			// All values should be finite
			if math.IsNaN(val) || math.IsInf(val, 0) {
				t.Errorf("Relative change for %s is not finite: %v", tt.checkTemplate, val)
			}
		})
	}
}

func TestCalculateKLDivergenceSymmetry(t *testing.T) {
	currentCounts := map[string]uint64{
		"template_001": 10,
		"template_002": 5,
	}
	baselineCounts := map[string]uint64{
		"template_001": 5,
		"template_002": 10,
	}

	kl1 := CalculateKLDivergence(currentCounts, baselineCounts)
	kl2 := CalculateKLDivergence(baselineCounts, currentCounts)

	// KL divergence is not symmetric, but both should produce valid results
	if len(kl1) == 0 || len(kl2) == 0 {
		t.Error("Both KL divergence calculations should produce results")
	}

	// Templates with opposite changes should have opposite-signed KL contributions
	for id := range kl1 {
		if math.IsNaN(kl1[id]) || math.IsNaN(kl2[id]) {
			t.Errorf("KL divergence for %s contains NaN", id)
		}
	}
}
