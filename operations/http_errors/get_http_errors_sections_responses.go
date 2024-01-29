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

package http_errors

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v6/models"
)

// GetHTTPErrorsSectionsOKCode is the HTTP code returned for type GetHTTPErrorsSectionsOK
const GetHTTPErrorsSectionsOKCode int = 200

/*
GetHTTPErrorsSectionsOK Successful operation

swagger:response getHttpErrorsSectionsOK
*/
type GetHTTPErrorsSectionsOK struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *GetHTTPErrorsSectionsOKBody `json:"body,omitempty"`
}

// NewGetHTTPErrorsSectionsOK creates GetHTTPErrorsSectionsOK with default headers values
func NewGetHTTPErrorsSectionsOK() *GetHTTPErrorsSectionsOK {

	return &GetHTTPErrorsSectionsOK{}
}

// WithConfigurationVersion adds the configurationVersion to the get Http errors sections o k response
func (o *GetHTTPErrorsSectionsOK) WithConfigurationVersion(configurationVersion string) *GetHTTPErrorsSectionsOK {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get Http errors sections o k response
func (o *GetHTTPErrorsSectionsOK) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get Http errors sections o k response
func (o *GetHTTPErrorsSectionsOK) WithPayload(payload *GetHTTPErrorsSectionsOKBody) *GetHTTPErrorsSectionsOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get Http errors sections o k response
func (o *GetHTTPErrorsSectionsOK) SetPayload(payload *GetHTTPErrorsSectionsOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetHTTPErrorsSectionsOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
GetHTTPErrorsSectionsDefault General Error

swagger:response getHttpErrorsSectionsDefault
*/
type GetHTTPErrorsSectionsDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetHTTPErrorsSectionsDefault creates GetHTTPErrorsSectionsDefault with default headers values
func NewGetHTTPErrorsSectionsDefault(code int) *GetHTTPErrorsSectionsDefault {
	if code <= 0 {
		code = 500
	}

	return &GetHTTPErrorsSectionsDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get HTTP errors sections default response
func (o *GetHTTPErrorsSectionsDefault) WithStatusCode(code int) *GetHTTPErrorsSectionsDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get HTTP errors sections default response
func (o *GetHTTPErrorsSectionsDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the get HTTP errors sections default response
func (o *GetHTTPErrorsSectionsDefault) WithConfigurationVersion(configurationVersion string) *GetHTTPErrorsSectionsDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get HTTP errors sections default response
func (o *GetHTTPErrorsSectionsDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get HTTP errors sections default response
func (o *GetHTTPErrorsSectionsDefault) WithPayload(payload *models.Error) *GetHTTPErrorsSectionsDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get HTTP errors sections default response
func (o *GetHTTPErrorsSectionsDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetHTTPErrorsSectionsDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
