package clickhouse

import (
	"errors"
	"testing"
)

// These small helpers exist because the original author avoided pulling
// in `strings` for some reason; they're now load-bearing for the
// schema-loader and the createTables error path. Pinning their behavior
// so a future strings-package migration is visibly equivalent.

func TestContains(t *testing.T) {
	cases := []struct {
		s, sub string
		want   bool
	}{
		{"", "", true},
		{"x", "", true},
		{"", "x", false},
		{"hello", "hello", true},
		{"hello world", "hello", true}, // prefix
		{"hello world", "world", true}, // suffix
		{"abc xyz def", "xyz", true},   // middle
		{"abc", "xyz", false},
		{"short", "longer-string", false},
	}
	for _, c := range cases {
		if got := contains(c.s, c.sub); got != c.want {
			t.Errorf("contains(%q, %q) = %v, want %v", c.s, c.sub, got, c.want)
		}
	}
}

func TestContainsMiddle(t *testing.T) {
	if !containsMiddle("abcdef", "cde") {
		t.Error("containsMiddle should find cde in abcdef")
	}
	if containsMiddle("abcdef", "xyz") {
		t.Error("containsMiddle should not find xyz in abcdef")
	}
}

func TestContainsError(t *testing.T) {
	if containsError(nil, "anything") {
		t.Error("nil err should never contain anything")
	}
	if !containsError(errors.New("table 'logs' UNKNOWN_TABLE error"), "UNKNOWN_TABLE") {
		t.Error("should detect UNKNOWN_TABLE substring")
	}
	if containsError(errors.New("disk full"), "UNKNOWN_TABLE") {
		t.Error("should not match unrelated error")
	}
}

func TestHasPrefix(t *testing.T) {
	if !hasPrefix("http://localhost", "http://") {
		t.Error("hasPrefix(http://) should match")
	}
	if hasPrefix("localhost", "http://") {
		t.Error("hasPrefix should reject non-matching prefix")
	}
	if !hasPrefix("abc", "") {
		t.Error("empty prefix should always match")
	}
	if hasPrefix("ab", "abc") {
		t.Error("prefix longer than string should not match")
	}
}

func TestHasSuffix(t *testing.T) {
	if !hasSuffix("hello;", ";") {
		t.Error("hasSuffix(;) should match")
	}
	if hasSuffix("hello", ";") {
		t.Error("hasSuffix should reject non-matching suffix")
	}
	if hasSuffix("ab", "abc") {
		t.Error("suffix longer than string should not match")
	}
}

func TestTrimSpace(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"", ""},
		{"   ", ""},
		{"\t\nhello\r ", "hello"},
		{"hello", "hello"},
		{"  hello  world  ", "hello  world"}, // only edges trimmed
	}
	for _, c := range cases {
		if got := trimSpace(c.in); got != c.want {
			t.Errorf("trimSpace(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestSplitLines(t *testing.T) {
	got := splitLines("a\nb\nc")
	want := []string{"a", "b", "c"}
	if len(got) != len(want) {
		t.Fatalf("expected %d lines, got %d: %v", len(want), len(got), got)
	}
	for i, line := range got {
		if line != want[i] {
			t.Errorf("line %d: got %q, want %q", i, line, want[i])
		}
	}
}

func TestSplitLinesNoTrailingNewline(t *testing.T) {
	got := splitLines("a\nb")
	if len(got) != 2 || got[0] != "a" || got[1] != "b" {
		t.Errorf("expected [a, b], got %v", got)
	}
}

func TestSplitLinesEmptyString(t *testing.T) {
	if got := splitLines(""); len(got) != 0 {
		t.Errorf("empty input should produce no lines, got %v", got)
	}
}

func TestSplitSQLStripsCommentsAndSplitsOnSemicolons(t *testing.T) {
	schema := `
-- comment line
CREATE TABLE foo (id Int32);

CREATE TABLE bar (
  id Int32,
  name String
);
`
	stmts := splitSQL(schema)
	if len(stmts) != 2 {
		t.Fatalf("expected 2 statements, got %d: %v", len(stmts), stmts)
	}
	if !contains(stmts[0], "CREATE TABLE foo") {
		t.Errorf("first stmt missing CREATE TABLE foo: %q", stmts[0])
	}
	if !contains(stmts[1], "CREATE TABLE bar") {
		t.Errorf("second stmt missing CREATE TABLE bar: %q", stmts[1])
	}
}

func TestSplitSQLAppendsTrailingNonSemicolonStatement(t *testing.T) {
	// Final statement without a trailing semicolon should still be
	// emitted (some schema files don't end with ;).
	stmts := splitSQL("SELECT 1")
	if len(stmts) != 1 || stmts[0] != "SELECT 1" {
		t.Errorf("expected [SELECT 1], got %v", stmts)
	}
}

func TestDeterministicID_StableAndKindSensitive(t *testing.T) {
	a := deterministicID("metric", "1", "Hover Plugin Demo", "CPU Usage", "cpu_usage")
	b := deterministicID("metric", "1", "Hover Plugin Demo", "CPU Usage", "cpu_usage")
	if a != b {
		t.Errorf("deterministicID not stable: %q vs %q", a, b)
	}
	if len(a) != 24 {
		t.Errorf("expected 24-char id, got %d (%q)", len(a), a)
	}
	c := deterministicID("mapping", "1", "Hover Plugin Demo", "CPU Usage", "cpu_usage")
	if a == c {
		t.Errorf("kind change should produce different id, got same %q", a)
	}
}

func TestDeterministicID_DelimiterConfusion(t *testing.T) {
	// "a|b"+"c" must not collide with "a"+"b|c" — the NUL separator
	// inside deterministicID is what guarantees this.
	x := deterministicID("k", "a|b", "c", "x", "x")
	y := deterministicID("k", "a", "b|c", "x", "x")
	if x == y {
		t.Errorf("delimiter confusion: %q == %q", x, y)
	}
}
