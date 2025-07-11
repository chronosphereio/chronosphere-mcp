// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// DataunstableLogCluster dataunstable log cluster
//
// swagger:model dataunstableLogCluster
type DataunstableLogCluster struct {

	// Filter expression that can be used to select logs matching the cluster.
	Filter string `json:"filter,omitempty"`

	// Number of logs present in this cluster.
	NumLogs int64 `json:"num_logs,omitempty"`

	// The pattern that represents this cluster.
	Pattern string `json:"pattern,omitempty"`
}

// Validate validates this dataunstable log cluster
func (m *DataunstableLogCluster) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this dataunstable log cluster based on context it is used
func (m *DataunstableLogCluster) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *DataunstableLogCluster) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DataunstableLogCluster) UnmarshalBinary(b []byte) error {
	var res DataunstableLogCluster
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
