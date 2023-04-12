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

	"github.com/haproxytech/client-native/v4/models"
)

// DeleteTableAcceptedCode is the HTTP code returned for type DeleteTableAccepted
const DeleteTableAcceptedCode int = 202

/*
DeleteTableAccepted Configuration change accepted and reload requested

swagger:response deleteTableAccepted
*/
type DeleteTableAccepted struct {
	/*ID of the requested reload

	 */
	ReloadID string `json:"Reload-ID"`
}

// NewDeleteTableAccepted creates DeleteTableAccepted with default headers values
func NewDeleteTableAccepted() *DeleteTableAccepted {

	return &DeleteTableAccepted{}
}

// WithReloadID adds the reloadId to the delete table accepted response
func (o *DeleteTableAccepted) WithReloadID(reloadID string) *DeleteTableAccepted {
	o.ReloadID = reloadID
	return o
}

// SetReloadID sets the reloadId to the delete table accepted response
func (o *DeleteTableAccepted) SetReloadID(reloadID string) {
	o.ReloadID = reloadID
}

// WriteResponse to the client
func (o *DeleteTableAccepted) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	// response header Reload-ID

	reloadID := o.ReloadID
	if reloadID != "" {
		rw.Header().Set("Reload-ID", reloadID)
	}

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(202)
}

// DeleteTableNoContentCode is the HTTP code returned for type DeleteTableNoContent
const DeleteTableNoContentCode int = 204

/*
DeleteTableNoContent Table deleted

swagger:response deleteTableNoContent
*/
type DeleteTableNoContent struct {
}

// NewDeleteTableNoContent creates DeleteTableNoContent with default headers values
func NewDeleteTableNoContent() *DeleteTableNoContent {

	return &DeleteTableNoContent{}
}

// WriteResponse to the client
func (o *DeleteTableNoContent) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(204)
}

// DeleteTableNotFoundCode is the HTTP code returned for type DeleteTableNotFound
const DeleteTableNotFoundCode int = 404

/*
DeleteTableNotFound The specified resource was not found

swagger:response deleteTableNotFound
*/
type DeleteTableNotFound struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDeleteTableNotFound creates DeleteTableNotFound with default headers values
func NewDeleteTableNotFound() *DeleteTableNotFound {

	return &DeleteTableNotFound{}
}

// WithConfigurationVersion adds the configurationVersion to the delete table not found response
func (o *DeleteTableNotFound) WithConfigurationVersion(configurationVersion string) *DeleteTableNotFound {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the delete table not found response
func (o *DeleteTableNotFound) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the delete table not found response
func (o *DeleteTableNotFound) WithPayload(payload *models.Error) *DeleteTableNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete table not found response
func (o *DeleteTableNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteTableNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
DeleteTableDefault General Error

swagger:response deleteTableDefault
*/
type DeleteTableDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDeleteTableDefault creates DeleteTableDefault with default headers values
func NewDeleteTableDefault(code int) *DeleteTableDefault {
	if code <= 0 {
		code = 500
	}

	return &DeleteTableDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the delete table default response
func (o *DeleteTableDefault) WithStatusCode(code int) *DeleteTableDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the delete table default response
func (o *DeleteTableDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the delete table default response
func (o *DeleteTableDefault) WithConfigurationVersion(configurationVersion string) *DeleteTableDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the delete table default response
func (o *DeleteTableDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the delete table default response
func (o *DeleteTableDefault) WithPayload(payload *models.Error) *DeleteTableDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete table default response
func (o *DeleteTableDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteTableDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
