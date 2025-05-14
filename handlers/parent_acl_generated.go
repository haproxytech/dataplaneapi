// Copyright 2019 HAProxy Technologies
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package handlers

import (
	"github.com/go-openapi/runtime/middleware"

	cnconstants "github.com/haproxytech/client-native/v6/configuration/parents"
	"github.com/haproxytech/dataplaneapi/operations/acl"
)

type (
	CreateACLBackendHandlerImpl  CreateACLHandlerImpl
	CreateACLFrontendHandlerImpl CreateACLHandlerImpl
	CreateACLFCGIAppHandlerImpl  CreateACLHandlerImpl
	CreateACLDefaultsHandlerImpl CreateACLHandlerImpl
)

type (
	GetACLBackendHandlerImpl  GetACLHandlerImpl
	GetACLFrontendHandlerImpl GetACLHandlerImpl
	GetACLFCGIAppHandlerImpl  GetACLHandlerImpl
	GetACLDefaultsHandlerImpl GetACLHandlerImpl
)

type (
	GetAllACLBackendHandlerImpl  GetAllACLHandlerImpl
	GetAllACLFrontendHandlerImpl GetAllACLHandlerImpl
	GetAllACLFCGIAppHandlerImpl  GetAllACLHandlerImpl
	GetAllACLDefaultsHandlerImpl GetAllACLHandlerImpl
)

type (
	DeleteACLBackendHandlerImpl  DeleteACLHandlerImpl
	DeleteACLFrontendHandlerImpl DeleteACLHandlerImpl
	DeleteACLFCGIAppHandlerImpl  DeleteACLHandlerImpl
	DeleteACLDefaultsHandlerImpl DeleteACLHandlerImpl
)

type (
	ReplaceACLBackendHandlerImpl  ReplaceACLHandlerImpl
	ReplaceACLFrontendHandlerImpl ReplaceACLHandlerImpl
	ReplaceACLFCGIAppHandlerImpl  ReplaceACLHandlerImpl
	ReplaceACLDefaultsHandlerImpl ReplaceACLHandlerImpl
)

type (
	ReplaceAllACLBackendHandlerImpl  ReplaceAllACLHandlerImpl
	ReplaceAllACLFrontendHandlerImpl ReplaceAllACLHandlerImpl
	ReplaceAllACLFCGIAppHandlerImpl  ReplaceAllACLHandlerImpl
	ReplaceAllACLDefaultsHandlerImpl ReplaceAllACLHandlerImpl
)

func (h *CreateACLBackendHandlerImpl) Handle(params acl.CreateACLBackendParams, principal interface{}) middleware.Responder {
	g := CreateACLHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *CreateACLFrontendHandlerImpl) Handle(params acl.CreateACLFrontendParams, principal interface{}) middleware.Responder {
	g := CreateACLHandlerImpl(*h)
	pg := acl.CreateACLBackendParams(params)
	return g.Handle(cnconstants.FrontendParentType, pg, principal)
}

func (h *CreateACLFCGIAppHandlerImpl) Handle(params acl.CreateACLFCGIAppParams, principal interface{}) middleware.Responder {
	g := CreateACLHandlerImpl(*h)
	pg := acl.CreateACLBackendParams(params)
	return g.Handle(cnconstants.FCGIAppParentType, pg, principal)
}

func (h *CreateACLDefaultsHandlerImpl) Handle(params acl.CreateACLDefaultsParams, principal interface{}) middleware.Responder {
	g := CreateACLHandlerImpl(*h)
	pg := acl.CreateACLBackendParams(params)
	return g.Handle(cnconstants.DefaultsParentType, pg, principal)
}

func (h *GetACLBackendHandlerImpl) Handle(params acl.GetACLBackendParams, principal interface{}) middleware.Responder {
	g := GetACLHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *GetACLFrontendHandlerImpl) Handle(params acl.GetACLFrontendParams, principal interface{}) middleware.Responder {
	g := GetACLHandlerImpl(*h)
	pg := acl.GetACLBackendParams(params)
	return g.Handle(cnconstants.FrontendParentType, pg, principal)
}

func (h *GetACLFCGIAppHandlerImpl) Handle(params acl.GetACLFCGIAppParams, principal interface{}) middleware.Responder {
	g := GetACLHandlerImpl(*h)
	pg := acl.GetACLBackendParams(params)
	return g.Handle(cnconstants.FCGIAppParentType, pg, principal)
}

func (h *GetACLDefaultsHandlerImpl) Handle(params acl.GetACLDefaultsParams, principal interface{}) middleware.Responder {
	g := GetACLHandlerImpl(*h)
	pg := acl.GetACLBackendParams(params)
	return g.Handle(cnconstants.DefaultsParentType, pg, principal)
}

func (h *GetAllACLBackendHandlerImpl) Handle(params acl.GetAllACLBackendParams, principal interface{}) middleware.Responder {
	g := GetAllACLHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *GetAllACLFrontendHandlerImpl) Handle(params acl.GetAllACLFrontendParams, principal interface{}) middleware.Responder {
	g := GetAllACLHandlerImpl(*h)
	pg := acl.GetAllACLBackendParams(params)
	return g.Handle(cnconstants.FrontendParentType, pg, principal)
}

func (h *GetAllACLFCGIAppHandlerImpl) Handle(params acl.GetAllACLFCGIAppParams, principal interface{}) middleware.Responder {
	g := GetAllACLHandlerImpl(*h)
	pg := acl.GetAllACLBackendParams(params)
	return g.Handle(cnconstants.FCGIAppParentType, pg, principal)
}

func (h *GetAllACLDefaultsHandlerImpl) Handle(params acl.GetAllACLDefaultsParams, principal interface{}) middleware.Responder {
	g := GetAllACLHandlerImpl(*h)
	pg := acl.GetAllACLBackendParams(params)
	return g.Handle(cnconstants.DefaultsParentType, pg, principal)
}

func (h *DeleteACLBackendHandlerImpl) Handle(params acl.DeleteACLBackendParams, principal interface{}) middleware.Responder {
	g := DeleteACLHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *DeleteACLFrontendHandlerImpl) Handle(params acl.DeleteACLFrontendParams, principal interface{}) middleware.Responder {
	g := DeleteACLHandlerImpl(*h)
	pg := acl.DeleteACLBackendParams(params)
	return g.Handle(cnconstants.FrontendParentType, pg, principal)
}

func (h *DeleteACLFCGIAppHandlerImpl) Handle(params acl.DeleteACLFCGIAppParams, principal interface{}) middleware.Responder {
	g := DeleteACLHandlerImpl(*h)
	pg := acl.DeleteACLBackendParams(params)
	return g.Handle(cnconstants.FCGIAppParentType, pg, principal)
}

func (h *DeleteACLDefaultsHandlerImpl) Handle(params acl.DeleteACLDefaultsParams, principal interface{}) middleware.Responder {
	g := DeleteACLHandlerImpl(*h)
	pg := acl.DeleteACLBackendParams(params)
	return g.Handle(cnconstants.DefaultsParentType, pg, principal)
}

func (h *ReplaceACLBackendHandlerImpl) Handle(params acl.ReplaceACLBackendParams, principal interface{}) middleware.Responder {
	g := ReplaceACLHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *ReplaceACLFrontendHandlerImpl) Handle(params acl.ReplaceACLFrontendParams, principal interface{}) middleware.Responder {
	g := ReplaceACLHandlerImpl(*h)
	pg := acl.ReplaceACLBackendParams(params)
	return g.Handle(cnconstants.FrontendParentType, pg, principal)
}

func (h *ReplaceACLFCGIAppHandlerImpl) Handle(params acl.ReplaceACLFCGIAppParams, principal interface{}) middleware.Responder {
	g := ReplaceACLHandlerImpl(*h)
	pg := acl.ReplaceACLBackendParams(params)
	return g.Handle(cnconstants.FCGIAppParentType, pg, principal)
}

func (h *ReplaceACLDefaultsHandlerImpl) Handle(params acl.ReplaceACLDefaultsParams, principal interface{}) middleware.Responder {
	g := ReplaceACLHandlerImpl(*h)
	pg := acl.ReplaceACLBackendParams(params)
	return g.Handle(cnconstants.DefaultsParentType, pg, principal)
}

func (h *ReplaceAllACLBackendHandlerImpl) Handle(params acl.ReplaceAllACLBackendParams, principal interface{}) middleware.Responder {
	g := ReplaceAllACLHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *ReplaceAllACLFrontendHandlerImpl) Handle(params acl.ReplaceAllACLFrontendParams, principal interface{}) middleware.Responder {
	g := ReplaceAllACLHandlerImpl(*h)
	pg := acl.ReplaceAllACLBackendParams(params)
	return g.Handle(cnconstants.FrontendParentType, pg, principal)
}

func (h *ReplaceAllACLFCGIAppHandlerImpl) Handle(params acl.ReplaceAllACLFCGIAppParams, principal interface{}) middleware.Responder {
	g := ReplaceAllACLHandlerImpl(*h)
	pg := acl.ReplaceAllACLBackendParams(params)
	return g.Handle(cnconstants.FCGIAppParentType, pg, principal)
}

func (h *ReplaceAllACLDefaultsHandlerImpl) Handle(params acl.ReplaceAllACLDefaultsParams, principal interface{}) middleware.Responder {
	g := ReplaceAllACLHandlerImpl(*h)
	pg := acl.ReplaceAllACLBackendParams(params)
	return g.Handle(cnconstants.DefaultsParentType, pg, principal)
}
