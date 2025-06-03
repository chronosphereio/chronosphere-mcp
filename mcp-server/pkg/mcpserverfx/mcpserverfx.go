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

package mcpserverfx

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/clientfx"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/mcpserver"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools"
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

	Config         *Config
	Logger         *zap.Logger
	ToolGroups     []tools.MCPTools `group:"mcp_tools"`
	TracerProvider *trace.TracerProvider
	MeterProvider  *metric.MeterProvider
}

type ToolsConfig struct {
	Disabled []string `yaml:"disabled"`
}

type Config struct {
	Transport    TransportConfig             `yaml:"transport"`
	Tools        *ToolsConfig                `yaml:"tools"`
	Chronosphere clientfx.ChronosphereConfig `yaml:"chronosphere" validate:"nonnil"`
}

func (c Config) Validate() error {
	if err := c.Chronosphere.Validate(); err != nil {
		return fmt.Errorf("chronosphere config validation failed: %w", err)
	}
	if err := c.Transport.Validate(); err != nil {
		return fmt.Errorf("transport config validation failed: %w", err)
	}
	return nil
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
			Logger:         p.Logger,
			ToolGroups:     p.ToolGroups,
			DisabledTools:  disabledTools,
			UseLogscale:    cfg.Chronosphere.UseLogscale,
			TracerProvider: p.TracerProvider,
			MeterProvider:  p.MeterProvider,
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
