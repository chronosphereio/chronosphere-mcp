// Code generated by go-swagger; DO NOT EDIT.

package s_l_o

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

// ReadSLOReader is a Reader for the ReadSLO structure.
type ReadSLOReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ReadSLOReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewReadSLOOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 404:
		result := NewReadSLONotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewReadSLOInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewReadSLODefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewReadSLOOK creates a ReadSLOOK with default headers values
func NewReadSLOOK() *ReadSLOOK {
	return &ReadSLOOK{}
}

/*
ReadSLOOK describes a response with status code 200, with default header values.

A successful response.
*/
type ReadSLOOK struct {
	Payload *models.Configv1ReadSLOResponse
}

// IsSuccess returns true when this read s l o o k response has a 2xx status code
func (o *ReadSLOOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this read s l o o k response has a 3xx status code
func (o *ReadSLOOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this read s l o o k response has a 4xx status code
func (o *ReadSLOOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this read s l o o k response has a 5xx status code
func (o *ReadSLOOK) IsServerError() bool {
	return false
}

// IsCode returns true when this read s l o o k response a status code equal to that given
func (o *ReadSLOOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the read s l o o k response
func (o *ReadSLOOK) Code() int {
	return 200
}

func (o *ReadSLOOK) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /api/v1/config/slos/{slug}][%d] readSLOOK %s", 200, payload)
}

func (o *ReadSLOOK) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /api/v1/config/slos/{slug}][%d] readSLOOK %s", 200, payload)
}

func (o *ReadSLOOK) GetPayload() *models.Configv1ReadSLOResponse {
	return o.Payload
}

func (o *ReadSLOOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Configv1ReadSLOResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewReadSLONotFound creates a ReadSLONotFound with default headers values
func NewReadSLONotFound() *ReadSLONotFound {
	return &ReadSLONotFound{}
}

/*
ReadSLONotFound describes a response with status code 404, with default header values.

Cannot read the SLO because the slug does not exist.
*/
type ReadSLONotFound struct {
	Payload *models.APIError
}

// IsSuccess returns true when this read s l o not found response has a 2xx status code
func (o *ReadSLONotFound) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this read s l o not found response has a 3xx status code
func (o *ReadSLONotFound) IsRedirect() bool {
	return false
}

// IsClientError returns true when this read s l o not found response has a 4xx status code
func (o *ReadSLONotFound) IsClientError() bool {
	return true
}

// IsServerError returns true when this read s l o not found response has a 5xx status code
func (o *ReadSLONotFound) IsServerError() bool {
	return false
}

// IsCode returns true when this read s l o not found response a status code equal to that given
func (o *ReadSLONotFound) IsCode(code int) bool {
	return code == 404
}

// Code gets the status code for the read s l o not found response
func (o *ReadSLONotFound) Code() int {
	return 404
}

func (o *ReadSLONotFound) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /api/v1/config/slos/{slug}][%d] readSLONotFound %s", 404, payload)
}

func (o *ReadSLONotFound) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /api/v1/config/slos/{slug}][%d] readSLONotFound %s", 404, payload)
}

func (o *ReadSLONotFound) GetPayload() *models.APIError {
	return o.Payload
}

func (o *ReadSLONotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.APIError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewReadSLOInternalServerError creates a ReadSLOInternalServerError with default headers values
func NewReadSLOInternalServerError() *ReadSLOInternalServerError {
	return &ReadSLOInternalServerError{}
}

/*
ReadSLOInternalServerError describes a response with status code 500, with default header values.

An unexpected error response.
*/
type ReadSLOInternalServerError struct {
	Payload *models.APIError
}

// IsSuccess returns true when this read s l o internal server error response has a 2xx status code
func (o *ReadSLOInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this read s l o internal server error response has a 3xx status code
func (o *ReadSLOInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this read s l o internal server error response has a 4xx status code
func (o *ReadSLOInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this read s l o internal server error response has a 5xx status code
func (o *ReadSLOInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this read s l o internal server error response a status code equal to that given
func (o *ReadSLOInternalServerError) IsCode(code int) bool {
	return code == 500
}

// Code gets the status code for the read s l o internal server error response
func (o *ReadSLOInternalServerError) Code() int {
	return 500
}

func (o *ReadSLOInternalServerError) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /api/v1/config/slos/{slug}][%d] readSLOInternalServerError %s", 500, payload)
}

func (o *ReadSLOInternalServerError) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /api/v1/config/slos/{slug}][%d] readSLOInternalServerError %s", 500, payload)
}

func (o *ReadSLOInternalServerError) GetPayload() *models.APIError {
	return o.Payload
}

func (o *ReadSLOInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.APIError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewReadSLODefault creates a ReadSLODefault with default headers values
func NewReadSLODefault(code int) *ReadSLODefault {
	return &ReadSLODefault{
		_statusCode: code,
	}
}

/*
ReadSLODefault describes a response with status code -1, with default header values.

An undefined error response.
*/
type ReadSLODefault struct {
	_statusCode int

	Payload models.GenericError
}

// IsSuccess returns true when this read s l o default response has a 2xx status code
func (o *ReadSLODefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this read s l o default response has a 3xx status code
func (o *ReadSLODefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this read s l o default response has a 4xx status code
func (o *ReadSLODefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this read s l o default response has a 5xx status code
func (o *ReadSLODefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this read s l o default response a status code equal to that given
func (o *ReadSLODefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the read s l o default response
func (o *ReadSLODefault) Code() int {
	return o._statusCode
}

func (o *ReadSLODefault) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /api/v1/config/slos/{slug}][%d] ReadSLO default %s", o._statusCode, payload)
}

func (o *ReadSLODefault) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /api/v1/config/slos/{slug}][%d] ReadSLO default %s", o._statusCode, payload)
}

func (o *ReadSLODefault) GetPayload() models.GenericError {
	return o.Payload
}

func (o *ReadSLODefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
