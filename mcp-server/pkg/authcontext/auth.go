// Package authcontext contains authorization/authentication utilities for the MCP server
package authcontext

import "context"

type sessionAPITokenKey struct{}

func SetSessionAPIToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, sessionAPITokenKey{}, token)
}

func FetchSessionAPIToken(ctx context.Context) string {
	sessionAPIToken, ok := ctx.Value(sessionAPITokenKey{}).(string)
	if !ok {
		return ""
	}
	return sessionAPIToken
}
