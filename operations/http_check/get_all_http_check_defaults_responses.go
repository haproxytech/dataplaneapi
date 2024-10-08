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

package http_check

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v6/models"
)

// GetAllHTTPCheckDefaultsOKCode is the HTTP code returned for type GetAllHTTPCheckDefaultsOK
const GetAllHTTPCheckDefaultsOKCode int = 200

/*
GetAllHTTPCheckDefaultsOK Successful operation

swagger:response getAllHttpCheckDefaultsOK
*/
type GetAllHTTPCheckDefaultsOK struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload models.HTTPChecks `json:"body,omitempty"`
}

// NewGetAllHTTPCheckDefaultsOK creates GetAllHTTPCheckDefaultsOK with default headers values
func NewGetAllHTTPCheckDefaultsOK() *GetAllHTTPCheckDefaultsOK {

	return &GetAllHTTPCheckDefaultsOK{}
}

// WithConfigurationVersion adds the configurationVersion to the get all Http check defaults o k response
func (o *GetAllHTTPCheckDefaultsOK) WithConfigurationVersion(configurationVersion string) *GetAllHTTPCheckDefaultsOK {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get all Http check defaults o k response
func (o *GetAllHTTPCheckDefaultsOK) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get all Http check defaults o k response
func (o *GetAllHTTPCheckDefaultsOK) WithPayload(payload models.HTTPChecks) *GetAllHTTPCheckDefaultsOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get all Http check defaults o k response
func (o *GetAllHTTPCheckDefaultsOK) SetPayload(payload models.HTTPChecks) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetAllHTTPCheckDefaultsOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	// response header Configuration-Version

	configurationVersion := o.ConfigurationVersion
	if configurationVersion != "" {
		rw.Header().Set("Configuration-Version", configurationVersion)
	}

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = models.HTTPChecks{}
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

/*
GetAllHTTPCheckDefaultsDefault General Error

swagger:response getAllHttpCheckDefaultsDefault
*/
type GetAllHTTPCheckDefaultsDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetAllHTTPCheckDefaultsDefault creates GetAllHTTPCheckDefaultsDefault with default headers values
func NewGetAllHTTPCheckDefaultsDefault(code int) *GetAllHTTPCheckDefaultsDefault {
	if code <= 0 {
		code = 500
	}

	return &GetAllHTTPCheckDefaultsDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get all HTTP check defaults default response
func (o *GetAllHTTPCheckDefaultsDefault) WithStatusCode(code int) *GetAllHTTPCheckDefaultsDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get all HTTP check defaults default response
func (o *GetAllHTTPCheckDefaultsDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the get all HTTP check defaults default response
func (o *GetAllHTTPCheckDefaultsDefault) WithConfigurationVersion(configurationVersion string) *GetAllHTTPCheckDefaultsDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get all HTTP check defaults default response
func (o *GetAllHTTPCheckDefaultsDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get all HTTP check defaults default response
func (o *GetAllHTTPCheckDefaultsDefault) WithPayload(payload *models.Error) *GetAllHTTPCheckDefaultsDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get all HTTP check defaults default response
func (o *GetAllHTTPCheckDefaultsDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetAllHTTPCheckDefaultsDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
