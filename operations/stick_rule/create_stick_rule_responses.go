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

package stick_rule

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v3/models"
)

// CreateStickRuleCreatedCode is the HTTP code returned for type CreateStickRuleCreated
const CreateStickRuleCreatedCode int = 201

/*
CreateStickRuleCreated Stick Rule created

swagger:response createStickRuleCreated
*/
type CreateStickRuleCreated struct {

	/*
	  In: Body
	*/
	Payload *models.StickRule `json:"body,omitempty"`
}

// NewCreateStickRuleCreated creates CreateStickRuleCreated with default headers values
func NewCreateStickRuleCreated() *CreateStickRuleCreated {

	return &CreateStickRuleCreated{}
}

// WithPayload adds the payload to the create stick rule created response
func (o *CreateStickRuleCreated) WithPayload(payload *models.StickRule) *CreateStickRuleCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create stick rule created response
func (o *CreateStickRuleCreated) SetPayload(payload *models.StickRule) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateStickRuleCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateStickRuleAcceptedCode is the HTTP code returned for type CreateStickRuleAccepted
const CreateStickRuleAcceptedCode int = 202

/*
CreateStickRuleAccepted Configuration change accepted and reload requested

swagger:response createStickRuleAccepted
*/
type CreateStickRuleAccepted struct {
	/*ID of the requested reload

	 */
	ReloadID string `json:"Reload-ID"`

	/*
	  In: Body
	*/
	Payload *models.StickRule `json:"body,omitempty"`
}

// NewCreateStickRuleAccepted creates CreateStickRuleAccepted with default headers values
func NewCreateStickRuleAccepted() *CreateStickRuleAccepted {

	return &CreateStickRuleAccepted{}
}

// WithReloadID adds the reloadId to the create stick rule accepted response
func (o *CreateStickRuleAccepted) WithReloadID(reloadID string) *CreateStickRuleAccepted {
	o.ReloadID = reloadID
	return o
}

// SetReloadID sets the reloadId to the create stick rule accepted response
func (o *CreateStickRuleAccepted) SetReloadID(reloadID string) {
	o.ReloadID = reloadID
}

// WithPayload adds the payload to the create stick rule accepted response
func (o *CreateStickRuleAccepted) WithPayload(payload *models.StickRule) *CreateStickRuleAccepted {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create stick rule accepted response
func (o *CreateStickRuleAccepted) SetPayload(payload *models.StickRule) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateStickRuleAccepted) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// CreateStickRuleBadRequestCode is the HTTP code returned for type CreateStickRuleBadRequest
const CreateStickRuleBadRequestCode int = 400

/*
CreateStickRuleBadRequest Bad request

swagger:response createStickRuleBadRequest
*/
type CreateStickRuleBadRequest struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateStickRuleBadRequest creates CreateStickRuleBadRequest with default headers values
func NewCreateStickRuleBadRequest() *CreateStickRuleBadRequest {

	return &CreateStickRuleBadRequest{}
}

// WithConfigurationVersion adds the configurationVersion to the create stick rule bad request response
func (o *CreateStickRuleBadRequest) WithConfigurationVersion(configurationVersion string) *CreateStickRuleBadRequest {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the create stick rule bad request response
func (o *CreateStickRuleBadRequest) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the create stick rule bad request response
func (o *CreateStickRuleBadRequest) WithPayload(payload *models.Error) *CreateStickRuleBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create stick rule bad request response
func (o *CreateStickRuleBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateStickRuleBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// CreateStickRuleConflictCode is the HTTP code returned for type CreateStickRuleConflict
const CreateStickRuleConflictCode int = 409

/*
CreateStickRuleConflict The specified resource already exists

swagger:response createStickRuleConflict
*/
type CreateStickRuleConflict struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateStickRuleConflict creates CreateStickRuleConflict with default headers values
func NewCreateStickRuleConflict() *CreateStickRuleConflict {

	return &CreateStickRuleConflict{}
}

// WithConfigurationVersion adds the configurationVersion to the create stick rule conflict response
func (o *CreateStickRuleConflict) WithConfigurationVersion(configurationVersion string) *CreateStickRuleConflict {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the create stick rule conflict response
func (o *CreateStickRuleConflict) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the create stick rule conflict response
func (o *CreateStickRuleConflict) WithPayload(payload *models.Error) *CreateStickRuleConflict {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create stick rule conflict response
func (o *CreateStickRuleConflict) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateStickRuleConflict) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
CreateStickRuleDefault General Error

swagger:response createStickRuleDefault
*/
type CreateStickRuleDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateStickRuleDefault creates CreateStickRuleDefault with default headers values
func NewCreateStickRuleDefault(code int) *CreateStickRuleDefault {
	if code <= 0 {
		code = 500
	}

	return &CreateStickRuleDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the create stick rule default response
func (o *CreateStickRuleDefault) WithStatusCode(code int) *CreateStickRuleDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the create stick rule default response
func (o *CreateStickRuleDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the create stick rule default response
func (o *CreateStickRuleDefault) WithConfigurationVersion(configurationVersion string) *CreateStickRuleDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the create stick rule default response
func (o *CreateStickRuleDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the create stick rule default response
func (o *CreateStickRuleDefault) WithPayload(payload *models.Error) *CreateStickRuleDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create stick rule default response
func (o *CreateStickRuleDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateStickRuleDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
