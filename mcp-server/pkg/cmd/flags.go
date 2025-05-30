package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

const (
	// ChronosphereOrgNameKey is the environment variable that specifies the Chronosphere customer organization
	ChronosphereOrgNameKey = "CHRONOSPHERE_ORG_NAME"
	// ChronosphereAPITokenKey is the environment variable that specifies the Chronosphere API token
	ChronosphereAPITokenKey = "CHRONOSPHERE_API_TOKEN"
	// LogscaleAPITokenKey is the environment variable that specifies the LogScale API token
	LogscaleAPITokenKey = "LOGSCALE_API_TOKEN"
	apiURLFormat        = "https://%s.chronosphere.io"
	logscaleURLFormat   = "https://%s.logs.chronosphere.io"
)

type Flags struct {
	apiToken         string
	apiTokenFileName string
	orgName          string
	apiURL           string
	ConfigFilePath   string
	VerboseLogging   bool
	UseLogScale      bool
}

// AddFlags adds client flags to a Cobra command.
func (f *Flags) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&f.ConfigFilePath, "config-file", "c", "", "The YAML file containing configuration parameters")
	cmd.Flags().BoolVarP(&f.VerboseLogging, "verbose", "v", false, "Whether verbose logging, including logging requests and responses, should be enabled")
	cmd.Flags().StringVar(&f.apiToken, "api-token", "", "The client API token used to authenticate to user. Mutally exclusive with --api-token-filename. If both --api-token and --api-token-filename are unset, the "+ChronosphereAPITokenKey+" environment variable is used.")
	cmd.Flags().StringVar(&f.apiTokenFileName, "api-token-filename", "", "A file containing the API token used for authentication. Mutally exclusive with --api-token. If both --api-token and --api-token-filename are unset, the "+ChronosphereAPITokenKey+" environment variable is used.")
	cmd.Flags().StringVar(&f.orgName, "org-name", "", "The name of your team's Chronosphere organization. Defaults to "+ChronosphereOrgNameKey+" environment variable.")
	cmd.Flags().BoolVar(&f.UseLogScale, "use-logscale", false, "Whether to use LogScale instead of chronosphere logs for log queries. If set, the LogScale API token must be provided via the LOGSCALE_API_TOKEN environment variable.")

	cmd.Flags().StringVar(&f.apiURL, "api-url", f.apiURL, "The URL of the Chronosphere API. Defaults to https://<organization>.chronosphere.io/api.")
	cmd.Flags().MarkHidden("api-url") //nolint:errcheck
}

func (f *Flags) GetLogscaleAPIToken() (string, error) {
	if f.UseLogScale {
		if token := os.Getenv(LogscaleAPITokenKey); token != "" {
			return token, nil
		}
		return "", errors.New("LogScale API token must be provided via the " + LogscaleAPITokenKey + " environment variable")
	}
	return "", nil
}

func (f *Flags) GetAPIToken() (string, error) {
	if f.apiToken != "" && f.apiTokenFileName != "" {
		return "", errors.New("only one of --api-token and --api-token-filename can be set")
	}

	if f.apiToken != "" {
		return f.apiToken, nil
	}

	if f.apiTokenFileName != "" {
		b, err := os.ReadFile(f.apiTokenFileName)
		if err != nil {
			return "", fmt.Errorf("reading api token from file %s: %w", f.apiTokenFileName, err)
		}
		return strings.TrimSpace(string(b)), nil
	}

	if key := os.Getenv(ChronosphereAPITokenKey); key != "" {
		return key, nil
	}

	return "", errors.New("api token must be provided as a flag, via the " + ChronosphereAPITokenKey + " environment variable, or by setting a file with the --api-token-filename flag")
}

func (f *Flags) GetAPIURL() (string, error) {
	if f.apiURL != "" {
		return f.apiURL, nil
	}
	if f.orgName == "" {
		if f.orgName = os.Getenv(ChronosphereOrgNameKey); f.orgName == "" {
			return "", errors.New("org name must be provided as a flag or via the " + ChronosphereOrgNameKey + " environment variable")
		}
	}
	return fmt.Sprintf(apiURLFormat, f.orgName), nil
}

func (f *Flags) GetLogscaleURL() (string, error) {
	if f.UseLogScale {
		if f.orgName == "" {
			if f.orgName = os.Getenv(ChronosphereOrgNameKey); f.orgName == "" {
				return "", errors.New("org name must be provided as a flag or via the " + ChronosphereOrgNameKey + " environment variable")
			}
		}
		return fmt.Sprintf(logscaleURLFormat, f.orgName), nil
	}
	return "", nil
}
