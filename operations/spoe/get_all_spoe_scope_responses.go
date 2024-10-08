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

package spoe

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v6/models"
)

// GetAllSpoeScopeOKCode is the HTTP code returned for type GetAllSpoeScopeOK
const GetAllSpoeScopeOKCode int = 200

/*
GetAllSpoeScopeOK Successful operation

swagger:response getAllSpoeScopeOK
*/
type GetAllSpoeScopeOK struct {
	/*Spoe configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload models.SpoeScopes `json:"body,omitempty"`
}

// NewGetAllSpoeScopeOK creates GetAllSpoeScopeOK with default headers values
func NewGetAllSpoeScopeOK() *GetAllSpoeScopeOK {

	return &GetAllSpoeScopeOK{}
}

// WithConfigurationVersion adds the configurationVersion to the get all spoe scope o k response
func (o *GetAllSpoeScopeOK) WithConfigurationVersion(configurationVersion string) *GetAllSpoeScopeOK {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get all spoe scope o k response
func (o *GetAllSpoeScopeOK) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get all spoe scope o k response
func (o *GetAllSpoeScopeOK) WithPayload(payload models.SpoeScopes) *GetAllSpoeScopeOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get all spoe scope o k response
func (o *GetAllSpoeScopeOK) SetPayload(payload models.SpoeScopes) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetAllSpoeScopeOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	// response header Configuration-Version

	configurationVersion := o.ConfigurationVersion
	if configurationVersion != "" {
		rw.Header().Set("Configuration-Version", configurationVersion)
	}

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = models.SpoeScopes{}
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

/*
GetAllSpoeScopeDefault General Error

swagger:response getAllSpoeScopeDefault
*/
type GetAllSpoeScopeDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetAllSpoeScopeDefault creates GetAllSpoeScopeDefault with default headers values
func NewGetAllSpoeScopeDefault(code int) *GetAllSpoeScopeDefault {
	if code <= 0 {
		code = 500
	}

	return &GetAllSpoeScopeDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get all spoe scope default response
func (o *GetAllSpoeScopeDefault) WithStatusCode(code int) *GetAllSpoeScopeDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get all spoe scope default response
func (o *GetAllSpoeScopeDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the get all spoe scope default response
func (o *GetAllSpoeScopeDefault) WithConfigurationVersion(configurationVersion string) *GetAllSpoeScopeDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get all spoe scope default response
func (o *GetAllSpoeScopeDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get all spoe scope default response
func (o *GetAllSpoeScopeDefault) WithPayload(payload *models.Error) *GetAllSpoeScopeDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get all spoe scope default response
func (o *GetAllSpoeScopeDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetAllSpoeScopeDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
