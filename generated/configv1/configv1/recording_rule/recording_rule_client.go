// Code generated by go-swagger; DO NOT EDIT.

package recording_rule

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

// New creates a new recording rule API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) ClientService {
	return &Client{transport: transport, formats: formats}
}

// New creates a new recording rule API client with basic auth credentials.
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

// New creates a new recording rule API client with a bearer token for authentication.
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
Client for recording rule API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

// ClientOption may be used to customize the behavior of Client methods.
type ClientOption func(*runtime.ClientOperation)

// ClientService is the interface for Client methods
type ClientService interface {
	CreateRecordingRule(params *CreateRecordingRuleParams, opts ...ClientOption) (*CreateRecordingRuleOK, error)

	DeleteRecordingRule(params *DeleteRecordingRuleParams, opts ...ClientOption) (*DeleteRecordingRuleOK, error)

	ListRecordingRules(params *ListRecordingRulesParams, opts ...ClientOption) (*ListRecordingRulesOK, error)

	ReadRecordingRule(params *ReadRecordingRuleParams, opts ...ClientOption) (*ReadRecordingRuleOK, error)

	UpdateRecordingRule(params *UpdateRecordingRuleParams, opts ...ClientOption) (*UpdateRecordingRuleOK, error)

	SetTransport(transport runtime.ClientTransport)
}

/*
CreateRecordingRule create recording rule API
*/
func (a *Client) CreateRecordingRule(params *CreateRecordingRuleParams, opts ...ClientOption) (*CreateRecordingRuleOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewCreateRecordingRuleParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "CreateRecordingRule",
		Method:             "POST",
		PathPattern:        "/api/v1/config/recording-rules",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &CreateRecordingRuleReader{formats: a.formats},
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
	success, ok := result.(*CreateRecordingRuleOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*CreateRecordingRuleDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
DeleteRecordingRule delete recording rule API
*/
func (a *Client) DeleteRecordingRule(params *DeleteRecordingRuleParams, opts ...ClientOption) (*DeleteRecordingRuleOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewDeleteRecordingRuleParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "DeleteRecordingRule",
		Method:             "DELETE",
		PathPattern:        "/api/v1/config/recording-rules/{slug}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &DeleteRecordingRuleReader{formats: a.formats},
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
	success, ok := result.(*DeleteRecordingRuleOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*DeleteRecordingRuleDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
ListRecordingRules list recording rules API
*/
func (a *Client) ListRecordingRules(params *ListRecordingRulesParams, opts ...ClientOption) (*ListRecordingRulesOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewListRecordingRulesParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "ListRecordingRules",
		Method:             "GET",
		PathPattern:        "/api/v1/config/recording-rules",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &ListRecordingRulesReader{formats: a.formats},
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
	success, ok := result.(*ListRecordingRulesOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*ListRecordingRulesDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
ReadRecordingRule read recording rule API
*/
func (a *Client) ReadRecordingRule(params *ReadRecordingRuleParams, opts ...ClientOption) (*ReadRecordingRuleOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewReadRecordingRuleParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "ReadRecordingRule",
		Method:             "GET",
		PathPattern:        "/api/v1/config/recording-rules/{slug}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &ReadRecordingRuleReader{formats: a.formats},
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
	success, ok := result.(*ReadRecordingRuleOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*ReadRecordingRuleDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
UpdateRecordingRule update recording rule API
*/
func (a *Client) UpdateRecordingRule(params *UpdateRecordingRuleParams, opts ...ClientOption) (*UpdateRecordingRuleOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewUpdateRecordingRuleParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "UpdateRecordingRule",
		Method:             "PUT",
		PathPattern:        "/api/v1/config/recording-rules/{slug}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &UpdateRecordingRuleReader{formats: a.formats},
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
	success, ok := result.(*UpdateRecordingRuleOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*UpdateRecordingRuleDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}
