// Package cmd contains the main mcp server command
package cmd

import (
	"os"
	"path"

	"github.com/spf13/cobra"
	"go.uber.org/config"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/chronosphereio/mcp-server/mcp-server/pkg/clientfx"
	pkgconfig "github.com/chronosphereio/mcp-server/mcp-server/pkg/config"
	"github.com/chronosphereio/mcp-server/mcp-server/pkg/mcpserverfx"
	"github.com/chronosphereio/mcp-server/mcp-server/pkg/toolsfx"
	"github.com/chronosphereio/mcp-server/pkg/links"
)

// New returns the root command.
func New() *cobra.Command {
	flags := &Flags{}
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
				fx.Provide(func() (*clientfx.APIConfig, error) {
					apiURL, err := flags.GetAPIURL()
					if err != nil {
						return nil, err
					}
					apiToken, err := flags.GetAPIToken()
					if err != nil {
						return nil, err
					}
					return &clientfx.APIConfig{
						APIURL:   apiURL,
						APIToken: apiToken,
					}, nil
				}),
				fx.Provide(func(apiConfig *clientfx.APIConfig) *links.Builder {
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
