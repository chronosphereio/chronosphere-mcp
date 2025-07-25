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

// Configv1ListMutingRulesResponse configv1 list muting rules response
//
// swagger:model configv1ListMutingRulesResponse
type Configv1ListMutingRulesResponse struct {

	// muting rules
	MutingRules []*Configv1MutingRule `json:"muting_rules"`

	// page
	Page *Configv1PageResult `json:"page,omitempty"`
}

// Validate validates this configv1 list muting rules response
func (m *Configv1ListMutingRulesResponse) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateMutingRules(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePage(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Configv1ListMutingRulesResponse) validateMutingRules(formats strfmt.Registry) error {
	if swag.IsZero(m.MutingRules) { // not required
		return nil
	}

	for i := 0; i < len(m.MutingRules); i++ {
		if swag.IsZero(m.MutingRules[i]) { // not required
			continue
		}

		if m.MutingRules[i] != nil {
			if err := m.MutingRules[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("muting_rules" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("muting_rules" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *Configv1ListMutingRulesResponse) validatePage(formats strfmt.Registry) error {
	if swag.IsZero(m.Page) { // not required
		return nil
	}

	if m.Page != nil {
		if err := m.Page.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("page")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("page")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this configv1 list muting rules response based on the context it is used
func (m *Configv1ListMutingRulesResponse) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateMutingRules(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidatePage(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Configv1ListMutingRulesResponse) contextValidateMutingRules(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.MutingRules); i++ {

		if m.MutingRules[i] != nil {

			if swag.IsZero(m.MutingRules[i]) { // not required
				return nil
			}

			if err := m.MutingRules[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("muting_rules" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("muting_rules" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *Configv1ListMutingRulesResponse) contextValidatePage(ctx context.Context, formats strfmt.Registry) error {

	if m.Page != nil {

		if swag.IsZero(m.Page) { // not required
			return nil
		}

		if err := m.Page.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("page")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("page")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Configv1ListMutingRulesResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Configv1ListMutingRulesResponse) UnmarshalBinary(b []byte) error {
	var res Configv1ListMutingRulesResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
