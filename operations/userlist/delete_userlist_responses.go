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

package userlist

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v6/models"
)

// DeleteUserlistAcceptedCode is the HTTP code returned for type DeleteUserlistAccepted
const DeleteUserlistAcceptedCode int = 202

/*
DeleteUserlistAccepted Configuration change accepted and reload requested

swagger:response deleteUserlistAccepted
*/
type DeleteUserlistAccepted struct {
	/*ID of the requested reload

	 */
	ReloadID string `json:"Reload-ID"`
}

// NewDeleteUserlistAccepted creates DeleteUserlistAccepted with default headers values
func NewDeleteUserlistAccepted() *DeleteUserlistAccepted {

	return &DeleteUserlistAccepted{}
}

// WithReloadID adds the reloadId to the delete userlist accepted response
func (o *DeleteUserlistAccepted) WithReloadID(reloadID string) *DeleteUserlistAccepted {
	o.ReloadID = reloadID
	return o
}

// SetReloadID sets the reloadId to the delete userlist accepted response
func (o *DeleteUserlistAccepted) SetReloadID(reloadID string) {
	o.ReloadID = reloadID
}

// WriteResponse to the client
func (o *DeleteUserlistAccepted) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	// response header Reload-ID

	reloadID := o.ReloadID
	if reloadID != "" {
		rw.Header().Set("Reload-ID", reloadID)
	}

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(202)
}

// DeleteUserlistNoContentCode is the HTTP code returned for type DeleteUserlistNoContent
const DeleteUserlistNoContentCode int = 204

/*
DeleteUserlistNoContent Userlist deleted

swagger:response deleteUserlistNoContent
*/
type DeleteUserlistNoContent struct {
}

// NewDeleteUserlistNoContent creates DeleteUserlistNoContent with default headers values
func NewDeleteUserlistNoContent() *DeleteUserlistNoContent {

	return &DeleteUserlistNoContent{}
}

// WriteResponse to the client
func (o *DeleteUserlistNoContent) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(204)
}

// DeleteUserlistNotFoundCode is the HTTP code returned for type DeleteUserlistNotFound
const DeleteUserlistNotFoundCode int = 404

/*
DeleteUserlistNotFound The specified resource was not found

swagger:response deleteUserlistNotFound
*/
type DeleteUserlistNotFound struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDeleteUserlistNotFound creates DeleteUserlistNotFound with default headers values
func NewDeleteUserlistNotFound() *DeleteUserlistNotFound {

	return &DeleteUserlistNotFound{}
}

// WithConfigurationVersion adds the configurationVersion to the delete userlist not found response
func (o *DeleteUserlistNotFound) WithConfigurationVersion(configurationVersion string) *DeleteUserlistNotFound {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the delete userlist not found response
func (o *DeleteUserlistNotFound) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the delete userlist not found response
func (o *DeleteUserlistNotFound) WithPayload(payload *models.Error) *DeleteUserlistNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete userlist not found response
func (o *DeleteUserlistNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteUserlistNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
DeleteUserlistDefault General Error

swagger:response deleteUserlistDefault
*/
type DeleteUserlistDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDeleteUserlistDefault creates DeleteUserlistDefault with default headers values
func NewDeleteUserlistDefault(code int) *DeleteUserlistDefault {
	if code <= 0 {
		code = 500
	}

	return &DeleteUserlistDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the delete userlist default response
func (o *DeleteUserlistDefault) WithStatusCode(code int) *DeleteUserlistDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the delete userlist default response
func (o *DeleteUserlistDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the delete userlist default response
func (o *DeleteUserlistDefault) WithConfigurationVersion(configurationVersion string) *DeleteUserlistDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the delete userlist default response
func (o *DeleteUserlistDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the delete userlist default response
func (o *DeleteUserlistDefault) WithPayload(payload *models.Error) *DeleteUserlistDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete userlist default response
func (o *DeleteUserlistDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteUserlistDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
