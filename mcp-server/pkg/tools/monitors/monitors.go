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

// Package monitors provides tools for querying Chronosphere monitor statuses.
package monitors

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"

	"github.com/chronosphereio/chronosphere-mcp/generated/statev1/models"
	"github.com/chronosphereio/chronosphere-mcp/generated/statev1/statev1"
	"github.com/chronosphereio/chronosphere-mcp/generated/statev1/statev1/monitor_statuses"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools/pkg/params"
)

var _ tools.MCPTools = (*Tools)(nil)

const sortByState = "SORT_BY_STATE"

// Monitor states reported by list_monitor_statuses.
const (
	stateAlerting = "ALERTING"
	stateMuted    = "MUTED"
)

// stateOrder orders states for SORT_BY_STATE: alerting first, then muted.
var stateOrder = map[string]int{
	stateAlerting: 0,
	stateMuted:    1,
}

// severityOrder orders severities for SORT_BY_STATE: critical first, then warn, then none.
var severityOrder = map[models.Datav1Severity]int{
	models.Datav1SeverityCRITICAL: 0,
	models.Datav1SeverityWARN:     1,
}

// Tools represents the monitor tools.
type Tools struct {
	logger *zap.Logger
	api    *statev1.StateV1API
}

// NewTools creates a new set of monitor tools.
func NewTools(api *statev1.StateV1API, logger *zap.Logger) (*Tools, error) {
	logger.Info("monitor status tool configured")

	return &Tools{
		logger: logger,
		api:    api,
	}, nil
}

func (t *Tools) GroupName() string {
	return "monitors"
}

// MonitorStatus is the status of one monitor as returned by list_monitor_statuses.
type MonitorStatus struct {
	MonitorSlug string `json:"monitor_slug"`
	// State is ALERTING, or MUTED when every firing alert is muted.
	State    string                        `json:"state"`
	Muted    bool                          `json:"muted,omitempty"`
	Severity models.Datav1Severity         `json:"severity,omitempty"`
	Signals  []*models.Statev1SignalStatus `json:"signals,omitempty"`
}

// ListMonitorStatusesResult is the result of list_monitor_statuses.
type ListMonitorStatusesResult struct {
	MonitorStatuses []*MonitorStatus `json:"monitor_statuses"`
}

// MCPTools returns the list of monitor tools.
func (t *Tools) MCPTools() []tools.MCPTool {
	return []tools.MCPTool{
		{
			Metadata: tools.NewMetadata("list_monitor_statuses",
				mcp.WithDescription("Lists the current status of monitors in Chronosphere. Returns every monitor matching the filters with a firing alert. Monitors that are not firing are omitted from the results. Each status carries its state (ALERTING or MUTED), its most severe alert severity, and optionally its firing signals."),
				params.WithStringArray("monitor_slugs",
					mcp.Description("Filter by monitor slug. If all filters are empty, return status for all firing monitors."),
				),
				params.WithStringArray("collection_slugs",
					mcp.Description("Filter monitor state by owning collection."),
				),
				params.WithStringArray("team_slugs",
					mcp.Description("Filter monitor state by owning team."),
				),
				params.WithStringArray("signal_labels",
					mcp.Description("Only count alerts whose signal labels match all of these matchers, each of the form name=value or name!=value."),
				),
				params.WithStringArray("monitor_labels",
					mcp.Description("Only count alerts whose monitor labels match all of these matchers, each of the form name=value or name!=value."),
				),
				mcp.WithString("severity",
					mcp.Description("Only count alerts of this severity. If unset, counts both."),
					mcp.Enum("warn", "critical"),
				),
				mcp.WithBoolean("include_signal_statuses",
					mcp.Description("Include each monitor's firing signals in the response."),
					mcp.DefaultBool(false),
				),
				mcp.WithString("sort_by",
					mcp.Description("Sort results by a specific field. Currently supports 'SORT_BY_STATE', which orders by state, then severity, then slug. Results are ordered by monitor slug by default."),
					mcp.Enum(sortByState),
				),
			),
			Handler: func(ctx context.Context, request mcp.CallToolRequest) (*tools.Result, error) {
				monitorSlugs, err := params.StringArray(request, "monitor_slugs", false, nil)
				if err != nil {
					return nil, err
				}

				collectionSlugs, err := params.StringArray(request, "collection_slugs", false, nil)
				if err != nil {
					return nil, err
				}

				teamSlugs, err := params.StringArray(request, "team_slugs", false, nil)
				if err != nil {
					return nil, err
				}

				signalLabels, err := params.StringArray(request, "signal_labels", false, nil)
				if err != nil {
					return nil, err
				}

				monitorLabels, err := params.StringArray(request, "monitor_labels", false, nil)
				if err != nil {
					return nil, err
				}

				severity, err := params.String(request, "severity", false, "")
				if err != nil {
					return nil, err
				}

				includeSignalStatuses, err := params.Bool(request, "include_signal_statuses", false, false)
				if err != nil {
					return nil, err
				}

				sortBy, err := params.String(request, "sort_by", false, "")
				if err != nil {
					return nil, err
				}
				if sortBy != "" && sortBy != sortByState {
					return nil, fmt.Errorf("invalid sort_by %q: must be %s", sortBy, sortByState)
				}

				body := &models.Statev1ListMonitorStatusesRequest{
					MonitorSlugs:    monitorSlugs,
					CollectionSlugs: collectionSlugs,
					TeamSlugs:       teamSlugs,
					IncludeSignals:  includeSignalStatuses,
				}

				body.Severity, err = parseSeverity(severity)
				if err != nil {
					return nil, err
				}

				labelFilter, err := labelFilter(signalLabels, monitorLabels)
				if err != nil {
					return nil, err
				}
				if labelFilter != nil {
					body.LabelFilters = []*models.Datav1LabelFilter{labelFilter}
				}

				t.logger.Info("list monitor statuses", zap.Any("request", body))

				resp, err := t.api.MonitorStatuses.ListMonitorStatuses(
					monitor_statuses.NewListMonitorStatusesParams().
						WithContext(ctx).
						WithBody(body),
				)
				if err != nil {
					var statusErr interface{ Code() int }
					if errors.As(err, &statusErr) && statusErr.Code() == http.StatusNotFound {
						return nil, errors.New("monitor status API is not enabled for this tenant")
					}
					return nil, fmt.Errorf("failed to list monitor statuses: %s", err)
				}

				var firing []*models.Statev1MonitorStatus
				if resp.Payload != nil {
					firing = resp.Payload.MonitorStatuses
				}

				return &tools.Result{
					JSONContent: &ListMonitorStatusesResult{
						MonitorStatuses: toMonitorStatuses(firing, sortBy),
					},
				}, nil
			},
		},
	}
}

// parseSeverity maps the tool's severity parameter onto the API's severity.
func parseSeverity(severity string) (models.Datav1Severity, error) {
	switch severity {
	case "":
		return "", nil
	case "warn":
		return models.Datav1SeverityWARN, nil
	case "critical":
		return models.Datav1SeverityCRITICAL, nil
	default:
		return "", fmt.Errorf("invalid severity %q: must be warn or critical", severity)
	}
}

// labelFilter combines signal and monitor label matchers into one SUBSET_MATCH
// clause, or returns nil when there are none.
func labelFilter(signalLabels, monitorLabels []string) (*models.Datav1LabelFilter, error) {
	var matchers []*models.LabelFilterLabelMatcher
	for _, group := range []struct {
		param  string
		target models.LabelMatcherFilterTarget
		values []string
	}{
		{param: "signal_labels", target: models.LabelMatcherFilterTargetSIGNAL, values: signalLabels},
		{param: "monitor_labels", target: models.LabelMatcherFilterTargetMONITOR, values: monitorLabels},
	} {
		for _, value := range group.values {
			matcher, err := parseLabelMatcher(value)
			if err != nil {
				return nil, fmt.Errorf("invalid %s: %w", group.param, err)
			}
			matchers = append(matchers, &models.LabelFilterLabelMatcher{
				FilterTarget: group.target,
				Matcher:      matcher,
			})
		}
	}
	if len(matchers) == 0 {
		return nil, nil
	}
	return &models.Datav1LabelFilter{
		MatchStrategy: models.LabelFilterMatchStrategySUBSETMATCH,
		LabelMatchers: matchers,
	}, nil
}

// parseLabelMatcher parses a matcher of the form name=value or name!=value.
func parseLabelMatcher(s string) (*models.CommonPromQLMatcher, error) {
	idx := strings.Index(s, "=")
	if idx < 0 {
		return nil, fmt.Errorf("%q must be of the form name=value or name!=value", s)
	}
	name, value := s[:idx], s[idx+1:]
	matchType := models.CommonPromQLMatcherTypeMatchEqual
	if strings.HasSuffix(name, "!") {
		name = strings.TrimSuffix(name, "!")
		matchType = models.CommonPromQLMatcherTypeMatchNotEqual
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, fmt.Errorf("%q has an empty label name", s)
	}
	return &models.CommonPromQLMatcher{
		Name:  name,
		Type:  matchType,
		Value: value,
	}, nil
}

// toMonitorStatuses converts the firing monitors into the tool's result and sorts.
func toMonitorStatuses(firing []*models.Statev1MonitorStatus, sortBy string) []*MonitorStatus {
	statuses := make([]*MonitorStatus, 0, len(firing))
	for _, f := range firing {
		if f == nil {
			continue
		}
		state := stateAlerting
		if f.Muted {
			state = stateMuted
		}
		statuses = append(statuses, &MonitorStatus{
			MonitorSlug: f.MonitorSlug,
			State:       state,
			Muted:       f.Muted,
			Severity:    f.Severity,
			Signals:     f.Signals,
		})
	}

	sort.SliceStable(statuses, func(i, j int) bool {
		a, b := statuses[i], statuses[j]
		if sortBy == sortByState {
			if stateOrder[a.State] != stateOrder[b.State] {
				return stateOrder[a.State] < stateOrder[b.State]
			}
			if sevOrd(a.Severity) != sevOrd(b.Severity) {
				return sevOrd(a.Severity) < sevOrd(b.Severity)
			}
		}
		return a.MonitorSlug < b.MonitorSlug
	})
	return statuses
}

func sevOrd(s models.Datav1Severity) int {
	if ord, ok := severityOrder[s]; ok {
		return ord
	}
	return len(severityOrder)
}
