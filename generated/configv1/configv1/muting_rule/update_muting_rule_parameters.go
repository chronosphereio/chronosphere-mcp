// Code generated by go-swagger; DO NOT EDIT.

package muting_rule

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

// NewUpdateMutingRuleParams creates a new UpdateMutingRuleParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewUpdateMutingRuleParams() *UpdateMutingRuleParams {
	return &UpdateMutingRuleParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewUpdateMutingRuleParamsWithTimeout creates a new UpdateMutingRuleParams object
// with the ability to set a timeout on a request.
func NewUpdateMutingRuleParamsWithTimeout(timeout time.Duration) *UpdateMutingRuleParams {
	return &UpdateMutingRuleParams{
		timeout: timeout,
	}
}

// NewUpdateMutingRuleParamsWithContext creates a new UpdateMutingRuleParams object
// with the ability to set a context for a request.
func NewUpdateMutingRuleParamsWithContext(ctx context.Context) *UpdateMutingRuleParams {
	return &UpdateMutingRuleParams{
		Context: ctx,
	}
}

// NewUpdateMutingRuleParamsWithHTTPClient creates a new UpdateMutingRuleParams object
// with the ability to set a custom HTTPClient for a request.
func NewUpdateMutingRuleParamsWithHTTPClient(client *http.Client) *UpdateMutingRuleParams {
	return &UpdateMutingRuleParams{
		HTTPClient: client,
	}
}

/*
UpdateMutingRuleParams contains all the parameters to send to the API endpoint

	for the update muting rule operation.

	Typically these are written to a http.Request.
*/
type UpdateMutingRuleParams struct {

	// Body.
	Body *models.ConfigV1UpdateMutingRuleBody

	// Slug.
	Slug string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the update muting rule params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *UpdateMutingRuleParams) WithDefaults() *UpdateMutingRuleParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the update muting rule params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *UpdateMutingRuleParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the update muting rule params
func (o *UpdateMutingRuleParams) WithTimeout(timeout time.Duration) *UpdateMutingRuleParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the update muting rule params
func (o *UpdateMutingRuleParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the update muting rule params
func (o *UpdateMutingRuleParams) WithContext(ctx context.Context) *UpdateMutingRuleParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the update muting rule params
func (o *UpdateMutingRuleParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the update muting rule params
func (o *UpdateMutingRuleParams) WithHTTPClient(client *http.Client) *UpdateMutingRuleParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the update muting rule params
func (o *UpdateMutingRuleParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the update muting rule params
func (o *UpdateMutingRuleParams) WithBody(body *models.ConfigV1UpdateMutingRuleBody) *UpdateMutingRuleParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the update muting rule params
func (o *UpdateMutingRuleParams) SetBody(body *models.ConfigV1UpdateMutingRuleBody) {
	o.Body = body
}

// WithSlug adds the slug to the update muting rule params
func (o *UpdateMutingRuleParams) WithSlug(slug string) *UpdateMutingRuleParams {
	o.SetSlug(slug)
	return o
}

// SetSlug adds the slug to the update muting rule params
func (o *UpdateMutingRuleParams) SetSlug(slug string) {
	o.Slug = slug
}

// WriteToRequest writes these params to a swagger request
func (o *UpdateMutingRuleParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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
