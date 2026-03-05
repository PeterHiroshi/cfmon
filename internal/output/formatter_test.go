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

func TestFormatJSON_Error(t *testing.T) {
	// Test with un-marshalable data (channel)
	ch := make(chan int)
	defer close(ch)

	_, err := FormatJSON(ch)
	if err == nil {
		t.Fatal("FormatJSON() with channel: error = nil, want error")
	}
}

func TestFormatJSON_ErrorWithFunc(t *testing.T) {
	// Test with function (also un-marshalable)
	fn := func() {}

	_, err := FormatJSON(fn)
	if err == nil {
		t.Fatal("FormatJSON() with function: error = nil, want error")
	}
}

func TestFormatTable_NilHeaders(t *testing.T) {
	rows := [][]string{
		{"value1", "value2"},
	}

	result := FormatTable(nil, rows)

	// Should return empty string with nil headers
	if result != "" {
		t.Errorf("FormatTable() with nil headers = %q, want empty string", result)
	}
}

func TestFormatTable_EmptyHeaders(t *testing.T) {
	headers := []string{}
	rows := [][]string{
		{"value1", "value2"},
	}

	result := FormatTable(headers, rows)

	// Should return empty string with empty headers
	if result != "" {
		t.Errorf("FormatTable() with empty headers = %q, want empty string", result)
	}
}

func TestFormatTable_MismatchedRowLengths(t *testing.T) {
	headers := []string{"Col1", "Col2", "Col3"}
	rows := [][]string{
		{"A", "B", "C"},        // matches header count
		{"X", "Y"},             // fewer columns
		{"1", "2", "3", "4"},   // more columns
	}

	result := FormatTable(headers, rows)

	// Should not crash and should contain all headers
	if !strings.Contains(result, "Col1") {
		t.Errorf("result missing Col1 header")
	}
	if !strings.Contains(result, "Col2") {
		t.Errorf("result missing Col2 header")
	}
	if !strings.Contains(result, "Col3") {
		t.Errorf("result missing Col3 header")
	}

	// Should contain data from rows
	if !strings.Contains(result, "A") {
		t.Errorf("result missing data from first row")
	}
	if !strings.Contains(result, "X") {
		t.Errorf("result missing data from second row")
	}
}

func TestFormatColoredTable_EmptyRows(t *testing.T) {
	headers := []string{"Header1", "Header2"}
	rows := [][]string{}

	result := FormatColoredTable(headers, rows, true)

	// Should show headers even with no rows
	if !strings.Contains(result, "Header1") {
		t.Errorf("result missing Header1")
	}
	if !strings.Contains(result, "Header2") {
		t.Errorf("result missing Header2")
	}

	// Should contain separator line
	if !strings.Contains(result, "-") {
		t.Errorf("result missing separator line")
	}
}

func TestFormatColoredTable_NilHeaders(t *testing.T) {
	rows := [][]string{
		{"value1", "value2"},
	}

	result := FormatColoredTable(nil, rows, true)

	// Should return empty string with nil headers
	if result != "" {
		t.Errorf("FormatColoredTable() with nil headers = %q, want empty string", result)
	}
}

func TestFormatColoredTable_MismatchedRowLengths(t *testing.T) {
	headers := []string{"Col1", "Col2", "Col3"}
	rows := [][]string{
		{"A", "B", "C"},
		{"X"},                // much shorter
		{"1", "2", "3", "4", "5"},  // much longer
	}

	result := FormatColoredTable(headers, rows, false)

	// Should not crash
	if len(result) == 0 {
		t.Error("result is empty, want non-empty")
	}

	// Should contain headers
	if !strings.Contains(result, "Col1") {
		t.Errorf("result missing Col1 header")
	}
}

func TestColorizeCell_CaseInsensitive(t *testing.T) {
	// Test that colorization is case-insensitive
	tests := []struct {
		input string
		desc  string
	}{
		{"ACTIVE", "uppercase active"},
		{"Active", "mixed case active"},
		{"ERROR", "uppercase error"},
		{"Error", "mixed case error"},
		{"WARNING", "uppercase warning"},
		{"Warning", "mixed case warning"},
	}

	headers := []string{"Status"}
	for _, tt := range tests {
		rows := [][]string{{tt.input}}
		result := FormatColoredTable(headers, rows, true)

		// Should contain the input (possibly with ANSI codes)
		if !strings.Contains(result, tt.input) {
			t.Errorf("result missing %s: %s", tt.desc, tt.input)
		}
	}
}
