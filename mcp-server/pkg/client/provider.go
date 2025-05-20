// Package client provides clients for the Chronosphere API.
package client

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/chronosphereio/chronoctl-core/src/cmd/pkg/client"
	"github.com/chronosphereio/chronoctl-core/src/cmd/pkg/transport"
	"github.com/go-openapi/runtime"
	openapiclient "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/prometheus/client_golang/api"

	"github.com/chronosphereio/mcp-server/generated/configv1/configv1"
	"github.com/chronosphereio/mcp-server/generated/dataunstable/dataunstable"
	"github.com/chronosphereio/mcp-server/generated/stateunstable/stateunstable"
	"github.com/chronosphereio/mcp-server/mcp-server/pkg/tools"
)

var (
	component = transport.Component("mcp-server")
)

type Provider struct {
	flags *client.Flags
}

func NewProvider(flags *client.Flags) (*Provider, error) {
	return &Provider{flags: flags}, nil
}

// ConfigV1Client creates a new client to hit configv1 APIs.
func (c *Provider) ConfigV1Client(session tools.Session) (*configv1.ConfigV1API, error) {
	t, err := c.transportForSession(session, "")
	if err != nil {
		return nil, fmt.Errorf("could not construct Chronosphere config v1 API client: %v", err)
	}
	return configv1.New(t, strfmt.Default), nil
}

// PrometheusDataClient creates a new client to hit Prometheus APIs.
func (c *Provider) PrometheusDataClient(session tools.Session) (api.Client, error) {
	return c.prometheusClientForBasePath(session, "/data/metrics")
}

// PrometheusPromClient creates a new client to hit Prometheus APIs.
func (c *Provider) PrometheusPromClient(session tools.Session) (api.Client, error) {
	return c.prometheusClientForBasePath(session, "/app/prom")
}

// DataUnstableClient creates a new client to hit data unstable APIs.
func (c *Provider) DataUnstableClient(session tools.Session) (*dataunstable.DataUnstableAPI, error) {
	t, err := c.transportForSession(session, "/")
	if err != nil {
		return nil, fmt.Errorf("could not construct Chronosphere data unstable API client: %v", err)
	}
	return dataunstable.New(t, strfmt.Default), nil
}

// StateUnstableClient creates a new client to hit state unstable APIs.
func (c *Provider) StateUnstableClient(session tools.Session) (*stateunstable.StateUnstableAPI, error) {
	t, err := c.transportForSession(session, "/")
	if err != nil {
		return nil, fmt.Errorf("could not construct Chronosphere state unstable API client: %v", err)
	}
	return stateunstable.New(t, strfmt.Default), nil
}

// TransportForSession creates a transport for the given session.
// This is exposed for use by tools that need to create clients that aren't provided directly.
func (c *Provider) TransportForSession(session tools.Session, basePath string) (*openapiclient.Runtime, error) {
	return c.transportForSession(session, basePath)
}

func (c *Provider) prometheusClientForBasePath(session tools.Session, basePath string) (api.Client, error) {
	t, err := c.transportForSession(session, basePath)
	if err != nil {
		return nil, fmt.Errorf("could not construct Chronosphere Prometheus API transport: %v", err)
	}

	scheme := "https"
	if c.flags.APIUrl != "" {
		parsed, err := url.Parse(c.flags.APIUrl)
		if err != nil {
			return nil, fmt.Errorf("could not parse API url given: %v", err)
		}
		scheme = parsed.Scheme
	}

	rt, err := c.prometheusRoundTripper(t.DefaultAuthentication, t.Transport)
	if err != nil {
		return nil, fmt.Errorf("could not construct Chronosphere Prometheus RoundTripper: %v", err)
	}

	cl, err := api.NewClient(api.Config{
		Address:      fmt.Sprintf("%s://%s/%s", scheme, t.Host, t.BasePath),
		RoundTripper: rt,
	})
	if err != nil {
		return nil, fmt.Errorf("could not construct Chronosphere Prometheus API client: %v", err)
	}
	return cl, nil
}

// prometheusRoundTripper populates the Chronosphere API token as a Bearer token, suitable for use with Prometheus APIs.
// The OpenAPI route typically populates the API token header in "Api-Token", which is different from Prometheus.
func (c *Provider) prometheusRoundTripper(openapiAuth runtime.ClientAuthInfoWriter, base http.RoundTripper) (http.RoundTripper, error) {
	testReq := &runtime.TestClientRequest{}
	if err := openapiAuth.AuthenticateRequest(testReq, strfmt.Default); err != nil {
		return nil, err
	}

	apiToken := ""
	if testReq.GetHeaderParams() != nil {
		if apiTokens := testReq.GetHeaderParams()["Api-Token"]; len(apiTokens) > 0 {
			apiToken = apiTokens[0]
		}
	}
	return roundTripperFn(func(req *http.Request) (*http.Response, error) {
		if apiToken != "" {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiToken))
		}
		return base.RoundTrip(req)
	}), nil
}

func (c *Provider) transportForSession(session tools.Session, basePath string) (*openapiclient.Runtime, error) {
	flagsCopy := *c.flags
	if session.APIToken != "" {
		flagsCopy.APIToken = session.APIToken
		flagsCopy.APITokenFilename = ""
	}
	return flagsCopy.Transport(component, basePath)
}

type roundTripperFn func(req *http.Request) (*http.Response, error)

var _ http.RoundTripper = roundTripperFn(nil)

func (fn roundTripperFn) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}
