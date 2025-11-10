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

// Package traces contains tools for querying Chronosphere traces.
package traces

import (
	"context"
	"fmt"

	"github.com/go-openapi/strfmt"
	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"

	"github.com/chronosphereio/chronosphere-mcp/generated/datav1/datav1"
	"github.com/chronosphereio/chronosphere-mcp/generated/datav1/datav1/version1"
	"github.com/chronosphereio/chronosphere-mcp/generated/datav1/models"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools/pkg/params"
)

var _ tools.MCPTools = (*Tools)(nil)

type Tools struct {
	logger *zap.Logger
	api    *datav1.DataV1API
}

func NewTools(
	api *datav1.DataV1API,
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

				queryParams := &version1.ListTracesParams{
					Context: ctx,
					Body: &models.Datav1ListTracesRequest{
						StartTime: strfmt.DateTime(timeRange.Start),
						EndTime:   strfmt.DateTime(timeRange.End),
					},
				}

				if len(traceIDs) > 0 {
					queryParams.Body.TraceIds = traceIDs
					queryParams.Body.QueryType = models.ListTracesRequestQueryTypeTRACEIDS
				} else {
					queryParams.Body.QueryType = models.ListTracesRequestQueryTypeSERVICEOPERATION
				}

				if service != "" {
					queryParams.Body.Service = service
				}
				if operation != "" {
					queryParams.Body.Operation = operation
				}

				resp, err := t.api.Version1.ListTraces(queryParams)
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
