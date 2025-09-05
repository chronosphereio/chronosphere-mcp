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
	"sort"
	"strconv"
	"time"

	"github.com/go-openapi/swag"

	"github.com/chronosphereio/chronosphere-mcp/generated/dataunstable/models"
)

// analyzeAlertsResponse processes the API response and provides analysis based on format preference
func (t *Tools) analyzeAlertsResponse(response *models.DataunstableAnalyzeAlertsResponse, responseFormat string) *AlertAnalysisSummary {
	if response == nil || len(response.Alerts) == 0 {
		return &AlertAnalysisSummary{
			Summary:        AlertSummary{},
			ClusteringInfo: ClusteringInsights{},
			Patterns:       AlertPatterns{},
			Alerts:         []AlertDetail{},
			AnalysisNotes:  []string{"No alerts found matching the specified criteria"},
		}
	}

	alerts := response.Alerts

	summary := t.calculateSummary(alerts)

	// Choose analysis depth based on response format
	switch responseFormat {
	case "compact":
		return t.createCompactResponse(alerts, summary)
	case "detailed":
		return t.createDetailedResponse(alerts, summary)
	default: // "standard"
		return t.createStandardResponse(alerts, summary)
	}
}

// calculateSummary generates high-level statistics about the alert set
func (t *Tools) calculateSummary(alerts []*models.AnalyzeAlertsResponseAlert) AlertSummary {
	if len(alerts) == 0 {
		return AlertSummary{}
	}

	uniqueMonitors := make(map[string]bool)
	mutedCount := 0
	ongoingCount := 0
	totalDuration := float64(0)
	maxDuration := float64(0)
	alertsWithDuration := 0

	for _, alert := range alerts {
		// Track unique monitors
		if alert.MonitorSlug != "" {
			uniqueMonitors[alert.MonitorSlug] = true
		}

		// Count muted alerts
		if alert.IsMuted {
			mutedCount++
		}

		// Count ongoing alerts (no end time)
		if swag.IsZero(alert.EndTime) {
			ongoingCount++
		}

		// Calculate duration statistics
		if !swag.IsZero(alert.StartTime) {
			startTime, err := time.Parse(time.RFC3339, alert.StartTime.String())
			if err == nil {
				var duration time.Duration
				if !swag.IsZero(alert.EndTime) {
					endTime, err := time.Parse(time.RFC3339, alert.EndTime.String())
					if err == nil {
						duration = endTime.Sub(startTime)
					}
				} else {
					// Ongoing alert - duration from start to now
					duration = time.Since(startTime)
				}

				durationMinutes := duration.Minutes()
				if durationMinutes > 0 {
					totalDuration += durationMinutes
					alertsWithDuration++
					if durationMinutes > maxDuration {
						maxDuration = durationMinutes
					}
				}
			}
		}
	}

	avgDuration := float64(0)
	if alertsWithDuration > 0 {
		avgDuration = totalDuration / float64(alertsWithDuration)
	}

	return AlertSummary{
		TotalAlerts:        len(alerts),
		UniqueMonitors:     len(uniqueMonitors),
		MutedAlerts:        mutedCount,
		OngoingAlerts:      ongoingCount,
		AvgDurationMinutes: avgDuration,
		MaxDurationMinutes: maxDuration,
	}
}

// analyzeClusteringInfo provides insights into alert groupings and relationships
func (t *Tools) analyzeClusteringInfo(alerts []*models.AnalyzeAlertsResponseAlert) ClusteringInsights {
	signalGroups := make(map[string]*SignalGroup)
	monitorLabelGroups := make(map[string]*MonitorLabelGroup)

	for _, alert := range alerts {
		// Group by signals hash
		if alert.SignalsHash != "" {
			hashStr := alert.SignalsHash
			if group, exists := signalGroups[hashStr]; exists {
				group.AlertCount++
			} else {
				// Create new group with common signals
				commonSignals := make(map[string]string)
				for _, signal := range alert.Signal {
					if signal != nil {
						commonSignals[signal.Name] = signal.Value
					}
				}
				signalGroups[hashStr] = &SignalGroup{
					SignalsHash:   hashStr,
					AlertCount:    1,
					CommonSignals: commonSignals,
				}
			}
		}

		// Group by monitor labels hash (if available)
		if len(alert.MonitorLabels) > 0 {
			// Create a simple hash from monitor labels for grouping
			labelsHash := t.createLabelsHash(alert.MonitorLabels)
			if group, exists := monitorLabelGroups[labelsHash]; exists {
				group.AlertCount++
			} else {
				commonLabels := make(map[string]string)
				for _, label := range alert.MonitorLabels {
					if label != nil {
						commonLabels[label.Name] = label.Value
					}
				}
				monitorLabelGroups[labelsHash] = &MonitorLabelGroup{
					MonitorLabelsHash: labelsHash,
					AlertCount:        1,
					CommonLabels:      commonLabels,
				}
			}
		}
	}

	// Convert maps to slices and sort by alert count
	var signalGroupsList []SignalGroup
	largestGroupSize := 0
	for _, group := range signalGroups {
		signalGroupsList = append(signalGroupsList, *group)
		if group.AlertCount > largestGroupSize {
			largestGroupSize = group.AlertCount
		}
	}
	sort.Slice(signalGroupsList, func(i, j int) bool {
		return signalGroupsList[i].AlertCount > signalGroupsList[j].AlertCount
	})

	var monitorLabelGroupsList []MonitorLabelGroup
	for _, group := range monitorLabelGroups {
		monitorLabelGroupsList = append(monitorLabelGroupsList, *group)
		if group.AlertCount > largestGroupSize {
			largestGroupSize = group.AlertCount
		}
	}
	sort.Slice(monitorLabelGroupsList, func(i, j int) bool {
		return monitorLabelGroupsList[i].AlertCount > monitorLabelGroupsList[j].AlertCount
	})

	// Generate correlation notes
	var correlationNotes []string
	if len(signalGroupsList) > 0 && signalGroupsList[0].AlertCount > 1 {
		correlationNotes = append(correlationNotes,
			"Multiple alerts share identical signals, suggesting potential incident correlation")
	}
	if len(monitorLabelGroupsList) > 0 && monitorLabelGroupsList[0].AlertCount > 1 {
		correlationNotes = append(correlationNotes,
			"Multiple alerts from monitors with similar labels detected")
	}

	return ClusteringInsights{
		TotalAlertGroups:    len(signalGroups) + len(monitorLabelGroups),
		LargestGroupSize:    largestGroupSize,
		SimilarSignalGroups: signalGroupsList,
		MonitorLabelGroups:  monitorLabelGroupsList,
		CorrelationNotes:    correlationNotes,
	}
}

// identifyPatterns finds interesting patterns in the alert data
func (t *Tools) identifyPatterns(alerts []*models.AnalyzeAlertsResponseAlert) AlertPatterns {
	// Track monitor frequency
	monitorCounts := make(map[string]*MonitorFrequency)

	// Track long-running alerts
	var longAlerts []LongAlert

	// Track recent mutes and notification failures
	var recentMutes []string
	var notificationFailures []string

	for _, alert := range alerts {
		// Monitor frequency
		if alert.MonitorSlug != "" {
			if freq, exists := monitorCounts[alert.MonitorSlug]; exists {
				freq.AlertCount++
			} else {
				monitorCounts[alert.MonitorSlug] = &MonitorFrequency{
					MonitorSlug: alert.MonitorSlug,
					MonitorName: alert.MonitorName,
					AlertCount:  1,
				}
			}
		}

		// Long-running alerts (>30 minutes)
		if !swag.IsZero(alert.StartTime) {
			startTime, err := time.Parse(time.RFC3339, alert.StartTime.String())
			if err == nil {
				var duration time.Duration
				isOngoing := false

				if !swag.IsZero(alert.EndTime) {
					endTime, err := time.Parse(time.RFC3339, alert.EndTime.String())
					if err == nil {
						duration = endTime.Sub(startTime)
					}
				} else {
					duration = time.Since(startTime)
					isOngoing = true
				}

				durationMinutes := duration.Minutes()
				if durationMinutes > 30 { // Alerts longer than 30 minutes
					longAlerts = append(longAlerts, LongAlert{
						MonitorSlug:     alert.MonitorSlug,
						AlertID:         alert.AlertID,
						DurationMinutes: durationMinutes,
						IsOngoing:       isOngoing,
					})
				}
			}
		}

		// Recent mutes
		if alert.IsMuted {
			recentMutes = append(recentMutes, alert.MonitorSlug)
		}

		// Note: Notification failure tracking would require additional fields in the API response
		// For now, we'll skip this analysis
	}

	// Convert monitor counts to sorted slice
	var monitorFreqs []MonitorFrequency
	for _, freq := range monitorCounts {
		monitorFreqs = append(monitorFreqs, *freq)
	}
	sort.Slice(monitorFreqs, func(i, j int) bool {
		return monitorFreqs[i].AlertCount > monitorFreqs[j].AlertCount
	})

	// Sort long alerts by duration
	sort.Slice(longAlerts, func(i, j int) bool {
		return longAlerts[i].DurationMinutes > longAlerts[j].DurationMinutes
	})

	// Deduplicate and limit arrays
	recentMutes = t.deduplicateStrings(recentMutes)
	notificationFailures = t.deduplicateStrings(notificationFailures)

	return AlertPatterns{
		MostFrequentMonitors: monitorFreqs,
		LongestAlerts:        longAlerts,
		RecentMutes:          recentMutes,
		NotificationFailures: notificationFailures,
	}
}

// convertToAlertDetails converts API response alerts to display format
func (t *Tools) convertToAlertDetails(alerts []*models.AnalyzeAlertsResponseAlert) []AlertDetail {
	var details []AlertDetail

	for _, alert := range alerts {
		detail := AlertDetail{
			AlertID:                alert.AlertID,
			MonitorSlug:            alert.MonitorSlug,
			MonitorName:            alert.MonitorName,
			IsMuted:                alert.IsMuted,
			Signal:                 alert.Signal,
			MonitorLabels:          alert.MonitorLabels,
			SignalMatchCount:       alert.SignalMatchCount,
			MonitorLabelMatchCount: alert.MonitorLabelMatchCount,
		}

		// Format timestamps
		if !swag.IsZero(alert.StartTime) {
			detail.StartTime = alert.StartTime.String()
		}
		if !swag.IsZero(alert.EndTime) {
			detail.EndTime = alert.EndTime.String()
		}

		// Calculate duration
		if !swag.IsZero(alert.StartTime) {
			startTime, err := time.Parse(time.RFC3339, alert.StartTime.String())
			if err == nil {
				var duration time.Duration
				if !swag.IsZero(alert.EndTime) {
					endTime, err := time.Parse(time.RFC3339, alert.EndTime.String())
					if err == nil {
						duration = endTime.Sub(startTime)
					}
				} else {
					duration = time.Since(startTime)
				}
				detail.DurationMinutes = duration.Minutes()
			}
		}

		// Get current severity (most recent from severity history)
		if len(alert.SeverityHistory) > 0 {
			detail.CurrentSeverity = alert.SeverityHistory[len(alert.SeverityHistory)-1].Severity
		}

		// Notification policy info
		if alert.NotificationPolicySlug != "" {
			detail.NotificationPolicy = alert.NotificationPolicySlug
		}

		details = append(details, detail)
	}

	return details
}

// generateAnalysisNotes creates helpful context for LLMs consuming the output
func (t *Tools) generateAnalysisNotes(summary AlertSummary, clustering ClusteringInsights, patterns AlertPatterns) []string {
	var notes []string

	// Summary insights
	if summary.TotalAlerts > 50 {
		notes = append(notes, "Large number of alerts detected - consider filtering for specific analysis")
	}
	if summary.OngoingAlerts > 0 {
		notes = append(notes, strconv.Itoa(summary.OngoingAlerts)+" alerts are still ongoing")
	}
	if summary.AvgDurationMinutes > 60 {
		notes = append(notes, "Alerts have unusually long average duration - may indicate monitoring tuning needed")
	}

	// Clustering insights
	if clustering.LargestGroupSize > 3 {
		notes = append(notes, "Large alert cluster detected - possible cascading failure or related incident")
	}

	// Pattern insights
	if len(patterns.MostFrequentMonitors) > 0 && patterns.MostFrequentMonitors[0].AlertCount > 5 {
		notes = append(notes, "Monitor '"+patterns.MostFrequentMonitors[0].MonitorSlug+"' is generating high alert volume")
	}
	if len(patterns.LongestAlerts) > 0 && patterns.LongestAlerts[0].DurationMinutes > 120 {
		notes = append(notes, "Extremely long-running alert detected - review monitor thresholds")
	}
	if len(patterns.NotificationFailures) > 0 {
		notes = append(notes, "Notification delivery failures detected - check notification configuration")
	}

	return notes
}

// Helper functions

// createLabelsHash creates a simple hash string from labels for grouping
func (t *Tools) createLabelsHash(labels []*models.DataunstableAnalyzeAlertsResponseLabel) string {
	if len(labels) == 0 {
		return ""
	}

	// Simple concatenation for grouping purposes
	var parts []string
	for _, label := range labels {
		if label != nil {
			parts = append(parts, label.Name+":"+label.Value)
		}
	}
	sort.Strings(parts) // Ensure consistent ordering

	result := ""
	for i, part := range parts {
		if i > 0 {
			result += ","
		}
		result += part
	}
	return result
}

// deduplicateStrings removes duplicates from a string slice
func (t *Tools) deduplicateStrings(input []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, str := range input {
		if !seen[str] {
			seen[str] = true
			result = append(result, str)
		}
	}

	return result
}

// createCompactResponse creates a minimal, token-efficient response with just key insights
func (t *Tools) createCompactResponse(alerts []*models.AnalyzeAlertsResponseAlert, summary AlertSummary) *AlertAnalysisSummary {
	// Only get top monitor patterns and minimal clustering info
	patterns := t.identifyPatterns(alerts)

	// Limit to top 5 monitors to reduce tokens
	if len(patterns.MostFrequentMonitors) > 5 {
		patterns.MostFrequentMonitors = patterns.MostFrequentMonitors[:5]
	}

	// Limit to top 3 long alerts
	if len(patterns.LongestAlerts) > 3 {
		patterns.LongestAlerts = patterns.LongestAlerts[:3]
	}

	// Simplified clustering - just counts, no detailed groups
	clustering := ClusteringInsights{
		TotalAlertGroups: t.countDistinctGroups(alerts),
		LargestGroupSize: t.findLargestGroupSize(alerts),
		// Skip detailed groups to save tokens
		SimilarSignalGroups: nil,
		MonitorLabelGroups:  nil,
		CorrelationNotes:    t.generateCompactCorrelationNotes(alerts),
	}

	// No full alert details - just essential IDs and monitor slugs
	compactAlerts := t.createCompactAlertList(alerts)

	// Concise analysis notes
	analysisNotes := t.generateCompactAnalysisNotes(summary, clustering, patterns)

	return &AlertAnalysisSummary{
		Summary:        summary,
		ClusteringInfo: clustering,
		Patterns:       patterns,
		Alerts:         compactAlerts,
		AnalysisNotes:  analysisNotes,
	}
}

// createStandardResponse creates the current default response
func (t *Tools) createStandardResponse(alerts []*models.AnalyzeAlertsResponseAlert, summary AlertSummary) *AlertAnalysisSummary {
	clustering := t.analyzeClusteringInfo(alerts)
	patterns := t.identifyPatterns(alerts)

	// Limit detailed groups to top 10 for standard mode
	if len(clustering.SimilarSignalGroups) > 10 {
		clustering.SimilarSignalGroups = clustering.SimilarSignalGroups[:10]
	}
	if len(clustering.MonitorLabelGroups) > 10 {
		clustering.MonitorLabelGroups = clustering.MonitorLabelGroups[:10]
	}

	// Create simplified alert details (no full signal/label arrays)
	alertDetails := t.createStandardAlertDetails(alerts)
	analysisNotes := t.generateAnalysisNotes(summary, clustering, patterns)

	return &AlertAnalysisSummary{
		Summary:        summary,
		ClusteringInfo: clustering,
		Patterns:       patterns,
		Alerts:         alertDetails,
		AnalysisNotes:  analysisNotes,
	}
}

// createDetailedResponse creates the full comprehensive response (current behavior)
func (t *Tools) createDetailedResponse(alerts []*models.AnalyzeAlertsResponseAlert, summary AlertSummary) *AlertAnalysisSummary {
	clustering := t.analyzeClusteringInfo(alerts)
	patterns := t.identifyPatterns(alerts)
	alertDetails := t.convertToAlertDetails(alerts)
	analysisNotes := t.generateAnalysisNotes(summary, clustering, patterns)

	return &AlertAnalysisSummary{
		Summary:        summary,
		ClusteringInfo: clustering,
		Patterns:       patterns,
		Alerts:         alertDetails,
		AnalysisNotes:  analysisNotes,
	}
}

// Helper functions for compact response

// countDistinctGroups counts unique signal and monitor label groups without creating full objects
func (t *Tools) countDistinctGroups(alerts []*models.AnalyzeAlertsResponseAlert) int {
	signalHashes := make(map[string]bool)
	labelHashes := make(map[string]bool)

	for _, alert := range alerts {
		if alert.SignalsHash != "" {
			signalHashes[alert.SignalsHash] = true
		}
		if len(alert.MonitorLabels) > 0 {
			labelHashes[t.createLabelsHash(alert.MonitorLabels)] = true
		}
	}

	return len(signalHashes) + len(labelHashes)
}

// findLargestGroupSize finds the largest cluster size without creating full group objects
func (t *Tools) findLargestGroupSize(alerts []*models.AnalyzeAlertsResponseAlert) int {
	signalCounts := make(map[string]int)
	labelCounts := make(map[string]int)

	for _, alert := range alerts {
		if alert.SignalsHash != "" {
			signalCounts[alert.SignalsHash]++
		}
		if len(alert.MonitorLabels) > 0 {
			labelCounts[t.createLabelsHash(alert.MonitorLabels)]++
		}
	}

	maxSize := 0
	for _, count := range signalCounts {
		if count > maxSize {
			maxSize = count
		}
	}
	for _, count := range labelCounts {
		if count > maxSize {
			maxSize = count
		}
	}

	return maxSize
}

// generateCompactCorrelationNotes creates minimal correlation insights
func (t *Tools) generateCompactCorrelationNotes(alerts []*models.AnalyzeAlertsResponseAlert) []string {
	largestGroupSize := t.findLargestGroupSize(alerts)

	var notes []string
	if largestGroupSize > 3 {
		notes = append(notes, "Large cluster detected - possible incident correlation")
	}

	return notes
}

// createCompactAlertList creates a minimal alert list with just IDs and monitor slugs
func (t *Tools) createCompactAlertList(alerts []*models.AnalyzeAlertsResponseAlert) []AlertDetail {
	var details []AlertDetail

	for _, alert := range alerts {
		detail := AlertDetail{
			AlertID:     alert.AlertID,
			MonitorSlug: alert.MonitorSlug,
			IsMuted:     alert.IsMuted,
		}

		// Calculate just duration, no full timestamps
		if !swag.IsZero(alert.StartTime) {
			startTime, err := time.Parse(time.RFC3339, alert.StartTime.String())
			if err == nil {
				var duration time.Duration
				if !swag.IsZero(alert.EndTime) {
					endTime, err := time.Parse(time.RFC3339, alert.EndTime.String())
					if err == nil {
						duration = endTime.Sub(startTime)
					}
				} else {
					duration = time.Since(startTime)
				}
				detail.DurationMinutes = duration.Minutes()
			}
		}

		details = append(details, detail)
	}

	return details
}

// createStandardAlertDetails creates alert details without full signal/label arrays
func (t *Tools) createStandardAlertDetails(alerts []*models.AnalyzeAlertsResponseAlert) []AlertDetail {
	var details []AlertDetail

	for _, alert := range alerts {
		detail := AlertDetail{
			AlertID:     alert.AlertID,
			MonitorSlug: alert.MonitorSlug,
			MonitorName: alert.MonitorName,
			IsMuted:     alert.IsMuted,
		}

		// Include timestamps but skip full signal/label arrays
		if !swag.IsZero(alert.StartTime) {
			detail.StartTime = alert.StartTime.String()
		}
		if !swag.IsZero(alert.EndTime) {
			detail.EndTime = alert.EndTime.String()
		}

		// Calculate duration
		if !swag.IsZero(alert.StartTime) {
			startTime, err := time.Parse(time.RFC3339, alert.StartTime.String())
			if err == nil {
				var duration time.Duration
				if !swag.IsZero(alert.EndTime) {
					endTime, err := time.Parse(time.RFC3339, alert.EndTime.String())
					if err == nil {
						duration = endTime.Sub(startTime)
					}
				} else {
					duration = time.Since(startTime)
				}
				detail.DurationMinutes = duration.Minutes()
			}
		}

		// Get current severity
		if len(alert.SeverityHistory) > 0 {
			detail.CurrentSeverity = alert.SeverityHistory[len(alert.SeverityHistory)-1].Severity
		}

		// Include notification policy but not full arrays
		if alert.NotificationPolicySlug != "" {
			detail.NotificationPolicy = alert.NotificationPolicySlug
		}

		details = append(details, detail)
	}

	return details
}

// generateCompactAnalysisNotes creates concise analysis insights
func (t *Tools) generateCompactAnalysisNotes(summary AlertSummary, clustering ClusteringInsights, patterns AlertPatterns) []string {
	var notes []string

	// Key insights only
	if summary.OngoingAlerts > 0 {
		notes = append(notes, fmt.Sprintf("%d ongoing alerts", summary.OngoingAlerts))
	}

	if clustering.LargestGroupSize > 3 {
		notes = append(notes, "Large alert cluster detected")
	}

	if len(patterns.MostFrequentMonitors) > 0 && patterns.MostFrequentMonitors[0].AlertCount > 5 {
		notes = append(notes, fmt.Sprintf("High volume from %s (%d alerts)",
			patterns.MostFrequentMonitors[0].MonitorSlug, patterns.MostFrequentMonitors[0].AlertCount))
	}

	if len(patterns.LongestAlerts) > 0 && patterns.LongestAlerts[0].DurationMinutes > 120 {
		notes = append(notes, "Long-running alerts detected")
	}

	return notes
}
