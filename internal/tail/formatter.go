package tail

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Formatter formats TailEvents for display
type Formatter struct {
	FormatType        string
	NoColor           bool
	IncludeLogs       bool
	IncludeExceptions bool
}

// NewFormatter creates a new Formatter with sensible defaults
func NewFormatter(format string, noColor bool) *Formatter {
	return &Formatter{
		FormatType:        format,
		NoColor:           noColor,
		IncludeLogs:       true,
		IncludeExceptions: true,
	}
}

// Format formats a TailEvent as a string
func (f *Formatter) Format(event TailEvent) string {
	switch f.FormatType {
	case "json":
		return f.formatJSON(event)
	case "compact":
		return f.formatCompact(event)
	default:
		return f.formatPretty(event)
	}
}

func (f *Formatter) formatJSON(event TailEvent) string {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Sprintf(`{"error":"marshal failed: %s"}`, err)
	}
	return string(data)
}

func (f *Formatter) formatCompact(event TailEvent) string {
	ts := event.Time().Format("15:04:05")
	status := event.Event.Response.Status
	method := event.Event.Request.Method
	url := event.Event.Request.URL

	statusStr := fmt.Sprintf("%d", status)
	if !f.NoColor {
		statusStr = f.colorizeStatus(status)
	}

	return fmt.Sprintf("[%s] %s %s %s", ts, statusStr, method, url)
}

func (f *Formatter) formatPretty(event TailEvent) string {
	var b strings.Builder

	ts := event.Time().Format("15:04:05.000")
	method := event.Event.Request.Method
	url := event.Event.Request.URL
	status := event.Event.Response.Status
	outcome := event.Outcome

	if f.NoColor {
		fmt.Fprintf(&b, "%s %s %s → %d (%s)\n", ts, method, url, status, outcome)
	} else {
		fmt.Fprintf(&b, "%s %s %s → %s (%s)\n",
			f.dim(ts), f.bold(method), f.cyan(url),
			f.colorizeStatus(status), f.colorizeOutcome(outcome))
	}

	if f.IncludeLogs && len(event.Logs) > 0 {
		for _, log := range event.Logs {
			msg := strings.Join(log.Message, " ")
			if f.NoColor {
				fmt.Fprintf(&b, "  [%s] %s\n", log.Level, msg)
			} else {
				fmt.Fprintf(&b, "  %s %s\n", f.colorizeLogLevel(log.Level), msg)
			}
		}
	}

	if f.IncludeExceptions && len(event.Exceptions) > 0 {
		for _, exc := range event.Exceptions {
			if f.NoColor {
				fmt.Fprintf(&b, "  ✗ %s: %s\n", exc.Name, exc.Message)
			} else {
				fmt.Fprintf(&b, "  %s %s: %s\n", f.red("✗"), f.red(exc.Name), exc.Message)
			}
		}
	}

	return b.String()
}

func (f *Formatter) colorizeStatus(status int) string {
	s := fmt.Sprintf("%d", status)
	switch {
	case status >= 200 && status < 300:
		return f.green(s)
	case status >= 300 && status < 400:
		return f.yellow(s)
	case status >= 400:
		return f.red(s)
	default:
		return s
	}
}

func (f *Formatter) colorizeOutcome(outcome string) string {
	switch outcome {
	case "ok":
		return f.green(outcome)
	case "exception", "exceededCpu", "exceededMemory":
		return f.red(outcome)
	case "canceled":
		return f.yellow(outcome)
	default:
		return outcome
	}
}

func (f *Formatter) colorizeLogLevel(level string) string {
	tag := fmt.Sprintf("[%s]", level)
	switch level {
	case "error":
		return f.red(tag)
	case "warn":
		return f.yellow(tag)
	case "debug":
		return f.dim(tag)
	default:
		return f.cyan(tag)
	}
}

func (f *Formatter) ansi(code, text string) string {
	return fmt.Sprintf("\033[%sm%s\033[0m", code, text)
}

func (f *Formatter) red(s string) string     { return f.ansi("31", s) }
func (f *Formatter) green(s string) string   { return f.ansi("32", s) }
func (f *Formatter) yellow(s string) string  { return f.ansi("33", s) }
func (f *Formatter) cyan(s string) string    { return f.ansi("36", s) }
func (f *Formatter) bold(s string) string    { return f.ansi("1", s) }
func (f *Formatter) dim(s string) string     { return f.ansi("2", s) }
