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

package alertsanalysis

import (
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/chronosphereio/chronosphere-mcp/generated/dataunstable/models"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools/pkg/params"
)

// buildMonitorSlugFilter builds a monitor slug filter from request parameters
func (t *Tools) buildMonitorSlugFilter(request mcp.CallToolRequest) (*models.AnalyzeAlertsRequestAlertFilter, error) {
	operation, err := params.String(request, "monitor_slug_operation", false, "")
	if err != nil {
		return nil, err
	}

	if operation == "" {
		return nil, nil // No filter specified
	}

	singleValue, err := params.String(request, "monitor_slug_value", false, "")
	if err != nil {
		return nil, err
	}

	multiValues, err := params.StringArray(request, "monitor_slug_values", false, nil)
	if err != nil {
		return nil, err
	}

	stringMatcher, err := buildStringMatcher(operation, singleValue, multiValues)
	if err != nil {
		return nil, fmt.Errorf("invalid monitor slug matcher: %w", err)
	}

	return &models.AnalyzeAlertsRequestAlertFilter{
		MonitorSlug: stringMatcher,
	}, nil
}

// buildMonitorLabelsFilter builds a monitor labels filter from request parameters
func (t *Tools) buildMonitorLabelsFilter(request mcp.CallToolRequest) (*models.AnalyzeAlertsRequestAlertFilter, error) {
	key, err := params.String(request, "monitor_label_key", false, "")
	if err != nil {
		return nil, err
	}

	if key == "" {
		return nil, nil // No filter specified
	}

	operation, err := params.String(request, "monitor_label_value_operation", false, "EQUAL")
	if err != nil {
		return nil, err
	}

	singleValue, err := params.String(request, "monitor_label_value", false, "")
	if err != nil {
		return nil, err
	}

	multiValues, err := params.StringArray(request, "monitor_label_values", false, nil)
	if err != nil {
		return nil, err
	}

	stringMatcher, err := buildStringMatcher(operation, singleValue, multiValues)
	if err != nil {
		return nil, fmt.Errorf("invalid monitor label value matcher: %w", err)
	}

	labelMatcher := &models.AnalyzeAlertsRequestLabelMatcher{
		Key:   key,
		Value: stringMatcher,
	}

	labelsMatcher := &models.AnalyzeAlertsRequestLabelsMatcher{
		Labels: []*models.AnalyzeAlertsRequestLabelMatcher{labelMatcher},
		// MinMatches defaults to 1 when not specified
	}

	return &models.AnalyzeAlertsRequestAlertFilter{
		MonitorLabels: labelsMatcher,
	}, nil
}

// buildSignalsFilter builds a signals filter from request parameters
func (t *Tools) buildSignalsFilter(request mcp.CallToolRequest) (*models.AnalyzeAlertsRequestAlertFilter, error) {
	key, err := params.String(request, "signal_key", false, "")
	if err != nil {
		return nil, err
	}

	if key == "" {
		return nil, nil // No filter specified
	}

	operation, err := params.String(request, "signal_value_operation", false, "EQUAL")
	if err != nil {
		return nil, err
	}

	singleValue, err := params.String(request, "signal_value", false, "")
	if err != nil {
		return nil, err
	}

	multiValues, err := params.StringArray(request, "signal_values", false, nil)
	if err != nil {
		return nil, err
	}

	stringMatcher, err := buildStringMatcher(operation, singleValue, multiValues)
	if err != nil {
		return nil, fmt.Errorf("invalid signal value matcher: %w", err)
	}

	labelMatcher := &models.AnalyzeAlertsRequestLabelMatcher{
		Key:   key,
		Value: stringMatcher,
	}

	labelsMatcher := &models.AnalyzeAlertsRequestLabelsMatcher{
		Labels: []*models.AnalyzeAlertsRequestLabelMatcher{labelMatcher},
	}

	return &models.AnalyzeAlertsRequestAlertFilter{
		Signals: labelsMatcher,
	}, nil
}

// buildNotificationPolicyFilter builds a notification policy filter from request parameters
func (t *Tools) buildNotificationPolicyFilter(request mcp.CallToolRequest) (*models.AnalyzeAlertsRequestAlertFilter, error) {
	operation, err := params.String(request, "notification_policy_operation", false, "")
	if err != nil {
		return nil, err
	}

	if operation == "" {
		return nil, nil // No filter specified
	}

	singleValue, err := params.String(request, "notification_policy_value", false, "")
	if err != nil {
		return nil, err
	}

	stringMatcher, err := buildStringMatcher(operation, singleValue, nil)
	if err != nil {
		return nil, fmt.Errorf("invalid notification policy matcher: %w", err)
	}

	return &models.AnalyzeAlertsRequestAlertFilter{
		NotificationPolicySlug: stringMatcher,
	}, nil
}

// buildMuteFilter builds a mute status filter from request parameters
func (t *Tools) buildMuteFilter(request mcp.CallToolRequest) (*models.AnalyzeAlertsRequestAlertFilter, error) {
	includeMuted, err := params.Bool(request, "include_muted", false, false)
	if err != nil {
		return nil, err
	}

	// If include_muted is true, we don't add a filter
	// If include_muted is false (default), we filter to only show non-muted alerts
	if includeMuted {
		return nil, nil
	}

	return &models.AnalyzeAlertsRequestAlertFilter{
		IsMuted: false,
	}, nil
}

// buildStringMatcher builds a string matcher from operation and values
func buildStringMatcher(operation, singleValue string, multiValues []string) (*models.AnalyzeAlertsRequestStringMatcher, error) {
	if operation == "" {
		return nil, fmt.Errorf("operation is required")
	}

	// Convert operation string to enum
	var op models.AnalyzeAlertsRequestStringMatcherOperation
	switch strings.ToUpper(operation) {
	case "EQUAL":
		op = models.AnalyzeAlertsRequestStringMatcherOperationEQUAL
	case "REGEX_EQUAL":
		op = models.AnalyzeAlertsRequestStringMatcherOperationREGEXEQUAL
	case "CONTAINS":
		op = models.AnalyzeAlertsRequestStringMatcherOperationCONTAINS
	case "NOT_EQUAL":
		op = models.AnalyzeAlertsRequestStringMatcherOperationNOTEQUAL
	case "NOT_REGEX_EQUAL":
		op = models.AnalyzeAlertsRequestStringMatcherOperationNOTREGEXEQUAL
	case "NOT_CONTAINS":
		op = models.AnalyzeAlertsRequestStringMatcherOperationNOTCONTAINS
	case "ANY_OF":
		op = models.AnalyzeAlertsRequestStringMatcherOperationANYOF
	case "NONE_OF":
		op = models.AnalyzeAlertsRequestStringMatcherOperationNONEOF
	default:
		return nil, fmt.Errorf("invalid operation: %s", operation)
	}

	matcher := &models.AnalyzeAlertsRequestStringMatcher{
		Operation: op,
	}

	// Set values based on operation type
	switch op {
	case models.AnalyzeAlertsRequestStringMatcherOperationANYOF,
		models.AnalyzeAlertsRequestStringMatcherOperationNONEOF:
		if len(multiValues) == 0 {
			return nil, fmt.Errorf("operation %s requires multiple values", operation)
		}
		matcher.MultiValue = multiValues
	default:
		if singleValue == "" {
			return nil, fmt.Errorf("operation %s requires a single value", operation)
		}
		matcher.SingleValue = singleValue
	}

	return matcher, nil
}
