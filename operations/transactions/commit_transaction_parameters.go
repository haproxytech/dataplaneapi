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

package transactions

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

// NewCommitTransactionParams creates a new CommitTransactionParams object
// with the default values initialized.
func NewCommitTransactionParams() CommitTransactionParams {

	var (
		// initialize parameters with default values

		forceReloadDefault   = bool(false)
		enableReverseDefault = bool(false)
	)

	return CommitTransactionParams{
		ForceReload:  &forceReloadDefault,
		EnableRevert: &enableReverseDefault,
	}
}

// CommitTransactionParams contains all the bound params for the commit transaction operation
// typically these are obtained from a http.Request
//
// swagger:parameters commitTransaction
type CommitTransactionParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*If set, do a force reload, do not wait for the configured reload-delay. Cannot be used when transaction is specified, as changes in transaction are not applied directly to configuration.
	  In: query
	  Default: false
	*/
	ForceReload *bool

	EnableRevert *bool
	/*Transaction id
	  Required: true
	  In: path
	*/
	ID string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewCommitTransactionParams() beforehand.
func (o *CommitTransactionParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	qForceReload, qhkForceReload, _ := qs.GetOK("force_reload")
	if err := o.bindForceReload(qForceReload, qhkForceReload, route.Formats); err != nil {
		res = append(res, err)
	}

	qEnableRevert, qhkEnableRevert, _ := qs.GetOK("enable_revert")
	if err := o.bindEnableRevert(qEnableRevert, qhkEnableRevert, route.Formats); err != nil {
		res = append(res, err)
	}

	rID, rhkID, _ := route.Params.GetOK("id")
	if err := o.bindID(rID, rhkID, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindForceReload binds and validates parameter ForceReload from query.
func (o *CommitTransactionParams) bindForceReload(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		// Default values have been previously initialized by NewCommitTransactionParams()
		return nil
	}

	value, err := swag.ConvertBool(raw)
	if err != nil {
		return errors.InvalidType("force_reload", "query", "bool", raw)
	}
	o.ForceReload = &value

	return nil
}

// bindEnableRevert binds and validates parameter EnableRevert from query.
func (o *CommitTransactionParams) bindEnableRevert(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		// Default values have been previously initialized by NewCommitTransactionParams()
		return nil
	}

	value, err := swag.ConvertBool(raw)
	if err != nil {
		return errors.InvalidType("enable_revert", "query", "bool", raw)
	}
	o.EnableRevert = &value

	return nil
}

// bindID binds and validates parameter ID from path.
func (o *CommitTransactionParams) bindID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.ID = raw

	return nil
}
