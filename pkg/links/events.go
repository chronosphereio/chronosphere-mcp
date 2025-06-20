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
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

// EventExplorerBuilder builds links to the Chronosphere event explorer
type EventExplorerBuilder struct {
	chronosphereURL string
	query           string
	start           *time.Time
	end             *time.Time
}

// EventExplorer creates a new EventExplorerBuilder for the specified chronosphereURL
func (b *Builder) EventExplorer() *EventExplorerBuilder {
	return &EventExplorerBuilder{
		chronosphereURL: b.chronosphereURL,
	}
}

// WithQuery sets the event query filter
func (b *EventExplorerBuilder) WithQuery(query string) *EventExplorerBuilder {
	b.query = query
	return b
}

// WithTimeRange sets the start and end time for the event query
func (b *EventExplorerBuilder) WithTimeRange(start, end time.Time) *EventExplorerBuilder {
	b.start = &start
	b.end = &end
	return b
}

// WithStartTime sets only the start time
func (b *EventExplorerBuilder) WithStartTime(start time.Time) *EventExplorerBuilder {
	b.start = &start
	return b
}

// WithEndTime sets only the end time
func (b *EventExplorerBuilder) WithEndTime(end time.Time) *EventExplorerBuilder {
	b.end = &end
	return b
}

// String builds and returns the complete URL
func (b *EventExplorerBuilder) String() string {
	baseURL := fmt.Sprintf("%s%s", b.chronosphereURL, "/events/explorer")

	params := url.Values{}

	// Set time range parameters using relative time format if provided
	if b.start != nil && b.end != nil {
		// Calculate relative time from end
		duration := b.end.Sub(*b.start)
		if duration > 0 {
			// Format duration as relative time (e.g., "1h", "30m", "2d")
			if duration.Hours() >= 24 {
				days := int(duration.Hours() / 24)
				params.Set("start", fmt.Sprintf("%dd", days))
			} else if duration.Hours() >= 1 {
				hours := int(duration.Hours())
				params.Set("start", fmt.Sprintf("%dh", hours))
			} else {
				minutes := int(duration.Minutes())
				params.Set("start", fmt.Sprintf("%dm", minutes))
			}
		}
	}

	// Set event query parameters if provided
	if b.query != "" {
		// Create the change_events_params JSON structure
		eventParams := map[string]interface{}{
			"query": b.query,
			"slug":  "",
		}
		if eventParamsJSON, err := json.Marshal(eventParams); err == nil {
			params.Set("change_events_params", string(eventParamsJSON))
		}
	}

	if len(params) > 0 {
		return fmt.Sprintf("%s?%s", baseURL, params.Encode())
	}

	return baseURL
}

// URL returns the built URL as a *url.URL
func (b *EventExplorerBuilder) URL() (*url.URL, error) {
	return url.Parse(b.String())
}
