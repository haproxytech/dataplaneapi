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

package spoe

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v4/models"
)

// ReplaceSpoeMessageOKCode is the HTTP code returned for type ReplaceSpoeMessageOK
const ReplaceSpoeMessageOKCode int = 200

/*
ReplaceSpoeMessageOK Spoe message replaced

swagger:response replaceSpoeMessageOK
*/
type ReplaceSpoeMessageOK struct {

	/*
	  In: Body
	*/
	Payload *models.SpoeMessage `json:"body,omitempty"`
}

// NewReplaceSpoeMessageOK creates ReplaceSpoeMessageOK with default headers values
func NewReplaceSpoeMessageOK() *ReplaceSpoeMessageOK {

	return &ReplaceSpoeMessageOK{}
}

// WithPayload adds the payload to the replace spoe message o k response
func (o *ReplaceSpoeMessageOK) WithPayload(payload *models.SpoeMessage) *ReplaceSpoeMessageOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace spoe message o k response
func (o *ReplaceSpoeMessageOK) SetPayload(payload *models.SpoeMessage) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceSpoeMessageOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// ReplaceSpoeMessageBadRequestCode is the HTTP code returned for type ReplaceSpoeMessageBadRequest
const ReplaceSpoeMessageBadRequestCode int = 400

/*
ReplaceSpoeMessageBadRequest Bad request

swagger:response replaceSpoeMessageBadRequest
*/
type ReplaceSpoeMessageBadRequest struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplaceSpoeMessageBadRequest creates ReplaceSpoeMessageBadRequest with default headers values
func NewReplaceSpoeMessageBadRequest() *ReplaceSpoeMessageBadRequest {

	return &ReplaceSpoeMessageBadRequest{}
}

// WithConfigurationVersion adds the configurationVersion to the replace spoe message bad request response
func (o *ReplaceSpoeMessageBadRequest) WithConfigurationVersion(configurationVersion string) *ReplaceSpoeMessageBadRequest {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the replace spoe message bad request response
func (o *ReplaceSpoeMessageBadRequest) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the replace spoe message bad request response
func (o *ReplaceSpoeMessageBadRequest) WithPayload(payload *models.Error) *ReplaceSpoeMessageBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace spoe message bad request response
func (o *ReplaceSpoeMessageBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceSpoeMessageBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// ReplaceSpoeMessageNotFoundCode is the HTTP code returned for type ReplaceSpoeMessageNotFound
const ReplaceSpoeMessageNotFoundCode int = 404

/*
ReplaceSpoeMessageNotFound The specified resource was not found

swagger:response replaceSpoeMessageNotFound
*/
type ReplaceSpoeMessageNotFound struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplaceSpoeMessageNotFound creates ReplaceSpoeMessageNotFound with default headers values
func NewReplaceSpoeMessageNotFound() *ReplaceSpoeMessageNotFound {

	return &ReplaceSpoeMessageNotFound{}
}

// WithConfigurationVersion adds the configurationVersion to the replace spoe message not found response
func (o *ReplaceSpoeMessageNotFound) WithConfigurationVersion(configurationVersion string) *ReplaceSpoeMessageNotFound {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the replace spoe message not found response
func (o *ReplaceSpoeMessageNotFound) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the replace spoe message not found response
func (o *ReplaceSpoeMessageNotFound) WithPayload(payload *models.Error) *ReplaceSpoeMessageNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace spoe message not found response
func (o *ReplaceSpoeMessageNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceSpoeMessageNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
ReplaceSpoeMessageDefault General Error

swagger:response replaceSpoeMessageDefault
*/
type ReplaceSpoeMessageDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplaceSpoeMessageDefault creates ReplaceSpoeMessageDefault with default headers values
func NewReplaceSpoeMessageDefault(code int) *ReplaceSpoeMessageDefault {
	if code <= 0 {
		code = 500
	}

	return &ReplaceSpoeMessageDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the replace spoe message default response
func (o *ReplaceSpoeMessageDefault) WithStatusCode(code int) *ReplaceSpoeMessageDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the replace spoe message default response
func (o *ReplaceSpoeMessageDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the replace spoe message default response
func (o *ReplaceSpoeMessageDefault) WithConfigurationVersion(configurationVersion string) *ReplaceSpoeMessageDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the replace spoe message default response
func (o *ReplaceSpoeMessageDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the replace spoe message default response
func (o *ReplaceSpoeMessageDefault) WithPayload(payload *models.Error) *ReplaceSpoeMessageDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace spoe message default response
func (o *ReplaceSpoeMessageDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceSpoeMessageDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
