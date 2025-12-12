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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPInboundContextFunc_DisabledTools(t *testing.T) {
	tests := []struct {
		name                  string
		headerValue           string
		expectedDisabledTools map[string]struct{}
	}{
		{
			name:        "single tool disabled",
			headerValue: "fetch_logs",
			expectedDisabledTools: map[string]struct{}{
				"fetch_logs": {},
			},
		},
		{
			name:        "multiple tools disabled",
			headerValue: "fetch_logs,fetch_metrics,fetch_traces",
			expectedDisabledTools: map[string]struct{}{
				"fetch_logs":    {},
				"fetch_metrics": {},
				"fetch_traces":  {},
			},
		},
		{
			name:        "tools with whitespace",
			headerValue: "fetch_logs, fetch_metrics , fetch_traces",
			expectedDisabledTools: map[string]struct{}{
				"fetch_logs":    {},
				"fetch_metrics": {},
				"fetch_traces":  {},
			},
		},
		{
			name:                  "empty header",
			headerValue:           "",
			expectedDisabledTools: nil,
		},
		{
			name:                  "only whitespace",
			headerValue:           "  ,  ,  ",
			expectedDisabledTools: map[string]struct{}{},
		},
		{
			name:        "mixed valid and empty entries",
			headerValue: "fetch_logs, , fetch_metrics",
			expectedDisabledTools: map[string]struct{}{
				"fetch_logs":    {},
				"fetch_metrics": {},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.headerValue != "" {
				req.Header.Set("X-Chrono-MCP-Disable-Tools", tt.headerValue)
			}

			ctx := HTTPInboundContextFunc(t.Context(), req)

			disabledTools := FetchDisabledTools(ctx)
			assert.Equal(t, tt.expectedDisabledTools, disabledTools)
		})
	}
}

func TestHTTPInboundContextFunc_Credentials(t *testing.T) {
	tests := []struct {
		name                string
		authHeader          string
		cookie              *http.Cookie
		expectedToken       string
		expectedCookieValue string
	}{
		{
			name:          "bearer token",
			authHeader:    "Bearer test-token-123",
			expectedToken: "test-token-123",
		},
		{
			name:          "token without bearer prefix",
			authHeader:    "test-token-456",
			expectedToken: "test-token-456",
		},
		{
			name: "access token cookie",
			cookie: &http.Cookie{
				Name:  "chrono-accesstoken",
				Value: "cookie-value-789",
			},
			expectedCookieValue: "cookie-value-789",
		},
		{
			name:       "both token and cookie",
			authHeader: "Bearer test-token-abc",
			cookie: &http.Cookie{
				Name:  "chrono-accesstoken",
				Value: "cookie-value-xyz",
			},
			expectedToken:       "test-token-abc",
			expectedCookieValue: "cookie-value-xyz",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			if tt.cookie != nil {
				req.AddCookie(tt.cookie)
			}

			ctx := HTTPInboundContextFunc(t.Context(), req)

			credentials := FetchSessionAPIToken(ctx)
			assert.Equal(t, tt.expectedToken, credentials.APIToken)
			assert.Equal(t, tt.expectedCookieValue, credentials.AccessTokenCookie)
		})
	}
}
