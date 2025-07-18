// Code generated by go-swagger; DO NOT EDIT.

package muting_rule

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

// DeleteMutingRuleReader is a Reader for the DeleteMutingRule structure.
type DeleteMutingRuleReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *DeleteMutingRuleReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewDeleteMutingRuleOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewDeleteMutingRuleBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewDeleteMutingRuleNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewDeleteMutingRuleInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewDeleteMutingRuleDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewDeleteMutingRuleOK creates a DeleteMutingRuleOK with default headers values
func NewDeleteMutingRuleOK() *DeleteMutingRuleOK {
	return &DeleteMutingRuleOK{}
}

/*
DeleteMutingRuleOK describes a response with status code 200, with default header values.

A successful response.
*/
type DeleteMutingRuleOK struct {
	Payload models.Configv1DeleteMutingRuleResponse
}

// IsSuccess returns true when this delete muting rule o k response has a 2xx status code
func (o *DeleteMutingRuleOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this delete muting rule o k response has a 3xx status code
func (o *DeleteMutingRuleOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this delete muting rule o k response has a 4xx status code
func (o *DeleteMutingRuleOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this delete muting rule o k response has a 5xx status code
func (o *DeleteMutingRuleOK) IsServerError() bool {
	return false
}

// IsCode returns true when this delete muting rule o k response a status code equal to that given
func (o *DeleteMutingRuleOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the delete muting rule o k response
func (o *DeleteMutingRuleOK) Code() int {
	return 200
}

func (o *DeleteMutingRuleOK) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[DELETE /api/v1/config/muting-rules/{slug}][%d] deleteMutingRuleOK %s", 200, payload)
}

func (o *DeleteMutingRuleOK) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[DELETE /api/v1/config/muting-rules/{slug}][%d] deleteMutingRuleOK %s", 200, payload)
}

func (o *DeleteMutingRuleOK) GetPayload() models.Configv1DeleteMutingRuleResponse {
	return o.Payload
}

func (o *DeleteMutingRuleOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDeleteMutingRuleBadRequest creates a DeleteMutingRuleBadRequest with default headers values
func NewDeleteMutingRuleBadRequest() *DeleteMutingRuleBadRequest {
	return &DeleteMutingRuleBadRequest{}
}

/*
DeleteMutingRuleBadRequest describes a response with status code 400, with default header values.

Cannot delete the MutingRule because it is in use.
*/
type DeleteMutingRuleBadRequest struct {
	Payload *models.APIError
}

// IsSuccess returns true when this delete muting rule bad request response has a 2xx status code
func (o *DeleteMutingRuleBadRequest) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this delete muting rule bad request response has a 3xx status code
func (o *DeleteMutingRuleBadRequest) IsRedirect() bool {
	return false
}

// IsClientError returns true when this delete muting rule bad request response has a 4xx status code
func (o *DeleteMutingRuleBadRequest) IsClientError() bool {
	return true
}

// IsServerError returns true when this delete muting rule bad request response has a 5xx status code
func (o *DeleteMutingRuleBadRequest) IsServerError() bool {
	return false
}

// IsCode returns true when this delete muting rule bad request response a status code equal to that given
func (o *DeleteMutingRuleBadRequest) IsCode(code int) bool {
	return code == 400
}

// Code gets the status code for the delete muting rule bad request response
func (o *DeleteMutingRuleBadRequest) Code() int {
	return 400
}

func (o *DeleteMutingRuleBadRequest) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[DELETE /api/v1/config/muting-rules/{slug}][%d] deleteMutingRuleBadRequest %s", 400, payload)
}

func (o *DeleteMutingRuleBadRequest) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[DELETE /api/v1/config/muting-rules/{slug}][%d] deleteMutingRuleBadRequest %s", 400, payload)
}

func (o *DeleteMutingRuleBadRequest) GetPayload() *models.APIError {
	return o.Payload
}

func (o *DeleteMutingRuleBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.APIError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDeleteMutingRuleNotFound creates a DeleteMutingRuleNotFound with default headers values
func NewDeleteMutingRuleNotFound() *DeleteMutingRuleNotFound {
	return &DeleteMutingRuleNotFound{}
}

/*
DeleteMutingRuleNotFound describes a response with status code 404, with default header values.

Cannot delete the MutingRule because the slug does not exist.
*/
type DeleteMutingRuleNotFound struct {
	Payload *models.APIError
}

// IsSuccess returns true when this delete muting rule not found response has a 2xx status code
func (o *DeleteMutingRuleNotFound) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this delete muting rule not found response has a 3xx status code
func (o *DeleteMutingRuleNotFound) IsRedirect() bool {
	return false
}

// IsClientError returns true when this delete muting rule not found response has a 4xx status code
func (o *DeleteMutingRuleNotFound) IsClientError() bool {
	return true
}

// IsServerError returns true when this delete muting rule not found response has a 5xx status code
func (o *DeleteMutingRuleNotFound) IsServerError() bool {
	return false
}

// IsCode returns true when this delete muting rule not found response a status code equal to that given
func (o *DeleteMutingRuleNotFound) IsCode(code int) bool {
	return code == 404
}

// Code gets the status code for the delete muting rule not found response
func (o *DeleteMutingRuleNotFound) Code() int {
	return 404
}

func (o *DeleteMutingRuleNotFound) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[DELETE /api/v1/config/muting-rules/{slug}][%d] deleteMutingRuleNotFound %s", 404, payload)
}

func (o *DeleteMutingRuleNotFound) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[DELETE /api/v1/config/muting-rules/{slug}][%d] deleteMutingRuleNotFound %s", 404, payload)
}

func (o *DeleteMutingRuleNotFound) GetPayload() *models.APIError {
	return o.Payload
}

func (o *DeleteMutingRuleNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.APIError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDeleteMutingRuleInternalServerError creates a DeleteMutingRuleInternalServerError with default headers values
func NewDeleteMutingRuleInternalServerError() *DeleteMutingRuleInternalServerError {
	return &DeleteMutingRuleInternalServerError{}
}

/*
DeleteMutingRuleInternalServerError describes a response with status code 500, with default header values.

An unexpected error response.
*/
type DeleteMutingRuleInternalServerError struct {
	Payload *models.APIError
}

// IsSuccess returns true when this delete muting rule internal server error response has a 2xx status code
func (o *DeleteMutingRuleInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this delete muting rule internal server error response has a 3xx status code
func (o *DeleteMutingRuleInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this delete muting rule internal server error response has a 4xx status code
func (o *DeleteMutingRuleInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this delete muting rule internal server error response has a 5xx status code
func (o *DeleteMutingRuleInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this delete muting rule internal server error response a status code equal to that given
func (o *DeleteMutingRuleInternalServerError) IsCode(code int) bool {
	return code == 500
}

// Code gets the status code for the delete muting rule internal server error response
func (o *DeleteMutingRuleInternalServerError) Code() int {
	return 500
}

func (o *DeleteMutingRuleInternalServerError) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[DELETE /api/v1/config/muting-rules/{slug}][%d] deleteMutingRuleInternalServerError %s", 500, payload)
}

func (o *DeleteMutingRuleInternalServerError) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[DELETE /api/v1/config/muting-rules/{slug}][%d] deleteMutingRuleInternalServerError %s", 500, payload)
}

func (o *DeleteMutingRuleInternalServerError) GetPayload() *models.APIError {
	return o.Payload
}

func (o *DeleteMutingRuleInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.APIError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDeleteMutingRuleDefault creates a DeleteMutingRuleDefault with default headers values
func NewDeleteMutingRuleDefault(code int) *DeleteMutingRuleDefault {
	return &DeleteMutingRuleDefault{
		_statusCode: code,
	}
}

/*
DeleteMutingRuleDefault describes a response with status code -1, with default header values.

An undefined error response.
*/
type DeleteMutingRuleDefault struct {
	_statusCode int

	Payload models.GenericError
}

// IsSuccess returns true when this delete muting rule default response has a 2xx status code
func (o *DeleteMutingRuleDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this delete muting rule default response has a 3xx status code
func (o *DeleteMutingRuleDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this delete muting rule default response has a 4xx status code
func (o *DeleteMutingRuleDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this delete muting rule default response has a 5xx status code
func (o *DeleteMutingRuleDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this delete muting rule default response a status code equal to that given
func (o *DeleteMutingRuleDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the delete muting rule default response
func (o *DeleteMutingRuleDefault) Code() int {
	return o._statusCode
}

func (o *DeleteMutingRuleDefault) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[DELETE /api/v1/config/muting-rules/{slug}][%d] DeleteMutingRule default %s", o._statusCode, payload)
}

func (o *DeleteMutingRuleDefault) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[DELETE /api/v1/config/muting-rules/{slug}][%d] DeleteMutingRule default %s", o._statusCode, payload)
}

func (o *DeleteMutingRuleDefault) GetPayload() models.GenericError {
	return o.Payload
}

func (o *DeleteMutingRuleDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
