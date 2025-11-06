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
	"encoding/json"
	"fmt"
	"unsafe"

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

// listTracesRequestBody wraps the generated model and provides correct JSON marshaling
// to work around go-swagger's broken allOf handling which generates an embedded anonymous
// struct that doesn't marshal correctly to JSON.
type listTracesRequestBody struct {
	*models.Datav1ListTracesRequest
	queryType models.ListTracesRequestQueryType
}

// MarshalBinary implements the encoding.BinaryMarshaler interface
func (m *listTracesRequestBody) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}

	// Create a properly structured JSON request
	type flatRequest struct {
		StartTime  strfmt.DateTime                          `json:"start_time,omitempty"`
		EndTime    strfmt.DateTime                          `json:"end_time,omitempty"`
		QueryType  models.ListTracesRequestQueryType        `json:"query_type,omitempty"`
		Service    string                                   `json:"service,omitempty"`
		Operation  string                                   `json:"operation,omitempty"`
		TraceIds   []string                                 `json:"trace_ids,omitempty"`
		TagFilters []*models.ListTracesRequestTagFilter     `json:"tag_filters,omitempty"`
	}

	flat := &flatRequest{
		StartTime:  m.Datav1ListTracesRequest.StartTime,
		EndTime:    m.Datav1ListTracesRequest.EndTime,
		QueryType:  m.queryType,
		Service:    m.Datav1ListTracesRequest.Service,
		Operation:  m.Datav1ListTracesRequest.Operation,
		TraceIds:   m.Datav1ListTracesRequest.TraceIds,
		TagFilters: m.Datav1ListTracesRequest.TagFilters,
	}

	return json.Marshal(flat)
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

				// Create wrapped request body with correct marshaling
				body := &listTracesRequestBody{
					Datav1ListTracesRequest: &models.Datav1ListTracesRequest{
						StartTime: strfmt.DateTime(timeRange.Start),
						EndTime:   strfmt.DateTime(timeRange.End),
						Service:   service,
						Operation: operation,
						TraceIds:  traceIDs,
					},
				}

				if len(traceIDs) > 0 {
					body.queryType = models.ListTracesRequestQueryTypeTRACEIDS
				} else {
					body.queryType = models.ListTracesRequestQueryTypeSERVICEOPERATION
				}

				queryParams := &version1.ListTracesParams{
					Context: ctx,
				}

				// Use unsafe to bypass type checking and inject our custom marshaler
				// This is necessary because go-swagger generates broken code for allOf with enums
				queryParams.Body = (*models.Datav1ListTracesRequest)(unsafe.Pointer(body))

				resp, err := t.api.Version1.ListTraces(queryParams, nil)
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
