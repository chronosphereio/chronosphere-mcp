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

// Package events contains tools for querying events.
package events

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"

	"github.com/chronosphereio/chronosphere-mcp/generated/dataunstable/dataunstable"
	"github.com/chronosphereio/chronosphere-mcp/generated/datav1/datav1"
	"github.com/chronosphereio/chronosphere-mcp/generated/datav1/datav1/version1"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools/pkg/params"
	"github.com/chronosphereio/chronosphere-mcp/pkg/links"
	"github.com/chronosphereio/chronosphere-mcp/pkg/ptr"
)

var _ tools.MCPTools = (*Tools)(nil)

type Tools struct {
	logger      *zap.Logger
	dataV1API   *datav1.DataV1API
	linkBuilder *links.Builder
}

// NewTools builds the events tools. The DataUnstable client is unused and is
// kept only so existing callers keep compiling.
func NewTools(
	dataV1API *datav1.DataV1API,
	_ *dataunstable.DataUnstableAPI,
	logger *zap.Logger,
	linkBuilder *links.Builder,
) (*Tools, error) {
	logger.Info("events tool configured")

	return &Tools{
		logger:      logger,
		dataV1API:   dataV1API,
		linkBuilder: linkBuilder,
	}, nil
}

func (t *Tools) GroupName() string {
	return "events"
}

func (t *Tools) MCPTools() []tools.MCPTool {
	return []tools.MCPTool{
		{
			Metadata: tools.NewMetadata("list_events",
				mcp.WithDescription("List events from a given query"),
				mcp.WithString("query",
					mcp.Description("The query to filter events e.g. categories, types, sources and arbitrary labels.")),
				params.WithTimeRange(),
			),
			Handler: func(ctx context.Context, request mcp.CallToolRequest) (*tools.Result, error) {
				timeRange, err := params.ParseTimeRange(request)
				if err != nil {
					return nil, err
				}
				query, err := params.String(request, "query", false, "")
				if err != nil {
					return nil, err
				}

				queryParams := &version1.ListEventsParams{
					Context:        ctx,
					HappenedAfter:  (*strfmt.DateTime)(ptr.To(timeRange.Start)),
					HappenedBefore: (*strfmt.DateTime)(ptr.To(timeRange.End)),
				}
				if query != "" {
					queryParams.Query = ptr.To(query)
				}

				resp, err := t.dataV1API.Version1.ListEvents(queryParams)
				if err != nil {
					return nil, fmt.Errorf("failed to list events: %s", err)
				}
				return &tools.Result{
					JSONContent: resp,
					ChronosphereLink: t.linkBuilder.EventExplorer().
						WithQuery(query).
						WithTimeRange(timeRange.Start, timeRange.End).
						String(),
				}, nil
			},
		},
		{
			Metadata: tools.NewMetadata("get_events_metadata",
				mcp.WithDescription(
					"List properties you can query on events: the distinct categories, sources, types and label names "+
						"of events that happened in the time range. Values are drawn from at most 10,000 of the most "+
						"recent events in the range; truncated is true when that limit was reached and the lists may be partial."),
				withFacetTimeRange(),
			),
			Handler: func(ctx context.Context, request mcp.CallToolRequest) (*tools.Result, error) {
				timeRange, err := parseFacetTimeRange(request)
				if err != nil {
					return nil, err
				}

				facets, err := t.collectEventFacets(ctx, timeRange, "", "", maxFacetEvents)
				if err != nil {
					return nil, err
				}

				return &tools.Result{
					JSONContent: struct {
						Categories []string `json:"categories"`
						Sources    []string `json:"sources"`
						Types      []string `json:"types"`
						LabelNames []string `json:"label_names"`
						Truncated  bool     `json:"truncated"`
					}{
						Categories: sortedKeys(facets.categories),
						Sources:    sortedKeys(facets.sources),
						Types:      sortedKeys(facets.types),
						LabelNames: sortedKeys(facets.labelNames),
						Truncated:  facets.truncated,
					},
					ChronosphereLink: t.linkBuilder.EventExplorer().
						WithTimeRange(timeRange.Start, timeRange.End).
						String(),
				}, nil
			},
		},
		{
			Metadata: tools.NewMetadata("list_events_label_values",
				mcp.WithDescription(
					"List the distinct values of a label across events that happened in the time range and match the "+
						"optional query. Values are drawn from at most 10,000 of the most recent matching events; "+
						"truncated is true when that limit was reached and the list may be partial."),
				mcp.WithString("label_name",
					mcp.Required(),
					mcp.Description("The label name to list values for."),
				),
				mcp.WithString("query",
					mcp.Description("Optional. The query to filter events e.g. categories, types, sources and arbitrary labels.")),
				withFacetTimeRange(),
			),
			Handler: func(ctx context.Context, request mcp.CallToolRequest) (*tools.Result, error) {
				labelName, err := params.String(request, "label_name", true, "")
				if err != nil {
					return nil, err
				}
				query, err := params.String(request, "query", false, "")
				if err != nil {
					return nil, err
				}
				timeRange, err := parseFacetTimeRange(request)
				if err != nil {
					return nil, err
				}

				facets, err := t.collectEventFacets(ctx, timeRange, query, labelName, maxFacetEvents)
				if err != nil {
					return nil, err
				}

				return &tools.Result{
					JSONContent: struct {
						LabelName string   `json:"label_name"`
						Values    []string `json:"values"`
						Truncated bool     `json:"truncated"`
					}{
						LabelName: labelName,
						Values:    sortedKeys(facets.labelValues),
						Truncated: facets.truncated,
					},
					ChronosphereLink: t.linkBuilder.EventExplorer().
						WithQuery(query).
						WithTimeRange(timeRange.Start, timeRange.End).
						String(),
				}, nil
			},
		},
	}
}

// maxFacetEvents bounds how many events the metadata tools read, matching the
// most ListEvents returns for a single query.
const maxFacetEvents = 10000

// defaultFacetWindow is how far back the metadata tools look when no start is given.
const defaultFacetWindow = 24 * time.Hour

func withFacetTimeRange() mcp.ToolOption {
	return func(tool *mcp.Tool) {
		mcp.WithString("start",
			mcp.Description("Optional. Start time in RFC3339 format or timestamp in seconds. Defaults to end - 24 hours."),
		)(tool)
		mcp.WithString("end",
			mcp.Description("Optional. End time in RFC3339 format or timestamp in seconds. Defaults to the current time."),
		)(tool)
	}
}

func parseFacetTimeRange(request mcp.CallToolRequest) (*params.TimeRange, error) {
	end, err := params.Time(request, "end", false, time.Now())
	if err != nil {
		return nil, err
	}
	start, err := params.Time(request, "start", false, end.Add(-defaultFacetWindow))
	if err != nil {
		return nil, err
	}
	return &params.TimeRange{Start: start, End: end}, nil
}

// eventFacets holds the distinct values seen across a set of events.
type eventFacets struct {
	categories  map[string]struct{}
	sources     map[string]struct{}
	types       map[string]struct{}
	labelNames  map[string]struct{}
	labelValues map[string]struct{}
	// truncated is true when maxEvents events were read, so more may exist.
	truncated bool
}

// collectEventFacets pages through ListEvents for the time range and query,
// reading at most maxEvents events, and collects their distinct categories,
// sources, types and label names. When labelName is set, it also collects the
// distinct values of that label.
func (t *Tools) collectEventFacets(
	ctx context.Context,
	timeRange *params.TimeRange,
	query string,
	labelName string,
	maxEvents int,
) (*eventFacets, error) {
	facets := &eventFacets{
		categories:  map[string]struct{}{},
		sources:     map[string]struct{}{},
		types:       map[string]struct{}{},
		labelNames:  map[string]struct{}{},
		labelValues: map[string]struct{}{},
	}

	var (
		read      int
		pageToken string
	)
	for {
		queryParams := &version1.ListEventsParams{
			Context:        ctx,
			HappenedAfter:  (*strfmt.DateTime)(ptr.To(timeRange.Start)),
			HappenedBefore: (*strfmt.DateTime)(ptr.To(timeRange.End)),
		}
		if query != "" {
			queryParams.Query = ptr.To(query)
		}
		if pageToken != "" {
			queryParams.PageToken = ptr.To(pageToken)
		}

		resp, err := t.dataV1API.Version1.ListEvents(queryParams)
		if err != nil {
			return nil, fmt.Errorf("failed to list events: %s", err)
		}

		for _, event := range resp.Payload.Events {
			if read >= maxEvents {
				break
			}
			read++
			if event == nil {
				continue
			}
			addNonEmpty(facets.categories, event.Category)
			addNonEmpty(facets.sources, event.Source)
			addNonEmpty(facets.types, event.Type)
			for name, value := range event.Labels {
				addNonEmpty(facets.labelNames, name)
				if labelName != "" && name == labelName {
					addNonEmpty(facets.labelValues, value)
				}
			}
		}

		if read >= maxEvents {
			facets.truncated = true
			return facets, nil
		}
		if resp.Payload.Page == nil || resp.Payload.Page.NextToken == "" || len(resp.Payload.Events) == 0 {
			return facets, nil
		}
		pageToken = resp.Payload.Page.NextToken
	}
}

func addNonEmpty(set map[string]struct{}, value string) {
	if value != "" {
		set[value] = struct{}{}
	}
}

// sortedKeys returns the set's members in sorted order, never nil.
func sortedKeys(set map[string]struct{}) []string {
	keys := make([]string, 0, len(set))
	for k := range set {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
