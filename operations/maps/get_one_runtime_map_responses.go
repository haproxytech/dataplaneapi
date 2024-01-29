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

package maps

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v6/models"
)

// GetOneRuntimeMapOKCode is the HTTP code returned for type GetOneRuntimeMapOK
const GetOneRuntimeMapOKCode int = 200

/*
GetOneRuntimeMapOK Successful operation

swagger:response getOneRuntimeMapOK
*/
type GetOneRuntimeMapOK struct {

	/*
	  In: Body
	*/
	Payload *models.Map `json:"body,omitempty"`
}

// NewGetOneRuntimeMapOK creates GetOneRuntimeMapOK with default headers values
func NewGetOneRuntimeMapOK() *GetOneRuntimeMapOK {

	return &GetOneRuntimeMapOK{}
}

// WithPayload adds the payload to the get one runtime map o k response
func (o *GetOneRuntimeMapOK) WithPayload(payload *models.Map) *GetOneRuntimeMapOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get one runtime map o k response
func (o *GetOneRuntimeMapOK) SetPayload(payload *models.Map) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetOneRuntimeMapOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetOneRuntimeMapNotFoundCode is the HTTP code returned for type GetOneRuntimeMapNotFound
const GetOneRuntimeMapNotFoundCode int = 404

/*
GetOneRuntimeMapNotFound The specified resource was not found

swagger:response getOneRuntimeMapNotFound
*/
type GetOneRuntimeMapNotFound struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetOneRuntimeMapNotFound creates GetOneRuntimeMapNotFound with default headers values
func NewGetOneRuntimeMapNotFound() *GetOneRuntimeMapNotFound {

	return &GetOneRuntimeMapNotFound{}
}

// WithConfigurationVersion adds the configurationVersion to the get one runtime map not found response
func (o *GetOneRuntimeMapNotFound) WithConfigurationVersion(configurationVersion string) *GetOneRuntimeMapNotFound {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get one runtime map not found response
func (o *GetOneRuntimeMapNotFound) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get one runtime map not found response
func (o *GetOneRuntimeMapNotFound) WithPayload(payload *models.Error) *GetOneRuntimeMapNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get one runtime map not found response
func (o *GetOneRuntimeMapNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetOneRuntimeMapNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
GetOneRuntimeMapDefault General Error

swagger:response getOneRuntimeMapDefault
*/
type GetOneRuntimeMapDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetOneRuntimeMapDefault creates GetOneRuntimeMapDefault with default headers values
func NewGetOneRuntimeMapDefault(code int) *GetOneRuntimeMapDefault {
	if code <= 0 {
		code = 500
	}

	return &GetOneRuntimeMapDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get one runtime map default response
func (o *GetOneRuntimeMapDefault) WithStatusCode(code int) *GetOneRuntimeMapDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get one runtime map default response
func (o *GetOneRuntimeMapDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the get one runtime map default response
func (o *GetOneRuntimeMapDefault) WithConfigurationVersion(configurationVersion string) *GetOneRuntimeMapDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get one runtime map default response
func (o *GetOneRuntimeMapDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get one runtime map default response
func (o *GetOneRuntimeMapDefault) WithPayload(payload *models.Error) *GetOneRuntimeMapDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get one runtime map default response
func (o *GetOneRuntimeMapDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetOneRuntimeMapDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
