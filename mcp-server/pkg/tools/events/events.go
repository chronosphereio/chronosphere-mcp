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

// Values of the field parameter of ListEventFieldValues.
const (
	eventFieldCategory    = "CATEGORY"
	eventFieldType        = "TYPE"
	eventFieldSource      = "SOURCE"
	eventFieldLabelName   = "LABEL_NAME"
	eventFieldLabelValue  = "LABEL_VALUE"
	eventFieldLensService = "LENS_SERVICE"
)

type Tools struct {
	logger      *zap.Logger
	dataV1API   *datav1.DataV1API
	linkBuilder *links.Builder
}

// NewTools builds the events tools. The DataUnstable client is unused and kept
// only so existing callers still compile.
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

func (t *Tools) listEventFieldValues(ctx context.Context, field string) ([]string, error) {
	resp, err := t.dataV1API.Version1.ListEventFieldValues(&version1.ListEventFieldValuesParams{
		Context: ctx,
		Field:   ptr.To(field),
	})
	if err != nil {
		return nil, err
	}
	return resp.Payload.Values, nil
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
				mcp.WithDescription("List properties you can query on events"),
			),
			Handler: func(ctx context.Context, _ mcp.CallToolRequest) (*tools.Result, error) {
				eventsMetadata := struct {
					Categories       []string `json:"categories"`
					Sources          []string `json:"sources"`
					Types            []string `json:"types"`
					LabelNames       []string `json:"label_names"`
					LensServiceNames []string `json:"lens_service_names"`
				}{}

				var err error
				eventsMetadata.Categories, err = t.listEventFieldValues(ctx, eventFieldCategory)
				if err != nil {
					return nil, fmt.Errorf("failed to get categories for events: %s", err)
				}
				eventsMetadata.LabelNames, err = t.listEventFieldValues(ctx, eventFieldLabelName)
				if err != nil {
					return nil, fmt.Errorf("failed to get label names for events: %s", err)
				}
				eventsMetadata.Sources, err = t.listEventFieldValues(ctx, eventFieldSource)
				if err != nil {
					return nil, fmt.Errorf("failed to get sources for events: %s", err)
				}
				eventsMetadata.Types, err = t.listEventFieldValues(ctx, eventFieldType)
				if err != nil {
					return nil, fmt.Errorf("failed to get types for events: %s", err)
				}
				eventsMetadata.LensServiceNames, err = t.listEventFieldValues(ctx, eventFieldLensService)
				if err != nil {
					return nil, fmt.Errorf("failed to get lens service names for events: %s", err)
				}
				return &tools.Result{
					JSONContent:      eventsMetadata,
					ChronosphereLink: t.linkBuilder.EventExplorer().String(),
				}, nil
			},
		},
		{
			Metadata: tools.NewMetadata("list_events_label_values",
				mcp.WithDescription("List values for a given label name"),
				mcp.WithString("label_name",
					mcp.Required(),
				),
			),
			Handler: func(ctx context.Context, request mcp.CallToolRequest) (*tools.Result, error) {
				labelName, err := params.String(request, "label_name", true, "")
				if err != nil {
					return nil, err
				}

				resp, err := t.dataV1API.Version1.ListEventFieldValues(&version1.ListEventFieldValuesParams{
					Context:   ctx,
					LabelName: ptr.To(labelName),
					Field:     ptr.To(eventFieldLabelValue),
				})
				if err != nil {
					return nil, fmt.Errorf("failed to list label values for events: %s", err)
				}
				return &tools.Result{
					JSONContent: resp,
					ChronosphereLink: t.linkBuilder.Custom("/api/v1/data/events:field-values").
						WithParam("field", eventFieldLabelValue).WithParam("label_name", labelName).String(),
				}, nil
			},
		},
	}
}
