// Code generated by go-swagger; DO NOT EDIT.

package derived_metric

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

// NewUpdateDerivedMetricParams creates a new UpdateDerivedMetricParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewUpdateDerivedMetricParams() *UpdateDerivedMetricParams {
	return &UpdateDerivedMetricParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewUpdateDerivedMetricParamsWithTimeout creates a new UpdateDerivedMetricParams object
// with the ability to set a timeout on a request.
func NewUpdateDerivedMetricParamsWithTimeout(timeout time.Duration) *UpdateDerivedMetricParams {
	return &UpdateDerivedMetricParams{
		timeout: timeout,
	}
}

// NewUpdateDerivedMetricParamsWithContext creates a new UpdateDerivedMetricParams object
// with the ability to set a context for a request.
func NewUpdateDerivedMetricParamsWithContext(ctx context.Context) *UpdateDerivedMetricParams {
	return &UpdateDerivedMetricParams{
		Context: ctx,
	}
}

// NewUpdateDerivedMetricParamsWithHTTPClient creates a new UpdateDerivedMetricParams object
// with the ability to set a custom HTTPClient for a request.
func NewUpdateDerivedMetricParamsWithHTTPClient(client *http.Client) *UpdateDerivedMetricParams {
	return &UpdateDerivedMetricParams{
		HTTPClient: client,
	}
}

/*
UpdateDerivedMetricParams contains all the parameters to send to the API endpoint

	for the update derived metric operation.

	Typically these are written to a http.Request.
*/
type UpdateDerivedMetricParams struct {

	// Body.
	Body *models.ConfigV1UpdateDerivedMetricBody

	// Slug.
	Slug string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the update derived metric params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *UpdateDerivedMetricParams) WithDefaults() *UpdateDerivedMetricParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the update derived metric params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *UpdateDerivedMetricParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the update derived metric params
func (o *UpdateDerivedMetricParams) WithTimeout(timeout time.Duration) *UpdateDerivedMetricParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the update derived metric params
func (o *UpdateDerivedMetricParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the update derived metric params
func (o *UpdateDerivedMetricParams) WithContext(ctx context.Context) *UpdateDerivedMetricParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the update derived metric params
func (o *UpdateDerivedMetricParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the update derived metric params
func (o *UpdateDerivedMetricParams) WithHTTPClient(client *http.Client) *UpdateDerivedMetricParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the update derived metric params
func (o *UpdateDerivedMetricParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the update derived metric params
func (o *UpdateDerivedMetricParams) WithBody(body *models.ConfigV1UpdateDerivedMetricBody) *UpdateDerivedMetricParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the update derived metric params
func (o *UpdateDerivedMetricParams) SetBody(body *models.ConfigV1UpdateDerivedMetricBody) {
	o.Body = body
}

// WithSlug adds the slug to the update derived metric params
func (o *UpdateDerivedMetricParams) WithSlug(slug string) *UpdateDerivedMetricParams {
	o.SetSlug(slug)
	return o
}

// SetSlug adds the slug to the update derived metric params
func (o *UpdateDerivedMetricParams) SetSlug(slug string) {
	o.Slug = slug
}

// WriteToRequest writes these params to a swagger request
func (o *UpdateDerivedMetricParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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
