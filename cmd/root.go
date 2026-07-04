package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/EdgarOrtegaRamirez/chartforge/chart"
	"github.com/EdgarOrtegaRamirez/chartforge/reader"
	"github.com/EdgarOrtegaRamirez/chartforge/render"
	"github.com/spf13/cobra"
)

var version = "0.1.0"

var rootCmd = &cobra.Command{
	Use:   "chartforge",
	Short: "A data visualization CLI tool",
	Long:  `ChartForge generates beautiful charts (bar, line, pie, scatter, histogram) from CSV and JSON data files.`,
}

var chartCmd = &cobra.Command{
	Use:   "chart [type] [file]",
	Short: "Generate a chart from data file",
	Long: `Generate a chart from a CSV or JSON data file.

Supported chart types: bar, line, pie, scatter, histogram

Examples:
  chartforge chart bar data.csv
  chartforge chart line data.json --output chart.svg
  chartforge chart pie data.csv --title "Sales by Region"
  chartforge chart histogram data.csv --bins 15
  cat data.csv | chartforge chart bar - --terminal`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		chartType := args[0]
		filePath := args[1]

		// Validate chart type
		validTypes := map[string]bool{"bar": true, "line": true, "pie": true, "scatter": true, "histogram": true, "hist": true}
		if !validTypes[chartType] {
			return fmt.Errorf("unsupported chart type: %s (valid: bar, line, pie, scatter, histogram)", chartType)
		}

		// Build config
		cfg := chart.DefaultConfig()
		cfg.Type = chart.ChartType(chartType)
		if chartType == "hist" {
			cfg.Type = chart.Hist
		}

		// Read flags
		title, _ := cmd.Flags().GetString("title")
		width, _ := cmd.Flags().GetInt("width")
		height, _ := cmd.Flags().GetInt("height")
		output, _ := cmd.Flags().GetString("output")
		terminalMode, _ := cmd.Flags().GetBool("terminal")
		bins, _ := cmd.Flags().GetInt("bins")
		precision, _ := cmd.Flags().GetInt("precision")
		noLegend, _ := cmd.Flags().GetBool("no-legend")
		noValues, _ := cmd.Flags().GetBool("no-values")

		if title != "" {
			cfg.Title = title
		}
		if width > 0 {
			cfg.Width = width
		}
		if height > 0 {
			cfg.Height = height
		}
		if bins > 0 {
			cfg.BinCount = bins
		}
		cfg.Precision = precision
		cfg.ShowLegend = !noLegend
		cfg.ShowValues = !noValues

		// Validate config
		if err := cfg.Validate(); err != nil {
			return err
		}

		// Read data
		var series []chart.Series
		var err error

		if filePath == "-" {
			// Read from stdin
			if strings.HasSuffix(output, ".svg") || (!terminalMode && output != "") {
				series, err = reader.ParseCSV(os.Stdin)
			} else {
				series, err = reader.ParseCSV(os.Stdin)
			}
		} else if strings.HasSuffix(filePath, ".json") {
			series, err = reader.ReadJSON(filePath)
		} else {
			series, err = reader.ReadCSV(filePath)
		}

		if err != nil {
			return fmt.Errorf("reading data: %w", err)
		}

		// Render
		var result string
		if terminalMode || output == "" {
			result = render.RenderTerminal(series, cfg)
		} else if strings.HasSuffix(output, ".svg") {
			result = render.RenderSVG(series, cfg)
		} else {
			// Default to SVG for file output
			result = render.RenderSVG(series, cfg)
		}

		// Write output
		if output != "" && output != "-" {
			return os.WriteFile(output, []byte(result), 0644)
		}

		fmt.Print(result)
		return nil
	},
}

var infoCmd = &cobra.Command{
	Use:   "info [file]",
	Short: "Show data file information",
	Long:  `Display information about a CSV or JSON data file including column names, row count, and basic statistics.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]

		var series []chart.Series
		var err error

		if filePath == "-" {
			series, err = reader.ParseCSV(os.Stdin)
		} else if strings.HasSuffix(filePath, ".json") {
			series, err = reader.ReadJSON(filePath)
		} else {
			series, err = reader.ReadCSV(filePath)
		}

		if err != nil {
			return fmt.Errorf("reading data: %w", err)
		}

		fmt.Printf("File: %s\n", filePath)
		fmt.Printf("Series: %d\n", len(series))
		fmt.Println()

		for i, s := range series {
			fmt.Printf("Series %d: %s\n", i+1, s.Name)
			fmt.Printf("  Points: %d\n", len(s.Points))
			if len(s.Points) > 0 {
				minVal, maxVal := s.Points[0].Value, s.Points[0].Value
				sum := 0.0
				for _, p := range s.Points {
					if p.Value < minVal {
						minVal = p.Value
					}
					if p.Value > maxVal {
						maxVal = p.Value
					}
					sum += p.Value
				}
				avg := sum / float64(len(s.Points))
				fmt.Printf("  Min: %.2f\n", minVal)
				fmt.Printf("  Max: %.2f\n", maxVal)
				fmt.Printf("  Avg: %.2f\n", avg)
				fmt.Printf("  Sum: %.2f\n", sum)
			}
			fmt.Println()
		}

		return nil
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("chartforge %s\n", version)
	},
}

func init() {
	// Chart command flags
	chartCmd.Flags().StringP("title", "t", "", "Chart title")
	chartCmd.Flags().IntP("width", "W", 80, "Chart width (terminal columns or SVG pixels/8)")
	chartCmd.Flags().IntP("height", "H", 20, "Chart height (terminal rows or SVG pixels/12)")
	chartCmd.Flags().StringP("output", "o", "", "Output file (e.g., chart.svg). Defaults to terminal.")
	chartCmd.Flags().BoolP("terminal", "T", false, "Force terminal output")
	chartCmd.Flags().IntP("bins", "b", 10, "Number of bins for histogram")
	chartCmd.Flags().IntP("precision", "p", 1, "Decimal precision for value labels")
	chartCmd.Flags().Bool("no-legend", false, "Hide legend")
	chartCmd.Flags().Bool("no-values", false, "Hide value labels")

	rootCmd.AddCommand(chartCmd)
	rootCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(versionCmd)
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
