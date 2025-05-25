// Package events contains tools for querying events.
package events

import (
	"context"
	"fmt"

	"github.com/chronosphereio/mcp-server/pkg/ptr"
	"github.com/go-openapi/strfmt"
	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"

	"github.com/chronosphereio/mcp-server/generated/dataunstable/dataunstable/data_unstable"
	"github.com/chronosphereio/mcp-server/generated/dataunstable/models"
	"github.com/chronosphereio/mcp-server/mcp-server/pkg/client"
	"github.com/chronosphereio/mcp-server/mcp-server/pkg/tools"
	"github.com/chronosphereio/mcp-server/mcp-server/pkg/tools/pkg/params"
)

var _ tools.MCPTools = (*Tools)(nil)

type Tools struct {
	logger         *zap.Logger
	clientProvider *client.Provider
}

func NewTools(
	clientProvider *client.Provider,
	logger *zap.Logger,
) (*Tools, error) {
	logger.Info("events tool configured")

	return &Tools{
		logger:         logger,
		clientProvider: clientProvider,
	}, nil
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
			Handler: func(_ context.Context, request mcp.CallToolRequest) (*tools.Result, error) {
				timeRange, err := params.ParseTimeRange(request)
				if err != nil {
					return nil, err
				}
				query, err := params.String(request, "query", false, "")
				if err != nil {
					return nil, err
				}

				queryParams := &data_unstable.ListEventsParams{
					HappenedAfter:  (*strfmt.DateTime)(ptr.To(timeRange.Start)),
					HappenedBefore: (*strfmt.DateTime)(ptr.To(timeRange.End)),
				}
				if query != "" {
					queryParams.Query = ptr.To(query)
				}

				api, err := t.clientProvider.DataUnstableClient()
				if err != nil {
					return nil, err
				}
				resp, err := api.DataUnstable.ListEvents(queryParams)
				if err != nil {
					return nil, fmt.Errorf("failed to list events: %s", err)
				}
				return &tools.Result{
					JSONContent: resp,
				}, nil
			},
		},
		{
			Metadata: tools.NewMetadata("get_events_metadata",
				mcp.WithDescription("List properties you can query on events"),
			),
			Handler: func(ctx context.Context, _ mcp.CallToolRequest) (*tools.Result, error) {
				api, err := t.clientProvider.DataUnstableClient()
				if err != nil {
					return nil, err
				}
				resp, err := api.DataUnstable.ListEventFieldValues(&data_unstable.ListEventFieldValuesParams{
					Context: ctx,
					Field:   ptr.To(string(models.DataunstableEventFieldCATEGORYEVENTFIELD)),
				})

				eventsMetadata := struct {
					Categories       []string `json:"categories"`
					Sources          []string `json:"sources"`
					Types            []string `json:"types"`
					LabelNames       []string `json:"label_names"`
					LensServiceNames []string `json:"lens_service_names"`
				}{}

				if err != nil {
					return nil, fmt.Errorf("failed to get categories for events: %s", err)
				}
				eventsMetadata.Categories = resp.Payload.Values

				resp, err = api.DataUnstable.ListEventFieldValues(&data_unstable.ListEventFieldValuesParams{
					Context: ctx,
					Field:   ptr.To(string(models.DataunstableEventFieldSOURCEEVENTFIELD)),
				})
				if err != nil {
					return nil, fmt.Errorf("failed to get sources for events: %s", err)
				}

				eventsMetadata.Sources = resp.Payload.Values

				resp, err = api.DataUnstable.ListEventFieldValues(&data_unstable.ListEventFieldValuesParams{
					Context: ctx,
					Field:   ptr.To(string(models.DataunstableEventFieldTYPEEVENTFIELD)),
				})
				if err != nil {
					return nil, fmt.Errorf("failed to get types for events: %s", err)
				}

				eventsMetadata.Types = resp.Payload.Values

				resp, err = api.DataUnstable.ListEventFieldValues(&data_unstable.ListEventFieldValuesParams{
					Context: ctx,
					Field:   ptr.To(string(models.DataunstableEventFieldLABELNAMEEVENTFIELD)),
				})
				if err != nil {
					return nil, fmt.Errorf("failed to get label names for events: %s", err)
				}

				eventsMetadata.LabelNames = resp.Payload.Values

				resp, err = api.DataUnstable.ListEventFieldValues(&data_unstable.ListEventFieldValuesParams{
					Context: ctx,
					Field:   ptr.To(string(models.DataunstableEventFieldLABELNAMEEVENTFIELD)),
				})
				if err != nil {
					return nil, fmt.Errorf("failed to get label names for events: %s", err)
				}

				eventsMetadata.LabelNames = resp.Payload.Values

				resp, err = api.DataUnstable.ListEventFieldValues(&data_unstable.ListEventFieldValuesParams{
					Context: ctx,
					Field:   ptr.To(string(models.DataunstableEventFieldLENSSERVICEEVENTFIELD)),
				})
				if err != nil {
					return nil, fmt.Errorf("failed to get label names for events: %s", err)
				}

				eventsMetadata.LensServiceNames = resp.Payload.Values
				return &tools.Result{
					JSONContent: eventsMetadata,
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

				api, err := t.clientProvider.DataUnstableClient()
				if err != nil {
					return nil, err
				}
				resp, err := api.DataUnstable.ListEventFieldValues(&data_unstable.ListEventFieldValuesParams{
					Context:   ctx,
					LabelName: ptr.To(labelName),
					Field:     ptr.To(string(models.DataunstableEventFieldLABELNAMEEVENTFIELD)),
				})
				if err != nil {
					return nil, fmt.Errorf("failed to list label values for events: %s", err)
				}
				return &tools.Result{
					JSONContent: resp,
				}, nil
			},
		},
	}
}
