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

// Package main provides the entry point for the agent application.
package main

import (
	"context"

	"go.uber.org/fx"

	"github.com/chronosphereio/chronosphere-mcp/agent/pkg/agentfx"
	"github.com/chronosphereio/chronosphere-mcp/agent/pkg/configfx"
)

func main() {
	app := fx.New(configfx.Module, agentfx.Module)
	if err := app.Start(context.Background()); err != nil {
		panic(err)
	}
	if err := app.Stop(context.Background()); err != nil {
		panic(err)
	}
}
