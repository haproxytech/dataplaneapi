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

package specification_openapiv3

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v3/models"
)

// GetOpenapiv3SpecificationOKCode is the HTTP code returned for type GetOpenapiv3SpecificationOK
const GetOpenapiv3SpecificationOKCode int = 200

/*
GetOpenapiv3SpecificationOK Success

swagger:response getOpenapiv3SpecificationOK
*/
type GetOpenapiv3SpecificationOK struct {

	/*
	  In: Body
	*/
	Payload interface{} `json:"body,omitempty"`
}

// NewGetOpenapiv3SpecificationOK creates GetOpenapiv3SpecificationOK with default headers values
func NewGetOpenapiv3SpecificationOK() *GetOpenapiv3SpecificationOK {

	return &GetOpenapiv3SpecificationOK{}
}

// WithPayload adds the payload to the get openapiv3 specification o k response
func (o *GetOpenapiv3SpecificationOK) WithPayload(payload interface{}) *GetOpenapiv3SpecificationOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get openapiv3 specification o k response
func (o *GetOpenapiv3SpecificationOK) SetPayload(payload interface{}) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetOpenapiv3SpecificationOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

/*
GetOpenapiv3SpecificationDefault General Error

swagger:response getOpenapiv3SpecificationDefault
*/
type GetOpenapiv3SpecificationDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetOpenapiv3SpecificationDefault creates GetOpenapiv3SpecificationDefault with default headers values
func NewGetOpenapiv3SpecificationDefault(code int) *GetOpenapiv3SpecificationDefault {
	if code <= 0 {
		code = 500
	}

	return &GetOpenapiv3SpecificationDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get openapiv3 specification default response
func (o *GetOpenapiv3SpecificationDefault) WithStatusCode(code int) *GetOpenapiv3SpecificationDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get openapiv3 specification default response
func (o *GetOpenapiv3SpecificationDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the get openapiv3 specification default response
func (o *GetOpenapiv3SpecificationDefault) WithConfigurationVersion(configurationVersion string) *GetOpenapiv3SpecificationDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get openapiv3 specification default response
func (o *GetOpenapiv3SpecificationDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get openapiv3 specification default response
func (o *GetOpenapiv3SpecificationDefault) WithPayload(payload *models.Error) *GetOpenapiv3SpecificationDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get openapiv3 specification default response
func (o *GetOpenapiv3SpecificationDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetOpenapiv3SpecificationDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
