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

package configapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/chronosphereio/chronosphere-mcp/generated/configv1/configv1"
	"github.com/chronosphereio/chronosphere-mcp/generated/configv1/models"
)

func TestUpdateDashboard(t *testing.T) {
	var capturedBody models.ConfigV1UpdateDashboardBody
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		assert.Equal(t, "/api/v1/config/dashboards/test-dashboard", r.URL.Path)
		require.NoError(t, json.NewDecoder(r.Body).Decode(&capturedBody))

		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{
			"dashboard": {
				"slug": "test-dashboard",
				"name": "Updated dashboard",
				"dashboard_json": "{\"widgets\":[]}"
			}
		}`))
		require.NoError(t, err)
	}))
	defer server.Close()

	tool := UpdateDashboard(newConfigV1TestClient(server))
	result, err := tool.Handler(t.Context(), mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "update_dashboard",
			Arguments: map[string]any{
				"slug": "test-dashboard",
				"dashboard": map[string]any{
					"name":            "Updated dashboard",
					"dashboard_json":  `{"widgets":[]}`,
					"collection_slug": "team-dashboards",
					"labels": map[string]any{
						"team": "platform",
					},
				},
				"dry_run":           true,
				"create_if_missing": true,
			},
		},
	})

	require.NoError(t, err)
	require.NotNil(t, result)
	response, ok := result.JSONContent.(*models.Configv1UpdateDashboardResponse)
	require.True(t, ok)
	require.NotNil(t, response.Dashboard)
	assert.Equal(t, "test-dashboard", response.Dashboard.Slug)

	require.NotNil(t, capturedBody.Dashboard)
	assert.Equal(t, "test-dashboard", capturedBody.Dashboard.Slug)
	assert.Equal(t, "Updated dashboard", capturedBody.Dashboard.Name)
	assert.Equal(t, `{"widgets":[]}`, capturedBody.Dashboard.DashboardJSON)
	assert.Equal(t, "team-dashboards", capturedBody.Dashboard.CollectionSlug)
	assert.Equal(t, map[string]string{"team": "platform"}, capturedBody.Dashboard.Labels)
	assert.True(t, capturedBody.DryRun)
	assert.True(t, capturedBody.CreateIfMissing)
}

func TestUpdateDashboardValidation(t *testing.T) {
	tool := UpdateDashboard(configv1.NewHTTPClient(nil))
	tests := []struct {
		name      string
		arguments map[string]any
		errorText string
	}{
		{
			name: "missing dashboard",
			arguments: map[string]any{
				"slug": "test-dashboard",
			},
			errorText: "missing required parameter: dashboard",
		},
		{
			name: "mismatched slug",
			arguments: map[string]any{
				"slug": "test-dashboard",
				"dashboard": map[string]any{
					"slug":           "other-dashboard",
					"name":           "Updated dashboard",
					"dashboard_json": "{}",
				},
			},
			errorText: `dashboard slug "other-dashboard" must match path slug "test-dashboard"`,
		},
		{
			name: "missing name",
			arguments: map[string]any{
				"slug": "test-dashboard",
				"dashboard": map[string]any{
					"dashboard_json": "{}",
				},
			},
			errorText: "dashboard.name cannot be empty",
		},
		{
			name: "missing dashboard JSON",
			arguments: map[string]any{
				"slug": "test-dashboard",
				"dashboard": map[string]any{
					"name": "Updated dashboard",
				},
			},
			errorText: "dashboard.dashboard_json cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tool.Handler(t.Context(), mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name:      "update_dashboard",
					Arguments: tt.arguments,
				},
			})
			require.ErrorContains(t, err, tt.errorText)
		})
	}
}

func TestUpdateDashboardAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(`{"code":400,"message":"invalid dashboard"}`))
		require.NoError(t, err)
	}))
	defer server.Close()

	tool := UpdateDashboard(newConfigV1TestClient(server))
	_, err := tool.Handler(t.Context(), mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "update_dashboard",
			Arguments: map[string]any{
				"slug": "test-dashboard",
				"dashboard": map[string]any{
					"name":           "Updated dashboard",
					"dashboard_json": "{}",
				},
			},
		},
	})

	require.ErrorContains(t, err, "failed to call UpdateDashboard")
	require.ErrorContains(t, err, "invalid dashboard")
}

func TestUpdateDashboardMetadata(t *testing.T) {
	tool := UpdateDashboard(configv1.NewHTTPClient(nil))

	assert.True(t, tool.Write)
	require.NotNil(t, tool.Metadata.Annotations.ReadOnlyHint)
	assert.False(t, *tool.Metadata.Annotations.ReadOnlyHint)
	require.NotNil(t, tool.Metadata.Annotations.DestructiveHint)
	assert.True(t, *tool.Metadata.Annotations.DestructiveHint)
	require.NotNil(t, tool.Metadata.Annotations.IdempotentHint)
	assert.True(t, *tool.Metadata.Annotations.IdempotentHint)
	require.NotNil(t, tool.Metadata.Annotations.OpenWorldHint)
	assert.True(t, *tool.Metadata.Annotations.OpenWorldHint)
}

func newConfigV1TestClient(server *httptest.Server) *configv1.ConfigV1API {
	transport := configv1.DefaultTransportConfig().
		WithHost(server.URL[7:]).
		WithSchemes([]string{"http"})
	return configv1.NewHTTPClientWithConfig(nil, transport)
}
