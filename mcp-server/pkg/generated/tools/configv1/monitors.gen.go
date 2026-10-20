package configv1

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"

	"github.com/chronosphereio/chronosphere-mcp/generated/configv1/configv1"
	"github.com/chronosphereio/chronosphere-mcp/generated/configv1/configv1/monitor"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools/pkg/params"
	"github.com/chronosphereio/chronosphere-mcp/pkg/ptr"
)

func GetMonitor(api *configv1.ConfigV1API, logger *zap.Logger) tools.MCPTool {
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

func ListMonitors(api *configv1.ConfigV1API, logger *zap.Logger) tools.MCPTool {
	return tools.MCPTool{
		Metadata: tools.NewMetadata("list_monitors",
			mcp.WithDescription("List monitors resources"),

			params.WithStringArray("bucket_slugs",
				mcp.Description("Filters results by bucket_slug, where any Monitor with a matching bucket_slug in the given list (and matches all other filters) will be returned."),
			),

			params.WithStringArray("collection_slugs",
				mcp.Description("Filters results by collection_slug, where any Monitor with a matching collection_slug in the given list (and matches all other filters) will be returned."),
			),

			params.WithStringArray("label_matchers",
				mcp.Description("Filter returned monitors by their labels. Each entry is `name=value` or `name!=value`, for example `team=payments` or `team!=payments`; the value may be double-quoted. `=~` and `!~` are not accepted. A monitor is returned only if every entry holds for it: `name=value` requires the monitor to carry that label with that value (it may carry other labels too), and `name!=value` excludes monitors whose label has that value. With an empty value, `name=` matches monitors that do not have the label and `name!=` matches monitors that have it with any value. Unlike the other list filters, whose values are OR'd, entries here are AND'd."),
			),

			mcp.WithString("name_contains",
				mcp.Description("Filter that returns monitors whose name contains this string. The filter is not case sensitive."),
			),

			params.WithStringArray("names",
				mcp.Description("Filters results by name, where any Monitor with a matching name in the given list (and matches all other filters) will be returned."),
			),

			mcp.WithNumber("page_max_size",
				mcp.Description("Sets the preferred number of items to return per page. If set to `0`, the server will use its default value. Regardless of the value specified, clients must never assume how many items will be returned."),
			),

			mcp.WithString("page_token",
				mcp.Description("An opaque page token that identifies which page the client should request. An empty value indicates the first page."),
			),

			params.WithStringArray("slugs",
				mcp.Description("Filters results by slug, where any Monitor with a matching slug in the given list (and matches all other filters) will be returned."),
			),

			params.WithStringArray("team_slugs",
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

			labelMatchers, err := params.StringArray(request, "label_matchers", false, nil)
			if err != nil {
				return nil, err
			}

			nameContains, err := params.String(request, "name_contains", false, "")
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

				LabelMatchers: labelMatchers,

				NameContains: ptr.To(nameContains),

				Names: names,

				PageMaxSize: ptr.To(int64(pageMaxSize)),

				PageToken: &pageToken,

				Slugs: slugs,

				TeamSlugs: teamSlugs,
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
