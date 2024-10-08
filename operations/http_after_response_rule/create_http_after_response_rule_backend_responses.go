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

// CreateHTTPAfterResponseRuleBackendCreatedCode is the HTTP code returned for type CreateHTTPAfterResponseRuleBackendCreated
const CreateHTTPAfterResponseRuleBackendCreatedCode int = 201

/*
CreateHTTPAfterResponseRuleBackendCreated HTTP Response Rule created

swagger:response createHttpAfterResponseRuleBackendCreated
*/
type CreateHTTPAfterResponseRuleBackendCreated struct {

	/*
	  In: Body
	*/
	Payload *models.HTTPAfterResponseRule `json:"body,omitempty"`
}

// NewCreateHTTPAfterResponseRuleBackendCreated creates CreateHTTPAfterResponseRuleBackendCreated with default headers values
func NewCreateHTTPAfterResponseRuleBackendCreated() *CreateHTTPAfterResponseRuleBackendCreated {

	return &CreateHTTPAfterResponseRuleBackendCreated{}
}

// WithPayload adds the payload to the create Http after response rule backend created response
func (o *CreateHTTPAfterResponseRuleBackendCreated) WithPayload(payload *models.HTTPAfterResponseRule) *CreateHTTPAfterResponseRuleBackendCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create Http after response rule backend created response
func (o *CreateHTTPAfterResponseRuleBackendCreated) SetPayload(payload *models.HTTPAfterResponseRule) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateHTTPAfterResponseRuleBackendCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateHTTPAfterResponseRuleBackendAcceptedCode is the HTTP code returned for type CreateHTTPAfterResponseRuleBackendAccepted
const CreateHTTPAfterResponseRuleBackendAcceptedCode int = 202

/*
CreateHTTPAfterResponseRuleBackendAccepted Configuration change accepted and reload requested

swagger:response createHttpAfterResponseRuleBackendAccepted
*/
type CreateHTTPAfterResponseRuleBackendAccepted struct {
	/*ID of the requested reload

	 */
	ReloadID string `json:"Reload-ID"`

	/*
	  In: Body
	*/
	Payload *models.HTTPAfterResponseRule `json:"body,omitempty"`
}

// NewCreateHTTPAfterResponseRuleBackendAccepted creates CreateHTTPAfterResponseRuleBackendAccepted with default headers values
func NewCreateHTTPAfterResponseRuleBackendAccepted() *CreateHTTPAfterResponseRuleBackendAccepted {

	return &CreateHTTPAfterResponseRuleBackendAccepted{}
}

// WithReloadID adds the reloadId to the create Http after response rule backend accepted response
func (o *CreateHTTPAfterResponseRuleBackendAccepted) WithReloadID(reloadID string) *CreateHTTPAfterResponseRuleBackendAccepted {
	o.ReloadID = reloadID
	return o
}

// SetReloadID sets the reloadId to the create Http after response rule backend accepted response
func (o *CreateHTTPAfterResponseRuleBackendAccepted) SetReloadID(reloadID string) {
	o.ReloadID = reloadID
}

// WithPayload adds the payload to the create Http after response rule backend accepted response
func (o *CreateHTTPAfterResponseRuleBackendAccepted) WithPayload(payload *models.HTTPAfterResponseRule) *CreateHTTPAfterResponseRuleBackendAccepted {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create Http after response rule backend accepted response
func (o *CreateHTTPAfterResponseRuleBackendAccepted) SetPayload(payload *models.HTTPAfterResponseRule) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateHTTPAfterResponseRuleBackendAccepted) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	// response header Reload-ID

	reloadID := o.ReloadID
	if reloadID != "" {
		rw.Header().Set("Reload-ID", reloadID)
	}

	rw.WriteHeader(202)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateHTTPAfterResponseRuleBackendBadRequestCode is the HTTP code returned for type CreateHTTPAfterResponseRuleBackendBadRequest
const CreateHTTPAfterResponseRuleBackendBadRequestCode int = 400

/*
CreateHTTPAfterResponseRuleBackendBadRequest Bad request

swagger:response createHttpAfterResponseRuleBackendBadRequest
*/
type CreateHTTPAfterResponseRuleBackendBadRequest struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateHTTPAfterResponseRuleBackendBadRequest creates CreateHTTPAfterResponseRuleBackendBadRequest with default headers values
func NewCreateHTTPAfterResponseRuleBackendBadRequest() *CreateHTTPAfterResponseRuleBackendBadRequest {

	return &CreateHTTPAfterResponseRuleBackendBadRequest{}
}

// WithConfigurationVersion adds the configurationVersion to the create Http after response rule backend bad request response
func (o *CreateHTTPAfterResponseRuleBackendBadRequest) WithConfigurationVersion(configurationVersion string) *CreateHTTPAfterResponseRuleBackendBadRequest {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the create Http after response rule backend bad request response
func (o *CreateHTTPAfterResponseRuleBackendBadRequest) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the create Http after response rule backend bad request response
func (o *CreateHTTPAfterResponseRuleBackendBadRequest) WithPayload(payload *models.Error) *CreateHTTPAfterResponseRuleBackendBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create Http after response rule backend bad request response
func (o *CreateHTTPAfterResponseRuleBackendBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateHTTPAfterResponseRuleBackendBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	// response header Configuration-Version

	configurationVersion := o.ConfigurationVersion
	if configurationVersion != "" {
		rw.Header().Set("Configuration-Version", configurationVersion)
	}

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateHTTPAfterResponseRuleBackendConflictCode is the HTTP code returned for type CreateHTTPAfterResponseRuleBackendConflict
const CreateHTTPAfterResponseRuleBackendConflictCode int = 409

/*
CreateHTTPAfterResponseRuleBackendConflict The specified resource already exists

swagger:response createHttpAfterResponseRuleBackendConflict
*/
type CreateHTTPAfterResponseRuleBackendConflict struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateHTTPAfterResponseRuleBackendConflict creates CreateHTTPAfterResponseRuleBackendConflict with default headers values
func NewCreateHTTPAfterResponseRuleBackendConflict() *CreateHTTPAfterResponseRuleBackendConflict {

	return &CreateHTTPAfterResponseRuleBackendConflict{}
}

// WithConfigurationVersion adds the configurationVersion to the create Http after response rule backend conflict response
func (o *CreateHTTPAfterResponseRuleBackendConflict) WithConfigurationVersion(configurationVersion string) *CreateHTTPAfterResponseRuleBackendConflict {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the create Http after response rule backend conflict response
func (o *CreateHTTPAfterResponseRuleBackendConflict) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the create Http after response rule backend conflict response
func (o *CreateHTTPAfterResponseRuleBackendConflict) WithPayload(payload *models.Error) *CreateHTTPAfterResponseRuleBackendConflict {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create Http after response rule backend conflict response
func (o *CreateHTTPAfterResponseRuleBackendConflict) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateHTTPAfterResponseRuleBackendConflict) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	// response header Configuration-Version

	configurationVersion := o.ConfigurationVersion
	if configurationVersion != "" {
		rw.Header().Set("Configuration-Version", configurationVersion)
	}

	rw.WriteHeader(409)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*
CreateHTTPAfterResponseRuleBackendDefault General Error

swagger:response createHttpAfterResponseRuleBackendDefault
*/
type CreateHTTPAfterResponseRuleBackendDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateHTTPAfterResponseRuleBackendDefault creates CreateHTTPAfterResponseRuleBackendDefault with default headers values
func NewCreateHTTPAfterResponseRuleBackendDefault(code int) *CreateHTTPAfterResponseRuleBackendDefault {
	if code <= 0 {
		code = 500
	}

	return &CreateHTTPAfterResponseRuleBackendDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the create HTTP after response rule backend default response
func (o *CreateHTTPAfterResponseRuleBackendDefault) WithStatusCode(code int) *CreateHTTPAfterResponseRuleBackendDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the create HTTP after response rule backend default response
func (o *CreateHTTPAfterResponseRuleBackendDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the create HTTP after response rule backend default response
func (o *CreateHTTPAfterResponseRuleBackendDefault) WithConfigurationVersion(configurationVersion string) *CreateHTTPAfterResponseRuleBackendDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the create HTTP after response rule backend default response
func (o *CreateHTTPAfterResponseRuleBackendDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the create HTTP after response rule backend default response
func (o *CreateHTTPAfterResponseRuleBackendDefault) WithPayload(payload *models.Error) *CreateHTTPAfterResponseRuleBackendDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create HTTP after response rule backend default response
func (o *CreateHTTPAfterResponseRuleBackendDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateHTTPAfterResponseRuleBackendDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
