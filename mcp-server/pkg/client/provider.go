// Package client provides clients for the Chronosphere API.
package client

import (
	"fmt"
	"net/http"
	"net/url"

	openapiclient "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/prometheus/client_golang/api"

	"github.com/chronosphereio/mcp-server/generated/configv1/configv1"
	"github.com/chronosphereio/mcp-server/generated/dataunstable/dataunstable"
	"github.com/chronosphereio/mcp-server/generated/stateunstable/stateunstable"
)

var (
	_component = Component("chrono-mcp")
)

type Provider struct {
	apiURL   string
	apiToken string
}

func NewProvider(apiURL, apiToken string) (*Provider, error) {
	return &Provider{
		apiURL:   apiURL,
		apiToken: apiToken,
	}, nil
}

// ConfigV1Client creates a new client to hit configv1 APIs.
func (c *Provider) ConfigV1Client() (*configv1.ConfigV1API, error) {
	t, err := c.transportForSession("")
	if err != nil {
		return nil, fmt.Errorf("could not construct Chronosphere config v1 API client: %v", err)
	}
	return configv1.New(t, strfmt.Default), nil
}

// PrometheusDataClient creates a new client to hit Prometheus APIs.
func (c *Provider) PrometheusDataClient() (api.Client, error) {
	return c.prometheusClientForBasePath("/data/metrics")
}

// PrometheusPromClient creates a new client to hit Prometheus APIs.
func (c *Provider) PrometheusPromClient() (api.Client, error) {
	return c.prometheusClientForBasePath("/app/prom")
}

// DataUnstableClient creates a new client to hit data unstable APIs.
func (c *Provider) DataUnstableClient() (*dataunstable.DataUnstableAPI, error) {
	t, err := c.transportForSession("/")
	if err != nil {
		return nil, fmt.Errorf("could not construct Chronosphere data unstable API client: %v", err)
	}
	return dataunstable.New(t, strfmt.Default), nil
}

// StateUnstableClient creates a new client to hit state unstable APIs.
func (c *Provider) StateUnstableClient() (*stateunstable.StateUnstableAPI, error) {
	t, err := c.transportForSession("/")
	if err != nil {
		return nil, fmt.Errorf("could not construct Chronosphere state unstable API client: %v", err)
	}
	return stateunstable.New(t, strfmt.Default), nil
}

// TransportForSession creates a transport for the given session.
// This is exposed for use by tools that need to create clients that aren't provided directly.
func (c *Provider) TransportForSession(basePath string) (*openapiclient.Runtime, error) {
	return c.transportForSession(basePath)
}

func (c *Provider) prometheusClientForBasePath(basePath string) (api.Client, error) {
	t, err := c.transportForSession(basePath)
	if err != nil {
		return nil, fmt.Errorf("could not construct Chronosphere Prometheus API transport: %v", err)
	}

	scheme := "https"
	if c.apiURL != "" {
		parsed, err := url.Parse(c.apiURL)
		if err != nil {
			return nil, fmt.Errorf("could not parse API url given: %v", err)
		}
		scheme = parsed.Scheme
	}

	rt := newRoundTripper(http.DefaultTransport, _component, c.apiToken)

	cl, err := api.NewClient(api.Config{
		Address:      fmt.Sprintf("%s://%s/%s", scheme, t.Host, t.BasePath),
		RoundTripper: rt,
	})
	if err != nil {
		return nil, fmt.Errorf("could not construct Chronosphere Prometheus API client: %v", err)
	}
	return cl, nil
}

func (c *Provider) transportForSession(basePath string) (*openapiclient.Runtime, error) {
	return newSwaggerRuntime(swaggerRuntimeConfig{
		component: _component,
		apiURL:    fmt.Sprintf("%s%s", c.apiURL, basePath),
		allowHTTP: false,
		authToken: c.apiToken,
	})
}
