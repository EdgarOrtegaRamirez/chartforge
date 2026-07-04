package reader

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/EdgarOrtegaRamirez/chartforge/chart"
)

// ReadCSV reads a CSV file and returns series suitable for charting.
// First column is used as labels, subsequent columns as separate series.
func ReadCSV(filename string) ([]chart.Series, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}
	defer f.Close()
	return ParseCSV(f)
}

// ParseCSV reads CSV data from an io.Reader and returns series.
func ParseCSV(r io.Reader) ([]chart.Series, error) {
	reader := csv.NewReader(r)
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("reading CSV: %w", err)
	}
	if len(records) < 2 {
		return nil, fmt.Errorf("CSV must have at least 2 rows (header + data), got %d", len(records))
	}

	header := records[0]
	if len(header) < 2 {
		return nil, fmt.Errorf("CSV must have at least 2 columns (labels + values), got %d", len(header))
	}

	// Create a series for each value column (index 1+)
	series := make([]chart.Series, len(header)-1)
	for i := 1; i < len(header); i++ {
		series[i-1] = chart.Series{
			Name:   strings.TrimSpace(header[i]),
			Points: make([]chart.DataPoint, 0, len(records)-1),
		}
	}

	for rowIdx, row := range records[1:] {
		if len(row) == 0 {
			continue
		}
		label := strings.TrimSpace(row[0])
		for colIdx := 1; colIdx < len(row) && colIdx < len(header); colIdx++ {
			val, err := parseFloat(row[colIdx])
			if err != nil {
				// Skip non-numeric values gracefully
				continue
			}
			series[colIdx-1].Points = append(series[colIdx-1].Points, chart.DataPoint{
				Label: label,
				Value: val,
			})
		}
		_ = rowIdx
	}

	// Filter out empty series
	var result []chart.Series
	for _, s := range series {
		if len(s.Points) > 0 {
			result = append(result, s)
		}
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no valid numeric data found in CSV")
	}

	return result, nil
}

// JSONRecord represents a flexible JSON data structure.
type JSONRecord struct {
	Label string  `json:"label"`
	Value float64 `json:"value"`
	Name  string  `json:"name"`
	// Support for flat arrays of numbers
}

// ReadJSON reads a JSON file and returns series suitable for charting.
// Supports formats:
// - [{"label": "A", "value": 10}, ...]
// - [{"label": "A", "series1": 10, "series2": 20}, ...]
// - {"series1": [{"label": "A", "value": 10}, ...], ...}
// - [10, 20, 30, ...] (simple array)
func ReadJSON(filename string) ([]chart.Series, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}
	defer f.Close()
	return ParseJSON(f)
}

// ParseJSON reads JSON data from an io.Reader and returns series.
func ParseJSON(r io.Reader) ([]chart.Series, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("reading JSON: %w", err)
	}

	// Try simple array of numbers
	var nums []float64
	if err := json.Unmarshal(data, &nums); err == nil && len(nums) > 0 {
		points := make([]chart.DataPoint, len(nums))
		for i, v := range nums {
			points[i] = chart.DataPoint{
				Label: strconv.Itoa(i + 1),
				Value: v,
			}
		}
		return []chart.Series{{Name: "values", Points: points}}, nil
	}

	// Try array of objects with label/value
	var records []map[string]interface{}
	if err := json.Unmarshal(data, &records); err != nil {
		return nil, fmt.Errorf("parsing JSON: %w", err)
	}
	if len(records) == 0 {
		return nil, fmt.Errorf("empty JSON array")
	}

	// Check if it's a map of series
	first := records[0]
	if _, ok := first["label"]; !ok {
		// Check if first record has "label" key at all
		// Try to interpret as {"seriesName": [...]}
		var seriesMap map[string]interface{}
		if err := json.Unmarshal(data, &seriesMap); err == nil {
			return parseSeriesMap(seriesMap)
		}
	}

	// Detect column names from all records
	colSet := make(map[string]bool)
	for _, rec := range records {
		for k := range rec {
			if k != "label" {
				colSet[k] = true
			}
		}
	}

	// If there's only a "value" key, simple series
	if _, hasValue := first["value"]; hasValue && len(colSet) == 1 {
		points := make([]chart.DataPoint, 0, len(records))
		for _, rec := range records {
			label := fmt.Sprintf("%v", rec["label"])
			if v, ok := rec["value"].(float64); ok {
				points = append(points, chart.DataPoint{Label: label, Value: v})
			}
		}
		return []chart.Series{{Name: "values", Points: points}}, nil
	}

	// Multiple value columns = multiple series
	colNames := make([]string, 0, len(colSet))
	for k := range colSet {
		colNames = append(colNames, k)
	}

	series := make([]chart.Series, len(colNames))
	for i, col := range colNames {
		series[i] = chart.Series{
			Name:   col,
			Points: make([]chart.DataPoint, 0, len(records)),
		}
		for _, rec := range records {
			label := fmt.Sprintf("%v", rec["label"])
			if v, ok := rec[col].(float64); ok {
				series[i].Points = append(series[i].Points, chart.DataPoint{Label: label, Value: v})
			}
		}
	}

	if len(series) == 0 {
		return nil, fmt.Errorf("no valid data found in JSON")
	}

	return series, nil
}

func parseSeriesMap(data map[string]interface{}) ([]chart.Series, error) {
	var result []chart.Series
	for name, raw := range data {
		arr, ok := raw.([]interface{})
		if !ok {
			continue
		}
		points := make([]chart.DataPoint, 0, len(arr))
		for _, item := range arr {
			switch v := item.(type) {
			case float64:
				idx := len(points) + 1
				points = append(points, chart.DataPoint{
					Label: strconv.Itoa(idx),
					Value: v,
				})
			case map[string]interface{}:
				label := fmt.Sprintf("%v", v["label"])
				val, _ := v["value"].(float64)
				points = append(points, chart.DataPoint{Label: label, Value: val})
			}
		}
		if len(points) > 0 {
			result = append(result, chart.Series{Name: name, Points: points})
		}
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("no valid series found in JSON map")
	}
	return result, nil
}

func parseFloat(s string) (float64, error) {
	s = strings.TrimSpace(s)
	if s == "" || s == "-" || s == "N/A" || s == "null" || s == "undefined" {
		return 0, fmt.Errorf("non-numeric value: %s", s)
	}
	return strconv.ParseFloat(s, 64)
}
