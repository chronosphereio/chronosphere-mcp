// Copyright 2026 Chronosphere Inc.
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

package mcpserver

import (
	"context"
	"encoding/json"
	"net/http"
	"sort"

	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"

	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/authcontext"
	"github.com/chronosphereio/chronosphere-mcp/pkg/version"
)

const (
	catalogPath          = "/mcp/catalog"
	catalogSchemaVersion = "1"
)

type Catalog struct {
	SchemaVersion     string                 `json:"schemaVersion"`
	Server            CatalogServer          `json:"server"`
	Tools             []mcp.Tool             `json:"tools"`
	Resources         []mcp.Resource         `json:"resources"`
	ResourceTemplates []mcp.ResourceTemplate `json:"resourceTemplates"`
	Prompts           []mcp.Prompt           `json:"prompts"`
}

type CatalogServer struct {
	Name         string                 `json:"name"`
	Version      string                 `json:"version"`
	Instructions string                 `json:"instructions"`
	Capabilities mcp.ServerCapabilities `json:"capabilities"`
}

func (s *Server) Catalog(ctx context.Context) Catalog {
	catalogTools := make([]mcp.Tool, len(s.tools))
	copy(catalogTools, s.tools)
	catalogTools = filterTools(ctx, catalogTools, s.writeToolNames, s.enableWrites)
	sort.Slice(catalogTools, func(i, j int) bool {
		return catalogTools[i].Name < catalogTools[j].Name
	})

	catalogResources := make([]mcp.Resource, len(s.resources))
	copy(catalogResources, s.resources)
	sort.Slice(catalogResources, func(i, j int) bool {
		return catalogResources[i].Name < catalogResources[j].Name
	})

	return Catalog{
		SchemaVersion: catalogSchemaVersion,
		Server: CatalogServer{
			Name:         "Chronosphere MCP Server",
			Version:      version.Version,
			Instructions: "",
			Capabilities: catalogCapabilities(),
		},
		Tools:             catalogTools,
		Resources:         catalogResources,
		ResourceTemplates: []mcp.ResourceTemplate{},
		Prompts:           []mcp.Prompt{},
	}
}

func catalogCapabilities() mcp.ServerCapabilities {
	return mcp.ServerCapabilities{
		Logging: &struct{}{},
		Prompts: &struct {
			ListChanged bool `json:"listChanged,omitempty"`
		}{
			ListChanged: true,
		},
		Resources: &struct {
			Subscribe   bool `json:"subscribe,omitempty"`
			ListChanged bool `json:"listChanged,omitempty"`
		}{
			Subscribe:   true,
			ListChanged: true,
		},
		Tools: &struct {
			ListChanged bool `json:"listChanged,omitempty"`
		}{
			ListChanged: true,
		},
	}
}

func (s *Server) CatalogHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Allow", "GET, HEAD")
		w.Header().Set("Cache-Control", "private, no-store")
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		ctx := authcontext.HTTPInboundContextFunc(r.Context(), r)
		body, err := json.Marshal(s.Catalog(ctx))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			if _, writeErr := w.Write([]byte(`{"error":"failed to encode catalog"}`)); writeErr != nil {
				s.logger.Error("failed to write catalog error", zap.Error(writeErr))
			}
			return
		}
		if r.Method == http.MethodHead {
			return
		}
		if _, err := w.Write(body); err != nil {
			s.logger.Error("failed to write catalog", zap.Error(err))
		}
	})
}

func (s *Server) HTTPHandler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/mcp", s.StreamableHTTPServer())
	mux.Handle(catalogPath, s.CatalogHandler())
	return mux
}
