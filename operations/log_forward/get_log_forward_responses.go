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

package log_forward

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v4/models"
)

// GetLogForwardOKCode is the HTTP code returned for type GetLogForwardOK
const GetLogForwardOKCode int = 200

/*GetLogForwardOK Successful operation

swagger:response getLogForwardOK
*/
type GetLogForwardOK struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *GetLogForwardOKBody `json:"body,omitempty"`
}

// NewGetLogForwardOK creates GetLogForwardOK with default headers values
func NewGetLogForwardOK() *GetLogForwardOK {

	return &GetLogForwardOK{}
}

// WithConfigurationVersion adds the configurationVersion to the get log forward o k response
func (o *GetLogForwardOK) WithConfigurationVersion(configurationVersion string) *GetLogForwardOK {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get log forward o k response
func (o *GetLogForwardOK) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get log forward o k response
func (o *GetLogForwardOK) WithPayload(payload *GetLogForwardOKBody) *GetLogForwardOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get log forward o k response
func (o *GetLogForwardOK) SetPayload(payload *GetLogForwardOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetLogForwardOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// GetLogForwardNotFoundCode is the HTTP code returned for type GetLogForwardNotFound
const GetLogForwardNotFoundCode int = 404

/*GetLogForwardNotFound The specified resource was not found

swagger:response getLogForwardNotFound
*/
type GetLogForwardNotFound struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetLogForwardNotFound creates GetLogForwardNotFound with default headers values
func NewGetLogForwardNotFound() *GetLogForwardNotFound {

	return &GetLogForwardNotFound{}
}

// WithConfigurationVersion adds the configurationVersion to the get log forward not found response
func (o *GetLogForwardNotFound) WithConfigurationVersion(configurationVersion string) *GetLogForwardNotFound {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get log forward not found response
func (o *GetLogForwardNotFound) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get log forward not found response
func (o *GetLogForwardNotFound) WithPayload(payload *models.Error) *GetLogForwardNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get log forward not found response
func (o *GetLogForwardNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetLogForwardNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

/*GetLogForwardDefault General Error

swagger:response getLogForwardDefault
*/
type GetLogForwardDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetLogForwardDefault creates GetLogForwardDefault with default headers values
func NewGetLogForwardDefault(code int) *GetLogForwardDefault {
	if code <= 0 {
		code = 500
	}

	return &GetLogForwardDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get log forward default response
func (o *GetLogForwardDefault) WithStatusCode(code int) *GetLogForwardDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get log forward default response
func (o *GetLogForwardDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the get log forward default response
func (o *GetLogForwardDefault) WithConfigurationVersion(configurationVersion string) *GetLogForwardDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get log forward default response
func (o *GetLogForwardDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get log forward default response
func (o *GetLogForwardDefault) WithPayload(payload *models.Error) *GetLogForwardDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get log forward default response
func (o *GetLogForwardDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetLogForwardDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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