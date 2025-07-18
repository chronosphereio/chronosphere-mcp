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

// ConfigV1UpdateDatasetBody config v1 update dataset body
//
// swagger:model ConfigV1UpdateDatasetBody
type ConfigV1UpdateDatasetBody struct {

	// If true, the Dataset will be created if it does not already exist, identified by slug. If false, an error will be returned if the Dataset does not already exist.
	CreateIfMissing bool `json:"create_if_missing,omitempty"`

	// dataset
	Dataset *Configv1Dataset `json:"dataset,omitempty"`

	// If true, the Dataset will not be created nor updated, and no response Dataset will be returned. The response will return an error if the given Dataset is invalid.
	DryRun bool `json:"dry_run,omitempty"`
}

// Validate validates this config v1 update dataset body
func (m *ConfigV1UpdateDatasetBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateDataset(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ConfigV1UpdateDatasetBody) validateDataset(formats strfmt.Registry) error {
	if swag.IsZero(m.Dataset) { // not required
		return nil
	}

	if m.Dataset != nil {
		if err := m.Dataset.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("dataset")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("dataset")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this config v1 update dataset body based on the context it is used
func (m *ConfigV1UpdateDatasetBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateDataset(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ConfigV1UpdateDatasetBody) contextValidateDataset(ctx context.Context, formats strfmt.Registry) error {

	if m.Dataset != nil {

		if swag.IsZero(m.Dataset) { // not required
			return nil
		}

		if err := m.Dataset.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("dataset")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("dataset")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ConfigV1UpdateDatasetBody) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ConfigV1UpdateDatasetBody) UnmarshalBinary(b []byte) error {
	var res ConfigV1UpdateDatasetBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
