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

package tcp_request_rule

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v6/models"
)

// ReplaceAllTCPRequestRuleDefaultsOKCode is the HTTP code returned for type ReplaceAllTCPRequestRuleDefaultsOK
const ReplaceAllTCPRequestRuleDefaultsOKCode int = 200

/*
ReplaceAllTCPRequestRuleDefaultsOK All TCP Request Rule lines replaced

swagger:response replaceAllTcpRequestRuleDefaultsOK
*/
type ReplaceAllTCPRequestRuleDefaultsOK struct {

	/*
	  In: Body
	*/
	Payload models.TCPRequestRules `json:"body,omitempty"`
}

// NewReplaceAllTCPRequestRuleDefaultsOK creates ReplaceAllTCPRequestRuleDefaultsOK with default headers values
func NewReplaceAllTCPRequestRuleDefaultsOK() *ReplaceAllTCPRequestRuleDefaultsOK {

	return &ReplaceAllTCPRequestRuleDefaultsOK{}
}

// WithPayload adds the payload to the replace all Tcp request rule defaults o k response
func (o *ReplaceAllTCPRequestRuleDefaultsOK) WithPayload(payload models.TCPRequestRules) *ReplaceAllTCPRequestRuleDefaultsOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace all Tcp request rule defaults o k response
func (o *ReplaceAllTCPRequestRuleDefaultsOK) SetPayload(payload models.TCPRequestRules) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceAllTCPRequestRuleDefaultsOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = models.TCPRequestRules{}
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// ReplaceAllTCPRequestRuleDefaultsAcceptedCode is the HTTP code returned for type ReplaceAllTCPRequestRuleDefaultsAccepted
const ReplaceAllTCPRequestRuleDefaultsAcceptedCode int = 202

/*
ReplaceAllTCPRequestRuleDefaultsAccepted Configuration change accepted and reload requested

swagger:response replaceAllTcpRequestRuleDefaultsAccepted
*/
type ReplaceAllTCPRequestRuleDefaultsAccepted struct {
	/*ID of the requested reload

	 */
	ReloadID string `json:"Reload-ID"`

	/*
	  In: Body
	*/
	Payload models.TCPRequestRules `json:"body,omitempty"`
}

// NewReplaceAllTCPRequestRuleDefaultsAccepted creates ReplaceAllTCPRequestRuleDefaultsAccepted with default headers values
func NewReplaceAllTCPRequestRuleDefaultsAccepted() *ReplaceAllTCPRequestRuleDefaultsAccepted {

	return &ReplaceAllTCPRequestRuleDefaultsAccepted{}
}

// WithReloadID adds the reloadId to the replace all Tcp request rule defaults accepted response
func (o *ReplaceAllTCPRequestRuleDefaultsAccepted) WithReloadID(reloadID string) *ReplaceAllTCPRequestRuleDefaultsAccepted {
	o.ReloadID = reloadID
	return o
}

// SetReloadID sets the reloadId to the replace all Tcp request rule defaults accepted response
func (o *ReplaceAllTCPRequestRuleDefaultsAccepted) SetReloadID(reloadID string) {
	o.ReloadID = reloadID
}

// WithPayload adds the payload to the replace all Tcp request rule defaults accepted response
func (o *ReplaceAllTCPRequestRuleDefaultsAccepted) WithPayload(payload models.TCPRequestRules) *ReplaceAllTCPRequestRuleDefaultsAccepted {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace all Tcp request rule defaults accepted response
func (o *ReplaceAllTCPRequestRuleDefaultsAccepted) SetPayload(payload models.TCPRequestRules) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceAllTCPRequestRuleDefaultsAccepted) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	// response header Reload-ID

	reloadID := o.ReloadID
	if reloadID != "" {
		rw.Header().Set("Reload-ID", reloadID)
	}

	rw.WriteHeader(202)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = models.TCPRequestRules{}
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// ReplaceAllTCPRequestRuleDefaultsBadRequestCode is the HTTP code returned for type ReplaceAllTCPRequestRuleDefaultsBadRequest
const ReplaceAllTCPRequestRuleDefaultsBadRequestCode int = 400

/*
ReplaceAllTCPRequestRuleDefaultsBadRequest Bad request

swagger:response replaceAllTcpRequestRuleDefaultsBadRequest
*/
type ReplaceAllTCPRequestRuleDefaultsBadRequest struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplaceAllTCPRequestRuleDefaultsBadRequest creates ReplaceAllTCPRequestRuleDefaultsBadRequest with default headers values
func NewReplaceAllTCPRequestRuleDefaultsBadRequest() *ReplaceAllTCPRequestRuleDefaultsBadRequest {

	return &ReplaceAllTCPRequestRuleDefaultsBadRequest{}
}

// WithConfigurationVersion adds the configurationVersion to the replace all Tcp request rule defaults bad request response
func (o *ReplaceAllTCPRequestRuleDefaultsBadRequest) WithConfigurationVersion(configurationVersion string) *ReplaceAllTCPRequestRuleDefaultsBadRequest {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the replace all Tcp request rule defaults bad request response
func (o *ReplaceAllTCPRequestRuleDefaultsBadRequest) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the replace all Tcp request rule defaults bad request response
func (o *ReplaceAllTCPRequestRuleDefaultsBadRequest) WithPayload(payload *models.Error) *ReplaceAllTCPRequestRuleDefaultsBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace all Tcp request rule defaults bad request response
func (o *ReplaceAllTCPRequestRuleDefaultsBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceAllTCPRequestRuleDefaultsBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
ReplaceAllTCPRequestRuleDefaultsDefault General Error

swagger:response replaceAllTcpRequestRuleDefaultsDefault
*/
type ReplaceAllTCPRequestRuleDefaultsDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplaceAllTCPRequestRuleDefaultsDefault creates ReplaceAllTCPRequestRuleDefaultsDefault with default headers values
func NewReplaceAllTCPRequestRuleDefaultsDefault(code int) *ReplaceAllTCPRequestRuleDefaultsDefault {
	if code <= 0 {
		code = 500
	}

	return &ReplaceAllTCPRequestRuleDefaultsDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the replace all TCP request rule defaults default response
func (o *ReplaceAllTCPRequestRuleDefaultsDefault) WithStatusCode(code int) *ReplaceAllTCPRequestRuleDefaultsDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the replace all TCP request rule defaults default response
func (o *ReplaceAllTCPRequestRuleDefaultsDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the replace all TCP request rule defaults default response
func (o *ReplaceAllTCPRequestRuleDefaultsDefault) WithConfigurationVersion(configurationVersion string) *ReplaceAllTCPRequestRuleDefaultsDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the replace all TCP request rule defaults default response
func (o *ReplaceAllTCPRequestRuleDefaultsDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the replace all TCP request rule defaults default response
func (o *ReplaceAllTCPRequestRuleDefaultsDefault) WithPayload(payload *models.Error) *ReplaceAllTCPRequestRuleDefaultsDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace all TCP request rule defaults default response
func (o *ReplaceAllTCPRequestRuleDefaultsDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceAllTCPRequestRuleDefaultsDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
