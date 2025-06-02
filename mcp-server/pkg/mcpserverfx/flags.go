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

package mcpserverfx

import (
	"github.com/spf13/cobra"
)

type Flags struct {
	apiToken       string
	orgName        string
	ConfigFilePath string
	VerboseLogging bool
	UseLogScale    bool
}

// AddFlags adds client flags to a Cobra command.
func (f *Flags) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&f.ConfigFilePath, "config-file", "c", "", "The YAML file containing configuration parameters")
	cmd.Flags().BoolVarP(&f.VerboseLogging, "verbose", "v", false, "Whether verbose logging, including logging requests and responses, should be enabled")
	cmd.Flags().StringVar(&f.apiToken, "api-token", "", "The client API token used to authenticate to user.")
	cmd.Flags().StringVar(&f.orgName, "org-name", "", "The name of your team's Chronosphere organization.")
	cmd.Flags().BoolVar(&f.UseLogScale, "use-logscale", false, "Whether to use LogScale instead of chronosphere logs for log queries. If set, the LogScale API token must be provided via the LOGSCALE_API_TOKEN environment variable.")
}
