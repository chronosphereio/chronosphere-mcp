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
	"encoding/json"
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetricExplorer(t *testing.T) {
	start := time.Unix(1747887246, 0)
	end := time.Unix(1747890846, 0)
	chronosphereURL := "https://rc.chronosphere.io"

	builder := NewBuilder(chronosphereURL)

	tests := []struct {
		name         string
		builder      *MetricExplorerBuilder
		expected     string
		expectedURLs []metricQuery
	}{
		{
			name:     "basic metric explorer",
			builder:  builder.MetricExplorer(),
			expected: "https://rc.chronosphere.io/metrics/explorer-v2",
		},
		{
			name:     "with single query",
			builder:  builder.MetricExplorer().WithQuery("up"),
			expected: "https://rc.chronosphere.io/metrics/explorer-v2?queries=%5B%7B%22kind%22%3A%22DataQuery%22%2C%22spec%22%3A%7B%22plugin%22%3A%7B%22kind%22%3A%22PrometheusTimeSeriesQuery%22%2C%22spec%22%3A%7B%22query%22%3A%22up%22%7D%7D%7D%7D%5D",
			expectedURLs: []metricQuery{
				{
					Kind: "DataQuery",
					Spec: dataQuerySpec{
						Plugin: metricQuerySpecPlugin{
							Kind: "PrometheusTimeSeriesQuery",
							Spec: promTimeseriesPluginSpec{
								Query: "up",
							},
						},
					},
				},
			},
		},
		{
			name:     "with complex query",
			builder:  builder.MetricExplorer().WithQuery(`rate(http_requests_total[5m])`),
			expected: "https://rc.chronosphere.io/metrics/explorer-v2?queries=%5B%7B%22kind%22%3A%22DataQuery%22%2C%22spec%22%3A%7B%22plugin%22%3A%7B%22kind%22%3A%22PrometheusTimeSeriesQuery%22%2C%22spec%22%3A%7B%22query%22%3A%22rate%28http_requests_total%5B5m%5D%29%22%7D%7D%7D%7D%5D",
			expectedURLs: []metricQuery{
				{
					Kind: "DataQuery",
					Spec: dataQuerySpec{
						Plugin: metricQuerySpecPlugin{
							Kind: "PrometheusTimeSeriesQuery",
							Spec: promTimeseriesPluginSpec{
								Query: "rate(http_requests_total[5m])",
							},
						},
					},
				},
			},
		},
		{
			name:     "with time range",
			builder:  builder.MetricExplorer().WithTimeRange(start, end),
			expected: "https://rc.chronosphere.io/metrics/explorer-v2?end=1747890846000&start=1747887246000",
		},
		{
			name:     "with start time only",
			builder:  builder.MetricExplorer().WithStartTime(start),
			expected: "https://rc.chronosphere.io/metrics/explorer-v2?start=1747887246000",
		},
		{
			name:     "with end time only",
			builder:  builder.MetricExplorer().WithEndTime(end),
			expected: "https://rc.chronosphere.io/metrics/explorer-v2?end=1747890846000",
		},
		{
			name:     "with query and time range",
			builder:  builder.MetricExplorer().WithQuery("cpu_usage").WithTimeRange(start, end),
			expected: "https://rc.chronosphere.io/metrics/explorer-v2?end=1747890846000&queries=%5B%7B%22kind%22%3A%22DataQuery%22%2C%22spec%22%3A%7B%22plugin%22%3A%7B%22kind%22%3A%22PrometheusTimeSeriesQuery%22%2C%22spec%22%3A%7B%22query%22%3A%22cpu_usage%22%7D%7D%7D%7D%5D&start=1747887246000",
			expectedURLs: []metricQuery{
				{
					Kind: "DataQuery",
					Spec: dataQuerySpec{
						Plugin: metricQuerySpecPlugin{
							Kind: "PrometheusTimeSeriesQuery",
							Spec: promTimeseriesPluginSpec{
								Query: "cpu_usage",
							},
						},
					},
				},
			},
		},
		{
			name:     "with query and start time only",
			builder:  builder.MetricExplorer().WithQuery("memory_usage").WithStartTime(start),
			expected: "https://rc.chronosphere.io/metrics/explorer-v2?queries=%5B%7B%22kind%22%3A%22DataQuery%22%2C%22spec%22%3A%7B%22plugin%22%3A%7B%22kind%22%3A%22PrometheusTimeSeriesQuery%22%2C%22spec%22%3A%7B%22query%22%3A%22memory_usage%22%7D%7D%7D%7D%5D&start=1747887246000",
			expectedURLs: []metricQuery{
				{
					Kind: "DataQuery",
					Spec: dataQuerySpec{
						Plugin: metricQuerySpecPlugin{
							Kind: "PrometheusTimeSeriesQuery",
							Spec: promTimeseriesPluginSpec{
								Query: "memory_usage",
							},
						},
					},
				},
			},
		},
		{
			name:     "with query and end time only",
			builder:  builder.MetricExplorer().WithQuery("disk_usage").WithEndTime(end),
			expected: "https://rc.chronosphere.io/metrics/explorer-v2?end=1747890846000&queries=%5B%7B%22kind%22%3A%22DataQuery%22%2C%22spec%22%3A%7B%22plugin%22%3A%7B%22kind%22%3A%22PrometheusTimeSeriesQuery%22%2C%22spec%22%3A%7B%22query%22%3A%22disk_usage%22%7D%7D%7D%7D%5D",
			expectedURLs: []metricQuery{
				{
					Kind: "DataQuery",
					Spec: dataQuerySpec{
						Plugin: metricQuerySpecPlugin{
							Kind: "PrometheusTimeSeriesQuery",
							Spec: promTimeseriesPluginSpec{
								Query: "disk_usage",
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.builder.String()
			assert.Equal(t, tt.expected, result)

			// Verify URL parsing works
			parsedURL, err := tt.builder.URL()
			require.NoError(t, err)
			assert.Equal(t, "/metrics/explorer-v2", parsedURL.Path)
			assert.Equal(t, chronosphereURL, parsedURL.Scheme+"://"+parsedURL.Host)

			// If expectedURLs is provided, verify the queries parameter
			if tt.expectedURLs != nil {
				queriesParam := parsedURL.Query().Get("queries")
				require.NotEmpty(t, queriesParam, "queries parameter should be present")

				var actualQueries []metricQuery
				err := json.Unmarshal([]byte(queriesParam), &actualQueries)
				require.NoError(t, err, "queries parameter should be valid JSON")
				assert.Equal(t, tt.expectedURLs, actualQueries)
			}
		})
	}
}

func TestMetricExplorer_ChainedCalls(t *testing.T) {
	start := time.Unix(1600000000, 0)
	end := time.Unix(1600003600, 0)
	builder := NewBuilder("https://test.chronosphere.io")

	tests := []struct {
		name     string
		buildURL func() *MetricExplorerBuilder
		expected string
	}{
		{
			name: "chain with query first",
			buildURL: func() *MetricExplorerBuilder {
				return builder.MetricExplorer().
					WithQuery("up").
					WithTimeRange(start, end)
			},
			expected: "https://test.chronosphere.io/metrics/explorer-v2?end=1600003600000&queries=%5B%7B%22kind%22%3A%22DataQuery%22%2C%22spec%22%3A%7B%22plugin%22%3A%7B%22kind%22%3A%22PrometheusTimeSeriesQuery%22%2C%22spec%22%3A%7B%22query%22%3A%22up%22%7D%7D%7D%7D%5D&start=1600000000000",
		},
		{
			name: "chain with time range first",
			buildURL: func() *MetricExplorerBuilder {
				return builder.MetricExplorer().
					WithTimeRange(start, end).
					WithQuery("cpu_usage")
			},
			expected: "https://test.chronosphere.io/metrics/explorer-v2?end=1600003600000&queries=%5B%7B%22kind%22%3A%22DataQuery%22%2C%22spec%22%3A%7B%22plugin%22%3A%7B%22kind%22%3A%22PrometheusTimeSeriesQuery%22%2C%22spec%22%3A%7B%22query%22%3A%22cpu_usage%22%7D%7D%7D%7D%5D&start=1600000000000",
		},
		{
			name: "separate start and end times",
			buildURL: func() *MetricExplorerBuilder {
				return builder.MetricExplorer().
					WithStartTime(start).
					WithEndTime(end).
					WithQuery("memory_usage")
			},
			expected: "https://test.chronosphere.io/metrics/explorer-v2?end=1600003600000&queries=%5B%7B%22kind%22%3A%22DataQuery%22%2C%22spec%22%3A%7B%22plugin%22%3A%7B%22kind%22%3A%22PrometheusTimeSeriesQuery%22%2C%22spec%22%3A%7B%22query%22%3A%22memory_usage%22%7D%7D%7D%7D%5D&start=1600000000000",
		},
		{
			name: "overwrite query",
			buildURL: func() *MetricExplorerBuilder {
				return builder.MetricExplorer().
					WithQuery("first_query").
					WithQuery("second_query")
			},
			expected: "https://test.chronosphere.io/metrics/explorer-v2?queries=%5B%7B%22kind%22%3A%22DataQuery%22%2C%22spec%22%3A%7B%22plugin%22%3A%7B%22kind%22%3A%22PrometheusTimeSeriesQuery%22%2C%22spec%22%3A%7B%22query%22%3A%22second_query%22%7D%7D%7D%7D%5D",
		},
		{
			name: "overwrite time range",
			buildURL: func() *MetricExplorerBuilder {
				laterStart := time.Unix(1600010000, 0)
				laterEnd := time.Unix(1600013600, 0)
				return builder.MetricExplorer().
					WithTimeRange(start, end).
					WithTimeRange(laterStart, laterEnd)
			},
			expected: "https://test.chronosphere.io/metrics/explorer-v2?end=1600013600000&start=1600010000000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.buildURL().String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMetricExplorer_EdgeCases(t *testing.T) {
	builder := NewBuilder("https://edge.chronosphere.io")

	tests := []struct {
		name    string
		builder *MetricExplorerBuilder
	}{
		{
			name:    "empty query",
			builder: builder.MetricExplorer().WithQuery(""),
		},
		{
			name:    "whitespace query",
			builder: builder.MetricExplorer().WithQuery("   "),
		},
		{
			name:    "special characters in query",
			builder: builder.MetricExplorer().WithQuery(`{__name__=~".*", job!=""}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Should not panic
			result := tt.builder.String()
			assert.NotEmpty(t, result)

			// Should be parseable as URL
			parsedURL, err := tt.builder.URL()
			require.NoError(t, err)
			assert.NotNil(t, parsedURL)
		})
	}
}

func TestMetricExplorer_TimeHandling(t *testing.T) {
	builder := NewBuilder("https://time.chronosphere.io")

	tests := []struct {
		name          string
		setupBuilder  func() *MetricExplorerBuilder
		expectedStart *time.Time
		expectedEnd   *time.Time
	}{
		{
			name: "zero time values",
			setupBuilder: func() *MetricExplorerBuilder {
				return builder.MetricExplorer().
					WithStartTime(time.Time{}).
					WithEndTime(time.Time{})
			},
			expectedStart: &time.Time{},
			expectedEnd:   &time.Time{},
		},
		{
			name: "unix epoch",
			setupBuilder: func() *MetricExplorerBuilder {
				epoch := time.Unix(0, 0)
				return builder.MetricExplorer().
					WithTimeRange(epoch, epoch)
			},
			expectedStart: func() *time.Time { t := time.Unix(0, 0); return &t }(),
			expectedEnd:   func() *time.Time { t := time.Unix(0, 0); return &t }(),
		},
		{
			name: "far future time",
			setupBuilder: func() *MetricExplorerBuilder {
				future := time.Date(2099, 12, 31, 23, 59, 59, 0, time.UTC)
				return builder.MetricExplorer().
					WithStartTime(future)
			},
			expectedStart: func() *time.Time { t := time.Date(2099, 12, 31, 23, 59, 59, 0, time.UTC); return &t }(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metricBuilder := tt.setupBuilder()
			result := metricBuilder.String()
			assert.NotEmpty(t, result)

			parsedURL, err := url.Parse(result)
			require.NoError(t, err)

			if tt.expectedStart != nil {
				startParam := parsedURL.Query().Get("start")
				expectedStartMs := tt.expectedStart.UnixMilli()
				assert.Equal(t, expectedStartMs, parseTimeParam(t, startParam))
			}

			if tt.expectedEnd != nil {
				endParam := parsedURL.Query().Get("end")
				expectedEndMs := tt.expectedEnd.UnixMilli()
				assert.Equal(t, expectedEndMs, parseTimeParam(t, endParam))
			}
		})
	}
}

func parseTimeParam(t *testing.T, param string) int64 {
	if param == "" {
		return 0
	}
	// Parse the parameter as a string representation of milliseconds
	var msInt int64
	_, err := fmt.Sscanf(param, "%d", &msInt)
	require.NoError(t, err)
	return msInt
}
