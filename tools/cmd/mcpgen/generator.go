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

package main

import (
	"bytes"
	_ "embed"
	"log"
	"regexp"
	"strings"
	"text/template"

	"github.com/chronosphereio/chronoctl-core/tools/pkg/clispec"
	"github.com/go-openapi/inflect"
)

type generatedFile struct {
	name    string
	content string
}

type ParameterSpec struct {
	Name        string
	GoName      string
	Description string
	// MCPParamFunc is the type for the mcp framework
	MCPParamFunc string
	ParseType    string
	DefaultValue string
	Required     bool
	// Field name of the swagger generated Go type.
	SwaggerGoFieldName string
	// Wither this is a scalar type which must be passed as a pointer.
	IsScalar bool
}

type ToolSpec struct {
	CamelName   string
	SnakeName   string
	Description string
	Verb        string
	Parameters  []ParameterSpec
	// APIParam is the Go type name of the generated API parameter.
	APIParam string
}

type Entity struct {
	Name      string
	PkgName   string
	ToolSpecs []ToolSpec
}

type entitySpec struct {
	Entity  Entity
	PkgName string
}

func generate(pkgName string, spec *clispec.Spec) ([]generatedFile, error) {
	var files []generatedFile
	specs := make(map[string]entitySpec, len(spec.Entities))
	allowList := entityAllowList()
	for _, entity := range spec.Entities {
		if len(allowList) > 0 {
			if _, ok := allowList[entity.Name]; !ok {
				continue
			}
		}
		specs[entity.Name] = entitySpec{
			Entity:  convertEntity(entity),
			PkgName: pkgName,
		}
	}

	for name, entity := range specs {
		contents, err := generateTool(entity)
		if err != nil {
			return nil, err
		}
		files = append(files, generatedFile{
			name:    name,
			content: contents,
		})
	}

	return files, nil
}

func convertEntity(entity *clispec.Entity) Entity {
	ent := Entity{
		Name:    camelCase(inflect.Singularize(entity.Name)),
		PkgName: entityPkgMap(inflect.Underscore(inflect.Singularize(entity.Name))),
	}

	ent.ToolSpecs = append(ent.ToolSpecs, convertToolSpec(inflect.Singularize(entity.Name), "Read", entity.Get))

	if entity.IsNotSingleton() {
		ent.ToolSpecs = append(ent.ToolSpecs, convertToolSpec(entity.Name, "List", entity.List))
	}
	return ent
}

// entityPkgMap maps special entity names to their actual package names in the generated swagger go code.
func entityPkgMap(entityPkg string) string {
	switch entityPkg {
	case "slo":
		return "s_l_o"
	}
	return entityPkg
}

func convertToolSpec(entityName string, action string, command *clispec.Command) ToolSpec {
	spec := ToolSpec{
		CamelName:   camelCase(command.Name) + camelCase(entityName),
		SnakeName:   inflect.Underscore(command.Name + "-" + entityName),
		Description: command.Description,
		Verb:        command.Verb,
		Parameters:  make([]ParameterSpec, 0, len(command.Parameters)),
		APIParam:    action + acronymReplace(camelCase(entityName)),
	}
	for _, param := range command.Parameters {
		spec.Parameters = append(spec.Parameters, ParameterSpec{
			Name:               inflect.Underscore(strings.Replace(param.Name, ".", "_", -1)),
			GoName:             uncapitalize(camelCase(param.Name)),
			SwaggerGoFieldName: acronymReplace(camelCase(param.Name)),
			Description:        param.Description,
			MCPParamFunc:       goTypeToMCPParamFunc(param.GoType),
			ParseType:          goTypeToParamParseType(param.GoType),
			DefaultValue:       getDefaultValue(param.GoType),
			Required:           param.Required,
			IsScalar:           isScalar(param.GoType),
		})
	}
	return spec
}

func isScalar(goType string) bool {
	switch goType {
	case "string", "bool", "int", "float":
		return true
	}
	return false
}

func getDefaultValue(goType string) string {
	switch goType {
	case "string":
		return `""`
	case "bool":
		return "false"
	case "int":
		return "0"
	case "[]string":
		return "nil"
	}
	log.Fatalf("unsupported type: %s", goType)
	return goType
}

// goTypeToParamParseType converts a Go type to the corresponding param parse function name.
func goTypeToParamParseType(goType string) string {
	switch goType {
	case "string":
		return "String"
	case "bool":
		return "Bool"
	case "int":
		return "Int"
	case "[]string":
		return "StringArray"
	}
	log.Fatalf("unsupported type: %s", goType)
	return goType
}

// goTypeToMCPParamFunc converts a Go type to the corresponding mcp.With${type}.
func goTypeToMCPParamFunc(goType string) string {
	switch goType {
	case "string":
		return "mcp.WithString"
	case "bool":
		return "mcp.WithBoolean"
	case "int":
		return "mcp.WithNumber"
	case "[]string":
		return "params.WithStringArray"
	}

	log.Fatalf("unsupported type: %s", goType)
	return goType
}

var (
	//go:embed tmpl/entity_tool.go.tmpl
	_entityToolTmpl     string
	_entityToolTemplate = template.Must(template.New("entity_tool").Parse(_entityToolTmpl))
)

func generateTool(entity entitySpec) (string, error) {
	var buffer bytes.Buffer
	if err := _entityToolTemplate.Execute(&buffer, entity); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func entityAllowList() map[string]any {
	m := map[string]any{}
	if *allowedEntities == "" {
		return m
	}
	for _, e := range strings.Split(*allowedEntities, ",") {
		m[e] = true
	}
	return m
}

func acronymReplace(str string) string {
	replaces := map[string]string{
		"Json": "JSON",
		"Id":   "ID",
		"Ui":   "UI",
		"Slo":  "SLO",
		"Slos": "SLOs",
	}
	for suffix, replace := range replaces {
		if strings.HasSuffix(str, suffix) {
			return strings.TrimSuffix(str, suffix) + replace
		}
	}
	return str
}

// camelCase capitalizes the first letter of each word in a string separated by -.
func camelCase(str string) string {
	re := regexp.MustCompile("[-_.]")

	parts := re.Split(str, -1)
	for i, part := range parts {
		parts[i] = inflect.Capitalize(part)
	}
	return strings.Join(parts, "")
}

func uncapitalize(str string) string {
	return strings.ToLower(str[:1]) + str[1:]
}
