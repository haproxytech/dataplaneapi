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

package server

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// CreateServerPeerHandlerFunc turns a function with the right signature into a create server peer handler
type CreateServerPeerHandlerFunc func(CreateServerPeerParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn CreateServerPeerHandlerFunc) Handle(params CreateServerPeerParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// CreateServerPeerHandler interface for that can handle valid create server peer params
type CreateServerPeerHandler interface {
	Handle(CreateServerPeerParams, interface{}) middleware.Responder
}

// NewCreateServerPeer creates a new http.Handler for the create server peer operation
func NewCreateServerPeer(ctx *middleware.Context, handler CreateServerPeerHandler) *CreateServerPeer {
	return &CreateServerPeer{Context: ctx, Handler: handler}
}

/*
	CreateServerPeer swagger:route POST /services/haproxy/configuration/peers/{parent_name}/servers Server createServerPeer

# Add a new server

Adds a new server in the specified backend in the configuration file.
*/
type CreateServerPeer struct {
	Context *middleware.Context
	Handler CreateServerPeerHandler
}

func (o *CreateServerPeer) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewCreateServerPeerParams()
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