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
	"errors"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func TestToolMetricsMiddlewareStatus(t *testing.T) {
	tests := []struct {
		name       string
		handler    func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error)
		wantStatus string
		wantPanic  bool
	}{
		{
			name: "success",
			handler: func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
				return mcp.NewToolResultText("ok"), nil
			},
			wantStatus: "success",
		},
		{
			name: "error",
			handler: func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
				return nil, errors.New("boom")
			},
			wantStatus: "error",
		},
		{
			name: "error result",
			handler: func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
				return mcp.NewToolResultError("boom"), nil
			},
			wantStatus: "error",
		},
		{
			name: "panic",
			handler: func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
				panic("boom")
			},
			wantStatus: "panic",
			wantPanic:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := sdkmetric.NewManualReader()
			handler := ToolMetricsMiddleware(sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader)))(tt.handler)
			request := mcp.CallToolRequest{Params: mcp.CallToolParams{Name: "test_tool"}}

			if tt.wantPanic {
				assert.Panics(t, func() {
					//nolint:errcheck // the handler panics before returning.
					handler(t.Context(), request)
				})
			} else {
				result, err := handler(t.Context(), request)
				if tt.wantStatus == "success" {
					require.NoError(t, err)
					assert.False(t, result.IsError)
				}
			}

			var rm metricdata.ResourceMetrics
			require.NoError(t, reader.Collect(t.Context(), &rm))
			require.Len(t, rm.ScopeMetrics, 1)
			require.Len(t, rm.ScopeMetrics[0].Metrics, 1)
			sum, ok := rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.Sum[int64])
			require.True(t, ok)
			require.Len(t, sum.DataPoints, 1)
			status, ok := sum.DataPoints[0].Attributes.Value("status")
			require.True(t, ok)
			assert.Equal(t, tt.wantStatus, status.AsString())
			assert.Equal(t, int64(1), sum.DataPoints[0].Value)
		})
	}
}
