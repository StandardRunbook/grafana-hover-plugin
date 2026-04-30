package clickhouse

import (
	"context"
	"testing"
	"time"
)

// The MockStore is what the API handler defaults to when ClickHouse is
// unreachable, so its happy-path behavior matters: returns templates,
// returns representative logs filtered to the requested IDs, and the
// no-op writer-side methods don't error.

func TestMockStore_GetTemplateCounts_ReturnsFiveTemplates(t *testing.T) {
	m := NewMockStore()
	counts, err := m.GetTemplateCounts(
		context.Background(),
		"org-1", "dashboard", "panel", "metric",
		time.Now().Add(-time.Hour), time.Now(),
	)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(counts) != 5 {
		t.Errorf("expected 5 mock templates, got %d", len(counts))
	}
	// The "anomaly" template (highest count, drives the demo) must be present.
	if counts["error_template_1"] != 150 {
		t.Errorf("expected error_template_1=150, got %d", counts["error_template_1"])
	}
}

func TestMockStore_GetRepresentativeLogs_FiltersToRequested(t *testing.T) {
	m := NewMockStore()
	logs, err := m.GetRepresentativeLogs(
		context.Background(),
		"org-1", "dashboard", "panel", "metric",
		[]string{"error_template_1", "info_template_1"},
		time.Now().Add(-time.Hour), time.Now(),
	)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	// Asked for two specific templates → must get exactly those keys.
	if len(logs) != 2 {
		t.Errorf("expected 2 templates back, got %d", len(logs))
	}
	if _, ok := logs["error_template_1"]; !ok {
		t.Error("missing error_template_1 in result")
	}
	if _, ok := logs["info_template_1"]; !ok {
		t.Error("missing info_template_1 in result")
	}
	if _, ok := logs["debug_template_1"]; ok {
		t.Error("got debug_template_1 even though not requested")
	}
}

func TestMockStore_GetRepresentativeLogs_UnknownTemplateIDsAreSkipped(t *testing.T) {
	m := NewMockStore()
	logs, err := m.GetRepresentativeLogs(
		context.Background(),
		"o", "d", "p", "m",
		[]string{"made-up-template-id"},
		time.Now(), time.Now(),
	)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(logs) != 0 {
		t.Errorf("expected 0 results for unknown template ids, got %d", len(logs))
	}
}

func TestMockStore_GetTemplateCountsByPrefix_DelegatesToGetTemplateCounts(t *testing.T) {
	// The mock has no real `message` to filter against, so the prefix
	// variant just returns the same data. Pin that contract — if we ever
	// switch to a smart mock, the analyzer's empty-baseline-bias tests
	// need to stay green.
	m := NewMockStore()
	a, _ := m.GetTemplateCounts(context.Background(), "o", "d", "p", "m", time.Now(), time.Now())
	b, _ := m.GetTemplateCountsByPrefix(
		context.Background(), "o", "d", "p", "m", "any-prefix", time.Now(), time.Now(),
	)
	if len(a) != len(b) {
		t.Errorf("prefix-filtered call returned different count: %d vs %d", len(a), len(b))
	}
}

func TestMockStore_EnsureMapping_NoOp(t *testing.T) {
	if err := NewMockStore().EnsureMapping(context.Background(), "o", "d", "p", "m"); err != nil {
		t.Errorf("EnsureMapping should be a no-op for the mock, got err=%v", err)
	}
}

func TestMockStore_VerifyAndClose_NoOp(t *testing.T) {
	m := NewMockStore()
	if err := m.VerifyTables(); err != nil {
		t.Errorf("VerifyTables should be a no-op, got err=%v", err)
	}
	if err := m.Close(); err != nil {
		t.Errorf("Close should be a no-op, got err=%v", err)
	}
}
