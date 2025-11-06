// Copyright 2025 Chronosphere Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package prometheus

import (
	"context"
	"errors"
	"image/color"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

type Renderer struct {
	DataAPI func() (v1.API, error)
}

// RendererOptions contains options for the Renderer.
type RendererOptions struct {
	api api.Client
}

func NewRenderer(
	opts RendererOptions,
) (*Renderer, error) {
	return &Renderer{
		DataAPI: func() (v1.API, error) {
			return v1.NewAPI(opts.api), nil
		},
	}, nil
}

// Dark theme colors - professional and easy on the eyes
var (
	darkBackground   = color.RGBA{R: 30, G: 30, B: 30, A: 255}    // Dark gray background
	darkGridColor    = color.RGBA{R: 60, G: 60, B: 60, A: 255}    // Subtle grid lines
	darkTextColor    = color.RGBA{R: 220, G: 220, B: 220, A: 255} // Light gray text
	darkSeriesColors = []color.Color{
		color.RGBA{R: 99, G: 179, B: 237, A: 255},  // Light blue
		color.RGBA{R: 46, G: 204, B: 113, A: 255},  // Green
		color.RGBA{R: 241, G: 196, B: 15, A: 255},  // Yellow
		color.RGBA{R: 231, G: 76, B: 60, A: 255},   // Red
		color.RGBA{R: 155, G: 89, B: 182, A: 255},  // Purple
		color.RGBA{R: 52, G: 152, B: 219, A: 255},  // Blue
		color.RGBA{R: 26, G: 188, B: 156, A: 255},  // Teal
		color.RGBA{R: 230, G: 126, B: 34, A: 255},  // Orange
		color.RGBA{R: 236, G: 240, B: 241, A: 255}, // Light gray
		color.RGBA{R: 149, G: 165, B: 166, A: 255}, // Gray
	}
)

// applyDarkTheme applies a modern dark theme to the plot
func applyDarkTheme(p *plot.Plot) {
	// Set background color
	p.BackgroundColor = darkBackground

	// Configure title styling
	p.Title.TextStyle.Color = darkTextColor
	p.Title.TextStyle.Font.Size = vg.Points(14)

	// Configure X-axis styling
	p.X.Label.TextStyle.Color = darkTextColor
	p.X.Label.TextStyle.Font.Size = vg.Points(11)
	p.X.Tick.Label.Color = darkTextColor
	p.X.Tick.Label.Font.Size = vg.Points(9)
	p.X.Tick.LineStyle.Color = darkGridColor
	p.X.LineStyle.Color = darkTextColor

	// Configure Y-axis styling
	p.Y.Label.TextStyle.Color = darkTextColor
	p.Y.Label.TextStyle.Font.Size = vg.Points(11)
	p.Y.Tick.Label.Color = darkTextColor
	p.Y.Tick.Label.Font.Size = vg.Points(9)
	p.Y.Tick.LineStyle.Color = darkGridColor
	p.Y.LineStyle.Color = darkTextColor

	// Configure legend styling
	p.Legend.TextStyle.Color = darkTextColor
	p.Legend.TextStyle.Font.Size = vg.Points(9)
}

// addStyledLines adds line plots with dark theme colors and proper styling
func addStyledLines(p *plot.Plot, pts []any) error {
	colorIndex := 0

	for i := 0; i < len(pts); {
		var name string
		var xyData plotter.XYer

		// Check if we have a name (legend entry)
		if i+1 < len(pts) {
			if n, ok := pts[i].(string); ok {
				name = n
				i++
			}
		}

		// Get the XY data
		if i < len(pts) {
			if xy, ok := pts[i].(plotter.XYer); ok {
				xyData = xy
				i++
			} else {
				return errors.New("expected plotter.XYer")
			}
		}

		// Create and style the line
		line, err := plotter.NewLine(xyData)
		if err != nil {
			return err
		}

		// Apply color from our dark theme palette
		line.Color = darkSeriesColors[colorIndex%len(darkSeriesColors)]
		line.Width = vg.Points(2)

		// Add to plot
		p.Add(line)

		// Add legend entry if we have a name
		if name != "" {
			p.Legend.Add(name, line)
		}

		colorIndex++
	}

	return nil
}

// RenderSeries renders a time series graph for the given series.
func (r *Renderer) RenderSeries(w io.Writer, series model.Matrix, ws, hs int, legend bool) error {
	p := plot.New()
	p.Legend.Top = true

	// Apply dark theme
	applyDarkTheme(p)

	pts := formatSeriesForRender(series, legend)

	if err := addStyledLines(p, pts); err != nil {
		return err
	}

	// Configure X-axis with ISO 8601 timestamp formatter and label
	p.X.Tick.Marker = &timeTickMarker{}
	p.X.Label.Text = "Time (UTC)"

	// Configure Y-axis label
	p.Y.Label.Text = detectYAxisLabel(series)

	p.Y.Max = p.Y.Max * 1.2
	c := vgimg.New(vg.Length(ws), vg.Length(hs))
	cpng := vgimg.PngCanvas{Canvas: c}
	cv := draw.New(cpng)
	p.Draw(cv)

	_, err := cpng.WriteTo(w)
	return err
}

// Render renders a time series graph for the given query and time range.
func (r *Renderer) Render(ctx context.Context, w io.Writer, query string, start, end time.Time, ws, hs int, legend bool) error {
	p := plot.New()
	p.Legend.Top = true

	// Apply dark theme
	applyDarkTheme(p)

	// Query the data
	api, err := r.DataAPI()
	if err != nil {
		return err
	}
	resp, _, err := api.QueryRange(ctx, query, v1.Range{
		Start: start,
		End:   end,
		Step:  60 * time.Second,
	})
	if err != nil {
		return err
	}
	matrix, ok := resp.(model.Matrix)
	if !ok {
		return errors.New("expected matrix")
	}

	timeseries := formatSeriesForRender(matrix, legend)

	if err := addStyledLines(p, timeseries); err != nil {
		return err
	}

	// Configure X-axis with ISO 8601 timestamp formatter and label
	p.X.Tick.Marker = &timeTickMarker{}
	p.X.Label.Text = "Time (UTC)"

	// Configure Y-axis label
	p.Y.Label.Text = detectYAxisLabel(matrix)

	p.Y.Max = p.Y.Max * 1.2
	c := vgimg.New(vg.Length(ws), vg.Length(hs))
	cpng := vgimg.PngCanvas{Canvas: c}
	cv := draw.New(cpng)
	p.Draw(cv)

	_, err = cpng.WriteTo(w)
	return err
}

// detectYAxisLabel attempts to detect an appropriate Y-axis label from the metric data
func detectYAxisLabel(series model.Matrix) string {
	if len(series) == 0 {
		return "Value"
	}

	// Get the first metric's __name__ label to determine the metric type
	firstMetric := series[0].Metric
	metricName := string(firstMetric["__name__"])

	// Common metric patterns and their units
	switch {
	case strings.Contains(metricName, "_bytes"):
		return "Bytes"
	case strings.Contains(metricName, "_seconds"):
		return "Seconds"
	case strings.Contains(metricName, "_milliseconds"):
		return "Milliseconds"
	case strings.Contains(metricName, "_ratio") || strings.Contains(metricName, "_percent"):
		return "Ratio"
	case strings.Contains(metricName, "_total") || strings.Contains(metricName, "_count"):
		return "Count"
	case strings.Contains(metricName, "_requests"):
		return "Requests"
	case strings.Contains(metricName, "cpu"):
		return "CPU Usage"
	case strings.Contains(metricName, "memory"):
		return "Memory"
	case strings.Contains(metricName, "latency"):
		return "Latency"
	case strings.Contains(metricName, "_rate"):
		return "Rate"
	default:
		// If we have a metric name, use it as the label
		if metricName != "" {
			return metricName
		}
		return "Value"
	}
}

func formatSeriesForRender(series model.Matrix, legend bool) []any {
	sort.Sort(series)
	plotters := make([]any, 0, len(series))
	for _, s := range series {
		pts := make(plotter.XYs, len(s.Values))
		for j, sample := range s.Values {
			pts[j].Y = float64(sample.Value)
			// Store Unix timestamp directly for proper time axis formatting
			pts[j].X = float64(sample.Timestamp.Unix())
		}
		if legend {
			plotters = append(plotters, formatMetric(s.Metric))
		}
		plotters = append(plotters, pts)
	}
	return plotters
}

func formatMetric(m model.Metric) string {
	ls := model.LabelSet(m)
	values := make([]string, len(ls))
	i := 0
	for _, v := range ls {
		values[i] = strings.ReplaceAll(string(v), "\n", " ")
		i++
	}
	return strings.Join(values, "|")
}

// timeTickMarker implements plot.Ticker to format X-axis ticks as ISO 8601 timestamps
type timeTickMarker struct{}

// Ticks returns tick marks for the given minVal and maxVal values
func (t *timeTickMarker) Ticks(minVal, maxVal float64) []plot.Tick {
	// Calculate a reasonable number of ticks based on the time range
	timeRange := maxVal - minVal
	var tickCount int
	var format string

	switch {
	case timeRange < 3600: // Less than 1 hour
		tickCount = 6
		format = "15:04:05" // HH:MM:SS
	case timeRange < 86400: // Less than 1 day
		tickCount = 8
		format = "15:04" // HH:MM
	case timeRange < 604800: // Less than 1 week
		tickCount = 7
		format = "01-02 15:04" // MM-DD HH:MM
	default: // 1 week or more
		tickCount = 10
		format = "2006-01-02" // YYYY-MM-DD
	}

	ticks := make([]plot.Tick, 0, tickCount)
	step := (maxVal - minVal) / float64(tickCount-1)

	for i := 0; i < tickCount; i++ {
		val := minVal + float64(i)*step
		timestamp := time.Unix(int64(val), 0).UTC()
		ticks = append(ticks, plot.Tick{
			Value: val,
			Label: timestamp.Format(format),
		})
	}

	// Add minor ticks between major ticks for better readability
	minorStep := step / 5
	for i := 0; i < tickCount-1; i++ {
		baseVal := minVal + float64(i)*step
		for j := 1; j < 5; j++ {
			minorVal := baseVal + float64(j)*minorStep
			ticks = append(ticks, plot.Tick{
				Value: minorVal,
				Label: "",
			})
		}
	}

	return ticks
}
