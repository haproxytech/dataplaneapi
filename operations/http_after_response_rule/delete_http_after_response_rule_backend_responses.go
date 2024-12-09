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

package http_after_response_rule

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v6/models"
)

// DeleteHTTPAfterResponseRuleBackendAcceptedCode is the HTTP code returned for type DeleteHTTPAfterResponseRuleBackendAccepted
const DeleteHTTPAfterResponseRuleBackendAcceptedCode int = 202

/*
DeleteHTTPAfterResponseRuleBackendAccepted Configuration change accepted and reload requested

swagger:response deleteHttpAfterResponseRuleBackendAccepted
*/
type DeleteHTTPAfterResponseRuleBackendAccepted struct {
	/*ID of the requested reload

	 */
	ReloadID string `json:"Reload-ID"`
}

// NewDeleteHTTPAfterResponseRuleBackendAccepted creates DeleteHTTPAfterResponseRuleBackendAccepted with default headers values
func NewDeleteHTTPAfterResponseRuleBackendAccepted() *DeleteHTTPAfterResponseRuleBackendAccepted {

	return &DeleteHTTPAfterResponseRuleBackendAccepted{}
}

// WithReloadID adds the reloadId to the delete Http after response rule backend accepted response
func (o *DeleteHTTPAfterResponseRuleBackendAccepted) WithReloadID(reloadID string) *DeleteHTTPAfterResponseRuleBackendAccepted {
	o.ReloadID = reloadID
	return o
}

// SetReloadID sets the reloadId to the delete Http after response rule backend accepted response
func (o *DeleteHTTPAfterResponseRuleBackendAccepted) SetReloadID(reloadID string) {
	o.ReloadID = reloadID
}

// WriteResponse to the client
func (o *DeleteHTTPAfterResponseRuleBackendAccepted) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	// response header Reload-ID

	reloadID := o.ReloadID
	if reloadID != "" {
		rw.Header().Set("Reload-ID", reloadID)
	}

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(202)
}

// DeleteHTTPAfterResponseRuleBackendNoContentCode is the HTTP code returned for type DeleteHTTPAfterResponseRuleBackendNoContent
const DeleteHTTPAfterResponseRuleBackendNoContentCode int = 204

/*
DeleteHTTPAfterResponseRuleBackendNoContent HTTP After Response Rule deleted

swagger:response deleteHttpAfterResponseRuleBackendNoContent
*/
type DeleteHTTPAfterResponseRuleBackendNoContent struct {
}

// NewDeleteHTTPAfterResponseRuleBackendNoContent creates DeleteHTTPAfterResponseRuleBackendNoContent with default headers values
func NewDeleteHTTPAfterResponseRuleBackendNoContent() *DeleteHTTPAfterResponseRuleBackendNoContent {

	return &DeleteHTTPAfterResponseRuleBackendNoContent{}
}

// WriteResponse to the client
func (o *DeleteHTTPAfterResponseRuleBackendNoContent) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(204)
}

// DeleteHTTPAfterResponseRuleBackendNotFoundCode is the HTTP code returned for type DeleteHTTPAfterResponseRuleBackendNotFound
const DeleteHTTPAfterResponseRuleBackendNotFoundCode int = 404

/*
DeleteHTTPAfterResponseRuleBackendNotFound The specified resource was not found

swagger:response deleteHttpAfterResponseRuleBackendNotFound
*/
type DeleteHTTPAfterResponseRuleBackendNotFound struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDeleteHTTPAfterResponseRuleBackendNotFound creates DeleteHTTPAfterResponseRuleBackendNotFound with default headers values
func NewDeleteHTTPAfterResponseRuleBackendNotFound() *DeleteHTTPAfterResponseRuleBackendNotFound {

	return &DeleteHTTPAfterResponseRuleBackendNotFound{}
}

// WithConfigurationVersion adds the configurationVersion to the delete Http after response rule backend not found response
func (o *DeleteHTTPAfterResponseRuleBackendNotFound) WithConfigurationVersion(configurationVersion string) *DeleteHTTPAfterResponseRuleBackendNotFound {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the delete Http after response rule backend not found response
func (o *DeleteHTTPAfterResponseRuleBackendNotFound) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the delete Http after response rule backend not found response
func (o *DeleteHTTPAfterResponseRuleBackendNotFound) WithPayload(payload *models.Error) *DeleteHTTPAfterResponseRuleBackendNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete Http after response rule backend not found response
func (o *DeleteHTTPAfterResponseRuleBackendNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteHTTPAfterResponseRuleBackendNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
DeleteHTTPAfterResponseRuleBackendDefault General Error

swagger:response deleteHttpAfterResponseRuleBackendDefault
*/
type DeleteHTTPAfterResponseRuleBackendDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDeleteHTTPAfterResponseRuleBackendDefault creates DeleteHTTPAfterResponseRuleBackendDefault with default headers values
func NewDeleteHTTPAfterResponseRuleBackendDefault(code int) *DeleteHTTPAfterResponseRuleBackendDefault {
	if code <= 0 {
		code = 500
	}

	return &DeleteHTTPAfterResponseRuleBackendDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the delete HTTP after response rule backend default response
func (o *DeleteHTTPAfterResponseRuleBackendDefault) WithStatusCode(code int) *DeleteHTTPAfterResponseRuleBackendDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the delete HTTP after response rule backend default response
func (o *DeleteHTTPAfterResponseRuleBackendDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the delete HTTP after response rule backend default response
func (o *DeleteHTTPAfterResponseRuleBackendDefault) WithConfigurationVersion(configurationVersion string) *DeleteHTTPAfterResponseRuleBackendDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the delete HTTP after response rule backend default response
func (o *DeleteHTTPAfterResponseRuleBackendDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the delete HTTP after response rule backend default response
func (o *DeleteHTTPAfterResponseRuleBackendDefault) WithPayload(payload *models.Error) *DeleteHTTPAfterResponseRuleBackendDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete HTTP after response rule backend default response
func (o *DeleteHTTPAfterResponseRuleBackendDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteHTTPAfterResponseRuleBackendDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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