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

// DataunstableGetEventsForMonitorResponse dataunstable get events for monitor response
//
// swagger:model dataunstableGetEventsForMonitorResponse
type DataunstableGetEventsForMonitorResponse struct {

	// The buckets and individual events making up a time range
	EventHistogramWithDetails []*DataunstableEventHistogramWithDetails `json:"event_histogram_with_details"`

	// The total number of events in the time range
	TotalEvents string `json:"total_events,omitempty"`
}

// Validate validates this dataunstable get events for monitor response
func (m *DataunstableGetEventsForMonitorResponse) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateEventHistogramWithDetails(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *DataunstableGetEventsForMonitorResponse) validateEventHistogramWithDetails(formats strfmt.Registry) error {
	if swag.IsZero(m.EventHistogramWithDetails) { // not required
		return nil
	}

	for i := 0; i < len(m.EventHistogramWithDetails); i++ {
		if swag.IsZero(m.EventHistogramWithDetails[i]) { // not required
			continue
		}

		if m.EventHistogramWithDetails[i] != nil {
			if err := m.EventHistogramWithDetails[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("event_histogram_with_details" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("event_histogram_with_details" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this dataunstable get events for monitor response based on the context it is used
func (m *DataunstableGetEventsForMonitorResponse) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateEventHistogramWithDetails(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *DataunstableGetEventsForMonitorResponse) contextValidateEventHistogramWithDetails(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.EventHistogramWithDetails); i++ {

		if m.EventHistogramWithDetails[i] != nil {

			if swag.IsZero(m.EventHistogramWithDetails[i]) { // not required
				return nil
			}

			if err := m.EventHistogramWithDetails[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("event_histogram_with_details" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("event_histogram_with_details" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *DataunstableGetEventsForMonitorResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DataunstableGetEventsForMonitorResponse) UnmarshalBinary(b []byte) error {
	var res DataunstableGetEventsForMonitorResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
