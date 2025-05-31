// Package mcpserverfx contains the fx module that provides the MCP server.
package mcpserverfx

import (
	"fmt"

	"go.uber.org/config"
	"go.uber.org/fx"
	"gopkg.in/validator.v2"

	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/clientfx"
)

const (
	apiURLFormat      = "https://%s.chronosphere.io"
	logscaleURLFormat = "https://%s.logs.chronosphere.io"
)

type configResult struct {
	fx.Out

	Config       *Config
	Chronosphere *clientfx.ChronosphereConfig
}

func parseConfig(cfgProvider config.Provider, flags *Flags) (configResult, error) {
	var cfg Config
	if err := cfgProvider.Get(configKey).Populate(&cfg); err != nil {
		return configResult{}, fmt.Errorf("failed to populate config: %w", err)
	}

	mergeConfig(&cfg, flags)

	if err := validator.Validate(cfg); err != nil {
		return configResult{}, fmt.Errorf("failed to validate config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return configResult{}, fmt.Errorf("failed to validate Chronosphere config: %w", err)
	}

	return configResult{
		Config:       &cfg,
		Chronosphere: &cfg.Chronosphere,
	}, nil
}

func mergeConfig(c *Config, f *Flags) {
	if f.UseLogScale {
		c.Chronosphere.UseLogscale = f.UseLogScale
	}

	if f.orgName != "" {
		c.Chronosphere.APIURL = fmt.Sprintf(apiURLFormat, f.orgName)
		c.Chronosphere.LogscaleURL = fmt.Sprintf(logscaleURLFormat, f.orgName)
	}
	if f.apiToken != "" {
		c.Chronosphere.APIToken = f.apiToken
	}
	if f.UseLogScale {
		c.Chronosphere.UseLogscale = f.UseLogScale
	}
}
