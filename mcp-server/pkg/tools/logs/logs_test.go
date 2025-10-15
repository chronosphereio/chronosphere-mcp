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

package logs

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/chronosphereio/chronosphere-mcp/generated/datav1/models"
)

func TestTrimLogEntries(t *testing.T) {
	tests := []struct {
		name         string
		payload      *models.Datav1QueryLogsRangeResponse
		limit        int
		offset       int
		wantTrimmed  int
		validateFunc func(t *testing.T, result *models.Datav1QueryLogsRangeResponse)
	}{
		{
			name:        "nil payload",
			payload:     nil,
			limit:       10,
			offset:      0,
			wantTrimmed: 0,
			validateFunc: func(t *testing.T, result *models.Datav1QueryLogsRangeResponse) {
				assert.Nil(t, result)
			},
		},
		{
			name: "zero limit and offset returns original",
			payload: &models.Datav1QueryLogsRangeResponse{
				GridData: &models.QueryLogsRangeResponseGridData{
					Rows: []*models.QueryLogsRangeResponseRow{{}, {}, {}},
				},
			},
			limit:       0,
			offset:      0,
			wantTrimmed: 0,
			validateFunc: func(t *testing.T, result *models.Datav1QueryLogsRangeResponse) {
				assert.Equal(t, 3, len(result.GridData.Rows))
			},
		},
		{
			name: "grid data - limit without offset",
			payload: &models.Datav1QueryLogsRangeResponse{
				GridData: &models.QueryLogsRangeResponseGridData{
					Columns: []*models.QueryLogsRangeResponseColumnMeta{{Name: "col1"}},
					Rows:    []*models.QueryLogsRangeResponseRow{{}, {}, {}, {}, {}},
				},
				Metadata: struct {
					models.QueryLogsRangeResponseMetadata
				}{},
			},
			limit:       3,
			offset:      0,
			wantTrimmed: 2,
			validateFunc: func(t *testing.T, result *models.Datav1QueryLogsRangeResponse) {
				assert.NotNil(t, result.GridData)
				assert.Equal(t, 3, len(result.GridData.Rows))
				assert.Equal(t, 1, len(result.GridData.Columns))
			},
		},
		{
			name: "grid data - offset and limit",
			payload: &models.Datav1QueryLogsRangeResponse{
				GridData: &models.QueryLogsRangeResponseGridData{
					Rows: []*models.QueryLogsRangeResponseRow{{}, {}, {}, {}, {}},
				},
			},
			limit:       2,
			offset:      2,
			wantTrimmed: 3,
			validateFunc: func(t *testing.T, result *models.Datav1QueryLogsRangeResponse) {
				assert.Equal(t, 2, len(result.GridData.Rows))
			},
		},
		{
			name: "grid data - offset beyond total rows",
			payload: &models.Datav1QueryLogsRangeResponse{
				GridData: &models.QueryLogsRangeResponseGridData{
					Rows: []*models.QueryLogsRangeResponseRow{{}, {}, {}},
				},
			},
			limit:       5,
			offset:      10,
			wantTrimmed: 3,
			validateFunc: func(t *testing.T, result *models.Datav1QueryLogsRangeResponse) {
				assert.Equal(t, 0, len(result.GridData.Rows))
			},
		},
		{
			name: "time series data - limit without offset",
			payload: &models.Datav1QueryLogsRangeResponse{
				TimeSeriesData: struct {
					models.QueryLogsRangeResponseTimeSeriesData
				}{
					QueryLogsRangeResponseTimeSeriesData: models.QueryLogsRangeResponseTimeSeriesData{
						GroupByDimensionNames: []string{"dim1"},
						Series: []*models.TimeSeriesDataTimeSeries{
							{
								AggregationName:        "count",
								GroupByDimensionValues: []string{"val1"},
								Buckets: []*models.TimeSeriesBucket{
									{}, {}, {}, {}, {},
								},
							},
						},
					},
				},
			},
			limit:       3,
			offset:      0,
			wantTrimmed: 2,
			validateFunc: func(t *testing.T, result *models.Datav1QueryLogsRangeResponse) {
				assert.NotNil(t, result.TimeSeriesData)
				assert.Equal(t, 1, len(result.TimeSeriesData.Series))
				assert.Equal(t, 3, len(result.TimeSeriesData.Series[0].Buckets))
				assert.Equal(t, []string{"dim1"}, result.TimeSeriesData.GroupByDimensionNames)
			},
		},
		{
			name: "time series data - multiple series",
			payload: &models.Datav1QueryLogsRangeResponse{
				TimeSeriesData: struct {
					models.QueryLogsRangeResponseTimeSeriesData
				}{
					QueryLogsRangeResponseTimeSeriesData: models.QueryLogsRangeResponseTimeSeriesData{
						Series: []*models.TimeSeriesDataTimeSeries{
							{
								Buckets: []*models.TimeSeriesBucket{
									{}, {}, {}, {},
								},
							},
							{
								Buckets: []*models.TimeSeriesBucket{
									{}, {}, {},
								},
							},
						},
					},
				},
			},
			limit:       2,
			offset:      1,
			wantTrimmed: 3,
			validateFunc: func(t *testing.T, result *models.Datav1QueryLogsRangeResponse) {
				assert.Equal(t, 2, len(result.TimeSeriesData.Series))
				assert.Equal(t, 2, len(result.TimeSeriesData.Series[0].Buckets))
				assert.Equal(t, 2, len(result.TimeSeriesData.Series[1].Buckets))
			},
		},
		{
			name: "preserve metadata",
			payload: &models.Datav1QueryLogsRangeResponse{
				Metadata: struct {
					models.QueryLogsRangeResponseMetadata
				}{
					QueryLogsRangeResponseMetadata: models.QueryLogsRangeResponseMetadata{
						LimitEnforced: true,
					},
				},
				GridData: &models.QueryLogsRangeResponseGridData{
					Rows: []*models.QueryLogsRangeResponseRow{{}, {}},
				},
			},
			limit:       1,
			offset:      0,
			wantTrimmed: 1,
			validateFunc: func(t *testing.T, result *models.Datav1QueryLogsRangeResponse) {
				assert.True(t, result.Metadata.LimitEnforced)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, trimmed := trimLogEntries(tt.payload, tt.limit, tt.offset)
			assert.Equal(t, tt.wantTrimmed, trimmed)
			tt.validateFunc(t, result)
		})
	}
}
