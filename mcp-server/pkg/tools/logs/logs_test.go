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
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"github.com/chronosphereio/chronosphere-mcp/generated/datav1/datav1"
	"github.com/chronosphereio/chronosphere-mcp/generated/datav1/datav1/version1"
	"github.com/chronosphereio/chronosphere-mcp/generated/datav1/models"
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

	tools, err := NewTools(client, logger, linkBuilder)
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

func TestQuoteField(t *testing.T) {
	tests := []struct {
		name    string
		field   string
		want    string
		wantErr bool
	}{
		{name: "bare", field: "service", want: "service"},
		{name: "dotted path", field: "kubernetes.namespace_name", want: "kubernetes.namespace_name"},
		{name: "hyphen and slash", field: "app/k8s-name", want: "app/k8s-name"},
		{name: "leading at", field: "@timestamp", want: "@timestamp"},
		{name: "space", field: "http status", want: "['http status']"},
		{name: "colon", field: "k8s:pod", want: "['k8s:pod']"},
		{name: "leading hyphen", field: "-x", want: "['-x']"},
		{name: "empty segment", field: "a..b", want: "['a..b']"},
		{name: "keyword", field: "by", want: "['by']"},
		{name: "quote", field: "it's", wantErr: true},
		{name: "pipe", field: "a|b", wantErr: true},
		{name: "empty", field: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := quoteField(tt.field)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetLogQuery(t *testing.T) {
	tests := []struct {
		name string
		id   string
		want string
	}{
		{
			name: "plain id",
			id:   "abc123",
			want: `logID = "abc123" | project logID, _timestamp, service, severity, message, payload`,
		},
		{
			name: "id with quote is escaped",
			id:   `a"b`,
			want: `logID = "a\"b" | project logID, _timestamp, service, severity, message, payload`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, getLogQuery(tt.id))
		})
	}
}

func TestHistogramQuery(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		step    time.Duration
		groupBy []string
		want    string
		wantErr bool
	}{
		{
			name: "no query, no group by",
			step: 36 * time.Second,
			want: "make-series count() step 36s",
		},
		{
			name:  "query with group by",
			query: `service="gateway"`,
			step:  time.Minute,
			groupBy: []string{
				"service", "severity",
			},
			want: `service="gateway" | make-series count() step 60s by service, severity`,
		},
		{
			name:    "quoted group by field",
			query:   "error",
			step:    90 * time.Second,
			groupBy: []string{"http status"},
			want:    "error | make-series count() step 90s by ['http status']",
		},
		{
			name: "sub-second step rounds up to 1s",
			step: 300 * time.Millisecond,
			want: "make-series count() step 1s",
		},
		{
			name: "fractional step truncates to whole seconds",
			step: 1500 * time.Millisecond,
			want: "make-series count() step 1s",
		},
		{
			name:    "unquotable group by field",
			step:    time.Minute,
			groupBy: []string{"a|b"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := histogramQuery(tt.query, tt.step, tt.groupBy)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFieldValuesQuery(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		field   string
		limit   int
		want    string
		wantErr bool
	}{
		{
			name:  "no query",
			field: "service",
			limit: 100,
			want:  "filter service exists | summarize n = count() by service | sort by n desc | limit 100",
		},
		{
			name:  "with query",
			query: `severity="ERROR"`,
			field: "kubernetes.pod_name",
			limit: 10,
			want:  `severity="ERROR" | filter kubernetes.pod_name exists | summarize n = count() by kubernetes.pod_name | sort by n desc | limit 10`,
		},
		{
			name:  "quoted field",
			field: "http status",
			limit: 5,
			want:  "filter ['http status'] exists | summarize n = count() by ['http status'] | sort by n desc | limit 5",
		},
		{
			name:    "unquotable field",
			field:   `a"b`,
			limit:   5,
			wantErr: true,
		},
		{
			name:    "non-positive limit",
			field:   "service",
			limit:   0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fieldValuesQuery(tt.query, tt.field, tt.limit)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

// fakeLogsServer serves a fixed QueryLogsRange response and records the query parameters of each request.
type fakeLogsServer struct {
	*httptest.Server
	requests []url.Values
}

func newFakeLogsServer(t *testing.T, body string) *fakeLogsServer {
	f := &fakeLogsServer{}
	f.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/data/logs:range", r.URL.Path)
		f.requests = append(f.requests, r.URL.Query())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(body)); err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	}))
	t.Cleanup(f.Close)
	return f
}

func newTestTools(t *testing.T, serverURL string) *Tools {
	transport := datav1.DefaultTransportConfig().
		WithHost(serverURL[len("http://"):]).
		WithSchemes([]string{"http"})
	client := datav1.NewHTTPClientWithConfig(nil, transport)
	tools, err := NewTools(client, zaptest.NewLogger(t), links.NewBuilder("https://test.chronosphere.io"))
	require.NoError(t, err)
	return tools
}

func callTool(t *testing.T, tools *Tools, name string, args map[string]any) (any, error) {
	for _, tool := range tools.MCPTools() {
		if tool.Metadata.Name != name {
			continue
		}
		req := mcp.CallToolRequest{}
		req.Params.Name = name
		req.Params.Arguments = args
		res, err := tool.Handler(t.Context(), req)
		if err != nil {
			return nil, err
		}
		return res.JSONContent, nil
	}
	t.Fatalf("tool %q not registered", name)
	return nil, nil
}

func toJSON(t *testing.T, v any) string {
	b, err := json.Marshal(v)
	require.NoError(t, err)
	return string(b)
}

func TestGetLogTool(t *testing.T) {
	server := newFakeLogsServer(t, `{
		"grid_data": {
			"columns": [
				{"name": "logID", "type": "STRING"},
				{"name": "severity", "type": "STRING"},
				{"name": "message", "type": "STRING"},
				{"name": "payload", "type": "STRING"}
			],
			"rows": [{"values": [
				{"string_value": "abc"},
				{"string_value": "ERROR"},
				{"string_value": "boom"},
				{"string_value": "{\"k\":\"v\"}"}
			]}]
		}
	}`)
	tools := newTestTools(t, server.URL)

	t.Run("default 7 day window", func(t *testing.T) {
		got, err := callTool(t, tools, "get_log", map[string]any{"id": "abc"})
		require.NoError(t, err)
		assert.JSONEq(t, `{"logID":"abc","severity":"ERROR","message":"boom","payload":"{\"k\":\"v\"}"}`, toJSON(t, got))

		req := server.requests[len(server.requests)-1]
		assert.Equal(t, getLogQuery("abc"), req.Get("query"))
		after, err := time.Parse(time.RFC3339, req.Get("time_range.after"))
		require.NoError(t, err)
		before, err := time.Parse(time.RFC3339, req.Get("time_range.before"))
		require.NoError(t, err)
		assert.InDelta(t, getLogDefaultLookback.Seconds(), before.Sub(after).Seconds(), 1)
	})

	t.Run("explicit window", func(t *testing.T) {
		_, err := callTool(t, tools, "get_log", map[string]any{
			"id":    "abc",
			"start": "2026-01-01T00:00:00Z",
			"end":   "2026-01-02T00:00:00Z",
		})
		require.NoError(t, err)
		req := server.requests[len(server.requests)-1]
		after, err := time.Parse(time.RFC3339, req.Get("time_range.after"))
		require.NoError(t, err)
		before, err := time.Parse(time.RFC3339, req.Get("time_range.before"))
		require.NoError(t, err)
		assert.True(t, after.Equal(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)))
		assert.True(t, before.Equal(time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)))
	})
}

func TestGetLogToolNotFound(t *testing.T) {
	server := newFakeLogsServer(t, `{"grid_data": {"columns": [{"name": "logID"}], "rows": []}}`)
	tools := newTestTools(t, server.URL)

	_, err := callTool(t, tools, "get_log", map[string]any{"id": "missing"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestGetLogHistogramTool(t *testing.T) {
	server := newFakeLogsServer(t, `{
		"time_series_data": {
			"group_by_dimension_names": ["service"],
			"series": [{
				"aggregation_name": "count",
				"group_by_dimension_values": ["gateway"],
				"buckets": [{"start_time": "2026-01-01T00:00:00Z", "end_time": "2026-01-01T00:01:00Z", "float_value": 3}]
			}]
		}
	}`)
	tools := newTestTools(t, server.URL)

	got, err := callTool(t, tools, "get_log_histogram", map[string]any{
		"query":    `severity="ERROR"`,
		"start":    "2026-01-01T00:00:00Z",
		"end":      "2026-01-01T01:40:00Z",
		"group_by": []any{"service"},
	})
	require.NoError(t, err)

	require.Len(t, server.requests, 1)
	assert.Equal(t, `severity="ERROR" | make-series count() step 60s by service`, server.requests[0].Get("query"))

	ts, ok := got.(*models.QueryLogsRangeResponseTimeSeriesData)
	require.True(t, ok)
	require.Len(t, ts.Series, 1)
	assert.Equal(t, []string{"gateway"}, ts.Series[0].GroupByDimensionValues)
	require.Len(t, ts.Series[0].Buckets, 1)
	assert.Equal(t, float64(3), ts.Series[0].Buckets[0].FloatValue)
}

func TestListLogFieldValuesTool(t *testing.T) {
	tests := []struct {
		name string
		body string
		want string
	}{
		{
			name: "value column first",
			body: `{"grid_data": {
				"columns": [{"name": "service", "type": "STRING"}, {"name": "n", "type": "FLOAT"}],
				"rows": [
					{"values": [{"string_value": "gateway"}, {"float_value": 12}]},
					{"values": [{"string_value": "api"}, {"float_value": 3}]}
				]
			}}`,
			want: `{"values":[{"value":"gateway","count":12},{"value":"api","count":3}]}`,
		},
		{
			name: "count column first",
			body: `{"grid_data": {
				"columns": [{"name": "n", "type": "FLOAT"}, {"name": "service", "type": "STRING"}],
				"rows": [{"values": [{"float_value": 7}, {"string_value": "gateway"}]}]
			}}`,
			want: `{"values":[{"value":"gateway","count":7}]}`,
		},
		{
			name: "numeric field values",
			body: `{"grid_data": {
				"columns": [{"name": "status", "type": "FLOAT"}, {"name": "n", "type": "FLOAT"}],
				"rows": [{"values": [{"float_value": 500}, {"float_value": 2}]}]
			}}`,
			want: `{"values":[{"value":"500","count":2}]}`,
		},
		{
			name: "no rows",
			body: `{"grid_data": {"columns": [{"name": "service"}, {"name": "n"}], "rows": []}}`,
			want: `{"values":[]}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := newFakeLogsServer(t, tt.body)
			tools := newTestTools(t, server.URL)

			got, err := callTool(t, tools, "list_log_field_values", map[string]any{
				"query":      `severity="ERROR"`,
				"field_name": "service",
				"limit":      float64(5),
			})
			require.NoError(t, err)
			assert.JSONEq(t, tt.want, toJSON(t, got))

			require.Len(t, server.requests, 1)
			assert.Equal(t,
				`severity="ERROR" | filter service exists | summarize n = count() by service | sort by n desc | limit 5`,
				server.requests[0].Get("query"))
		})
	}
}

func TestListLogFieldNamesRetired(t *testing.T) {
	const withQuery = "list_log_field_names is no longer supported.\n" +
		"To approximate it:\n" +
		"1. Call query_logs_range with query `service=\"gateway\" | project logID, payload | limit 20` over the same time range.\n" +
		"2. Collect the keys that appear across the returned payloads. Fields present on more of the sampled logs are\n" +
		"   the common ones.\n" +
		"3. To confirm a candidate field and see its values, call list_log_field_values with that field_name.\n" +
		"This finds fields in the raw log body. It can miss labels added by the ingestion pipeline and OTLP attributes\n" +
		"that are not in the payload, so if you know a field name, try list_log_field_values directly."
	const withoutQuery = "list_log_field_names is no longer supported.\n" +
		"To approximate it:\n" +
		"1. Call query_logs_range with query `project logID, payload | limit 20` over the same time range.\n" +
		"2. Collect the keys that appear across the returned payloads. Fields present on more of the sampled logs are\n" +
		"   the common ones.\n" +
		"3. To confirm a candidate field and see its values, call list_log_field_values with that field_name.\n" +
		"This finds fields in the raw log body. It can miss labels added by the ingestion pipeline and OTLP attributes\n" +
		"that are not in the payload, so if you know a field name, try list_log_field_values directly."

	tests := []struct {
		name    string
		args    map[string]any
		wantErr string
	}{
		{
			name:    "with query",
			args:    map[string]any{"query": `service="gateway"`, "start": "2026-01-01T00:00:00Z", "limit": float64(10)},
			wantErr: withQuery,
		},
		{
			name:    "without query",
			args:    map[string]any{},
			wantErr: withoutQuery,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := newFakeLogsServer(t, `{}`)
			tools := newTestTools(t, server.URL)

			var found bool
			for _, tool := range tools.MCPTools() {
				if tool.Metadata.Name != "list_log_field_names" {
					continue
				}
				found = true
				assert.True(t, strings.HasPrefix(tool.Metadata.Description, "Deprecated: no longer supported, always returns an error."))
				for _, p := range []string{"query", "start", "end", "limit"} {
					assert.Contains(t, tool.Metadata.InputSchema.Properties, p)
				}
			}
			require.True(t, found, "list_log_field_names must stay registered")

			_, err := callTool(t, tools, "list_log_field_names", tt.args)
			require.Error(t, err)
			assert.Equal(t, tt.wantErr, err.Error())
			assert.Empty(t, server.requests, "list_log_field_names must not call the API")
		})
	}
}

func TestCreateCompactSummary(t *testing.T) {
	resp := &models.Datav1QueryLogsRangeResponse{
		GridData: &models.QueryLogsRangeResponseGridData{
			Columns: []*models.QueryLogsRangeResponseColumnMeta{{Name: "logID"}, {Name: "message"}},
			Rows: []*models.QueryLogsRangeResponseRow{
				{Values: []*models.QueryLogsRangeResponseRowValue{{StringValue: "a"}, {StringValue: "x|y"}}},
			},
		},
		Metadata: &models.QueryLogsRangeResponseMetadata{LimitEnforced: true},
	}

	got := createCompactSummary(resp)
	assert.Equal(t, 1, got.TotalLogs)
	assert.Equal(t, []string{"logID", "message"}, got.Columns)
	assert.Equal(t, []string{"logID|message", "a|x￨y"}, got.Logs)
	assert.Equal(t, "Found 1 log entries. Columns: logID, message. All entries shown below in pipe-separated format", got.Summary)
	assert.True(t, got.Metadata.LimitEnforced)

	encoded := toJSON(t, got)
	assert.NotContains(t, encoded, "available_fields")
	assert.NotContains(t, encoded, "field_suggestion")
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trimmed, _ := trimLogEntries(tt.resp, 10, 0)
			summary := createCompactSummary(trimmed)
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

	server := newFakeLogsServer(t, `{"metadata":{"page":{}},"time_series_data":{}}`)
	tools := newTestTools(t, server.URL)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := map[string]any{"query": `"nonexistent" | make-series count() step 5m`}
			for k, v := range tt.args {
				args[k] = v
			}
			got, err := callTool(t, tools, "query_logs_range", args)
			require.NoError(t, err)
			summary, ok := got.(*CompactLogSummary)
			require.True(t, ok)
			assert.Equal(t, "No log data returned", summary.Summary)
		})
	}
}
