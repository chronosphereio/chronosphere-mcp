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
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"

	"github.com/chronosphereio/chronosphere-mcp/generated/datav1/datav1"
	"github.com/chronosphereio/chronosphere-mcp/generated/datav1/datav1/version1"
	"github.com/chronosphereio/chronosphere-mcp/generated/datav1/models"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools"
	"github.com/chronosphereio/chronosphere-mcp/mcp-server/pkg/tools/pkg/params"
	"github.com/chronosphereio/chronosphere-mcp/pkg/links"
)

var _ tools.MCPTools = (*Tools)(nil)

// getLogDefaultLookback is how far back get_log searches when no time range is given.
const getLogDefaultLookback = 7 * 24 * time.Hour

type Tools struct {
	logger      *zap.Logger
	dataV1API   *datav1.DataV1API
	linkBuilder *links.Builder
}

// NewTools creates the log tools.
func NewTools(
	dataV1API *datav1.DataV1API,
	logger *zap.Logger,
	linkBuilder *links.Builder,
) (*Tools, error) {
	logger.Info("logging tool configured")

	return &Tools{
		logger:      logger,
		dataV1API:   dataV1API,
		linkBuilder: linkBuilder,
	}, nil
}

func (t *Tools) GroupName() string {
	return "logs"
}

// CompactLogSummary represents a compact version of log query results
type CompactLogSummary struct {
	Summary        string                                       `json:"summary"`
	TotalLogs      int                                          `json:"total_logs"`
	Columns        []string                                     `json:"columns,omitempty"`
	Logs           []string                                     `json:"logs,omitempty"`
	TimeSeriesData *models.QueryLogsRangeResponseTimeSeriesData `json:"time_series_data,omitempty"`
	Metadata       *models.QueryLogsRangeResponseMetadata       `json:"metadata,omitempty"`
}

// LogFieldValue is one value of a log field and the number of matching logs that have it.
type LogFieldValue struct {
	Value string `json:"value"`
	Count int64  `json:"count"`
}

// LogFieldValues is the result of list_log_field_values.
type LogFieldValues struct {
	Values []LogFieldValue `json:"values"`
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
the next page. Consult the Log Query Syntax resource for more details on query syntax.
To see every field of a log, pass its logID to the get_log tool.`),
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

				queryParams := &version1.QueryLogsRangeParams{
					Context:         ctx,
					Query:           &query,
					TimeRangeAfter:  (*strfmt.DateTime)(&timeRange.Start),
					TimeRangeBefore: (*strfmt.DateTime)(&timeRange.End),
					PageToken:       &pageToken,
				}
				t.logger.Info("query logs range", zap.Any("params", queryParams))

				resp, err := t.dataV1API.Version1.QueryLogsRange(queryParams)
				if err != nil {
					return nil, fmt.Errorf("failed to query logs range: %s", err)
				}

				resp.Payload, _ = trimLogEntries(resp.Payload, limit, offset)
				jsonContent := createCompactSummary(resp.Payload)

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
				mcp.WithDescription(`Get a full log message by its ID. The ID is the unique identifier for the log.
Searches the last 7 days unless start or end is set; set them to find an older log.`),
				mcp.WithString("id",
					mcp.Description("ID of the log message. This is the logID field in the log message."),
				),
				mcp.WithString("start",
					mcp.Description("Optional. Start time in RFC3339 format or timestamp in seconds. Defaults to end - 7 days."),
				),
				mcp.WithString("end",
					mcp.Description("Optional. End time in RFC3339 format or timestamp in seconds. Defaults to the current time."),
				),
			),
			Handler: func(ctx context.Context, request mcp.CallToolRequest) (*tools.Result, error) {
				id, err := params.String(request, "id", true, "")
				if err != nil {
					return nil, err
				}

				end, err := params.Time(request, "end", false, time.Now())
				if err != nil {
					return nil, err
				}
				start, err := params.Time(request, "start", false, end.Add(-getLogDefaultLookback))
				if err != nil {
					return nil, err
				}

				query := getLogQuery(id)
				resp, err := t.queryLogsRange(ctx, query, start, end)
				if err != nil {
					return nil, fmt.Errorf("failed to query log by ID: %s", err)
				}

				log, err := firstGridRow(resp)
				if err != nil {
					return nil, fmt.Errorf("log %q: %w", id, err)
				}

				return &tools.Result{
					JSONContent: log,
					ChronosphereLink: t.linkBuilder.LogExplorer().
						WithQuery(fmt.Sprintf("logID=%q", id)).
						WithTimeRange(start, end).
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

				// 100 buckets across the range.
				histQuery, err := histogramQuery(query, timeRange.End.Sub(timeRange.Start)/100, groupBy)
				if err != nil {
					return nil, err
				}

				resp, err := t.queryLogsRange(ctx, histQuery, timeRange.Start, timeRange.End)
				if err != nil {
					return nil, fmt.Errorf("failed to get log histogram: %s", err)
				}

				return &tools.Result{
					JSONContent: resp.TimeSeriesData,
					ChronosphereLink: t.linkBuilder.LogExplorer().
						WithQuery(query).
						WithTimeRange(timeRange.Start, timeRange.End).
						String(),
				}, nil
			},
		},
		{
			Metadata: tools.NewMetadata("list_log_field_names",
				mcp.WithDescription("Deprecated: no longer supported, always returns an error. List field names of logs"),
				withLogQueryParam(),
				params.WithTimeRange(),
				mcp.WithNumber("limit",
					mcp.Description("Maximum number of field names to return. Default is 100."),
				),
			),
			Handler: func(_ context.Context, request mcp.CallToolRequest) (*tools.Result, error) {
				query, err := params.String(request, "query", false, "")
				if err != nil {
					return nil, err
				}
				return nil, listLogFieldNamesRetiredError(query)
			},
		},
		{
			Metadata: tools.NewMetadata("list_log_field_values",
				mcp.WithDescription(`List field values of logs, most frequent first, with the number of matching logs that have each value.
To see which fields a log has, pass its logID to the get_log tool.`),
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

				valuesQuery, err := fieldValuesQuery(query, fieldName, limit)
				if err != nil {
					return nil, err
				}

				resp, err := t.queryLogsRange(ctx, valuesQuery, timeRange.Start, timeRange.End)
				if err != nil {
					return nil, fmt.Errorf("failed to list log field values: %s", err)
				}

				values, err := fieldValuesFromGrid(resp)
				if err != nil {
					return nil, fmt.Errorf("failed to list log field values: %w", err)
				}

				return &tools.Result{
					JSONContent: values,
					ChronosphereLink: t.linkBuilder.LogExplorer().
						WithQuery(query).
						WithTimeRange(timeRange.Start, timeRange.End).
						String(),
				}, nil
			},
		},
	}
}

// queryLogsRange runs a log query over [start, end] and returns the response payload.
func (t *Tools) queryLogsRange(ctx context.Context, query string, start, end time.Time) (*models.Datav1QueryLogsRangeResponse, error) {
	queryParams := &version1.QueryLogsRangeParams{
		Context:         ctx,
		Query:           &query,
		TimeRangeAfter:  (*strfmt.DateTime)(&start),
		TimeRangeBefore: (*strfmt.DateTime)(&end),
	}
	t.logger.Info("query logs range", zap.Any("params", queryParams))

	resp, err := t.dataV1API.Version1.QueryLogsRange(queryParams)
	if err != nil {
		return nil, err
	}
	if resp.Payload == nil {
		return &models.Datav1QueryLogsRangeResponse{}, nil
	}
	return resp.Payload, nil
}

// listLogFieldNamesRetiredError is the error list_log_field_names returns on every call.
// query is the caller's query, used to make the suggested query_logs_range call concrete.
func listLogFieldNamesRetiredError(query string) error {
	sampleQuery := "project logID, payload | limit 20"
	if q := strings.TrimSpace(query); q != "" {
		sampleQuery = q + " | " + sampleQuery
	}
	return errors.New(`list_log_field_names is no longer supported.
To approximate it:
1. Call query_logs_range with query ` + "`" + sampleQuery + "`" + ` over the same time range.
2. Collect the keys that appear across the returned payloads. Fields present on more of the sampled logs are
   the common ones.
3. To confirm a candidate field and see its values, call list_log_field_values with that field_name.
This finds fields in the raw log body. It can miss labels added by the ingestion pipeline and OTLP attributes
that are not in the payload, so if you know a field name, try list_log_field_values directly.`)
}

var (
	// bareFieldRE matches field names the query grammar accepts unquoted: dot-separated
	// segments of [A-Za-z0-9_/@-] that do not start with "-".
	bareFieldRE = regexp.MustCompile(`^[A-Za-z0-9_/@][A-Za-z0-9_/@-]*(\.[A-Za-z0-9_/@][A-Za-z0-9_/@-]*)*$`)
	// bracketFieldRE matches field names the grammar accepts inside ['...'].
	bracketFieldRE = regexp.MustCompile(`^[A-Za-z0-9_/@.: -]+$`)
	// queryKeywords are words that would be read as part of the query rather than as a field name.
	queryKeywords = map[string]struct{}{
		"and": {}, "or": {}, "NOT": {}, "by": {}, "step": {}, "exists": {}, "asc": {}, "desc": {}, "true": {}, "false": {},
	}
)

// quoteField renders a field name as a query identifier, using the grammar's ['...']
// quoting for names that are not bare identifiers.
func quoteField(name string) (string, error) {
	if _, ok := queryKeywords[name]; !ok && bareFieldRE.MatchString(name) {
		return name, nil
	}
	if bracketFieldRE.MatchString(name) {
		return "['" + name + "']", nil
	}
	return "", fmt.Errorf("field name %q contains characters that cannot be used in a log query", name)
}

// pipe appends operator to query, or returns operator alone when query is empty.
func pipe(query, operator string) string {
	if q := strings.TrimSpace(query); q != "" {
		return q + " | " + operator
	}
	return operator
}

// getLogQuery returns the query that selects one log by ID with its full content.
func getLogQuery(id string) string {
	return fmt.Sprintf("logID = %q | project logID, _timestamp, service, severity, message, payload", id)
}

// histogramQuery returns a query that counts logs matching query in buckets of step, grouped by groupBy.
func histogramQuery(query string, step time.Duration, groupBy []string) (string, error) {
	stepSeconds := int64(step / time.Second)
	if stepSeconds < 1 {
		stepSeconds = 1
	}
	op := fmt.Sprintf("make-series count() step %ds", stepSeconds)
	if len(groupBy) > 0 {
		fields := make([]string, len(groupBy))
		for i, f := range groupBy {
			quoted, err := quoteField(f)
			if err != nil {
				return "", err
			}
			fields[i] = quoted
		}
		op += " by " + strings.Join(fields, ", ")
	}
	return pipe(query, op), nil
}

// fieldValuesQuery returns a query that counts the logs matching query per value of field,
// most frequent first.
func fieldValuesQuery(query, field string, limit int) (string, error) {
	quoted, err := quoteField(field)
	if err != nil {
		return "", err
	}
	if limit <= 0 {
		return "", fmt.Errorf("limit must be positive, got %d", limit)
	}
	return pipe(query, fmt.Sprintf("filter %s exists | summarize n = count() by %s | sort by n desc | limit %d", quoted, quoted, limit)), nil
}

// rowValueString renders a grid cell as a string.
func rowValueString(val *models.QueryLogsRangeResponseRowValue) string {
	switch {
	case val == nil:
		return ""
	case val.StringValue != "":
		return val.StringValue
	case val.FloatValue != 0:
		return fmt.Sprintf("%g", val.FloatValue)
	default:
		return fmt.Sprintf("%t", val.BoolValue)
	}
}

// rowValue returns a grid cell as the Go value of its column's type.
func rowValue(col *models.QueryLogsRangeResponseColumnMeta, val *models.QueryLogsRangeResponseRowValue) any {
	if val == nil {
		return nil
	}
	if col != nil {
		switch col.Type {
		case models.QueryLogsRangeResponseColumnMetaTypeSTRING:
			return val.StringValue
		case models.QueryLogsRangeResponseColumnMetaTypeFLOAT:
			return val.FloatValue
		case models.QueryLogsRangeResponseColumnMetaTypeBOOLEAN:
			return val.BoolValue
		}
	}
	return rowValueString(val)
}

// firstGridRow returns the first row of a grid response keyed by column name.
func firstGridRow(resp *models.Datav1QueryLogsRangeResponse) (map[string]any, error) {
	if resp.GridData == nil || len(resp.GridData.Rows) == 0 || resp.GridData.Rows[0] == nil {
		return nil, errors.New("not found in the time range")
	}
	row := make(map[string]any, len(resp.GridData.Columns))
	for i, col := range resp.GridData.Columns {
		if col == nil || i >= len(resp.GridData.Rows[0].Values) {
			continue
		}
		row[col.Name] = rowValue(col, resp.GridData.Rows[0].Values[i])
	}
	return row, nil
}

// fieldValuesFromGrid reshapes the (value, n) grid of a fieldValuesQuery into a values list.
func fieldValuesFromGrid(resp *models.Datav1QueryLogsRangeResponse) (*LogFieldValues, error) {
	result := &LogFieldValues{Values: []LogFieldValue{}}
	if resp.GridData == nil {
		return result, nil
	}
	cols := resp.GridData.Columns
	if len(cols) != 2 {
		return nil, fmt.Errorf("expected 2 columns in the response, got %d", len(cols))
	}
	// The count column is the one aliased n; if the field itself is named n, take the last.
	countIdx := 1
	if cols[0] != nil && cols[0].Name == "n" && (cols[1] == nil || cols[1].Name != "n") {
		countIdx = 0
	}
	valueIdx := 1 - countIdx

	for _, row := range resp.GridData.Rows {
		if row == nil || len(row.Values) != 2 {
			continue
		}
		var count int64
		if v := row.Values[countIdx]; v != nil {
			count = int64(v.FloatValue)
		}
		result.Values = append(result.Values, LogFieldValue{
			Value: rowValueString(row.Values[valueIdx]),
			Count: count,
		})
	}
	return result, nil
}

// trimLogEntries trims log entries to at most `limit` rows with an offset. It returns a trimmed payload along with the total number of entries trimmed.
func trimLogEntries(payload *models.Datav1QueryLogsRangeResponse, limit int, offset int) (*models.Datav1QueryLogsRangeResponse, int) {
	if payload == nil || (limit == 0 && offset == 0) {
		return payload, 0
	}

	trimmed := &models.Datav1QueryLogsRangeResponse{
		Metadata: payload.Metadata,
	}

	totalTrimmed := 0

	if payload.GridData != nil {
		trimmed.GridData = &models.QueryLogsRangeResponseGridData{
			Columns: payload.GridData.Columns,
		}

		totalRows := len(payload.GridData.Rows)
		if offset >= totalRows {
			trimmed.GridData.Rows = []*models.QueryLogsRangeResponseRow{}
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
	} else if payload.TimeSeriesData.GroupByDimensionNames != nil || payload.TimeSeriesData.Series != nil {
		trimmed.TimeSeriesData = &models.QueryLogsRangeResponseTimeSeriesData{
			GroupByDimensionNames: payload.TimeSeriesData.GroupByDimensionNames,
			Series:                make([]*models.TimeSeriesDataTimeSeries, len(payload.TimeSeriesData.Series)),
		}

		for i, series := range payload.TimeSeriesData.Series {
			trimmedSeries := &models.TimeSeriesDataTimeSeries{
				AggregationName:        series.AggregationName,
				GroupByDimensionValues: series.GroupByDimensionValues,
			}

			totalBuckets := len(series.Buckets)
			if offset >= totalBuckets {
				trimmedSeries.Buckets = []*models.TimeSeriesBucket{}
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

// createCompactSummary creates a compact summary of log query results from the stable API
func createCompactSummary(resp *models.Datav1QueryLogsRangeResponse) *CompactLogSummary {
	summary := &CompactLogSummary{}

	// Extract metadata
	if resp.Metadata != nil {
		summary.Metadata = &models.QueryLogsRangeResponseMetadata{
			LimitEnforced: resp.Metadata.LimitEnforced,
			Page:          resp.Metadata.Page,
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
						// Replace pipe characters with Unicode pipe-like character to avoid conflicts
						values[j] = strings.ReplaceAll(rowValueString(val), "|", "￨")
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
		summary.Summary = strings.Join(summaryParts, ". ")
	} else if resp.TimeSeriesData.GroupByDimensionNames != nil || resp.TimeSeriesData.Series != nil {
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

func withLogQueryParam() mcp.ToolOption {
	return mcp.WithString("query",
		mcp.Description(`Query to filter logs e.g. service="gateway" AND level="ERROR"`),
	)
}
