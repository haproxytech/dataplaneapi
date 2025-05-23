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

package s_s_l_runtime

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v6/models"
)

// DeleteCertNoContentCode is the HTTP code returned for type DeleteCertNoContent
const DeleteCertNoContentCode int = 204

/*
DeleteCertNoContent File deleted

swagger:response deleteCertNoContent
*/
type DeleteCertNoContent struct {
}

// NewDeleteCertNoContent creates DeleteCertNoContent with default headers values
func NewDeleteCertNoContent() *DeleteCertNoContent {

	return &DeleteCertNoContent{}
}

// WriteResponse to the client
func (o *DeleteCertNoContent) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(204)
}

// DeleteCertBadRequestCode is the HTTP code returned for type DeleteCertBadRequest
const DeleteCertBadRequestCode int = 400

/*
DeleteCertBadRequest Bad request

swagger:response deleteCertBadRequest
*/
type DeleteCertBadRequest struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDeleteCertBadRequest creates DeleteCertBadRequest with default headers values
func NewDeleteCertBadRequest() *DeleteCertBadRequest {

	return &DeleteCertBadRequest{}
}

// WithConfigurationVersion adds the configurationVersion to the delete cert bad request response
func (o *DeleteCertBadRequest) WithConfigurationVersion(configurationVersion string) *DeleteCertBadRequest {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the delete cert bad request response
func (o *DeleteCertBadRequest) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the delete cert bad request response
func (o *DeleteCertBadRequest) WithPayload(payload *models.Error) *DeleteCertBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete cert bad request response
func (o *DeleteCertBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteCertBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
DeleteCertDefault General Error

swagger:response deleteCertDefault
*/
type DeleteCertDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDeleteCertDefault creates DeleteCertDefault with default headers values
func NewDeleteCertDefault(code int) *DeleteCertDefault {
	if code <= 0 {
		code = 500
	}

	return &DeleteCertDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the delete cert default response
func (o *DeleteCertDefault) WithStatusCode(code int) *DeleteCertDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the delete cert default response
func (o *DeleteCertDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the delete cert default response
func (o *DeleteCertDefault) WithConfigurationVersion(configurationVersion string) *DeleteCertDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the delete cert default response
func (o *DeleteCertDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the delete cert default response
func (o *DeleteCertDefault) WithPayload(payload *models.Error) *DeleteCertDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete cert default response
func (o *DeleteCertDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteCertDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
