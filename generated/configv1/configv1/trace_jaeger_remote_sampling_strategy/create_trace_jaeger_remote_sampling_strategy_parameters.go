// Code generated by go-swagger; DO NOT EDIT.

package trace_jaeger_remote_sampling_strategy

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

// NewCreateTraceJaegerRemoteSamplingStrategyParams creates a new CreateTraceJaegerRemoteSamplingStrategyParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewCreateTraceJaegerRemoteSamplingStrategyParams() *CreateTraceJaegerRemoteSamplingStrategyParams {
	return &CreateTraceJaegerRemoteSamplingStrategyParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewCreateTraceJaegerRemoteSamplingStrategyParamsWithTimeout creates a new CreateTraceJaegerRemoteSamplingStrategyParams object
// with the ability to set a timeout on a request.
func NewCreateTraceJaegerRemoteSamplingStrategyParamsWithTimeout(timeout time.Duration) *CreateTraceJaegerRemoteSamplingStrategyParams {
	return &CreateTraceJaegerRemoteSamplingStrategyParams{
		timeout: timeout,
	}
}

// NewCreateTraceJaegerRemoteSamplingStrategyParamsWithContext creates a new CreateTraceJaegerRemoteSamplingStrategyParams object
// with the ability to set a context for a request.
func NewCreateTraceJaegerRemoteSamplingStrategyParamsWithContext(ctx context.Context) *CreateTraceJaegerRemoteSamplingStrategyParams {
	return &CreateTraceJaegerRemoteSamplingStrategyParams{
		Context: ctx,
	}
}

// NewCreateTraceJaegerRemoteSamplingStrategyParamsWithHTTPClient creates a new CreateTraceJaegerRemoteSamplingStrategyParams object
// with the ability to set a custom HTTPClient for a request.
func NewCreateTraceJaegerRemoteSamplingStrategyParamsWithHTTPClient(client *http.Client) *CreateTraceJaegerRemoteSamplingStrategyParams {
	return &CreateTraceJaegerRemoteSamplingStrategyParams{
		HTTPClient: client,
	}
}

/*
CreateTraceJaegerRemoteSamplingStrategyParams contains all the parameters to send to the API endpoint

	for the create trace jaeger remote sampling strategy operation.

	Typically these are written to a http.Request.
*/
type CreateTraceJaegerRemoteSamplingStrategyParams struct {

	// Body.
	Body *models.Configv1CreateTraceJaegerRemoteSamplingStrategyRequest

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the create trace jaeger remote sampling strategy params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *CreateTraceJaegerRemoteSamplingStrategyParams) WithDefaults() *CreateTraceJaegerRemoteSamplingStrategyParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the create trace jaeger remote sampling strategy params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *CreateTraceJaegerRemoteSamplingStrategyParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the create trace jaeger remote sampling strategy params
func (o *CreateTraceJaegerRemoteSamplingStrategyParams) WithTimeout(timeout time.Duration) *CreateTraceJaegerRemoteSamplingStrategyParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the create trace jaeger remote sampling strategy params
func (o *CreateTraceJaegerRemoteSamplingStrategyParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the create trace jaeger remote sampling strategy params
func (o *CreateTraceJaegerRemoteSamplingStrategyParams) WithContext(ctx context.Context) *CreateTraceJaegerRemoteSamplingStrategyParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the create trace jaeger remote sampling strategy params
func (o *CreateTraceJaegerRemoteSamplingStrategyParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the create trace jaeger remote sampling strategy params
func (o *CreateTraceJaegerRemoteSamplingStrategyParams) WithHTTPClient(client *http.Client) *CreateTraceJaegerRemoteSamplingStrategyParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the create trace jaeger remote sampling strategy params
func (o *CreateTraceJaegerRemoteSamplingStrategyParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the create trace jaeger remote sampling strategy params
func (o *CreateTraceJaegerRemoteSamplingStrategyParams) WithBody(body *models.Configv1CreateTraceJaegerRemoteSamplingStrategyRequest) *CreateTraceJaegerRemoteSamplingStrategyParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the create trace jaeger remote sampling strategy params
func (o *CreateTraceJaegerRemoteSamplingStrategyParams) SetBody(body *models.Configv1CreateTraceJaegerRemoteSamplingStrategyRequest) {
	o.Body = body
}

// WriteToRequest writes these params to a swagger request
func (o *CreateTraceJaegerRemoteSamplingStrategyParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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
