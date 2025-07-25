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
	"errors"
	"os"
	"sync"
	"time"

	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/mcpserver"

	"github.com/mark3labs/mcp-go/server"
	"go.uber.org/zap"
)

type TransportConfig struct {
	StdioTransport *StdioTransportConfig `yaml:"stdio"`
	SSETransport   *SSETransportConfig   `yaml:"sse"`
	HTTPTransport  *HTTPTransportConfig  `yaml:"http"`
}

func (t TransportConfig) Validate() error {
	if !t.HTTPTransport.IsEnabled() && !t.SSETransport.IsEnabled() && !t.StdioTransport.IsEnabled() {
		return errors.New("at least one transport must be enabled")
	}

	return nil
}

type SSETransportConfig struct {
	Enabled bool   `yaml:"enabled"`
	Address string `yaml:"address" validate:"nonzero"`
	BaseURL string `yaml:"baseURL" validate:"nonzero"`
}

func (s *SSETransportConfig) IsEnabled() bool {
	return s != nil && s.Enabled
}

type HTTPTransportConfig struct {
	Enabled bool   `yaml:"enabled"`
	Address string `yaml:"address" validate:"nonzero"`
}

func (h *HTTPTransportConfig) IsEnabled() bool {
	return h != nil && h.Enabled
}

type StdioTransportConfig struct {
	Enabled bool `yaml:"enabled"`
}

func (s *StdioTransportConfig) IsEnabled() bool {
	return s != nil && s.Enabled
}

type Transports struct {
	logger        *zap.Logger
	server        *mcpserver.Server
	stdio         *StdioTransportConfig
	sse           *SSETransportConfig
	http          *HTTPTransportConfig
	sseServer     *server.SSEServer
	httpServer    *server.StreamableHTTPServer
	cancelContext context.CancelFunc
	wg            sync.WaitGroup
}

func NewTransports(
	opts mcpserver.Options,
	logger *zap.Logger,
	transportConfig *TransportConfig,
) (*Transports, error) {
	if err := transportConfig.Validate(); err != nil {
		return nil, err
	}

	s, err := mcpserver.NewServer(opts, logger)
	if err != nil {
		return nil, err
	}

	return &Transports{
		logger:        logger,
		server:        s,
		stdio:         transportConfig.StdioTransport,
		sse:           transportConfig.SSETransport,
		http:          transportConfig.HTTPTransport,
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
				t.logger.DPanic("error serving stdio transport", zap.Error(err))
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
				t.logger.DPanic("error serving sse transport", zap.Error(err))
			}
		}()
	}

	if t.http.IsEnabled() {
		t.wg.Add(1)
		t.httpServer = t.server.StreamableHTTPServer()
		go func() {
			defer t.wg.Done()

			t.logger.Info("serving streamable http transport",
				zap.String("address", t.http.Address))
			if err := t.httpServer.Start(t.http.Address); err != nil {
				t.logger.Error("error serving streamable http transport", zap.Error(err))
			}
		}()
	}
}

func (t *Transports) Close(ctx context.Context) {
	t.cancelContext()

	if t.sse.IsEnabled() {
		ctxShutdown, cancel := context.WithTimeout(ctx, 1*time.Second)
		defer cancel()
		if err := t.sseServer.Shutdown(ctxShutdown); err != nil {
			t.logger.Error("error shutting down sse server", zap.Error(err))
		}
	}

	if t.http.IsEnabled() {
		ctxShutdown, cancel := context.WithTimeout(ctx, 1*time.Second)
		defer cancel()
		if err := t.httpServer.Shutdown(ctxShutdown); err != nil {
			t.logger.Error("error shutting down http server", zap.Error(err))
		}
	}

	t.wg.Wait()
}
