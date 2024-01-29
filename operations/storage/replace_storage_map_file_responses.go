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

	"github.com/haproxytech/client-native/v6/models"
)

// ReplaceStorageMapFileAcceptedCode is the HTTP code returned for type ReplaceStorageMapFileAccepted
const ReplaceStorageMapFileAcceptedCode int = 202

/*
ReplaceStorageMapFileAccepted Configuration change accepted and reload requested

swagger:response replaceStorageMapFileAccepted
*/
type ReplaceStorageMapFileAccepted struct {
	/*ID of the requested reload

	 */
	ReloadID string `json:"Reload-ID"`
}

// NewReplaceStorageMapFileAccepted creates ReplaceStorageMapFileAccepted with default headers values
func NewReplaceStorageMapFileAccepted() *ReplaceStorageMapFileAccepted {

	return &ReplaceStorageMapFileAccepted{}
}

// WithReloadID adds the reloadId to the replace storage map file accepted response
func (o *ReplaceStorageMapFileAccepted) WithReloadID(reloadID string) *ReplaceStorageMapFileAccepted {
	o.ReloadID = reloadID
	return o
}

// SetReloadID sets the reloadId to the replace storage map file accepted response
func (o *ReplaceStorageMapFileAccepted) SetReloadID(reloadID string) {
	o.ReloadID = reloadID
}

// WriteResponse to the client
func (o *ReplaceStorageMapFileAccepted) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	// response header Reload-ID

	reloadID := o.ReloadID
	if reloadID != "" {
		rw.Header().Set("Reload-ID", reloadID)
	}

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(202)
}

// ReplaceStorageMapFileNoContentCode is the HTTP code returned for type ReplaceStorageMapFileNoContent
const ReplaceStorageMapFileNoContentCode int = 204

/*
ReplaceStorageMapFileNoContent Map file replaced

swagger:response replaceStorageMapFileNoContent
*/
type ReplaceStorageMapFileNoContent struct {
}

// NewReplaceStorageMapFileNoContent creates ReplaceStorageMapFileNoContent with default headers values
func NewReplaceStorageMapFileNoContent() *ReplaceStorageMapFileNoContent {

	return &ReplaceStorageMapFileNoContent{}
}

// WriteResponse to the client
func (o *ReplaceStorageMapFileNoContent) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(204)
}

// ReplaceStorageMapFileBadRequestCode is the HTTP code returned for type ReplaceStorageMapFileBadRequest
const ReplaceStorageMapFileBadRequestCode int = 400

/*
ReplaceStorageMapFileBadRequest Bad request

swagger:response replaceStorageMapFileBadRequest
*/
type ReplaceStorageMapFileBadRequest struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplaceStorageMapFileBadRequest creates ReplaceStorageMapFileBadRequest with default headers values
func NewReplaceStorageMapFileBadRequest() *ReplaceStorageMapFileBadRequest {

	return &ReplaceStorageMapFileBadRequest{}
}

// WithConfigurationVersion adds the configurationVersion to the replace storage map file bad request response
func (o *ReplaceStorageMapFileBadRequest) WithConfigurationVersion(configurationVersion string) *ReplaceStorageMapFileBadRequest {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the replace storage map file bad request response
func (o *ReplaceStorageMapFileBadRequest) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the replace storage map file bad request response
func (o *ReplaceStorageMapFileBadRequest) WithPayload(payload *models.Error) *ReplaceStorageMapFileBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace storage map file bad request response
func (o *ReplaceStorageMapFileBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceStorageMapFileBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// ReplaceStorageMapFileNotFoundCode is the HTTP code returned for type ReplaceStorageMapFileNotFound
const ReplaceStorageMapFileNotFoundCode int = 404

/*
ReplaceStorageMapFileNotFound The specified resource was not found

swagger:response replaceStorageMapFileNotFound
*/
type ReplaceStorageMapFileNotFound struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplaceStorageMapFileNotFound creates ReplaceStorageMapFileNotFound with default headers values
func NewReplaceStorageMapFileNotFound() *ReplaceStorageMapFileNotFound {

	return &ReplaceStorageMapFileNotFound{}
}

// WithConfigurationVersion adds the configurationVersion to the replace storage map file not found response
func (o *ReplaceStorageMapFileNotFound) WithConfigurationVersion(configurationVersion string) *ReplaceStorageMapFileNotFound {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the replace storage map file not found response
func (o *ReplaceStorageMapFileNotFound) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the replace storage map file not found response
func (o *ReplaceStorageMapFileNotFound) WithPayload(payload *models.Error) *ReplaceStorageMapFileNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace storage map file not found response
func (o *ReplaceStorageMapFileNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceStorageMapFileNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
ReplaceStorageMapFileDefault General Error

swagger:response replaceStorageMapFileDefault
*/
type ReplaceStorageMapFileDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplaceStorageMapFileDefault creates ReplaceStorageMapFileDefault with default headers values
func NewReplaceStorageMapFileDefault(code int) *ReplaceStorageMapFileDefault {
	if code <= 0 {
		code = 500
	}

	return &ReplaceStorageMapFileDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the replace storage map file default response
func (o *ReplaceStorageMapFileDefault) WithStatusCode(code int) *ReplaceStorageMapFileDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the replace storage map file default response
func (o *ReplaceStorageMapFileDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the replace storage map file default response
func (o *ReplaceStorageMapFileDefault) WithConfigurationVersion(configurationVersion string) *ReplaceStorageMapFileDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the replace storage map file default response
func (o *ReplaceStorageMapFileDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the replace storage map file default response
func (o *ReplaceStorageMapFileDefault) WithPayload(payload *models.Error) *ReplaceStorageMapFileDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace storage map file default response
func (o *ReplaceStorageMapFileDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceStorageMapFileDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
