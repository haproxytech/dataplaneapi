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

// ReplaceTCPRequestRuleFrontendOKCode is the HTTP code returned for type ReplaceTCPRequestRuleFrontendOK
const ReplaceTCPRequestRuleFrontendOKCode int = 200

/*
ReplaceTCPRequestRuleFrontendOK TCP Request Rule replaced

swagger:response replaceTcpRequestRuleFrontendOK
*/
type ReplaceTCPRequestRuleFrontendOK struct {

	/*
	  In: Body
	*/
	Payload *models.TCPRequestRule `json:"body,omitempty"`
}

// NewReplaceTCPRequestRuleFrontendOK creates ReplaceTCPRequestRuleFrontendOK with default headers values
func NewReplaceTCPRequestRuleFrontendOK() *ReplaceTCPRequestRuleFrontendOK {

	return &ReplaceTCPRequestRuleFrontendOK{}
}

// WithPayload adds the payload to the replace Tcp request rule frontend o k response
func (o *ReplaceTCPRequestRuleFrontendOK) WithPayload(payload *models.TCPRequestRule) *ReplaceTCPRequestRuleFrontendOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace Tcp request rule frontend o k response
func (o *ReplaceTCPRequestRuleFrontendOK) SetPayload(payload *models.TCPRequestRule) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceTCPRequestRuleFrontendOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// ReplaceTCPRequestRuleFrontendAcceptedCode is the HTTP code returned for type ReplaceTCPRequestRuleFrontendAccepted
const ReplaceTCPRequestRuleFrontendAcceptedCode int = 202

/*
ReplaceTCPRequestRuleFrontendAccepted Configuration change accepted and reload requested

swagger:response replaceTcpRequestRuleFrontendAccepted
*/
type ReplaceTCPRequestRuleFrontendAccepted struct {
	/*ID of the requested reload

	 */
	ReloadID string `json:"Reload-ID"`

	/*
	  In: Body
	*/
	Payload *models.TCPRequestRule `json:"body,omitempty"`
}

// NewReplaceTCPRequestRuleFrontendAccepted creates ReplaceTCPRequestRuleFrontendAccepted with default headers values
func NewReplaceTCPRequestRuleFrontendAccepted() *ReplaceTCPRequestRuleFrontendAccepted {

	return &ReplaceTCPRequestRuleFrontendAccepted{}
}

// WithReloadID adds the reloadId to the replace Tcp request rule frontend accepted response
func (o *ReplaceTCPRequestRuleFrontendAccepted) WithReloadID(reloadID string) *ReplaceTCPRequestRuleFrontendAccepted {
	o.ReloadID = reloadID
	return o
}

// SetReloadID sets the reloadId to the replace Tcp request rule frontend accepted response
func (o *ReplaceTCPRequestRuleFrontendAccepted) SetReloadID(reloadID string) {
	o.ReloadID = reloadID
}

// WithPayload adds the payload to the replace Tcp request rule frontend accepted response
func (o *ReplaceTCPRequestRuleFrontendAccepted) WithPayload(payload *models.TCPRequestRule) *ReplaceTCPRequestRuleFrontendAccepted {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace Tcp request rule frontend accepted response
func (o *ReplaceTCPRequestRuleFrontendAccepted) SetPayload(payload *models.TCPRequestRule) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceTCPRequestRuleFrontendAccepted) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// ReplaceTCPRequestRuleFrontendBadRequestCode is the HTTP code returned for type ReplaceTCPRequestRuleFrontendBadRequest
const ReplaceTCPRequestRuleFrontendBadRequestCode int = 400

/*
ReplaceTCPRequestRuleFrontendBadRequest Bad request

swagger:response replaceTcpRequestRuleFrontendBadRequest
*/
type ReplaceTCPRequestRuleFrontendBadRequest struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplaceTCPRequestRuleFrontendBadRequest creates ReplaceTCPRequestRuleFrontendBadRequest with default headers values
func NewReplaceTCPRequestRuleFrontendBadRequest() *ReplaceTCPRequestRuleFrontendBadRequest {

	return &ReplaceTCPRequestRuleFrontendBadRequest{}
}

// WithConfigurationVersion adds the configurationVersion to the replace Tcp request rule frontend bad request response
func (o *ReplaceTCPRequestRuleFrontendBadRequest) WithConfigurationVersion(configurationVersion string) *ReplaceTCPRequestRuleFrontendBadRequest {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the replace Tcp request rule frontend bad request response
func (o *ReplaceTCPRequestRuleFrontendBadRequest) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the replace Tcp request rule frontend bad request response
func (o *ReplaceTCPRequestRuleFrontendBadRequest) WithPayload(payload *models.Error) *ReplaceTCPRequestRuleFrontendBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace Tcp request rule frontend bad request response
func (o *ReplaceTCPRequestRuleFrontendBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceTCPRequestRuleFrontendBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// ReplaceTCPRequestRuleFrontendNotFoundCode is the HTTP code returned for type ReplaceTCPRequestRuleFrontendNotFound
const ReplaceTCPRequestRuleFrontendNotFoundCode int = 404

/*
ReplaceTCPRequestRuleFrontendNotFound The specified resource was not found

swagger:response replaceTcpRequestRuleFrontendNotFound
*/
type ReplaceTCPRequestRuleFrontendNotFound struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplaceTCPRequestRuleFrontendNotFound creates ReplaceTCPRequestRuleFrontendNotFound with default headers values
func NewReplaceTCPRequestRuleFrontendNotFound() *ReplaceTCPRequestRuleFrontendNotFound {

	return &ReplaceTCPRequestRuleFrontendNotFound{}
}

// WithConfigurationVersion adds the configurationVersion to the replace Tcp request rule frontend not found response
func (o *ReplaceTCPRequestRuleFrontendNotFound) WithConfigurationVersion(configurationVersion string) *ReplaceTCPRequestRuleFrontendNotFound {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the replace Tcp request rule frontend not found response
func (o *ReplaceTCPRequestRuleFrontendNotFound) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the replace Tcp request rule frontend not found response
func (o *ReplaceTCPRequestRuleFrontendNotFound) WithPayload(payload *models.Error) *ReplaceTCPRequestRuleFrontendNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace Tcp request rule frontend not found response
func (o *ReplaceTCPRequestRuleFrontendNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceTCPRequestRuleFrontendNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
ReplaceTCPRequestRuleFrontendDefault General Error

swagger:response replaceTcpRequestRuleFrontendDefault
*/
type ReplaceTCPRequestRuleFrontendDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplaceTCPRequestRuleFrontendDefault creates ReplaceTCPRequestRuleFrontendDefault with default headers values
func NewReplaceTCPRequestRuleFrontendDefault(code int) *ReplaceTCPRequestRuleFrontendDefault {
	if code <= 0 {
		code = 500
	}

	return &ReplaceTCPRequestRuleFrontendDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the replace TCP request rule frontend default response
func (o *ReplaceTCPRequestRuleFrontendDefault) WithStatusCode(code int) *ReplaceTCPRequestRuleFrontendDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the replace TCP request rule frontend default response
func (o *ReplaceTCPRequestRuleFrontendDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the replace TCP request rule frontend default response
func (o *ReplaceTCPRequestRuleFrontendDefault) WithConfigurationVersion(configurationVersion string) *ReplaceTCPRequestRuleFrontendDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the replace TCP request rule frontend default response
func (o *ReplaceTCPRequestRuleFrontendDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the replace TCP request rule frontend default response
func (o *ReplaceTCPRequestRuleFrontendDefault) WithPayload(payload *models.Error) *ReplaceTCPRequestRuleFrontendDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace TCP request rule frontend default response
func (o *ReplaceTCPRequestRuleFrontendDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceTCPRequestRuleFrontendDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
