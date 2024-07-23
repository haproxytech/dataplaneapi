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

package bind

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v6/models"
)

// GetAllBindLogForwardOKCode is the HTTP code returned for type GetAllBindLogForwardOK
const GetAllBindLogForwardOKCode int = 200

/*
GetAllBindLogForwardOK Successful operation

swagger:response getAllBindLogForwardOK
*/
type GetAllBindLogForwardOK struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload models.Binds `json:"body,omitempty"`
}

// NewGetAllBindLogForwardOK creates GetAllBindLogForwardOK with default headers values
func NewGetAllBindLogForwardOK() *GetAllBindLogForwardOK {

	return &GetAllBindLogForwardOK{}
}

// WithConfigurationVersion adds the configurationVersion to the get all bind log forward o k response
func (o *GetAllBindLogForwardOK) WithConfigurationVersion(configurationVersion string) *GetAllBindLogForwardOK {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get all bind log forward o k response
func (o *GetAllBindLogForwardOK) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get all bind log forward o k response
func (o *GetAllBindLogForwardOK) WithPayload(payload models.Binds) *GetAllBindLogForwardOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get all bind log forward o k response
func (o *GetAllBindLogForwardOK) SetPayload(payload models.Binds) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetAllBindLogForwardOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	// response header Configuration-Version

	configurationVersion := o.ConfigurationVersion
	if configurationVersion != "" {
		rw.Header().Set("Configuration-Version", configurationVersion)
	}

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = models.Binds{}
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

/*
GetAllBindLogForwardDefault General Error

swagger:response getAllBindLogForwardDefault
*/
type GetAllBindLogForwardDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetAllBindLogForwardDefault creates GetAllBindLogForwardDefault with default headers values
func NewGetAllBindLogForwardDefault(code int) *GetAllBindLogForwardDefault {
	if code <= 0 {
		code = 500
	}

	return &GetAllBindLogForwardDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get all bind log forward default response
func (o *GetAllBindLogForwardDefault) WithStatusCode(code int) *GetAllBindLogForwardDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get all bind log forward default response
func (o *GetAllBindLogForwardDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the get all bind log forward default response
func (o *GetAllBindLogForwardDefault) WithConfigurationVersion(configurationVersion string) *GetAllBindLogForwardDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get all bind log forward default response
func (o *GetAllBindLogForwardDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get all bind log forward default response
func (o *GetAllBindLogForwardDefault) WithPayload(payload *models.Error) *GetAllBindLogForwardDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get all bind log forward default response
func (o *GetAllBindLogForwardDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetAllBindLogForwardDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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