// Package traces contains tools for querying Chronosphere traces.
package traces

import (
	"context"
	"fmt"

	"github.com/go-openapi/strfmt"
	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"

	"github.com/chronosphereio/chronosphere-mcp/generated/dataunstable/dataunstable"
	"github.com/chronosphereio/chronosphere-mcp/generated/dataunstable/dataunstable/data_unstable"
	"github.com/chronosphereio/chronosphere-mcp/generated/dataunstable/models"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools/pkg/params"
)

var _ tools.MCPTools = (*Tools)(nil)

type Tools struct {
	logger *zap.Logger
	api    *dataunstable.DataUnstableAPI
}

func NewTools(
	api *dataunstable.DataUnstableAPI,
	logger *zap.Logger,
) (*Tools, error) {
	logger.Info("events tool configured")

	return &Tools{
		logger: logger,
		api:    api,
	}, nil
}

func (t *Tools) GroupName() string {
	return "traces"
}

// MCPTools returns a list of MCP tools related to traces.
func (t *Tools) MCPTools() []tools.MCPTool {
	return []tools.MCPTool{
		{
			Metadata: tools.NewMetadata("list_traces",
				mcp.WithDescription("List traces from a given query"),
				mcp.WithString("service",
					mcp.Description("Optional. Service to filter traces. Can not be used with trace_ids"),
				),
				mcp.WithString("operation",
					mcp.Description("Optional. Operation to filter traces. Can not be used with trace_ids"),
				),
				params.WithStringArray("trace_ids",
					mcp.Description("Optional. Trace IDs to filter traces. Can not be used with service or operation"),
				),
				params.WithTimeRange(),
			),
			Handler: func(ctx context.Context, request mcp.CallToolRequest) (*tools.Result, error) {
				timeRange, err := params.ParseTimeRange(request)
				if err != nil {
					return nil, err
				}

				service, err := params.String(request, "service", false, "")
				if err != nil {
					return nil, err
				}
				operation, err := params.String(request, "operation", false, "")
				if err != nil {
					return nil, err
				}
				traceIDs, err := params.StringArray(request, "trace_ids", false, nil)
				if err != nil {
					return nil, err
				}

				if len(traceIDs) > 0 && (service != "" || operation != "") {
					return nil, fmt.Errorf("trace_ids can not be used with service or operation")
				}

				queryParams := &data_unstable.ListTracesParams{
					Context: ctx,
					Body: &models.Datav1ListTracesRequest{
						StartTime: strfmt.DateTime(timeRange.Start),
						EndTime:   strfmt.DateTime(timeRange.End),
					},
				}

				if len(traceIDs) > 0 {
					queryParams.Body.TraceIds = traceIDs
					queryParams.Body.QueryType = "TRACE_IDS"
				} else {
					queryParams.Body.QueryType = "SERVICE_OPERATION"
				}

				if service != "" {
					queryParams.Body.Service = service
				}
				if operation != "" {
					queryParams.Body.Operation = operation
				}

				resp, err := t.api.DataUnstable.ListTraces(queryParams)
				if err != nil {
					return nil, fmt.Errorf("failed to list traces: %s", err)
				}
				return &tools.Result{
					JSONContent: resp,
				}, nil
			},
		},
	}
}
