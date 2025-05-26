package mcpserverfx

import (
	"context"
	"fmt"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/chronosphereio/mcp-server/mcp-server/pkg/mcpserver"
	"github.com/chronosphereio/mcp-server/mcp-server/pkg/tools"
)

// Module registers the server.
var Module = fx.Options(
	fx.Provide(parseConfig),
	fx.Invoke(invoke),
)

const configKey = "server"

type params struct {
	fx.In
	LifeCycle fx.Lifecycle

	Config     *Config
	Logger     *zap.Logger
	ToolGroups []tools.MCPTools `group:"mcp_tools"`
}

type ToolsConfig struct {
	Disabled []string `yaml:"disabled"`
}

type Config struct {
	Transport TransportConfig `yaml:"transport"`
	Tools     *ToolsConfig    `yaml:"tools"`
}

func invoke(p params) (*Transports, error) {
	cfg := p.Config

	disabledTools := make(map[string]struct{})
	if cfg.Tools != nil {
		for _, name := range cfg.Tools.Disabled {
			disabledTools[name] = struct{}{}
		}
	}
	transports, err := NewTransports(
		mcpserver.Options{
			Logger:        p.Logger,
			ToolGroups:    p.ToolGroups,
			DisabledTools: disabledTools,
		},
		p.Logger,
		&cfg.Transport,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create transports: %v", err)
	}

	p.LifeCycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go transports.Start(context.Background())
			return nil
		},
		OnStop: func(ctx context.Context) error {
			transports.Close(ctx)
			return nil
		},
	})

	return transports, nil
}
