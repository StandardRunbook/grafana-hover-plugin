package analyzer

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/StandardRunbook/grafana-hover-plugin/internal/clickhouse"
)

// fakeStore is a Store implementation hand-rolled for tests so we can:
//   - vary what GetTemplateCounts returns per (window) call,
//   - track EnsureMapping invocations,
//   - inject errors on demand.
type fakeStore struct {
	current        map[string]uint64
	baseline       map[string]uint64
	prefixCounts   map[string]uint64
	representative map[string][]string
	ensureCalls    int
	tcCalls        int

	templateCountsErr error
	representativeErr error
	ensureMappingErr  error
}

func (f *fakeStore) GetTemplateCounts(ctx context.Context, org, dashboard, panel, metric string, start, end time.Time) (map[string]uint64, error) {
	if f.templateCountsErr != nil {
		return nil, f.templateCountsErr
	}
	// AnalyzeLogs calls GetTemplateCounts twice in this order: baseline first,
	// then current. Track call number and route accordingly.
	f.tcCalls++
	if f.tcCalls == 1 {
		return cloneMap(f.baseline), nil
	}
	return cloneMap(f.current), nil
}

func (f *fakeStore) GetTemplateCountsByPrefix(ctx context.Context, org, dashboard, panel, metric, prefix string, start, end time.Time) (map[string]uint64, error) {
	return cloneMap(f.prefixCounts), nil
}

func (f *fakeStore) GetRepresentativeLogs(ctx context.Context, org, dashboard, panel, metric string, ids []string, start, end time.Time) (map[string][]string, error) {
	if f.representativeErr != nil {
		return nil, f.representativeErr
	}
	out := make(map[string][]string)
	for _, id := range ids {
		if logs, ok := f.representative[id]; ok {
			out[id] = logs
		}
	}
	return out, nil
}

func (f *fakeStore) EnsureMapping(ctx context.Context, org, dashboard, panel, metric string) error {
	f.ensureCalls++
	return f.ensureMappingErr
}

func (f *fakeStore) VerifyTables() error { return nil }
func (f *fakeStore) Close() error        { return nil }

func cloneMap(m map[string]uint64) map[string]uint64 {
	out := make(map[string]uint64, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

// satisfy the Store interface
var _ clickhouse.Store = (*fakeStore)(nil)

func TestAnalyzeLogs_HappyPath_RanksByJSContribution(t *testing.T) {
	// Baseline has stable distribution; current has a spike on template
	// `incident` that didn't exist before. JS divergence should rank
	// `incident` first.
	store := &fakeStore{
		baseline: map[string]uint64{
			"normal_a": 100,
			"normal_b": 80,
			"normal_c": 50,
		},
		current: map[string]uint64{
			"normal_a": 100,
			"normal_b": 80,
			"normal_c": 50,
			"incident": 200, // brand new
		},
		representative: map[string][]string{
			"incident": {"CRITICAL: db pool exhausted on shard-2"},
			"normal_a": {"normal log a"},
			"normal_b": {"normal log b"},
			"normal_c": {"normal log c"},
		},
	}
	la := NewLogAnalyzerWithStore(store)
	groups, err := la.AnalyzeLogs(
		context.Background(),
		"org-1", "CPU Usage", "CPU Usage", "cpu_usage",
		time.Now().Add(-1*time.Minute), time.Now(),
	)
	if err != nil {
		t.Fatalf("AnalyzeLogs: %v", err)
	}
	if len(groups) == 0 {
		t.Fatal("expected at least one log group, got 0")
	}
	if groups[0].TemplateID != "incident" {
		t.Errorf("expected incident ranked first, got %q (full: %+v)", groups[0].TemplateID, groups)
	}
	// EnsureMapping must be called exactly once per AnalyzeLogs call.
	if store.ensureCalls != 1 {
		t.Errorf("EnsureMapping calls = %d, want 1", store.ensureCalls)
	}
}

func TestAnalyzeLogs_EmptyBaselineFallback_AppliesPrefixBias(t *testing.T) {
	// Empty baseline + non-empty prefix-filtered current → analyzer
	// uses the prefix counts so hovering metric_name surfaces relevant
	// templates instead of whatever dominates the stream.
	store := &fakeStore{
		baseline: map[string]uint64{},
		current: map[string]uint64{
			"cpu_template":  10,
			"user_template": 100, // dominates raw counts
		},
		prefixCounts: map[string]uint64{
			"cpu_template": 10,
		},
		representative: map[string][]string{
			"cpu_template":  {"cpu_usage: 67%"},
			"user_template": {"User foo logged in"},
		},
	}
	la := NewLogAnalyzerWithStore(store)
	groups, err := la.AnalyzeLogs(
		context.Background(),
		"org-1", "CPU Usage", "CPU Usage", "cpu_usage",
		time.Now().Add(-1*time.Minute), time.Now(),
	)
	if err != nil {
		t.Fatalf("AnalyzeLogs: %v", err)
	}
	if len(groups) != 1 {
		t.Fatalf("expected 1 group via prefix-bias (cpu_template), got %d", len(groups))
	}
	if groups[0].TemplateID != "cpu_template" {
		t.Errorf("expected cpu_template selected by prefix bias, got %q", groups[0].TemplateID)
	}
}

func TestAnalyzeLogs_BothEmpty_ReturnsEmpty(t *testing.T) {
	store := &fakeStore{
		baseline: map[string]uint64{},
		current:  map[string]uint64{},
	}
	la := NewLogAnalyzerWithStore(store)
	groups, err := la.AnalyzeLogs(
		context.Background(),
		"o", "d", "p", "m",
		time.Now().Add(-time.Minute), time.Now(),
	)
	if err != nil {
		t.Fatalf("AnalyzeLogs: %v", err)
	}
	if len(groups) != 0 {
		t.Errorf("expected empty result for empty windows, got %d groups", len(groups))
	}
}

func TestAnalyzeLogs_BaselineFetchError(t *testing.T) {
	store := &fakeStore{
		templateCountsErr: errors.New("clickhouse unreachable"),
	}
	la := NewLogAnalyzerWithStore(store)
	_, err := la.AnalyzeLogs(
		context.Background(), "o", "d", "p", "m",
		time.Now().Add(-time.Minute), time.Now(),
	)
	if err == nil {
		t.Fatal("expected error to propagate from GetTemplateCounts")
	}
}

func TestAnalyzeLogs_RepresentativeFetchError(t *testing.T) {
	store := &fakeStore{
		baseline:          map[string]uint64{"a": 5},
		current:           map[string]uint64{"a": 5, "b": 100}, // b only in current
		representativeErr: errors.New("kaboom"),
	}
	la := NewLogAnalyzerWithStore(store)
	_, err := la.AnalyzeLogs(
		context.Background(), "o", "d", "p", "m",
		time.Now().Add(-time.Minute), time.Now(),
	)
	if err == nil {
		t.Fatal("expected error to propagate from GetRepresentativeLogs")
	}
}

func TestAnalyzeLogs_EnsureMappingErrorIsLoggedNotFatal(t *testing.T) {
	// EnsureMapping failures must not fail the analysis — auto-
	// registration is best-effort.
	store := &fakeStore{
		baseline: map[string]uint64{"a": 5},
		current:  map[string]uint64{"a": 6, "b": 10},
		representative: map[string][]string{
			"a": {"log a"},
			"b": {"log b"},
		},
		ensureMappingErr: errors.New("insert failed"),
	}
	la := NewLogAnalyzerWithStore(store)
	groups, err := la.AnalyzeLogs(
		context.Background(), "o", "d", "p", "m",
		time.Now().Add(-time.Minute), time.Now(),
	)
	if err != nil {
		t.Fatalf("AnalyzeLogs should not fail when EnsureMapping errors: %v", err)
	}
	if len(groups) == 0 {
		t.Fatal("expected groups even when EnsureMapping fails")
	}
}

func TestAnalyzeLogs_Close(t *testing.T) {
	store := &fakeStore{}
	la := NewLogAnalyzerWithStore(store)
	if err := la.Close(); err != nil {
		t.Errorf("Close: %v", err)
	}
}

func TestAnalyzeLogs_VerifyTables(t *testing.T) {
	store := &fakeStore{}
	la := NewLogAnalyzerWithStore(store)
	if err := la.VerifyTables(); err != nil {
		t.Errorf("VerifyTables: %v", err)
	}
}

func TestAnalyzeRawLogs_StubPath(t *testing.T) {
	// AnalyzeRawLogs is a stub for the standalone-deployment path that
	// hasn't been wired up yet. Pin the contract: returns one group with
	// the placeholder strings + zero divergence.
	la := NewLogAnalyzerWithStore(&fakeStore{})
	groups, err := la.AnalyzeRawLogs(
		context.Background(),
		"loki", "o", "d", "p", "m",
		time.Now().Add(-time.Minute), time.Now(),
	)
	if err != nil {
		t.Fatalf("AnalyzeRawLogs: %v", err)
	}
	if len(groups) != 1 {
		t.Fatalf("expected 1 stub group, got %d", len(groups))
	}
	if groups[0].TemplateID != "stub" {
		t.Errorf("stub template id = %q, want %q", groups[0].TemplateID, "stub")
	}
	// Backend echoed back through the placeholder strings.
	found := false
	for _, line := range groups[0].RepresentativeLogs {
		if line == "Backend: loki" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected backend name in representative logs, got: %v", groups[0].RepresentativeLogs)
	}
}

func TestAnalyzeRawLogsWithExtractor_StubPath(t *testing.T) {
	la := NewLogAnalyzerWithStore(&fakeStore{})
	groups, err := la.AnalyzeRawLogsWithExtractor(
		context.Background(),
		"clickhouse", "o", "d", "p", "m",
		time.Now().Add(-time.Minute), time.Now(),
		nil, nil,
	)
	if err != nil {
		t.Fatalf("AnalyzeRawLogsWithExtractor: %v", err)
	}
	if len(groups) != 1 || groups[0].TemplateID != "stub" {
		t.Errorf("expected one stub group, got %+v", groups)
	}
}
