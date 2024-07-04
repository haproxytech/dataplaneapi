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

package acl_runtime

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v6/models"
)

// DeleteServicesHaproxyRuntimeAclsParentNameEntriesIDNoContentCode is the HTTP code returned for type DeleteServicesHaproxyRuntimeAclsParentNameEntriesIDNoContent
const DeleteServicesHaproxyRuntimeAclsParentNameEntriesIDNoContentCode int = 204

/*
DeleteServicesHaproxyRuntimeAclsParentNameEntriesIDNoContent Successful operation

swagger:response deleteServicesHaproxyRuntimeAclsParentNameEntriesIdNoContent
*/
type DeleteServicesHaproxyRuntimeAclsParentNameEntriesIDNoContent struct {
}

// NewDeleteServicesHaproxyRuntimeAclsParentNameEntriesIDNoContent creates DeleteServicesHaproxyRuntimeAclsParentNameEntriesIDNoContent with default headers values
func NewDeleteServicesHaproxyRuntimeAclsParentNameEntriesIDNoContent() *DeleteServicesHaproxyRuntimeAclsParentNameEntriesIDNoContent {

	return &DeleteServicesHaproxyRuntimeAclsParentNameEntriesIDNoContent{}
}

// WriteResponse to the client
func (o *DeleteServicesHaproxyRuntimeAclsParentNameEntriesIDNoContent) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(204)
}

// DeleteServicesHaproxyRuntimeAclsParentNameEntriesIDBadRequestCode is the HTTP code returned for type DeleteServicesHaproxyRuntimeAclsParentNameEntriesIDBadRequest
const DeleteServicesHaproxyRuntimeAclsParentNameEntriesIDBadRequestCode int = 400

/*
DeleteServicesHaproxyRuntimeAclsParentNameEntriesIDBadRequest Bad request

swagger:response deleteServicesHaproxyRuntimeAclsParentNameEntriesIdBadRequest
*/
type DeleteServicesHaproxyRuntimeAclsParentNameEntriesIDBadRequest struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDeleteServicesHaproxyRuntimeAclsParentNameEntriesIDBadRequest creates DeleteServicesHaproxyRuntimeAclsParentNameEntriesIDBadRequest with default headers values
func NewDeleteServicesHaproxyRuntimeAclsParentNameEntriesIDBadRequest() *DeleteServicesHaproxyRuntimeAclsParentNameEntriesIDBadRequest {

	return &DeleteServicesHaproxyRuntimeAclsParentNameEntriesIDBadRequest{}
}

// WithConfigurationVersion adds the configurationVersion to the delete services haproxy runtime acls parent name entries Id bad request response
func (o *DeleteServicesHaproxyRuntimeAclsParentNameEntriesIDBadRequest) WithConfigurationVersion(configurationVersion string) *DeleteServicesHaproxyRuntimeAclsParentNameEntriesIDBadRequest {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the delete services haproxy runtime acls parent name entries Id bad request response
func (o *DeleteServicesHaproxyRuntimeAclsParentNameEntriesIDBadRequest) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the delete services haproxy runtime acls parent name entries Id bad request response
func (o *DeleteServicesHaproxyRuntimeAclsParentNameEntriesIDBadRequest) WithPayload(payload *models.Error) *DeleteServicesHaproxyRuntimeAclsParentNameEntriesIDBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete services haproxy runtime acls parent name entries Id bad request response
func (o *DeleteServicesHaproxyRuntimeAclsParentNameEntriesIDBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteServicesHaproxyRuntimeAclsParentNameEntriesIDBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// DeleteServicesHaproxyRuntimeAclsParentNameEntriesIDNotFoundCode is the HTTP code returned for type DeleteServicesHaproxyRuntimeAclsParentNameEntriesIDNotFound
const DeleteServicesHaproxyRuntimeAclsParentNameEntriesIDNotFoundCode int = 404

/*
DeleteServicesHaproxyRuntimeAclsParentNameEntriesIDNotFound The specified resource was not found

swagger:response deleteServicesHaproxyRuntimeAclsParentNameEntriesIdNotFound
*/
type DeleteServicesHaproxyRuntimeAclsParentNameEntriesIDNotFound struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDeleteServicesHaproxyRuntimeAclsParentNameEntriesIDNotFound creates DeleteServicesHaproxyRuntimeAclsParentNameEntriesIDNotFound with default headers values
func NewDeleteServicesHaproxyRuntimeAclsParentNameEntriesIDNotFound() *DeleteServicesHaproxyRuntimeAclsParentNameEntriesIDNotFound {

	return &DeleteServicesHaproxyRuntimeAclsParentNameEntriesIDNotFound{}
}

// WithConfigurationVersion adds the configurationVersion to the delete services haproxy runtime acls parent name entries Id not found response
func (o *DeleteServicesHaproxyRuntimeAclsParentNameEntriesIDNotFound) WithConfigurationVersion(configurationVersion string) *DeleteServicesHaproxyRuntimeAclsParentNameEntriesIDNotFound {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the delete services haproxy runtime acls parent name entries Id not found response
func (o *DeleteServicesHaproxyRuntimeAclsParentNameEntriesIDNotFound) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the delete services haproxy runtime acls parent name entries Id not found response
func (o *DeleteServicesHaproxyRuntimeAclsParentNameEntriesIDNotFound) WithPayload(payload *models.Error) *DeleteServicesHaproxyRuntimeAclsParentNameEntriesIDNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete services haproxy runtime acls parent name entries Id not found response
func (o *DeleteServicesHaproxyRuntimeAclsParentNameEntriesIDNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteServicesHaproxyRuntimeAclsParentNameEntriesIDNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
DeleteServicesHaproxyRuntimeAclsParentNameEntriesIDDefault General Error

swagger:response deleteServicesHaproxyRuntimeAclsParentNameEntriesIdDefault
*/
type DeleteServicesHaproxyRuntimeAclsParentNameEntriesIDDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDeleteServicesHaproxyRuntimeAclsParentNameEntriesIDDefault creates DeleteServicesHaproxyRuntimeAclsParentNameEntriesIDDefault with default headers values
func NewDeleteServicesHaproxyRuntimeAclsParentNameEntriesIDDefault(code int) *DeleteServicesHaproxyRuntimeAclsParentNameEntriesIDDefault {
	if code <= 0 {
		code = 500
	}

	return &DeleteServicesHaproxyRuntimeAclsParentNameEntriesIDDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the delete services haproxy runtime acls parent name entries ID default response
func (o *DeleteServicesHaproxyRuntimeAclsParentNameEntriesIDDefault) WithStatusCode(code int) *DeleteServicesHaproxyRuntimeAclsParentNameEntriesIDDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the delete services haproxy runtime acls parent name entries ID default response
func (o *DeleteServicesHaproxyRuntimeAclsParentNameEntriesIDDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the delete services haproxy runtime acls parent name entries ID default response
func (o *DeleteServicesHaproxyRuntimeAclsParentNameEntriesIDDefault) WithConfigurationVersion(configurationVersion string) *DeleteServicesHaproxyRuntimeAclsParentNameEntriesIDDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the delete services haproxy runtime acls parent name entries ID default response
func (o *DeleteServicesHaproxyRuntimeAclsParentNameEntriesIDDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the delete services haproxy runtime acls parent name entries ID default response
func (o *DeleteServicesHaproxyRuntimeAclsParentNameEntriesIDDefault) WithPayload(payload *models.Error) *DeleteServicesHaproxyRuntimeAclsParentNameEntriesIDDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete services haproxy runtime acls parent name entries ID default response
func (o *DeleteServicesHaproxyRuntimeAclsParentNameEntriesIDDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteServicesHaproxyRuntimeAclsParentNameEntriesIDDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
