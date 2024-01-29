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

package resolver

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v6/models"
)

// GetResolverOKCode is the HTTP code returned for type GetResolverOK
const GetResolverOKCode int = 200

/*
GetResolverOK Successful operation

swagger:response getResolverOK
*/
type GetResolverOK struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *GetResolverOKBody `json:"body,omitempty"`
}

// NewGetResolverOK creates GetResolverOK with default headers values
func NewGetResolverOK() *GetResolverOK {

	return &GetResolverOK{}
}

// WithConfigurationVersion adds the configurationVersion to the get resolver o k response
func (o *GetResolverOK) WithConfigurationVersion(configurationVersion string) *GetResolverOK {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get resolver o k response
func (o *GetResolverOK) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get resolver o k response
func (o *GetResolverOK) WithPayload(payload *GetResolverOKBody) *GetResolverOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get resolver o k response
func (o *GetResolverOK) SetPayload(payload *GetResolverOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetResolverOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// GetResolverNotFoundCode is the HTTP code returned for type GetResolverNotFound
const GetResolverNotFoundCode int = 404

/*
GetResolverNotFound The specified resource was not found

swagger:response getResolverNotFound
*/
type GetResolverNotFound struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetResolverNotFound creates GetResolverNotFound with default headers values
func NewGetResolverNotFound() *GetResolverNotFound {

	return &GetResolverNotFound{}
}

// WithConfigurationVersion adds the configurationVersion to the get resolver not found response
func (o *GetResolverNotFound) WithConfigurationVersion(configurationVersion string) *GetResolverNotFound {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get resolver not found response
func (o *GetResolverNotFound) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get resolver not found response
func (o *GetResolverNotFound) WithPayload(payload *models.Error) *GetResolverNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get resolver not found response
func (o *GetResolverNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetResolverNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
GetResolverDefault General Error

swagger:response getResolverDefault
*/
type GetResolverDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetResolverDefault creates GetResolverDefault with default headers values
func NewGetResolverDefault(code int) *GetResolverDefault {
	if code <= 0 {
		code = 500
	}

	return &GetResolverDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get resolver default response
func (o *GetResolverDefault) WithStatusCode(code int) *GetResolverDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get resolver default response
func (o *GetResolverDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the get resolver default response
func (o *GetResolverDefault) WithConfigurationVersion(configurationVersion string) *GetResolverDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get resolver default response
func (o *GetResolverDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get resolver default response
func (o *GetResolverDefault) WithPayload(payload *models.Error) *GetResolverDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get resolver default response
func (o *GetResolverDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetResolverDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
