package configv1

import (
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"

	"github.com/chronosphereio/mcp-server/generated/configv1/configv1/s_l_o"
	"github.com/chronosphereio/mcp-server/mcp-server/pkg/client"
	"github.com/chronosphereio/mcp-server/mcp-server/pkg/ptr"
	"github.com/chronosphereio/mcp-server/mcp-server/pkg/tools"
	"github.com/chronosphereio/mcp-server/mcp-server/pkg/tools/pkg/params"
)

func GetSlo(clientProvider *client.Provider, logger *zap.Logger) tools.MCPTool {
	return tools.MCPTool{
		Metadata: tools.NewMetadata("get_slo",
			mcp.WithDescription("Get slos resource"),

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

			queryParams := &s_l_o.ReadSLOParams{
				Slug: slug,
			}

			api, err := clientProvider.ConfigV1Client(session)
			if err != nil {
				return nil, err
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

func ListSlos(clientProvider *client.Provider, logger *zap.Logger) tools.MCPTool {
	return tools.MCPTool{
		Metadata: tools.NewMetadata("list_slos",
			mcp.WithDescription("List slos resources"),

			mcp.WithArray("collection_slugs",
				mcp.Description(""),
			),

			mcp.WithArray("names",
				mcp.Description("Filters results by name, where any SLO with a matching name in the given list (and matches all other filters) is returned."),
			),

			mcp.WithNumber("page.max_size",
				mcp.Description("Page size preference (i.e. how many items are returned in the next page). If zero, the server will use a default. Regardless of what size is given, clients must never assume how many items will be returned."),
			),

			mcp.WithString("page.token",
				mcp.Description("Opaque page token identifying which page to request. An empty token identifies the first page."),
			),

			mcp.WithArray("service_slugs",
				mcp.Description(""),
			),

			mcp.WithArray("slugs",
				mcp.Description("Filters results by slug, where any SLO with a matching slug in the given list (and matches all other filters) is returned."),
			),
		),
		Handler: func(session tools.Session, request mcp.CallToolRequest) (*tools.Result, error) {
			collectionSlugs, err := params.StringArray(request, "collection_slugs", false, nil)
			if err != nil {
				return nil, err
			}

			names, err := params.StringArray(request, "names", false, nil)
			if err != nil {
				return nil, err
			}

			pageMaxSize, err := params.Int(request, "page.max_size", false, 0)
			if err != nil {
				return nil, err
			}

			pageToken, err := params.String(request, "page.token", false, "")
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
				CollectionSlugs: collectionSlugs,

				Names: names,

				PageMaxSize: ptr.To(int64(pageMaxSize)),

				PageToken: &pageToken,

				ServiceSlugs: serviceSlugs,

				Slugs: slugs,
			}

			api, err := clientProvider.ConfigV1Client(session)
			if err != nil {
				return nil, err
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
