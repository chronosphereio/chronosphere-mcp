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
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"

	"github.com/chronosphereio/chronosphere-mcp/generated/dataunstable/models"
)

func TestAnalyzeAlertsResponseEdgeCases(t *testing.T) {
	tools := &Tools{logger: zaptest.NewLogger(t)}

	tests := []struct {
		name     string
		response *models.DataunstableAnalyzeAlertsResponse
		validate func(t *testing.T, result *AlertAnalysisSummary)
	}{
		{
			name:     "nil response",
			response: nil,
			validate: func(t *testing.T, result *AlertAnalysisSummary) {
				assert.NotNil(t, result)
				assert.Equal(t, 0, result.Summary.TotalAlerts)
				assert.Len(t, result.AnalysisNotes, 1)
				assert.Contains(t, result.AnalysisNotes[0], "No alerts found")
			},
		},
		{
			name: "empty alerts array",
			response: &models.DataunstableAnalyzeAlertsResponse{
				Alerts: []*models.AnalyzeAlertsResponseAlert{},
			},
			validate: func(t *testing.T, result *AlertAnalysisSummary) {
				assert.Equal(t, 0, result.Summary.TotalAlerts)
				assert.Empty(t, result.Alerts)
				assert.Len(t, result.AnalysisNotes, 1)
				assert.Contains(t, result.AnalysisNotes[0], "No alerts found")
			},
		},
		{
			name: "alerts with missing fields",
			response: &models.DataunstableAnalyzeAlertsResponse{
				Alerts: []*models.AnalyzeAlertsResponseAlert{
					{
						// Only AlertID populated
						AlertID: "minimal-alert",
					},
					{
						// No fields populated (should handle gracefully)
					},
				},
			},
			validate: func(t *testing.T, result *AlertAnalysisSummary) {
				assert.Equal(t, 2, result.Summary.TotalAlerts)
				assert.Equal(t, 0, result.Summary.UniqueMonitors) // Empty monitor slugs aren't counted
				assert.Len(t, result.Alerts, 2)

				// Verify minimal alert is handled correctly
				found := false
				for _, alert := range result.Alerts {
					if alert.AlertID == "minimal-alert" {
						found = true
						assert.Empty(t, alert.MonitorSlug)
						// strfmt.DateTime zero value becomes a formatted timestamp
						break
					}
				}
				assert.True(t, found, "Should contain minimal alert")
			},
		},
		{
			name: "alerts with empty timestamps",
			response: &models.DataunstableAnalyzeAlertsResponse{
				Alerts: []*models.AnalyzeAlertsResponseAlert{
					{
						AlertID: "empty-timestamp-alert",
						// StartTime and EndTime are empty (zero values)
					},
				},
			},
			validate: func(t *testing.T, result *AlertAnalysisSummary) {
				assert.Equal(t, 1, result.Summary.TotalAlerts)
				assert.Len(t, result.Alerts, 1)
				assert.Equal(t, "empty-timestamp-alert", result.Alerts[0].AlertID)
				// Duration calculation should handle empty timestamps gracefully
				assert.Equal(t, float64(0), result.Alerts[0].DurationMinutes)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tools.analyzeAlertsResponse(tt.response, "standard")
			tt.validate(t, result)
		})
	}
}

func TestDeduplicateStrings(t *testing.T) {
	tools := &Tools{logger: zaptest.NewLogger(t)}

	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "empty slice",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "no duplicates",
			input:    []string{"a", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "with duplicates",
			input:    []string{"a", "b", "a", "c", "b"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "all same",
			input:    []string{"same", "same", "same"},
			expected: []string{"same"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tools.deduplicateStrings(tt.input)
			assert.ElementsMatch(t, tt.expected, result)
		})
	}
}

func TestCreateLabelsHash(t *testing.T) {
	tools := &Tools{logger: zaptest.NewLogger(t)}

	tests := []struct {
		name     string
		labels   []*models.DataunstableAnalyzeAlertsResponseLabel
		expected string
	}{
		{
			name:     "empty labels",
			labels:   []*models.DataunstableAnalyzeAlertsResponseLabel{},
			expected: "",
		},
		{
			name:     "nil labels",
			labels:   nil,
			expected: "",
		},
		{
			name: "single label",
			labels: []*models.DataunstableAnalyzeAlertsResponseLabel{
				{Name: "service", Value: "api"},
			},
			expected: "service:api",
		},
		{
			name: "multiple labels sorted",
			labels: []*models.DataunstableAnalyzeAlertsResponseLabel{
				{Name: "env", Value: "prod"},
				{Name: "service", Value: "api"},
			},
			expected: "env:prod,service:api", // Should be sorted alphabetically
		},
		{
			name: "labels with nil entries",
			labels: []*models.DataunstableAnalyzeAlertsResponseLabel{
				{Name: "service", Value: "api"},
				nil, // Should be handled gracefully
				{Name: "env", Value: "prod"},
			},
			expected: "env:prod,service:api",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tools.createLabelsHash(tt.labels)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGroupNameAndToolRegistration(t *testing.T) {
	tools := &Tools{logger: zaptest.NewLogger(t)}

	// Test GroupName
	assert.Equal(t, "alertsanalysis", tools.GroupName())

	// Test MCPTools returns exactly one tool
	mcpTools := tools.MCPTools()
	assert.Len(t, mcpTools, 1)
	assert.Equal(t, "analyze_alerts", mcpTools[0].Metadata.Name)
	assert.NotNil(t, mcpTools[0].Handler)

	// Verify tool description mentions key capabilities
	description := mcpTools[0].Metadata.Description
	assert.Contains(t, description, "clustering")
	assert.Contains(t, description, "filtering")
	assert.Contains(t, description, "max_duration_seconds defaults to 3600")
}
