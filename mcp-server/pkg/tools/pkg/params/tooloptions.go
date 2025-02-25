package params

import (
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

func WithTimeRange() mcp.ToolOption {
	return func(tool *mcp.Tool) {
		mcp.WithString("start",
			mcp.Description("Optional. Start time in RFC3339 format or timestamp in seconds. Defaults to end - 1 hour."),
		)(tool)
		mcp.WithString("end",
			mcp.Description("Optional. End time in RFC3339 format or timestamp in seconds. Defaults to the current time."),
		)(tool)
	}
}

type TimeRange struct {
	Start time.Time
	End   time.Time
}

func ParseTimeRange(request mcp.CallToolRequest) (*TimeRange, error) {
	end, err := Time(request, "end", false, time.Now())
	if err != nil {
		return nil, err
	}
	start, err := Time(request, "start", false, end.Add(-time.Hour))
	if err != nil {
		return nil, err
	}
	return &TimeRange{
		Start: start,
		End:   end,
	}, nil
}
