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

package http_after_response_rule

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

	"github.com/haproxytech/client-native/v4/models"
)

// GetHTTPAfterResponseRulesHandlerFunc turns a function with the right signature into a get HTTP after response rules handler
type GetHTTPAfterResponseRulesHandlerFunc func(GetHTTPAfterResponseRulesParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn GetHTTPAfterResponseRulesHandlerFunc) Handle(params GetHTTPAfterResponseRulesParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// GetHTTPAfterResponseRulesHandler interface for that can handle valid get HTTP after response rules params
type GetHTTPAfterResponseRulesHandler interface {
	Handle(GetHTTPAfterResponseRulesParams, interface{}) middleware.Responder
}

// NewGetHTTPAfterResponseRules creates a new http.Handler for the get HTTP after response rules operation
func NewGetHTTPAfterResponseRules(ctx *middleware.Context, handler GetHTTPAfterResponseRulesHandler) *GetHTTPAfterResponseRules {
	return &GetHTTPAfterResponseRules{Context: ctx, Handler: handler}
}

/* GetHTTPAfterResponseRules swagger:route GET /services/haproxy/configuration/http_after_response_rules HTTPAfterResponseRule getHttpAfterResponseRules

Return an array of all HTTP After Response Rules

Returns all HTTP After Response Rules that are configured in specified parent.

*/
type GetHTTPAfterResponseRules struct {
	Context *middleware.Context
	Handler GetHTTPAfterResponseRulesHandler
}

func (o *GetHTTPAfterResponseRules) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewGetHTTPAfterResponseRulesParams()
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

// GetHTTPAfterResponseRulesOKBody get HTTP after response rules o k body
//
// swagger:model GetHTTPAfterResponseRulesOKBody
type GetHTTPAfterResponseRulesOKBody struct {

	// version
	Version int64 `json:"_version,omitempty"`

	// data
	// Required: true
	Data models.HTTPAfterResponseRules `json:"data"`
}

// Validate validates this get HTTP after response rules o k body
func (o *GetHTTPAfterResponseRulesOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateData(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetHTTPAfterResponseRulesOKBody) validateData(formats strfmt.Registry) error {

	if err := validate.Required("getHttpAfterResponseRulesOK"+"."+"data", "body", o.Data); err != nil {
		return err
	}

	if err := o.Data.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("getHttpAfterResponseRulesOK" + "." + "data")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("getHttpAfterResponseRulesOK" + "." + "data")
		}
		return err
	}

	return nil
}

// ContextValidate validate this get HTTP after response rules o k body based on the context it is used
func (o *GetHTTPAfterResponseRulesOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateData(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetHTTPAfterResponseRulesOKBody) contextValidateData(ctx context.Context, formats strfmt.Registry) error {

	if err := o.Data.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("getHttpAfterResponseRulesOK" + "." + "data")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("getHttpAfterResponseRulesOK" + "." + "data")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *GetHTTPAfterResponseRulesOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetHTTPAfterResponseRulesOKBody) UnmarshalBinary(b []byte) error {
	var res GetHTTPAfterResponseRulesOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}