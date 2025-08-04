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
				mcp.WithDescription(`Executes a Prometheus PromQL query over a specified time range and returns time series data points as JSON.

Supports standard PromQL syntax plus Chronosphere custom functions:
- cardinality_estimate(vector): Estimates element count in vector, useful for understanding metric cardinality. Example: cardinality_estimate(http_requests_total) by (service)
- head_avg/head_max/head_min/head_sum(query, n): Returns top n series by aggregation value. Example: head_avg(cpu_usage{}, 5) for top 5 CPU users
- tail_avg/tail_max/tail_min/tail_sum(query, n): Returns bottom n series by aggregation value. Example: tail_avg(memory_usage{}, 3) for lowest memory users  
- sum_per_second(range_vector): Calculates per-second rate for delta counters. Example: sum_per_second(http_request_count{}[5m])

Returns raw time series data - use query_prometheus_instant for single point-in-time values.`),
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
				mcp.WithDescription(`Returns the complete time series (full label sets with all key-value pairs) that match the given selectors. Each result shows the exact combination of labels for an active time series. Use this tool only when you need to see the actual label combinations that exist.

IMPORTANT: This tool returns a lot of data and can overwhelm context windows. For most use cases, prefer:
- list_prometheus_label_names with selector {__name__="metric_name"} to find what labels are available on a metric
- list_prometheus_label_values to find what values a specific label has
- cardinality_estimate() function in query tools to estimate series count without returning all data

Use list_prometheus_series only when you specifically need to see exact label combinations, such as debugging specific series or finding unusual label patterns.

Example usage:
- Find specific series combinations: ["{__name__=\"http_requests_total\", service=\"api\"}"]
- Debug label patterns: ["{job=\"kubernetes-pods\"}"] (use sparingly)`),
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
				mcp.WithDescription(`Returns the list of values for a specific label name, optionally filtered by selectors. Use this tool when you know the label name and want to discover what values it has across your metrics.

Common use cases:
- Find metric names for a service: label_name="__name__" with selectors=["{service=\"api\"}"]
- Find all HTTP status codes: label_name="status" with selectors=["{__name__=\"http_requests_total\"}"]
- Discover all services: label_name="service" (optionally with additional selectors)
- Find all methods for an API: label_name="method" with selectors=["{service=\"api\"}"]
- Get all Kubernetes namespaces: label_name="kubernetes_namespace"

Example usage:
- Metrics for a service: label_name="__name__", selectors=["{service=\"web-server\"}"]
- HTTP status codes: label_name="status", selectors=["{__name__=\"http_requests_total\"}"]
- Available services: label_name="service", selectors=["{job=\"kubernetes-pods\"}"]
- API endpoints: label_name="endpoint", selectors=["{service=\"web-api\"}"]

Supports pagination with limit/offset. To find what label names are available first, use list_prometheus_label_names.`),
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
				mcp.WithDescription(`Returns the list of label names (keys) available on metrics that match the given selectors. Use this tool when you need to discover what labels are available on specific metrics or services.

Example usage:
- Labels on a metric: selectors=["{__name__=\"http_requests_total\"}"]
- Labels for a service: selectors=["{service=\"web-server\"}"]
- Labels in a namespace: selectors=["{kubernetes_namespace=\"production\"}"]

This returns only the label keys (e.g., "method", "status", "service"). To get the values for a specific label, use list_prometheus_label_values.`),
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
