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

package user

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v6/models"
)

// GetUserOKCode is the HTTP code returned for type GetUserOK
const GetUserOKCode int = 200

/*
GetUserOK Successful operation

swagger:response getUserOK
*/
type GetUserOK struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *GetUserOKBody `json:"body,omitempty"`
}

// NewGetUserOK creates GetUserOK with default headers values
func NewGetUserOK() *GetUserOK {

	return &GetUserOK{}
}

// WithConfigurationVersion adds the configurationVersion to the get user o k response
func (o *GetUserOK) WithConfigurationVersion(configurationVersion string) *GetUserOK {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get user o k response
func (o *GetUserOK) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get user o k response
func (o *GetUserOK) WithPayload(payload *GetUserOKBody) *GetUserOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get user o k response
func (o *GetUserOK) SetPayload(payload *GetUserOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetUserOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// GetUserNotFoundCode is the HTTP code returned for type GetUserNotFound
const GetUserNotFoundCode int = 404

/*
GetUserNotFound The specified resource already exists

swagger:response getUserNotFound
*/
type GetUserNotFound struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetUserNotFound creates GetUserNotFound with default headers values
func NewGetUserNotFound() *GetUserNotFound {

	return &GetUserNotFound{}
}

// WithConfigurationVersion adds the configurationVersion to the get user not found response
func (o *GetUserNotFound) WithConfigurationVersion(configurationVersion string) *GetUserNotFound {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get user not found response
func (o *GetUserNotFound) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get user not found response
func (o *GetUserNotFound) WithPayload(payload *models.Error) *GetUserNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get user not found response
func (o *GetUserNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetUserNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
GetUserDefault General Error

swagger:response getUserDefault
*/
type GetUserDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetUserDefault creates GetUserDefault with default headers values
func NewGetUserDefault(code int) *GetUserDefault {
	if code <= 0 {
		code = 500
	}

	return &GetUserDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get user default response
func (o *GetUserDefault) WithStatusCode(code int) *GetUserDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get user default response
func (o *GetUserDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the get user default response
func (o *GetUserDefault) WithConfigurationVersion(configurationVersion string) *GetUserDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get user default response
func (o *GetUserDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get user default response
func (o *GetUserDefault) WithPayload(payload *models.Error) *GetUserDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get user default response
func (o *GetUserDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetUserDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
