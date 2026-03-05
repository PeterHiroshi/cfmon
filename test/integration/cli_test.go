//go:build integration
// +build integration

package integration

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestCLIFullFlow tests the complete cfmon CLI workflow
func TestCLIFullFlow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Build the binary first
	binPath := filepath.Join(t.TempDir(), "cfmon")
	buildCmd := exec.Command("go", "build", "-o", binPath, "../..")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build cfmon: %v", err)
	}

	// Create temp config directory
	configDir := t.TempDir()
	configPath := filepath.Join(configDir, "config.yaml")

	// Set environment for tests
	env := append(os.Environ(),
		fmt.Sprintf("CFMON_CONFIG=%s", configPath),
	)

	tests := []struct {
		name       string
		args       []string
		wantOutput []string
		wantErr    bool
	}{
		{
			name:       "version command",
			args:       []string{"--version"},
			wantOutput: []string{"cfmon version"},
			wantErr:    false,
		},
		{
			name:       "help command",
			args:       []string{"help"},
			wantOutput: []string{"cfmon - Cloudflare Workers/Containers CLI Monitor", "USAGE", "CORE COMMANDS"},
			wantErr:    false,
		},
		{
			name:       "doctor command without token",
			args:       []string{"doctor"},
			wantOutput: []string{"System Diagnostics", "Go Runtime", "Operating System"},
			wantErr:    false, // Doctor doesn't fail, just reports issues
		},
		{
			name:       "config path command",
			args:       []string{"config", "path"},
			wantOutput: []string{configPath},
			wantErr:    false,
		},
		{
			name:       "config show without config",
			args:       []string{"config", "show"},
			wantOutput: []string{"Configuration", "not found"},
			wantErr:    false,
		},
		{
			name:       "containers list without token",
			args:       []string{"containers", "list", "test-account"},
			wantOutput: []string{"no API token"},
			wantErr:    true,
		},
		{
			name:       "workers list without token",
			args:       []string{"workers", "list", "test-account"},
			wantOutput: []string{"no API token"},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binPath, tt.args...)
			cmd.Env = env

			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()

			// Check error expectation
			if (err != nil) != tt.wantErr {
				t.Errorf("Command error = %v, wantErr %v\nStderr: %s", err, tt.wantErr, stderr.String())
			}

			// Check output contains expected strings
			output := stdout.String() + stderr.String()
			for _, want := range tt.wantOutput {
				if !strings.Contains(output, want) {
					t.Errorf("Output missing %q\nGot: %s", want, output)
				}
			}
		})
	}
}

// TestCLIWithMockToken tests commands that require a token
func TestCLIWithMockToken(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Build the binary
	binPath := filepath.Join(t.TempDir(), "cfmon")
	buildCmd := exec.Command("go", "build", "-o", binPath, "../..")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build cfmon: %v", err)
	}

	// Create config with mock token
	configDir := t.TempDir()
	configPath := filepath.Join(configDir, "config.yaml")

	// Write mock config
	configContent := `token: mock-test-token-1234567890
api_endpoint: https://api.cloudflare.com/client/v4
default_format: table`

	if err := os.WriteFile(configPath, []byte(configContent), 0600); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	env := append(os.Environ(),
		fmt.Sprintf("CFMON_CONFIG=%s", configPath),
	)

	tests := []struct {
		name       string
		args       []string
		wantOutput []string
	}{
		{
			name:       "config show with token",
			args:       []string{"config", "show"},
			wantOutput: []string{"API Token:", "mock-test-token"},
		},
		{
			name:       "doctor with token",
			args:       []string{"doctor"},
			wantOutput: []string{"API Token", "Configured"},
		},
		{
			name:       "containers list with sorting",
			args:       []string{"containers", "list", "test-account", "--sort", "cpu"},
			wantOutput: []string{"listing containers"}, // Will fail at API call but shows it parses flags
		},
		{
			name:       "workers list with filter",
			args:       []string{"workers", "list", "test-account", "--filter", "api", "--limit", "5"},
			wantOutput: []string{"listing workers"}, // Will fail at API call but shows it parses flags
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binPath, tt.args...)
			cmd.Env = env

			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			// We expect some commands to fail (API calls) but we're testing parsing
			_ = cmd.Run()

			output := stdout.String() + stderr.String()
			for _, want := range tt.wantOutput {
				if !strings.Contains(output, want) {
					t.Errorf("Output missing %q\nGot: %s", want, output)
				}
			}
		})
	}
}

// TestCLICompletions tests shell completion generation
func TestCLICompletions(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	binPath := filepath.Join(t.TempDir(), "cfmon")
	buildCmd := exec.Command("go", "build", "-o", binPath, "../..")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build cfmon: %v", err)
	}

	shells := []string{"bash", "zsh", "fish", "powershell"}

	for _, shell := range shells {
		t.Run(fmt.Sprintf("completion_%s", shell), func(t *testing.T) {
			cmd := exec.Command(binPath, "completion", shell)

			var stdout bytes.Buffer
			cmd.Stdout = &stdout

			if err := cmd.Run(); err != nil {
				t.Errorf("Failed to generate %s completion: %v", shell, err)
			}

			output := stdout.String()
			if len(output) < 100 {
				t.Errorf("%s completion output too short: %d bytes", shell, len(output))
			}

			// Check for shell-specific patterns
			switch shell {
			case "bash":
				if !strings.Contains(output, "complete") {
					t.Errorf("Bash completion missing 'complete' command")
				}
			case "zsh":
				if !strings.Contains(output, "compdef") {
					t.Errorf("Zsh completion missing 'compdef'")
				}
			case "fish":
				if !strings.Contains(output, "complete -c cfmon") {
					t.Errorf("Fish completion missing 'complete -c cfmon'")
				}
			case "powershell":
				if !strings.Contains(output, "Register-ArgumentCompleter") {
					t.Errorf("PowerShell completion missing 'Register-ArgumentCompleter'")
				}
			}
		})
	}
}

// TestCLITimeout tests the timeout flag
func TestCLITimeout(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	binPath := filepath.Join(t.TempDir(), "cfmon")
	buildCmd := exec.Command("go", "build", "-o", binPath, "../..")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build cfmon: %v", err)
	}

	// Test various timeout values
	timeouts := []string{"1s", "30s", "1m", "100ms"}

	for _, timeout := range timeouts {
		t.Run(fmt.Sprintf("timeout_%s", timeout), func(t *testing.T) {
			cmd := exec.Command(binPath, "doctor", "--timeout", timeout, "-v")

			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			// Set a deadline for the command
			done := make(chan error, 1)
			go func() {
				done <- cmd.Run()
			}()

			select {
			case <-done:
				// Command completed
				output := stderr.String()
				if strings.Contains(output, "Debug:") && strings.Contains(output, "timeout") {
					// Verbose mode should show timeout info
					t.Logf("Timeout %s parsed correctly", timeout)
				}
			case <-time.After(5 * time.Second):
				cmd.Process.Kill()
				t.Errorf("Command timed out with --timeout %s", timeout)
			}
		})
	}
}

// TestCLIVerboseMode tests verbose output
func TestCLIVerboseMode(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	binPath := filepath.Join(t.TempDir(), "cfmon")
	buildCmd := exec.Command("go", "build", "-o", binPath, "../..")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build cfmon: %v", err)
	}

	// Create config with token
	configDir := t.TempDir()
	configPath := filepath.Join(configDir, "config.yaml")
	os.WriteFile(configPath, []byte("token: test-token"), 0600)

	env := append(os.Environ(),
		fmt.Sprintf("CFMON_CONFIG=%s", configPath),
	)

	tests := []struct {
		name       string
		args       []string
		wantDebug  bool
	}{
		{
			name:      "doctor with verbose",
			args:      []string{"doctor", "-v"},
			wantDebug: true,
		},
		{
			name:      "doctor without verbose",
			args:      []string{"doctor"},
			wantDebug: false,
		},
		{
			name:      "containers list with verbose",
			args:      []string{"containers", "list", "test", "--verbose"},
			wantDebug: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binPath, tt.args...)
			cmd.Env = env

			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			_ = cmd.Run() // Might fail due to API calls, but we're checking output

			output := stderr.String()
			hasDebug := strings.Contains(output, "Debug:")

			if hasDebug != tt.wantDebug {
				t.Errorf("Verbose output: got debug=%v, want debug=%v\nStderr: %s",
					hasDebug, tt.wantDebug, output)
			}
		})
	}
}