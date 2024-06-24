package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	// Create a temporary config file
	content := []byte(`
targets:
  - name: "Test1"
    url: "https://test1.com"
    interval: "10s"
    timeout: "5s"
  - name: "Test2"
    url: "https://test2.com"
    interval: "1m"
    timeout: "3s"
server:
  port: 9090
`)
	tmpfile, err := os.CreateTemp("", "config.*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Test loading the config
	cfg, err := Load(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Check if the config is parsed correctly
	if len(cfg.Targets) != 2 {
		t.Errorf("Expected 2 targets, got %d", len(cfg.Targets))
	}

	if cfg.Targets[0].Name != "Test1" {
		t.Errorf("Expected first target name to be 'Test1', got '%s'", cfg.Targets[0].Name)
	}

	if cfg.Targets[0].Interval != 10*time.Second {
		t.Errorf("Expected first target interval to be 10s, got %v", cfg.Targets[0].Interval)
	}

	if cfg.Server.Port != 9090 {
		t.Errorf("Expected server port to be 9090, got %d", cfg.Server.Port)
	}
}
