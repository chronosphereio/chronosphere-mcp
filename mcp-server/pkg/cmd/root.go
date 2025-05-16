package cmd

import (
	"os"
	"path"

	chronoctlclient "github.com/chronosphereio/chronoctl-core/src/cmd/pkg/client"
	"github.com/spf13/cobra"
	"go.uber.org/config"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/chronosphereio/mcp-server/mcp-server/pkg/client"
	pkgconfig "github.com/chronosphereio/mcp-server/mcp-server/pkg/config"
	"github.com/chronosphereio/mcp-server/mcp-server/pkg/mcpserverfx"
	"github.com/chronosphereio/mcp-server/mcp-server/pkg/toolsfx"
)

// New returns the root command.
func New() *cobra.Command {
	flags := chronoctlclient.NewClientFlags()
	configFilePath := ""
	verboseLogging := false

	cmd := &cobra.Command{
		Use:   "mcp-server",
		Short: "mcp-server provides an MCP server to AI applications",
		Long:  "mcp-server provides an MCP server to AI applications",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// If command parsing works, let's silence usage so errors RunE errors
			// don't display usage (adding unnecessary noise to the output).
			cmd.SilenceUsage = true
		},
		Run: func(*cobra.Command, []string) {
			app := fx.New(
				fx.Provide(func() (config.Provider, error) {
					return pkgconfig.ParseFile(configFilePath)
				}),
				fx.Provide(func() (*client.Provider, error) {
					return client.NewProvider(flags)
				}),
				fx.Provide(func() (*zap.Logger, error) {
					return provideLogger(verboseLogging)
				}),
				toolsfx.Module,
				mcpserverfx.Module,
			)
			app.Run()
		},
	}

	cmd.Flags().StringVarP(&configFilePath, "config-file", "c", "", "The YAML file containing configuration parameters")
	cmd.Flags().BoolVarP(&verboseLogging, "verbose", "v", false, "Whether verbose logging, including logging requests and responses, should be enabled")
	flags.AddFlags(cmd)
	return cmd
}

func provideLogger(verboseLogging bool) (*zap.Logger, error) {
	tmpDir := os.TempDir()
	loggerCfg := zap.NewProductionConfig()
	loggerCfg.OutputPaths = []string{"stdout", path.Join(tmpDir, "mcp_server.log")}
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
