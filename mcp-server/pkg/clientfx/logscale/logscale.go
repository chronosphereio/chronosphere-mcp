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

// Package logscale contains LogScale query client
package logscale

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	contentType      = "application/json"
	repositoriesPath = "/api/v1/repositories"
	queryPathSuffix  = "query"
)

// Client is a logscale query API client
type Client interface {
	// Query runs a query against a given logscale repository
	Query(ctx context.Context, query, repository string, start, end time.Time) ([]map[string]any, error)
}

// Client communicates with LogScale query API over HTTP.
type client struct {
	tenantURL *url.URL
	client    *http.Client
	apiToken  string
}

// Options are options for the logscale query client.
type Options struct {
	URL string // required
	// APIToken for logscale (this is different from the Chronosphere API token).
	APIToken string
	// Transport is an optional RoundTripper that is used as the transport for the http request.
	// If not set, the http.DefaultTransport is used.
	Transport http.RoundTripper
}

func (o *Options) validate() error {
	if o.URL == "" {
		return errors.New("tenant URL is required")
	}
	return nil
}

// New makes a new HTTP logscale query client.
func New(opts *Options) (Client, error) {
	if err := opts.validate(); err != nil {
		return nil, err
	}

	// Configure transport
	rt := opts.Transport
	if rt == nil {
		rt = http.DefaultTransport
	}

	tenantURL, err := url.Parse(opts.URL)
	if err != nil {
		return nil, err
	}
	return &client{
		tenantURL: tenantURL,
		client:    &http.Client{Transport: rt},
		apiToken:  opts.APIToken,
	}, nil
}

func (c *client) Query(ctx context.Context, query, repository string, start, end time.Time) ([]map[string]any, error) {
	if repository == "" {
		return nil, errors.New("repository must be set")
	}

	queryURL := c.tenantURL.JoinPath(queryPath(repository)...).String()
	b, err := json.Marshal(queryRequest{
		Query: query,
		Start: start.UnixMilli(),
		End:   end.UnixMilli(),
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", queryURL, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Accept", contentType)
	req.Header.Set("Authorization", "Bearer "+c.apiToken)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() //nolint:errcheck

	err = checkResponseError(resp)
	if err != nil {
		return nil, err
	}

	var result []map[string]any
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response body: %w", err)
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("unable to unmarshal response body: %w", err)
	}
	return result, nil
}

func queryPath(repository string) []string {
	return []string{
		repositoriesPath,
		repository,
		queryPathSuffix,
	}
}

type queryRequest struct {
	Query string `json:"queryString"`
	Start int64  `json:"start"`
	End   int64  `json:"end"`
}

// checkResponseError checks the http response's status code and returns an error if the code is not 200
func checkResponseError(resp *http.Response) error {
	if resp.StatusCode != http.StatusOK {
		all, err := io.ReadAll(resp.Body)
		if err == nil {
			return fmt.Errorf("server returned non-200 status: %d %s %s", resp.StatusCode, resp.Status, string(all))
		}
		return fmt.Errorf("server returned non-200 status: %d %s", resp.StatusCode, resp.Status)
	}
	return nil
}
