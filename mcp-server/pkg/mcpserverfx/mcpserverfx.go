package mcpserverfx

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/mark3labs/mcp-go/server"
	"go.uber.org/config"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/chronosphereio/mcp-server/mcp-server/pkg/client"
	"github.com/chronosphereio/mcp-server/mcp-server/pkg/mcpserver"
	"github.com/chronosphereio/mcp-server/mcp-server/pkg/tools"
	"github.com/chronosphereio/mcp-server/mcp-server/pkg/tools/prometheus"
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
	ClientProvider *client.Provider
	ConfigProvider config.Provider
	Logger         *zap.Logger
	ToolGroups     []tools.MCPTools `group:"mcp_tools"`
}

type ToolsConfig struct {
	Disabled []string `yaml:"disabled"`
}

type Config struct {
	Prometheus *PrometheusConfig `yaml:"prometheus"`
	Transport  TransportConfig   `yaml:"transport"`
	Tools      *ToolsConfig      `yaml:"tools"`
}

type PrometheusConfig struct {
	UseLoopback bool `yaml:"useLoopback"`
}

func (c PrometheusConfig) Options() (*prometheus.Options, error) {
	return &prometheus.Options{
		UseLoopback: c.UseLoopback,
	}, nil
}

type TransportConfig struct {
	StdioTransport *StdioTransportConfig `yaml:"stdio"`
	SSETransport   *SSETransportConfig   `yaml:"sse"`
}

type StdioTransportConfig struct {
	Enabled bool `yaml:"enabled"`
}

type SSETransportConfig struct {
	Enabled bool   `yaml:"enabled"`
	Address string `yaml:"address" validate:"nonzero"`
	BaseURL string `yaml:"baseURL" validate:"nonzero"`
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
		cfg.Transport.StdioTransport,
		cfg.Transport.SSETransport)
	if err != nil {
		return nil, fmt.Errorf("failed to create transports: %v", err)
	}

	p.LifeCycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return transports.Start(ctx)
		},
		OnStop: func(ctx context.Context) error {
			transports.Close(ctx)
			return nil
		},
	})

	return transports, nil
}

type Transports struct {
	logger    *zap.Logger
	server    *mcpserver.Server
	stdio     *StdioTransportConfig
	stdioDone chan struct{}
	sse       *SSETransportConfig
	sseServer *server.SSEServer
}

func NewTransports(
	opts mcpserver.Options,
	logger *zap.Logger,
	stdio *StdioTransportConfig,
	sse *SSETransportConfig,
) (*Transports, error) {
	if (stdio == nil || !stdio.Enabled) && (sse == nil || !sse.Enabled) {
		return nil, errors.New("at least one transport must be enabled")
	}

	server, err := mcpserver.NewServer(opts, logger)
	if err != nil {
		return nil, err
	}

	return &Transports{
		logger: logger,
		server: server,
		stdio:  stdio,
		sse:    sse,
	}, nil
}

func (t *Transports) Start(ctx context.Context) error {
	if t.stdio != nil && t.stdio.Enabled {
		t.stdioDone = make(chan struct{})
		go func() {
			t.logger.Info("serving stdio transport")
			if err := t.server.StdioServer().Listen(ctx, os.Stdin, os.Stdout); err != nil {
				t.logger.Error("error serving stdio transport", zap.Error(err))
			}
			t.stdioDone <- struct{}{}
		}()
	}
	if t.sse != nil && t.sse.Enabled {
		t.sseServer = t.server.SSEServer(t.sse.Address, t.sse.BaseURL)
		go func() {
			t.logger.Info("serving sse transport",
				zap.String("address", t.sse.Address),
				zap.String("baseURL", t.sse.BaseURL))
			if err := t.sseServer.Start(t.sse.Address); err != nil {
				t.logger.Error("error serving sse transport", zap.Error(err))
			}
		}()
	}
	return nil
}

func (t *Transports) Close(ctx context.Context) {
	var wg sync.WaitGroup

	if t.stdio != nil && t.stdio.Enabled {
		wg.Add(1)
		go func() {
			<-t.stdioDone
			wg.Done()
		}()
	}

	if t.sse != nil && t.sse.Enabled {
		wg.Add(1)
		go func() {
			ctxShutdown, cancel := context.WithTimeout(ctx, 1*time.Second)
			defer cancel()
			if err := t.sseServer.Shutdown(ctxShutdown); err != nil {
				t.logger.Error("error shutting down sse server", zap.Error(err))
			}
			wg.Done()
		}()
	}

	wg.Wait()
}
