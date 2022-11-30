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

package service_discovery

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v3/models"
)

// GetConsulsOKCode is the HTTP code returned for type GetConsulsOK
const GetConsulsOKCode int = 200

/*
GetConsulsOK Successful operation

swagger:response getConsulsOK
*/
type GetConsulsOK struct {

	/*
	  In: Body
	*/
	Payload *GetConsulsOKBody `json:"body,omitempty"`
}

// NewGetConsulsOK creates GetConsulsOK with default headers values
func NewGetConsulsOK() *GetConsulsOK {

	return &GetConsulsOK{}
}

// WithPayload adds the payload to the get consuls o k response
func (o *GetConsulsOK) WithPayload(payload *GetConsulsOKBody) *GetConsulsOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get consuls o k response
func (o *GetConsulsOK) SetPayload(payload *GetConsulsOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetConsulsOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*
GetConsulsDefault General Error

swagger:response getConsulsDefault
*/
type GetConsulsDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetConsulsDefault creates GetConsulsDefault with default headers values
func NewGetConsulsDefault(code int) *GetConsulsDefault {
	if code <= 0 {
		code = 500
	}

	return &GetConsulsDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get consuls default response
func (o *GetConsulsDefault) WithStatusCode(code int) *GetConsulsDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get consuls default response
func (o *GetConsulsDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the get consuls default response
func (o *GetConsulsDefault) WithConfigurationVersion(configurationVersion string) *GetConsulsDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get consuls default response
func (o *GetConsulsDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get consuls default response
func (o *GetConsulsDefault) WithPayload(payload *models.Error) *GetConsulsDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get consuls default response
func (o *GetConsulsDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetConsulsDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
