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

// Package mcpserverfx contains the fx module that provides the MCP server.
package mcpserverfx

import (
	"fmt"

	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools"
	"go.uber.org/config"
	"go.uber.org/fx"
	"gopkg.in/validator.v2"

	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/clientfx"
)

const (
	apiURLFormat = "https://%s.chronosphere.io"
)

type configResult struct {
	fx.Out

	Config       *Config
	ToolsConfig  *tools.Config
	Chronosphere *clientfx.ChronosphereConfig
}

func parseConfig(cfgProvider config.Provider, flags *Flags) (configResult, error) {
	var cfg Config
	if err := cfgProvider.Get(configKey).Populate(&cfg); err != nil {
		return configResult{}, fmt.Errorf("failed to populate config: %w", err)
	}

	mergeConfig(&cfg, flags)

	if cfg.Tools == nil {
		cfg.Tools = &tools.Config{}
	}

	if err := validator.Validate(cfg); err != nil {
		return configResult{}, fmt.Errorf("failed to validate config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return configResult{}, fmt.Errorf("failed to validate Chronosphere config: %w", err)
	}

	return configResult{
		Config:       &cfg,
		ToolsConfig:  cfg.Tools,
		Chronosphere: &cfg.Chronosphere,
	}, nil
}

func mergeConfig(c *Config, f *Flags) {
	if f.orgName != "" {
		c.Chronosphere.APIURL = fmt.Sprintf(apiURLFormat, f.orgName)
	}
	if f.apiToken != "" {
		c.Chronosphere.APIToken = f.apiToken
	}
}
