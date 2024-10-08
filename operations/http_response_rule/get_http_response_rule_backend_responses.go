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

package http_response_rule

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v6/models"
)

// GetHTTPResponseRuleBackendOKCode is the HTTP code returned for type GetHTTPResponseRuleBackendOK
const GetHTTPResponseRuleBackendOKCode int = 200

/*
GetHTTPResponseRuleBackendOK Successful operation

swagger:response getHttpResponseRuleBackendOK
*/
type GetHTTPResponseRuleBackendOK struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.HTTPResponseRule `json:"body,omitempty"`
}

// NewGetHTTPResponseRuleBackendOK creates GetHTTPResponseRuleBackendOK with default headers values
func NewGetHTTPResponseRuleBackendOK() *GetHTTPResponseRuleBackendOK {

	return &GetHTTPResponseRuleBackendOK{}
}

// WithConfigurationVersion adds the configurationVersion to the get Http response rule backend o k response
func (o *GetHTTPResponseRuleBackendOK) WithConfigurationVersion(configurationVersion string) *GetHTTPResponseRuleBackendOK {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get Http response rule backend o k response
func (o *GetHTTPResponseRuleBackendOK) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get Http response rule backend o k response
func (o *GetHTTPResponseRuleBackendOK) WithPayload(payload *models.HTTPResponseRule) *GetHTTPResponseRuleBackendOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get Http response rule backend o k response
func (o *GetHTTPResponseRuleBackendOK) SetPayload(payload *models.HTTPResponseRule) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetHTTPResponseRuleBackendOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// GetHTTPResponseRuleBackendNotFoundCode is the HTTP code returned for type GetHTTPResponseRuleBackendNotFound
const GetHTTPResponseRuleBackendNotFoundCode int = 404

/*
GetHTTPResponseRuleBackendNotFound The specified resource was not found

swagger:response getHttpResponseRuleBackendNotFound
*/
type GetHTTPResponseRuleBackendNotFound struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetHTTPResponseRuleBackendNotFound creates GetHTTPResponseRuleBackendNotFound with default headers values
func NewGetHTTPResponseRuleBackendNotFound() *GetHTTPResponseRuleBackendNotFound {

	return &GetHTTPResponseRuleBackendNotFound{}
}

// WithConfigurationVersion adds the configurationVersion to the get Http response rule backend not found response
func (o *GetHTTPResponseRuleBackendNotFound) WithConfigurationVersion(configurationVersion string) *GetHTTPResponseRuleBackendNotFound {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get Http response rule backend not found response
func (o *GetHTTPResponseRuleBackendNotFound) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get Http response rule backend not found response
func (o *GetHTTPResponseRuleBackendNotFound) WithPayload(payload *models.Error) *GetHTTPResponseRuleBackendNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get Http response rule backend not found response
func (o *GetHTTPResponseRuleBackendNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetHTTPResponseRuleBackendNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
GetHTTPResponseRuleBackendDefault General Error

swagger:response getHttpResponseRuleBackendDefault
*/
type GetHTTPResponseRuleBackendDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetHTTPResponseRuleBackendDefault creates GetHTTPResponseRuleBackendDefault with default headers values
func NewGetHTTPResponseRuleBackendDefault(code int) *GetHTTPResponseRuleBackendDefault {
	if code <= 0 {
		code = 500
	}

	return &GetHTTPResponseRuleBackendDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get HTTP response rule backend default response
func (o *GetHTTPResponseRuleBackendDefault) WithStatusCode(code int) *GetHTTPResponseRuleBackendDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get HTTP response rule backend default response
func (o *GetHTTPResponseRuleBackendDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the get HTTP response rule backend default response
func (o *GetHTTPResponseRuleBackendDefault) WithConfigurationVersion(configurationVersion string) *GetHTTPResponseRuleBackendDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get HTTP response rule backend default response
func (o *GetHTTPResponseRuleBackendDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get HTTP response rule backend default response
func (o *GetHTTPResponseRuleBackendDefault) WithPayload(payload *models.Error) *GetHTTPResponseRuleBackendDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get HTTP response rule backend default response
func (o *GetHTTPResponseRuleBackendDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetHTTPResponseRuleBackendDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
