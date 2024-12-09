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

// ReplaceAllQUICInitialRuleDefaultsOKCode is the HTTP code returned for type ReplaceAllQUICInitialRuleDefaultsOK
const ReplaceAllQUICInitialRuleDefaultsOKCode int = 200

/*
ReplaceAllQUICInitialRuleDefaultsOK All TTP After Response Rules lines replaced

swagger:response replaceAllQuicInitialRuleDefaultsOK
*/
type ReplaceAllQUICInitialRuleDefaultsOK struct {

	/*
	  In: Body
	*/
	Payload models.QUICInitialRules `json:"body,omitempty"`
}

// NewReplaceAllQUICInitialRuleDefaultsOK creates ReplaceAllQUICInitialRuleDefaultsOK with default headers values
func NewReplaceAllQUICInitialRuleDefaultsOK() *ReplaceAllQUICInitialRuleDefaultsOK {

	return &ReplaceAllQUICInitialRuleDefaultsOK{}
}

// WithPayload adds the payload to the replace all Quic initial rule defaults o k response
func (o *ReplaceAllQUICInitialRuleDefaultsOK) WithPayload(payload models.QUICInitialRules) *ReplaceAllQUICInitialRuleDefaultsOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace all Quic initial rule defaults o k response
func (o *ReplaceAllQUICInitialRuleDefaultsOK) SetPayload(payload models.QUICInitialRules) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceAllQUICInitialRuleDefaultsOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = models.QUICInitialRules{}
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// ReplaceAllQUICInitialRuleDefaultsAcceptedCode is the HTTP code returned for type ReplaceAllQUICInitialRuleDefaultsAccepted
const ReplaceAllQUICInitialRuleDefaultsAcceptedCode int = 202

/*
ReplaceAllQUICInitialRuleDefaultsAccepted Configuration change accepted and reload requested

swagger:response replaceAllQuicInitialRuleDefaultsAccepted
*/
type ReplaceAllQUICInitialRuleDefaultsAccepted struct {
	/*ID of the requested reload

	 */
	ReloadID string `json:"Reload-ID"`

	/*
	  In: Body
	*/
	Payload models.QUICInitialRules `json:"body,omitempty"`
}

// NewReplaceAllQUICInitialRuleDefaultsAccepted creates ReplaceAllQUICInitialRuleDefaultsAccepted with default headers values
func NewReplaceAllQUICInitialRuleDefaultsAccepted() *ReplaceAllQUICInitialRuleDefaultsAccepted {

	return &ReplaceAllQUICInitialRuleDefaultsAccepted{}
}

// WithReloadID adds the reloadId to the replace all Quic initial rule defaults accepted response
func (o *ReplaceAllQUICInitialRuleDefaultsAccepted) WithReloadID(reloadID string) *ReplaceAllQUICInitialRuleDefaultsAccepted {
	o.ReloadID = reloadID
	return o
}

// SetReloadID sets the reloadId to the replace all Quic initial rule defaults accepted response
func (o *ReplaceAllQUICInitialRuleDefaultsAccepted) SetReloadID(reloadID string) {
	o.ReloadID = reloadID
}

// WithPayload adds the payload to the replace all Quic initial rule defaults accepted response
func (o *ReplaceAllQUICInitialRuleDefaultsAccepted) WithPayload(payload models.QUICInitialRules) *ReplaceAllQUICInitialRuleDefaultsAccepted {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace all Quic initial rule defaults accepted response
func (o *ReplaceAllQUICInitialRuleDefaultsAccepted) SetPayload(payload models.QUICInitialRules) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceAllQUICInitialRuleDefaultsAccepted) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	// response header Reload-ID

	reloadID := o.ReloadID
	if reloadID != "" {
		rw.Header().Set("Reload-ID", reloadID)
	}

	rw.WriteHeader(202)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = models.QUICInitialRules{}
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// ReplaceAllQUICInitialRuleDefaultsBadRequestCode is the HTTP code returned for type ReplaceAllQUICInitialRuleDefaultsBadRequest
const ReplaceAllQUICInitialRuleDefaultsBadRequestCode int = 400

/*
ReplaceAllQUICInitialRuleDefaultsBadRequest Bad request

swagger:response replaceAllQuicInitialRuleDefaultsBadRequest
*/
type ReplaceAllQUICInitialRuleDefaultsBadRequest struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplaceAllQUICInitialRuleDefaultsBadRequest creates ReplaceAllQUICInitialRuleDefaultsBadRequest with default headers values
func NewReplaceAllQUICInitialRuleDefaultsBadRequest() *ReplaceAllQUICInitialRuleDefaultsBadRequest {

	return &ReplaceAllQUICInitialRuleDefaultsBadRequest{}
}

// WithConfigurationVersion adds the configurationVersion to the replace all Quic initial rule defaults bad request response
func (o *ReplaceAllQUICInitialRuleDefaultsBadRequest) WithConfigurationVersion(configurationVersion string) *ReplaceAllQUICInitialRuleDefaultsBadRequest {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the replace all Quic initial rule defaults bad request response
func (o *ReplaceAllQUICInitialRuleDefaultsBadRequest) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the replace all Quic initial rule defaults bad request response
func (o *ReplaceAllQUICInitialRuleDefaultsBadRequest) WithPayload(payload *models.Error) *ReplaceAllQUICInitialRuleDefaultsBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace all Quic initial rule defaults bad request response
func (o *ReplaceAllQUICInitialRuleDefaultsBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceAllQUICInitialRuleDefaultsBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
ReplaceAllQUICInitialRuleDefaultsDefault General Error

swagger:response replaceAllQuicInitialRuleDefaultsDefault
*/
type ReplaceAllQUICInitialRuleDefaultsDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplaceAllQUICInitialRuleDefaultsDefault creates ReplaceAllQUICInitialRuleDefaultsDefault with default headers values
func NewReplaceAllQUICInitialRuleDefaultsDefault(code int) *ReplaceAllQUICInitialRuleDefaultsDefault {
	if code <= 0 {
		code = 500
	}

	return &ReplaceAllQUICInitialRuleDefaultsDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the replace all QUIC initial rule defaults default response
func (o *ReplaceAllQUICInitialRuleDefaultsDefault) WithStatusCode(code int) *ReplaceAllQUICInitialRuleDefaultsDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the replace all QUIC initial rule defaults default response
func (o *ReplaceAllQUICInitialRuleDefaultsDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the replace all QUIC initial rule defaults default response
func (o *ReplaceAllQUICInitialRuleDefaultsDefault) WithConfigurationVersion(configurationVersion string) *ReplaceAllQUICInitialRuleDefaultsDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the replace all QUIC initial rule defaults default response
func (o *ReplaceAllQUICInitialRuleDefaultsDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the replace all QUIC initial rule defaults default response
func (o *ReplaceAllQUICInitialRuleDefaultsDefault) WithPayload(payload *models.Error) *ReplaceAllQUICInitialRuleDefaultsDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace all QUIC initial rule defaults default response
func (o *ReplaceAllQUICInitialRuleDefaultsDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceAllQUICInitialRuleDefaultsDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
