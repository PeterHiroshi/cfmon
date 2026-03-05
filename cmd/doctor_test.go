package cmd

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestDoctorChecks(t *testing.T) {
	t.Run("checkConfigFile", func(t *testing.T) {
		// Create temp directory for config
		tempDir := t.TempDir()
		configPath := filepath.Join(tempDir, ".cfmon", "config.yaml")

		// Test when config doesn't exist
		result := checkConfigFile()
		if result.status != "warn" && result.status != "pass" {
			t.Errorf("expected warn or pass status for missing config, got %s", result.status)
		}

		// Create config file
		os.MkdirAll(filepath.Dir(configPath), 0755)
		os.WriteFile(configPath, []byte("token: test-token"), 0600)

		// Test with valid config
		os.Setenv("CFMON_CONFIG", configPath)
		defer os.Unsetenv("CFMON_CONFIG")

		result = checkConfigFile()
		if result.status != "pass" {
			t.Errorf("expected pass status for valid config, got %s: %s", result.status, result.message)
		}
	})

	t.Run("checkAPIToken", func(t *testing.T) {
		// Test with no token
		origToken := token
		token = ""
		os.Unsetenv("CFMON_TOKEN")
		defer func() { token = origToken }()

		result := checkAPIToken()
		if result.status != "fail" {
			t.Errorf("expected fail status for no token, got %s", result.status)
		}

		// Test with environment token
		os.Setenv("CFMON_TOKEN", "test-environment-token-12345678901234567890")
		defer os.Unsetenv("CFMON_TOKEN")

		result = checkAPIToken()
		if result.status != "pass" {
			t.Errorf("expected pass status for env token, got %s", result.status)
		}

		// Test with command line token
		os.Unsetenv("CFMON_TOKEN")
		token = "test-cli-token-12345678901234567890"

		result = checkAPIToken()
		if result.status != "pass" {
			t.Errorf("expected pass status for CLI token, got %s", result.status)
		}
	})

	t.Run("checkNetworkConnectivity", func(t *testing.T) {
		// Create test server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		// This test would need to mock the actual network call
		// For now, we just verify the function doesn't panic
		result := checkNetworkConnectivity()
		if result.name != "Network Connectivity" {
			t.Errorf("unexpected check name: %s", result.name)
		}
	})

	t.Run("checkFilePermissions", func(t *testing.T) {
		// Create temp directory
		tempDir := t.TempDir()
		os.Setenv("CFMON_CONFIG", filepath.Join(tempDir, "config.yaml"))
		defer os.Unsetenv("CFMON_CONFIG")

		result := checkFilePermissions()
		if result.status != "pass" {
			t.Errorf("expected pass status for writable directory, got %s: %s", result.status, result.message)
		}
	})
}

func TestDisplayCheckResults(t *testing.T) {
	checks := []checkResult{
		{name: "Test Pass", status: "pass", message: "All good"},
		{name: "Test Fail", status: "fail", message: "Something wrong"},
		{name: "Test Warn", status: "warn", message: "Warning message"},
		{name: "Test Skip", status: "skip", message: "Skipped"},
	}

	// Just verify it doesn't panic
	displayCheckResults(checks)
}

func TestMin(t *testing.T) {
	tests := []struct {
		a, b, want int
	}{
		{1, 2, 1},
		{5, 3, 3},
		{4, 4, 4},
		{-1, 0, -1},
	}

	for _, tt := range tests {
		got := min(tt.a, tt.b)
		if got != tt.want {
			t.Errorf("min(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
		}
	}
}