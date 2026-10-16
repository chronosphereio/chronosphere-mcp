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

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"github.com/chronosphereio/chronosphere-mcp/generated/statev1/models"
	"github.com/chronosphereio/chronosphere-mcp/generated/statev1/statev1"
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
		name           string
		firing         []*models.Statev1MonitorStatus
		requestedSlugs []string
		sortBy         string
		want           []*MonitorStatus
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
			name:           "requested slugs that are not firing pass",
			firing:         []*models.Statev1MonitorStatus{firing("b", false, models.Datav1SeverityWARN)},
			requestedSlugs: []string{"c", "b", "a", "c"},
			want: []*MonitorStatus{
				{MonitorSlug: "a", State: statePass},
				{MonitorSlug: "b", State: stateAlerting, Severity: models.Datav1SeverityWARN},
				{MonitorSlug: "c", State: statePass},
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
			requestedSlugs: []string{"aaa-pass"},
			sortBy:         sortByState,
			want: []*MonitorStatus{
				{MonitorSlug: "z-crit", State: stateAlerting, Severity: models.Datav1SeverityCRITICAL},
				{MonitorSlug: "a-warn", State: stateAlerting, Severity: models.Datav1SeverityWARN},
				{MonitorSlug: "b-warn", State: stateAlerting, Severity: models.Datav1SeverityWARN},
				{MonitorSlug: "muted-crit", State: stateMuted, Muted: true, Severity: models.Datav1SeverityCRITICAL},
				{MonitorSlug: "aaa-pass", State: statePass},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, toMonitorStatuses(tt.firing, tt.requestedSlugs, tt.sortBy))
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
					{MonitorSlug: "quiet", State: statePass},
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
			monitorTools, err := NewTools(client, zaptest.NewLogger(t))
			require.NoError(t, err)

			mcpTools := monitorTools.MCPTools()
			require.Len(t, mcpTools, 1)
			request := mcp.CallToolRequest{}
			request.Params.Name = "list_monitor_statuses"
			request.Params.Arguments = tt.args

			result, err := mcpTools[0].Handler(t.Context(), request)
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
