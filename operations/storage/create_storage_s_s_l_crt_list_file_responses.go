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

// CreateStorageSSLCrtListFileCreatedCode is the HTTP code returned for type CreateStorageSSLCrtListFileCreated
const CreateStorageSSLCrtListFileCreatedCode int = 201

/*
CreateStorageSSLCrtListFileCreated Certificate list created

swagger:response createStorageSSLCrtListFileCreated
*/
type CreateStorageSSLCrtListFileCreated struct {

	/*
	  In: Body
	*/
	Payload *models.SslCrtListFile `json:"body,omitempty"`
}

// NewCreateStorageSSLCrtListFileCreated creates CreateStorageSSLCrtListFileCreated with default headers values
func NewCreateStorageSSLCrtListFileCreated() *CreateStorageSSLCrtListFileCreated {

	return &CreateStorageSSLCrtListFileCreated{}
}

// WithPayload adds the payload to the create storage s s l crt list file created response
func (o *CreateStorageSSLCrtListFileCreated) WithPayload(payload *models.SslCrtListFile) *CreateStorageSSLCrtListFileCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create storage s s l crt list file created response
func (o *CreateStorageSSLCrtListFileCreated) SetPayload(payload *models.SslCrtListFile) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateStorageSSLCrtListFileCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateStorageSSLCrtListFileAcceptedCode is the HTTP code returned for type CreateStorageSSLCrtListFileAccepted
const CreateStorageSSLCrtListFileAcceptedCode int = 202

/*
CreateStorageSSLCrtListFileAccepted Certificate list created requested

swagger:response createStorageSSLCrtListFileAccepted
*/
type CreateStorageSSLCrtListFileAccepted struct {
	/*ID of the requested reload

	 */
	ReloadID string `json:"Reload-ID"`

	/*
	  In: Body
	*/
	Payload *models.SslCrtListFile `json:"body,omitempty"`
}

// NewCreateStorageSSLCrtListFileAccepted creates CreateStorageSSLCrtListFileAccepted with default headers values
func NewCreateStorageSSLCrtListFileAccepted() *CreateStorageSSLCrtListFileAccepted {

	return &CreateStorageSSLCrtListFileAccepted{}
}

// WithReloadID adds the reloadId to the create storage s s l crt list file accepted response
func (o *CreateStorageSSLCrtListFileAccepted) WithReloadID(reloadID string) *CreateStorageSSLCrtListFileAccepted {
	o.ReloadID = reloadID
	return o
}

// SetReloadID sets the reloadId to the create storage s s l crt list file accepted response
func (o *CreateStorageSSLCrtListFileAccepted) SetReloadID(reloadID string) {
	o.ReloadID = reloadID
}

// WithPayload adds the payload to the create storage s s l crt list file accepted response
func (o *CreateStorageSSLCrtListFileAccepted) WithPayload(payload *models.SslCrtListFile) *CreateStorageSSLCrtListFileAccepted {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create storage s s l crt list file accepted response
func (o *CreateStorageSSLCrtListFileAccepted) SetPayload(payload *models.SslCrtListFile) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateStorageSSLCrtListFileAccepted) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// CreateStorageSSLCrtListFileBadRequestCode is the HTTP code returned for type CreateStorageSSLCrtListFileBadRequest
const CreateStorageSSLCrtListFileBadRequestCode int = 400

/*
CreateStorageSSLCrtListFileBadRequest Bad request

swagger:response createStorageSSLCrtListFileBadRequest
*/
type CreateStorageSSLCrtListFileBadRequest struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateStorageSSLCrtListFileBadRequest creates CreateStorageSSLCrtListFileBadRequest with default headers values
func NewCreateStorageSSLCrtListFileBadRequest() *CreateStorageSSLCrtListFileBadRequest {

	return &CreateStorageSSLCrtListFileBadRequest{}
}

// WithConfigurationVersion adds the configurationVersion to the create storage s s l crt list file bad request response
func (o *CreateStorageSSLCrtListFileBadRequest) WithConfigurationVersion(configurationVersion string) *CreateStorageSSLCrtListFileBadRequest {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the create storage s s l crt list file bad request response
func (o *CreateStorageSSLCrtListFileBadRequest) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the create storage s s l crt list file bad request response
func (o *CreateStorageSSLCrtListFileBadRequest) WithPayload(payload *models.Error) *CreateStorageSSLCrtListFileBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create storage s s l crt list file bad request response
func (o *CreateStorageSSLCrtListFileBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateStorageSSLCrtListFileBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// CreateStorageSSLCrtListFileConflictCode is the HTTP code returned for type CreateStorageSSLCrtListFileConflict
const CreateStorageSSLCrtListFileConflictCode int = 409

/*
CreateStorageSSLCrtListFileConflict The specified resource already exists

swagger:response createStorageSSLCrtListFileConflict
*/
type CreateStorageSSLCrtListFileConflict struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateStorageSSLCrtListFileConflict creates CreateStorageSSLCrtListFileConflict with default headers values
func NewCreateStorageSSLCrtListFileConflict() *CreateStorageSSLCrtListFileConflict {

	return &CreateStorageSSLCrtListFileConflict{}
}

// WithConfigurationVersion adds the configurationVersion to the create storage s s l crt list file conflict response
func (o *CreateStorageSSLCrtListFileConflict) WithConfigurationVersion(configurationVersion string) *CreateStorageSSLCrtListFileConflict {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the create storage s s l crt list file conflict response
func (o *CreateStorageSSLCrtListFileConflict) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the create storage s s l crt list file conflict response
func (o *CreateStorageSSLCrtListFileConflict) WithPayload(payload *models.Error) *CreateStorageSSLCrtListFileConflict {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create storage s s l crt list file conflict response
func (o *CreateStorageSSLCrtListFileConflict) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateStorageSSLCrtListFileConflict) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	// response header Configuration-Version

	configurationVersion := o.ConfigurationVersion
	if configurationVersion != "" {
		rw.Header().Set("Configuration-Version", configurationVersion)
	}

	rw.WriteHeader(409)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*
CreateStorageSSLCrtListFileDefault General Error

swagger:response createStorageSSLCrtListFileDefault
*/
type CreateStorageSSLCrtListFileDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateStorageSSLCrtListFileDefault creates CreateStorageSSLCrtListFileDefault with default headers values
func NewCreateStorageSSLCrtListFileDefault(code int) *CreateStorageSSLCrtListFileDefault {
	if code <= 0 {
		code = 500
	}

	return &CreateStorageSSLCrtListFileDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the create storage s s l crt list file default response
func (o *CreateStorageSSLCrtListFileDefault) WithStatusCode(code int) *CreateStorageSSLCrtListFileDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the create storage s s l crt list file default response
func (o *CreateStorageSSLCrtListFileDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the create storage s s l crt list file default response
func (o *CreateStorageSSLCrtListFileDefault) WithConfigurationVersion(configurationVersion string) *CreateStorageSSLCrtListFileDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the create storage s s l crt list file default response
func (o *CreateStorageSSLCrtListFileDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the create storage s s l crt list file default response
func (o *CreateStorageSSLCrtListFileDefault) WithPayload(payload *models.Error) *CreateStorageSSLCrtListFileDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create storage s s l crt list file default response
func (o *CreateStorageSSLCrtListFileDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateStorageSSLCrtListFileDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
