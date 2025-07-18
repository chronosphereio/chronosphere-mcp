// Code generated by go-swagger; DO NOT EDIT.

package otel_metrics_ingestion

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

// DeleteOtelMetricsIngestionReader is a Reader for the DeleteOtelMetricsIngestion structure.
type DeleteOtelMetricsIngestionReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *DeleteOtelMetricsIngestionReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewDeleteOtelMetricsIngestionOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewDeleteOtelMetricsIngestionBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewDeleteOtelMetricsIngestionNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewDeleteOtelMetricsIngestionInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewDeleteOtelMetricsIngestionDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewDeleteOtelMetricsIngestionOK creates a DeleteOtelMetricsIngestionOK with default headers values
func NewDeleteOtelMetricsIngestionOK() *DeleteOtelMetricsIngestionOK {
	return &DeleteOtelMetricsIngestionOK{}
}

/*
DeleteOtelMetricsIngestionOK describes a response with status code 200, with default header values.

A successful response.
*/
type DeleteOtelMetricsIngestionOK struct {
	Payload models.Configv1DeleteOtelMetricsIngestionResponse
}

// IsSuccess returns true when this delete otel metrics ingestion o k response has a 2xx status code
func (o *DeleteOtelMetricsIngestionOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this delete otel metrics ingestion o k response has a 3xx status code
func (o *DeleteOtelMetricsIngestionOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this delete otel metrics ingestion o k response has a 4xx status code
func (o *DeleteOtelMetricsIngestionOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this delete otel metrics ingestion o k response has a 5xx status code
func (o *DeleteOtelMetricsIngestionOK) IsServerError() bool {
	return false
}

// IsCode returns true when this delete otel metrics ingestion o k response a status code equal to that given
func (o *DeleteOtelMetricsIngestionOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the delete otel metrics ingestion o k response
func (o *DeleteOtelMetricsIngestionOK) Code() int {
	return 200
}

func (o *DeleteOtelMetricsIngestionOK) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[DELETE /api/v1/config/otel-metrics-ingestion][%d] deleteOtelMetricsIngestionOK %s", 200, payload)
}

func (o *DeleteOtelMetricsIngestionOK) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[DELETE /api/v1/config/otel-metrics-ingestion][%d] deleteOtelMetricsIngestionOK %s", 200, payload)
}

func (o *DeleteOtelMetricsIngestionOK) GetPayload() models.Configv1DeleteOtelMetricsIngestionResponse {
	return o.Payload
}

func (o *DeleteOtelMetricsIngestionOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDeleteOtelMetricsIngestionBadRequest creates a DeleteOtelMetricsIngestionBadRequest with default headers values
func NewDeleteOtelMetricsIngestionBadRequest() *DeleteOtelMetricsIngestionBadRequest {
	return &DeleteOtelMetricsIngestionBadRequest{}
}

/*
DeleteOtelMetricsIngestionBadRequest describes a response with status code 400, with default header values.

Cannot delete the OtelMetricsIngestion because it is in use.
*/
type DeleteOtelMetricsIngestionBadRequest struct {
	Payload *models.APIError
}

// IsSuccess returns true when this delete otel metrics ingestion bad request response has a 2xx status code
func (o *DeleteOtelMetricsIngestionBadRequest) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this delete otel metrics ingestion bad request response has a 3xx status code
func (o *DeleteOtelMetricsIngestionBadRequest) IsRedirect() bool {
	return false
}

// IsClientError returns true when this delete otel metrics ingestion bad request response has a 4xx status code
func (o *DeleteOtelMetricsIngestionBadRequest) IsClientError() bool {
	return true
}

// IsServerError returns true when this delete otel metrics ingestion bad request response has a 5xx status code
func (o *DeleteOtelMetricsIngestionBadRequest) IsServerError() bool {
	return false
}

// IsCode returns true when this delete otel metrics ingestion bad request response a status code equal to that given
func (o *DeleteOtelMetricsIngestionBadRequest) IsCode(code int) bool {
	return code == 400
}

// Code gets the status code for the delete otel metrics ingestion bad request response
func (o *DeleteOtelMetricsIngestionBadRequest) Code() int {
	return 400
}

func (o *DeleteOtelMetricsIngestionBadRequest) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[DELETE /api/v1/config/otel-metrics-ingestion][%d] deleteOtelMetricsIngestionBadRequest %s", 400, payload)
}

func (o *DeleteOtelMetricsIngestionBadRequest) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[DELETE /api/v1/config/otel-metrics-ingestion][%d] deleteOtelMetricsIngestionBadRequest %s", 400, payload)
}

func (o *DeleteOtelMetricsIngestionBadRequest) GetPayload() *models.APIError {
	return o.Payload
}

func (o *DeleteOtelMetricsIngestionBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.APIError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDeleteOtelMetricsIngestionNotFound creates a DeleteOtelMetricsIngestionNotFound with default headers values
func NewDeleteOtelMetricsIngestionNotFound() *DeleteOtelMetricsIngestionNotFound {
	return &DeleteOtelMetricsIngestionNotFound{}
}

/*
DeleteOtelMetricsIngestionNotFound describes a response with status code 404, with default header values.

Cannot delete the OtelMetricsIngestion because the slug does not exist.
*/
type DeleteOtelMetricsIngestionNotFound struct {
	Payload *models.APIError
}

// IsSuccess returns true when this delete otel metrics ingestion not found response has a 2xx status code
func (o *DeleteOtelMetricsIngestionNotFound) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this delete otel metrics ingestion not found response has a 3xx status code
func (o *DeleteOtelMetricsIngestionNotFound) IsRedirect() bool {
	return false
}

// IsClientError returns true when this delete otel metrics ingestion not found response has a 4xx status code
func (o *DeleteOtelMetricsIngestionNotFound) IsClientError() bool {
	return true
}

// IsServerError returns true when this delete otel metrics ingestion not found response has a 5xx status code
func (o *DeleteOtelMetricsIngestionNotFound) IsServerError() bool {
	return false
}

// IsCode returns true when this delete otel metrics ingestion not found response a status code equal to that given
func (o *DeleteOtelMetricsIngestionNotFound) IsCode(code int) bool {
	return code == 404
}

// Code gets the status code for the delete otel metrics ingestion not found response
func (o *DeleteOtelMetricsIngestionNotFound) Code() int {
	return 404
}

func (o *DeleteOtelMetricsIngestionNotFound) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[DELETE /api/v1/config/otel-metrics-ingestion][%d] deleteOtelMetricsIngestionNotFound %s", 404, payload)
}

func (o *DeleteOtelMetricsIngestionNotFound) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[DELETE /api/v1/config/otel-metrics-ingestion][%d] deleteOtelMetricsIngestionNotFound %s", 404, payload)
}

func (o *DeleteOtelMetricsIngestionNotFound) GetPayload() *models.APIError {
	return o.Payload
}

func (o *DeleteOtelMetricsIngestionNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.APIError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDeleteOtelMetricsIngestionInternalServerError creates a DeleteOtelMetricsIngestionInternalServerError with default headers values
func NewDeleteOtelMetricsIngestionInternalServerError() *DeleteOtelMetricsIngestionInternalServerError {
	return &DeleteOtelMetricsIngestionInternalServerError{}
}

/*
DeleteOtelMetricsIngestionInternalServerError describes a response with status code 500, with default header values.

An unexpected error response.
*/
type DeleteOtelMetricsIngestionInternalServerError struct {
	Payload *models.APIError
}

// IsSuccess returns true when this delete otel metrics ingestion internal server error response has a 2xx status code
func (o *DeleteOtelMetricsIngestionInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this delete otel metrics ingestion internal server error response has a 3xx status code
func (o *DeleteOtelMetricsIngestionInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this delete otel metrics ingestion internal server error response has a 4xx status code
func (o *DeleteOtelMetricsIngestionInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this delete otel metrics ingestion internal server error response has a 5xx status code
func (o *DeleteOtelMetricsIngestionInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this delete otel metrics ingestion internal server error response a status code equal to that given
func (o *DeleteOtelMetricsIngestionInternalServerError) IsCode(code int) bool {
	return code == 500
}

// Code gets the status code for the delete otel metrics ingestion internal server error response
func (o *DeleteOtelMetricsIngestionInternalServerError) Code() int {
	return 500
}

func (o *DeleteOtelMetricsIngestionInternalServerError) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[DELETE /api/v1/config/otel-metrics-ingestion][%d] deleteOtelMetricsIngestionInternalServerError %s", 500, payload)
}

func (o *DeleteOtelMetricsIngestionInternalServerError) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[DELETE /api/v1/config/otel-metrics-ingestion][%d] deleteOtelMetricsIngestionInternalServerError %s", 500, payload)
}

func (o *DeleteOtelMetricsIngestionInternalServerError) GetPayload() *models.APIError {
	return o.Payload
}

func (o *DeleteOtelMetricsIngestionInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.APIError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDeleteOtelMetricsIngestionDefault creates a DeleteOtelMetricsIngestionDefault with default headers values
func NewDeleteOtelMetricsIngestionDefault(code int) *DeleteOtelMetricsIngestionDefault {
	return &DeleteOtelMetricsIngestionDefault{
		_statusCode: code,
	}
}

/*
DeleteOtelMetricsIngestionDefault describes a response with status code -1, with default header values.

An undefined error response.
*/
type DeleteOtelMetricsIngestionDefault struct {
	_statusCode int

	Payload models.GenericError
}

// IsSuccess returns true when this delete otel metrics ingestion default response has a 2xx status code
func (o *DeleteOtelMetricsIngestionDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this delete otel metrics ingestion default response has a 3xx status code
func (o *DeleteOtelMetricsIngestionDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this delete otel metrics ingestion default response has a 4xx status code
func (o *DeleteOtelMetricsIngestionDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this delete otel metrics ingestion default response has a 5xx status code
func (o *DeleteOtelMetricsIngestionDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this delete otel metrics ingestion default response a status code equal to that given
func (o *DeleteOtelMetricsIngestionDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the delete otel metrics ingestion default response
func (o *DeleteOtelMetricsIngestionDefault) Code() int {
	return o._statusCode
}

func (o *DeleteOtelMetricsIngestionDefault) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[DELETE /api/v1/config/otel-metrics-ingestion][%d] DeleteOtelMetricsIngestion default %s", o._statusCode, payload)
}

func (o *DeleteOtelMetricsIngestionDefault) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[DELETE /api/v1/config/otel-metrics-ingestion][%d] DeleteOtelMetricsIngestion default %s", o._statusCode, payload)
}

func (o *DeleteOtelMetricsIngestionDefault) GetPayload() models.GenericError {
	return o.Payload
}

func (o *DeleteOtelMetricsIngestionDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
