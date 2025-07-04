// Code generated by go-swagger; DO NOT EDIT.

package collection

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

// DeleteCollectionReader is a Reader for the DeleteCollection structure.
type DeleteCollectionReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *DeleteCollectionReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewDeleteCollectionOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewDeleteCollectionBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewDeleteCollectionNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewDeleteCollectionInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewDeleteCollectionDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewDeleteCollectionOK creates a DeleteCollectionOK with default headers values
func NewDeleteCollectionOK() *DeleteCollectionOK {
	return &DeleteCollectionOK{}
}

/*
DeleteCollectionOK describes a response with status code 200, with default header values.

A successful response.
*/
type DeleteCollectionOK struct {
	Payload models.Configv1DeleteCollectionResponse
}

// IsSuccess returns true when this delete collection o k response has a 2xx status code
func (o *DeleteCollectionOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this delete collection o k response has a 3xx status code
func (o *DeleteCollectionOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this delete collection o k response has a 4xx status code
func (o *DeleteCollectionOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this delete collection o k response has a 5xx status code
func (o *DeleteCollectionOK) IsServerError() bool {
	return false
}

// IsCode returns true when this delete collection o k response a status code equal to that given
func (o *DeleteCollectionOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the delete collection o k response
func (o *DeleteCollectionOK) Code() int {
	return 200
}

func (o *DeleteCollectionOK) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[DELETE /api/v1/config/collections/{slug}][%d] deleteCollectionOK %s", 200, payload)
}

func (o *DeleteCollectionOK) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[DELETE /api/v1/config/collections/{slug}][%d] deleteCollectionOK %s", 200, payload)
}

func (o *DeleteCollectionOK) GetPayload() models.Configv1DeleteCollectionResponse {
	return o.Payload
}

func (o *DeleteCollectionOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDeleteCollectionBadRequest creates a DeleteCollectionBadRequest with default headers values
func NewDeleteCollectionBadRequest() *DeleteCollectionBadRequest {
	return &DeleteCollectionBadRequest{}
}

/*
DeleteCollectionBadRequest describes a response with status code 400, with default header values.

Cannot delete the Collection because it is in use.
*/
type DeleteCollectionBadRequest struct {
	Payload *models.APIError
}

// IsSuccess returns true when this delete collection bad request response has a 2xx status code
func (o *DeleteCollectionBadRequest) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this delete collection bad request response has a 3xx status code
func (o *DeleteCollectionBadRequest) IsRedirect() bool {
	return false
}

// IsClientError returns true when this delete collection bad request response has a 4xx status code
func (o *DeleteCollectionBadRequest) IsClientError() bool {
	return true
}

// IsServerError returns true when this delete collection bad request response has a 5xx status code
func (o *DeleteCollectionBadRequest) IsServerError() bool {
	return false
}

// IsCode returns true when this delete collection bad request response a status code equal to that given
func (o *DeleteCollectionBadRequest) IsCode(code int) bool {
	return code == 400
}

// Code gets the status code for the delete collection bad request response
func (o *DeleteCollectionBadRequest) Code() int {
	return 400
}

func (o *DeleteCollectionBadRequest) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[DELETE /api/v1/config/collections/{slug}][%d] deleteCollectionBadRequest %s", 400, payload)
}

func (o *DeleteCollectionBadRequest) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[DELETE /api/v1/config/collections/{slug}][%d] deleteCollectionBadRequest %s", 400, payload)
}

func (o *DeleteCollectionBadRequest) GetPayload() *models.APIError {
	return o.Payload
}

func (o *DeleteCollectionBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.APIError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDeleteCollectionNotFound creates a DeleteCollectionNotFound with default headers values
func NewDeleteCollectionNotFound() *DeleteCollectionNotFound {
	return &DeleteCollectionNotFound{}
}

/*
DeleteCollectionNotFound describes a response with status code 404, with default header values.

Cannot delete the Collection because the slug does not exist.
*/
type DeleteCollectionNotFound struct {
	Payload *models.APIError
}

// IsSuccess returns true when this delete collection not found response has a 2xx status code
func (o *DeleteCollectionNotFound) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this delete collection not found response has a 3xx status code
func (o *DeleteCollectionNotFound) IsRedirect() bool {
	return false
}

// IsClientError returns true when this delete collection not found response has a 4xx status code
func (o *DeleteCollectionNotFound) IsClientError() bool {
	return true
}

// IsServerError returns true when this delete collection not found response has a 5xx status code
func (o *DeleteCollectionNotFound) IsServerError() bool {
	return false
}

// IsCode returns true when this delete collection not found response a status code equal to that given
func (o *DeleteCollectionNotFound) IsCode(code int) bool {
	return code == 404
}

// Code gets the status code for the delete collection not found response
func (o *DeleteCollectionNotFound) Code() int {
	return 404
}

func (o *DeleteCollectionNotFound) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[DELETE /api/v1/config/collections/{slug}][%d] deleteCollectionNotFound %s", 404, payload)
}

func (o *DeleteCollectionNotFound) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[DELETE /api/v1/config/collections/{slug}][%d] deleteCollectionNotFound %s", 404, payload)
}

func (o *DeleteCollectionNotFound) GetPayload() *models.APIError {
	return o.Payload
}

func (o *DeleteCollectionNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.APIError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDeleteCollectionInternalServerError creates a DeleteCollectionInternalServerError with default headers values
func NewDeleteCollectionInternalServerError() *DeleteCollectionInternalServerError {
	return &DeleteCollectionInternalServerError{}
}

/*
DeleteCollectionInternalServerError describes a response with status code 500, with default header values.

An unexpected error response.
*/
type DeleteCollectionInternalServerError struct {
	Payload *models.APIError
}

// IsSuccess returns true when this delete collection internal server error response has a 2xx status code
func (o *DeleteCollectionInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this delete collection internal server error response has a 3xx status code
func (o *DeleteCollectionInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this delete collection internal server error response has a 4xx status code
func (o *DeleteCollectionInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this delete collection internal server error response has a 5xx status code
func (o *DeleteCollectionInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this delete collection internal server error response a status code equal to that given
func (o *DeleteCollectionInternalServerError) IsCode(code int) bool {
	return code == 500
}

// Code gets the status code for the delete collection internal server error response
func (o *DeleteCollectionInternalServerError) Code() int {
	return 500
}

func (o *DeleteCollectionInternalServerError) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[DELETE /api/v1/config/collections/{slug}][%d] deleteCollectionInternalServerError %s", 500, payload)
}

func (o *DeleteCollectionInternalServerError) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[DELETE /api/v1/config/collections/{slug}][%d] deleteCollectionInternalServerError %s", 500, payload)
}

func (o *DeleteCollectionInternalServerError) GetPayload() *models.APIError {
	return o.Payload
}

func (o *DeleteCollectionInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.APIError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDeleteCollectionDefault creates a DeleteCollectionDefault with default headers values
func NewDeleteCollectionDefault(code int) *DeleteCollectionDefault {
	return &DeleteCollectionDefault{
		_statusCode: code,
	}
}

/*
DeleteCollectionDefault describes a response with status code -1, with default header values.

An undefined error response.
*/
type DeleteCollectionDefault struct {
	_statusCode int

	Payload models.GenericError
}

// IsSuccess returns true when this delete collection default response has a 2xx status code
func (o *DeleteCollectionDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this delete collection default response has a 3xx status code
func (o *DeleteCollectionDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this delete collection default response has a 4xx status code
func (o *DeleteCollectionDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this delete collection default response has a 5xx status code
func (o *DeleteCollectionDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this delete collection default response a status code equal to that given
func (o *DeleteCollectionDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the delete collection default response
func (o *DeleteCollectionDefault) Code() int {
	return o._statusCode
}

func (o *DeleteCollectionDefault) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[DELETE /api/v1/config/collections/{slug}][%d] DeleteCollection default %s", o._statusCode, payload)
}

func (o *DeleteCollectionDefault) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[DELETE /api/v1/config/collections/{slug}][%d] DeleteCollection default %s", o._statusCode, payload)
}

func (o *DeleteCollectionDefault) GetPayload() models.GenericError {
	return o.Payload
}

func (o *DeleteCollectionDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
