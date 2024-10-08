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

package http_after_response_rule

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v6/models"
)

// GetAllHTTPAfterResponseRuleBackendOKCode is the HTTP code returned for type GetAllHTTPAfterResponseRuleBackendOK
const GetAllHTTPAfterResponseRuleBackendOKCode int = 200

/*
GetAllHTTPAfterResponseRuleBackendOK Successful operation

swagger:response getAllHttpAfterResponseRuleBackendOK
*/
type GetAllHTTPAfterResponseRuleBackendOK struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload models.HTTPAfterResponseRules `json:"body,omitempty"`
}

// NewGetAllHTTPAfterResponseRuleBackendOK creates GetAllHTTPAfterResponseRuleBackendOK with default headers values
func NewGetAllHTTPAfterResponseRuleBackendOK() *GetAllHTTPAfterResponseRuleBackendOK {

	return &GetAllHTTPAfterResponseRuleBackendOK{}
}

// WithConfigurationVersion adds the configurationVersion to the get all Http after response rule backend o k response
func (o *GetAllHTTPAfterResponseRuleBackendOK) WithConfigurationVersion(configurationVersion string) *GetAllHTTPAfterResponseRuleBackendOK {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get all Http after response rule backend o k response
func (o *GetAllHTTPAfterResponseRuleBackendOK) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get all Http after response rule backend o k response
func (o *GetAllHTTPAfterResponseRuleBackendOK) WithPayload(payload models.HTTPAfterResponseRules) *GetAllHTTPAfterResponseRuleBackendOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get all Http after response rule backend o k response
func (o *GetAllHTTPAfterResponseRuleBackendOK) SetPayload(payload models.HTTPAfterResponseRules) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetAllHTTPAfterResponseRuleBackendOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	// response header Configuration-Version

	configurationVersion := o.ConfigurationVersion
	if configurationVersion != "" {
		rw.Header().Set("Configuration-Version", configurationVersion)
	}

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = models.HTTPAfterResponseRules{}
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

/*
GetAllHTTPAfterResponseRuleBackendDefault General Error

swagger:response getAllHttpAfterResponseRuleBackendDefault
*/
type GetAllHTTPAfterResponseRuleBackendDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetAllHTTPAfterResponseRuleBackendDefault creates GetAllHTTPAfterResponseRuleBackendDefault with default headers values
func NewGetAllHTTPAfterResponseRuleBackendDefault(code int) *GetAllHTTPAfterResponseRuleBackendDefault {
	if code <= 0 {
		code = 500
	}

	return &GetAllHTTPAfterResponseRuleBackendDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get all HTTP after response rule backend default response
func (o *GetAllHTTPAfterResponseRuleBackendDefault) WithStatusCode(code int) *GetAllHTTPAfterResponseRuleBackendDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get all HTTP after response rule backend default response
func (o *GetAllHTTPAfterResponseRuleBackendDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the get all HTTP after response rule backend default response
func (o *GetAllHTTPAfterResponseRuleBackendDefault) WithConfigurationVersion(configurationVersion string) *GetAllHTTPAfterResponseRuleBackendDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get all HTTP after response rule backend default response
func (o *GetAllHTTPAfterResponseRuleBackendDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get all HTTP after response rule backend default response
func (o *GetAllHTTPAfterResponseRuleBackendDefault) WithPayload(payload *models.Error) *GetAllHTTPAfterResponseRuleBackendDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get all HTTP after response rule backend default response
func (o *GetAllHTTPAfterResponseRuleBackendDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetAllHTTPAfterResponseRuleBackendDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
