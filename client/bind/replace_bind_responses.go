// Code generated by go-swagger; DO NOT EDIT.

package bind

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/haproxytech/models"
)

// ReplaceBindReader is a Reader for the ReplaceBind structure.
type ReplaceBindReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ReplaceBindReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewReplaceBindOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	case 202:
		result := NewReplaceBindAccepted()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	case 400:
		result := NewReplaceBindBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	case 404:
		result := NewReplaceBindNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		result := NewReplaceBindDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewReplaceBindOK creates a ReplaceBindOK with default headers values
func NewReplaceBindOK() *ReplaceBindOK {
	return &ReplaceBindOK{}
}

/*ReplaceBindOK handles this case with default header values.

Bind replaced
*/
type ReplaceBindOK struct {
	Payload *models.Bind
}

func (o *ReplaceBindOK) Error() string {
	return fmt.Sprintf("[PUT /services/haproxy/configuration/binds/{name}][%d] replaceBindOK  %+v", 200, o.Payload)
}

func (o *ReplaceBindOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Bind)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewReplaceBindAccepted creates a ReplaceBindAccepted with default headers values
func NewReplaceBindAccepted() *ReplaceBindAccepted {
	return &ReplaceBindAccepted{}
}

/*ReplaceBindAccepted handles this case with default header values.

Configuration change accepted and reload requested
*/
type ReplaceBindAccepted struct {
	/*ID of the requested reload
	 */
	ReloadID string

	Payload *models.Bind
}

func (o *ReplaceBindAccepted) Error() string {
	return fmt.Sprintf("[PUT /services/haproxy/configuration/binds/{name}][%d] replaceBindAccepted  %+v", 202, o.Payload)
}

func (o *ReplaceBindAccepted) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response header Reload-ID
	o.ReloadID = response.GetHeader("Reload-ID")

	o.Payload = new(models.Bind)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewReplaceBindBadRequest creates a ReplaceBindBadRequest with default headers values
func NewReplaceBindBadRequest() *ReplaceBindBadRequest {
	return &ReplaceBindBadRequest{}
}

/*ReplaceBindBadRequest handles this case with default header values.

Bad request
*/
type ReplaceBindBadRequest struct {
	Payload *models.Error
}

func (o *ReplaceBindBadRequest) Error() string {
	return fmt.Sprintf("[PUT /services/haproxy/configuration/binds/{name}][%d] replaceBindBadRequest  %+v", 400, o.Payload)
}

func (o *ReplaceBindBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewReplaceBindNotFound creates a ReplaceBindNotFound with default headers values
func NewReplaceBindNotFound() *ReplaceBindNotFound {
	return &ReplaceBindNotFound{}
}

/*ReplaceBindNotFound handles this case with default header values.

The specified resource was not found
*/
type ReplaceBindNotFound struct {
	Payload *models.Error
}

func (o *ReplaceBindNotFound) Error() string {
	return fmt.Sprintf("[PUT /services/haproxy/configuration/binds/{name}][%d] replaceBindNotFound  %+v", 404, o.Payload)
}

func (o *ReplaceBindNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewReplaceBindDefault creates a ReplaceBindDefault with default headers values
func NewReplaceBindDefault(code int) *ReplaceBindDefault {
	return &ReplaceBindDefault{
		_statusCode: code,
	}
}

/*ReplaceBindDefault handles this case with default header values.

General Error
*/
type ReplaceBindDefault struct {
	_statusCode int

	Payload *models.Error
}

// Code gets the status code for the replace bind default response
func (o *ReplaceBindDefault) Code() int {
	return o._statusCode
}

func (o *ReplaceBindDefault) Error() string {
	return fmt.Sprintf("[PUT /services/haproxy/configuration/binds/{name}][%d] replaceBind default  %+v", o._statusCode, o.Payload)
}

func (o *ReplaceBindDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
