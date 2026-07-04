package reader

import (
	"strings"
	"testing"
)

func TestParseCSV(t *testing.T) {
	csv := `name,value
Alice,10
Bob,20
Charlie,30`

	series, err := ParseCSV(strings.NewReader(csv))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(series) != 1 {
		t.Fatalf("expected 1 series, got %d", len(series))
	}

	if series[0].Name != "value" {
		t.Errorf("expected series name 'value', got %q", series[0].Name)
	}

	if len(series[0].Points) != 3 {
		t.Fatalf("expected 3 points, got %d", len(series[0].Points))
	}

	expected := []struct {
		label string
		value float64
	}{
		{"Alice", 10},
		{"Bob", 20},
		{"Charlie", 30},
	}

	for i, exp := range expected {
		if series[0].Points[i].Label != exp.label {
			t.Errorf("point %d: expected label %q, got %q", i, exp.label, series[0].Points[i].Label)
		}
		if series[0].Points[i].Value != exp.value {
			t.Errorf("point %d: expected value %f, got %f", i, exp.value, series[0].Points[i].Value)
		}
	}
}

func TestParseCSVMultiSeries(t *testing.T) {
	csv := `month,revenue,costs
Jan,100,80
Feb,120,90
Mar,150,100`

	series, err := ParseCSV(strings.NewReader(csv))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(series) != 2 {
		t.Fatalf("expected 2 series, got %d", len(series))
	}

	if series[0].Name != "revenue" {
		t.Errorf("expected first series name 'revenue', got %q", series[0].Name)
	}
	if series[1].Name != "costs" {
		t.Errorf("expected second series name 'costs', got %q", series[1].Name)
	}
}

func TestParseCSVEmpty(t *testing.T) {
	_, err := ParseCSV(strings.NewReader(""))
	if err == nil {
		t.Error("expected error for empty CSV")
	}
}

func TestParseCSVHeaderOnly(t *testing.T) {
	csv := `name,value`
	_, err := ParseCSV(strings.NewReader(csv))
	if err == nil {
		t.Error("expected error for header-only CSV")
	}
}

func TestParseJSONSimple(t *testing.T) {
	json := `[
		{"label": "A", "value": 10},
		{"label": "B", "value": 20},
		{"label": "C", "value": 30}
	]`

	series, err := ParseJSON(strings.NewReader(json))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(series) != 1 {
		t.Fatalf("expected 1 series, got %d", len(series))
	}

	if len(series[0].Points) != 3 {
		t.Fatalf("expected 3 points, got %d", len(series[0].Points))
	}
}

func TestParseJSONMultiSeries(t *testing.T) {
	json := `[
		{"label": "Q1", "revenue": 100, "costs": 80},
		{"label": "Q2", "revenue": 150, "costs": 90}
	]`

	series, err := ParseJSON(strings.NewReader(json))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(series) != 2 {
		t.Fatalf("expected 2 series, got %d", len(series))
	}
}

func TestParseJSONArray(t *testing.T) {
	json := `[10, 20, 30, 40, 50]`

	series, err := ParseJSON(strings.NewReader(json))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(series) != 1 {
		t.Fatalf("expected 1 series, got %d", len(series))
	}

	if len(series[0].Points) != 5 {
		t.Fatalf("expected 5 points, got %d", len(series[0].Points))
	}
}

func TestParseJSONEmpty(t *testing.T) {
	_, err := ParseJSON(strings.NewReader("[]"))
	if err == nil {
		t.Error("expected error for empty JSON")
	}
}
