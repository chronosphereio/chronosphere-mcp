// Code generated by go-swagger; DO NOT EDIT.

package trace_behavior

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

// ReadTraceBehaviorReader is a Reader for the ReadTraceBehavior structure.
type ReadTraceBehaviorReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ReadTraceBehaviorReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewReadTraceBehaviorOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 404:
		result := NewReadTraceBehaviorNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewReadTraceBehaviorInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewReadTraceBehaviorDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewReadTraceBehaviorOK creates a ReadTraceBehaviorOK with default headers values
func NewReadTraceBehaviorOK() *ReadTraceBehaviorOK {
	return &ReadTraceBehaviorOK{}
}

/*
ReadTraceBehaviorOK describes a response with status code 200, with default header values.

A successful response.
*/
type ReadTraceBehaviorOK struct {
	Payload *models.Configv1ReadTraceBehaviorResponse
}

// IsSuccess returns true when this read trace behavior o k response has a 2xx status code
func (o *ReadTraceBehaviorOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this read trace behavior o k response has a 3xx status code
func (o *ReadTraceBehaviorOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this read trace behavior o k response has a 4xx status code
func (o *ReadTraceBehaviorOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this read trace behavior o k response has a 5xx status code
func (o *ReadTraceBehaviorOK) IsServerError() bool {
	return false
}

// IsCode returns true when this read trace behavior o k response a status code equal to that given
func (o *ReadTraceBehaviorOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the read trace behavior o k response
func (o *ReadTraceBehaviorOK) Code() int {
	return 200
}

func (o *ReadTraceBehaviorOK) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /api/v1/config/trace-behaviors/{slug}][%d] readTraceBehaviorOK %s", 200, payload)
}

func (o *ReadTraceBehaviorOK) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /api/v1/config/trace-behaviors/{slug}][%d] readTraceBehaviorOK %s", 200, payload)
}

func (o *ReadTraceBehaviorOK) GetPayload() *models.Configv1ReadTraceBehaviorResponse {
	return o.Payload
}

func (o *ReadTraceBehaviorOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Configv1ReadTraceBehaviorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewReadTraceBehaviorNotFound creates a ReadTraceBehaviorNotFound with default headers values
func NewReadTraceBehaviorNotFound() *ReadTraceBehaviorNotFound {
	return &ReadTraceBehaviorNotFound{}
}

/*
ReadTraceBehaviorNotFound describes a response with status code 404, with default header values.

Cannot read the TraceBehavior because the slug does not exist.
*/
type ReadTraceBehaviorNotFound struct {
	Payload *models.APIError
}

// IsSuccess returns true when this read trace behavior not found response has a 2xx status code
func (o *ReadTraceBehaviorNotFound) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this read trace behavior not found response has a 3xx status code
func (o *ReadTraceBehaviorNotFound) IsRedirect() bool {
	return false
}

// IsClientError returns true when this read trace behavior not found response has a 4xx status code
func (o *ReadTraceBehaviorNotFound) IsClientError() bool {
	return true
}

// IsServerError returns true when this read trace behavior not found response has a 5xx status code
func (o *ReadTraceBehaviorNotFound) IsServerError() bool {
	return false
}

// IsCode returns true when this read trace behavior not found response a status code equal to that given
func (o *ReadTraceBehaviorNotFound) IsCode(code int) bool {
	return code == 404
}

// Code gets the status code for the read trace behavior not found response
func (o *ReadTraceBehaviorNotFound) Code() int {
	return 404
}

func (o *ReadTraceBehaviorNotFound) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /api/v1/config/trace-behaviors/{slug}][%d] readTraceBehaviorNotFound %s", 404, payload)
}

func (o *ReadTraceBehaviorNotFound) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /api/v1/config/trace-behaviors/{slug}][%d] readTraceBehaviorNotFound %s", 404, payload)
}

func (o *ReadTraceBehaviorNotFound) GetPayload() *models.APIError {
	return o.Payload
}

func (o *ReadTraceBehaviorNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.APIError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewReadTraceBehaviorInternalServerError creates a ReadTraceBehaviorInternalServerError with default headers values
func NewReadTraceBehaviorInternalServerError() *ReadTraceBehaviorInternalServerError {
	return &ReadTraceBehaviorInternalServerError{}
}

/*
ReadTraceBehaviorInternalServerError describes a response with status code 500, with default header values.

An unexpected error response.
*/
type ReadTraceBehaviorInternalServerError struct {
	Payload *models.APIError
}

// IsSuccess returns true when this read trace behavior internal server error response has a 2xx status code
func (o *ReadTraceBehaviorInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this read trace behavior internal server error response has a 3xx status code
func (o *ReadTraceBehaviorInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this read trace behavior internal server error response has a 4xx status code
func (o *ReadTraceBehaviorInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this read trace behavior internal server error response has a 5xx status code
func (o *ReadTraceBehaviorInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this read trace behavior internal server error response a status code equal to that given
func (o *ReadTraceBehaviorInternalServerError) IsCode(code int) bool {
	return code == 500
}

// Code gets the status code for the read trace behavior internal server error response
func (o *ReadTraceBehaviorInternalServerError) Code() int {
	return 500
}

func (o *ReadTraceBehaviorInternalServerError) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /api/v1/config/trace-behaviors/{slug}][%d] readTraceBehaviorInternalServerError %s", 500, payload)
}

func (o *ReadTraceBehaviorInternalServerError) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /api/v1/config/trace-behaviors/{slug}][%d] readTraceBehaviorInternalServerError %s", 500, payload)
}

func (o *ReadTraceBehaviorInternalServerError) GetPayload() *models.APIError {
	return o.Payload
}

func (o *ReadTraceBehaviorInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.APIError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewReadTraceBehaviorDefault creates a ReadTraceBehaviorDefault with default headers values
func NewReadTraceBehaviorDefault(code int) *ReadTraceBehaviorDefault {
	return &ReadTraceBehaviorDefault{
		_statusCode: code,
	}
}

/*
ReadTraceBehaviorDefault describes a response with status code -1, with default header values.

An undefined error response.
*/
type ReadTraceBehaviorDefault struct {
	_statusCode int

	Payload models.GenericError
}

// IsSuccess returns true when this read trace behavior default response has a 2xx status code
func (o *ReadTraceBehaviorDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this read trace behavior default response has a 3xx status code
func (o *ReadTraceBehaviorDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this read trace behavior default response has a 4xx status code
func (o *ReadTraceBehaviorDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this read trace behavior default response has a 5xx status code
func (o *ReadTraceBehaviorDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this read trace behavior default response a status code equal to that given
func (o *ReadTraceBehaviorDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the read trace behavior default response
func (o *ReadTraceBehaviorDefault) Code() int {
	return o._statusCode
}

func (o *ReadTraceBehaviorDefault) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /api/v1/config/trace-behaviors/{slug}][%d] ReadTraceBehavior default %s", o._statusCode, payload)
}

func (o *ReadTraceBehaviorDefault) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /api/v1/config/trace-behaviors/{slug}][%d] ReadTraceBehavior default %s", o._statusCode, payload)
}

func (o *ReadTraceBehaviorDefault) GetPayload() models.GenericError {
	return o.Payload
}

func (o *ReadTraceBehaviorDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
