// Package main is the entry point to starting the mcp-server.
package main

import (
	"fmt"
	"os"

	"github.com/chronosphereio/mcp-server/mcp-server/pkg/cmd"
	"github.com/chronosphereio/mcp-server/pkg/version"
)

func main() {
	cmd := cmd.New()
	cmd.Version = fmt.Sprintf("%s (%s) %s", version.Version, version.GitCommit, version.BuildDate)
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
