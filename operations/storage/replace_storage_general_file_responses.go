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

package storage

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v3/models"
)

// ReplaceStorageGeneralFileAcceptedCode is the HTTP code returned for type ReplaceStorageGeneralFileAccepted
const ReplaceStorageGeneralFileAcceptedCode int = 202

/*
ReplaceStorageGeneralFileAccepted Configuration change accepted and reload requested

swagger:response replaceStorageGeneralFileAccepted
*/
type ReplaceStorageGeneralFileAccepted struct {
	/*ID of the requested reload

	 */
	ReloadID string `json:"Reload-ID"`
}

// NewReplaceStorageGeneralFileAccepted creates ReplaceStorageGeneralFileAccepted with default headers values
func NewReplaceStorageGeneralFileAccepted() *ReplaceStorageGeneralFileAccepted {

	return &ReplaceStorageGeneralFileAccepted{}
}

// WithReloadID adds the reloadId to the replace storage general file accepted response
func (o *ReplaceStorageGeneralFileAccepted) WithReloadID(reloadID string) *ReplaceStorageGeneralFileAccepted {
	o.ReloadID = reloadID
	return o
}

// SetReloadID sets the reloadId to the replace storage general file accepted response
func (o *ReplaceStorageGeneralFileAccepted) SetReloadID(reloadID string) {
	o.ReloadID = reloadID
}

// WriteResponse to the client
func (o *ReplaceStorageGeneralFileAccepted) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	// response header Reload-ID

	reloadID := o.ReloadID
	if reloadID != "" {
		rw.Header().Set("Reload-ID", reloadID)
	}

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(202)
}

// ReplaceStorageGeneralFileNoContentCode is the HTTP code returned for type ReplaceStorageGeneralFileNoContent
const ReplaceStorageGeneralFileNoContentCode int = 204

/*
ReplaceStorageGeneralFileNoContent General use file replaced

swagger:response replaceStorageGeneralFileNoContent
*/
type ReplaceStorageGeneralFileNoContent struct {
}

// NewReplaceStorageGeneralFileNoContent creates ReplaceStorageGeneralFileNoContent with default headers values
func NewReplaceStorageGeneralFileNoContent() *ReplaceStorageGeneralFileNoContent {

	return &ReplaceStorageGeneralFileNoContent{}
}

// WriteResponse to the client
func (o *ReplaceStorageGeneralFileNoContent) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(204)
}

// ReplaceStorageGeneralFileBadRequestCode is the HTTP code returned for type ReplaceStorageGeneralFileBadRequest
const ReplaceStorageGeneralFileBadRequestCode int = 400

/*
ReplaceStorageGeneralFileBadRequest Bad request

swagger:response replaceStorageGeneralFileBadRequest
*/
type ReplaceStorageGeneralFileBadRequest struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplaceStorageGeneralFileBadRequest creates ReplaceStorageGeneralFileBadRequest with default headers values
func NewReplaceStorageGeneralFileBadRequest() *ReplaceStorageGeneralFileBadRequest {

	return &ReplaceStorageGeneralFileBadRequest{}
}

// WithConfigurationVersion adds the configurationVersion to the replace storage general file bad request response
func (o *ReplaceStorageGeneralFileBadRequest) WithConfigurationVersion(configurationVersion string) *ReplaceStorageGeneralFileBadRequest {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the replace storage general file bad request response
func (o *ReplaceStorageGeneralFileBadRequest) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the replace storage general file bad request response
func (o *ReplaceStorageGeneralFileBadRequest) WithPayload(payload *models.Error) *ReplaceStorageGeneralFileBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace storage general file bad request response
func (o *ReplaceStorageGeneralFileBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceStorageGeneralFileBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// ReplaceStorageGeneralFileNotFoundCode is the HTTP code returned for type ReplaceStorageGeneralFileNotFound
const ReplaceStorageGeneralFileNotFoundCode int = 404

/*
ReplaceStorageGeneralFileNotFound The specified resource was not found

swagger:response replaceStorageGeneralFileNotFound
*/
type ReplaceStorageGeneralFileNotFound struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplaceStorageGeneralFileNotFound creates ReplaceStorageGeneralFileNotFound with default headers values
func NewReplaceStorageGeneralFileNotFound() *ReplaceStorageGeneralFileNotFound {

	return &ReplaceStorageGeneralFileNotFound{}
}

// WithConfigurationVersion adds the configurationVersion to the replace storage general file not found response
func (o *ReplaceStorageGeneralFileNotFound) WithConfigurationVersion(configurationVersion string) *ReplaceStorageGeneralFileNotFound {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the replace storage general file not found response
func (o *ReplaceStorageGeneralFileNotFound) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the replace storage general file not found response
func (o *ReplaceStorageGeneralFileNotFound) WithPayload(payload *models.Error) *ReplaceStorageGeneralFileNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace storage general file not found response
func (o *ReplaceStorageGeneralFileNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceStorageGeneralFileNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
ReplaceStorageGeneralFileDefault General Error

swagger:response replaceStorageGeneralFileDefault
*/
type ReplaceStorageGeneralFileDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplaceStorageGeneralFileDefault creates ReplaceStorageGeneralFileDefault with default headers values
func NewReplaceStorageGeneralFileDefault(code int) *ReplaceStorageGeneralFileDefault {
	if code <= 0 {
		code = 500
	}

	return &ReplaceStorageGeneralFileDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the replace storage general file default response
func (o *ReplaceStorageGeneralFileDefault) WithStatusCode(code int) *ReplaceStorageGeneralFileDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the replace storage general file default response
func (o *ReplaceStorageGeneralFileDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the replace storage general file default response
func (o *ReplaceStorageGeneralFileDefault) WithConfigurationVersion(configurationVersion string) *ReplaceStorageGeneralFileDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the replace storage general file default response
func (o *ReplaceStorageGeneralFileDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the replace storage general file default response
func (o *ReplaceStorageGeneralFileDefault) WithPayload(payload *models.Error) *ReplaceStorageGeneralFileDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace storage general file default response
func (o *ReplaceStorageGeneralFileDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceStorageGeneralFileDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
