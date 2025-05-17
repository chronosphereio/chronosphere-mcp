// Package main provides the entry point for the agent application.
package main

import (
	"context"

	"go.uber.org/fx"

	"github.com/chronosphereio/mcp-server/agent/pkg/agentfx"
	"github.com/chronosphereio/mcp-server/agent/pkg/configfx"
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
