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

package traces

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// ReplaceTracesHandlerFunc turns a function with the right signature into a replace traces handler
type ReplaceTracesHandlerFunc func(ReplaceTracesParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn ReplaceTracesHandlerFunc) Handle(params ReplaceTracesParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// ReplaceTracesHandler interface for that can handle valid replace traces params
type ReplaceTracesHandler interface {
	Handle(ReplaceTracesParams, interface{}) middleware.Responder
}

// NewReplaceTraces creates a new http.Handler for the replace traces operation
func NewReplaceTraces(ctx *middleware.Context, handler ReplaceTracesHandler) *ReplaceTraces {
	return &ReplaceTraces{Context: ctx, Handler: handler}
}

/*
	ReplaceTraces swagger:route PUT /services/haproxy/configuration/traces Traces replaceTraces

# Replace traces

Replace the traces section contents
*/
type ReplaceTraces struct {
	Context *middleware.Context
	Handler ReplaceTracesHandler
}

func (o *ReplaceTraces) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewReplaceTracesParams()
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
