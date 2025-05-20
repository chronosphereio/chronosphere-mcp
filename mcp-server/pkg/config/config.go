// Package config parses a config file
package config

import (
	"bytes"
	"os"

	"go.uber.org/config"
)

func ParseFile(path string) (config.Provider, error) {
	data, err := os.ReadFile(path) //nolint:gosec
	if err != nil {
		return nil, err
	}

	opts := []config.YAMLOption{
		config.Expand(os.LookupEnv),
		config.Source(bytes.NewReader(data)),
	}
	provider, err := config.NewYAML(opts...)
	if err != nil {
		return nil, err
	}
	return provider, nil
}
