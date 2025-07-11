// Code generated by go-swagger; DO NOT EDIT.

package trace_tail_sampling_rules

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

// UpdateTraceTailSamplingRulesReader is a Reader for the UpdateTraceTailSamplingRules structure.
type UpdateTraceTailSamplingRulesReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *UpdateTraceTailSamplingRulesReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewUpdateTraceTailSamplingRulesOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewUpdateTraceTailSamplingRulesBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewUpdateTraceTailSamplingRulesNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewUpdateTraceTailSamplingRulesInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewUpdateTraceTailSamplingRulesDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewUpdateTraceTailSamplingRulesOK creates a UpdateTraceTailSamplingRulesOK with default headers values
func NewUpdateTraceTailSamplingRulesOK() *UpdateTraceTailSamplingRulesOK {
	return &UpdateTraceTailSamplingRulesOK{}
}

/*
UpdateTraceTailSamplingRulesOK describes a response with status code 200, with default header values.

A successful response containing the updated TraceTailSamplingRules.
*/
type UpdateTraceTailSamplingRulesOK struct {
	Payload *models.Configv1UpdateTraceTailSamplingRulesResponse
}

// IsSuccess returns true when this update trace tail sampling rules o k response has a 2xx status code
func (o *UpdateTraceTailSamplingRulesOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this update trace tail sampling rules o k response has a 3xx status code
func (o *UpdateTraceTailSamplingRulesOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this update trace tail sampling rules o k response has a 4xx status code
func (o *UpdateTraceTailSamplingRulesOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this update trace tail sampling rules o k response has a 5xx status code
func (o *UpdateTraceTailSamplingRulesOK) IsServerError() bool {
	return false
}

// IsCode returns true when this update trace tail sampling rules o k response a status code equal to that given
func (o *UpdateTraceTailSamplingRulesOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the update trace tail sampling rules o k response
func (o *UpdateTraceTailSamplingRulesOK) Code() int {
	return 200
}

func (o *UpdateTraceTailSamplingRulesOK) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[PUT /api/v1/config/trace-tail-sampling-rules][%d] updateTraceTailSamplingRulesOK %s", 200, payload)
}

func (o *UpdateTraceTailSamplingRulesOK) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[PUT /api/v1/config/trace-tail-sampling-rules][%d] updateTraceTailSamplingRulesOK %s", 200, payload)
}

func (o *UpdateTraceTailSamplingRulesOK) GetPayload() *models.Configv1UpdateTraceTailSamplingRulesResponse {
	return o.Payload
}

func (o *UpdateTraceTailSamplingRulesOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Configv1UpdateTraceTailSamplingRulesResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateTraceTailSamplingRulesBadRequest creates a UpdateTraceTailSamplingRulesBadRequest with default headers values
func NewUpdateTraceTailSamplingRulesBadRequest() *UpdateTraceTailSamplingRulesBadRequest {
	return &UpdateTraceTailSamplingRulesBadRequest{}
}

/*
UpdateTraceTailSamplingRulesBadRequest describes a response with status code 400, with default header values.

Cannot update the TraceTailSamplingRules because the request is invalid.
*/
type UpdateTraceTailSamplingRulesBadRequest struct {
	Payload *models.APIError
}

// IsSuccess returns true when this update trace tail sampling rules bad request response has a 2xx status code
func (o *UpdateTraceTailSamplingRulesBadRequest) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this update trace tail sampling rules bad request response has a 3xx status code
func (o *UpdateTraceTailSamplingRulesBadRequest) IsRedirect() bool {
	return false
}

// IsClientError returns true when this update trace tail sampling rules bad request response has a 4xx status code
func (o *UpdateTraceTailSamplingRulesBadRequest) IsClientError() bool {
	return true
}

// IsServerError returns true when this update trace tail sampling rules bad request response has a 5xx status code
func (o *UpdateTraceTailSamplingRulesBadRequest) IsServerError() bool {
	return false
}

// IsCode returns true when this update trace tail sampling rules bad request response a status code equal to that given
func (o *UpdateTraceTailSamplingRulesBadRequest) IsCode(code int) bool {
	return code == 400
}

// Code gets the status code for the update trace tail sampling rules bad request response
func (o *UpdateTraceTailSamplingRulesBadRequest) Code() int {
	return 400
}

func (o *UpdateTraceTailSamplingRulesBadRequest) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[PUT /api/v1/config/trace-tail-sampling-rules][%d] updateTraceTailSamplingRulesBadRequest %s", 400, payload)
}

func (o *UpdateTraceTailSamplingRulesBadRequest) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[PUT /api/v1/config/trace-tail-sampling-rules][%d] updateTraceTailSamplingRulesBadRequest %s", 400, payload)
}

func (o *UpdateTraceTailSamplingRulesBadRequest) GetPayload() *models.APIError {
	return o.Payload
}

func (o *UpdateTraceTailSamplingRulesBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.APIError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateTraceTailSamplingRulesNotFound creates a UpdateTraceTailSamplingRulesNotFound with default headers values
func NewUpdateTraceTailSamplingRulesNotFound() *UpdateTraceTailSamplingRulesNotFound {
	return &UpdateTraceTailSamplingRulesNotFound{}
}

/*
UpdateTraceTailSamplingRulesNotFound describes a response with status code 404, with default header values.

Cannot update the TraceTailSamplingRules because TraceTailSamplingRules has not been created.
*/
type UpdateTraceTailSamplingRulesNotFound struct {
	Payload *models.APIError
}

// IsSuccess returns true when this update trace tail sampling rules not found response has a 2xx status code
func (o *UpdateTraceTailSamplingRulesNotFound) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this update trace tail sampling rules not found response has a 3xx status code
func (o *UpdateTraceTailSamplingRulesNotFound) IsRedirect() bool {
	return false
}

// IsClientError returns true when this update trace tail sampling rules not found response has a 4xx status code
func (o *UpdateTraceTailSamplingRulesNotFound) IsClientError() bool {
	return true
}

// IsServerError returns true when this update trace tail sampling rules not found response has a 5xx status code
func (o *UpdateTraceTailSamplingRulesNotFound) IsServerError() bool {
	return false
}

// IsCode returns true when this update trace tail sampling rules not found response a status code equal to that given
func (o *UpdateTraceTailSamplingRulesNotFound) IsCode(code int) bool {
	return code == 404
}

// Code gets the status code for the update trace tail sampling rules not found response
func (o *UpdateTraceTailSamplingRulesNotFound) Code() int {
	return 404
}

func (o *UpdateTraceTailSamplingRulesNotFound) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[PUT /api/v1/config/trace-tail-sampling-rules][%d] updateTraceTailSamplingRulesNotFound %s", 404, payload)
}

func (o *UpdateTraceTailSamplingRulesNotFound) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[PUT /api/v1/config/trace-tail-sampling-rules][%d] updateTraceTailSamplingRulesNotFound %s", 404, payload)
}

func (o *UpdateTraceTailSamplingRulesNotFound) GetPayload() *models.APIError {
	return o.Payload
}

func (o *UpdateTraceTailSamplingRulesNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.APIError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateTraceTailSamplingRulesInternalServerError creates a UpdateTraceTailSamplingRulesInternalServerError with default headers values
func NewUpdateTraceTailSamplingRulesInternalServerError() *UpdateTraceTailSamplingRulesInternalServerError {
	return &UpdateTraceTailSamplingRulesInternalServerError{}
}

/*
UpdateTraceTailSamplingRulesInternalServerError describes a response with status code 500, with default header values.

An unexpected error response.
*/
type UpdateTraceTailSamplingRulesInternalServerError struct {
	Payload *models.APIError
}

// IsSuccess returns true when this update trace tail sampling rules internal server error response has a 2xx status code
func (o *UpdateTraceTailSamplingRulesInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this update trace tail sampling rules internal server error response has a 3xx status code
func (o *UpdateTraceTailSamplingRulesInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this update trace tail sampling rules internal server error response has a 4xx status code
func (o *UpdateTraceTailSamplingRulesInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this update trace tail sampling rules internal server error response has a 5xx status code
func (o *UpdateTraceTailSamplingRulesInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this update trace tail sampling rules internal server error response a status code equal to that given
func (o *UpdateTraceTailSamplingRulesInternalServerError) IsCode(code int) bool {
	return code == 500
}

// Code gets the status code for the update trace tail sampling rules internal server error response
func (o *UpdateTraceTailSamplingRulesInternalServerError) Code() int {
	return 500
}

func (o *UpdateTraceTailSamplingRulesInternalServerError) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[PUT /api/v1/config/trace-tail-sampling-rules][%d] updateTraceTailSamplingRulesInternalServerError %s", 500, payload)
}

func (o *UpdateTraceTailSamplingRulesInternalServerError) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[PUT /api/v1/config/trace-tail-sampling-rules][%d] updateTraceTailSamplingRulesInternalServerError %s", 500, payload)
}

func (o *UpdateTraceTailSamplingRulesInternalServerError) GetPayload() *models.APIError {
	return o.Payload
}

func (o *UpdateTraceTailSamplingRulesInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.APIError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateTraceTailSamplingRulesDefault creates a UpdateTraceTailSamplingRulesDefault with default headers values
func NewUpdateTraceTailSamplingRulesDefault(code int) *UpdateTraceTailSamplingRulesDefault {
	return &UpdateTraceTailSamplingRulesDefault{
		_statusCode: code,
	}
}

/*
UpdateTraceTailSamplingRulesDefault describes a response with status code -1, with default header values.

An undefined error response.
*/
type UpdateTraceTailSamplingRulesDefault struct {
	_statusCode int

	Payload models.GenericError
}

// IsSuccess returns true when this update trace tail sampling rules default response has a 2xx status code
func (o *UpdateTraceTailSamplingRulesDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this update trace tail sampling rules default response has a 3xx status code
func (o *UpdateTraceTailSamplingRulesDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this update trace tail sampling rules default response has a 4xx status code
func (o *UpdateTraceTailSamplingRulesDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this update trace tail sampling rules default response has a 5xx status code
func (o *UpdateTraceTailSamplingRulesDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this update trace tail sampling rules default response a status code equal to that given
func (o *UpdateTraceTailSamplingRulesDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the update trace tail sampling rules default response
func (o *UpdateTraceTailSamplingRulesDefault) Code() int {
	return o._statusCode
}

func (o *UpdateTraceTailSamplingRulesDefault) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[PUT /api/v1/config/trace-tail-sampling-rules][%d] UpdateTraceTailSamplingRules default %s", o._statusCode, payload)
}

func (o *UpdateTraceTailSamplingRulesDefault) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[PUT /api/v1/config/trace-tail-sampling-rules][%d] UpdateTraceTailSamplingRules default %s", o._statusCode, payload)
}

func (o *UpdateTraceTailSamplingRulesDefault) GetPayload() models.GenericError {
	return o.Payload
}

func (o *UpdateTraceTailSamplingRulesDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
