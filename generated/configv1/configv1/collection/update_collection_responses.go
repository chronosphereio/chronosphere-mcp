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

// UpdateCollectionReader is a Reader for the UpdateCollection structure.
type UpdateCollectionReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *UpdateCollectionReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewUpdateCollectionOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewUpdateCollectionBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewUpdateCollectionNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 409:
		result := NewUpdateCollectionConflict()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewUpdateCollectionInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewUpdateCollectionDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewUpdateCollectionOK creates a UpdateCollectionOK with default headers values
func NewUpdateCollectionOK() *UpdateCollectionOK {
	return &UpdateCollectionOK{}
}

/*
UpdateCollectionOK describes a response with status code 200, with default header values.

A successful response containing the updated Collection.
*/
type UpdateCollectionOK struct {
	Payload *models.Configv1UpdateCollectionResponse
}

// IsSuccess returns true when this update collection o k response has a 2xx status code
func (o *UpdateCollectionOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this update collection o k response has a 3xx status code
func (o *UpdateCollectionOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this update collection o k response has a 4xx status code
func (o *UpdateCollectionOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this update collection o k response has a 5xx status code
func (o *UpdateCollectionOK) IsServerError() bool {
	return false
}

// IsCode returns true when this update collection o k response a status code equal to that given
func (o *UpdateCollectionOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the update collection o k response
func (o *UpdateCollectionOK) Code() int {
	return 200
}

func (o *UpdateCollectionOK) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[PUT /api/v1/config/collections/{slug}][%d] updateCollectionOK %s", 200, payload)
}

func (o *UpdateCollectionOK) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[PUT /api/v1/config/collections/{slug}][%d] updateCollectionOK %s", 200, payload)
}

func (o *UpdateCollectionOK) GetPayload() *models.Configv1UpdateCollectionResponse {
	return o.Payload
}

func (o *UpdateCollectionOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Configv1UpdateCollectionResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateCollectionBadRequest creates a UpdateCollectionBadRequest with default headers values
func NewUpdateCollectionBadRequest() *UpdateCollectionBadRequest {
	return &UpdateCollectionBadRequest{}
}

/*
UpdateCollectionBadRequest describes a response with status code 400, with default header values.

Cannot update the Collection because the request is invalid.
*/
type UpdateCollectionBadRequest struct {
	Payload *models.APIError
}

// IsSuccess returns true when this update collection bad request response has a 2xx status code
func (o *UpdateCollectionBadRequest) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this update collection bad request response has a 3xx status code
func (o *UpdateCollectionBadRequest) IsRedirect() bool {
	return false
}

// IsClientError returns true when this update collection bad request response has a 4xx status code
func (o *UpdateCollectionBadRequest) IsClientError() bool {
	return true
}

// IsServerError returns true when this update collection bad request response has a 5xx status code
func (o *UpdateCollectionBadRequest) IsServerError() bool {
	return false
}

// IsCode returns true when this update collection bad request response a status code equal to that given
func (o *UpdateCollectionBadRequest) IsCode(code int) bool {
	return code == 400
}

// Code gets the status code for the update collection bad request response
func (o *UpdateCollectionBadRequest) Code() int {
	return 400
}

func (o *UpdateCollectionBadRequest) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[PUT /api/v1/config/collections/{slug}][%d] updateCollectionBadRequest %s", 400, payload)
}

func (o *UpdateCollectionBadRequest) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[PUT /api/v1/config/collections/{slug}][%d] updateCollectionBadRequest %s", 400, payload)
}

func (o *UpdateCollectionBadRequest) GetPayload() *models.APIError {
	return o.Payload
}

func (o *UpdateCollectionBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.APIError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateCollectionNotFound creates a UpdateCollectionNotFound with default headers values
func NewUpdateCollectionNotFound() *UpdateCollectionNotFound {
	return &UpdateCollectionNotFound{}
}

/*
UpdateCollectionNotFound describes a response with status code 404, with default header values.

Cannot update the Collection because the slug does not exist.
*/
type UpdateCollectionNotFound struct {
	Payload *models.APIError
}

// IsSuccess returns true when this update collection not found response has a 2xx status code
func (o *UpdateCollectionNotFound) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this update collection not found response has a 3xx status code
func (o *UpdateCollectionNotFound) IsRedirect() bool {
	return false
}

// IsClientError returns true when this update collection not found response has a 4xx status code
func (o *UpdateCollectionNotFound) IsClientError() bool {
	return true
}

// IsServerError returns true when this update collection not found response has a 5xx status code
func (o *UpdateCollectionNotFound) IsServerError() bool {
	return false
}

// IsCode returns true when this update collection not found response a status code equal to that given
func (o *UpdateCollectionNotFound) IsCode(code int) bool {
	return code == 404
}

// Code gets the status code for the update collection not found response
func (o *UpdateCollectionNotFound) Code() int {
	return 404
}

func (o *UpdateCollectionNotFound) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[PUT /api/v1/config/collections/{slug}][%d] updateCollectionNotFound %s", 404, payload)
}

func (o *UpdateCollectionNotFound) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[PUT /api/v1/config/collections/{slug}][%d] updateCollectionNotFound %s", 404, payload)
}

func (o *UpdateCollectionNotFound) GetPayload() *models.APIError {
	return o.Payload
}

func (o *UpdateCollectionNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.APIError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateCollectionConflict creates a UpdateCollectionConflict with default headers values
func NewUpdateCollectionConflict() *UpdateCollectionConflict {
	return &UpdateCollectionConflict{}
}

/*
UpdateCollectionConflict describes a response with status code 409, with default header values.

Cannot update the Collection because there is a conflict with an existing Collection.
*/
type UpdateCollectionConflict struct {
	Payload *models.APIError
}

// IsSuccess returns true when this update collection conflict response has a 2xx status code
func (o *UpdateCollectionConflict) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this update collection conflict response has a 3xx status code
func (o *UpdateCollectionConflict) IsRedirect() bool {
	return false
}

// IsClientError returns true when this update collection conflict response has a 4xx status code
func (o *UpdateCollectionConflict) IsClientError() bool {
	return true
}

// IsServerError returns true when this update collection conflict response has a 5xx status code
func (o *UpdateCollectionConflict) IsServerError() bool {
	return false
}

// IsCode returns true when this update collection conflict response a status code equal to that given
func (o *UpdateCollectionConflict) IsCode(code int) bool {
	return code == 409
}

// Code gets the status code for the update collection conflict response
func (o *UpdateCollectionConflict) Code() int {
	return 409
}

func (o *UpdateCollectionConflict) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[PUT /api/v1/config/collections/{slug}][%d] updateCollectionConflict %s", 409, payload)
}

func (o *UpdateCollectionConflict) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[PUT /api/v1/config/collections/{slug}][%d] updateCollectionConflict %s", 409, payload)
}

func (o *UpdateCollectionConflict) GetPayload() *models.APIError {
	return o.Payload
}

func (o *UpdateCollectionConflict) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.APIError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateCollectionInternalServerError creates a UpdateCollectionInternalServerError with default headers values
func NewUpdateCollectionInternalServerError() *UpdateCollectionInternalServerError {
	return &UpdateCollectionInternalServerError{}
}

/*
UpdateCollectionInternalServerError describes a response with status code 500, with default header values.

An unexpected error response.
*/
type UpdateCollectionInternalServerError struct {
	Payload *models.APIError
}

// IsSuccess returns true when this update collection internal server error response has a 2xx status code
func (o *UpdateCollectionInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this update collection internal server error response has a 3xx status code
func (o *UpdateCollectionInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this update collection internal server error response has a 4xx status code
func (o *UpdateCollectionInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this update collection internal server error response has a 5xx status code
func (o *UpdateCollectionInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this update collection internal server error response a status code equal to that given
func (o *UpdateCollectionInternalServerError) IsCode(code int) bool {
	return code == 500
}

// Code gets the status code for the update collection internal server error response
func (o *UpdateCollectionInternalServerError) Code() int {
	return 500
}

func (o *UpdateCollectionInternalServerError) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[PUT /api/v1/config/collections/{slug}][%d] updateCollectionInternalServerError %s", 500, payload)
}

func (o *UpdateCollectionInternalServerError) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[PUT /api/v1/config/collections/{slug}][%d] updateCollectionInternalServerError %s", 500, payload)
}

func (o *UpdateCollectionInternalServerError) GetPayload() *models.APIError {
	return o.Payload
}

func (o *UpdateCollectionInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.APIError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateCollectionDefault creates a UpdateCollectionDefault with default headers values
func NewUpdateCollectionDefault(code int) *UpdateCollectionDefault {
	return &UpdateCollectionDefault{
		_statusCode: code,
	}
}

/*
UpdateCollectionDefault describes a response with status code -1, with default header values.

An undefined error response.
*/
type UpdateCollectionDefault struct {
	_statusCode int

	Payload models.GenericError
}

// IsSuccess returns true when this update collection default response has a 2xx status code
func (o *UpdateCollectionDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this update collection default response has a 3xx status code
func (o *UpdateCollectionDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this update collection default response has a 4xx status code
func (o *UpdateCollectionDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this update collection default response has a 5xx status code
func (o *UpdateCollectionDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this update collection default response a status code equal to that given
func (o *UpdateCollectionDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the update collection default response
func (o *UpdateCollectionDefault) Code() int {
	return o._statusCode
}

func (o *UpdateCollectionDefault) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[PUT /api/v1/config/collections/{slug}][%d] UpdateCollection default %s", o._statusCode, payload)
}

func (o *UpdateCollectionDefault) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[PUT /api/v1/config/collections/{slug}][%d] UpdateCollection default %s", o._statusCode, payload)
}

func (o *UpdateCollectionDefault) GetPayload() models.GenericError {
	return o.Payload
}

func (o *UpdateCollectionDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
