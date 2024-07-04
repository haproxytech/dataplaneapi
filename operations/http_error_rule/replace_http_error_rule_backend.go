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

package http_error_rule

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// ReplaceHTTPErrorRuleBackendHandlerFunc turns a function with the right signature into a replace HTTP error rule backend handler
type ReplaceHTTPErrorRuleBackendHandlerFunc func(ReplaceHTTPErrorRuleBackendParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn ReplaceHTTPErrorRuleBackendHandlerFunc) Handle(params ReplaceHTTPErrorRuleBackendParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// ReplaceHTTPErrorRuleBackendHandler interface for that can handle valid replace HTTP error rule backend params
type ReplaceHTTPErrorRuleBackendHandler interface {
	Handle(ReplaceHTTPErrorRuleBackendParams, interface{}) middleware.Responder
}

// NewReplaceHTTPErrorRuleBackend creates a new http.Handler for the replace HTTP error rule backend operation
func NewReplaceHTTPErrorRuleBackend(ctx *middleware.Context, handler ReplaceHTTPErrorRuleBackendHandler) *ReplaceHTTPErrorRuleBackend {
	return &ReplaceHTTPErrorRuleBackend{Context: ctx, Handler: handler}
}

/*
	ReplaceHTTPErrorRuleBackend swagger:route PUT /services/haproxy/configuration/backends/{parent_name}/http_error_rules/{index} HTTPErrorRule replaceHttpErrorRuleBackend

# Replace a HTTP Error Rule

Replaces a HTTP Error Rule configuration by its index in the specified parent.
*/
type ReplaceHTTPErrorRuleBackend struct {
	Context *middleware.Context
	Handler ReplaceHTTPErrorRuleBackendHandler
}

func (o *ReplaceHTTPErrorRuleBackend) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewReplaceHTTPErrorRuleBackendParams()
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
