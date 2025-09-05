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
	"github.com/stretchr/testify/require"

	"github.com/chronosphereio/chronosphere-mcp/generated/dataunstable/models"
)

func TestBuildStringMatcher(t *testing.T) {
	tests := []struct {
		name           string
		operation      string
		singleValue    string
		multiValues    []string
		expectError    bool
		expectedOp     models.AnalyzeAlertsRequestStringMatcherOperation
		expectedSingle string
		expectedMulti  []string
	}{
		{
			name:           "equal operation",
			operation:      "EQUAL",
			singleValue:    "test-monitor",
			expectedOp:     models.AnalyzeAlertsRequestStringMatcherOperationEQUAL,
			expectedSingle: "test-monitor",
		},
		{
			name:           "contains operation",
			operation:      "CONTAINS",
			singleValue:    "database",
			expectedOp:     models.AnalyzeAlertsRequestStringMatcherOperationCONTAINS,
			expectedSingle: "database",
		},
		{
			name:           "regex operation",
			operation:      "REGEX_EQUAL",
			singleValue:    "api-.*",
			expectedOp:     models.AnalyzeAlertsRequestStringMatcherOperationREGEXEQUAL,
			expectedSingle: "api-.*",
		},
		{
			name:          "any_of operation",
			operation:     "ANY_OF",
			multiValues:   []string{"api", "web", "database"},
			expectedOp:    models.AnalyzeAlertsRequestStringMatcherOperationANYOF,
			expectedMulti: []string{"api", "web", "database"},
		},
		{
			name:          "none_of operation",
			operation:     "NONE_OF",
			multiValues:   []string{"test", "staging"},
			expectedOp:    models.AnalyzeAlertsRequestStringMatcherOperationNONEOF,
			expectedMulti: []string{"test", "staging"},
		},
		{
			name:           "case insensitive operation",
			operation:      "equal", // lowercase
			singleValue:    "test",
			expectedOp:     models.AnalyzeAlertsRequestStringMatcherOperationEQUAL,
			expectedSingle: "test",
		},
		{
			name:        "invalid operation",
			operation:   "INVALID_OP",
			singleValue: "test",
			expectError: true,
		},
		{
			name:        "empty operation",
			operation:   "",
			singleValue: "test",
			expectError: true,
		},
		{
			name:        "any_of without multi values",
			operation:   "ANY_OF",
			singleValue: "test",
			expectError: true,
		},
		{
			name:        "equal without single value",
			operation:   "EQUAL",
			multiValues: []string{"test"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := buildStringMatcher(tt.operation, tt.singleValue, tt.multiValues)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, result)

			assert.Equal(t, tt.expectedOp, result.Operation)

			if tt.expectedSingle != "" {
				assert.Equal(t, tt.expectedSingle, result.SingleValue)
				assert.Empty(t, result.MultiValue)
			}

			if len(tt.expectedMulti) > 0 {
				assert.Equal(t, tt.expectedMulti, result.MultiValue)
				assert.Empty(t, result.SingleValue)
			}
		})
	}
}

func TestBuildStringMatcherOperationCoverage(t *testing.T) {
	// Ensure all operations are properly mapped
	operations := []struct {
		input    string
		expected models.AnalyzeAlertsRequestStringMatcherOperation
	}{
		{"EQUAL", models.AnalyzeAlertsRequestStringMatcherOperationEQUAL},
		{"REGEX_EQUAL", models.AnalyzeAlertsRequestStringMatcherOperationREGEXEQUAL},
		{"CONTAINS", models.AnalyzeAlertsRequestStringMatcherOperationCONTAINS},
		{"NOT_EQUAL", models.AnalyzeAlertsRequestStringMatcherOperationNOTEQUAL},
		{"NOT_REGEX_EQUAL", models.AnalyzeAlertsRequestStringMatcherOperationNOTREGEXEQUAL},
		{"NOT_CONTAINS", models.AnalyzeAlertsRequestStringMatcherOperationNOTCONTAINS},
		{"ANY_OF", models.AnalyzeAlertsRequestStringMatcherOperationANYOF},
		{"NONE_OF", models.AnalyzeAlertsRequestStringMatcherOperationNONEOF},
	}

	for _, op := range operations {
		t.Run(op.input, func(t *testing.T) {
			var result *models.AnalyzeAlertsRequestStringMatcher
			var err error

			if op.expected == models.AnalyzeAlertsRequestStringMatcherOperationANYOF ||
				op.expected == models.AnalyzeAlertsRequestStringMatcherOperationNONEOF {
				result, err = buildStringMatcher(op.input, "", []string{"test"})
			} else {
				result, err = buildStringMatcher(op.input, "test", nil)
			}

			require.NoError(t, err)
			assert.Equal(t, op.expected, result.Operation)
		})
	}
}
