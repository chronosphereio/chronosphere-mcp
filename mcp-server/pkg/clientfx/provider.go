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

// Package clientfx provides clients for the Chronosphere API.
package clientfx

import (
	"fmt"
	"net/http"

	openapiclient "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/prometheus/client_golang/api"
	"go.uber.org/fx"

	"github.com/chronosphereio/chronosphere-mcp/generated/configv1/configv1"
	"github.com/chronosphereio/chronosphere-mcp/generated/dataunstable/dataunstable"
	"github.com/chronosphereio/chronosphere-mcp/generated/datav1/datav1"
	"github.com/chronosphereio/chronosphere-mcp/generated/stateunstable/stateunstable"
	"github.com/chronosphereio/chronosphere-mcp/generated/statev1/statev1"
)

var (
	_component = Component("chrono-mcp")
	Module     = fx.Provide(
		NewProvider,
	)
)

type ChronosphereConfig struct {
	APIURL   string `yaml:"apiURL" validate:"nonzero"`
	APIToken string `yaml:"apiToken"`
}

func (c *ChronosphereConfig) Validate() error {
	if c.APIURL == "" {
		return fmt.Errorf("apiURL must be set")
	}
	return nil
}

type Provider struct {
	fx.Out

	PrometheusData api.Client `fx:"PrometheusData"`
	ConfigV1       *configv1.ConfigV1API
	DataUnstable   *dataunstable.DataUnstableAPI
	DataV1         *datav1.DataV1API
	StateV1        *statev1.StateV1API
	StateUnstable  *stateunstable.StateUnstableAPI
}

func NewProvider(apiConfig *ChronosphereConfig) (Provider, error) {
	t, err := transportForSession(apiConfig, "")
	if err != nil {
		return Provider{}, fmt.Errorf("could not construct Chronosphere config v1 API client: %v", err)
	}

	promClient, err := prometheusClientForBasePath(apiConfig, "/data/metrics")
	if err != nil {
		return Provider{}, fmt.Errorf("could not create Prometheus data client: %v", err)
	}

	return Provider{
		PrometheusData: promClient,
		ConfigV1:       configv1.New(t, strfmt.Default),
		DataUnstable:   dataunstable.New(t, strfmt.Default),
		DataV1:         datav1.New(t, strfmt.Default),
		StateV1:        statev1.New(t, strfmt.Default),
		StateUnstable:  stateunstable.New(t, strfmt.Default),
	}, nil
}

func prometheusClientForBasePath(config *ChronosphereConfig, basePath string) (api.Client, error) {
	rt := newRoundTripper(http.DefaultTransport, _component, config.APIToken)

	cl, err := api.NewClient(api.Config{
		Address:      fmt.Sprintf("%s%s", config.APIURL, basePath),
		RoundTripper: rt,
	})
	if err != nil {
		return nil, fmt.Errorf("could not construct Chronosphere Prometheus API client: %v", err)
	}
	return cl, nil
}

func transportForSession(config *ChronosphereConfig, basePath string) (*openapiclient.Runtime, error) {
	return newSwaggerRuntime(swaggerRuntimeConfig{
		component: _component,
		apiURL:    fmt.Sprintf("%s%s", config.APIURL, basePath),
		allowHTTP: false,
		authToken: config.APIToken,
	})
}
