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

package process_manager

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v5/models"
)

// GetProgramsOKCode is the HTTP code returned for type GetProgramsOK
const GetProgramsOKCode int = 200

/*
GetProgramsOK Successful operation

swagger:response getProgramsOK
*/
type GetProgramsOK struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *GetProgramsOKBody `json:"body,omitempty"`
}

// NewGetProgramsOK creates GetProgramsOK with default headers values
func NewGetProgramsOK() *GetProgramsOK {

	return &GetProgramsOK{}
}

// WithConfigurationVersion adds the configurationVersion to the get programs o k response
func (o *GetProgramsOK) WithConfigurationVersion(configurationVersion string) *GetProgramsOK {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get programs o k response
func (o *GetProgramsOK) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get programs o k response
func (o *GetProgramsOK) WithPayload(payload *GetProgramsOKBody) *GetProgramsOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get programs o k response
func (o *GetProgramsOK) SetPayload(payload *GetProgramsOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetProgramsOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

/*
GetProgramsDefault General Error

swagger:response getProgramsDefault
*/
type GetProgramsDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetProgramsDefault creates GetProgramsDefault with default headers values
func NewGetProgramsDefault(code int) *GetProgramsDefault {
	if code <= 0 {
		code = 500
	}

	return &GetProgramsDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get programs default response
func (o *GetProgramsDefault) WithStatusCode(code int) *GetProgramsDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get programs default response
func (o *GetProgramsDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the get programs default response
func (o *GetProgramsDefault) WithConfigurationVersion(configurationVersion string) *GetProgramsDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get programs default response
func (o *GetProgramsDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get programs default response
func (o *GetProgramsDefault) WithPayload(payload *models.Error) *GetProgramsDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get programs default response
func (o *GetProgramsDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetProgramsDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
