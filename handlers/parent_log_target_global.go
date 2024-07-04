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
	CreateLogTargetGlobalHandlerImpl CreateLogTargetHandlerImpl
)

type (
	GetLogTargetGlobalHandlerImpl GetLogTargetHandlerImpl
)

type (
	GetAllLogTargetGlobalHandlerImpl GetAllLogTargetHandlerImpl
)

type (
	DeleteLogTargetGlobalHandlerImpl DeleteLogTargetHandlerImpl
)

type (
	ReplaceLogTargetGlobalHandlerImpl ReplaceLogTargetHandlerImpl
)

type (
	ReplaceAllLogTargetGlobalHandlerImpl ReplaceAllLogTargetHandlerImpl
)

func (h *CreateLogTargetGlobalHandlerImpl) Handle(params log_target.CreateLogTargetGlobalParams, principal interface{}) middleware.Responder {
	g := CreateLogTargetHandlerImpl(*h)
	paramsG := log_target.CreateLogTargetBackendParams{
		Data:          params.Data,
		HTTPRequest:   params.HTTPRequest,
		ForceReload:   params.ForceReload,
		Index:         params.Index,
		ParentName:    "", // empty for Global
		TransactionID: params.TransactionID,
		Version:       params.Version,
	}
	return g.Handle(cnconstants.GlobalParentType, paramsG, principal)
}

func (h *GetLogTargetGlobalHandlerImpl) Handle(params log_target.GetLogTargetGlobalParams, principal interface{}) middleware.Responder {
	g := GetLogTargetHandlerImpl(*h)
	paramsG := log_target.GetLogTargetBackendParams{
		HTTPRequest:   params.HTTPRequest,
		Index:         params.Index,
		ParentName:    "", // empty for Global
		TransactionID: params.TransactionID,
	}
	return g.Handle(cnconstants.GlobalParentType, paramsG, principal)
}

func (h *GetAllLogTargetGlobalHandlerImpl) Handle(params log_target.GetAllLogTargetGlobalParams, principal interface{}) middleware.Responder {
	g := GetAllLogTargetHandlerImpl(*h)
	paramsG := log_target.GetAllLogTargetBackendParams{
		HTTPRequest:   params.HTTPRequest,
		ParentName:    "", // empty for Global
		TransactionID: params.TransactionID,
	}
	return g.Handle(cnconstants.GlobalParentType, paramsG, principal)
}

func (h *DeleteLogTargetGlobalHandlerImpl) Handle(params log_target.DeleteLogTargetGlobalParams, principal interface{}) middleware.Responder {
	g := DeleteLogTargetHandlerImpl(*h)
	paramsG := log_target.DeleteLogTargetBackendParams{
		HTTPRequest:   params.HTTPRequest,
		ForceReload:   params.ForceReload,
		Index:         params.Index,
		ParentName:    "", // empty for Global
		TransactionID: params.TransactionID,
		Version:       params.Version,
	}
	return g.Handle(cnconstants.GlobalParentType, paramsG, principal)
}

func (h *ReplaceLogTargetGlobalHandlerImpl) Handle(params log_target.ReplaceLogTargetGlobalParams, principal interface{}) middleware.Responder {
	g := ReplaceLogTargetHandlerImpl(*h)
	paramsG := log_target.ReplaceLogTargetBackendParams{
		Data:          params.Data,
		HTTPRequest:   params.HTTPRequest,
		ForceReload:   params.ForceReload,
		Index:         params.Index,
		ParentName:    "", // empty for Global
		TransactionID: params.TransactionID,
		Version:       params.Version,
	}
	return g.Handle(cnconstants.GlobalParentType, paramsG, principal)
}

func (h *ReplaceAllLogTargetGlobalHandlerImpl) Handle(params log_target.ReplaceAllLogTargetGlobalParams, principal interface{}) middleware.Responder {
	g := ReplaceAllLogTargetHandlerImpl(*h)
	paramsG := log_target.ReplaceAllLogTargetBackendParams{
		Data:          params.Data,
		HTTPRequest:   params.HTTPRequest,
		ForceReload:   params.ForceReload,
		ParentName:    "", // empty for Global
		TransactionID: params.TransactionID,
		Version:       params.Version,
	}
	return g.Handle(cnconstants.GlobalParentType, paramsG, principal)
}
