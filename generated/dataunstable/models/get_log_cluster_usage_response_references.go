// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// GetLogClusterUsageResponseReferences not used initially
//
// swagger:model GetLogClusterUsageResponseReferences
type GetLogClusterUsageResponseReferences struct {

	// dashboards
	Dashboards string `json:"dashboards,omitempty"`

	// monitors
	Monitors string `json:"monitors,omitempty"`

	// saved searches
	SavedSearches string `json:"saved_searches,omitempty"`
}

// Validate validates this get log cluster usage response references
func (m *GetLogClusterUsageResponseReferences) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this get log cluster usage response references based on context it is used
func (m *GetLogClusterUsageResponseReferences) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *GetLogClusterUsageResponseReferences) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *GetLogClusterUsageResponseReferences) UnmarshalBinary(b []byte) error {
	var res GetLogClusterUsageResponseReferences
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
