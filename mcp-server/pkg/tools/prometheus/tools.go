package prometheus

import (
	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"

	"github.com/chronosphereio/mcp-server/mcp-server/pkg/client"
	"github.com/chronosphereio/mcp-server/mcp-server/pkg/tools"
	"github.com/chronosphereio/mcp-server/mcp-server/pkg/tools/pkg/params"
)

var _ tools.MCPTools = (*Tools)(nil)

type Tools struct {
	logger   *zap.Logger
	renderer *Renderer
}

// NewTools creates a new Tools instance.
func NewTools(clientProvider *client.Provider, logger *zap.Logger) (*Tools, error) {
	renderer, err := NewRenderer(RendererOptions{
		ClientProvider: clientProvider,
	})
	if err != nil {
		return nil, err
	}

	logger.Info("prometheus tool configured")

	return &Tools{
		logger:   logger,
		renderer: renderer,
	}, nil
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
			),
			Handler: t.queryPrometheusRange,
		},
		{
			Metadata: tools.NewMetadata("query_prometheus_instant",
				mcp.WithDescription("Evaluates a Prometheus instant query at a single point in time"),
				mcp.WithString("expression",
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
				mcp.WithArray("selectors",
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
				mcp.WithArray("selectors",
					mcp.Description("Repeated series selector argument that selects the series to return. This is an array of strings. Optional."),
				),
				params.WithTimeRange(),
			),
			Handler: t.listPrometheusLabelValues,
		},
		{
			Metadata: tools.NewMetadata("list_prometheus_label_names",
				mcp.WithDescription("returns the list of all label names that (optionally) match a set of selectors and/or were present during a given time range."),
				mcp.WithArray("selectors",
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
