// Package clientfx provides clients for the Chronosphere API.
package clientfx

import (
	"fmt"
	"net/http"

	openapiclient "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/prometheus/client_golang/api"
	"go.uber.org/fx"

	"github.com/chronosphereio/mcp-server/generated/configv1/configv1"
	"github.com/chronosphereio/mcp-server/generated/dataunstable/dataunstable"
	"github.com/chronosphereio/mcp-server/generated/stateunstable/stateunstable"
)

var (
	_component = Component("chrono-mcp")
	Module     = fx.Provide(
		NewProvider,
	)
)

type APIConfig struct {
	APIURL   string
	APIToken string
}

type Provider struct {
	fx.Out

	PrometheusData api.Client `fx:"PrometheusData"`
	ConfigV1       *configv1.ConfigV1API
	DataUnstable   *dataunstable.DataUnstableAPI
	StateUnstable  *stateunstable.StateUnstableAPI
}

func NewProvider(apiConfig *APIConfig) (Provider, error) {
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
		StateUnstable:  stateunstable.New(t, strfmt.Default),
	}, nil
}

func prometheusClientForBasePath(config *APIConfig, basePath string) (api.Client, error) {
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

func transportForSession(config *APIConfig, basePath string) (*openapiclient.Runtime, error) {
	return newSwaggerRuntime(swaggerRuntimeConfig{
		component: _component,
		apiURL:    fmt.Sprintf("%s%s", config.APIURL, basePath),
		allowHTTP: false,
		authToken: config.APIToken,
	})
}
