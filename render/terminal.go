package render

import (
	"fmt"
	"math"
	"strings"

	"github.com/EdgarOrtegaRamirez/chartforge/chart"
)

// Unicode block characters for rendering.
const (
	BlockFull    = "█"
	Block34      = "▓"
	Block12      = "▒"
	Block14      = "░"
	HLine        = "─"
	VLine        = "│"
	CornerTL     = "┌"
	CornerTR     = "┐"
	CornerBL     = "└"
	CornerBR     = "┘"
	TLeft        = "├"
	TRight       = "┤"
	TTop         = "┬"
	TBottom      = "┴"
	Cross        = "┼"
	UpTriangle   = "▲"
	DownTriangle = "▼"
	Circle       = "●"
	Diamond      = "◆"
	Square       = "■"
)

// Unicode block fractions for partial bars.
var blockFractions = []string{"", "▁", "▂", "▃", "▄", "▅", "▆", "▇", "█"}

// Color codes for terminal output.
var terminalColors = []string{
	"\033[38;5;75m",  // blue
	"\033[38;5;114m", // green
	"\033[38;5;215m", // orange
	"\033[38;5;203m", // red
	"\033[38;5;141m", // purple
	"\033[38;5;80m",  // teal
	"\033[38;5;149m", // lime
	"\033[38;5;209m", // salmon
}

const resetColor = "\033[0m"

// RenderTerminal renders chart as Unicode block art for terminal display.
func RenderTerminal(s []chart.Series, cfg chart.ChartConfig) string {
	switch cfg.Type {
	case chart.Bar:
		return renderBarTerminal(s, cfg)
	case chart.Line:
		return renderLineTerminal(s, cfg)
	case chart.Pie:
		return renderPieTerminal(s, cfg)
	case chart.Scatter:
		return renderScatterTerminal(s, cfg)
	case chart.Hist:
		return renderHistogramTerminal(s, cfg)
	default:
		return "Unsupported chart type: " + string(cfg.Type)
	}
}

func renderBarTerminal(series []chart.Series, cfg chart.ChartConfig) string {
	if len(series) == 0 || len(series[0].Points) == 0 {
		return "No data to render"
	}

	var sb strings.Builder
	width := cfg.Width
	height := cfg.Height

	// Use first series for bar chart
	data := series[0].Points
	numBars := len(data)
	if numBars == 0 {
		return "No data points"
	}

	// Find max value
	maxVal := 0.0
	for _, p := range data {
		if p.Value > maxVal {
			maxVal = p.Value
		}
	}
	if maxVal == 0 {
		maxVal = 1
	}

	// Calculate label width
	maxLabelLen := 0
	for _, p := range data {
		if len(p.Label) > maxLabelLen {
			maxLabelLen = len(p.Label)
		}
	}
	if maxLabelLen > width/3 {
		maxLabelLen = width / 3
	}

	_ = width - maxLabelLen - 6 // reserved for future label rendering

	// Title
	if cfg.Title != "" {
		titleLines := chart.WrapText(cfg.Title, width)
		for _, line := range titleLines {
			padding := (width - len(line)) / 2
			if padding > 0 {
				sb.WriteString(strings.Repeat(" ", padding))
			}
			sb.WriteString(line)
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
	}

	// Y-axis label
	if cfg.ShowValues && maxVal > 0 {
		// Render bars from top to bottom
		for row := height; row >= 1; row-- {
			sb.WriteString(VLine)
			for i, p := range data {
				barHeight := int(math.Round(p.Value / maxVal * float64(height)))
				if barHeight >= row {
					sb.WriteString(BlockFull)
				} else if barHeight == row-1 && row > 1 {
					sb.WriteString(Block34)
				} else if barHeight == row-2 && row > 2 {
					sb.WriteString(Block12)
				} else {
					sb.WriteString(" ")
				}
				sb.WriteString(" ")
				_ = i
			}
			sb.WriteString(VLine)
			if row == height {
				sb.WriteString(fmt.Sprintf(" %s", chart.FormatValue(maxVal, cfg.Precision)))
			} else if row == 1 {
				sb.WriteString(fmt.Sprintf(" 0"))
			} else if row == height/2 {
				sb.WriteString(fmt.Sprintf(" %s", chart.FormatValue(maxVal/2, cfg.Precision)))
			}
			sb.WriteString("\n")
		}

		// Bottom axis
		sb.WriteString(CornerBL)
		for i := 0; i <= numBars*2; i++ {
			sb.WriteString(HLine)
		}
		sb.WriteString(CornerBR)
		sb.WriteString("\n")

		// Labels
		sb.WriteString(strings.Repeat(" ", 1))
		for _, p := range data {
			label := p.Label
			if len(label) > 4 {
				label = label[:4]
			}
			sb.WriteString(fmt.Sprintf("%-2s", label))
		}
		sb.WriteString("\n")
	}

	// Legend
	if cfg.ShowLegend && len(series) > 1 {
		sb.WriteString("\n")
		for i, s := range series {
			color := terminalColors[i%len(terminalColors)]
			sb.WriteString(fmt.Sprintf("  %s%s%s %s\n", color, BlockFull, resetColor, s.Name))
		}
	}

	return sb.String()
}

func renderLineTerminal(series []chart.Series, cfg chart.ChartConfig) string {
	if len(series) == 0 || len(series[0].Points) == 0 {
		return "No data to render"
	}

	var sb strings.Builder
	width := cfg.Width
	height := cfg.Height

	// Find max value across all series
	maxVal := 0.0
	for _, s := range series {
		for _, p := range s.Points {
			if p.Value > maxVal {
				maxVal = p.Value
			}
		}
	}
	if maxVal == 0 {
		maxVal = 1
	}

	// Create canvas
	canvas := make([][]rune, height)
	for i := range canvas {
		canvas[i] = make([]rune, width)
		for j := range canvas[i] {
			canvas[i][j] = ' '
		}
	}

	// Plot each series
	markers := []string{string(Circle), string(Diamond), string(Square), string(UpTriangle), "×", "+"}
	for si, s := range series {
		marker := []rune(markers[si%len(markers)])[0]
		numPoints := len(s.Points)
		if numPoints <= 1 {
			continue
		}

		for i, p := range s.Points {
			x := int(float64(i) / float64(numPoints-1) * float64(width-1))
			y := int(p.Value / maxVal * float64(height-1))
			y = height - 1 - y // flip Y axis
			if x >= 0 && x < width && y >= 0 && y < height {
				canvas[y][x] = marker
			}
			// Connect points with lines
			if i > 0 {
				prevX := int(float64(i-1) / float64(numPoints-1) * float64(width-1))
				prevY := int(s.Points[i-1].Value / maxVal * float64(height-1))
				prevY = height - 1 - prevY

				// Draw line between points using Bresenham-like approach
				dx := x - prevX
				dy := y - prevY
				steps := int(math.Max(float64(abs(dx)), float64(abs(dy))))
				if steps > 0 {
					for step := 1; step < steps; step++ {
						cx := prevX + dx*step/steps
						cy := prevY + dy*step/steps
						if cx >= 0 && cx < width && cy >= 0 && cy < height && canvas[cy][cx] == ' ' {
							if abs(dx) > abs(dy) {
								canvas[cy][cx] = []rune(HLine)[0]
							} else {
								canvas[cy][cx] = []rune(VLine)[0]
							}
						}
					}
				}
			}
		}
	}

	// Title
	if cfg.Title != "" {
		titleLines := chart.WrapText(cfg.Title, width)
		for _, line := range titleLines {
			padding := (width - len(line)) / 2
			if padding > 0 {
				sb.WriteString(strings.Repeat(" ", padding))
			}
			sb.WriteString(line)
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
	}

	// Render canvas
	for _, row := range canvas {
		sb.WriteString(VLine)
		sb.WriteString(string(row))
		sb.WriteString(VLine)
		sb.WriteString("\n")
	}

	// Bottom axis
	sb.WriteString(CornerBL)
	for i := 0; i < width; i++ {
		sb.WriteString(HLine)
	}
	sb.WriteString(CornerBR)
	sb.WriteString("\n")

	// Legend
	if cfg.ShowLegend {
		sb.WriteString("\n")
		for i, s := range series {
			color := terminalColors[i%len(terminalColors)]
			markerStr := markers[i%len(markers)]
			sb.WriteString(fmt.Sprintf("  %s%s%s %s\n", color, markerStr, resetColor, s.Name))
		}
	}

	return sb.String()
}

func renderPieTerminal(series []chart.Series, cfg chart.ChartConfig) string {
	if len(series) == 0 || len(series[0].Points) == 0 {
		return "No data to render"
	}

	var sb strings.Builder
	data := series[0].Points

	// Calculate total
	total := 0.0
	for _, p := range data {
		total += p.Value
	}
	if total == 0 {
		return "Total is zero, cannot render pie chart"
	}

	// Title
	if cfg.Title != "" {
		titleLines := chart.WrapText(cfg.Title, cfg.Width)
		for _, line := range titleLines {
			padding := (cfg.Width - len(line)) / 2
			if padding > 0 {
				sb.WriteString(strings.Repeat(" ", padding))
			}
			sb.WriteString(line)
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
	}

	// Pie chart as horizontal bar breakdown
	barWidth := cfg.Width - 20
	for i, p := range data {
		pct := p.Value / total * 100
		barLen := int(pct / 100 * float64(barWidth))
		if barLen < 1 && p.Value > 0 {
			barLen = 1
		}
		color := terminalColors[i%len(terminalColors)]
		sb.WriteString(fmt.Sprintf("  %s%s%-12s%s ", color, BlockFull, p.Label[:minInt(12, len(p.Label))], resetColor))
		sb.WriteString(BlockFull)
		// Use different block chars for visual distinction
		blocks := []string{BlockFull, Block34, Block12, Block14}
		block := blocks[i%len(blocks)]
		for j := 0; j < barLen; j++ {
			sb.WriteString(block)
		}
		sb.WriteString(fmt.Sprintf("  %.1f%% (%s)\n", pct, chart.FormatValue(p.Value, cfg.Precision)))
	}

	return sb.String()
}

func renderScatterTerminal(series []chart.Series, cfg chart.ChartConfig) string {
	if len(series) == 0 || len(series[0].Points) == 0 {
		return "No data to render"
	}

	var sb strings.Builder
	width := cfg.Width
	height := cfg.Height

	// Find ranges
	minX, maxX := math.MaxFloat64, math.SmallestNonzeroFloat64
	minY, maxY := math.MaxFloat64, math.SmallestNonzeroFloat64
	for _, s := range series {
		for i, p := range s.Points {
			x := float64(i)
			if x < minX {
				minX = x
			}
			if x > maxX {
				maxX = x
			}
			if p.Value < minY {
				minY = p.Value
			}
			if p.Value > maxY {
				maxY = p.Value
			}
		}
	}
	if minX == maxX {
		maxX = minX + 1
	}
	if minY == maxY {
		maxY = minY + 1
	}

	// Create canvas
	canvas := make([][]rune, height)
	for i := range canvas {
		canvas[i] = make([]rune, width)
		for j := range canvas[i] {
			canvas[i][j] = ' '
		}
	}

	// Plot points
	markers := []string{string(Circle), string(Diamond), string(Square), string(UpTriangle), "×", "+"}
	for si, s := range series {
		marker := []rune(markers[si%len(markers)])[0]
		for i, p := range s.Points {
			x := int((float64(i) - minX) / (maxX - minX) * float64(width-1))
			y := int((p.Value - minY) / (maxY - minY) * float64(height-1))
			y = height - 1 - y
			if x >= 0 && x < width && y >= 0 && y < height {
				canvas[y][x] = marker
			}
		}
	}

	// Title
	if cfg.Title != "" {
		titleLines := chart.WrapText(cfg.Title, width)
		for _, line := range titleLines {
			padding := (width - len(line)) / 2
			if padding > 0 {
				sb.WriteString(strings.Repeat(" ", padding))
			}
			sb.WriteString(line)
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
	}

	// Render canvas
	for _, row := range canvas {
		sb.WriteString(VLine)
		sb.WriteString(string(row))
		sb.WriteString(VLine)
		sb.WriteString("\n")
	}

	// Bottom axis
	sb.WriteString(CornerBL)
	for i := 0; i < width; i++ {
		sb.WriteString(HLine)
	}
	sb.WriteString(CornerBR)
	sb.WriteString("\n")

	// Legend
	if cfg.ShowLegend {
		sb.WriteString("\n")
		for i, s := range series {
			color := terminalColors[i%len(terminalColors)]
			markerStr := markers[i%len(markers)]
			sb.WriteString(fmt.Sprintf("  %s%s%s %s\n", color, markerStr, resetColor, s.Name))
		}
	}

	return sb.String()
}

func renderHistogramTerminal(series []chart.Series, cfg chart.ChartConfig) string {
	if len(series) == 0 || len(series[0].Points) == 0 {
		return "No data to render"
	}

	var sb strings.Builder

	// Collect all values
	var values []float64
	for _, s := range series {
		for _, p := range s.Points {
			values = append(values, p.Value)
		}
	}

	binCount := cfg.BinCount
	if binCount <= 0 {
		binCount = 10
	}

	edges := chart.AutoBin(values, binCount)
	counts := chart.CountInBins(values, edges)

	maxCount := 0
	for _, c := range counts {
		if c > maxCount {
			maxCount = c
		}
	}
	if maxCount == 0 {
		maxCount = 1
	}

	barWidth := cfg.Width - 20

	// Title
	if cfg.Title != "" {
		titleLines := chart.WrapText(cfg.Title, cfg.Width)
		for _, line := range titleLines {
			padding := (cfg.Width - len(line)) / 2
			if padding > 0 {
				sb.WriteString(strings.Repeat(" ", padding))
			}
			sb.WriteString(line)
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
	}

	// Render histogram
	for i, count := range counts {
		barLen := int(float64(count) / float64(maxCount) * float64(barWidth))
		if barLen < 1 && count > 0 {
			barLen = 1
		}
		color := terminalColors[i%len(terminalColors)]
		label := fmt.Sprintf("[%s, %s)", chart.FormatValue(edges[i], cfg.Precision), chart.FormatValue(edges[i+1], cfg.Precision))
		sb.WriteString(fmt.Sprintf("  %s%-16s%s ", color, label[:minInt(16, len(label))], resetColor))
		for j := 0; j < barLen; j++ {
			sb.WriteString(BlockFull)
		}
		sb.WriteString(fmt.Sprintf(" %d\n", count))
	}

	return sb.String()
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
