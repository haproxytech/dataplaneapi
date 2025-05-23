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

package http_request_rule

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v6/models"
)

// ReplaceHTTPRequestRuleDefaultsOKCode is the HTTP code returned for type ReplaceHTTPRequestRuleDefaultsOK
const ReplaceHTTPRequestRuleDefaultsOKCode int = 200

/*
ReplaceHTTPRequestRuleDefaultsOK HTTP Request Rule replaced

swagger:response replaceHttpRequestRuleDefaultsOK
*/
type ReplaceHTTPRequestRuleDefaultsOK struct {

	/*
	  In: Body
	*/
	Payload *models.HTTPRequestRule `json:"body,omitempty"`
}

// NewReplaceHTTPRequestRuleDefaultsOK creates ReplaceHTTPRequestRuleDefaultsOK with default headers values
func NewReplaceHTTPRequestRuleDefaultsOK() *ReplaceHTTPRequestRuleDefaultsOK {

	return &ReplaceHTTPRequestRuleDefaultsOK{}
}

// WithPayload adds the payload to the replace Http request rule defaults o k response
func (o *ReplaceHTTPRequestRuleDefaultsOK) WithPayload(payload *models.HTTPRequestRule) *ReplaceHTTPRequestRuleDefaultsOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace Http request rule defaults o k response
func (o *ReplaceHTTPRequestRuleDefaultsOK) SetPayload(payload *models.HTTPRequestRule) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceHTTPRequestRuleDefaultsOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// ReplaceHTTPRequestRuleDefaultsAcceptedCode is the HTTP code returned for type ReplaceHTTPRequestRuleDefaultsAccepted
const ReplaceHTTPRequestRuleDefaultsAcceptedCode int = 202

/*
ReplaceHTTPRequestRuleDefaultsAccepted Configuration change accepted and reload requested

swagger:response replaceHttpRequestRuleDefaultsAccepted
*/
type ReplaceHTTPRequestRuleDefaultsAccepted struct {
	/*ID of the requested reload

	 */
	ReloadID string `json:"Reload-ID"`

	/*
	  In: Body
	*/
	Payload *models.HTTPRequestRule `json:"body,omitempty"`
}

// NewReplaceHTTPRequestRuleDefaultsAccepted creates ReplaceHTTPRequestRuleDefaultsAccepted with default headers values
func NewReplaceHTTPRequestRuleDefaultsAccepted() *ReplaceHTTPRequestRuleDefaultsAccepted {

	return &ReplaceHTTPRequestRuleDefaultsAccepted{}
}

// WithReloadID adds the reloadId to the replace Http request rule defaults accepted response
func (o *ReplaceHTTPRequestRuleDefaultsAccepted) WithReloadID(reloadID string) *ReplaceHTTPRequestRuleDefaultsAccepted {
	o.ReloadID = reloadID
	return o
}

// SetReloadID sets the reloadId to the replace Http request rule defaults accepted response
func (o *ReplaceHTTPRequestRuleDefaultsAccepted) SetReloadID(reloadID string) {
	o.ReloadID = reloadID
}

// WithPayload adds the payload to the replace Http request rule defaults accepted response
func (o *ReplaceHTTPRequestRuleDefaultsAccepted) WithPayload(payload *models.HTTPRequestRule) *ReplaceHTTPRequestRuleDefaultsAccepted {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace Http request rule defaults accepted response
func (o *ReplaceHTTPRequestRuleDefaultsAccepted) SetPayload(payload *models.HTTPRequestRule) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceHTTPRequestRuleDefaultsAccepted) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// ReplaceHTTPRequestRuleDefaultsBadRequestCode is the HTTP code returned for type ReplaceHTTPRequestRuleDefaultsBadRequest
const ReplaceHTTPRequestRuleDefaultsBadRequestCode int = 400

/*
ReplaceHTTPRequestRuleDefaultsBadRequest Bad request

swagger:response replaceHttpRequestRuleDefaultsBadRequest
*/
type ReplaceHTTPRequestRuleDefaultsBadRequest struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplaceHTTPRequestRuleDefaultsBadRequest creates ReplaceHTTPRequestRuleDefaultsBadRequest with default headers values
func NewReplaceHTTPRequestRuleDefaultsBadRequest() *ReplaceHTTPRequestRuleDefaultsBadRequest {

	return &ReplaceHTTPRequestRuleDefaultsBadRequest{}
}

// WithConfigurationVersion adds the configurationVersion to the replace Http request rule defaults bad request response
func (o *ReplaceHTTPRequestRuleDefaultsBadRequest) WithConfigurationVersion(configurationVersion string) *ReplaceHTTPRequestRuleDefaultsBadRequest {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the replace Http request rule defaults bad request response
func (o *ReplaceHTTPRequestRuleDefaultsBadRequest) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the replace Http request rule defaults bad request response
func (o *ReplaceHTTPRequestRuleDefaultsBadRequest) WithPayload(payload *models.Error) *ReplaceHTTPRequestRuleDefaultsBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace Http request rule defaults bad request response
func (o *ReplaceHTTPRequestRuleDefaultsBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceHTTPRequestRuleDefaultsBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// ReplaceHTTPRequestRuleDefaultsNotFoundCode is the HTTP code returned for type ReplaceHTTPRequestRuleDefaultsNotFound
const ReplaceHTTPRequestRuleDefaultsNotFoundCode int = 404

/*
ReplaceHTTPRequestRuleDefaultsNotFound The specified resource was not found

swagger:response replaceHttpRequestRuleDefaultsNotFound
*/
type ReplaceHTTPRequestRuleDefaultsNotFound struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplaceHTTPRequestRuleDefaultsNotFound creates ReplaceHTTPRequestRuleDefaultsNotFound with default headers values
func NewReplaceHTTPRequestRuleDefaultsNotFound() *ReplaceHTTPRequestRuleDefaultsNotFound {

	return &ReplaceHTTPRequestRuleDefaultsNotFound{}
}

// WithConfigurationVersion adds the configurationVersion to the replace Http request rule defaults not found response
func (o *ReplaceHTTPRequestRuleDefaultsNotFound) WithConfigurationVersion(configurationVersion string) *ReplaceHTTPRequestRuleDefaultsNotFound {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the replace Http request rule defaults not found response
func (o *ReplaceHTTPRequestRuleDefaultsNotFound) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the replace Http request rule defaults not found response
func (o *ReplaceHTTPRequestRuleDefaultsNotFound) WithPayload(payload *models.Error) *ReplaceHTTPRequestRuleDefaultsNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace Http request rule defaults not found response
func (o *ReplaceHTTPRequestRuleDefaultsNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceHTTPRequestRuleDefaultsNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
ReplaceHTTPRequestRuleDefaultsDefault General Error

swagger:response replaceHttpRequestRuleDefaultsDefault
*/
type ReplaceHTTPRequestRuleDefaultsDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplaceHTTPRequestRuleDefaultsDefault creates ReplaceHTTPRequestRuleDefaultsDefault with default headers values
func NewReplaceHTTPRequestRuleDefaultsDefault(code int) *ReplaceHTTPRequestRuleDefaultsDefault {
	if code <= 0 {
		code = 500
	}

	return &ReplaceHTTPRequestRuleDefaultsDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the replace HTTP request rule defaults default response
func (o *ReplaceHTTPRequestRuleDefaultsDefault) WithStatusCode(code int) *ReplaceHTTPRequestRuleDefaultsDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the replace HTTP request rule defaults default response
func (o *ReplaceHTTPRequestRuleDefaultsDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the replace HTTP request rule defaults default response
func (o *ReplaceHTTPRequestRuleDefaultsDefault) WithConfigurationVersion(configurationVersion string) *ReplaceHTTPRequestRuleDefaultsDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the replace HTTP request rule defaults default response
func (o *ReplaceHTTPRequestRuleDefaultsDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the replace HTTP request rule defaults default response
func (o *ReplaceHTTPRequestRuleDefaultsDefault) WithPayload(payload *models.Error) *ReplaceHTTPRequestRuleDefaultsDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace HTTP request rule defaults default response
func (o *ReplaceHTTPRequestRuleDefaultsDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceHTTPRequestRuleDefaultsDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
