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

package traces

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"github.com/chronosphereio/chronosphere-mcp/generated/datav1/datav1"
	"github.com/chronosphereio/chronosphere-mcp/generated/datav1/models"
)

func TestTrimTraces(t *testing.T) {
	tests := []struct {
		name         string
		payload      *models.Datav1ListTracesResponse
		limit        int
		offset       int
		wantTrimmed  int
		validateFunc func(t *testing.T, result *models.Datav1ListTracesResponse)
	}{
		{
			name:        "nil payload",
			payload:     nil,
			limit:       10,
			offset:      0,
			wantTrimmed: 0,
			validateFunc: func(t *testing.T, result *models.Datav1ListTracesResponse) {
				assert.Nil(t, result)
			},
		},
		{
			name: "zero limit and offset returns original",
			payload: &models.Datav1ListTracesResponse{
				Traces: []*models.V1TracesData{{}, {}, {}},
			},
			limit:       0,
			offset:      0,
			wantTrimmed: 0,
			validateFunc: func(t *testing.T, result *models.Datav1ListTracesResponse) {
				assert.Equal(t, 3, len(result.Traces))
			},
		},
		{
			name: "limit without offset",
			payload: &models.Datav1ListTracesResponse{
				Traces: []*models.V1TracesData{{}, {}, {}, {}, {}},
			},
			limit:       3,
			offset:      0,
			wantTrimmed: 2,
			validateFunc: func(t *testing.T, result *models.Datav1ListTracesResponse) {
				assert.NotNil(t, result)
				assert.Equal(t, 3, len(result.Traces))
			},
		},
		{
			name: "offset and limit - pagination",
			payload: &models.Datav1ListTracesResponse{
				Traces: []*models.V1TracesData{{}, {}, {}, {}, {}},
			},
			limit:       2,
			offset:      2,
			wantTrimmed: 3,
			validateFunc: func(t *testing.T, result *models.Datav1ListTracesResponse) {
				assert.Equal(t, 2, len(result.Traces))
			},
		},
		{
			name: "offset beyond total traces",
			payload: &models.Datav1ListTracesResponse{
				Traces: []*models.V1TracesData{{}, {}, {}},
			},
			limit:       5,
			offset:      10,
			wantTrimmed: 3,
			validateFunc: func(t *testing.T, result *models.Datav1ListTracesResponse) {
				assert.Equal(t, 0, len(result.Traces))
			},
		},
		{
			name: "limit exceeds remaining traces",
			payload: &models.Datav1ListTracesResponse{
				Traces: []*models.V1TracesData{{}, {}, {}, {}, {}},
			},
			limit:       10,
			offset:      3,
			wantTrimmed: 3,
			validateFunc: func(t *testing.T, result *models.Datav1ListTracesResponse) {
				assert.Equal(t, 2, len(result.Traces))
			},
		},
		{
			name: "single trace with offset 0 and limit 1",
			payload: &models.Datav1ListTracesResponse{
				Traces: []*models.V1TracesData{{}, {}, {}},
			},
			limit:       1,
			offset:      0,
			wantTrimmed: 2,
			validateFunc: func(t *testing.T, result *models.Datav1ListTracesResponse) {
				assert.Equal(t, 1, len(result.Traces))
			},
		},
		{
			name: "nil traces slice",
			payload: &models.Datav1ListTracesResponse{
				Traces: nil,
			},
			limit:       5,
			offset:      0,
			wantTrimmed: 0,
			validateFunc: func(t *testing.T, result *models.Datav1ListTracesResponse) {
				assert.Equal(t, result, &models.Datav1ListTracesResponse{
					Traces: nil,
				})
			},
		},
		{
			name: "empty traces slice",
			payload: &models.Datav1ListTracesResponse{
				Traces: []*models.V1TracesData{},
			},
			limit:       5,
			offset:      0,
			wantTrimmed: 0,
			validateFunc: func(t *testing.T, result *models.Datav1ListTracesResponse) {
				assert.Equal(t, 0, len(result.Traces))
			},
		},
		{
			name: "only offset, no limit",
			payload: &models.Datav1ListTracesResponse{
				Traces: []*models.V1TracesData{{}, {}, {}, {}, {}},
			},
			limit:       0,
			offset:      2,
			wantTrimmed: 2,
			validateFunc: func(t *testing.T, result *models.Datav1ListTracesResponse) {
				// With limit=0, it should return all traces from offset onwards
				assert.Equal(t, 3, len(result.Traces))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, trimmed := trimTraces(tt.payload, tt.limit, tt.offset)
			assert.Equal(t, tt.wantTrimmed, trimmed)
			tt.validateFunc(t, result)
		})
	}
}

// TestListTracesQueryTypeFields tests that StartTime and EndTime are only set
// for SERVICE_OPERATION query type, not for TRACE_IDS query type.
// This is a regression test for https://github.com/chronosphereio/chronosphere-mcp/issues/125
func TestListTracesQueryTypeFields(t *testing.T) {
	tests := []struct {
		name               string
		requestArgs        map[string]any
		expectStartEndTime bool
		expectQueryType    string
		expectTraceIDs     []string
		expectService      string
		expectOperation    string
	}{
		{
			name: "trace_ids query type should not include StartTime/EndTime",
			requestArgs: map[string]any{
				"trace_ids":  []any{"trace-1", "trace-2"},
				"start_time": "2025-01-01T00:00:00Z",
				"end_time":   "2025-01-01T01:00:00Z",
			},
			expectStartEndTime: false,
			expectQueryType:    string(models.ListTracesRequestQueryTypeTRACEIDS),
			expectTraceIDs:     []string{"trace-1", "trace-2"},
		},
		{
			name: "service/operation query type should include StartTime/EndTime",
			requestArgs: map[string]any{
				"service":    "my-service",
				"operation":  "my-operation",
				"start_time": "2025-01-01T00:00:00Z",
				"end_time":   "2025-01-01T01:00:00Z",
			},
			expectStartEndTime: true,
			expectQueryType:    string(models.ListTracesRequestQueryTypeSERVICEOPERATION),
			expectService:      "my-service",
			expectOperation:    "my-operation",
		},
		{
			name: "service only query type should include StartTime/EndTime",
			requestArgs: map[string]any{
				"service":    "my-service",
				"start_time": "2025-01-01T00:00:00Z",
				"end_time":   "2025-01-01T01:00:00Z",
			},
			expectStartEndTime: true,
			expectQueryType:    string(models.ListTracesRequestQueryTypeSERVICEOPERATION),
			expectService:      "my-service",
		},
		{
			name: "operation only query type should include StartTime/EndTime",
			requestArgs: map[string]any{
				"operation":  "my-operation",
				"start_time": "2025-01-01T00:00:00Z",
				"end_time":   "2025-01-01T01:00:00Z",
			},
			expectStartEndTime: true,
			expectQueryType:    string(models.ListTracesRequestQueryTypeSERVICEOPERATION),
			expectOperation:    "my-operation",
		},
		{
			name: "no filters query type should include StartTime/EndTime",
			requestArgs: map[string]any{
				"start_time": "2025-01-01T00:00:00Z",
				"end_time":   "2025-01-01T01:00:00Z",
			},
			expectStartEndTime: true,
			expectQueryType:    string(models.ListTracesRequestQueryTypeSERVICEOPERATION),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var capturedRequest *models.Datav1ListTracesRequest

			// Create a test server that captures the request body
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Capture the request body
				body, err := io.ReadAll(r.Body)
				require.NoError(t, err)

				capturedRequest = &models.Datav1ListTracesRequest{}
				err = json.Unmarshal(body, capturedRequest)
				require.NoError(t, err)

				// Return a valid response
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, err = w.Write([]byte(`{"traces": []}`))
				require.NoError(t, err)
			}))
			defer server.Close()

			// Create client pointing to test server
			transport := datav1.DefaultTransportConfig().
				WithHost(server.URL[7:]). // Remove "http://"
				WithSchemes([]string{"http"})

			client := datav1.NewHTTPClientWithConfig(nil, transport)
			logger := zaptest.NewLogger(t)

			tools, err := NewTools(client, logger)
			require.NoError(t, err)

			// Get the list_traces tool
			mcpTools := tools.MCPTools()
			require.Len(t, mcpTools, 1)
			listTracesTool := mcpTools[0]

			// Create the request
			request := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name:      "list_traces",
					Arguments: tt.requestArgs,
				},
			}

			// Call the handler
			_, err = listTracesTool.Handler(context.Background(), request)
			require.NoError(t, err)

			// Verify the captured request
			require.NotNil(t, capturedRequest)
			assert.Equal(t, tt.expectQueryType, string(capturedRequest.QueryType))

			// Check StartTime and EndTime
			if tt.expectStartEndTime {
				assert.NotEqual(t, time.Time{}, time.Time(capturedRequest.StartTime),
					"StartTime should be set for %s query type", tt.expectQueryType)
				assert.NotEqual(t, time.Time{}, time.Time(capturedRequest.EndTime),
					"EndTime should be set for %s query type", tt.expectQueryType)
			} else {
				assert.Equal(t, time.Time{}, time.Time(capturedRequest.StartTime),
					"StartTime should NOT be set for %s query type", tt.expectQueryType)
				assert.Equal(t, time.Time{}, time.Time(capturedRequest.EndTime),
					"EndTime should NOT be set for %s query type", tt.expectQueryType)
			}

			// Check trace_ids
			if len(tt.expectTraceIDs) > 0 {
				assert.Equal(t, tt.expectTraceIDs, capturedRequest.TraceIds)
			}

			// Check service
			if tt.expectService != "" {
				assert.Equal(t, tt.expectService, capturedRequest.Service)
			}

			// Check operation
			if tt.expectOperation != "" {
				assert.Equal(t, tt.expectOperation, capturedRequest.Operation)
			}
		})
	}
}
