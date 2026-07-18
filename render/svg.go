package render

import (
	"fmt"
	"math"
	"strings"

	"github.com/EdgarOrtegaRamirez/chartforge/chart"
)

// SVGColors for chart series.
var SVGColors = []string{
	"#4FC3F7", "#81C784", "#FFB74D", "#E57373",
	"#BA68C8", "#4DD0E1", "#AED581", "#FF8A65",
	"#7986CB", "#F06292", "#4DB6AC", "#DCE775",
}

// RenderSVG generates an SVG chart.
func RenderSVG(s []chart.Series, cfg chart.ChartConfig) string {
	switch cfg.Type {
	case chart.Bar:
		return renderBarSVG(s, cfg)
	case chart.Line:
		return renderLineSVG(s, cfg)
	case chart.Pie:
		return renderPieSVG(s, cfg)
	case chart.Scatter:
		return renderScatterSVG(s, cfg)
	case chart.Hist:
		return renderHistogramSVG(s, cfg)
	default:
		return "<!-- Unsupported chart type -->"
	}
}

func renderBarSVG(series []chart.Series, cfg chart.ChartConfig) string {
	if len(series) == 0 || len(series[0].Points) == 0 {
		return "<!-- No data -->"
	}

	data := series[0].Points
	padding := SVGPadding{Top: 40, Right: 20, Bottom: 60, Left: 60}
	width := cfg.Width * 8 // scale up for SVG
	height := cfg.Height * 12

	chartW := width - padding.Left - padding.Right
	chartH := height - padding.Top - padding.Bottom

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

	numBars := len(data)
	barWidth := chartW / numBars
	if barWidth < 2 {
		barWidth = 2
	}
	barGap := barWidth / 4
	if barGap > 4 {
		barGap = 4
	}
	actualBarWidth := barWidth - barGap

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`, width, height, width, height))
	sb.WriteString(`<style>`)
	sb.WriteString(`text { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; }`)
	sb.WriteString(`</style>`)

	// Background
	sb.WriteString(fmt.Sprintf(`<rect width="%d" height="%d" fill="#1e1e2e" rx="8"/>`, width, height))

	// Title
	if cfg.Title != "" {
		sb.WriteString(fmt.Sprintf(`<text x="%d" y="25" text-anchor="middle" fill="#cdd6f4" font-size="16" font-weight="bold">%s</text>`, width/2, escapeXML(cfg.Title)))
	}

	// Y-axis gridlines
	for i := 0; i <= 4; i++ {
		y := padding.Top + chartH - (chartH * i / 4)
		val := maxVal * float64(i) / 4
		sb.WriteString(fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#313244" stroke-width="1"/>`, padding.Left, y, width-padding.Right, y))
		sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" text-anchor="end" fill="#a6adc8" font-size="10">%s</text>`, padding.Left-5, y+4, chart.FormatValue(val, cfg.Precision)))
	}

	// Bars
	for i, p := range data {
		barH := p.Value / maxVal * float64(chartH)
		x := padding.Left + i*barWidth + barGap/2
		y := float64(padding.Top+chartH) - barH
		color := SVGColors[i%len(SVGColors)]

		sb.WriteString(fmt.Sprintf(`<rect x="%d" y="%d" width="%d" height="%d" fill="%s" rx="2" opacity="0.9">`, x, int(y), actualBarWidth, int(barH), color))
		sb.WriteString(fmt.Sprintf(`<title>%s: %s</title>`, escapeXML(p.Label), chart.FormatValue(p.Value, cfg.Precision)))
		sb.WriteString(`</rect>`)

		// Value label on top
		if cfg.ShowValues {
			sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" text-anchor="middle" fill="#cdd6f4" font-size="9">%s</text>`, x+actualBarWidth/2, int(y)-4, chart.FormatValue(p.Value, cfg.Precision)))
		}

		// X-axis label
		label := p.Label
		if len(label) > 8 {
			label = label[:8] + "…"
		}
		sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" text-anchor="middle" fill="#a6adc8" font-size="10" transform="rotate(-30 %d %d)">%s</text>`, x+actualBarWidth/2, height-padding.Bottom+15, x+actualBarWidth/2, height-padding.Bottom+15, escapeXML(label)))
	}

	// Axes
	sb.WriteString(fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#585b70" stroke-width="1.5"/>`, padding.Left, padding.Top, padding.Left, padding.Top+chartH))
	sb.WriteString(fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#585b70" stroke-width="1.5"/>`, padding.Left, padding.Top+chartH, width-padding.Right, padding.Top+chartH))

	sb.WriteString(`</svg>`)
	return sb.String()
}

func renderLineSVG(series []chart.Series, cfg chart.ChartConfig) string {
	if len(series) == 0 || len(series[0].Points) == 0 {
		return "<!-- No data -->"
	}

	padding := SVGPadding{Top: 40, Right: 20, Bottom: 60, Left: 60}
	width := cfg.Width * 8
	height := cfg.Height * 12

	chartW := width - padding.Left - padding.Right
	chartH := height - padding.Top - padding.Bottom

	// Find ranges
	maxVal := 0.0
	numPoints := 0
	for _, s := range series {
		if len(s.Points) > numPoints {
			numPoints = len(s.Points)
		}
		for _, p := range s.Points {
			if p.Value > maxVal {
				maxVal = p.Value
			}
		}
	}
	if maxVal == 0 {
		maxVal = 1
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`, width, height, width, height))
	sb.WriteString(`<style>text { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; }</style>`)
	sb.WriteString(fmt.Sprintf(`<rect width="%d" height="%d" fill="#1e1e2e" rx="8"/>`, width, height))

	// Title
	if cfg.Title != "" {
		sb.WriteString(fmt.Sprintf(`<text x="%d" y="25" text-anchor="middle" fill="#cdd6f4" font-size="16" font-weight="bold">%s</text>`, width/2, escapeXML(cfg.Title)))
	}

	// Gridlines
	for i := 0; i <= 4; i++ {
		y := padding.Top + chartH - (chartH * i / 4)
		val := maxVal * float64(i) / 4
		sb.WriteString(fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#313244" stroke-width="1"/>`, padding.Left, y, width-padding.Right, y))
		sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" text-anchor="end" fill="#a6adc8" font-size="10">%s</text>`, padding.Left-5, y+4, chart.FormatValue(val, cfg.Precision)))
	}

	// Lines and points
	for si, s := range series {
		color := SVGColors[si%len(SVGColors)]
		points := s.Points
		if len(points) < 2 {
			continue
		}

		// Build path
		var pathD strings.Builder
		for i, p := range points {
			x := float64(padding.Left) + float64(i)/float64(numPoints-1)*float64(chartW)
			y := float64(padding.Top) + float64(chartH) - p.Value/maxVal*float64(chartH)
			if i == 0 {
				pathD.WriteString(fmt.Sprintf("M %.1f %.1f", x, y))
			} else {
				pathD.WriteString(fmt.Sprintf(" L %.1f %.1f", x, y))
			}
		}

		// Area fill (semi-transparent)
		areaPath := pathD.String()
		lastX := float64(padding.Left) + float64(len(points)-1)/float64(numPoints-1)*float64(chartW)
		areaPath += fmt.Sprintf(" L %.1f %.1f L %.1f %.1f Z", lastX, float64(padding.Top+chartH), float64(padding.Left), float64(padding.Top+chartH))
		sb.WriteString(fmt.Sprintf(`<path d="%s" fill="%s" opacity="0.15"/>`, areaPath, color))

		// Line
		sb.WriteString(fmt.Sprintf(`<path d="%s" fill="none" stroke="%s" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"/>`, pathD.String(), color))

		// Points
		for i, p := range points {
			x := float64(padding.Left) + float64(i)/float64(numPoints-1)*float64(chartW)
			y := float64(padding.Top) + float64(chartH) - p.Value/maxVal*float64(chartH)
			sb.WriteString(fmt.Sprintf(`<circle cx="%.1f" cy="%.1f" r="4" fill="%s" stroke="#1e1e2e" stroke-width="1.5">`, x, y, color))
			sb.WriteString(fmt.Sprintf(`<title>%s: %s</title>`, escapeXML(p.Label), chart.FormatValue(p.Value, cfg.Precision)))
			sb.WriteString(`</circle>`)
		}
	}

	// X-axis labels
	if numPoints > 0 {
		labelStep := 1
		if numPoints > 20 {
			labelStep = numPoints / 10
		}
		for i := 0; i < numPoints; i += labelStep {
			if len(series[0].Points) > i {
				label := series[0].Points[i].Label
				if len(label) > 8 {
					label = label[:8] + "…"
				}
				x := float64(padding.Left) + float64(i)/float64(numPoints-1)*float64(chartW)
				sb.WriteString(fmt.Sprintf(`<text x="%.0f" y="%d" text-anchor="middle" fill="#a6adc8" font-size="10" transform="rotate(-30 %.0f %d)">%s</text>`, x, height-padding.Bottom+15, x, height-padding.Bottom+15, escapeXML(label)))
			}
		}
	}

	// Axes
	sb.WriteString(fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#585b70" stroke-width="1.5"/>`, padding.Left, padding.Top, padding.Left, padding.Top+chartH))
	sb.WriteString(fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#585b70" stroke-width="1.5"/>`, padding.Left, padding.Top+chartH, width-padding.Right, padding.Top+chartH))

	// Legend
	if cfg.ShowLegend && len(series) > 1 {
		legendY := 25
		for i, s := range series {
			color := SVGColors[i%len(SVGColors)]
			lx := padding.Left + i*100
			sb.WriteString(fmt.Sprintf(`<rect x="%d" y="%d" width="12" height="12" fill="%s" rx="2"/>`, lx, legendY-10, color))
			sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" fill="#cdd6f4" font-size="11">%s</text>`, lx+16, legendY, escapeXML(s.Name)))
		}
	}

	sb.WriteString(`</svg>`)
	return sb.String()
}

func renderPieSVG(series []chart.Series, cfg chart.ChartConfig) string {
	if len(series) == 0 || len(series[0].Points) == 0 {
		return "<!-- No data -->"
	}

	data := series[0].Points
	width := cfg.Width * 8
	height := cfg.Height * 12

	total := 0.0
	for _, p := range data {
		total += p.Value
	}
	if total == 0 {
		return "<!-- Total is zero -->"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`, width, height, width, height))
	sb.WriteString(`<style>text { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; }</style>`)
	sb.WriteString(fmt.Sprintf(`<rect width="%d" height="%d" fill="#1e1e2e" rx="8"/>`, width, height))

	// Title
	if cfg.Title != "" {
		sb.WriteString(fmt.Sprintf(`<text x="%d" y="25" text-anchor="middle" fill="#cdd6f4" font-size="16" font-weight="bold">%s</text>`, width/2, escapeXML(cfg.Title)))
	}

	// Pie chart
	cx := float64(width) / 2
	cy := float64(height)/2 + 10
	radius := math.Min(float64(width), float64(height))/2 - 60

	startAngle := -math.Pi / 2
	for i, p := range data {
		sweepAngle := p.Value / total * 2 * math.Pi
		endAngle := startAngle + sweepAngle
		color := SVGColors[i%len(SVGColors)]

		// Arc path
		x1 := cx + radius*math.Cos(startAngle)
		y1 := cy + radius*math.Sin(startAngle)
		x2 := cx + radius*math.Cos(endAngle)
		y2 := cy + radius*math.Sin(endAngle)

		largeArc := 0
		if sweepAngle > math.Pi {
			largeArc = 1
		}

		d := fmt.Sprintf("M %.1f %.1f L %.1f %.1f A %.1f %.1f 0 %d 1 %.1f %.1f Z", cx, cy, x1, y1, radius, radius, largeArc, x2, y2)
		sb.WriteString(fmt.Sprintf(`<path d="%s" fill="%s" stroke="#1e1e2e" stroke-width="2">`, d, color))
		sb.WriteString(fmt.Sprintf(`<title>%s: %.1f%%</title>`, escapeXML(p.Label), p.Value/total*100))
		sb.WriteString(`</path>`)

		// Label
		pct := p.Value / total * 100
		if pct > 3 { // Only show label if slice is large enough
			labelAngle := startAngle + sweepAngle/2
			labelR := radius * 0.65
			lx := cx + labelR*math.Cos(labelAngle)
			ly := cy + labelR*math.Sin(labelAngle)
			sb.WriteString(fmt.Sprintf(`<text x="%.1f" y="%.1f" text-anchor="middle" fill="#1e1e2e" font-size="11" font-weight="bold">%.0f%%</text>`, lx, ly+4, pct))
		}

		startAngle = endAngle
	}

	// Legend below
	legendY := height - 30
	legendX := 20
	for i, p := range data {
		color := SVGColors[i%len(SVGColors)]
		pct := p.Value / total * 100
		label := p.Label
		if len(label) > 10 {
			label = label[:10] + "…"
		}
		sb.WriteString(fmt.Sprintf(`<rect x="%d" y="%d" width="10" height="10" fill="%s" rx="2"/>`, legendX, legendY, color))
		sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" fill="#cdd6f4" font-size="10">%s (%.1f%%)</text>`, legendX+14, legendY+9, escapeXML(label), pct))
		legendX += len(label)*7 + 50
		if legendX > width-80 {
			legendX = 20
			legendY += 18
		}
	}

	sb.WriteString(`</svg>`)
	return sb.String()
}

func renderScatterSVG(series []chart.Series, cfg chart.ChartConfig) string {
	if len(series) == 0 || len(series[0].Points) == 0 {
		return "<!-- No data -->"
	}

	padding := SVGPadding{Top: 40, Right: 20, Bottom: 60, Left: 60}
	width := cfg.Width * 8
	height := cfg.Height * 12

	chartW := width - padding.Left - padding.Right
	chartH := height - padding.Top - padding.Bottom

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

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`, width, height, width, height))
	sb.WriteString(`<style>text { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; }</style>`)
	sb.WriteString(fmt.Sprintf(`<rect width="%d" height="%d" fill="#1e1e2e" rx="8"/>`, width, height))

	// Title
	if cfg.Title != "" {
		sb.WriteString(fmt.Sprintf(`<text x="%d" y="25" text-anchor="middle" fill="#cdd6f4" font-size="16" font-weight="bold">%s</text>`, width/2, escapeXML(cfg.Title)))
	}

	// Gridlines
	for i := 0; i <= 4; i++ {
		y := padding.Top + chartH - (chartH * i / 4)
		val := minY + (maxY-minY)*float64(i)/4
		sb.WriteString(fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#313244" stroke-width="1"/>`, padding.Left, y, width-padding.Right, y))
		sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" text-anchor="end" fill="#a6adc8" font-size="10">%s</text>`, padding.Left-5, y+4, chart.FormatValue(val, cfg.Precision)))
	}

	// Points
	for si, s := range series {
		color := SVGColors[si%len(SVGColors)]
		for i, p := range s.Points {
			x := float64(padding.Left) + (float64(i)-minX)/(maxX-minX)*float64(chartW)
			y := float64(padding.Top) + float64(chartH) - (p.Value-minY)/(maxY-minY)*float64(chartH)
			sb.WriteString(fmt.Sprintf(`<circle cx="%.1f" cy="%.1f" r="5" fill="%s" opacity="0.8" stroke="#1e1e2e" stroke-width="1">`, x, y, color))
			sb.WriteString(fmt.Sprintf(`<title>%s: %s</title>`, escapeXML(p.Label), chart.FormatValue(p.Value, cfg.Precision)))
			sb.WriteString(`</circle>`)
		}
	}

	// Axes
	sb.WriteString(fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#585b70" stroke-width="1.5"/>`, padding.Left, padding.Top, padding.Left, padding.Top+chartH))
	sb.WriteString(fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#585b70" stroke-width="1.5"/>`, padding.Left, padding.Top+chartH, width-padding.Right, padding.Top+chartH))

	// Legend
	if cfg.ShowLegend && len(series) > 1 {
		for i, s := range series {
			color := SVGColors[i%len(SVGColors)]
			lx := padding.Left + i*100
			sb.WriteString(fmt.Sprintf(`<circle cx="%d" cy="20" r="5" fill="%s"/>`, lx, color))
			sb.WriteString(fmt.Sprintf(`<text x="%d" y="24" fill="#cdd6f4" font-size="11">%s</text>`, lx+10, escapeXML(s.Name)))
		}
	}

	sb.WriteString(`</svg>`)
	return sb.String()
}

func renderHistogramSVG(series []chart.Series, cfg chart.ChartConfig) string {
	if len(series) == 0 || len(series[0].Points) == 0 {
		return "<!-- No data -->"
	}

	padding := SVGPadding{Top: 40, Right: 20, Bottom: 60, Left: 60}
	width := cfg.Width * 8
	height := cfg.Height * 12

	chartW := width - padding.Left - padding.Right
	chartH := height - padding.Top - padding.Bottom

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

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`, width, height, width, height))
	sb.WriteString(`<style>text { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; }</style>`)
	sb.WriteString(fmt.Sprintf(`<rect width="%d" height="%d" fill="#1e1e2e" rx="8"/>`, width, height))

	// Title
	if cfg.Title != "" {
		sb.WriteString(fmt.Sprintf(`<text x="%d" y="25" text-anchor="middle" fill="#cdd6f4" font-size="16" font-weight="bold">%s</text>`, width/2, escapeXML(cfg.Title)))
	}

	// Gridlines
	for i := 0; i <= 4; i++ {
		y := padding.Top + chartH - (chartH * i / 4)
		val := maxCount * i / 4
		sb.WriteString(fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#313244" stroke-width="1"/>`, padding.Left, y, width-padding.Right, y))
		sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" text-anchor="end" fill="#a6adc8" font-size="10">%d</text>`, padding.Left-5, y+4, val))
	}

	// Bars
	barWidth := chartW / len(counts)
	for i, count := range counts {
		barH := float64(count) / float64(maxCount) * float64(chartH)
		x := padding.Left + i*barWidth
		y := float64(padding.Top+chartH) - barH
		color := SVGColors[i%len(SVGColors)]

		sb.WriteString(fmt.Sprintf(`<rect x="%d" y="%.1f" width="%d" height="%.1f" fill="%s" opacity="0.85" rx="2">`, x, y, barWidth-2, barH, color))
		sb.WriteString(fmt.Sprintf(`<title>[%s, %s): %d</title>`, chart.FormatValue(edges[i], cfg.Precision), chart.FormatValue(edges[i+1], cfg.Precision), count))
		sb.WriteString(`</rect>`)

		// Count label
		if count > 0 {
			sb.WriteString(fmt.Sprintf(`<text x="%d" y="%.1f" text-anchor="middle" fill="#cdd6f4" font-size="10">%d</text>`, x+barWidth/2, y-4, count))
		}
	}

	// Axes
	sb.WriteString(fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#585b70" stroke-width="1.5"/>`, padding.Left, padding.Top, padding.Left, padding.Top+chartH))
	sb.WriteString(fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#585b70" stroke-width="1.5"/>`, padding.Left, padding.Top+chartH, width-padding.Right, padding.Top+chartH))

	sb.WriteString(`</svg>`)
	return sb.String()
}

// SVGPadding holds SVG chart padding.
type SVGPadding struct {
	Top, Right, Bottom, Left int
}

func escapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}
