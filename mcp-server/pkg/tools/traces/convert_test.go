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

package traces

import (
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/chronosphereio/chronosphere-mcp/generated/datav1/models"
)

func TestBytesToHex(t *testing.T) {
	tests := []struct {
		name  string
		input strfmt.Base64
		want  string
	}{
		{
			name:  "valid 16-byte trace ID",
			input: strfmt.Base64{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10},
			want:  "0102030405060708090a0b0c0d0e0f10",
		},
		{
			name:  "valid 8-byte span ID",
			input: strfmt.Base64{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
			want:  "0102030405060708",
		},
		{
			name:  "empty input",
			input: strfmt.Base64{},
			want:  "",
		},
		{
			name:  "nil input",
			input: nil,
			want:  "",
		},
		{
			name:  "all zeros - 16 bytes",
			input: strfmt.Base64{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			want:  "00000000000000000000000000000000",
		},
		{
			name:  "all zeros - 8 bytes",
			input: strfmt.Base64{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			want:  "0000000000000000",
		},
		{
			name:  "all 0xFF - 16 bytes",
			input: strfmt.Base64{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
			want:  "ffffffffffffffffffffffffffffffff",
		},
		{
			name:  "mixed bytes",
			input: strfmt.Base64{0xde, 0xad, 0xbe, 0xef, 0xca, 0xfe, 0xba, 0xbe},
			want:  "deadbeefcafebabe",
		},
		{
			name:  "arbitrary length (not validated)",
			input: strfmt.Base64{0x01, 0x02, 0x03},
			want:  "010203",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := bytesToHex(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestConvertToHexResponse(t *testing.T) {
	tests := []struct {
		name     string
		input    *models.Datav1ListTracesResponse
		validate func(t *testing.T, result *ListTracesResponse)
	}{
		{
			name:  "nil response",
			input: nil,
			validate: func(t *testing.T, result *ListTracesResponse) {
				assert.Nil(t, result)
			},
		},
		{
			name: "empty traces",
			input: &models.Datav1ListTracesResponse{
				Traces: []*models.V1TracesData{},
			},
			validate: func(t *testing.T, result *ListTracesResponse) {
				assert.NotNil(t, result)
				assert.Empty(t, result.Traces)
			},
		},
		{
			name: "single span with valid IDs",
			input: &models.Datav1ListTracesResponse{
				Traces: []*models.V1TracesData{
					{
						ResourceSpans: []*models.V1ResourceSpans{
							{
								ScopeSpans: []*models.V1ScopeSpans{
									{
										Spans: []*models.V1Span{
											{
												TraceID:      strfmt.Base64{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10},
												SpanID:       strfmt.Base64{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
												ParentSpanID: strfmt.Base64{0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10},
												Name:         "test-span",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			validate: func(t *testing.T, result *ListTracesResponse) {
				require.NotNil(t, result)
				require.Len(t, result.Traces, 1)
				require.Len(t, result.Traces[0].ResourceSpans, 1)
				require.Len(t, result.Traces[0].ResourceSpans[0].ScopeSpans, 1)
				require.Len(t, result.Traces[0].ResourceSpans[0].ScopeSpans[0].Spans, 1)

				span := result.Traces[0].ResourceSpans[0].ScopeSpans[0].Spans[0]
				assert.Equal(t, "0102030405060708090a0b0c0d0e0f10", span.TraceID)
				assert.Equal(t, "0102030405060708", span.SpanID)
				assert.Equal(t, "090a0b0c0d0e0f10", span.ParentSpanID)
				assert.Equal(t, "test-span", span.Name)
			},
		},
		{
			name: "span with empty parent (root span)",
			input: &models.Datav1ListTracesResponse{
				Traces: []*models.V1TracesData{
					{
						ResourceSpans: []*models.V1ResourceSpans{
							{
								ScopeSpans: []*models.V1ScopeSpans{
									{
										Spans: []*models.V1Span{
											{
												TraceID:      strfmt.Base64{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10},
												SpanID:       strfmt.Base64{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
												ParentSpanID: strfmt.Base64{},
												Name:         "root-span",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			validate: func(t *testing.T, result *ListTracesResponse) {
				span := result.Traces[0].ResourceSpans[0].ScopeSpans[0].Spans[0]
				assert.Equal(t, "", span.ParentSpanID)
				assert.Equal(t, "root-span", span.Name)
			},
		},
		{
			name: "span with links",
			input: &models.Datav1ListTracesResponse{
				Traces: []*models.V1TracesData{
					{
						ResourceSpans: []*models.V1ResourceSpans{
							{
								ScopeSpans: []*models.V1ScopeSpans{
									{
										Spans: []*models.V1Span{
											{
												TraceID: strfmt.Base64{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10},
												SpanID:  strfmt.Base64{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
												Links: []*models.SpanLink{
													{
														TraceID: strfmt.Base64{0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f},
														SpanID:  strfmt.Base64{0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			validate: func(t *testing.T, result *ListTracesResponse) {
				span := result.Traces[0].ResourceSpans[0].ScopeSpans[0].Spans[0]
				require.Len(t, span.Links, 1)
				assert.Equal(t, "101112131415161718191a1b1c1d1e1f", span.Links[0].TraceID)
				assert.Equal(t, "2021222324252627", span.Links[0].SpanID)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertToHexResponse(tt.input)
			tt.validate(t, result)
		})
	}
}

func TestConvertSpanToHex(t *testing.T) {
	t.Run("nil span", func(t *testing.T) {
		result := convertSpanToHex(nil)
		assert.Nil(t, result)
	})

	t.Run("valid span", func(t *testing.T) {
		span := &models.V1Span{
			TraceID:           strfmt.Base64{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10},
			SpanID:            strfmt.Base64{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
			ParentSpanID:      strfmt.Base64{0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10},
			Name:              "test-span",
			StartTimeUnixNano: "1234567890",
			EndTimeUnixNano:   "1234567900",
		}

		result := convertSpanToHex(span)
		require.NotNil(t, result)
		assert.Equal(t, "0102030405060708090a0b0c0d0e0f10", result.TraceID)
		assert.Equal(t, "0102030405060708", result.SpanID)
		assert.Equal(t, "090a0b0c0d0e0f10", result.ParentSpanID)
		assert.Equal(t, "test-span", result.Name)
		assert.Equal(t, "1234567890", result.StartTimeUnixNano)
		assert.Equal(t, "1234567900", result.EndTimeUnixNano)
	})
}

func TestConvertLinkToHex(t *testing.T) {
	t.Run("nil link", func(t *testing.T) {
		result := convertLinkToHex(nil)
		assert.Nil(t, result)
	})

	t.Run("valid link", func(t *testing.T) {
		link := &models.SpanLink{
			TraceID:    strfmt.Base64{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10},
			SpanID:     strfmt.Base64{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
			TraceState: "state=value",
		}

		result := convertLinkToHex(link)
		require.NotNil(t, result)
		assert.Equal(t, "0102030405060708090a0b0c0d0e0f10", result.TraceID)
		assert.Equal(t, "0102030405060708", result.SpanID)
		assert.Equal(t, "state=value", result.TraceState)
	})
}
