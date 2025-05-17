// Package mcpserverfx contains the fx module that provides the MCP server.
package mcpserverfx

import (
	"fmt"

	"go.uber.org/config"
	"gopkg.in/validator.v2"
)

func parseConfig(cfgProvider config.Provider) (*Config, error) {
	var cfg Config
	if err := cfgProvider.Get(configKey).Populate(&cfg); err != nil {
		return nil, fmt.Errorf("failed to populate config: %v", err)
	}
	if err := validator.Validate(cfg); err != nil {
		return nil, fmt.Errorf("failed to validate config: %v", err)
	}
	if err := validator.Validate(cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
