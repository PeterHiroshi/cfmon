package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestCompletionCmd_Bash(t *testing.T) {
	resetGlobalFlags()

	// Capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Execute completion bash
	rootCmd.SetArgs([]string{"completion", "bash"})

	err := rootCmd.Execute()

	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("completion bash error = %v, want nil", err)
	}

	// Read captured output
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Should contain bash completion script markers
	if !strings.Contains(output, "bash") && !strings.Contains(output, "complete") {
		t.Errorf("bash completion output doesn't look like a bash script: %q", output[:minInt(100, len(output))])
	}
}

func TestCompletionCmd_Zsh(t *testing.T) {
	resetGlobalFlags()

	// Capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Execute completion zsh
	rootCmd.SetArgs([]string{"completion", "zsh"})

	err := rootCmd.Execute()

	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("completion zsh error = %v, want nil", err)
	}

	// Read captured output
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Should contain zsh completion script markers
	if !strings.Contains(output, "compdef") && !strings.Contains(output, "zsh") {
		t.Errorf("zsh completion output doesn't look like a zsh script")
	}
}

func TestCompletionCmd_Fish(t *testing.T) {
	resetGlobalFlags()

	// Capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Execute completion fish
	rootCmd.SetArgs([]string{"completion", "fish"})

	err := rootCmd.Execute()

	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("completion fish error = %v, want nil", err)
	}

	// Read captured output
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Should contain fish completion script markers
	if !strings.Contains(output, "complete") && !strings.Contains(output, "fish") {
		t.Errorf("fish completion output doesn't look like a fish script")
	}
}

func TestCompletionCmd_PowerShell(t *testing.T) {
	resetGlobalFlags()

	// Capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Execute completion powershell
	rootCmd.SetArgs([]string{"completion", "powershell"})

	err := rootCmd.Execute()

	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("completion powershell error = %v, want nil", err)
	}

	// Read captured output
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Should contain PowerShell script markers
	if !strings.Contains(output, "Register-ArgumentCompleter") && !strings.Contains(output, "PowerShell") {
		t.Errorf("powershell completion output doesn't look like a PowerShell script")
	}
}

func TestCompletionCmd_MissingShellArg(t *testing.T) {
	resetGlobalFlags()

	// Execute completion without shell argument
	rootCmd.SetArgs([]string{"completion"})

	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("completion without shell arg: error = nil, want error")
	}
}

func TestCompletionCmd_InvalidShell(t *testing.T) {
	resetGlobalFlags()

	// Execute completion with invalid shell
	rootCmd.SetArgs([]string{"completion", "invalid-shell"})

	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("completion with invalid shell: error = nil, want error")
	}

	if !strings.Contains(err.Error(), "invalid") && !strings.Contains(err.Error(), "arg") {
		t.Errorf("error message = %q, should mention invalid argument", err.Error())
	}
}

func TestCompletionCmd_TooManyArgs(t *testing.T) {
	resetGlobalFlags()

	// Execute completion with too many arguments
	rootCmd.SetArgs([]string{"completion", "bash", "extra"})

	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("completion with too many args: error = nil, want error")
	}
}

func TestCompletionCmd_ValidArgs(t *testing.T) {
	resetGlobalFlags()

	// Get the completion command
	var completionCmd *cobra.Command
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "completion" {
			completionCmd = cmd
			break
		}
	}

	if completionCmd == nil {
		t.Fatal("completion command not found")
	}

	// Check that ValidArgs contains the expected shells
	expectedShells := []string{"bash", "zsh", "fish", "powershell"}
	for _, shell := range expectedShells {
		found := false
		for _, validArg := range completionCmd.ValidArgs {
			if validArg == shell {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("shell %q not found in ValidArgs: %v", shell, completionCmd.ValidArgs)
		}
	}
}

func TestCompletionCmd_CommandRegistered(t *testing.T) {
	resetGlobalFlags()

	// Verify completion command is registered with root
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "completion" {
			found = true
			break
		}
	}

	if !found {
		t.Error("completion command not found in root commands")
	}
}

func TestCompletionCmd_ShortDescription(t *testing.T) {
	resetGlobalFlags()

	// Verify completion command has a short description
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "completion" {
			if cmd.Short == "" {
				t.Error("completion command has empty short description")
			}
			if !strings.Contains(strings.ToLower(cmd.Short), "completion") && !strings.Contains(strings.ToLower(cmd.Short), "shell") {
				t.Errorf("completion command short description doesn't mention completion/shell: %q", cmd.Short)
			}
			return
		}
	}

	t.Fatal("completion command not found")
}

// Helper function for minInt (Go 1.20 doesn't have built-in min)
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
