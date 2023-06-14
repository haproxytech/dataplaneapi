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

	"github.com/haproxytech/client-native/v5/models"
)

// ReplaceHTTPAfterResponseRuleOKCode is the HTTP code returned for type ReplaceHTTPAfterResponseRuleOK
const ReplaceHTTPAfterResponseRuleOKCode int = 200

/*
ReplaceHTTPAfterResponseRuleOK HTTP After Response Rule replaced

swagger:response replaceHttpAfterResponseRuleOK
*/
type ReplaceHTTPAfterResponseRuleOK struct {

	/*
	  In: Body
	*/
	Payload *models.HTTPAfterResponseRule `json:"body,omitempty"`
}

// NewReplaceHTTPAfterResponseRuleOK creates ReplaceHTTPAfterResponseRuleOK with default headers values
func NewReplaceHTTPAfterResponseRuleOK() *ReplaceHTTPAfterResponseRuleOK {

	return &ReplaceHTTPAfterResponseRuleOK{}
}

// WithPayload adds the payload to the replace Http after response rule o k response
func (o *ReplaceHTTPAfterResponseRuleOK) WithPayload(payload *models.HTTPAfterResponseRule) *ReplaceHTTPAfterResponseRuleOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace Http after response rule o k response
func (o *ReplaceHTTPAfterResponseRuleOK) SetPayload(payload *models.HTTPAfterResponseRule) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceHTTPAfterResponseRuleOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// ReplaceHTTPAfterResponseRuleAcceptedCode is the HTTP code returned for type ReplaceHTTPAfterResponseRuleAccepted
const ReplaceHTTPAfterResponseRuleAcceptedCode int = 202

/*
ReplaceHTTPAfterResponseRuleAccepted Configuration change accepted and reload requested

swagger:response replaceHttpAfterResponseRuleAccepted
*/
type ReplaceHTTPAfterResponseRuleAccepted struct {
	/*ID of the requested reload

	 */
	ReloadID string `json:"Reload-ID"`

	/*
	  In: Body
	*/
	Payload *models.HTTPAfterResponseRule `json:"body,omitempty"`
}

// NewReplaceHTTPAfterResponseRuleAccepted creates ReplaceHTTPAfterResponseRuleAccepted with default headers values
func NewReplaceHTTPAfterResponseRuleAccepted() *ReplaceHTTPAfterResponseRuleAccepted {

	return &ReplaceHTTPAfterResponseRuleAccepted{}
}

// WithReloadID adds the reloadId to the replace Http after response rule accepted response
func (o *ReplaceHTTPAfterResponseRuleAccepted) WithReloadID(reloadID string) *ReplaceHTTPAfterResponseRuleAccepted {
	o.ReloadID = reloadID
	return o
}

// SetReloadID sets the reloadId to the replace Http after response rule accepted response
func (o *ReplaceHTTPAfterResponseRuleAccepted) SetReloadID(reloadID string) {
	o.ReloadID = reloadID
}

// WithPayload adds the payload to the replace Http after response rule accepted response
func (o *ReplaceHTTPAfterResponseRuleAccepted) WithPayload(payload *models.HTTPAfterResponseRule) *ReplaceHTTPAfterResponseRuleAccepted {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace Http after response rule accepted response
func (o *ReplaceHTTPAfterResponseRuleAccepted) SetPayload(payload *models.HTTPAfterResponseRule) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceHTTPAfterResponseRuleAccepted) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// ReplaceHTTPAfterResponseRuleBadRequestCode is the HTTP code returned for type ReplaceHTTPAfterResponseRuleBadRequest
const ReplaceHTTPAfterResponseRuleBadRequestCode int = 400

/*
ReplaceHTTPAfterResponseRuleBadRequest Bad request

swagger:response replaceHttpAfterResponseRuleBadRequest
*/
type ReplaceHTTPAfterResponseRuleBadRequest struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplaceHTTPAfterResponseRuleBadRequest creates ReplaceHTTPAfterResponseRuleBadRequest with default headers values
func NewReplaceHTTPAfterResponseRuleBadRequest() *ReplaceHTTPAfterResponseRuleBadRequest {

	return &ReplaceHTTPAfterResponseRuleBadRequest{}
}

// WithConfigurationVersion adds the configurationVersion to the replace Http after response rule bad request response
func (o *ReplaceHTTPAfterResponseRuleBadRequest) WithConfigurationVersion(configurationVersion string) *ReplaceHTTPAfterResponseRuleBadRequest {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the replace Http after response rule bad request response
func (o *ReplaceHTTPAfterResponseRuleBadRequest) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the replace Http after response rule bad request response
func (o *ReplaceHTTPAfterResponseRuleBadRequest) WithPayload(payload *models.Error) *ReplaceHTTPAfterResponseRuleBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace Http after response rule bad request response
func (o *ReplaceHTTPAfterResponseRuleBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceHTTPAfterResponseRuleBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// ReplaceHTTPAfterResponseRuleNotFoundCode is the HTTP code returned for type ReplaceHTTPAfterResponseRuleNotFound
const ReplaceHTTPAfterResponseRuleNotFoundCode int = 404

/*
ReplaceHTTPAfterResponseRuleNotFound The specified resource was not found

swagger:response replaceHttpAfterResponseRuleNotFound
*/
type ReplaceHTTPAfterResponseRuleNotFound struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplaceHTTPAfterResponseRuleNotFound creates ReplaceHTTPAfterResponseRuleNotFound with default headers values
func NewReplaceHTTPAfterResponseRuleNotFound() *ReplaceHTTPAfterResponseRuleNotFound {

	return &ReplaceHTTPAfterResponseRuleNotFound{}
}

// WithConfigurationVersion adds the configurationVersion to the replace Http after response rule not found response
func (o *ReplaceHTTPAfterResponseRuleNotFound) WithConfigurationVersion(configurationVersion string) *ReplaceHTTPAfterResponseRuleNotFound {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the replace Http after response rule not found response
func (o *ReplaceHTTPAfterResponseRuleNotFound) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the replace Http after response rule not found response
func (o *ReplaceHTTPAfterResponseRuleNotFound) WithPayload(payload *models.Error) *ReplaceHTTPAfterResponseRuleNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace Http after response rule not found response
func (o *ReplaceHTTPAfterResponseRuleNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceHTTPAfterResponseRuleNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
ReplaceHTTPAfterResponseRuleDefault General Error

swagger:response replaceHttpAfterResponseRuleDefault
*/
type ReplaceHTTPAfterResponseRuleDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplaceHTTPAfterResponseRuleDefault creates ReplaceHTTPAfterResponseRuleDefault with default headers values
func NewReplaceHTTPAfterResponseRuleDefault(code int) *ReplaceHTTPAfterResponseRuleDefault {
	if code <= 0 {
		code = 500
	}

	return &ReplaceHTTPAfterResponseRuleDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the replace HTTP after response rule default response
func (o *ReplaceHTTPAfterResponseRuleDefault) WithStatusCode(code int) *ReplaceHTTPAfterResponseRuleDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the replace HTTP after response rule default response
func (o *ReplaceHTTPAfterResponseRuleDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the replace HTTP after response rule default response
func (o *ReplaceHTTPAfterResponseRuleDefault) WithConfigurationVersion(configurationVersion string) *ReplaceHTTPAfterResponseRuleDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the replace HTTP after response rule default response
func (o *ReplaceHTTPAfterResponseRuleDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the replace HTTP after response rule default response
func (o *ReplaceHTTPAfterResponseRuleDefault) WithPayload(payload *models.Error) *ReplaceHTTPAfterResponseRuleDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace HTTP after response rule default response
func (o *ReplaceHTTPAfterResponseRuleDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceHTTPAfterResponseRuleDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
