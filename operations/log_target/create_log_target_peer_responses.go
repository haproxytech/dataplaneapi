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

package log_target

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v6/models"
)

// CreateLogTargetPeerCreatedCode is the HTTP code returned for type CreateLogTargetPeerCreated
const CreateLogTargetPeerCreatedCode int = 201

/*
CreateLogTargetPeerCreated Log Target created

swagger:response createLogTargetPeerCreated
*/
type CreateLogTargetPeerCreated struct {

	/*
	  In: Body
	*/
	Payload *models.LogTarget `json:"body,omitempty"`
}

// NewCreateLogTargetPeerCreated creates CreateLogTargetPeerCreated with default headers values
func NewCreateLogTargetPeerCreated() *CreateLogTargetPeerCreated {

	return &CreateLogTargetPeerCreated{}
}

// WithPayload adds the payload to the create log target peer created response
func (o *CreateLogTargetPeerCreated) WithPayload(payload *models.LogTarget) *CreateLogTargetPeerCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create log target peer created response
func (o *CreateLogTargetPeerCreated) SetPayload(payload *models.LogTarget) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateLogTargetPeerCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateLogTargetPeerAcceptedCode is the HTTP code returned for type CreateLogTargetPeerAccepted
const CreateLogTargetPeerAcceptedCode int = 202

/*
CreateLogTargetPeerAccepted Configuration change accepted and reload requested

swagger:response createLogTargetPeerAccepted
*/
type CreateLogTargetPeerAccepted struct {
	/*ID of the requested reload

	 */
	ReloadID string `json:"Reload-ID"`

	/*
	  In: Body
	*/
	Payload *models.LogTarget `json:"body,omitempty"`
}

// NewCreateLogTargetPeerAccepted creates CreateLogTargetPeerAccepted with default headers values
func NewCreateLogTargetPeerAccepted() *CreateLogTargetPeerAccepted {

	return &CreateLogTargetPeerAccepted{}
}

// WithReloadID adds the reloadId to the create log target peer accepted response
func (o *CreateLogTargetPeerAccepted) WithReloadID(reloadID string) *CreateLogTargetPeerAccepted {
	o.ReloadID = reloadID
	return o
}

// SetReloadID sets the reloadId to the create log target peer accepted response
func (o *CreateLogTargetPeerAccepted) SetReloadID(reloadID string) {
	o.ReloadID = reloadID
}

// WithPayload adds the payload to the create log target peer accepted response
func (o *CreateLogTargetPeerAccepted) WithPayload(payload *models.LogTarget) *CreateLogTargetPeerAccepted {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create log target peer accepted response
func (o *CreateLogTargetPeerAccepted) SetPayload(payload *models.LogTarget) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateLogTargetPeerAccepted) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// CreateLogTargetPeerBadRequestCode is the HTTP code returned for type CreateLogTargetPeerBadRequest
const CreateLogTargetPeerBadRequestCode int = 400

/*
CreateLogTargetPeerBadRequest Bad request

swagger:response createLogTargetPeerBadRequest
*/
type CreateLogTargetPeerBadRequest struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateLogTargetPeerBadRequest creates CreateLogTargetPeerBadRequest with default headers values
func NewCreateLogTargetPeerBadRequest() *CreateLogTargetPeerBadRequest {

	return &CreateLogTargetPeerBadRequest{}
}

// WithConfigurationVersion adds the configurationVersion to the create log target peer bad request response
func (o *CreateLogTargetPeerBadRequest) WithConfigurationVersion(configurationVersion string) *CreateLogTargetPeerBadRequest {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the create log target peer bad request response
func (o *CreateLogTargetPeerBadRequest) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the create log target peer bad request response
func (o *CreateLogTargetPeerBadRequest) WithPayload(payload *models.Error) *CreateLogTargetPeerBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create log target peer bad request response
func (o *CreateLogTargetPeerBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateLogTargetPeerBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// CreateLogTargetPeerConflictCode is the HTTP code returned for type CreateLogTargetPeerConflict
const CreateLogTargetPeerConflictCode int = 409

/*
CreateLogTargetPeerConflict The specified resource already exists

swagger:response createLogTargetPeerConflict
*/
type CreateLogTargetPeerConflict struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateLogTargetPeerConflict creates CreateLogTargetPeerConflict with default headers values
func NewCreateLogTargetPeerConflict() *CreateLogTargetPeerConflict {

	return &CreateLogTargetPeerConflict{}
}

// WithConfigurationVersion adds the configurationVersion to the create log target peer conflict response
func (o *CreateLogTargetPeerConflict) WithConfigurationVersion(configurationVersion string) *CreateLogTargetPeerConflict {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the create log target peer conflict response
func (o *CreateLogTargetPeerConflict) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the create log target peer conflict response
func (o *CreateLogTargetPeerConflict) WithPayload(payload *models.Error) *CreateLogTargetPeerConflict {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create log target peer conflict response
func (o *CreateLogTargetPeerConflict) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateLogTargetPeerConflict) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
CreateLogTargetPeerDefault General Error

swagger:response createLogTargetPeerDefault
*/
type CreateLogTargetPeerDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateLogTargetPeerDefault creates CreateLogTargetPeerDefault with default headers values
func NewCreateLogTargetPeerDefault(code int) *CreateLogTargetPeerDefault {
	if code <= 0 {
		code = 500
	}

	return &CreateLogTargetPeerDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the create log target peer default response
func (o *CreateLogTargetPeerDefault) WithStatusCode(code int) *CreateLogTargetPeerDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the create log target peer default response
func (o *CreateLogTargetPeerDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the create log target peer default response
func (o *CreateLogTargetPeerDefault) WithConfigurationVersion(configurationVersion string) *CreateLogTargetPeerDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the create log target peer default response
func (o *CreateLogTargetPeerDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the create log target peer default response
func (o *CreateLogTargetPeerDefault) WithPayload(payload *models.Error) *CreateLogTargetPeerDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create log target peer default response
func (o *CreateLogTargetPeerDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateLogTargetPeerDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
