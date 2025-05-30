// Package logscale provides tools for querying LogScale.
package logscale

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"

	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/clientfx/logscale"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools/pkg/params"
	"github.com/chronosphereio/chronosphere-mcp/pkg/links"
)

var _ tools.MCPTools = (*Tools)(nil)

type Tools struct {
	logger      *zap.Logger
	client      logscale.Client
	linkBuilder *links.Builder
}

func NewTools(client logscale.Client, logger *zap.Logger, linkBuilder *links.Builder) (*Tools, error) {
	logger.Info("logscale tool configured")

	return &Tools{
		logger:      logger,
		client:      client,
		linkBuilder: linkBuilder,
	}, nil
}

func (t *Tools) GroupName() string {
	return "logscale"
}

func (t *Tools) MCPTools() []tools.MCPTool {
	return []tools.MCPTool{
		{
			Metadata: tools.NewMetadata("query_logscale",
				mcp.WithDescription(`Query LogScale repository with a given query string. LogScale uses its own query language for searching and analyzing logs.`),
				mcp.WithString("query",
					mcp.Description("LogScale query string to execute"),
				),
				mcp.WithString("repository",
					mcp.Description("LogScale repository name to query"),
				),
				params.WithTimeRange(),
			),
			Handler: func(ctx context.Context, request mcp.CallToolRequest) (*tools.Result, error) {
				query, err := params.String(request, "query", true, "")
				if err != nil {
					return nil, err
				}

				repository, err := params.String(request, "repository", true, "")
				if err != nil {
					return nil, err
				}

				timeRange, err := params.ParseTimeRange(request)
				if err != nil {
					return nil, err
				}

				t.logger.Info("querying logscale",
					zap.String("query", query),
					zap.String("repository", repository),
					zap.Time("start", timeRange.Start),
					zap.Time("end", timeRange.End))

				results, err := t.client.Query(ctx, query, repository, timeRange.Start, timeRange.End)
				if err != nil {
					return nil, fmt.Errorf("failed to query logscale: %w", err)
				}

				return &tools.Result{
					JSONContent: results,
				}, nil
			},
		},
	}
}
