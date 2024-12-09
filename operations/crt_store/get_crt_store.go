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

// GetCrtStoreHandlerFunc turns a function with the right signature into a get crt store handler
type GetCrtStoreHandlerFunc func(GetCrtStoreParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn GetCrtStoreHandlerFunc) Handle(params GetCrtStoreParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// GetCrtStoreHandler interface for that can handle valid get crt store params
type GetCrtStoreHandler interface {
	Handle(GetCrtStoreParams, interface{}) middleware.Responder
}

// NewGetCrtStore creates a new http.Handler for the get crt store operation
func NewGetCrtStore(ctx *middleware.Context, handler GetCrtStoreHandler) *GetCrtStore {
	return &GetCrtStore{Context: ctx, Handler: handler}
}

/*
	GetCrtStore swagger:route GET /services/haproxy/configuration/crt_stores/{name} CrtStore getCrtStore

# Return a Certificate Store

Returns crt_store section by its name
*/
type GetCrtStore struct {
	Context *middleware.Context
	Handler GetCrtStoreHandler
}

func (o *GetCrtStore) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewGetCrtStoreParams()
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