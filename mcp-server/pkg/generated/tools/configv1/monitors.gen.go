package configv1

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"

	"github.com/chronosphereio/mcp-server/generated/configv1/configv1/monitor"
	"github.com/chronosphereio/mcp-server/mcp-server/pkg/client"
	"github.com/chronosphereio/mcp-server/mcp-server/pkg/tools"
	"github.com/chronosphereio/mcp-server/mcp-server/pkg/tools/pkg/params"
	"github.com/chronosphereio/mcp-server/pkg/ptr"
)

func GetMonitor(clientProvider *client.Provider, logger *zap.Logger) tools.MCPTool {
	return tools.MCPTool{
		Metadata: tools.NewMetadata("get_monitor",
			mcp.WithDescription("Get monitors resource"),

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

			queryParams := &monitor.ReadMonitorParams{
				Context: ctx,

				Slug: slug,
			}

			api, err := clientProvider.ConfigV1Client()
			if err != nil {
				return nil, err
			}
			resp, err := api.Monitor.ReadMonitor(queryParams)
			if err != nil {
				return nil, fmt.Errorf("failed to call ReadMonitor: %s", err)
			}
			return &tools.Result{
				JSONContent: resp,
			}, nil
		},
	}
}

func ListMonitors(clientProvider *client.Provider, logger *zap.Logger) tools.MCPTool {
	return tools.MCPTool{
		Metadata: tools.NewMetadata("list_monitors",
			mcp.WithDescription("List monitors resources"),

			mcp.WithArray("bucket_slugs",
				mcp.Description("Filters results by bucket_slug, where any Monitor with a matching bucket_slug in the given list (and matches all other filters) is returned."),
			),

			mcp.WithArray("collection_slugs",
				mcp.Description("Filters results by collection_slug, where any Monitor with a matching collection_slug in the given list (and matches all other filters) is returned."),
			),

			mcp.WithArray("names",
				mcp.Description("Filters results by name, where any Monitor with a matching name in the given list (and matches all other filters) is returned."),
			),

			mcp.WithNumber("page_max_size",
				mcp.Description("Page size preference (i.e. how many items are returned in the next page). If zero, the server will use a default. Regardless of what size is given, clients must never assume how many items will be returned."),
			),

			mcp.WithString("page_token",
				mcp.Description("Opaque page token identifying which page to request. An empty token identifies the first page."),
			),

			mcp.WithArray("slugs",
				mcp.Description("Filters results by slug, where any Monitor with a matching slug in the given list (and matches all other filters) is returned."),
			),

			mcp.WithArray("team_slugs",
				mcp.Description("Filter returned monitors by the teams that own the collections that they belong to."),
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

			teamSlugs, err := params.StringArray(request, "team_slugs", false, nil)
			if err != nil {
				return nil, err
			}

			queryParams := &monitor.ListMonitorsParams{
				Context: ctx,

				BucketSlugs: bucketSlugs,

				CollectionSlugs: collectionSlugs,

				Names: names,

				PageMaxSize: ptr.To(int64(pageMaxSize)),

				PageToken: &pageToken,

				Slugs: slugs,

				TeamSlugs: teamSlugs,
			}

			api, err := clientProvider.ConfigV1Client()
			if err != nil {
				return nil, err
			}
			resp, err := api.Monitor.ListMonitors(queryParams)
			if err != nil {
				return nil, fmt.Errorf("failed to call ListMonitors: %s", err)
			}
			return &tools.Result{
				JSONContent: resp,
			}, nil
		},
	}
}
