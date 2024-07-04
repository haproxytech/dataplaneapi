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

package bind

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v6/models"
)

// ReplaceBindLogForwardOKCode is the HTTP code returned for type ReplaceBindLogForwardOK
const ReplaceBindLogForwardOKCode int = 200

/*
ReplaceBindLogForwardOK Bind replaced

swagger:response replaceBindLogForwardOK
*/
type ReplaceBindLogForwardOK struct {

	/*
	  In: Body
	*/
	Payload *models.Bind `json:"body,omitempty"`
}

// NewReplaceBindLogForwardOK creates ReplaceBindLogForwardOK with default headers values
func NewReplaceBindLogForwardOK() *ReplaceBindLogForwardOK {

	return &ReplaceBindLogForwardOK{}
}

// WithPayload adds the payload to the replace bind log forward o k response
func (o *ReplaceBindLogForwardOK) WithPayload(payload *models.Bind) *ReplaceBindLogForwardOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace bind log forward o k response
func (o *ReplaceBindLogForwardOK) SetPayload(payload *models.Bind) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceBindLogForwardOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// ReplaceBindLogForwardAcceptedCode is the HTTP code returned for type ReplaceBindLogForwardAccepted
const ReplaceBindLogForwardAcceptedCode int = 202

/*
ReplaceBindLogForwardAccepted Configuration change accepted and reload requested

swagger:response replaceBindLogForwardAccepted
*/
type ReplaceBindLogForwardAccepted struct {
	/*ID of the requested reload

	 */
	ReloadID string `json:"Reload-ID"`

	/*
	  In: Body
	*/
	Payload *models.Bind `json:"body,omitempty"`
}

// NewReplaceBindLogForwardAccepted creates ReplaceBindLogForwardAccepted with default headers values
func NewReplaceBindLogForwardAccepted() *ReplaceBindLogForwardAccepted {

	return &ReplaceBindLogForwardAccepted{}
}

// WithReloadID adds the reloadId to the replace bind log forward accepted response
func (o *ReplaceBindLogForwardAccepted) WithReloadID(reloadID string) *ReplaceBindLogForwardAccepted {
	o.ReloadID = reloadID
	return o
}

// SetReloadID sets the reloadId to the replace bind log forward accepted response
func (o *ReplaceBindLogForwardAccepted) SetReloadID(reloadID string) {
	o.ReloadID = reloadID
}

// WithPayload adds the payload to the replace bind log forward accepted response
func (o *ReplaceBindLogForwardAccepted) WithPayload(payload *models.Bind) *ReplaceBindLogForwardAccepted {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace bind log forward accepted response
func (o *ReplaceBindLogForwardAccepted) SetPayload(payload *models.Bind) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceBindLogForwardAccepted) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// ReplaceBindLogForwardBadRequestCode is the HTTP code returned for type ReplaceBindLogForwardBadRequest
const ReplaceBindLogForwardBadRequestCode int = 400

/*
ReplaceBindLogForwardBadRequest Bad request

swagger:response replaceBindLogForwardBadRequest
*/
type ReplaceBindLogForwardBadRequest struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplaceBindLogForwardBadRequest creates ReplaceBindLogForwardBadRequest with default headers values
func NewReplaceBindLogForwardBadRequest() *ReplaceBindLogForwardBadRequest {

	return &ReplaceBindLogForwardBadRequest{}
}

// WithConfigurationVersion adds the configurationVersion to the replace bind log forward bad request response
func (o *ReplaceBindLogForwardBadRequest) WithConfigurationVersion(configurationVersion string) *ReplaceBindLogForwardBadRequest {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the replace bind log forward bad request response
func (o *ReplaceBindLogForwardBadRequest) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the replace bind log forward bad request response
func (o *ReplaceBindLogForwardBadRequest) WithPayload(payload *models.Error) *ReplaceBindLogForwardBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace bind log forward bad request response
func (o *ReplaceBindLogForwardBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceBindLogForwardBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// ReplaceBindLogForwardNotFoundCode is the HTTP code returned for type ReplaceBindLogForwardNotFound
const ReplaceBindLogForwardNotFoundCode int = 404

/*
ReplaceBindLogForwardNotFound The specified resource was not found

swagger:response replaceBindLogForwardNotFound
*/
type ReplaceBindLogForwardNotFound struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplaceBindLogForwardNotFound creates ReplaceBindLogForwardNotFound with default headers values
func NewReplaceBindLogForwardNotFound() *ReplaceBindLogForwardNotFound {

	return &ReplaceBindLogForwardNotFound{}
}

// WithConfigurationVersion adds the configurationVersion to the replace bind log forward not found response
func (o *ReplaceBindLogForwardNotFound) WithConfigurationVersion(configurationVersion string) *ReplaceBindLogForwardNotFound {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the replace bind log forward not found response
func (o *ReplaceBindLogForwardNotFound) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the replace bind log forward not found response
func (o *ReplaceBindLogForwardNotFound) WithPayload(payload *models.Error) *ReplaceBindLogForwardNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace bind log forward not found response
func (o *ReplaceBindLogForwardNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceBindLogForwardNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
ReplaceBindLogForwardDefault General Error

swagger:response replaceBindLogForwardDefault
*/
type ReplaceBindLogForwardDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplaceBindLogForwardDefault creates ReplaceBindLogForwardDefault with default headers values
func NewReplaceBindLogForwardDefault(code int) *ReplaceBindLogForwardDefault {
	if code <= 0 {
		code = 500
	}

	return &ReplaceBindLogForwardDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the replace bind log forward default response
func (o *ReplaceBindLogForwardDefault) WithStatusCode(code int) *ReplaceBindLogForwardDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the replace bind log forward default response
func (o *ReplaceBindLogForwardDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the replace bind log forward default response
func (o *ReplaceBindLogForwardDefault) WithConfigurationVersion(configurationVersion string) *ReplaceBindLogForwardDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the replace bind log forward default response
func (o *ReplaceBindLogForwardDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the replace bind log forward default response
func (o *ReplaceBindLogForwardDefault) WithPayload(payload *models.Error) *ReplaceBindLogForwardDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace bind log forward default response
func (o *ReplaceBindLogForwardDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceBindLogForwardDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
