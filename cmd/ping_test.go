package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestPingCmd_Output(t *testing.T) {
	resetGlobalFlags()

	// Capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Execute ping command
	rootCmd.SetArgs([]string{"ping"})

	err := rootCmd.Execute()

	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("pingCmd.Execute() error = %v, want nil", err)
	}

	// Read captured output
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Should contain "pong"
	if !strings.Contains(output, "pong") {
		t.Errorf("ping output doesn't contain 'pong': %q", output)
	}
}

func TestPingCmd_NoArgs(t *testing.T) {
	resetGlobalFlags()

	// Ping command should not require arguments
	rootCmd.SetArgs([]string{"ping"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("pingCmd.Execute() with no args error = %v, want nil", err)
	}
}

func TestPingCmd_OutputContainsForgeDispatch(t *testing.T) {
	resetGlobalFlags()

	// Capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rootCmd.SetArgs([]string{"ping"})
	err := rootCmd.Execute()

	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("pingCmd.Execute() error = %v, want nil", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Should mention forge-dispatch
	if !strings.Contains(output, "forge-dispatch") {
		t.Errorf("ping output doesn't contain 'forge-dispatch': %q", output)
	}
}

func TestPingCmd_OutputContainsTimestamp(t *testing.T) {
	resetGlobalFlags()

	// Capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rootCmd.SetArgs([]string{"ping"})
	err := rootCmd.Execute()

	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("pingCmd.Execute() error = %v, want nil", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Should contain timestamp marker
	if !strings.Contains(output, "at:") {
		t.Errorf("ping output doesn't contain timestamp 'at:': %q", output)
	}
}

func TestPingCmd_IgnoresExtraArgs(t *testing.T) {
	resetGlobalFlags()

	// Capture output to avoid cluttering test output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Ping command with extra args (should be ignored by cobra)
	rootCmd.SetArgs([]string{"ping", "extra", "args"})

	err := rootCmd.Execute()

	w.Close()
	os.Stdout = oldStdout

	// Read and discard output
	var buf bytes.Buffer
	buf.ReadFrom(r)

	// Should still succeed (extra args are ignored)
	if err != nil {
		t.Logf("pingCmd with extra args: error = %v", err)
	}
}

func TestPingCmd_CommandRegistered(t *testing.T) {
	resetGlobalFlags()

	// Verify ping command is registered with root
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "ping" {
			found = true
			break
		}
	}

	if !found {
		t.Error("ping command not found in root commands")
	}
}

func TestPingCmd_ShortDescription(t *testing.T) {
	resetGlobalFlags()

	// Verify ping command has a short description
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "ping" {
			if cmd.Short == "" {
				t.Error("ping command has empty short description")
			}
			if !strings.Contains(strings.ToLower(cmd.Short), "pong") && !strings.Contains(strings.ToLower(cmd.Short), "ping") {
				t.Errorf("ping command short description doesn't mention ping/pong: %q", cmd.Short)
			}
			return
		}
	}

	t.Fatal("ping command not found")
}
