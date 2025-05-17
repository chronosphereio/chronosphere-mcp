// Package resources provides the MCP resources for the server.
package resources

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type MCPResource struct {
	Resource mcp.Resource
	Handler  server.ResourceHandlerFunc
}

type MCPResources interface {
	MCPResources() []mcp.Resource
}
