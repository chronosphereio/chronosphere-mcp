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
	"encoding/hex"

	"github.com/chronosphereio/chronosphere-mcp/generated/datav1/models"
	"github.com/go-openapi/strfmt"
)

// ListTracesResponse mirrors Datav1ListTracesResponse with hex-encoded IDs.
type ListTracesResponse struct {
	Traces []*Data `json:"traces"`
}

// Data mirrors V1TracesData with hex-encoded IDs.
type Data struct {
	ResourceSpans []*ResourceSpans `json:"resource_spans"`
}

// ResourceSpans mirrors V1ResourceSpans.
type ResourceSpans struct {
	Resource   *models.V1Resource `json:"resource,omitempty"`
	SchemaURL  string             `json:"schema_url,omitempty"`
	ScopeSpans []*ScopeSpans      `json:"scope_spans"`
}

// ScopeSpans mirrors V1ScopeSpans.
type ScopeSpans struct {
	Scope     *models.V1InstrumentationScope `json:"scope,omitempty"`
	SchemaURL string                         `json:"schema_url,omitempty"`
	Spans     []*Span                        `json:"spans"`
}

// Span mirrors V1Span but with hex-encoded ID fields.
type Span struct {
	// Converted fields
	TraceID      string `json:"trace_id"`       // 32-char hex (16 bytes)
	SpanID       string `json:"span_id"`        // 16-char hex (8 bytes)
	ParentSpanID string `json:"parent_span_id"` // 16-char hex, empty if root

	// Preserved fields from V1Span
	Name                   string                `json:"name,omitempty"`
	Kind                   models.SpanSpanKind   `json:"kind,omitempty"`
	StartTimeUnixNano      string                `json:"start_time_unix_nano,omitempty"`
	EndTimeUnixNano        string                `json:"end_time_unix_nano,omitempty"`
	Attributes             []*models.V1KeyValue  `json:"attributes"`
	DroppedAttributesCount int64                 `json:"dropped_attributes_count,omitempty"`
	Events                 []*models.V1SpanEvent `json:"events"`
	DroppedEventsCount     int64                 `json:"dropped_events_count,omitempty"`
	Links                  []*SpanLink           `json:"links"`
	DroppedLinksCount      int64                 `json:"dropped_links_count,omitempty"`
	Status                 *models.Tracev1Status `json:"status,omitempty"`
	TraceState             string                `json:"trace_state,omitempty"`
	Flags                  int64                 `json:"flags,omitempty"`
}

// SpanLink mirrors SpanLink with hex-encoded IDs.
type SpanLink struct {
	TraceID                string               `json:"trace_id"`
	SpanID                 string               `json:"span_id"`
	TraceState             string               `json:"trace_state,omitempty"`
	Attributes             []*models.V1KeyValue `json:"attributes"`
	DroppedAttributesCount int64                `json:"dropped_attributes_count,omitempty"`
	Flags                  int64                `json:"flags,omitempty"`
}

// convertToHexResponse converts API response with base64 IDs to hex-encoded IDs.
func convertToHexResponse(apiResp *models.Datav1ListTracesResponse) *ListTracesResponse {
	if apiResp == nil {
		return nil
	}

	result := &ListTracesResponse{
		Traces: make([]*Data, len(apiResp.Traces)),
	}

	// Iterate through all levels
	for i, traceData := range apiResp.Traces {
		if traceData == nil {
			continue
		}

		hexTraceData := &Data{
			ResourceSpans: make([]*ResourceSpans, len(traceData.ResourceSpans)),
		}

		for j, rs := range traceData.ResourceSpans {
			if rs == nil {
				continue
			}

			hexRS := &ResourceSpans{
				Resource:   rs.Resource,
				SchemaURL:  rs.SchemaURL,
				ScopeSpans: make([]*ScopeSpans, len(rs.ScopeSpans)),
			}

			for k, ss := range rs.ScopeSpans {
				if ss == nil {
					continue
				}

				hexSS := &ScopeSpans{
					Scope:     ss.Scope,
					SchemaURL: ss.SchemaURL,
					Spans:     make([]*Span, len(ss.Spans)),
				}

				for l, span := range ss.Spans {
					if span == nil {
						continue
					}

					hexSS.Spans[l] = convertSpanToHex(span)
				}

				hexRS.ScopeSpans[k] = hexSS
			}

			hexTraceData.ResourceSpans[j] = hexRS
		}

		result.Traces[i] = hexTraceData
	}

	return result
}

// convertSpanToHex converts a single span to hex format.
func convertSpanToHex(span *models.V1Span) *Span {
	if span == nil {
		return nil
	}

	// Convert IDs
	traceID := bytesToHex(span.TraceID)
	spanID := bytesToHex(span.SpanID)
	parentSpanID := bytesToHex(span.ParentSpanID)

	// Convert links
	hexLinks := make([]*SpanLink, len(span.Links))
	for i, link := range span.Links {
		if link == nil {
			continue
		}
		hexLink := convertLinkToHex(link)
		hexLinks[i] = hexLink
	}

	return &Span{
		TraceID:                traceID,
		SpanID:                 spanID,
		ParentSpanID:           parentSpanID,
		Name:                   span.Name,
		Kind:                   span.Kind,
		StartTimeUnixNano:      span.StartTimeUnixNano,
		EndTimeUnixNano:        span.EndTimeUnixNano,
		Attributes:             span.Attributes,
		DroppedAttributesCount: span.DroppedAttributesCount,
		Events:                 span.Events,
		DroppedEventsCount:     span.DroppedEventsCount,
		Links:                  hexLinks,
		DroppedLinksCount:      span.DroppedLinksCount,
		Status:                 span.Status,
		TraceState:             span.TraceState,
		Flags:                  span.Flags,
	}
}

// convertLinkToHex converts a span link to hex format.
func convertLinkToHex(link *models.SpanLink) *SpanLink {
	if link == nil {
		return nil
	}

	traceID := bytesToHex(link.TraceID)
	spanID := bytesToHex(link.SpanID)

	return &SpanLink{
		TraceID:                traceID,
		SpanID:                 spanID,
		TraceState:             link.TraceState,
		Attributes:             link.Attributes,
		DroppedAttributesCount: link.DroppedAttributesCount,
		Flags:                  link.Flags,
	}
}

// bytesToHex converts a byte array to a hexadecimal string.
// Returns an empty string if the input is empty or nil.
func bytesToHex(bytes strfmt.Base64) string {
	if len(bytes) == 0 {
		return ""
	}

	// Convert to lowercase hex string
	return hex.EncodeToString(bytes)
}
