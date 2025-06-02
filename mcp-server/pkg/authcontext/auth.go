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

// SetSessionAPIToken sets the session API token in the context.
func SetSessionAPIToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, sessionAPITokenKey{}, token)
}

// FetchSessionAPIToken retrieves the session API token from the context.
func FetchSessionAPIToken(ctx context.Context) string {
	sessionAPIToken, ok := ctx.Value(sessionAPITokenKey{}).(string)
	if !ok {
		return ""
	}
	return sessionAPIToken
}
