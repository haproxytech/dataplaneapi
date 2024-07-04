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

package http_request_rule

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v6/models"
)

// GetHTTPRequestRuleFrontendOKCode is the HTTP code returned for type GetHTTPRequestRuleFrontendOK
const GetHTTPRequestRuleFrontendOKCode int = 200

/*
GetHTTPRequestRuleFrontendOK Successful operation

swagger:response getHttpRequestRuleFrontendOK
*/
type GetHTTPRequestRuleFrontendOK struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.HTTPRequestRule `json:"body,omitempty"`
}

// NewGetHTTPRequestRuleFrontendOK creates GetHTTPRequestRuleFrontendOK with default headers values
func NewGetHTTPRequestRuleFrontendOK() *GetHTTPRequestRuleFrontendOK {

	return &GetHTTPRequestRuleFrontendOK{}
}

// WithConfigurationVersion adds the configurationVersion to the get Http request rule frontend o k response
func (o *GetHTTPRequestRuleFrontendOK) WithConfigurationVersion(configurationVersion string) *GetHTTPRequestRuleFrontendOK {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get Http request rule frontend o k response
func (o *GetHTTPRequestRuleFrontendOK) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get Http request rule frontend o k response
func (o *GetHTTPRequestRuleFrontendOK) WithPayload(payload *models.HTTPRequestRule) *GetHTTPRequestRuleFrontendOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get Http request rule frontend o k response
func (o *GetHTTPRequestRuleFrontendOK) SetPayload(payload *models.HTTPRequestRule) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetHTTPRequestRuleFrontendOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// GetHTTPRequestRuleFrontendNotFoundCode is the HTTP code returned for type GetHTTPRequestRuleFrontendNotFound
const GetHTTPRequestRuleFrontendNotFoundCode int = 404

/*
GetHTTPRequestRuleFrontendNotFound The specified resource was not found

swagger:response getHttpRequestRuleFrontendNotFound
*/
type GetHTTPRequestRuleFrontendNotFound struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetHTTPRequestRuleFrontendNotFound creates GetHTTPRequestRuleFrontendNotFound with default headers values
func NewGetHTTPRequestRuleFrontendNotFound() *GetHTTPRequestRuleFrontendNotFound {

	return &GetHTTPRequestRuleFrontendNotFound{}
}

// WithConfigurationVersion adds the configurationVersion to the get Http request rule frontend not found response
func (o *GetHTTPRequestRuleFrontendNotFound) WithConfigurationVersion(configurationVersion string) *GetHTTPRequestRuleFrontendNotFound {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get Http request rule frontend not found response
func (o *GetHTTPRequestRuleFrontendNotFound) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get Http request rule frontend not found response
func (o *GetHTTPRequestRuleFrontendNotFound) WithPayload(payload *models.Error) *GetHTTPRequestRuleFrontendNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get Http request rule frontend not found response
func (o *GetHTTPRequestRuleFrontendNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetHTTPRequestRuleFrontendNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
GetHTTPRequestRuleFrontendDefault General Error

swagger:response getHttpRequestRuleFrontendDefault
*/
type GetHTTPRequestRuleFrontendDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetHTTPRequestRuleFrontendDefault creates GetHTTPRequestRuleFrontendDefault with default headers values
func NewGetHTTPRequestRuleFrontendDefault(code int) *GetHTTPRequestRuleFrontendDefault {
	if code <= 0 {
		code = 500
	}

	return &GetHTTPRequestRuleFrontendDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get HTTP request rule frontend default response
func (o *GetHTTPRequestRuleFrontendDefault) WithStatusCode(code int) *GetHTTPRequestRuleFrontendDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get HTTP request rule frontend default response
func (o *GetHTTPRequestRuleFrontendDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the get HTTP request rule frontend default response
func (o *GetHTTPRequestRuleFrontendDefault) WithConfigurationVersion(configurationVersion string) *GetHTTPRequestRuleFrontendDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get HTTP request rule frontend default response
func (o *GetHTTPRequestRuleFrontendDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get HTTP request rule frontend default response
func (o *GetHTTPRequestRuleFrontendDefault) WithPayload(payload *models.Error) *GetHTTPRequestRuleFrontendDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get HTTP request rule frontend default response
func (o *GetHTTPRequestRuleFrontendDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetHTTPRequestRuleFrontendDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
