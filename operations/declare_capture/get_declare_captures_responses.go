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

package declare_capture

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v3/models"
)

// GetDeclareCapturesOKCode is the HTTP code returned for type GetDeclareCapturesOK
const GetDeclareCapturesOKCode int = 200

/*
GetDeclareCapturesOK Successful operation

swagger:response getDeclareCapturesOK
*/
type GetDeclareCapturesOK struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *GetDeclareCapturesOKBody `json:"body,omitempty"`
}

// NewGetDeclareCapturesOK creates GetDeclareCapturesOK with default headers values
func NewGetDeclareCapturesOK() *GetDeclareCapturesOK {

	return &GetDeclareCapturesOK{}
}

// WithConfigurationVersion adds the configurationVersion to the get declare captures o k response
func (o *GetDeclareCapturesOK) WithConfigurationVersion(configurationVersion string) *GetDeclareCapturesOK {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get declare captures o k response
func (o *GetDeclareCapturesOK) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get declare captures o k response
func (o *GetDeclareCapturesOK) WithPayload(payload *GetDeclareCapturesOKBody) *GetDeclareCapturesOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get declare captures o k response
func (o *GetDeclareCapturesOK) SetPayload(payload *GetDeclareCapturesOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetDeclareCapturesOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
GetDeclareCapturesDefault General Error

swagger:response getDeclareCapturesDefault
*/
type GetDeclareCapturesDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetDeclareCapturesDefault creates GetDeclareCapturesDefault with default headers values
func NewGetDeclareCapturesDefault(code int) *GetDeclareCapturesDefault {
	if code <= 0 {
		code = 500
	}

	return &GetDeclareCapturesDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get declare captures default response
func (o *GetDeclareCapturesDefault) WithStatusCode(code int) *GetDeclareCapturesDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get declare captures default response
func (o *GetDeclareCapturesDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the get declare captures default response
func (o *GetDeclareCapturesDefault) WithConfigurationVersion(configurationVersion string) *GetDeclareCapturesDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get declare captures default response
func (o *GetDeclareCapturesDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get declare captures default response
func (o *GetDeclareCapturesDefault) WithPayload(payload *models.Error) *GetDeclareCapturesDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get declare captures default response
func (o *GetDeclareCapturesDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetDeclareCapturesDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
