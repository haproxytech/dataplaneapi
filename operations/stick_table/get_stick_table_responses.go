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

package stick_table

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v3/models"
)

// GetStickTableOKCode is the HTTP code returned for type GetStickTableOK
const GetStickTableOKCode int = 200

/*
GetStickTableOK Successful operation

swagger:response getStickTableOK
*/
type GetStickTableOK struct {

	/*
	  In: Body
	*/
	Payload *models.StickTable `json:"body,omitempty"`
}

// NewGetStickTableOK creates GetStickTableOK with default headers values
func NewGetStickTableOK() *GetStickTableOK {

	return &GetStickTableOK{}
}

// WithPayload adds the payload to the get stick table o k response
func (o *GetStickTableOK) WithPayload(payload *models.StickTable) *GetStickTableOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get stick table o k response
func (o *GetStickTableOK) SetPayload(payload *models.StickTable) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetStickTableOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetStickTableNotFoundCode is the HTTP code returned for type GetStickTableNotFound
const GetStickTableNotFoundCode int = 404

/*
GetStickTableNotFound The specified resource was not found

swagger:response getStickTableNotFound
*/
type GetStickTableNotFound struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetStickTableNotFound creates GetStickTableNotFound with default headers values
func NewGetStickTableNotFound() *GetStickTableNotFound {

	return &GetStickTableNotFound{}
}

// WithConfigurationVersion adds the configurationVersion to the get stick table not found response
func (o *GetStickTableNotFound) WithConfigurationVersion(configurationVersion string) *GetStickTableNotFound {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get stick table not found response
func (o *GetStickTableNotFound) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get stick table not found response
func (o *GetStickTableNotFound) WithPayload(payload *models.Error) *GetStickTableNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get stick table not found response
func (o *GetStickTableNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetStickTableNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
GetStickTableDefault General Error

swagger:response getStickTableDefault
*/
type GetStickTableDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetStickTableDefault creates GetStickTableDefault with default headers values
func NewGetStickTableDefault(code int) *GetStickTableDefault {
	if code <= 0 {
		code = 500
	}

	return &GetStickTableDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get stick table default response
func (o *GetStickTableDefault) WithStatusCode(code int) *GetStickTableDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get stick table default response
func (o *GetStickTableDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the get stick table default response
func (o *GetStickTableDefault) WithConfigurationVersion(configurationVersion string) *GetStickTableDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get stick table default response
func (o *GetStickTableDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get stick table default response
func (o *GetStickTableDefault) WithPayload(payload *models.Error) *GetStickTableDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get stick table default response
func (o *GetStickTableDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetStickTableDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
