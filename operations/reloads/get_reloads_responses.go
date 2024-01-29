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

package reloads

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v6/models"
)

// GetReloadsOKCode is the HTTP code returned for type GetReloadsOK
const GetReloadsOKCode int = 200

/*
GetReloadsOK Success

swagger:response getReloadsOK
*/
type GetReloadsOK struct {

	/*
	  In: Body
	*/
	Payload models.Reloads `json:"body,omitempty"`
}

// NewGetReloadsOK creates GetReloadsOK with default headers values
func NewGetReloadsOK() *GetReloadsOK {

	return &GetReloadsOK{}
}

// WithPayload adds the payload to the get reloads o k response
func (o *GetReloadsOK) WithPayload(payload models.Reloads) *GetReloadsOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get reloads o k response
func (o *GetReloadsOK) SetPayload(payload models.Reloads) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetReloadsOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = models.Reloads{}
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

/*
GetReloadsDefault General Error

swagger:response getReloadsDefault
*/
type GetReloadsDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetReloadsDefault creates GetReloadsDefault with default headers values
func NewGetReloadsDefault(code int) *GetReloadsDefault {
	if code <= 0 {
		code = 500
	}

	return &GetReloadsDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get reloads default response
func (o *GetReloadsDefault) WithStatusCode(code int) *GetReloadsDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get reloads default response
func (o *GetReloadsDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the get reloads default response
func (o *GetReloadsDefault) WithConfigurationVersion(configurationVersion string) *GetReloadsDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get reloads default response
func (o *GetReloadsDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get reloads default response
func (o *GetReloadsDefault) WithPayload(payload *models.Error) *GetReloadsDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get reloads default response
func (o *GetReloadsDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetReloadsDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
