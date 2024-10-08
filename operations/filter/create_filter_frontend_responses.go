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

package filter

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v6/models"
)

// CreateFilterFrontendCreatedCode is the HTTP code returned for type CreateFilterFrontendCreated
const CreateFilterFrontendCreatedCode int = 201

/*
CreateFilterFrontendCreated Filter created

swagger:response createFilterFrontendCreated
*/
type CreateFilterFrontendCreated struct {

	/*
	  In: Body
	*/
	Payload *models.Filter `json:"body,omitempty"`
}

// NewCreateFilterFrontendCreated creates CreateFilterFrontendCreated with default headers values
func NewCreateFilterFrontendCreated() *CreateFilterFrontendCreated {

	return &CreateFilterFrontendCreated{}
}

// WithPayload adds the payload to the create filter frontend created response
func (o *CreateFilterFrontendCreated) WithPayload(payload *models.Filter) *CreateFilterFrontendCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create filter frontend created response
func (o *CreateFilterFrontendCreated) SetPayload(payload *models.Filter) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateFilterFrontendCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateFilterFrontendAcceptedCode is the HTTP code returned for type CreateFilterFrontendAccepted
const CreateFilterFrontendAcceptedCode int = 202

/*
CreateFilterFrontendAccepted Configuration change accepted and reload requested

swagger:response createFilterFrontendAccepted
*/
type CreateFilterFrontendAccepted struct {
	/*ID of the requested reload

	 */
	ReloadID string `json:"Reload-ID"`

	/*
	  In: Body
	*/
	Payload *models.Filter `json:"body,omitempty"`
}

// NewCreateFilterFrontendAccepted creates CreateFilterFrontendAccepted with default headers values
func NewCreateFilterFrontendAccepted() *CreateFilterFrontendAccepted {

	return &CreateFilterFrontendAccepted{}
}

// WithReloadID adds the reloadId to the create filter frontend accepted response
func (o *CreateFilterFrontendAccepted) WithReloadID(reloadID string) *CreateFilterFrontendAccepted {
	o.ReloadID = reloadID
	return o
}

// SetReloadID sets the reloadId to the create filter frontend accepted response
func (o *CreateFilterFrontendAccepted) SetReloadID(reloadID string) {
	o.ReloadID = reloadID
}

// WithPayload adds the payload to the create filter frontend accepted response
func (o *CreateFilterFrontendAccepted) WithPayload(payload *models.Filter) *CreateFilterFrontendAccepted {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create filter frontend accepted response
func (o *CreateFilterFrontendAccepted) SetPayload(payload *models.Filter) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateFilterFrontendAccepted) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// CreateFilterFrontendBadRequestCode is the HTTP code returned for type CreateFilterFrontendBadRequest
const CreateFilterFrontendBadRequestCode int = 400

/*
CreateFilterFrontendBadRequest Bad request

swagger:response createFilterFrontendBadRequest
*/
type CreateFilterFrontendBadRequest struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateFilterFrontendBadRequest creates CreateFilterFrontendBadRequest with default headers values
func NewCreateFilterFrontendBadRequest() *CreateFilterFrontendBadRequest {

	return &CreateFilterFrontendBadRequest{}
}

// WithConfigurationVersion adds the configurationVersion to the create filter frontend bad request response
func (o *CreateFilterFrontendBadRequest) WithConfigurationVersion(configurationVersion string) *CreateFilterFrontendBadRequest {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the create filter frontend bad request response
func (o *CreateFilterFrontendBadRequest) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the create filter frontend bad request response
func (o *CreateFilterFrontendBadRequest) WithPayload(payload *models.Error) *CreateFilterFrontendBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create filter frontend bad request response
func (o *CreateFilterFrontendBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateFilterFrontendBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// CreateFilterFrontendConflictCode is the HTTP code returned for type CreateFilterFrontendConflict
const CreateFilterFrontendConflictCode int = 409

/*
CreateFilterFrontendConflict The specified resource already exists

swagger:response createFilterFrontendConflict
*/
type CreateFilterFrontendConflict struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateFilterFrontendConflict creates CreateFilterFrontendConflict with default headers values
func NewCreateFilterFrontendConflict() *CreateFilterFrontendConflict {

	return &CreateFilterFrontendConflict{}
}

// WithConfigurationVersion adds the configurationVersion to the create filter frontend conflict response
func (o *CreateFilterFrontendConflict) WithConfigurationVersion(configurationVersion string) *CreateFilterFrontendConflict {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the create filter frontend conflict response
func (o *CreateFilterFrontendConflict) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the create filter frontend conflict response
func (o *CreateFilterFrontendConflict) WithPayload(payload *models.Error) *CreateFilterFrontendConflict {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create filter frontend conflict response
func (o *CreateFilterFrontendConflict) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateFilterFrontendConflict) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
CreateFilterFrontendDefault General Error

swagger:response createFilterFrontendDefault
*/
type CreateFilterFrontendDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateFilterFrontendDefault creates CreateFilterFrontendDefault with default headers values
func NewCreateFilterFrontendDefault(code int) *CreateFilterFrontendDefault {
	if code <= 0 {
		code = 500
	}

	return &CreateFilterFrontendDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the create filter frontend default response
func (o *CreateFilterFrontendDefault) WithStatusCode(code int) *CreateFilterFrontendDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the create filter frontend default response
func (o *CreateFilterFrontendDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the create filter frontend default response
func (o *CreateFilterFrontendDefault) WithConfigurationVersion(configurationVersion string) *CreateFilterFrontendDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the create filter frontend default response
func (o *CreateFilterFrontendDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the create filter frontend default response
func (o *CreateFilterFrontendDefault) WithPayload(payload *models.Error) *CreateFilterFrontendDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create filter frontend default response
func (o *CreateFilterFrontendDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateFilterFrontendDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
