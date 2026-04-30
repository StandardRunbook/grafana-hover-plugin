package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestServerConfigGetAddress(t *testing.T) {
	s := ServerConfig{Host: "0.0.0.0", Port: 8080}
	if got := s.GetAddress(); got != "0.0.0.0:8080" {
		t.Errorf("GetAddress() = %q, want %q", got, "0.0.0.0:8080")
	}
}

func TestLoadFromMissingConfigUsesDefaults(t *testing.T) {
	// Run from a temp dir where there's no config.toml in `.`, `..`, or
	// the absolute fallback path. Load() should silently use defaults
	// (the function logs "config: no file found in search paths").
	tmpDir := t.TempDir()
	cwd, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer os.Chdir(cwd)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error on no-config path: %v", err)
	}
	// Defaults baked into Load()
	if cfg.Server.Host != "127.0.0.1" {
		t.Errorf("default host = %q, want 127.0.0.1", cfg.Server.Host)
	}
	if cfg.Server.Port != 8080 {
		t.Errorf("default port = %d, want 8080", cfg.Server.Port)
	}
	if cfg.ClickHouse.URL != "localhost:9000" {
		t.Errorf("default clickhouse url = %q, want localhost:9000", cfg.ClickHouse.URL)
	}
	if cfg.ClickHouse.Database != "default" {
		t.Errorf("default database = %q, want default", cfg.ClickHouse.Database)
	}
}

func TestLoadReadsConfigToml(t *testing.T) {
	// Drop a config.toml in the cwd and verify Load picks it up.
	tmpDir := t.TempDir()
	cwd, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer os.Chdir(cwd)

	body := `[server]
host = "10.0.0.1"
port = 9999

[clickhouse]
url = "ch.internal:9000"
database = "metrics"
user = "reader"
password = "shh"
`
	if err := os.WriteFile(filepath.Join(tmpDir, "config.toml"), []byte(body), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}
	if cfg.Server.Host != "10.0.0.1" {
		t.Errorf("host = %q, want 10.0.0.1", cfg.Server.Host)
	}
	if cfg.Server.Port != 9999 {
		t.Errorf("port = %d, want 9999", cfg.Server.Port)
	}
	if cfg.ClickHouse.URL != "ch.internal:9000" {
		t.Errorf("ch url = %q, want ch.internal:9000", cfg.ClickHouse.URL)
	}
	if cfg.ClickHouse.User != "reader" {
		t.Errorf("ch user = %q, want reader", cfg.ClickHouse.User)
	}
}

func TestServerConfigGetAddress_FormatPattern(t *testing.T) {
	cases := []struct {
		host string
		port int
		want string
	}{
		{"localhost", 8080, "localhost:8080"},
		{"0.0.0.0", 0, "0.0.0.0:0"},
		{"::1", 8080, "::1:8080"}, // raw IPv6 — caller's responsibility to wrap in []
	}
	for _, c := range cases {
		s := ServerConfig{Host: c.host, Port: c.port}
		if got := s.GetAddress(); got != c.want {
			t.Errorf("GetAddress(%q, %d) = %q, want %q", c.host, c.port, got, c.want)
		}
	}
}
