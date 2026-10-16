package configv1

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"

	"github.com/chronosphereio/chronosphere-mcp/generated/configv1/configv1"
	"github.com/chronosphereio/chronosphere-mcp/generated/configv1/configv1/rollup_rule"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools/pkg/params"
	"github.com/chronosphereio/chronosphere-mcp/pkg/ptr"
)

func GetRollupRule(api *configv1.ConfigV1API, logger *zap.Logger) tools.MCPTool {
	return tools.MCPTool{
		Metadata: tools.NewMetadata("get_rollup_rule",
			mcp.WithDescription("Get rollup-rules resource"),

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

			queryParams := &rollup_rule.ReadRollupRuleParams{
				Context: ctx,

				Slug: slug,
			}

			resp, err := api.RollupRule.ReadRollupRule(queryParams)
			if err != nil {
				return nil, fmt.Errorf("failed to call ReadRollupRule: %s", err)
			}
			return &tools.Result{
				JSONContent: resp,
			}, nil
		},
	}
}

func ListRollupRules(api *configv1.ConfigV1API, logger *zap.Logger) tools.MCPTool {
	return tools.MCPTool{
		Metadata: tools.NewMetadata("list_rollup_rules",
			mcp.WithDescription("List rollup-rules resources"),

			params.WithStringArray("bucket_slugs",
				mcp.Description("Filters results by bucket_slug, where any RollupRule with a matching bucket_slug in the given list (and matches all other filters) will be returned."),
			),

			params.WithStringArray("matches_series",
				mcp.Description("Restricts the returned rollup rules to those whose `filters` would accept a series carrying these labels. Each entry is `name=value`, for example `__name__=http_requests_total`. A rule matches when every supplied label satisfies the rule's glob filter for that label, so a rule with no filter on a label accepts any value of it. Filters on labels not supplied here are ignored. Values are concrete: a value containing glob characters such as `*` or `{a,b}` is rejected. Currently only the metric name may be supplied, as a single `__name__=<metric>` entry; other labels are rejected."),
			),

			params.WithStringArray("names",
				mcp.Description("Filters results by name, where any RollupRule with a matching name in the given list (and matches all other filters) will be returned."),
			),

			mcp.WithNumber("page_max_size",
				mcp.Description("Sets the preferred number of items to return per page. If set to `0`, the server will use its default value. Regardless of the value specified, clients must never assume how many items will be returned."),
			),

			mcp.WithString("page_token",
				mcp.Description("An opaque page token that identifies which page the client should request. An empty value indicates the first page."),
			),

			params.WithStringArray("slugs",
				mcp.Description("Filters results by slug, where any RollupRule with a matching slug in the given list (and matches all other filters) will be returned."),
			),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*tools.Result, error) {
			bucketSlugs, err := params.StringArray(request, "bucket_slugs", false, nil)
			if err != nil {
				return nil, err
			}

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

			queryParams := &rollup_rule.ListRollupRulesParams{
				Context: ctx,

				BucketSlugs: bucketSlugs,

				MatchesSeries: matchesSeries,

				Names: names,

				PageMaxSize: ptr.To(int64(pageMaxSize)),

				PageToken: &pageToken,

				Slugs: slugs,
			}

			resp, err := api.RollupRule.ListRollupRules(queryParams)
			if err != nil {
				return nil, fmt.Errorf("failed to call ListRollupRules: %s", err)
			}
			return &tools.Result{
				JSONContent: resp,
			}, nil
		},
	}
}
