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

package log_target

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v6/models"
)

// GetLogTargetFrontendOKCode is the HTTP code returned for type GetLogTargetFrontendOK
const GetLogTargetFrontendOKCode int = 200

/*
GetLogTargetFrontendOK Successful operation

swagger:response getLogTargetFrontendOK
*/
type GetLogTargetFrontendOK struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.LogTarget `json:"body,omitempty"`
}

// NewGetLogTargetFrontendOK creates GetLogTargetFrontendOK with default headers values
func NewGetLogTargetFrontendOK() *GetLogTargetFrontendOK {

	return &GetLogTargetFrontendOK{}
}

// WithConfigurationVersion adds the configurationVersion to the get log target frontend o k response
func (o *GetLogTargetFrontendOK) WithConfigurationVersion(configurationVersion string) *GetLogTargetFrontendOK {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get log target frontend o k response
func (o *GetLogTargetFrontendOK) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get log target frontend o k response
func (o *GetLogTargetFrontendOK) WithPayload(payload *models.LogTarget) *GetLogTargetFrontendOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get log target frontend o k response
func (o *GetLogTargetFrontendOK) SetPayload(payload *models.LogTarget) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetLogTargetFrontendOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// GetLogTargetFrontendNotFoundCode is the HTTP code returned for type GetLogTargetFrontendNotFound
const GetLogTargetFrontendNotFoundCode int = 404

/*
GetLogTargetFrontendNotFound The specified resource was not found

swagger:response getLogTargetFrontendNotFound
*/
type GetLogTargetFrontendNotFound struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetLogTargetFrontendNotFound creates GetLogTargetFrontendNotFound with default headers values
func NewGetLogTargetFrontendNotFound() *GetLogTargetFrontendNotFound {

	return &GetLogTargetFrontendNotFound{}
}

// WithConfigurationVersion adds the configurationVersion to the get log target frontend not found response
func (o *GetLogTargetFrontendNotFound) WithConfigurationVersion(configurationVersion string) *GetLogTargetFrontendNotFound {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get log target frontend not found response
func (o *GetLogTargetFrontendNotFound) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get log target frontend not found response
func (o *GetLogTargetFrontendNotFound) WithPayload(payload *models.Error) *GetLogTargetFrontendNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get log target frontend not found response
func (o *GetLogTargetFrontendNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetLogTargetFrontendNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
GetLogTargetFrontendDefault General Error

swagger:response getLogTargetFrontendDefault
*/
type GetLogTargetFrontendDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetLogTargetFrontendDefault creates GetLogTargetFrontendDefault with default headers values
func NewGetLogTargetFrontendDefault(code int) *GetLogTargetFrontendDefault {
	if code <= 0 {
		code = 500
	}

	return &GetLogTargetFrontendDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get log target frontend default response
func (o *GetLogTargetFrontendDefault) WithStatusCode(code int) *GetLogTargetFrontendDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get log target frontend default response
func (o *GetLogTargetFrontendDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the get log target frontend default response
func (o *GetLogTargetFrontendDefault) WithConfigurationVersion(configurationVersion string) *GetLogTargetFrontendDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get log target frontend default response
func (o *GetLogTargetFrontendDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get log target frontend default response
func (o *GetLogTargetFrontendDefault) WithPayload(payload *models.Error) *GetLogTargetFrontendDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get log target frontend default response
func (o *GetLogTargetFrontendDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetLogTargetFrontendDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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