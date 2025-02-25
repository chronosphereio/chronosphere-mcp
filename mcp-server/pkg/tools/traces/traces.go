package traces

import (
	"fmt"

	"github.com/go-openapi/strfmt"
	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"

	"github.com/chronosphereio/mcp-server/generated/dataunstable/dataunstable/data_unstable"
	"github.com/chronosphereio/mcp-server/mcp-server/pkg/client"
	"github.com/chronosphereio/mcp-server/mcp-server/pkg/ptr"
	"github.com/chronosphereio/mcp-server/mcp-server/pkg/tools"
	"github.com/chronosphereio/mcp-server/mcp-server/pkg/tools/pkg/params"
)

var _ tools.MCPTools = (*Tools)(nil)

type Tools struct {
	logger         *zap.Logger
	clientProvider *client.Provider
}

func NewTools(
	clientProvider *client.Provider,
	logger *zap.Logger,
) (*Tools, error) {
	logger.Info("events tool configured")

	return &Tools{
		logger:         logger,
		clientProvider: clientProvider,
	}, nil
}

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
				mcp.WithArray("trace_ids",
					mcp.Description("Optional. Trace IDs to filter traces. Can not be used with service or operation"),
				),
				params.WithTimeRange(),
			),
			Handler: func(session tools.Session, request mcp.CallToolRequest) (*tools.Result, error) {
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
					StartTime: (*strfmt.DateTime)(ptr.To(timeRange.Start)),
					EndTime:   (*strfmt.DateTime)(ptr.To(timeRange.End)),
				}

				if len(traceIDs) > 0 {
					queryParams.TraceIds = traceIDs
					queryParams.QueryType = ptr.To("TRACE_IDS")
				} else {
					queryParams.QueryType = ptr.To("SERVICE_OPERATION")
				}

				if service != "" {
					queryParams.Service = &service
				}
				if operation != "" {
					queryParams.Operation = &operation
				}

				api, err := t.clientProvider.DataUnstableClient(session)
				if err != nil {
					return nil, err
				}
				resp, err := api.ListTraces(queryParams)
				if err != nil {
					return nil, fmt.Errorf("failed to list traces: %s", err)
				}
				return &tools.Result{
					JsonContent: resp,
				}, nil
			},
		},
	}
}
