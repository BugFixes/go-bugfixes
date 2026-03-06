package bugfixes_test

import (
	"testing"

	bugfixes "github.com/bugfixes/go-bugfixes"
)

func TestConfigMerge(t *testing.T) {
	base := bugfixes.Config{
		Server:      "https://base.example",
		AgentKey:    "base-key",
		AgentSecret: "base-secret",
		LogLevel:    "warn",
	}

	merged := base.Merge(bugfixes.Config{
		AgentKey:  "override-key",
		LocalOnly: true,
	})

	if merged.Server != "https://base.example" {
		t.Fatalf("expected base server, got %q", merged.Server)
	}
	if merged.AgentKey != "override-key" {
		t.Fatalf("expected override key, got %q", merged.AgentKey)
	}
	if merged.AgentSecret != "base-secret" {
		t.Fatalf("expected base secret, got %q", merged.AgentSecret)
	}
	if merged.LogLevel != "warn" {
		t.Fatalf("expected base log level, got %q", merged.LogLevel)
	}
	if !merged.LocalOnly {
		t.Fatal("expected LocalOnly to be true")
	}
}

func TestConfigMerge_ResetLocalOnly(t *testing.T) {
	base := bugfixes.Config{
		LocalOnly: true,
	}

	merged := base.Merge(bugfixes.Config{
		LocalOnly: false,
	})

	if merged.LocalOnly {
		t.Fatal("expected LocalOnly to be reset to false by override")
	}
}

func TestSetDefaultConfig(t *testing.T) {
	t.Cleanup(bugfixes.ResetDefaultConfig)

	bugfixes.SetDefaultConfig(bugfixes.Config{
		Server:      "https://config.example",
		AgentKey:    "key",
		AgentSecret: "secret",
	})

	cfg := bugfixes.GetDefaultConfig()
	if cfg.Server != "https://config.example" {
		t.Fatalf("expected configured server, got %q", cfg.Server)
	}
	if cfg.AgentKey != "key" {
		t.Fatalf("expected configured key, got %q", cfg.AgentKey)
	}
	if cfg.AgentSecret != "secret" {
		t.Fatalf("expected configured secret, got %q", cfg.AgentSecret)
	}
}
