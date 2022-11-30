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
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v3/models"
)

// GetOneStorageSSLCertificateOKCode is the HTTP code returned for type GetOneStorageSSLCertificateOK
const GetOneStorageSSLCertificateOKCode int = 200

/*
GetOneStorageSSLCertificateOK Successful operation

swagger:response getOneStorageSSLCertificateOK
*/
type GetOneStorageSSLCertificateOK struct {

	/*
	  In: Body
	*/
	Payload *models.SslCertificate `json:"body,omitempty"`
}

// NewGetOneStorageSSLCertificateOK creates GetOneStorageSSLCertificateOK with default headers values
func NewGetOneStorageSSLCertificateOK() *GetOneStorageSSLCertificateOK {

	return &GetOneStorageSSLCertificateOK{}
}

// WithPayload adds the payload to the get one storage s s l certificate o k response
func (o *GetOneStorageSSLCertificateOK) WithPayload(payload *models.SslCertificate) *GetOneStorageSSLCertificateOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get one storage s s l certificate o k response
func (o *GetOneStorageSSLCertificateOK) SetPayload(payload *models.SslCertificate) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetOneStorageSSLCertificateOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetOneStorageSSLCertificateNotFoundCode is the HTTP code returned for type GetOneStorageSSLCertificateNotFound
const GetOneStorageSSLCertificateNotFoundCode int = 404

/*
GetOneStorageSSLCertificateNotFound The specified resource was not found

swagger:response getOneStorageSSLCertificateNotFound
*/
type GetOneStorageSSLCertificateNotFound struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetOneStorageSSLCertificateNotFound creates GetOneStorageSSLCertificateNotFound with default headers values
func NewGetOneStorageSSLCertificateNotFound() *GetOneStorageSSLCertificateNotFound {

	return &GetOneStorageSSLCertificateNotFound{}
}

// WithConfigurationVersion adds the configurationVersion to the get one storage s s l certificate not found response
func (o *GetOneStorageSSLCertificateNotFound) WithConfigurationVersion(configurationVersion string) *GetOneStorageSSLCertificateNotFound {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get one storage s s l certificate not found response
func (o *GetOneStorageSSLCertificateNotFound) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get one storage s s l certificate not found response
func (o *GetOneStorageSSLCertificateNotFound) WithPayload(payload *models.Error) *GetOneStorageSSLCertificateNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get one storage s s l certificate not found response
func (o *GetOneStorageSSLCertificateNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetOneStorageSSLCertificateNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
GetOneStorageSSLCertificateDefault General Error

swagger:response getOneStorageSSLCertificateDefault
*/
type GetOneStorageSSLCertificateDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetOneStorageSSLCertificateDefault creates GetOneStorageSSLCertificateDefault with default headers values
func NewGetOneStorageSSLCertificateDefault(code int) *GetOneStorageSSLCertificateDefault {
	if code <= 0 {
		code = 500
	}

	return &GetOneStorageSSLCertificateDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get one storage s s l certificate default response
func (o *GetOneStorageSSLCertificateDefault) WithStatusCode(code int) *GetOneStorageSSLCertificateDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get one storage s s l certificate default response
func (o *GetOneStorageSSLCertificateDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the get one storage s s l certificate default response
func (o *GetOneStorageSSLCertificateDefault) WithConfigurationVersion(configurationVersion string) *GetOneStorageSSLCertificateDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get one storage s s l certificate default response
func (o *GetOneStorageSSLCertificateDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get one storage s s l certificate default response
func (o *GetOneStorageSSLCertificateDefault) WithPayload(payload *models.Error) *GetOneStorageSSLCertificateDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get one storage s s l certificate default response
func (o *GetOneStorageSSLCertificateDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetOneStorageSSLCertificateDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
