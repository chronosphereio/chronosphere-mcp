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

// Package authcontext contains authorization/authentication utilities for the MCP server
package authcontext

import "context"

type sessionAPITokenKey struct{}

type SessionCredentials struct {
	APIToken          string
	AccessTokenCookie string
}

func (s SessionCredentials) IsEmpty() bool {
	return s.APIToken == "" && s.AccessTokenCookie == ""
}

// SetSessionCredentials sets the session credentials in the context.
func SetSessionCredentials(ctx context.Context, credentials SessionCredentials) context.Context {
	return context.WithValue(ctx, sessionAPITokenKey{}, credentials)
}

// FetchSessionAPIToken retrieves the session API token from the context.
func FetchSessionAPIToken(ctx context.Context) SessionCredentials {
	credentials, ok := ctx.Value(sessionAPITokenKey{}).(SessionCredentials)
	if !ok {
		return SessionCredentials{}
	}
	return credentials
}
