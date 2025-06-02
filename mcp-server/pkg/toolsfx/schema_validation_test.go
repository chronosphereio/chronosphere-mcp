// Copyright 2025 Chronosphere Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package toolsfx_test

import (
	"encoding/json"
	"testing"

	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"

	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/clientfx"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/toolsfx"
	"github.com/chronosphereio/chronosphere-mcp/pkg/links"
)

// getAllTools returns all tools from all tool packages for testing using fx dependency injection
func getAllTools(t *testing.T) []tools.MCPTool {
	var allTools []tools.MCPTool

	// Create test app with fx to initialize tools the same way as the MCP server
	app := fxtest.New(t,
		// Provide test logger
		fx.Provide(func() *zap.Logger {
			return zap.NewNop()
		}),
		// Provide test API config (with empty values for testing)
		fx.Provide(func() *clientfx.ChronosphereConfig {
			return &clientfx.ChronosphereConfig{
				APIURL:      "https://test.chronosphere.io",
				APIToken:    "test-token",
				LogscaleURL: "https://test.logs.chronosphere.io",
			}
		}),
		// Provide links builder
		fx.Provide(func() *links.Builder {
			return &links.Builder{}
		}),
		// Include client module for API clients
		clientfx.Module,
		// Include tools module to get all tools via fx
		toolsfx.Module,
		// Populate the tools into our test variable
		fx.Invoke(func(p struct {
			fx.In
			ToolGroups []tools.MCPTools `group:"mcp_tools"`
		}) {
			// Flatten all tool groups into a single slice
			for _, group := range p.ToolGroups {
				allTools = append(allTools, group.MCPTools()...)
			}
		}),
	)

	// Start the app to initialize dependencies
	app.RequireStart()
	t.Cleanup(func() {
		app.RequireStop()
	})

	return allTools
}

func TestToolParameterSchemaValidation(t *testing.T) {
	allTools := getAllTools(t)

	require.NotEmpty(t, allTools, "No tools found for testing")

	for _, tool := range allTools {
		t.Run(tool.Metadata.Name, func(t *testing.T) {
			validateToolSchema(t, tool)
		})
	}
}

func validateToolSchema(t *testing.T, tool tools.MCPTool) {
	// Check if we have either InputSchema or RawInputSchema
	var schema spec.Schema
	var err error

	if tool.Metadata.RawInputSchema != nil {
		err = json.Unmarshal(tool.Metadata.RawInputSchema, &schema)
		require.NoError(t, err, "Tool %s has invalid RawInputSchema JSON", tool.Metadata.Name)
	} else {
		// Convert InputSchema to JSON and back to get a generic map
		schemaBytes, err := json.Marshal(tool.Metadata.InputSchema)
		require.NoError(t, err, "Tool %s InputSchema cannot be marshaled to JSON", tool.Metadata.Name)

		err = json.Unmarshal(schemaBytes, &schema)
		require.NoError(t, err, "Tool %s InputSchema cannot be unmarshaled from JSON", tool.Metadata.Name)
	}

	// Validate the schema
	validateJSONSchema(t, tool.Metadata.Name, &schema, "")
}

func validateJSONSchema(t *testing.T, toolName string, schema *spec.Schema, path string) {
	// Basic schema validation
	assert.NotNil(t, schema, "Tool %s: Schema at path '%s' is nil", toolName, path)
	if schema == nil {
		return
	}

	if path == "" && len(schema.Properties) == 0 {
		// Some tools may legitimately have no parameters (e.g., metadata endpoints)
		t.Logf("Tool %s: Root schema has no properties defined (this may be intentional)", toolName)
	}
	if len(schema.Type) == 0 {
		assert.Fail(t, "Tool %s: Schema type at path '%s' is empty", toolName, path)
	}
	if len(schema.Type) > 1 {
		assert.Fail(t, "Tool %s: Schema with multiple types not suppported at path '%s': %v", toolName, path, schema.Type)
	}

	switch schema.Type[0] {
	case "array":
		validateArraySchema(t, toolName, schema, path)
	case "object":
		validateObjectSchema(t, toolName, schema, path)
	case "string", "number", "integer", "boolean":
		// Basic types are fine
	default:
		assert.Fail(t, "Unknown schema type", "Tool %s: Schema at path '%s' has unknown type: %s", toolName, path, schema.Type[0])
	}

	// Recursively validate nested schemas
	for propName, propSchema := range schema.Properties {
		newPath := path
		if newPath != "" {
			newPath += "."
		}
		newPath += propName
		validateJSONSchema(t, toolName, &propSchema, newPath)
	}

	// Validate array items if present
	if schema.Items != nil && schema.Items.Schema != nil {
		newPath := path + ".items"
		validateJSONSchema(t, toolName, schema.Items.Schema, newPath)
	}
}

func validateArraySchema(t *testing.T, toolName string, schema *spec.Schema, path string) {
	// Arrays must have an 'items' property
	require.NotNil(t, schema.Items, "Tool %s: Items must be defined for array at path %q", toolName, path)
	require.NotNil(t, schema.Items.Schema, "Tool %s: Items schema must be defined for array at path %q", toolName, path)
	require.NotEmpty(t, schema.Items.Schema.Type, "Tool %s: Items schema must have a type defined for array at path %q", toolName, path)

	// Perform detailed array validation (previously in validateArraySchemaDetailed)
	items := schema.Items.Schema
	itemTypeStr := items.Type[0]

	// Validate based on item type
	switch itemTypeStr {
	case "string":
		// String arrays are common and should be well-defined
		t.Logf("Tool %s: Found string array at path '%s' ✓", toolName, path)
	case "object":
		// Object arrays should have properties defined
		if len(items.Properties) == 0 {
			assert.Fail(t, "Object array missing properties", "Tool %s: Object array items at path '%s' have no properties defined", toolName, path)
		} else {
			t.Logf("Tool %s: Found object array at path '%s' ✓", toolName, path)
		}
	case "number", "integer", "boolean":
		t.Logf("Tool %s: Found %s array at path '%s'", toolName, itemTypeStr, path)
	default:
		assert.Fail(t, "Unknown array item type", "Tool %s: Array items at path '%s' have unknown type: %s", toolName, path, itemTypeStr)
	}
}

func validateObjectSchema(t *testing.T, toolName string, schema *spec.Schema, path string) {
	// Objects should have properties (though it's not strictly required)
	if len(schema.Properties) == 0 {
		// This might be OK for some schemas, but let's warn
		t.Logf("Tool %s: Object schema at path '%s' has no properties defined", toolName, path)
		return
	}

	// Check that each property is a valid schema
	for propName, propSchema := range schema.Properties {
		// Each property should have a type
		if len(propSchema.Type) == 0 {
			assert.Fail(t, "Tool %s: Object property '%s' at path '%s' is missing 'type' field", toolName, propName, path)
		}
	}
}
