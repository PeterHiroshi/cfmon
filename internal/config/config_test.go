package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_EmptyConfig(t *testing.T) {
	// Create a temp directory for testing
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Load should return empty config if file doesn't exist
	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load() error = %v, want nil", err)
	}

	if cfg.Token != "" {
		t.Errorf("Token = %q, want empty string", cfg.Token)
	}
}

func TestSave_CreateConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	cfg := &Config{
		Token: "test-token-123",
	}

	err := Save(configPath, cfg)
	if err != nil {
		t.Fatalf("Save() error = %v, want nil", err)
	}

	// Verify file was created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Errorf("config file was not created")
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	original := &Config{
		Token: "my-secret-token",
	}

	// Save
	err := Save(configPath, original)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Load
	loaded, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if loaded.Token != original.Token {
		t.Errorf("Token = %q, want %q", loaded.Token, original.Token)
	}
}
