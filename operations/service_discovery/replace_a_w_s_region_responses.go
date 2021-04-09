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

package service_discovery

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/swag"

	"github.com/haproxytech/client-native/v2/models"
)

// ReplaceAWSRegionOKCode is the HTTP code returned for type ReplaceAWSRegionOK
const ReplaceAWSRegionOKCode int = 200

/*ReplaceAWSRegionOK Resource updated

swagger:response replaceAWSRegionOK
*/
type ReplaceAWSRegionOK struct {

	/*
	  In: Body
	*/
	Payload *models.AwsRegion `json:"body,omitempty"`
}

// NewReplaceAWSRegionOK creates ReplaceAWSRegionOK with default headers values
func NewReplaceAWSRegionOK() *ReplaceAWSRegionOK {

	return &ReplaceAWSRegionOK{}
}

// WithPayload adds the payload to the replace a w s region o k response
func (o *ReplaceAWSRegionOK) WithPayload(payload *models.AwsRegion) *ReplaceAWSRegionOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace a w s region o k response
func (o *ReplaceAWSRegionOK) SetPayload(payload *models.AwsRegion) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceAWSRegionOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// ReplaceAWSRegionBadRequestCode is the HTTP code returned for type ReplaceAWSRegionBadRequest
const ReplaceAWSRegionBadRequestCode int = 400

/*ReplaceAWSRegionBadRequest Bad request

swagger:response replaceAWSRegionBadRequest
*/
type ReplaceAWSRegionBadRequest struct {
	/*Configuration file version

	  Default: 0
	*/
	ConfigurationVersion int64 `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplaceAWSRegionBadRequest creates ReplaceAWSRegionBadRequest with default headers values
func NewReplaceAWSRegionBadRequest() *ReplaceAWSRegionBadRequest {

	var (
		// initialize headers with default values

		configurationVersionDefault = int64(0)
	)

	return &ReplaceAWSRegionBadRequest{

		ConfigurationVersion: configurationVersionDefault,
	}
}

// WithConfigurationVersion adds the configurationVersion to the replace a w s region bad request response
func (o *ReplaceAWSRegionBadRequest) WithConfigurationVersion(configurationVersion int64) *ReplaceAWSRegionBadRequest {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the replace a w s region bad request response
func (o *ReplaceAWSRegionBadRequest) SetConfigurationVersion(configurationVersion int64) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the replace a w s region bad request response
func (o *ReplaceAWSRegionBadRequest) WithPayload(payload *models.Error) *ReplaceAWSRegionBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace a w s region bad request response
func (o *ReplaceAWSRegionBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceAWSRegionBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	// response header Configuration-Version

	configurationVersion := swag.FormatInt64(o.ConfigurationVersion)
	if configurationVersion != "" {
		rw.Header().Set("Configuration-Version", configurationVersion)
	}

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// ReplaceAWSRegionNotFoundCode is the HTTP code returned for type ReplaceAWSRegionNotFound
const ReplaceAWSRegionNotFoundCode int = 404

/*ReplaceAWSRegionNotFound The specified resource was not found

swagger:response replaceAWSRegionNotFound
*/
type ReplaceAWSRegionNotFound struct {
	/*Configuration file version

	  Default: 0
	*/
	ConfigurationVersion int64 `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplaceAWSRegionNotFound creates ReplaceAWSRegionNotFound with default headers values
func NewReplaceAWSRegionNotFound() *ReplaceAWSRegionNotFound {

	var (
		// initialize headers with default values

		configurationVersionDefault = int64(0)
	)

	return &ReplaceAWSRegionNotFound{

		ConfigurationVersion: configurationVersionDefault,
	}
}

// WithConfigurationVersion adds the configurationVersion to the replace a w s region not found response
func (o *ReplaceAWSRegionNotFound) WithConfigurationVersion(configurationVersion int64) *ReplaceAWSRegionNotFound {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the replace a w s region not found response
func (o *ReplaceAWSRegionNotFound) SetConfigurationVersion(configurationVersion int64) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the replace a w s region not found response
func (o *ReplaceAWSRegionNotFound) WithPayload(payload *models.Error) *ReplaceAWSRegionNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace a w s region not found response
func (o *ReplaceAWSRegionNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceAWSRegionNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	// response header Configuration-Version

	configurationVersion := swag.FormatInt64(o.ConfigurationVersion)
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

/*ReplaceAWSRegionDefault General Error

swagger:response replaceAWSRegionDefault
*/
type ReplaceAWSRegionDefault struct {
	_statusCode int
	/*Configuration file version

	  Default: 0
	*/
	ConfigurationVersion int64 `json:"Configuration-Version"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewReplaceAWSRegionDefault creates ReplaceAWSRegionDefault with default headers values
func NewReplaceAWSRegionDefault(code int) *ReplaceAWSRegionDefault {
	if code <= 0 {
		code = 500
	}

	var (
		// initialize headers with default values

		configurationVersionDefault = int64(0)
	)

	return &ReplaceAWSRegionDefault{
		_statusCode: code,

		ConfigurationVersion: configurationVersionDefault,
	}
}

// WithStatusCode adds the status to the replace a w s region default response
func (o *ReplaceAWSRegionDefault) WithStatusCode(code int) *ReplaceAWSRegionDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the replace a w s region default response
func (o *ReplaceAWSRegionDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithConfigurationVersion adds the configurationVersion to the replace a w s region default response
func (o *ReplaceAWSRegionDefault) WithConfigurationVersion(configurationVersion int64) *ReplaceAWSRegionDefault {
	o.ConfigurationVersion = configurationVersion
	return o
}

// SetConfigurationVersion sets the configurationVersion to the replace a w s region default response
func (o *ReplaceAWSRegionDefault) SetConfigurationVersion(configurationVersion int64) {
	o.ConfigurationVersion = configurationVersion
}

// WithPayload adds the payload to the replace a w s region default response
func (o *ReplaceAWSRegionDefault) WithPayload(payload *models.Error) *ReplaceAWSRegionDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the replace a w s region default response
func (o *ReplaceAWSRegionDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReplaceAWSRegionDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	// response header Configuration-Version

	configurationVersion := swag.FormatInt64(o.ConfigurationVersion)
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