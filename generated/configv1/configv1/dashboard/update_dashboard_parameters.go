// Code generated by go-swagger; DO NOT EDIT.

package dashboard

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

// NewUpdateDashboardParams creates a new UpdateDashboardParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewUpdateDashboardParams() *UpdateDashboardParams {
	return &UpdateDashboardParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewUpdateDashboardParamsWithTimeout creates a new UpdateDashboardParams object
// with the ability to set a timeout on a request.
func NewUpdateDashboardParamsWithTimeout(timeout time.Duration) *UpdateDashboardParams {
	return &UpdateDashboardParams{
		timeout: timeout,
	}
}

// NewUpdateDashboardParamsWithContext creates a new UpdateDashboardParams object
// with the ability to set a context for a request.
func NewUpdateDashboardParamsWithContext(ctx context.Context) *UpdateDashboardParams {
	return &UpdateDashboardParams{
		Context: ctx,
	}
}

// NewUpdateDashboardParamsWithHTTPClient creates a new UpdateDashboardParams object
// with the ability to set a custom HTTPClient for a request.
func NewUpdateDashboardParamsWithHTTPClient(client *http.Client) *UpdateDashboardParams {
	return &UpdateDashboardParams{
		HTTPClient: client,
	}
}

/*
UpdateDashboardParams contains all the parameters to send to the API endpoint

	for the update dashboard operation.

	Typically these are written to a http.Request.
*/
type UpdateDashboardParams struct {

	// Body.
	Body *models.ConfigV1UpdateDashboardBody

	// Slug.
	Slug string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the update dashboard params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *UpdateDashboardParams) WithDefaults() *UpdateDashboardParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the update dashboard params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *UpdateDashboardParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the update dashboard params
func (o *UpdateDashboardParams) WithTimeout(timeout time.Duration) *UpdateDashboardParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the update dashboard params
func (o *UpdateDashboardParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the update dashboard params
func (o *UpdateDashboardParams) WithContext(ctx context.Context) *UpdateDashboardParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the update dashboard params
func (o *UpdateDashboardParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the update dashboard params
func (o *UpdateDashboardParams) WithHTTPClient(client *http.Client) *UpdateDashboardParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the update dashboard params
func (o *UpdateDashboardParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the update dashboard params
func (o *UpdateDashboardParams) WithBody(body *models.ConfigV1UpdateDashboardBody) *UpdateDashboardParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the update dashboard params
func (o *UpdateDashboardParams) SetBody(body *models.ConfigV1UpdateDashboardBody) {
	o.Body = body
}

// WithSlug adds the slug to the update dashboard params
func (o *UpdateDashboardParams) WithSlug(slug string) *UpdateDashboardParams {
	o.SetSlug(slug)
	return o
}

// SetSlug adds the slug to the update dashboard params
func (o *UpdateDashboardParams) SetSlug(slug string) {
	o.Slug = slug
}

// WriteToRequest writes these params to a swagger request
func (o *UpdateDashboardParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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
