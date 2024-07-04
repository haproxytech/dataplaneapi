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
	"github.com/haproxytech/dataplaneapi/operations/filter"
)

type (
	CreateFilterBackendHandlerImpl  CreateFilterHandlerImpl
	CreateFilterFrontendHandlerImpl CreateFilterHandlerImpl
)

type (
	GetFilterBackendHandlerImpl  GetFilterHandlerImpl
	GetFilterFrontendHandlerImpl GetFilterHandlerImpl
)

type (
	GetAllFilterBackendHandlerImpl  GetAllFilterHandlerImpl
	GetAllFilterFrontendHandlerImpl GetAllFilterHandlerImpl
)

type (
	DeleteFilterBackendHandlerImpl  DeleteFilterHandlerImpl
	DeleteFilterFrontendHandlerImpl DeleteFilterHandlerImpl
)

type (
	ReplaceFilterBackendHandlerImpl  ReplaceFilterHandlerImpl
	ReplaceFilterFrontendHandlerImpl ReplaceFilterHandlerImpl
)

type (
	ReplaceAllFilterBackendHandlerImpl  ReplaceAllFilterHandlerImpl
	ReplaceAllFilterFrontendHandlerImpl ReplaceAllFilterHandlerImpl
)

func (h *CreateFilterBackendHandlerImpl) Handle(params filter.CreateFilterBackendParams, principal interface{}) middleware.Responder {
	g := CreateFilterHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *CreateFilterFrontendHandlerImpl) Handle(params filter.CreateFilterFrontendParams, principal interface{}) middleware.Responder {
	g := CreateFilterHandlerImpl(*h)
	pg := filter.CreateFilterBackendParams(params)
	return g.Handle(cnconstants.FrontendParentType, pg, principal)
}

func (h *GetFilterBackendHandlerImpl) Handle(params filter.GetFilterBackendParams, principal interface{}) middleware.Responder {
	g := GetFilterHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *GetFilterFrontendHandlerImpl) Handle(params filter.GetFilterFrontendParams, principal interface{}) middleware.Responder {
	g := GetFilterHandlerImpl(*h)
	pg := filter.GetFilterBackendParams(params)
	return g.Handle(cnconstants.FrontendParentType, pg, principal)
}

func (h *GetAllFilterBackendHandlerImpl) Handle(params filter.GetAllFilterBackendParams, principal interface{}) middleware.Responder {
	g := GetAllFilterHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *GetAllFilterFrontendHandlerImpl) Handle(params filter.GetAllFilterFrontendParams, principal interface{}) middleware.Responder {
	g := GetAllFilterHandlerImpl(*h)
	pg := filter.GetAllFilterBackendParams(params)
	return g.Handle(cnconstants.FrontendParentType, pg, principal)
}

func (h *DeleteFilterBackendHandlerImpl) Handle(params filter.DeleteFilterBackendParams, principal interface{}) middleware.Responder {
	g := DeleteFilterHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *DeleteFilterFrontendHandlerImpl) Handle(params filter.DeleteFilterFrontendParams, principal interface{}) middleware.Responder {
	g := DeleteFilterHandlerImpl(*h)
	pg := filter.DeleteFilterBackendParams(params)
	return g.Handle(cnconstants.FrontendParentType, pg, principal)
}

func (h *ReplaceFilterBackendHandlerImpl) Handle(params filter.ReplaceFilterBackendParams, principal interface{}) middleware.Responder {
	g := ReplaceFilterHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *ReplaceFilterFrontendHandlerImpl) Handle(params filter.ReplaceFilterFrontendParams, principal interface{}) middleware.Responder {
	g := ReplaceFilterHandlerImpl(*h)
	pg := filter.ReplaceFilterBackendParams(params)
	return g.Handle(cnconstants.FrontendParentType, pg, principal)
}

func (h *ReplaceAllFilterBackendHandlerImpl) Handle(params filter.ReplaceAllFilterBackendParams, principal interface{}) middleware.Responder {
	g := ReplaceAllFilterHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *ReplaceAllFilterFrontendHandlerImpl) Handle(params filter.ReplaceAllFilterFrontendParams, principal interface{}) middleware.Responder {
	g := ReplaceAllFilterHandlerImpl(*h)
	pg := filter.ReplaceAllFilterBackendParams(params)
	return g.Handle(cnconstants.FrontendParentType, pg, principal)
}
