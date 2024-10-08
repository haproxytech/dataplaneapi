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

package server

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v6/models"
)

// GetAllRuntimeServerOKCode is the HTTP code returned for type GetAllRuntimeServerOK
const GetAllRuntimeServerOKCode int = 200

/*
GetAllRuntimeServerOK Successful operation

swagger:response getAllRuntimeServerOK
*/
type GetAllRuntimeServerOK struct {

	/*
	  In: Body
	*/
	Payload models.RuntimeServers `json:"body,omitempty"`
}

// NewGetAllRuntimeServerOK creates GetAllRuntimeServerOK with default headers values
func NewGetAllRuntimeServerOK() *GetAllRuntimeServerOK {

	return &GetAllRuntimeServerOK{}
}

// WithPayload adds the payload to the get all runtime server o k response
func (o *GetAllRuntimeServerOK) WithPayload(payload models.RuntimeServers) *GetAllRuntimeServerOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get all runtime server o k response
func (o *GetAllRuntimeServerOK) SetPayload(payload models.RuntimeServers) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetAllRuntimeServerOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = models.RuntimeServers{}
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

/*
GetAllRuntimeServerDefault General Error

swagger:response getAllRuntimeServerDefault
*/
type GetAllRuntimeServerDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetAllRuntimeServerDefault creates GetAllRuntimeServerDefault with default headers values
func NewGetAllRuntimeServerDefault(code int) *GetAllRuntimeServerDefault {
	if code <= 0 {
		code = 500
	}

	return &GetAllRuntimeServerDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get all runtime server default response
func (o *GetAllRuntimeServerDefault) WithStatusCode(code int) *GetAllRuntimeServerDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get all runtime server default response
func (o *GetAllRuntimeServerDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the get all runtime server default response
func (o *GetAllRuntimeServerDefault) WithConfigurationVersion(configurationVersion string) *GetAllRuntimeServerDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get all runtime server default response
func (o *GetAllRuntimeServerDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get all runtime server default response
func (o *GetAllRuntimeServerDefault) WithPayload(payload *models.Error) *GetAllRuntimeServerDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get all runtime server default response
func (o *GetAllRuntimeServerDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetAllRuntimeServerDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
