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

// Package configapi provides tools for the Chronosphere config API.
package configapi

import (
	"go.uber.org/zap"

	configv1client "github.com/chronosphereio/chronosphere-mcp/generated/configv1/configv1"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/generated/tools/configv1"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools"
)

var _ tools.MCPTools = (*Tools)(nil)

// Tools provides tools for the Chronosphere config API.
type Tools struct {
	logger *zap.Logger
	client *configv1client.ConfigV1API
	config *tools.Config
}

func NewTools(
	client *configv1client.ConfigV1API,
	logger *zap.Logger,
	config *tools.Config,
) (*Tools, error) {
	logger.Info("events tool configured")

	return &Tools{
		logger: logger,
		client: client,
		config: config,
	}, nil
}

func (t *Tools) GroupName() string {
	return "configapi"
}

func (t *Tools) MCPTools() []tools.MCPTool {
	mcpTools := []tools.MCPTool{
		configv1.GetDashboard(t.client, t.logger),
		configv1.ListDashboards(t.client, t.logger),
		configv1.GetDropRule(t.client, t.logger),
		configv1.ListDropRules(t.client, t.logger),
		configv1.GetMappingRule(t.client, t.logger),
		configv1.ListMappingRules(t.client, t.logger),
		configv1.GetMonitor(t.client, t.logger),
		configv1.ListMonitors(t.client, t.logger),
		configv1.GetNotificationPolicy(t.client, t.logger),
		configv1.ListNotificationPolicies(t.client, t.logger),
		configv1.GetRecordingRule(t.client, t.logger),
		configv1.ListRecordingRules(t.client, t.logger),
		configv1.GetRollupRule(t.client, t.logger),
		configv1.ListRollupRules(t.client, t.logger),
		configv1.GetSlo(t.client, t.logger),
		configv1.ListSlos(t.client, t.logger),
	}
	if t.config.EnableClassicDashboards {
		mcpTools = append(mcpTools,
			configv1.GetClassicDashboard(t.client, t.logger),
			configv1.ListClassicDashboards(t.client, t.logger),
		)
	}
	return mcpTools
}
