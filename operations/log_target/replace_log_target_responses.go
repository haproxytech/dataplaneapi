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

package log_target

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v4/models"
)

// ReplaceLogTargetOKCode is the HTTP code returned for type ReplaceLogTargetOK
const ReplaceLogTargetOKCode int = 200

/*
ReplaceLogTargetOK Log Target replaced

swagger:response replaceLogTargetOK
*/
type ReplaceLogTargetOK struct {

	/*
	  In: Body
	*/
	Payload *models.LogTarget `json:"body,omitempty"`
}

// NewReplaceLogTargetOK creates ReplaceLogTargetOK with default headers values
func NewReplaceLogTargetOK() *ReplaceLogTargetOK {

	return &ReplaceLogTargetOK{}
}

// WithPayload adds the payload to the replace log target o k response
func (o *ReplaceLogTargetOK) WithPayload(payload *models.LogTarget) *ReplaceLogTargetOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace log target o k response
func (o *ReplaceLogTargetOK) SetPayload(payload *models.LogTarget) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceLogTargetOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// ReplaceLogTargetAcceptedCode is the HTTP code returned for type ReplaceLogTargetAccepted
const ReplaceLogTargetAcceptedCode int = 202

/*
ReplaceLogTargetAccepted Configuration change accepted and reload requested

swagger:response replaceLogTargetAccepted
*/
type ReplaceLogTargetAccepted struct {
	/*ID of the requested reload

	 */
	ReloadID string `json:"Reload-ID"`

	/*
	  In: Body
	*/
	Payload *models.LogTarget `json:"body,omitempty"`
}

// NewReplaceLogTargetAccepted creates ReplaceLogTargetAccepted with default headers values
func NewReplaceLogTargetAccepted() *ReplaceLogTargetAccepted {

	return &ReplaceLogTargetAccepted{}
}

// WithReloadID adds the reloadId to the replace log target accepted response
func (o *ReplaceLogTargetAccepted) WithReloadID(reloadID string) *ReplaceLogTargetAccepted {
	o.ReloadID = reloadID
	return o
}

// SetReloadID sets the reloadId to the replace log target accepted response
func (o *ReplaceLogTargetAccepted) SetReloadID(reloadID string) {
	o.ReloadID = reloadID
}

// WithPayload adds the payload to the replace log target accepted response
func (o *ReplaceLogTargetAccepted) WithPayload(payload *models.LogTarget) *ReplaceLogTargetAccepted {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace log target accepted response
func (o *ReplaceLogTargetAccepted) SetPayload(payload *models.LogTarget) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceLogTargetAccepted) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// ReplaceLogTargetBadRequestCode is the HTTP code returned for type ReplaceLogTargetBadRequest
const ReplaceLogTargetBadRequestCode int = 400

/*
ReplaceLogTargetBadRequest Bad request

swagger:response replaceLogTargetBadRequest
*/
type ReplaceLogTargetBadRequest struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplaceLogTargetBadRequest creates ReplaceLogTargetBadRequest with default headers values
func NewReplaceLogTargetBadRequest() *ReplaceLogTargetBadRequest {

	return &ReplaceLogTargetBadRequest{}
}

// WithConfigurationVersion adds the configurationVersion to the replace log target bad request response
func (o *ReplaceLogTargetBadRequest) WithConfigurationVersion(configurationVersion string) *ReplaceLogTargetBadRequest {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the replace log target bad request response
func (o *ReplaceLogTargetBadRequest) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the replace log target bad request response
func (o *ReplaceLogTargetBadRequest) WithPayload(payload *models.Error) *ReplaceLogTargetBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace log target bad request response
func (o *ReplaceLogTargetBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceLogTargetBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// ReplaceLogTargetNotFoundCode is the HTTP code returned for type ReplaceLogTargetNotFound
const ReplaceLogTargetNotFoundCode int = 404

/*
ReplaceLogTargetNotFound The specified resource was not found

swagger:response replaceLogTargetNotFound
*/
type ReplaceLogTargetNotFound struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplaceLogTargetNotFound creates ReplaceLogTargetNotFound with default headers values
func NewReplaceLogTargetNotFound() *ReplaceLogTargetNotFound {

	return &ReplaceLogTargetNotFound{}
}

// WithConfigurationVersion adds the configurationVersion to the replace log target not found response
func (o *ReplaceLogTargetNotFound) WithConfigurationVersion(configurationVersion string) *ReplaceLogTargetNotFound {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the replace log target not found response
func (o *ReplaceLogTargetNotFound) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the replace log target not found response
func (o *ReplaceLogTargetNotFound) WithPayload(payload *models.Error) *ReplaceLogTargetNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace log target not found response
func (o *ReplaceLogTargetNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceLogTargetNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	// response header Configuration-Version

	configurationVersion := o.ConfigurationVersion
	if configurationVersion != "" {
		rw.Header().Set("Configuration-Version", configurationVersion)
	}

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*
ReplaceLogTargetDefault General Error

swagger:response replaceLogTargetDefault
*/
type ReplaceLogTargetDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplaceLogTargetDefault creates ReplaceLogTargetDefault with default headers values
func NewReplaceLogTargetDefault(code int) *ReplaceLogTargetDefault {
	if code <= 0 {
		code = 500
	}

	return &ReplaceLogTargetDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the replace log target default response
func (o *ReplaceLogTargetDefault) WithStatusCode(code int) *ReplaceLogTargetDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the replace log target default response
func (o *ReplaceLogTargetDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the replace log target default response
func (o *ReplaceLogTargetDefault) WithConfigurationVersion(configurationVersion string) *ReplaceLogTargetDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the replace log target default response
func (o *ReplaceLogTargetDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the replace log target default response
func (o *ReplaceLogTargetDefault) WithPayload(payload *models.Error) *ReplaceLogTargetDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace log target default response
func (o *ReplaceLogTargetDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceLogTargetDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
