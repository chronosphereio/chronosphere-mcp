// Code generated by go-swagger; DO NOT EDIT.

package monitor

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

// ListMonitorsReader is a Reader for the ListMonitors structure.
type ListMonitorsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ListMonitorsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewListMonitorsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 500:
		result := NewListMonitorsInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewListMonitorsDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewListMonitorsOK creates a ListMonitorsOK with default headers values
func NewListMonitorsOK() *ListMonitorsOK {
	return &ListMonitorsOK{}
}

/*
ListMonitorsOK describes a response with status code 200, with default header values.

A successful response.
*/
type ListMonitorsOK struct {
	Payload *models.Configv1ListMonitorsResponse
}

// IsSuccess returns true when this list monitors o k response has a 2xx status code
func (o *ListMonitorsOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this list monitors o k response has a 3xx status code
func (o *ListMonitorsOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this list monitors o k response has a 4xx status code
func (o *ListMonitorsOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this list monitors o k response has a 5xx status code
func (o *ListMonitorsOK) IsServerError() bool {
	return false
}

// IsCode returns true when this list monitors o k response a status code equal to that given
func (o *ListMonitorsOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the list monitors o k response
func (o *ListMonitorsOK) Code() int {
	return 200
}

func (o *ListMonitorsOK) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /api/v1/config/monitors][%d] listMonitorsOK %s", 200, payload)
}

func (o *ListMonitorsOK) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /api/v1/config/monitors][%d] listMonitorsOK %s", 200, payload)
}

func (o *ListMonitorsOK) GetPayload() *models.Configv1ListMonitorsResponse {
	return o.Payload
}

func (o *ListMonitorsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Configv1ListMonitorsResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewListMonitorsInternalServerError creates a ListMonitorsInternalServerError with default headers values
func NewListMonitorsInternalServerError() *ListMonitorsInternalServerError {
	return &ListMonitorsInternalServerError{}
}

/*
ListMonitorsInternalServerError describes a response with status code 500, with default header values.

An unexpected error response.
*/
type ListMonitorsInternalServerError struct {
	Payload *models.APIError
}

// IsSuccess returns true when this list monitors internal server error response has a 2xx status code
func (o *ListMonitorsInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this list monitors internal server error response has a 3xx status code
func (o *ListMonitorsInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this list monitors internal server error response has a 4xx status code
func (o *ListMonitorsInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this list monitors internal server error response has a 5xx status code
func (o *ListMonitorsInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this list monitors internal server error response a status code equal to that given
func (o *ListMonitorsInternalServerError) IsCode(code int) bool {
	return code == 500
}

// Code gets the status code for the list monitors internal server error response
func (o *ListMonitorsInternalServerError) Code() int {
	return 500
}

func (o *ListMonitorsInternalServerError) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /api/v1/config/monitors][%d] listMonitorsInternalServerError %s", 500, payload)
}

func (o *ListMonitorsInternalServerError) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /api/v1/config/monitors][%d] listMonitorsInternalServerError %s", 500, payload)
}

func (o *ListMonitorsInternalServerError) GetPayload() *models.APIError {
	return o.Payload
}

func (o *ListMonitorsInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.APIError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewListMonitorsDefault creates a ListMonitorsDefault with default headers values
func NewListMonitorsDefault(code int) *ListMonitorsDefault {
	return &ListMonitorsDefault{
		_statusCode: code,
	}
}

/*
ListMonitorsDefault describes a response with status code -1, with default header values.

An undefined error response.
*/
type ListMonitorsDefault struct {
	_statusCode int

	Payload models.GenericError
}

// IsSuccess returns true when this list monitors default response has a 2xx status code
func (o *ListMonitorsDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this list monitors default response has a 3xx status code
func (o *ListMonitorsDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this list monitors default response has a 4xx status code
func (o *ListMonitorsDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this list monitors default response has a 5xx status code
func (o *ListMonitorsDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this list monitors default response a status code equal to that given
func (o *ListMonitorsDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the list monitors default response
func (o *ListMonitorsDefault) Code() int {
	return o._statusCode
}

func (o *ListMonitorsDefault) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /api/v1/config/monitors][%d] ListMonitors default %s", o._statusCode, payload)
}

func (o *ListMonitorsDefault) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /api/v1/config/monitors][%d] ListMonitors default %s", o._statusCode, payload)
}

func (o *ListMonitorsDefault) GetPayload() models.GenericError {
	return o.Payload
}

func (o *ListMonitorsDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
