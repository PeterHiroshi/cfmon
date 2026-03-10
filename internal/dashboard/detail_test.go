package dashboard

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/PeterHiroshi/cfmon/internal/api"
)

func TestDetailOpenWithEnter(t *testing.T) {
	m := Model{
		width: 80, height: 24,
		activeTab:   TabWorkers,
		selectedRow: 0,
		data: &DashboardData{
			Workers: []api.Worker{{Name: "w1", Status: "active"}},
		},
	}
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	updated := newModel.(Model)
	if !updated.showDetail {
		t.Error("Enter should open detail view")
	}
}

func TestDetailNotOnOverview(t *testing.T) {
	m := Model{
		width: 80, height: 24,
		activeTab: TabOverview,
	}
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	updated := newModel.(Model)
	if updated.showDetail {
		t.Error("Enter on Overview should not open detail")
	}
}

func TestDetailCloseWithEsc(t *testing.T) {
	m := Model{
		width: 80, height: 24,
		activeTab:  TabWorkers,
		showDetail: true,
		data: &DashboardData{
			Workers: []api.Worker{{Name: "w1", Status: "active"}},
		},
	}
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	updated := newModel.(Model)
	if updated.showDetail {
		t.Error("Esc should close detail view")
	}
}

func TestDetailTabSwitchClosesDetail(t *testing.T) {
	m := Model{
		width: 80, height: 24,
		activeTab:  TabWorkers,
		showDetail: true,
		data: &DashboardData{
			Workers: []api.Worker{{Name: "w1", Status: "active"}},
		},
	}
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	updated := newModel.(Model)
	if updated.showDetail {
		t.Error("Tab switch should close detail view")
	}
	if updated.activeTab != TabContainers {
		t.Errorf("should switch to Containers tab, got %d", updated.activeTab)
	}
}

func TestRenderWorkerDetail(t *testing.T) {
	m := Model{
		width: 80, height: 24,
		activeTab:   TabWorkers,
		selectedRow: 0,
		showDetail:  true,
		data: &DashboardData{
			Workers: []api.Worker{
				{ID: "w-123", Name: "api-gateway", Status: "active", Requests: 1000, Errors: 5, CPUMS: 12, SuccessRate: 99.5},
			},
		},
	}
	result := m.renderWorkerDetail()
	checks := []string{"api-gateway", "w-123", "active", "1000", "12"}
	for _, c := range checks {
		if !strings.Contains(result, c) {
			t.Errorf("worker detail should contain %q", c)
		}
	}
}

func TestRenderContainerDetail(t *testing.T) {
	m := Model{
		width: 80, height: 24,
		activeTab:   TabContainers,
		selectedRow: 0,
		showDetail:  true,
		data: &DashboardData{
			Containers: []api.Container{
				{ID: "c-456", Name: "web-app", Status: "running", CPUMS: 500, MemoryMB: 64, Requests: 2000},
			},
		},
	}
	result := m.renderContainerDetail()
	checks := []string{"web-app", "c-456", "running", "500", "64", "2000"}
	for _, c := range checks {
		if !strings.Contains(result, c) {
			t.Errorf("container detail should contain %q", c)
		}
	}
}

func TestDetailSuppressesJK(t *testing.T) {
	m := Model{
		width: 80, height: 24,
		activeTab:   TabWorkers,
		showDetail:  true,
		selectedRow: 2,
		data: &DashboardData{
			Workers: make([]api.Worker, 10),
		},
	}
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	updated := newModel.(Model)
	if updated.selectedRow != 2 {
		t.Errorf("j should not move selectedRow in detail view, got %d", updated.selectedRow)
	}
}
