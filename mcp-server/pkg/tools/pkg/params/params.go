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

// Package params provides utilities for parsing arguments from MCP tool calls
package params

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

// ParamResult is a generic type that holds either a successfully parsed value or an error result
type ParamResult[T any] struct {
	Value  T
	Result *mcp.CallToolResult
	Err    error
}

// IsError returns true if the parsing resulted in an error
func (pr ParamResult[T]) IsError() bool {
	return pr.Result != nil || pr.Err != nil
}

// ParamParser defines a generic function to convert interface{} to a specific type
type ParamParser[T any] func(param interface{}) (T, error)

// parse is a generic function to parse parameters of any type from the request
// It returns a ParamResult which contains either the parsed value or error information
func parse[T any](request mcp.CallToolRequest, key string, required bool, defaultValue T, parser ParamParser[T]) (T, error) {
	// Check if the parameter exists
	param, ok := request.GetArguments()[key]
	if !ok {
		if required {
			var zero T
			return zero, fmt.Errorf("missing required parameter: %s", key)
		}
		return defaultValue, nil
	}

	// Handle null/nil values
	if param == nil || param == "" {
		if required {
			var zero T
			return zero, fmt.Errorf("parameter %s cannot be null or empty", key)
		}
		return defaultValue, nil
	}

	return parseObject(parser, param, key)
}

func parseObject[T any](parser ParamParser[T], param interface{}, key string) (T, error) {
	// Use the provided parser function to convert to the target type
	value, err := parser(param)
	if err != nil {
		var zero T
		return zero, fmt.Errorf("invalid parameter %s: %s", key, err.Error())
	}
	return value, nil
}

// Parser functions for different types

func stringParser(param interface{}) (string, error) {
	str, ok := param.(string)
	if !ok {
		return "", fmt.Errorf("must be a string, got %T", param)
	}
	return str, nil
}

func intParser(param interface{}) (int, error) {
	switch v := param.(type) {
	case int:
		return v, nil
	case int64:
		return int(v), nil
	case float64:
		return int(v), nil
	default:
		return 0, fmt.Errorf("must be a number, got %T", param)
	}
}

func floatParser(param interface{}) (float64, error) {
	switch v := param.(type) {
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case float64:
		return v, nil
	default:
		return 0, fmt.Errorf("must be a number, got %T", param)
	}
}

func boolParser(param interface{}) (bool, error) {
	b, ok := param.(bool)
	if !ok {
		return false, fmt.Errorf("must be a boolean, got %T", param)
	}
	return b, nil
}

func jsonRoundtripParser[T any](param interface{}) (T, error) {
	bs, err := json.Marshal(param)
	if err != nil {
		var zero T
		return zero, fmt.Errorf("failed to serialize to JSON while parsing object")
	}
	var base T
	if err = json.Unmarshal(bs, &base); err != nil {
		var zero T
		return zero, fmt.Errorf("failed to parse object: %s", err.Error())
	}
	return base, nil
}

var timeOffsetRe = regexp.MustCompile(`^(?:now)?(?:-([0-9]+)([dhm]))?$`)

func parseTimeOffsetString(input string, now time.Time) (time.Time, error) {
	matches := timeOffsetRe.FindStringSubmatch(strings.TrimSpace(input))
	if matches == nil {
		return time.Time{}, fmt.Errorf("invalid time string: %s", input)
	}

	if matches[1] == "" && matches[2] == "" {
		// Just "now"
		return now, nil
	}

	amountStr, unit := matches[1], matches[2]
	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid number: %v", err)
	}

	switch unit {
	case "d":
		return now.AddDate(0, 0, -amount), nil
	case "h":
		return now.Add(-time.Duration(amount) * time.Hour), nil
	case "m":
		return now.Add(-time.Duration(amount) * time.Minute), nil
	default:
		return time.Time{}, fmt.Errorf("unsupported unit: %s", unit)
	}
}

func timeParser(param interface{}) (time.Time, error) {
	timeStr, ok := param.(string)
	if !ok {
		return time.Time{}, fmt.Errorf("must be a string, got %T", param)
	}
	return parseTime(timeStr, time.Now())
}

func parseTime(timeStr string, now time.Time) (time.Time, error) {
	t, err := parseTimeOffsetString(timeStr, now)
	if err == nil {
		return t, nil
	}

	if unixSec, err := strconv.ParseInt(timeStr, 10, 64); err == nil {
		return time.Unix(unixSec, 0), nil
	}

	t, err = time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid RFC3339 time format: %w", err)
	}

	return t, nil
}

// Public API with concise function names

// String parses a string parameter and returns the value and any tool result error
func String(request mcp.CallToolRequest, key string, required bool, defaultValue string) (string, error) {
	return parse(request, key, required, defaultValue, stringParser)
}

// Int parses an integer parameter and returns the value and any tool result error
func Int(request mcp.CallToolRequest, key string, required bool, defaultValue int) (int, error) {
	return parse(request, key, required, defaultValue, intParser)
}

// Float parses a float64 parameter and returns the value and any tool result error
func Float(request mcp.CallToolRequest, key string, required bool, defaultValue float64) (float64, error) {
	return parse(request, key, required, defaultValue, floatParser)
}

// Bool parses a boolean parameter and returns the value and any tool result error
func Bool(request mcp.CallToolRequest, key string, required bool, defaultValue bool) (bool, error) {
	return parse(request, key, required, defaultValue, boolParser)
}

// Object parses an object parameter and returns the value and any tool result error
func Object[T any](request mcp.CallToolRequest, key string, required bool, defaultValue T) (T, error) {
	return parse[T](request, key, required, defaultValue, func(param interface{}) (T, error) {
		return jsonRoundtripParser[T](param)
	})
}

// ObjectArray parses an array of objects parameter and returns the value and any tool result error
func ObjectArray[T any](request mcp.CallToolRequest, key string, required bool) ([]T, error) {
	args := request.GetArguments()
	arrayRaw, ok := args[key]
	if !ok {
		if required {
			return nil, fmt.Errorf("request field %s required", key)
		}
		return nil, nil
	}

	sliceValue, err := parseObject(func(param interface{}) ([]T, error) {
		return jsonRoundtripParser[[]T](param)
	}, arrayRaw, key)
	if err != nil {
		return nil, err
	}

	if len(sliceValue) == 0 && required {
		return nil, fmt.Errorf("missing required parameter: %s", key)
	}
	return sliceValue, nil
}

// Time parses an RFC3339 time parameter and returns the value and any tool result error
func Time(request mcp.CallToolRequest, key string, required bool, defaultValue time.Time) (time.Time, error) {
	return parse(request, key, required, defaultValue, timeParser)
}

// stringArrayParser converts an interface{} to a []string
func stringArrayParser(param interface{}) ([]string, error) {
	switch v := param.(type) {
	case []interface{}:
		// Convert []interface{} to []string
		result := make([]string, len(v))
		for i, item := range v {
			// Try to convert each item to a string
			str, ok := item.(string)
			if !ok {
				return nil, fmt.Errorf("item at index %d must be a string, got %T", i, item)
			}
			result[i] = str
		}
		return result, nil

	case []string:
		// Already a []string, return as is
		return v, nil

	default:
		return nil, fmt.Errorf("must be a string array, got %T", param)
	}
}

func WithStringArray(name string, opts ...mcp.PropertyOption) mcp.ToolOption {
	opts = append(opts, mcp.Items(map[string]any{
		"type": "string",
	}))
	return mcp.WithArray(name, opts...)
}

// StringArray parses a string array parameter and returns the value and any tool result error
func StringArray(request mcp.CallToolRequest, key string, required bool, defaultValue []string) ([]string, error) {
	return parse(request, key, required, defaultValue, stringArrayParser)
}
