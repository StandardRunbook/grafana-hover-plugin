package analyzer

import (
	"reflect"
	"testing"
)

// biasCandidates feeds the empty-baseline relevance bias: probes are
// tried longest-first, stopping at the first match. The behavior matters
// for hover-time UX:
//   - cpu_usage   → matches "cpu_usage:" log lines on the literal probe
//   - error_rate  → no literal match against "ERROR:", but the "error"
//                   underscore-segment does match (case-insensitive)
//   - latency_p99 → falls back to "latency"
// If this drift breaks, hovering "error_rate" silently returns
// whatever template happens to dominate the stream — the most
// noticeable demo regression we've had.
func TestBiasCandidates(t *testing.T) {
	cases := []struct {
		in   string
		want []string
	}{
		{in: "", want: nil},
		{in: "cpu", want: []string{"cpu"}},
		{in: "cpu_usage", want: []string{"cpu_usage", "cpu"}},
		{in: "error_rate", want: []string{"error_rate", "error"}},
		{in: "latency_p99", want: []string{"latency_p99", "latency"}},
		// Three segments: full, then leading-2, then leading-1 (if >= 3 chars).
		{in: "request_duration_p99", want: []string{
			"request_duration_p99",
			"request_duration",
			"request",
		}},
		// Short leading segment (2 chars) is dropped — too noisy to be useful.
		{in: "io_bytes", want: []string{"io_bytes"}},
	}

	for _, tc := range cases {
		t.Run(tc.in, func(t *testing.T) {
			got := biasCandidates(tc.in)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("biasCandidates(%q) = %v, want %v", tc.in, got, tc.want)
			}
		})
	}
}
