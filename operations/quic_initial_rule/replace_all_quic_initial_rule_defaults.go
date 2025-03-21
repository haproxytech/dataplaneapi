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

package quic_initial_rule

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// ReplaceAllQUICInitialRuleDefaultsHandlerFunc turns a function with the right signature into a replace all QUIC initial rule defaults handler
type ReplaceAllQUICInitialRuleDefaultsHandlerFunc func(ReplaceAllQUICInitialRuleDefaultsParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn ReplaceAllQUICInitialRuleDefaultsHandlerFunc) Handle(params ReplaceAllQUICInitialRuleDefaultsParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// ReplaceAllQUICInitialRuleDefaultsHandler interface for that can handle valid replace all QUIC initial rule defaults params
type ReplaceAllQUICInitialRuleDefaultsHandler interface {
	Handle(ReplaceAllQUICInitialRuleDefaultsParams, interface{}) middleware.Responder
}

// NewReplaceAllQUICInitialRuleDefaults creates a new http.Handler for the replace all QUIC initial rule defaults operation
func NewReplaceAllQUICInitialRuleDefaults(ctx *middleware.Context, handler ReplaceAllQUICInitialRuleDefaultsHandler) *ReplaceAllQUICInitialRuleDefaults {
	return &ReplaceAllQUICInitialRuleDefaults{Context: ctx, Handler: handler}
}

/*
	ReplaceAllQUICInitialRuleDefaults swagger:route PUT /services/haproxy/configuration/defaults/{parent_name}/quic_initial_rules QUICInitialRule replaceAllQuicInitialRuleDefaults

# Replace an QUIC Initial rules list

Replaces a whole list of QUIC Initial rules with the list given in parameter
*/
type ReplaceAllQUICInitialRuleDefaults struct {
	Context *middleware.Context
	Handler ReplaceAllQUICInitialRuleDefaultsHandler
}

func (o *ReplaceAllQUICInitialRuleDefaults) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewReplaceAllQUICInitialRuleDefaultsParams()
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
