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
	"io"
	"sort"
	"strings"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
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

// RenderSeries renders a time series graph for the given series.
func (r *Renderer) RenderSeries(w io.Writer, series model.Matrix, ws, hs int, legend bool) error {
	p := plot.New()
	p.Legend.Top = true

	pts := formatSeriesForRender(series, legend)

	if err := plotutil.AddLinePoints(p, pts...); err != nil {
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

	if err := plotutil.AddLinePoints(p, timeseries...); err != nil {
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

// Ticks returns tick marks for the given min and max values
func (t *timeTickMarker) Ticks(min, max float64) []plot.Tick {
	// Calculate a reasonable number of ticks based on the time range
	timeRange := max - min
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
	step := (max - min) / float64(tickCount-1)

	for i := 0; i < tickCount; i++ {
		val := min + float64(i)*step
		timestamp := time.Unix(int64(val), 0).UTC()
		ticks = append(ticks, plot.Tick{
			Value: val,
			Label: timestamp.Format(format),
		})
	}

	// Add minor ticks between major ticks for better readability
	minorStep := step / 5
	for i := 0; i < tickCount-1; i++ {
		baseVal := min + float64(i)*step
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
