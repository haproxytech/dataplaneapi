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

package http_error_rule

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v6/models"
)

// DeleteHTTPErrorRuleAcceptedCode is the HTTP code returned for type DeleteHTTPErrorRuleAccepted
const DeleteHTTPErrorRuleAcceptedCode int = 202

/*
DeleteHTTPErrorRuleAccepted Configuration change accepted and reload requested

swagger:response deleteHttpErrorRuleAccepted
*/
type DeleteHTTPErrorRuleAccepted struct {
	/*ID of the requested reload

	 */
	ReloadID string `json:"Reload-ID"`
}

// NewDeleteHTTPErrorRuleAccepted creates DeleteHTTPErrorRuleAccepted with default headers values
func NewDeleteHTTPErrorRuleAccepted() *DeleteHTTPErrorRuleAccepted {

	return &DeleteHTTPErrorRuleAccepted{}
}

// WithReloadID adds the reloadId to the delete Http error rule accepted response
func (o *DeleteHTTPErrorRuleAccepted) WithReloadID(reloadID string) *DeleteHTTPErrorRuleAccepted {
	o.ReloadID = reloadID
	return o
}

// SetReloadID sets the reloadId to the delete Http error rule accepted response
func (o *DeleteHTTPErrorRuleAccepted) SetReloadID(reloadID string) {
	o.ReloadID = reloadID
}

// WriteResponse to the client
func (o *DeleteHTTPErrorRuleAccepted) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	// response header Reload-ID

	reloadID := o.ReloadID
	if reloadID != "" {
		rw.Header().Set("Reload-ID", reloadID)
	}

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(202)
}

// DeleteHTTPErrorRuleNoContentCode is the HTTP code returned for type DeleteHTTPErrorRuleNoContent
const DeleteHTTPErrorRuleNoContentCode int = 204

/*
DeleteHTTPErrorRuleNoContent HTTP Error Rule deleted

swagger:response deleteHttpErrorRuleNoContent
*/
type DeleteHTTPErrorRuleNoContent struct {
}

// NewDeleteHTTPErrorRuleNoContent creates DeleteHTTPErrorRuleNoContent with default headers values
func NewDeleteHTTPErrorRuleNoContent() *DeleteHTTPErrorRuleNoContent {

	return &DeleteHTTPErrorRuleNoContent{}
}

// WriteResponse to the client
func (o *DeleteHTTPErrorRuleNoContent) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(204)
}

// DeleteHTTPErrorRuleNotFoundCode is the HTTP code returned for type DeleteHTTPErrorRuleNotFound
const DeleteHTTPErrorRuleNotFoundCode int = 404

/*
DeleteHTTPErrorRuleNotFound The specified resource was not found

swagger:response deleteHttpErrorRuleNotFound
*/
type DeleteHTTPErrorRuleNotFound struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDeleteHTTPErrorRuleNotFound creates DeleteHTTPErrorRuleNotFound with default headers values
func NewDeleteHTTPErrorRuleNotFound() *DeleteHTTPErrorRuleNotFound {

	return &DeleteHTTPErrorRuleNotFound{}
}

// WithConfigurationVersion adds the configurationVersion to the delete Http error rule not found response
func (o *DeleteHTTPErrorRuleNotFound) WithConfigurationVersion(configurationVersion string) *DeleteHTTPErrorRuleNotFound {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the delete Http error rule not found response
func (o *DeleteHTTPErrorRuleNotFound) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the delete Http error rule not found response
func (o *DeleteHTTPErrorRuleNotFound) WithPayload(payload *models.Error) *DeleteHTTPErrorRuleNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete Http error rule not found response
func (o *DeleteHTTPErrorRuleNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteHTTPErrorRuleNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
DeleteHTTPErrorRuleDefault General Error

swagger:response deleteHttpErrorRuleDefault
*/
type DeleteHTTPErrorRuleDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDeleteHTTPErrorRuleDefault creates DeleteHTTPErrorRuleDefault with default headers values
func NewDeleteHTTPErrorRuleDefault(code int) *DeleteHTTPErrorRuleDefault {
	if code <= 0 {
		code = 500
	}

	return &DeleteHTTPErrorRuleDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the delete HTTP error rule default response
func (o *DeleteHTTPErrorRuleDefault) WithStatusCode(code int) *DeleteHTTPErrorRuleDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the delete HTTP error rule default response
func (o *DeleteHTTPErrorRuleDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the delete HTTP error rule default response
func (o *DeleteHTTPErrorRuleDefault) WithConfigurationVersion(configurationVersion string) *DeleteHTTPErrorRuleDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the delete HTTP error rule default response
func (o *DeleteHTTPErrorRuleDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the delete HTTP error rule default response
func (o *DeleteHTTPErrorRuleDefault) WithPayload(payload *models.Error) *DeleteHTTPErrorRuleDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete HTTP error rule default response
func (o *DeleteHTTPErrorRuleDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteHTTPErrorRuleDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
