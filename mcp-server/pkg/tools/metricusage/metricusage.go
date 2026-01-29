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

// Package metricusage provides tools for querying metric usage data from the State V1 API.
// These tools help identify unused or underutilized metrics that could be candidates for dropping.
package metricusage

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"

	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools/pkg/params"
	"github.com/chronosphereio/chronosphere-mcp/pkg/links"
)

var _ tools.MCPTools = (*Tools)(nil)

// StateV1Client is an HTTP client for the State V1 API.
type StateV1Client struct {
	httpClient *http.Client
	baseURL    string
}

// Tools provides MCP tools for metric usage analysis.
type Tools struct {
	logger      *zap.Logger
	client      *StateV1Client
	linkBuilder *links.Builder
}

// NewStateV1Client creates a new State V1 API client.
func NewStateV1Client(httpClient *http.Client, baseURL string) *StateV1Client {
	return &StateV1Client{
		httpClient: httpClient,
		baseURL:    baseURL,
	}
}

// NewTools creates new metric usage tools.
func NewTools(
	stateV1Client *StateV1Client,
	logger *zap.Logger,
	linkBuilder *links.Builder,
) (*Tools, error) {
	logger.Info("metric usage tool configured")
	return &Tools{
		logger:      logger,
		client:      stateV1Client,
		linkBuilder: linkBuilder,
	}, nil
}

func (t *Tools) GroupName() string {
	return "metric_usage"
}

func (t *Tools) MCPTools() []tools.MCPTool {
	return []tools.MCPTool{
		{
			Metadata: tools.NewMetadata("list_metric_usages_by_metric_name",
				mcp.WithDescription(`Lists metric usage statistics grouped by metric name. Use this to find unused or underutilized metrics that could be dropped to reduce costs.

Response fields:
- metric_name: The name of the metric
- dpps: Data points per second (ingestion rate)
- usage.total_references: Number of dashboards, monitors, recording rules, drop rules, and aggregation rules referencing this metric
- usage.total_query_executions: Number of times this metric was queried (in dashboards, explorer, or external queries)
- usage.total_unique_users: Number of unique users who queried this metric
- usage.utility_score: A computed score indicating how valuable this metric is (higher = more valuable)
- usage.reference_counts_by_type: Breakdown of references by type (dashboards, monitors, recording_rules, drop_rules, aggregation_rules)
- usage.query_execution_counts_by_type: Breakdown of query executions by source (explorer, dashboard, external)

Metrics with low utility_score, zero references, and zero query executions are good candidates for dropping.`),
				mcp.WithNumber("page_max_size",
					mcp.Description("Maximum number of results per page. Default is server-defined.")),
				mcp.WithString("page_token",
					mcp.Description("Pagination token for fetching the next page.")),
				mcp.WithString("order_by",
					mcp.Description("Field to order results by. One of: VALUABLE, DPPS, UTILITY, REFERENCES, EXECUTIONS, UNIQUE_VALUES, UNIQUE_USERS")),
				mcp.WithBoolean("order_ascending",
					mcp.Description("If true, sort in ascending order. Default is false (descending).")),
				mcp.WithNumber("lookback_secs",
					mcp.Description("Time range in seconds for query execution counts. Default is 2592000 (30 days).")),
				mcp.WithString("metric_name_glob",
					mcp.Description("Glob pattern to filter metrics by name (e.g., 'http_*' or '*_total').")),
				mcp.WithBoolean("include_counts_by_type",
					mcp.Description("If true, include detailed breakdowns in reference_counts_by_type and query_execution_counts_by_type.")),
			),
			Handler: t.listMetricUsagesByMetricName,
		},
		{
			Metadata: tools.NewMetadata("list_metric_usages_by_label_name",
				mcp.WithDescription(`Lists metric usage statistics grouped by label name. Use this to find unused or high-cardinality labels that could be dropped.

Response fields:
- label_name: The name of the label
- dpps: Data points per second (contribution to ingestion rate)
- total_unique_values: Number of unique values for this label (indicates cardinality)
- percent_of_series_with_label_name: Percentage of all series that have this label (0-100)
- usage.total_references: Number of dashboards, monitors, recording rules, drop rules, and aggregation rules using this label
- usage.total_query_executions: Number of times this label was used in queries
- usage.total_unique_users: Number of unique users who queried with this label
- usage.utility_score: A computed score indicating how valuable this label is
- usage.reference_counts_by_type: Breakdown of references by type
- usage.query_execution_counts_by_type: Breakdown of query executions by source

Labels with low utility_score, zero references, and high unique_values are good candidates for dropping (high cardinality, low value).`),
				mcp.WithNumber("page_max_size",
					mcp.Description("Maximum number of results per page. Default is server-defined.")),
				mcp.WithString("page_token",
					mcp.Description("Pagination token for fetching the next page.")),
				mcp.WithString("order_by",
					mcp.Description("Field to order results by. One of: VALUABLE, DPPS, UTILITY, REFERENCES, EXECUTIONS, UNIQUE_VALUES, UNIQUE_USERS")),
				mcp.WithBoolean("order_ascending",
					mcp.Description("If true, sort in ascending order. Default is false (descending).")),
				mcp.WithNumber("lookback_secs",
					mcp.Description("Time range in seconds for query execution counts. Default is 2592000 (30 days).")),
				mcp.WithString("label_name_glob",
					mcp.Description("Glob pattern to filter labels by name (e.g., 'kubernetes_*').")),
				mcp.WithBoolean("include_counts_by_type",
					mcp.Description("If true, include detailed breakdowns in reference_counts_by_type and query_execution_counts_by_type.")),
			),
			Handler: t.listMetricUsagesByLabelName,
		},
		{
			Metadata: tools.NewMetadata("list_rule_evaluations",
				mcp.WithDescription(`Lists rule evaluation issues for monitors and recording rules. Use this to identify rules that are failing or having problems.

Response fields:
- rule_slug: Unique identifier of the rule
- rule_type: Type of rule - either MONITOR or RECORDING
- detected_at: Timestamp when the issue was detected (issues are aggregated over 5 minute windows)
- count: Number of evaluation issues in the last 5 minutes
- message: Detailed error message about what went wrong

Common issues include query timeouts, invalid PromQL, or missing metrics.`),
				mcp.WithNumber("page_max_size",
					mcp.Description("Maximum number of results per page. Default is server-defined.")),
				mcp.WithString("page_token",
					mcp.Description("Pagination token for fetching the next page.")),
			),
			Handler: t.listRuleEvaluations,
		},
	}
}

func (t *Tools) listMetricUsagesByMetricName(ctx context.Context, request mcp.CallToolRequest) (*tools.Result, error) {
	query := url.Values{}

	if v, err := params.Int(request, "page_max_size", false, 0); err != nil {
		return nil, err
	} else if v > 0 {
		query.Set("page.max_size", strconv.Itoa(v))
	}

	if v, err := params.String(request, "page_token", false, ""); err != nil {
		return nil, err
	} else if v != "" {
		query.Set("page.token", v)
	}

	if v, err := params.String(request, "order_by", false, ""); err != nil {
		return nil, err
	} else if v != "" {
		query.Set("order.by", v)
	}

	if v, err := params.Bool(request, "order_ascending", false, false); err != nil {
		return nil, err
	} else if v {
		query.Set("order.ascending", "true")
	}

	if v, err := params.Int(request, "lookback_secs", false, 0); err != nil {
		return nil, err
	} else if v > 0 {
		query.Set("lookback_secs", strconv.Itoa(v))
	}

	if v, err := params.String(request, "metric_name_glob", false, ""); err != nil {
		return nil, err
	} else if v != "" {
		query.Set("metric_name_glob", v)
	}

	if v, err := params.Bool(request, "include_counts_by_type", false, false); err != nil {
		return nil, err
	} else if v {
		query.Set("include_counts_by_type", "true")
	}

	resp, err := t.client.get(ctx, "/api/v1/state/metric-usages-by-metric-name", query)
	if err != nil {
		return nil, fmt.Errorf("failed to list metric usages by metric name: %w", err)
	}

	return &tools.Result{
		JSONContent: resp,
		ChronosphereLink: t.linkBuilder.Custom("/api/v1/state/metric-usages-by-metric-name").
			WithParams(query).
			String(),
	}, nil
}

func (t *Tools) listMetricUsagesByLabelName(ctx context.Context, request mcp.CallToolRequest) (*tools.Result, error) {
	query := url.Values{}

	if v, err := params.Int(request, "page_max_size", false, 0); err != nil {
		return nil, err
	} else if v > 0 {
		query.Set("page.max_size", strconv.Itoa(v))
	}

	if v, err := params.String(request, "page_token", false, ""); err != nil {
		return nil, err
	} else if v != "" {
		query.Set("page.token", v)
	}

	if v, err := params.String(request, "order_by", false, ""); err != nil {
		return nil, err
	} else if v != "" {
		query.Set("order.by", v)
	}

	if v, err := params.Bool(request, "order_ascending", false, false); err != nil {
		return nil, err
	} else if v {
		query.Set("order.ascending", "true")
	}

	if v, err := params.Int(request, "lookback_secs", false, 0); err != nil {
		return nil, err
	} else if v > 0 {
		query.Set("lookback_secs", strconv.Itoa(v))
	}

	if v, err := params.String(request, "label_name_glob", false, ""); err != nil {
		return nil, err
	} else if v != "" {
		query.Set("label_name_glob", v)
	}

	if v, err := params.Bool(request, "include_counts_by_type", false, false); err != nil {
		return nil, err
	} else if v {
		query.Set("include_counts_by_type", "true")
	}

	resp, err := t.client.get(ctx, "/api/v1/state/metric-usages-by-label-name", query)
	if err != nil {
		return nil, fmt.Errorf("failed to list metric usages by label name: %w", err)
	}

	return &tools.Result{
		JSONContent: resp,
		ChronosphereLink: t.linkBuilder.Custom("/api/v1/state/metric-usages-by-label-name").
			WithParams(query).
			String(),
	}, nil
}

func (t *Tools) listRuleEvaluations(ctx context.Context, request mcp.CallToolRequest) (*tools.Result, error) {
	query := url.Values{}

	if v, err := params.Int(request, "page_max_size", false, 0); err != nil {
		return nil, err
	} else if v > 0 {
		query.Set("page.max_size", strconv.Itoa(v))
	}

	if v, err := params.String(request, "page_token", false, ""); err != nil {
		return nil, err
	} else if v != "" {
		query.Set("page.token", v)
	}

	resp, err := t.client.get(ctx, "/api/v1/state/rule-evaluations", query)
	if err != nil {
		return nil, fmt.Errorf("failed to list rule evaluations: %w", err)
	}

	return &tools.Result{
		JSONContent: resp,
		ChronosphereLink: t.linkBuilder.Custom("/api/v1/state/rule-evaluations").
			WithParams(query).
			String(),
	}, nil
}

// get makes a GET request to the State V1 API.
func (c *StateV1Client) get(ctx context.Context, path string, query url.Values) (any, error) {
	reqURL := c.baseURL + path
	if len(query) > 0 {
		reqURL += "?" + query.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var result any
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return result, nil
}
