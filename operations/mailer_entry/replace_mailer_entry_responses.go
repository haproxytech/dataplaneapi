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

	"github.com/haproxytech/client-native/v5/models"
)

// ReplaceMailerEntryOKCode is the HTTP code returned for type ReplaceMailerEntryOK
const ReplaceMailerEntryOKCode int = 200

/*
ReplaceMailerEntryOK MailerEntry replaced

swagger:response replaceMailerEntryOK
*/
type ReplaceMailerEntryOK struct {

	/*
	  In: Body
	*/
	Payload *models.MailerEntry `json:"body,omitempty"`
}

// NewReplaceMailerEntryOK creates ReplaceMailerEntryOK with default headers values
func NewReplaceMailerEntryOK() *ReplaceMailerEntryOK {

	return &ReplaceMailerEntryOK{}
}

// WithPayload adds the payload to the replace mailer entry o k response
func (o *ReplaceMailerEntryOK) WithPayload(payload *models.MailerEntry) *ReplaceMailerEntryOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace mailer entry o k response
func (o *ReplaceMailerEntryOK) SetPayload(payload *models.MailerEntry) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceMailerEntryOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// ReplaceMailerEntryAcceptedCode is the HTTP code returned for type ReplaceMailerEntryAccepted
const ReplaceMailerEntryAcceptedCode int = 202

/*
ReplaceMailerEntryAccepted Configuration change accepted and reload requested

swagger:response replaceMailerEntryAccepted
*/
type ReplaceMailerEntryAccepted struct {
	/*ID of the requested reload

	 */
	ReloadID string `json:"Reload-ID"`

	/*
	  In: Body
	*/
	Payload *models.MailerEntry `json:"body,omitempty"`
}

// NewReplaceMailerEntryAccepted creates ReplaceMailerEntryAccepted with default headers values
func NewReplaceMailerEntryAccepted() *ReplaceMailerEntryAccepted {

	return &ReplaceMailerEntryAccepted{}
}

// WithReloadID adds the reloadId to the replace mailer entry accepted response
func (o *ReplaceMailerEntryAccepted) WithReloadID(reloadID string) *ReplaceMailerEntryAccepted {
	o.ReloadID = reloadID
	return o
}

// SetReloadID sets the reloadId to the replace mailer entry accepted response
func (o *ReplaceMailerEntryAccepted) SetReloadID(reloadID string) {
	o.ReloadID = reloadID
}

// WithPayload adds the payload to the replace mailer entry accepted response
func (o *ReplaceMailerEntryAccepted) WithPayload(payload *models.MailerEntry) *ReplaceMailerEntryAccepted {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace mailer entry accepted response
func (o *ReplaceMailerEntryAccepted) SetPayload(payload *models.MailerEntry) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceMailerEntryAccepted) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// ReplaceMailerEntryBadRequestCode is the HTTP code returned for type ReplaceMailerEntryBadRequest
const ReplaceMailerEntryBadRequestCode int = 400

/*
ReplaceMailerEntryBadRequest Bad request

swagger:response replaceMailerEntryBadRequest
*/
type ReplaceMailerEntryBadRequest struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplaceMailerEntryBadRequest creates ReplaceMailerEntryBadRequest with default headers values
func NewReplaceMailerEntryBadRequest() *ReplaceMailerEntryBadRequest {

	return &ReplaceMailerEntryBadRequest{}
}

// WithConfigurationVersion adds the configurationVersion to the replace mailer entry bad request response
func (o *ReplaceMailerEntryBadRequest) WithConfigurationVersion(configurationVersion string) *ReplaceMailerEntryBadRequest {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the replace mailer entry bad request response
func (o *ReplaceMailerEntryBadRequest) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the replace mailer entry bad request response
func (o *ReplaceMailerEntryBadRequest) WithPayload(payload *models.Error) *ReplaceMailerEntryBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace mailer entry bad request response
func (o *ReplaceMailerEntryBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceMailerEntryBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// ReplaceMailerEntryNotFoundCode is the HTTP code returned for type ReplaceMailerEntryNotFound
const ReplaceMailerEntryNotFoundCode int = 404

/*
ReplaceMailerEntryNotFound The specified resource was not found

swagger:response replaceMailerEntryNotFound
*/
type ReplaceMailerEntryNotFound struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplaceMailerEntryNotFound creates ReplaceMailerEntryNotFound with default headers values
func NewReplaceMailerEntryNotFound() *ReplaceMailerEntryNotFound {

	return &ReplaceMailerEntryNotFound{}
}

// WithConfigurationVersion adds the configurationVersion to the replace mailer entry not found response
func (o *ReplaceMailerEntryNotFound) WithConfigurationVersion(configurationVersion string) *ReplaceMailerEntryNotFound {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the replace mailer entry not found response
func (o *ReplaceMailerEntryNotFound) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the replace mailer entry not found response
func (o *ReplaceMailerEntryNotFound) WithPayload(payload *models.Error) *ReplaceMailerEntryNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace mailer entry not found response
func (o *ReplaceMailerEntryNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceMailerEntryNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
ReplaceMailerEntryDefault General Error

swagger:response replaceMailerEntryDefault
*/
type ReplaceMailerEntryDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplaceMailerEntryDefault creates ReplaceMailerEntryDefault with default headers values
func NewReplaceMailerEntryDefault(code int) *ReplaceMailerEntryDefault {
	if code <= 0 {
		code = 500
	}

	return &ReplaceMailerEntryDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the replace mailer entry default response
func (o *ReplaceMailerEntryDefault) WithStatusCode(code int) *ReplaceMailerEntryDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the replace mailer entry default response
func (o *ReplaceMailerEntryDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the replace mailer entry default response
func (o *ReplaceMailerEntryDefault) WithConfigurationVersion(configurationVersion string) *ReplaceMailerEntryDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the replace mailer entry default response
func (o *ReplaceMailerEntryDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the replace mailer entry default response
func (o *ReplaceMailerEntryDefault) WithPayload(payload *models.Error) *ReplaceMailerEntryDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace mailer entry default response
func (o *ReplaceMailerEntryDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceMailerEntryDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
