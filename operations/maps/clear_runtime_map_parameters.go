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
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// NewClearRuntimeMapParams creates a new ClearRuntimeMapParams object
// with the default values initialized.
func NewClearRuntimeMapParams() ClearRuntimeMapParams {

	var (
		// initialize parameters with default values

		forceSyncDefault = bool(false)
	)

	return ClearRuntimeMapParams{
		ForceSync: &forceSyncDefault,
	}
}

// ClearRuntimeMapParams contains all the bound params for the clear runtime map operation
// typically these are obtained from a http.Request
//
// swagger:parameters clearRuntimeMap
type ClearRuntimeMapParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*If true, deletes file from disk
	  In: query
	*/
	ForceDelete *bool
	/*If true, immediately syncs changes to disk
	  In: query
	  Default: false
	*/
	ForceSync *bool
	/*Map file name
	  Required: true
	  In: path
	*/
	Name string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewClearRuntimeMapParams() beforehand.
func (o *ClearRuntimeMapParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	qForceDelete, qhkForceDelete, _ := qs.GetOK("forceDelete")
	if err := o.bindForceDelete(qForceDelete, qhkForceDelete, route.Formats); err != nil {
		res = append(res, err)
	}

	qForceSync, qhkForceSync, _ := qs.GetOK("force_sync")
	if err := o.bindForceSync(qForceSync, qhkForceSync, route.Formats); err != nil {
		res = append(res, err)
	}

	rName, rhkName, _ := route.Params.GetOK("name")
	if err := o.bindName(rName, rhkName, route.Formats); err != nil {
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindForceDelete binds and validates parameter ForceDelete from query.
func (o *ClearRuntimeMapParams) bindForceDelete(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false

	if raw == "" { // empty values pass all other validations
		return nil
	}

	value, err := swag.ConvertBool(raw)
	if err != nil {
		return errors.InvalidType("forceDelete", "query", "bool", raw)
	}
	o.ForceDelete = &value

	return nil
}

// bindForceSync binds and validates parameter ForceSync from query.
func (o *ClearRuntimeMapParams) bindForceSync(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false

	if raw == "" { // empty values pass all other validations
		// Default values have been previously initialized by NewClearRuntimeMapParams()
		return nil
	}

	value, err := swag.ConvertBool(raw)
	if err != nil {
		return errors.InvalidType("force_sync", "query", "bool", raw)
	}
	o.ForceSync = &value

	return nil
}

// bindName binds and validates parameter Name from path.
func (o *ClearRuntimeMapParams) bindName(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route
	o.Name = raw

	return nil
}
