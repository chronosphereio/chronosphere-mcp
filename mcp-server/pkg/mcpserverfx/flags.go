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
