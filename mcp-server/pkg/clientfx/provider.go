// Package clientfx provides clients for the Chronosphere API.
package clientfx

import (
	"fmt"
	"net/http"
	"net/url"

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
	promClient, err := PrometheusDataClient(apiConfig)
	if err != nil {
		return Provider{}, fmt.Errorf("could not create Prometheus data client: %v", err)
	}
	configV1, err := ConfigV1Client(apiConfig)
	if err != nil {
		return Provider{}, fmt.Errorf("could not create config v1 client: %v", err)
	}
	dataUnstableClient, err := DataUnstableClient(apiConfig)
	if err != nil {
		return Provider{}, fmt.Errorf("could not create data unstable client: %v", err)
	}
	stateUnstableClient, err := StateUnstableClient(apiConfig)
	if err != nil {
		return Provider{}, fmt.Errorf("could not create state unstable client: %v", err)
	}
	return Provider{
		PrometheusData: promClient,
		ConfigV1:       configV1,
		DataUnstable:   dataUnstableClient,
		StateUnstable:  stateUnstableClient,
	}, nil
}

// ConfigV1Client creates a new client to hit configv1 APIs.
func ConfigV1Client(config *APIConfig) (*configv1.ConfigV1API, error) {
	t, err := transportForSession(config, "")
	if err != nil {
		return nil, fmt.Errorf("could not construct Chronosphere config v1 API client: %v", err)
	}
	return configv1.New(t, strfmt.Default), nil
}

// PrometheusDataClient creates a new client to hit Prometheus APIs.
func PrometheusDataClient(config *APIConfig) (api.Client, error) {
	return prometheusClientForBasePath(config, "/data/metrics")
}

// DataUnstableClient creates a new client to hit data unstable APIs.
func DataUnstableClient(config *APIConfig) (*dataunstable.DataUnstableAPI, error) {
	t, err := transportForSession(config, "/")
	if err != nil {
		return nil, fmt.Errorf("could not construct Chronosphere data unstable API client: %v", err)
	}
	return dataunstable.New(t, strfmt.Default), nil
}

// StateUnstableClient creates a new client to hit state unstable APIs.
func StateUnstableClient(config *APIConfig) (*stateunstable.StateUnstableAPI, error) {
	t, err := transportForSession(config, "/")
	if err != nil {
		return nil, fmt.Errorf("could not construct Chronosphere state unstable API client: %v", err)
	}
	return stateunstable.New(t, strfmt.Default), nil
}

func prometheusClientForBasePath(config *APIConfig, basePath string) (api.Client, error) {
	t, err := transportForSession(config, basePath)
	if err != nil {
		return nil, fmt.Errorf("could not construct Chronosphere Prometheus API transport: %v", err)
	}

	scheme := "https"
	if config.APIURL != "" {
		parsed, err := url.Parse(config.APIURL)
		if err != nil {
			return nil, fmt.Errorf("could not parse API url given: %v", err)
		}
		scheme = parsed.Scheme
	}

	rt := newRoundTripper(http.DefaultTransport, _component, config.APIToken)

	cl, err := api.NewClient(api.Config{
		Address:      fmt.Sprintf("%s://%s/%s", scheme, t.Host, t.BasePath),
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
