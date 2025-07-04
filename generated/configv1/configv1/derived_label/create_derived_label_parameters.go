// Code generated by go-swagger; DO NOT EDIT.

package derived_label

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"

	"github.com/chronosphereio/chronosphere-mcp/generated/configv1/models"
)

// NewCreateDerivedLabelParams creates a new CreateDerivedLabelParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewCreateDerivedLabelParams() *CreateDerivedLabelParams {
	return &CreateDerivedLabelParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewCreateDerivedLabelParamsWithTimeout creates a new CreateDerivedLabelParams object
// with the ability to set a timeout on a request.
func NewCreateDerivedLabelParamsWithTimeout(timeout time.Duration) *CreateDerivedLabelParams {
	return &CreateDerivedLabelParams{
		timeout: timeout,
	}
}

// NewCreateDerivedLabelParamsWithContext creates a new CreateDerivedLabelParams object
// with the ability to set a context for a request.
func NewCreateDerivedLabelParamsWithContext(ctx context.Context) *CreateDerivedLabelParams {
	return &CreateDerivedLabelParams{
		Context: ctx,
	}
}

// NewCreateDerivedLabelParamsWithHTTPClient creates a new CreateDerivedLabelParams object
// with the ability to set a custom HTTPClient for a request.
func NewCreateDerivedLabelParamsWithHTTPClient(client *http.Client) *CreateDerivedLabelParams {
	return &CreateDerivedLabelParams{
		HTTPClient: client,
	}
}

/*
CreateDerivedLabelParams contains all the parameters to send to the API endpoint

	for the create derived label operation.

	Typically these are written to a http.Request.
*/
type CreateDerivedLabelParams struct {

	// Body.
	Body *models.Configv1CreateDerivedLabelRequest

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the create derived label params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *CreateDerivedLabelParams) WithDefaults() *CreateDerivedLabelParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the create derived label params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *CreateDerivedLabelParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the create derived label params
func (o *CreateDerivedLabelParams) WithTimeout(timeout time.Duration) *CreateDerivedLabelParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the create derived label params
func (o *CreateDerivedLabelParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the create derived label params
func (o *CreateDerivedLabelParams) WithContext(ctx context.Context) *CreateDerivedLabelParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the create derived label params
func (o *CreateDerivedLabelParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the create derived label params
func (o *CreateDerivedLabelParams) WithHTTPClient(client *http.Client) *CreateDerivedLabelParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the create derived label params
func (o *CreateDerivedLabelParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the create derived label params
func (o *CreateDerivedLabelParams) WithBody(body *models.Configv1CreateDerivedLabelRequest) *CreateDerivedLabelParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the create derived label params
func (o *CreateDerivedLabelParams) SetBody(body *models.Configv1CreateDerivedLabelRequest) {
	o.Body = body
}

// WriteToRequest writes these params to a swagger request
func (o *CreateDerivedLabelParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error
	if o.Body != nil {
		if err := r.SetBodyParam(o.Body); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
