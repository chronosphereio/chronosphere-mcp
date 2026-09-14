// Copyright 2026 Chronosphere Inc.
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

package mcpserver

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mark3labs/mcp-go/client"
	clienttransport "github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"

	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools"
)

func TestCatalogToolsMatchMCPDiscovery(t *testing.T) {
	testServer := newCatalogTestServer(t)
	httpServer := httptest.NewServer(testServer.HTTPHandler())
	t.Cleanup(httpServer.Close)

	tests := []struct {
		name          string
		headers       map[string]string
		expectedTools []string
	}{
		{
			name: "write tool hidden without request header",
			headers: map[string]string{
				"Authorization": "Bearer catalog-secret",
			},
			expectedTools: []string{"read_tool"},
		},
		{
			name: "disabled tool and write preference applied",
			headers: map[string]string{
				"Authorization":              "Bearer catalog-secret",
				"X-Chrono-MCP-Disable-Tools": "read_tool",
				"X-Chrono-MCP-Enable-Writes": "true",
			},
			expectedTools: []string{"write_tool"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initializeResult, discoveredTools := discoverMCP(t, httpServer.URL+"/mcp", tt.headers)
			catalog, body := fetchCatalog(t, httpServer.URL+catalogPath, http.MethodGet, tt.headers)

			require.Equal(t, discoveredTools, catalog.Tools)
			assert.Equal(t, initializeResult.ServerInfo.Name, catalog.Server.Name)
			assert.Equal(t, initializeResult.ServerInfo.Version, catalog.Server.Version)
			assert.Equal(t, initializeResult.Instructions, catalog.Server.Instructions)
			assert.Equal(t, initializeResult.Capabilities, catalog.Server.Capabilities)
			require.Len(t, catalog.Tools, len(tt.expectedTools))
			for i, expectedName := range tt.expectedTools {
				assert.Equal(t, expectedName, catalog.Tools[i].Name)
			}
			assert.NotContains(t, string(body), "catalog-secret")
		})
	}
}

func TestCatalogRoute(t *testing.T) {
	testServer := newCatalogTestServer(t)
	httpServer := httptest.NewServer(testServer.HTTPHandler())
	t.Cleanup(httpServer.Close)
	url := httpServer.URL + catalogPath

	headers := map[string]string{
		"Authorization": "Bearer route-secret",
		"Cookie":        "chrono-accesstoken=cookie-secret",
	}
	firstCatalog, firstBody := fetchCatalog(t, url, http.MethodGet, headers)
	secondCatalog, secondBody := fetchCatalog(t, url, http.MethodGet, headers)

	assert.Equal(t, firstBody, secondBody)
	assert.Equal(t, firstCatalog, secondCatalog)
	assert.Equal(t, catalogSchemaVersion, firstCatalog.SchemaVersion)
	assert.Empty(t, firstCatalog.Server.Instructions)
	assert.Len(t, firstCatalog.Resources, 1)
	assert.Empty(t, firstCatalog.ResourceTemplates)
	assert.Empty(t, firstCatalog.Prompts)
	require.Len(t, firstCatalog.Tools, 1)
	assert.NotEmpty(t, firstCatalog.Tools[0].InputSchema.Properties)
	assert.NotEmpty(t, firstCatalog.Tools[0].OutputSchema.Properties)
	assert.Contains(t, string(firstBody), `"instructions":""`)
	assert.Contains(t, string(firstBody), `"resourceTemplates":[]`)
	assert.Contains(t, string(firstBody), `"prompts":[]`)
	assert.NotContains(t, string(firstBody), "route-secret")
	assert.NotContains(t, string(firstBody), "cookie-secret")

	request, err := http.NewRequestWithContext(t.Context(), http.MethodHead, url, nil)
	require.NoError(t, err)
	response, err := http.DefaultClient.Do(request)
	require.NoError(t, err)
	headBody, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	require.NoError(t, response.Body.Close())
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "application/json", response.Header.Get("Content-Type"))
	assert.Equal(t, "private, no-store", response.Header.Get("Cache-Control"))
	assert.Empty(t, headBody)

	for _, method := range []string{
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodOptions,
	} {
		t.Run(method, func(t *testing.T) {
			request, err := http.NewRequestWithContext(t.Context(), method, url, nil)
			require.NoError(t, err)
			response, err := http.DefaultClient.Do(request)
			require.NoError(t, err)
			require.NoError(t, response.Body.Close())
			assert.Equal(t, http.StatusMethodNotAllowed, response.StatusCode)
			assert.Equal(t, "GET, HEAD", response.Header.Get("Allow"))
			assert.Equal(t, "application/json", response.Header.Get("Content-Type"))
			assert.Equal(t, "private, no-store", response.Header.Get("Cache-Control"))
		})
	}
}

func newCatalogTestServer(t *testing.T) *Server {
	t.Helper()

	testTools := []tools.MCPTool{
		{
			Metadata: tools.NewMetadata(
				"read_tool",
				mcp.WithDescription("Reads a value."),
				mcp.WithString("query", mcp.Required()),
				mcp.WithOutputSchema[struct {
					Result string `json:"result"`
				}](),
			),
			Handler: func(_ context.Context, _ mcp.CallToolRequest) (*tools.Result, error) {
				return &tools.Result{TextContent: "read"}, nil
			},
		},
		{
			Metadata: tools.NewMetadata("write_tool", mcp.WithDescription("Writes a value.")),
			Handler: func(_ context.Context, _ mcp.CallToolRequest) (*tools.Result, error) {
				return &tools.Result{TextContent: "written"}, nil
			},
			Write: true,
		},
		{
			Metadata: tools.NewMetadata("disabled_tool", mcp.WithDescription("Never registered.")),
			Handler: func(_ context.Context, _ mcp.CallToolRequest) (*tools.Result, error) {
				return &tools.Result{TextContent: "disabled"}, nil
			},
		},
	}
	logger := zap.NewNop()
	s, err := NewServer(Options{
		Logger:         logger,
		EnableWrites:   true,
		DisabledTools:  map[string]struct{}{"disabled_tool": {}},
		ToolGroups:     []tools.MCPTools{&mockToolGroup{tools: testTools}},
		TracerProvider: trace.NewTracerProvider(),
		MeterProvider:  metric.NewMeterProvider(),
	}, logger)
	require.NoError(t, err)
	return s
}

func discoverMCP(
	t *testing.T,
	url string,
	headers map[string]string,
) (*mcp.InitializeResult, []mcp.Tool) {
	t.Helper()

	mcpClient, err := client.NewStreamableHttpClient(url, clienttransport.WithHTTPHeaders(headers))
	require.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, mcpClient.Close())
	})
	require.NoError(t, mcpClient.Start(t.Context()))
	initializeResult, err := mcpClient.Initialize(t.Context(), mcp.InitializeRequest{
		Params: mcp.InitializeParams{
			ProtocolVersion: mcp.LATEST_PROTOCOL_VERSION,
			ClientInfo: mcp.Implementation{
				Name:    "catalog-test",
				Version: "1",
			},
		},
	})
	require.NoError(t, err)
	result, err := mcpClient.ListTools(t.Context(), mcp.ListToolsRequest{})
	require.NoError(t, err)
	return initializeResult, result.Tools
}

func fetchCatalog(
	t *testing.T,
	url string,
	method string,
	headers map[string]string,
) (Catalog, []byte) {
	t.Helper()

	request, err := http.NewRequestWithContext(t.Context(), method, url, nil)
	require.NoError(t, err)
	for name, value := range headers {
		request.Header.Set(name, value)
	}
	response, err := http.DefaultClient.Do(request)
	require.NoError(t, err)
	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	require.NoError(t, response.Body.Close())
	require.Equal(t, http.StatusOK, response.StatusCode)
	require.Equal(t, "application/json", response.Header.Get("Content-Type"))
	require.Equal(t, "private, no-store", response.Header.Get("Cache-Control"))

	var catalog Catalog
	require.NoError(t, json.Unmarshal(body, &catalog))
	return catalog, body
}
