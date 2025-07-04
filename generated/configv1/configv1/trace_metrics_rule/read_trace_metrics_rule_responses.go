// Code generated by go-swagger; DO NOT EDIT.

package trace_metrics_rule

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

// ReadTraceMetricsRuleReader is a Reader for the ReadTraceMetricsRule structure.
type ReadTraceMetricsRuleReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ReadTraceMetricsRuleReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewReadTraceMetricsRuleOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 404:
		result := NewReadTraceMetricsRuleNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewReadTraceMetricsRuleInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewReadTraceMetricsRuleDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewReadTraceMetricsRuleOK creates a ReadTraceMetricsRuleOK with default headers values
func NewReadTraceMetricsRuleOK() *ReadTraceMetricsRuleOK {
	return &ReadTraceMetricsRuleOK{}
}

/*
ReadTraceMetricsRuleOK describes a response with status code 200, with default header values.

A successful response.
*/
type ReadTraceMetricsRuleOK struct {
	Payload *models.Configv1ReadTraceMetricsRuleResponse
}

// IsSuccess returns true when this read trace metrics rule o k response has a 2xx status code
func (o *ReadTraceMetricsRuleOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this read trace metrics rule o k response has a 3xx status code
func (o *ReadTraceMetricsRuleOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this read trace metrics rule o k response has a 4xx status code
func (o *ReadTraceMetricsRuleOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this read trace metrics rule o k response has a 5xx status code
func (o *ReadTraceMetricsRuleOK) IsServerError() bool {
	return false
}

// IsCode returns true when this read trace metrics rule o k response a status code equal to that given
func (o *ReadTraceMetricsRuleOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the read trace metrics rule o k response
func (o *ReadTraceMetricsRuleOK) Code() int {
	return 200
}

func (o *ReadTraceMetricsRuleOK) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /api/v1/config/trace-metrics-rules/{slug}][%d] readTraceMetricsRuleOK %s", 200, payload)
}

func (o *ReadTraceMetricsRuleOK) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /api/v1/config/trace-metrics-rules/{slug}][%d] readTraceMetricsRuleOK %s", 200, payload)
}

func (o *ReadTraceMetricsRuleOK) GetPayload() *models.Configv1ReadTraceMetricsRuleResponse {
	return o.Payload
}

func (o *ReadTraceMetricsRuleOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Configv1ReadTraceMetricsRuleResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewReadTraceMetricsRuleNotFound creates a ReadTraceMetricsRuleNotFound with default headers values
func NewReadTraceMetricsRuleNotFound() *ReadTraceMetricsRuleNotFound {
	return &ReadTraceMetricsRuleNotFound{}
}

/*
ReadTraceMetricsRuleNotFound describes a response with status code 404, with default header values.

Cannot read the TraceMetricsRule because the slug does not exist.
*/
type ReadTraceMetricsRuleNotFound struct {
	Payload *models.APIError
}

// IsSuccess returns true when this read trace metrics rule not found response has a 2xx status code
func (o *ReadTraceMetricsRuleNotFound) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this read trace metrics rule not found response has a 3xx status code
func (o *ReadTraceMetricsRuleNotFound) IsRedirect() bool {
	return false
}

// IsClientError returns true when this read trace metrics rule not found response has a 4xx status code
func (o *ReadTraceMetricsRuleNotFound) IsClientError() bool {
	return true
}

// IsServerError returns true when this read trace metrics rule not found response has a 5xx status code
func (o *ReadTraceMetricsRuleNotFound) IsServerError() bool {
	return false
}

// IsCode returns true when this read trace metrics rule not found response a status code equal to that given
func (o *ReadTraceMetricsRuleNotFound) IsCode(code int) bool {
	return code == 404
}

// Code gets the status code for the read trace metrics rule not found response
func (o *ReadTraceMetricsRuleNotFound) Code() int {
	return 404
}

func (o *ReadTraceMetricsRuleNotFound) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /api/v1/config/trace-metrics-rules/{slug}][%d] readTraceMetricsRuleNotFound %s", 404, payload)
}

func (o *ReadTraceMetricsRuleNotFound) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /api/v1/config/trace-metrics-rules/{slug}][%d] readTraceMetricsRuleNotFound %s", 404, payload)
}

func (o *ReadTraceMetricsRuleNotFound) GetPayload() *models.APIError {
	return o.Payload
}

func (o *ReadTraceMetricsRuleNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.APIError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewReadTraceMetricsRuleInternalServerError creates a ReadTraceMetricsRuleInternalServerError with default headers values
func NewReadTraceMetricsRuleInternalServerError() *ReadTraceMetricsRuleInternalServerError {
	return &ReadTraceMetricsRuleInternalServerError{}
}

/*
ReadTraceMetricsRuleInternalServerError describes a response with status code 500, with default header values.

An unexpected error response.
*/
type ReadTraceMetricsRuleInternalServerError struct {
	Payload *models.APIError
}

// IsSuccess returns true when this read trace metrics rule internal server error response has a 2xx status code
func (o *ReadTraceMetricsRuleInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this read trace metrics rule internal server error response has a 3xx status code
func (o *ReadTraceMetricsRuleInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this read trace metrics rule internal server error response has a 4xx status code
func (o *ReadTraceMetricsRuleInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this read trace metrics rule internal server error response has a 5xx status code
func (o *ReadTraceMetricsRuleInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this read trace metrics rule internal server error response a status code equal to that given
func (o *ReadTraceMetricsRuleInternalServerError) IsCode(code int) bool {
	return code == 500
}

// Code gets the status code for the read trace metrics rule internal server error response
func (o *ReadTraceMetricsRuleInternalServerError) Code() int {
	return 500
}

func (o *ReadTraceMetricsRuleInternalServerError) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /api/v1/config/trace-metrics-rules/{slug}][%d] readTraceMetricsRuleInternalServerError %s", 500, payload)
}

func (o *ReadTraceMetricsRuleInternalServerError) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /api/v1/config/trace-metrics-rules/{slug}][%d] readTraceMetricsRuleInternalServerError %s", 500, payload)
}

func (o *ReadTraceMetricsRuleInternalServerError) GetPayload() *models.APIError {
	return o.Payload
}

func (o *ReadTraceMetricsRuleInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.APIError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewReadTraceMetricsRuleDefault creates a ReadTraceMetricsRuleDefault with default headers values
func NewReadTraceMetricsRuleDefault(code int) *ReadTraceMetricsRuleDefault {
	return &ReadTraceMetricsRuleDefault{
		_statusCode: code,
	}
}

/*
ReadTraceMetricsRuleDefault describes a response with status code -1, with default header values.

An undefined error response.
*/
type ReadTraceMetricsRuleDefault struct {
	_statusCode int

	Payload models.GenericError
}

// IsSuccess returns true when this read trace metrics rule default response has a 2xx status code
func (o *ReadTraceMetricsRuleDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this read trace metrics rule default response has a 3xx status code
func (o *ReadTraceMetricsRuleDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this read trace metrics rule default response has a 4xx status code
func (o *ReadTraceMetricsRuleDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this read trace metrics rule default response has a 5xx status code
func (o *ReadTraceMetricsRuleDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this read trace metrics rule default response a status code equal to that given
func (o *ReadTraceMetricsRuleDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the read trace metrics rule default response
func (o *ReadTraceMetricsRuleDefault) Code() int {
	return o._statusCode
}

func (o *ReadTraceMetricsRuleDefault) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /api/v1/config/trace-metrics-rules/{slug}][%d] ReadTraceMetricsRule default %s", o._statusCode, payload)
}

func (o *ReadTraceMetricsRuleDefault) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /api/v1/config/trace-metrics-rules/{slug}][%d] ReadTraceMetricsRule default %s", o._statusCode, payload)
}

func (o *ReadTraceMetricsRuleDefault) GetPayload() models.GenericError {
	return o.Payload
}

func (o *ReadTraceMetricsRuleDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
