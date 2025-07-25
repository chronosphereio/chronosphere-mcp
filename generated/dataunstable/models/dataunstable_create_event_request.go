// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// DataunstableCreateEventRequest dataunstable create event request
//
// swagger:model dataunstableCreateEventRequest
type DataunstableCreateEventRequest struct {

	// event
	Event *DataunstableEvent `json:"event,omitempty"`
}

// Validate validates this dataunstable create event request
func (m *DataunstableCreateEventRequest) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateEvent(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *DataunstableCreateEventRequest) validateEvent(formats strfmt.Registry) error {
	if swag.IsZero(m.Event) { // not required
		return nil
	}

	if m.Event != nil {
		if err := m.Event.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("event")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("event")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this dataunstable create event request based on the context it is used
func (m *DataunstableCreateEventRequest) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateEvent(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *DataunstableCreateEventRequest) contextValidateEvent(ctx context.Context, formats strfmt.Registry) error {

	if m.Event != nil {

		if swag.IsZero(m.Event) { // not required
			return nil
		}

		if err := m.Event.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("event")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("event")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *DataunstableCreateEventRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DataunstableCreateEventRequest) UnmarshalBinary(b []byte) error {
	var res DataunstableCreateEventRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
