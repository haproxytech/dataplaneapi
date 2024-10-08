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

package tcp_check

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// DeleteTCPCheckDefaultsHandlerFunc turns a function with the right signature into a delete TCP check defaults handler
type DeleteTCPCheckDefaultsHandlerFunc func(DeleteTCPCheckDefaultsParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn DeleteTCPCheckDefaultsHandlerFunc) Handle(params DeleteTCPCheckDefaultsParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// DeleteTCPCheckDefaultsHandler interface for that can handle valid delete TCP check defaults params
type DeleteTCPCheckDefaultsHandler interface {
	Handle(DeleteTCPCheckDefaultsParams, interface{}) middleware.Responder
}

// NewDeleteTCPCheckDefaults creates a new http.Handler for the delete TCP check defaults operation
func NewDeleteTCPCheckDefaults(ctx *middleware.Context, handler DeleteTCPCheckDefaultsHandler) *DeleteTCPCheckDefaults {
	return &DeleteTCPCheckDefaults{Context: ctx, Handler: handler}
}

/*
	DeleteTCPCheckDefaults swagger:route DELETE /services/haproxy/configuration/defaults/{parent_name}/tcp_checks/{index} TCPCheck deleteTcpCheckDefaults

# Delete a TCP check

Deletes a TCP check configuration by it's index from the specified parent.
*/
type DeleteTCPCheckDefaults struct {
	Context *middleware.Context
	Handler DeleteTCPCheckDefaultsHandler
}

func (o *DeleteTCPCheckDefaults) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewDeleteTCPCheckDefaultsParams()
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
