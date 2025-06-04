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
	"strconv"
	"time"
)

// MetricExplorerBuilder builds links to the Chronosphere log explorer
type MetricExplorerBuilder struct {
	chronosphereURL string
	queries         []string
	start           *time.Time
	end             *time.Time
}

// MetricExplorer creates a new MetricExplorerBuilder for the specified chronosphereURL
func (b *Builder) MetricExplorer() *MetricExplorerBuilder {
	return &MetricExplorerBuilder{
		chronosphereURL: b.chronosphereURL,
	}
}

// WithQuery sets the promQL query
func (b *MetricExplorerBuilder) WithQuery(query string) *MetricExplorerBuilder {
	b.queries = []string{query}
	return b
}

// WithTimeRange sets the start and end time for the query
func (b *MetricExplorerBuilder) WithTimeRange(start, end time.Time) *MetricExplorerBuilder {
	b.start = &start
	b.end = &end
	return b
}

// WithStartTime sets only the start time
func (b *MetricExplorerBuilder) WithStartTime(start time.Time) *MetricExplorerBuilder {
	b.start = &start
	return b
}

// WithEndTime sets only the end time
func (b *MetricExplorerBuilder) WithEndTime(end time.Time) *MetricExplorerBuilder {
	b.end = &end
	return b
}

type metricQuery struct {
	Kind string        `json:"kind"`
	Spec dataQuerySpec `json:"spec"`
}

type dataQuerySpec struct {
	Plugin metricQuerySpecPlugin `json:"plugin"`
}

type metricQuerySpecPlugin struct {
	Kind string                   `json:"kind"`
	Spec promTimeseriesPluginSpec `json:"spec"`
}

type promTimeseriesPluginSpec struct {
	Query string `json:"query"`
}

// String builds and returns the complete URL
func (b *MetricExplorerBuilder) String() string {
	baseURL := fmt.Sprintf("%s%s", b.chronosphereURL, "/metrics/explorer-v2")
	params := url.Values{}

	querySpecs := make([]metricQuery, len(b.queries))
	for i, query := range b.queries {
		querySpecs[i] = metricQuery{
			Kind: "DataQuery",
			Spec: dataQuerySpec{
				Plugin: metricQuerySpecPlugin{
					Kind: "PrometheusTimeSeriesQuery",
					Spec: promTimeseriesPluginSpec{
						Query: query,
					},
				},
			},
		}
	}
	if len(b.queries) > 0 {
		queryJSON, _ := json.Marshal(querySpecs) // nolint:errcheck
		params.Set("queries", string(queryJSON))
	}

	if b.start != nil {
		params.Set("start", strconv.FormatInt(b.start.UnixMilli(), 10))
	}

	if b.end != nil {
		params.Set("end", strconv.FormatInt(b.end.UnixMilli(), 10))
	}

	if len(params) > 0 {
		return fmt.Sprintf("%s?%s", baseURL, params.Encode())
	}

	return baseURL
}

// URL returns the built URL as a *url.URL
func (b *MetricExplorerBuilder) URL() (*url.URL, error) {
	return url.Parse(b.String())
}
