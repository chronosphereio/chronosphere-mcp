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

package events

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"github.com/chronosphereio/chronosphere-mcp/generated/datav1/datav1"
	"github.com/chronosphereio/chronosphere-mcp/generated/datav1/datav1/version1"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools"
	"github.com/chronosphereio/chronosphere-mcp/pkg/links"
	"github.com/chronosphereio/chronosphere-mcp/pkg/ptr"
)

// TestListEventsAPICall tests that ListEvents can be called successfully.
func TestListEventsAPICall(t *testing.T) {
	// Create a test server that returns a valid response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// Return a minimal valid response
		if _, err := w.Write([]byte(`{
			"events": [],
			"page": {"token": ""}
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

	eventTools, err := NewTools(client, nil, logger, linkBuilder)
	require.NoError(t, err)

	// Create query params
	now := time.Now()
	queryParams := &version1.ListEventsParams{
		Context:        t.Context(),
		HappenedAfter:  (*strfmt.DateTime)(ptr.To(now)),
		HappenedBefore: (*strfmt.DateTime)(ptr.To(now)),
		Query:          ptr.To("test query"),
	}

	// Call the API
	resp, err := eventTools.dataV1API.Version1.ListEvents(queryParams)

	// We expect success since our test server returns a valid response
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Payload)
}

func TestEventFieldValuesTools(t *testing.T) {
	valuesByField := map[string][]string{
		"CATEGORY":     {"deploy", "alert"},
		"TYPE":         {"rollout"},
		"SOURCE":       {"github"},
		"LABEL_NAME":   {"cluster", "service"},
		"LABEL_VALUE":  {"prod", "staging"},
		"LENS_SERVICE": {"gateway"},
	}

	tests := []struct {
		name          string
		tool          string
		args          map[string]any
		wantQueries   []map[string]string
		wantJSON      string
		wantChronLink string
	}{
		{
			name: "get_events_metadata lists each field once",
			tool: "get_events_metadata",
			wantQueries: []map[string]string{
				{"field": "CATEGORY"},
				{"field": "LABEL_NAME"},
				{"field": "SOURCE"},
				{"field": "TYPE"},
				{"field": "LENS_SERVICE"},
			},
			wantJSON: `{
				"categories": ["deploy", "alert"],
				"sources": ["github"],
				"types": ["rollout"],
				"label_names": ["cluster", "service"],
				"lens_service_names": ["gateway"]
			}`,
			wantChronLink: "https://test.chronosphere.io/events/explorer",
		},
		{
			name: "list_events_label_values lists values of one label",
			tool: "list_events_label_values",
			args: map[string]any{"label_name": "cluster"},
			wantQueries: []map[string]string{
				{"field": "LABEL_VALUE", "label_name": "cluster"},
			},
			wantJSON:      `{"Payload": {"values": ["prod", "staging"]}}`,
			wantChronLink: "https://test.chronosphere.io/api/v1/data/events:field-values?field=LABEL_VALUE&label_name=cluster",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotQueries []map[string]string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				assert.Equal(t, "/api/v1/data/events:field-values", r.URL.Path)
				query := map[string]string{}
				for k := range r.URL.Query() {
					query[k] = r.URL.Query().Get(k)
				}
				gotQueries = append(gotQueries, query)

				w.Header().Set("Content-Type", "application/json")
				assert.NoError(t, json.NewEncoder(w).Encode(map[string][]string{
					"values": valuesByField[query["field"]],
				}))
			}))
			defer server.Close()

			transport := datav1.DefaultTransportConfig().
				WithHost(server.URL[7:]).
				WithSchemes([]string{"http"})
			eventTools, err := NewTools(
				datav1.NewHTTPClientWithConfig(nil, transport),
				nil,
				zaptest.NewLogger(t),
				links.NewBuilder("https://test.chronosphere.io"),
			)
			require.NoError(t, err)

			var tool *tools.MCPTool
			for _, mt := range eventTools.MCPTools() {
				if mt.Metadata.Name == tt.tool {
					tool = &mt
				}
			}
			require.NotNil(t, tool)

			result, err := tool.Handler(t.Context(), mcp.CallToolRequest{
				Params: mcp.CallToolParams{Name: tt.tool, Arguments: tt.args},
			})
			require.NoError(t, err)

			assert.Equal(t, tt.wantQueries, gotQueries)
			gotJSON, err := json.Marshal(result.JSONContent)
			require.NoError(t, err)
			assert.JSONEq(t, tt.wantJSON, string(gotJSON))
			assert.Equal(t, tt.wantChronLink, result.ChronosphereLink)
		})
	}
}
