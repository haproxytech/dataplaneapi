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

package acl

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v5/models"
)

// DeleteACLAcceptedCode is the HTTP code returned for type DeleteACLAccepted
const DeleteACLAcceptedCode int = 202

/*
DeleteACLAccepted Configuration change accepted and reload requested

swagger:response deleteAclAccepted
*/
type DeleteACLAccepted struct {
	/*ID of the requested reload

	 */
	ReloadID string `json:"Reload-ID"`
}

// NewDeleteACLAccepted creates DeleteACLAccepted with default headers values
func NewDeleteACLAccepted() *DeleteACLAccepted {

	return &DeleteACLAccepted{}
}

// WithReloadID adds the reloadId to the delete Acl accepted response
func (o *DeleteACLAccepted) WithReloadID(reloadID string) *DeleteACLAccepted {
	o.ReloadID = reloadID
	return o
}

// SetReloadID sets the reloadId to the delete Acl accepted response
func (o *DeleteACLAccepted) SetReloadID(reloadID string) {
	o.ReloadID = reloadID
}

// WriteResponse to the client
func (o *DeleteACLAccepted) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	// response header Reload-ID

	reloadID := o.ReloadID
	if reloadID != "" {
		rw.Header().Set("Reload-ID", reloadID)
	}

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(202)
}

// DeleteACLNoContentCode is the HTTP code returned for type DeleteACLNoContent
const DeleteACLNoContentCode int = 204

/*
DeleteACLNoContent ACL line deleted

swagger:response deleteAclNoContent
*/
type DeleteACLNoContent struct {
}

// NewDeleteACLNoContent creates DeleteACLNoContent with default headers values
func NewDeleteACLNoContent() *DeleteACLNoContent {

	return &DeleteACLNoContent{}
}

// WriteResponse to the client
func (o *DeleteACLNoContent) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(204)
}

// DeleteACLNotFoundCode is the HTTP code returned for type DeleteACLNotFound
const DeleteACLNotFoundCode int = 404

/*
DeleteACLNotFound The specified resource was not found

swagger:response deleteAclNotFound
*/
type DeleteACLNotFound struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDeleteACLNotFound creates DeleteACLNotFound with default headers values
func NewDeleteACLNotFound() *DeleteACLNotFound {

	return &DeleteACLNotFound{}
}

// WithConfigurationVersion adds the configurationVersion to the delete Acl not found response
func (o *DeleteACLNotFound) WithConfigurationVersion(configurationVersion string) *DeleteACLNotFound {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the delete Acl not found response
func (o *DeleteACLNotFound) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the delete Acl not found response
func (o *DeleteACLNotFound) WithPayload(payload *models.Error) *DeleteACLNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete Acl not found response
func (o *DeleteACLNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteACLNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
DeleteACLDefault General Error

swagger:response deleteAclDefault
*/
type DeleteACLDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDeleteACLDefault creates DeleteACLDefault with default headers values
func NewDeleteACLDefault(code int) *DeleteACLDefault {
	if code <= 0 {
		code = 500
	}

	return &DeleteACLDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the delete Acl default response
func (o *DeleteACLDefault) WithStatusCode(code int) *DeleteACLDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the delete Acl default response
func (o *DeleteACLDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the delete Acl default response
func (o *DeleteACLDefault) WithConfigurationVersion(configurationVersion string) *DeleteACLDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the delete Acl default response
func (o *DeleteACLDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the delete Acl default response
func (o *DeleteACLDefault) WithPayload(payload *models.Error) *DeleteACLDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete Acl default response
func (o *DeleteACLDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteACLDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
