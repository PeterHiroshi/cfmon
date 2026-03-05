package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/PeterHiroshi/cfmon/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage cfmon configuration",
	Long:  `Manage cfmon configuration including viewing current settings and configuration file location.`,
	Example: `  # Show current configuration
  cfmon config show

  # Show configuration file path
  cfmon config path`,
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Display current configuration",
	Long:  `Display the current cfmon configuration with sensitive values masked for security.`,
	Run:   runConfigShow,
}

var configPathCmd = &cobra.Command{
	Use:   "path",
	Short: "Show configuration file location",
	Long:  `Display the path to the cfmon configuration file.`,
	Run:   runConfigPath,
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configPathCmd)
}

func runConfigShow(cmd *cobra.Command, args []string) {
	cfg := config.New()
	configPath := cfg.GetConfigPath()

	// Load configuration
	if err := cfg.Load(); err != nil {
		if os.IsNotExist(err) {
			fmt.Println(colorize("Configuration", "cyan", true))
			fmt.Println(strings.Repeat("-", 50))
			fmt.Printf("Config File: %s (not found)\n", configPath)
			fmt.Println("\nNo configuration found. Run 'cfmon login <token>' to set up.")
			return
		}
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Display configuration
	fmt.Println(colorize("cfmon Configuration", "cyan", true))
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println()

	// Config file location
	fmt.Printf("%s %s\n", colorize("Config File:", "yellow", true), configPath)
	fmt.Println()

	// Authentication
	fmt.Println(colorize("Authentication", "yellow", true))
	fmt.Println(strings.Repeat("-", 30))

	// Check for token in different locations
	displayToken := ""
	tokenSource := ""

	if envToken := os.Getenv("CFMON_TOKEN"); envToken != "" {
		displayToken = maskToken(envToken)
		tokenSource = "environment variable (CFMON_TOKEN)"
	} else if token != "" {
		displayToken = maskToken(token)
		tokenSource = "command line flag (--token)"
	} else if cfg.Token != "" {
		displayToken = maskToken(cfg.Token)
		tokenSource = "config file"
	}

	if displayToken != "" {
		fmt.Printf("  API Token:    %s\n", displayToken)
		fmt.Printf("  Source:       %s\n", tokenSource)
	} else {
		fmt.Printf("  API Token:    %s\n", colorize("Not configured", "red", false))
		fmt.Printf("  %s Run 'cfmon login <token>' to configure\n", colorize("→", "yellow", false))
	}

	// API Settings
	fmt.Println()
	fmt.Println(colorize("API Settings", "yellow", true))
	fmt.Println(strings.Repeat("-", 30))

	if cfg.APIEndpoint != "" {
		fmt.Printf("  API Endpoint: %s\n", cfg.APIEndpoint)
	} else {
		fmt.Printf("  API Endpoint: %s (default)\n", "https://api.cloudflare.com/client/v4")
	}

	timeoutValue := timeout
	if timeoutValue == 0 {
		timeoutValue = 30 * time.Second
	}
	fmt.Printf("  Timeout:      %s\n", timeoutValue)

	// Output Settings
	fmt.Println()
	fmt.Println(colorize("Output Settings", "yellow", true))
	fmt.Println(strings.Repeat("-", 30))

	outputFormat := format
	if outputFormat == "" {
		if cfg.DefaultFormat != "" {
			outputFormat = cfg.DefaultFormat
		} else {
			outputFormat = "table"
		}
	}
	fmt.Printf("  Format:       %s\n", outputFormat)
	fmt.Printf("  Colored:      %t\n", !noColor && os.Getenv("CFMON_NO_COLOR") == "")
	fmt.Printf("  Verbose:      %t\n", verbose)

	// Environment Variables
	fmt.Println()
	fmt.Println(colorize("Environment Variables", "yellow", true))
	fmt.Println(strings.Repeat("-", 30))

	envVars := []struct {
		name  string
		value string
	}{
		{"CFMON_TOKEN", os.Getenv("CFMON_TOKEN")},
		{"CFMON_CONFIG", os.Getenv("CFMON_CONFIG")},
		{"CFMON_FORMAT", os.Getenv("CFMON_FORMAT")},
		{"CFMON_NO_COLOR", os.Getenv("CFMON_NO_COLOR")},
	}

	hasEnvVars := false
	for _, ev := range envVars {
		if ev.value != "" {
			hasEnvVars = true
			displayValue := ev.value
			if ev.name == "CFMON_TOKEN" {
				displayValue = maskToken(ev.value)
			}
			fmt.Printf("  %s: %s\n", ev.name, displayValue)
		}
	}

	if !hasEnvVars {
		fmt.Printf("  %s\n", colorize("(none set)", "gray", false))
	}

	// Additional info
	fmt.Println()
	fmt.Println(colorize("Additional Information", "yellow", true))
	fmt.Println(strings.Repeat("-", 30))
	fmt.Printf("  Version:      %s\n", Version)
	fmt.Printf("  Build Time:   %s\n", BuildTime)
	fmt.Printf("  Git Commit:   %s\n", GitCommit)

	// Tips
	fmt.Println()
	fmt.Println(colorize("Tips:", "cyan", true))
	fmt.Println("  • Use 'cfmon doctor' to check system health")
	fmt.Println("  • Use 'cfmon help' for command information")
	fmt.Println("  • Set CFMON_TOKEN environment variable to avoid using --token flag")
}

func runConfigPath(cmd *cobra.Command, args []string) {
	cfg := config.New()
	configPath := cfg.GetConfigPath()

	if format == "json" {
		fmt.Printf(`{"config_path": "%s"}`+"\n", configPath)
	} else {
		fmt.Println(configPath)
	}
}

func maskToken(token string) string {
	if len(token) <= 8 {
		return "****"
	}
	return token[:4] + strings.Repeat("*", len(token)-8) + token[len(token)-4:]
}