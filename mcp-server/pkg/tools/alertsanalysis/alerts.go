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

// Package alertsanalysis provides tools for analyzing Chronosphere alerts.
package alertsanalysis

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-openapi/strfmt"
	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"

	"github.com/chronosphereio/chronosphere-mcp/generated/dataunstable/dataunstable"
	"github.com/chronosphereio/chronosphere-mcp/generated/dataunstable/dataunstable/data_unstable"
	"github.com/chronosphereio/chronosphere-mcp/generated/dataunstable/models"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools/pkg/params"
)

var _ tools.MCPTools = (*Tools)(nil)

type Tools struct {
	logger          *zap.Logger
	dataUnstableAPI *dataunstable.DataUnstableAPI
}

func NewTools(
	dataUnstableAPI *dataunstable.DataUnstableAPI,
	logger *zap.Logger,
) (*Tools, error) {
	logger.Info("alerts tool configured")

	return &Tools{
		logger:          logger,
		dataUnstableAPI: dataUnstableAPI,
	}, nil
}

func (t *Tools) GroupName() string {
	return "alertsanalysis"
}

// AlertAnalysisSummary represents a comprehensive analysis of alert data
type AlertAnalysisSummary struct {
	Summary        AlertSummary       `json:"summary"`
	ClusteringInfo ClusteringInsights `json:"clustering_insights"`
	Patterns       AlertPatterns      `json:"patterns"`
	Alerts         []AlertDetail      `json:"alerts"`
	AnalysisNotes  []string           `json:"analysis_notes,omitempty"`
}

// AlertSummary provides high-level statistics about the alert set
type AlertSummary struct {
	TotalAlerts        int     `json:"total_alerts"`
	UniqueMonitors     int     `json:"unique_monitors"`
	MutedAlerts        int     `json:"muted_alerts"`
	OngoingAlerts      int     `json:"ongoing_alerts"`
	AvgDurationMinutes float64 `json:"avg_duration_minutes"`
	MaxDurationMinutes float64 `json:"max_duration_minutes"`
}

// ClusteringInsights provides information about alert relationships and groupings
type ClusteringInsights struct {
	TotalAlertGroups    int                 `json:"total_alert_groups"`
	LargestGroupSize    int                 `json:"largest_group_size"`
	SimilarSignalGroups []SignalGroup       `json:"similar_signal_groups"`
	MonitorLabelGroups  []MonitorLabelGroup `json:"monitor_label_groups"`
	CorrelationNotes    []string            `json:"correlation_notes,omitempty"`
}

// SignalGroup represents alerts grouped by similar signals
type SignalGroup struct {
	SignalsHash   string            `json:"signals_hash"`
	AlertCount    int               `json:"alert_count"`
	CommonSignals map[string]string `json:"common_signals"`
}

// MonitorLabelGroup represents alerts grouped by similar monitor labels
type MonitorLabelGroup struct {
	MonitorLabelsHash string            `json:"monitor_labels_hash"`
	AlertCount        int               `json:"alert_count"`
	CommonLabels      map[string]string `json:"common_labels"`
}

// AlertPatterns identifies interesting patterns in the alert data
type AlertPatterns struct {
	MostFrequentMonitors []MonitorFrequency `json:"most_frequent_monitors"`
	LongestAlerts        []LongAlert        `json:"longest_alerts"`
	RecentMutes          []string           `json:"recent_mutes,omitempty"`
	NotificationFailures []string           `json:"notification_failures,omitempty"`
}

// MonitorFrequency tracks monitor alert frequency
type MonitorFrequency struct {
	MonitorSlug string `json:"monitor_slug"`
	MonitorName string `json:"monitor_name,omitempty"`
	AlertCount  int    `json:"alert_count"`
}

// LongAlert represents alerts with notable duration
type LongAlert struct {
	MonitorSlug     string  `json:"monitor_slug"`
	AlertID         string  `json:"alert_id"`
	DurationMinutes float64 `json:"duration_minutes"`
	IsOngoing       bool    `json:"is_ongoing"`
}

// AlertDetail represents essential alert information for the response
type AlertDetail struct {
	AlertID                string                                           `json:"alert_id"`
	MonitorSlug            string                                           `json:"monitor_slug"`
	MonitorName            string                                           `json:"monitor_name,omitempty"`
	IsMuted                bool                                             `json:"is_muted"`
	StartTime              string                                           `json:"start_time,omitempty"`
	EndTime                string                                           `json:"end_time,omitempty"`
	DurationMinutes        float64                                          `json:"duration_minutes,omitempty"`
	CurrentSeverity        string                                           `json:"current_severity,omitempty"`
	Signal                 []*models.DataunstableAnalyzeAlertsResponseLabel `json:"signal,omitempty"`
	MonitorLabels          []*models.DataunstableAnalyzeAlertsResponseLabel `json:"monitor_labels,omitempty"`
	NotificationPolicy     string                                           `json:"notification_policy,omitempty"`
	SignalMatchCount       int32                                            `json:"signal_match_count,omitempty"`
	MonitorLabelMatchCount int32                                            `json:"monitor_label_match_count,omitempty"`
}

func (t *Tools) MCPTools() []tools.MCPTool {
	return []tools.MCPTool{
		{
			Metadata: tools.NewMetadata("analyze_alerts",
				mcp.WithDescription(`Analyze alert lifecycles with comprehensive filtering and clustering insights.

This tool queries alert data with sophisticated filtering options and provides observability insights 
including alert patterns, clustering information, and correlation analysis. 

Note: max_duration_seconds defaults to 3600 (1 hour) to align with backend time range limits and 
prevent analysis of extremely long-running alerts that may not be actionable.`),
				params.WithTimeRange(),

				// String matching filters
				mcp.WithString("monitor_slug_operation",
					mcp.Description("String operation for monitor slug: EQUAL, CONTAINS, REGEX_EQUAL, NOT_EQUAL, NOT_CONTAINS, NOT_REGEX_EQUAL, ANY_OF, NONE_OF"),
				),
				mcp.WithString("monitor_slug_value",
					mcp.Description("Monitor slug value to match (for single-value operations)"),
				),
				mcp.WithString("response_format",
					mcp.Description("Response detail level: 'compact' (summary + key insights only, ~50% fewer tokens), 'standard' (default), 'detailed' (full data). Default: compact"),
				),
				params.WithStringArray("monitor_slug_values",
					mcp.Description("Multiple monitor slug values (for ANY_OF, NONE_OF operations)"),
				),

				// Label-based filtering
				mcp.WithString("monitor_label_key",
					mcp.Description("Monitor label key to filter on"),
				),
				mcp.WithString("monitor_label_value_operation",
					mcp.Description("String operation for monitor label value: EQUAL, CONTAINS, REGEX_EQUAL, etc."),
				),
				mcp.WithString("monitor_label_value",
					mcp.Description("Monitor label value to match"),
				),
				params.WithStringArray("monitor_label_values",
					mcp.Description("Multiple monitor label values (for ANY_OF, NONE_OF operations)"),
				),

				// Signal filtering
				mcp.WithString("signal_key",
					mcp.Description("Signal/alert label key to filter on"),
				),
				mcp.WithString("signal_value_operation",
					mcp.Description("String operation for signal value: EQUAL, CONTAINS, REGEX_EQUAL, etc."),
				),
				mcp.WithString("signal_value",
					mcp.Description("Signal/alert label value to match"),
				),
				params.WithStringArray("signal_values",
					mcp.Description("Multiple signal values (for ANY_OF, NONE_OF operations)"),
				),

				// Notification policy filtering
				mcp.WithString("notification_policy_operation",
					mcp.Description("String operation for notification policy slug"),
				),
				mcp.WithString("notification_policy_value",
					mcp.Description("Notification policy slug to match"),
				),

				// Other filters
				mcp.WithBoolean("include_muted",
					mcp.Description("Include muted alerts (default: false - muted alerts excluded by default)"),
				),
				mcp.WithNumber("min_duration_seconds",
					mcp.Description("Minimum alert duration filter in seconds"),
				),
				mcp.WithNumber("max_duration_seconds",
					mcp.Description("Maximum alert duration filter in seconds (default: 3600 to align with 1-hour backend limits)"),
				),

				// Behavior controls
				mcp.WithString("filters_combine_mode",
					mcp.Description("How to combine multiple alert filters: 'AND' or 'OR' (default: AND)"),
				),
				mcp.WithNumber("limit",
					mcp.Description("Maximum number of alerts to return (default: 50)"),
				),
			),
			Handler: func(ctx context.Context, request mcp.CallToolRequest) (*tools.Result, error) {
				return t.handleAnalyzeAlerts(ctx, request)
			},
		},
	}
}

func (t *Tools) handleAnalyzeAlerts(ctx context.Context, request mcp.CallToolRequest) (*tools.Result, error) {
	// Parse time range
	timeRange, err := params.ParseTimeRange(request)
	if err != nil {
		return nil, err
	}

	// Build the request
	analyzeReq, err := t.buildAnalyzeAlertsRequest(ctx, request, timeRange)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}

	t.logger.Info("analyze alerts request", zap.Any("params", analyzeReq))

	// Call the API
	resp, err := t.dataUnstableAPI.DataUnstable.AnalyzeAlerts(analyzeReq)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze alerts: %w", err)
	}

	// Get response format preference
	responseFormat, err := params.String(request, "response_format", false, "compact")
	if err != nil {
		return nil, fmt.Errorf("invalid response_format parameter: %w", err)
	}

	// Process and analyze the response
	analysis := t.analyzeAlertsResponse(resp.Payload, responseFormat)

	return &tools.Result{
		JSONContent: analysis,
	}, nil
}

func (t *Tools) buildAnalyzeAlertsRequest(ctx context.Context, request mcp.CallToolRequest, timeRange *params.TimeRange) (*data_unstable.AnalyzeAlertsParams, error) {
	// Build time filter
	timeFilter := &models.AnalyzeAlertsRequestTimeFilter{
		StartTime: strfmt.DateTime(timeRange.Start),
		EndTime:   strfmt.DateTime(timeRange.End),
	}

	// Note: Duration filtering would need to be handled client-side since the API
	// doesn't expose min/max duration filters in the current swagger spec

	// Build alert filters
	var alertFilters []*models.AnalyzeAlertsRequestAlertFilter

	// Monitor slug filter
	monitorFilter, err := t.buildMonitorSlugFilter(request)
	if err != nil {
		return nil, fmt.Errorf("invalid monitor slug filter: %w", err)
	}
	if monitorFilter != nil {
		alertFilters = append(alertFilters, monitorFilter)
	}

	// Monitor labels filter
	monitorLabelsFilter, err := t.buildMonitorLabelsFilter(request)
	if err != nil {
		return nil, fmt.Errorf("invalid monitor labels filter: %w", err)
	}
	if monitorLabelsFilter != nil {
		alertFilters = append(alertFilters, monitorLabelsFilter)
	}

	// Signals filter
	signalsFilter, err := t.buildSignalsFilter(request)
	if err != nil {
		return nil, fmt.Errorf("invalid signals filter: %w", err)
	}
	if signalsFilter != nil {
		alertFilters = append(alertFilters, signalsFilter)
	}

	// Notification policy filter
	notificationPolicyFilter, err := t.buildNotificationPolicyFilter(request)
	if err != nil {
		return nil, fmt.Errorf("invalid notification policy filter: %w", err)
	}
	if notificationPolicyFilter != nil {
		alertFilters = append(alertFilters, notificationPolicyFilter)
	}

	// Mute filter
	muteFilter, err := t.buildMuteFilter(request)
	if err != nil {
		return nil, fmt.Errorf("invalid mute filter: %w", err)
	}
	if muteFilter != nil {
		alertFilters = append(alertFilters, muteFilter)
	}

	// Combine mode
	combineMode := models.AnalyzeAlertsRequestAlertFiltersCombineModeALERTFILTERSOR
	combineModeStr, err := params.String(request, "filters_combine_mode", false, "AND")
	if err != nil {
		return nil, err
	}
	if strings.ToUpper(combineModeStr) == "AND" {
		combineMode = "" // Default is AND, represented as empty/unset
	}

	// Build the request
	analyzeReq := &data_unstable.AnalyzeAlertsParams{
		Context: ctx, // Use the request context for proper cancellation
		Body: &models.DataunstableAnalyzeAlertsRequest{
			TimeFilter:              timeFilter,
			AlertFilters:            alertFilters,
			AlertFiltersCombineMode: combineMode,
		},
	}

	return analyzeReq, nil
}

// Helper functions for building specific filters will be implemented next...
