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

package crt_load

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v6/models"
)

// CreateCrtLoadCreatedCode is the HTTP code returned for type CreateCrtLoadCreated
const CreateCrtLoadCreatedCode int = 201

/*
CreateCrtLoadCreated Certificate load entry created

swagger:response createCrtLoadCreated
*/
type CreateCrtLoadCreated struct {

	/*
	  In: Body
	*/
	Payload *models.CrtLoad `json:"body,omitempty"`
}

// NewCreateCrtLoadCreated creates CreateCrtLoadCreated with default headers values
func NewCreateCrtLoadCreated() *CreateCrtLoadCreated {

	return &CreateCrtLoadCreated{}
}

// WithPayload adds the payload to the create crt load created response
func (o *CreateCrtLoadCreated) WithPayload(payload *models.CrtLoad) *CreateCrtLoadCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create crt load created response
func (o *CreateCrtLoadCreated) SetPayload(payload *models.CrtLoad) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateCrtLoadCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateCrtLoadAcceptedCode is the HTTP code returned for type CreateCrtLoadAccepted
const CreateCrtLoadAcceptedCode int = 202

/*
CreateCrtLoadAccepted Configuration change accepted and reload requested

swagger:response createCrtLoadAccepted
*/
type CreateCrtLoadAccepted struct {
	/*ID of the requested reload

	 */
	ReloadID string `json:"Reload-ID"`

	/*
	  In: Body
	*/
	Payload *models.CrtLoad `json:"body,omitempty"`
}

// NewCreateCrtLoadAccepted creates CreateCrtLoadAccepted with default headers values
func NewCreateCrtLoadAccepted() *CreateCrtLoadAccepted {

	return &CreateCrtLoadAccepted{}
}

// WithReloadID adds the reloadId to the create crt load accepted response
func (o *CreateCrtLoadAccepted) WithReloadID(reloadID string) *CreateCrtLoadAccepted {
	o.ReloadID = reloadID
	return o
}

// SetReloadID sets the reloadId to the create crt load accepted response
func (o *CreateCrtLoadAccepted) SetReloadID(reloadID string) {
	o.ReloadID = reloadID
}

// WithPayload adds the payload to the create crt load accepted response
func (o *CreateCrtLoadAccepted) WithPayload(payload *models.CrtLoad) *CreateCrtLoadAccepted {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create crt load accepted response
func (o *CreateCrtLoadAccepted) SetPayload(payload *models.CrtLoad) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateCrtLoadAccepted) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// CreateCrtLoadBadRequestCode is the HTTP code returned for type CreateCrtLoadBadRequest
const CreateCrtLoadBadRequestCode int = 400

/*
CreateCrtLoadBadRequest Bad request

swagger:response createCrtLoadBadRequest
*/
type CreateCrtLoadBadRequest struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateCrtLoadBadRequest creates CreateCrtLoadBadRequest with default headers values
func NewCreateCrtLoadBadRequest() *CreateCrtLoadBadRequest {

	return &CreateCrtLoadBadRequest{}
}

// WithConfigurationVersion adds the configurationVersion to the create crt load bad request response
func (o *CreateCrtLoadBadRequest) WithConfigurationVersion(configurationVersion string) *CreateCrtLoadBadRequest {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the create crt load bad request response
func (o *CreateCrtLoadBadRequest) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the create crt load bad request response
func (o *CreateCrtLoadBadRequest) WithPayload(payload *models.Error) *CreateCrtLoadBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create crt load bad request response
func (o *CreateCrtLoadBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateCrtLoadBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// CreateCrtLoadConflictCode is the HTTP code returned for type CreateCrtLoadConflict
const CreateCrtLoadConflictCode int = 409

/*
CreateCrtLoadConflict The specified resource already exists

swagger:response createCrtLoadConflict
*/
type CreateCrtLoadConflict struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateCrtLoadConflict creates CreateCrtLoadConflict with default headers values
func NewCreateCrtLoadConflict() *CreateCrtLoadConflict {

	return &CreateCrtLoadConflict{}
}

// WithConfigurationVersion adds the configurationVersion to the create crt load conflict response
func (o *CreateCrtLoadConflict) WithConfigurationVersion(configurationVersion string) *CreateCrtLoadConflict {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the create crt load conflict response
func (o *CreateCrtLoadConflict) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the create crt load conflict response
func (o *CreateCrtLoadConflict) WithPayload(payload *models.Error) *CreateCrtLoadConflict {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create crt load conflict response
func (o *CreateCrtLoadConflict) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateCrtLoadConflict) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
CreateCrtLoadDefault General Error

swagger:response createCrtLoadDefault
*/
type CreateCrtLoadDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateCrtLoadDefault creates CreateCrtLoadDefault with default headers values
func NewCreateCrtLoadDefault(code int) *CreateCrtLoadDefault {
	if code <= 0 {
		code = 500
	}

	return &CreateCrtLoadDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the create crt load default response
func (o *CreateCrtLoadDefault) WithStatusCode(code int) *CreateCrtLoadDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the create crt load default response
func (o *CreateCrtLoadDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the create crt load default response
func (o *CreateCrtLoadDefault) WithConfigurationVersion(configurationVersion string) *CreateCrtLoadDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the create crt load default response
func (o *CreateCrtLoadDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the create crt load default response
func (o *CreateCrtLoadDefault) WithPayload(payload *models.Error) *CreateCrtLoadDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create crt load default response
func (o *CreateCrtLoadDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateCrtLoadDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
