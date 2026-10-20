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

			params.WithStringArray("matches_series",
				mcp.Description("Restricts the returned drop rules to those whose `filters` would accept a series carrying these labels. Each entry is `name=value`, for example `__name__=http_requests_total`. A rule matches when every supplied label satisfies the rule's glob filter for that label, so a rule with no filter on a label accepts any value of it. Filters on labels not supplied here are ignored. Values are concrete: a value containing glob characters such as `*` or `{a,b}` is rejected. Currently only the metric name may be supplied, as a single `__name__=<metric>` entry; other labels are rejected."),
			),

			params.WithStringArray("names",
				mcp.Description("Filters results by name, where any DropRule with a matching name in the given list (and matches all other filters) will be returned."),
			),

			mcp.WithNumber("page_max_size",
				mcp.Description("Sets the preferred number of items to return per page. If set to `0`, the server will use its default value. Regardless of the value specified, clients must never assume how many items will be returned."),
			),

			mcp.WithString("page_token",
				mcp.Description("An opaque page token that identifies which page the client should request. An empty value indicates the first page."),
			),

			params.WithStringArray("slugs",
				mcp.Description("Filters results by slug, where any DropRule with a matching slug in the given list (and matches all other filters) will be returned."),
			),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*tools.Result, error) {
			matchesSeries, err := params.StringArray(request, "matches_series", false, nil)
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

			queryParams := &drop_rule.ListDropRulesParams{
				Context: ctx,

				MatchesSeries: matchesSeries,

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
