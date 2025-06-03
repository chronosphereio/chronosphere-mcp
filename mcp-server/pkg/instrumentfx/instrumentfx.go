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
	"go.uber.org/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Options(
	fx.Provide(NewInstrument),
)

type Input struct {
	fx.In

	ConfigProvider config.Provider
}

type Result struct {
	fx.Out

	Logger *zap.Logger
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

	return Result{
		Logger: logger,
	}, nil
}
