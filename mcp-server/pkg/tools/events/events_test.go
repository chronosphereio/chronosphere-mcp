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

package events

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"github.com/chronosphereio/chronosphere-mcp/generated/datav1/datav1"
	"github.com/chronosphereio/chronosphere-mcp/generated/datav1/datav1/version1"
	"github.com/chronosphereio/chronosphere-mcp/pkg/links"
	"github.com/chronosphereio/chronosphere-mcp/pkg/ptr"
)

// TestListEventsAPICall tests that ListEvents can be called successfully.
func TestListEventsAPICall(t *testing.T) {
	// Create a test server that returns a valid response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// Return a minimal valid response
		if _, err := w.Write([]byte(`{
			"events": [],
			"page": {"token": ""}
		}`)); err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	}))
	defer server.Close()

	// Create client pointing to test server
	transport := datav1.DefaultTransportConfig().
		WithHost(server.URL[7:]). // Remove "http://"
		WithSchemes([]string{"http"})

	client := datav1.NewHTTPClientWithConfig(nil, transport)
	logger := zaptest.NewLogger(t)
	linkBuilder := links.NewBuilder("https://test.chronosphere.io")

	tools, err := NewTools(client, nil, logger, linkBuilder)
	require.NoError(t, err)

	// Create query params
	now := time.Now()
	queryParams := &version1.ListEventsParams{
		Context:        t.Context(),
		HappenedAfter:  (*strfmt.DateTime)(ptr.To(now)),
		HappenedBefore: (*strfmt.DateTime)(ptr.To(now)),
		Query:          ptr.To("test query"),
	}

	// Call the API
	resp, err := tools.dataV1API.Version1.ListEvents(queryParams)

	// We expect success since our test server returns a valid response
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Payload)
}
