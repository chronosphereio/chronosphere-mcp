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
type disabledToolsKey struct{}
type writesEnabledKey struct{}

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

// SetDisabledTools sets the disabled tools set in the context.
func SetDisabledTools(ctx context.Context, disabledTools map[string]struct{}) context.Context {
	return context.WithValue(ctx, disabledToolsKey{}, disabledTools)
}

// FetchDisabledTools retrieves the disabled tools set from the context.
func FetchDisabledTools(ctx context.Context) map[string]struct{} {
	disabledTools, ok := ctx.Value(disabledToolsKey{}).(map[string]struct{})
	if !ok {
		return nil
	}
	return disabledTools
}

// SetWritesEnabled records whether write tools are enabled for the request.
func SetWritesEnabled(ctx context.Context, enabled bool) context.Context {
	return context.WithValue(ctx, writesEnabledKey{}, enabled)
}

// FetchWritesEnabled reports whether write tools are enabled for the request.
func FetchWritesEnabled(ctx context.Context) bool {
	enabled, ok := ctx.Value(writesEnabledKey{}).(bool)
	return ok && enabled
}

// WritesEnabled reports whether the server and, when present, request both enable writes.
// HTTP and SSE middleware always records a request preference. Stdio has no request
// preference, so the server setting alone controls writes for that transport.
func WritesEnabled(ctx context.Context, serverEnabled bool) bool {
	if !serverEnabled {
		return false
	}
	requestEnabled, hasRequestPreference := ctx.Value(writesEnabledKey{}).(bool)
	return !hasRequestPreference || requestEnabled
}
