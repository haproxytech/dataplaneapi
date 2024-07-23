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

// ReplaceBindPeerOKCode is the HTTP code returned for type ReplaceBindPeerOK
const ReplaceBindPeerOKCode int = 200

/*
ReplaceBindPeerOK Bind replaced

swagger:response replaceBindPeerOK
*/
type ReplaceBindPeerOK struct {

	/*
	  In: Body
	*/
	Payload *models.Bind `json:"body,omitempty"`
}

// NewReplaceBindPeerOK creates ReplaceBindPeerOK with default headers values
func NewReplaceBindPeerOK() *ReplaceBindPeerOK {

	return &ReplaceBindPeerOK{}
}

// WithPayload adds the payload to the replace bind peer o k response
func (o *ReplaceBindPeerOK) WithPayload(payload *models.Bind) *ReplaceBindPeerOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace bind peer o k response
func (o *ReplaceBindPeerOK) SetPayload(payload *models.Bind) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceBindPeerOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// ReplaceBindPeerAcceptedCode is the HTTP code returned for type ReplaceBindPeerAccepted
const ReplaceBindPeerAcceptedCode int = 202

/*
ReplaceBindPeerAccepted Configuration change accepted and reload requested

swagger:response replaceBindPeerAccepted
*/
type ReplaceBindPeerAccepted struct {
	/*ID of the requested reload

	 */
	ReloadID string `json:"Reload-ID"`

	/*
	  In: Body
	*/
	Payload *models.Bind `json:"body,omitempty"`
}

// NewReplaceBindPeerAccepted creates ReplaceBindPeerAccepted with default headers values
func NewReplaceBindPeerAccepted() *ReplaceBindPeerAccepted {

	return &ReplaceBindPeerAccepted{}
}

// WithReloadID adds the reloadId to the replace bind peer accepted response
func (o *ReplaceBindPeerAccepted) WithReloadID(reloadID string) *ReplaceBindPeerAccepted {
	o.ReloadID = reloadID
	return o
}

// SetReloadID sets the reloadId to the replace bind peer accepted response
func (o *ReplaceBindPeerAccepted) SetReloadID(reloadID string) {
	o.ReloadID = reloadID
}

// WithPayload adds the payload to the replace bind peer accepted response
func (o *ReplaceBindPeerAccepted) WithPayload(payload *models.Bind) *ReplaceBindPeerAccepted {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace bind peer accepted response
func (o *ReplaceBindPeerAccepted) SetPayload(payload *models.Bind) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceBindPeerAccepted) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// ReplaceBindPeerBadRequestCode is the HTTP code returned for type ReplaceBindPeerBadRequest
const ReplaceBindPeerBadRequestCode int = 400

/*
ReplaceBindPeerBadRequest Bad request

swagger:response replaceBindPeerBadRequest
*/
type ReplaceBindPeerBadRequest struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplaceBindPeerBadRequest creates ReplaceBindPeerBadRequest with default headers values
func NewReplaceBindPeerBadRequest() *ReplaceBindPeerBadRequest {

	return &ReplaceBindPeerBadRequest{}
}

// WithConfigurationVersion adds the configurationVersion to the replace bind peer bad request response
func (o *ReplaceBindPeerBadRequest) WithConfigurationVersion(configurationVersion string) *ReplaceBindPeerBadRequest {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the replace bind peer bad request response
func (o *ReplaceBindPeerBadRequest) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the replace bind peer bad request response
func (o *ReplaceBindPeerBadRequest) WithPayload(payload *models.Error) *ReplaceBindPeerBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace bind peer bad request response
func (o *ReplaceBindPeerBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceBindPeerBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// ReplaceBindPeerNotFoundCode is the HTTP code returned for type ReplaceBindPeerNotFound
const ReplaceBindPeerNotFoundCode int = 404

/*
ReplaceBindPeerNotFound The specified resource was not found

swagger:response replaceBindPeerNotFound
*/
type ReplaceBindPeerNotFound struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplaceBindPeerNotFound creates ReplaceBindPeerNotFound with default headers values
func NewReplaceBindPeerNotFound() *ReplaceBindPeerNotFound {

	return &ReplaceBindPeerNotFound{}
}

// WithConfigurationVersion adds the configurationVersion to the replace bind peer not found response
func (o *ReplaceBindPeerNotFound) WithConfigurationVersion(configurationVersion string) *ReplaceBindPeerNotFound {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the replace bind peer not found response
func (o *ReplaceBindPeerNotFound) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the replace bind peer not found response
func (o *ReplaceBindPeerNotFound) WithPayload(payload *models.Error) *ReplaceBindPeerNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace bind peer not found response
func (o *ReplaceBindPeerNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceBindPeerNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
ReplaceBindPeerDefault General Error

swagger:response replaceBindPeerDefault
*/
type ReplaceBindPeerDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplaceBindPeerDefault creates ReplaceBindPeerDefault with default headers values
func NewReplaceBindPeerDefault(code int) *ReplaceBindPeerDefault {
	if code <= 0 {
		code = 500
	}

	return &ReplaceBindPeerDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the replace bind peer default response
func (o *ReplaceBindPeerDefault) WithStatusCode(code int) *ReplaceBindPeerDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the replace bind peer default response
func (o *ReplaceBindPeerDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the replace bind peer default response
func (o *ReplaceBindPeerDefault) WithConfigurationVersion(configurationVersion string) *ReplaceBindPeerDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the replace bind peer default response
func (o *ReplaceBindPeerDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the replace bind peer default response
func (o *ReplaceBindPeerDefault) WithPayload(payload *models.Error) *ReplaceBindPeerDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace bind peer default response
func (o *ReplaceBindPeerDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceBindPeerDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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