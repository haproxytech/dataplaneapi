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
	"github.com/haproxytech/dataplaneapi/operations/server"
)

type (
	CreateServerBackendHandlerImpl CreateServerHandlerImpl
	CreateServerPeerHandlerImpl    CreateServerHandlerImpl
	CreateServerRingHandlerImpl    CreateServerHandlerImpl
)

type (
	GetServerBackendHandlerImpl GetServerHandlerImpl
	GetServerPeerHandlerImpl    GetServerHandlerImpl
	GetServerRingHandlerImpl    GetServerHandlerImpl
)

type (
	GetAllServerBackendHandlerImpl GetAllServerHandlerImpl
	GetAllServerPeerHandlerImpl    GetAllServerHandlerImpl
	GetAllServerRingHandlerImpl    GetAllServerHandlerImpl
)

type (
	DeleteServerBackendHandlerImpl DeleteServerHandlerImpl
	DeleteServerPeerHandlerImpl    DeleteServerHandlerImpl
	DeleteServerRingHandlerImpl    DeleteServerHandlerImpl
)

type (
	ReplaceServerBackendHandlerImpl ReplaceServerHandlerImpl
	ReplaceServerPeerHandlerImpl    ReplaceServerHandlerImpl
	ReplaceServerRingHandlerImpl    ReplaceServerHandlerImpl
)

func (h *CreateServerBackendHandlerImpl) Handle(params server.CreateServerBackendParams, principal interface{}) middleware.Responder {
	g := CreateServerHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *CreateServerPeerHandlerImpl) Handle(params server.CreateServerPeerParams, principal interface{}) middleware.Responder {
	g := CreateServerHandlerImpl(*h)
	pg := server.CreateServerBackendParams(params)
	return g.Handle(cnconstants.PeerParentType, pg, principal)
}

func (h *CreateServerRingHandlerImpl) Handle(params server.CreateServerRingParams, principal interface{}) middleware.Responder {
	g := CreateServerHandlerImpl(*h)
	pg := server.CreateServerBackendParams(params)
	return g.Handle(cnconstants.RingParentType, pg, principal)
}

func (h *GetServerBackendHandlerImpl) Handle(params server.GetServerBackendParams, principal interface{}) middleware.Responder {
	g := GetServerHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *GetServerPeerHandlerImpl) Handle(params server.GetServerPeerParams, principal interface{}) middleware.Responder {
	g := GetServerHandlerImpl(*h)
	pg := server.GetServerBackendParams(params)
	return g.Handle(cnconstants.PeerParentType, pg, principal)
}

func (h *GetServerRingHandlerImpl) Handle(params server.GetServerRingParams, principal interface{}) middleware.Responder {
	g := GetServerHandlerImpl(*h)
	pg := server.GetServerBackendParams(params)
	return g.Handle(cnconstants.RingParentType, pg, principal)
}

func (h *GetAllServerBackendHandlerImpl) Handle(params server.GetAllServerBackendParams, principal interface{}) middleware.Responder {
	g := GetAllServerHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *GetAllServerPeerHandlerImpl) Handle(params server.GetAllServerPeerParams, principal interface{}) middleware.Responder {
	g := GetAllServerHandlerImpl(*h)
	pg := server.GetAllServerBackendParams(params)
	return g.Handle(cnconstants.PeerParentType, pg, principal)
}

func (h *GetAllServerRingHandlerImpl) Handle(params server.GetAllServerRingParams, principal interface{}) middleware.Responder {
	g := GetAllServerHandlerImpl(*h)
	pg := server.GetAllServerBackendParams(params)
	return g.Handle(cnconstants.RingParentType, pg, principal)
}

func (h *DeleteServerBackendHandlerImpl) Handle(params server.DeleteServerBackendParams, principal interface{}) middleware.Responder {
	g := DeleteServerHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *DeleteServerPeerHandlerImpl) Handle(params server.DeleteServerPeerParams, principal interface{}) middleware.Responder {
	g := DeleteServerHandlerImpl(*h)
	pg := server.DeleteServerBackendParams(params)
	return g.Handle(cnconstants.PeerParentType, pg, principal)
}

func (h *DeleteServerRingHandlerImpl) Handle(params server.DeleteServerRingParams, principal interface{}) middleware.Responder {
	g := DeleteServerHandlerImpl(*h)
	pg := server.DeleteServerBackendParams(params)
	return g.Handle(cnconstants.RingParentType, pg, principal)
}

func (h *ReplaceServerBackendHandlerImpl) Handle(params server.ReplaceServerBackendParams, principal interface{}) middleware.Responder {
	g := ReplaceServerHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *ReplaceServerPeerHandlerImpl) Handle(params server.ReplaceServerPeerParams, principal interface{}) middleware.Responder {
	g := ReplaceServerHandlerImpl(*h)
	pg := server.ReplaceServerBackendParams(params)
	return g.Handle(cnconstants.PeerParentType, pg, principal)
}

func (h *ReplaceServerRingHandlerImpl) Handle(params server.ReplaceServerRingParams, principal interface{}) middleware.Responder {
	g := ReplaceServerHandlerImpl(*h)
	pg := server.ReplaceServerBackendParams(params)
	return g.Handle(cnconstants.RingParentType, pg, principal)
}
