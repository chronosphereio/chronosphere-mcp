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

package params

import (
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

// WithTimeRange adds start and end time parameters to the tool.
func WithTimeRange() mcp.ToolOption {
	return func(tool *mcp.Tool) {
		mcp.WithString("start",
			mcp.Description("Optional. Start time in RFC3339 format or timestamp in seconds. Defaults to end - 1 hour."),
		)(tool)
		mcp.WithString("end",
			mcp.Description("Optional. End time in RFC3339 format or timestamp in seconds. Defaults to the current time."),
		)(tool)
	}
}

type TimeRange struct {
	Start time.Time
	End   time.Time
}

// ParseTimeRange parses start and end time parameters from the request.
func ParseTimeRange(request mcp.CallToolRequest) (*TimeRange, error) {
	end, err := Time(request, "end", false, time.Now())
	if err != nil {
		return nil, err
	}
	start, err := Time(request, "start", false, end.Add(-time.Hour))
	if err != nil {
		return nil, err
	}
	return &TimeRange{
		Start: start,
		End:   end,
	}, nil
}
