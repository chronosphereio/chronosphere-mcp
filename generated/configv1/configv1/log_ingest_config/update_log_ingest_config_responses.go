// Code generated by go-swagger; DO NOT EDIT.

package log_ingest_config

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/chronosphereio/chronosphere-mcp/generated/configv1/models"
)

// UpdateLogIngestConfigReader is a Reader for the UpdateLogIngestConfig structure.
type UpdateLogIngestConfigReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *UpdateLogIngestConfigReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewUpdateLogIngestConfigOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewUpdateLogIngestConfigBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewUpdateLogIngestConfigNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewUpdateLogIngestConfigInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewUpdateLogIngestConfigDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewUpdateLogIngestConfigOK creates a UpdateLogIngestConfigOK with default headers values
func NewUpdateLogIngestConfigOK() *UpdateLogIngestConfigOK {
	return &UpdateLogIngestConfigOK{}
}

/*
UpdateLogIngestConfigOK describes a response with status code 200, with default header values.

A successful response containing the updated LogIngestConfig.
*/
type UpdateLogIngestConfigOK struct {
	Payload *models.Configv1UpdateLogIngestConfigResponse
}

// IsSuccess returns true when this update log ingest config o k response has a 2xx status code
func (o *UpdateLogIngestConfigOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this update log ingest config o k response has a 3xx status code
func (o *UpdateLogIngestConfigOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this update log ingest config o k response has a 4xx status code
func (o *UpdateLogIngestConfigOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this update log ingest config o k response has a 5xx status code
func (o *UpdateLogIngestConfigOK) IsServerError() bool {
	return false
}

// IsCode returns true when this update log ingest config o k response a status code equal to that given
func (o *UpdateLogIngestConfigOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the update log ingest config o k response
func (o *UpdateLogIngestConfigOK) Code() int {
	return 200
}

func (o *UpdateLogIngestConfigOK) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[PUT /api/v1/config/log-ingest-config][%d] updateLogIngestConfigOK %s", 200, payload)
}

func (o *UpdateLogIngestConfigOK) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[PUT /api/v1/config/log-ingest-config][%d] updateLogIngestConfigOK %s", 200, payload)
}

func (o *UpdateLogIngestConfigOK) GetPayload() *models.Configv1UpdateLogIngestConfigResponse {
	return o.Payload
}

func (o *UpdateLogIngestConfigOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Configv1UpdateLogIngestConfigResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateLogIngestConfigBadRequest creates a UpdateLogIngestConfigBadRequest with default headers values
func NewUpdateLogIngestConfigBadRequest() *UpdateLogIngestConfigBadRequest {
	return &UpdateLogIngestConfigBadRequest{}
}

/*
UpdateLogIngestConfigBadRequest describes a response with status code 400, with default header values.

Cannot update the LogIngestConfig because the request is invalid.
*/
type UpdateLogIngestConfigBadRequest struct {
	Payload *models.APIError
}

// IsSuccess returns true when this update log ingest config bad request response has a 2xx status code
func (o *UpdateLogIngestConfigBadRequest) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this update log ingest config bad request response has a 3xx status code
func (o *UpdateLogIngestConfigBadRequest) IsRedirect() bool {
	return false
}

// IsClientError returns true when this update log ingest config bad request response has a 4xx status code
func (o *UpdateLogIngestConfigBadRequest) IsClientError() bool {
	return true
}

// IsServerError returns true when this update log ingest config bad request response has a 5xx status code
func (o *UpdateLogIngestConfigBadRequest) IsServerError() bool {
	return false
}

// IsCode returns true when this update log ingest config bad request response a status code equal to that given
func (o *UpdateLogIngestConfigBadRequest) IsCode(code int) bool {
	return code == 400
}

// Code gets the status code for the update log ingest config bad request response
func (o *UpdateLogIngestConfigBadRequest) Code() int {
	return 400
}

func (o *UpdateLogIngestConfigBadRequest) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[PUT /api/v1/config/log-ingest-config][%d] updateLogIngestConfigBadRequest %s", 400, payload)
}

func (o *UpdateLogIngestConfigBadRequest) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[PUT /api/v1/config/log-ingest-config][%d] updateLogIngestConfigBadRequest %s", 400, payload)
}

func (o *UpdateLogIngestConfigBadRequest) GetPayload() *models.APIError {
	return o.Payload
}

func (o *UpdateLogIngestConfigBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.APIError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateLogIngestConfigNotFound creates a UpdateLogIngestConfigNotFound with default headers values
func NewUpdateLogIngestConfigNotFound() *UpdateLogIngestConfigNotFound {
	return &UpdateLogIngestConfigNotFound{}
}

/*
UpdateLogIngestConfigNotFound describes a response with status code 404, with default header values.

Cannot update the LogIngestConfig because LogIngestConfig has not been created.
*/
type UpdateLogIngestConfigNotFound struct {
	Payload *models.APIError
}

// IsSuccess returns true when this update log ingest config not found response has a 2xx status code
func (o *UpdateLogIngestConfigNotFound) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this update log ingest config not found response has a 3xx status code
func (o *UpdateLogIngestConfigNotFound) IsRedirect() bool {
	return false
}

// IsClientError returns true when this update log ingest config not found response has a 4xx status code
func (o *UpdateLogIngestConfigNotFound) IsClientError() bool {
	return true
}

// IsServerError returns true when this update log ingest config not found response has a 5xx status code
func (o *UpdateLogIngestConfigNotFound) IsServerError() bool {
	return false
}

// IsCode returns true when this update log ingest config not found response a status code equal to that given
func (o *UpdateLogIngestConfigNotFound) IsCode(code int) bool {
	return code == 404
}

// Code gets the status code for the update log ingest config not found response
func (o *UpdateLogIngestConfigNotFound) Code() int {
	return 404
}

func (o *UpdateLogIngestConfigNotFound) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[PUT /api/v1/config/log-ingest-config][%d] updateLogIngestConfigNotFound %s", 404, payload)
}

func (o *UpdateLogIngestConfigNotFound) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[PUT /api/v1/config/log-ingest-config][%d] updateLogIngestConfigNotFound %s", 404, payload)
}

func (o *UpdateLogIngestConfigNotFound) GetPayload() *models.APIError {
	return o.Payload
}

func (o *UpdateLogIngestConfigNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.APIError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateLogIngestConfigInternalServerError creates a UpdateLogIngestConfigInternalServerError with default headers values
func NewUpdateLogIngestConfigInternalServerError() *UpdateLogIngestConfigInternalServerError {
	return &UpdateLogIngestConfigInternalServerError{}
}

/*
UpdateLogIngestConfigInternalServerError describes a response with status code 500, with default header values.

An unexpected error response.
*/
type UpdateLogIngestConfigInternalServerError struct {
	Payload *models.APIError
}

// IsSuccess returns true when this update log ingest config internal server error response has a 2xx status code
func (o *UpdateLogIngestConfigInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this update log ingest config internal server error response has a 3xx status code
func (o *UpdateLogIngestConfigInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this update log ingest config internal server error response has a 4xx status code
func (o *UpdateLogIngestConfigInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this update log ingest config internal server error response has a 5xx status code
func (o *UpdateLogIngestConfigInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this update log ingest config internal server error response a status code equal to that given
func (o *UpdateLogIngestConfigInternalServerError) IsCode(code int) bool {
	return code == 500
}

// Code gets the status code for the update log ingest config internal server error response
func (o *UpdateLogIngestConfigInternalServerError) Code() int {
	return 500
}

func (o *UpdateLogIngestConfigInternalServerError) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[PUT /api/v1/config/log-ingest-config][%d] updateLogIngestConfigInternalServerError %s", 500, payload)
}

func (o *UpdateLogIngestConfigInternalServerError) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[PUT /api/v1/config/log-ingest-config][%d] updateLogIngestConfigInternalServerError %s", 500, payload)
}

func (o *UpdateLogIngestConfigInternalServerError) GetPayload() *models.APIError {
	return o.Payload
}

func (o *UpdateLogIngestConfigInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.APIError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateLogIngestConfigDefault creates a UpdateLogIngestConfigDefault with default headers values
func NewUpdateLogIngestConfigDefault(code int) *UpdateLogIngestConfigDefault {
	return &UpdateLogIngestConfigDefault{
		_statusCode: code,
	}
}

/*
UpdateLogIngestConfigDefault describes a response with status code -1, with default header values.

An undefined error response.
*/
type UpdateLogIngestConfigDefault struct {
	_statusCode int

	Payload models.GenericError
}

// IsSuccess returns true when this update log ingest config default response has a 2xx status code
func (o *UpdateLogIngestConfigDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this update log ingest config default response has a 3xx status code
func (o *UpdateLogIngestConfigDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this update log ingest config default response has a 4xx status code
func (o *UpdateLogIngestConfigDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this update log ingest config default response has a 5xx status code
func (o *UpdateLogIngestConfigDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this update log ingest config default response a status code equal to that given
func (o *UpdateLogIngestConfigDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the update log ingest config default response
func (o *UpdateLogIngestConfigDefault) Code() int {
	return o._statusCode
}

func (o *UpdateLogIngestConfigDefault) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[PUT /api/v1/config/log-ingest-config][%d] UpdateLogIngestConfig default %s", o._statusCode, payload)
}

func (o *UpdateLogIngestConfigDefault) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[PUT /api/v1/config/log-ingest-config][%d] UpdateLogIngestConfig default %s", o._statusCode, payload)
}

func (o *UpdateLogIngestConfigDefault) GetPayload() models.GenericError {
	return o.Payload
}

func (o *UpdateLogIngestConfigDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
