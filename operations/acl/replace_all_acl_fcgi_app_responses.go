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

package acl

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v6/models"
)

// ReplaceAllACLFCGIAppOKCode is the HTTP code returned for type ReplaceAllACLFCGIAppOK
const ReplaceAllACLFCGIAppOKCode int = 200

/*
ReplaceAllACLFCGIAppOK All ACL lines replaced

swagger:response replaceAllAclFcgiAppOK
*/
type ReplaceAllACLFCGIAppOK struct {

	/*
	  In: Body
	*/
	Payload models.Acls `json:"body,omitempty"`
}

// NewReplaceAllACLFCGIAppOK creates ReplaceAllACLFCGIAppOK with default headers values
func NewReplaceAllACLFCGIAppOK() *ReplaceAllACLFCGIAppOK {

	return &ReplaceAllACLFCGIAppOK{}
}

// WithPayload adds the payload to the replace all Acl Fcgi app o k response
func (o *ReplaceAllACLFCGIAppOK) WithPayload(payload models.Acls) *ReplaceAllACLFCGIAppOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace all Acl Fcgi app o k response
func (o *ReplaceAllACLFCGIAppOK) SetPayload(payload models.Acls) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceAllACLFCGIAppOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = models.Acls{}
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// ReplaceAllACLFCGIAppAcceptedCode is the HTTP code returned for type ReplaceAllACLFCGIAppAccepted
const ReplaceAllACLFCGIAppAcceptedCode int = 202

/*
ReplaceAllACLFCGIAppAccepted Configuration change accepted and reload requested

swagger:response replaceAllAclFcgiAppAccepted
*/
type ReplaceAllACLFCGIAppAccepted struct {
	/*ID of the requested reload

	 */
	ReloadID string `json:"Reload-ID"`

	/*
	  In: Body
	*/
	Payload models.Acls `json:"body,omitempty"`
}

// NewReplaceAllACLFCGIAppAccepted creates ReplaceAllACLFCGIAppAccepted with default headers values
func NewReplaceAllACLFCGIAppAccepted() *ReplaceAllACLFCGIAppAccepted {

	return &ReplaceAllACLFCGIAppAccepted{}
}

// WithReloadID adds the reloadId to the replace all Acl Fcgi app accepted response
func (o *ReplaceAllACLFCGIAppAccepted) WithReloadID(reloadID string) *ReplaceAllACLFCGIAppAccepted {
	o.ReloadID = reloadID
	return o
}

// SetReloadID sets the reloadId to the replace all Acl Fcgi app accepted response
func (o *ReplaceAllACLFCGIAppAccepted) SetReloadID(reloadID string) {
	o.ReloadID = reloadID
}

// WithPayload adds the payload to the replace all Acl Fcgi app accepted response
func (o *ReplaceAllACLFCGIAppAccepted) WithPayload(payload models.Acls) *ReplaceAllACLFCGIAppAccepted {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace all Acl Fcgi app accepted response
func (o *ReplaceAllACLFCGIAppAccepted) SetPayload(payload models.Acls) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceAllACLFCGIAppAccepted) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	// response header Reload-ID

	reloadID := o.ReloadID
	if reloadID != "" {
		rw.Header().Set("Reload-ID", reloadID)
	}

	rw.WriteHeader(202)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = models.Acls{}
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// ReplaceAllACLFCGIAppBadRequestCode is the HTTP code returned for type ReplaceAllACLFCGIAppBadRequest
const ReplaceAllACLFCGIAppBadRequestCode int = 400

/*
ReplaceAllACLFCGIAppBadRequest Bad request

swagger:response replaceAllAclFcgiAppBadRequest
*/
type ReplaceAllACLFCGIAppBadRequest struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplaceAllACLFCGIAppBadRequest creates ReplaceAllACLFCGIAppBadRequest with default headers values
func NewReplaceAllACLFCGIAppBadRequest() *ReplaceAllACLFCGIAppBadRequest {

	return &ReplaceAllACLFCGIAppBadRequest{}
}

// WithConfigurationVersion adds the configurationVersion to the replace all Acl Fcgi app bad request response
func (o *ReplaceAllACLFCGIAppBadRequest) WithConfigurationVersion(configurationVersion string) *ReplaceAllACLFCGIAppBadRequest {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the replace all Acl Fcgi app bad request response
func (o *ReplaceAllACLFCGIAppBadRequest) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the replace all Acl Fcgi app bad request response
func (o *ReplaceAllACLFCGIAppBadRequest) WithPayload(payload *models.Error) *ReplaceAllACLFCGIAppBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace all Acl Fcgi app bad request response
func (o *ReplaceAllACLFCGIAppBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceAllACLFCGIAppBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

/*
ReplaceAllACLFCGIAppDefault General Error

swagger:response replaceAllAclFcgiAppDefault
*/
type ReplaceAllACLFCGIAppDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplaceAllACLFCGIAppDefault creates ReplaceAllACLFCGIAppDefault with default headers values
func NewReplaceAllACLFCGIAppDefault(code int) *ReplaceAllACLFCGIAppDefault {
	if code <= 0 {
		code = 500
	}

	return &ReplaceAllACLFCGIAppDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the replace all Acl FCGI app default response
func (o *ReplaceAllACLFCGIAppDefault) WithStatusCode(code int) *ReplaceAllACLFCGIAppDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the replace all Acl FCGI app default response
func (o *ReplaceAllACLFCGIAppDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the replace all Acl FCGI app default response
func (o *ReplaceAllACLFCGIAppDefault) WithConfigurationVersion(configurationVersion string) *ReplaceAllACLFCGIAppDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the replace all Acl FCGI app default response
func (o *ReplaceAllACLFCGIAppDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the replace all Acl FCGI app default response
func (o *ReplaceAllACLFCGIAppDefault) WithPayload(payload *models.Error) *ReplaceAllACLFCGIAppDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace all Acl FCGI app default response
func (o *ReplaceAllACLFCGIAppDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceAllACLFCGIAppDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
