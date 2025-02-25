package main

import (
	"os"

	"github.com/chronosphereio/mcp-server/mcp-server/pkg/cmd"
)

func main() {
	cmd := cmd.New()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
