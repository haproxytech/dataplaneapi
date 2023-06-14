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

package userlist

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v5/models"
)

// CreateUserlistCreatedCode is the HTTP code returned for type CreateUserlistCreated
const CreateUserlistCreatedCode int = 201

/*
CreateUserlistCreated Userlist created

swagger:response createUserlistCreated
*/
type CreateUserlistCreated struct {

	/*
	  In: Body
	*/
	Payload *models.Userlist `json:"body,omitempty"`
}

// NewCreateUserlistCreated creates CreateUserlistCreated with default headers values
func NewCreateUserlistCreated() *CreateUserlistCreated {

	return &CreateUserlistCreated{}
}

// WithPayload adds the payload to the create userlist created response
func (o *CreateUserlistCreated) WithPayload(payload *models.Userlist) *CreateUserlistCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create userlist created response
func (o *CreateUserlistCreated) SetPayload(payload *models.Userlist) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateUserlistCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateUserlistAcceptedCode is the HTTP code returned for type CreateUserlistAccepted
const CreateUserlistAcceptedCode int = 202

/*
CreateUserlistAccepted Configuration change accepted and reload requested

swagger:response createUserlistAccepted
*/
type CreateUserlistAccepted struct {
	/*ID of the requested reload

	 */
	ReloadID string `json:"Reload-ID"`

	/*
	  In: Body
	*/
	Payload *models.Userlist `json:"body,omitempty"`
}

// NewCreateUserlistAccepted creates CreateUserlistAccepted with default headers values
func NewCreateUserlistAccepted() *CreateUserlistAccepted {

	return &CreateUserlistAccepted{}
}

// WithReloadID adds the reloadId to the create userlist accepted response
func (o *CreateUserlistAccepted) WithReloadID(reloadID string) *CreateUserlistAccepted {
	o.ReloadID = reloadID
	return o
}

// SetReloadID sets the reloadId to the create userlist accepted response
func (o *CreateUserlistAccepted) SetReloadID(reloadID string) {
	o.ReloadID = reloadID
}

// WithPayload adds the payload to the create userlist accepted response
func (o *CreateUserlistAccepted) WithPayload(payload *models.Userlist) *CreateUserlistAccepted {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create userlist accepted response
func (o *CreateUserlistAccepted) SetPayload(payload *models.Userlist) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateUserlistAccepted) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// CreateUserlistBadRequestCode is the HTTP code returned for type CreateUserlistBadRequest
const CreateUserlistBadRequestCode int = 400

/*
CreateUserlistBadRequest Bad request

swagger:response createUserlistBadRequest
*/
type CreateUserlistBadRequest struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateUserlistBadRequest creates CreateUserlistBadRequest with default headers values
func NewCreateUserlistBadRequest() *CreateUserlistBadRequest {

	return &CreateUserlistBadRequest{}
}

// WithConfigurationVersion adds the configurationVersion to the create userlist bad request response
func (o *CreateUserlistBadRequest) WithConfigurationVersion(configurationVersion string) *CreateUserlistBadRequest {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the create userlist bad request response
func (o *CreateUserlistBadRequest) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the create userlist bad request response
func (o *CreateUserlistBadRequest) WithPayload(payload *models.Error) *CreateUserlistBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create userlist bad request response
func (o *CreateUserlistBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateUserlistBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// CreateUserlistConflictCode is the HTTP code returned for type CreateUserlistConflict
const CreateUserlistConflictCode int = 409

/*
CreateUserlistConflict The specified resource already exists

swagger:response createUserlistConflict
*/
type CreateUserlistConflict struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateUserlistConflict creates CreateUserlistConflict with default headers values
func NewCreateUserlistConflict() *CreateUserlistConflict {

	return &CreateUserlistConflict{}
}

// WithConfigurationVersion adds the configurationVersion to the create userlist conflict response
func (o *CreateUserlistConflict) WithConfigurationVersion(configurationVersion string) *CreateUserlistConflict {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the create userlist conflict response
func (o *CreateUserlistConflict) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the create userlist conflict response
func (o *CreateUserlistConflict) WithPayload(payload *models.Error) *CreateUserlistConflict {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create userlist conflict response
func (o *CreateUserlistConflict) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateUserlistConflict) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
CreateUserlistDefault General Error

swagger:response createUserlistDefault
*/
type CreateUserlistDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateUserlistDefault creates CreateUserlistDefault with default headers values
func NewCreateUserlistDefault(code int) *CreateUserlistDefault {
	if code <= 0 {
		code = 500
	}

	return &CreateUserlistDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the create userlist default response
func (o *CreateUserlistDefault) WithStatusCode(code int) *CreateUserlistDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the create userlist default response
func (o *CreateUserlistDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the create userlist default response
func (o *CreateUserlistDefault) WithConfigurationVersion(configurationVersion string) *CreateUserlistDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the create userlist default response
func (o *CreateUserlistDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the create userlist default response
func (o *CreateUserlistDefault) WithPayload(payload *models.Error) *CreateUserlistDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create userlist default response
func (o *CreateUserlistDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateUserlistDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
