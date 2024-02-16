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

	"github.com/haproxytech/client-native/v6/models"
)

// GetConsulOKCode is the HTTP code returned for type GetConsulOK
const GetConsulOKCode int = 200

/*
GetConsulOK Successful operation

swagger:response getConsulOK
*/
type GetConsulOK struct {

	/*
	  In: Body
	*/
	Payload *models.Consul `json:"body,omitempty"`
}

// NewGetConsulOK creates GetConsulOK with default headers values
func NewGetConsulOK() *GetConsulOK {

	return &GetConsulOK{}
}

// WithPayload adds the payload to the get consul o k response
func (o *GetConsulOK) WithPayload(payload *models.Consul) *GetConsulOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get consul o k response
func (o *GetConsulOK) SetPayload(payload *models.Consul) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetConsulOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetConsulNotFoundCode is the HTTP code returned for type GetConsulNotFound
const GetConsulNotFoundCode int = 404

/*
GetConsulNotFound The specified resource was not found

swagger:response getConsulNotFound
*/
type GetConsulNotFound struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetConsulNotFound creates GetConsulNotFound with default headers values
func NewGetConsulNotFound() *GetConsulNotFound {

	return &GetConsulNotFound{}
}

// WithConfigurationVersion adds the configurationVersion to the get consul not found response
func (o *GetConsulNotFound) WithConfigurationVersion(configurationVersion string) *GetConsulNotFound {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get consul not found response
func (o *GetConsulNotFound) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get consul not found response
func (o *GetConsulNotFound) WithPayload(payload *models.Error) *GetConsulNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get consul not found response
func (o *GetConsulNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetConsulNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
GetConsulDefault General Error

swagger:response getConsulDefault
*/
type GetConsulDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetConsulDefault creates GetConsulDefault with default headers values
func NewGetConsulDefault(code int) *GetConsulDefault {
	if code <= 0 {
		code = 500
	}

	return &GetConsulDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get consul default response
func (o *GetConsulDefault) WithStatusCode(code int) *GetConsulDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get consul default response
func (o *GetConsulDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the get consul default response
func (o *GetConsulDefault) WithConfigurationVersion(configurationVersion string) *GetConsulDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get consul default response
func (o *GetConsulDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get consul default response
func (o *GetConsulDefault) WithPayload(payload *models.Error) *GetConsulDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get consul default response
func (o *GetConsulDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetConsulDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
