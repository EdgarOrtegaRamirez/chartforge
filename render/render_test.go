package render

import (
	"strings"
	"testing"

	"github.com/EdgarOrtegaRamirez/chartforge/chart"
)

func TestRenderBarTerminal(t *testing.T) {
	series := []chart.Series{
		{
			Name: "test",
			Points: []chart.DataPoint{
				{Label: "A", Value: 10},
				{Label: "B", Value: 20},
				{Label: "C", Value: 30},
			},
		},
	}

	cfg := chart.DefaultConfig()
	cfg.Type = chart.Bar
	cfg.Width = 40
	cfg.Height = 10

	result := RenderTerminal(series, cfg)
	if result == "" {
		t.Error("expected non-empty result")
	}
	if !strings.Contains(result, "A") || !strings.Contains(result, "B") || !strings.Contains(result, "C") {
		t.Error("expected labels A, B, C in output")
	}
}

func TestRenderLineTerminal(t *testing.T) {
	series := []chart.Series{
		{
			Name: "line",
			Points: []chart.DataPoint{
				{Label: "1", Value: 5},
				{Label: "2", Value: 15},
				{Label: "3", Value: 10},
			},
		},
	}

	cfg := chart.DefaultConfig()
	cfg.Type = chart.Line
	cfg.Width = 40
	cfg.Height = 10

	result := RenderTerminal(series, cfg)
	if result == "" {
		t.Error("expected non-empty result")
	}
}

func TestRenderPieTerminal(t *testing.T) {
	series := []chart.Series{
		{
			Name: "pie",
			Points: []chart.DataPoint{
				{Label: "Slice1", Value: 30},
				{Label: "Slice2", Value: 70},
			},
		},
	}

	cfg := chart.DefaultConfig()
	cfg.Type = chart.Pie
	cfg.Width = 60

	result := RenderTerminal(series, cfg)
	if result == "" {
		t.Error("expected non-empty result")
	}
	if !strings.Contains(result, "Slice1") || !strings.Contains(result, "Slice2") {
		t.Error("expected slice labels in output")
	}
}

func TestRenderScatterTerminal(t *testing.T) {
	series := []chart.Series{
		{
			Name: "scatter",
			Points: []chart.DataPoint{
				{Label: "A", Value: 10},
				{Label: "B", Value: 25},
				{Label: "C", Value: 15},
			},
		},
	}

	cfg := chart.DefaultConfig()
	cfg.Type = chart.Scatter
	cfg.Width = 40
	cfg.Height = 10

	result := RenderTerminal(series, cfg)
	if result == "" {
		t.Error("expected non-empty result")
	}
}

func TestRenderHistogramTerminal(t *testing.T) {
	series := []chart.Series{
		{
			Name: "hist",
			Points: []chart.DataPoint{
				{Label: "1", Value: 5},
				{Label: "2", Value: 15},
				{Label: "3", Value: 10},
				{Label: "4", Value: 25},
				{Label: "5", Value: 20},
			},
		},
	}

	cfg := chart.DefaultConfig()
	cfg.Type = chart.Hist
	cfg.Width = 60
	cfg.BinCount = 3

	result := RenderTerminal(series, cfg)
	if result == "" {
		t.Error("expected non-empty result")
	}
}

func TestRenderSVGBar(t *testing.T) {
	series := []chart.Series{
		{
			Name: "test",
			Points: []chart.DataPoint{
				{Label: "A", Value: 10},
				{Label: "B", Value: 20},
			},
		},
	}

	cfg := chart.DefaultConfig()
	cfg.Type = chart.Bar
	cfg.Width = 40
	cfg.Height = 10

	result := RenderSVG(series, cfg)
	if !strings.Contains(result, "<svg") {
		t.Error("expected SVG output")
	}
	if !strings.Contains(result, "Monthly") {
		// No title set, that's fine
	}
	if !strings.Contains(result, "A") || !strings.Contains(result, "B") {
		t.Error("expected labels in SVG output")
	}
}

func TestRenderSVGLine(t *testing.T) {
	series := []chart.Series{
		{
			Name: "line",
			Points: []chart.DataPoint{
				{Label: "1", Value: 5},
				{Label: "2", Value: 15},
				{Label: "3", Value: 10},
			},
		},
	}

	cfg := chart.DefaultConfig()
	cfg.Type = chart.Line
	cfg.Width = 40
	cfg.Height = 10

	result := RenderSVG(series, cfg)
	if !strings.Contains(result, "<svg") {
		t.Error("expected SVG output")
	}
	if !strings.Contains(result, "path") {
		t.Error("expected path elements in line chart SVG")
	}
}

func TestRenderSVGEmpty(t *testing.T) {
	series := []chart.Series{}
	cfg := chart.DefaultConfig()

	result := RenderSVG(series, cfg)
	if !strings.Contains(result, "No data") {
		t.Error("expected 'No data' message for empty series")
	}
}

func TestRenderTerminalEmpty(t *testing.T) {
	series := []chart.Series{}
	cfg := chart.DefaultConfig()

	result := RenderTerminal(series, cfg)
	if !strings.Contains(result, "No data") {
		t.Error("expected 'No data' message for empty series")
	}
}

func TestRenderTerminalTitle(t *testing.T) {
	series := []chart.Series{
		{
			Name: "test",
			Points: []chart.DataPoint{
				{Label: "A", Value: 10},
			},
		},
	}

	cfg := chart.DefaultConfig()
	cfg.Type = chart.Bar
	cfg.Width = 40
	cfg.Height = 10
	cfg.Title = "My Chart"

	result := RenderTerminal(series, cfg)
	if !strings.Contains(result, "My Chart") {
		t.Error("expected title 'My Chart' in output")
	}
}

func TestRenderSVGTitle(t *testing.T) {
	series := []chart.Series{
		{
			Name: "test",
			Points: []chart.DataPoint{
				{Label: "A", Value: 10},
			},
		},
	}

	cfg := chart.DefaultConfig()
	cfg.Type = chart.Bar
	cfg.Width = 40
	cfg.Height = 10
	cfg.Title = "My Chart"

	result := RenderSVG(series, cfg)
	if !strings.Contains(result, "My Chart") {
		t.Error("expected title 'My Chart' in SVG output")
	}
}

func TestEscapeXML(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "hello"},
		{"a < b", "a &lt; b"},
		{"a & b", "a &amp; b"},
		{`"quoted"`, "&quot;quoted&quot;"},
	}

	for _, tt := range tests {
		result := escapeXML(tt.input)
		if result != tt.expected {
			t.Errorf("escapeXML(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}
