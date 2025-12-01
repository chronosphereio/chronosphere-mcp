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

package mcpserver

import (
	"context"
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/authcontext"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools"
)

func TestLoggingTool_mustHandle(t *testing.T) {
	tests := []struct {
		name            string
		sessionAPIToken string
		tool            tools.MCPTool
		expectedContent []mcp.Content
		expectedMeta    map[string]any
	}{
		{
			name:            "successful JSON response",
			sessionAPIToken: "test-token",
			tool: tools.MCPTool{
				Handler: func(_ context.Context, _ mcp.CallToolRequest) (*tools.Result, error) {
					return &tools.Result{
						JSONContent: map[string]string{"key": "value"},
					}, nil
				},
			},
			expectedContent: []mcp.Content{
				mcp.NewTextContent(`{"key":"value"}`),
			},
		},
		{
			name:            "handler error",
			sessionAPIToken: "test-token",
			tool: tools.MCPTool{
				Handler: func(_ context.Context, _ mcp.CallToolRequest) (*tools.Result, error) {
					return &tools.Result{}, fmt.Errorf("handler error")
				},
			},
			expectedContent: mcp.NewToolResultError("handler error").Content,
		},
		{
			name:            "response with Chronosphere link",
			sessionAPIToken: "test-token",
			tool: tools.MCPTool{
				Handler: func(_ context.Context, _ mcp.CallToolRequest) (*tools.Result, error) {
					return &tools.Result{
						ChronosphereLink: "https://chronosphere.io",
						JSONContent:      map[string]string{"data": "test"},
					}, nil
				},
			},
			expectedContent: []mcp.Content{
				mcp.NewTextContent("link to chronosphere: https://chronosphere.io"),
				mcp.NewTextContent(`{"data":"test"}`),
			},
		},
		{
			name:            "response with image content",
			sessionAPIToken: "test-token",
			tool: tools.MCPTool{
				Handler: func(_ context.Context, _ mcp.CallToolRequest) (*tools.Result, error) {
					return &tools.Result{
						ImageContent: []byte("test-image-data"),
					}, nil
				},
			},
			expectedContent: []mcp.Content{
				mcp.NewImageContent(base64.StdEncoding.EncodeToString([]byte("test-image-data")), "image/png"),
			},
		},
		{
			name:            "response with metadata",
			sessionAPIToken: "test-token",
			tool: tools.MCPTool{
				Handler: func(_ context.Context, _ mcp.CallToolRequest) (*tools.Result, error) {
					return &tools.Result{
						JSONContent: map[string]string{"key": "value"},
						Meta:        map[string]any{"meta-key": "meta-value"},
					}, nil
				},
			},
			expectedContent: []mcp.Content{
				mcp.NewTextContent(`{"key":"value"}`),
			},
			expectedMeta: map[string]any{"meta-key": "meta-value"},
		},
		{
			name:            "response with text content",
			sessionAPIToken: "test-token",
			tool: tools.MCPTool{
				Handler: func(_ context.Context, _ mcp.CallToolRequest) (*tools.Result, error) {
					return &tools.Result{
						TextContent: "CSV formatted data\nrow1,col1,col2\nrow2,val1,val2",
					}, nil
				},
			},
			expectedContent: []mcp.Content{
				mcp.NewTextContent("CSV formatted data\nrow1,col1,col2\nrow2,val1,val2"),
			},
		},
		{
			name:            "response with text content and metadata",
			sessionAPIToken: "test-token",
			tool: tools.MCPTool{
				Handler: func(_ context.Context, _ mcp.CallToolRequest) (*tools.Result, error) {
					return &tools.Result{
						TextContent: "CSV data",
						Meta:        map[string]any{"total_series": 10},
					}, nil
				},
			},
			expectedContent: []mcp.Content{
				mcp.NewTextContent("CSV data"),
			},
			expectedMeta: map[string]any{"total_series": 10},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := zap.NewNop()
			lt := &loggingTool{
				logger: logger,
				tool:   tt.tool,
			}

			ctx := authcontext.SetSessionCredentials(t.Context(), authcontext.SessionCredentials{
				APIToken: tt.sessionAPIToken,
			})
			result := lt.mustHandle(ctx, mcp.CallToolRequest{})

			assert.Equal(t, len(tt.expectedContent), len(result.Content), "content length mismatch")

			for i, content := range result.Content {
				assert.Equal(t, tt.expectedContent[i], content, "content mismatch at index %d", i)
			}

			assert.Equal(t, tt.expectedMeta, result.Meta, "metadata mismatch")
		})
	}
}
