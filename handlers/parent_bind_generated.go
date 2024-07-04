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
	"github.com/haproxytech/dataplaneapi/operations/bind"
)

type (
	CreateBindFrontendHandlerImpl   CreateBindHandlerImpl
	CreateBindLogForwardHandlerImpl CreateBindHandlerImpl
	CreateBindPeerHandlerImpl       CreateBindHandlerImpl
)

type (
	GetBindFrontendHandlerImpl   GetBindHandlerImpl
	GetBindLogForwardHandlerImpl GetBindHandlerImpl
	GetBindPeerHandlerImpl       GetBindHandlerImpl
)

type (
	GetAllBindFrontendHandlerImpl   GetAllBindHandlerImpl
	GetAllBindLogForwardHandlerImpl GetAllBindHandlerImpl
	GetAllBindPeerHandlerImpl       GetAllBindHandlerImpl
)

type (
	DeleteBindFrontendHandlerImpl   DeleteBindHandlerImpl
	DeleteBindLogForwardHandlerImpl DeleteBindHandlerImpl
	DeleteBindPeerHandlerImpl       DeleteBindHandlerImpl
)

type (
	ReplaceBindFrontendHandlerImpl   ReplaceBindHandlerImpl
	ReplaceBindLogForwardHandlerImpl ReplaceBindHandlerImpl
	ReplaceBindPeerHandlerImpl       ReplaceBindHandlerImpl
)

func (h *CreateBindFrontendHandlerImpl) Handle(params bind.CreateBindFrontendParams, principal interface{}) middleware.Responder {
	g := CreateBindHandlerImpl(*h)
	return g.Handle(cnconstants.FrontendParentType, params, principal)
}

func (h *CreateBindLogForwardHandlerImpl) Handle(params bind.CreateBindLogForwardParams, principal interface{}) middleware.Responder {
	g := CreateBindHandlerImpl(*h)
	pg := bind.CreateBindFrontendParams(params)
	return g.Handle(cnconstants.LogForwardParentType, pg, principal)
}

func (h *CreateBindPeerHandlerImpl) Handle(params bind.CreateBindPeerParams, principal interface{}) middleware.Responder {
	g := CreateBindHandlerImpl(*h)
	pg := bind.CreateBindFrontendParams(params)
	return g.Handle(cnconstants.PeerParentType, pg, principal)
}

func (h *GetBindFrontendHandlerImpl) Handle(params bind.GetBindFrontendParams, principal interface{}) middleware.Responder {
	g := GetBindHandlerImpl(*h)
	return g.Handle(cnconstants.FrontendParentType, params, principal)
}

func (h *GetBindLogForwardHandlerImpl) Handle(params bind.GetBindLogForwardParams, principal interface{}) middleware.Responder {
	g := GetBindHandlerImpl(*h)
	pg := bind.GetBindFrontendParams(params)
	return g.Handle(cnconstants.LogForwardParentType, pg, principal)
}

func (h *GetBindPeerHandlerImpl) Handle(params bind.GetBindPeerParams, principal interface{}) middleware.Responder {
	g := GetBindHandlerImpl(*h)
	pg := bind.GetBindFrontendParams(params)
	return g.Handle(cnconstants.PeerParentType, pg, principal)
}

func (h *GetAllBindFrontendHandlerImpl) Handle(params bind.GetAllBindFrontendParams, principal interface{}) middleware.Responder {
	g := GetAllBindHandlerImpl(*h)
	return g.Handle(cnconstants.FrontendParentType, params, principal)
}

func (h *GetAllBindLogForwardHandlerImpl) Handle(params bind.GetAllBindLogForwardParams, principal interface{}) middleware.Responder {
	g := GetAllBindHandlerImpl(*h)
	pg := bind.GetAllBindFrontendParams(params)
	return g.Handle(cnconstants.LogForwardParentType, pg, principal)
}

func (h *GetAllBindPeerHandlerImpl) Handle(params bind.GetAllBindPeerParams, principal interface{}) middleware.Responder {
	g := GetAllBindHandlerImpl(*h)
	pg := bind.GetAllBindFrontendParams(params)
	return g.Handle(cnconstants.PeerParentType, pg, principal)
}

func (h *DeleteBindFrontendHandlerImpl) Handle(params bind.DeleteBindFrontendParams, principal interface{}) middleware.Responder {
	g := DeleteBindHandlerImpl(*h)
	return g.Handle(cnconstants.FrontendParentType, params, principal)
}

func (h *DeleteBindLogForwardHandlerImpl) Handle(params bind.DeleteBindLogForwardParams, principal interface{}) middleware.Responder {
	g := DeleteBindHandlerImpl(*h)
	pg := bind.DeleteBindFrontendParams(params)
	return g.Handle(cnconstants.LogForwardParentType, pg, principal)
}

func (h *DeleteBindPeerHandlerImpl) Handle(params bind.DeleteBindPeerParams, principal interface{}) middleware.Responder {
	g := DeleteBindHandlerImpl(*h)
	pg := bind.DeleteBindFrontendParams(params)
	return g.Handle(cnconstants.PeerParentType, pg, principal)
}

func (h *ReplaceBindFrontendHandlerImpl) Handle(params bind.ReplaceBindFrontendParams, principal interface{}) middleware.Responder {
	g := ReplaceBindHandlerImpl(*h)
	return g.Handle(cnconstants.FrontendParentType, params, principal)
}

func (h *ReplaceBindLogForwardHandlerImpl) Handle(params bind.ReplaceBindLogForwardParams, principal interface{}) middleware.Responder {
	g := ReplaceBindHandlerImpl(*h)
	pg := bind.ReplaceBindFrontendParams(params)
	return g.Handle(cnconstants.LogForwardParentType, pg, principal)
}

func (h *ReplaceBindPeerHandlerImpl) Handle(params bind.ReplaceBindPeerParams, principal interface{}) middleware.Responder {
	g := ReplaceBindHandlerImpl(*h)
	pg := bind.ReplaceBindFrontendParams(params)
	return g.Handle(cnconstants.PeerParentType, pg, principal)
}
