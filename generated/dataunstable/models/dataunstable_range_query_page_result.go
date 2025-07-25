// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// DataunstableRangeQueryPageResult dataunstable range query page result
//
// swagger:model dataunstableRangeQueryPageResult
type DataunstableRangeQueryPageResult struct {

	// Opaque page token which identifies the next page of items which the
	// client should request. An empty next_token indicates that there are no
	// more items to return.
	NextToken string `json:"next_token,omitempty"`
}

// Validate validates this dataunstable range query page result
func (m *DataunstableRangeQueryPageResult) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this dataunstable range query page result based on context it is used
func (m *DataunstableRangeQueryPageResult) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *DataunstableRangeQueryPageResult) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DataunstableRangeQueryPageResult) UnmarshalBinary(b []byte) error {
	var res DataunstableRangeQueryPageResult
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
