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

// PollListLogsQueryReader is a Reader for the PollListLogsQuery structure.
type PollListLogsQueryReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *PollListLogsQueryReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewPollListLogsQueryOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewPollListLogsQueryDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewPollListLogsQueryOK creates a PollListLogsQueryOK with default headers values
func NewPollListLogsQueryOK() *PollListLogsQueryOK {
	return &PollListLogsQueryOK{}
}

/*
PollListLogsQueryOK describes a response with status code 200, with default header values.

A successful response.
*/
type PollListLogsQueryOK struct {
	Payload *models.DataunstablePollListLogsResponse
}

// IsSuccess returns true when this poll list logs query o k response has a 2xx status code
func (o *PollListLogsQueryOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this poll list logs query o k response has a 3xx status code
func (o *PollListLogsQueryOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this poll list logs query o k response has a 4xx status code
func (o *PollListLogsQueryOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this poll list logs query o k response has a 5xx status code
func (o *PollListLogsQueryOK) IsServerError() bool {
	return false
}

// IsCode returns true when this poll list logs query o k response a status code equal to that given
func (o *PollListLogsQueryOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the poll list logs query o k response
func (o *PollListLogsQueryOK) Code() int {
	return 200
}

func (o *PollListLogsQueryOK) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /api/unstable/data/logs:list-poll][%d] pollListLogsQueryOK %s", 200, payload)
}

func (o *PollListLogsQueryOK) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /api/unstable/data/logs:list-poll][%d] pollListLogsQueryOK %s", 200, payload)
}

func (o *PollListLogsQueryOK) GetPayload() *models.DataunstablePollListLogsResponse {
	return o.Payload
}

func (o *PollListLogsQueryOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.DataunstablePollListLogsResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPollListLogsQueryDefault creates a PollListLogsQueryDefault with default headers values
func NewPollListLogsQueryDefault(code int) *PollListLogsQueryDefault {
	return &PollListLogsQueryDefault{
		_statusCode: code,
	}
}

/*
PollListLogsQueryDefault describes a response with status code -1, with default header values.

An unexpected error response.
*/
type PollListLogsQueryDefault struct {
	_statusCode int

	Payload *models.GooglerpcStatus
}

// IsSuccess returns true when this poll list logs query default response has a 2xx status code
func (o *PollListLogsQueryDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this poll list logs query default response has a 3xx status code
func (o *PollListLogsQueryDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this poll list logs query default response has a 4xx status code
func (o *PollListLogsQueryDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this poll list logs query default response has a 5xx status code
func (o *PollListLogsQueryDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this poll list logs query default response a status code equal to that given
func (o *PollListLogsQueryDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the poll list logs query default response
func (o *PollListLogsQueryDefault) Code() int {
	return o._statusCode
}

func (o *PollListLogsQueryDefault) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /api/unstable/data/logs:list-poll][%d] PollListLogsQuery default %s", o._statusCode, payload)
}

func (o *PollListLogsQueryDefault) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /api/unstable/data/logs:list-poll][%d] PollListLogsQuery default %s", o._statusCode, payload)
}

func (o *PollListLogsQueryDefault) GetPayload() *models.GooglerpcStatus {
	return o.Payload
}

func (o *PollListLogsQueryDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.GooglerpcStatus)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
