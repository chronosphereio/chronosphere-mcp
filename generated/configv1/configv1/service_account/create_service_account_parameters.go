// Code generated by go-swagger; DO NOT EDIT.

package service_account

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

// NewCreateServiceAccountParams creates a new CreateServiceAccountParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewCreateServiceAccountParams() *CreateServiceAccountParams {
	return &CreateServiceAccountParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewCreateServiceAccountParamsWithTimeout creates a new CreateServiceAccountParams object
// with the ability to set a timeout on a request.
func NewCreateServiceAccountParamsWithTimeout(timeout time.Duration) *CreateServiceAccountParams {
	return &CreateServiceAccountParams{
		timeout: timeout,
	}
}

// NewCreateServiceAccountParamsWithContext creates a new CreateServiceAccountParams object
// with the ability to set a context for a request.
func NewCreateServiceAccountParamsWithContext(ctx context.Context) *CreateServiceAccountParams {
	return &CreateServiceAccountParams{
		Context: ctx,
	}
}

// NewCreateServiceAccountParamsWithHTTPClient creates a new CreateServiceAccountParams object
// with the ability to set a custom HTTPClient for a request.
func NewCreateServiceAccountParamsWithHTTPClient(client *http.Client) *CreateServiceAccountParams {
	return &CreateServiceAccountParams{
		HTTPClient: client,
	}
}

/*
CreateServiceAccountParams contains all the parameters to send to the API endpoint

	for the create service account operation.

	Typically these are written to a http.Request.
*/
type CreateServiceAccountParams struct {

	// Body.
	Body *models.Configv1CreateServiceAccountRequest

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the create service account params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *CreateServiceAccountParams) WithDefaults() *CreateServiceAccountParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the create service account params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *CreateServiceAccountParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the create service account params
func (o *CreateServiceAccountParams) WithTimeout(timeout time.Duration) *CreateServiceAccountParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the create service account params
func (o *CreateServiceAccountParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the create service account params
func (o *CreateServiceAccountParams) WithContext(ctx context.Context) *CreateServiceAccountParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the create service account params
func (o *CreateServiceAccountParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the create service account params
func (o *CreateServiceAccountParams) WithHTTPClient(client *http.Client) *CreateServiceAccountParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the create service account params
func (o *CreateServiceAccountParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the create service account params
func (o *CreateServiceAccountParams) WithBody(body *models.Configv1CreateServiceAccountRequest) *CreateServiceAccountParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the create service account params
func (o *CreateServiceAccountParams) SetBody(body *models.Configv1CreateServiceAccountRequest) {
	o.Body = body
}

// WriteToRequest writes these params to a swagger request
func (o *CreateServiceAccountParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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
