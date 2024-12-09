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

package quic_initial_rule

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v6/models"
)

// ReplaceQUICInitialRuleFrontendOKCode is the HTTP code returned for type ReplaceQUICInitialRuleFrontendOK
const ReplaceQUICInitialRuleFrontendOKCode int = 200

/*
ReplaceQUICInitialRuleFrontendOK QUIC Initial Rule replaced

swagger:response replaceQuicInitialRuleFrontendOK
*/
type ReplaceQUICInitialRuleFrontendOK struct {

	/*
	  In: Body
	*/
	Payload *models.QUICInitialRule `json:"body,omitempty"`
}

// NewReplaceQUICInitialRuleFrontendOK creates ReplaceQUICInitialRuleFrontendOK with default headers values
func NewReplaceQUICInitialRuleFrontendOK() *ReplaceQUICInitialRuleFrontendOK {

	return &ReplaceQUICInitialRuleFrontendOK{}
}

// WithPayload adds the payload to the replace Quic initial rule frontend o k response
func (o *ReplaceQUICInitialRuleFrontendOK) WithPayload(payload *models.QUICInitialRule) *ReplaceQUICInitialRuleFrontendOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace Quic initial rule frontend o k response
func (o *ReplaceQUICInitialRuleFrontendOK) SetPayload(payload *models.QUICInitialRule) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceQUICInitialRuleFrontendOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// ReplaceQUICInitialRuleFrontendAcceptedCode is the HTTP code returned for type ReplaceQUICInitialRuleFrontendAccepted
const ReplaceQUICInitialRuleFrontendAcceptedCode int = 202

/*
ReplaceQUICInitialRuleFrontendAccepted Configuration change accepted and reload requested

swagger:response replaceQuicInitialRuleFrontendAccepted
*/
type ReplaceQUICInitialRuleFrontendAccepted struct {
	/*ID of the requested reload

	 */
	ReloadID string `json:"Reload-ID"`

	/*
	  In: Body
	*/
	Payload *models.QUICInitialRule `json:"body,omitempty"`
}

// NewReplaceQUICInitialRuleFrontendAccepted creates ReplaceQUICInitialRuleFrontendAccepted with default headers values
func NewReplaceQUICInitialRuleFrontendAccepted() *ReplaceQUICInitialRuleFrontendAccepted {

	return &ReplaceQUICInitialRuleFrontendAccepted{}
}

// WithReloadID adds the reloadId to the replace Quic initial rule frontend accepted response
func (o *ReplaceQUICInitialRuleFrontendAccepted) WithReloadID(reloadID string) *ReplaceQUICInitialRuleFrontendAccepted {
	o.ReloadID = reloadID
	return o
}

// SetReloadID sets the reloadId to the replace Quic initial rule frontend accepted response
func (o *ReplaceQUICInitialRuleFrontendAccepted) SetReloadID(reloadID string) {
	o.ReloadID = reloadID
}

// WithPayload adds the payload to the replace Quic initial rule frontend accepted response
func (o *ReplaceQUICInitialRuleFrontendAccepted) WithPayload(payload *models.QUICInitialRule) *ReplaceQUICInitialRuleFrontendAccepted {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace Quic initial rule frontend accepted response
func (o *ReplaceQUICInitialRuleFrontendAccepted) SetPayload(payload *models.QUICInitialRule) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceQUICInitialRuleFrontendAccepted) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// ReplaceQUICInitialRuleFrontendBadRequestCode is the HTTP code returned for type ReplaceQUICInitialRuleFrontendBadRequest
const ReplaceQUICInitialRuleFrontendBadRequestCode int = 400

/*
ReplaceQUICInitialRuleFrontendBadRequest Bad request

swagger:response replaceQuicInitialRuleFrontendBadRequest
*/
type ReplaceQUICInitialRuleFrontendBadRequest struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplaceQUICInitialRuleFrontendBadRequest creates ReplaceQUICInitialRuleFrontendBadRequest with default headers values
func NewReplaceQUICInitialRuleFrontendBadRequest() *ReplaceQUICInitialRuleFrontendBadRequest {

	return &ReplaceQUICInitialRuleFrontendBadRequest{}
}

// WithConfigurationVersion adds the configurationVersion to the replace Quic initial rule frontend bad request response
func (o *ReplaceQUICInitialRuleFrontendBadRequest) WithConfigurationVersion(configurationVersion string) *ReplaceQUICInitialRuleFrontendBadRequest {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the replace Quic initial rule frontend bad request response
func (o *ReplaceQUICInitialRuleFrontendBadRequest) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the replace Quic initial rule frontend bad request response
func (o *ReplaceQUICInitialRuleFrontendBadRequest) WithPayload(payload *models.Error) *ReplaceQUICInitialRuleFrontendBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace Quic initial rule frontend bad request response
func (o *ReplaceQUICInitialRuleFrontendBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceQUICInitialRuleFrontendBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// ReplaceQUICInitialRuleFrontendNotFoundCode is the HTTP code returned for type ReplaceQUICInitialRuleFrontendNotFound
const ReplaceQUICInitialRuleFrontendNotFoundCode int = 404

/*
ReplaceQUICInitialRuleFrontendNotFound The specified resource was not found

swagger:response replaceQuicInitialRuleFrontendNotFound
*/
type ReplaceQUICInitialRuleFrontendNotFound struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplaceQUICInitialRuleFrontendNotFound creates ReplaceQUICInitialRuleFrontendNotFound with default headers values
func NewReplaceQUICInitialRuleFrontendNotFound() *ReplaceQUICInitialRuleFrontendNotFound {

	return &ReplaceQUICInitialRuleFrontendNotFound{}
}

// WithConfigurationVersion adds the configurationVersion to the replace Quic initial rule frontend not found response
func (o *ReplaceQUICInitialRuleFrontendNotFound) WithConfigurationVersion(configurationVersion string) *ReplaceQUICInitialRuleFrontendNotFound {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the replace Quic initial rule frontend not found response
func (o *ReplaceQUICInitialRuleFrontendNotFound) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the replace Quic initial rule frontend not found response
func (o *ReplaceQUICInitialRuleFrontendNotFound) WithPayload(payload *models.Error) *ReplaceQUICInitialRuleFrontendNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace Quic initial rule frontend not found response
func (o *ReplaceQUICInitialRuleFrontendNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceQUICInitialRuleFrontendNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
ReplaceQUICInitialRuleFrontendDefault General Error

swagger:response replaceQuicInitialRuleFrontendDefault
*/
type ReplaceQUICInitialRuleFrontendDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplaceQUICInitialRuleFrontendDefault creates ReplaceQUICInitialRuleFrontendDefault with default headers values
func NewReplaceQUICInitialRuleFrontendDefault(code int) *ReplaceQUICInitialRuleFrontendDefault {
	if code <= 0 {
		code = 500
	}

	return &ReplaceQUICInitialRuleFrontendDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the replace QUIC initial rule frontend default response
func (o *ReplaceQUICInitialRuleFrontendDefault) WithStatusCode(code int) *ReplaceQUICInitialRuleFrontendDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the replace QUIC initial rule frontend default response
func (o *ReplaceQUICInitialRuleFrontendDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the replace QUIC initial rule frontend default response
func (o *ReplaceQUICInitialRuleFrontendDefault) WithConfigurationVersion(configurationVersion string) *ReplaceQUICInitialRuleFrontendDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the replace QUIC initial rule frontend default response
func (o *ReplaceQUICInitialRuleFrontendDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the replace QUIC initial rule frontend default response
func (o *ReplaceQUICInitialRuleFrontendDefault) WithPayload(payload *models.Error) *ReplaceQUICInitialRuleFrontendDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace QUIC initial rule frontend default response
func (o *ReplaceQUICInitialRuleFrontendDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceQUICInitialRuleFrontendDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
