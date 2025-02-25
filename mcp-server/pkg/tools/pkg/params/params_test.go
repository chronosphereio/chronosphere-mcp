package params

import (
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/require"
)

func TestString(t *testing.T) {
	testCases := []struct {
		name         string
		args         map[string]interface{}
		required     bool
		defaultValue string
		wantErr      string
		wantValue    string
	}{
		{
			name:      "valid string",
			args:      map[string]interface{}{"name": "test value"},
			required:  true,
			wantValue: "test value",
		},
		{
			name:     "missing required parameter",
			args:     map[string]interface{}{},
			required: true,
			wantErr:  "missing required parameter: name",
		},
		{
			name:         "missing optional parameter",
			args:         map[string]interface{}{},
			required:     false,
			defaultValue: "default",
			wantValue:    "default",
		},
		{
			name:         "null value for required parameter",
			args:         map[string]interface{}{"name": nil},
			required:     true,
			defaultValue: "",
			wantErr:      "parameter name cannot be null",
			wantValue:    "",
		},
		{
			name:         "null value for optional parameter",
			args:         map[string]interface{}{"name": nil},
			required:     false,
			defaultValue: "default",
			wantValue:    "default",
		},
		{
			name:     "invalid type (number instead of string)",
			args:     map[string]interface{}{"name": 12345},
			required: true,
			wantErr:  "must be a string",
		},
		{
			name:     "invalid type (boolean instead of string)",
			args:     map[string]interface{}{"name": true},
			required: true,
			wantErr:  "must be a string",
		},
		{
			name:         "empty string is invalid if required",
			args:         map[string]interface{}{"name": ""},
			required:     true,
			defaultValue: "default",
			wantErr:      "empty",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			request := createRequestWithArgs(tc.args)
			result, toolErr := String(request, "name", tc.required, tc.defaultValue)

			if tc.wantErr != "" {
				requireErrorResultContains(t, toolErr, tc.wantErr)
				return
			}

			// Check result value if no error expected
			require.NoError(t, toolErr)
			require.Equal(t, tc.wantValue, result, "result value mismatch")
		})
	}
}

func TestInt(t *testing.T) {
	testCases := []struct {
		name         string
		args         map[string]interface{}
		required     bool
		defaultValue int
		wantErr      string
		wantValue    int
	}{
		{
			name:      "valid int",
			args:      map[string]interface{}{"count": 42},
			required:  true,
			wantValue: 42,
		},
		{
			name:         "valid int64",
			args:         map[string]interface{}{"count": int64(42)},
			required:     true,
			defaultValue: 0,
			wantValue:    42,
		},
		{
			name:         "valid float64 (converted to int)",
			args:         map[string]interface{}{"count": 42.75},
			required:     true,
			defaultValue: 0,
			wantValue:    42, // truncated
		},
		{
			name:         "missing required parameter",
			args:         map[string]interface{}{},
			required:     true,
			defaultValue: 0,
			wantErr:      "missing required parameter: count",
			wantValue:    0,
		},
		{
			name:         "missing optional parameter",
			args:         map[string]interface{}{},
			required:     false,
			defaultValue: 100,
			wantValue:    100,
		},
		{
			name:         "null value for required parameter",
			args:         map[string]interface{}{"count": nil},
			required:     true,
			defaultValue: 0,
			wantErr:      "parameter count cannot be null",
			wantValue:    0,
		},
		{
			name:         "invalid type (string instead of int)",
			args:         map[string]interface{}{"count": "42"},
			required:     true,
			defaultValue: 0,
			wantErr:      "must be a number",
			wantValue:    0,
		},
		{
			name:         "invalid type (boolean instead of int)",
			args:         map[string]interface{}{"count": true},
			required:     true,
			defaultValue: 0,
			wantErr:      "must be a number",
			wantValue:    0,
		},
		{
			name:         "zero is valid",
			args:         map[string]interface{}{"count": 0},
			required:     true,
			defaultValue: 100,
			wantValue:    0,
		},
		{
			name:         "negative number is valid",
			args:         map[string]interface{}{"count": -42},
			required:     true,
			defaultValue: 0,
			wantValue:    -42,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			request := createRequestWithArgs(tc.args)
			result, toolErr := Int(request, "count", tc.required, tc.defaultValue)

			if tc.wantErr != "" {
				requireErrorResultContains(t, toolErr, tc.wantErr)
				return
			}

			require.Nil(t, toolErr, "expected no error but got one")
			require.Equal(t, tc.wantValue, result, "result value mismatch")
		})
	}
}

func TestFloat(t *testing.T) {
	testCases := []struct {
		name         string
		args         map[string]interface{}
		required     bool
		defaultValue float64
		wantErr      string
		wantValue    float64
	}{
		{
			name:         "valid float64",
			args:         map[string]interface{}{"value": 42.75},
			required:     true,
			defaultValue: 0,
			wantValue:    42.75,
		},
		{
			name:         "valid int (converted to float64)",
			args:         map[string]interface{}{"value": 42},
			required:     true,
			defaultValue: 0,
			wantValue:    42.0,
		},
		{
			name:         "valid int64 (converted to float64)",
			args:         map[string]interface{}{"value": int64(42)},
			required:     true,
			defaultValue: 0,
			wantValue:    42.0,
		},
		{
			name:         "missing required parameter",
			args:         map[string]interface{}{},
			required:     true,
			defaultValue: 0,
			wantErr:      "missing required parameter: value",
			wantValue:    0,
		},
		{
			name:         "missing optional parameter",
			args:         map[string]interface{}{},
			required:     false,
			defaultValue: 99.9,
			wantValue:    99.9,
		},
		{
			name:         "null value for required parameter",
			args:         map[string]interface{}{"value": nil},
			required:     true,
			defaultValue: 0,
			wantErr:      "parameter value cannot be null",
			wantValue:    0,
		},
		{
			name:         "invalid type (string instead of float)",
			args:         map[string]interface{}{"value": "42.75"},
			required:     true,
			defaultValue: 0,
			wantErr:      "must be a number",
			wantValue:    0,
		},
		{
			name:         "zero is valid",
			args:         map[string]interface{}{"value": 0.0},
			required:     true,
			defaultValue: 99.9,
			wantValue:    0.0,
		},
		{
			name:         "negative number is valid",
			args:         map[string]interface{}{"value": -42.75},
			required:     true,
			defaultValue: 0,
			wantValue:    -42.75,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			request := createRequestWithArgs(tc.args)
			result, toolErr := Float(request, "value", tc.required, tc.defaultValue)

			if tc.wantErr != "" {
				requireErrorResultContains(t, toolErr, tc.wantErr)
				return
			}

			require.Nil(t, toolErr, "expected no error but got one")
			require.Equal(t, tc.wantValue, result, "result value mismatch")
		})
	}
}

func TestBool(t *testing.T) {
	testCases := []struct {
		name         string
		args         map[string]interface{}
		required     bool
		defaultValue bool
		wantErr      string
		wantValue    bool
	}{
		{
			name:         "valid bool (true)",
			args:         map[string]interface{}{"enabled": true},
			required:     true,
			defaultValue: false,
			wantValue:    true,
		},
		{
			name:         "valid bool (false)",
			args:         map[string]interface{}{"enabled": false},
			required:     true,
			defaultValue: true,
			wantValue:    false,
		},
		{
			name:         "missing required parameter",
			args:         map[string]interface{}{},
			required:     true,
			defaultValue: false,
			wantErr:      "missing required parameter: enabled",
			wantValue:    false,
		},
		{
			name:         "missing optional parameter",
			args:         map[string]interface{}{},
			required:     false,
			defaultValue: true,
			wantValue:    true,
		},
		{
			name:         "null value for required parameter",
			args:         map[string]interface{}{"enabled": nil},
			required:     true,
			defaultValue: false,
			wantErr:      "parameter enabled cannot be null",
			wantValue:    false,
		},
		{
			name:         "invalid type (string instead of bool)",
			args:         map[string]interface{}{"enabled": "true"},
			required:     true,
			defaultValue: false,
			wantErr:      "must be a boolean",
			wantValue:    false,
		},
		{
			name:         "invalid type (number instead of bool)",
			args:         map[string]interface{}{"enabled": 1},
			required:     true,
			defaultValue: false,
			wantErr:      "must be a boolean",
			wantValue:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			request := createRequestWithArgs(tc.args)
			result, toolErr := Bool(request, "enabled", tc.required, tc.defaultValue)

			if tc.wantErr != "" {
				requireErrorResultContains(t, toolErr, tc.wantErr)
				return
			}

			require.Nil(t, toolErr, "expected no error but got one")
			require.Equal(t, tc.wantValue, result, "result value mismatch")
		})
	}
}

func TestTime(t *testing.T) {
	utc := time.UTC
	defaultTime := time.Date(2023, 1, 1, 0, 0, 0, 0, utc)

	testCases := []struct {
		name         string
		args         map[string]interface{}
		required     bool
		defaultValue time.Time
		wantErr      string
		wantValue    time.Time
	}{
		{
			name:      "valid RFC3339 time",
			args:      map[string]interface{}{"timestamp": "2023-04-15T16:30:00Z"},
			required:  true,
			wantValue: time.Date(2023, 4, 15, 16, 30, 0, 0, utc),
		},
		{
			name:      "valid RFC3339 time with timezone",
			args:      map[string]interface{}{"timestamp": "2023-04-15T16:30:00+02:00"},
			required:  true,
			wantValue: time.Date(2023, 4, 15, 16, 30, 0, 0, time.FixedZone("", 7200)),
		},
		{
			name:      "valid RFC3339 time with nanoseconds",
			args:      map[string]interface{}{"timestamp": "2023-04-15T16:30:00.123456789Z"},
			required:  true,
			wantValue: time.Date(2023, 4, 15, 16, 30, 0, 123456789, utc),
		},
		{
			name:     "missing required parameter",
			args:     map[string]interface{}{},
			required: true,
			wantErr:  "missing required parameter: timestamp",
		},
		{
			name:         "missing optional parameter",
			args:         map[string]interface{}{},
			defaultValue: defaultTime,
			required:     false,
			wantValue:    defaultTime,
		},
		{
			name:     "null value for required parameter",
			args:     map[string]interface{}{"timestamp": nil},
			required: true,
			wantErr:  "parameter timestamp cannot be null",
		},
		{
			name:     "invalid type (number instead of string)",
			args:     map[string]interface{}{"timestamp": 12345},
			required: true,
			wantErr:  "must be a string",
		},
		{
			name:     "invalid time format",
			args:     map[string]interface{}{"timestamp": "2023/04/15 16:30:00"},
			required: true,
			wantErr:  "invalid RFC3339 time format",
		},
		{
			name:     "invalid RFC3339 string",
			args:     map[string]interface{}{"timestamp": "not-a-time"},
			required: true,
			wantErr:  "invalid RFC3339 time format",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			request := createRequestWithArgs(tc.args)
			result, toolErr := Time(request, "timestamp", tc.required, tc.defaultValue)

			if tc.wantErr != "" {
				requireErrorResultContains(t, toolErr, tc.wantErr)
				return
			}

			require.Nil(t, toolErr, "expected no error but got one")
			require.Equal(t, tc.wantValue, result, "result value mismatch")
		})
	}
}

func TestParseTime(t *testing.T) {
	fixedNow := time.Date(2025, 4, 25, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		input    string
		expected time.Time
		wantErr  bool
	}{
		{
			name:     "now",
			input:    "now",
			expected: fixedNow,
		},
		{
			name:     "now minus 1 day",
			input:    "now-1d",
			expected: fixedNow.AddDate(0, 0, -1),
		},
		{
			name:     "minus 1 day",
			input:    "-1d",
			expected: fixedNow.AddDate(0, 0, -1),
		},
		{
			name:     "now minus 2 hours",
			input:    "now-2h",
			expected: fixedNow.Add(-2 * time.Hour),
		},
		{
			name:     "minus 2 hours",
			input:    "-2h",
			expected: fixedNow.Add(-2 * time.Hour),
		},
		{
			name:     "now minus 5 minutes",
			input:    "now-5m",
			expected: fixedNow.Add(-5 * time.Minute),
		},
		{
			name:     "minus 5 minutes",
			input:    "-5m",
			expected: fixedNow.Add(-5 * time.Minute),
		},
		{
			name:     "unix timestamp",
			input:    "1714070400",
			expected: time.Unix(1714070400, 0),
		},
		{
			name:     "rFC3339 timestamp",
			input:    "2025-04-25T05:00:00Z",
			expected: time.Date(2025, 4, 25, 5, 0, 0, 0, time.UTC),
		},
		{
			name:    "invalid input",
			input:   "invalid-time",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := parseTime(tc.input, fixedNow)

			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.expected, got)
		})
	}
}

type testObject struct {
	A string `json:"a"`
	B string `json:"b"`
}

func TestObject(t *testing.T) {
	testCases := []struct {
		name         string
		args         map[string]interface{}
		required     bool
		defaultValue testObject
		wantErr      string
		wantValue    testObject
	}{
		{
			name: "single field populated",
			args: map[string]interface{}{"value": map[string]string{
				"a": "populated",
			}},
			required:  true,
			wantValue: testObject{A: "populated"},
		},
		{
			name: "both fields populated",
			args: map[string]interface{}{"value": map[string]string{
				"a": "populated",
				"b": "populated-2",
			}},
			required:  true,
			wantValue: testObject{A: "populated", B: "populated-2"},
		},
		{
			name:     "empty",
			args:     map[string]interface{}{"value": map[string]string{}},
			required: true,
		},
		{
			name:     "missing but required",
			args:     map[string]interface{}{"missing": map[string]string{}},
			required: true,
			wantErr:  "required",
		},
		{
			name:     "missing but okay",
			args:     map[string]interface{}{"missing": map[string]string{}},
			required: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			request := createRequestWithArgs(tc.args)
			result, toolErr := Object[testObject](request, "value", tc.required, tc.defaultValue)

			if tc.wantErr != "" {
				requireErrorResultContains(t, toolErr, tc.wantErr)
				return
			}

			require.Nil(t, toolErr, "expected no error but got one")
			require.Equal(t, tc.wantValue, result, "result value mismatch")
		})
	}
}

func TestObjectArray(t *testing.T) {
	testCases := []struct {
		name         string
		args         map[string]interface{}
		required     bool
		defaultValue []testObject
		wantErr      string
		wantValue    []testObject
	}{
		{
			name: "single element populated",
			args: map[string]interface{}{"value": []map[string]string{
				{
					"a": "populated",
				},
			}},
			required:  true,
			wantValue: []testObject{{A: "populated"}},
		},
		{
			name: "multiple elements populated",
			args: map[string]interface{}{"value": []map[string]string{
				{
					"a": "populated",
				},
				{
					"b": "populated-2",
				},
			}},
			required:  true,
			wantValue: []testObject{{A: "populated"}, {B: "populated-2"}},
		},
		{
			name:      "no elements populated but not required",
			args:      map[string]interface{}{"value": []map[string]string{}},
			wantValue: []testObject{},
		},
		{
			name:     "no elements populated but required",
			args:     map[string]interface{}{"value": []map[string]string{}},
			required: true,
			wantErr:  "required",
		},
		{
			name:     "no key present",
			args:     map[string]interface{}{"missing": []map[string]string{}},
			required: true,
			wantErr:  "required",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			request := createRequestWithArgs(tc.args)
			result, toolErr := ObjectArray[testObject](request, "value", tc.required)

			if tc.wantErr != "" {
				requireErrorResultContains(t, toolErr, tc.wantErr)
				return
			}

			require.Nil(t, toolErr, "expected no error but got one: %s", toolErr)
			require.Equal(t, tc.wantValue, result, "result value mismatch")
		})
	}
}

// Helper function to create a CallToolRequest with specified arguments
func createRequestWithArgs(args map[string]interface{}) mcp.CallToolRequest {
	request := mcp.CallToolRequest{}
	request.Params.Arguments = args
	request.Params.Name = "test_tool"
	return request
}

func requireErrorResultContains(t *testing.T, errResult error, substr string) {
	require.NotNil(t, errResult)
	require.ErrorContains(t, errResult, substr)
}
