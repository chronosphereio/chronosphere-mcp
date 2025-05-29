// Package monitors provides tools for querying Chronosphere monitor statuses.
package monitors

import (
	"context"
	"fmt"

	"github.com/chronosphereio/chronosphere-mcp/generated/stateunstable/stateunstable"
	"github.com/chronosphereio/chronosphere-mcp/pkg/ptr"
	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"

	"github.com/chronosphereio/chronosphere-mcp/generated/stateunstable/stateunstable/state_unstable"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools/pkg/params"
)

var _ tools.MCPTools = (*Tools)(nil)

// Tools represents the monitor tools.
type Tools struct {
	logger *zap.Logger
	api    *stateunstable.StateUnstableAPI
}

// NewTools creates a new set of monitor tools.
func NewTools(api *stateunstable.StateUnstableAPI, logger *zap.Logger) (*Tools, error) {
	logger.Info("monitor status tool configured")

	return &Tools{
		logger: logger,
		api:    api,
	}, nil
}

// Using the StateUnstableClient method from the Provider

// MCPTools returns the list of monitor tools.
func (t *Tools) MCPTools() []tools.MCPTool {
	return []tools.MCPTool{
		{
			Metadata: tools.NewMetadata("list_monitor_statuses",
				mcp.WithDescription("Lists the current status of monitors in Chronosphere. Returns monitor statuses with alert states and optional signal and series details."),
				params.WithStringArray("monitor_slugs",
					mcp.Description("Filter by monitor slug. If all filters are empty, return status for all monitors."),
				),
				params.WithStringArray("collection_slugs",
					mcp.Description("Filter monitor state by owning collection."),
				),
				params.WithStringArray("team_slugs",
					mcp.Description("Filter monitor state by owning team."),
				),
				mcp.WithBoolean("include_signal_statuses",
					mcp.Description("Include signal status details in the response."),
					mcp.DefaultBool(false),
				),
				mcp.WithBoolean("include_series_statuses",
					mcp.Description("Include series status details in the response. Requires include_signal_statuses to be true."),
					mcp.DefaultBool(false),
				),
				mcp.WithString("sort_by",
					mcp.Description("Sort results by a specific field. Currently supports 'SORT_BY_STATE'."),
					mcp.Enum("SORT_BY_STATE"),
				),
			),
			Handler: func(ctx context.Context, request mcp.CallToolRequest) (*tools.Result, error) {
				// Parse parameters
				monitorSlugs, err := params.StringArray(request, "monitor_slugs", false, nil)
				if err != nil {
					return nil, err
				}

				collectionSlugs, err := params.StringArray(request, "collection_slugs", false, nil)
				if err != nil {
					return nil, err
				}

				teamSlugs, err := params.StringArray(request, "team_slugs", false, nil)
				if err != nil {
					return nil, err
				}

				includeSignalStatuses, err := params.Bool(request, "include_signal_statuses", false, false)
				if err != nil {
					return nil, err
				}

				includeSeriesStatuses, err := params.Bool(request, "include_series_statuses", false, false)
				if err != nil {
					return nil, err
				}

				sortBy, err := params.String(request, "sort_by", false, "")
				if err != nil {
					return nil, err
				}

				// Construct query parameters
				queryParams := state_unstable.NewListMonitorStatusesParams().
					WithContext(ctx)

				if len(monitorSlugs) > 0 {
					queryParams.SetMonitorSlugs(monitorSlugs)
				}

				if len(collectionSlugs) > 0 {
					queryParams.SetCollectionSlugs(collectionSlugs)
				}

				if len(teamSlugs) > 0 {
					queryParams.SetTeamSlugs(teamSlugs)
				}

				if includeSignalStatuses {
					queryParams.SetIncludeSignalStatuses(ptr.To(includeSignalStatuses))
				}

				if includeSeriesStatuses {
					if !includeSignalStatuses {
						return nil, fmt.Errorf("include_signal_statuses must be true for include_series_statuses to be true")
					}
					queryParams.SetIncludeSeriesStatuses(ptr.To(includeSeriesStatuses))
				}

				if sortBy != "" {
					queryParams.SetSortBy(ptr.To(sortBy))
				}

				t.logger.Info("list monitor statuses", zap.Any("params", queryParams))

				resp, err := t.api.StateUnstable.ListMonitorStatuses(queryParams)
				if err != nil {
					return nil, fmt.Errorf("failed to list monitor statuses: %s", err)
				}

				return &tools.Result{
					JSONContent: resp,
				}, nil
			},
		},
	}
}
