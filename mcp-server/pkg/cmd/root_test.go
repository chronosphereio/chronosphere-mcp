package cmd

import (
	"testing"
	"time"

	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/mcpserverfx"
	"go.uber.org/fx/fxtest"
)

func TestFxStartup(t *testing.T) {
	tests := []struct {
		configFile string
	}{
		{
			configFile: "../../../config.yaml",
		},
		{
			configFile: "../../../config.http.yaml",
		},
		{
			configFile: "../../../config.sse.yaml",
		},
		{
			configFile: "testdata/config.nil-tools.yaml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.configFile, func(t *testing.T) {
			flags := &mcpserverfx.Flags{
				ConfigFilePath: tt.configFile,
			}
			app := fxtest.New(t, allModules(flags)...).RequireStart()
			// It can take a small amount of time for the http server to start up, so wait a second
			// before calling shutdown otherwise we could end up shutting down before the server is initialized.
			time.Sleep(1 * time.Second)
			app.RequireStop()
		})
	}
}
