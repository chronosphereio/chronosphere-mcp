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
