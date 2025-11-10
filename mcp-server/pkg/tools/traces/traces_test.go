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

package traces

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/chronosphereio/chronosphere-mcp/generated/datav1/models"
)

func TestTrimTraces(t *testing.T) {
	tests := []struct {
		name         string
		payload      *models.Datav1ListTracesResponse
		limit        int
		offset       int
		wantTrimmed  int
		validateFunc func(t *testing.T, result *models.Datav1ListTracesResponse)
	}{
		{
			name:        "nil payload",
			payload:     nil,
			limit:       10,
			offset:      0,
			wantTrimmed: 0,
			validateFunc: func(t *testing.T, result *models.Datav1ListTracesResponse) {
				assert.Nil(t, result)
			},
		},
		{
			name: "zero limit and offset returns original",
			payload: &models.Datav1ListTracesResponse{
				Traces: []*models.V1TracesData{{}, {}, {}},
			},
			limit:       0,
			offset:      0,
			wantTrimmed: 0,
			validateFunc: func(t *testing.T, result *models.Datav1ListTracesResponse) {
				assert.Equal(t, 3, len(result.Traces))
			},
		},
		{
			name: "limit without offset",
			payload: &models.Datav1ListTracesResponse{
				Traces: []*models.V1TracesData{{}, {}, {}, {}, {}},
			},
			limit:       3,
			offset:      0,
			wantTrimmed: 2,
			validateFunc: func(t *testing.T, result *models.Datav1ListTracesResponse) {
				assert.NotNil(t, result)
				assert.Equal(t, 3, len(result.Traces))
			},
		},
		{
			name: "offset and limit - pagination",
			payload: &models.Datav1ListTracesResponse{
				Traces: []*models.V1TracesData{{}, {}, {}, {}, {}},
			},
			limit:       2,
			offset:      2,
			wantTrimmed: 3,
			validateFunc: func(t *testing.T, result *models.Datav1ListTracesResponse) {
				assert.Equal(t, 2, len(result.Traces))
			},
		},
		{
			name: "offset beyond total traces",
			payload: &models.Datav1ListTracesResponse{
				Traces: []*models.V1TracesData{{}, {}, {}},
			},
			limit:       5,
			offset:      10,
			wantTrimmed: 3,
			validateFunc: func(t *testing.T, result *models.Datav1ListTracesResponse) {
				assert.Equal(t, 0, len(result.Traces))
			},
		},
		{
			name: "limit exceeds remaining traces",
			payload: &models.Datav1ListTracesResponse{
				Traces: []*models.V1TracesData{{}, {}, {}, {}, {}},
			},
			limit:       10,
			offset:      3,
			wantTrimmed: 3,
			validateFunc: func(t *testing.T, result *models.Datav1ListTracesResponse) {
				assert.Equal(t, 2, len(result.Traces))
			},
		},
		{
			name: "single trace with offset 0 and limit 1",
			payload: &models.Datav1ListTracesResponse{
				Traces: []*models.V1TracesData{{}, {}, {}},
			},
			limit:       1,
			offset:      0,
			wantTrimmed: 2,
			validateFunc: func(t *testing.T, result *models.Datav1ListTracesResponse) {
				assert.Equal(t, 1, len(result.Traces))
			},
		},
		{
			name: "nil traces slice",
			payload: &models.Datav1ListTracesResponse{
				Traces: nil,
			},
			limit:       5,
			offset:      0,
			wantTrimmed: 0,
			validateFunc: func(t *testing.T, result *models.Datav1ListTracesResponse) {
				assert.Equal(t, result, &models.Datav1ListTracesResponse{
					Traces: nil,
				})
			},
		},
		{
			name: "empty traces slice",
			payload: &models.Datav1ListTracesResponse{
				Traces: []*models.V1TracesData{},
			},
			limit:       5,
			offset:      0,
			wantTrimmed: 0,
			validateFunc: func(t *testing.T, result *models.Datav1ListTracesResponse) {
				assert.Equal(t, 0, len(result.Traces))
			},
		},
		{
			name: "only offset, no limit",
			payload: &models.Datav1ListTracesResponse{
				Traces: []*models.V1TracesData{{}, {}, {}, {}, {}},
			},
			limit:       0,
			offset:      2,
			wantTrimmed: 2,
			validateFunc: func(t *testing.T, result *models.Datav1ListTracesResponse) {
				// With limit=0, it should return all traces from offset onwards
				assert.Equal(t, 3, len(result.Traces))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, trimmed := trimTraces(tt.payload, tt.limit, tt.offset)
			assert.Equal(t, tt.wantTrimmed, trimmed)
			tt.validateFunc(t, result)
		})
	}
}
