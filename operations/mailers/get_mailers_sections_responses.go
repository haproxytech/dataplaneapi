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

package mailers

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v6/models"
)

// GetMailersSectionsOKCode is the HTTP code returned for type GetMailersSectionsOK
const GetMailersSectionsOKCode int = 200

/*
GetMailersSectionsOK Successful operation

swagger:response getMailersSectionsOK
*/
type GetMailersSectionsOK struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *GetMailersSectionsOKBody `json:"body,omitempty"`
}

// NewGetMailersSectionsOK creates GetMailersSectionsOK with default headers values
func NewGetMailersSectionsOK() *GetMailersSectionsOK {

	return &GetMailersSectionsOK{}
}

// WithConfigurationVersion adds the configurationVersion to the get mailers sections o k response
func (o *GetMailersSectionsOK) WithConfigurationVersion(configurationVersion string) *GetMailersSectionsOK {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get mailers sections o k response
func (o *GetMailersSectionsOK) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get mailers sections o k response
func (o *GetMailersSectionsOK) WithPayload(payload *GetMailersSectionsOKBody) *GetMailersSectionsOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get mailers sections o k response
func (o *GetMailersSectionsOK) SetPayload(payload *GetMailersSectionsOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetMailersSectionsOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
GetMailersSectionsDefault General Error

swagger:response getMailersSectionsDefault
*/
type GetMailersSectionsDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetMailersSectionsDefault creates GetMailersSectionsDefault with default headers values
func NewGetMailersSectionsDefault(code int) *GetMailersSectionsDefault {
	if code <= 0 {
		code = 500
	}

	return &GetMailersSectionsDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get mailers sections default response
func (o *GetMailersSectionsDefault) WithStatusCode(code int) *GetMailersSectionsDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get mailers sections default response
func (o *GetMailersSectionsDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the get mailers sections default response
func (o *GetMailersSectionsDefault) WithConfigurationVersion(configurationVersion string) *GetMailersSectionsDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get mailers sections default response
func (o *GetMailersSectionsDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get mailers sections default response
func (o *GetMailersSectionsDefault) WithPayload(payload *models.Error) *GetMailersSectionsDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get mailers sections default response
func (o *GetMailersSectionsDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetMailersSectionsDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
