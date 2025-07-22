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
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/prometheus/client_golang/api"
	"go.uber.org/zap"

	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools/pkg/params"
	"github.com/chronosphereio/chronosphere-mcp/pkg/links"
)

var _ tools.MCPTools = (*Tools)(nil)

type Tools struct {
	logger      *zap.Logger
	renderer    *Renderer
	linkBuilder *links.Builder
}

// NewTools creates a new Tools instance.
func NewTools(api api.Client, logger *zap.Logger, linkBuilder *links.Builder) (*Tools, error) {
	renderer, err := NewRenderer(RendererOptions{
		api: api,
	})
	if err != nil {
		return nil, err
	}

	logger.Info("prometheus tool configured")

	return &Tools{
		logger:      logger,
		renderer:    renderer,
		linkBuilder: linkBuilder,
	}, nil
}

func (t *Tools) GroupName() string {
	return "metrics"
}

// MCPTools returns the list of tools that this package provides.
func (t *Tools) MCPTools() []tools.MCPTool {
	// Metrics-related tools
	return []tools.MCPTool{
		{
			Metadata: tools.NewMetadata("render_prometheus_range_query",
				mcp.WithDescription("Evaluates a Prometheus expression query over a range of time and renders it as a PNG image."),
				mcp.WithString("query",
					mcp.Description("Prometheus PromQL expression query string"),
					mcp.Required(),
				),
				params.WithTimeRange(),
				mcp.WithNumber("step_seconds",
					mcp.Description("Optional. Query resolution step width, in number of seconds."),
					mcp.DefaultNumber(60),
				),
			),
			Handler: t.renderPrometheusRangeQuery,
		},
		{
			Metadata: tools.NewMetadata("query_prometheus_range",
				mcp.WithDescription("Evaluates a Prometheus expression query over a range of time"),
				mcp.WithString("query",
					mcp.Description("Prometheus PromQL expression query string"),
					mcp.Required(),
				),
				params.WithTimeRange(),
				mcp.WithNumber("step_seconds",
					mcp.Description("Query resolution step width, in number of seconds. Optional."),
					mcp.DefaultNumber(60),
				),
				mcp.WithNumber("limit",
					mcp.Description("Maximum number of time series to return. Default is 100. Set to 0 for no limit."),
					mcp.DefaultNumber(100),
				),
				mcp.WithNumber("offset",
					mcp.Description("Number of time series to skip before returning results. Default is 0."),
					mcp.DefaultNumber(0),
				),
			),
			Handler: t.queryPrometheusRange,
		},
		{
			Metadata: tools.NewMetadata("query_prometheus_instant",
				mcp.WithDescription("Evaluates a Prometheus instant query at a single point in time"),
				mcp.WithString("query",
					mcp.Description("The PromQL expression to query"),
					mcp.Required(),
				),
				mcp.WithString("time",
					mcp.Description("RFC3389 representation of the time. Optional. Defaults to the current time"),
				),
			),
			Handler: t.queryPrometheusInstant,
		},
		{
			Metadata: tools.NewMetadata("list_prometheus_series",
				mcp.WithDescription("Returns the list of time series that match a certain label set."),
				params.WithStringArray("selectors",
					mcp.Description("Repeated series selector arguments that select the series to return. At least one selector must be provided."),
					mcp.Required(),
				),
				params.WithTimeRange(),
			),
			Handler: t.listPrometheusSeries,
		},
		{
			Metadata: tools.NewMetadata("list_prometheus_label_values",
				mcp.WithDescription("returns the list of all label values that (optionally) match a set of selectors and/or were present during a given time range."),
				mcp.WithString("label_name",
					mcp.Description("The label name for which values should be returned. Required."),
					mcp.Required(),
				),
				params.WithStringArray("selectors",
					mcp.Description("Repeated series selector argument that selects the series to return. This is an array of PromQL query strings. Optional."),
				),
				params.WithTimeRange(),
				mcp.WithNumber("limit",
					mcp.Description("Maximum number of label values to return. Default is 100. Set to 0 for no limit."),
					mcp.DefaultNumber(100),
				),
				mcp.WithNumber("offset",
					mcp.Description("Number of label values to skip before returning results. Default is 0."),
					mcp.DefaultNumber(0),
				),
			),
			Handler: t.listPrometheusLabelValues,
		},
		{
			Metadata: tools.NewMetadata("list_prometheus_label_names",
				mcp.WithDescription("returns the list of all label names that (optionally) match a set of selectors and/or were present during a given time range."),
				params.WithStringArray("selectors",
					mcp.Description("Repeated series selector argument that selects the series to return. This is an array of strings. Optional."),
				),
				params.WithTimeRange(),
			),
			Handler: t.listPrometheusLabelNames,
		},
		{
			Metadata: tools.NewMetadata("list_prometheus_series_metadata",
				mcp.WithDescription(
					"Returns metadata for the given metric in the form of \"type\", \"help\", and \"unit\". "+
						"If the metric does not exist or there is no metadata for this metric, return an empty object."),
				mcp.WithString("metric", mcp.Description("The metric name for which metadata should be returned. Required.")),
			),
			Handler: t.listPrometheusSeriesMetadata,
		},
	}
}
