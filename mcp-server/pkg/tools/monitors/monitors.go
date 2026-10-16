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

// Package monitors provides tools for querying Chronosphere monitor statuses and alerts.
package monitors

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"net/http"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/go-openapi/strfmt"
	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"

	"github.com/chronosphereio/chronosphere-mcp/generated/datav1/datav1"
	"github.com/chronosphereio/chronosphere-mcp/generated/datav1/datav1/version1"
	datav1models "github.com/chronosphereio/chronosphere-mcp/generated/datav1/models"
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
	logger   *zap.Logger
	stateAPI *statev1.StateV1API
	dataAPI  *datav1.DataV1API
}

// NewTools creates a new set of monitor tools.
func NewTools(stateAPI *statev1.StateV1API, dataAPI *datav1.DataV1API, logger *zap.Logger) (*Tools, error) {
	logger.Info("monitor tools configured")

	return &Tools{
		logger:   logger,
		stateAPI: stateAPI,
		dataAPI:  dataAPI,
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

				resp, err := t.stateAPI.MonitorStatuses.ListMonitorStatuses(
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
		{
			Metadata: tools.NewMetadata("list_alerts",
				mcp.WithDescription("Lists alerts raised by monitors and SLOs in Chronosphere, newest first by default. The five search filters (monitor_slugs, monitor_name_contains, slo_slugs, slo_name_contains and alert_title_contains) are OR'ed with each other: an alert matches the search if it matches any one of them, so adding a search filter widens the results. The search as a whole is AND'ed with every other filter. Returns total_count (matching alerts across all pages), the alerts on this page, and next_page_token when there are more."),
				params.WithTimeRange(),
				mcp.WithString("time_filter",
					mcp.Description("Which alerts the time range selects: 'triggered' (default) returns alerts that started firing in the range, 'active' returns alerts that were firing at any point in the range."),
					mcp.Enum(timeFilterTriggered, timeFilterActive),
				),
				mcp.WithString("firing_status",
					mcp.Description("Return only alerts that are still firing or that have resolved. If unset, returns both."),
					mcp.Enum(slices.Sorted(maps.Keys(firingStatuses))...),
				),
				mcp.WithString("severity",
					mcp.Description("Return only alerts of this severity. If unset, returns both."),
					mcp.Enum(slices.Sorted(maps.Keys(alertSeverities))...),
				),
				mcp.WithString("mute_status",
					mcp.Description("Return only muted or only unmuted alerts. If unset, returns both."),
					mcp.Enum(slices.Sorted(maps.Keys(muteStatuses))...),
				),
				mcp.WithString("source_type",
					mcp.Description("Return only alerts raised by monitors or only alerts raised by SLOs. If unset, returns both."),
					mcp.Enum(slices.Sorted(maps.Keys(sourceTypes))...),
				),
				params.WithStringArray("monitor_slugs",
					mcp.Description("Search filter: match alerts raised by monitors with these slugs. OR'ed with the other search filters."),
				),
				mcp.WithString("monitor_name_contains",
					mcp.Description("Search filter: match alerts raised by monitors whose name or slug contains this string, case-insensitively. OR'ed with the other search filters."),
				),
				params.WithStringArray("slo_slugs",
					mcp.Description("Search filter: match alerts raised by SLOs with these slugs. OR'ed with the other search filters."),
				),
				mcp.WithString("slo_name_contains",
					mcp.Description("Search filter: match alerts raised by SLOs whose name or slug contains this string, case-insensitively. OR'ed with the other search filters."),
				),
				mcp.WithString("alert_title_contains",
					mcp.Description("Search filter: match alerts whose title contains this string, case-insensitively. OR'ed with the other search filters."),
				),
				params.WithStringArray("collection_slugs",
					mcp.Description("Return only alerts whose monitor or SLO belongs to one of these collections."),
				),
				params.WithStringArray("team_slugs",
					mcp.Description("Return only alerts owned by one of these teams."),
				),
				params.WithStringArray("notifier_slugs",
					mcp.Description("Return only alerts a notification was attempted against one of these notifiers for, including failed sends."),
				),
				params.WithStringArray("signal_labels",
					mcp.Description("Return only alerts whose signal labels match all of these matchers, each of the form name=value or name!=value."),
				),
				params.WithStringArray("monitor_labels",
					mcp.Description("Return only alerts whose monitor or SLO labels match all of these matchers, each of the form name=value or name!=value."),
				),
				params.WithStringArray("alert_ids",
					mcp.Description("Return only the alerts with these IDs."),
				),
				mcp.WithString("sort",
					mcp.Description("Order by trigger time: 'newest' (default) or 'oldest' first."),
					mcp.Enum(slices.Sorted(maps.Keys(sortOrders))...),
				),
				mcp.WithNumber("page_size",
					mcp.Description("Maximum number of alerts per page. Default is server-defined."),
				),
				mcp.WithString("page_token",
					mcp.Description("Pagination token from a previous response's next_page_token."),
				),
			),
			Handler: t.listAlerts,
		},
	}
}

// Values of list_alerts' time_filter parameter.
const (
	timeFilterTriggered = "triggered"
	timeFilterActive    = "active"
)

var (
	firingStatuses = map[string]datav1models.Datav1FiringStatus{
		"firing":   datav1models.Datav1FiringStatusFIRING,
		"resolved": datav1models.Datav1FiringStatusRESOLVED,
	}
	alertSeverities = map[string]datav1models.Datav1Severity{
		"warn":     datav1models.Datav1SeverityWARN,
		"critical": datav1models.Datav1SeverityCRITICAL,
	}
	muteStatuses = map[string]datav1models.Datav1MuteStatus{
		"muted":     datav1models.Datav1MuteStatusMUTED,
		"not_muted": datav1models.Datav1MuteStatusNOTMUTED,
	}
	sourceTypes = map[string]datav1models.Datav1SourceType{
		"monitor": datav1models.Datav1SourceTypeMONITOR,
		"slo":     datav1models.Datav1SourceTypeSLO,
	}
	sortOrders = map[string]datav1models.ListAlertsRequestSortOrder{
		"newest": datav1models.ListAlertsRequestSortOrderTRIGGEREDDESC,
		"oldest": datav1models.ListAlertsRequestSortOrderTRIGGEREDASC,
	}
)

// listAlertsArgs holds list_alerts' parameters as the caller sent them.
type listAlertsArgs struct {
	timeRange           params.TimeRange
	timeFilter          string
	firingStatus        string
	severity            string
	muteStatus          string
	sourceType          string
	monitorSlugs        []string
	monitorNameContains string
	sloSlugs            []string
	sloNameContains     string
	alertTitleContains  string
	collectionSlugs     []string
	teamSlugs           []string
	notifierSlugs       []string
	signalLabels        []string
	monitorLabels       []string
	alertIDs            []string
	sort                string
	pageSize            int
	pageToken           string
}

// Alert is one alert as returned by list_alerts.
type Alert struct {
	*datav1models.Datav1Alert
	// ResolvedAt shadows the embedded field so that a firing alert omits it
	// rather than reporting the zero time.
	ResolvedAt *strfmt.DateTime `json:"resolved_at,omitempty"`
}

// ListAlertsResult is the result of list_alerts.
type ListAlertsResult struct {
	TotalCount    int64    `json:"total_count"`
	Alerts        []*Alert `json:"alerts"`
	NextPageToken string   `json:"next_page_token,omitempty"`
}

func (t *Tools) listAlerts(ctx context.Context, request mcp.CallToolRequest) (*tools.Result, error) {
	args, err := parseListAlertsArgs(request)
	if err != nil {
		return nil, err
	}
	body, err := listAlertsRequest(args)
	if err != nil {
		return nil, err
	}

	t.logger.Info("list alerts", zap.Any("request", body))

	resp, err := t.dataAPI.Version1.ListAlerts(
		version1.NewListAlertsParams().
			WithContext(ctx).
			WithBody(body),
	)
	if err != nil {
		var statusErr interface{ Code() int }
		if errors.As(err, &statusErr) && statusErr.Code() == http.StatusNotFound {
			return nil, errors.New("alerts API is not enabled for this tenant")
		}
		return nil, fmt.Errorf("failed to list alerts: %s", err)
	}

	result, err := toListAlertsResult(resp.Payload)
	if err != nil {
		return nil, err
	}
	return &tools.Result{JSONContent: result}, nil
}

// parseListAlertsArgs reads list_alerts' parameters from the request.
func parseListAlertsArgs(request mcp.CallToolRequest) (listAlertsArgs, error) {
	var (
		args listAlertsArgs
		err  error
	)
	str := func(key string) string {
		if err != nil {
			return ""
		}
		var v string
		v, err = params.String(request, key, false, "")
		return v
	}
	strs := func(key string) []string {
		if err != nil {
			return nil
		}
		var v []string
		v, err = params.StringArray(request, key, false, nil)
		return v
	}

	args.timeFilter = str("time_filter")
	args.firingStatus = str("firing_status")
	args.severity = str("severity")
	args.muteStatus = str("mute_status")
	args.sourceType = str("source_type")
	args.monitorSlugs = strs("monitor_slugs")
	args.monitorNameContains = str("monitor_name_contains")
	args.sloSlugs = strs("slo_slugs")
	args.sloNameContains = str("slo_name_contains")
	args.alertTitleContains = str("alert_title_contains")
	args.collectionSlugs = strs("collection_slugs")
	args.teamSlugs = strs("team_slugs")
	args.notifierSlugs = strs("notifier_slugs")
	args.signalLabels = strs("signal_labels")
	args.monitorLabels = strs("monitor_labels")
	args.alertIDs = strs("alert_ids")
	args.sort = str("sort")
	args.pageToken = str("page_token")
	if err != nil {
		return listAlertsArgs{}, err
	}

	if args.pageSize, err = params.Int(request, "page_size", false, 0); err != nil {
		return listAlertsArgs{}, err
	}
	timeRange, err := params.ParseTimeRange(request)
	if err != nil {
		return listAlertsArgs{}, err
	}
	args.timeRange = *timeRange
	return args, nil
}

// listAlertsRequest maps list_alerts' parameters onto the ListAlerts request.
func listAlertsRequest(args listAlertsArgs) (*datav1models.Datav1ListAlertsRequest, error) {
	filter := &datav1models.Datav1AlertsFilter{
		MonitorSlugs:          args.monitorSlugs,
		MonitorNameContains:   args.monitorNameContains,
		SloSlugs:              args.sloSlugs,
		SloNameContains:       args.sloNameContains,
		AlertTitleContains:    args.alertTitleContains,
		SourceCollectionSlugs: args.collectionSlugs,
		SourceTeamSlugs:       args.teamSlugs,
		NotifierSlugs:         args.notifierSlugs,
		AlertIds:              args.alertIDs,
	}

	timeRange := &datav1models.Datav1TimeRange{
		Start: strfmt.DateTime(args.timeRange.Start),
		End:   strfmt.DateTime(args.timeRange.End),
	}
	switch args.timeFilter {
	case "", timeFilterTriggered:
		filter.TriggerTimeRange = timeRange
	case timeFilterActive:
		filter.ActiveTimeRange = timeRange
	default:
		return nil, fmt.Errorf("invalid time_filter %q: must be %s or %s", args.timeFilter, timeFilterTriggered, timeFilterActive)
	}

	var err error
	if filter.FiringStatus, err = parseEnum("firing_status", args.firingStatus, firingStatuses); err != nil {
		return nil, err
	}
	if filter.Severity, err = parseEnum("severity", args.severity, alertSeverities); err != nil {
		return nil, err
	}
	if filter.MuteStatus, err = parseEnum("mute_status", args.muteStatus, muteStatuses); err != nil {
		return nil, err
	}
	if filter.SourceType, err = parseEnum("source_type", args.sourceType, sourceTypes); err != nil {
		return nil, err
	}

	labelFilter, err := alertLabelFilter(args.signalLabels, args.monitorLabels)
	if err != nil {
		return nil, err
	}
	if labelFilter != nil {
		filter.LabelFilters = []*datav1models.Datav1LabelFilter{labelFilter}
	}

	sortOrder := args.sort
	if sortOrder == "" {
		sortOrder = "newest"
	}
	req := &datav1models.Datav1ListAlertsRequest{Filter: filter}
	if req.SortOrder, err = parseEnum("sort", sortOrder, sortOrders); err != nil {
		return nil, err
	}

	if args.pageSize < 0 {
		return nil, fmt.Errorf("invalid page_size %d: must not be negative", args.pageSize)
	}
	if args.pageSize > 0 || args.pageToken != "" {
		req.Page = &datav1models.Configv1PageParams{
			MaxSize: int64(args.pageSize),
			Token:   args.pageToken,
		}
	}
	return req, nil
}

// parseEnum maps a tool parameter onto its API enum value. An empty value maps
// to the enum's zero value, which leaves the filter unset.
func parseEnum[T ~string](param, value string, values map[string]T) (T, error) {
	if value == "" {
		return "", nil
	}
	v, ok := values[value]
	if !ok {
		return "", fmt.Errorf("invalid %s %q: must be one of %s", param, value, strings.Join(slices.Sorted(maps.Keys(values)), ", "))
	}
	return v, nil
}

// alertLabelFilter is labelFilter for the Data V1 types ListAlerts takes.
func alertLabelFilter(signalLabels, monitorLabels []string) (*datav1models.Datav1LabelFilter, error) {
	filter, err := labelFilter(signalLabels, monitorLabels)
	if err != nil || filter == nil {
		return nil, err
	}
	matchers := make([]*datav1models.LabelFilterLabelMatcher, 0, len(filter.LabelMatchers))
	for _, m := range filter.LabelMatchers {
		matchers = append(matchers, &datav1models.LabelFilterLabelMatcher{
			FilterTarget: datav1models.LabelMatcherFilterTarget(m.FilterTarget),
			Matcher: &datav1models.CommonPromQLMatcher{
				Name:  m.Matcher.Name,
				Type:  datav1models.CommonPromQLMatcherType(m.Matcher.Type),
				Value: m.Matcher.Value,
			},
		})
	}
	return &datav1models.Datav1LabelFilter{
		MatchStrategy: datav1models.LabelFilterMatchStrategy(filter.MatchStrategy),
		LabelMatchers: matchers,
	}, nil
}

// toListAlertsResult converts the ListAlerts response into the tool's result.
func toListAlertsResult(resp *datav1models.Datav1ListAlertsResponse) (*ListAlertsResult, error) {
	result := &ListAlertsResult{Alerts: []*Alert{}}
	if resp == nil {
		return result, nil
	}
	if resp.TotalCount != "" {
		count, err := strconv.ParseInt(resp.TotalCount, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid total_count %q in response: %w", resp.TotalCount, err)
		}
		result.TotalCount = count
	}
	if resp.Page != nil {
		result.NextPageToken = resp.Page.NextToken
	}
	for _, a := range resp.Alerts {
		if a == nil {
			continue
		}
		alert := &Alert{Datav1Alert: a}
		if !a.ResolvedAt.IsZero() {
			resolvedAt := a.ResolvedAt
			alert.ResolvedAt = &resolvedAt
		}
		result.Alerts = append(result.Alerts, alert)
	}
	return result, nil
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
