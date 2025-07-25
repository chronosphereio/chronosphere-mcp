// Code generated by go-swagger; DO NOT EDIT.

package notification_policy

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

// New creates a new notification policy API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) ClientService {
	return &Client{transport: transport, formats: formats}
}

// New creates a new notification policy API client with basic auth credentials.
// It takes the following parameters:
// - host: http host (github.com).
// - basePath: any base path for the API client ("/v1", "/v3").
// - scheme: http scheme ("http", "https").
// - user: user for basic authentication header.
// - password: password for basic authentication header.
func NewClientWithBasicAuth(host, basePath, scheme, user, password string) ClientService {
	transport := httptransport.New(host, basePath, []string{scheme})
	transport.DefaultAuthentication = httptransport.BasicAuth(user, password)
	return &Client{transport: transport, formats: strfmt.Default}
}

// New creates a new notification policy API client with a bearer token for authentication.
// It takes the following parameters:
// - host: http host (github.com).
// - basePath: any base path for the API client ("/v1", "/v3").
// - scheme: http scheme ("http", "https").
// - bearerToken: bearer token for Bearer authentication header.
func NewClientWithBearerToken(host, basePath, scheme, bearerToken string) ClientService {
	transport := httptransport.New(host, basePath, []string{scheme})
	transport.DefaultAuthentication = httptransport.BearerToken(bearerToken)
	return &Client{transport: transport, formats: strfmt.Default}
}

/*
Client for notification policy API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

// ClientOption may be used to customize the behavior of Client methods.
type ClientOption func(*runtime.ClientOperation)

// ClientService is the interface for Client methods
type ClientService interface {
	CreateNotificationPolicy(params *CreateNotificationPolicyParams, opts ...ClientOption) (*CreateNotificationPolicyOK, error)

	DeleteNotificationPolicy(params *DeleteNotificationPolicyParams, opts ...ClientOption) (*DeleteNotificationPolicyOK, error)

	ListNotificationPolicies(params *ListNotificationPoliciesParams, opts ...ClientOption) (*ListNotificationPoliciesOK, error)

	ReadNotificationPolicy(params *ReadNotificationPolicyParams, opts ...ClientOption) (*ReadNotificationPolicyOK, error)

	UpdateNotificationPolicy(params *UpdateNotificationPolicyParams, opts ...ClientOption) (*UpdateNotificationPolicyOK, error)

	SetTransport(transport runtime.ClientTransport)
}

/*
CreateNotificationPolicy create notification policy API
*/
func (a *Client) CreateNotificationPolicy(params *CreateNotificationPolicyParams, opts ...ClientOption) (*CreateNotificationPolicyOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewCreateNotificationPolicyParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "CreateNotificationPolicy",
		Method:             "POST",
		PathPattern:        "/api/v1/config/notification-policies",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &CreateNotificationPolicyReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*CreateNotificationPolicyOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*CreateNotificationPolicyDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
DeleteNotificationPolicy delete notification policy API
*/
func (a *Client) DeleteNotificationPolicy(params *DeleteNotificationPolicyParams, opts ...ClientOption) (*DeleteNotificationPolicyOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewDeleteNotificationPolicyParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "DeleteNotificationPolicy",
		Method:             "DELETE",
		PathPattern:        "/api/v1/config/notification-policies/{slug}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &DeleteNotificationPolicyReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*DeleteNotificationPolicyOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*DeleteNotificationPolicyDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
ListNotificationPolicies list notification policies API
*/
func (a *Client) ListNotificationPolicies(params *ListNotificationPoliciesParams, opts ...ClientOption) (*ListNotificationPoliciesOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewListNotificationPoliciesParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "ListNotificationPolicies",
		Method:             "GET",
		PathPattern:        "/api/v1/config/notification-policies",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &ListNotificationPoliciesReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*ListNotificationPoliciesOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*ListNotificationPoliciesDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
ReadNotificationPolicy read notification policy API
*/
func (a *Client) ReadNotificationPolicy(params *ReadNotificationPolicyParams, opts ...ClientOption) (*ReadNotificationPolicyOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewReadNotificationPolicyParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "ReadNotificationPolicy",
		Method:             "GET",
		PathPattern:        "/api/v1/config/notification-policies/{slug}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &ReadNotificationPolicyReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*ReadNotificationPolicyOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*ReadNotificationPolicyDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
UpdateNotificationPolicy update notification policy API
*/
func (a *Client) UpdateNotificationPolicy(params *UpdateNotificationPolicyParams, opts ...ClientOption) (*UpdateNotificationPolicyOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewUpdateNotificationPolicyParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "UpdateNotificationPolicy",
		Method:             "PUT",
		PathPattern:        "/api/v1/config/notification-policies/{slug}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &UpdateNotificationPolicyReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*UpdateNotificationPolicyOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*UpdateNotificationPolicyDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}
