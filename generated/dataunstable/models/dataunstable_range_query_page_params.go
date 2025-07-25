// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// DataunstableRangeQueryPageParams dataunstable range query page params
//
// swagger:model dataunstableRangeQueryPageParams
type DataunstableRangeQueryPageParams struct {

	// Opaque page token identifying which page to request. An empty token
	// identifies the first page.
	Token string `json:"token,omitempty"`
}

// Validate validates this dataunstable range query page params
func (m *DataunstableRangeQueryPageParams) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this dataunstable range query page params based on context it is used
func (m *DataunstableRangeQueryPageParams) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *DataunstableRangeQueryPageParams) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DataunstableRangeQueryPageParams) UnmarshalBinary(b []byte) error {
	var res DataunstableRangeQueryPageParams
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
