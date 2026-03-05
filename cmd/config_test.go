package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/PeterHiroshi/cfmon/internal/config"
)

func TestMaskToken(t *testing.T) {
	tests := []struct {
		name     string
		token    string
		expected string
	}{
		{
			name:     "short token",
			token:    "short",
			expected: "****",
		},
		{
			name:     "medium token",
			token:    "12345678",
			expected: "****",
		},
		{
			name:     "long token",
			token:    "1234567890abcdef",
			expected: "1234********cdef",
		},
		{
			name:     "very long token",
			token:    "abcdefghijklmnopqrstuvwxyz123456",
			expected: "abcd************************3456",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := maskToken(tt.token)
			if result != tt.expected {
				t.Errorf("maskToken(%q) = %q, want %q", tt.token, result, tt.expected)
			}
		})
	}
}

func TestConfigPathCommand(t *testing.T) {
	// Create temp config
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	// Set environment variable
	os.Setenv("CFMON_CONFIG", configPath)
	defer os.Unsetenv("CFMON_CONFIG")

	// Test that GetConfigPath returns the right path
	cfg := config.New()
	gotPath := cfg.GetConfigPath()

	if gotPath != configPath {
		t.Errorf("GetConfigPath() = %q, want %q", gotPath, configPath)
	}
}

func TestConfigShowCommand(t *testing.T) {
	// Setup temp config
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	// Create config file
	configContent := `token: test-token-1234567890
api_endpoint: https://api.test.com
default_format: json`

	err := os.WriteFile(configPath, []byte(configContent), 0600)
	if err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}

	// Set environment
	os.Setenv("CFMON_CONFIG", configPath)
	defer os.Unsetenv("CFMON_CONFIG")

	// Load and verify config
	cfg := config.New()
	err = cfg.Load()
	if err != nil {
		t.Errorf("Failed to load config: %v", err)
	}

	if cfg.Token != "test-token-1234567890" {
		t.Errorf("Token = %q, want %q", cfg.Token, "test-token-1234567890")
	}

	if cfg.APIEndpoint != "https://api.test.com" {
		t.Errorf("APIEndpoint = %q, want %q", cfg.APIEndpoint, "https://api.test.com")
	}

	if cfg.DefaultFormat != "json" {
		t.Errorf("DefaultFormat = %q, want %q", cfg.DefaultFormat, "json")
	}
}