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

package s_s_l_front_use

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v6/models"
)

// DeleteSSLFrontUseAcceptedCode is the HTTP code returned for type DeleteSSLFrontUseAccepted
const DeleteSSLFrontUseAcceptedCode int = 202

/*
DeleteSSLFrontUseAccepted Configuration change accepted and reload requested

swagger:response deleteSSLFrontUseAccepted
*/
type DeleteSSLFrontUseAccepted struct {
	/*ID of the requested reload

	 */
	ReloadID string `json:"Reload-ID"`
}

// NewDeleteSSLFrontUseAccepted creates DeleteSSLFrontUseAccepted with default headers values
func NewDeleteSSLFrontUseAccepted() *DeleteSSLFrontUseAccepted {

	return &DeleteSSLFrontUseAccepted{}
}

// WithReloadID adds the reloadId to the delete s s l front use accepted response
func (o *DeleteSSLFrontUseAccepted) WithReloadID(reloadID string) *DeleteSSLFrontUseAccepted {
	o.ReloadID = reloadID
	return o
}

// SetReloadID sets the reloadId to the delete s s l front use accepted response
func (o *DeleteSSLFrontUseAccepted) SetReloadID(reloadID string) {
	o.ReloadID = reloadID
}

// WriteResponse to the client
func (o *DeleteSSLFrontUseAccepted) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	// response header Reload-ID

	reloadID := o.ReloadID
	if reloadID != "" {
		rw.Header().Set("Reload-ID", reloadID)
	}

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(202)
}

// DeleteSSLFrontUseNoContentCode is the HTTP code returned for type DeleteSSLFrontUseNoContent
const DeleteSSLFrontUseNoContentCode int = 204

/*
DeleteSSLFrontUseNoContent SSLFrontUse deleted

swagger:response deleteSSLFrontUseNoContent
*/
type DeleteSSLFrontUseNoContent struct {
}

// NewDeleteSSLFrontUseNoContent creates DeleteSSLFrontUseNoContent with default headers values
func NewDeleteSSLFrontUseNoContent() *DeleteSSLFrontUseNoContent {

	return &DeleteSSLFrontUseNoContent{}
}

// WriteResponse to the client
func (o *DeleteSSLFrontUseNoContent) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(204)
}

// DeleteSSLFrontUseNotFoundCode is the HTTP code returned for type DeleteSSLFrontUseNotFound
const DeleteSSLFrontUseNotFoundCode int = 404

/*
DeleteSSLFrontUseNotFound The specified resource was not found

swagger:response deleteSSLFrontUseNotFound
*/
type DeleteSSLFrontUseNotFound struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDeleteSSLFrontUseNotFound creates DeleteSSLFrontUseNotFound with default headers values
func NewDeleteSSLFrontUseNotFound() *DeleteSSLFrontUseNotFound {

	return &DeleteSSLFrontUseNotFound{}
}

// WithConfigurationVersion adds the configurationVersion to the delete s s l front use not found response
func (o *DeleteSSLFrontUseNotFound) WithConfigurationVersion(configurationVersion string) *DeleteSSLFrontUseNotFound {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the delete s s l front use not found response
func (o *DeleteSSLFrontUseNotFound) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the delete s s l front use not found response
func (o *DeleteSSLFrontUseNotFound) WithPayload(payload *models.Error) *DeleteSSLFrontUseNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete s s l front use not found response
func (o *DeleteSSLFrontUseNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteSSLFrontUseNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
DeleteSSLFrontUseDefault General Error

swagger:response deleteSSLFrontUseDefault
*/
type DeleteSSLFrontUseDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDeleteSSLFrontUseDefault creates DeleteSSLFrontUseDefault with default headers values
func NewDeleteSSLFrontUseDefault(code int) *DeleteSSLFrontUseDefault {
	if code <= 0 {
		code = 500
	}

	return &DeleteSSLFrontUseDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the delete s s l front use default response
func (o *DeleteSSLFrontUseDefault) WithStatusCode(code int) *DeleteSSLFrontUseDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the delete s s l front use default response
func (o *DeleteSSLFrontUseDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the delete s s l front use default response
func (o *DeleteSSLFrontUseDefault) WithConfigurationVersion(configurationVersion string) *DeleteSSLFrontUseDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the delete s s l front use default response
func (o *DeleteSSLFrontUseDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the delete s s l front use default response
func (o *DeleteSSLFrontUseDefault) WithPayload(payload *models.Error) *DeleteSSLFrontUseDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete s s l front use default response
func (o *DeleteSSLFrontUseDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteSSLFrontUseDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
