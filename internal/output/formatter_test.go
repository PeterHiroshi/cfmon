package output

import (
	"encoding/json"
	"strings"
	"testing"
)

type testData struct {
	ID   string
	Name string
	CPU  int
}

func TestFormatJSON(t *testing.T) {
	data := []testData{
		{ID: "1", Name: "test1", CPU: 100},
		{ID: "2", Name: "test2", CPU: 200},
	}

	result, err := FormatJSON(data)
	if err != nil {
		t.Fatalf("FormatJSON() error = %v", err)
	}

	// Verify it's valid JSON
	var parsed []testData
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		t.Fatalf("result is not valid JSON: %v", err)
	}

	if len(parsed) != 2 {
		t.Errorf("parsed length = %d, want 2", len(parsed))
	}
}

func TestFormatTable(t *testing.T) {
	headers := []string{"ID", "Name", "CPU"}
	rows := [][]string{
		{"1", "test1", "100"},
		{"2", "test2", "200"},
	}

	result := FormatTable(headers, rows)

	// Verify headers are present
	if !strings.Contains(result, "ID") {
		t.Errorf("result missing ID header")
	}
	if !strings.Contains(result, "Name") {
		t.Errorf("result missing Name header")
	}
	if !strings.Contains(result, "CPU") {
		t.Errorf("result missing CPU header")
	}

	// Verify data rows are present
	if !strings.Contains(result, "test1") {
		t.Errorf("result missing test1 data")
	}
	if !strings.Contains(result, "test2") {
		t.Errorf("result missing test2 data")
	}
}

func TestFormatTable_Empty(t *testing.T) {
	headers := []string{"ID", "Name"}
	rows := [][]string{}

	result := FormatTable(headers, rows)

	// Should still show headers
	if !strings.Contains(result, "ID") {
		t.Errorf("result missing ID header")
	}
}

func TestFormatColoredTable_WithColors(t *testing.T) {
	headers := []string{"Status", "Name", "Value"}
	rows := [][]string{
		{"active", "Service 1", "100"},
		{"error", "Service 2", "200"},
		{"warning", "Service 3", "300"},
	}

	result := FormatColoredTable(headers, rows, true)

	// Verify ANSI escape codes are present (colors enabled)
	if !strings.Contains(result, "\x1b[") {
		t.Errorf("result should contain ANSI escape codes when colors enabled")
	}

	// Verify data is present
	if !strings.Contains(result, "Service 1") {
		t.Errorf("result missing Service 1 data")
	}
}

func TestFormatColoredTable_WithoutColors(t *testing.T) {
	headers := []string{"Status", "Name", "Value"}
	rows := [][]string{
		{"active", "Service 1", "100"},
		{"error", "Service 2", "200"},
	}

	result := FormatColoredTable(headers, rows, false)

	// Verify no ANSI escape codes are present (colors disabled)
	if strings.Contains(result, "\x1b[") {
		t.Errorf("result should not contain ANSI escape codes when colors disabled")
	}

	// Verify data is still present
	if !strings.Contains(result, "Service 1") {
		t.Errorf("result missing Service 1 data")
	}
	if !strings.Contains(result, "active") {
		t.Errorf("result missing status data")
	}
}

func TestFormatColoredTable_ColorMapping(t *testing.T) {
	headers := []string{"Status"}
	rows := [][]string{
		{"active"},
		{"healthy"},
		{"error"},
		{"failed"},
		{"warning"},
		{"pending"},
		{"normal"},
	}

	result := FormatColoredTable(headers, rows, true)

	// Just verify the function runs and produces output
	if len(result) == 0 {
		t.Errorf("result should not be empty")
	}

	// Verify all status values are present
	for _, row := range rows {
		if !strings.Contains(result, row[0]) {
			t.Errorf("result missing status: %s", row[0])
		}
	}
}
