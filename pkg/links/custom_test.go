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

package links

import (
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCustomBuilder(t *testing.T) {
	chronosphereURL := "https://custom.chronosphere.io"
	builder := NewBuilder(chronosphereURL)

	tests := []struct {
		name     string
		builder  *CustomBuilder
		expected string
	}{
		{
			name:     "basic custom path",
			builder:  builder.Custom("/api/v1/test"),
			expected: "https://custom.chronosphere.io/api/v1/test?",
		},
		{
			name:     "with single parameter",
			builder:  builder.Custom("/api/data").WithParam("key", "value"),
			expected: "https://custom.chronosphere.io/api/data?key=value",
		},
		{
			name: "with multiple parameters",
			builder: builder.Custom("/search").
				WithParam("query", "test").
				WithParam("limit", "10"),
			expected: "https://custom.chronosphere.io/search?limit=10&query=test",
		},
		{
			name: "with special characters in parameters",
			builder: builder.Custom("/logs").
				WithParam("filter", `service="gateway" AND level="ERROR"`).
				WithParam("format", "json"),
			expected: "https://custom.chronosphere.io/logs?filter=service%3D%22gateway%22+AND+level%3D%22ERROR%22&format=json",
		},
		{
			name:     "with empty parameter value",
			builder:  builder.Custom("/endpoint").WithParam("empty", ""),
			expected: "https://custom.chronosphere.io/endpoint?",
		},
		{
			name: "with duplicate parameter keys",
			builder: builder.Custom("/multi").
				WithParam("tag", "value1").
				WithParam("tag", "value2"),
			expected: "https://custom.chronosphere.io/multi?tag=value1&tag=value2",
		},
		{
			name:     "with root path",
			builder:  builder.Custom("/"),
			expected: "https://custom.chronosphere.io/?",
		},
		{
			name:     "with empty path",
			builder:  builder.Custom(""),
			expected: "https://custom.chronosphere.io?",
		},
		{
			name: "with unicode characters",
			builder: builder.Custom("/unicode").
				WithParam("name", "测试").
				WithParam("emoji", "🚀"),
			expected: "https://custom.chronosphere.io/unicode?emoji=%F0%9F%9A%80&name=%E6%B5%8B%E8%AF%95",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.builder.String()
			assert.Equal(t, tt.expected, result)

			// Verify URL is parseable
			parsedURL, err := url.Parse(result)
			require.NoError(t, err)
			assert.Equal(t, chronosphereURL, parsedURL.Scheme+"://"+parsedURL.Host)
		})
	}
}

func TestCustomBuilder_WithTimeSec(t *testing.T) {
	builder := NewBuilder("https://time.chronosphere.io")

	tests := []struct {
		name     string
		builder  *CustomBuilder
		expected string
	}{
		{
			name: "with unix epoch",
			builder: builder.Custom("/time").
				WithTimeSec("timestamp", time.Unix(0, 0)),
			expected: "https://time.chronosphere.io/time?timestamp=0",
		},
		{
			name: "with specific timestamp",
			builder: builder.Custom("/events").
				WithTimeSec("start", time.Unix(1600000000, 0)),
			expected: "https://time.chronosphere.io/events?start=1600000000",
		},
		{
			name: "with millisecond precision",
			builder: builder.Custom("/metrics").
				WithTimeSec("time", time.Unix(1600000000, 500000000)), // 500ms
			expected: "https://time.chronosphere.io/metrics?time=1600000000",
		},
		{
			name: "with start and end times",
			builder: builder.Custom("/range").
				WithTimeSec("start", time.Unix(1600000000, 0)).
				WithTimeSec("end", time.Unix(1600003600, 0)),
			expected: "https://time.chronosphere.io/range?end=1600003600&start=1600000000",
		},
		{
			name: "with time and regular parameters",
			builder: builder.Custom("/query").
				WithParam("q", "cpu_usage").
				WithTimeSec("timestamp", time.Unix(1700000000, 0)).
				WithParam("format", "json"),
			expected: "https://time.chronosphere.io/query?format=json&q=cpu_usage&timestamp=1700000000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.builder.String()
			assert.Equal(t, tt.expected, result)

			// Verify the time parameters are correctly formatted
			parsedURL, err := url.Parse(result)
			require.NoError(t, err)

			// Check that time parameters are integers (seconds)
			for key, values := range parsedURL.Query() {
				if key == "start" || key == "end" || key == "timestamp" || key == "time" {
					for _, value := range values {
						// Parse as seconds (not milliseconds like metrics explorer)
						var seconds int64
						_, err := fmt.Sscanf(value, "%d", &seconds)
						require.NoError(t, err, "Time parameter %s=%s should be valid seconds", key, value)
						assert.GreaterOrEqual(t, seconds, int64(0), "Time parameter %s=%s should be non-negative", key, value)
					}
				}
			}
		})
	}
}

func TestCustomBuilder_ChainedCalls(t *testing.T) {
	builder := NewBuilder("https://chain.chronosphere.io")

	tests := []struct {
		name     string
		buildURL func() *CustomBuilder
		expected string
	}{
		{
			name: "parameters then time",
			buildURL: func() *CustomBuilder {
				return builder.Custom("/api").
					WithParam("service", "gateway").
					WithTimeSec("start", time.Unix(1500000000, 0)).
					WithParam("limit", "100")
			},
			expected: "https://chain.chronosphere.io/api?limit=100&service=gateway&start=1500000000",
		},
		{
			name: "time then parameters",
			buildURL: func() *CustomBuilder {
				return builder.Custom("/metrics").
					WithTimeSec("timestamp", time.Unix(1500000000, 0)).
					WithParam("metric", "cpu").
					WithParam("format", "prometheus")
			},
			expected: "https://chain.chronosphere.io/metrics?format=prometheus&metric=cpu&timestamp=1500000000",
		},
		{
			name: "multiple time parameters",
			buildURL: func() *CustomBuilder {
				start := time.Unix(1500000000, 0)
				end := time.Unix(1500003600, 0)
				checkpoint := time.Unix(1500001800, 0)
				return builder.Custom("/analysis").
					WithTimeSec("start", start).
					WithTimeSec("end", end).
					WithTimeSec("checkpoint", checkpoint)
			},
			expected: "https://chain.chronosphere.io/analysis?checkpoint=1500001800&end=1500003600&start=1500000000",
		},
		{
			name: "override parameter values",
			buildURL: func() *CustomBuilder {
				return builder.Custom("/config").
					WithParam("env", "dev").
					WithParam("env", "prod")
			},
			expected: "https://chain.chronosphere.io/config?env=dev&env=prod",
		},
		{
			name: "complex mixed chain",
			buildURL: func() *CustomBuilder {
				return builder.Custom("/complex").
					WithParam("query", "error").
					WithTimeSec("start", time.Unix(1600000000, 0)).
					WithParam("service", "api").
					WithParam("service", "gateway").
					WithTimeSec("end", time.Unix(1600003600, 0)).
					WithParam("limit", "50")
			},
			expected: "https://chain.chronosphere.io/complex?end=1600003600&limit=50&query=error&service=api&service=gateway&start=1600000000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.buildURL().String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCustomBuilder_EdgeCases(t *testing.T) {
	builder := NewBuilder("https://edge.chronosphere.io")

	tests := []struct {
		name    string
		builder *CustomBuilder
	}{
		{
			name:    "path with query params",
			builder: builder.Custom("/api?existing=param"),
		},
		{
			name:    "path with fragment",
			builder: builder.Custom("/page#section"),
		},
		{
			name:    "path with spaces",
			builder: builder.Custom("/path with spaces"),
		},
		{
			name: "parameter with spaces and special chars",
			builder: builder.Custom("/search").
				WithParam("query", "hello world & more!"),
		},
		{
			name: "empty parameter key",
			builder: builder.Custom("/test").
				WithParam("", "value"),
		},
		{
			name: "parameter with newlines",
			builder: builder.Custom("/multiline").
				WithParam("text", "line1\nline2\r\nline3"),
		},
		{
			name: "very long parameter value",
			builder: builder.Custom("/long").
				WithParam("data", string(make([]byte, 1000))),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Should not panic
			result := tt.builder.String()
			assert.NotEmpty(t, result)

			// Should be parseable as URL (though may not be valid)
			parsedURL, err := url.Parse(result)
			require.NoError(t, err)
			assert.NotNil(t, parsedURL)
		})
	}
}

func TestCustomBuilder_TimeHandling(t *testing.T) {
	builder := NewBuilder("https://timetest.chronosphere.io")

	tests := []struct {
		name          string
		setupBuilder  func() *CustomBuilder
		expectedParam string
		expectedValue int64
	}{
		{
			name: "zero time",
			setupBuilder: func() *CustomBuilder {
				return builder.Custom("/zero").
					WithTimeSec("time", time.Time{})
			},
			expectedParam: "time",
			expectedValue: time.Time{}.Unix(),
		},
		{
			name: "unix epoch",
			setupBuilder: func() *CustomBuilder {
				return builder.Custom("/epoch").
					WithTimeSec("timestamp", time.Unix(0, 0))
			},
			expectedParam: "timestamp",
			expectedValue: 0,
		},
		{
			name: "far future",
			setupBuilder: func() *CustomBuilder {
				future := time.Date(2099, 12, 31, 23, 59, 59, 0, time.UTC)
				return builder.Custom("/future").
					WithTimeSec("when", future)
			},
			expectedParam: "when",
			expectedValue: time.Date(2099, 12, 31, 23, 59, 59, 0, time.UTC).Unix(),
		},
		{
			name: "millisecond truncation",
			setupBuilder: func() *CustomBuilder {
				// Time with 750ms - WithTimeSec converts via UnixMilli()/1000
				timeWithMs := time.Unix(1600000000, 750000000)
				return builder.Custom("/truncate").
					WithTimeSec("truncated", timeWithMs)
			},
			expectedParam: "truncated",
			expectedValue: time.Unix(1600000000, 750000000).UnixMilli() / 1000, // Uses same logic as WithTimeSec
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			customBuilder := tt.setupBuilder()
			result := customBuilder.String()
			assert.NotEmpty(t, result)

			parsedURL, err := url.Parse(result)
			require.NoError(t, err)

			actualValue := parsedURL.Query().Get(tt.expectedParam)
			require.NotEmpty(t, actualValue, "Expected parameter %s should be present", tt.expectedParam)

			// Parse the timestamp directly as seconds (not milliseconds like metrics explorer)
			var actualUnix int64
			_, err = fmt.Sscanf(actualValue, "%d", &actualUnix)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedValue, actualUnix)
		})
	}
}

func TestCustomBuilder_ParameterTypes(t *testing.T) {
	builder := NewBuilder("https://types.chronosphere.io")

	tests := []struct {
		name     string
		builder  *CustomBuilder
		expected map[string][]string
	}{
		{
			name: "string parameters",
			builder: builder.Custom("/strings").
				WithParam("name", "test").
				WithParam("value", "123").
				WithParam("special", "a&b=c"),
			expected: map[string][]string{
				"name":    {"test"},
				"value":   {"123"},
				"special": {"a&b=c"},
			},
		},
		{
			name: "empty and whitespace",
			builder: builder.Custom("/empty").
				WithParam("empty", "").
				WithParam("space", " ").
				WithParam("tabs", "\t\t"),
			expected: map[string][]string{
				"empty": nil,
				"space": {" "},
				"tabs":  {"\t\t"},
			},
		},
		{
			name: "numeric strings",
			builder: builder.Custom("/numbers").
				WithParam("int", "42").
				WithParam("float", "3.14").
				WithParam("negative", "-100"),
			expected: map[string][]string{
				"int":      {"42"},
				"float":    {"3.14"},
				"negative": {"-100"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.builder.String()
			parsedURL, err := url.Parse(result)
			require.NoError(t, err)

			actualParams := parsedURL.Query()
			for key, expectedValues := range tt.expected {
				actualValues := actualParams[key]
				assert.Equal(t, expectedValues, actualValues, "Parameter %s values mismatch", key)
			}
		})
	}
}
