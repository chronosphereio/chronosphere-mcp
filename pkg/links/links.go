// Package links provides builders for constructing links to Chronosphere resources
package links

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// LogExplorerBuilder builds links to the Chronosphere log explorer
type LogExplorerBuilder struct {
	chronosphereURL string
	query           string
	start           *time.Time
	end             *time.Time
	visualization   string
}

// LogExplorer creates a new LogExplorerBuilder for the specified chronosphereURL
func LogExplorer(chronosphereURL string) *LogExplorerBuilder {
	return &LogExplorerBuilder{
		chronosphereURL: chronosphereURL,
		visualization:   "list", // default visualization
	}
}

// WithQuery sets the log query filter
func (b *LogExplorerBuilder) WithQuery(query string) *LogExplorerBuilder {
	b.query = query
	return b
}

// WithTimeRange sets the start and end time for the log query
func (b *LogExplorerBuilder) WithTimeRange(start, end time.Time) *LogExplorerBuilder {
	b.start = &start
	b.end = &end
	return b
}

// WithStartTime sets only the start time
func (b *LogExplorerBuilder) WithStartTime(start time.Time) *LogExplorerBuilder {
	b.start = &start
	return b
}

// WithEndTime sets only the end time
func (b *LogExplorerBuilder) WithEndTime(end time.Time) *LogExplorerBuilder {
	b.end = &end
	return b
}

// WithVisualization sets the visualization type (default: "list")
func (b *LogExplorerBuilder) WithVisualization(viz string) *LogExplorerBuilder {
	b.visualization = viz
	return b
}

// String builds and returns the complete URL
func (b *LogExplorerBuilder) String() string {
	baseURL := fmt.Sprintf("%s%s", b.chronosphereURL, "/logs/explorer")

	params := url.Values{}
	params.Set("visualization", b.visualization)

	if b.query != "" {
		params.Set("query", b.query)
	}

	if b.start != nil {
		params.Set("start", strconv.FormatInt(b.start.UnixMilli(), 10))
	}

	if b.end != nil {
		params.Set("end", strconv.FormatInt(b.end.UnixMilli(), 10))
	}

	if len(params) > 0 {
		return fmt.Sprintf("%s?%s", baseURL, params.Encode())
	}

	return baseURL
}

// URL returns the built URL as a *url.URL
func (b *LogExplorerBuilder) URL() (*url.URL, error) {
	return url.Parse(b.String())
}
