package output

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fatih/color"
)

// FormatJSON converts data to pretty-printed JSON
func FormatJSON(data interface{}) (string, error) {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// FormatTable formats data as a table
func FormatTable(headers []string, rows [][]string) string {
	if len(headers) == 0 {
		return ""
	}

	// Calculate column widths
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}

	for _, row := range rows {
		for i, cell := range row {
			if i < len(widths) && len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	var sb strings.Builder

	// Print headers
	for i, h := range headers {
		if i > 0 {
			sb.WriteString("  ")
		}
		sb.WriteString(fmt.Sprintf("%-*s", widths[i], h))
	}
	sb.WriteString("\n")

	// Print separator
	for i, w := range widths {
		if i > 0 {
			sb.WriteString("  ")
		}
		sb.WriteString(strings.Repeat("-", w))
	}
	sb.WriteString("\n")

	// Print rows
	for _, row := range rows {
		for i, cell := range row {
			if i > 0 {
				sb.WriteString("  ")
			}
			if i < len(widths) {
				sb.WriteString(fmt.Sprintf("%-*s", widths[i], cell))
			}
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// FormatColoredTable formats data as a table with color support
func FormatColoredTable(headers []string, rows [][]string, enableColors bool) string {
	if len(headers) == 0 {
		return ""
	}

	// Disable colors if requested
	if !enableColors {
		color.NoColor = true
	} else {
		color.NoColor = false
	}

	// Calculate column widths (using plain text length)
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}

	for _, row := range rows {
		for i, cell := range row {
			if i < len(widths) && len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	var sb strings.Builder

	// Print headers
	for i, h := range headers {
		if i > 0 {
			sb.WriteString("  ")
		}
		sb.WriteString(fmt.Sprintf("%-*s", widths[i], h))
	}
	sb.WriteString("\n")

	// Print separator
	for i, w := range widths {
		if i > 0 {
			sb.WriteString("  ")
		}
		sb.WriteString(strings.Repeat("-", w))
	}
	sb.WriteString("\n")

	// Print rows with colors
	for _, row := range rows {
		for i, cell := range row {
			if i > 0 {
				sb.WriteString("  ")
			}
			if i < len(widths) {
				// Apply color based on cell content
				coloredCell := colorizeCell(cell)
				// Calculate padding based on original cell length
				padding := widths[i] - len(cell)
				sb.WriteString(coloredCell)
				if padding > 0 {
					sb.WriteString(strings.Repeat(" ", padding))
				}
			}
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// colorizeCell applies color to a cell based on its content
func colorizeCell(cell string) string {
	lower := strings.ToLower(cell)

	// Green for healthy/active status
	if lower == "active" || lower == "healthy" || lower == "true" {
		return color.GreenString(cell)
	}

	// Red for errors/failures
	if lower == "error" || lower == "failed" || lower == "false" || lower == "invalid" {
		return color.RedString(cell)
	}

	// Yellow for warnings/pending
	if lower == "warning" || lower == "pending" || lower == "degraded" {
		return color.YellowString(cell)
	}

	// Default: no color
	return cell
}
