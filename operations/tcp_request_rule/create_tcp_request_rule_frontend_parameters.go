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
	"io"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/haproxytech/client-native/v6/models"
)

// NewCreateTCPRequestRuleFrontendParams creates a new CreateTCPRequestRuleFrontendParams object
// with the default values initialized.
func NewCreateTCPRequestRuleFrontendParams() CreateTCPRequestRuleFrontendParams {

	var (
		// initialize parameters with default values

		forceReloadDefault = bool(false)
	)

	return CreateTCPRequestRuleFrontendParams{
		ForceReload: &forceReloadDefault,
	}
}

// CreateTCPRequestRuleFrontendParams contains all the bound params for the create TCP request rule frontend operation
// typically these are obtained from a http.Request
//
// swagger:parameters createTCPRequestRuleFrontend
type CreateTCPRequestRuleFrontendParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  Required: true
	  In: body
	*/
	Data *models.TCPRequestRule
	/*If set, do a force reload, do not wait for the configured reload-delay. Cannot be used when transaction is specified, as changes in transaction are not applied directly to configuration.
	  In: query
	  Default: false
	*/
	ForceReload *bool
	/*TCP Request Rule Index
	  Required: true
	  In: path
	*/
	Index int64
	/*Parent name
	  Required: true
	  In: path
	*/
	ParentName string
	/*ID of the transaction where we want to add the operation. Cannot be used when version is specified.
	  In: query
	*/
	TransactionID *string
	/*Version used for checking configuration version. Cannot be used when transaction is specified, transaction has it's own version.
	  In: query
	*/
	Version *int64
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewCreateTCPRequestRuleFrontendParams() beforehand.
func (o *CreateTCPRequestRuleFrontendParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	if runtime.HasBody(r) {
		defer r.Body.Close()
		var body models.TCPRequestRule
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

	qForceReload, qhkForceReload, _ := qs.GetOK("force_reload")
	if err := o.bindForceReload(qForceReload, qhkForceReload, route.Formats); err != nil {
		res = append(res, err)
	}

	rIndex, rhkIndex, _ := route.Params.GetOK("index")
	if err := o.bindIndex(rIndex, rhkIndex, route.Formats); err != nil {
		res = append(res, err)
	}

	rParentName, rhkParentName, _ := route.Params.GetOK("parent_name")
	if err := o.bindParentName(rParentName, rhkParentName, route.Formats); err != nil {
		res = append(res, err)
	}

	qTransactionID, qhkTransactionID, _ := qs.GetOK("transaction_id")
	if err := o.bindTransactionID(qTransactionID, qhkTransactionID, route.Formats); err != nil {
		res = append(res, err)
	}

	qVersion, qhkVersion, _ := qs.GetOK("version")
	if err := o.bindVersion(qVersion, qhkVersion, route.Formats); err != nil {
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindForceReload binds and validates parameter ForceReload from query.
func (o *CreateTCPRequestRuleFrontendParams) bindForceReload(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false

	if raw == "" { // empty values pass all other validations
		// Default values have been previously initialized by NewCreateTCPRequestRuleFrontendParams()
		return nil
	}

	value, err := swag.ConvertBool(raw)
	if err != nil {
		return errors.InvalidType("force_reload", "query", "bool", raw)
	}
	o.ForceReload = &value

	return nil
}

// bindIndex binds and validates parameter Index from path.
func (o *CreateTCPRequestRuleFrontendParams) bindIndex(rawData []string, hasKey bool, formats strfmt.Registry) error {
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

// bindParentName binds and validates parameter ParentName from path.
func (o *CreateTCPRequestRuleFrontendParams) bindParentName(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route
	o.ParentName = raw

	return nil
}

// bindTransactionID binds and validates parameter TransactionID from query.
func (o *CreateTCPRequestRuleFrontendParams) bindTransactionID(rawData []string, hasKey bool, formats strfmt.Registry) error {
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

// bindVersion binds and validates parameter Version from query.
func (o *CreateTCPRequestRuleFrontendParams) bindVersion(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false

	if raw == "" { // empty values pass all other validations
		return nil
	}

	value, err := swag.ConvertInt64(raw)
	if err != nil {
		return errors.InvalidType("version", "query", "int64", raw)
	}
	o.Version = &value

	return nil
}
