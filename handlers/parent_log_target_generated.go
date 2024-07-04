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
	"github.com/haproxytech/dataplaneapi/operations/log_target"
)

type (
	CreateLogTargetBackendHandlerImpl    CreateLogTargetHandlerImpl
	CreateLogTargetFrontendHandlerImpl   CreateLogTargetHandlerImpl
	CreateLogTargetDefaultsHandlerImpl   CreateLogTargetHandlerImpl
	CreateLogTargetPeerHandlerImpl       CreateLogTargetHandlerImpl
	CreateLogTargetLogForwardHandlerImpl CreateLogTargetHandlerImpl
)

type (
	GetLogTargetBackendHandlerImpl    GetLogTargetHandlerImpl
	GetLogTargetFrontendHandlerImpl   GetLogTargetHandlerImpl
	GetLogTargetDefaultsHandlerImpl   GetLogTargetHandlerImpl
	GetLogTargetPeerHandlerImpl       GetLogTargetHandlerImpl
	GetLogTargetLogForwardHandlerImpl GetLogTargetHandlerImpl
)

type (
	GetAllLogTargetBackendHandlerImpl    GetAllLogTargetHandlerImpl
	GetAllLogTargetFrontendHandlerImpl   GetAllLogTargetHandlerImpl
	GetAllLogTargetDefaultsHandlerImpl   GetAllLogTargetHandlerImpl
	GetAllLogTargetPeerHandlerImpl       GetAllLogTargetHandlerImpl
	GetAllLogTargetLogForwardHandlerImpl GetAllLogTargetHandlerImpl
)

type (
	DeleteLogTargetBackendHandlerImpl    DeleteLogTargetHandlerImpl
	DeleteLogTargetFrontendHandlerImpl   DeleteLogTargetHandlerImpl
	DeleteLogTargetDefaultsHandlerImpl   DeleteLogTargetHandlerImpl
	DeleteLogTargetPeerHandlerImpl       DeleteLogTargetHandlerImpl
	DeleteLogTargetLogForwardHandlerImpl DeleteLogTargetHandlerImpl
)

type (
	ReplaceLogTargetBackendHandlerImpl    ReplaceLogTargetHandlerImpl
	ReplaceLogTargetFrontendHandlerImpl   ReplaceLogTargetHandlerImpl
	ReplaceLogTargetDefaultsHandlerImpl   ReplaceLogTargetHandlerImpl
	ReplaceLogTargetPeerHandlerImpl       ReplaceLogTargetHandlerImpl
	ReplaceLogTargetLogForwardHandlerImpl ReplaceLogTargetHandlerImpl
)

type (
	ReplaceAllLogTargetBackendHandlerImpl    ReplaceAllLogTargetHandlerImpl
	ReplaceAllLogTargetFrontendHandlerImpl   ReplaceAllLogTargetHandlerImpl
	ReplaceAllLogTargetDefaultsHandlerImpl   ReplaceAllLogTargetHandlerImpl
	ReplaceAllLogTargetPeerHandlerImpl       ReplaceAllLogTargetHandlerImpl
	ReplaceAllLogTargetLogForwardHandlerImpl ReplaceAllLogTargetHandlerImpl
)

func (h *CreateLogTargetBackendHandlerImpl) Handle(params log_target.CreateLogTargetBackendParams, principal interface{}) middleware.Responder {
	g := CreateLogTargetHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *CreateLogTargetFrontendHandlerImpl) Handle(params log_target.CreateLogTargetFrontendParams, principal interface{}) middleware.Responder {
	g := CreateLogTargetHandlerImpl(*h)
	pg := log_target.CreateLogTargetBackendParams(params)
	return g.Handle(cnconstants.FrontendParentType, pg, principal)
}

func (h *CreateLogTargetDefaultsHandlerImpl) Handle(params log_target.CreateLogTargetDefaultsParams, principal interface{}) middleware.Responder {
	g := CreateLogTargetHandlerImpl(*h)
	pg := log_target.CreateLogTargetBackendParams(params)
	return g.Handle(cnconstants.DefaultsParentType, pg, principal)
}

func (h *CreateLogTargetPeerHandlerImpl) Handle(params log_target.CreateLogTargetPeerParams, principal interface{}) middleware.Responder {
	g := CreateLogTargetHandlerImpl(*h)
	pg := log_target.CreateLogTargetBackendParams(params)
	return g.Handle(cnconstants.PeerParentType, pg, principal)
}

func (h *CreateLogTargetLogForwardHandlerImpl) Handle(params log_target.CreateLogTargetLogForwardParams, principal interface{}) middleware.Responder {
	g := CreateLogTargetHandlerImpl(*h)
	pg := log_target.CreateLogTargetBackendParams(params)
	return g.Handle(cnconstants.LogForwardParentType, pg, principal)
}

func (h *GetLogTargetBackendHandlerImpl) Handle(params log_target.GetLogTargetBackendParams, principal interface{}) middleware.Responder {
	g := GetLogTargetHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *GetLogTargetFrontendHandlerImpl) Handle(params log_target.GetLogTargetFrontendParams, principal interface{}) middleware.Responder {
	g := GetLogTargetHandlerImpl(*h)
	pg := log_target.GetLogTargetBackendParams(params)
	return g.Handle(cnconstants.FrontendParentType, pg, principal)
}

func (h *GetLogTargetDefaultsHandlerImpl) Handle(params log_target.GetLogTargetDefaultsParams, principal interface{}) middleware.Responder {
	g := GetLogTargetHandlerImpl(*h)
	pg := log_target.GetLogTargetBackendParams(params)
	return g.Handle(cnconstants.DefaultsParentType, pg, principal)
}

func (h *GetLogTargetPeerHandlerImpl) Handle(params log_target.GetLogTargetPeerParams, principal interface{}) middleware.Responder {
	g := GetLogTargetHandlerImpl(*h)
	pg := log_target.GetLogTargetBackendParams(params)
	return g.Handle(cnconstants.PeerParentType, pg, principal)
}

func (h *GetLogTargetLogForwardHandlerImpl) Handle(params log_target.GetLogTargetLogForwardParams, principal interface{}) middleware.Responder {
	g := GetLogTargetHandlerImpl(*h)
	pg := log_target.GetLogTargetBackendParams(params)
	return g.Handle(cnconstants.LogForwardParentType, pg, principal)
}

func (h *GetAllLogTargetBackendHandlerImpl) Handle(params log_target.GetAllLogTargetBackendParams, principal interface{}) middleware.Responder {
	g := GetAllLogTargetHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *GetAllLogTargetFrontendHandlerImpl) Handle(params log_target.GetAllLogTargetFrontendParams, principal interface{}) middleware.Responder {
	g := GetAllLogTargetHandlerImpl(*h)
	pg := log_target.GetAllLogTargetBackendParams(params)
	return g.Handle(cnconstants.FrontendParentType, pg, principal)
}

func (h *GetAllLogTargetDefaultsHandlerImpl) Handle(params log_target.GetAllLogTargetDefaultsParams, principal interface{}) middleware.Responder {
	g := GetAllLogTargetHandlerImpl(*h)
	pg := log_target.GetAllLogTargetBackendParams(params)
	return g.Handle(cnconstants.DefaultsParentType, pg, principal)
}

func (h *GetAllLogTargetPeerHandlerImpl) Handle(params log_target.GetAllLogTargetPeerParams, principal interface{}) middleware.Responder {
	g := GetAllLogTargetHandlerImpl(*h)
	pg := log_target.GetAllLogTargetBackendParams(params)
	return g.Handle(cnconstants.PeerParentType, pg, principal)
}

func (h *GetAllLogTargetLogForwardHandlerImpl) Handle(params log_target.GetAllLogTargetLogForwardParams, principal interface{}) middleware.Responder {
	g := GetAllLogTargetHandlerImpl(*h)
	pg := log_target.GetAllLogTargetBackendParams(params)
	return g.Handle(cnconstants.LogForwardParentType, pg, principal)
}

func (h *DeleteLogTargetBackendHandlerImpl) Handle(params log_target.DeleteLogTargetBackendParams, principal interface{}) middleware.Responder {
	g := DeleteLogTargetHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *DeleteLogTargetFrontendHandlerImpl) Handle(params log_target.DeleteLogTargetFrontendParams, principal interface{}) middleware.Responder {
	g := DeleteLogTargetHandlerImpl(*h)
	pg := log_target.DeleteLogTargetBackendParams(params)
	return g.Handle(cnconstants.FrontendParentType, pg, principal)
}

func (h *DeleteLogTargetDefaultsHandlerImpl) Handle(params log_target.DeleteLogTargetDefaultsParams, principal interface{}) middleware.Responder {
	g := DeleteLogTargetHandlerImpl(*h)
	pg := log_target.DeleteLogTargetBackendParams(params)
	return g.Handle(cnconstants.DefaultsParentType, pg, principal)
}

func (h *DeleteLogTargetPeerHandlerImpl) Handle(params log_target.DeleteLogTargetPeerParams, principal interface{}) middleware.Responder {
	g := DeleteLogTargetHandlerImpl(*h)
	pg := log_target.DeleteLogTargetBackendParams(params)
	return g.Handle(cnconstants.PeerParentType, pg, principal)
}

func (h *DeleteLogTargetLogForwardHandlerImpl) Handle(params log_target.DeleteLogTargetLogForwardParams, principal interface{}) middleware.Responder {
	g := DeleteLogTargetHandlerImpl(*h)
	pg := log_target.DeleteLogTargetBackendParams(params)
	return g.Handle(cnconstants.LogForwardParentType, pg, principal)
}

func (h *ReplaceLogTargetBackendHandlerImpl) Handle(params log_target.ReplaceLogTargetBackendParams, principal interface{}) middleware.Responder {
	g := ReplaceLogTargetHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *ReplaceLogTargetFrontendHandlerImpl) Handle(params log_target.ReplaceLogTargetFrontendParams, principal interface{}) middleware.Responder {
	g := ReplaceLogTargetHandlerImpl(*h)
	pg := log_target.ReplaceLogTargetBackendParams(params)
	return g.Handle(cnconstants.FrontendParentType, pg, principal)
}

func (h *ReplaceLogTargetDefaultsHandlerImpl) Handle(params log_target.ReplaceLogTargetDefaultsParams, principal interface{}) middleware.Responder {
	g := ReplaceLogTargetHandlerImpl(*h)
	pg := log_target.ReplaceLogTargetBackendParams(params)
	return g.Handle(cnconstants.DefaultsParentType, pg, principal)
}

func (h *ReplaceLogTargetPeerHandlerImpl) Handle(params log_target.ReplaceLogTargetPeerParams, principal interface{}) middleware.Responder {
	g := ReplaceLogTargetHandlerImpl(*h)
	pg := log_target.ReplaceLogTargetBackendParams(params)
	return g.Handle(cnconstants.PeerParentType, pg, principal)
}

func (h *ReplaceLogTargetLogForwardHandlerImpl) Handle(params log_target.ReplaceLogTargetLogForwardParams, principal interface{}) middleware.Responder {
	g := ReplaceLogTargetHandlerImpl(*h)
	pg := log_target.ReplaceLogTargetBackendParams(params)
	return g.Handle(cnconstants.LogForwardParentType, pg, principal)
}

func (h *ReplaceAllLogTargetBackendHandlerImpl) Handle(params log_target.ReplaceAllLogTargetBackendParams, principal interface{}) middleware.Responder {
	g := ReplaceAllLogTargetHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *ReplaceAllLogTargetFrontendHandlerImpl) Handle(params log_target.ReplaceAllLogTargetFrontendParams, principal interface{}) middleware.Responder {
	g := ReplaceAllLogTargetHandlerImpl(*h)
	pg := log_target.ReplaceAllLogTargetBackendParams(params)
	return g.Handle(cnconstants.FrontendParentType, pg, principal)
}

func (h *ReplaceAllLogTargetDefaultsHandlerImpl) Handle(params log_target.ReplaceAllLogTargetDefaultsParams, principal interface{}) middleware.Responder {
	g := ReplaceAllLogTargetHandlerImpl(*h)
	pg := log_target.ReplaceAllLogTargetBackendParams(params)
	return g.Handle(cnconstants.DefaultsParentType, pg, principal)
}

func (h *ReplaceAllLogTargetPeerHandlerImpl) Handle(params log_target.ReplaceAllLogTargetPeerParams, principal interface{}) middleware.Responder {
	g := ReplaceAllLogTargetHandlerImpl(*h)
	pg := log_target.ReplaceAllLogTargetBackendParams(params)
	return g.Handle(cnconstants.PeerParentType, pg, principal)
}

func (h *ReplaceAllLogTargetLogForwardHandlerImpl) Handle(params log_target.ReplaceAllLogTargetLogForwardParams, principal interface{}) middleware.Responder {
	g := ReplaceAllLogTargetHandlerImpl(*h)
	pg := log_target.ReplaceAllLogTargetBackendParams(params)
	return g.Handle(cnconstants.LogForwardParentType, pg, principal)
}
