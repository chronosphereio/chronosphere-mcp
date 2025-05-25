package authcontext

import "net/http"

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

	if authToken := FetchSessionAPIToken(req.Context()); authToken != "" {
		// forward the api token from context if available
		req2.Header.Set("Authorization", "Bearer "+authToken)
	} else if r.token != "" {
		req2.Header.Set("Authorization", "Bearer "+r.token)
	}

	return r.transport.RoundTrip(req2)
}
