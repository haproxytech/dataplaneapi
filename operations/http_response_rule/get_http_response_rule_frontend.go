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

// GetHTTPResponseRuleFrontendHandlerFunc turns a function with the right signature into a get HTTP response rule frontend handler
type GetHTTPResponseRuleFrontendHandlerFunc func(GetHTTPResponseRuleFrontendParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn GetHTTPResponseRuleFrontendHandlerFunc) Handle(params GetHTTPResponseRuleFrontendParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// GetHTTPResponseRuleFrontendHandler interface for that can handle valid get HTTP response rule frontend params
type GetHTTPResponseRuleFrontendHandler interface {
	Handle(GetHTTPResponseRuleFrontendParams, interface{}) middleware.Responder
}

// NewGetHTTPResponseRuleFrontend creates a new http.Handler for the get HTTP response rule frontend operation
func NewGetHTTPResponseRuleFrontend(ctx *middleware.Context, handler GetHTTPResponseRuleFrontendHandler) *GetHTTPResponseRuleFrontend {
	return &GetHTTPResponseRuleFrontend{Context: ctx, Handler: handler}
}

/*
	GetHTTPResponseRuleFrontend swagger:route GET /services/haproxy/configuration/frontends/{parent_name}/http_response_rules/{index} HTTPResponseRule getHttpResponseRuleFrontend

# Return one HTTP Response Rule

Returns one HTTP Response Rule configuration by it's index in the specified parent.
*/
type GetHTTPResponseRuleFrontend struct {
	Context *middleware.Context
	Handler GetHTTPResponseRuleFrontendHandler
}

func (o *GetHTTPResponseRuleFrontend) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewGetHTTPResponseRuleFrontendParams()
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
