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

// Package logs provides logs related resources for the MCP server.
package logs

import (
	"context"
	_ "embed"

	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/resources"
	"github.com/mark3labs/mcp-go/mcp"
)

var (
	//go:embed data/log-query-syntax.mdx
	_logQuerySyntaxMDX string
)

// Resources returns a list of MCP resources related to logs.
func Resources() []resources.MCPResource {
	querySyntaxURL := "file://chronosphere/docs/logs/syntax.md"
	return []resources.MCPResource{
		{
			Resource: mcp.NewResource(querySyntaxURL, "Log Query Syntax",
				mcp.WithResourceDescription("Documentation for the log query syntax in mdx format. Use this to help construct or interpret log queries."),
				mcp.WithMIMEType("text/markdown"),
			),
			Handler: func(_ context.Context, _ mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
				return []mcp.ResourceContents{
					mcp.TextResourceContents{
						Text: _logQuerySyntaxMDX,
					},
				}, nil
			},
		},
	}
}
