# AGENTS.md — ChartForge

## Project Overview

ChartForge is a Go CLI tool that generates charts (bar, line, pie, scatter, histogram) from CSV and JSON data, outputting to the terminal (Unicode blocks + ANSI) or as SVG files.

## Build & Test

```bash
# Build
cd /root/workspace/chartforge
go build -o chartforge .

# Run all tests
go test ./...

# Run specific package tests
go test ./chart/
go test ./reader/
go test ./render/

# Lint
go vet ./...
```

## Architecture

- **`chart/`** — Core types (`ChartType`, `DataPoint`, `Series`, `ChartConfig`), validation, utility functions (WrapText, FormatValue, HistogramBins)
- **`reader/`** — Input parsing: `ParseCSV(io.Reader)` and `ParseJSON(io.Reader)` returning `[]Series`
- **`render/terminal.go`** — Terminal rendering using Unicode block characters (█ ▓ ▒ ░) with ANSI 256-color palette
- **`render/svg.go`** — SVG rendering with Catppuccin Mocha dark theme (#1e1e2e background)
- **`cmd/root.go`** — Cobra CLI: `chart` subcommand (bar/line/pie/scatter/histogram) and `info` subcommand
- **`main.go`** — Entry point, calls `cmd.Execute()`

## Key Design Decisions

1. **No external deps** beyond `github.com/spf13/cobra` — everything else uses stdlib
2. **Flag shorthands**: `-W` for width, `-H` for height (lowercase `-w`/`-h` conflict with cobra builtins)
3. **Multi-series**: JSON objects with numeric array values become separate series; CSV supports single series
4. **SVG theme**: Catppuccin Mocha palette — soft pastels on dark background
5. **Terminal bars**: Use Unicode block elements (█ ▓ ▒ ░) with gradient effects for visual richness

## Common Tasks

### Add a new chart type
1. Add constant to `chart/chart.go` (e.g., `ChartRadar ChartType = "radar"`)
2. Add to `Validate()` allowed types
3. Add rendering in `render/terminal.go` and `render/svg.go`
4. Add cobra `validArgs` in `cmd/root.go`

### Modify colors
- Terminal colors: ANSI 256-color codes in `render/terminal.go` — `colorPalette` array
- SVG colors: `colors` slice in `render/svg.go` and background `#1e1e2e`

### Add new input format
1. Create parser in `reader/reader.go` (e.g., `ParseYAML`)
2. Add format detection logic in `cmd/root.go` `runChart()`
3. Add tests in `reader/reader_test.go`

## Testing Strategy

- **chart_test.go**: Unit tests for config defaults, validation, binning, formatting
- **reader_test.go**: Tests CSV/JSON parsing with edge cases (empty files, single values, multi-series)
- **render_test.go**: Tests that terminal/SVG renderers produce non-empty output for valid configs

## CI

GitHub Actions workflow at `.github/workflows/ci.yml` runs `go test ./...` and `go vet ./...` on push to main.
