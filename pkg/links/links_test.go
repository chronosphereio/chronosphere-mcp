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

package links

import (
	"testing"
	"time"
)

func TestLogExplorer(t *testing.T) {
	start := time.Unix(1747887246, 0)
	end := time.Unix(1747890846, 0)

	builder := NewBuilder("https://rc.chronosphere.io")
	tests := []struct {
		name     string
		builder  *LogExplorerBuilder
		expected string
	}{
		{
			name:     "basic log explorer",
			builder:  builder.LogExplorer(),
			expected: "https://rc.chronosphere.io/logs/explorer?visualization=list",
		},
		{
			name:     "with query",
			builder:  builder.LogExplorer().WithQuery(`service="chronogateway"`),
			expected: "https://rc.chronosphere.io/logs/explorer?query=service%3D%22chronogateway%22&visualization=list",
		},
		{
			name:     "with time range and query",
			builder:  builder.LogExplorer().WithTimeRange(start, end).WithQuery(`service="chronogateway"`),
			expected: "https://rc.chronosphere.io/logs/explorer?end=1747890846000&query=service%3D%22chronogateway%22&start=1747887246000&visualization=list",
		},
		{
			name:     "with different visualization",
			builder:  builder.LogExplorer().WithVisualization("chart"),
			expected: "https://rc.chronosphere.io/logs/explorer?visualization=chart",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.builder.String()
			if result != tt.expected {
				t.Errorf("LogExplorer() = %v, want %v", result, tt.expected)
			}
		})
	}
}
