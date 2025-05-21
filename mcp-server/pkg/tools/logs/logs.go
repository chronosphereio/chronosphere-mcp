// Package logs provides tools for querying Chronosphere logs.
package logs

import (
	"fmt"

	"github.com/go-openapi/strfmt"
	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"

	"github.com/chronosphereio/mcp-server/generated/dataunstable/dataunstable/data_unstable"
	"github.com/chronosphereio/mcp-server/mcp-server/pkg/client"
	"github.com/chronosphereio/mcp-server/mcp-server/pkg/ptr"
	"github.com/chronosphereio/mcp-server/mcp-server/pkg/tools"
	"github.com/chronosphereio/mcp-server/mcp-server/pkg/tools/pkg/params"
)

var _ tools.MCPTools = (*Tools)(nil)

type Tools struct {
	logger         *zap.Logger
	clientProvider *client.Provider
}

func NewTools(clientProvider *client.Provider, logger *zap.Logger) (*Tools, error) {
	logger.Info("logging tool configured")

	return &Tools{
		logger:         logger,
		clientProvider: clientProvider,
	}, nil
}

func (t *Tools) MCPTools() []tools.MCPTool {
	return []tools.MCPTool{
		{
			Metadata: tools.NewMetadata("list_logs",
				mcp.WithDescription(`List logs from a query. Since log results are quite large, will only receive the first page of logs. use the pageToken to fetch the next page. Consult the Log Query Syntax resource for more details on query syntax`),
				withLogQueryParam(),
				params.WithTimeRange(),
				mcp.WithNumber("page_max_size",
					mcp.Description("Maximum number of logs to return. Default is 10. Use page toke to fetch more pages."),
					mcp.DefaultNumber(10)),
				mcp.WithString("page_token",
					mcp.Description("Page token to fetch the next page of logs"),
				),
			),
			Handler: func(session tools.Session, request mcp.CallToolRequest) (*tools.Result, error) {
				query, err := params.String(request, "query", true, "")
				if err != nil {
					return nil, err
				}

				timeRange, err := params.ParseTimeRange(request)
				if err != nil {
					return nil, err
				}

				pageSize, err := params.Int(request, "page_max_size", false, 10)
				if err != nil {
					return nil, err
				}

				pageToken, err := params.String(request, "page_token", false, "")
				if err != nil {
					return nil, err
				}

				logsAPI, err := t.clientProvider.DataUnstableClient(session)
				if err != nil {
					return nil, err
				}

				queryParams := &data_unstable.ListLogsParams{
					Context:                 session.Context,
					LogFilterQuery:          &query,
					LogFilterHappenedAfter:  (*strfmt.DateTime)(&timeRange.Start),
					LogFilterHappenedBefore: (*strfmt.DateTime)(&timeRange.End),
					PageMaxSize:             ptr.To(int64(pageSize)),
					PageToken:               &pageToken,
				}
				t.logger.Info("list logs", zap.Any("params", queryParams))

				resp, err := logsAPI.DataUnstable.ListLogs(queryParams)
				if err != nil {
					return nil, fmt.Errorf("failed to list logs: %s", err)
				}

				// TODO: summarize logs before returning.
				return &tools.Result{
					JSONContent: resp,
				}, nil
			},
		},
		{
			Metadata: tools.NewMetadata("query_logs_range",
				mcp.WithDescription(`Execute a range query for logs. This endpoint is more efficient and returns logs as either time series or grid data. By default, only the timestamp, message, severity and service fields are returned, you may want to request other properties by using the "project" function e.g. "<query> | project logID,timestamp,message,severity,service". Since log results are quite large, will only receive the first page of logs. use the pageToken to fetch the next page. Consult the Log Query Syntax resource for more details on query syntax`),
				withLogQueryParam(),
				params.WithTimeRange(),
				mcp.WithString("page_token",
					mcp.Description("Page token to fetch the next page of logs. An empty token identifies the first page."),
				),
			),
			Handler: func(session tools.Session, request mcp.CallToolRequest) (*tools.Result, error) {
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

				logsAPI, err := t.clientProvider.DataUnstableClient(session)
				if err != nil {
					return nil, err
				}

				queryParams := &data_unstable.GetRangeQueryParams{
					Context:                       session.Context,
					Query:                         &query,
					TimestampFilterHappenedAfter:  (*strfmt.DateTime)(&timeRange.Start),
					TimestampFilterHappenedBefore: (*strfmt.DateTime)(&timeRange.End),
					PageToken:                     &pageToken,
				}
				t.logger.Info("get range query", zap.Any("params", queryParams))

				resp, err := logsAPI.DataUnstable.GetRangeQuery(queryParams)
				if err != nil {
					return nil, fmt.Errorf("failed to get range query: %s", err)
				}

				return &tools.Result{
					JSONContent: resp,
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
			Handler: func(session tools.Session, request mcp.CallToolRequest) (*tools.Result, error) {
				id, err := params.String(request, "id", true, "")
				if err != nil {
					return nil, err
				}
				queryParams := &data_unstable.ListLogsParams{
					Context:        session.Context,
					LogFilterQuery: ptr.To(fmt.Sprintf(`logID=%q`, id)),
				}
				logsAPI, err := t.clientProvider.DataUnstableClient(session)
				if err != nil {
					return nil, err
				}
				resp, err := logsAPI.DataUnstable.ListLogs(queryParams)
				if err != nil {
					return nil, fmt.Errorf("failed to query log by ID: %s", err)
				}
				return &tools.Result{
					JSONContent: resp,
				}, nil
			},
		},
		{
			Metadata: tools.NewMetadata("get_log_histogram",
				mcp.WithDescription("Get histogram of logs from a given query"),
				withLogQueryParam(),
				params.WithTimeRange(),
				mcp.WithArray("group_by",
					mcp.Description(`Log fields to group results within each bucket. May be "service", "severity" or any label name, or unset to have one group per bucket.`),
				),
			),
			Handler: func(session tools.Session, request mcp.CallToolRequest) (*tools.Result, error) {
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
					Context:                 session.Context,
					LogFilterQuery:          &query,
					LogFilterHappenedAfter:  (*strfmt.DateTime)(&timeRange.Start),
					LogFilterHappenedBefore: (*strfmt.DateTime)(&timeRange.End),
					StepSize:                ptr.To(fmt.Sprintf("%f.0s", stepSize.Seconds())),
					GroupByFieldNames:       groupBy,
				}

				t.logger.Info("get log histogram", zap.Any("params", queryParams))

				logsAPI, err := t.clientProvider.DataUnstableClient(session)
				if err != nil {
					return nil, err
				}
				resp, err := logsAPI.DataUnstable.GetLogHistogram(queryParams)
				if err != nil {
					return nil, fmt.Errorf("failed to get log histogram: %s", err)
				}

				return &tools.Result{
					JSONContent: resp,
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
					mcp.DefaultNumber(100),
				),
			),
			Handler: func(session tools.Session, request mcp.CallToolRequest) (*tools.Result, error) {
				query, err := params.String(request, "query", false, "")
				if err != nil {
					return nil, err
				}

				timeRange, err := params.ParseTimeRange(request)
				if err != nil {
					return nil, err
				}

				limit, err := params.Int(request, "limit", false, 100)
				if err != nil {
					return nil, err
				}

				queryParams := &data_unstable.ListLogFieldNamesParams{
					Context:                 session.Context,
					LogFilterQuery:          &query,
					LogFilterHappenedAfter:  (*strfmt.DateTime)(&timeRange.Start),
					LogFilterHappenedBefore: (*strfmt.DateTime)(&timeRange.End),
					Limit:                   ptr.To(int64(limit)),
				}

				t.logger.Info("list log field names", zap.Any("params", queryParams))

				logsAPI, err := t.clientProvider.DataUnstableClient(session)
				if err != nil {
					return nil, err
				}
				resp, err := logsAPI.DataUnstable.ListLogFieldNames(queryParams)
				if err != nil {
					return nil, fmt.Errorf("failed to list log field names: %s", err)
				}

				return &tools.Result{
					JSONContent: resp,
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
				),
				mcp.WithNumber("limit",
					mcp.Description("Maximum number of field values to return. Default is 100."),
					mcp.DefaultNumber(100),
				),
			),
			Handler: func(session tools.Session, request mcp.CallToolRequest) (*tools.Result, error) {
				query, err := params.String(request, "query", false, "")
				if err != nil {
					return nil, err
				}

				timeRange, err := params.ParseTimeRange(request)
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
					Context:                 session.Context,
					LogFilterQuery:          &query,
					LogFilterHappenedAfter:  (*strfmt.DateTime)(&timeRange.Start),
					LogFilterHappenedBefore: (*strfmt.DateTime)(&timeRange.End),
					FieldName:               ptr.To(fieldName),
					Limit:                   ptr.To(int64(limit)),
				}

				t.logger.Info("list log field values", zap.Any("params", queryParams))

				logsAPI, err := t.clientProvider.DataUnstableClient(session)
				if err != nil {
					return nil, err
				}
				resp, err := logsAPI.DataUnstable.ListLogFieldValues(queryParams)
				if err != nil {
					return nil, fmt.Errorf("failed to list log field values: %s", err)
				}

				return &tools.Result{
					JSONContent: resp,
				}, nil
			},
			// TODO: Implement autocomplete requests.
		},
	}
}

func withLogQueryParam() mcp.ToolOption {
	return mcp.WithString("query",
		mcp.Description(`Query to filter logs e.g. service="gateway" AND level="ERROR"`),
	)
}
