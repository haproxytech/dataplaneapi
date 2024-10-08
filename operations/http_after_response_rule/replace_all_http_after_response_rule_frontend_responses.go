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

package http_after_response_rule

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v6/models"
)

// ReplaceAllHTTPAfterResponseRuleFrontendOKCode is the HTTP code returned for type ReplaceAllHTTPAfterResponseRuleFrontendOK
const ReplaceAllHTTPAfterResponseRuleFrontendOKCode int = 200

/*
ReplaceAllHTTPAfterResponseRuleFrontendOK All TTP After Response Rules lines replaced

swagger:response replaceAllHttpAfterResponseRuleFrontendOK
*/
type ReplaceAllHTTPAfterResponseRuleFrontendOK struct {

	/*
	  In: Body
	*/
	Payload models.HTTPAfterResponseRules `json:"body,omitempty"`
}

// NewReplaceAllHTTPAfterResponseRuleFrontendOK creates ReplaceAllHTTPAfterResponseRuleFrontendOK with default headers values
func NewReplaceAllHTTPAfterResponseRuleFrontendOK() *ReplaceAllHTTPAfterResponseRuleFrontendOK {

	return &ReplaceAllHTTPAfterResponseRuleFrontendOK{}
}

// WithPayload adds the payload to the replace all Http after response rule frontend o k response
func (o *ReplaceAllHTTPAfterResponseRuleFrontendOK) WithPayload(payload models.HTTPAfterResponseRules) *ReplaceAllHTTPAfterResponseRuleFrontendOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace all Http after response rule frontend o k response
func (o *ReplaceAllHTTPAfterResponseRuleFrontendOK) SetPayload(payload models.HTTPAfterResponseRules) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceAllHTTPAfterResponseRuleFrontendOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = models.HTTPAfterResponseRules{}
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// ReplaceAllHTTPAfterResponseRuleFrontendAcceptedCode is the HTTP code returned for type ReplaceAllHTTPAfterResponseRuleFrontendAccepted
const ReplaceAllHTTPAfterResponseRuleFrontendAcceptedCode int = 202

/*
ReplaceAllHTTPAfterResponseRuleFrontendAccepted Configuration change accepted and reload requested

swagger:response replaceAllHttpAfterResponseRuleFrontendAccepted
*/
type ReplaceAllHTTPAfterResponseRuleFrontendAccepted struct {
	/*ID of the requested reload

	 */
	ReloadID string `json:"Reload-ID"`

	/*
	  In: Body
	*/
	Payload models.HTTPAfterResponseRules `json:"body,omitempty"`
}

// NewReplaceAllHTTPAfterResponseRuleFrontendAccepted creates ReplaceAllHTTPAfterResponseRuleFrontendAccepted with default headers values
func NewReplaceAllHTTPAfterResponseRuleFrontendAccepted() *ReplaceAllHTTPAfterResponseRuleFrontendAccepted {

	return &ReplaceAllHTTPAfterResponseRuleFrontendAccepted{}
}

// WithReloadID adds the reloadId to the replace all Http after response rule frontend accepted response
func (o *ReplaceAllHTTPAfterResponseRuleFrontendAccepted) WithReloadID(reloadID string) *ReplaceAllHTTPAfterResponseRuleFrontendAccepted {
	o.ReloadID = reloadID
	return o
}

// SetReloadID sets the reloadId to the replace all Http after response rule frontend accepted response
func (o *ReplaceAllHTTPAfterResponseRuleFrontendAccepted) SetReloadID(reloadID string) {
	o.ReloadID = reloadID
}

// WithPayload adds the payload to the replace all Http after response rule frontend accepted response
func (o *ReplaceAllHTTPAfterResponseRuleFrontendAccepted) WithPayload(payload models.HTTPAfterResponseRules) *ReplaceAllHTTPAfterResponseRuleFrontendAccepted {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace all Http after response rule frontend accepted response
func (o *ReplaceAllHTTPAfterResponseRuleFrontendAccepted) SetPayload(payload models.HTTPAfterResponseRules) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceAllHTTPAfterResponseRuleFrontendAccepted) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	// response header Reload-ID

	reloadID := o.ReloadID
	if reloadID != "" {
		rw.Header().Set("Reload-ID", reloadID)
	}

	rw.WriteHeader(202)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = models.HTTPAfterResponseRules{}
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// ReplaceAllHTTPAfterResponseRuleFrontendBadRequestCode is the HTTP code returned for type ReplaceAllHTTPAfterResponseRuleFrontendBadRequest
const ReplaceAllHTTPAfterResponseRuleFrontendBadRequestCode int = 400

/*
ReplaceAllHTTPAfterResponseRuleFrontendBadRequest Bad request

swagger:response replaceAllHttpAfterResponseRuleFrontendBadRequest
*/
type ReplaceAllHTTPAfterResponseRuleFrontendBadRequest struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplaceAllHTTPAfterResponseRuleFrontendBadRequest creates ReplaceAllHTTPAfterResponseRuleFrontendBadRequest with default headers values
func NewReplaceAllHTTPAfterResponseRuleFrontendBadRequest() *ReplaceAllHTTPAfterResponseRuleFrontendBadRequest {

	return &ReplaceAllHTTPAfterResponseRuleFrontendBadRequest{}
}

// WithConfigurationVersion adds the configurationVersion to the replace all Http after response rule frontend bad request response
func (o *ReplaceAllHTTPAfterResponseRuleFrontendBadRequest) WithConfigurationVersion(configurationVersion string) *ReplaceAllHTTPAfterResponseRuleFrontendBadRequest {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the replace all Http after response rule frontend bad request response
func (o *ReplaceAllHTTPAfterResponseRuleFrontendBadRequest) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the replace all Http after response rule frontend bad request response
func (o *ReplaceAllHTTPAfterResponseRuleFrontendBadRequest) WithPayload(payload *models.Error) *ReplaceAllHTTPAfterResponseRuleFrontendBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace all Http after response rule frontend bad request response
func (o *ReplaceAllHTTPAfterResponseRuleFrontendBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceAllHTTPAfterResponseRuleFrontendBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
ReplaceAllHTTPAfterResponseRuleFrontendDefault General Error

swagger:response replaceAllHttpAfterResponseRuleFrontendDefault
*/
type ReplaceAllHTTPAfterResponseRuleFrontendDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplaceAllHTTPAfterResponseRuleFrontendDefault creates ReplaceAllHTTPAfterResponseRuleFrontendDefault with default headers values
func NewReplaceAllHTTPAfterResponseRuleFrontendDefault(code int) *ReplaceAllHTTPAfterResponseRuleFrontendDefault {
	if code <= 0 {
		code = 500
	}

	return &ReplaceAllHTTPAfterResponseRuleFrontendDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the replace all HTTP after response rule frontend default response
func (o *ReplaceAllHTTPAfterResponseRuleFrontendDefault) WithStatusCode(code int) *ReplaceAllHTTPAfterResponseRuleFrontendDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the replace all HTTP after response rule frontend default response
func (o *ReplaceAllHTTPAfterResponseRuleFrontendDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the replace all HTTP after response rule frontend default response
func (o *ReplaceAllHTTPAfterResponseRuleFrontendDefault) WithConfigurationVersion(configurationVersion string) *ReplaceAllHTTPAfterResponseRuleFrontendDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the replace all HTTP after response rule frontend default response
func (o *ReplaceAllHTTPAfterResponseRuleFrontendDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the replace all HTTP after response rule frontend default response
func (o *ReplaceAllHTTPAfterResponseRuleFrontendDefault) WithPayload(payload *models.Error) *ReplaceAllHTTPAfterResponseRuleFrontendDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace all HTTP after response rule frontend default response
func (o *ReplaceAllHTTPAfterResponseRuleFrontendDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceAllHTTPAfterResponseRuleFrontendDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
