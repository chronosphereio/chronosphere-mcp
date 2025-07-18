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

// ServiceAccountMetricsRestriction service account metrics restriction
//
// swagger:model ServiceAccountMetricsRestriction
type ServiceAccountMetricsRestriction struct {

	// Optional labels which further restricts the service account to only read
	// or write metrics with the given label names and values.
	Labels map[string]string `json:"labels,omitempty"`

	// permission
	Permission MetricsRestrictionPermission `json:"permission,omitempty"`
}

// Validate validates this service account metrics restriction
func (m *ServiceAccountMetricsRestriction) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validatePermission(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ServiceAccountMetricsRestriction) validatePermission(formats strfmt.Registry) error {
	if swag.IsZero(m.Permission) { // not required
		return nil
	}

	if err := m.Permission.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("permission")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("permission")
		}
		return err
	}

	return nil
}

// ContextValidate validate this service account metrics restriction based on the context it is used
func (m *ServiceAccountMetricsRestriction) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidatePermission(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ServiceAccountMetricsRestriction) contextValidatePermission(ctx context.Context, formats strfmt.Registry) error {

	if swag.IsZero(m.Permission) { // not required
		return nil
	}

	if err := m.Permission.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("permission")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("permission")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ServiceAccountMetricsRestriction) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ServiceAccountMetricsRestriction) UnmarshalBinary(b []byte) error {
	var res ServiceAccountMetricsRestriction
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
