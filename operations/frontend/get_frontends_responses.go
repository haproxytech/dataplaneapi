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

package frontend

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v3/models"
)

// GetFrontendsOKCode is the HTTP code returned for type GetFrontendsOK
const GetFrontendsOKCode int = 200

/*
GetFrontendsOK Successful operation

swagger:response getFrontendsOK
*/
type GetFrontendsOK struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *GetFrontendsOKBody `json:"body,omitempty"`
}

// NewGetFrontendsOK creates GetFrontendsOK with default headers values
func NewGetFrontendsOK() *GetFrontendsOK {

	return &GetFrontendsOK{}
}

// WithConfigurationVersion adds the configurationVersion to the get frontends o k response
func (o *GetFrontendsOK) WithConfigurationVersion(configurationVersion string) *GetFrontendsOK {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get frontends o k response
func (o *GetFrontendsOK) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get frontends o k response
func (o *GetFrontendsOK) WithPayload(payload *GetFrontendsOKBody) *GetFrontendsOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get frontends o k response
func (o *GetFrontendsOK) SetPayload(payload *GetFrontendsOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetFrontendsOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

/*
GetFrontendsDefault General Error

swagger:response getFrontendsDefault
*/
type GetFrontendsDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetFrontendsDefault creates GetFrontendsDefault with default headers values
func NewGetFrontendsDefault(code int) *GetFrontendsDefault {
	if code <= 0 {
		code = 500
	}

	return &GetFrontendsDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get frontends default response
func (o *GetFrontendsDefault) WithStatusCode(code int) *GetFrontendsDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get frontends default response
func (o *GetFrontendsDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the get frontends default response
func (o *GetFrontendsDefault) WithConfigurationVersion(configurationVersion string) *GetFrontendsDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get frontends default response
func (o *GetFrontendsDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get frontends default response
func (o *GetFrontendsDefault) WithPayload(payload *models.Error) *GetFrontendsDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get frontends default response
func (o *GetFrontendsDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetFrontendsDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
