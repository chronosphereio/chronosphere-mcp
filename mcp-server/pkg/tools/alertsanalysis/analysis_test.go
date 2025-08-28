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
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"github.com/chronosphereio/chronosphere-mcp/generated/dataunstable/models"
)

func TestCalculateSummary(t *testing.T) {
	tools := &Tools{logger: zaptest.NewLogger(t)}

	now := time.Now()
	oneHourAgo := now.Add(-1 * time.Hour)
	twoHoursAgo := now.Add(-2 * time.Hour)

	tests := []struct {
		name     string
		alerts   []*models.AnalyzeAlertsResponseAlert
		expected AlertSummary
	}{
		{
			name:     "empty alerts",
			alerts:   []*models.AnalyzeAlertsResponseAlert{},
			expected: AlertSummary{},
		},
		{
			name: "single ongoing alert",
			alerts: []*models.AnalyzeAlertsResponseAlert{
				{
					AlertID:     "alert-1",
					MonitorSlug: "cpu-high",
					IsMuted:     false,
					StartTime:   strfmt.DateTime(oneHourAgo),
					// No EndTime = ongoing
				},
			},
			expected: AlertSummary{
				TotalAlerts:        1,
				UniqueMonitors:     1,
				MutedAlerts:        0,
				OngoingAlerts:      1,
				AvgDurationMinutes: 60, // approximately 1 hour
				MaxDurationMinutes: 60,
			},
		},
		{
			name: "mixed resolved and ongoing alerts",
			alerts: []*models.AnalyzeAlertsResponseAlert{
				{
					AlertID:     "alert-1",
					MonitorSlug: "cpu-high",
					IsMuted:     true,
					StartTime:   strfmt.DateTime(twoHoursAgo),
					EndTime:     strfmt.DateTime(oneHourAgo),
				},
				{
					AlertID:     "alert-2",
					MonitorSlug: "memory-high",
					IsMuted:     false,
					StartTime:   strfmt.DateTime(oneHourAgo),
					// ongoing
				},
				{
					AlertID:     "alert-3",
					MonitorSlug: "cpu-high", // same monitor as alert-1
					IsMuted:     false,
					StartTime:   strfmt.DateTime(now.Add(-30 * time.Minute)),
					EndTime:     strfmt.DateTime(now.Add(-15 * time.Minute)),
				},
			},
			expected: AlertSummary{
				TotalAlerts:        3,
				UniqueMonitors:     2, // cpu-high and memory-high
				MutedAlerts:        1,
				OngoingAlerts:      1,
				AvgDurationMinutes: 45, // (60 + 60 + 15) / 3 = 45
				MaxDurationMinutes: 60,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tools.calculateSummary(tt.alerts)

			assert.Equal(t, tt.expected.TotalAlerts, result.TotalAlerts)
			assert.Equal(t, tt.expected.UniqueMonitors, result.UniqueMonitors)
			assert.Equal(t, tt.expected.MutedAlerts, result.MutedAlerts)
			assert.Equal(t, tt.expected.OngoingAlerts, result.OngoingAlerts)

			// Duration calculations are approximate due to timing, so we use ranges
			if tt.expected.AvgDurationMinutes > 0 {
				assert.InDelta(t, tt.expected.AvgDurationMinutes, result.AvgDurationMinutes, 5) // ±5 minutes
			}
			if tt.expected.MaxDurationMinutes > 0 {
				assert.InDelta(t, tt.expected.MaxDurationMinutes, result.MaxDurationMinutes, 5)
			}
		})
	}
}

func TestAnalyzeClusteringInfo(t *testing.T) {
	tools := &Tools{logger: zaptest.NewLogger(t)}

	tests := []struct {
		name     string
		alerts   []*models.AnalyzeAlertsResponseAlert
		expected ClusteringInsights
	}{
		{
			name:   "empty alerts",
			alerts: []*models.AnalyzeAlertsResponseAlert{},
			expected: ClusteringInsights{
				TotalAlertGroups:    0,
				LargestGroupSize:    0,
				SimilarSignalGroups: []SignalGroup{},
				MonitorLabelGroups:  []MonitorLabelGroup{},
			},
		},
		{
			name: "alerts with similar signals",
			alerts: []*models.AnalyzeAlertsResponseAlert{
				{
					AlertID:     "alert-1",
					SignalsHash: "hash123",
					Signal: []*models.DataunstableAnalyzeAlertsResponseLabel{
						{Name: "service", Value: "api"},
						{Name: "env", Value: "prod"},
					},
				},
				{
					AlertID:     "alert-2",
					SignalsHash: "hash123", // same hash
					Signal: []*models.DataunstableAnalyzeAlertsResponseLabel{
						{Name: "service", Value: "api"},
						{Name: "env", Value: "prod"},
					},
				},
				{
					AlertID:     "alert-3",
					SignalsHash: "hash456", // different hash
					Signal: []*models.DataunstableAnalyzeAlertsResponseLabel{
						{Name: "service", Value: "web"},
					},
				},
			},
			expected: ClusteringInsights{
				TotalAlertGroups: 2, // Two signal groups
				LargestGroupSize: 2, // hash123 has 2 alerts
				SimilarSignalGroups: []SignalGroup{
					{
						SignalsHash: "hash123",
						AlertCount:  2,
						CommonSignals: map[string]string{
							"service": "api",
							"env":     "prod",
						},
					},
					{
						SignalsHash: "hash456",
						AlertCount:  1,
						CommonSignals: map[string]string{
							"service": "web",
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tools.analyzeClusteringInfo(tt.alerts)

			assert.Equal(t, tt.expected.TotalAlertGroups, result.TotalAlertGroups)
			assert.Equal(t, tt.expected.LargestGroupSize, result.LargestGroupSize)

			// For signal groups, verify key properties (order may vary due to sorting)
			assert.Len(t, result.SimilarSignalGroups, len(tt.expected.SimilarSignalGroups))

			if len(tt.expected.SimilarSignalGroups) > 0 {
				// Find the largest group and verify its properties
				var largestGroup *SignalGroup
				for i := range result.SimilarSignalGroups {
					if largestGroup == nil || result.SimilarSignalGroups[i].AlertCount > largestGroup.AlertCount {
						largestGroup = &result.SimilarSignalGroups[i]
					}
				}

				if largestGroup != nil {
					expectedLargest := tt.expected.SimilarSignalGroups[0] // assumes first is largest in test data
					assert.Equal(t, expectedLargest.AlertCount, largestGroup.AlertCount)
					assert.Equal(t, expectedLargest.SignalsHash, largestGroup.SignalsHash)
				}
			}
		})
	}
}

func TestIdentifyPatterns(t *testing.T) {
	tools := &Tools{logger: zaptest.NewLogger(t)}

	now := time.Now()
	oneHourAgo := now.Add(-1 * time.Hour)
	twoHoursAgo := now.Add(-2 * time.Hour)

	tests := []struct {
		name     string
		alerts   []*models.AnalyzeAlertsResponseAlert
		validate func(t *testing.T, patterns AlertPatterns)
	}{
		{
			name:   "empty alerts",
			alerts: []*models.AnalyzeAlertsResponseAlert{},
			validate: func(t *testing.T, patterns AlertPatterns) {
				assert.Empty(t, patterns.MostFrequentMonitors)
				assert.Empty(t, patterns.LongestAlerts)
				assert.Empty(t, patterns.RecentMutes)
			},
		},
		{
			name: "monitor frequency analysis",
			alerts: []*models.AnalyzeAlertsResponseAlert{
				{AlertID: "1", MonitorSlug: "cpu-high", MonitorName: "CPU High"},
				{AlertID: "2", MonitorSlug: "cpu-high", MonitorName: "CPU High"},
				{AlertID: "3", MonitorSlug: "memory-high", MonitorName: "Memory High"},
				{AlertID: "4", MonitorSlug: "cpu-high", MonitorName: "CPU High"},
			},
			validate: func(t *testing.T, patterns AlertPatterns) {
				require.Len(t, patterns.MostFrequentMonitors, 2)

				// Should be sorted by frequency, cpu-high first (3 alerts)
				assert.Equal(t, "cpu-high", patterns.MostFrequentMonitors[0].MonitorSlug)
				assert.Equal(t, "CPU High", patterns.MostFrequentMonitors[0].MonitorName)
				assert.Equal(t, 3, patterns.MostFrequentMonitors[0].AlertCount)

				assert.Equal(t, "memory-high", patterns.MostFrequentMonitors[1].MonitorSlug)
				assert.Equal(t, 1, patterns.MostFrequentMonitors[1].AlertCount)
			},
		},
		{
			name: "long running alerts detection",
			alerts: []*models.AnalyzeAlertsResponseAlert{
				{
					AlertID:     "short-alert",
					MonitorSlug: "quick-monitor",
					StartTime:   strfmt.DateTime(now.Add(-5 * time.Minute)),
					EndTime:     strfmt.DateTime(now),
				},
				{
					AlertID:     "long-alert",
					MonitorSlug: "slow-monitor",
					StartTime:   strfmt.DateTime(twoHoursAgo),
					EndTime:     strfmt.DateTime(oneHourAgo),
				},
				{
					AlertID:     "ongoing-long",
					MonitorSlug: "ongoing-monitor",
					StartTime:   strfmt.DateTime(now.Add(-45 * time.Minute)),
					// No end time = ongoing
				},
			},
			validate: func(t *testing.T, patterns AlertPatterns) {
				// Should detect 2 long alerts (>30 minutes)
				require.Len(t, patterns.LongestAlerts, 2)

				// Should be sorted by duration, longest first
				assert.Equal(t, "long-alert", patterns.LongestAlerts[0].AlertID)
				assert.InDelta(t, 60, patterns.LongestAlerts[0].DurationMinutes, 5) // ~60 minutes
				assert.False(t, patterns.LongestAlerts[0].IsOngoing)

				assert.Equal(t, "ongoing-long", patterns.LongestAlerts[1].AlertID)
				assert.InDelta(t, 45, patterns.LongestAlerts[1].DurationMinutes, 5) // ~45 minutes
				assert.True(t, patterns.LongestAlerts[1].IsOngoing)
			},
		},
		{
			name: "muted alerts tracking",
			alerts: []*models.AnalyzeAlertsResponseAlert{
				{AlertID: "1", MonitorSlug: "monitor-1", IsMuted: true},
				{AlertID: "2", MonitorSlug: "monitor-2", IsMuted: false},
				{AlertID: "3", MonitorSlug: "monitor-1", IsMuted: true}, // duplicate slug
			},
			validate: func(t *testing.T, patterns AlertPatterns) {
				require.Len(t, patterns.RecentMutes, 1) // deduplicated
				assert.Contains(t, patterns.RecentMutes, "monitor-1")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tools.identifyPatterns(tt.alerts)
			tt.validate(t, result)
		})
	}
}

func TestGenerateAnalysisNotes(t *testing.T) {
	tools := &Tools{logger: zaptest.NewLogger(t)}

	tests := []struct {
		name       string
		summary    AlertSummary
		clustering ClusteringInsights
		patterns   AlertPatterns
		validate   func(t *testing.T, notes []string)
	}{
		{
			name: "high alert volume",
			summary: AlertSummary{
				TotalAlerts: 75, // > 50
			},
			validate: func(t *testing.T, notes []string) {
				assert.Contains(t, notes[0], "Large number of alerts detected")
			},
		},
		{
			name: "ongoing alerts warning",
			summary: AlertSummary{
				OngoingAlerts: 5,
			},
			validate: func(t *testing.T, notes []string) {
				found := false
				for _, note := range notes {
					if assert.Contains(t, note, "5 alerts are still ongoing") {
						found = true
						break
					}
				}
				assert.True(t, found, "Should contain ongoing alerts note")
			},
		},
		{
			name: "large cluster detection",
			clustering: ClusteringInsights{
				LargestGroupSize: 5, // > 3
			},
			validate: func(t *testing.T, notes []string) {
				found := false
				for _, note := range notes {
					if assert.Contains(t, note, "Large alert cluster detected") {
						found = true
						break
					}
				}
				assert.True(t, found, "Should contain clustering note")
			},
		},
		{
			name: "noisy monitor detection",
			patterns: AlertPatterns{
				MostFrequentMonitors: []MonitorFrequency{
					{MonitorSlug: "noisy-monitor", AlertCount: 10}, // > 5
				},
			},
			validate: func(t *testing.T, notes []string) {
				found := false
				for _, note := range notes {
					if assert.Contains(t, note, "Monitor 'noisy-monitor' is generating high alert volume") {
						found = true
						break
					}
				}
				assert.True(t, found, "Should contain noisy monitor note")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tools.generateAnalysisNotes(tt.summary, tt.clustering, tt.patterns)
			tt.validate(t, result)
		})
	}
}
