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

package cache

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v3/models"
)

// GetCacheOKCode is the HTTP code returned for type GetCacheOK
const GetCacheOKCode int = 200

/*
GetCacheOK Successful operation

swagger:response getCacheOK
*/
type GetCacheOK struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *GetCacheOKBody `json:"body,omitempty"`
}

// NewGetCacheOK creates GetCacheOK with default headers values
func NewGetCacheOK() *GetCacheOK {

	return &GetCacheOK{}
}

// WithConfigurationVersion adds the configurationVersion to the get cache o k response
func (o *GetCacheOK) WithConfigurationVersion(configurationVersion string) *GetCacheOK {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get cache o k response
func (o *GetCacheOK) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get cache o k response
func (o *GetCacheOK) WithPayload(payload *GetCacheOKBody) *GetCacheOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get cache o k response
func (o *GetCacheOK) SetPayload(payload *GetCacheOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetCacheOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// GetCacheNotFoundCode is the HTTP code returned for type GetCacheNotFound
const GetCacheNotFoundCode int = 404

/*
GetCacheNotFound The specified resource was not found

swagger:response getCacheNotFound
*/
type GetCacheNotFound struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetCacheNotFound creates GetCacheNotFound with default headers values
func NewGetCacheNotFound() *GetCacheNotFound {

	return &GetCacheNotFound{}
}

// WithConfigurationVersion adds the configurationVersion to the get cache not found response
func (o *GetCacheNotFound) WithConfigurationVersion(configurationVersion string) *GetCacheNotFound {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get cache not found response
func (o *GetCacheNotFound) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get cache not found response
func (o *GetCacheNotFound) WithPayload(payload *models.Error) *GetCacheNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get cache not found response
func (o *GetCacheNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetCacheNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
GetCacheDefault General Error

swagger:response getCacheDefault
*/
type GetCacheDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetCacheDefault creates GetCacheDefault with default headers values
func NewGetCacheDefault(code int) *GetCacheDefault {
	if code <= 0 {
		code = 500
	}

	return &GetCacheDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get cache default response
func (o *GetCacheDefault) WithStatusCode(code int) *GetCacheDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get cache default response
func (o *GetCacheDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the get cache default response
func (o *GetCacheDefault) WithConfigurationVersion(configurationVersion string) *GetCacheDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get cache default response
func (o *GetCacheDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get cache default response
func (o *GetCacheDefault) WithPayload(payload *models.Error) *GetCacheDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get cache default response
func (o *GetCacheDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetCacheDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
