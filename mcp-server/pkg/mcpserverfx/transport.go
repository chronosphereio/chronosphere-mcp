package mcpserverfx

import (
	"context"
	"errors"
	"os"
	"sync"
	"time"

	"github.com/chronosphereio/mcp-server/mcp-server/pkg/mcpserver"

	"github.com/mark3labs/mcp-go/server"
	"go.uber.org/zap"
)

type TransportConfig struct {
	StdioTransport *StdioTransportConfig `yaml:"stdio"`
	SSETransport   *SSETransportConfig   `yaml:"sse"`
}

type StdioTransportConfig struct {
	Enabled bool `yaml:"enabled"`
}

type Transports struct {
	logger        *zap.Logger
	server        *mcpserver.Server
	stdio         *StdioTransportConfig
	sse           *SSETransportConfig
	sseServer     *server.SSEServer
	cancelContext context.CancelFunc
	wg            sync.WaitGroup
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
		logger:        logger,
		server:        server,
		stdio:         stdio,
		sse:           sse,
		cancelContext: func() {},
	}, nil
}

func (t *Transports) Start(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	t.cancelContext = cancel

	if t.stdio != nil && t.stdio.Enabled {
		t.wg.Add(1)
		go func() {
			defer t.wg.Done()
			t.logger.Info("serving stdio transport")
			if err := t.server.StdioServer().Listen(ctx, os.Stdin, os.Stdout); err != nil {
				t.logger.Error("error serving stdio transport", zap.Error(err))
			}
		}()
	}
	if t.sse != nil && t.sse.Enabled {
		t.wg.Add(1)
		t.sseServer = t.server.SSEServer(t.sse.BaseURL)

		go func() {
			defer t.wg.Done()
			t.logger.Info("serving sse transport",
				zap.String("address", t.sse.Address),
				zap.String("baseURL", t.sse.BaseURL))
			if err := t.sseServer.Start(t.sse.Address); err != nil {
				t.logger.Error("error serving sse transport", zap.Error(err))
			}
		}()
	}
}

func (t *Transports) Close(ctx context.Context) {
	t.cancelContext()

	if t.sse != nil && t.sse.Enabled {
		ctxShutdown, cancel := context.WithTimeout(ctx, 1*time.Second)
		defer cancel()
		if err := t.sseServer.Shutdown(ctxShutdown); err != nil {
			t.logger.Error("error shutting down sse server", zap.Error(err))
		}
	}

	t.wg.Wait()
}
