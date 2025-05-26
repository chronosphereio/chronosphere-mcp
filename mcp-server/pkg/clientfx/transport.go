package clientfx

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	xswagger "github.com/chronosphereio/chronoctl-core/src/x/swagger"
	httpruntime "github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"

	"github.com/chronosphereio/mcp-server/mcp-server/pkg/authcontext"
	"github.com/chronosphereio/mcp-server/pkg/version"
)

// Component is a value that indicates the part of the CLI that is invoking an
// API. This is used to set the User-Agent header when making requests to the Chronosphere API.
type Component string

// swaggerRuntimeConfig is a struct that contains the configuration for creating a new HTTP transport
type swaggerRuntimeConfig struct {
	component Component
	apiURL    string
	allowHTTP bool
	authToken string
}

func newRoundTripper(base http.RoundTripper, component Component, authToken string) http.RoundTripper {
	return roundTripperFn(func(req *http.Request) (*http.Response, error) {
		req = req.Clone(req.Context())
		req.Header.Set("User-Agent", fmt.Sprintf("%s/%v-%v", component, version.Version, version.GitCommit))
		return authcontext.NewRoundTripper(base, authToken).RoundTrip(req)
	})
}

// New creates a new HTTP transport that can communicate with the Chronosphere API.
func newSwaggerRuntime(config swaggerRuntimeConfig) (*httptransport.Runtime, error) {
	apiURL, err := url.Parse(config.apiURL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse URL of Chronosphere URL: %v", err)
	}

	var schemes []string
	switch apiURL.Scheme {
	case "https":
		schemes = []string{"https"}
	case "http":
		if !config.allowHTTP {
			return nil, errors.New("the scheme of the API URL is HTTP but the --allow-http flag was not set")
		}
		schemes = []string{"http"}
	default:
		// If the client didn't specify a scheme in the URL just use the API client's default
		// by passing an empty slice to the transport constructor.
	}

	transport := httptransport.New(apiURL.Host, apiURL.Path, schemes)

	transport.Transport = xswagger.WithRequestIDTrailerTransport(
		newRoundTripper(transport.Transport, config.component, config.authToken),
	)
	transport.Consumers[httpruntime.JSONMime] = xswagger.JSONConsumer()
	transport.Consumers[httpruntime.HTMLMime] = xswagger.TextConsumer()
	transport.Consumers[httpruntime.TextMime] = xswagger.TextConsumer()
	transport.Consumers["*/*"] = xswagger.TextConsumer() // backup, default consumer.

	return transport, nil
}

type roundTripperFn func(req *http.Request) (*http.Response, error)

var _ http.RoundTripper = roundTripperFn(nil)

func (fn roundTripperFn) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}
