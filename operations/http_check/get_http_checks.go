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

package http_check

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

// GetHTTPChecksHandlerFunc turns a function with the right signature into a get HTTP checks handler
type GetHTTPChecksHandlerFunc func(GetHTTPChecksParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn GetHTTPChecksHandlerFunc) Handle(params GetHTTPChecksParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// GetHTTPChecksHandler interface for that can handle valid get HTTP checks params
type GetHTTPChecksHandler interface {
	Handle(GetHTTPChecksParams, interface{}) middleware.Responder
}

// NewGetHTTPChecks creates a new http.Handler for the get HTTP checks operation
func NewGetHTTPChecks(ctx *middleware.Context, handler GetHTTPChecksHandler) *GetHTTPChecks {
	return &GetHTTPChecks{Context: ctx, Handler: handler}
}

/*
	GetHTTPChecks swagger:route GET /services/haproxy/configuration/http_checks HTTPCheck getHttpChecks

# Return an array of HTTP checks

Returns all HTTP checks that are configured in specified parent.
*/
type GetHTTPChecks struct {
	Context *middleware.Context
	Handler GetHTTPChecksHandler
}

func (o *GetHTTPChecks) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewGetHTTPChecksParams()
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

// GetHTTPChecksOKBody get HTTP checks o k body
//
// swagger:model GetHTTPChecksOKBody
type GetHTTPChecksOKBody struct {

	// version
	Version int64 `json:"_version,omitempty"`

	// data
	// Required: true
	Data models.HTTPChecks `json:"data"`
}

// Validate validates this get HTTP checks o k body
func (o *GetHTTPChecksOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateData(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetHTTPChecksOKBody) validateData(formats strfmt.Registry) error {

	if err := validate.Required("getHttpChecksOK"+"."+"data", "body", o.Data); err != nil {
		return err
	}

	if err := o.Data.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("getHttpChecksOK" + "." + "data")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("getHttpChecksOK" + "." + "data")
		}
		return err
	}

	return nil
}

// ContextValidate validate this get HTTP checks o k body based on the context it is used
func (o *GetHTTPChecksOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateData(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetHTTPChecksOKBody) contextValidateData(ctx context.Context, formats strfmt.Registry) error {

	if err := o.Data.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("getHttpChecksOK" + "." + "data")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("getHttpChecksOK" + "." + "data")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *GetHTTPChecksOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetHTTPChecksOKBody) UnmarshalBinary(b []byte) error {
	var res GetHTTPChecksOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
