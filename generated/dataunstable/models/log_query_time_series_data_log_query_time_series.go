// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// LogQueryTimeSeriesDataLogQueryTimeSeries log query time series data log query time series
//
// swagger:model LogQueryTimeSeriesDataLogQueryTimeSeries
type LogQueryTimeSeriesDataLogQueryTimeSeries struct {

	// aggregation_name is by default the name of the aggregation used to calculate
	// the values.
	// In the future we may support aliasing in query grammar.
	AggregationName string `json:"aggregation_name,omitempty"`

	// buckets
	Buckets []*LogQueryTimeSeriesLogQueryTimeSeriesBucket `json:"buckets"`

	// group by dimension values
	GroupByDimensionValues []string `json:"group_by_dimension_values"`
}

// Validate validates this log query time series data log query time series
func (m *LogQueryTimeSeriesDataLogQueryTimeSeries) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateBuckets(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *LogQueryTimeSeriesDataLogQueryTimeSeries) validateBuckets(formats strfmt.Registry) error {
	if swag.IsZero(m.Buckets) { // not required
		return nil
	}

	for i := 0; i < len(m.Buckets); i++ {
		if swag.IsZero(m.Buckets[i]) { // not required
			continue
		}

		if m.Buckets[i] != nil {
			if err := m.Buckets[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("buckets" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("buckets" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this log query time series data log query time series based on the context it is used
func (m *LogQueryTimeSeriesDataLogQueryTimeSeries) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateBuckets(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *LogQueryTimeSeriesDataLogQueryTimeSeries) contextValidateBuckets(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Buckets); i++ {

		if m.Buckets[i] != nil {

			if swag.IsZero(m.Buckets[i]) { // not required
				return nil
			}

			if err := m.Buckets[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("buckets" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("buckets" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *LogQueryTimeSeriesDataLogQueryTimeSeries) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *LogQueryTimeSeriesDataLogQueryTimeSeries) UnmarshalBinary(b []byte) error {
	var res LogQueryTimeSeriesDataLogQueryTimeSeries
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
