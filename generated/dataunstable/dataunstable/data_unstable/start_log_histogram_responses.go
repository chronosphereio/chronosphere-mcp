// Code generated by go-swagger; DO NOT EDIT.

package data_unstable

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/chronosphereio/chronosphere-mcp/generated/dataunstable/models"
)

// StartLogHistogramReader is a Reader for the StartLogHistogram structure.
type StartLogHistogramReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *StartLogHistogramReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewStartLogHistogramOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewStartLogHistogramDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewStartLogHistogramOK creates a StartLogHistogramOK with default headers values
func NewStartLogHistogramOK() *StartLogHistogramOK {
	return &StartLogHistogramOK{}
}

/*
StartLogHistogramOK describes a response with status code 200, with default header values.

A successful response.
*/
type StartLogHistogramOK struct {
	Payload *models.DataunstableStartLogHistogramResponse
}

// IsSuccess returns true when this start log histogram o k response has a 2xx status code
func (o *StartLogHistogramOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this start log histogram o k response has a 3xx status code
func (o *StartLogHistogramOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this start log histogram o k response has a 4xx status code
func (o *StartLogHistogramOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this start log histogram o k response has a 5xx status code
func (o *StartLogHistogramOK) IsServerError() bool {
	return false
}

// IsCode returns true when this start log histogram o k response a status code equal to that given
func (o *StartLogHistogramOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the start log histogram o k response
func (o *StartLogHistogramOK) Code() int {
	return 200
}

func (o *StartLogHistogramOK) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /api/unstable/data/logs:histogram-start][%d] startLogHistogramOK %s", 200, payload)
}

func (o *StartLogHistogramOK) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /api/unstable/data/logs:histogram-start][%d] startLogHistogramOK %s", 200, payload)
}

func (o *StartLogHistogramOK) GetPayload() *models.DataunstableStartLogHistogramResponse {
	return o.Payload
}

func (o *StartLogHistogramOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.DataunstableStartLogHistogramResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewStartLogHistogramDefault creates a StartLogHistogramDefault with default headers values
func NewStartLogHistogramDefault(code int) *StartLogHistogramDefault {
	return &StartLogHistogramDefault{
		_statusCode: code,
	}
}

/*
StartLogHistogramDefault describes a response with status code -1, with default header values.

An unexpected error response.
*/
type StartLogHistogramDefault struct {
	_statusCode int

	Payload *models.GooglerpcStatus
}

// IsSuccess returns true when this start log histogram default response has a 2xx status code
func (o *StartLogHistogramDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this start log histogram default response has a 3xx status code
func (o *StartLogHistogramDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this start log histogram default response has a 4xx status code
func (o *StartLogHistogramDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this start log histogram default response has a 5xx status code
func (o *StartLogHistogramDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this start log histogram default response a status code equal to that given
func (o *StartLogHistogramDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the start log histogram default response
func (o *StartLogHistogramDefault) Code() int {
	return o._statusCode
}

func (o *StartLogHistogramDefault) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /api/unstable/data/logs:histogram-start][%d] StartLogHistogram default %s", o._statusCode, payload)
}

func (o *StartLogHistogramDefault) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /api/unstable/data/logs:histogram-start][%d] StartLogHistogram default %s", o._statusCode, payload)
}

func (o *StartLogHistogramDefault) GetPayload() *models.GooglerpcStatus {
	return o.Payload
}

func (o *StartLogHistogramDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.GooglerpcStatus)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
