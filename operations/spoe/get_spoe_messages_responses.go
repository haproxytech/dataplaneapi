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

package spoe

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v3/models"
)

// GetSpoeMessagesOKCode is the HTTP code returned for type GetSpoeMessagesOK
const GetSpoeMessagesOKCode int = 200

/*
GetSpoeMessagesOK Successful operation

swagger:response getSpoeMessagesOK
*/
type GetSpoeMessagesOK struct {
	/*Spoe configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *GetSpoeMessagesOKBody `json:"body,omitempty"`
}

// NewGetSpoeMessagesOK creates GetSpoeMessagesOK with default headers values
func NewGetSpoeMessagesOK() *GetSpoeMessagesOK {

	return &GetSpoeMessagesOK{}
}

// WithConfigurationVersion adds the configurationVersion to the get spoe messages o k response
func (o *GetSpoeMessagesOK) WithConfigurationVersion(configurationVersion string) *GetSpoeMessagesOK {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get spoe messages o k response
func (o *GetSpoeMessagesOK) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get spoe messages o k response
func (o *GetSpoeMessagesOK) WithPayload(payload *GetSpoeMessagesOKBody) *GetSpoeMessagesOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get spoe messages o k response
func (o *GetSpoeMessagesOK) SetPayload(payload *GetSpoeMessagesOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetSpoeMessagesOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
GetSpoeMessagesDefault General Error

swagger:response getSpoeMessagesDefault
*/
type GetSpoeMessagesDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetSpoeMessagesDefault creates GetSpoeMessagesDefault with default headers values
func NewGetSpoeMessagesDefault(code int) *GetSpoeMessagesDefault {
	if code <= 0 {
		code = 500
	}

	return &GetSpoeMessagesDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get spoe messages default response
func (o *GetSpoeMessagesDefault) WithStatusCode(code int) *GetSpoeMessagesDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get spoe messages default response
func (o *GetSpoeMessagesDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the get spoe messages default response
func (o *GetSpoeMessagesDefault) WithConfigurationVersion(configurationVersion string) *GetSpoeMessagesDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get spoe messages default response
func (o *GetSpoeMessagesDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get spoe messages default response
func (o *GetSpoeMessagesDefault) WithPayload(payload *models.Error) *GetSpoeMessagesDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get spoe messages default response
func (o *GetSpoeMessagesDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetSpoeMessagesDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
