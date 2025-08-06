// Package main generates a markdown table of available MCP tools by parsing Go source files.
package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Tool struct {
	Name        string
	Description string
	GroupName   string
}

func main() {
	// Find all tool packages
	toolsDir := "mcp-server/pkg/tools"
	generatedToolsDir := "mcp-server/pkg/generated/tools"

	var tools []Tool

	// Parse regular tools
	err := filepath.Walk(toolsDir, func(path string, _ os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip test files and the main tools.go file at the root level
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") || strings.HasSuffix(path, "pkg/tools/tools.go") {
			return nil
		}

		// Parse the Go file
		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			log.Printf("Error parsing %s: %v", path, err)
			return nil
		}

		// Extract tools from this file
		fileTools := extractToolsFromFile(node, path)
		tools = append(tools, fileTools...)

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	// Parse generated tools
	err = filepath.Walk(generatedToolsDir, func(path string, _ os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only process .gen.go files
		if !strings.HasSuffix(path, ".gen.go") {
			return nil
		}

		// Parse the Go file
		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			log.Printf("Error parsing %s: %v", path, err)
			return nil
		}

		// Extract tools from generated files
		generatedTools := extractGeneratedTools(node, path)
		tools = append(tools, generatedTools...)

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	// Sort tools by group name, then by tool name
	sort.Slice(tools, func(i, j int) bool {
		if tools[i].GroupName != tools[j].GroupName {
			return tools[i].GroupName < tools[j].GroupName
		}
		return tools[i].Name < tools[j].Name
	})

	// Generate markdown table
	generateMarkdownTable(tools)
}

func extractToolsFromFile(node *ast.File, _ string) []Tool {
	var tools []Tool
	var groupName string

	// First, find the GroupName method to get the group name
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			if x.Name.Name == "GroupName" && x.Recv != nil {
				// Find the return statement
				for _, stmt := range x.Body.List {
					if retStmt, ok := stmt.(*ast.ReturnStmt); ok {
						if len(retStmt.Results) > 0 {
							if lit, ok := retStmt.Results[0].(*ast.BasicLit); ok && lit.Kind == token.STRING {
								var err error
								groupName, err = strconv.Unquote(lit.Value)
								if err != nil {
									log.Printf("Error unquoting group name: %v", err)
								}
							}
						}
					}
				}
			}
		}
		return true
	})

	// Then find the MCPTools method and extract tool definitions
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			if x.Name.Name == "MCPTools" && x.Recv != nil {
				tools = append(tools, extractToolsFromMCPToolsMethod(x, groupName)...)
			}
		}
		return true
	})

	return tools
}

func extractToolsFromMCPToolsMethod(funcDecl *ast.FuncDecl, groupName string) []Tool {
	var tools []Tool

	// Look for return statements with composite literals
	ast.Inspect(funcDecl, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.ReturnStmt:
			for _, result := range x.Results {
				if compLit, ok := result.(*ast.CompositeLit); ok {
					// This should be []tools.MCPTool
					for _, elt := range compLit.Elts {
						if toolLit, ok := elt.(*ast.CompositeLit); ok {
							tool := extractToolFromCompositeLiteral(toolLit, groupName)
							if tool.Name != "" {
								tools = append(tools, tool)
							}
						}
					}
				}
			}
		}
		return true
	})

	return tools
}

func extractToolFromCompositeLiteral(compLit *ast.CompositeLit, groupName string) Tool {
	var tool Tool
	tool.GroupName = groupName

	// Look for Metadata field
	for _, elt := range compLit.Elts {
		if kv, ok := elt.(*ast.KeyValueExpr); ok {
			if ident, ok := kv.Key.(*ast.Ident); ok && ident.Name == "Metadata" {
				// Extract metadata
				tool = extractMetadata(kv.Value, groupName)
			}
		}
	}

	return tool
}

func extractMetadata(expr ast.Expr, groupName string) Tool {
	var tool Tool
	tool.GroupName = groupName

	// Look for tools.NewMetadata call
	if callExpr, ok := expr.(*ast.CallExpr); ok {
		if len(callExpr.Args) > 0 {
			// First argument should be the tool name
			if lit, ok := callExpr.Args[0].(*ast.BasicLit); ok && lit.Kind == token.STRING {
				var err error
				tool.Name, err = strconv.Unquote(lit.Value)
				if err != nil {
					log.Printf("Error unquoting tool name: %v", err)
				}
			}

			// Look for mcp.WithDescription in the arguments
			for _, arg := range callExpr.Args {
				if innerCall, ok := arg.(*ast.CallExpr); ok {
					if sel, ok := innerCall.Fun.(*ast.SelectorExpr); ok {
						if sel.Sel.Name == "WithDescription" && len(innerCall.Args) > 0 {
							if lit, ok := innerCall.Args[0].(*ast.BasicLit); ok && lit.Kind == token.STRING {
								desc, err := strconv.Unquote(lit.Value)
								if err != nil {
									log.Printf("Error unquoting description: %v", err)
									continue
								}
								// Clean up the description - remove extra whitespace and newlines
								desc = cleanDescription(desc)
								tool.Description = desc
							}
						}
					}
				}
			}
		}
	}

	return tool
}

func cleanDescription(desc string) string {
	// Remove excessive whitespace and normalize line breaks
	desc = regexp.MustCompile(`\s+`).ReplaceAllString(desc, " ")
	desc = strings.TrimSpace(desc)

	// Truncate very long descriptions and add ellipsis
	if len(desc) > 200 {
		desc = desc[:197] + "..."
	}

	return desc
}

func extractGeneratedTools(node *ast.File, _ string) []Tool {
	var tools []Tool

	// For generated tools, look for function declarations that return tools.MCPTool
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			// Check if this function returns tools.MCPTool
			if x.Type.Results != nil && len(x.Type.Results.List) == 1 {
				if sel, ok := x.Type.Results.List[0].Type.(*ast.SelectorExpr); ok {
					if ident, ok := sel.X.(*ast.Ident); ok && ident.Name == "tools" && sel.Sel.Name == "MCPTool" {
						tool := extractToolFromGeneratedFunction(x)
						if tool.Name != "" {
							tools = append(tools, tool)
						}
					}
				}
			}
		}
		return true
	})

	return tools
}

func extractToolFromGeneratedFunction(funcDecl *ast.FuncDecl) Tool {
	var tool Tool
	tool.GroupName = "configapi" // All generated tools are configapi tools

	// Look for return statements with tools.MCPTool struct
	ast.Inspect(funcDecl, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.ReturnStmt:
			for _, result := range x.Results {
				if compLit, ok := result.(*ast.CompositeLit); ok {
					// This should be tools.MCPTool
					tool = extractToolFromCompositeLiteral(compLit, "configapi")
				}
			}
		}
		return true
	})

	return tool
}

func generateMarkdownTable(tools []Tool) {
	writeString("## Available MCP Tools\n")
	writeString("\n")
	writeString("| Group | Tool Name | Description |\n")
	writeString("|-------|-----------|-------------|\n")

	for _, tool := range tools {
		// Escape pipe characters in description
		description := strings.ReplaceAll(tool.Description, "|", "\\|")
		writeString("| " + tool.GroupName + " | " + tool.Name + " | " + description + " |\n")
	}
	writeString("\n")
}

func writeString(s string) {
	if _, err := os.Stdout.WriteString(s); err != nil {
		log.Printf("Error writing to stdout: %v", err)
	}
}
