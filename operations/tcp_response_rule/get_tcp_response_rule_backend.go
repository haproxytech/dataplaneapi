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

package tcp_response_rule

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// GetTCPResponseRuleBackendHandlerFunc turns a function with the right signature into a get TCP response rule backend handler
type GetTCPResponseRuleBackendHandlerFunc func(GetTCPResponseRuleBackendParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn GetTCPResponseRuleBackendHandlerFunc) Handle(params GetTCPResponseRuleBackendParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// GetTCPResponseRuleBackendHandler interface for that can handle valid get TCP response rule backend params
type GetTCPResponseRuleBackendHandler interface {
	Handle(GetTCPResponseRuleBackendParams, interface{}) middleware.Responder
}

// NewGetTCPResponseRuleBackend creates a new http.Handler for the get TCP response rule backend operation
func NewGetTCPResponseRuleBackend(ctx *middleware.Context, handler GetTCPResponseRuleBackendHandler) *GetTCPResponseRuleBackend {
	return &GetTCPResponseRuleBackend{Context: ctx, Handler: handler}
}

/*
	GetTCPResponseRuleBackend swagger:route GET /services/haproxy/configuration/backends/{parent_name}/tcp_response_rules/{index} TCPResponseRule getTcpResponseRuleBackend

# Return one TCP Response Rule

Returns one TCP Response Rule configuration by it's index in the specified backend.
*/
type GetTCPResponseRuleBackend struct {
	Context *middleware.Context
	Handler GetTCPResponseRuleBackendHandler
}

func (o *GetTCPResponseRuleBackend) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewGetTCPResponseRuleBackendParams()
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
