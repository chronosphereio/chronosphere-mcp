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

// Package tools provides the tools for the MCP server.
package tools

import (
	"context"
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"
)

type Config struct {
	Disabled                []string `yaml:"disabled"`
	EnableClassicDashboards bool     `yaml:"enableClassicDashboards"`
}

type Result struct {
	// Exactly 1 of these *Content fields should be populated if this is not an error response.
	ImageContent []byte
	JSONContent  any

	// Optional.
	Meta map[string]any
	// Optional. If set, an additional text content will be returned with this link in the response.
	ChronosphereLink string
}

type Handler func(ctx context.Context, request mcp.CallToolRequest) (*Result, error)

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
	Context context.Context
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
	GroupName() string
	MCPTools() []MCPTool
}
