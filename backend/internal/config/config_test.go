package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeConfig(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}
	return path
}

func TestDefaultsWhenFieldsMissing(t *testing.T) {
	path := writeConfig(t, "http:\n  port: 9000\n")

	conf, err := Init(path)
	if err != nil {
		t.Fatalf("failed to init config: %v", err)
	}

	if conf.LogLevel != "info" {
		t.Errorf("LogLevel = %q, want %q", conf.LogLevel, "info")
	}
	if conf.Http.Port != "9000" {
		t.Errorf("Http.Port = %q, want %q", conf.Http.Port, "9000")
	}
	if conf.Postgres.MaxOpenConns != 25 {
		t.Errorf("Postgres.MaxOpenConns = %d, want 25", conf.Postgres.MaxOpenConns)
	}
	if conf.Postgres.MaxIdleConns != 10 {
		t.Errorf("Postgres.MaxIdleConns = %d, want 10", conf.Postgres.MaxIdleConns)
	}
	if conf.Postgres.ConnMaxLifetime != 5*time.Minute {
		t.Errorf("Postgres.ConnMaxLifetime = %v, want 5m", conf.Postgres.ConnMaxLifetime)
	}
	if conf.Postgres.ConnMaxIdleTime != 5*time.Minute {
		t.Errorf("Postgres.ConnMaxIdleTime = %v, want 5m", conf.Postgres.ConnMaxIdleTime)
	}
	if len(conf.Http.TrustedProxies) != 0 {
		t.Errorf("Http.TrustedProxies = %v, want empty", conf.Http.TrustedProxies)
	}
}

func TestParsesYamlValues(t *testing.T) {
	content := `log_level: debug
http:
  port: 9001
  trusted_proxies:
    - 192.168.1.0/24
    - 10.0.0.0/8
postgres:
  max_open_conns: 50
  max_idle_conns: 20
  conn_max_lifetime: 10m
  conn_max_idle_time: 1m
`
	path := writeConfig(t, content)

	conf, err := Init(path)
	if err != nil {
		t.Fatalf("failed to init config: %v", err)
	}

	if conf.LogLevel != "debug" {
		t.Errorf("LogLevel = %q, want %q", conf.LogLevel, "debug")
	}
	if conf.Http.Port != "9001" {
		t.Errorf("Http.Port = %q, want %q", conf.Http.Port, "9001")
	}
	wantProxies := []string{"192.168.1.0/24", "10.0.0.0/8"}
	if len(conf.Http.TrustedProxies) != len(wantProxies) {
		t.Fatalf("TrustedProxies = %v, want %v", conf.Http.TrustedProxies, wantProxies)
	}
	for i, p := range wantProxies {
		if conf.Http.TrustedProxies[i] != p {
			t.Errorf("TrustedProxies[%d] = %q, want %q", i, conf.Http.TrustedProxies[i], p)
		}
	}
	if conf.Postgres.MaxOpenConns != 50 {
		t.Errorf("MaxOpenConns = %d, want 50", conf.Postgres.MaxOpenConns)
	}
	if conf.Postgres.MaxIdleConns != 20 {
		t.Errorf("MaxIdleConns = %d, want 20", conf.Postgres.MaxIdleConns)
	}
	if conf.Postgres.ConnMaxLifetime != 10*time.Minute {
		t.Errorf("ConnMaxLifetime = %v, want 10m", conf.Postgres.ConnMaxLifetime)
	}
	if conf.Postgres.ConnMaxIdleTime != time.Minute {
		t.Errorf("ConnMaxIdleTime = %v, want 1m", conf.Postgres.ConnMaxIdleTime)
	}
}

func TestEnvOverridesYaml(t *testing.T) {
	t.Setenv("POSTGRES_MAX_OPEN_CONNS", "100")
	t.Setenv("HOST", "10.0.0.5")

	path := writeConfig(t, "http:\n  port: 9000\n")

	conf, err := Init(path)
	if err != nil {
		t.Fatalf("failed to init config: %v", err)
	}

	if conf.Postgres.MaxOpenConns != 100 {
		t.Errorf("MaxOpenConns = %d, want 100", conf.Postgres.MaxOpenConns)
	}
	if conf.Http.Host != "10.0.0.5" {
		t.Errorf("Http.Host = %q, want %q", conf.Http.Host, "10.0.0.5")
	}
}
