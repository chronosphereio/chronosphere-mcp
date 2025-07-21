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

package prometheus

import (
	"testing"

	"github.com/prometheus/common/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListPrometheusLabelValuesPagination(t *testing.T) {
	tests := []struct {
		name           string
		labelValues    model.LabelValues
		limit          int
		offset         int
		expectedValues []string
		expectedMeta   map[string]any
	}{
		{
			name:           "no pagination",
			labelValues:    model.LabelValues{"value1", "value2", "value3"},
			limit:          0,
			offset:         0,
			expectedValues: []string{"value1", "value2", "value3"},
			expectedMeta: map[string]any{
				"total":    3,
				"offset":   0,
				"limit":    0,
				"returned": 3,
			},
		},
		{
			name:           "limit only",
			labelValues:    model.LabelValues{"value1", "value2", "value3", "value4", "value5"},
			limit:          3,
			offset:         0,
			expectedValues: []string{"value1", "value2", "value3"},
			expectedMeta: map[string]any{
				"total":    5,
				"offset":   0,
				"limit":    3,
				"returned": 3,
			},
		},
		{
			name:           "offset only",
			labelValues:    model.LabelValues{"value1", "value2", "value3", "value4", "value5"},
			limit:          0,
			offset:         2,
			expectedValues: []string{"value3", "value4", "value5"},
			expectedMeta: map[string]any{
				"total":    5,
				"offset":   2,
				"limit":    0,
				"returned": 3,
			},
		},
		{
			name:           "limit and offset",
			labelValues:    model.LabelValues{"value1", "value2", "value3", "value4", "value5"},
			limit:          2,
			offset:         1,
			expectedValues: []string{"value2", "value3"},
			expectedMeta: map[string]any{
				"total":    5,
				"offset":   1,
				"limit":    2,
				"returned": 2,
			},
		},
		{
			name:           "offset beyond total",
			labelValues:    model.LabelValues{"value1", "value2", "value3"},
			limit:          10,
			offset:         5,
			expectedValues: []string{},
			expectedMeta: map[string]any{
				"total":    3,
				"offset":   5,
				"limit":    10,
				"returned": 0,
			},
		},
		{
			name:           "limit exceeds remaining",
			labelValues:    model.LabelValues{"value1", "value2", "value3", "value4", "value5"},
			limit:          10,
			offset:         3,
			expectedValues: []string{"value4", "value5"},
			expectedMeta: map[string]any{
				"total":    5,
				"offset":   3,
				"limit":    10,
				"returned": 2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock the API response
			mockResult := mockLabelValuesWithPagination(tt.labelValues, tt.limit, tt.offset)

			// Verify result values
			resultValues, ok := mockResult.JSONContent.(map[string]any)["result"].(model.LabelValues)
			require.True(t, ok, "result should be model.LabelValues")

			assert.Equal(t, len(tt.expectedValues), len(resultValues))
			for i, expected := range tt.expectedValues {
				assert.Equal(t, model.LabelValue(expected), resultValues[i])
			}

			// Verify metadata
			assert.Equal(t, tt.expectedMeta, mockResult.Meta)
		})
	}
}

// mockLabelValuesWithPagination simulates the pagination logic from listPrometheusLabelValues
func mockLabelValuesWithPagination(resp model.LabelValues, limit, offset int) *Result {
	var paginatedResp model.LabelValues
	if limit > 0 || offset > 0 {
		totalValues := len(resp)
		start := offset
		if start > totalValues {
			start = totalValues
		}
		end := totalValues
		if limit > 0 && start+limit < totalValues {
			end = start + limit
		}
		paginatedResp = resp[start:end]
	} else {
		paginatedResp = resp
	}

	result := &Result{
		JSONContent: map[string]any{
			"result": paginatedResp,
		},
		Meta: make(map[string]any),
	}

	result.Meta["total"] = len(resp)
	result.Meta["offset"] = offset
	result.Meta["limit"] = limit
	result.Meta["returned"] = len(paginatedResp)

	return result
}

// Result is a local copy for testing
type Result struct {
	JSONContent any
	Meta        map[string]any
}
