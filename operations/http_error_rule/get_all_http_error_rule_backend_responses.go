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

package http_error_rule

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v6/models"
)

// GetAllHTTPErrorRuleBackendOKCode is the HTTP code returned for type GetAllHTTPErrorRuleBackendOK
const GetAllHTTPErrorRuleBackendOKCode int = 200

/*
GetAllHTTPErrorRuleBackendOK Successful operation

swagger:response getAllHttpErrorRuleBackendOK
*/
type GetAllHTTPErrorRuleBackendOK struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload models.HTTPErrorRules `json:"body,omitempty"`
}

// NewGetAllHTTPErrorRuleBackendOK creates GetAllHTTPErrorRuleBackendOK with default headers values
func NewGetAllHTTPErrorRuleBackendOK() *GetAllHTTPErrorRuleBackendOK {

	return &GetAllHTTPErrorRuleBackendOK{}
}

// WithConfigurationVersion adds the configurationVersion to the get all Http error rule backend o k response
func (o *GetAllHTTPErrorRuleBackendOK) WithConfigurationVersion(configurationVersion string) *GetAllHTTPErrorRuleBackendOK {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get all Http error rule backend o k response
func (o *GetAllHTTPErrorRuleBackendOK) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get all Http error rule backend o k response
func (o *GetAllHTTPErrorRuleBackendOK) WithPayload(payload models.HTTPErrorRules) *GetAllHTTPErrorRuleBackendOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get all Http error rule backend o k response
func (o *GetAllHTTPErrorRuleBackendOK) SetPayload(payload models.HTTPErrorRules) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetAllHTTPErrorRuleBackendOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	// response header Configuration-Version

	configurationVersion := o.ConfigurationVersion
	if configurationVersion != "" {
		rw.Header().Set("Configuration-Version", configurationVersion)
	}

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = models.HTTPErrorRules{}
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

/*
GetAllHTTPErrorRuleBackendDefault General Error

swagger:response getAllHttpErrorRuleBackendDefault
*/
type GetAllHTTPErrorRuleBackendDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetAllHTTPErrorRuleBackendDefault creates GetAllHTTPErrorRuleBackendDefault with default headers values
func NewGetAllHTTPErrorRuleBackendDefault(code int) *GetAllHTTPErrorRuleBackendDefault {
	if code <= 0 {
		code = 500
	}

	return &GetAllHTTPErrorRuleBackendDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get all HTTP error rule backend default response
func (o *GetAllHTTPErrorRuleBackendDefault) WithStatusCode(code int) *GetAllHTTPErrorRuleBackendDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get all HTTP error rule backend default response
func (o *GetAllHTTPErrorRuleBackendDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the get all HTTP error rule backend default response
func (o *GetAllHTTPErrorRuleBackendDefault) WithConfigurationVersion(configurationVersion string) *GetAllHTTPErrorRuleBackendDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get all HTTP error rule backend default response
func (o *GetAllHTTPErrorRuleBackendDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get all HTTP error rule backend default response
func (o *GetAllHTTPErrorRuleBackendDefault) WithPayload(payload *models.Error) *GetAllHTTPErrorRuleBackendDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get all HTTP error rule backend default response
func (o *GetAllHTTPErrorRuleBackendDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetAllHTTPErrorRuleBackendDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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