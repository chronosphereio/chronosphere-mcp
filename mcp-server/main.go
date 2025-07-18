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

// Package main is the entry point to starting the mcp-server.
package main

import (
	"fmt"
	"os"

	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/cmd"
	"github.com/chronosphereio/chronosphere-mcp/pkg/version"
)

func main() {
	cmd := cmd.New()
	cmd.Version = fmt.Sprintf("%s (%s) %s", version.Version, version.GitCommit, version.BuildDate)
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
