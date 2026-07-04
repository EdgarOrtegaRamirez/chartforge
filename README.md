# ChartForge 🔨

A fast, beautiful CLI tool that generates charts from CSV and JSON data — right in your terminal or as SVG files.

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![CI](https://github.com/EdgarOrtegaRamirez/chartforge/actions/workflows/ci.yml/badge.svg)](https://github.com/EdgarOrtegaRamirez/chartforge/actions/workflows/ci.yml)

## Features

- **5 chart types**: bar, line, pie, scatter, histogram
- **2 output modes**: terminal (with Unicode blocks & ANSI colors) and SVG
- **Multi-series support**: plot multiple datasets from a single JSON file
- **Catppuccin Mocha theme**: gorgeous dark-themed SVG output
- **Zero dependencies** beyond the Go standard library + cobra
- **Fast**: single static binary, instant rendering

## Quick Start

```bash
# Install
go install github.com/EdgarOrtegaRamirez/chartforge@latest

# Or build from source
git clone https://github.com/EdgarOrtegaRamirez/chartforge.git
cd chartforge
go build -o chartforge .
```

### Terminal Chart

```bash
# Create some data
echo -e "Month,Sales\nJan,100\nFeb,150\nMar,200\nApr,175\nMay,300\nJun,250" > data.csv

# Render in terminal
./chartforge chart bar data.csv --title "Monthly Sales"

# Line chart
./chartforge chart line data.csv --title "Sales Trend"

# Pie chart
./chartforge chart pie data.csv --title "Revenue Distribution"

# Scatter plot
./chartforge chart scatter data.csv --title "Sales vs Month"

# Histogram (auto-binned)
./chartforge chart histogram data.csv --bins 5 --title "Sales Distribution"
```

### SVG Output

```bash
# Generate SVG file
./chartforge chart bar data.csv --title "Monthly Sales" -o chart.svg

# Custom dimensions
./chartforge chart line data.csv -W 1200 -H 600 -o trend.svg
```

### Multi-Series (JSON)

```bash
cat > financials.json << 'EOF'
{
  "profit": [10, 15, 13, 17, 20],
  "revenue": [50, 65, 58, 72, 80],
  "costs": [40, 50, 45, 55, 60]
}
EOF

./chartforge chart line financials.json --title "Quarterly Financials" -l
./chartforge chart bar financials.json --title "Financial Overview" -l -v
```

### Data Inspection

```bash
./chartforge info data.csv
```

## CLI Reference

```
chartforge chart [type] [file] [flags]
chartforge info [file]
```

### Chart Types

| Type | Description |
|------|-------------|
| `bar` | Vertical bar chart with Unicode block characters |
| `line` | Line chart with data point markers |
| `pie` | Pie chart with percentage labels |
| `scatter` | Scatter plot with markers |
| `histogram` | Frequency distribution with auto or manual binning |

### Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--title` | `-t` | `"Chart"` | Chart title |
| `--width` | `-W` | `60` | Chart width (terminal columns or SVG px) |
| `--height` | `-H` | `20` | Chart height (terminal rows or SVG px) |
| `--output` | `-o` | (none) | Output SVG file path |
| `--terminal` | `-T` | `false` | Force terminal output |
| `--bins` | `-b` | `10` | Number of histogram bins |
| `--precision` | `-p` | `1` | Decimal places for values |
| `--legend` | `-l` | `false` | Show legend for multi-series |
| `--values` | `-v` | `false` | Show values on bars/points |

## Input Formats

### CSV

```csv
Category,Value
Widgets,150
Gadgets,200
Gizmos,100
Doohickeys,75
```

- First column: labels
- Remaining columns: numeric values (first column used by default for single-series)

### JSON

**Single series** (array):
```json
[10, 20, 30, 25, 15]
```

**Multi-series** (object):
```json
{
  "sales": [100, 150, 200],
  "returns": [10, 15, 8]
}
```

## Examples

### Bar Chart (Terminal)
```
  ┌─ Monthly Sales ──────────────────────────────────────┐
  │                                        ████████ 300  │
  │                               ████████ 200           │
  │                      ████████ 150                    │
  │             ████████                                 │
  │    ████████ 100                                      │
  │                                        ████████ 250  │
  │             ████████ 175                             │
  └──────────────────────────────────────────────────────┘
     Jan   Feb   Mar   Apr   May   Jun
```

### SVG Output
Generates a dark-themed SVG with the Catppuccin Mocha color palette, perfect for embedding in documentation, presentations, or websites.

## Architecture

```
chartforge/
├── main.go              # Entry point
├── cmd/
│   └── root.go          # Cobra CLI commands & flags
├── chart/
│   └── chart.go         # Core types, validation, utilities
├── reader/
│   └── reader.go        # CSV & JSON parsing
├── render/
│   ├── terminal.go      # Terminal rendering (Unicode + ANSI)
│   └── svg.go           # SVG rendering (Catppuccin theme)
└── tests/
    └── testdata/        # Test fixtures
```

## Contributing

1. Fork the repo
2. Create a feature branch (`git checkout -b feat/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feat/amazing-feature`)
5. Open a Pull Request

## License

[MIT](LICENSE) © 2025 Edgar Ortega Ramirez
