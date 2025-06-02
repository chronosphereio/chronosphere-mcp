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

// Package configfx contains a fx module for config provider.
package configfx

import (
	"bytes"
	"errors"
	"flag"
	"log"
	"os"
	"strings"

	"go.uber.org/config"
	"go.uber.org/fx"
)

// ConfigData is the contents of a config file.
type ConfigData string
type Inputs []string

// Module provides a config provider.
// By default this takes in a config file path from "-f" flag through command
// line, which is the convention used by most of the services and tools we have.
// It also accepts input from environment variable "CONFIG_FILE",
// which will override the value taken from "-f".
var Module = fx.Provide(provide)

type result struct {
	fx.Out

	ConfigProvider config.Provider
	Inputs         Inputs
}

func provide() (result, error) {
	cfgPath := flag.String("f", "", "configuration file")
	inputPath := flag.String("i", "", "input file")

	flag.Parse()
	if cfgPath == nil || len(*cfgPath) == 0 {
		flag.PrintDefaults()
		return result{}, errors.New("config file path is required")
	}

	if inputPath == nil || len(*inputPath) == 0 {
		flag.PrintDefaults()
		return result{}, errors.New("input file path is required")
	}

	opts := []config.YAMLOption{config.Expand(os.LookupEnv)}

	log.Printf("Loading config file from %q\n", *cfgPath)

	provider, err := configProvider(cfgPath, opts)
	if err != nil {
		return result{}, err
	}

	inputs, err := readInputs(inputPath)
	if err != nil {
		return result{}, err
	}

	return result{
		ConfigProvider: provider,
		Inputs:         inputs,
	}, nil
}

func readInputs(inputsPath *string) ([]string, error) {
	data, err := os.ReadFile(*inputsPath) //nolint:gosec
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	linesFiltered := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		linesFiltered = append(linesFiltered, line)
	}
	return linesFiltered, nil
}

func configProvider(cfgPath *string, opts []config.YAMLOption) (*config.YAML, error) {
	data, err := os.ReadFile(*cfgPath) //nolint:gosec
	if err != nil {
		return nil, err
	}

	opts = append(opts, config.Source(bytes.NewReader(data)))

	provider, err := config.NewYAML(opts...)
	if err != nil {
		return nil, err
	}
	return provider, nil
}
