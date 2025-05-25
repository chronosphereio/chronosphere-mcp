package prometheus

import (
	"context"
	"errors"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/chronosphereio/mcp-server/mcp-server/pkg/client"
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
	PromAPI func() (v1.API, error)
}

// RendererOptions contains options for the Renderer.
type RendererOptions struct {
	ClientProvider *client.Provider
}

func NewRenderer(
	opts RendererOptions,
) (*Renderer, error) {
	return &Renderer{
		DataAPI: func() (v1.API, error) {
			c, err := opts.ClientProvider.PrometheusDataClient()
			if err != nil {
				return nil, err
			}
			return v1.NewAPI(c), nil
		},
		PromAPI: func() (v1.API, error) {
			c, err := opts.ClientProvider.PrometheusPromClient()
			if err != nil {
				return nil, err
			}
			return v1.NewAPI(c), nil
		},
	}, nil
}

// RenderSeries renders a time series graph for the given series.
func (r *Renderer) RenderSeries(w io.Writer, series model.Matrix, ws, hs int, legend bool) error {
	plot := plot.New()
	plot.Legend.Top = true

	pts := formatSeriesForRender(series, legend)

	if err := plotutil.AddLinePoints(plot, pts...); err != nil {
		return err
	}

	plot.Y.Max = plot.Y.Max * 1.2
	c := vgimg.New(vg.Length(ws), vg.Length(hs))
	cpng := vgimg.PngCanvas{Canvas: c}
	cv := draw.New(cpng)
	plot.Draw(cv)

	_, err := cpng.WriteTo(w)
	return err
}

// Render renders a time series graph for the given query and time range.
func (r *Renderer) Render(ctx context.Context, w io.Writer, query string, start, end time.Time, ws, hs int, legend bool) error {
	plot := plot.New()
	plot.Legend.Top = true
	timeseries, err := r.queryRange(ctx, query, start, end, legend)
	if err != nil {
		return err
	}

	if err := plotutil.AddLinePoints(plot, timeseries...); err != nil {
		return err
	}

	plot.Y.Max = plot.Y.Max * 1.2
	c := vgimg.New(vg.Length(ws), vg.Length(hs))
	cpng := vgimg.PngCanvas{Canvas: c}
	cv := draw.New(cpng)
	plot.Draw(cv)

	_, err = cpng.WriteTo(w)
	return err
}

// returns name, plotter.XYer, name1, plotter.XYer ...
func (r *Renderer) queryRange(ctx context.Context, query string, start, end time.Time, legend bool) ([]interface{}, error) {
	api, err := r.DataAPI()
	if err != nil {
		return nil, err
	}
	resp, _, err := api.QueryRange(ctx, query, v1.Range{
		Start: start,
		End:   end,
		Step:  60 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	matrix, ok := resp.(model.Matrix)
	if !ok {
		return nil, errors.New("expected matrix")
	}
	return formatSeriesForRender(matrix, legend), nil
}

func formatSeriesForRender(series model.Matrix, legend bool) []any {
	sort.Sort(series)
	plotters := make([]any, 0, len(series))
	for _, s := range series {
		pts := make(plotter.XYs, len(s.Values))
		for j, sample := range s.Values {
			pts[j].Y = float64(sample.Value)
			pts[j].X = float64(sample.Timestamp.Unix()) / 60
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
		values[i] = strings.Replace(string(v), "\n", " ", -1)
		i++
	}
	return strings.Join(values, "|")
}
