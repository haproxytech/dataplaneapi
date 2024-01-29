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

package http_response_rule

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"context"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"

	"github.com/haproxytech/client-native/v6/models"
)

// GetHTTPResponseRulesHandlerFunc turns a function with the right signature into a get HTTP response rules handler
type GetHTTPResponseRulesHandlerFunc func(GetHTTPResponseRulesParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn GetHTTPResponseRulesHandlerFunc) Handle(params GetHTTPResponseRulesParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// GetHTTPResponseRulesHandler interface for that can handle valid get HTTP response rules params
type GetHTTPResponseRulesHandler interface {
	Handle(GetHTTPResponseRulesParams, interface{}) middleware.Responder
}

// NewGetHTTPResponseRules creates a new http.Handler for the get HTTP response rules operation
func NewGetHTTPResponseRules(ctx *middleware.Context, handler GetHTTPResponseRulesHandler) *GetHTTPResponseRules {
	return &GetHTTPResponseRules{Context: ctx, Handler: handler}
}

/*
	GetHTTPResponseRules swagger:route GET /services/haproxy/configuration/http_response_rules HTTPResponseRule getHttpResponseRules

# Return an array of all HTTP Response Rules

Returns all HTTP Response Rules that are configured in specified parent.
*/
type GetHTTPResponseRules struct {
	Context *middleware.Context
	Handler GetHTTPResponseRulesHandler
}

func (o *GetHTTPResponseRules) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewGetHTTPResponseRulesParams()
	uprinc, aCtx, err := o.Context.Authorize(r, route)
	if err != nil {
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}
	if aCtx != nil {
		*r = *aCtx
	}
	var principal interface{}
	if uprinc != nil {
		principal = uprinc.(interface{}) // this is really a interface{}, I promise
	}

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params, principal) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}

// GetHTTPResponseRulesOKBody get HTTP response rules o k body
//
// swagger:model GetHTTPResponseRulesOKBody
type GetHTTPResponseRulesOKBody struct {

	// version
	Version int64 `json:"_version,omitempty"`

	// data
	// Required: true
	Data models.HTTPResponseRules `json:"data"`
}

// Validate validates this get HTTP response rules o k body
func (o *GetHTTPResponseRulesOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateData(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetHTTPResponseRulesOKBody) validateData(formats strfmt.Registry) error {

	if err := validate.Required("getHttpResponseRulesOK"+"."+"data", "body", o.Data); err != nil {
		return err
	}

	if err := o.Data.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("getHttpResponseRulesOK" + "." + "data")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("getHttpResponseRulesOK" + "." + "data")
		}
		return err
	}

	return nil
}

// ContextValidate validate this get HTTP response rules o k body based on the context it is used
func (o *GetHTTPResponseRulesOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateData(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetHTTPResponseRulesOKBody) contextValidateData(ctx context.Context, formats strfmt.Registry) error {

	if err := o.Data.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("getHttpResponseRulesOK" + "." + "data")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("getHttpResponseRulesOK" + "." + "data")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *GetHTTPResponseRulesOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetHTTPResponseRulesOKBody) UnmarshalBinary(b []byte) error {
	var res GetHTTPResponseRulesOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
