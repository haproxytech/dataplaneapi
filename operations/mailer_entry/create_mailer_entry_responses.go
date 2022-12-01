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

package mailer_entry

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v4/models"
)

// CreateMailerEntryCreatedCode is the HTTP code returned for type CreateMailerEntryCreated
const CreateMailerEntryCreatedCode int = 201

/*
CreateMailerEntryCreated MailerEntry created

swagger:response createMailerEntryCreated
*/
type CreateMailerEntryCreated struct {

	/*
	  In: Body
	*/
	Payload *models.MailerEntry `json:"body,omitempty"`
}

// NewCreateMailerEntryCreated creates CreateMailerEntryCreated with default headers values
func NewCreateMailerEntryCreated() *CreateMailerEntryCreated {

	return &CreateMailerEntryCreated{}
}

// WithPayload adds the payload to the create mailer entry created response
func (o *CreateMailerEntryCreated) WithPayload(payload *models.MailerEntry) *CreateMailerEntryCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create mailer entry created response
func (o *CreateMailerEntryCreated) SetPayload(payload *models.MailerEntry) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateMailerEntryCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateMailerEntryAcceptedCode is the HTTP code returned for type CreateMailerEntryAccepted
const CreateMailerEntryAcceptedCode int = 202

/*
CreateMailerEntryAccepted Configuration change accepted and reload requested

swagger:response createMailerEntryAccepted
*/
type CreateMailerEntryAccepted struct {
	/*ID of the requested reload

	 */
	ReloadID string `json:"Reload-ID"`

	/*
	  In: Body
	*/
	Payload *models.MailerEntry `json:"body,omitempty"`
}

// NewCreateMailerEntryAccepted creates CreateMailerEntryAccepted with default headers values
func NewCreateMailerEntryAccepted() *CreateMailerEntryAccepted {

	return &CreateMailerEntryAccepted{}
}

// WithReloadID adds the reloadId to the create mailer entry accepted response
func (o *CreateMailerEntryAccepted) WithReloadID(reloadID string) *CreateMailerEntryAccepted {
	o.ReloadID = reloadID
	return o
}

// SetReloadID sets the reloadId to the create mailer entry accepted response
func (o *CreateMailerEntryAccepted) SetReloadID(reloadID string) {
	o.ReloadID = reloadID
}

// WithPayload adds the payload to the create mailer entry accepted response
func (o *CreateMailerEntryAccepted) WithPayload(payload *models.MailerEntry) *CreateMailerEntryAccepted {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create mailer entry accepted response
func (o *CreateMailerEntryAccepted) SetPayload(payload *models.MailerEntry) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateMailerEntryAccepted) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// CreateMailerEntryBadRequestCode is the HTTP code returned for type CreateMailerEntryBadRequest
const CreateMailerEntryBadRequestCode int = 400

/*
CreateMailerEntryBadRequest Bad request

swagger:response createMailerEntryBadRequest
*/
type CreateMailerEntryBadRequest struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateMailerEntryBadRequest creates CreateMailerEntryBadRequest with default headers values
func NewCreateMailerEntryBadRequest() *CreateMailerEntryBadRequest {

	return &CreateMailerEntryBadRequest{}
}

// WithConfigurationVersion adds the configurationVersion to the create mailer entry bad request response
func (o *CreateMailerEntryBadRequest) WithConfigurationVersion(configurationVersion string) *CreateMailerEntryBadRequest {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the create mailer entry bad request response
func (o *CreateMailerEntryBadRequest) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the create mailer entry bad request response
func (o *CreateMailerEntryBadRequest) WithPayload(payload *models.Error) *CreateMailerEntryBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create mailer entry bad request response
func (o *CreateMailerEntryBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateMailerEntryBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// CreateMailerEntryConflictCode is the HTTP code returned for type CreateMailerEntryConflict
const CreateMailerEntryConflictCode int = 409

/*
CreateMailerEntryConflict The specified resource already exists

swagger:response createMailerEntryConflict
*/
type CreateMailerEntryConflict struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateMailerEntryConflict creates CreateMailerEntryConflict with default headers values
func NewCreateMailerEntryConflict() *CreateMailerEntryConflict {

	return &CreateMailerEntryConflict{}
}

// WithConfigurationVersion adds the configurationVersion to the create mailer entry conflict response
func (o *CreateMailerEntryConflict) WithConfigurationVersion(configurationVersion string) *CreateMailerEntryConflict {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the create mailer entry conflict response
func (o *CreateMailerEntryConflict) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the create mailer entry conflict response
func (o *CreateMailerEntryConflict) WithPayload(payload *models.Error) *CreateMailerEntryConflict {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create mailer entry conflict response
func (o *CreateMailerEntryConflict) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateMailerEntryConflict) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
CreateMailerEntryDefault General Error

swagger:response createMailerEntryDefault
*/
type CreateMailerEntryDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateMailerEntryDefault creates CreateMailerEntryDefault with default headers values
func NewCreateMailerEntryDefault(code int) *CreateMailerEntryDefault {
	if code <= 0 {
		code = 500
	}

	return &CreateMailerEntryDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the create mailer entry default response
func (o *CreateMailerEntryDefault) WithStatusCode(code int) *CreateMailerEntryDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the create mailer entry default response
func (o *CreateMailerEntryDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the create mailer entry default response
func (o *CreateMailerEntryDefault) WithConfigurationVersion(configurationVersion string) *CreateMailerEntryDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the create mailer entry default response
func (o *CreateMailerEntryDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the create mailer entry default response
func (o *CreateMailerEntryDefault) WithPayload(payload *models.Error) *CreateMailerEntryDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create mailer entry default response
func (o *CreateMailerEntryDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateMailerEntryDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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