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

package s_s_l_runtime

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v6/models"
)

// GetAllCrtListsOKCode is the HTTP code returned for type GetAllCrtListsOK
const GetAllCrtListsOKCode int = 200

/*
GetAllCrtListsOK Successful operation

swagger:response getAllCrtListsOK
*/
type GetAllCrtListsOK struct {

	/*
	  In: Body
	*/
	Payload models.SslCrtLists `json:"body,omitempty"`
}

// NewGetAllCrtListsOK creates GetAllCrtListsOK with default headers values
func NewGetAllCrtListsOK() *GetAllCrtListsOK {

	return &GetAllCrtListsOK{}
}

// WithPayload adds the payload to the get all crt lists o k response
func (o *GetAllCrtListsOK) WithPayload(payload models.SslCrtLists) *GetAllCrtListsOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get all crt lists o k response
func (o *GetAllCrtListsOK) SetPayload(payload models.SslCrtLists) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetAllCrtListsOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = models.SslCrtLists{}
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

/*
GetAllCrtListsDefault General Error

swagger:response getAllCrtListsDefault
*/
type GetAllCrtListsDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetAllCrtListsDefault creates GetAllCrtListsDefault with default headers values
func NewGetAllCrtListsDefault(code int) *GetAllCrtListsDefault {
	if code <= 0 {
		code = 500
	}

	return &GetAllCrtListsDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get all crt lists default response
func (o *GetAllCrtListsDefault) WithStatusCode(code int) *GetAllCrtListsDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get all crt lists default response
func (o *GetAllCrtListsDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the get all crt lists default response
func (o *GetAllCrtListsDefault) WithConfigurationVersion(configurationVersion string) *GetAllCrtListsDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get all crt lists default response
func (o *GetAllCrtListsDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get all crt lists default response
func (o *GetAllCrtListsDefault) WithPayload(payload *models.Error) *GetAllCrtListsDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get all crt lists default response
func (o *GetAllCrtListsDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetAllCrtListsDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
