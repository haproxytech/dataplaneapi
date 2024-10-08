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

package tcp_check

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v6/models"
)

// DeleteTCPCheckDefaultsAcceptedCode is the HTTP code returned for type DeleteTCPCheckDefaultsAccepted
const DeleteTCPCheckDefaultsAcceptedCode int = 202

/*
DeleteTCPCheckDefaultsAccepted Configuration change accepted and reload requested

swagger:response deleteTcpCheckDefaultsAccepted
*/
type DeleteTCPCheckDefaultsAccepted struct {
	/*ID of the requested reload

	 */
	ReloadID string `json:"Reload-ID"`
}

// NewDeleteTCPCheckDefaultsAccepted creates DeleteTCPCheckDefaultsAccepted with default headers values
func NewDeleteTCPCheckDefaultsAccepted() *DeleteTCPCheckDefaultsAccepted {

	return &DeleteTCPCheckDefaultsAccepted{}
}

// WithReloadID adds the reloadId to the delete Tcp check defaults accepted response
func (o *DeleteTCPCheckDefaultsAccepted) WithReloadID(reloadID string) *DeleteTCPCheckDefaultsAccepted {
	o.ReloadID = reloadID
	return o
}

// SetReloadID sets the reloadId to the delete Tcp check defaults accepted response
func (o *DeleteTCPCheckDefaultsAccepted) SetReloadID(reloadID string) {
	o.ReloadID = reloadID
}

// WriteResponse to the client
func (o *DeleteTCPCheckDefaultsAccepted) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	// response header Reload-ID

	reloadID := o.ReloadID
	if reloadID != "" {
		rw.Header().Set("Reload-ID", reloadID)
	}

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(202)
}

// DeleteTCPCheckDefaultsNoContentCode is the HTTP code returned for type DeleteTCPCheckDefaultsNoContent
const DeleteTCPCheckDefaultsNoContentCode int = 204

/*
DeleteTCPCheckDefaultsNoContent TCP check deleted

swagger:response deleteTcpCheckDefaultsNoContent
*/
type DeleteTCPCheckDefaultsNoContent struct {
}

// NewDeleteTCPCheckDefaultsNoContent creates DeleteTCPCheckDefaultsNoContent with default headers values
func NewDeleteTCPCheckDefaultsNoContent() *DeleteTCPCheckDefaultsNoContent {

	return &DeleteTCPCheckDefaultsNoContent{}
}

// WriteResponse to the client
func (o *DeleteTCPCheckDefaultsNoContent) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(204)
}

// DeleteTCPCheckDefaultsNotFoundCode is the HTTP code returned for type DeleteTCPCheckDefaultsNotFound
const DeleteTCPCheckDefaultsNotFoundCode int = 404

/*
DeleteTCPCheckDefaultsNotFound The specified resource was not found

swagger:response deleteTcpCheckDefaultsNotFound
*/
type DeleteTCPCheckDefaultsNotFound struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDeleteTCPCheckDefaultsNotFound creates DeleteTCPCheckDefaultsNotFound with default headers values
func NewDeleteTCPCheckDefaultsNotFound() *DeleteTCPCheckDefaultsNotFound {

	return &DeleteTCPCheckDefaultsNotFound{}
}

// WithConfigurationVersion adds the configurationVersion to the delete Tcp check defaults not found response
func (o *DeleteTCPCheckDefaultsNotFound) WithConfigurationVersion(configurationVersion string) *DeleteTCPCheckDefaultsNotFound {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the delete Tcp check defaults not found response
func (o *DeleteTCPCheckDefaultsNotFound) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the delete Tcp check defaults not found response
func (o *DeleteTCPCheckDefaultsNotFound) WithPayload(payload *models.Error) *DeleteTCPCheckDefaultsNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete Tcp check defaults not found response
func (o *DeleteTCPCheckDefaultsNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteTCPCheckDefaultsNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
DeleteTCPCheckDefaultsDefault General Error

swagger:response deleteTcpCheckDefaultsDefault
*/
type DeleteTCPCheckDefaultsDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDeleteTCPCheckDefaultsDefault creates DeleteTCPCheckDefaultsDefault with default headers values
func NewDeleteTCPCheckDefaultsDefault(code int) *DeleteTCPCheckDefaultsDefault {
	if code <= 0 {
		code = 500
	}

	return &DeleteTCPCheckDefaultsDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the delete TCP check defaults default response
func (o *DeleteTCPCheckDefaultsDefault) WithStatusCode(code int) *DeleteTCPCheckDefaultsDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the delete TCP check defaults default response
func (o *DeleteTCPCheckDefaultsDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the delete TCP check defaults default response
func (o *DeleteTCPCheckDefaultsDefault) WithConfigurationVersion(configurationVersion string) *DeleteTCPCheckDefaultsDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the delete TCP check defaults default response
func (o *DeleteTCPCheckDefaultsDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the delete TCP check defaults default response
func (o *DeleteTCPCheckDefaultsDefault) WithPayload(payload *models.Error) *DeleteTCPCheckDefaultsDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete TCP check defaults default response
func (o *DeleteTCPCheckDefaultsDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteTCPCheckDefaultsDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
