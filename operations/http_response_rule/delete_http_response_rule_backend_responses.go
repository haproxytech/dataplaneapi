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

package http_response_rule

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v6/models"
)

// DeleteHTTPResponseRuleBackendAcceptedCode is the HTTP code returned for type DeleteHTTPResponseRuleBackendAccepted
const DeleteHTTPResponseRuleBackendAcceptedCode int = 202

/*
DeleteHTTPResponseRuleBackendAccepted Configuration change accepted and reload requested

swagger:response deleteHttpResponseRuleBackendAccepted
*/
type DeleteHTTPResponseRuleBackendAccepted struct {
	/*ID of the requested reload

	 */
	ReloadID string `json:"Reload-ID"`
}

// NewDeleteHTTPResponseRuleBackendAccepted creates DeleteHTTPResponseRuleBackendAccepted with default headers values
func NewDeleteHTTPResponseRuleBackendAccepted() *DeleteHTTPResponseRuleBackendAccepted {

	return &DeleteHTTPResponseRuleBackendAccepted{}
}

// WithReloadID adds the reloadId to the delete Http response rule backend accepted response
func (o *DeleteHTTPResponseRuleBackendAccepted) WithReloadID(reloadID string) *DeleteHTTPResponseRuleBackendAccepted {
	o.ReloadID = reloadID
	return o
}

// SetReloadID sets the reloadId to the delete Http response rule backend accepted response
func (o *DeleteHTTPResponseRuleBackendAccepted) SetReloadID(reloadID string) {
	o.ReloadID = reloadID
}

// WriteResponse to the client
func (o *DeleteHTTPResponseRuleBackendAccepted) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	// response header Reload-ID

	reloadID := o.ReloadID
	if reloadID != "" {
		rw.Header().Set("Reload-ID", reloadID)
	}

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(202)
}

// DeleteHTTPResponseRuleBackendNoContentCode is the HTTP code returned for type DeleteHTTPResponseRuleBackendNoContent
const DeleteHTTPResponseRuleBackendNoContentCode int = 204

/*
DeleteHTTPResponseRuleBackendNoContent HTTP Response Rule deleted

swagger:response deleteHttpResponseRuleBackendNoContent
*/
type DeleteHTTPResponseRuleBackendNoContent struct {
}

// NewDeleteHTTPResponseRuleBackendNoContent creates DeleteHTTPResponseRuleBackendNoContent with default headers values
func NewDeleteHTTPResponseRuleBackendNoContent() *DeleteHTTPResponseRuleBackendNoContent {

	return &DeleteHTTPResponseRuleBackendNoContent{}
}

// WriteResponse to the client
func (o *DeleteHTTPResponseRuleBackendNoContent) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(204)
}

// DeleteHTTPResponseRuleBackendNotFoundCode is the HTTP code returned for type DeleteHTTPResponseRuleBackendNotFound
const DeleteHTTPResponseRuleBackendNotFoundCode int = 404

/*
DeleteHTTPResponseRuleBackendNotFound The specified resource was not found

swagger:response deleteHttpResponseRuleBackendNotFound
*/
type DeleteHTTPResponseRuleBackendNotFound struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDeleteHTTPResponseRuleBackendNotFound creates DeleteHTTPResponseRuleBackendNotFound with default headers values
func NewDeleteHTTPResponseRuleBackendNotFound() *DeleteHTTPResponseRuleBackendNotFound {

	return &DeleteHTTPResponseRuleBackendNotFound{}
}

// WithConfigurationVersion adds the configurationVersion to the delete Http response rule backend not found response
func (o *DeleteHTTPResponseRuleBackendNotFound) WithConfigurationVersion(configurationVersion string) *DeleteHTTPResponseRuleBackendNotFound {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the delete Http response rule backend not found response
func (o *DeleteHTTPResponseRuleBackendNotFound) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the delete Http response rule backend not found response
func (o *DeleteHTTPResponseRuleBackendNotFound) WithPayload(payload *models.Error) *DeleteHTTPResponseRuleBackendNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete Http response rule backend not found response
func (o *DeleteHTTPResponseRuleBackendNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteHTTPResponseRuleBackendNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
DeleteHTTPResponseRuleBackendDefault General Error

swagger:response deleteHttpResponseRuleBackendDefault
*/
type DeleteHTTPResponseRuleBackendDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDeleteHTTPResponseRuleBackendDefault creates DeleteHTTPResponseRuleBackendDefault with default headers values
func NewDeleteHTTPResponseRuleBackendDefault(code int) *DeleteHTTPResponseRuleBackendDefault {
	if code <= 0 {
		code = 500
	}

	return &DeleteHTTPResponseRuleBackendDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the delete HTTP response rule backend default response
func (o *DeleteHTTPResponseRuleBackendDefault) WithStatusCode(code int) *DeleteHTTPResponseRuleBackendDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the delete HTTP response rule backend default response
func (o *DeleteHTTPResponseRuleBackendDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the delete HTTP response rule backend default response
func (o *DeleteHTTPResponseRuleBackendDefault) WithConfigurationVersion(configurationVersion string) *DeleteHTTPResponseRuleBackendDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the delete HTTP response rule backend default response
func (o *DeleteHTTPResponseRuleBackendDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the delete HTTP response rule backend default response
func (o *DeleteHTTPResponseRuleBackendDefault) WithPayload(payload *models.Error) *DeleteHTTPResponseRuleBackendDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete HTTP response rule backend default response
func (o *DeleteHTTPResponseRuleBackendDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteHTTPResponseRuleBackendDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
