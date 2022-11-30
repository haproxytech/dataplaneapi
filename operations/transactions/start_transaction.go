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
// Editing this file might prove futile when you re-run the generate command

import (
	"context"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// StartTransactionHandlerFunc turns a function with the right signature into a start transaction handler
type StartTransactionHandlerFunc func(StartTransactionParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn StartTransactionHandlerFunc) Handle(params StartTransactionParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// StartTransactionHandler interface for that can handle valid start transaction params
type StartTransactionHandler interface {
	Handle(StartTransactionParams, interface{}) middleware.Responder
}

// NewStartTransaction creates a new http.Handler for the start transaction operation
func NewStartTransaction(ctx *middleware.Context, handler StartTransactionHandler) *StartTransaction {
	return &StartTransaction{Context: ctx, Handler: handler}
}

/*
	StartTransaction swagger:route POST /services/haproxy/transactions Transactions startTransaction

# Start a new transaction

Starts a new transaction and returns it's id
*/
type StartTransaction struct {
	Context *middleware.Context
	Handler StartTransactionHandler
}

func (o *StartTransaction) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewStartTransactionParams()
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

// StartTransactionTooManyRequestsBody start transaction too many requests body
// Example: {"code":429,"message":"cannot start a new transaction, reached the maximum amount of 20 active transactions available"}
//
// swagger:model StartTransactionTooManyRequestsBody
type StartTransactionTooManyRequestsBody struct {

	// code
	Code int64 `json:"code,omitempty"`

	// message
	Message string `json:"message,omitempty"`
}

// Validate validates this start transaction too many requests body
func (o *StartTransactionTooManyRequestsBody) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this start transaction too many requests body based on context it is used
func (o *StartTransactionTooManyRequestsBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *StartTransactionTooManyRequestsBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *StartTransactionTooManyRequestsBody) UnmarshalBinary(b []byte) error {
	var res StartTransactionTooManyRequestsBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
