// Code generated by go-swagger; DO NOT EDIT.

package dataset

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

// NewUpdateDatasetParams creates a new UpdateDatasetParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewUpdateDatasetParams() *UpdateDatasetParams {
	return &UpdateDatasetParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewUpdateDatasetParamsWithTimeout creates a new UpdateDatasetParams object
// with the ability to set a timeout on a request.
func NewUpdateDatasetParamsWithTimeout(timeout time.Duration) *UpdateDatasetParams {
	return &UpdateDatasetParams{
		timeout: timeout,
	}
}

// NewUpdateDatasetParamsWithContext creates a new UpdateDatasetParams object
// with the ability to set a context for a request.
func NewUpdateDatasetParamsWithContext(ctx context.Context) *UpdateDatasetParams {
	return &UpdateDatasetParams{
		Context: ctx,
	}
}

// NewUpdateDatasetParamsWithHTTPClient creates a new UpdateDatasetParams object
// with the ability to set a custom HTTPClient for a request.
func NewUpdateDatasetParamsWithHTTPClient(client *http.Client) *UpdateDatasetParams {
	return &UpdateDatasetParams{
		HTTPClient: client,
	}
}

/*
UpdateDatasetParams contains all the parameters to send to the API endpoint

	for the update dataset operation.

	Typically these are written to a http.Request.
*/
type UpdateDatasetParams struct {

	// Body.
	Body *models.ConfigV1UpdateDatasetBody

	// Slug.
	Slug string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the update dataset params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *UpdateDatasetParams) WithDefaults() *UpdateDatasetParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the update dataset params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *UpdateDatasetParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the update dataset params
func (o *UpdateDatasetParams) WithTimeout(timeout time.Duration) *UpdateDatasetParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the update dataset params
func (o *UpdateDatasetParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the update dataset params
func (o *UpdateDatasetParams) WithContext(ctx context.Context) *UpdateDatasetParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the update dataset params
func (o *UpdateDatasetParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the update dataset params
func (o *UpdateDatasetParams) WithHTTPClient(client *http.Client) *UpdateDatasetParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the update dataset params
func (o *UpdateDatasetParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the update dataset params
func (o *UpdateDatasetParams) WithBody(body *models.ConfigV1UpdateDatasetBody) *UpdateDatasetParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the update dataset params
func (o *UpdateDatasetParams) SetBody(body *models.ConfigV1UpdateDatasetBody) {
	o.Body = body
}

// WithSlug adds the slug to the update dataset params
func (o *UpdateDatasetParams) WithSlug(slug string) *UpdateDatasetParams {
	o.SetSlug(slug)
	return o
}

// SetSlug adds the slug to the update dataset params
func (o *UpdateDatasetParams) SetSlug(slug string) {
	o.Slug = slug
}

// WriteToRequest writes these params to a swagger request
func (o *UpdateDatasetParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error
	if o.Body != nil {
		if err := r.SetBodyParam(o.Body); err != nil {
			return err
		}
	}

	// path param slug
	if err := r.SetPathParam("slug", o.Slug); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
