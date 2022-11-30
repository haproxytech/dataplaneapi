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

package peer_entry

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"context"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/haproxytech/client-native/v3/models"
)

// GetPeerEntryHandlerFunc turns a function with the right signature into a get peer entry handler
type GetPeerEntryHandlerFunc func(GetPeerEntryParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn GetPeerEntryHandlerFunc) Handle(params GetPeerEntryParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// GetPeerEntryHandler interface for that can handle valid get peer entry params
type GetPeerEntryHandler interface {
	Handle(GetPeerEntryParams, interface{}) middleware.Responder
}

// NewGetPeerEntry creates a new http.Handler for the get peer entry operation
func NewGetPeerEntry(ctx *middleware.Context, handler GetPeerEntryHandler) *GetPeerEntry {
	return &GetPeerEntry{Context: ctx, Handler: handler}
}

/*
	GetPeerEntry swagger:route GET /services/haproxy/configuration/peer_entries/{name} PeerEntry getPeerEntry

Return one peer_entry

Returns one peer_entry configuration by it's name in the specified peer section.
*/
type GetPeerEntry struct {
	Context *middleware.Context
	Handler GetPeerEntryHandler
}

func (o *GetPeerEntry) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewGetPeerEntryParams()
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

// GetPeerEntryOKBody get peer entry o k body
//
// swagger:model GetPeerEntryOKBody
type GetPeerEntryOKBody struct {

	// version
	Version int64 `json:"_version,omitempty"`

	// data
	Data *models.PeerEntry `json:"data,omitempty"`
}

// Validate validates this get peer entry o k body
func (o *GetPeerEntryOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateData(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetPeerEntryOKBody) validateData(formats strfmt.Registry) error {
	if swag.IsZero(o.Data) { // not required
		return nil
	}

	if o.Data != nil {
		if err := o.Data.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("getPeerEntryOK" + "." + "data")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("getPeerEntryOK" + "." + "data")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this get peer entry o k body based on the context it is used
func (o *GetPeerEntryOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateData(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetPeerEntryOKBody) contextValidateData(ctx context.Context, formats strfmt.Registry) error {

	if o.Data != nil {
		if err := o.Data.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("getPeerEntryOK" + "." + "data")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("getPeerEntryOK" + "." + "data")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (o *GetPeerEntryOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetPeerEntryOKBody) UnmarshalBinary(b []byte) error {
	var res GetPeerEntryOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
