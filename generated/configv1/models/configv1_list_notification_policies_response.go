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

// Configv1ListNotificationPoliciesResponse configv1 list notification policies response
//
// swagger:model configv1ListNotificationPoliciesResponse
type Configv1ListNotificationPoliciesResponse struct {

	// notification policies
	NotificationPolicies []*Configv1NotificationPolicy `json:"notification_policies"`

	// page
	Page *Configv1PageResult `json:"page,omitempty"`
}

// Validate validates this configv1 list notification policies response
func (m *Configv1ListNotificationPoliciesResponse) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateNotificationPolicies(formats); err != nil {
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

func (m *Configv1ListNotificationPoliciesResponse) validateNotificationPolicies(formats strfmt.Registry) error {
	if swag.IsZero(m.NotificationPolicies) { // not required
		return nil
	}

	for i := 0; i < len(m.NotificationPolicies); i++ {
		if swag.IsZero(m.NotificationPolicies[i]) { // not required
			continue
		}

		if m.NotificationPolicies[i] != nil {
			if err := m.NotificationPolicies[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("notification_policies" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("notification_policies" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *Configv1ListNotificationPoliciesResponse) validatePage(formats strfmt.Registry) error {
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

// ContextValidate validate this configv1 list notification policies response based on the context it is used
func (m *Configv1ListNotificationPoliciesResponse) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateNotificationPolicies(ctx, formats); err != nil {
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

func (m *Configv1ListNotificationPoliciesResponse) contextValidateNotificationPolicies(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.NotificationPolicies); i++ {

		if m.NotificationPolicies[i] != nil {

			if swag.IsZero(m.NotificationPolicies[i]) { // not required
				return nil
			}

			if err := m.NotificationPolicies[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("notification_policies" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("notification_policies" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *Configv1ListNotificationPoliciesResponse) contextValidatePage(ctx context.Context, formats strfmt.Registry) error {

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
func (m *Configv1ListNotificationPoliciesResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Configv1ListNotificationPoliciesResponse) UnmarshalBinary(b []byte) error {
	var res Configv1ListNotificationPoliciesResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
