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

package group

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v6/models"
)

// DeleteGroupAcceptedCode is the HTTP code returned for type DeleteGroupAccepted
const DeleteGroupAcceptedCode int = 202

/*
DeleteGroupAccepted Configuration change accepted and reload requested

swagger:response deleteGroupAccepted
*/
type DeleteGroupAccepted struct {
	/*ID of the requested reload

	 */
	ReloadID string `json:"Reload-ID"`
}

// NewDeleteGroupAccepted creates DeleteGroupAccepted with default headers values
func NewDeleteGroupAccepted() *DeleteGroupAccepted {

	return &DeleteGroupAccepted{}
}

// WithReloadID adds the reloadId to the delete group accepted response
func (o *DeleteGroupAccepted) WithReloadID(reloadID string) *DeleteGroupAccepted {
	o.ReloadID = reloadID
	return o
}

// SetReloadID sets the reloadId to the delete group accepted response
func (o *DeleteGroupAccepted) SetReloadID(reloadID string) {
	o.ReloadID = reloadID
}

// WriteResponse to the client
func (o *DeleteGroupAccepted) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	// response header Reload-ID

	reloadID := o.ReloadID
	if reloadID != "" {
		rw.Header().Set("Reload-ID", reloadID)
	}

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(202)
}

// DeleteGroupNoContentCode is the HTTP code returned for type DeleteGroupNoContent
const DeleteGroupNoContentCode int = 204

/*
DeleteGroupNoContent Group deleted

swagger:response deleteGroupNoContent
*/
type DeleteGroupNoContent struct {
}

// NewDeleteGroupNoContent creates DeleteGroupNoContent with default headers values
func NewDeleteGroupNoContent() *DeleteGroupNoContent {

	return &DeleteGroupNoContent{}
}

// WriteResponse to the client
func (o *DeleteGroupNoContent) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(204)
}

// DeleteGroupNotFoundCode is the HTTP code returned for type DeleteGroupNotFound
const DeleteGroupNotFoundCode int = 404

/*
DeleteGroupNotFound The specified resource was not found

swagger:response deleteGroupNotFound
*/
type DeleteGroupNotFound struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDeleteGroupNotFound creates DeleteGroupNotFound with default headers values
func NewDeleteGroupNotFound() *DeleteGroupNotFound {

	return &DeleteGroupNotFound{}
}

// WithConfigurationVersion adds the configurationVersion to the delete group not found response
func (o *DeleteGroupNotFound) WithConfigurationVersion(configurationVersion string) *DeleteGroupNotFound {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the delete group not found response
func (o *DeleteGroupNotFound) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the delete group not found response
func (o *DeleteGroupNotFound) WithPayload(payload *models.Error) *DeleteGroupNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete group not found response
func (o *DeleteGroupNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteGroupNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
DeleteGroupDefault General Error

swagger:response deleteGroupDefault
*/
type DeleteGroupDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDeleteGroupDefault creates DeleteGroupDefault with default headers values
func NewDeleteGroupDefault(code int) *DeleteGroupDefault {
	if code <= 0 {
		code = 500
	}

	return &DeleteGroupDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the delete group default response
func (o *DeleteGroupDefault) WithStatusCode(code int) *DeleteGroupDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the delete group default response
func (o *DeleteGroupDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the delete group default response
func (o *DeleteGroupDefault) WithConfigurationVersion(configurationVersion string) *DeleteGroupDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the delete group default response
func (o *DeleteGroupDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the delete group default response
func (o *DeleteGroupDefault) WithPayload(payload *models.Error) *DeleteGroupDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete group default response
func (o *DeleteGroupDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteGroupDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
