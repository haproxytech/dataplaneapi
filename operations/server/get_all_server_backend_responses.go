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

// GetAllServerBackendOKCode is the HTTP code returned for type GetAllServerBackendOK
const GetAllServerBackendOKCode int = 200

/*
GetAllServerBackendOK Successful operation

swagger:response getAllServerBackendOK
*/
type GetAllServerBackendOK struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload models.Servers `json:"body,omitempty"`
}

// NewGetAllServerBackendOK creates GetAllServerBackendOK with default headers values
func NewGetAllServerBackendOK() *GetAllServerBackendOK {

	return &GetAllServerBackendOK{}
}

// WithConfigurationVersion adds the configurationVersion to the get all server backend o k response
func (o *GetAllServerBackendOK) WithConfigurationVersion(configurationVersion string) *GetAllServerBackendOK {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get all server backend o k response
func (o *GetAllServerBackendOK) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get all server backend o k response
func (o *GetAllServerBackendOK) WithPayload(payload models.Servers) *GetAllServerBackendOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get all server backend o k response
func (o *GetAllServerBackendOK) SetPayload(payload models.Servers) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetAllServerBackendOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	// response header Configuration-Version

	configurationVersion := o.ConfigurationVersion
	if configurationVersion != "" {
		rw.Header().Set("Configuration-Version", configurationVersion)
	}

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = models.Servers{}
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

/*
GetAllServerBackendDefault General Error

swagger:response getAllServerBackendDefault
*/
type GetAllServerBackendDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetAllServerBackendDefault creates GetAllServerBackendDefault with default headers values
func NewGetAllServerBackendDefault(code int) *GetAllServerBackendDefault {
	if code <= 0 {
		code = 500
	}

	return &GetAllServerBackendDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get all server backend default response
func (o *GetAllServerBackendDefault) WithStatusCode(code int) *GetAllServerBackendDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get all server backend default response
func (o *GetAllServerBackendDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the get all server backend default response
func (o *GetAllServerBackendDefault) WithConfigurationVersion(configurationVersion string) *GetAllServerBackendDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get all server backend default response
func (o *GetAllServerBackendDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get all server backend default response
func (o *GetAllServerBackendDefault) WithPayload(payload *models.Error) *GetAllServerBackendDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get all server backend default response
func (o *GetAllServerBackendDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetAllServerBackendDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
