package chart

import (
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Type != Bar {
		t.Errorf("expected Bar chart type, got %v", cfg.Type)
	}
	if cfg.Width != 80 {
		t.Errorf("expected width 80, got %d", cfg.Width)
	}
	if cfg.Height != 20 {
		t.Errorf("expected height 20, got %d", cfg.Height)
	}
	if cfg.Precision != 1 {
		t.Errorf("expected precision 1, got %d", cfg.Precision)
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     ChartConfig
		wantErr bool
	}{
		{"valid bar", ChartConfig{Type: Bar, Width: 80, Height: 20, Precision: 1}, false},
		{"valid line", ChartConfig{Type: Line, Width: 100, Height: 30, Precision: 2}, false},
		{"valid pie", ChartConfig{Type: Pie, Width: 60, Height: 15, Precision: 0}, false},
		{"invalid type", ChartConfig{Type: "invalid", Width: 80, Height: 20, Precision: 1}, true},
		{"width too small", ChartConfig{Type: Bar, Width: 10, Height: 20, Precision: 1}, true},
		{"width too large", ChartConfig{Type: Bar, Width: 600, Height: 20, Precision: 1}, true},
		{"height too small", ChartConfig{Type: Bar, Width: 80, Height: 2, Precision: 1}, true},
		{"height too large", ChartConfig{Type: Bar, Width: 80, Height: 300, Precision: 1}, true},
		{"precision too low", ChartConfig{Type: Bar, Width: 80, Height: 20, Precision: -1}, true},
		{"precision too high", ChartConfig{Type: Bar, Width: 80, Height: 20, Precision: 10}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValueRange(t *testing.T) {
	series := []Series{
		{
			Name: "test",
			Points: []DataPoint{
				{Label: "a", Value: 10},
				{Label: "b", Value: 50},
				{Label: "c", Value: 25},
			},
		},
	}

	minVal, maxVal := ValueRange(series)
	if minVal != 10 {
		t.Errorf("expected min 10, got %f", minVal)
	}
	if maxVal != 50 {
		t.Errorf("expected max 50, got %f", maxVal)
	}
}

func TestValueRangeEmpty(t *testing.T) {
	minVal, maxVal := ValueRange(nil)
	if minVal != 0 || maxVal != 0 {
		t.Errorf("expected 0, 0 for empty series, got %f, %f", minVal, maxVal)
	}
}

func TestAutoBin(t *testing.T) {
	values := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	edges := AutoBin(values, 5)
	if len(edges) != 6 {
		t.Errorf("expected 6 edges, got %d", len(edges))
	}
	if edges[0] != 1 {
		t.Errorf("expected first edge 1, got %f", edges[0])
	}
	if edges[5] != 10 {
		t.Errorf("expected last edge 10, got %f", edges[5])
	}
}

func TestAutoBinEmpty(t *testing.T) {
	edges := AutoBin(nil, 5)
	if edges != nil {
		t.Errorf("expected nil for empty values, got %v", edges)
	}
}

func TestCountInBins(t *testing.T) {
	values := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	edges := []float64{1, 3, 5, 7, 10}
	counts := CountInBins(values, edges)

	expected := []int{2, 2, 2, 4} // [1,3): 1,2; [3,5): 3,4; [5,7): 5,6; [7,10]: 7,8,9,10
	if len(counts) != len(expected) {
		t.Fatalf("expected %d bins, got %d", len(expected), len(counts))
	}
	for i, c := range counts {
		if c != expected[i] {
			t.Errorf("bin %d: expected %d, got %d", i, expected[i], c)
		}
	}
}

func TestFormatValue(t *testing.T) {
	tests := []struct {
		value     float64
		precision int
		expected  string
	}{
		{3.14159, 2, "3.14"},
		{3.14159, 0, "3"},
		{42.0, 1, "42.0"},
		{0.123456, 4, "0.1235"},
	}

	for _, tt := range tests {
		result := FormatValue(tt.value, tt.precision)
		if result != tt.expected {
			t.Errorf("FormatValue(%f, %d) = %s, want %s", tt.value, tt.precision, result, tt.expected)
		}
	}
}

func TestWrapText(t *testing.T) {
	tests := []struct {
		text     string
		maxWidth int
		expected int // number of lines
	}{
		{"short", 10, 1},
		{"this is a longer text that should be wrapped", 10, 5},
		{"", 10, 1},
		{"word", 0, 1},
	}

	for _, tt := range tests {
		lines := WrapText(tt.text, tt.maxWidth)
		if len(lines) != tt.expected {
			t.Errorf("WrapText(%q, %d) returned %d lines, want %d", tt.text, tt.maxWidth, len(lines), tt.expected)
		}
	}
}
