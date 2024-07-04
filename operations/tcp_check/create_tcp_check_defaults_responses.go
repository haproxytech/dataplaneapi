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

package tcp_check

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v6/models"
)

// CreateTCPCheckDefaultsCreatedCode is the HTTP code returned for type CreateTCPCheckDefaultsCreated
const CreateTCPCheckDefaultsCreatedCode int = 201

/*
CreateTCPCheckDefaultsCreated TCP check created

swagger:response createTcpCheckDefaultsCreated
*/
type CreateTCPCheckDefaultsCreated struct {

	/*
	  In: Body
	*/
	Payload *models.TCPCheck `json:"body,omitempty"`
}

// NewCreateTCPCheckDefaultsCreated creates CreateTCPCheckDefaultsCreated with default headers values
func NewCreateTCPCheckDefaultsCreated() *CreateTCPCheckDefaultsCreated {

	return &CreateTCPCheckDefaultsCreated{}
}

// WithPayload adds the payload to the create Tcp check defaults created response
func (o *CreateTCPCheckDefaultsCreated) WithPayload(payload *models.TCPCheck) *CreateTCPCheckDefaultsCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create Tcp check defaults created response
func (o *CreateTCPCheckDefaultsCreated) SetPayload(payload *models.TCPCheck) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateTCPCheckDefaultsCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateTCPCheckDefaultsAcceptedCode is the HTTP code returned for type CreateTCPCheckDefaultsAccepted
const CreateTCPCheckDefaultsAcceptedCode int = 202

/*
CreateTCPCheckDefaultsAccepted Configuration change accepted and reload requested

swagger:response createTcpCheckDefaultsAccepted
*/
type CreateTCPCheckDefaultsAccepted struct {
	/*ID of the requested reload

	 */
	ReloadID string `json:"Reload-ID"`

	/*
	  In: Body
	*/
	Payload *models.TCPCheck `json:"body,omitempty"`
}

// NewCreateTCPCheckDefaultsAccepted creates CreateTCPCheckDefaultsAccepted with default headers values
func NewCreateTCPCheckDefaultsAccepted() *CreateTCPCheckDefaultsAccepted {

	return &CreateTCPCheckDefaultsAccepted{}
}

// WithReloadID adds the reloadId to the create Tcp check defaults accepted response
func (o *CreateTCPCheckDefaultsAccepted) WithReloadID(reloadID string) *CreateTCPCheckDefaultsAccepted {
	o.ReloadID = reloadID
	return o
}

// SetReloadID sets the reloadId to the create Tcp check defaults accepted response
func (o *CreateTCPCheckDefaultsAccepted) SetReloadID(reloadID string) {
	o.ReloadID = reloadID
}

// WithPayload adds the payload to the create Tcp check defaults accepted response
func (o *CreateTCPCheckDefaultsAccepted) WithPayload(payload *models.TCPCheck) *CreateTCPCheckDefaultsAccepted {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create Tcp check defaults accepted response
func (o *CreateTCPCheckDefaultsAccepted) SetPayload(payload *models.TCPCheck) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateTCPCheckDefaultsAccepted) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// CreateTCPCheckDefaultsBadRequestCode is the HTTP code returned for type CreateTCPCheckDefaultsBadRequest
const CreateTCPCheckDefaultsBadRequestCode int = 400

/*
CreateTCPCheckDefaultsBadRequest Bad request

swagger:response createTcpCheckDefaultsBadRequest
*/
type CreateTCPCheckDefaultsBadRequest struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateTCPCheckDefaultsBadRequest creates CreateTCPCheckDefaultsBadRequest with default headers values
func NewCreateTCPCheckDefaultsBadRequest() *CreateTCPCheckDefaultsBadRequest {

	return &CreateTCPCheckDefaultsBadRequest{}
}

// WithConfigurationVersion adds the configurationVersion to the create Tcp check defaults bad request response
func (o *CreateTCPCheckDefaultsBadRequest) WithConfigurationVersion(configurationVersion string) *CreateTCPCheckDefaultsBadRequest {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the create Tcp check defaults bad request response
func (o *CreateTCPCheckDefaultsBadRequest) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the create Tcp check defaults bad request response
func (o *CreateTCPCheckDefaultsBadRequest) WithPayload(payload *models.Error) *CreateTCPCheckDefaultsBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create Tcp check defaults bad request response
func (o *CreateTCPCheckDefaultsBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateTCPCheckDefaultsBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// CreateTCPCheckDefaultsConflictCode is the HTTP code returned for type CreateTCPCheckDefaultsConflict
const CreateTCPCheckDefaultsConflictCode int = 409

/*
CreateTCPCheckDefaultsConflict The specified resource already exists

swagger:response createTcpCheckDefaultsConflict
*/
type CreateTCPCheckDefaultsConflict struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateTCPCheckDefaultsConflict creates CreateTCPCheckDefaultsConflict with default headers values
func NewCreateTCPCheckDefaultsConflict() *CreateTCPCheckDefaultsConflict {

	return &CreateTCPCheckDefaultsConflict{}
}

// WithConfigurationVersion adds the configurationVersion to the create Tcp check defaults conflict response
func (o *CreateTCPCheckDefaultsConflict) WithConfigurationVersion(configurationVersion string) *CreateTCPCheckDefaultsConflict {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the create Tcp check defaults conflict response
func (o *CreateTCPCheckDefaultsConflict) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the create Tcp check defaults conflict response
func (o *CreateTCPCheckDefaultsConflict) WithPayload(payload *models.Error) *CreateTCPCheckDefaultsConflict {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create Tcp check defaults conflict response
func (o *CreateTCPCheckDefaultsConflict) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateTCPCheckDefaultsConflict) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
CreateTCPCheckDefaultsDefault General Error

swagger:response createTcpCheckDefaultsDefault
*/
type CreateTCPCheckDefaultsDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateTCPCheckDefaultsDefault creates CreateTCPCheckDefaultsDefault with default headers values
func NewCreateTCPCheckDefaultsDefault(code int) *CreateTCPCheckDefaultsDefault {
	if code <= 0 {
		code = 500
	}

	return &CreateTCPCheckDefaultsDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the create TCP check defaults default response
func (o *CreateTCPCheckDefaultsDefault) WithStatusCode(code int) *CreateTCPCheckDefaultsDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the create TCP check defaults default response
func (o *CreateTCPCheckDefaultsDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the create TCP check defaults default response
func (o *CreateTCPCheckDefaultsDefault) WithConfigurationVersion(configurationVersion string) *CreateTCPCheckDefaultsDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the create TCP check defaults default response
func (o *CreateTCPCheckDefaultsDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the create TCP check defaults default response
func (o *CreateTCPCheckDefaultsDefault) WithPayload(payload *models.Error) *CreateTCPCheckDefaultsDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create TCP check defaults default response
func (o *CreateTCPCheckDefaultsDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateTCPCheckDefaultsDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
