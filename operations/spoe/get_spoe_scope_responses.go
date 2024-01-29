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

	"github.com/haproxytech/client-native/v6/models"
)

// GetSpoeScopeOKCode is the HTTP code returned for type GetSpoeScopeOK
const GetSpoeScopeOKCode int = 200

/*
GetSpoeScopeOK Successful operation

swagger:response getSpoeScopeOK
*/
type GetSpoeScopeOK struct {
	/*Spoe configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *GetSpoeScopeOKBody `json:"body,omitempty"`
}

// NewGetSpoeScopeOK creates GetSpoeScopeOK with default headers values
func NewGetSpoeScopeOK() *GetSpoeScopeOK {

	return &GetSpoeScopeOK{}
}

// WithConfigurationVersion adds the configurationVersion to the get spoe scope o k response
func (o *GetSpoeScopeOK) WithConfigurationVersion(configurationVersion string) *GetSpoeScopeOK {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get spoe scope o k response
func (o *GetSpoeScopeOK) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get spoe scope o k response
func (o *GetSpoeScopeOK) WithPayload(payload *GetSpoeScopeOKBody) *GetSpoeScopeOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get spoe scope o k response
func (o *GetSpoeScopeOK) SetPayload(payload *GetSpoeScopeOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetSpoeScopeOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// GetSpoeScopeNotFoundCode is the HTTP code returned for type GetSpoeScopeNotFound
const GetSpoeScopeNotFoundCode int = 404

/*
GetSpoeScopeNotFound The specified resource was not found

swagger:response getSpoeScopeNotFound
*/
type GetSpoeScopeNotFound struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetSpoeScopeNotFound creates GetSpoeScopeNotFound with default headers values
func NewGetSpoeScopeNotFound() *GetSpoeScopeNotFound {

	return &GetSpoeScopeNotFound{}
}

// WithConfigurationVersion adds the configurationVersion to the get spoe scope not found response
func (o *GetSpoeScopeNotFound) WithConfigurationVersion(configurationVersion string) *GetSpoeScopeNotFound {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get spoe scope not found response
func (o *GetSpoeScopeNotFound) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get spoe scope not found response
func (o *GetSpoeScopeNotFound) WithPayload(payload *models.Error) *GetSpoeScopeNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get spoe scope not found response
func (o *GetSpoeScopeNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetSpoeScopeNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
GetSpoeScopeDefault General Error

swagger:response getSpoeScopeDefault
*/
type GetSpoeScopeDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetSpoeScopeDefault creates GetSpoeScopeDefault with default headers values
func NewGetSpoeScopeDefault(code int) *GetSpoeScopeDefault {
	if code <= 0 {
		code = 500
	}

	return &GetSpoeScopeDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get spoe scope default response
func (o *GetSpoeScopeDefault) WithStatusCode(code int) *GetSpoeScopeDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get spoe scope default response
func (o *GetSpoeScopeDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the get spoe scope default response
func (o *GetSpoeScopeDefault) WithConfigurationVersion(configurationVersion string) *GetSpoeScopeDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get spoe scope default response
func (o *GetSpoeScopeDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get spoe scope default response
func (o *GetSpoeScopeDefault) WithPayload(payload *models.Error) *GetSpoeScopeDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get spoe scope default response
func (o *GetSpoeScopeDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetSpoeScopeDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
