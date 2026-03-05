package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(helpCmd)
}

var helpCmd = &cobra.Command{
	Use:   "help [command]",
	Short: "Display help information for cfmon",
	Long:  `Display detailed help information for cfmon commands with examples and usage patterns.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			displayMainHelp()
		} else {
			// Show help for specific command
			helpCmd, _, err := rootCmd.Find(args)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Unknown command: %s\n", strings.Join(args, " "))
				fmt.Fprintln(os.Stderr, "Run 'cfmon help' for usage.")
				os.Exit(1)
			}
			helpCmd.Help()
		}
	},
}

func displayMainHelp() {
	help := `
` + colorize("cfmon - Cloudflare Workers/Containers CLI Monitor", "cyan", true) + `

` + colorize("USAGE", "yellow", true) + `
  cfmon [command] [flags]
  cfmon [command] [subcommand] [flags]

` + colorize("CORE COMMANDS", "yellow", true) + `
  ` + colorize("login", "green", false) + `         Store your Cloudflare API token
  ` + colorize("doctor", "green", false) + `        Check system configuration and connectivity
  ` + colorize("help", "green", false) + `          Display help information

` + colorize("RESOURCE COMMANDS", "yellow", true) + `
  ` + colorize("containers", "green", false) + `    Manage and monitor Cloudflare Containers
    list         List all containers with resource usage
    status       Get detailed status of a specific container

  ` + colorize("workers", "green", false) + `       Manage and monitor Cloudflare Workers
    list         List all workers with metrics
    status       Get detailed status of a specific worker

` + colorize("CONFIGURATION COMMANDS", "yellow", true) + `
  ` + colorize("config", "green", false) + `        Manage cfmon configuration
    show         Display current configuration
    path         Show configuration file location

  ` + colorize("completion", "green", false) + `    Generate shell completion scripts
    bash         Generate bash completion script
    zsh          Generate zsh completion script
    fish         Generate fish completion script
    powershell   Generate PowerShell completion script

` + colorize("GLOBAL FLAGS", "yellow", true) + `
  -f, --format string    Output format (table, json) (default "table")
  -v, --verbose          Enable verbose debug output
      --timeout duration API request timeout (default 30s)
      --no-color         Disable colored output
      --config string    Config file path (default $HOME/.cfmon/config.yaml)
      --token string     Override API token from config
  -h, --help            Help for cfmon
      --version         Version for cfmon

` + colorize("EXAMPLES", "yellow", true) + `
  ` + colorize("# Initial setup", "blue", false) + `
  cfmon login <your-cloudflare-api-token>
  cfmon doctor

  ` + colorize("# List resources", "blue", false) + `
  cfmon containers list
  cfmon containers list --format json
  cfmon containers list --filter "prod" --limit 5
  cfmon containers list --sort cpu

  ` + colorize("# Get status", "blue", false) + `
  cfmon workers status my-worker
  cfmon containers status container-id --json

  ` + colorize("# Configuration", "blue", false) + `
  cfmon config show
  cfmon config path

  ` + colorize("# Shell completion", "blue", false) + `
  cfmon completion bash > /etc/bash_completion.d/cfmon
  cfmon completion zsh > "${fpath[1]}/_cfmon"

` + colorize("ENVIRONMENT VARIABLES", "yellow", true) + `
  CFMON_TOKEN       Cloudflare API token
  CFMON_CONFIG      Path to config file
  CFMON_FORMAT      Default output format
  CFMON_NO_COLOR    Disable colors (set to any value)

` + colorize("LEARN MORE", "yellow", true) + `
  Use 'cfmon <command> --help' for more information about a command.
  Documentation: https://github.com/PeterHiroshi/cfmon
  Report issues: https://github.com/PeterHiroshi/cfmon/issues

` + colorize("QUICK TIPS", "yellow", true) + `
  • Use ` + colorize("--json", "cyan", false) + ` flag for scripting and automation
  • Use ` + colorize("--filter", "cyan", false) + ` to search by name substring
  • Use ` + colorize("--sort", "cyan", false) + ` to order results (name, cpu, memory, requests)
  • Use ` + colorize("--limit", "cyan", false) + ` to show only top N results
  • Add ` + colorize("-v", "cyan", false) + ` flag for debug output when troubleshooting
`

	fmt.Print(help)
}

// colorize adds ANSI color codes to text if colors are enabled
func colorize(text, color string, bold bool) string {
	if noColor || os.Getenv("CFMON_NO_COLOR") != "" {
		return text
	}

	var code string
	switch color {
	case "red":
		code = "31"
	case "green":
		code = "32"
	case "yellow":
		code = "33"
	case "blue":
		code = "34"
	case "magenta":
		code = "35"
	case "cyan":
		code = "36"
	default:
		return text
	}

	if bold {
		return fmt.Sprintf("\033[1;%sm%s\033[0m", code, text)
	}
	return fmt.Sprintf("\033[%sm%s\033[0m", code, text)
}