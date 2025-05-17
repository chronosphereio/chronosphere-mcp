// Package configapi provides tools for the Chronosphere config API.
package configapi

import (
	"go.uber.org/zap"

	"github.com/chronosphereio/mcp-server/mcp-server/pkg/client"
	"github.com/chronosphereio/mcp-server/mcp-server/pkg/generated/tools/configv1"
	"github.com/chronosphereio/mcp-server/mcp-server/pkg/tools"
)

var _ tools.MCPTools = (*Tools)(nil)

// Tools provides tools for the Chronosphere config API.
type Tools struct {
	logger         *zap.Logger
	clientProvider *client.Provider
}

func NewTools(
	clientProvider *client.Provider,
	logger *zap.Logger,
) (*Tools, error) {
	logger.Info("events tool configured")

	return &Tools{
		logger:         logger,
		clientProvider: clientProvider,
	}, nil
}

func (t *Tools) MCPTools() []tools.MCPTool {
	return []tools.MCPTool{
		configv1.GetMonitor(t.clientProvider, t.logger),
		configv1.ListMonitors(t.clientProvider, t.logger),
		configv1.GetDashboard(t.clientProvider, t.logger),
		configv1.ListDashboards(t.clientProvider, t.logger),
		configv1.GetSlo(t.clientProvider, t.logger),
		configv1.ListSlos(t.clientProvider, t.logger),
	}
}
