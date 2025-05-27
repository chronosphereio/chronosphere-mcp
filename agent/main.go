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
