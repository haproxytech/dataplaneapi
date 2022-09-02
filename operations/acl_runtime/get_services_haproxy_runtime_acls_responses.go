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

package acl_runtime

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/haproxytech/client-native/v4/models"
)

// GetServicesHaproxyRuntimeAclsOKCode is the HTTP code returned for type GetServicesHaproxyRuntimeAclsOK
const GetServicesHaproxyRuntimeAclsOKCode int = 200

/*
GetServicesHaproxyRuntimeAclsOK Successful operation

swagger:response getServicesHaproxyRuntimeAclsOK
*/
type GetServicesHaproxyRuntimeAclsOK struct {

	/*
	  In: Body
	*/
	Payload models.ACLFiles `json:"body,omitempty"`
}

// NewGetServicesHaproxyRuntimeAclsOK creates GetServicesHaproxyRuntimeAclsOK with default headers values
func NewGetServicesHaproxyRuntimeAclsOK() *GetServicesHaproxyRuntimeAclsOK {

	return &GetServicesHaproxyRuntimeAclsOK{}
}

// WithPayload adds the payload to the get services haproxy runtime acls o k response
func (o *GetServicesHaproxyRuntimeAclsOK) WithPayload(payload models.ACLFiles) *GetServicesHaproxyRuntimeAclsOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get services haproxy runtime acls o k response
func (o *GetServicesHaproxyRuntimeAclsOK) SetPayload(payload models.ACLFiles) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetServicesHaproxyRuntimeAclsOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = models.ACLFiles{}
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

/*
GetServicesHaproxyRuntimeAclsDefault General Error

swagger:response getServicesHaproxyRuntimeAclsDefault
*/
type GetServicesHaproxyRuntimeAclsDefault struct {
	_statusCode int
	/*Configuration file version

	 */
	ConfigurationVersion string `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetServicesHaproxyRuntimeAclsDefault creates GetServicesHaproxyRuntimeAclsDefault with default headers values
func NewGetServicesHaproxyRuntimeAclsDefault(code int) *GetServicesHaproxyRuntimeAclsDefault {
	if code <= 0 {
		code = 500
	}

	return &GetServicesHaproxyRuntimeAclsDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get services haproxy runtime acls default response
func (o *GetServicesHaproxyRuntimeAclsDefault) WithStatusCode(code int) *GetServicesHaproxyRuntimeAclsDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get services haproxy runtime acls default response
func (o *GetServicesHaproxyRuntimeAclsDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the get services haproxy runtime acls default response
func (o *GetServicesHaproxyRuntimeAclsDefault) WithConfigurationVersion(configurationVersion string) *GetServicesHaproxyRuntimeAclsDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the get services haproxy runtime acls default response
func (o *GetServicesHaproxyRuntimeAclsDefault) SetConfigurationVersion(configurationVersion string) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the get services haproxy runtime acls default response
func (o *GetServicesHaproxyRuntimeAclsDefault) WithPayload(payload *models.Error) *GetServicesHaproxyRuntimeAclsDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get services haproxy runtime acls default response
func (o *GetServicesHaproxyRuntimeAclsDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetServicesHaproxyRuntimeAclsDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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
