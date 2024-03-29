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

package peer_entry

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// GetPeerEntryHandlerFunc turns a function with the right signature into a get peer entry handler
type GetPeerEntryHandlerFunc func(GetPeerEntryParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn GetPeerEntryHandlerFunc) Handle(params GetPeerEntryParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// GetPeerEntryHandler interface for that can handle valid get peer entry params
type GetPeerEntryHandler interface {
	Handle(GetPeerEntryParams, interface{}) middleware.Responder
}

// NewGetPeerEntry creates a new http.Handler for the get peer entry operation
func NewGetPeerEntry(ctx *middleware.Context, handler GetPeerEntryHandler) *GetPeerEntry {
	return &GetPeerEntry{Context: ctx, Handler: handler}
}

/*
	GetPeerEntry swagger:route GET /services/haproxy/configuration/peer_entries/{name} PeerEntry getPeerEntry

Return one peer_entry

Returns one peer_entry configuration by it's name in the specified peer section.
*/
type GetPeerEntry struct {
	Context *middleware.Context
	Handler GetPeerEntryHandler
}

func (o *GetPeerEntry) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewGetPeerEntryParams()
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
