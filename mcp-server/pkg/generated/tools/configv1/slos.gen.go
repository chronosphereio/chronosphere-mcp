package configv1

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"

	"github.com/chronosphereio/chronosphere-mcp/generated/configv1/configv1"
	"github.com/chronosphereio/chronosphere-mcp/generated/configv1/configv1/s_l_o"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools/pkg/params"
	"github.com/chronosphereio/chronosphere-mcp/pkg/ptr"
)

func GetSlo(api *configv1.ConfigV1API, logger *zap.Logger) tools.MCPTool {
	return tools.MCPTool{
		Metadata: tools.NewMetadata("get_slo",
			mcp.WithDescription("Get slos resource"),

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

			queryParams := &s_l_o.ReadSLOParams{
				Context: ctx,

				Slug: slug,
			}

			resp, err := api.Slo.ReadSLO(queryParams)
			if err != nil {
				return nil, fmt.Errorf("failed to call ReadSLO: %s", err)
			}
			return &tools.Result{
				JSONContent: resp,
			}, nil
		},
	}
}

func ListSlos(api *configv1.ConfigV1API, logger *zap.Logger) tools.MCPTool {
	return tools.MCPTool{
		Metadata: tools.NewMetadata("list_slos",
			mcp.WithDescription("List slos resources"),

			params.WithStringArray("collection_slugs",
				mcp.Description(""),
			),

			params.WithStringArray("names",
				mcp.Description("Filters results by name, where any SLO with a matching name in the given list (and matches all other filters) is returned."),
			),

			mcp.WithNumber("page_max_size",
				mcp.Description("Page size preference (i.e. how many items are returned in the next page). If zero, the server will use a default. Regardless of what size is given, clients must never assume how many items will be returned."),
			),

			mcp.WithString("page_token",
				mcp.Description("Opaque page token identifying which page to request. An empty token identifies the first page."),
			),

			params.WithStringArray("service_slugs",
				mcp.Description(""),
			),

			params.WithStringArray("slugs",
				mcp.Description("Filters results by slug, where any SLO with a matching slug in the given list (and matches all other filters) is returned."),
			),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*tools.Result, error) {
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

			serviceSlugs, err := params.StringArray(request, "service_slugs", false, nil)
			if err != nil {
				return nil, err
			}

			slugs, err := params.StringArray(request, "slugs", false, nil)
			if err != nil {
				return nil, err
			}

			queryParams := &s_l_o.ListSLOsParams{
				Context: ctx,

				CollectionSlugs: collectionSlugs,

				Names: names,

				PageMaxSize: ptr.To(int64(pageMaxSize)),

				PageToken: &pageToken,

				ServiceSlugs: serviceSlugs,

				Slugs: slugs,
			}

			resp, err := api.Slo.ListSLOs(queryParams)
			if err != nil {
				return nil, fmt.Errorf("failed to call ListSLOs: %s", err)
			}
			return &tools.Result{
				JSONContent: resp,
			}, nil
		},
	}
}
