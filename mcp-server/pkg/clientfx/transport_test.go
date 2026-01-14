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

package clientfx

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type captureRoundTripper struct {
	t              *testing.T
	expectedHeader string
}

func (c captureRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	actual := req.Header.Get("traceparent")
	if actual != c.expectedHeader {
		c.t.Fatalf("expected traceparent %q, got %q", c.expectedHeader, actual)
	}
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader("ok")),
		Header:     make(http.Header),
		Request:    req,
	}, nil
}

func TestRoundTripperInjectsTraceContext(t *testing.T) {
	prevPropagator := otel.GetTextMapPropagator()
	t.Cleanup(func() {
		otel.SetTextMapPropagator(prevPropagator)
	})
	propagator := propagation.TraceContext{}
	otel.SetTextMapPropagator(propagator)

	tp := sdktrace.NewTracerProvider()
	tracer := tp.Tracer("test")
	ctx, span := tracer.Start(t.Context(), "test-span")
	t.Cleanup(func() { span.End() })

	expectedCarrier := propagation.HeaderCarrier{}
	propagator.Inject(ctx, expectedCarrier)
	expected := expectedCarrier.Get("traceparent")

	rt := newRoundTripper(captureRoundTripper{
		t:              t,
		expectedHeader: expected,
	}, "chrono-mcp", "")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://example.com", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	if _, err := rt.RoundTrip(req); err != nil {
		t.Fatalf("round trip failed: %v", err)
	}
}
