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

func TestLoad_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "invalid.yaml")

	// Write invalid YAML
	err := os.WriteFile(configPath, []byte("invalid: yaml: content: ["), 0600)
	if err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	// Load should return an error
	_, err = Load(configPath)
	if err == nil {
		t.Fatal("Load() with invalid YAML: error = nil, want error")
	}
}

func TestLoad_PermissionError(t *testing.T) {
	// Skip this test if running as root (which can read files with any permissions)
	if os.Geteuid() == 0 {
		t.Skip("Skipping permission test when running as root")
	}

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "unreadable.yaml")

	// Create a file with no read permissions
	cfg := &Config{Token: "test"}
	err := Save(configPath, cfg)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Remove read permissions
	err = os.Chmod(configPath, 0000)
	if err != nil {
		t.Fatalf("Chmod() error = %v", err)
	}

	// Restore permissions after test
	defer os.Chmod(configPath, 0600)

	// Load should return an error
	_, err = Load(configPath)
	if err == nil {
		t.Fatal("Load() with unreadable file: error = nil, want error")
	}
}

func TestSave_InvalidPath(t *testing.T) {
	// Try to save to a path that's a file instead of a directory
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "file.txt")

	// Create a file
	err := os.WriteFile(filePath, []byte("test"), 0600)
	if err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	// Try to save config to file/config.yaml (invalid path)
	invalidPath := filepath.Join(filePath, "config.yaml")
	cfg := &Config{Token: "test"}

	err = Save(invalidPath, cfg)
	if err == nil {
		t.Fatal("Save() with invalid path: error = nil, want error")
	}
}

func TestSave_FilePermissions(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	cfg := &Config{Token: "secret-token"}
	err := Save(configPath, cfg)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Check file permissions
	info, err := os.Stat(configPath)
	if err != nil {
		t.Fatalf("Stat() error = %v", err)
	}

	mode := info.Mode().Perm()
	if mode != 0600 {
		t.Errorf("file permissions = %o, want 0600", mode)
	}
}
