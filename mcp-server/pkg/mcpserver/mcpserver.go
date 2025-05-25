// Package mcpserver contains the http server that serves the MCP API.
package mcpserver

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"go.uber.org/zap"

	"github.com/chronosphereio/mcp-server/mcp-server/pkg/authcontext"
	"github.com/chronosphereio/mcp-server/mcp-server/pkg/resources"
	logresources "github.com/chronosphereio/mcp-server/mcp-server/pkg/resources/logs"
	"github.com/chronosphereio/mcp-server/mcp-server/pkg/tools"
)

type Server struct {
	server *server.MCPServer
	logger *zap.Logger
	opts   Options
}

type Options struct {
	Logger        *zap.Logger
	DisabledTools map[string]struct{}
	ToolGroups    []tools.MCPTools
}

func NewServer(
	opts Options,
	logger *zap.Logger,
) (*Server, error) {
	s := &Server{
		server: server.NewMCPServer(
			"Chronosphere MCP Server",
			"0.0.1",
			server.WithResourceCapabilities(true, true),
			server.WithPromptCapabilities(true),
			server.WithLogging(),
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
				server.WithHTTPContextFunc(func(ctx context.Context, r *http.Request) context.Context {
					authValue := strings.ReplaceAll(r.Header.Get("Authorization"), "Bearer ", "")
					return authcontext.SetSessionAPIToken(ctx, authValue)
				}),
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

	resultBytes, err := json.Marshal(resp.JSONContent)
	if err != nil {
		return mcp.NewToolResultError(fmt.Errorf("failed to serialize content: %s", err).Error())
	}

	toolResult.Content = append(toolResult.Content, mcp.NewTextContent(string(resultBytes)))
	if len(resp.Meta) > 0 {
		toolResult.Meta = resp.Meta
	}
	return &toolResult
}
