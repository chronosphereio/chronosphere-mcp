package configv1

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"

	"github.com/chronosphereio/chronosphere-mcp/generated/configv1/configv1"
	"github.com/chronosphereio/chronosphere-mcp/generated/configv1/configv1/classic_dashboard"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools/pkg/params"
	"github.com/chronosphereio/chronosphere-mcp/pkg/ptr"
)

func GetClassicDashboard(api *configv1.ConfigV1API, logger *zap.Logger) tools.MCPTool {
	return tools.MCPTool{
		Metadata: tools.NewMetadata("get_classic_dashboard",
			mcp.WithDescription("Get classic-dashboards resource"),

			mcp.WithString("slug",
				mcp.Description(""),
				mcp.Required(),
			),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*tools.Result, error) {
			slug, err := params.String(request, "slug", true, "")
			if err != nil {
				return nil, err
			}

			queryParams := &classic_dashboard.ReadClassicDashboardParams{
				Context: ctx,

				Slug: slug,
			}

			resp, err := api.ClassicDashboard.ReadClassicDashboard(queryParams)
			if err != nil {
				return nil, fmt.Errorf("failed to call ReadClassicDashboard: %s", err)
			}
			return &tools.Result{
				JSONContent: resp,
			}, nil
		},
	}
}

func ListClassicDashboards(api *configv1.ConfigV1API, logger *zap.Logger) tools.MCPTool {
	return tools.MCPTool{
		Metadata: tools.NewMetadata("list_classic_dashboards",
			mcp.WithDescription("List classic-dashboards resources"),

			params.WithStringArray("bucket_slugs",
				mcp.Description("Filters results by bucket_slug, where any ClassicDashboard with a matching bucket_slug in the given list (and matches all other filters) is returned."),
			),

			params.WithStringArray("collection_slugs",
				mcp.Description("Filters results by collection_slug, where any ClassicDashboard with a matching collection_slug in the given list (and matches all other filters) is returned."),
			),

			mcp.WithBoolean("include_dashboard_json",
				mcp.Description("Optional flag to populate the dashboard_json of the returned dashboards. By default, dashboard_json will be left empty."),
			),

			params.WithStringArray("names",
				mcp.Description("Filters results by name, where any ClassicDashboard with a matching name in the given list (and matches all other filters) is returned."),
			),

			mcp.WithNumber("page_max_size",
				mcp.Description("Page size preference (i.e. how many items are returned in the next page). If zero, the server will use a default. Regardless of what size is given, clients must never assume how many items will be returned."),
			),

			mcp.WithString("page_token",
				mcp.Description("Opaque page token identifying which page to request. An empty token identifies the first page."),
			),

			params.WithStringArray("slugs",
				mcp.Description("Filters results by slug, where any ClassicDashboard with a matching slug in the given list (and matches all other filters) is returned."),
			),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*tools.Result, error) {
			bucketSlugs, err := params.StringArray(request, "bucket_slugs", false, nil)
			if err != nil {
				return nil, err
			}

			collectionSlugs, err := params.StringArray(request, "collection_slugs", false, nil)
			if err != nil {
				return nil, err
			}

			includeDashboardJson, err := params.Bool(request, "include_dashboard_json", false, false)
			if err != nil {
				return nil, err
			}

			names, err := params.StringArray(request, "names", false, nil)
			if err != nil {
				return nil, err
			}

			pageMaxSize, err := params.Int(request, "page_max_size", false, 0)
			if err != nil {
				return nil, err
			}

			pageToken, err := params.String(request, "page_token", false, "")
			if err != nil {
				return nil, err
			}

			slugs, err := params.StringArray(request, "slugs", false, nil)
			if err != nil {
				return nil, err
			}

			queryParams := &classic_dashboard.ListClassicDashboardsParams{
				Context: ctx,

				BucketSlugs: bucketSlugs,

				CollectionSlugs: collectionSlugs,

				IncludeDashboardJSON: ptr.To(includeDashboardJson),

				Names: names,

				PageMaxSize: ptr.To(int64(pageMaxSize)),

				PageToken: &pageToken,

				Slugs: slugs,
			}

			resp, err := api.ClassicDashboard.ListClassicDashboards(queryParams)
			if err != nil {
				return nil, fmt.Errorf("failed to call ListClassicDashboards: %s", err)
			}
			return &tools.Result{
				JSONContent: resp,
			}, nil
		},
	}
}
