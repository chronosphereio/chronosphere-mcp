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
	"net/url"
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
	mcptools "github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools/pkg/params"
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

	tools, err := NewTools(client, nil, logger, linkBuilder)
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
	resp, err := tools.dataV1API.Version1.ListEvents(queryParams)

	// We expect success since our test server returns a valid response
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Payload)
}

// fakeEventPages serves ListEvents from pages keyed by page token, recording each request's query.
func fakeEventPages(t *testing.T, pages map[string]models.Datav1ListEventsResponse) (*datav1.DataV1API, *[]url.Values) {
	t.Helper()
	var requests []url.Values
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, r.URL.Query())
		page, ok := pages[r.URL.Query().Get("page.token")]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(page))
	}))
	t.Cleanup(server.Close)

	transport := datav1.DefaultTransportConfig().
		WithHost(server.URL[len("http://"):]).
		WithSchemes([]string{"http"})
	return datav1.NewHTTPClientWithConfig(nil, transport), &requests
}

func event(category, source, typ string, labels map[string]string) *models.Datav1Event {
	return &models.Datav1Event{Category: category, Source: source, Type: typ, Labels: labels}
}

func TestCollectEventFacets(t *testing.T) {
	twoPages := map[string]models.Datav1ListEventsResponse{
		"": {
			Events: []*models.Datav1Event{
				event("deploys", "github", "commit", map[string]string{"service": "api", "env": "prod"}),
				event("alerts", "", "firing", map[string]string{"service": "web"}),
			},
			Page: &models.Configv1PageResult{NextToken: "p2"},
		},
		"p2": {
			Events: []*models.Datav1Event{
				event("deploys", "argo", "", map[string]string{"service": "api", "region": "us"}),
				nil,
			},
		},
	}

	tests := []struct {
		name             string
		pages            map[string]models.Datav1ListEventsResponse
		labelName        string
		maxEvents        int
		wantCategories   []string
		wantSources      []string
		wantTypes        []string
		wantLabelNames   []string
		wantLabelValues  []string
		wantTruncated    bool
		wantRequestCount int
	}{
		{
			name:             "multi page distinct values",
			pages:            twoPages,
			labelName:        "service",
			maxEvents:        100,
			wantCategories:   []string{"alerts", "deploys"},
			wantSources:      []string{"argo", "github"},
			wantTypes:        []string{"commit", "firing"},
			wantLabelNames:   []string{"env", "region", "service"},
			wantLabelValues:  []string{"api", "web"},
			wantRequestCount: 2,
		},
		{
			name:             "absent label has no values",
			pages:            twoPages,
			labelName:        "missing",
			maxEvents:        100,
			wantCategories:   []string{"alerts", "deploys"},
			wantSources:      []string{"argo", "github"},
			wantTypes:        []string{"commit", "firing"},
			wantLabelNames:   []string{"env", "region", "service"},
			wantLabelValues:  []string{},
			wantRequestCount: 2,
		},
		{
			name:             "truncates at max events without reading further pages",
			pages:            twoPages,
			labelName:        "service",
			maxEvents:        2,
			wantCategories:   []string{"alerts", "deploys"},
			wantSources:      []string{"github"},
			wantTypes:        []string{"commit", "firing"},
			wantLabelNames:   []string{"env", "service"},
			wantLabelValues:  []string{"api", "web"},
			wantTruncated:    true,
			wantRequestCount: 1,
		},
		{
			name:             "truncates mid page",
			pages:            twoPages,
			maxEvents:        1,
			wantCategories:   []string{"deploys"},
			wantSources:      []string{"github"},
			wantTypes:        []string{"commit"},
			wantLabelNames:   []string{"env", "service"},
			wantLabelValues:  []string{},
			wantTruncated:    true,
			wantRequestCount: 1,
		},
		{
			name:             "no events",
			pages:            map[string]models.Datav1ListEventsResponse{"": {Events: []*models.Datav1Event{}}},
			maxEvents:        100,
			wantCategories:   []string{},
			wantSources:      []string{},
			wantTypes:        []string{},
			wantLabelNames:   []string{},
			wantLabelValues:  []string{},
			wantRequestCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, requests := fakeEventPages(t, tt.pages)
			tools, err := NewTools(client, nil, zaptest.NewLogger(t), links.NewBuilder("https://test.chronosphere.io"))
			require.NoError(t, err)

			now := time.Now()
			facets, err := tools.collectEventFacets(t.Context(),
				&params.TimeRange{Start: now.Add(-time.Hour), End: now}, `type = "x"`, tt.labelName, tt.maxEvents)
			require.NoError(t, err)

			assert.Equal(t, tt.wantCategories, sortedKeys(facets.categories))
			assert.Equal(t, tt.wantSources, sortedKeys(facets.sources))
			assert.Equal(t, tt.wantTypes, sortedKeys(facets.types))
			assert.Equal(t, tt.wantLabelNames, sortedKeys(facets.labelNames))
			assert.Equal(t, tt.wantLabelValues, sortedKeys(facets.labelValues))
			assert.Equal(t, tt.wantTruncated, facets.truncated)
			require.Len(t, *requests, tt.wantRequestCount)
			for _, req := range *requests {
				assert.Equal(t, `type = "x"`, req.Get("query"))
			}
		})
	}
}

func TestCollectEventFacetsError(t *testing.T) {
	client, _ := fakeEventPages(t, map[string]models.Datav1ListEventsResponse{
		"": {Events: []*models.Datav1Event{event("deploys", "", "", nil)}, Page: &models.Configv1PageResult{NextToken: "gone"}},
	})
	tools, err := NewTools(client, nil, zaptest.NewLogger(t), links.NewBuilder("https://test.chronosphere.io"))
	require.NoError(t, err)

	now := time.Now()
	_, err = tools.collectEventFacets(t.Context(), &params.TimeRange{Start: now.Add(-time.Hour), End: now}, "", "", 100)
	require.ErrorContains(t, err, "failed to list events")
}

func findTool(t *testing.T, tools *Tools, name string) mcptools.MCPTool {
	t.Helper()
	for _, tool := range tools.MCPTools() {
		if tool.Metadata.Name == name {
			return tool
		}
	}
	t.Fatalf("tool %q not found", name)
	return mcptools.MCPTool{}
}

func TestEventMetadataTools(t *testing.T) {
	pages := map[string]models.Datav1ListEventsResponse{
		"": {
			Events: []*models.Datav1Event{
				event("deploys", "github", "commit", map[string]string{"service": "api"}),
				event("alerts", "pagerduty", "firing", map[string]string{"service": "web", "env": "prod"}),
			},
		},
	}

	tests := []struct {
		name      string
		tool      string
		args      map[string]any
		wantJSON  string
		wantQuery string
		wantStart string
		wantEnd   string
	}{
		{
			name: "get_events_metadata defaults to the last 24 hours",
			tool: "get_events_metadata",
			args: map[string]any{"end": "2026-01-02T00:00:00Z"},
			wantJSON: `{"categories":["alerts","deploys"],"sources":["github","pagerduty"],` +
				`"types":["commit","firing"],"label_names":["env","service"],"truncated":false}`,
			wantStart: "2026-01-01T00:00:00.000Z",
			wantEnd:   "2026-01-02T00:00:00.000Z",
		},
		{
			name: "list_events_label_values with query and range",
			tool: "list_events_label_values",
			args: map[string]any{
				"label_name": "service",
				"query":      `category = "deploys"`,
				"start":      "2026-01-01T10:00:00Z",
				"end":        "2026-01-01T12:00:00Z",
			},
			wantJSON:  `{"label_name":"service","values":["api","web"],"truncated":false}`,
			wantQuery: `category = "deploys"`,
			wantStart: "2026-01-01T10:00:00.000Z",
			wantEnd:   "2026-01-01T12:00:00.000Z",
		},
		{
			name: "list_events_label_values for an absent label",
			tool: "list_events_label_values",
			args: map[string]any{
				"label_name": "missing",
				"end":        "2026-01-02T00:00:00Z",
			},
			wantJSON:  `{"label_name":"missing","values":[],"truncated":false}`,
			wantStart: "2026-01-01T00:00:00.000Z",
			wantEnd:   "2026-01-02T00:00:00.000Z",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, requests := fakeEventPages(t, pages)
			tools, err := NewTools(client, nil, zaptest.NewLogger(t), links.NewBuilder("https://test.chronosphere.io"))
			require.NoError(t, err)

			result, err := findTool(t, tools, tt.tool).Handler(t.Context(), mcp.CallToolRequest{
				Params: mcp.CallToolParams{Name: tt.tool, Arguments: tt.args},
			})
			require.NoError(t, err)

			got, err := json.Marshal(result.JSONContent)
			require.NoError(t, err)
			assert.JSONEq(t, tt.wantJSON, string(got))

			require.Len(t, *requests, 1)
			req := (*requests)[0]
			assert.Equal(t, tt.wantQuery, req.Get("query"))
			assert.Equal(t, tt.wantStart, req.Get("happened_after"))
			assert.Equal(t, tt.wantEnd, req.Get("happened_before"))
		})
	}
}
