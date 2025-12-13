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

// Package mcpserver contains the http server that serves the MCP API.
package mcpserver

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/authcontext"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/instrumentfx"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/resources"
	logresources "github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/resources/logs"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools"
	"github.com/chronosphereio/chronosphere-mcp/pkg/version"
)

type Server struct {
	server *server.MCPServer
	logger *zap.Logger
	opts   Options
}

type Options struct {
	Logger         *zap.Logger
	DisabledTools  map[string]struct{}
	ToolGroups     []tools.MCPTools
	UseLogscale    bool
	TracerProvider trace.TracerProvider
	MeterProvider  *metric.MeterProvider
}

func NewServer(
	opts Options,
	logger *zap.Logger,
) (*Server, error) {
	// Build server options
	serverOptions := []server.ServerOption{
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
		server.WithToolHandlerMiddleware(instrumentfx.ToolTracingMiddleware(opts.TracerProvider)),
		server.WithToolHandlerMiddleware(instrumentfx.ToolMetricsMiddleware(opts.MeterProvider)),
		server.WithResourceHandlerMiddleware(instrumentfx.ResourceTracingMiddleware(opts.TracerProvider)),
		server.WithResourceHandlerMiddleware(instrumentfx.ResourceMetricsMiddleware(opts.MeterProvider)),
		server.WithLogging(),
		// Filter tools based on disabled tools from request context
		server.WithToolFilter(func(ctx context.Context, tools []mcp.Tool) []mcp.Tool {
			disabledTools := authcontext.FetchDisabledTools(ctx)
			if len(disabledTools) == 0 {
				return tools
			}
			filtered := make([]mcp.Tool, 0, len(tools))
			for _, tool := range tools {
				if _, disabled := disabledTools[tool.Name]; !disabled {
					filtered = append(filtered, tool)
				}
			}
			return filtered
		}),
	}

	s := &Server{
		server: server.NewMCPServer(
			"Chronosphere MCP Server",
			version.Version,
			serverOptions...,
		),
		logger: logger,
		opts:   opts,
	}

	var (
		resources []resources.MCPResource
	)

	resources = append(resources, logresources.Resources()...)

	// Register all tools.
	for _, group := range opts.ToolGroups {
		if opts.UseLogscale && group.GroupName() == "logs" {
			// If we're using Logscale, we skip the logs group as it is already handled by the log resources.
			continue
		} else if !opts.UseLogscale && group.GroupName() == "logscale" {
			// If we're not using Logscale, we skip the logscale group as it is not needed.
			continue
		}
		for _, tool := range group.MCPTools() {
			if _, ok := opts.DisabledTools[tool.Metadata.Name]; ok {
				continue
			}

			wrapper := &loggingTool{
				logger: logger,
				tool:   tool,
			}
			s.server.AddTool(tool.MCPGoTool(), wrapper.handle)
		}
	}

	for _, resource := range resources {
		s.server.AddResource(resource.Resource, resource.Handler)
	}

	return s, nil
}

func (s *Server) StdioServer() *server.StdioServer {
	return server.NewStdioServer(s.server)
}

func (s *Server) SSEServer(baseURL string, options ...server.SSEOption) *server.SSEServer {
	return server.NewSSEServer(s.server,
		append(options,
			[]server.SSEOption{
				server.WithBaseURL(baseURL),
				server.WithSSEContextFunc(authcontext.HTTPInboundContextFunc),
			}...)...)
}

func (s *Server) StreamableHTTPServer(options ...server.StreamableHTTPOption) *server.StreamableHTTPServer {
	return server.NewStreamableHTTPServer(s.server,
		append(options,
			[]server.StreamableHTTPOption{
				server.WithHTTPContextFunc(authcontext.HTTPInboundContextFunc),
			}...)...)
}

var _ server.ToolHandlerFunc = (*loggingTool)(nil).handle

type loggingTool struct {
	logger *zap.Logger
	tool   tools.MCPTool
}

func (t *loggingTool) handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	t.logger.Info("received request for tool",
		zap.String("method", request.Method),
		zap.String("tool_name", t.tool.Metadata.Name))
	t.logger.Debug("raw request for tool",
		zap.String("method", request.Method),
		zap.String("tool_name", t.tool.Metadata.Name),
		zap.Any("request", request.Request))

	// Always wrap error responses in a proper MCP response
	resp := t.mustHandle(ctx, request)

	t.logger.Info("received response from tool",
		zap.String("method", request.Method),
		zap.String("tool_name", t.tool.Metadata.Name))
	t.logger.Debug("raw response from tool",
		zap.String("method", request.Method),
		zap.String("tool_name", t.tool.Metadata.Name),
		zap.Any("request", request),
		zap.Any("response", resp))
	return resp, nil
}

func (t *loggingTool) mustHandle(ctx context.Context, request mcp.CallToolRequest) *mcp.CallToolResult {
	resp, err := t.tool.Handler(ctx, request)
	if err != nil {
		t.logger.Info("received error from handler",
			zap.String("method", request.Method),
			zap.Error(err),
			zap.Any("request", request))
		return mcp.NewToolResultError(err.Error())
	}

	var toolResult mcp.CallToolResult
	if resp.ChronosphereLink != "" {
		toolResult.Content = append(toolResult.Content, mcp.NewTextContent("link to chronosphere: "+resp.ChronosphereLink))
	}
	if len(resp.ImageContent) > 0 {
		encoded := base64.StdEncoding.EncodeToString(resp.ImageContent)
		toolResult.Content = append(toolResult.Content, mcp.NewImageContent(encoded, "image/png"))
		return &toolResult
	}

	if resp.TextContent != "" {
		toolResult.Content = append(toolResult.Content, mcp.NewTextContent(resp.TextContent))
		if len(resp.Meta) > 0 {
			toolResult.Meta = mcp.NewMetaFromMap(resp.Meta)
		}
		return &toolResult
	}

	resultBytes, err := json.Marshal(resp.JSONContent)
	if err != nil {
		return mcp.NewToolResultError(fmt.Errorf("failed to serialize content: %s", err).Error())
	}

	toolResult.Content = append(toolResult.Content, mcp.NewTextContent(string(resultBytes)))
	if len(resp.Meta) > 0 {
		toolResult.Meta = mcp.NewMetaFromMap(resp.Meta)
	}
	return &toolResult
}
