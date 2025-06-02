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

// Package cmd contains the main mcp server command
package cmd

import (
	"os"
	"path"

	"github.com/spf13/cobra"
	"go.uber.org/config"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/clientfx"
	pkgconfig "github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/config"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/mcpserverfx"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/toolsfx"
	"github.com/chronosphereio/chronosphere-mcp/pkg/links"
)

// New returns the root command.
func New() *cobra.Command {
	flags := &mcpserverfx.Flags{}
	cmd := &cobra.Command{
		Use:   "chronomcp",
		Short: "chronomcp provides an MCP server to AI applications",
		Long:  "chronomcp provides an MCP server to AI applications",
		PersistentPreRun: func(cmd *cobra.Command, _ []string) {
			// If command parsing works, let's silence usage so errors RunE errors
			// don't display usage (adding unnecessary noise to the output).
			cmd.SilenceUsage = true
		},
		Run: func(*cobra.Command, []string) {
			app := fx.New(
				fx.Provide(func() (config.Provider, error) {
					return pkgconfig.ParseFile(flags.ConfigFilePath)
				}),
				fx.Provide(func() *mcpserverfx.Flags {
					return flags
				}),
				fx.Provide(func(apiConfig *clientfx.ChronosphereConfig) *links.Builder {
					return links.NewBuilder(apiConfig.APIURL)
				}),
				clientfx.Module,
				fx.Provide(func() (*zap.Logger, error) {
					return provideLogger(flags.VerboseLogging)
				}),
				toolsfx.Module,
				mcpserverfx.Module,
			)
			app.Run()
		},
	}

	flags.AddFlags(cmd)
	return cmd
}

func provideLogger(verboseLogging bool) (*zap.Logger, error) {
	tmpDir := os.TempDir()
	loggerCfg := zap.NewProductionConfig()
	loggerCfg.OutputPaths = []string{"stderr", path.Join(tmpDir, "mcp_server.log")}
	loggerCfg.ErrorOutputPaths = []string{"stderr", path.Join(tmpDir, "mcp_server_error.log")}
	if verboseLogging {
		loggerCfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}
	logger, err := loggerCfg.Build()
	if err != nil {
		return nil, err
	}
	logger.Info("starting MCP server with logger",
		zap.Strings("log", loggerCfg.OutputPaths),
		zap.Strings("errorLog", loggerCfg.ErrorOutputPaths))
	return logger, nil
}
