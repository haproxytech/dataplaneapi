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

package http_error_rule

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v6/models"
)

// ReplaceAllHTTPErrorRuleDefaultsOKCode is the HTTP code returned for type ReplaceAllHTTPErrorRuleDefaultsOK
const ReplaceAllHTTPErrorRuleDefaultsOKCode int = 200

/*
ReplaceAllHTTPErrorRuleDefaultsOK All HTTP Error Rules lines replaced

swagger:response replaceAllHttpErrorRuleDefaultsOK
*/
type ReplaceAllHTTPErrorRuleDefaultsOK struct {

	/*
	  In: Body
	*/
	Payload models.HTTPErrorRules `json:"body,omitempty"`
}

// NewReplaceAllHTTPErrorRuleDefaultsOK creates ReplaceAllHTTPErrorRuleDefaultsOK with default headers values
func NewReplaceAllHTTPErrorRuleDefaultsOK() *ReplaceAllHTTPErrorRuleDefaultsOK {

	return &ReplaceAllHTTPErrorRuleDefaultsOK{}
}

// WithPayload adds the payload to the replace all Http error rule defaults o k response
func (o *ReplaceAllHTTPErrorRuleDefaultsOK) WithPayload(payload models.HTTPErrorRules) *ReplaceAllHTTPErrorRuleDefaultsOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace all Http error rule defaults o k response
func (o *ReplaceAllHTTPErrorRuleDefaultsOK) SetPayload(payload models.HTTPErrorRules) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceAllHTTPErrorRuleDefaultsOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = models.HTTPErrorRules{}
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// ReplaceAllHTTPErrorRuleDefaultsAcceptedCode is the HTTP code returned for type ReplaceAllHTTPErrorRuleDefaultsAccepted
const ReplaceAllHTTPErrorRuleDefaultsAcceptedCode int = 202

/*
ReplaceAllHTTPErrorRuleDefaultsAccepted Configuration change accepted and reload requested

swagger:response replaceAllHttpErrorRuleDefaultsAccepted
*/
type ReplaceAllHTTPErrorRuleDefaultsAccepted struct {
	/*ID of the requested reload

	 */
	ReloadID string `json:"Reload-ID"`

	/*
	  In: Body
	*/
	Payload models.HTTPErrorRules `json:"body,omitempty"`
}

// NewReplaceAllHTTPErrorRuleDefaultsAccepted creates ReplaceAllHTTPErrorRuleDefaultsAccepted with default headers values
func NewReplaceAllHTTPErrorRuleDefaultsAccepted() *ReplaceAllHTTPErrorRuleDefaultsAccepted {

	return &ReplaceAllHTTPErrorRuleDefaultsAccepted{}
}

// WithReloadID adds the reloadId to the replace all Http error rule defaults accepted response
func (o *ReplaceAllHTTPErrorRuleDefaultsAccepted) WithReloadID(reloadID string) *ReplaceAllHTTPErrorRuleDefaultsAccepted {
	o.ReloadID = reloadID
	return o
}

// SetReloadID sets the reloadId to the replace all Http error rule defaults accepted response
func (o *ReplaceAllHTTPErrorRuleDefaultsAccepted) SetReloadID(reloadID string) {
	o.ReloadID = reloadID
}

// WithPayload adds the payload to the replace all Http error rule defaults accepted response
func (o *ReplaceAllHTTPErrorRuleDefaultsAccepted) WithPayload(payload models.HTTPErrorRules) *ReplaceAllHTTPErrorRuleDefaultsAccepted {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace all Http error rule defaults accepted response
func (o *ReplaceAllHTTPErrorRuleDefaultsAccepted) SetPayload(payload models.HTTPErrorRules) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceAllHTTPErrorRuleDefaultsAccepted) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	// response header Reload-ID

	reloadID := o.ReloadID
	if reloadID != "" {
		rw.Header().Set("Reload-ID", reloadID)
	}

	rw.WriteHeader(202)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = models.HTTPErrorRules{}
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// ReplaceAllHTTPErrorRuleDefaultsBadRequestCode is the HTTP code returned for type ReplaceAllHTTPErrorRuleDefaultsBadRequest
const ReplaceAllHTTPErrorRuleDefaultsBadRequestCode int = 400

/*
ReplaceAllHTTPErrorRuleDefaultsBadRequest Bad request

swagger:response replaceAllHttpErrorRuleDefaultsBadRequest
*/
type ReplaceAllHTTPErrorRuleDefaultsBadRequest struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplaceAllHTTPErrorRuleDefaultsBadRequest creates ReplaceAllHTTPErrorRuleDefaultsBadRequest with default headers values
func NewReplaceAllHTTPErrorRuleDefaultsBadRequest() *ReplaceAllHTTPErrorRuleDefaultsBadRequest {

	return &ReplaceAllHTTPErrorRuleDefaultsBadRequest{}
}

// WithConfigurationVersion adds the configurationVersion to the replace all Http error rule defaults bad request response
func (o *ReplaceAllHTTPErrorRuleDefaultsBadRequest) WithConfigurationVersion(configurationVersion string) *ReplaceAllHTTPErrorRuleDefaultsBadRequest {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the replace all Http error rule defaults bad request response
func (o *ReplaceAllHTTPErrorRuleDefaultsBadRequest) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the replace all Http error rule defaults bad request response
func (o *ReplaceAllHTTPErrorRuleDefaultsBadRequest) WithPayload(payload *models.Error) *ReplaceAllHTTPErrorRuleDefaultsBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace all Http error rule defaults bad request response
func (o *ReplaceAllHTTPErrorRuleDefaultsBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceAllHTTPErrorRuleDefaultsBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
ReplaceAllHTTPErrorRuleDefaultsDefault General Error

swagger:response replaceAllHttpErrorRuleDefaultsDefault
*/
type ReplaceAllHTTPErrorRuleDefaultsDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplaceAllHTTPErrorRuleDefaultsDefault creates ReplaceAllHTTPErrorRuleDefaultsDefault with default headers values
func NewReplaceAllHTTPErrorRuleDefaultsDefault(code int) *ReplaceAllHTTPErrorRuleDefaultsDefault {
	if code <= 0 {
		code = 500
	}

	return &ReplaceAllHTTPErrorRuleDefaultsDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the replace all HTTP error rule defaults default response
func (o *ReplaceAllHTTPErrorRuleDefaultsDefault) WithStatusCode(code int) *ReplaceAllHTTPErrorRuleDefaultsDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the replace all HTTP error rule defaults default response
func (o *ReplaceAllHTTPErrorRuleDefaultsDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the replace all HTTP error rule defaults default response
func (o *ReplaceAllHTTPErrorRuleDefaultsDefault) WithConfigurationVersion(configurationVersion string) *ReplaceAllHTTPErrorRuleDefaultsDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the replace all HTTP error rule defaults default response
func (o *ReplaceAllHTTPErrorRuleDefaultsDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the replace all HTTP error rule defaults default response
func (o *ReplaceAllHTTPErrorRuleDefaultsDefault) WithPayload(payload *models.Error) *ReplaceAllHTTPErrorRuleDefaultsDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace all HTTP error rule defaults default response
func (o *ReplaceAllHTTPErrorRuleDefaultsDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceAllHTTPErrorRuleDefaultsDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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