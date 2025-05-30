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

// CreateHTTPAfterResponseRuleDefaultsCreatedCode is the HTTP code returned for type CreateHTTPAfterResponseRuleDefaultsCreated
const CreateHTTPAfterResponseRuleDefaultsCreatedCode int = 201

/*
CreateHTTPAfterResponseRuleDefaultsCreated HTTP Response Rule created

swagger:response createHttpAfterResponseRuleDefaultsCreated
*/
type CreateHTTPAfterResponseRuleDefaultsCreated struct {

	/*
	  In: Body
	*/
	Payload *models.HTTPAfterResponseRule `json:"body,omitempty"`
}

// NewCreateHTTPAfterResponseRuleDefaultsCreated creates CreateHTTPAfterResponseRuleDefaultsCreated with default headers values
func NewCreateHTTPAfterResponseRuleDefaultsCreated() *CreateHTTPAfterResponseRuleDefaultsCreated {

	return &CreateHTTPAfterResponseRuleDefaultsCreated{}
}

// WithPayload adds the payload to the create Http after response rule defaults created response
func (o *CreateHTTPAfterResponseRuleDefaultsCreated) WithPayload(payload *models.HTTPAfterResponseRule) *CreateHTTPAfterResponseRuleDefaultsCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create Http after response rule defaults created response
func (o *CreateHTTPAfterResponseRuleDefaultsCreated) SetPayload(payload *models.HTTPAfterResponseRule) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateHTTPAfterResponseRuleDefaultsCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateHTTPAfterResponseRuleDefaultsAcceptedCode is the HTTP code returned for type CreateHTTPAfterResponseRuleDefaultsAccepted
const CreateHTTPAfterResponseRuleDefaultsAcceptedCode int = 202

/*
CreateHTTPAfterResponseRuleDefaultsAccepted Configuration change accepted and reload requested

swagger:response createHttpAfterResponseRuleDefaultsAccepted
*/
type CreateHTTPAfterResponseRuleDefaultsAccepted struct {
	/*ID of the requested reload

	 */
	ReloadID string `json:"Reload-ID"`

	/*
	  In: Body
	*/
	Payload *models.HTTPAfterResponseRule `json:"body,omitempty"`
}

// NewCreateHTTPAfterResponseRuleDefaultsAccepted creates CreateHTTPAfterResponseRuleDefaultsAccepted with default headers values
func NewCreateHTTPAfterResponseRuleDefaultsAccepted() *CreateHTTPAfterResponseRuleDefaultsAccepted {

	return &CreateHTTPAfterResponseRuleDefaultsAccepted{}
}

// WithReloadID adds the reloadId to the create Http after response rule defaults accepted response
func (o *CreateHTTPAfterResponseRuleDefaultsAccepted) WithReloadID(reloadID string) *CreateHTTPAfterResponseRuleDefaultsAccepted {
	o.ReloadID = reloadID
	return o
}

// SetReloadID sets the reloadId to the create Http after response rule defaults accepted response
func (o *CreateHTTPAfterResponseRuleDefaultsAccepted) SetReloadID(reloadID string) {
	o.ReloadID = reloadID
}

// WithPayload adds the payload to the create Http after response rule defaults accepted response
func (o *CreateHTTPAfterResponseRuleDefaultsAccepted) WithPayload(payload *models.HTTPAfterResponseRule) *CreateHTTPAfterResponseRuleDefaultsAccepted {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create Http after response rule defaults accepted response
func (o *CreateHTTPAfterResponseRuleDefaultsAccepted) SetPayload(payload *models.HTTPAfterResponseRule) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateHTTPAfterResponseRuleDefaultsAccepted) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// CreateHTTPAfterResponseRuleDefaultsBadRequestCode is the HTTP code returned for type CreateHTTPAfterResponseRuleDefaultsBadRequest
const CreateHTTPAfterResponseRuleDefaultsBadRequestCode int = 400

/*
CreateHTTPAfterResponseRuleDefaultsBadRequest Bad request

swagger:response createHttpAfterResponseRuleDefaultsBadRequest
*/
type CreateHTTPAfterResponseRuleDefaultsBadRequest struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateHTTPAfterResponseRuleDefaultsBadRequest creates CreateHTTPAfterResponseRuleDefaultsBadRequest with default headers values
func NewCreateHTTPAfterResponseRuleDefaultsBadRequest() *CreateHTTPAfterResponseRuleDefaultsBadRequest {

	return &CreateHTTPAfterResponseRuleDefaultsBadRequest{}
}

// WithConfigurationVersion adds the configurationVersion to the create Http after response rule defaults bad request response
func (o *CreateHTTPAfterResponseRuleDefaultsBadRequest) WithConfigurationVersion(configurationVersion string) *CreateHTTPAfterResponseRuleDefaultsBadRequest {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the create Http after response rule defaults bad request response
func (o *CreateHTTPAfterResponseRuleDefaultsBadRequest) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the create Http after response rule defaults bad request response
func (o *CreateHTTPAfterResponseRuleDefaultsBadRequest) WithPayload(payload *models.Error) *CreateHTTPAfterResponseRuleDefaultsBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create Http after response rule defaults bad request response
func (o *CreateHTTPAfterResponseRuleDefaultsBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateHTTPAfterResponseRuleDefaultsBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// CreateHTTPAfterResponseRuleDefaultsConflictCode is the HTTP code returned for type CreateHTTPAfterResponseRuleDefaultsConflict
const CreateHTTPAfterResponseRuleDefaultsConflictCode int = 409

/*
CreateHTTPAfterResponseRuleDefaultsConflict The specified resource already exists

swagger:response createHttpAfterResponseRuleDefaultsConflict
*/
type CreateHTTPAfterResponseRuleDefaultsConflict struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateHTTPAfterResponseRuleDefaultsConflict creates CreateHTTPAfterResponseRuleDefaultsConflict with default headers values
func NewCreateHTTPAfterResponseRuleDefaultsConflict() *CreateHTTPAfterResponseRuleDefaultsConflict {

	return &CreateHTTPAfterResponseRuleDefaultsConflict{}
}

// WithConfigurationVersion adds the configurationVersion to the create Http after response rule defaults conflict response
func (o *CreateHTTPAfterResponseRuleDefaultsConflict) WithConfigurationVersion(configurationVersion string) *CreateHTTPAfterResponseRuleDefaultsConflict {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the create Http after response rule defaults conflict response
func (o *CreateHTTPAfterResponseRuleDefaultsConflict) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the create Http after response rule defaults conflict response
func (o *CreateHTTPAfterResponseRuleDefaultsConflict) WithPayload(payload *models.Error) *CreateHTTPAfterResponseRuleDefaultsConflict {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create Http after response rule defaults conflict response
func (o *CreateHTTPAfterResponseRuleDefaultsConflict) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateHTTPAfterResponseRuleDefaultsConflict) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
CreateHTTPAfterResponseRuleDefaultsDefault General Error

swagger:response createHttpAfterResponseRuleDefaultsDefault
*/
type CreateHTTPAfterResponseRuleDefaultsDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateHTTPAfterResponseRuleDefaultsDefault creates CreateHTTPAfterResponseRuleDefaultsDefault with default headers values
func NewCreateHTTPAfterResponseRuleDefaultsDefault(code int) *CreateHTTPAfterResponseRuleDefaultsDefault {
	if code <= 0 {
		code = 500
	}

	return &CreateHTTPAfterResponseRuleDefaultsDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the create HTTP after response rule defaults default response
func (o *CreateHTTPAfterResponseRuleDefaultsDefault) WithStatusCode(code int) *CreateHTTPAfterResponseRuleDefaultsDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the create HTTP after response rule defaults default response
func (o *CreateHTTPAfterResponseRuleDefaultsDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the create HTTP after response rule defaults default response
func (o *CreateHTTPAfterResponseRuleDefaultsDefault) WithConfigurationVersion(configurationVersion string) *CreateHTTPAfterResponseRuleDefaultsDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the create HTTP after response rule defaults default response
func (o *CreateHTTPAfterResponseRuleDefaultsDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the create HTTP after response rule defaults default response
func (o *CreateHTTPAfterResponseRuleDefaultsDefault) WithPayload(payload *models.Error) *CreateHTTPAfterResponseRuleDefaultsDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create HTTP after response rule defaults default response
func (o *CreateHTTPAfterResponseRuleDefaultsDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateHTTPAfterResponseRuleDefaultsDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
