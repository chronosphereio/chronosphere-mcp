// Copyright 2025 Chronosphere Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package configapi

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/chronosphereio/chronosphere-mcp/generated/configv1/configv1"
	dashboardapi "github.com/chronosphereio/chronosphere-mcp/generated/configv1/configv1/dashboard"
	"github.com/chronosphereio/chronosphere-mcp/generated/configv1/models"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools/pkg/params"
)

// UpdateDashboard returns a tool that replaces a dashboard identified by slug.
func UpdateDashboard(api *configv1.ConfigV1API) tools.MCPTool {
	return tools.MCPTool{
		Metadata: tools.NewMetadata("update_dashboard",
			mcp.WithDescription("Replace a dashboard identified by slug. Use dry_run to validate without saving."),
			mcp.WithReadOnlyHintAnnotation(false),
			mcp.WithDestructiveHintAnnotation(true),
			mcp.WithIdempotentHintAnnotation(true),
			mcp.WithOpenWorldHintAnnotation(true),
			mcp.WithString("slug",
				mcp.Description("Slug of the dashboard to update."),
				mcp.Required(),
			),
			mcp.WithObject("dashboard",
				mcp.Description("Complete replacement dashboard. Include all fields that should be retained."),
				mcp.Properties(map[string]any{
					"slug": map[string]any{
						"type":        "string",
						"description": "Dashboard slug. It cannot be changed after creation.",
					},
					"name": map[string]any{
						"type":        "string",
						"description": "Dashboard name.",
					},
					"dashboard_json": map[string]any{
						"type":        "string",
						"description": "Raw JSON string containing the dashboard definition.",
					},
					"collection_slug": map[string]any{
						"type":        "string",
						"description": "Optional slug of the collection containing the dashboard.",
					},
					"collection": map[string]any{
						"type":        "object",
						"description": "Optional collection reference.",
						"properties": map[string]any{
							"slug": map[string]any{
								"type": "string",
							},
							"type": map[string]any{
								"type": "string",
								"enum": []string{"SIMPLE", "SERVICE"},
							},
						},
						"additionalProperties": false,
					},
					"labels": map[string]any{
						"type":        "object",
						"description": "Dashboard labels.",
						"additionalProperties": map[string]any{
							"type": "string",
						},
					},
				}),
				mcp.AdditionalProperties(false),
				mcp.Required(),
			),
			mcp.WithBoolean("dry_run",
				mcp.Description("Validate the dashboard without creating or updating it."),
			),
			mcp.WithBoolean("create_if_missing",
				mcp.Description("Create the dashboard when the path slug does not exist."),
			),
		),
		Write: true,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*tools.Result, error) {
			slug, err := params.String(request, "slug", true, "")
			if err != nil {
				return nil, err
			}
			dashboard, err := params.Object(request, "dashboard", true, models.Configv1Dashboard{})
			if err != nil {
				return nil, err
			}
			if dashboard.Slug != "" && dashboard.Slug != slug {
				return nil, fmt.Errorf("dashboard slug %q must match path slug %q", dashboard.Slug, slug)
			}
			dashboard.Slug = slug
			if dashboard.Name == "" {
				return nil, fmt.Errorf("dashboard.name cannot be empty")
			}
			if dashboard.DashboardJSON == "" {
				return nil, fmt.Errorf("dashboard.dashboard_json cannot be empty")
			}
			dryRun, err := params.Bool(request, "dry_run", false, false)
			if err != nil {
				return nil, err
			}
			createIfMissing, err := params.Bool(request, "create_if_missing", false, false)
			if err != nil {
				return nil, err
			}

			resp, err := api.Dashboard.UpdateDashboard(&dashboardapi.UpdateDashboardParams{
				Context: ctx,
				Slug:    slug,
				Body: &models.ConfigV1UpdateDashboardBody{
					Dashboard:       &dashboard,
					DryRun:          dryRun,
					CreateIfMissing: createIfMissing,
				},
			})
			if err != nil {
				return nil, fmt.Errorf("failed to call UpdateDashboard: %w", err)
			}
			return &tools.Result{
				JSONContent: resp.Payload,
			}, nil
		},
	}
}
