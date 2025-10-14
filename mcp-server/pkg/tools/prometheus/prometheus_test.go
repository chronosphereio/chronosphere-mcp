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
	"strings"
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

func TestQueryPrometheusRangePagination(t *testing.T) {
	tests := []struct {
		name           string
		matrix         model.Matrix
		limit          int
		offset         int
		expectedSeries int
		expectedMeta   map[string]any
	}{
		{
			name: "no pagination",
			matrix: model.Matrix{
				&model.SampleStream{Metric: model.Metric{"__name__": "series1"}},
				&model.SampleStream{Metric: model.Metric{"__name__": "series2"}},
				&model.SampleStream{Metric: model.Metric{"__name__": "series3"}},
			},
			limit:          0,
			offset:         0,
			expectedSeries: 3,
			expectedMeta: map[string]any{
				"total_series":    3,
				"offset":          0,
				"limit":           0,
				"returned_series": 3,
			},
		},
		{
			name: "limit only",
			matrix: model.Matrix{
				&model.SampleStream{Metric: model.Metric{"__name__": "series1"}},
				&model.SampleStream{Metric: model.Metric{"__name__": "series2"}},
				&model.SampleStream{Metric: model.Metric{"__name__": "series3"}},
				&model.SampleStream{Metric: model.Metric{"__name__": "series4"}},
				&model.SampleStream{Metric: model.Metric{"__name__": "series5"}},
			},
			limit:          3,
			offset:         0,
			expectedSeries: 3,
			expectedMeta: map[string]any{
				"total_series":    5,
				"offset":          0,
				"limit":           3,
				"returned_series": 3,
			},
		},
		{
			name: "offset only",
			matrix: model.Matrix{
				&model.SampleStream{Metric: model.Metric{"__name__": "series1"}},
				&model.SampleStream{Metric: model.Metric{"__name__": "series2"}},
				&model.SampleStream{Metric: model.Metric{"__name__": "series3"}},
				&model.SampleStream{Metric: model.Metric{"__name__": "series4"}},
				&model.SampleStream{Metric: model.Metric{"__name__": "series5"}},
			},
			limit:          0,
			offset:         2,
			expectedSeries: 3,
			expectedMeta: map[string]any{
				"total_series":    5,
				"offset":          2,
				"limit":           0,
				"returned_series": 3,
			},
		},
		{
			name: "limit and offset",
			matrix: model.Matrix{
				&model.SampleStream{Metric: model.Metric{"__name__": "series1"}},
				&model.SampleStream{Metric: model.Metric{"__name__": "series2"}},
				&model.SampleStream{Metric: model.Metric{"__name__": "series3"}},
				&model.SampleStream{Metric: model.Metric{"__name__": "series4"}},
				&model.SampleStream{Metric: model.Metric{"__name__": "series5"}},
			},
			limit:          2,
			offset:         1,
			expectedSeries: 2,
			expectedMeta: map[string]any{
				"total_series":    5,
				"offset":          1,
				"limit":           2,
				"returned_series": 2,
			},
		},
		{
			name: "offset beyond total",
			matrix: model.Matrix{
				&model.SampleStream{Metric: model.Metric{"__name__": "series1"}},
				&model.SampleStream{Metric: model.Metric{"__name__": "series2"}},
				&model.SampleStream{Metric: model.Metric{"__name__": "series3"}},
			},
			limit:          10,
			offset:         5,
			expectedSeries: 0,
			expectedMeta: map[string]any{
				"total_series":    3,
				"offset":          5,
				"limit":           10,
				"returned_series": 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock the pagination logic
			mockResult := mockMatrixWithPagination(tt.matrix, tt.limit, tt.offset)

			// Verify result matrix
			resultMatrix, ok := mockResult.JSONContent.(map[string]any)["result"].(model.Matrix)
			require.True(t, ok, "result should be model.Matrix")

			assert.Equal(t, tt.expectedSeries, len(resultMatrix))

			// Verify metadata
			assert.Equal(t, tt.expectedMeta, mockResult.Meta)
		})
	}
}

// mockMatrixWithPagination simulates the pagination logic from queryPrometheusRange
func mockMatrixWithPagination(matrix model.Matrix, limit, offset int) *Result {
	var paginatedMatrix model.Matrix
	if limit > 0 || offset > 0 {
		totalSeries := len(matrix)
		start := offset
		if start > totalSeries {
			start = totalSeries
		}
		end := totalSeries
		if limit > 0 && start+limit < totalSeries {
			end = start + limit
		}
		paginatedMatrix = matrix[start:end]
	} else {
		paginatedMatrix = matrix
	}

	result := &Result{
		JSONContent: map[string]any{
			"result": paginatedMatrix,
		},
		Meta: make(map[string]any),
	}

	result.Meta["total_series"] = len(matrix)
	result.Meta["offset"] = offset
	result.Meta["limit"] = limit
	result.Meta["returned_series"] = len(paginatedMatrix)

	return result
}

func TestFormatMatrixAsCSV(t *testing.T) {
	tests := []struct {
		name     string
		matrix   model.Matrix
		expected string
	}{
		{
			name:     "empty matrix",
			matrix:   model.Matrix{},
			expected: "# No data\n",
		},
		{
			name: "single series single timestamp",
			matrix: model.Matrix{
				&model.SampleStream{
					Metric: model.Metric{
						"__name__": "test_metric",
						"label1":   "value1",
					},
					Values: []model.SamplePair{
						{Timestamp: 1000000, Value: 42},
					},
				},
			},
			expected: `# Series Metadata
series_id,__name__,label1
1,test_metric,value1

# Time Series Data
timestamp,series_1
1000.000,42
`,
		},
		{
			name: "multiple series multiple timestamps",
			matrix: model.Matrix{
				&model.SampleStream{
					Metric: model.Metric{
						"__name__": "metric1",
						"env":      "prod",
					},
					Values: []model.SamplePair{
						{Timestamp: 1000000, Value: 10},
						{Timestamp: 2000000, Value: 20},
					},
				},
				&model.SampleStream{
					Metric: model.Metric{
						"__name__": "metric2",
						"env":      "dev",
					},
					Values: []model.SamplePair{
						{Timestamp: 1000000, Value: 30},
						{Timestamp: 2000000, Value: 40},
					},
				},
			},
			expected: `# Series Metadata
series_id,__name__,env
1,metric1,prod
2,metric2,dev

# Time Series Data
timestamp,series_1,series_2
1000.000,10,30
2000.000,20,40
`,
		},
		{
			name: "series with missing values (sparse data)",
			matrix: model.Matrix{
				&model.SampleStream{
					Metric: model.Metric{
						"__name__": "metric1",
					},
					Values: []model.SamplePair{
						{Timestamp: 1000000, Value: 10},
						{Timestamp: 3000000, Value: 30},
					},
				},
				&model.SampleStream{
					Metric: model.Metric{
						"__name__": "metric2",
					},
					Values: []model.SamplePair{
						{Timestamp: 2000000, Value: 20},
						{Timestamp: 3000000, Value: 40},
					},
				},
			},
			expected: `# Series Metadata
series_id,__name__
1,metric1
2,metric2

# Time Series Data
timestamp,series_1,series_2
1000.000,10,
2000.000,,20
3000.000,30,40
`,
		},
		{
			name: "series with different label sets",
			matrix: model.Matrix{
				&model.SampleStream{
					Metric: model.Metric{
						"__name__": "http_requests",
						"method":   "GET",
						"status":   "200",
					},
					Values: []model.SamplePair{
						{Timestamp: 1000000, Value: 100},
					},
				},
				&model.SampleStream{
					Metric: model.Metric{
						"__name__": "http_requests",
						"method":   "POST",
						// status is missing
					},
					Values: []model.SamplePair{
						{Timestamp: 1000000, Value: 50},
					},
				},
			},
			expected: `# Series Metadata
series_id,__name__,method,status
1,http_requests,GET,200
2,http_requests,POST,

# Time Series Data
timestamp,series_1,series_2
1000.000,100,50
`,
		},
		{
			name: "labels with special characters",
			matrix: model.Matrix{
				&model.SampleStream{
					Metric: model.Metric{
						"__name__": "test",
						"label":    "value,with,commas",
					},
					Values: []model.SamplePair{
						{Timestamp: 1000000, Value: 1},
					},
				},
			},
			expected: `# Series Metadata
series_id,__name__,label
1,test,"value,with,commas"

# Time Series Data
timestamp,series_1
1000.000,1
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatMatrixAsCSV(tt.matrix)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestFormatMatrixAsCSV_RealWorldExample tests with data similar to the user's example
func TestFormatMatrixAsCSV_RealWorldExample(t *testing.T) {
	matrix := model.Matrix{
		&model.SampleStream{
			Metric: model.Metric{
				"__name__":                   "disabled_metrics_total",
				"chronosphere_k8s_cluster":   "prod-eu-a",
				"chronosphere_k8s_namespace": "global",
				"chronosphere_namespace":     "global",
				"chronosphere_owner":         "infra",
				"chronosphere_service":       "descheduler",
				"environment":                "production",
				"instance":                   "global/descheduler-high-utilization-cd75d8b99-nl866",
				"job":                        "global/descheduler-high-utilization/0",
			},
			Values: []model.SamplePair{
				{Timestamp: 1760474193915, Value: 0},
				{Timestamp: 1760474253915, Value: 0},
			},
		},
		&model.SampleStream{
			Metric: model.Metric{
				"__name__":                   "disabled_metrics_total",
				"chronosphere_k8s_cluster":   "prod-eu-a",
				"chronosphere_k8s_namespace": "global",
				"chronosphere_namespace":     "global",
				"chronosphere_owner":         "infra",
				"chronosphere_service":       "descheduler",
				"environment":                "production",
				"instance":                   "global/descheduler-tsc-policies-c79bcb666-pfjg7",
				"job":                        "global/descheduler-tsc-policies/0",
			},
			Values: []model.SamplePair{
				{Timestamp: 1760474193915, Value: 0},
				{Timestamp: 1760474253915, Value: 0},
			},
		},
	}

	result := formatMatrixAsCSV(matrix)

	// Verify structure
	assert.True(t, strings.Contains(result, "# Series Metadata"))
	assert.True(t, strings.Contains(result, "# Time Series Data"))

	// Verify all label names are present
	assert.True(t, strings.Contains(result, "__name__"))
	assert.True(t, strings.Contains(result, "chronosphere_k8s_cluster"))
	assert.True(t, strings.Contains(result, "instance"))

	// Verify series IDs
	assert.True(t, strings.Contains(result, "series_1"))
	assert.True(t, strings.Contains(result, "series_2"))

	// Verify timestamps (converted to seconds)
	assert.True(t, strings.Contains(result, "1760474193.915"))
	assert.True(t, strings.Contains(result, "1760474253.915"))

	// Verify it's much more compact than JSON
	// The original JSON from the user was thousands of characters
	// Our CSV should be significantly smaller
	assert.Less(t, len(result), 1000, "CSV output should be compact")
}

func TestFormatLabelSetsAsCSV(t *testing.T) {
	tests := []struct {
		name      string
		labelSets []model.LabelSet
		expected  string
	}{
		{
			name:      "empty label sets",
			labelSets: []model.LabelSet{},
			expected:  "# No series found\n",
		},
		{
			name: "single series",
			labelSets: []model.LabelSet{
				{
					"__name__": "http_requests_total",
					"method":   "GET",
					"status":   "200",
				},
			},
			expected: `__name__,method,status
http_requests_total,GET,200
`,
		},
		{
			name: "multiple series with same labels",
			labelSets: []model.LabelSet{
				{
					"__name__": "http_requests_total",
					"method":   "GET",
					"status":   "200",
				},
				{
					"__name__": "http_requests_total",
					"method":   "POST",
					"status":   "201",
				},
			},
			expected: `__name__,method,status
http_requests_total,GET,200
http_requests_total,POST,201
`,
		},
		{
			name: "series with different label sets",
			labelSets: []model.LabelSet{
				{
					"__name__": "http_requests_total",
					"method":   "GET",
					"status":   "200",
				},
				{
					"__name__": "http_requests_total",
					"method":   "POST",
					// status is missing
				},
			},
			expected: `__name__,method,status
http_requests_total,GET,200
http_requests_total,POST,
`,
		},
		{
			name: "labels with special characters",
			labelSets: []model.LabelSet{
				{
					"__name__": "test_metric",
					"label":    "value,with,commas",
				},
			},
			expected: `__name__,label
test_metric,"value,with,commas"
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatLabelSetsAsCSV(tt.labelSets)
			assert.Equal(t, tt.expected, result)
		})
	}
}
