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

package discovery

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v5/models"
)

// GetStatsEndpointsOKCode is the HTTP code returned for type GetStatsEndpointsOK
const GetStatsEndpointsOKCode int = 200

/*
GetStatsEndpointsOK Success

swagger:response getStatsEndpointsOK
*/
type GetStatsEndpointsOK struct {

	/*
	  In: Body
	*/
	Payload models.Endpoints `json:"body,omitempty"`
}

// NewGetStatsEndpointsOK creates GetStatsEndpointsOK with default headers values
func NewGetStatsEndpointsOK() *GetStatsEndpointsOK {

	return &GetStatsEndpointsOK{}
}

// WithPayload adds the payload to the get stats endpoints o k response
func (o *GetStatsEndpointsOK) WithPayload(payload models.Endpoints) *GetStatsEndpointsOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get stats endpoints o k response
func (o *GetStatsEndpointsOK) SetPayload(payload models.Endpoints) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetStatsEndpointsOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = models.Endpoints{}
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

/*
GetStatsEndpointsDefault General Error

swagger:response getStatsEndpointsDefault
*/
type GetStatsEndpointsDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetStatsEndpointsDefault creates GetStatsEndpointsDefault with default headers values
func NewGetStatsEndpointsDefault(code int) *GetStatsEndpointsDefault {
	if code <= 0 {
		code = 500
	}

	return &GetStatsEndpointsDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get stats endpoints default response
func (o *GetStatsEndpointsDefault) WithStatusCode(code int) *GetStatsEndpointsDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get stats endpoints default response
func (o *GetStatsEndpointsDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the get stats endpoints default response
func (o *GetStatsEndpointsDefault) WithConfigurationVersion(configurationVersion string) *GetStatsEndpointsDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get stats endpoints default response
func (o *GetStatsEndpointsDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get stats endpoints default response
func (o *GetStatsEndpointsDefault) WithPayload(payload *models.Error) *GetStatsEndpointsDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get stats endpoints default response
func (o *GetStatsEndpointsDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetStatsEndpointsDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
