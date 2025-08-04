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

// Package logs provides tools for querying Chronosphere logs.
package logs

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-openapi/strfmt"
	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"

	"github.com/chronosphereio/chronosphere-mcp/generated/dataunstable/dataunstable"
	"github.com/chronosphereio/chronosphere-mcp/generated/dataunstable/dataunstable/data_unstable"
	"github.com/chronosphereio/chronosphere-mcp/generated/dataunstable/models"
	"github.com/chronosphereio/chronosphere-mcp/generated/datav1/datav1"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools/pkg/params"
	"github.com/chronosphereio/chronosphere-mcp/pkg/links"
	"github.com/chronosphereio/chronosphere-mcp/pkg/ptr"
)

var _ tools.MCPTools = (*Tools)(nil)

type Tools struct {
	logger          *zap.Logger
	dataV1API       *datav1.DataV1API
	dataUnstableAPI *dataunstable.DataUnstableAPI
	linkBuilder     *links.Builder
}

func NewTools(
	dataV1API *datav1.DataV1API,
	dataUnstableAPI *dataunstable.DataUnstableAPI,
	logger *zap.Logger,
	linkBuilder *links.Builder,
) (*Tools, error) {
	logger.Info("logging tool configured")

	return &Tools{
		logger:          logger,
		dataV1API:       dataV1API,
		dataUnstableAPI: dataUnstableAPI,
		linkBuilder:     linkBuilder,
	}, nil
}

func (t *Tools) GroupName() string {
	return "logs"
}

// CompactLogSummary represents a compact version of log query results
type CompactLogSummary struct {
	Summary         string                                          `json:"summary"`
	TotalLogs       int                                             `json:"total_logs"`
	Columns         []string                                        `json:"columns,omitempty"`
	AvailableFields []string                                        `json:"available_fields,omitempty"`
	Logs            []string                                        `json:"logs,omitempty"`
	TimeSeriesData  *models.DataunstableLogQueryTimeSeriesData      `json:"time_series_data,omitempty"`
	Metadata        *models.GetRangeQueryResponseRangeQueryMetadata `json:"metadata,omitempty"`
	FieldSuggestion string                                          `json:"field_suggestion,omitempty"`
}

// createCompactSummary creates a compact summary of log query results
func (t *Tools) createCompactSummary(ctx context.Context, query string, timeRange *params.TimeRange, resp *models.DataunstableGetRangeQueryResponse) *CompactLogSummary {
	summary := &CompactLogSummary{
		Metadata: resp.Metadata,
	}

	// Fetch available field names to help with query refinement
	fieldParams := &data_unstable.ListLogFieldNamesParams{
		Context:                 ctx,
		LogFilterQuery:          &query,
		LogFilterHappenedAfter:  (*strfmt.DateTime)(&timeRange.Start),
		LogFilterHappenedBefore: (*strfmt.DateTime)(&timeRange.End),
		Limit:                   ptr.To(int64(50)), // Limit to 50 fields to avoid overwhelming
	}
	if fieldResp, fieldErr := t.dataUnstableAPI.DataUnstable.ListLogFieldNames(fieldParams); fieldErr == nil && fieldResp.Payload != nil {
		if fieldResp.Payload.Suggestions != nil {
			// Extract field names from suggestions
			fieldNames := make([]string, len(fieldResp.Payload.Suggestions))
			for i, suggestion := range fieldResp.Payload.Suggestions {
				if suggestion != nil {
					fieldNames[i] = suggestion.Value
				}
			}
			summary.AvailableFields = fieldNames
			// Create a helpful suggestion for using the project keyword
			if len(summary.AvailableFields) > 4 {
				summary.FieldSuggestion = fmt.Sprintf("To get more details, try adding '| project logID,_timestamp,severity,message,%s' to your query to include additional fields. You can inspect full individual logs using the logID and the get_log tool.", strings.Join(summary.AvailableFields[:4], ","))
			}
		}
	}

	if resp.GridData != nil {
		// Extract column names
		columnNames := make([]string, len(resp.GridData.Columns))
		for i, col := range resp.GridData.Columns {
			if col != nil {
				columnNames[i] = col.Name
			}
		}
		summary.Columns = columnNames
		summary.TotalLogs = len(resp.GridData.Rows)

		// Create pipe-separated logs (all rows)
		var logLines []string
		// Add header row
		logLines = append(logLines, strings.Join(columnNames, "|"))

		// Add all data rows
		for _, row := range resp.GridData.Rows {
			if row != nil {
				values := make([]string, len(columnNames))
				for j, val := range row.Values {
					if j < len(columnNames) && val != nil {
						// Extract the actual value based on type and escape pipes
						var strVal string
						if val.StringValue != "" {
							strVal = val.StringValue
						} else if val.FloatValue != 0 {
							strVal = fmt.Sprintf("%g", val.FloatValue)
						} else {
							strVal = fmt.Sprintf("%t", val.BoolValue)
						}
						// Replace pipe characters with Unicode pipe-like character to avoid conflicts
						values[j] = strings.ReplaceAll(strVal, "|", "￨")
					}
				}
				logLines = append(logLines, strings.Join(values, "|"))
			}
		}
		summary.Logs = logLines

		// Create summary text
		var summaryParts []string
		summaryParts = append(summaryParts, fmt.Sprintf("Found %d log entries", summary.TotalLogs))
		if len(columnNames) > 0 {
			summaryParts = append(summaryParts, fmt.Sprintf("Columns: %s", strings.Join(columnNames, ", ")))
		}
		if summary.TotalLogs > 0 {
			summaryParts = append(summaryParts, "All entries shown below in pipe-separated format")
		}
		if len(summary.AvailableFields) > len(columnNames) {
			summaryParts = append(summaryParts, fmt.Sprintf("%d additional fields available", len(summary.AvailableFields)-len(columnNames)))
		}
		summary.Summary = strings.Join(summaryParts, ". ")
	} else if resp.TimeSeriesData != nil {
		// Include full time series data since it's typically not verbose
		summary.TimeSeriesData = resp.TimeSeriesData
		summary.TotalLogs = 0
		if resp.TimeSeriesData.Series != nil {
			for _, series := range resp.TimeSeriesData.Series {
				if series.Buckets != nil {
					summary.TotalLogs += len(series.Buckets)
				}
			}
		}
		summary.Summary = fmt.Sprintf("Time series data with %d data points across %d series", summary.TotalLogs, len(resp.TimeSeriesData.Series))
	} else {
		summary.Summary = "No log data returned"
	}

	return summary
}

func (t *Tools) MCPTools() []tools.MCPTool {
	return []tools.MCPTool{
		{
			Metadata: tools.NewMetadata("query_logs_range",
				mcp.WithDescription(`Execute a range query for logs.
This endpoint returns logs as either timeSeries or gridData. It may return a large amount of data,
so be careful putting the result of this direction into context. Use offset and limit parameters to
limit the amount of data returned.

By default, only the timestamp, message, severity and service fields are returned, 
you may want to request other properties by using the "project" function e.g. 

"<query> | project logID,timestamp,message,severity,service"

Since log results are quite large, will only receive the first page of logs. use the pageToken to fetch
the next page. Consult the Log Query Syntax resource for more details on query syntax`),
				withLogQueryParam(),
				params.WithTimeRange(),
				mcp.WithString("page_token",
					mcp.Description("Page token to fetch the next page of logs. An empty token identifies the first page."),
				),
				mcp.WithNumber("limit",
					mcp.Description("limit the number of results to return")),
				mcp.WithNumber("offset",
					mcp.Description("skip `offset` number of results before returning")),
			),
			Handler: func(ctx context.Context, request mcp.CallToolRequest) (*tools.Result, error) {
				query, err := params.String(request, "query", true, "")
				if err != nil {
					return nil, err
				}

				timeRange, err := params.ParseTimeRange(request)
				if err != nil {
					return nil, err
				}

				pageToken, err := params.String(request, "page_token", false, "")
				if err != nil {
					return nil, err
				}

				limit, err := params.Int(request, "limit", false, 0)
				if err != nil {
					return nil, err
				}
				offset, err := params.Int(request, "offset", false, 0)
				if err != nil {
					return nil, err
				}

				queryParams := &data_unstable.GetRangeQueryParams{
					Context:                       ctx,
					Query:                         &query,
					TimestampFilterHappenedAfter:  (*strfmt.DateTime)(&timeRange.Start),
					TimestampFilterHappenedBefore: (*strfmt.DateTime)(&timeRange.End),
					PageToken:                     &pageToken,
				}
				t.logger.Info("get range query", zap.Any("params", queryParams))

				resp, err := t.dataUnstableAPI.DataUnstable.GetRangeQuery(queryParams)
				if err != nil {
					return nil, fmt.Errorf("failed to get range query: %s", err)
				}

				resp.Payload, _ = trimLogEntries(resp.Payload, limit, offset)
				jsonContent := t.createCompactSummary(ctx, query, timeRange, resp.Payload)

				return &tools.Result{
					JSONContent: jsonContent,
					ChronosphereLink: t.linkBuilder.LogExplorer().
						WithQuery(query).
						WithTimeRange(timeRange.Start, timeRange.End).
						String(),
				}, nil
			},
		},
		{
			Metadata: tools.NewMetadata("get_log",
				mcp.WithDescription(`Get a full log message by its ID. The ID is the unique identifier for the log.`),
				mcp.WithString("id",
					mcp.Description("ID of the log message. This is the logID field in the log message."),
				),
			),
			Handler: func(ctx context.Context, request mcp.CallToolRequest) (*tools.Result, error) {
				id, err := params.String(request, "id", true, "")
				if err != nil {
					return nil, err
				}
				queryParams := &data_unstable.ListLogsParams{
					Context:        ctx,
					LogFilterQuery: ptr.To(fmt.Sprintf(`logID=%q`, id)),
				}

				resp, err := t.dataUnstableAPI.DataUnstable.ListLogs(queryParams)
				if err != nil {
					return nil, fmt.Errorf("failed to query log by ID: %s", err)
				}
				return &tools.Result{
					JSONContent: resp,
					ChronosphereLink: t.linkBuilder.LogExplorer().
						WithQuery(fmt.Sprintf("logID=%q", id)).
						String(),
				}, nil
			},
		},
		{
			Metadata: tools.NewMetadata("get_log_histogram",
				mcp.WithDescription("Get histogram of logs from a given query"),
				withLogQueryParam(),
				params.WithTimeRange(),
				params.WithStringArray("group_by",
					mcp.Description(`Log fields to group results within each bucket. May be "service", "severity" or any label name, or unset to have one group per bucket.`),
				),
			),
			Handler: func(ctx context.Context, request mcp.CallToolRequest) (*tools.Result, error) {
				timeRange, err := params.ParseTimeRange(request)
				if err != nil {
					return nil, err
				}

				query, err := params.String(request, "query", false, "")
				if err != nil {
					return nil, err
				}

				groupBy, err := params.StringArray(request, "group_by", false, nil)
				if err != nil {
					return nil, err
				}

				// Calculate step size for 100 buckets
				stepSize := timeRange.End.Sub(timeRange.Start) / 100
				queryParams := &data_unstable.GetLogHistogramParams{
					Context:                 ctx,
					LogFilterQuery:          &query,
					LogFilterHappenedAfter:  (*strfmt.DateTime)(&timeRange.Start),
					LogFilterHappenedBefore: (*strfmt.DateTime)(&timeRange.End),
					StepSize:                ptr.To(fmt.Sprintf("%.1fs", stepSize.Seconds())),
					GroupByFieldNames:       groupBy,
				}

				t.logger.Info("get log histogram", zap.Any("params", queryParams))

				resp, err := t.dataUnstableAPI.DataUnstable.GetLogHistogram(queryParams)
				if err != nil {
					return nil, fmt.Errorf("failed to get log histogram: %s", err)
				}

				return &tools.Result{
					JSONContent: resp,
					ChronosphereLink: t.linkBuilder.LogExplorer().
						WithQuery(query).
						WithTimeRange(timeRange.Start, timeRange.End).
						String(),
				}, nil
			},
		},
		{
			Metadata: tools.NewMetadata("list_log_field_names",
				mcp.WithDescription("List field names of logs"),
				withLogQueryParam(),
				params.WithTimeRange(),
				mcp.WithNumber("limit",
					mcp.Description("Maximum number of field names to return. Default is 100."),
				),
			),
			Handler: func(ctx context.Context, request mcp.CallToolRequest) (*tools.Result, error) {
				timeRange, err := params.ParseTimeRange(request)
				if err != nil {
					return nil, err
				}

				query, err := params.String(request, "query", false, "")
				if err != nil {
					return nil, err
				}

				limit, err := params.Int(request, "limit", false, 100)
				if err != nil {
					return nil, err
				}

				queryParams := &data_unstable.ListLogFieldNamesParams{
					Context:                 ctx,
					LogFilterQuery:          &query,
					LogFilterHappenedAfter:  (*strfmt.DateTime)(&timeRange.Start),
					LogFilterHappenedBefore: (*strfmt.DateTime)(&timeRange.End),
					Limit:                   ptr.To(int64(limit)),
				}

				resp, err := t.dataUnstableAPI.DataUnstable.ListLogFieldNames(queryParams)
				if err != nil {
					return nil, fmt.Errorf("failed to list log field names: %s", err)
				}

				return &tools.Result{
					JSONContent: resp,
					ChronosphereLink: t.linkBuilder.LogExplorer().
						WithQuery(query).
						WithTimeRange(timeRange.Start, timeRange.End).
						String(),
				}, nil
			},
		},
		{
			Metadata: tools.NewMetadata("list_log_field_values",
				mcp.WithDescription("List field values of logs"),
				withLogQueryParam(),
				params.WithTimeRange(),
				mcp.WithString("field_name",
					mcp.Description("Field name for listing values"),
					mcp.Required(),
				),
				mcp.WithNumber("limit",
					mcp.Description("Maximum number of field values to return. Default is 100."),
				),
			),
			Handler: func(ctx context.Context, request mcp.CallToolRequest) (*tools.Result, error) {
				timeRange, err := params.ParseTimeRange(request)
				if err != nil {
					return nil, err
				}

				query, err := params.String(request, "query", false, "")
				if err != nil {
					return nil, err
				}

				fieldName, err := params.String(request, "field_name", true, "")
				if err != nil {
					return nil, err
				}

				limit, err := params.Int(request, "limit", false, 100)
				if err != nil {
					return nil, err
				}

				queryParams := &data_unstable.ListLogFieldValuesParams{
					Context:                 ctx,
					LogFilterQuery:          &query,
					LogFilterHappenedAfter:  (*strfmt.DateTime)(&timeRange.Start),
					LogFilterHappenedBefore: (*strfmt.DateTime)(&timeRange.End),
					FieldName:               ptr.To(fieldName),
					Limit:                   ptr.To(int64(limit)),
				}

				resp, err := t.dataUnstableAPI.DataUnstable.ListLogFieldValues(queryParams)
				if err != nil {
					return nil, fmt.Errorf("failed to list log field values: %s", err)
				}

				return &tools.Result{
					JSONContent: resp,
					ChronosphereLink: t.linkBuilder.LogExplorer().
						WithQuery(query).
						WithTimeRange(timeRange.Start, timeRange.End).
						String(),
				}, nil
			},
		},
	}
}

// trimLogEntries trims log entries to at most `limit` rows with an offset. It returns a trimmed payload along with the total number of entries trimmed.
func trimLogEntries(payload *models.DataunstableGetRangeQueryResponse, limit int, offset int) (*models.DataunstableGetRangeQueryResponse, int) {
	if payload == nil || (limit == 0 && offset == 0) {
		return payload, 0
	}

	trimmed := &models.DataunstableGetRangeQueryResponse{
		Metadata: payload.Metadata,
	}

	totalTrimmed := 0

	if payload.GridData != nil {
		trimmed.GridData = &models.DataunstableLogQueryGridData{
			Columns: payload.GridData.Columns,
		}

		totalRows := len(payload.GridData.Rows)
		if offset >= totalRows {
			trimmed.GridData.Rows = []*models.DataunstableRow{}
			totalTrimmed = totalRows
		} else {
			end := totalRows
			if limit > 0 {
				end = offset + limit
				if end > totalRows {
					end = totalRows
				}
			}
			trimmed.GridData.Rows = payload.GridData.Rows[offset:end]
			totalTrimmed = offset + (totalRows - end)
		}
	} else if payload.TimeSeriesData != nil {
		trimmed.TimeSeriesData = &models.DataunstableLogQueryTimeSeriesData{
			GroupByDimensionNames: payload.TimeSeriesData.GroupByDimensionNames,
			Series:                make([]*models.LogQueryTimeSeriesDataLogQueryTimeSeries, len(payload.TimeSeriesData.Series)),
		}

		for i, series := range payload.TimeSeriesData.Series {
			trimmedSeries := &models.LogQueryTimeSeriesDataLogQueryTimeSeries{
				AggregationName:        series.AggregationName,
				GroupByDimensionValues: series.GroupByDimensionValues,
			}

			totalBuckets := len(series.Buckets)
			if offset >= totalBuckets {
				trimmedSeries.Buckets = []*models.LogQueryTimeSeriesLogQueryTimeSeriesBucket{}
				totalTrimmed += totalBuckets
			} else {
				end := totalBuckets
				if limit > 0 {
					end = offset + limit
					if end > totalBuckets {
						end = totalBuckets
					}
				}
				trimmedSeries.Buckets = series.Buckets[offset:end]
				totalTrimmed += offset + (totalBuckets - end)
			}

			trimmed.TimeSeriesData.Series[i] = trimmedSeries
		}
	}

	return trimmed, totalTrimmed
}

func withLogQueryParam() mcp.ToolOption {
	return mcp.WithString("query",
		mcp.Description(`Query to filter logs e.g. service="gateway" AND level="ERROR"`),
	)
}
