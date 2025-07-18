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

// DataunstableLogQueryTimeSeriesData dataunstable log query time series data
//
// swagger:model dataunstableLogQueryTimeSeriesData
type DataunstableLogQueryTimeSeriesData struct {

	// The names of the dimensions by which the results are grouped by.
	GroupByDimensionNames []string `json:"group_by_dimension_names"`

	// series
	Series []*LogQueryTimeSeriesDataLogQueryTimeSeries `json:"series"`
}

// Validate validates this dataunstable log query time series data
func (m *DataunstableLogQueryTimeSeriesData) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateSeries(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *DataunstableLogQueryTimeSeriesData) validateSeries(formats strfmt.Registry) error {
	if swag.IsZero(m.Series) { // not required
		return nil
	}

	for i := 0; i < len(m.Series); i++ {
		if swag.IsZero(m.Series[i]) { // not required
			continue
		}

		if m.Series[i] != nil {
			if err := m.Series[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("series" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("series" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this dataunstable log query time series data based on the context it is used
func (m *DataunstableLogQueryTimeSeriesData) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateSeries(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *DataunstableLogQueryTimeSeriesData) contextValidateSeries(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Series); i++ {

		if m.Series[i] != nil {

			if swag.IsZero(m.Series[i]) { // not required
				return nil
			}

			if err := m.Series[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("series" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("series" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *DataunstableLogQueryTimeSeriesData) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DataunstableLogQueryTimeSeriesData) UnmarshalBinary(b []byte) error {
	var res DataunstableLogQueryTimeSeriesData
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
