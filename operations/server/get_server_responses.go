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

	"github.com/haproxytech/client-native/v3/models"
)

// GetServerOKCode is the HTTP code returned for type GetServerOK
const GetServerOKCode int = 200

/*
GetServerOK Successful operation

swagger:response getServerOK
*/
type GetServerOK struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *GetServerOKBody `json:"body,omitempty"`
}

// NewGetServerOK creates GetServerOK with default headers values
func NewGetServerOK() *GetServerOK {

	return &GetServerOK{}
}

// WithConfigurationVersion adds the configurationVersion to the get server o k response
func (o *GetServerOK) WithConfigurationVersion(configurationVersion string) *GetServerOK {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get server o k response
func (o *GetServerOK) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get server o k response
func (o *GetServerOK) WithPayload(payload *GetServerOKBody) *GetServerOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get server o k response
func (o *GetServerOK) SetPayload(payload *GetServerOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetServerOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	// response header Configuration-Version

	configurationVersion := o.ConfigurationVersion
	if configurationVersion != "" {
		rw.Header().Set("Configuration-Version", configurationVersion)
	}

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetServerNotFoundCode is the HTTP code returned for type GetServerNotFound
const GetServerNotFoundCode int = 404

/*
GetServerNotFound The specified resource was not found

swagger:response getServerNotFound
*/
type GetServerNotFound struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetServerNotFound creates GetServerNotFound with default headers values
func NewGetServerNotFound() *GetServerNotFound {

	return &GetServerNotFound{}
}

// WithConfigurationVersion adds the configurationVersion to the get server not found response
func (o *GetServerNotFound) WithConfigurationVersion(configurationVersion string) *GetServerNotFound {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get server not found response
func (o *GetServerNotFound) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get server not found response
func (o *GetServerNotFound) WithPayload(payload *models.Error) *GetServerNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get server not found response
func (o *GetServerNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetServerNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
GetServerDefault General Error

swagger:response getServerDefault
*/
type GetServerDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetServerDefault creates GetServerDefault with default headers values
func NewGetServerDefault(code int) *GetServerDefault {
	if code <= 0 {
		code = 500
	}

	return &GetServerDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get server default response
func (o *GetServerDefault) WithStatusCode(code int) *GetServerDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get server default response
func (o *GetServerDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the get server default response
func (o *GetServerDefault) WithConfigurationVersion(configurationVersion string) *GetServerDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get server default response
func (o *GetServerDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get server default response
func (o *GetServerDefault) WithPayload(payload *models.Error) *GetServerDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get server default response
func (o *GetServerDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetServerDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
