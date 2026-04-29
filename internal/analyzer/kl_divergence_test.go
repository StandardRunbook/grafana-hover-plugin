package analyzer

import (
	"math"
	"testing"
)

func TestCalculateJSDivergence(t *testing.T) {
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
			result := CalculateJSDivergence(tt.currentCounts, tt.baselineCounts)

			// Check that expected templates have non-zero divergence contribution.
			for _, templateID := range tt.expectNonZero {
				if val, ok := result[templateID]; !ok || math.Abs(val) < 1e-9 {
					t.Errorf("Expected non-zero JS divergence for %s, got %v", templateID, val)
				}
			}

			// Check that the specified template has the highest contribution.
			if tt.expectHigherFor != "" {
				maxJS := -math.MaxFloat64
				maxTemplate := ""
				for id, js := range result {
					if js > maxJS {
						maxJS = js
						maxTemplate = id
					}
				}
				if maxTemplate != tt.expectHigherFor {
					t.Errorf("Expected highest JS divergence for %s, got %s (JS: %v)", tt.expectHigherFor, maxTemplate, maxJS)
				}
			}

			// All values must be finite and non-negative (JS is bounded [0, 1]).
			for id, val := range result {
				if math.IsNaN(val) || math.IsInf(val, 0) {
					t.Errorf("JS divergence for %s is not finite: %v", id, val)
				}
				if val < 0 {
					t.Errorf("JS divergence for %s is negative: %v (JS must be >= 0)", id, val)
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

// JS divergence is symmetric: D_JS(P || Q) == D_JS(Q || P). KL isn't —
// this test was originally named TestCalculateKLDivergenceSymmetry but
// only checked finiteness, which masked the algorithm switch. Now it
// asserts the actual symmetry property the algorithm gives us.
func TestCalculateJSDivergenceSymmetry(t *testing.T) {
	currentCounts := map[string]uint64{
		"template_001": 10,
		"template_002": 5,
	}
	baselineCounts := map[string]uint64{
		"template_001": 5,
		"template_002": 10,
	}

	js1 := CalculateJSDivergence(currentCounts, baselineCounts)
	js2 := CalculateJSDivergence(baselineCounts, currentCounts)

	if len(js1) == 0 || len(js2) == 0 {
		t.Fatal("Both JS divergence calculations should produce results")
	}

	const tol = 1e-9
	for id, v1 := range js1 {
		v2, ok := js2[id]
		if !ok {
			t.Errorf("template %s present in js1 but missing from js2", id)
			continue
		}
		if math.IsNaN(v1) || math.IsNaN(v2) {
			t.Errorf("JS divergence for %s contains NaN: js1=%v js2=%v", id, v1, v2)
			continue
		}
		if math.Abs(v1-v2) > tol {
			t.Errorf("JS not symmetric for %s: D(P,Q)=%v vs D(Q,P)=%v (diff %v)", id, v1, v2, math.Abs(v1-v2))
		}
	}
}
