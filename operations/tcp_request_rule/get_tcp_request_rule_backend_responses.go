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

package tcp_request_rule

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v6/models"
)

// GetTCPRequestRuleBackendOKCode is the HTTP code returned for type GetTCPRequestRuleBackendOK
const GetTCPRequestRuleBackendOKCode int = 200

/*
GetTCPRequestRuleBackendOK Successful operation

swagger:response getTcpRequestRuleBackendOK
*/
type GetTCPRequestRuleBackendOK struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.TCPRequestRule `json:"body,omitempty"`
}

// NewGetTCPRequestRuleBackendOK creates GetTCPRequestRuleBackendOK with default headers values
func NewGetTCPRequestRuleBackendOK() *GetTCPRequestRuleBackendOK {

	return &GetTCPRequestRuleBackendOK{}
}

// WithConfigurationVersion adds the configurationVersion to the get Tcp request rule backend o k response
func (o *GetTCPRequestRuleBackendOK) WithConfigurationVersion(configurationVersion string) *GetTCPRequestRuleBackendOK {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get Tcp request rule backend o k response
func (o *GetTCPRequestRuleBackendOK) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get Tcp request rule backend o k response
func (o *GetTCPRequestRuleBackendOK) WithPayload(payload *models.TCPRequestRule) *GetTCPRequestRuleBackendOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get Tcp request rule backend o k response
func (o *GetTCPRequestRuleBackendOK) SetPayload(payload *models.TCPRequestRule) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetTCPRequestRuleBackendOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// GetTCPRequestRuleBackendNotFoundCode is the HTTP code returned for type GetTCPRequestRuleBackendNotFound
const GetTCPRequestRuleBackendNotFoundCode int = 404

/*
GetTCPRequestRuleBackendNotFound The specified resource was not found

swagger:response getTcpRequestRuleBackendNotFound
*/
type GetTCPRequestRuleBackendNotFound struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetTCPRequestRuleBackendNotFound creates GetTCPRequestRuleBackendNotFound with default headers values
func NewGetTCPRequestRuleBackendNotFound() *GetTCPRequestRuleBackendNotFound {

	return &GetTCPRequestRuleBackendNotFound{}
}

// WithConfigurationVersion adds the configurationVersion to the get Tcp request rule backend not found response
func (o *GetTCPRequestRuleBackendNotFound) WithConfigurationVersion(configurationVersion string) *GetTCPRequestRuleBackendNotFound {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get Tcp request rule backend not found response
func (o *GetTCPRequestRuleBackendNotFound) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get Tcp request rule backend not found response
func (o *GetTCPRequestRuleBackendNotFound) WithPayload(payload *models.Error) *GetTCPRequestRuleBackendNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get Tcp request rule backend not found response
func (o *GetTCPRequestRuleBackendNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetTCPRequestRuleBackendNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
GetTCPRequestRuleBackendDefault General Error

swagger:response getTcpRequestRuleBackendDefault
*/
type GetTCPRequestRuleBackendDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetTCPRequestRuleBackendDefault creates GetTCPRequestRuleBackendDefault with default headers values
func NewGetTCPRequestRuleBackendDefault(code int) *GetTCPRequestRuleBackendDefault {
	if code <= 0 {
		code = 500
	}

	return &GetTCPRequestRuleBackendDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get TCP request rule backend default response
func (o *GetTCPRequestRuleBackendDefault) WithStatusCode(code int) *GetTCPRequestRuleBackendDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get TCP request rule backend default response
func (o *GetTCPRequestRuleBackendDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the get TCP request rule backend default response
func (o *GetTCPRequestRuleBackendDefault) WithConfigurationVersion(configurationVersion string) *GetTCPRequestRuleBackendDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get TCP request rule backend default response
func (o *GetTCPRequestRuleBackendDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get TCP request rule backend default response
func (o *GetTCPRequestRuleBackendDefault) WithPayload(payload *models.Error) *GetTCPRequestRuleBackendDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get TCP request rule backend default response
func (o *GetTCPRequestRuleBackendDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetTCPRequestRuleBackendDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
