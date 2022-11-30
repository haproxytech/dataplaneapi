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

package storage

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"io"
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v3/models"
)

// GetOneStorageGeneralFileOKCode is the HTTP code returned for type GetOneStorageGeneralFileOK
const GetOneStorageGeneralFileOKCode int = 200

/*
GetOneStorageGeneralFileOK Successful operation

swagger:response getOneStorageGeneralFileOK
*/
type GetOneStorageGeneralFileOK struct {

	/*
	  In: Body
	*/
	Payload io.ReadCloser `json:"body,omitempty"`
}

// NewGetOneStorageGeneralFileOK creates GetOneStorageGeneralFileOK with default headers values
func NewGetOneStorageGeneralFileOK() *GetOneStorageGeneralFileOK {

	return &GetOneStorageGeneralFileOK{}
}

// WithPayload adds the payload to the get one storage general file o k response
func (o *GetOneStorageGeneralFileOK) WithPayload(payload io.ReadCloser) *GetOneStorageGeneralFileOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get one storage general file o k response
func (o *GetOneStorageGeneralFileOK) SetPayload(payload io.ReadCloser) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetOneStorageGeneralFileOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// GetOneStorageGeneralFileNotFoundCode is the HTTP code returned for type GetOneStorageGeneralFileNotFound
const GetOneStorageGeneralFileNotFoundCode int = 404

/*
GetOneStorageGeneralFileNotFound The specified resource was not found

swagger:response getOneStorageGeneralFileNotFound
*/
type GetOneStorageGeneralFileNotFound struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetOneStorageGeneralFileNotFound creates GetOneStorageGeneralFileNotFound with default headers values
func NewGetOneStorageGeneralFileNotFound() *GetOneStorageGeneralFileNotFound {

	return &GetOneStorageGeneralFileNotFound{}
}

// WithConfigurationVersion adds the configurationVersion to the get one storage general file not found response
func (o *GetOneStorageGeneralFileNotFound) WithConfigurationVersion(configurationVersion string) *GetOneStorageGeneralFileNotFound {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get one storage general file not found response
func (o *GetOneStorageGeneralFileNotFound) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get one storage general file not found response
func (o *GetOneStorageGeneralFileNotFound) WithPayload(payload *models.Error) *GetOneStorageGeneralFileNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get one storage general file not found response
func (o *GetOneStorageGeneralFileNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetOneStorageGeneralFileNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
GetOneStorageGeneralFileDefault General Error

swagger:response getOneStorageGeneralFileDefault
*/
type GetOneStorageGeneralFileDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetOneStorageGeneralFileDefault creates GetOneStorageGeneralFileDefault with default headers values
func NewGetOneStorageGeneralFileDefault(code int) *GetOneStorageGeneralFileDefault {
	if code <= 0 {
		code = 500
	}

	return &GetOneStorageGeneralFileDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get one storage general file default response
func (o *GetOneStorageGeneralFileDefault) WithStatusCode(code int) *GetOneStorageGeneralFileDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get one storage general file default response
func (o *GetOneStorageGeneralFileDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the get one storage general file default response
func (o *GetOneStorageGeneralFileDefault) WithConfigurationVersion(configurationVersion string) *GetOneStorageGeneralFileDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get one storage general file default response
func (o *GetOneStorageGeneralFileDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get one storage general file default response
func (o *GetOneStorageGeneralFileDefault) WithPayload(payload *models.Error) *GetOneStorageGeneralFileDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get one storage general file default response
func (o *GetOneStorageGeneralFileDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetOneStorageGeneralFileDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
