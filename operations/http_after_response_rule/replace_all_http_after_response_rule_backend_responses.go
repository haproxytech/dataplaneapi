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

// ReplaceAllHTTPAfterResponseRuleBackendOKCode is the HTTP code returned for type ReplaceAllHTTPAfterResponseRuleBackendOK
const ReplaceAllHTTPAfterResponseRuleBackendOKCode int = 200

/*
ReplaceAllHTTPAfterResponseRuleBackendOK All TTP After Response Rules lines replaced

swagger:response replaceAllHttpAfterResponseRuleBackendOK
*/
type ReplaceAllHTTPAfterResponseRuleBackendOK struct {

	/*
	  In: Body
	*/
	Payload models.HTTPAfterResponseRules `json:"body,omitempty"`
}

// NewReplaceAllHTTPAfterResponseRuleBackendOK creates ReplaceAllHTTPAfterResponseRuleBackendOK with default headers values
func NewReplaceAllHTTPAfterResponseRuleBackendOK() *ReplaceAllHTTPAfterResponseRuleBackendOK {

	return &ReplaceAllHTTPAfterResponseRuleBackendOK{}
}

// WithPayload adds the payload to the replace all Http after response rule backend o k response
func (o *ReplaceAllHTTPAfterResponseRuleBackendOK) WithPayload(payload models.HTTPAfterResponseRules) *ReplaceAllHTTPAfterResponseRuleBackendOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace all Http after response rule backend o k response
func (o *ReplaceAllHTTPAfterResponseRuleBackendOK) SetPayload(payload models.HTTPAfterResponseRules) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceAllHTTPAfterResponseRuleBackendOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// ReplaceAllHTTPAfterResponseRuleBackendAcceptedCode is the HTTP code returned for type ReplaceAllHTTPAfterResponseRuleBackendAccepted
const ReplaceAllHTTPAfterResponseRuleBackendAcceptedCode int = 202

/*
ReplaceAllHTTPAfterResponseRuleBackendAccepted Configuration change accepted and reload requested

swagger:response replaceAllHttpAfterResponseRuleBackendAccepted
*/
type ReplaceAllHTTPAfterResponseRuleBackendAccepted struct {
	/*ID of the requested reload

	 */
	ReloadID string `json:"Reload-ID"`

	/*
	  In: Body
	*/
	Payload models.HTTPAfterResponseRules `json:"body,omitempty"`
}

// NewReplaceAllHTTPAfterResponseRuleBackendAccepted creates ReplaceAllHTTPAfterResponseRuleBackendAccepted with default headers values
func NewReplaceAllHTTPAfterResponseRuleBackendAccepted() *ReplaceAllHTTPAfterResponseRuleBackendAccepted {

	return &ReplaceAllHTTPAfterResponseRuleBackendAccepted{}
}

// WithReloadID adds the reloadId to the replace all Http after response rule backend accepted response
func (o *ReplaceAllHTTPAfterResponseRuleBackendAccepted) WithReloadID(reloadID string) *ReplaceAllHTTPAfterResponseRuleBackendAccepted {
	o.ReloadID = reloadID
	return o
}

// SetReloadID sets the reloadId to the replace all Http after response rule backend accepted response
func (o *ReplaceAllHTTPAfterResponseRuleBackendAccepted) SetReloadID(reloadID string) {
	o.ReloadID = reloadID
}

// WithPayload adds the payload to the replace all Http after response rule backend accepted response
func (o *ReplaceAllHTTPAfterResponseRuleBackendAccepted) WithPayload(payload models.HTTPAfterResponseRules) *ReplaceAllHTTPAfterResponseRuleBackendAccepted {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace all Http after response rule backend accepted response
func (o *ReplaceAllHTTPAfterResponseRuleBackendAccepted) SetPayload(payload models.HTTPAfterResponseRules) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceAllHTTPAfterResponseRuleBackendAccepted) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// ReplaceAllHTTPAfterResponseRuleBackendBadRequestCode is the HTTP code returned for type ReplaceAllHTTPAfterResponseRuleBackendBadRequest
const ReplaceAllHTTPAfterResponseRuleBackendBadRequestCode int = 400

/*
ReplaceAllHTTPAfterResponseRuleBackendBadRequest Bad request

swagger:response replaceAllHttpAfterResponseRuleBackendBadRequest
*/
type ReplaceAllHTTPAfterResponseRuleBackendBadRequest struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplaceAllHTTPAfterResponseRuleBackendBadRequest creates ReplaceAllHTTPAfterResponseRuleBackendBadRequest with default headers values
func NewReplaceAllHTTPAfterResponseRuleBackendBadRequest() *ReplaceAllHTTPAfterResponseRuleBackendBadRequest {

	return &ReplaceAllHTTPAfterResponseRuleBackendBadRequest{}
}

// WithConfigurationVersion adds the configurationVersion to the replace all Http after response rule backend bad request response
func (o *ReplaceAllHTTPAfterResponseRuleBackendBadRequest) WithConfigurationVersion(configurationVersion string) *ReplaceAllHTTPAfterResponseRuleBackendBadRequest {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the replace all Http after response rule backend bad request response
func (o *ReplaceAllHTTPAfterResponseRuleBackendBadRequest) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the replace all Http after response rule backend bad request response
func (o *ReplaceAllHTTPAfterResponseRuleBackendBadRequest) WithPayload(payload *models.Error) *ReplaceAllHTTPAfterResponseRuleBackendBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace all Http after response rule backend bad request response
func (o *ReplaceAllHTTPAfterResponseRuleBackendBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceAllHTTPAfterResponseRuleBackendBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
ReplaceAllHTTPAfterResponseRuleBackendDefault General Error

swagger:response replaceAllHttpAfterResponseRuleBackendDefault
*/
type ReplaceAllHTTPAfterResponseRuleBackendDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplaceAllHTTPAfterResponseRuleBackendDefault creates ReplaceAllHTTPAfterResponseRuleBackendDefault with default headers values
func NewReplaceAllHTTPAfterResponseRuleBackendDefault(code int) *ReplaceAllHTTPAfterResponseRuleBackendDefault {
	if code <= 0 {
		code = 500
	}

	return &ReplaceAllHTTPAfterResponseRuleBackendDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the replace all HTTP after response rule backend default response
func (o *ReplaceAllHTTPAfterResponseRuleBackendDefault) WithStatusCode(code int) *ReplaceAllHTTPAfterResponseRuleBackendDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the replace all HTTP after response rule backend default response
func (o *ReplaceAllHTTPAfterResponseRuleBackendDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the replace all HTTP after response rule backend default response
func (o *ReplaceAllHTTPAfterResponseRuleBackendDefault) WithConfigurationVersion(configurationVersion string) *ReplaceAllHTTPAfterResponseRuleBackendDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the replace all HTTP after response rule backend default response
func (o *ReplaceAllHTTPAfterResponseRuleBackendDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the replace all HTTP after response rule backend default response
func (o *ReplaceAllHTTPAfterResponseRuleBackendDefault) WithPayload(payload *models.Error) *ReplaceAllHTTPAfterResponseRuleBackendDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace all HTTP after response rule backend default response
func (o *ReplaceAllHTTPAfterResponseRuleBackendDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceAllHTTPAfterResponseRuleBackendDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
