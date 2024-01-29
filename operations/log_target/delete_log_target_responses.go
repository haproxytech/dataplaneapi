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

package log_target

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v6/models"
)

// DeleteLogTargetAcceptedCode is the HTTP code returned for type DeleteLogTargetAccepted
const DeleteLogTargetAcceptedCode int = 202

/*
DeleteLogTargetAccepted Configuration change accepted and reload requested

swagger:response deleteLogTargetAccepted
*/
type DeleteLogTargetAccepted struct {
	/*ID of the requested reload

	 */
	ReloadID string `json:"Reload-ID"`
}

// NewDeleteLogTargetAccepted creates DeleteLogTargetAccepted with default headers values
func NewDeleteLogTargetAccepted() *DeleteLogTargetAccepted {

	return &DeleteLogTargetAccepted{}
}

// WithReloadID adds the reloadId to the delete log target accepted response
func (o *DeleteLogTargetAccepted) WithReloadID(reloadID string) *DeleteLogTargetAccepted {
	o.ReloadID = reloadID
	return o
}

// SetReloadID sets the reloadId to the delete log target accepted response
func (o *DeleteLogTargetAccepted) SetReloadID(reloadID string) {
	o.ReloadID = reloadID
}

// WriteResponse to the client
func (o *DeleteLogTargetAccepted) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	// response header Reload-ID

	reloadID := o.ReloadID
	if reloadID != "" {
		rw.Header().Set("Reload-ID", reloadID)
	}

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(202)
}

// DeleteLogTargetNoContentCode is the HTTP code returned for type DeleteLogTargetNoContent
const DeleteLogTargetNoContentCode int = 204

/*
DeleteLogTargetNoContent Log Target deleted

swagger:response deleteLogTargetNoContent
*/
type DeleteLogTargetNoContent struct {
}

// NewDeleteLogTargetNoContent creates DeleteLogTargetNoContent with default headers values
func NewDeleteLogTargetNoContent() *DeleteLogTargetNoContent {

	return &DeleteLogTargetNoContent{}
}

// WriteResponse to the client
func (o *DeleteLogTargetNoContent) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(204)
}

// DeleteLogTargetNotFoundCode is the HTTP code returned for type DeleteLogTargetNotFound
const DeleteLogTargetNotFoundCode int = 404

/*
DeleteLogTargetNotFound The specified resource was not found

swagger:response deleteLogTargetNotFound
*/
type DeleteLogTargetNotFound struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDeleteLogTargetNotFound creates DeleteLogTargetNotFound with default headers values
func NewDeleteLogTargetNotFound() *DeleteLogTargetNotFound {

	return &DeleteLogTargetNotFound{}
}

// WithConfigurationVersion adds the configurationVersion to the delete log target not found response
func (o *DeleteLogTargetNotFound) WithConfigurationVersion(configurationVersion string) *DeleteLogTargetNotFound {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the delete log target not found response
func (o *DeleteLogTargetNotFound) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the delete log target not found response
func (o *DeleteLogTargetNotFound) WithPayload(payload *models.Error) *DeleteLogTargetNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete log target not found response
func (o *DeleteLogTargetNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteLogTargetNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
DeleteLogTargetDefault General Error

swagger:response deleteLogTargetDefault
*/
type DeleteLogTargetDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDeleteLogTargetDefault creates DeleteLogTargetDefault with default headers values
func NewDeleteLogTargetDefault(code int) *DeleteLogTargetDefault {
	if code <= 0 {
		code = 500
	}

	return &DeleteLogTargetDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the delete log target default response
func (o *DeleteLogTargetDefault) WithStatusCode(code int) *DeleteLogTargetDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the delete log target default response
func (o *DeleteLogTargetDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the delete log target default response
func (o *DeleteLogTargetDefault) WithConfigurationVersion(configurationVersion string) *DeleteLogTargetDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the delete log target default response
func (o *DeleteLogTargetDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the delete log target default response
func (o *DeleteLogTargetDefault) WithPayload(payload *models.Error) *DeleteLogTargetDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete log target default response
func (o *DeleteLogTargetDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteLogTargetDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
