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

package tcp_response_rule

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v6/models"
)

// GetAllTCPResponseRuleBackendOKCode is the HTTP code returned for type GetAllTCPResponseRuleBackendOK
const GetAllTCPResponseRuleBackendOKCode int = 200

/*
GetAllTCPResponseRuleBackendOK Successful operation

swagger:response getAllTcpResponseRuleBackendOK
*/
type GetAllTCPResponseRuleBackendOK struct {
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload models.TCPResponseRules `json:"body,omitempty"`
}

// NewGetAllTCPResponseRuleBackendOK creates GetAllTCPResponseRuleBackendOK with default headers values
func NewGetAllTCPResponseRuleBackendOK() *GetAllTCPResponseRuleBackendOK {

	return &GetAllTCPResponseRuleBackendOK{}
}

// WithConfigurationVersion adds the configurationVersion to the get all Tcp response rule backend o k response
func (o *GetAllTCPResponseRuleBackendOK) WithConfigurationVersion(configurationVersion string) *GetAllTCPResponseRuleBackendOK {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get all Tcp response rule backend o k response
func (o *GetAllTCPResponseRuleBackendOK) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get all Tcp response rule backend o k response
func (o *GetAllTCPResponseRuleBackendOK) WithPayload(payload models.TCPResponseRules) *GetAllTCPResponseRuleBackendOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get all Tcp response rule backend o k response
func (o *GetAllTCPResponseRuleBackendOK) SetPayload(payload models.TCPResponseRules) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetAllTCPResponseRuleBackendOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	// response header Configuration-Version

	configurationVersion := o.ConfigurationVersion
	if configurationVersion != "" {
		rw.Header().Set("Configuration-Version", configurationVersion)
	}

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = models.TCPResponseRules{}
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

/*
GetAllTCPResponseRuleBackendDefault General Error

swagger:response getAllTcpResponseRuleBackendDefault
*/
type GetAllTCPResponseRuleBackendDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetAllTCPResponseRuleBackendDefault creates GetAllTCPResponseRuleBackendDefault with default headers values
func NewGetAllTCPResponseRuleBackendDefault(code int) *GetAllTCPResponseRuleBackendDefault {
	if code <= 0 {
		code = 500
	}

	return &GetAllTCPResponseRuleBackendDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get all TCP response rule backend default response
func (o *GetAllTCPResponseRuleBackendDefault) WithStatusCode(code int) *GetAllTCPResponseRuleBackendDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get all TCP response rule backend default response
func (o *GetAllTCPResponseRuleBackendDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the get all TCP response rule backend default response
func (o *GetAllTCPResponseRuleBackendDefault) WithConfigurationVersion(configurationVersion string) *GetAllTCPResponseRuleBackendDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get all TCP response rule backend default response
func (o *GetAllTCPResponseRuleBackendDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get all TCP response rule backend default response
func (o *GetAllTCPResponseRuleBackendDefault) WithPayload(payload *models.Error) *GetAllTCPResponseRuleBackendDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get all TCP response rule backend default response
func (o *GetAllTCPResponseRuleBackendDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetAllTCPResponseRuleBackendDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
