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

package maps

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v5/models"
)

// AddMapEntryCreatedCode is the HTTP code returned for type AddMapEntryCreated
const AddMapEntryCreatedCode int = 201

/*
AddMapEntryCreated Map entry created

swagger:response addMapEntryCreated
*/
type AddMapEntryCreated struct {

	/*
	  In: Body
	*/
	Payload *models.MapEntry `json:"body,omitempty"`
}

// NewAddMapEntryCreated creates AddMapEntryCreated with default headers values
func NewAddMapEntryCreated() *AddMapEntryCreated {

	return &AddMapEntryCreated{}
}

// WithPayload adds the payload to the add map entry created response
func (o *AddMapEntryCreated) WithPayload(payload *models.MapEntry) *AddMapEntryCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the add map entry created response
func (o *AddMapEntryCreated) SetPayload(payload *models.MapEntry) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AddMapEntryCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// AddMapEntryBadRequestCode is the HTTP code returned for type AddMapEntryBadRequest
const AddMapEntryBadRequestCode int = 400

/*
AddMapEntryBadRequest Bad request

swagger:response addMapEntryBadRequest
*/
type AddMapEntryBadRequest struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewAddMapEntryBadRequest creates AddMapEntryBadRequest with default headers values
func NewAddMapEntryBadRequest() *AddMapEntryBadRequest {

	return &AddMapEntryBadRequest{}
}

// WithConfigurationVersion adds the configurationVersion to the add map entry bad request response
func (o *AddMapEntryBadRequest) WithConfigurationVersion(configurationVersion string) *AddMapEntryBadRequest {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the add map entry bad request response
func (o *AddMapEntryBadRequest) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the add map entry bad request response
func (o *AddMapEntryBadRequest) WithPayload(payload *models.Error) *AddMapEntryBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the add map entry bad request response
func (o *AddMapEntryBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AddMapEntryBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// AddMapEntryConflictCode is the HTTP code returned for type AddMapEntryConflict
const AddMapEntryConflictCode int = 409

/*
AddMapEntryConflict The specified resource already exists

swagger:response addMapEntryConflict
*/
type AddMapEntryConflict struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewAddMapEntryConflict creates AddMapEntryConflict with default headers values
func NewAddMapEntryConflict() *AddMapEntryConflict {

	return &AddMapEntryConflict{}
}

// WithConfigurationVersion adds the configurationVersion to the add map entry conflict response
func (o *AddMapEntryConflict) WithConfigurationVersion(configurationVersion string) *AddMapEntryConflict {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the add map entry conflict response
func (o *AddMapEntryConflict) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the add map entry conflict response
func (o *AddMapEntryConflict) WithPayload(payload *models.Error) *AddMapEntryConflict {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the add map entry conflict response
func (o *AddMapEntryConflict) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AddMapEntryConflict) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
AddMapEntryDefault General Error

swagger:response addMapEntryDefault
*/
type AddMapEntryDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewAddMapEntryDefault creates AddMapEntryDefault with default headers values
func NewAddMapEntryDefault(code int) *AddMapEntryDefault {
	if code <= 0 {
		code = 500
	}

	return &AddMapEntryDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the add map entry default response
func (o *AddMapEntryDefault) WithStatusCode(code int) *AddMapEntryDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the add map entry default response
func (o *AddMapEntryDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the add map entry default response
func (o *AddMapEntryDefault) WithConfigurationVersion(configurationVersion string) *AddMapEntryDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the add map entry default response
func (o *AddMapEntryDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the add map entry default response
func (o *AddMapEntryDefault) WithPayload(payload *models.Error) *AddMapEntryDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the add map entry default response
func (o *AddMapEntryDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AddMapEntryDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
