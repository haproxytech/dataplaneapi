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

// ReplaceHTTPAfterResponseRuleBackendOKCode is the HTTP code returned for type ReplaceHTTPAfterResponseRuleBackendOK
const ReplaceHTTPAfterResponseRuleBackendOKCode int = 200

/*
ReplaceHTTPAfterResponseRuleBackendOK HTTP After Response Rule replaced

swagger:response replaceHttpAfterResponseRuleBackendOK
*/
type ReplaceHTTPAfterResponseRuleBackendOK struct {

	/*
	  In: Body
	*/
	Payload *models.HTTPAfterResponseRule `json:"body,omitempty"`
}

// NewReplaceHTTPAfterResponseRuleBackendOK creates ReplaceHTTPAfterResponseRuleBackendOK with default headers values
func NewReplaceHTTPAfterResponseRuleBackendOK() *ReplaceHTTPAfterResponseRuleBackendOK {

	return &ReplaceHTTPAfterResponseRuleBackendOK{}
}

// WithPayload adds the payload to the replace Http after response rule backend o k response
func (o *ReplaceHTTPAfterResponseRuleBackendOK) WithPayload(payload *models.HTTPAfterResponseRule) *ReplaceHTTPAfterResponseRuleBackendOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace Http after response rule backend o k response
func (o *ReplaceHTTPAfterResponseRuleBackendOK) SetPayload(payload *models.HTTPAfterResponseRule) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceHTTPAfterResponseRuleBackendOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// ReplaceHTTPAfterResponseRuleBackendAcceptedCode is the HTTP code returned for type ReplaceHTTPAfterResponseRuleBackendAccepted
const ReplaceHTTPAfterResponseRuleBackendAcceptedCode int = 202

/*
ReplaceHTTPAfterResponseRuleBackendAccepted Configuration change accepted and reload requested

swagger:response replaceHttpAfterResponseRuleBackendAccepted
*/
type ReplaceHTTPAfterResponseRuleBackendAccepted struct {
	/*ID of the requested reload

	 */
	ReloadID string `json:"Reload-ID"`

	/*
	  In: Body
	*/
	Payload *models.HTTPAfterResponseRule `json:"body,omitempty"`
}

// NewReplaceHTTPAfterResponseRuleBackendAccepted creates ReplaceHTTPAfterResponseRuleBackendAccepted with default headers values
func NewReplaceHTTPAfterResponseRuleBackendAccepted() *ReplaceHTTPAfterResponseRuleBackendAccepted {

	return &ReplaceHTTPAfterResponseRuleBackendAccepted{}
}

// WithReloadID adds the reloadId to the replace Http after response rule backend accepted response
func (o *ReplaceHTTPAfterResponseRuleBackendAccepted) WithReloadID(reloadID string) *ReplaceHTTPAfterResponseRuleBackendAccepted {
	o.ReloadID = reloadID
	return o
}

// SetReloadID sets the reloadId to the replace Http after response rule backend accepted response
func (o *ReplaceHTTPAfterResponseRuleBackendAccepted) SetReloadID(reloadID string) {
	o.ReloadID = reloadID
}

// WithPayload adds the payload to the replace Http after response rule backend accepted response
func (o *ReplaceHTTPAfterResponseRuleBackendAccepted) WithPayload(payload *models.HTTPAfterResponseRule) *ReplaceHTTPAfterResponseRuleBackendAccepted {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace Http after response rule backend accepted response
func (o *ReplaceHTTPAfterResponseRuleBackendAccepted) SetPayload(payload *models.HTTPAfterResponseRule) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceHTTPAfterResponseRuleBackendAccepted) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// ReplaceHTTPAfterResponseRuleBackendBadRequestCode is the HTTP code returned for type ReplaceHTTPAfterResponseRuleBackendBadRequest
const ReplaceHTTPAfterResponseRuleBackendBadRequestCode int = 400

/*
ReplaceHTTPAfterResponseRuleBackendBadRequest Bad request

swagger:response replaceHttpAfterResponseRuleBackendBadRequest
*/
type ReplaceHTTPAfterResponseRuleBackendBadRequest struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplaceHTTPAfterResponseRuleBackendBadRequest creates ReplaceHTTPAfterResponseRuleBackendBadRequest with default headers values
func NewReplaceHTTPAfterResponseRuleBackendBadRequest() *ReplaceHTTPAfterResponseRuleBackendBadRequest {

	return &ReplaceHTTPAfterResponseRuleBackendBadRequest{}
}

// WithConfigurationVersion adds the configurationVersion to the replace Http after response rule backend bad request response
func (o *ReplaceHTTPAfterResponseRuleBackendBadRequest) WithConfigurationVersion(configurationVersion string) *ReplaceHTTPAfterResponseRuleBackendBadRequest {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the replace Http after response rule backend bad request response
func (o *ReplaceHTTPAfterResponseRuleBackendBadRequest) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the replace Http after response rule backend bad request response
func (o *ReplaceHTTPAfterResponseRuleBackendBadRequest) WithPayload(payload *models.Error) *ReplaceHTTPAfterResponseRuleBackendBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace Http after response rule backend bad request response
func (o *ReplaceHTTPAfterResponseRuleBackendBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceHTTPAfterResponseRuleBackendBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// ReplaceHTTPAfterResponseRuleBackendNotFoundCode is the HTTP code returned for type ReplaceHTTPAfterResponseRuleBackendNotFound
const ReplaceHTTPAfterResponseRuleBackendNotFoundCode int = 404

/*
ReplaceHTTPAfterResponseRuleBackendNotFound The specified resource was not found

swagger:response replaceHttpAfterResponseRuleBackendNotFound
*/
type ReplaceHTTPAfterResponseRuleBackendNotFound struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplaceHTTPAfterResponseRuleBackendNotFound creates ReplaceHTTPAfterResponseRuleBackendNotFound with default headers values
func NewReplaceHTTPAfterResponseRuleBackendNotFound() *ReplaceHTTPAfterResponseRuleBackendNotFound {

	return &ReplaceHTTPAfterResponseRuleBackendNotFound{}
}

// WithConfigurationVersion adds the configurationVersion to the replace Http after response rule backend not found response
func (o *ReplaceHTTPAfterResponseRuleBackendNotFound) WithConfigurationVersion(configurationVersion string) *ReplaceHTTPAfterResponseRuleBackendNotFound {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the replace Http after response rule backend not found response
func (o *ReplaceHTTPAfterResponseRuleBackendNotFound) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the replace Http after response rule backend not found response
func (o *ReplaceHTTPAfterResponseRuleBackendNotFound) WithPayload(payload *models.Error) *ReplaceHTTPAfterResponseRuleBackendNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace Http after response rule backend not found response
func (o *ReplaceHTTPAfterResponseRuleBackendNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceHTTPAfterResponseRuleBackendNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
ReplaceHTTPAfterResponseRuleBackendDefault General Error

swagger:response replaceHttpAfterResponseRuleBackendDefault
*/
type ReplaceHTTPAfterResponseRuleBackendDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplaceHTTPAfterResponseRuleBackendDefault creates ReplaceHTTPAfterResponseRuleBackendDefault with default headers values
func NewReplaceHTTPAfterResponseRuleBackendDefault(code int) *ReplaceHTTPAfterResponseRuleBackendDefault {
	if code <= 0 {
		code = 500
	}

	return &ReplaceHTTPAfterResponseRuleBackendDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the replace HTTP after response rule backend default response
func (o *ReplaceHTTPAfterResponseRuleBackendDefault) WithStatusCode(code int) *ReplaceHTTPAfterResponseRuleBackendDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the replace HTTP after response rule backend default response
func (o *ReplaceHTTPAfterResponseRuleBackendDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the replace HTTP after response rule backend default response
func (o *ReplaceHTTPAfterResponseRuleBackendDefault) WithConfigurationVersion(configurationVersion string) *ReplaceHTTPAfterResponseRuleBackendDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the replace HTTP after response rule backend default response
func (o *ReplaceHTTPAfterResponseRuleBackendDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the replace HTTP after response rule backend default response
func (o *ReplaceHTTPAfterResponseRuleBackendDefault) WithPayload(payload *models.Error) *ReplaceHTTPAfterResponseRuleBackendDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace HTTP after response rule backend default response
func (o *ReplaceHTTPAfterResponseRuleBackendDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceHTTPAfterResponseRuleBackendDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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