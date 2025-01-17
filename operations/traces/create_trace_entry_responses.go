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

package traces

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v6/models"
)

// CreateTraceEntryCreatedCode is the HTTP code returned for type CreateTraceEntryCreated
const CreateTraceEntryCreatedCode int = 201

/*
CreateTraceEntryCreated Trace entry added

swagger:response createTraceEntryCreated
*/
type CreateTraceEntryCreated struct {

	/*
	  In: Body
	*/
	Payload *models.TraceEntry `json:"body,omitempty"`
}

// NewCreateTraceEntryCreated creates CreateTraceEntryCreated with default headers values
func NewCreateTraceEntryCreated() *CreateTraceEntryCreated {

	return &CreateTraceEntryCreated{}
}

// WithPayload adds the payload to the create trace entry created response
func (o *CreateTraceEntryCreated) WithPayload(payload *models.TraceEntry) *CreateTraceEntryCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create trace entry created response
func (o *CreateTraceEntryCreated) SetPayload(payload *models.TraceEntry) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateTraceEntryCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateTraceEntryAcceptedCode is the HTTP code returned for type CreateTraceEntryAccepted
const CreateTraceEntryAcceptedCode int = 202

/*
CreateTraceEntryAccepted Configuration change accepted and reload requested

swagger:response createTraceEntryAccepted
*/
type CreateTraceEntryAccepted struct {
	/*ID of the requested reload

	 */
	ReloadID string `json:"Reload-ID"`

	/*
	  In: Body
	*/
	Payload *models.TraceEntry `json:"body,omitempty"`
}

// NewCreateTraceEntryAccepted creates CreateTraceEntryAccepted with default headers values
func NewCreateTraceEntryAccepted() *CreateTraceEntryAccepted {

	return &CreateTraceEntryAccepted{}
}

// WithReloadID adds the reloadId to the create trace entry accepted response
func (o *CreateTraceEntryAccepted) WithReloadID(reloadID string) *CreateTraceEntryAccepted {
	o.ReloadID = reloadID
	return o
}

// SetReloadID sets the reloadId to the create trace entry accepted response
func (o *CreateTraceEntryAccepted) SetReloadID(reloadID string) {
	o.ReloadID = reloadID
}

// WithPayload adds the payload to the create trace entry accepted response
func (o *CreateTraceEntryAccepted) WithPayload(payload *models.TraceEntry) *CreateTraceEntryAccepted {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create trace entry accepted response
func (o *CreateTraceEntryAccepted) SetPayload(payload *models.TraceEntry) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateTraceEntryAccepted) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// CreateTraceEntryBadRequestCode is the HTTP code returned for type CreateTraceEntryBadRequest
const CreateTraceEntryBadRequestCode int = 400

/*
CreateTraceEntryBadRequest Bad request

swagger:response createTraceEntryBadRequest
*/
type CreateTraceEntryBadRequest struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateTraceEntryBadRequest creates CreateTraceEntryBadRequest with default headers values
func NewCreateTraceEntryBadRequest() *CreateTraceEntryBadRequest {

	return &CreateTraceEntryBadRequest{}
}

// WithConfigurationVersion adds the configurationVersion to the create trace entry bad request response
func (o *CreateTraceEntryBadRequest) WithConfigurationVersion(configurationVersion string) *CreateTraceEntryBadRequest {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the create trace entry bad request response
func (o *CreateTraceEntryBadRequest) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the create trace entry bad request response
func (o *CreateTraceEntryBadRequest) WithPayload(payload *models.Error) *CreateTraceEntryBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create trace entry bad request response
func (o *CreateTraceEntryBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateTraceEntryBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// CreateTraceEntryConflictCode is the HTTP code returned for type CreateTraceEntryConflict
const CreateTraceEntryConflictCode int = 409

/*
CreateTraceEntryConflict The specified resource already exists

swagger:response createTraceEntryConflict
*/
type CreateTraceEntryConflict struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateTraceEntryConflict creates CreateTraceEntryConflict with default headers values
func NewCreateTraceEntryConflict() *CreateTraceEntryConflict {

	return &CreateTraceEntryConflict{}
}

// WithConfigurationVersion adds the configurationVersion to the create trace entry conflict response
func (o *CreateTraceEntryConflict) WithConfigurationVersion(configurationVersion string) *CreateTraceEntryConflict {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the create trace entry conflict response
func (o *CreateTraceEntryConflict) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the create trace entry conflict response
func (o *CreateTraceEntryConflict) WithPayload(payload *models.Error) *CreateTraceEntryConflict {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create trace entry conflict response
func (o *CreateTraceEntryConflict) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateTraceEntryConflict) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
CreateTraceEntryDefault General Error

swagger:response createTraceEntryDefault
*/
type CreateTraceEntryDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateTraceEntryDefault creates CreateTraceEntryDefault with default headers values
func NewCreateTraceEntryDefault(code int) *CreateTraceEntryDefault {
	if code <= 0 {
		code = 500
	}

	return &CreateTraceEntryDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the create trace entry default response
func (o *CreateTraceEntryDefault) WithStatusCode(code int) *CreateTraceEntryDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the create trace entry default response
func (o *CreateTraceEntryDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the create trace entry default response
func (o *CreateTraceEntryDefault) WithConfigurationVersion(configurationVersion string) *CreateTraceEntryDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the create trace entry default response
func (o *CreateTraceEntryDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the create trace entry default response
func (o *CreateTraceEntryDefault) WithPayload(payload *models.Error) *CreateTraceEntryDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create trace entry default response
func (o *CreateTraceEntryDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateTraceEntryDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
