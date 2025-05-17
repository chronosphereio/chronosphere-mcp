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
