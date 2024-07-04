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
	"github.com/haproxytech/dataplaneapi/operations/tcp_check"
)

type (
	CreateTCPCheckBackendHandlerImpl  CreateTCPCheckHandlerImpl
	CreateTCPCheckDefaultsHandlerImpl CreateTCPCheckHandlerImpl
)

type (
	GetTCPCheckBackendHandlerImpl  GetTCPCheckHandlerImpl
	GetTCPCheckDefaultsHandlerImpl GetTCPCheckHandlerImpl
)

type (
	GetAllTCPCheckBackendHandlerImpl  GetAllTCPCheckHandlerImpl
	GetAllTCPCheckDefaultsHandlerImpl GetAllTCPCheckHandlerImpl
)

type (
	DeleteTCPCheckBackendHandlerImpl  DeleteTCPCheckHandlerImpl
	DeleteTCPCheckDefaultsHandlerImpl DeleteTCPCheckHandlerImpl
)

type (
	ReplaceTCPCheckBackendHandlerImpl  ReplaceTCPCheckHandlerImpl
	ReplaceTCPCheckDefaultsHandlerImpl ReplaceTCPCheckHandlerImpl
)

type (
	ReplaceAllTCPCheckBackendHandlerImpl  ReplaceAllTCPCheckHandlerImpl
	ReplaceAllTCPCheckDefaultsHandlerImpl ReplaceAllTCPCheckHandlerImpl
)

func (h *CreateTCPCheckBackendHandlerImpl) Handle(params tcp_check.CreateTCPCheckBackendParams, principal interface{}) middleware.Responder {
	g := CreateTCPCheckHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *CreateTCPCheckDefaultsHandlerImpl) Handle(params tcp_check.CreateTCPCheckDefaultsParams, principal interface{}) middleware.Responder {
	g := CreateTCPCheckHandlerImpl(*h)
	pg := tcp_check.CreateTCPCheckBackendParams(params)
	return g.Handle(cnconstants.DefaultsParentType, pg, principal)
}

func (h *GetTCPCheckBackendHandlerImpl) Handle(params tcp_check.GetTCPCheckBackendParams, principal interface{}) middleware.Responder {
	g := GetTCPCheckHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *GetTCPCheckDefaultsHandlerImpl) Handle(params tcp_check.GetTCPCheckDefaultsParams, principal interface{}) middleware.Responder {
	g := GetTCPCheckHandlerImpl(*h)
	pg := tcp_check.GetTCPCheckBackendParams(params)
	return g.Handle(cnconstants.DefaultsParentType, pg, principal)
}

func (h *GetAllTCPCheckBackendHandlerImpl) Handle(params tcp_check.GetAllTCPCheckBackendParams, principal interface{}) middleware.Responder {
	g := GetAllTCPCheckHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *GetAllTCPCheckDefaultsHandlerImpl) Handle(params tcp_check.GetAllTCPCheckDefaultsParams, principal interface{}) middleware.Responder {
	g := GetAllTCPCheckHandlerImpl(*h)
	pg := tcp_check.GetAllTCPCheckBackendParams(params)
	return g.Handle(cnconstants.DefaultsParentType, pg, principal)
}

func (h *DeleteTCPCheckBackendHandlerImpl) Handle(params tcp_check.DeleteTCPCheckBackendParams, principal interface{}) middleware.Responder {
	g := DeleteTCPCheckHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *DeleteTCPCheckDefaultsHandlerImpl) Handle(params tcp_check.DeleteTCPCheckDefaultsParams, principal interface{}) middleware.Responder {
	g := DeleteTCPCheckHandlerImpl(*h)
	pg := tcp_check.DeleteTCPCheckBackendParams(params)
	return g.Handle(cnconstants.DefaultsParentType, pg, principal)
}

func (h *ReplaceTCPCheckBackendHandlerImpl) Handle(params tcp_check.ReplaceTCPCheckBackendParams, principal interface{}) middleware.Responder {
	g := ReplaceTCPCheckHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *ReplaceTCPCheckDefaultsHandlerImpl) Handle(params tcp_check.ReplaceTCPCheckDefaultsParams, principal interface{}) middleware.Responder {
	g := ReplaceTCPCheckHandlerImpl(*h)
	pg := tcp_check.ReplaceTCPCheckBackendParams(params)
	return g.Handle(cnconstants.DefaultsParentType, pg, principal)
}

func (h *ReplaceAllTCPCheckBackendHandlerImpl) Handle(params tcp_check.ReplaceAllTCPCheckBackendParams, principal interface{}) middleware.Responder {
	g := ReplaceAllTCPCheckHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *ReplaceAllTCPCheckDefaultsHandlerImpl) Handle(params tcp_check.ReplaceAllTCPCheckDefaultsParams, principal interface{}) middleware.Responder {
	g := ReplaceAllTCPCheckHandlerImpl(*h)
	pg := tcp_check.ReplaceAllTCPCheckBackendParams(params)
	return g.Handle(cnconstants.DefaultsParentType, pg, principal)
}
