// Package tools provides the tools for the MCP server.
package tools

import (
	"context"
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"
)

type Result struct {
	// Exactly 1 of these *Content fields should be populated if this is not an error response.
	ImageContent []byte
	JSONContent  any

	// Optional.
	Meta map[string]any
}

type Handler func(session Session, request mcp.CallToolRequest) (*Result, error)

type Metadata struct {
	// The name of the tool.
	Name string `json:"name"`
	// A human-readable description of the tool.
	Description string `json:"description,omitempty"`
	// A JSON Schema object defining the expected parameters for the tool.
	InputSchema mcp.ToolInputSchema `json:"inputSchema"`
	// Alternative to InputSchema - allows arbitrary JSON Schema to be provided
	RawInputSchema json.RawMessage `json:"-"` // Hide this from JSON marshaling
	// Optional properties describing tool behavior
	Annotations mcp.ToolAnnotation `json:"annotations"`
}

type Session struct {
	APIToken string
	Context  context.Context
}

type MCPTool struct {
	Metadata Metadata
	Handler  Handler
}

func NewMetadata(name string, opts ...mcp.ToolOption) Metadata {
	mcpTool := mcp.NewTool(name, opts...)
	return Metadata{
		Name:           name,
		Description:    mcpTool.Description,
		InputSchema:    mcpTool.InputSchema,
		RawInputSchema: mcpTool.RawInputSchema,
		Annotations:    mcpTool.Annotations,
	}
}

func (t MCPTool) MCPGoTool() mcp.Tool {
	return mcp.Tool{
		Name:           t.Metadata.Name,
		Description:    t.Metadata.Description,
		InputSchema:    t.Metadata.InputSchema,
		RawInputSchema: t.Metadata.RawInputSchema,
		Annotations:    t.Metadata.Annotations,
	}
}

type MCPTools interface {
	MCPTools() []MCPTool
}
