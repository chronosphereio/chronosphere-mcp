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

package logs

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"github.com/chronosphereio/chronosphere-mcp/generated/dataunstable/dataunstable"
	"github.com/chronosphereio/chronosphere-mcp/generated/datav1/datav1"
	"github.com/chronosphereio/chronosphere-mcp/generated/datav1/datav1/version1"
	"github.com/chronosphereio/chronosphere-mcp/generated/datav1/models"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools/pkg/params"
	"github.com/chronosphereio/chronosphere-mcp/pkg/links"
)

func TestTrimLogEntries(t *testing.T) {
	tests := []struct {
		name         string
		payload      *models.Datav1QueryLogsRangeResponse
		limit        int
		offset       int
		wantTrimmed  int
		validateFunc func(t *testing.T, result *models.Datav1QueryLogsRangeResponse)
	}{
		{
			name:        "nil payload",
			payload:     nil,
			limit:       10,
			offset:      0,
			wantTrimmed: 0,
			validateFunc: func(t *testing.T, result *models.Datav1QueryLogsRangeResponse) {
				assert.Nil(t, result)
			},
		},
		{
			name: "zero limit and offset returns original",
			payload: &models.Datav1QueryLogsRangeResponse{
				GridData: &models.QueryLogsRangeResponseGridData{
					Rows: []*models.QueryLogsRangeResponseRow{{}, {}, {}},
				},
			},
			limit:       0,
			offset:      0,
			wantTrimmed: 0,
			validateFunc: func(t *testing.T, result *models.Datav1QueryLogsRangeResponse) {
				assert.Equal(t, 3, len(result.GridData.Rows))
			},
		},
		{
			name: "grid data - limit without offset",
			payload: &models.Datav1QueryLogsRangeResponse{
				GridData: &models.QueryLogsRangeResponseGridData{
					Columns: []*models.QueryLogsRangeResponseColumnMeta{{Name: "col1"}},
					Rows:    []*models.QueryLogsRangeResponseRow{{}, {}, {}, {}, {}},
				},
				Metadata: &models.QueryLogsRangeResponseMetadata{},
			},
			limit:       3,
			offset:      0,
			wantTrimmed: 2,
			validateFunc: func(t *testing.T, result *models.Datav1QueryLogsRangeResponse) {
				assert.NotNil(t, result.GridData)
				assert.Equal(t, 3, len(result.GridData.Rows))
				assert.Equal(t, 1, len(result.GridData.Columns))
			},
		},
		{
			name: "grid data - offset and limit",
			payload: &models.Datav1QueryLogsRangeResponse{
				GridData: &models.QueryLogsRangeResponseGridData{
					Rows: []*models.QueryLogsRangeResponseRow{{}, {}, {}, {}, {}},
				},
			},
			limit:       2,
			offset:      2,
			wantTrimmed: 3,
			validateFunc: func(t *testing.T, result *models.Datav1QueryLogsRangeResponse) {
				assert.Equal(t, 2, len(result.GridData.Rows))
			},
		},
		{
			name: "grid data - offset beyond total rows",
			payload: &models.Datav1QueryLogsRangeResponse{
				GridData: &models.QueryLogsRangeResponseGridData{
					Rows: []*models.QueryLogsRangeResponseRow{{}, {}, {}},
				},
			},
			limit:       5,
			offset:      10,
			wantTrimmed: 3,
			validateFunc: func(t *testing.T, result *models.Datav1QueryLogsRangeResponse) {
				assert.Equal(t, 0, len(result.GridData.Rows))
			},
		},
		{
			name: "time series data - limit without offset",
			payload: &models.Datav1QueryLogsRangeResponse{
				TimeSeriesData: &models.QueryLogsRangeResponseTimeSeriesData{
					GroupByDimensionNames: []string{"dim1"},
					Series: []*models.TimeSeriesDataTimeSeries{
						{
							AggregationName:        "count",
							GroupByDimensionValues: []string{"val1"},
							Buckets: []*models.TimeSeriesBucket{
								{}, {}, {}, {}, {},
							},
						},
					},
				},
			},
			limit:       3,
			offset:      0,
			wantTrimmed: 2,
			validateFunc: func(t *testing.T, result *models.Datav1QueryLogsRangeResponse) {
				assert.NotNil(t, result.TimeSeriesData)
				assert.Equal(t, 1, len(result.TimeSeriesData.Series))
				assert.Equal(t, 3, len(result.TimeSeriesData.Series[0].Buckets))
				assert.Equal(t, []string{"dim1"}, result.TimeSeriesData.GroupByDimensionNames)
			},
		},
		{
			name: "time series data - multiple series",
			payload: &models.Datav1QueryLogsRangeResponse{
				TimeSeriesData: &models.QueryLogsRangeResponseTimeSeriesData{
					Series: []*models.TimeSeriesDataTimeSeries{
						{
							Buckets: []*models.TimeSeriesBucket{
								{}, {}, {}, {},
							},
						},
						{
							Buckets: []*models.TimeSeriesBucket{
								{}, {}, {},
							},
						},
					},
				},
			},
			limit:       2,
			offset:      1,
			wantTrimmed: 3,
			validateFunc: func(t *testing.T, result *models.Datav1QueryLogsRangeResponse) {
				assert.Equal(t, 2, len(result.TimeSeriesData.Series))
				assert.Equal(t, 2, len(result.TimeSeriesData.Series[0].Buckets))
				assert.Equal(t, 2, len(result.TimeSeriesData.Series[1].Buckets))
			},
		},
		{
			name: "neither grid nor time series data",
			payload: &models.Datav1QueryLogsRangeResponse{
				Metadata: &models.QueryLogsRangeResponseMetadata{},
			},
			limit:       5,
			offset:      0,
			wantTrimmed: 0,
			validateFunc: func(t *testing.T, result *models.Datav1QueryLogsRangeResponse) {
				assert.Nil(t, result.GridData)
				assert.Nil(t, result.TimeSeriesData)
				assert.NotNil(t, result.Metadata)
			},
		},
		{
			name: "empty time series data",
			payload: &models.Datav1QueryLogsRangeResponse{
				TimeSeriesData: &models.QueryLogsRangeResponseTimeSeriesData{},
			},
			limit:       5,
			offset:      1,
			wantTrimmed: 0,
			validateFunc: func(t *testing.T, result *models.Datav1QueryLogsRangeResponse) {
				require.NotNil(t, result.TimeSeriesData)
				assert.Empty(t, result.TimeSeriesData.Series)
			},
		},
		{
			name: "preserve metadata",
			payload: &models.Datav1QueryLogsRangeResponse{
				Metadata: &models.QueryLogsRangeResponseMetadata{
					LimitEnforced: true,
				},
				GridData: &models.QueryLogsRangeResponseGridData{
					Rows: []*models.QueryLogsRangeResponseRow{{}, {}},
				},
			},
			limit:       1,
			offset:      0,
			wantTrimmed: 1,
			validateFunc: func(t *testing.T, result *models.Datav1QueryLogsRangeResponse) {
				assert.True(t, result.Metadata.LimitEnforced)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, trimmed := trimLogEntries(tt.payload, tt.limit, tt.offset)
			assert.Equal(t, tt.wantTrimmed, trimmed)
			tt.validateFunc(t, result)
		})
	}
}

// TestQueryLogsRangeAPICall tests that QueryLogsRange can be called successfully.
func TestQueryLogsRangeAPICall(t *testing.T) {
	// Create a test server that returns a valid response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// Return a minimal valid response
		if _, err := w.Write([]byte(`{
			"gridData": {
				"columns": [{"name": "message"}],
				"rows": []
			},
			"metadata": {
				"page": {"token": ""}
			}
		}`)); err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	}))
	defer server.Close()

	// Create client pointing to test server
	transport := datav1.DefaultTransportConfig().
		WithHost(server.URL[7:]). // Remove "http://"
		WithSchemes([]string{"http"})

	client := datav1.NewHTTPClientWithConfig(nil, transport)
	logger := zaptest.NewLogger(t)
	linkBuilder := links.NewBuilder("https://test.chronosphere.io")

	tools, err := NewTools(client, nil, logger, linkBuilder)
	require.NoError(t, err)

	// Create query params
	now := time.Now()
	query := "test query"
	pageToken := ""

	queryParams := &version1.QueryLogsRangeParams{
		Context:         t.Context(),
		Query:           &query,
		TimeRangeAfter:  (*strfmt.DateTime)(&now),
		TimeRangeBefore: (*strfmt.DateTime)(&now),
		PageToken:       &pageToken,
	}

	// Call the API
	resp, err := tools.dataV1API.Version1.QueryLogsRange(queryParams)

	// We expect success since our test server returns a valid response
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Payload)
}

// newTestTools returns Tools whose v1 and unstable clients both talk to a test
// server that answers logs:range with rangeBody and 404s everything else.
func newTestTools(t *testing.T, rangeBody string) *Tools {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "logs:range") {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(rangeBody))
		assert.NoError(t, err)
	}))
	t.Cleanup(server.Close)

	host := strings.TrimPrefix(server.URL, "http://")
	dataV1 := datav1.NewHTTPClientWithConfig(nil, datav1.DefaultTransportConfig().WithHost(host).WithSchemes([]string{"http"}))
	dataUnstable := dataunstable.NewHTTPClientWithConfig(nil, dataunstable.DefaultTransportConfig().WithHost(host).WithSchemes([]string{"http"}))

	tools, err := NewTools(dataV1, dataUnstable, zaptest.NewLogger(t), links.NewBuilder("https://test.chronosphere.io"))
	require.NoError(t, err)
	return tools
}

func TestCreateCompactSummaryNoData(t *testing.T) {
	tests := []struct {
		name string
		resp *models.Datav1QueryLogsRangeResponse
	}{
		{
			name: "nil response",
			resp: nil,
		},
		{
			name: "neither grid nor time series data",
			resp: &models.Datav1QueryLogsRangeResponse{
				Metadata: &models.QueryLogsRangeResponseMetadata{},
			},
		},
		{
			name: "no metadata",
			resp: &models.Datav1QueryLogsRangeResponse{},
		},
		{
			name: "empty time series data",
			resp: &models.Datav1QueryLogsRangeResponse{
				Metadata:       &models.QueryLogsRangeResponseMetadata{},
				TimeSeriesData: &models.QueryLogsRangeResponseTimeSeriesData{},
			},
		},
	}

	tools := newTestTools(t, `{}`)
	now := time.Now()
	timeRange := &params.TimeRange{Start: now.Add(-time.Hour), End: now}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trimmed, _ := trimLogEntries(tt.resp, 10, 0)
			summary := tools.createCompactSummary(t.Context(), "test", timeRange, trimmed)
			assert.Equal(t, "No log data returned", summary.Summary)
			assert.Equal(t, 0, summary.TotalLogs)
		})
	}
}

// TestQueryLogsRangeEmptyMakeSeries covers the response the logs API returns
// for a make-series query that matches no logs. With a limit or offset, the
// tool used to drop time_series_data while trimming and then panic.
func TestQueryLogsRangeEmptyMakeSeries(t *testing.T) {
	tests := []struct {
		name string
		args map[string]any
	}{
		{
			name: "no limit",
			args: map[string]any{},
		},
		{
			name: "limit",
			args: map[string]any{"limit": float64(10)},
		},
		{
			name: "offset",
			args: map[string]any{"offset": float64(1)},
		},
	}

	tools := newTestTools(t, `{"metadata":{"page":{}},"time_series_data":{}}`)
	handler := tools.MCPTools()[0].Handler

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := map[string]any{"query": `"nonexistent" | make-series count() step 5m`}
			for k, v := range tt.args {
				args[k] = v
			}
			result, err := handler(t.Context(), mcp.CallToolRequest{
				Params: mcp.CallToolParams{Name: "query_logs_range", Arguments: args},
			})
			require.NoError(t, err)
			summary, ok := result.JSONContent.(*CompactLogSummary)
			require.True(t, ok)
			assert.Equal(t, "No log data returned", summary.Summary)
		})
	}
}
