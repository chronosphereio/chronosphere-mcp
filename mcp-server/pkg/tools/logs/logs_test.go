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

	"github.com/chronosphereio/chronosphere-mcp/generated/dataunstable/models"
)

func TestTrimLogEntries(t *testing.T) {
	tests := []struct {
		name         string
		payload      *models.DataunstableGetRangeQueryResponse
		limit        int
		offset       int
		wantTrimmed  int
		validateFunc func(t *testing.T, result *models.DataunstableGetRangeQueryResponse)
	}{
		{
			name:        "nil payload",
			payload:     nil,
			limit:       10,
			offset:      0,
			wantTrimmed: 0,
			validateFunc: func(t *testing.T, result *models.DataunstableGetRangeQueryResponse) {
				assert.Nil(t, result)
			},
		},
		{
			name: "zero limit and offset returns original",
			payload: &models.DataunstableGetRangeQueryResponse{
				GridData: &models.DataunstableLogQueryGridData{
					Rows: []*models.DataunstableRow{{}, {}, {}},
				},
			},
			limit:       0,
			offset:      0,
			wantTrimmed: 0,
			validateFunc: func(t *testing.T, result *models.DataunstableGetRangeQueryResponse) {
				assert.Equal(t, 3, len(result.GridData.Rows))
			},
		},
		{
			name: "grid data - limit without offset",
			payload: &models.DataunstableGetRangeQueryResponse{
				GridData: &models.DataunstableLogQueryGridData{
					Columns: []*models.DataunstableColumnMeta{{Name: "col1"}},
					Rows:    []*models.DataunstableRow{{}, {}, {}, {}, {}},
				},
				Metadata: &models.GetRangeQueryResponseRangeQueryMetadata{},
			},
			limit:       3,
			offset:      0,
			wantTrimmed: 2,
			validateFunc: func(t *testing.T, result *models.DataunstableGetRangeQueryResponse) {
				assert.NotNil(t, result.GridData)
				assert.Equal(t, 3, len(result.GridData.Rows))
				assert.Equal(t, 1, len(result.GridData.Columns))
				assert.NotNil(t, result.Metadata)
			},
		},
		{
			name: "grid data - offset and limit",
			payload: &models.DataunstableGetRangeQueryResponse{
				GridData: &models.DataunstableLogQueryGridData{
					Rows: []*models.DataunstableRow{{}, {}, {}, {}, {}},
				},
			},
			limit:       2,
			offset:      2,
			wantTrimmed: 3,
			validateFunc: func(t *testing.T, result *models.DataunstableGetRangeQueryResponse) {
				assert.Equal(t, 2, len(result.GridData.Rows))
			},
		},
		{
			name: "grid data - offset beyond total rows",
			payload: &models.DataunstableGetRangeQueryResponse{
				GridData: &models.DataunstableLogQueryGridData{
					Rows: []*models.DataunstableRow{{}, {}, {}},
				},
			},
			limit:       5,
			offset:      10,
			wantTrimmed: 3,
			validateFunc: func(t *testing.T, result *models.DataunstableGetRangeQueryResponse) {
				assert.Equal(t, 0, len(result.GridData.Rows))
			},
		},
		{
			name: "time series data - limit without offset",
			payload: &models.DataunstableGetRangeQueryResponse{
				TimeSeriesData: &models.DataunstableLogQueryTimeSeriesData{
					GroupByDimensionNames: []string{"dim1"},
					Series: []*models.LogQueryTimeSeriesDataLogQueryTimeSeries{
						{
							AggregationName:        "count",
							GroupByDimensionValues: []string{"val1"},
							Buckets: []*models.LogQueryTimeSeriesLogQueryTimeSeriesBucket{
								{}, {}, {}, {}, {},
							},
						},
					},
				},
			},
			limit:       3,
			offset:      0,
			wantTrimmed: 2,
			validateFunc: func(t *testing.T, result *models.DataunstableGetRangeQueryResponse) {
				assert.NotNil(t, result.TimeSeriesData)
				assert.Equal(t, 1, len(result.TimeSeriesData.Series))
				assert.Equal(t, 3, len(result.TimeSeriesData.Series[0].Buckets))
				assert.Equal(t, []string{"dim1"}, result.TimeSeriesData.GroupByDimensionNames)
			},
		},
		{
			name: "time series data - multiple series",
			payload: &models.DataunstableGetRangeQueryResponse{
				TimeSeriesData: &models.DataunstableLogQueryTimeSeriesData{
					Series: []*models.LogQueryTimeSeriesDataLogQueryTimeSeries{
						{
							Buckets: []*models.LogQueryTimeSeriesLogQueryTimeSeriesBucket{
								{}, {}, {}, {},
							},
						},
						{
							Buckets: []*models.LogQueryTimeSeriesLogQueryTimeSeriesBucket{
								{}, {}, {},
							},
						},
					},
				},
			},
			limit:       2,
			offset:      1,
			wantTrimmed: 3,
			validateFunc: func(t *testing.T, result *models.DataunstableGetRangeQueryResponse) {
				assert.Equal(t, 2, len(result.TimeSeriesData.Series))
				assert.Equal(t, 2, len(result.TimeSeriesData.Series[0].Buckets))
				assert.Equal(t, 2, len(result.TimeSeriesData.Series[1].Buckets))
			},
		},
		{
			name: "preserve metadata",
			payload: &models.DataunstableGetRangeQueryResponse{
				Metadata: &models.GetRangeQueryResponseRangeQueryMetadata{
					LimitEnforced: true,
				},
				GridData: &models.DataunstableLogQueryGridData{
					Rows: []*models.DataunstableRow{{}, {}},
				},
			},
			limit:       1,
			offset:      0,
			wantTrimmed: 1,
			validateFunc: func(t *testing.T, result *models.DataunstableGetRangeQueryResponse) {
				assert.NotNil(t, result.Metadata)
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
