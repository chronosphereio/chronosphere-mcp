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

package monitors

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"github.com/chronosphereio/chronosphere-mcp/generated/datav1/datav1"
	datav1models "github.com/chronosphereio/chronosphere-mcp/generated/datav1/models"
	"github.com/chronosphereio/chronosphere-mcp/generated/statev1/models"
	"github.com/chronosphereio/chronosphere-mcp/generated/statev1/statev1"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools/pkg/params"
)

func TestParseLabelMatcher(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *models.CommonPromQLMatcher
		wantErr string
	}{
		{
			name:  "equal",
			input: "service=api",
			want:  &models.CommonPromQLMatcher{Name: "service", Type: models.CommonPromQLMatcherTypeMatchEqual, Value: "api"},
		},
		{
			name:  "not equal",
			input: "env!=prod",
			want:  &models.CommonPromQLMatcher{Name: "env", Type: models.CommonPromQLMatcherTypeMatchNotEqual, Value: "prod"},
		},
		{
			name:  "value containing equals",
			input: "query=a=b",
			want:  &models.CommonPromQLMatcher{Name: "query", Type: models.CommonPromQLMatcherTypeMatchEqual, Value: "a=b"},
		},
		{
			name:  "empty value",
			input: "service=",
			want:  &models.CommonPromQLMatcher{Name: "service", Type: models.CommonPromQLMatcherTypeMatchEqual, Value: ""},
		},
		{
			name:    "no operator",
			input:   "service",
			wantErr: "must be of the form",
		},
		{
			name:    "empty name",
			input:   "=api",
			wantErr: "empty label name",
		},
		{
			name:    "empty name with not equal",
			input:   "!=api",
			wantErr: "empty label name",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseLabelMatcher(tt.input)
			if tt.wantErr != "" {
				require.ErrorContains(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestLabelFilter(t *testing.T) {
	tests := []struct {
		name          string
		signalLabels  []string
		monitorLabels []string
		want          *models.Datav1LabelFilter
		wantErr       string
	}{
		{
			name: "no labels",
			want: nil,
		},
		{
			name:          "signal and monitor labels in one clause",
			signalLabels:  []string{"service=api"},
			monitorLabels: []string{"team!=infra"},
			want: &models.Datav1LabelFilter{
				MatchStrategy: models.LabelFilterMatchStrategySUBSETMATCH,
				LabelMatchers: []*models.LabelFilterLabelMatcher{
					{
						FilterTarget: models.LabelMatcherFilterTargetSIGNAL,
						Matcher:      &models.CommonPromQLMatcher{Name: "service", Type: models.CommonPromQLMatcherTypeMatchEqual, Value: "api"},
					},
					{
						FilterTarget: models.LabelMatcherFilterTargetMONITOR,
						Matcher:      &models.CommonPromQLMatcher{Name: "team", Type: models.CommonPromQLMatcherTypeMatchNotEqual, Value: "infra"},
					},
				},
			},
		},
		{
			name:          "invalid monitor label",
			monitorLabels: []string{"team"},
			wantErr:       "invalid monitor_labels",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := labelFilter(tt.signalLabels, tt.monitorLabels)
			if tt.wantErr != "" {
				require.ErrorContains(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseSeverity(t *testing.T) {
	tests := []struct {
		input   string
		want    models.Datav1Severity
		wantErr bool
	}{
		{input: "", want: ""},
		{input: "warn", want: models.Datav1SeverityWARN},
		{input: "critical", want: models.Datav1SeverityCRITICAL},
		{input: "info", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := parseSeverity(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestToMonitorStatuses(t *testing.T) {
	firing := func(slug string, muted bool, severity models.Datav1Severity) *models.Statev1MonitorStatus {
		return &models.Statev1MonitorStatus{MonitorSlug: slug, Muted: muted, Severity: severity}
	}
	tests := []struct {
		name   string
		firing []*models.Statev1MonitorStatus
		sortBy string
		want   []*MonitorStatus
	}{
		{
			name:   "no results",
			firing: nil,
			want:   []*MonitorStatus{},
		},
		{
			name: "firing monitors sorted by slug",
			firing: []*models.Statev1MonitorStatus{
				firing("b", false, models.Datav1SeverityWARN),
				firing("a", true, models.Datav1SeverityCRITICAL),
			},
			want: []*MonitorStatus{
				{MonitorSlug: "a", State: stateMuted, Muted: true, Severity: models.Datav1SeverityCRITICAL},
				{MonitorSlug: "b", State: stateAlerting, Severity: models.Datav1SeverityWARN},
			},
		},
		{
			name: "sort by state then severity then slug",
			firing: []*models.Statev1MonitorStatus{
				firing("muted-crit", true, models.Datav1SeverityCRITICAL),
				firing("b-warn", false, models.Datav1SeverityWARN),
				firing("z-crit", false, models.Datav1SeverityCRITICAL),
				firing("a-warn", false, models.Datav1SeverityWARN),
			},
			sortBy: sortByState,
			want: []*MonitorStatus{
				{MonitorSlug: "z-crit", State: stateAlerting, Severity: models.Datav1SeverityCRITICAL},
				{MonitorSlug: "a-warn", State: stateAlerting, Severity: models.Datav1SeverityWARN},
				{MonitorSlug: "b-warn", State: stateAlerting, Severity: models.Datav1SeverityWARN},
				{MonitorSlug: "muted-crit", State: stateMuted, Muted: true, Severity: models.Datav1SeverityCRITICAL},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, toMonitorStatuses(tt.firing, tt.sortBy))
		})
	}
}

func TestListMonitorStatuses(t *testing.T) {
	tests := []struct {
		name        string
		args        map[string]any
		status      int
		response    string
		wantBody    string
		wantResult  *ListMonitorStatusesResult
		wantErr     string
		skipRequest bool
	}{
		{
			name: "maps request and response",
			args: map[string]any{
				"monitor_slugs":           []any{"firing", "quiet"},
				"team_slugs":              []any{"team-a"},
				"signal_labels":           []any{"service=api"},
				"monitor_labels":          []any{"env!=dev"},
				"severity":                "critical",
				"include_signal_statuses": true,
				"sort_by":                 sortByState,
			},
			status: http.StatusOK,
			response: `{"monitor_statuses": [{
				"monitor_slug": "firing",
				"severity": "CRITICAL",
				"signals": [{"signal_labels": {"service": "api"}, "severities": [{"severity": "CRITICAL"}, {"severity": "WARN", "muted": true}]}]
			}]}`,
			wantBody: `{
				"monitor_slugs": ["firing", "quiet"],
				"collection_slugs": null,
				"team_slugs": ["team-a"],
				"include_signals": true,
				"severity": "CRITICAL",
				"label_filters": [{
					"match_strategy": "SUBSET_MATCH",
					"label_matchers": [
						{"filter_target": "SIGNAL", "matcher": {"name": "service", "type": "MatchEqual", "value": "api"}},
						{"filter_target": "MONITOR", "matcher": {"name": "env", "type": "MatchNotEqual", "value": "dev"}}
					]
				}]
			}`,
			wantResult: &ListMonitorStatusesResult{
				MonitorStatuses: []*MonitorStatus{
					{
						MonitorSlug: "firing",
						State:       stateAlerting,
						Severity:    models.Datav1SeverityCRITICAL,
						Signals: []*models.Statev1SignalStatus{{
							SignalLabels: map[string]string{"service": "api"},
							Severities: []*models.Statev1SeverityStatus{
								{Severity: models.Datav1SeverityCRITICAL},
								{Severity: models.Datav1SeverityWARN, Muted: true},
							},
						}},
					},
				},
			},
		},
		{
			name:     "no filters",
			args:     map[string]any{},
			status:   http.StatusOK,
			response: `{}`,
			wantBody: `{
				"monitor_slugs": null,
				"collection_slugs": null,
				"team_slugs": null,
				"label_filters": null
			}`,
			wantResult: &ListMonitorStatusesResult{MonitorStatuses: []*MonitorStatus{}},
		},
		{
			name:     "not found means not enabled",
			args:     map[string]any{},
			status:   http.StatusNotFound,
			response: `{"message": "not found"}`,
			wantErr:  "monitor status API is not enabled for this tenant",
		},
		{
			name:     "other errors are reported",
			args:     map[string]any{},
			status:   http.StatusInternalServerError,
			response: `{"message": "boom"}`,
			wantErr:  "failed to list monitor statuses",
		},
		{
			name:        "invalid label matcher",
			args:        map[string]any{"signal_labels": []any{"service"}},
			wantErr:     "invalid signal_labels",
			skipRequest: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var requested bool
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requested = true
				assert.Equal(t, http.MethodPost, r.Method)
				assert.Equal(t, "/api/v1/state/monitor-statuses", r.URL.Path)
				body, err := io.ReadAll(r.Body)
				assert.NoError(t, err)
				if tt.wantBody != "" {
					assert.JSONEq(t, tt.wantBody, string(body))
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.status)
				_, err = w.Write([]byte(tt.response))
				assert.NoError(t, err)
			}))
			defer server.Close()

			client := statev1.NewHTTPClientWithConfig(nil, statev1.DefaultTransportConfig().
				WithHost(strings.TrimPrefix(server.URL, "http://")).
				WithSchemes([]string{"http"}))
			monitorTools, err := NewTools(client, nil, zaptest.NewLogger(t))
			require.NoError(t, err)

			request := mcp.CallToolRequest{}
			request.Params.Name = "list_monitor_statuses"
			request.Params.Arguments = tt.args

			result, err := findTool(t, monitorTools, "list_monitor_statuses").Handler(t.Context(), request)
			assert.Equal(t, !tt.skipRequest, requested)
			if tt.wantErr != "" {
				require.ErrorContains(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)

			// Compare through JSON, which is what the tool returns to the caller.
			got, err := json.Marshal(result.JSONContent)
			require.NoError(t, err)
			want, err := json.Marshal(tt.wantResult)
			require.NoError(t, err)
			assert.JSONEq(t, string(want), string(got))
		})
	}
}

func findTool(t *testing.T, monitorTools *Tools, name string) tools.MCPTool {
	t.Helper()
	for _, tool := range monitorTools.MCPTools() {
		if tool.Metadata.Name == name {
			return tool
		}
	}
	require.Failf(t, "tool not found", "no tool named %q", name)
	return tools.MCPTool{}
}

func TestListAlertsRequest(t *testing.T) {
	start := time.Date(2026, 1, 2, 3, 0, 0, 0, time.UTC)
	end := start.Add(time.Hour)
	timeRange := params.TimeRange{Start: start, End: end}
	window := &datav1models.Datav1TimeRange{Start: strfmt.DateTime(start), End: strfmt.DateTime(end)}

	tests := []struct {
		name    string
		args    listAlertsArgs
		want    *datav1models.Datav1ListAlertsRequest
		wantErr string
	}{
		{
			name: "defaults to triggered window and newest first",
			args: listAlertsArgs{timeRange: timeRange},
			want: &datav1models.Datav1ListAlertsRequest{
				Filter:    &datav1models.Datav1AlertsFilter{TriggerTimeRange: window},
				SortOrder: datav1models.ListAlertsRequestSortOrderTRIGGEREDDESC,
			},
		},
		{
			name: "active time filter sets the active window",
			args: listAlertsArgs{timeRange: timeRange, timeFilter: "active", sort: "oldest"},
			want: &datav1models.Datav1ListAlertsRequest{
				Filter:    &datav1models.Datav1AlertsFilter{ActiveTimeRange: window},
				SortOrder: datav1models.ListAlertsRequestSortOrderTRIGGEREDASC,
			},
		},
		{
			name: "every filter",
			args: listAlertsArgs{
				timeRange:           timeRange,
				timeFilter:          "triggered",
				firingStatus:        "resolved",
				severity:            "critical",
				muteStatus:          "not_muted",
				sourceType:          "slo",
				monitorSlugs:        []string{"mon"},
				monitorNameContains: "api",
				sloSlugs:            []string{"slo"},
				sloNameContains:     "latency",
				alertTitleContains:  "down",
				collectionSlugs:     []string{"coll"},
				teamSlugs:           []string{"team"},
				notifierSlugs:       []string{"pager"},
				signalLabels:        []string{"service=api"},
				monitorLabels:       []string{"env!=dev"},
				alertIDs:            []string{"id-1"},
				pageSize:            10,
				pageToken:           "next",
			},
			want: &datav1models.Datav1ListAlertsRequest{
				Filter: &datav1models.Datav1AlertsFilter{
					TriggerTimeRange:      window,
					FiringStatus:          datav1models.Datav1FiringStatusRESOLVED,
					Severity:              datav1models.Datav1SeverityCRITICAL,
					MuteStatus:            datav1models.Datav1MuteStatusNOTMUTED,
					SourceType:            datav1models.Datav1SourceTypeSLO,
					MonitorSlugs:          []string{"mon"},
					MonitorNameContains:   "api",
					SloSlugs:              []string{"slo"},
					SloNameContains:       "latency",
					AlertTitleContains:    "down",
					SourceCollectionSlugs: []string{"coll"},
					SourceTeamSlugs:       []string{"team"},
					NotifierSlugs:         []string{"pager"},
					AlertIds:              []string{"id-1"},
					LabelFilters: []*datav1models.Datav1LabelFilter{{
						MatchStrategy: datav1models.LabelFilterMatchStrategySUBSETMATCH,
						LabelMatchers: []*datav1models.LabelFilterLabelMatcher{
							{
								FilterTarget: datav1models.LabelMatcherFilterTargetSIGNAL,
								Matcher:      &datav1models.CommonPromQLMatcher{Name: "service", Type: datav1models.CommonPromQLMatcherTypeMatchEqual, Value: "api"},
							},
							{
								FilterTarget: datav1models.LabelMatcherFilterTargetMONITOR,
								Matcher:      &datav1models.CommonPromQLMatcher{Name: "env", Type: datav1models.CommonPromQLMatcherTypeMatchNotEqual, Value: "dev"},
							},
						},
					}},
				},
				SortOrder: datav1models.ListAlertsRequestSortOrderTRIGGEREDDESC,
				Page:      &datav1models.Configv1PageParams{MaxSize: 10, Token: "next"},
			},
		},
		{
			name:    "invalid time filter",
			args:    listAlertsArgs{timeFilter: "resolved"},
			wantErr: `invalid time_filter "resolved"`,
		},
		{
			name:    "invalid firing status",
			args:    listAlertsArgs{firingStatus: "FIRING"},
			wantErr: `invalid firing_status "FIRING": must be one of firing, resolved`,
		},
		{
			name:    "invalid severity",
			args:    listAlertsArgs{severity: "info"},
			wantErr: `invalid severity "info"`,
		},
		{
			name:    "invalid mute status",
			args:    listAlertsArgs{muteStatus: "unmuted"},
			wantErr: `invalid mute_status "unmuted"`,
		},
		{
			name:    "invalid source type",
			args:    listAlertsArgs{sourceType: "dashboard"},
			wantErr: `invalid source_type "dashboard"`,
		},
		{
			name:    "invalid sort",
			args:    listAlertsArgs{sort: "asc"},
			wantErr: `invalid sort "asc"`,
		},
		{
			name:    "invalid label matcher",
			args:    listAlertsArgs{monitorLabels: []string{"env=~dev.*", "team"}},
			wantErr: "invalid monitor_labels",
		},
		{
			name:    "negative page size",
			args:    listAlertsArgs{pageSize: -1},
			wantErr: "invalid page_size",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := listAlertsRequest(tt.args)
			if tt.wantErr != "" {
				require.ErrorContains(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestListAlerts(t *testing.T) {
	tests := []struct {
		name        string
		args        map[string]any
		status      int
		response    string
		wantBody    string
		wantResult  string
		wantErr     string
		skipRequest bool
	}{
		{
			name: "pages and maps the response",
			args: map[string]any{
				"start":         "2026-01-02T03:00:00Z",
				"end":           "2026-01-02T04:00:00Z",
				"monitor_slugs": []any{"api-errors"},
				"page_size":     2,
				"page_token":    "page-2",
			},
			status: http.StatusOK,
			response: `{
				"alerts": [
					{"alert_id": "a1", "monitor_slug": "api-errors", "severity": "CRITICAL", "firing_status": "FIRING", "triggered_at": "2026-01-02T03:30:00Z"},
					{"alert_id": "a2", "monitor_slug": "api-errors", "severity": "WARN", "firing_status": "RESOLVED", "triggered_at": "2026-01-02T03:10:00Z", "resolved_at": "2026-01-02T03:20:00Z"}
				],
				"page": {"next_token": "page-3"},
				"total_count": "5"
			}`,
			wantBody: `{
				"filter": {
					"trigger_time_range": {"start": "2026-01-02T03:00:00.000Z", "end": "2026-01-02T04:00:00.000Z"},
					"monitor_slugs": ["api-errors"],
					"alert_ids": null,
					"external_connections": null,
					"label_filters": null,
					"notifier_slugs": null,
					"slo_slugs": null,
					"source_collection_slugs": null,
					"source_team_slugs": null
				},
				"page": {"max_size": 2, "token": "page-2"},
				"sort_order": "TRIGGERED_DESC"
			}`,
			wantResult: `{
				"total_count": 5,
				"alerts": [
					{"alert_id": "a1", "monitor_slug": "api-errors", "severity": "CRITICAL", "firing_status": "FIRING", "triggered_at": "2026-01-02T03:30:00.000Z", "external_connections": null, "notifier_slugs": null},
					{"alert_id": "a2", "monitor_slug": "api-errors", "severity": "WARN", "firing_status": "RESOLVED", "triggered_at": "2026-01-02T03:10:00.000Z", "resolved_at": "2026-01-02T03:20:00.000Z", "external_connections": null, "notifier_slugs": null}
				],
				"next_page_token": "page-3"
			}`,
		},
		{
			name:       "last page",
			args:       map[string]any{},
			status:     http.StatusOK,
			response:   `{"alerts": [], "page": {}, "total_count": "0"}`,
			wantResult: `{"total_count": 0, "alerts": []}`,
		},
		{
			name:     "not found means not enabled",
			args:     map[string]any{},
			status:   http.StatusNotFound,
			response: `{"message": "not found"}`,
			wantErr:  "alerts API is not enabled for this tenant",
		},
		{
			name:     "other errors are reported",
			args:     map[string]any{},
			status:   http.StatusInternalServerError,
			response: `{"message": "boom"}`,
			wantErr:  "failed to list alerts",
		},
		{
			name:     "malformed total count",
			args:     map[string]any{},
			status:   http.StatusOK,
			response: `{"total_count": "many"}`,
			wantErr:  "invalid total_count",
		},
		{
			name:        "invalid enum sends no request",
			args:        map[string]any{"severity": "info"},
			wantErr:     `invalid severity "info"`,
			skipRequest: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var requested bool
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requested = true
				assert.Equal(t, http.MethodPost, r.Method)
				assert.Equal(t, "/api/v1/data/alerts", r.URL.Path)
				body, err := io.ReadAll(r.Body)
				assert.NoError(t, err)
				if tt.wantBody != "" {
					assert.JSONEq(t, tt.wantBody, string(body))
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.status)
				_, err = w.Write([]byte(tt.response))
				assert.NoError(t, err)
			}))
			defer server.Close()

			client := datav1.NewHTTPClientWithConfig(nil, datav1.DefaultTransportConfig().
				WithHost(strings.TrimPrefix(server.URL, "http://")).
				WithSchemes([]string{"http"}))
			monitorTools, err := NewTools(nil, client, zaptest.NewLogger(t))
			require.NoError(t, err)

			request := mcp.CallToolRequest{}
			request.Params.Name = "list_alerts"
			request.Params.Arguments = tt.args

			result, err := findTool(t, monitorTools, "list_alerts").Handler(t.Context(), request)
			assert.Equal(t, !tt.skipRequest, requested)
			if tt.wantErr != "" {
				require.ErrorContains(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)

			got, err := json.Marshal(result.JSONContent)
			require.NoError(t, err)
			assert.JSONEq(t, tt.wantResult, string(got))
		})
	}
}
