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

package quic_initial_rule

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v6/models"
)

// CreateQUICInitialRuleFrontendCreatedCode is the HTTP code returned for type CreateQUICInitialRuleFrontendCreated
const CreateQUICInitialRuleFrontendCreatedCode int = 201

/*
CreateQUICInitialRuleFrontendCreated HTTP Response Rule created

swagger:response createQuicInitialRuleFrontendCreated
*/
type CreateQUICInitialRuleFrontendCreated struct {

	/*
	  In: Body
	*/
	Payload *models.QUICInitialRule `json:"body,omitempty"`
}

// NewCreateQUICInitialRuleFrontendCreated creates CreateQUICInitialRuleFrontendCreated with default headers values
func NewCreateQUICInitialRuleFrontendCreated() *CreateQUICInitialRuleFrontendCreated {

	return &CreateQUICInitialRuleFrontendCreated{}
}

// WithPayload adds the payload to the create Quic initial rule frontend created response
func (o *CreateQUICInitialRuleFrontendCreated) WithPayload(payload *models.QUICInitialRule) *CreateQUICInitialRuleFrontendCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create Quic initial rule frontend created response
func (o *CreateQUICInitialRuleFrontendCreated) SetPayload(payload *models.QUICInitialRule) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateQUICInitialRuleFrontendCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateQUICInitialRuleFrontendAcceptedCode is the HTTP code returned for type CreateQUICInitialRuleFrontendAccepted
const CreateQUICInitialRuleFrontendAcceptedCode int = 202

/*
CreateQUICInitialRuleFrontendAccepted Configuration change accepted and reload requested

swagger:response createQuicInitialRuleFrontendAccepted
*/
type CreateQUICInitialRuleFrontendAccepted struct {
	/*ID of the requested reload

	 */
	ReloadID string `json:"Reload-ID"`

	/*
	  In: Body
	*/
	Payload *models.QUICInitialRule `json:"body,omitempty"`
}

// NewCreateQUICInitialRuleFrontendAccepted creates CreateQUICInitialRuleFrontendAccepted with default headers values
func NewCreateQUICInitialRuleFrontendAccepted() *CreateQUICInitialRuleFrontendAccepted {

	return &CreateQUICInitialRuleFrontendAccepted{}
}

// WithReloadID adds the reloadId to the create Quic initial rule frontend accepted response
func (o *CreateQUICInitialRuleFrontendAccepted) WithReloadID(reloadID string) *CreateQUICInitialRuleFrontendAccepted {
	o.ReloadID = reloadID
	return o
}

// SetReloadID sets the reloadId to the create Quic initial rule frontend accepted response
func (o *CreateQUICInitialRuleFrontendAccepted) SetReloadID(reloadID string) {
	o.ReloadID = reloadID
}

// WithPayload adds the payload to the create Quic initial rule frontend accepted response
func (o *CreateQUICInitialRuleFrontendAccepted) WithPayload(payload *models.QUICInitialRule) *CreateQUICInitialRuleFrontendAccepted {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create Quic initial rule frontend accepted response
func (o *CreateQUICInitialRuleFrontendAccepted) SetPayload(payload *models.QUICInitialRule) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateQUICInitialRuleFrontendAccepted) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// CreateQUICInitialRuleFrontendBadRequestCode is the HTTP code returned for type CreateQUICInitialRuleFrontendBadRequest
const CreateQUICInitialRuleFrontendBadRequestCode int = 400

/*
CreateQUICInitialRuleFrontendBadRequest Bad request

swagger:response createQuicInitialRuleFrontendBadRequest
*/
type CreateQUICInitialRuleFrontendBadRequest struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateQUICInitialRuleFrontendBadRequest creates CreateQUICInitialRuleFrontendBadRequest with default headers values
func NewCreateQUICInitialRuleFrontendBadRequest() *CreateQUICInitialRuleFrontendBadRequest {

	return &CreateQUICInitialRuleFrontendBadRequest{}
}

// WithConfigurationVersion adds the configurationVersion to the create Quic initial rule frontend bad request response
func (o *CreateQUICInitialRuleFrontendBadRequest) WithConfigurationVersion(configurationVersion string) *CreateQUICInitialRuleFrontendBadRequest {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the create Quic initial rule frontend bad request response
func (o *CreateQUICInitialRuleFrontendBadRequest) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the create Quic initial rule frontend bad request response
func (o *CreateQUICInitialRuleFrontendBadRequest) WithPayload(payload *models.Error) *CreateQUICInitialRuleFrontendBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create Quic initial rule frontend bad request response
func (o *CreateQUICInitialRuleFrontendBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateQUICInitialRuleFrontendBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// CreateQUICInitialRuleFrontendConflictCode is the HTTP code returned for type CreateQUICInitialRuleFrontendConflict
const CreateQUICInitialRuleFrontendConflictCode int = 409

/*
CreateQUICInitialRuleFrontendConflict The specified resource already exists

swagger:response createQuicInitialRuleFrontendConflict
*/
type CreateQUICInitialRuleFrontendConflict struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateQUICInitialRuleFrontendConflict creates CreateQUICInitialRuleFrontendConflict with default headers values
func NewCreateQUICInitialRuleFrontendConflict() *CreateQUICInitialRuleFrontendConflict {

	return &CreateQUICInitialRuleFrontendConflict{}
}

// WithConfigurationVersion adds the configurationVersion to the create Quic initial rule frontend conflict response
func (o *CreateQUICInitialRuleFrontendConflict) WithConfigurationVersion(configurationVersion string) *CreateQUICInitialRuleFrontendConflict {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the create Quic initial rule frontend conflict response
func (o *CreateQUICInitialRuleFrontendConflict) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the create Quic initial rule frontend conflict response
func (o *CreateQUICInitialRuleFrontendConflict) WithPayload(payload *models.Error) *CreateQUICInitialRuleFrontendConflict {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create Quic initial rule frontend conflict response
func (o *CreateQUICInitialRuleFrontendConflict) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateQUICInitialRuleFrontendConflict) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
CreateQUICInitialRuleFrontendDefault General Error

swagger:response createQuicInitialRuleFrontendDefault
*/
type CreateQUICInitialRuleFrontendDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateQUICInitialRuleFrontendDefault creates CreateQUICInitialRuleFrontendDefault with default headers values
func NewCreateQUICInitialRuleFrontendDefault(code int) *CreateQUICInitialRuleFrontendDefault {
	if code <= 0 {
		code = 500
	}

	return &CreateQUICInitialRuleFrontendDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the create QUIC initial rule frontend default response
func (o *CreateQUICInitialRuleFrontendDefault) WithStatusCode(code int) *CreateQUICInitialRuleFrontendDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the create QUIC initial rule frontend default response
func (o *CreateQUICInitialRuleFrontendDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the create QUIC initial rule frontend default response
func (o *CreateQUICInitialRuleFrontendDefault) WithConfigurationVersion(configurationVersion string) *CreateQUICInitialRuleFrontendDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the create QUIC initial rule frontend default response
func (o *CreateQUICInitialRuleFrontendDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the create QUIC initial rule frontend default response
func (o *CreateQUICInitialRuleFrontendDefault) WithPayload(payload *models.Error) *CreateQUICInitialRuleFrontendDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create QUIC initial rule frontend default response
func (o *CreateQUICInitialRuleFrontendDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateQUICInitialRuleFrontendDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
