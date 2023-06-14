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

// AddPayloadRuntimeMapCreatedCode is the HTTP code returned for type AddPayloadRuntimeMapCreated
const AddPayloadRuntimeMapCreatedCode int = 201

/*
AddPayloadRuntimeMapCreated Map payload added

swagger:response addPayloadRuntimeMapCreated
*/
type AddPayloadRuntimeMapCreated struct {

	/*
	  In: Body
	*/
	Payload models.MapEntries `json:"body,omitempty"`
}

// NewAddPayloadRuntimeMapCreated creates AddPayloadRuntimeMapCreated with default headers values
func NewAddPayloadRuntimeMapCreated() *AddPayloadRuntimeMapCreated {

	return &AddPayloadRuntimeMapCreated{}
}

// WithPayload adds the payload to the add payload runtime map created response
func (o *AddPayloadRuntimeMapCreated) WithPayload(payload models.MapEntries) *AddPayloadRuntimeMapCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the add payload runtime map created response
func (o *AddPayloadRuntimeMapCreated) SetPayload(payload models.MapEntries) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AddPayloadRuntimeMapCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = models.MapEntries{}
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// AddPayloadRuntimeMapBadRequestCode is the HTTP code returned for type AddPayloadRuntimeMapBadRequest
const AddPayloadRuntimeMapBadRequestCode int = 400

/*
AddPayloadRuntimeMapBadRequest Bad request

swagger:response addPayloadRuntimeMapBadRequest
*/
type AddPayloadRuntimeMapBadRequest struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewAddPayloadRuntimeMapBadRequest creates AddPayloadRuntimeMapBadRequest with default headers values
func NewAddPayloadRuntimeMapBadRequest() *AddPayloadRuntimeMapBadRequest {

	return &AddPayloadRuntimeMapBadRequest{}
}

// WithConfigurationVersion adds the configurationVersion to the add payload runtime map bad request response
func (o *AddPayloadRuntimeMapBadRequest) WithConfigurationVersion(configurationVersion string) *AddPayloadRuntimeMapBadRequest {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the add payload runtime map bad request response
func (o *AddPayloadRuntimeMapBadRequest) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the add payload runtime map bad request response
func (o *AddPayloadRuntimeMapBadRequest) WithPayload(payload *models.Error) *AddPayloadRuntimeMapBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the add payload runtime map bad request response
func (o *AddPayloadRuntimeMapBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AddPayloadRuntimeMapBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

/*
AddPayloadRuntimeMapDefault General Error

swagger:response addPayloadRuntimeMapDefault
*/
type AddPayloadRuntimeMapDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewAddPayloadRuntimeMapDefault creates AddPayloadRuntimeMapDefault with default headers values
func NewAddPayloadRuntimeMapDefault(code int) *AddPayloadRuntimeMapDefault {
	if code <= 0 {
		code = 500
	}

	return &AddPayloadRuntimeMapDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the add payload runtime map default response
func (o *AddPayloadRuntimeMapDefault) WithStatusCode(code int) *AddPayloadRuntimeMapDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the add payload runtime map default response
func (o *AddPayloadRuntimeMapDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the add payload runtime map default response
func (o *AddPayloadRuntimeMapDefault) WithConfigurationVersion(configurationVersion string) *AddPayloadRuntimeMapDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the add payload runtime map default response
func (o *AddPayloadRuntimeMapDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the add payload runtime map default response
func (o *AddPayloadRuntimeMapDefault) WithPayload(payload *models.Error) *AddPayloadRuntimeMapDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the add payload runtime map default response
func (o *AddPayloadRuntimeMapDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AddPayloadRuntimeMapDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
