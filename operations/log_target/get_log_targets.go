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

package log_target

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

	"github.com/haproxytech/client-native/v3/models"
)

// GetLogTargetsHandlerFunc turns a function with the right signature into a get log targets handler
type GetLogTargetsHandlerFunc func(GetLogTargetsParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn GetLogTargetsHandlerFunc) Handle(params GetLogTargetsParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// GetLogTargetsHandler interface for that can handle valid get log targets params
type GetLogTargetsHandler interface {
	Handle(GetLogTargetsParams, interface{}) middleware.Responder
}

// NewGetLogTargets creates a new http.Handler for the get log targets operation
func NewGetLogTargets(ctx *middleware.Context, handler GetLogTargetsHandler) *GetLogTargets {
	return &GetLogTargets{Context: ctx, Handler: handler}
}

/*
	GetLogTargets swagger:route GET /services/haproxy/configuration/log_targets LogTarget getLogTargets

# Return an array of all Log Targets

Returns all Log Targets that are configured in specified parent.
*/
type GetLogTargets struct {
	Context *middleware.Context
	Handler GetLogTargetsHandler
}

func (o *GetLogTargets) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewGetLogTargetsParams()
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

// GetLogTargetsOKBody get log targets o k body
//
// swagger:model GetLogTargetsOKBody
type GetLogTargetsOKBody struct {

	// version
	Version int64 `json:"_version,omitempty"`

	// data
	// Required: true
	Data models.LogTargets `json:"data"`
}

// Validate validates this get log targets o k body
func (o *GetLogTargetsOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateData(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetLogTargetsOKBody) validateData(formats strfmt.Registry) error {

	if err := validate.Required("getLogTargetsOK"+"."+"data", "body", o.Data); err != nil {
		return err
	}

	if err := o.Data.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("getLogTargetsOK" + "." + "data")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("getLogTargetsOK" + "." + "data")
		}
		return err
	}

	return nil
}

// ContextValidate validate this get log targets o k body based on the context it is used
func (o *GetLogTargetsOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateData(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetLogTargetsOKBody) contextValidateData(ctx context.Context, formats strfmt.Registry) error {

	if err := o.Data.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("getLogTargetsOK" + "." + "data")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("getLogTargetsOK" + "." + "data")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *GetLogTargetsOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetLogTargetsOKBody) UnmarshalBinary(b []byte) error {
	var res GetLogTargetsOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
