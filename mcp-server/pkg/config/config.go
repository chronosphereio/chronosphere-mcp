// Package config parses a config file
package config

import (
	"bytes"
	"log"
	"os"

	"go.uber.org/config"
)

func ParseFile(path string) (config.Provider, error) {
	log.Printf("Loading config file from %q\n", path)

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
