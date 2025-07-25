// Code generated by go-swagger; DO NOT EDIT.

package data_unstable

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
)

// NewCancelLogQueryParams creates a new CancelLogQueryParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewCancelLogQueryParams() *CancelLogQueryParams {
	return &CancelLogQueryParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewCancelLogQueryParamsWithTimeout creates a new CancelLogQueryParams object
// with the ability to set a timeout on a request.
func NewCancelLogQueryParamsWithTimeout(timeout time.Duration) *CancelLogQueryParams {
	return &CancelLogQueryParams{
		timeout: timeout,
	}
}

// NewCancelLogQueryParamsWithContext creates a new CancelLogQueryParams object
// with the ability to set a context for a request.
func NewCancelLogQueryParamsWithContext(ctx context.Context) *CancelLogQueryParams {
	return &CancelLogQueryParams{
		Context: ctx,
	}
}

// NewCancelLogQueryParamsWithHTTPClient creates a new CancelLogQueryParams object
// with the ability to set a custom HTTPClient for a request.
func NewCancelLogQueryParamsWithHTTPClient(client *http.Client) *CancelLogQueryParams {
	return &CancelLogQueryParams{
		HTTPClient: client,
	}
}

/*
CancelLogQueryParams contains all the parameters to send to the API endpoint

	for the cancel log query operation.

	Typically these are written to a http.Request.
*/
type CancelLogQueryParams struct {

	// QueryID.
	QueryID *string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the cancel log query params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *CancelLogQueryParams) WithDefaults() *CancelLogQueryParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the cancel log query params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *CancelLogQueryParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the cancel log query params
func (o *CancelLogQueryParams) WithTimeout(timeout time.Duration) *CancelLogQueryParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the cancel log query params
func (o *CancelLogQueryParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the cancel log query params
func (o *CancelLogQueryParams) WithContext(ctx context.Context) *CancelLogQueryParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the cancel log query params
func (o *CancelLogQueryParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the cancel log query params
func (o *CancelLogQueryParams) WithHTTPClient(client *http.Client) *CancelLogQueryParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the cancel log query params
func (o *CancelLogQueryParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithQueryID adds the queryID to the cancel log query params
func (o *CancelLogQueryParams) WithQueryID(queryID *string) *CancelLogQueryParams {
	o.SetQueryID(queryID)
	return o
}

// SetQueryID adds the queryId to the cancel log query params
func (o *CancelLogQueryParams) SetQueryID(queryID *string) {
	o.QueryID = queryID
}

// WriteToRequest writes these params to a swagger request
func (o *CancelLogQueryParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.QueryID != nil {

		// query param query_id
		var qrQueryID string

		if o.QueryID != nil {
			qrQueryID = *o.QueryID
		}
		qQueryID := qrQueryID
		if qQueryID != "" {

			if err := r.SetQueryParam("query_id", qQueryID); err != nil {
				return err
			}
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
