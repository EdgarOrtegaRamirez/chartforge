package chart

import (
	"fmt"
	"math"
	"strings"
)

// ChartType represents the type of chart to generate.
type ChartType string

const (
	Bar     ChartType = "bar"
	Line    ChartType = "line"
	Pie     ChartType = "pie"
	Scatter ChartType = "scatter"
	Hist    ChartType = "histogram"
)

// DataPoint represents a single data point with label and value.
type DataPoint struct {
	Label string
	Value float64
}

// Series represents a collection of data points.
type Series struct {
	Name   string
	Points []DataPoint
}

// ChartConfig holds configuration for chart generation.
type ChartConfig struct {
	Type        ChartType
	Title       string
	Width       int
	Height      int
	Colors      []string
	ShowLegend  bool
	ShowGrid    bool
	ShowValues  bool
	LabelX      string
	LabelY      string
	BinCount    int // for histogram
	Precision   int // decimal precision for value labels
}

// DefaultConfig returns a ChartConfig with sensible defaults.
func DefaultConfig() ChartConfig {
	return ChartConfig{
		Type:       Bar,
		Width:      80,
		Height:     20,
		Colors:     []string{"#4FC3F7", "#81C784", "#FFB74D", "#E57373", "#BA68C8", "#4DD0E1", "#AED581", "#FF8A65"},
		ShowLegend: true,
		ShowGrid:   true,
		ShowValues: true,
		Precision:  1,
		BinCount:   10,
	}
}

// Validate checks the chart configuration for validity.
func (c *ChartConfig) Validate() error {
	validTypes := map[ChartType]bool{Bar: true, Line: true, Pie: true, Scatter: true, Hist: true}
	if !validTypes[c.Type] {
		return fmt.Errorf("unsupported chart type: %s (valid: bar, line, pie, scatter, histogram)", c.Type)
	}
	if c.Width < 20 || c.Width > 500 {
		return fmt.Errorf("width must be between 20 and 500, got %d", c.Width)
	}
	if c.Height < 5 || c.Height > 200 {
		return fmt.Errorf("height must be between 5 and 200, got %d", c.Height)
	}
	if c.Precision < 0 || c.Precision > 6 {
		return fmt.Errorf("precision must be between 0 and 6, got %d", c.Precision)
	}
	return nil
}

// ValueRange returns the min and max values across all series.
func ValueRange(series []Series) (float64, float64) {
	minVal := math.MaxFloat64
	maxVal := math.SmallestNonzeroFloat64
	for _, s := range series {
		for _, p := range s.Points {
			if p.Value < minVal {
				minVal = p.Value
			}
			if p.Value > maxVal {
				maxVal = p.Value
			}
		}
	}
	if minVal == math.MaxFloat64 {
		return 0, 0
	}
	return minVal, maxVal
}

// AutoBin calculates optimal bin edges for a histogram.
func AutoBin(values []float64, binCount int) []float64 {
	if len(values) == 0 || binCount <= 0 {
		return nil
	}
	minVal := math.MaxFloat64
	maxVal := math.SmallestNonzeroFloat64
	for _, v := range values {
		if v < minVal {
			minVal = v
		}
		if v > maxVal {
			maxVal = v
		}
	}
	if minVal == maxVal {
		return []float64{minVal - 0.5, minVal + 0.5}
	}
	step := (maxVal - minVal) / float64(binCount)
	edges := make([]float64, binCount+1)
	for i := 0; i <= binCount; i++ {
		edges[i] = minVal + float64(i)*step
	}
	return edges
}

// CountInBins counts how many values fall into each bin.
func CountInBins(values []float64, edges []float64) []int {
	counts := make([]int, len(edges)-1)
	for _, v := range values {
		for i := 0; i < len(edges)-1; i++ {
			if v >= edges[i] && (i == len(edges)-2 || v < edges[i+1]) {
				counts[i]++
				break
			}
		}
	}
	return counts
}

// FormatValue formats a float64 with the given precision.
func FormatValue(v float64, precision int) string {
	return fmt.Sprintf("%.*f", precision, v)
}

// WrapText wraps text to fit within maxWidth, returning multiple lines.
func WrapText(text string, maxWidth int) []string {
	if maxWidth <= 0 {
		return []string{text}
	}
	if len(text) <= maxWidth {
		return []string{text}
	}
	var lines []string
	for len(text) > maxWidth {
		idx := strings.LastIndex(text[:maxWidth], " ")
		if idx <= 0 {
			idx = maxWidth
		}
		lines = append(lines, text[:idx])
		text = strings.TrimSpace(text[idx:])
	}
	if len(text) > 0 {
		lines = append(lines, text)
	}
	return lines
}
