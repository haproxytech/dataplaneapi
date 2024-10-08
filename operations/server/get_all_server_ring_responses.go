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

// GetAllServerRingOKCode is the HTTP code returned for type GetAllServerRingOK
const GetAllServerRingOKCode int = 200

/*
GetAllServerRingOK Successful operation

swagger:response getAllServerRingOK
*/
type GetAllServerRingOK struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload models.Servers `json:"body,omitempty"`
}

// NewGetAllServerRingOK creates GetAllServerRingOK with default headers values
func NewGetAllServerRingOK() *GetAllServerRingOK {

	return &GetAllServerRingOK{}
}

// WithConfigurationVersion adds the configurationVersion to the get all server ring o k response
func (o *GetAllServerRingOK) WithConfigurationVersion(configurationVersion string) *GetAllServerRingOK {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get all server ring o k response
func (o *GetAllServerRingOK) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get all server ring o k response
func (o *GetAllServerRingOK) WithPayload(payload models.Servers) *GetAllServerRingOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get all server ring o k response
func (o *GetAllServerRingOK) SetPayload(payload models.Servers) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetAllServerRingOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
GetAllServerRingDefault General Error

swagger:response getAllServerRingDefault
*/
type GetAllServerRingDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetAllServerRingDefault creates GetAllServerRingDefault with default headers values
func NewGetAllServerRingDefault(code int) *GetAllServerRingDefault {
	if code <= 0 {
		code = 500
	}

	return &GetAllServerRingDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get all server ring default response
func (o *GetAllServerRingDefault) WithStatusCode(code int) *GetAllServerRingDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get all server ring default response
func (o *GetAllServerRingDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the get all server ring default response
func (o *GetAllServerRingDefault) WithConfigurationVersion(configurationVersion string) *GetAllServerRingDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get all server ring default response
func (o *GetAllServerRingDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get all server ring default response
func (o *GetAllServerRingDefault) WithPayload(payload *models.Error) *GetAllServerRingDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get all server ring default response
func (o *GetAllServerRingDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetAllServerRingDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
