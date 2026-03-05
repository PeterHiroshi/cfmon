package cmd

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/PeterHiroshi/cfmon/internal/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(doctorCmd)
}

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check system configuration and connectivity",
	Long: `The doctor command performs a comprehensive system check to ensure cfmon is properly configured.

It checks:
  • Go runtime version
  • Configuration file status
  • API token validity
  • Cloudflare API connectivity
  • Network connectivity
  • File permissions`,
	Example: `  # Run system diagnostics
  cfmon doctor

  # Run with verbose output
  cfmon doctor -v`,
	Run: runDoctor,
}

type checkResult struct {
	name    string
	status  string
	message string
}

func runDoctor(cmd *cobra.Command, args []string) {
	fmt.Println(colorize("cfmon System Diagnostics", "cyan", true))
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println()

	checks := []checkResult{}

	// Check 1: Go Version
	goVersion := runtime.Version()
	checks = append(checks, checkResult{
		name:    "Go Runtime",
		status:  "pass",
		message: goVersion,
	})

	// Check 2: cfmon Version
	checks = append(checks, checkResult{
		name:    "cfmon Version",
		status:  "pass",
		message: fmt.Sprintf("%s (commit: %s, built: %s)", Version, GitCommit, BuildTime),
	})

	// Check 3: Operating System
	checks = append(checks, checkResult{
		name:    "Operating System",
		status:  "pass",
		message: fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	})

	// Check 4: Config File
	configCheck := checkConfigFile()
	checks = append(checks, configCheck)

	// Check 5: API Token
	tokenCheck := checkAPIToken()
	checks = append(checks, tokenCheck)

	// Check 6: Network Connectivity
	networkCheck := checkNetworkConnectivity()
	checks = append(checks, networkCheck)

	// Check 7: Cloudflare API
	if tokenCheck.status == "pass" && networkCheck.status == "pass" {
		apiCheck := checkCloudflareAPI()
		checks = append(checks, apiCheck)
	} else {
		checks = append(checks, checkResult{
			name:    "Cloudflare API",
			status:  "skip",
			message: "Skipped (requires valid token and network)",
		})
	}

	// Check 8: File Permissions
	permCheck := checkFilePermissions()
	checks = append(checks, permCheck)

	// Display results
	displayCheckResults(checks)

	// Summary
	fmt.Println()
	failCount := 0
	warnCount := 0
	for _, check := range checks {
		if check.status == "fail" {
			failCount++
		} else if check.status == "warn" {
			warnCount++
		}
	}

	if failCount > 0 {
		fmt.Println(colorize(fmt.Sprintf("✗ %d check(s) failed", failCount), "red", true))
		fmt.Println("\nPlease address the issues above and run 'cfmon doctor' again.")
		os.Exit(1)
	} else if warnCount > 0 {
		fmt.Println(colorize(fmt.Sprintf("⚠ %d warning(s)", warnCount), "yellow", true))
		fmt.Println("\nSome checks have warnings but cfmon should work.")
	} else {
		fmt.Println(colorize("✓ All checks passed!", "green", true))
		fmt.Println("\ncfmon is properly configured and ready to use.")
	}

	// Tips
	fmt.Println("\n" + colorize("Next Steps:", "cyan", true))
	fmt.Println("  • Run 'cfmon containers list' to list containers")
	fmt.Println("  • Run 'cfmon workers list' to list workers")
	fmt.Println("  • Run 'cfmon help' for more commands")
}

func checkConfigFile() checkResult {
	cfg := config.New()
	configPath := cfg.GetConfigPath()

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return checkResult{
			name:    "Config File",
			status:  "warn",
			message: fmt.Sprintf("Not found at %s (will be created on login)", configPath),
		}
	}

	// Try to load config
	if err := cfg.Load(); err != nil {
		return checkResult{
			name:    "Config File",
			status:  "fail",
			message: fmt.Sprintf("Error loading config: %v", err),
		}
	}

	return checkResult{
		name:    "Config File",
		status:  "pass",
		message: configPath,
	}
}

func checkAPIToken() checkResult {
	cfg := config.New()

	// Check environment variable first
	if envToken := os.Getenv("CFMON_TOKEN"); envToken != "" {
		if len(envToken) < 20 {
			return checkResult{
				name:    "API Token",
				status:  "fail",
				message: "Invalid token in CFMON_TOKEN environment variable",
			}
		}
		return checkResult{
			name:    "API Token",
			status:  "pass",
			message: "Set via CFMON_TOKEN environment variable",
		}
	}

	// Check config file
	if err := cfg.Load(); err == nil && cfg.Token != "" {
		// Mask the token for security
		maskedToken := cfg.Token[:4] + "..." + cfg.Token[len(cfg.Token)-4:]
		return checkResult{
			name:    "API Token",
			status:  "pass",
			message: fmt.Sprintf("Configured (%s)", maskedToken),
		}
	}

	// Check command line flag
	if token != "" {
		maskedToken := token[:min(4, len(token))] + "..."
		return checkResult{
			name:    "API Token",
			status:  "pass",
			message: fmt.Sprintf("Set via --token flag (%s)", maskedToken),
		}
	}

	return checkResult{
		name:    "API Token",
		status:  "fail",
		message: "Not configured. Run 'cfmon login <token>' to set token",
	}
}

func checkNetworkConnectivity() checkResult {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get("https://api.cloudflare.com")
	if err != nil {
		return checkResult{
			name:    "Network Connectivity",
			status:  "fail",
			message: fmt.Sprintf("Cannot reach api.cloudflare.com: %v", err),
		}
	}
	defer resp.Body.Close()

	return checkResult{
		name:    "Network Connectivity",
		status:  "pass",
		message: "Connected to api.cloudflare.com",
	}
}

func checkCloudflareAPI() checkResult {
	cfg := config.New()
	cfg.Load()

	apiToken := cfg.Token
	if token != "" {
		apiToken = token
	}
	if envToken := os.Getenv("CFMON_TOKEN"); envToken != "" {
		apiToken = envToken
	}

	if apiToken == "" {
		return checkResult{
			name:    "Cloudflare API",
			status:  "skip",
			message: "No token configured",
		}
	}

	// Make a simple API call to verify token
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", "https://api.cloudflare.com/client/v4/user/tokens/verify", nil)
	if err != nil {
		return checkResult{
			name:    "Cloudflare API",
			status:  "fail",
			message: fmt.Sprintf("Failed to create request: %v", err),
		}
	}

	req.Header.Set("Authorization", "Bearer "+apiToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return checkResult{
			name:    "Cloudflare API",
			status:  "fail",
			message: fmt.Sprintf("API request failed: %v", err),
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return checkResult{
			name:    "Cloudflare API",
			status:  "pass",
			message: "Token verified successfully",
		}
	} else if resp.StatusCode == 401 || resp.StatusCode == 403 {
		return checkResult{
			name:    "Cloudflare API",
			status:  "fail",
			message: "Invalid or expired token. Run 'cfmon login <token>' with a valid token",
		}
	} else {
		return checkResult{
			name:    "Cloudflare API",
			status:  "warn",
			message: fmt.Sprintf("Unexpected response (HTTP %d)", resp.StatusCode),
		}
	}
}

func checkFilePermissions() checkResult {
	cfg := config.New()
	configDir := filepath.Dir(cfg.GetConfigPath())

	// Check if we can write to config directory
	testFile := filepath.Join(configDir, ".cfmon_test")

	// Try to create the directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return checkResult{
			name:    "File Permissions",
			status:  "fail",
			message: fmt.Sprintf("Cannot create config directory: %v", err),
		}
	}

	// Try to write a test file
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		return checkResult{
			name:    "File Permissions",
			status:  "fail",
			message: fmt.Sprintf("Cannot write to config directory: %v", err),
		}
	}

	// Clean up test file
	os.Remove(testFile)

	return checkResult{
		name:    "File Permissions",
		status:  "pass",
		message: fmt.Sprintf("Can write to %s", configDir),
	}
}

func displayCheckResults(checks []checkResult) {
	maxNameLen := 0
	for _, check := range checks {
		if len(check.name) > maxNameLen {
			maxNameLen = len(check.name)
		}
	}

	for _, check := range checks {
		statusIcon := ""
		statusColor := ""

		switch check.status {
		case "pass":
			statusIcon = "✓"
			statusColor = "green"
		case "fail":
			statusIcon = "✗"
			statusColor = "red"
		case "warn":
			statusIcon = "⚠"
			statusColor = "yellow"
		case "skip":
			statusIcon = "○"
			statusColor = "gray"
		}

		fmt.Printf("  %s %-*s  %s\n",
			colorize(statusIcon, statusColor, false),
			maxNameLen,
			check.name,
			check.message)

		if verbose && check.status == "fail" {
			// Add verbose debug info for failed checks
			fmt.Printf("    %s\n", colorize("Debug: Check with -v flag for more details", "gray", false))
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}