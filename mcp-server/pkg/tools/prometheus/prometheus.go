// Package prometheus provides the tools for fetching and rendering prometheus metrics.
package prometheus

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools/pkg/params"
	"github.com/mark3labs/mcp-go/mcp"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

func (t *Tools) listPrometheusSeries(ctx context.Context, request mcp.CallToolRequest) (*tools.Result, error) {
	selectors, errResult := params.StringArray(request, "selectors", true, nil)
	if errResult != nil {
		return nil, errResult
	}

	if len(selectors) == 0 {
		return nil, fmt.Errorf("at least one selector must be provided")
	}

	timeRange, err := params.ParseTimeRange(request)
	if err != nil {
		return nil, err
	}

	api, err := t.renderer.DataAPI()
	if err != nil {
		return nil, err
	}
	resp, warnings, err := api.Series(ctx, selectors, timeRange.Start, timeRange.End)
	if err != nil {
		return nil, fmt.Errorf("failed to get series: %s", err)
	}

	return promJSONResponse(resp, warnings)
}

func (t *Tools) queryPrometheusRange(ctx context.Context, request mcp.CallToolRequest) (*tools.Result, error) {
	query, err := params.String(request, "query", true, "")
	if err != nil {
		return nil, err
	}

	timeRange, err := params.ParseTimeRange(request)
	if err != nil {
		return nil, err
	}
	stepSeconds, err := params.Float(request, "step_seconds", false, 60)
	if err != nil {
		return nil, err
	}

	step := time.Duration(stepSeconds) * time.Second
	if step <= 0 {
		return nil, fmt.Errorf("stepSeconds must be a positive integer")
	}

	api, err := t.renderer.DataAPI()
	if err != nil {
		return nil, err
	}
	resp, warnings, err := api.QueryRange(ctx, query, v1.Range{
		Start: timeRange.Start,
		End:   timeRange.End,
		Step:  step,
	})
	if err != nil {
		return nil, err
	}
	matrix, ok := resp.(model.Matrix)
	if !ok {
		return nil, fmt.Errorf("unexpected result from prometheus server")
	}

	return promJSONResponse(matrix, warnings)
}

func (t *Tools) renderPrometheusRangeQuery(ctx context.Context, request mcp.CallToolRequest) (*tools.Result, error) {
	query, err := params.String(request, "query", true, "")
	if err != nil {
		return nil, err
	}
	timeRange, err := params.ParseTimeRange(request)
	if err != nil {
		return nil, err
	}

	step, err := params.Int(request, "step_seconds", false, 60)
	if err != nil {
		return nil, err
	}

	api, err := t.renderer.DataAPI()
	if err != nil {
		return nil, err
	}
	resp, _, err := api.QueryRange(ctx, query, v1.Range{
		Start: timeRange.Start,
		End:   timeRange.End,
		Step:  time.Duration(step) * time.Second,
	})
	if err != nil {
		return nil, err
	}
	matrix, ok := resp.(model.Matrix)
	if !ok {
		return nil, fmt.Errorf("unexpected result from prometheus server")
	}

	data, err := t.renderPrometheusPNG(matrix)
	if err != nil {
		return nil, fmt.Errorf("failed to render query: %s", err)
	}

	return &tools.Result{
		ImageContent: data,
	}, nil
}

func (t *Tools) renderPrometheusPNG(
	matrix model.Matrix,
) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	err := t.renderer.RenderSeries(buf, matrix, 1024, 768, true)
	if err != nil {
		return nil, fmt.Errorf("unable to render data: %w", err)
	}
	return buf.Bytes(), nil
}

func (t *Tools) queryPrometheusInstant(ctx context.Context, request mcp.CallToolRequest) (*tools.Result, error) {
	expression, errResult := params.String(request, "expression", true, "")
	if errResult != nil {
		return nil, errResult
	}
	evalTime, errResult := params.Time(request, "time", false, time.Now())
	if errResult != nil {
		return nil, errResult
	}
	api, err := t.renderer.DataAPI()
	if err != nil {
		return nil, err
	}
	resp, warnings, err := api.Query(ctx, expression, evalTime)
	if err != nil {
		return nil, fmt.Errorf("failed to query prometheus: %s", err)
	}

	return promJSONResponse(resp, warnings)
}

func (t *Tools) listPrometheusLabelValues(ctx context.Context, request mcp.CallToolRequest) (*tools.Result, error) {
	labelName, errResult := params.String(request, "label_name", true, "")
	if errResult != nil {
		return nil, errResult
	}
	selectors, errResult := params.StringArray(request, "selectors", false, nil)
	if errResult != nil {
		return nil, errResult
	}
	timeRange, err := params.ParseTimeRange(request)
	if err != nil {
		return nil, err
	}
	api, err := t.renderer.DataAPI()
	if err != nil {
		return nil, err
	}
	resp, warnings, err := api.LabelValues(ctx, labelName, selectors, timeRange.Start, timeRange.End)
	if err != nil {
		return nil, fmt.Errorf("failed to get label names: %s", err)
	}

	return promJSONResponse(resp, warnings)
}

func (t *Tools) listPrometheusLabelNames(ctx context.Context, request mcp.CallToolRequest) (*tools.Result, error) {
	selectors, errResult := params.StringArray(request, "selectors", false, nil)
	if errResult != nil {
		return nil, errResult
	}

	timeRange, err := params.ParseTimeRange(request)
	if err != nil {
		return nil, err
	}

	api, err := t.renderer.DataAPI()
	if err != nil {
		return nil, err
	}
	// TODO: should we set a limit?
	resp, warnings, err := api.LabelNames(ctx, selectors, timeRange.End, timeRange.End)
	if err != nil {
		return nil, fmt.Errorf("failed to get label names: %s", err)
	}
	return promJSONResponse(resp, warnings)
}

func (t *Tools) listPrometheusSeriesMetadata(ctx context.Context, request mcp.CallToolRequest) (*tools.Result, error) {
	metric, err := params.String(request, "metric", true, "")
	if err != nil {
		return nil, err
	}
	api, err := t.renderer.DataAPI()
	if err != nil {
		return nil, err
	}
	resp, err := api.Metadata(ctx, metric, "1")
	if err != nil {
		return nil, fmt.Errorf("failed to get series metadata: %s", err)
	}
	respMetric, ok := resp[metric]
	if !ok || len(respMetric) == 0 {
		return &tools.Result{
			JSONContent: map[string]any{
				"metadata": map[string]v1.Metadata{},
			},
		}, nil
	}
	return &tools.Result{
		JSONContent: map[string]any{
			"metadata": respMetric[0],
		},
	}, nil
}

func promJSONResponse(resp any, warnings []string) (*tools.Result, error) {
	result := &tools.Result{
		JSONContent: map[string]any{
			"result": resp,
		},
	}
	if len(warnings) > 0 {
		result.Meta = map[string]any{
			"warnings": warnings,
		}
	}
	return result, nil
}
