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

// CreateSLOReader is a Reader for the CreateSLO structure.
type CreateSLOReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *CreateSLOReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewCreateSLOOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewCreateSLOBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 409:
		result := NewCreateSLOConflict()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewCreateSLOInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewCreateSLODefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewCreateSLOOK creates a CreateSLOOK with default headers values
func NewCreateSLOOK() *CreateSLOOK {
	return &CreateSLOOK{}
}

/*
CreateSLOOK describes a response with status code 200, with default header values.

A successful response containing the created SLO.
*/
type CreateSLOOK struct {
	Payload *models.Configv1CreateSLOResponse
}

// IsSuccess returns true when this create s l o o k response has a 2xx status code
func (o *CreateSLOOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this create s l o o k response has a 3xx status code
func (o *CreateSLOOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this create s l o o k response has a 4xx status code
func (o *CreateSLOOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this create s l o o k response has a 5xx status code
func (o *CreateSLOOK) IsServerError() bool {
	return false
}

// IsCode returns true when this create s l o o k response a status code equal to that given
func (o *CreateSLOOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the create s l o o k response
func (o *CreateSLOOK) Code() int {
	return 200
}

func (o *CreateSLOOK) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[POST /api/v1/config/slos][%d] createSLOOK %s", 200, payload)
}

func (o *CreateSLOOK) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[POST /api/v1/config/slos][%d] createSLOOK %s", 200, payload)
}

func (o *CreateSLOOK) GetPayload() *models.Configv1CreateSLOResponse {
	return o.Payload
}

func (o *CreateSLOOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Configv1CreateSLOResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewCreateSLOBadRequest creates a CreateSLOBadRequest with default headers values
func NewCreateSLOBadRequest() *CreateSLOBadRequest {
	return &CreateSLOBadRequest{}
}

/*
CreateSLOBadRequest describes a response with status code 400, with default header values.

Cannot create the SLO because the request is invalid.
*/
type CreateSLOBadRequest struct {
	Payload *models.APIError
}

// IsSuccess returns true when this create s l o bad request response has a 2xx status code
func (o *CreateSLOBadRequest) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this create s l o bad request response has a 3xx status code
func (o *CreateSLOBadRequest) IsRedirect() bool {
	return false
}

// IsClientError returns true when this create s l o bad request response has a 4xx status code
func (o *CreateSLOBadRequest) IsClientError() bool {
	return true
}

// IsServerError returns true when this create s l o bad request response has a 5xx status code
func (o *CreateSLOBadRequest) IsServerError() bool {
	return false
}

// IsCode returns true when this create s l o bad request response a status code equal to that given
func (o *CreateSLOBadRequest) IsCode(code int) bool {
	return code == 400
}

// Code gets the status code for the create s l o bad request response
func (o *CreateSLOBadRequest) Code() int {
	return 400
}

func (o *CreateSLOBadRequest) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[POST /api/v1/config/slos][%d] createSLOBadRequest %s", 400, payload)
}

func (o *CreateSLOBadRequest) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[POST /api/v1/config/slos][%d] createSLOBadRequest %s", 400, payload)
}

func (o *CreateSLOBadRequest) GetPayload() *models.APIError {
	return o.Payload
}

func (o *CreateSLOBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.APIError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewCreateSLOConflict creates a CreateSLOConflict with default headers values
func NewCreateSLOConflict() *CreateSLOConflict {
	return &CreateSLOConflict{}
}

/*
CreateSLOConflict describes a response with status code 409, with default header values.

Cannot create the SLO because there is a conflict with an existing SLO.
*/
type CreateSLOConflict struct {
	Payload *models.APIError
}

// IsSuccess returns true when this create s l o conflict response has a 2xx status code
func (o *CreateSLOConflict) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this create s l o conflict response has a 3xx status code
func (o *CreateSLOConflict) IsRedirect() bool {
	return false
}

// IsClientError returns true when this create s l o conflict response has a 4xx status code
func (o *CreateSLOConflict) IsClientError() bool {
	return true
}

// IsServerError returns true when this create s l o conflict response has a 5xx status code
func (o *CreateSLOConflict) IsServerError() bool {
	return false
}

// IsCode returns true when this create s l o conflict response a status code equal to that given
func (o *CreateSLOConflict) IsCode(code int) bool {
	return code == 409
}

// Code gets the status code for the create s l o conflict response
func (o *CreateSLOConflict) Code() int {
	return 409
}

func (o *CreateSLOConflict) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[POST /api/v1/config/slos][%d] createSLOConflict %s", 409, payload)
}

func (o *CreateSLOConflict) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[POST /api/v1/config/slos][%d] createSLOConflict %s", 409, payload)
}

func (o *CreateSLOConflict) GetPayload() *models.APIError {
	return o.Payload
}

func (o *CreateSLOConflict) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.APIError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewCreateSLOInternalServerError creates a CreateSLOInternalServerError with default headers values
func NewCreateSLOInternalServerError() *CreateSLOInternalServerError {
	return &CreateSLOInternalServerError{}
}

/*
CreateSLOInternalServerError describes a response with status code 500, with default header values.

An unexpected error response.
*/
type CreateSLOInternalServerError struct {
	Payload *models.APIError
}

// IsSuccess returns true when this create s l o internal server error response has a 2xx status code
func (o *CreateSLOInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this create s l o internal server error response has a 3xx status code
func (o *CreateSLOInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this create s l o internal server error response has a 4xx status code
func (o *CreateSLOInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this create s l o internal server error response has a 5xx status code
func (o *CreateSLOInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this create s l o internal server error response a status code equal to that given
func (o *CreateSLOInternalServerError) IsCode(code int) bool {
	return code == 500
}

// Code gets the status code for the create s l o internal server error response
func (o *CreateSLOInternalServerError) Code() int {
	return 500
}

func (o *CreateSLOInternalServerError) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[POST /api/v1/config/slos][%d] createSLOInternalServerError %s", 500, payload)
}

func (o *CreateSLOInternalServerError) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[POST /api/v1/config/slos][%d] createSLOInternalServerError %s", 500, payload)
}

func (o *CreateSLOInternalServerError) GetPayload() *models.APIError {
	return o.Payload
}

func (o *CreateSLOInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.APIError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewCreateSLODefault creates a CreateSLODefault with default headers values
func NewCreateSLODefault(code int) *CreateSLODefault {
	return &CreateSLODefault{
		_statusCode: code,
	}
}

/*
CreateSLODefault describes a response with status code -1, with default header values.

An undefined error response.
*/
type CreateSLODefault struct {
	_statusCode int

	Payload models.GenericError
}

// IsSuccess returns true when this create s l o default response has a 2xx status code
func (o *CreateSLODefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this create s l o default response has a 3xx status code
func (o *CreateSLODefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this create s l o default response has a 4xx status code
func (o *CreateSLODefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this create s l o default response has a 5xx status code
func (o *CreateSLODefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this create s l o default response a status code equal to that given
func (o *CreateSLODefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the create s l o default response
func (o *CreateSLODefault) Code() int {
	return o._statusCode
}

func (o *CreateSLODefault) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[POST /api/v1/config/slos][%d] CreateSLO default %s", o._statusCode, payload)
}

func (o *CreateSLODefault) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[POST /api/v1/config/slos][%d] CreateSLO default %s", o._statusCode, payload)
}

func (o *CreateSLODefault) GetPayload() models.GenericError {
	return o.Payload
}

func (o *CreateSLODefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
