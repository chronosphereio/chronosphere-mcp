// Package configapi provides tools for the Chronosphere config API.
package configapi

import (
	"go.uber.org/zap"

	configv1client "github.com/chronosphereio/mcp-server/generated/configv1/configv1"
	"github.com/chronosphereio/mcp-server/mcp-server/pkg/generated/tools/configv1"
	"github.com/chronosphereio/mcp-server/mcp-server/pkg/tools"
)

var _ tools.MCPTools = (*Tools)(nil)

// Tools provides tools for the Chronosphere config API.
type Tools struct {
	logger *zap.Logger
	client *configv1client.ConfigV1API
}

func NewTools(
	client *configv1client.ConfigV1API,
	logger *zap.Logger,
) (*Tools, error) {
	logger.Info("events tool configured")

	return &Tools{
		logger: logger,
		client: client,
	}, nil
}

func (t *Tools) MCPTools() []tools.MCPTool {
	return []tools.MCPTool{
		configv1.GetMonitor(t.client, t.logger),
		configv1.ListMonitors(t.client, t.logger),
		configv1.GetDashboard(t.client, t.logger),
		configv1.ListDashboards(t.client, t.logger),
		configv1.GetSlo(t.client, t.logger),
		configv1.ListSlos(t.client, t.logger),
	}
}
