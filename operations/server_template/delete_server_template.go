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

package server_template

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// DeleteServerTemplateHandlerFunc turns a function with the right signature into a delete server template handler
type DeleteServerTemplateHandlerFunc func(DeleteServerTemplateParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn DeleteServerTemplateHandlerFunc) Handle(params DeleteServerTemplateParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// DeleteServerTemplateHandler interface for that can handle valid delete server template params
type DeleteServerTemplateHandler interface {
	Handle(DeleteServerTemplateParams, interface{}) middleware.Responder
}

// NewDeleteServerTemplate creates a new http.Handler for the delete server template operation
func NewDeleteServerTemplate(ctx *middleware.Context, handler DeleteServerTemplateHandler) *DeleteServerTemplate {
	return &DeleteServerTemplate{Context: ctx, Handler: handler}
}

/*
	DeleteServerTemplate swagger:route DELETE /services/haproxy/configuration/backends/{parent_name}/server_templates/{prefix} ServerTemplate deleteServerTemplate

# Delete a server template

Deletes a server template configuration by it's prefix in the specified backend.
*/
type DeleteServerTemplate struct {
	Context *middleware.Context
	Handler DeleteServerTemplateHandler
}

func (o *DeleteServerTemplate) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewDeleteServerTemplateParams()
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
