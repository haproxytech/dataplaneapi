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

package maps

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"io"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"

	"github.com/haproxytech/client-native/v4/models"
)

// NewAddMapEntryParams creates a new AddMapEntryParams object
// with the default values initialized.
func NewAddMapEntryParams() AddMapEntryParams {

	var (
		// initialize parameters with default values

		forceSyncDefault = bool(false)
	)

	return AddMapEntryParams{
		ForceSync: &forceSyncDefault,
	}
}

// AddMapEntryParams contains all the bound params for the add map entry operation
// typically these are obtained from a http.Request
//
// swagger:parameters addMapEntry
type AddMapEntryParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  Required: true
	  In: body
	*/
	Data *models.MapEntry
	/*If true, immediately syncs changes to disk
	  In: query
	  Default: false
	*/
	ForceSync *bool
	/*Mapfile attribute storage_name
	  Required: true
	  In: query
	*/
	Map string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewAddMapEntryParams() beforehand.
func (o *AddMapEntryParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	if runtime.HasBody(r) {
		defer r.Body.Close()
		var body models.MapEntry
		if err := route.Consumer.Consume(r.Body, &body); err != nil {
			if err == io.EOF {
				res = append(res, errors.Required("data", "body", ""))
			} else {
				res = append(res, errors.NewParseError("data", "body", "", err))
			}
		} else {
			// validate body object
			if err := body.Validate(route.Formats); err != nil {
				res = append(res, err)
			}

			if len(res) == 0 {
				o.Data = &body
			}
		}
	} else {
		res = append(res, errors.Required("data", "body", ""))
	}

	qForceSync, qhkForceSync, _ := qs.GetOK("force_sync")
	if err := o.bindForceSync(qForceSync, qhkForceSync, route.Formats); err != nil {
		res = append(res, err)
	}

	qMap, qhkMap, _ := qs.GetOK("map")
	if err := o.bindMap(qMap, qhkMap, route.Formats); err != nil {
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindForceSync binds and validates parameter ForceSync from query.
func (o *AddMapEntryParams) bindForceSync(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false

	if raw == "" { // empty values pass all other validations
		// Default values have been previously initialized by NewAddMapEntryParams()
		return nil
	}

	value, err := swag.ConvertBool(raw)
	if err != nil {
		return errors.InvalidType("force_sync", "query", "bool", raw)
	}
	o.ForceSync = &value

	return nil
}

// bindMap binds and validates parameter Map from query.
func (o *AddMapEntryParams) bindMap(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("map", "query", rawData)
	}
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// AllowEmptyValue: false

	if err := validate.RequiredString("map", "query", raw); err != nil {
		return err
	}
	o.Map = raw

	return nil
}
