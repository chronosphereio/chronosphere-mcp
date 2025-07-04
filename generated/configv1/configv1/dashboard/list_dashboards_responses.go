// Code generated by go-swagger; DO NOT EDIT.

package dashboard

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

// ListDashboardsReader is a Reader for the ListDashboards structure.
type ListDashboardsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ListDashboardsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewListDashboardsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 500:
		result := NewListDashboardsInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewListDashboardsDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewListDashboardsOK creates a ListDashboardsOK with default headers values
func NewListDashboardsOK() *ListDashboardsOK {
	return &ListDashboardsOK{}
}

/*
ListDashboardsOK describes a response with status code 200, with default header values.

A successful response.
*/
type ListDashboardsOK struct {
	Payload *models.Configv1ListDashboardsResponse
}

// IsSuccess returns true when this list dashboards o k response has a 2xx status code
func (o *ListDashboardsOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this list dashboards o k response has a 3xx status code
func (o *ListDashboardsOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this list dashboards o k response has a 4xx status code
func (o *ListDashboardsOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this list dashboards o k response has a 5xx status code
func (o *ListDashboardsOK) IsServerError() bool {
	return false
}

// IsCode returns true when this list dashboards o k response a status code equal to that given
func (o *ListDashboardsOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the list dashboards o k response
func (o *ListDashboardsOK) Code() int {
	return 200
}

func (o *ListDashboardsOK) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /api/v1/config/dashboards][%d] listDashboardsOK %s", 200, payload)
}

func (o *ListDashboardsOK) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /api/v1/config/dashboards][%d] listDashboardsOK %s", 200, payload)
}

func (o *ListDashboardsOK) GetPayload() *models.Configv1ListDashboardsResponse {
	return o.Payload
}

func (o *ListDashboardsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Configv1ListDashboardsResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewListDashboardsInternalServerError creates a ListDashboardsInternalServerError with default headers values
func NewListDashboardsInternalServerError() *ListDashboardsInternalServerError {
	return &ListDashboardsInternalServerError{}
}

/*
ListDashboardsInternalServerError describes a response with status code 500, with default header values.

An unexpected error response.
*/
type ListDashboardsInternalServerError struct {
	Payload *models.APIError
}

// IsSuccess returns true when this list dashboards internal server error response has a 2xx status code
func (o *ListDashboardsInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this list dashboards internal server error response has a 3xx status code
func (o *ListDashboardsInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this list dashboards internal server error response has a 4xx status code
func (o *ListDashboardsInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this list dashboards internal server error response has a 5xx status code
func (o *ListDashboardsInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this list dashboards internal server error response a status code equal to that given
func (o *ListDashboardsInternalServerError) IsCode(code int) bool {
	return code == 500
}

// Code gets the status code for the list dashboards internal server error response
func (o *ListDashboardsInternalServerError) Code() int {
	return 500
}

func (o *ListDashboardsInternalServerError) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /api/v1/config/dashboards][%d] listDashboardsInternalServerError %s", 500, payload)
}

func (o *ListDashboardsInternalServerError) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /api/v1/config/dashboards][%d] listDashboardsInternalServerError %s", 500, payload)
}

func (o *ListDashboardsInternalServerError) GetPayload() *models.APIError {
	return o.Payload
}

func (o *ListDashboardsInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.APIError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewListDashboardsDefault creates a ListDashboardsDefault with default headers values
func NewListDashboardsDefault(code int) *ListDashboardsDefault {
	return &ListDashboardsDefault{
		_statusCode: code,
	}
}

/*
ListDashboardsDefault describes a response with status code -1, with default header values.

An undefined error response.
*/
type ListDashboardsDefault struct {
	_statusCode int

	Payload models.GenericError
}

// IsSuccess returns true when this list dashboards default response has a 2xx status code
func (o *ListDashboardsDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this list dashboards default response has a 3xx status code
func (o *ListDashboardsDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this list dashboards default response has a 4xx status code
func (o *ListDashboardsDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this list dashboards default response has a 5xx status code
func (o *ListDashboardsDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this list dashboards default response a status code equal to that given
func (o *ListDashboardsDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the list dashboards default response
func (o *ListDashboardsDefault) Code() int {
	return o._statusCode
}

func (o *ListDashboardsDefault) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /api/v1/config/dashboards][%d] ListDashboards default %s", o._statusCode, payload)
}

func (o *ListDashboardsDefault) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /api/v1/config/dashboards][%d] ListDashboards default %s", o._statusCode, payload)
}

func (o *ListDashboardsDefault) GetPayload() models.GenericError {
	return o.Payload
}

func (o *ListDashboardsDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
