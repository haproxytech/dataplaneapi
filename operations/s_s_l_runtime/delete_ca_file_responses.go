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

// DeleteCaFileNoContentCode is the HTTP code returned for type DeleteCaFileNoContent
const DeleteCaFileNoContentCode int = 204

/*
DeleteCaFileNoContent SSL CA deleted

swagger:response deleteCaFileNoContent
*/
type DeleteCaFileNoContent struct {
}

// NewDeleteCaFileNoContent creates DeleteCaFileNoContent with default headers values
func NewDeleteCaFileNoContent() *DeleteCaFileNoContent {

	return &DeleteCaFileNoContent{}
}

// WriteResponse to the client
func (o *DeleteCaFileNoContent) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(204)
}

// DeleteCaFileBadRequestCode is the HTTP code returned for type DeleteCaFileBadRequest
const DeleteCaFileBadRequestCode int = 400

/*
DeleteCaFileBadRequest Bad request

swagger:response deleteCaFileBadRequest
*/
type DeleteCaFileBadRequest struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDeleteCaFileBadRequest creates DeleteCaFileBadRequest with default headers values
func NewDeleteCaFileBadRequest() *DeleteCaFileBadRequest {

	return &DeleteCaFileBadRequest{}
}

// WithConfigurationVersion adds the configurationVersion to the delete ca file bad request response
func (o *DeleteCaFileBadRequest) WithConfigurationVersion(configurationVersion string) *DeleteCaFileBadRequest {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the delete ca file bad request response
func (o *DeleteCaFileBadRequest) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the delete ca file bad request response
func (o *DeleteCaFileBadRequest) WithPayload(payload *models.Error) *DeleteCaFileBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete ca file bad request response
func (o *DeleteCaFileBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteCaFileBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
DeleteCaFileDefault General Error

swagger:response deleteCaFileDefault
*/
type DeleteCaFileDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDeleteCaFileDefault creates DeleteCaFileDefault with default headers values
func NewDeleteCaFileDefault(code int) *DeleteCaFileDefault {
	if code <= 0 {
		code = 500
	}

	return &DeleteCaFileDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the delete ca file default response
func (o *DeleteCaFileDefault) WithStatusCode(code int) *DeleteCaFileDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the delete ca file default response
func (o *DeleteCaFileDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the delete ca file default response
func (o *DeleteCaFileDefault) WithConfigurationVersion(configurationVersion string) *DeleteCaFileDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the delete ca file default response
func (o *DeleteCaFileDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the delete ca file default response
func (o *DeleteCaFileDefault) WithPayload(payload *models.Error) *DeleteCaFileDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete ca file default response
func (o *DeleteCaFileDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteCaFileDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
