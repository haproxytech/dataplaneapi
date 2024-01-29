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

package ring

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v6/models"
)

// GetRingOKCode is the HTTP code returned for type GetRingOK
const GetRingOKCode int = 200

/*
GetRingOK Successful operation

swagger:response getRingOK
*/
type GetRingOK struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *GetRingOKBody `json:"body,omitempty"`
}

// NewGetRingOK creates GetRingOK with default headers values
func NewGetRingOK() *GetRingOK {

	return &GetRingOK{}
}

// WithConfigurationVersion adds the configurationVersion to the get ring o k response
func (o *GetRingOK) WithConfigurationVersion(configurationVersion string) *GetRingOK {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get ring o k response
func (o *GetRingOK) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get ring o k response
func (o *GetRingOK) WithPayload(payload *GetRingOKBody) *GetRingOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get ring o k response
func (o *GetRingOK) SetPayload(payload *GetRingOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetRingOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// GetRingNotFoundCode is the HTTP code returned for type GetRingNotFound
const GetRingNotFoundCode int = 404

/*
GetRingNotFound The specified resource was not found

swagger:response getRingNotFound
*/
type GetRingNotFound struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetRingNotFound creates GetRingNotFound with default headers values
func NewGetRingNotFound() *GetRingNotFound {

	return &GetRingNotFound{}
}

// WithConfigurationVersion adds the configurationVersion to the get ring not found response
func (o *GetRingNotFound) WithConfigurationVersion(configurationVersion string) *GetRingNotFound {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get ring not found response
func (o *GetRingNotFound) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get ring not found response
func (o *GetRingNotFound) WithPayload(payload *models.Error) *GetRingNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get ring not found response
func (o *GetRingNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetRingNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
GetRingDefault General Error

swagger:response getRingDefault
*/
type GetRingDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetRingDefault creates GetRingDefault with default headers values
func NewGetRingDefault(code int) *GetRingDefault {
	if code <= 0 {
		code = 500
	}

	return &GetRingDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get ring default response
func (o *GetRingDefault) WithStatusCode(code int) *GetRingDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get ring default response
func (o *GetRingDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the get ring default response
func (o *GetRingDefault) WithConfigurationVersion(configurationVersion string) *GetRingDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get ring default response
func (o *GetRingDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get ring default response
func (o *GetRingDefault) WithPayload(payload *models.Error) *GetRingDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get ring default response
func (o *GetRingDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetRingDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
