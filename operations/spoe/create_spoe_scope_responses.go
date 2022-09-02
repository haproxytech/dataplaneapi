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

package spoe

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v4/models"
)

// CreateSpoeScopeCreatedCode is the HTTP code returned for type CreateSpoeScopeCreated
const CreateSpoeScopeCreatedCode int = 201

/*
CreateSpoeScopeCreated Spoe scope created

swagger:response createSpoeScopeCreated
*/
type CreateSpoeScopeCreated struct {

	/*
	  In: Body
	*/
	Payload models.SpoeScope `json:"body,omitempty"`
}

// NewCreateSpoeScopeCreated creates CreateSpoeScopeCreated with default headers values
func NewCreateSpoeScopeCreated() *CreateSpoeScopeCreated {

	return &CreateSpoeScopeCreated{}
}

// WithPayload adds the payload to the create spoe scope created response
func (o *CreateSpoeScopeCreated) WithPayload(payload models.SpoeScope) *CreateSpoeScopeCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create spoe scope created response
func (o *CreateSpoeScopeCreated) SetPayload(payload models.SpoeScope) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateSpoeScopeCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// CreateSpoeScopeBadRequestCode is the HTTP code returned for type CreateSpoeScopeBadRequest
const CreateSpoeScopeBadRequestCode int = 400

/*
CreateSpoeScopeBadRequest Bad request

swagger:response createSpoeScopeBadRequest
*/
type CreateSpoeScopeBadRequest struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateSpoeScopeBadRequest creates CreateSpoeScopeBadRequest with default headers values
func NewCreateSpoeScopeBadRequest() *CreateSpoeScopeBadRequest {

	return &CreateSpoeScopeBadRequest{}
}

// WithConfigurationVersion adds the configurationVersion to the create spoe scope bad request response
func (o *CreateSpoeScopeBadRequest) WithConfigurationVersion(configurationVersion string) *CreateSpoeScopeBadRequest {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the create spoe scope bad request response
func (o *CreateSpoeScopeBadRequest) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the create spoe scope bad request response
func (o *CreateSpoeScopeBadRequest) WithPayload(payload *models.Error) *CreateSpoeScopeBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create spoe scope bad request response
func (o *CreateSpoeScopeBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateSpoeScopeBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// CreateSpoeScopeConflictCode is the HTTP code returned for type CreateSpoeScopeConflict
const CreateSpoeScopeConflictCode int = 409

/*
CreateSpoeScopeConflict The specified resource already exists

swagger:response createSpoeScopeConflict
*/
type CreateSpoeScopeConflict struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateSpoeScopeConflict creates CreateSpoeScopeConflict with default headers values
func NewCreateSpoeScopeConflict() *CreateSpoeScopeConflict {

	return &CreateSpoeScopeConflict{}
}

// WithConfigurationVersion adds the configurationVersion to the create spoe scope conflict response
func (o *CreateSpoeScopeConflict) WithConfigurationVersion(configurationVersion string) *CreateSpoeScopeConflict {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the create spoe scope conflict response
func (o *CreateSpoeScopeConflict) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the create spoe scope conflict response
func (o *CreateSpoeScopeConflict) WithPayload(payload *models.Error) *CreateSpoeScopeConflict {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create spoe scope conflict response
func (o *CreateSpoeScopeConflict) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateSpoeScopeConflict) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
CreateSpoeScopeDefault General Error

swagger:response createSpoeScopeDefault
*/
type CreateSpoeScopeDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateSpoeScopeDefault creates CreateSpoeScopeDefault with default headers values
func NewCreateSpoeScopeDefault(code int) *CreateSpoeScopeDefault {
	if code <= 0 {
		code = 500
	}

	return &CreateSpoeScopeDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the create spoe scope default response
func (o *CreateSpoeScopeDefault) WithStatusCode(code int) *CreateSpoeScopeDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the create spoe scope default response
func (o *CreateSpoeScopeDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the create spoe scope default response
func (o *CreateSpoeScopeDefault) WithConfigurationVersion(configurationVersion string) *CreateSpoeScopeDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the create spoe scope default response
func (o *CreateSpoeScopeDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the create spoe scope default response
func (o *CreateSpoeScopeDefault) WithPayload(payload *models.Error) *CreateSpoeScopeDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create spoe scope default response
func (o *CreateSpoeScopeDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateSpoeScopeDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
