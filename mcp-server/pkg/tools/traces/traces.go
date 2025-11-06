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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"

	"github.com/go-openapi/strfmt"
	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"

	"github.com/chronosphereio/chronosphere-mcp/generated/datav1/datav1"
	"github.com/chronosphereio/chronosphere-mcp/generated/datav1/models"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools/pkg/params"
)

var _ tools.MCPTools = (*Tools)(nil)

type Tools struct {
	logger     *zap.Logger
	api        *datav1.DataV1API
	httpClient *http.Client
	baseURL    string
}

func NewTools(
	api *datav1.DataV1API,
	logger *zap.Logger,
) (*Tools, error) {
	logger.Info("traces tool configured")

	// Extract HTTP client and base URL from the transport for direct API calls
	// This is needed to work around go-swagger's broken allOf handling
	//
	// The transport is typically *httptransport.Runtime which has unexported fields,
	// so we need to use reflection to extract the HTTP client
	httpClient := http.DefaultClient
	baseURL := "https://" + datav1.DefaultHost + datav1.DefaultBasePath

	// Try to extract the client using reflection from the httptransport.Runtime
	transportValue := reflect.ValueOf(api.Transport)
	if transportValue.Kind() == reflect.Ptr {
		transportElem := transportValue.Elem()

		// Try to get the "client" field
		clientField := transportElem.FieldByName("Client")
		if clientField.IsValid() && clientField.CanInterface() {
			if client, ok := clientField.Interface().(*http.Client); ok {
				httpClient = client
			}
		}

		// Try to get the host
		hostField := transportElem.FieldByName("Host")
		if hostField.IsValid() && hostField.CanInterface() {
			if host, ok := hostField.Interface().(string); ok && host != "" {
				baseURL = "https://" + host + datav1.DefaultBasePath
			}
		}
	}

	return &Tools{
		logger:     logger,
		api:        api,
		httpClient: httpClient,
		baseURL:    baseURL,
	}, nil
}

func (t *Tools) GroupName() string {
	return "traces"
}

// listTracesRequestBody wraps the generated model and provides correct JSON marshaling
// to work around go-swagger's broken allOf handling which generates an embedded anonymous
// struct that doesn't marshal correctly to JSON.
type listTracesRequestBody struct {
	models.Datav1ListTracesRequest
	queryType models.ListTracesRequestQueryType
}

// MarshalBinary implements the encoding.BinaryMarshaler interface
func (m *listTracesRequestBody) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}

	// Create a properly structured JSON request with the query_type as a plain string
	// instead of the broken embedded struct that go-swagger generates
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

// Validate implements the validation interface to satisfy go-swagger requirements
func (m *listTracesRequestBody) Validate(formats strfmt.Registry) error {
	if m == nil {
		return nil
	}
	return m.Datav1ListTracesRequest.Validate(formats)
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

				// Determine query type
				var queryType models.ListTracesRequestQueryType
				if len(traceIDs) > 0 {
					queryType = models.ListTracesRequestQueryTypeTRACEIDS
				} else {
					queryType = models.ListTracesRequestQueryTypeSERVICEOPERATION
				}

				// Create request body with correct marshaling
				// We marshal it ourselves to work around go-swagger's broken allOf handling
				body := &listTracesRequestBody{
					Datav1ListTracesRequest: models.Datav1ListTracesRequest{
						StartTime: strfmt.DateTime(timeRange.Start),
						EndTime:   strfmt.DateTime(timeRange.End),
						Service:   service,
						Operation: operation,
						TraceIds:  traceIDs,
					},
					queryType: queryType,
				}

				// Marshal the body ourselves to ensure correct JSON format
				bodyJSON, err := body.MarshalBinary()
				if err != nil {
					return nil, fmt.Errorf("failed to marshal request: %w", err)
				}

				// Make a direct HTTP call since we can't use the generated client
				// with our custom marshaling due to go-swagger's broken allOf handling
				httpReq, err := http.NewRequestWithContext(ctx, "POST",
					t.baseURL+"/traces",
					bytes.NewReader(bodyJSON))
				if err != nil {
					return nil, fmt.Errorf("failed to create request: %w", err)
				}
				httpReq.Header.Set("Content-Type", "application/json")

				// Make the request
				httpResp, err := t.httpClient.Do(httpReq)
				if err != nil {
					return nil, fmt.Errorf("failed to execute request: %w", err)
				}
				defer httpResp.Body.Close()

				if httpResp.StatusCode != http.StatusOK {
					bodyBytes, _ := io.ReadAll(httpResp.Body)
					return nil, fmt.Errorf("failed to list traces: %s (status %d)", string(bodyBytes), httpResp.StatusCode)
				}

				// Parse the response
				var listResp models.Datav1ListTracesResponse
				if err := json.NewDecoder(httpResp.Body).Decode(&listResp); err != nil {
					return nil, fmt.Errorf("failed to decode response: %w", err)
				}

				return &tools.Result{
					JSONContent: &listResp,
				}, nil
			},
		},
	}
}
