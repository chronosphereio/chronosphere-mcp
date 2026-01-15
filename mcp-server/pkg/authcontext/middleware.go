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

package authcontext

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

const (
	_chronoAccessTokenHeaderName = "chrono-accesstoken"
	_disableToolsHeaderName      = "X-Chrono-MCP-Disable-Tools"
)

// HTTPInboundContextFunc extracts the Authorization header from the HTTP request and sets it in the context.
func HTTPInboundContextFunc(ctx context.Context, r *http.Request) context.Context {
	ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(r.Header))
	authValue := strings.ReplaceAll(r.Header.Get("Authorization"), "Bearer ", "")
	cookie, err := r.Cookie(_chronoAccessTokenHeaderName)
	var cookieValue string
	if err == nil {
		cookieValue = cookie.Value
	}
	ctx = SetSessionCredentials(ctx, SessionCredentials{
		APIToken:          authValue,
		AccessTokenCookie: cookieValue,
	})

	// Extract disabled tools from header
	disabledToolsHeader := r.Header.Get(_disableToolsHeaderName)
	if disabledToolsHeader != "" {
		disabledTools := make(map[string]struct{})
		for _, tool := range strings.Split(disabledToolsHeader, ",") {
			tool = strings.TrimSpace(tool)
			if tool != "" {
				disabledTools[tool] = struct{}{}
			}
		}
		ctx = SetDisabledTools(ctx, disabledTools)
	}

	return ctx
}

// RoundTripper wraps an http.RoundTripper and adds an Authorization header.
type RoundTripper struct {
	token     string
	transport http.RoundTripper
}

func NewRoundTripper(base http.RoundTripper, token string) *RoundTripper {
	return &RoundTripper{
		token:     token,
		transport: base,
	}
}

func (r *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// Clone the request to avoid modifying the original.
	req2 := req.Clone(req.Context())

	if credentials := FetchSessionAPIToken(req.Context()); !credentials.IsEmpty() {
		// forward the api token & access token cookie from context if available
		req2.Header.Set("Authorization", "Bearer "+credentials.APIToken)
		req2.Header.Set("Cookie", fmt.Sprintf("%s=%s", _chronoAccessTokenHeaderName, credentials.AccessTokenCookie))
	} else if r.token != "" {
		req2.Header.Set("Authorization", "Bearer "+r.token)
	}

	return r.transport.RoundTrip(req2)
}
