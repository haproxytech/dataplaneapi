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

// ReplaceHTTPRequestRuleHandlerFunc turns a function with the right signature into a replace HTTP request rule handler
type ReplaceHTTPRequestRuleHandlerFunc func(ReplaceHTTPRequestRuleParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn ReplaceHTTPRequestRuleHandlerFunc) Handle(params ReplaceHTTPRequestRuleParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// ReplaceHTTPRequestRuleHandler interface for that can handle valid replace HTTP request rule params
type ReplaceHTTPRequestRuleHandler interface {
	Handle(ReplaceHTTPRequestRuleParams, interface{}) middleware.Responder
}

// NewReplaceHTTPRequestRule creates a new http.Handler for the replace HTTP request rule operation
func NewReplaceHTTPRequestRule(ctx *middleware.Context, handler ReplaceHTTPRequestRuleHandler) *ReplaceHTTPRequestRule {
	return &ReplaceHTTPRequestRule{Context: ctx, Handler: handler}
}

/*
	ReplaceHTTPRequestRule swagger:route PUT /services/haproxy/configuration/http_request_rules/{index} HTTPRequestRule replaceHttpRequestRule

# Replace a HTTP Request Rule

Replaces a HTTP Request Rule configuration by it's index in the specified parent.
*/
type ReplaceHTTPRequestRule struct {
	Context *middleware.Context
	Handler ReplaceHTTPRequestRuleHandler
}

func (o *ReplaceHTTPRequestRule) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewReplaceHTTPRequestRuleParams()
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
