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

	"github.com/haproxytech/client-native/v4/models"
)

// GetHTTPResponseRuleOKCode is the HTTP code returned for type GetHTTPResponseRuleOK
const GetHTTPResponseRuleOKCode int = 200

/*
GetHTTPResponseRuleOK Successful operation

swagger:response getHttpResponseRuleOK
*/
type GetHTTPResponseRuleOK struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *GetHTTPResponseRuleOKBody `json:"body,omitempty"`
}

// NewGetHTTPResponseRuleOK creates GetHTTPResponseRuleOK with default headers values
func NewGetHTTPResponseRuleOK() *GetHTTPResponseRuleOK {

	return &GetHTTPResponseRuleOK{}
}

// WithConfigurationVersion adds the configurationVersion to the get Http response rule o k response
func (o *GetHTTPResponseRuleOK) WithConfigurationVersion(configurationVersion string) *GetHTTPResponseRuleOK {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get Http response rule o k response
func (o *GetHTTPResponseRuleOK) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get Http response rule o k response
func (o *GetHTTPResponseRuleOK) WithPayload(payload *GetHTTPResponseRuleOKBody) *GetHTTPResponseRuleOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get Http response rule o k response
func (o *GetHTTPResponseRuleOK) SetPayload(payload *GetHTTPResponseRuleOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetHTTPResponseRuleOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// GetHTTPResponseRuleNotFoundCode is the HTTP code returned for type GetHTTPResponseRuleNotFound
const GetHTTPResponseRuleNotFoundCode int = 404

/*
GetHTTPResponseRuleNotFound The specified resource was not found

swagger:response getHttpResponseRuleNotFound
*/
type GetHTTPResponseRuleNotFound struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetHTTPResponseRuleNotFound creates GetHTTPResponseRuleNotFound with default headers values
func NewGetHTTPResponseRuleNotFound() *GetHTTPResponseRuleNotFound {

	return &GetHTTPResponseRuleNotFound{}
}

// WithConfigurationVersion adds the configurationVersion to the get Http response rule not found response
func (o *GetHTTPResponseRuleNotFound) WithConfigurationVersion(configurationVersion string) *GetHTTPResponseRuleNotFound {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get Http response rule not found response
func (o *GetHTTPResponseRuleNotFound) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get Http response rule not found response
func (o *GetHTTPResponseRuleNotFound) WithPayload(payload *models.Error) *GetHTTPResponseRuleNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get Http response rule not found response
func (o *GetHTTPResponseRuleNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetHTTPResponseRuleNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
GetHTTPResponseRuleDefault General Error

swagger:response getHttpResponseRuleDefault
*/
type GetHTTPResponseRuleDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetHTTPResponseRuleDefault creates GetHTTPResponseRuleDefault with default headers values
func NewGetHTTPResponseRuleDefault(code int) *GetHTTPResponseRuleDefault {
	if code <= 0 {
		code = 500
	}

	return &GetHTTPResponseRuleDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get HTTP response rule default response
func (o *GetHTTPResponseRuleDefault) WithStatusCode(code int) *GetHTTPResponseRuleDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get HTTP response rule default response
func (o *GetHTTPResponseRuleDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the get HTTP response rule default response
func (o *GetHTTPResponseRuleDefault) WithConfigurationVersion(configurationVersion string) *GetHTTPResponseRuleDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get HTTP response rule default response
func (o *GetHTTPResponseRuleDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get HTTP response rule default response
func (o *GetHTTPResponseRuleDefault) WithPayload(payload *models.Error) *GetHTTPResponseRuleDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get HTTP response rule default response
func (o *GetHTTPResponseRuleDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetHTTPResponseRuleDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
