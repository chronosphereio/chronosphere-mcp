package configv1

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"

	"github.com/chronosphereio/chronosphere-mcp/generated/configv1/configv1"
	"github.com/chronosphereio/chronosphere-mcp/generated/configv1/configv1/recording_rule"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools/pkg/params"
	"github.com/chronosphereio/chronosphere-mcp/pkg/ptr"
)

func GetRecordingRule(api *configv1.ConfigV1API, logger *zap.Logger) tools.MCPTool {
	return tools.MCPTool{
		Metadata: tools.NewMetadata("get_recording_rule",
			mcp.WithDescription("Get recording-rules resource"),

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

			queryParams := &recording_rule.ReadRecordingRuleParams{
				Context: ctx,

				Slug: slug,
			}

			resp, err := api.RecordingRule.ReadRecordingRule(queryParams)
			if err != nil {
				return nil, fmt.Errorf("failed to call ReadRecordingRule: %s", err)
			}
			return &tools.Result{
				JSONContent: resp,
			}, nil
		},
	}
}

func ListRecordingRules(api *configv1.ConfigV1API, logger *zap.Logger) tools.MCPTool {
	return tools.MCPTool{
		Metadata: tools.NewMetadata("list_recording_rules",
			mcp.WithDescription("List recording-rules resources"),

			params.WithStringArray("bucket_slugs",
				mcp.Description("The execution_groups filter cannot be used when a bucket_slug filter is provided."),
			),

			params.WithStringArray("execution_groups",
				mcp.Description("The bucket_slugs filter cannot be used when an execution_group filter is provided."),
			),

			params.WithStringArray("names",
				mcp.Description("Filters results by name, where any RecordingRule with a matching name in the given list (and matches all other filters) is returned."),
			),

			mcp.WithNumber("page_max_size",
				mcp.Description("Page size preference (i.e. how many items are returned in the next page). If zero, the server will use a default. Regardless of what size is given, clients must never assume how many items will be returned."),
			),

			mcp.WithString("page_token",
				mcp.Description("Opaque page token identifying which page to request. An empty token identifies the first page."),
			),

			params.WithStringArray("slugs",
				mcp.Description("Filters results by slug, where any RecordingRule with a matching slug in the given list (and matches all other filters) is returned."),
			),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*tools.Result, error) {
			bucketSlugs, err := params.StringArray(request, "bucket_slugs", false, nil)
			if err != nil {
				return nil, err
			}

			executionGroups, err := params.StringArray(request, "execution_groups", false, nil)
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

			queryParams := &recording_rule.ListRecordingRulesParams{
				Context: ctx,

				BucketSlugs: bucketSlugs,

				ExecutionGroups: executionGroups,

				Names: names,

				PageMaxSize: ptr.To(int64(pageMaxSize)),

				PageToken: &pageToken,

				Slugs: slugs,
			}

			resp, err := api.RecordingRule.ListRecordingRules(queryParams)
			if err != nil {
				return nil, fmt.Errorf("failed to call ListRecordingRules: %s", err)
			}
			return &tools.Result{
				JSONContent: resp,
			}, nil
		},
	}
}
