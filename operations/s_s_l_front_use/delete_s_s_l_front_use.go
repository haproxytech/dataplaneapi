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

package s_s_l_front_use

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// DeleteSSLFrontUseHandlerFunc turns a function with the right signature into a delete s s l front use handler
type DeleteSSLFrontUseHandlerFunc func(DeleteSSLFrontUseParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn DeleteSSLFrontUseHandlerFunc) Handle(params DeleteSSLFrontUseParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// DeleteSSLFrontUseHandler interface for that can handle valid delete s s l front use params
type DeleteSSLFrontUseHandler interface {
	Handle(DeleteSSLFrontUseParams, interface{}) middleware.Responder
}

// NewDeleteSSLFrontUse creates a new http.Handler for the delete s s l front use operation
func NewDeleteSSLFrontUse(ctx *middleware.Context, handler DeleteSSLFrontUseHandler) *DeleteSSLFrontUse {
	return &DeleteSSLFrontUse{Context: ctx, Handler: handler}
}

/*
	DeleteSSLFrontUse swagger:route DELETE /services/haproxy/configuration/frontends/{parent_name}/ssl_front_uses/{index} SSLFrontUse deleteSSLFrontUse

# Delete an SSLFrontUse

Deletes an SSLFrontUse configuration by its index in the specified frontend.
*/
type DeleteSSLFrontUse struct {
	Context *middleware.Context
	Handler DeleteSSLFrontUseHandler
}

func (o *DeleteSSLFrontUse) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewDeleteSSLFrontUseParams()
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
