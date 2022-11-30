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

package backend_switching_rule

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// NewGetBackendSwitchingRuleParams creates a new GetBackendSwitchingRuleParams object
//
// There are no default values defined in the spec.
func NewGetBackendSwitchingRuleParams() GetBackendSwitchingRuleParams {

	return GetBackendSwitchingRuleParams{}
}

// GetBackendSwitchingRuleParams contains all the bound params for the get backend switching rule operation
// typically these are obtained from a http.Request
//
// swagger:parameters getBackendSwitchingRule
type GetBackendSwitchingRuleParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*Frontend name
	  Required: true
	  In: query
	*/
	Frontend string
	/*Switching Rule Index
	  Required: true
	  In: path
	*/
	Index int64
	/*ID of the transaction where we want to add the operation. Cannot be used when version is specified.
	  In: query
	*/
	TransactionID *string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewGetBackendSwitchingRuleParams() beforehand.
func (o *GetBackendSwitchingRuleParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	qFrontend, qhkFrontend, _ := qs.GetOK("frontend")
	if err := o.bindFrontend(qFrontend, qhkFrontend, route.Formats); err != nil {
		res = append(res, err)
	}

	rIndex, rhkIndex, _ := route.Params.GetOK("index")
	if err := o.bindIndex(rIndex, rhkIndex, route.Formats); err != nil {
		res = append(res, err)
	}

	qTransactionID, qhkTransactionID, _ := qs.GetOK("transaction_id")
	if err := o.bindTransactionID(qTransactionID, qhkTransactionID, route.Formats); err != nil {
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindFrontend binds and validates parameter Frontend from query.
func (o *GetBackendSwitchingRuleParams) bindFrontend(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("frontend", "query", rawData)
	}
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// AllowEmptyValue: false

	if err := validate.RequiredString("frontend", "query", raw); err != nil {
		return err
	}
	o.Frontend = raw

	return nil
}

// bindIndex binds and validates parameter Index from path.
func (o *GetBackendSwitchingRuleParams) bindIndex(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	value, err := swag.ConvertInt64(raw)
	if err != nil {
		return errors.InvalidType("index", "path", "int64", raw)
	}
	o.Index = value

	return nil
}

// bindTransactionID binds and validates parameter TransactionID from query.
func (o *GetBackendSwitchingRuleParams) bindTransactionID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false

	if raw == "" { // empty values pass all other validations
		return nil
	}
	o.TransactionID = &raw

	return nil
}
