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

// Package toolsfx contains the fx module that provides all tools.
package toolsfx

import (
	"go.uber.org/fx"

	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools/configapi"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools/events"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools/logs"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools/metricusage"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools/monitors"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools/prometheus"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools/traces"
)

var Module = fx.Provide(
	annotateAsTool(
		configapi.NewTools,
		events.NewTools,
		logs.NewTools,
		metricusage.NewTools,
		monitors.NewTools,
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
