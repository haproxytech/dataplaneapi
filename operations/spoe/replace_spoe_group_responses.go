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

	"github.com/haproxytech/client-native/v5/models"
)

// ReplaceSpoeGroupOKCode is the HTTP code returned for type ReplaceSpoeGroupOK
const ReplaceSpoeGroupOKCode int = 200

/*
ReplaceSpoeGroupOK Spoe groups replaced

swagger:response replaceSpoeGroupOK
*/
type ReplaceSpoeGroupOK struct {

	/*
	  In: Body
	*/
	Payload *models.SpoeGroup `json:"body,omitempty"`
}

// NewReplaceSpoeGroupOK creates ReplaceSpoeGroupOK with default headers values
func NewReplaceSpoeGroupOK() *ReplaceSpoeGroupOK {

	return &ReplaceSpoeGroupOK{}
}

// WithPayload adds the payload to the replace spoe group o k response
func (o *ReplaceSpoeGroupOK) WithPayload(payload *models.SpoeGroup) *ReplaceSpoeGroupOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace spoe group o k response
func (o *ReplaceSpoeGroupOK) SetPayload(payload *models.SpoeGroup) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceSpoeGroupOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// ReplaceSpoeGroupBadRequestCode is the HTTP code returned for type ReplaceSpoeGroupBadRequest
const ReplaceSpoeGroupBadRequestCode int = 400

/*
ReplaceSpoeGroupBadRequest Bad request

swagger:response replaceSpoeGroupBadRequest
*/
type ReplaceSpoeGroupBadRequest struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplaceSpoeGroupBadRequest creates ReplaceSpoeGroupBadRequest with default headers values
func NewReplaceSpoeGroupBadRequest() *ReplaceSpoeGroupBadRequest {

	return &ReplaceSpoeGroupBadRequest{}
}

// WithConfigurationVersion adds the configurationVersion to the replace spoe group bad request response
func (o *ReplaceSpoeGroupBadRequest) WithConfigurationVersion(configurationVersion string) *ReplaceSpoeGroupBadRequest {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the replace spoe group bad request response
func (o *ReplaceSpoeGroupBadRequest) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the replace spoe group bad request response
func (o *ReplaceSpoeGroupBadRequest) WithPayload(payload *models.Error) *ReplaceSpoeGroupBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace spoe group bad request response
func (o *ReplaceSpoeGroupBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceSpoeGroupBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// ReplaceSpoeGroupNotFoundCode is the HTTP code returned for type ReplaceSpoeGroupNotFound
const ReplaceSpoeGroupNotFoundCode int = 404

/*
ReplaceSpoeGroupNotFound The specified resource was not found

swagger:response replaceSpoeGroupNotFound
*/
type ReplaceSpoeGroupNotFound struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplaceSpoeGroupNotFound creates ReplaceSpoeGroupNotFound with default headers values
func NewReplaceSpoeGroupNotFound() *ReplaceSpoeGroupNotFound {

	return &ReplaceSpoeGroupNotFound{}
}

// WithConfigurationVersion adds the configurationVersion to the replace spoe group not found response
func (o *ReplaceSpoeGroupNotFound) WithConfigurationVersion(configurationVersion string) *ReplaceSpoeGroupNotFound {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the replace spoe group not found response
func (o *ReplaceSpoeGroupNotFound) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the replace spoe group not found response
func (o *ReplaceSpoeGroupNotFound) WithPayload(payload *models.Error) *ReplaceSpoeGroupNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace spoe group not found response
func (o *ReplaceSpoeGroupNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceSpoeGroupNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
ReplaceSpoeGroupDefault General Error

swagger:response replaceSpoeGroupDefault
*/
type ReplaceSpoeGroupDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplaceSpoeGroupDefault creates ReplaceSpoeGroupDefault with default headers values
func NewReplaceSpoeGroupDefault(code int) *ReplaceSpoeGroupDefault {
	if code <= 0 {
		code = 500
	}

	return &ReplaceSpoeGroupDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the replace spoe group default response
func (o *ReplaceSpoeGroupDefault) WithStatusCode(code int) *ReplaceSpoeGroupDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the replace spoe group default response
func (o *ReplaceSpoeGroupDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the replace spoe group default response
func (o *ReplaceSpoeGroupDefault) WithConfigurationVersion(configurationVersion string) *ReplaceSpoeGroupDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the replace spoe group default response
func (o *ReplaceSpoeGroupDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the replace spoe group default response
func (o *ReplaceSpoeGroupDefault) WithPayload(payload *models.Error) *ReplaceSpoeGroupDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace spoe group default response
func (o *ReplaceSpoeGroupDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceSpoeGroupDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
