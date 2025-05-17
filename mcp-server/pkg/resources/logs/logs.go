// Package logs provides logs related resources for the MCP server.
package logs

import (
	"context"
	_ "embed"

	"github.com/chronosphereio/mcp-server/mcp-server/pkg/resources"
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
