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

// Package instrumentfx is a fx module that provides logs, metrics and tracing integration.
package instrumentfx

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.uber.org/config"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/chronosphereio/chronosphere-mcp/pkg/version"
)

var Module = fx.Options(
	fx.Provide(NewInstrument),
	fx.Invoke(initializeTracing),
)

type Input struct {
	fx.In

	ConfigProvider config.Provider
}

type Result struct {
	fx.Out

	Logger         *zap.Logger
	TracerProvider *trace.TracerProvider
}

type Config struct {
	ServiceName string     `yaml:"serviceName"`
	Logs        zap.Config `yaml:"logs"`
}

func defaultConfig() Config {
	return Config{
		ServiceName: "chrono-mcp",
		Logs:        zap.NewProductionConfig(),
	}
}

func NewInstrument(
	params Input,
) (Result, error) {
	cfg := defaultConfig()
	if err := params.ConfigProvider.Get("instrument").Populate(&cfg); err != nil {
		return Result{}, err
	}

	logger, err := cfg.Logs.Build()
	if err != nil {
		return Result{}, err
	}

	logger.Info("Instrument config",
		zap.String("serviceName", cfg.ServiceName),
		zap.String("logLevel", cfg.Logs.Level.String()),
		zap.Strings("output", cfg.Logs.OutputPaths),
		zap.Strings("errorOutput", cfg.Logs.ErrorOutputPaths),
	)

	traceProvider, err := newTracerProvider()
	if err != nil {
		return Result{}, fmt.Errorf("failed to create tracer provider: %w", err)
	}

	return Result{
		Logger:         logger,
		TracerProvider: traceProvider,
	}, nil
}

// newTracerProvider creates a new OpenTelemetry TracerProvider.
func newTracerProvider() (*trace.TracerProvider, error) {
	// Create resource with service information
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("chrono-mcp"),
			semconv.ServiceVersion(version.Version),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Create TracerProvider with default configuration
	tp := trace.NewTracerProvider(
		trace.WithResource(res),
		trace.WithSampler(trace.AlwaysSample()),
	)

	return tp, nil
}

// initializeTracing sets up the global tracer provider and handles lifecycle.
func initializeTracing(lc fx.Lifecycle, logger *zap.Logger, tp *trace.TracerProvider) {
	otel.SetTracerProvider(tp)
	logger.Info("OpenTelemetry tracing initialized")

	// Handle shutdown on application stop
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			logger.Info("Shutting down OpenTelemetry tracer provider")
			if err := tp.Shutdown(ctx); err != nil {
				logger.Error("Failed to shutdown OpenTelemetry tracer provider", zap.Error(err))
				return err
			}
			return nil
		},
	})
}
