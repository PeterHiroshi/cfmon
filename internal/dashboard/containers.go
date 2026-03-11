package dashboard

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const defaultBarWidth = 15

func (m Model) renderContainers() string {
	if m.data == nil || len(m.data.Containers) == 0 {
		return emptyMessageStyle.Render("No containers found")
	}

	containers := m.data.Containers
	if m.filterText != "" {
		containers = filterContainers(containers, m.filterText)
	}
	if len(containers) == 0 {
		return emptyMessageStyle.Render("No containers match filter")
	}

	var b strings.Builder

	if m.filterMode {
		b.WriteString("Filter: " + m.filterInput.View())
		b.WriteString("\n\n")
	} else if m.filterText != "" {
		b.WriteString(filterActiveStyle.Render("filter: " + m.filterText))
		b.WriteString("\n\n")
	}

	// Find max CPU and memory for bar scaling
	maxCPU := 1000
	maxMem := 128
	for _, c := range containers {
		if c.CPUMS > maxCPU {
			maxCPU = c.CPUMS
		}
		if c.MemoryMB > maxMem {
			maxMem = c.MemoryMB
		}
	}

	// Column widths
	nameW, statusW, cpuW, memW := 20, 10, 10, 14

	// Header
	header := fmt.Sprintf("%-*s %-*s %*s %*s %-*s %-*s",
		nameW, "Name", statusW, "Status", cpuW, "CPU (ms)", memW, "Memory (MB)",
		defaultBarWidth+2, "CPU Bar", defaultBarWidth+2, "Memory Bar")
	b.WriteString(tableHeaderStyle.Render(header))
	b.WriteString("\n")

	// Calculate visible rows
	visibleRows := m.height - 12
	if visibleRows < 3 {
		visibleRows = 3
	}

	maxScroll := len(containers) - visibleRows
	if maxScroll < 0 {
		maxScroll = 0
	}
	offset := m.scrollOffset
	if offset > maxScroll {
		offset = maxScroll
	}

	end := offset + visibleRows
	if end > len(containers) {
		end = len(containers)
	}

	// Totals
	var totalCPU, totalMem int
	for _, c := range containers {
		totalCPU += c.CPUMS
		totalMem += c.MemoryMB
	}
	avgMem := totalMem / len(containers)

	// Rows
	for i, c := range containers[offset:end] {
		statusStyled := lipgloss.NewStyle().Foreground(lipgloss.Color(statusColor(c.Status))).Render(fmt.Sprintf("%-*s", statusW, c.Status))
		cpuBar := renderBar(c.CPUMS, maxCPU, defaultBarWidth)
		memBar := renderBar(c.MemoryMB, maxMem, defaultBarWidth)

		row := fmt.Sprintf("%-*s %s %*d %*d %s %s",
			nameW, truncate(c.Name, nameW), statusStyled, cpuW, c.CPUMS, memW, c.MemoryMB, cpuBar, memBar)
		if offset+i == m.selectedRow {
			b.WriteString(selectedRowStyle.Render(row))
		} else {
			b.WriteString(tableRowStyle.Render(row))
		}
		b.WriteString("\n")
	}

	// Scroll indicator
	if len(containers) > visibleRows {
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render(
			fmt.Sprintf("  showing %d-%d of %d (j/k to scroll)", offset+1, end, len(containers))))
		b.WriteString("\n")
	}

	// Totals row
	totals := fmt.Sprintf("%-*s %-*s %*d %*d (avg: %d)",
		nameW, "TOTAL", statusW, "", cpuW, totalCPU, memW, totalMem, avgMem)
	b.WriteString(tableTotalsStyle.Render(totals))

	return b.String()
}

func renderBar(value, max, width int) string {
	if max <= 0 {
		return strings.Repeat(string(gaugeEmptyChar), width)
	}
	filled := value * width / max
	if filled > width {
		filled = width
	}
	if filled < 0 {
		filled = 0
	}
	empty := width - filled
	return strings.Repeat(string(gaugeFillChar), filled) + strings.Repeat(string(gaugeEmptyChar), empty)
}
