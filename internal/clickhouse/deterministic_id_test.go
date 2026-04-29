package clickhouse

import (
	"strings"
	"testing"
)

// EnsureMapping relies on deterministicID for idempotency: two calls
// with the same inputs must produce the same ID so duplicate INSERTs
// are squashed by ReplacingMergeTree, while different inputs (kind,
// org, dashboard, panel, metric) must produce distinct IDs so unrelated
// panels don't collide on the same row. If this drifts, auto-
// registration silently overwrites neighboring panels' mappings.
func TestDeterministicIDStability(t *testing.T) {
	a := deterministicID("metric", "1", "Hover Plugin Demo", "CPU Usage", "cpu_usage")
	b := deterministicID("metric", "1", "Hover Plugin Demo", "CPU Usage", "cpu_usage")
	if a != b {
		t.Errorf("same inputs should produce same id: a=%q b=%q", a, b)
	}
	if !isHexLower(a) || len(a) != 24 {
		t.Errorf("id should be 24 lowercase hex chars, got %q (len %d)", a, len(a))
	}
}

func TestDeterministicIDSeparation(t *testing.T) {
	// kind separation
	if deterministicID("metric", "1", "d", "p", "m") == deterministicID("mapping", "1", "d", "p", "m") {
		t.Error("kind should affect id (metric vs mapping collided)")
	}
	// org separation — same panel under different orgs must not collide
	if deterministicID("metric", "1", "d", "p", "m") == deterministicID("metric", "2", "d", "p", "m") {
		t.Error("org_id should affect id (collision across orgs)")
	}
	// Each component independently affects the output
	base := deterministicID("metric", "1", "d", "p", "m")
	for _, alt := range []string{
		deterministicID("metric", "1", "d2", "p", "m"),
		deterministicID("metric", "1", "d", "p2", "m"),
		deterministicID("metric", "1", "d", "p", "m2"),
	} {
		if alt == base {
			t.Errorf("changing one component should change id; got %q == %q", alt, base)
		}
	}
}

// Guards against a delimiter-confusion attack: parts joined naively
// could let "a|b" + "c" collide with "a" + "b|c". The implementation
// uses a NUL separator, which is illegal in our string inputs, so
// these distinct logical inputs must hash to distinct ids.
func TestDeterministicIDDelimiterConfusion(t *testing.T) {
	x := deterministicID("metric", "1", "a", "b|c", "m")
	y := deterministicID("metric", "1", "a|b", "c", "m")
	if x == y {
		t.Errorf("delimiter confusion: %q == %q", x, y)
	}
}

func isHexLower(s string) bool {
	const hex = "0123456789abcdef"
	for _, c := range s {
		if !strings.ContainsRune(hex, c) {
			return false
		}
	}
	return true
}
