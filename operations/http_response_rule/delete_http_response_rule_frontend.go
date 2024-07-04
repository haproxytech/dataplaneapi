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
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// DeleteHTTPResponseRuleFrontendHandlerFunc turns a function with the right signature into a delete HTTP response rule frontend handler
type DeleteHTTPResponseRuleFrontendHandlerFunc func(DeleteHTTPResponseRuleFrontendParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn DeleteHTTPResponseRuleFrontendHandlerFunc) Handle(params DeleteHTTPResponseRuleFrontendParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// DeleteHTTPResponseRuleFrontendHandler interface for that can handle valid delete HTTP response rule frontend params
type DeleteHTTPResponseRuleFrontendHandler interface {
	Handle(DeleteHTTPResponseRuleFrontendParams, interface{}) middleware.Responder
}

// NewDeleteHTTPResponseRuleFrontend creates a new http.Handler for the delete HTTP response rule frontend operation
func NewDeleteHTTPResponseRuleFrontend(ctx *middleware.Context, handler DeleteHTTPResponseRuleFrontendHandler) *DeleteHTTPResponseRuleFrontend {
	return &DeleteHTTPResponseRuleFrontend{Context: ctx, Handler: handler}
}

/*
	DeleteHTTPResponseRuleFrontend swagger:route DELETE /services/haproxy/configuration/frontends/{parent_name}/http_response_rules/{index} HTTPResponseRule deleteHttpResponseRuleFrontend

# Delete a HTTP Response Rule

Deletes a HTTP Response Rule configuration by it's index from the specified parent.
*/
type DeleteHTTPResponseRuleFrontend struct {
	Context *middleware.Context
	Handler DeleteHTTPResponseRuleFrontendHandler
}

func (o *DeleteHTTPResponseRuleFrontend) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewDeleteHTTPResponseRuleFrontendParams()
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
