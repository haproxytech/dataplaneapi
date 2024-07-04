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
	"github.com/haproxytech/dataplaneapi/operations/http_check"
)

type (
	CreateHTTPCheckBackendHandlerImpl  CreateHTTPCheckHandlerImpl
	CreateHTTPCheckDefaultsHandlerImpl CreateHTTPCheckHandlerImpl
)

type (
	GetHTTPCheckBackendHandlerImpl  GetHTTPCheckHandlerImpl
	GetHTTPCheckDefaultsHandlerImpl GetHTTPCheckHandlerImpl
)

type (
	GetAllHTTPCheckBackendHandlerImpl  GetAllHTTPCheckHandlerImpl
	GetAllHTTPCheckDefaultsHandlerImpl GetAllHTTPCheckHandlerImpl
)

type (
	DeleteHTTPCheckBackendHandlerImpl  DeleteHTTPCheckHandlerImpl
	DeleteHTTPCheckDefaultsHandlerImpl DeleteHTTPCheckHandlerImpl
)

type (
	ReplaceHTTPCheckBackendHandlerImpl  ReplaceHTTPCheckHandlerImpl
	ReplaceHTTPCheckDefaultsHandlerImpl ReplaceHTTPCheckHandlerImpl
)

type (
	ReplaceAllHTTPCheckBackendHandlerImpl  ReplaceAllHTTPCheckHandlerImpl
	ReplaceAllHTTPCheckDefaultsHandlerImpl ReplaceAllHTTPCheckHandlerImpl
)

func (h *CreateHTTPCheckBackendHandlerImpl) Handle(params http_check.CreateHTTPCheckBackendParams, principal interface{}) middleware.Responder {
	g := CreateHTTPCheckHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *CreateHTTPCheckDefaultsHandlerImpl) Handle(params http_check.CreateHTTPCheckDefaultsParams, principal interface{}) middleware.Responder {
	g := CreateHTTPCheckHandlerImpl(*h)
	pg := http_check.CreateHTTPCheckBackendParams(params)
	return g.Handle(cnconstants.DefaultsParentType, pg, principal)
}

func (h *GetHTTPCheckBackendHandlerImpl) Handle(params http_check.GetHTTPCheckBackendParams, principal interface{}) middleware.Responder {
	g := GetHTTPCheckHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *GetHTTPCheckDefaultsHandlerImpl) Handle(params http_check.GetHTTPCheckDefaultsParams, principal interface{}) middleware.Responder {
	g := GetHTTPCheckHandlerImpl(*h)
	pg := http_check.GetHTTPCheckBackendParams(params)
	return g.Handle(cnconstants.DefaultsParentType, pg, principal)
}

func (h *GetAllHTTPCheckBackendHandlerImpl) Handle(params http_check.GetAllHTTPCheckBackendParams, principal interface{}) middleware.Responder {
	g := GetAllHTTPCheckHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *GetAllHTTPCheckDefaultsHandlerImpl) Handle(params http_check.GetAllHTTPCheckDefaultsParams, principal interface{}) middleware.Responder {
	g := GetAllHTTPCheckHandlerImpl(*h)
	pg := http_check.GetAllHTTPCheckBackendParams(params)
	return g.Handle(cnconstants.DefaultsParentType, pg, principal)
}

func (h *DeleteHTTPCheckBackendHandlerImpl) Handle(params http_check.DeleteHTTPCheckBackendParams, principal interface{}) middleware.Responder {
	g := DeleteHTTPCheckHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *DeleteHTTPCheckDefaultsHandlerImpl) Handle(params http_check.DeleteHTTPCheckDefaultsParams, principal interface{}) middleware.Responder {
	g := DeleteHTTPCheckHandlerImpl(*h)
	pg := http_check.DeleteHTTPCheckBackendParams(params)
	return g.Handle(cnconstants.DefaultsParentType, pg, principal)
}

func (h *ReplaceHTTPCheckBackendHandlerImpl) Handle(params http_check.ReplaceHTTPCheckBackendParams, principal interface{}) middleware.Responder {
	g := ReplaceHTTPCheckHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *ReplaceHTTPCheckDefaultsHandlerImpl) Handle(params http_check.ReplaceHTTPCheckDefaultsParams, principal interface{}) middleware.Responder {
	g := ReplaceHTTPCheckHandlerImpl(*h)
	pg := http_check.ReplaceHTTPCheckBackendParams(params)
	return g.Handle(cnconstants.DefaultsParentType, pg, principal)
}

func (h *ReplaceAllHTTPCheckBackendHandlerImpl) Handle(params http_check.ReplaceAllHTTPCheckBackendParams, principal interface{}) middleware.Responder {
	g := ReplaceAllHTTPCheckHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *ReplaceAllHTTPCheckDefaultsHandlerImpl) Handle(params http_check.ReplaceAllHTTPCheckDefaultsParams, principal interface{}) middleware.Responder {
	g := ReplaceAllHTTPCheckHandlerImpl(*h)
	pg := http_check.ReplaceAllHTTPCheckBackendParams(params)
	return g.Handle(cnconstants.DefaultsParentType, pg, principal)
}
