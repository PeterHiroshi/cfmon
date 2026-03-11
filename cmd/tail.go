package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/PeterHiroshi/cfmon/internal/api"
	"github.com/PeterHiroshi/cfmon/internal/config"
	"github.com/PeterHiroshi/cfmon/internal/tail"
	"github.com/spf13/cobra"
)

var (
	tailFormat            string
	tailStatus            string
	tailMethod            string
	tailSearch            string
	tailIP                string
	tailSampleRate        float64
	tailHeaders           []string
	tailSince             string
	tailMaxEvents         int
	tailIncludeLogs       bool
	tailIncludeExceptions bool
)

var tailCmd = &cobra.Command{
	Use:   "tail [account-id] <worker-name>",
	Short: "Stream real-time logs from a Worker",
	Long: `Stream real-time logs from a Cloudflare Worker via the Tail API.

Connects via WebSocket to stream live request logs, console output,
and exceptions. Supports filtering, multiple output formats, and
auto-reconnect.`,
	Example: `  # Tail a worker (using default account)
  cfmon tail my-worker

  # Tail with account ID
  cfmon tail abc123 my-worker

  # Pretty output with search filter
  cfmon tail my-worker --search "error" --format pretty

  # Compact one-line output
  cfmon tail my-worker --format compact

  # JSON output for piping
  cfmon tail my-worker --format json | jq .

  # Filter by status and method
  cfmon tail my-worker --status error --method POST

  # Limit to 100 events
  cfmon tail my-worker -n 100

  # Sample 10% of traffic
  cfmon tail my-worker --sample-rate 0.1`,
	Args: cobra.RangeArgs(1, 2),
	RunE: runTail,
}

func init() {
	rootCmd.AddCommand(tailCmd)

	tailCmd.Flags().StringVarP(&tailFormat, "format", "f", "pretty", "Output format: pretty, json, compact")
	tailCmd.Flags().StringVarP(&tailStatus, "status", "s", "", `Filter by status: "ok", "error", or specific codes like "500"`)
	tailCmd.Flags().StringVar(&tailMethod, "method", "", "Filter by HTTP method (GET, POST, etc.)")
	tailCmd.Flags().StringVar(&tailSearch, "search", "", "Filter logs containing this string")
	tailCmd.Flags().StringVar(&tailIP, "ip", "", "Filter by client IP address")
	tailCmd.Flags().Float64Var(&tailSampleRate, "sample-rate", 1.0, "Sampling rate 0.0-1.0")
	tailCmd.Flags().StringArrayVarP(&tailHeaders, "header", "H", nil, "Filter by header (key:value)")
	tailCmd.Flags().StringVar(&tailSince, "since", "", `Only show events after this duration ago (e.g. "5m", "1h")`)
	tailCmd.Flags().IntVarP(&tailMaxEvents, "max-events", "n", 0, "Stop after N events (0 = unlimited)")
	tailCmd.Flags().BoolVar(&tailIncludeLogs, "include-logs", true, "Show console.log() output")
	tailCmd.Flags().BoolVar(&tailIncludeExceptions, "include-exceptions", true, "Show exceptions")
}

func runTail(cmd *cobra.Command, args []string) error {
	var accountID, workerName string

	if len(args) == 2 {
		accountID = args[0]
		workerName = args[1]
	} else {
		workerName = args[0]

		configPath := cfgFile
		if configPath == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("getting home directory: %w", err)
			}
			configPath = filepath.Join(home, ".cfmon", "config.yaml")
		}

		cfg, err := config.Load(configPath)
		if err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("loading config: %w", err)
		}

		if cfg == nil || cfg.DefaultAccountID == "" {
			return fmt.Errorf("no account ID provided and no default account set. Use 'cfmon accounts set-default <account-id>' to set a default")
		}

		accountID = cfg.DefaultAccountID
	}

	if tailSampleRate < 0 || tailSampleRate > 1 {
		return fmt.Errorf("sample rate must be between 0.0 and 1.0, got %.2f", tailSampleRate)
	}

	apiToken, err := getAPIToken()
	if err != nil {
		return err
	}

	filter := api.TailFilter{
		SamplingRate: tailSampleRate,
	}

	if tailStatus != "" {
		filter.Status = []string{tailStatus}
	}
	if tailMethod != "" {
		filter.Method = []string{strings.ToUpper(tailMethod)}
	}
	if tailIP != "" {
		filter.ClientIP = []string{tailIP}
	}
	if len(tailHeaders) > 0 {
		filter.Headers = make(map[string]string)
		for _, h := range tailHeaders {
			parts := strings.SplitN(h, ":", 2)
			if len(parts) == 2 {
				filter.Headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
		}
	}

	var sinceDuration time.Duration
	if tailSince != "" {
		sinceDuration, err = time.ParseDuration(tailSince)
		if err != nil {
			return fmt.Errorf("invalid --since duration: %w", err)
		}
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "Debug: Tailing worker %q in account %s\n", workerName, accountID)
	}

	client := api.NewClient(apiToken)
	if timeout > 0 {
		client.SetTimeout(timeout)
	}

	session, err := client.CreateTail(accountID, workerName, filter)
	if err != nil {
		return fmt.Errorf("creating tail session: %w", err)
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "Debug: Tail session created: %s\n", session.ID)
		fmt.Fprintf(os.Stderr, "Debug: WebSocket URL: %s\n", session.URL)
	}

	cleanup := func() {
		if verbose {
			fmt.Fprintf(os.Stderr, "Debug: Deleting tail session %s\n", session.ID)
		}
		if err := client.DeleteTail(accountID, workerName, session.ID); err != nil {
			if verbose {
				fmt.Fprintf(os.Stderr, "Debug: Failed to delete tail: %v\n", err)
			}
		}
	}

	formatter := tail.NewFormatter(tailFormat, noColor)
	formatter.IncludeLogs = tailIncludeLogs
	formatter.IncludeExceptions = tailIncludeExceptions

	engine := tail.NewEngine(tail.EngineConfig{
		WebSocketURL: session.URL,
		MaxEvents:    tailMaxEvents,
		Search:       tailSearch,
		Since:        sinceDuration,
		OnEvent: func(event tail.TailEvent) {
			fmt.Print(formatter.Format(event))
		},
		OnError: func(err error) {
			fmt.Fprintf(os.Stderr, "WebSocket error: %v\n", err)
		},
	})

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		fmt.Fprintf(os.Stderr, "\nShutting down...\n")
		engine.Stop()
	}()

	if !quiet {
		fmt.Fprintf(os.Stderr, "Tailing %s... (Ctrl+C to stop)\n", workerName)
	}

	engine.Run()

	cleanup()

	return nil
}
