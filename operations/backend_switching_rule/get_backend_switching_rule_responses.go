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

package backend_switching_rule

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v6/models"
)

// GetBackendSwitchingRuleOKCode is the HTTP code returned for type GetBackendSwitchingRuleOK
const GetBackendSwitchingRuleOKCode int = 200

/*
GetBackendSwitchingRuleOK Successful operation

swagger:response getBackendSwitchingRuleOK
*/
type GetBackendSwitchingRuleOK struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *GetBackendSwitchingRuleOKBody `json:"body,omitempty"`
}

// NewGetBackendSwitchingRuleOK creates GetBackendSwitchingRuleOK with default headers values
func NewGetBackendSwitchingRuleOK() *GetBackendSwitchingRuleOK {

	return &GetBackendSwitchingRuleOK{}
}

// WithConfigurationVersion adds the configurationVersion to the get backend switching rule o k response
func (o *GetBackendSwitchingRuleOK) WithConfigurationVersion(configurationVersion string) *GetBackendSwitchingRuleOK {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get backend switching rule o k response
func (o *GetBackendSwitchingRuleOK) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get backend switching rule o k response
func (o *GetBackendSwitchingRuleOK) WithPayload(payload *GetBackendSwitchingRuleOKBody) *GetBackendSwitchingRuleOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get backend switching rule o k response
func (o *GetBackendSwitchingRuleOK) SetPayload(payload *GetBackendSwitchingRuleOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetBackendSwitchingRuleOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// GetBackendSwitchingRuleNotFoundCode is the HTTP code returned for type GetBackendSwitchingRuleNotFound
const GetBackendSwitchingRuleNotFoundCode int = 404

/*
GetBackendSwitchingRuleNotFound The specified resource was not found

swagger:response getBackendSwitchingRuleNotFound
*/
type GetBackendSwitchingRuleNotFound struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetBackendSwitchingRuleNotFound creates GetBackendSwitchingRuleNotFound with default headers values
func NewGetBackendSwitchingRuleNotFound() *GetBackendSwitchingRuleNotFound {

	return &GetBackendSwitchingRuleNotFound{}
}

// WithConfigurationVersion adds the configurationVersion to the get backend switching rule not found response
func (o *GetBackendSwitchingRuleNotFound) WithConfigurationVersion(configurationVersion string) *GetBackendSwitchingRuleNotFound {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get backend switching rule not found response
func (o *GetBackendSwitchingRuleNotFound) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get backend switching rule not found response
func (o *GetBackendSwitchingRuleNotFound) WithPayload(payload *models.Error) *GetBackendSwitchingRuleNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get backend switching rule not found response
func (o *GetBackendSwitchingRuleNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetBackendSwitchingRuleNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
GetBackendSwitchingRuleDefault General Error

swagger:response getBackendSwitchingRuleDefault
*/
type GetBackendSwitchingRuleDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetBackendSwitchingRuleDefault creates GetBackendSwitchingRuleDefault with default headers values
func NewGetBackendSwitchingRuleDefault(code int) *GetBackendSwitchingRuleDefault {
	if code <= 0 {
		code = 500
	}

	return &GetBackendSwitchingRuleDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get backend switching rule default response
func (o *GetBackendSwitchingRuleDefault) WithStatusCode(code int) *GetBackendSwitchingRuleDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get backend switching rule default response
func (o *GetBackendSwitchingRuleDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the get backend switching rule default response
func (o *GetBackendSwitchingRuleDefault) WithConfigurationVersion(configurationVersion string) *GetBackendSwitchingRuleDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get backend switching rule default response
func (o *GetBackendSwitchingRuleDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get backend switching rule default response
func (o *GetBackendSwitchingRuleDefault) WithPayload(payload *models.Error) *GetBackendSwitchingRuleDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get backend switching rule default response
func (o *GetBackendSwitchingRuleDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetBackendSwitchingRuleDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
