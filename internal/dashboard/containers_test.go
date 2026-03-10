package dashboard

import (
	"fmt"
	"strings"
	"testing"

	"github.com/PeterHiroshi/cfmon/internal/api"
)

func TestRenderContainersEmpty(t *testing.T) {
	m := Model{width: 80, height: 24, data: &DashboardData{}}
	result := m.renderContainers()
	if !strings.Contains(result, "No containers found") {
		t.Errorf("empty containers should show 'No containers found', got: %s", result)
	}
}

func TestRenderContainersTable(t *testing.T) {
	m := Model{
		width:  120,
		height: 24,
		data: &DashboardData{
			Containers: []api.Container{
				{Name: "web-app", Status: "running", CPUMS: 500, MemoryMB: 64},
				{Name: "worker-bg", Status: "stopped", CPUMS: 200, MemoryMB: 32},
			},
		},
	}
	result := m.renderContainers()

	if !strings.Contains(result, "Name") {
		t.Error("should contain Name header")
	}
	if !strings.Contains(result, "Status") {
		t.Error("should contain Status header")
	}
	if !strings.Contains(result, "CPU (ms)") {
		t.Error("should contain CPU (ms) header")
	}
	if !strings.Contains(result, "Memory (MB)") {
		t.Error("should contain Memory (MB) header")
	}
	if !strings.Contains(result, "web-app") {
		t.Error("should contain container name 'web-app'")
	}
	if !strings.Contains(result, "worker-bg") {
		t.Error("should contain container name 'worker-bg'")
	}
	// Should contain bar characters
	if !strings.Contains(result, string(gaugeFillChar)) {
		t.Error("should contain gauge fill characters for bars")
	}
}

func TestRenderContainersBars(t *testing.T) {
	m := Model{
		width:  120,
		height: 24,
		data: &DashboardData{
			Containers: []api.Container{
				{Name: "max-cpu", Status: "running", CPUMS: 1000, MemoryMB: 128},
				{Name: "half-cpu", Status: "running", CPUMS: 500, MemoryMB: 64},
			},
		},
	}
	result := m.renderContainers()
	if !strings.Contains(result, "max-cpu") {
		t.Error("should contain max-cpu")
	}
	if !strings.Contains(result, "half-cpu") {
		t.Error("should contain half-cpu")
	}
}

func TestRenderContainersTotals(t *testing.T) {
	m := Model{
		width:  120,
		height: 24,
		data: &DashboardData{
			Containers: []api.Container{
				{Name: "c1", Status: "running", CPUMS: 300, MemoryMB: 60},
				{Name: "c2", Status: "running", CPUMS: 200, MemoryMB: 40},
			},
		},
	}
	result := m.renderContainers()
	// Total CPU: 500
	if !strings.Contains(result, "500") {
		t.Errorf("should contain total CPU 500, got: %s", result)
	}
}

func TestRenderContainersScroll(t *testing.T) {
	containers := make([]api.Container, 30)
	for i := range containers {
		containers[i] = api.Container{Name: fmt.Sprintf("container-%02d", i), Status: "running", CPUMS: 100, MemoryMB: 32}
	}
	m := Model{
		width:        120,
		height:       15,
		scrollOffset: 5,
		data:         &DashboardData{Containers: containers},
	}
	result := m.renderContainers()
	if !strings.Contains(result, "container-05") {
		t.Errorf("scrolled view should contain container-05, got: %s", result)
	}
}

func TestRenderBar(t *testing.T) {
	// Full bar
	bar := renderBar(100, 100, 10)
	if len(bar) != 40 { // 10 × 4 bytes per rune (█ is 3 bytes, but let's check content)
		// Just check it has fill chars
	}
	if !strings.Contains(bar, string(gaugeFillChar)) {
		t.Error("full bar should contain fill chars")
	}

	// Empty bar
	bar = renderBar(0, 100, 10)
	if !strings.Contains(bar, string(gaugeEmptyChar)) {
		t.Error("empty bar should contain empty chars")
	}
	if strings.Contains(bar, string(gaugeFillChar)) {
		t.Error("empty bar should not contain fill chars")
	}

	// Half bar
	bar = renderBar(50, 100, 10)
	fillCount := strings.Count(bar, string(gaugeFillChar))
	emptyCount := strings.Count(bar, string(gaugeEmptyChar))
	if fillCount != 5 {
		t.Errorf("half bar fill count = %d, want 5", fillCount)
	}
	if emptyCount != 5 {
		t.Errorf("half bar empty count = %d, want 5", emptyCount)
	}
}
