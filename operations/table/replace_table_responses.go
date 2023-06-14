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

package table

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v5/models"
)

// ReplaceTableOKCode is the HTTP code returned for type ReplaceTableOK
const ReplaceTableOKCode int = 200

/*
ReplaceTableOK Table replaced

swagger:response replaceTableOK
*/
type ReplaceTableOK struct {

	/*
	  In: Body
	*/
	Payload *models.Table `json:"body,omitempty"`
}

// NewReplaceTableOK creates ReplaceTableOK with default headers values
func NewReplaceTableOK() *ReplaceTableOK {

	return &ReplaceTableOK{}
}

// WithPayload adds the payload to the replace table o k response
func (o *ReplaceTableOK) WithPayload(payload *models.Table) *ReplaceTableOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace table o k response
func (o *ReplaceTableOK) SetPayload(payload *models.Table) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceTableOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// ReplaceTableAcceptedCode is the HTTP code returned for type ReplaceTableAccepted
const ReplaceTableAcceptedCode int = 202

/*
ReplaceTableAccepted Configuration change accepted and reload requested

swagger:response replaceTableAccepted
*/
type ReplaceTableAccepted struct {
	/*ID of the requested reload

	 */
	ReloadID string `json:"Reload-ID"`

	/*
	  In: Body
	*/
	Payload *models.Table `json:"body,omitempty"`
}

// NewReplaceTableAccepted creates ReplaceTableAccepted with default headers values
func NewReplaceTableAccepted() *ReplaceTableAccepted {

	return &ReplaceTableAccepted{}
}

// WithReloadID adds the reloadId to the replace table accepted response
func (o *ReplaceTableAccepted) WithReloadID(reloadID string) *ReplaceTableAccepted {
	o.ReloadID = reloadID
	return o
}

// SetReloadID sets the reloadId to the replace table accepted response
func (o *ReplaceTableAccepted) SetReloadID(reloadID string) {
	o.ReloadID = reloadID
}

// WithPayload adds the payload to the replace table accepted response
func (o *ReplaceTableAccepted) WithPayload(payload *models.Table) *ReplaceTableAccepted {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace table accepted response
func (o *ReplaceTableAccepted) SetPayload(payload *models.Table) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceTableAccepted) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// ReplaceTableBadRequestCode is the HTTP code returned for type ReplaceTableBadRequest
const ReplaceTableBadRequestCode int = 400

/*
ReplaceTableBadRequest Bad request

swagger:response replaceTableBadRequest
*/
type ReplaceTableBadRequest struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplaceTableBadRequest creates ReplaceTableBadRequest with default headers values
func NewReplaceTableBadRequest() *ReplaceTableBadRequest {

	return &ReplaceTableBadRequest{}
}

// WithConfigurationVersion adds the configurationVersion to the replace table bad request response
func (o *ReplaceTableBadRequest) WithConfigurationVersion(configurationVersion string) *ReplaceTableBadRequest {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the replace table bad request response
func (o *ReplaceTableBadRequest) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the replace table bad request response
func (o *ReplaceTableBadRequest) WithPayload(payload *models.Error) *ReplaceTableBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace table bad request response
func (o *ReplaceTableBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceTableBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// ReplaceTableNotFoundCode is the HTTP code returned for type ReplaceTableNotFound
const ReplaceTableNotFoundCode int = 404

/*
ReplaceTableNotFound The specified resource was not found

swagger:response replaceTableNotFound
*/
type ReplaceTableNotFound struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplaceTableNotFound creates ReplaceTableNotFound with default headers values
func NewReplaceTableNotFound() *ReplaceTableNotFound {

	return &ReplaceTableNotFound{}
}

// WithConfigurationVersion adds the configurationVersion to the replace table not found response
func (o *ReplaceTableNotFound) WithConfigurationVersion(configurationVersion string) *ReplaceTableNotFound {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the replace table not found response
func (o *ReplaceTableNotFound) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the replace table not found response
func (o *ReplaceTableNotFound) WithPayload(payload *models.Error) *ReplaceTableNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace table not found response
func (o *ReplaceTableNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceTableNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
ReplaceTableDefault General Error

swagger:response replaceTableDefault
*/
type ReplaceTableDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplaceTableDefault creates ReplaceTableDefault with default headers values
func NewReplaceTableDefault(code int) *ReplaceTableDefault {
	if code <= 0 {
		code = 500
	}

	return &ReplaceTableDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the replace table default response
func (o *ReplaceTableDefault) WithStatusCode(code int) *ReplaceTableDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the replace table default response
func (o *ReplaceTableDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the replace table default response
func (o *ReplaceTableDefault) WithConfigurationVersion(configurationVersion string) *ReplaceTableDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the replace table default response
func (o *ReplaceTableDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the replace table default response
func (o *ReplaceTableDefault) WithPayload(payload *models.Error) *ReplaceTableDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace table default response
func (o *ReplaceTableDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceTableDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
