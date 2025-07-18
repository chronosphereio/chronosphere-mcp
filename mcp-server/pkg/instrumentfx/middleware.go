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

package instrumentfx

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

const (
	instrumentationName = "github.com/chronosphereio/chronosphere-mcp/mcp-server/tools"
)

func ToolTracingMiddleware(tp trace.TracerProvider) server.ToolHandlerMiddleware {
	tracer := tp.Tracer(instrumentationName)
	return func(next server.ToolHandlerFunc) server.ToolHandlerFunc {
		return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			// Start a new span for the tool call
			ctx, span := tracer.Start(ctx, "ToolCall",
				trace.WithAttributes(
					attribute.String("tool.method", request.Method),
				),
				trace.WithSpanKind(trace.SpanKindClient),
			)
			defer span.End()

			result, err := next(ctx, request)
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, "tool call resulted in error")
			} else {
				span.SetStatus(codes.Ok, "tool call completed")
			}
			return result, err
		}
	}
}

func ToolMetricsMiddleware(mp metric.MeterProvider) server.ToolHandlerMiddleware {
	meter := mp.Meter(instrumentationName)

	// Create the toolcall_total counter
	toolCallCounter, err := meter.Int64Counter(
		"toolcall_total",
		metric.WithDescription("Total number of tool calls"),
		metric.WithUnit("1"),
	)
	if err != nil {
		// If counter creation fails, return a no-op middleware
		return func(next server.ToolHandlerFunc) server.ToolHandlerFunc {
			return next
		}
	}

	return func(next server.ToolHandlerFunc) server.ToolHandlerFunc {
		return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			result, err := next(ctx, request)

			// Determine status based on error
			status := "success"
			if err != nil {
				status = "error"
			}

			// Record the metric with labels
			toolCallCounter.Add(ctx, 1,
				metric.WithAttributes(
					attribute.String("service", "chrono-mcp"),
					attribute.String("method", request.Method),
					attribute.String("status", status),
				),
			)

			return result, err
		}
	}
}
