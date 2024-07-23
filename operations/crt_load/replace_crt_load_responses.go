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

package crt_load

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v6/models"
)

// ReplaceCrtLoadOKCode is the HTTP code returned for type ReplaceCrtLoadOK
const ReplaceCrtLoadOKCode int = 200

/*
ReplaceCrtLoadOK Certificate load entry replaced

swagger:response replaceCrtLoadOK
*/
type ReplaceCrtLoadOK struct {

	/*
	  In: Body
	*/
	Payload *models.CrtLoad `json:"body,omitempty"`
}

// NewReplaceCrtLoadOK creates ReplaceCrtLoadOK with default headers values
func NewReplaceCrtLoadOK() *ReplaceCrtLoadOK {

	return &ReplaceCrtLoadOK{}
}

// WithPayload adds the payload to the replace crt load o k response
func (o *ReplaceCrtLoadOK) WithPayload(payload *models.CrtLoad) *ReplaceCrtLoadOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace crt load o k response
func (o *ReplaceCrtLoadOK) SetPayload(payload *models.CrtLoad) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceCrtLoadOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// ReplaceCrtLoadAcceptedCode is the HTTP code returned for type ReplaceCrtLoadAccepted
const ReplaceCrtLoadAcceptedCode int = 202

/*
ReplaceCrtLoadAccepted Configuration change accepted and reload requested

swagger:response replaceCrtLoadAccepted
*/
type ReplaceCrtLoadAccepted struct {
	/*ID of the requested reload

	 */
	ReloadID string `json:"Reload-ID"`

	/*
	  In: Body
	*/
	Payload *models.CrtLoad `json:"body,omitempty"`
}

// NewReplaceCrtLoadAccepted creates ReplaceCrtLoadAccepted with default headers values
func NewReplaceCrtLoadAccepted() *ReplaceCrtLoadAccepted {

	return &ReplaceCrtLoadAccepted{}
}

// WithReloadID adds the reloadId to the replace crt load accepted response
func (o *ReplaceCrtLoadAccepted) WithReloadID(reloadID string) *ReplaceCrtLoadAccepted {
	o.ReloadID = reloadID
	return o
}

// SetReloadID sets the reloadId to the replace crt load accepted response
func (o *ReplaceCrtLoadAccepted) SetReloadID(reloadID string) {
	o.ReloadID = reloadID
}

// WithPayload adds the payload to the replace crt load accepted response
func (o *ReplaceCrtLoadAccepted) WithPayload(payload *models.CrtLoad) *ReplaceCrtLoadAccepted {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace crt load accepted response
func (o *ReplaceCrtLoadAccepted) SetPayload(payload *models.CrtLoad) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceCrtLoadAccepted) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// ReplaceCrtLoadBadRequestCode is the HTTP code returned for type ReplaceCrtLoadBadRequest
const ReplaceCrtLoadBadRequestCode int = 400

/*
ReplaceCrtLoadBadRequest Bad request

swagger:response replaceCrtLoadBadRequest
*/
type ReplaceCrtLoadBadRequest struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplaceCrtLoadBadRequest creates ReplaceCrtLoadBadRequest with default headers values
func NewReplaceCrtLoadBadRequest() *ReplaceCrtLoadBadRequest {

	return &ReplaceCrtLoadBadRequest{}
}

// WithConfigurationVersion adds the configurationVersion to the replace crt load bad request response
func (o *ReplaceCrtLoadBadRequest) WithConfigurationVersion(configurationVersion string) *ReplaceCrtLoadBadRequest {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the replace crt load bad request response
func (o *ReplaceCrtLoadBadRequest) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the replace crt load bad request response
func (o *ReplaceCrtLoadBadRequest) WithPayload(payload *models.Error) *ReplaceCrtLoadBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace crt load bad request response
func (o *ReplaceCrtLoadBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceCrtLoadBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// ReplaceCrtLoadNotFoundCode is the HTTP code returned for type ReplaceCrtLoadNotFound
const ReplaceCrtLoadNotFoundCode int = 404

/*
ReplaceCrtLoadNotFound The specified resource was not found

swagger:response replaceCrtLoadNotFound
*/
type ReplaceCrtLoadNotFound struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplaceCrtLoadNotFound creates ReplaceCrtLoadNotFound with default headers values
func NewReplaceCrtLoadNotFound() *ReplaceCrtLoadNotFound {

	return &ReplaceCrtLoadNotFound{}
}

// WithConfigurationVersion adds the configurationVersion to the replace crt load not found response
func (o *ReplaceCrtLoadNotFound) WithConfigurationVersion(configurationVersion string) *ReplaceCrtLoadNotFound {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the replace crt load not found response
func (o *ReplaceCrtLoadNotFound) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the replace crt load not found response
func (o *ReplaceCrtLoadNotFound) WithPayload(payload *models.Error) *ReplaceCrtLoadNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace crt load not found response
func (o *ReplaceCrtLoadNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceCrtLoadNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
ReplaceCrtLoadDefault General Error

swagger:response replaceCrtLoadDefault
*/
type ReplaceCrtLoadDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplaceCrtLoadDefault creates ReplaceCrtLoadDefault with default headers values
func NewReplaceCrtLoadDefault(code int) *ReplaceCrtLoadDefault {
	if code <= 0 {
		code = 500
	}

	return &ReplaceCrtLoadDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the replace crt load default response
func (o *ReplaceCrtLoadDefault) WithStatusCode(code int) *ReplaceCrtLoadDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the replace crt load default response
func (o *ReplaceCrtLoadDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the replace crt load default response
func (o *ReplaceCrtLoadDefault) WithConfigurationVersion(configurationVersion string) *ReplaceCrtLoadDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the replace crt load default response
func (o *ReplaceCrtLoadDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the replace crt load default response
func (o *ReplaceCrtLoadDefault) WithPayload(payload *models.Error) *ReplaceCrtLoadDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace crt load default response
func (o *ReplaceCrtLoadDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceCrtLoadDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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