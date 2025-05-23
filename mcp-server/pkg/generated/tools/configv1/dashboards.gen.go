package configv1

import (
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"

	"github.com/chronosphereio/mcp-server/generated/configv1/configv1/dashboard"
	"github.com/chronosphereio/mcp-server/mcp-server/pkg/client"
	"github.com/chronosphereio/mcp-server/mcp-server/pkg/tools"
	"github.com/chronosphereio/mcp-server/mcp-server/pkg/tools/pkg/params"
	"github.com/chronosphereio/mcp-server/pkg/ptr"
)

func GetDashboard(clientProvider *client.Provider, logger *zap.Logger) tools.MCPTool {
	return tools.MCPTool{
		Metadata: tools.NewMetadata("get_dashboard",
			mcp.WithDescription("Get dashboards resource"),

			mcp.WithString("slug",
				mcp.Description(""),
				mcp.Required(),
			),
		),
		Handler: func(session tools.Session, request mcp.CallToolRequest) (*tools.Result, error) {
			slug, err := params.String(request, "slug", true, "")
			if err != nil {
				return nil, err
			}

			queryParams := &dashboard.ReadDashboardParams{
				Slug: slug,
			}

			api, err := clientProvider.ConfigV1Client(session)
			if err != nil {
				return nil, err
			}
			resp, err := api.Dashboard.ReadDashboard(queryParams)
			if err != nil {
				return nil, fmt.Errorf("failed to call ReadDashboard: %s", err)
			}
			return &tools.Result{
				JSONContent: resp,
			}, nil
		},
	}
}

func ListDashboards(clientProvider *client.Provider, logger *zap.Logger) tools.MCPTool {
	return tools.MCPTool{
		Metadata: tools.NewMetadata("list_dashboards",
			mcp.WithDescription("List dashboards resources"),

			mcp.WithArray("collection_slugs",
				mcp.Description("Filters results by collection_slug, where any Dashboard with a matching collection_slug in the given list (and matches all other filters) is returned."),
			),

			mcp.WithBoolean("include_dashboard_json",
				mcp.Description("Optional flag to populate the dashboard_json of the returned dashboards. By default, dashboard_json will be left empty."),
			),

			mcp.WithArray("names",
				mcp.Description("Filters results by name, where any Dashboard with a matching name in the given list (and matches all other filters) is returned."),
			),

			mcp.WithNumber("page_max_size",
				mcp.Description("Page size preference (i.e. how many items are returned in the next page). If zero, the server will use a default. Regardless of what size is given, clients must never assume how many items will be returned."),
			),

			mcp.WithString("page_token",
				mcp.Description("Opaque page token identifying which page to request. An empty token identifies the first page."),
			),

			mcp.WithArray("slugs",
				mcp.Description("Filters results by slug, where any Dashboard with a matching slug in the given list (and matches all other filters) is returned."),
			),
		),
		Handler: func(session tools.Session, request mcp.CallToolRequest) (*tools.Result, error) {
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

			queryParams := &dashboard.ListDashboardsParams{
				CollectionSlugs: collectionSlugs,

				IncludeDashboardJSON: ptr.To(includeDashboardJson),

				Names: names,

				PageMaxSize: ptr.To(int64(pageMaxSize)),

				PageToken: &pageToken,

				Slugs: slugs,
			}

			api, err := clientProvider.ConfigV1Client(session)
			if err != nil {
				return nil, err
			}
			resp, err := api.Dashboard.ListDashboards(queryParams)
			if err != nil {
				return nil, fmt.Errorf("failed to call ListDashboards: %s", err)
			}
			return &tools.Result{
				JSONContent: resp,
			}, nil
		},
	}
}
