package configv1

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"

	"github.com/chronosphereio/chronosphere-mcp/generated/configv1/configv1"
	"github.com/chronosphereio/chronosphere-mcp/generated/configv1/configv1/drop_rule"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools/pkg/params"
	"github.com/chronosphereio/chronosphere-mcp/pkg/ptr"
)

func GetDropRule(api *configv1.ConfigV1API, logger *zap.Logger) tools.MCPTool {
	return tools.MCPTool{
		Metadata: tools.NewMetadata("get_drop_rule",
			mcp.WithDescription("Get drop-rules resource"),

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

			queryParams := &drop_rule.ReadDropRuleParams{
				Context: ctx,

				Slug: slug,
			}

			resp, err := api.DropRule.ReadDropRule(queryParams)
			if err != nil {
				return nil, fmt.Errorf("failed to call ReadDropRule: %s", err)
			}
			return &tools.Result{
				JSONContent: resp,
			}, nil
		},
	}
}

func ListDropRules(api *configv1.ConfigV1API, logger *zap.Logger) tools.MCPTool {
	return tools.MCPTool{
		Metadata: tools.NewMetadata("list_drop_rules",
			mcp.WithDescription("List drop-rules resources"),

			params.WithStringArray("names",
				mcp.Description("Filters results by name, where any DropRule with a matching name in the given list (and matches all other filters) is returned."),
			),

			mcp.WithNumber("page_max_size",
				mcp.Description("Page size preference (i.e. how many items are returned in the next page). If zero, the server will use a default. Regardless of what size is given, clients must never assume how many items will be returned."),
			),

			mcp.WithString("page_token",
				mcp.Description("Opaque page token identifying which page to request. An empty token identifies the first page."),
			),

			params.WithStringArray("slugs",
				mcp.Description("Filters results by slug, where any DropRule with a matching slug in the given list (and matches all other filters) is returned."),
			),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*tools.Result, error) {
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

			queryParams := &drop_rule.ListDropRulesParams{
				Context: ctx,

				Names: names,

				PageMaxSize: ptr.To(int64(pageMaxSize)),

				PageToken: &pageToken,

				Slugs: slugs,
			}

			resp, err := api.DropRule.ListDropRules(queryParams)
			if err != nil {
				return nil, fmt.Errorf("failed to call ListDropRules: %s", err)
			}
			return &tools.Result{
				JSONContent: resp,
			}, nil
		},
	}
}
