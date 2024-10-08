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

package crt_store

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// GetCrtStoresHandlerFunc turns a function with the right signature into a get crt stores handler
type GetCrtStoresHandlerFunc func(GetCrtStoresParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn GetCrtStoresHandlerFunc) Handle(params GetCrtStoresParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// GetCrtStoresHandler interface for that can handle valid get crt stores params
type GetCrtStoresHandler interface {
	Handle(GetCrtStoresParams, interface{}) middleware.Responder
}

// NewGetCrtStores creates a new http.Handler for the get crt stores operation
func NewGetCrtStores(ctx *middleware.Context, handler GetCrtStoresHandler) *GetCrtStores {
	return &GetCrtStores{Context: ctx, Handler: handler}
}

/*
	GetCrtStores swagger:route GET /services/haproxy/configuration/crt_stores CrtStore getCrtStores

# Return all the Certificate Stores

Returns an array of all the configured crt_store sections in HAProxy
*/
type GetCrtStores struct {
	Context *middleware.Context
	Handler GetCrtStoresHandler
}

func (o *GetCrtStores) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewGetCrtStoresParams()
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
