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

package http_request_rule

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// DeleteHTTPRequestRuleFrontendHandlerFunc turns a function with the right signature into a delete HTTP request rule frontend handler
type DeleteHTTPRequestRuleFrontendHandlerFunc func(DeleteHTTPRequestRuleFrontendParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn DeleteHTTPRequestRuleFrontendHandlerFunc) Handle(params DeleteHTTPRequestRuleFrontendParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// DeleteHTTPRequestRuleFrontendHandler interface for that can handle valid delete HTTP request rule frontend params
type DeleteHTTPRequestRuleFrontendHandler interface {
	Handle(DeleteHTTPRequestRuleFrontendParams, interface{}) middleware.Responder
}

// NewDeleteHTTPRequestRuleFrontend creates a new http.Handler for the delete HTTP request rule frontend operation
func NewDeleteHTTPRequestRuleFrontend(ctx *middleware.Context, handler DeleteHTTPRequestRuleFrontendHandler) *DeleteHTTPRequestRuleFrontend {
	return &DeleteHTTPRequestRuleFrontend{Context: ctx, Handler: handler}
}

/*
	DeleteHTTPRequestRuleFrontend swagger:route DELETE /services/haproxy/configuration/frontends/{parent_name}/http_request_rules/{index} HTTPRequestRule deleteHttpRequestRuleFrontend

# Delete a HTTP Request Rule

Deletes a HTTP Request Rule configuration by it's index from the specified parent.
*/
type DeleteHTTPRequestRuleFrontend struct {
	Context *middleware.Context
	Handler DeleteHTTPRequestRuleFrontendHandler
}

func (o *DeleteHTTPRequestRuleFrontend) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewDeleteHTTPRequestRuleFrontendParams()
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