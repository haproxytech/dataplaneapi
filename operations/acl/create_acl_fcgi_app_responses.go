// Code generated by go-swagger; DO NOT EDIT.

// Copyright 2019 HAProxy Technologies
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package acl

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v6/models"
)

// CreateACLFCGIAppCreatedCode is the HTTP code returned for type CreateACLFCGIAppCreated
const CreateACLFCGIAppCreatedCode int = 201

/*
CreateACLFCGIAppCreated ACL line created

swagger:response createAclFcgiAppCreated
*/
type CreateACLFCGIAppCreated struct {

	/*
	  In: Body
	*/
	Payload *models.ACL `json:"body,omitempty"`
}

// NewCreateACLFCGIAppCreated creates CreateACLFCGIAppCreated with default headers values
func NewCreateACLFCGIAppCreated() *CreateACLFCGIAppCreated {

	return &CreateACLFCGIAppCreated{}
}

// WithPayload adds the payload to the create Acl Fcgi app created response
func (o *CreateACLFCGIAppCreated) WithPayload(payload *models.ACL) *CreateACLFCGIAppCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create Acl Fcgi app created response
func (o *CreateACLFCGIAppCreated) SetPayload(payload *models.ACL) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateACLFCGIAppCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateACLFCGIAppAcceptedCode is the HTTP code returned for type CreateACLFCGIAppAccepted
const CreateACLFCGIAppAcceptedCode int = 202

/*
CreateACLFCGIAppAccepted Configuration change accepted and reload requested

swagger:response createAclFcgiAppAccepted
*/
type CreateACLFCGIAppAccepted struct {
	/*ID of the requested reload

	 */
	ReloadID string `json:"Reload-ID"`

	/*
	  In: Body
	*/
	Payload *models.ACL `json:"body,omitempty"`
}

// NewCreateACLFCGIAppAccepted creates CreateACLFCGIAppAccepted with default headers values
func NewCreateACLFCGIAppAccepted() *CreateACLFCGIAppAccepted {

	return &CreateACLFCGIAppAccepted{}
}

// WithReloadID adds the reloadId to the create Acl Fcgi app accepted response
func (o *CreateACLFCGIAppAccepted) WithReloadID(reloadID string) *CreateACLFCGIAppAccepted {
	o.ReloadID = reloadID
	return o
}

// SetReloadID sets the reloadId to the create Acl Fcgi app accepted response
func (o *CreateACLFCGIAppAccepted) SetReloadID(reloadID string) {
	o.ReloadID = reloadID
}

// WithPayload adds the payload to the create Acl Fcgi app accepted response
func (o *CreateACLFCGIAppAccepted) WithPayload(payload *models.ACL) *CreateACLFCGIAppAccepted {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create Acl Fcgi app accepted response
func (o *CreateACLFCGIAppAccepted) SetPayload(payload *models.ACL) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateACLFCGIAppAccepted) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	// response header Reload-ID

	reloadID := o.ReloadID
	if reloadID != "" {
		rw.Header().Set("Reload-ID", reloadID)
	}

	rw.WriteHeader(202)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateACLFCGIAppBadRequestCode is the HTTP code returned for type CreateACLFCGIAppBadRequest
const CreateACLFCGIAppBadRequestCode int = 400

/*
CreateACLFCGIAppBadRequest Bad request

swagger:response createAclFcgiAppBadRequest
*/
type CreateACLFCGIAppBadRequest struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateACLFCGIAppBadRequest creates CreateACLFCGIAppBadRequest with default headers values
func NewCreateACLFCGIAppBadRequest() *CreateACLFCGIAppBadRequest {

	return &CreateACLFCGIAppBadRequest{}
}

// WithConfigurationVersion adds the configurationVersion to the create Acl Fcgi app bad request response
func (o *CreateACLFCGIAppBadRequest) WithConfigurationVersion(configurationVersion string) *CreateACLFCGIAppBadRequest {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the create Acl Fcgi app bad request response
func (o *CreateACLFCGIAppBadRequest) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the create Acl Fcgi app bad request response
func (o *CreateACLFCGIAppBadRequest) WithPayload(payload *models.Error) *CreateACLFCGIAppBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create Acl Fcgi app bad request response
func (o *CreateACLFCGIAppBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateACLFCGIAppBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	// response header Configuration-Version

	configurationVersion := o.ConfigurationVersion
	if configurationVersion != "" {
		rw.Header().Set("Configuration-Version", configurationVersion)
	}

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateACLFCGIAppConflictCode is the HTTP code returned for type CreateACLFCGIAppConflict
const CreateACLFCGIAppConflictCode int = 409

/*
CreateACLFCGIAppConflict The specified resource already exists

swagger:response createAclFcgiAppConflict
*/
type CreateACLFCGIAppConflict struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateACLFCGIAppConflict creates CreateACLFCGIAppConflict with default headers values
func NewCreateACLFCGIAppConflict() *CreateACLFCGIAppConflict {

	return &CreateACLFCGIAppConflict{}
}

// WithConfigurationVersion adds the configurationVersion to the create Acl Fcgi app conflict response
func (o *CreateACLFCGIAppConflict) WithConfigurationVersion(configurationVersion string) *CreateACLFCGIAppConflict {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the create Acl Fcgi app conflict response
func (o *CreateACLFCGIAppConflict) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the create Acl Fcgi app conflict response
func (o *CreateACLFCGIAppConflict) WithPayload(payload *models.Error) *CreateACLFCGIAppConflict {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create Acl Fcgi app conflict response
func (o *CreateACLFCGIAppConflict) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateACLFCGIAppConflict) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	// response header Configuration-Version

	configurationVersion := o.ConfigurationVersion
	if configurationVersion != "" {
		rw.Header().Set("Configuration-Version", configurationVersion)
	}

	rw.WriteHeader(409)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*
CreateACLFCGIAppDefault General Error

swagger:response createAclFcgiAppDefault
*/
type CreateACLFCGIAppDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateACLFCGIAppDefault creates CreateACLFCGIAppDefault with default headers values
func NewCreateACLFCGIAppDefault(code int) *CreateACLFCGIAppDefault {
	if code <= 0 {
		code = 500
	}

	return &CreateACLFCGIAppDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the create Acl FCGI app default response
func (o *CreateACLFCGIAppDefault) WithStatusCode(code int) *CreateACLFCGIAppDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the create Acl FCGI app default response
func (o *CreateACLFCGIAppDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the create Acl FCGI app default response
func (o *CreateACLFCGIAppDefault) WithConfigurationVersion(configurationVersion string) *CreateACLFCGIAppDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the create Acl FCGI app default response
func (o *CreateACLFCGIAppDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the create Acl FCGI app default response
func (o *CreateACLFCGIAppDefault) WithPayload(payload *models.Error) *CreateACLFCGIAppDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create Acl FCGI app default response
func (o *CreateACLFCGIAppDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateACLFCGIAppDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	// response header Configuration-Version

	configurationVersion := o.ConfigurationVersion
	if configurationVersion != "" {
		rw.Header().Set("Configuration-Version", configurationVersion)
	}

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
