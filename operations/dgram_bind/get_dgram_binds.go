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

package dgram_bind

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

// GetDgramBindsHandlerFunc turns a function with the right signature into a get dgram binds handler
type GetDgramBindsHandlerFunc func(GetDgramBindsParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn GetDgramBindsHandlerFunc) Handle(params GetDgramBindsParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// GetDgramBindsHandler interface for that can handle valid get dgram binds params
type GetDgramBindsHandler interface {
	Handle(GetDgramBindsParams, interface{}) middleware.Responder
}

// NewGetDgramBinds creates a new http.Handler for the get dgram binds operation
func NewGetDgramBinds(ctx *middleware.Context, handler GetDgramBindsHandler) *GetDgramBinds {
	return &GetDgramBinds{Context: ctx, Handler: handler}
}

/*
	GetDgramBinds swagger:route GET /services/haproxy/configuration/dgram_binds DgramBind getDgramBinds

# Return an array of dgram binds

Returns an array of all dgram binds that are configured in specified log forward.
*/
type GetDgramBinds struct {
	Context *middleware.Context
	Handler GetDgramBindsHandler
}

func (o *GetDgramBinds) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewGetDgramBindsParams()
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

// GetDgramBindsOKBody get dgram binds o k body
//
// swagger:model GetDgramBindsOKBody
type GetDgramBindsOKBody struct {

	// version
	Version int64 `json:"_version,omitempty"`

	// data
	// Required: true
	Data models.DgramBinds `json:"data"`
}

// Validate validates this get dgram binds o k body
func (o *GetDgramBindsOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateData(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetDgramBindsOKBody) validateData(formats strfmt.Registry) error {

	if err := validate.Required("getDgramBindsOK"+"."+"data", "body", o.Data); err != nil {
		return err
	}

	if err := o.Data.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("getDgramBindsOK" + "." + "data")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("getDgramBindsOK" + "." + "data")
		}
		return err
	}

	return nil
}

// ContextValidate validate this get dgram binds o k body based on the context it is used
func (o *GetDgramBindsOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateData(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetDgramBindsOKBody) contextValidateData(ctx context.Context, formats strfmt.Registry) error {

	if err := o.Data.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("getDgramBindsOK" + "." + "data")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("getDgramBindsOK" + "." + "data")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *GetDgramBindsOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetDgramBindsOKBody) UnmarshalBinary(b []byte) error {
	var res GetDgramBindsOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
