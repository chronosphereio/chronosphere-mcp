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
	"encoding/json"
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/chronosphereio/chronosphere-mcp/generated/datav1/models"
)

func TestListTracesRequestQueryTypeSerialization(t *testing.T) {
	tests := []struct {
		name              string
		queryType         models.ListTracesRequestQueryType
		expectedJSONValue string
	}{
		{
			name:              "SERVICE_OPERATION query type",
			queryType:         models.ListTracesRequestQueryTypeSERVICEOPERATION,
			expectedJSONValue: `"SERVICE_OPERATION"`,
		},
		{
			name:              "TRACE_IDS query type",
			queryType:         models.ListTracesRequestQueryTypeTRACEIDS,
			expectedJSONValue: `"TRACE_IDS"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a wrapped request body using the same pattern as the fixed code
			now := time.Now()
			body := &listTracesRequestBody{
				Datav1ListTracesRequest: &models.Datav1ListTracesRequest{
					StartTime: strfmt.DateTime(now.Add(-1 * time.Hour)),
					EndTime:   strfmt.DateTime(now),
				},
				queryType: tt.queryType,
			}

			// Marshal to JSON using our custom marshaler
			jsonBytes, err := body.MarshalBinary()
			require.NoError(t, err)

			t.Logf("Generated JSON: %s", string(jsonBytes))

			// Verify the JSON contains the correct query_type value
			var result map[string]interface{}
			err = json.Unmarshal(jsonBytes, &result)
			require.NoError(t, err)

			// Check that query_type is present
			queryTypeValue, exists := result["query_type"]
			require.True(t, exists, "query_type field should exist in JSON")

			t.Logf("query_type value: %+v (type: %T)", queryTypeValue, queryTypeValue)

			// The value should be a string, not an object
			queryTypeStr, ok := queryTypeValue.(string)
			require.True(t, ok, "query_type should be a string, got %T: %+v", queryTypeValue, queryTypeValue)

			assert.Equal(t, string(tt.queryType), queryTypeStr)
		})
	}
}
