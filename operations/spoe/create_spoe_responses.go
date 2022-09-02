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

// CreateSpoeCreatedCode is the HTTP code returned for type CreateSpoeCreated
const CreateSpoeCreatedCode int = 201

/*
CreateSpoeCreated SPOE file created with its entries

swagger:response createSpoeCreated
*/
type CreateSpoeCreated struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewCreateSpoeCreated creates CreateSpoeCreated with default headers values
func NewCreateSpoeCreated() *CreateSpoeCreated {

	return &CreateSpoeCreated{}
}

// WithPayload adds the payload to the create spoe created response
func (o *CreateSpoeCreated) WithPayload(payload string) *CreateSpoeCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create spoe created response
func (o *CreateSpoeCreated) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateSpoeCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// CreateSpoeBadRequestCode is the HTTP code returned for type CreateSpoeBadRequest
const CreateSpoeBadRequestCode int = 400

/*
CreateSpoeBadRequest Bad request

swagger:response createSpoeBadRequest
*/
type CreateSpoeBadRequest struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateSpoeBadRequest creates CreateSpoeBadRequest with default headers values
func NewCreateSpoeBadRequest() *CreateSpoeBadRequest {

	return &CreateSpoeBadRequest{}
}

// WithConfigurationVersion adds the configurationVersion to the create spoe bad request response
func (o *CreateSpoeBadRequest) WithConfigurationVersion(configurationVersion string) *CreateSpoeBadRequest {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the create spoe bad request response
func (o *CreateSpoeBadRequest) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the create spoe bad request response
func (o *CreateSpoeBadRequest) WithPayload(payload *models.Error) *CreateSpoeBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create spoe bad request response
func (o *CreateSpoeBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateSpoeBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// CreateSpoeConflictCode is the HTTP code returned for type CreateSpoeConflict
const CreateSpoeConflictCode int = 409

/*
CreateSpoeConflict The specified resource already exists

swagger:response createSpoeConflict
*/
type CreateSpoeConflict struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateSpoeConflict creates CreateSpoeConflict with default headers values
func NewCreateSpoeConflict() *CreateSpoeConflict {

	return &CreateSpoeConflict{}
}

// WithConfigurationVersion adds the configurationVersion to the create spoe conflict response
func (o *CreateSpoeConflict) WithConfigurationVersion(configurationVersion string) *CreateSpoeConflict {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the create spoe conflict response
func (o *CreateSpoeConflict) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the create spoe conflict response
func (o *CreateSpoeConflict) WithPayload(payload *models.Error) *CreateSpoeConflict {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create spoe conflict response
func (o *CreateSpoeConflict) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateSpoeConflict) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
CreateSpoeDefault General Error

swagger:response createSpoeDefault
*/
type CreateSpoeDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateSpoeDefault creates CreateSpoeDefault with default headers values
func NewCreateSpoeDefault(code int) *CreateSpoeDefault {
	if code <= 0 {
		code = 500
	}

	return &CreateSpoeDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the create spoe default response
func (o *CreateSpoeDefault) WithStatusCode(code int) *CreateSpoeDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the create spoe default response
func (o *CreateSpoeDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the create spoe default response
func (o *CreateSpoeDefault) WithConfigurationVersion(configurationVersion string) *CreateSpoeDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the create spoe default response
func (o *CreateSpoeDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the create spoe default response
func (o *CreateSpoeDefault) WithPayload(payload *models.Error) *CreateSpoeDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create spoe default response
func (o *CreateSpoeDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateSpoeDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
