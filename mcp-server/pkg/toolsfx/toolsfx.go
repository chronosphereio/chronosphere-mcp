package toolsfx

import (
	"go.uber.org/fx"

	"github.com/chronosphereio/mcp-server/mcp-server/pkg/tools"
	"github.com/chronosphereio/mcp-server/mcp-server/pkg/tools/configapi"
	"github.com/chronosphereio/mcp-server/mcp-server/pkg/tools/events"
	"github.com/chronosphereio/mcp-server/mcp-server/pkg/tools/logs"
	"github.com/chronosphereio/mcp-server/mcp-server/pkg/tools/prometheus"
	"github.com/chronosphereio/mcp-server/mcp-server/pkg/tools/traces"
)

var Module = fx.Provide(
	annotateAsTool(
		configapi.NewTools,
		events.NewTools,
		logs.NewTools,
		prometheus.NewTools,
		traces.NewTools,
	)...,
)

func annotateAsTool(constructors ...any) []any {
	annotated := make([]any, len(constructors))
	for i, constructor := range constructors {
		annotated[i] = fx.Annotate(
			constructor,
			// This is a value group.
			// See https://uber-go.github.io/fx/value-groups/feed.html
			fx.ResultTags(`group:"mcp_tools"`),
			// Provide hint that this is an interface so it can be grouped into the right tool group.
			// See https://uber-go.github.io/fx/faq.html#why-does-fxsupply-not-accept-interfaces
			fx.As(new(tools.MCPTools)),
		)
	}
	return annotated
}
