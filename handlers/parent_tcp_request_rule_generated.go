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
	"github.com/haproxytech/dataplaneapi/operations/tcp_request_rule"
)

type (
	CreateTCPRequestRuleBackendHandlerImpl  CreateTCPRequestRuleHandlerImpl
	CreateTCPRequestRuleFrontendHandlerImpl CreateTCPRequestRuleHandlerImpl
)

type (
	GetTCPRequestRuleBackendHandlerImpl  GetTCPRequestRuleHandlerImpl
	GetTCPRequestRuleFrontendHandlerImpl GetTCPRequestRuleHandlerImpl
)

type (
	GetAllTCPRequestRuleBackendHandlerImpl  GetAllTCPRequestRuleHandlerImpl
	GetAllTCPRequestRuleFrontendHandlerImpl GetAllTCPRequestRuleHandlerImpl
)

type (
	DeleteTCPRequestRuleBackendHandlerImpl  DeleteTCPRequestRuleHandlerImpl
	DeleteTCPRequestRuleFrontendHandlerImpl DeleteTCPRequestRuleHandlerImpl
)

type (
	ReplaceTCPRequestRuleBackendHandlerImpl  ReplaceTCPRequestRuleHandlerImpl
	ReplaceTCPRequestRuleFrontendHandlerImpl ReplaceTCPRequestRuleHandlerImpl
)

type (
	ReplaceAllTCPRequestRuleBackendHandlerImpl  ReplaceAllTCPRequestRuleHandlerImpl
	ReplaceAllTCPRequestRuleFrontendHandlerImpl ReplaceAllTCPRequestRuleHandlerImpl
)

func (h *CreateTCPRequestRuleBackendHandlerImpl) Handle(params tcp_request_rule.CreateTCPRequestRuleBackendParams, principal interface{}) middleware.Responder {
	g := CreateTCPRequestRuleHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *CreateTCPRequestRuleFrontendHandlerImpl) Handle(params tcp_request_rule.CreateTCPRequestRuleFrontendParams, principal interface{}) middleware.Responder {
	g := CreateTCPRequestRuleHandlerImpl(*h)
	pg := tcp_request_rule.CreateTCPRequestRuleBackendParams(params)
	return g.Handle(cnconstants.FrontendParentType, pg, principal)
}

func (h *GetTCPRequestRuleBackendHandlerImpl) Handle(params tcp_request_rule.GetTCPRequestRuleBackendParams, principal interface{}) middleware.Responder {
	g := GetTCPRequestRuleHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *GetTCPRequestRuleFrontendHandlerImpl) Handle(params tcp_request_rule.GetTCPRequestRuleFrontendParams, principal interface{}) middleware.Responder {
	g := GetTCPRequestRuleHandlerImpl(*h)
	pg := tcp_request_rule.GetTCPRequestRuleBackendParams(params)
	return g.Handle(cnconstants.FrontendParentType, pg, principal)
}

func (h *GetAllTCPRequestRuleBackendHandlerImpl) Handle(params tcp_request_rule.GetAllTCPRequestRuleBackendParams, principal interface{}) middleware.Responder {
	g := GetAllTCPRequestRuleHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *GetAllTCPRequestRuleFrontendHandlerImpl) Handle(params tcp_request_rule.GetAllTCPRequestRuleFrontendParams, principal interface{}) middleware.Responder {
	g := GetAllTCPRequestRuleHandlerImpl(*h)
	pg := tcp_request_rule.GetAllTCPRequestRuleBackendParams(params)
	return g.Handle(cnconstants.FrontendParentType, pg, principal)
}

func (h *DeleteTCPRequestRuleBackendHandlerImpl) Handle(params tcp_request_rule.DeleteTCPRequestRuleBackendParams, principal interface{}) middleware.Responder {
	g := DeleteTCPRequestRuleHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *DeleteTCPRequestRuleFrontendHandlerImpl) Handle(params tcp_request_rule.DeleteTCPRequestRuleFrontendParams, principal interface{}) middleware.Responder {
	g := DeleteTCPRequestRuleHandlerImpl(*h)
	pg := tcp_request_rule.DeleteTCPRequestRuleBackendParams(params)
	return g.Handle(cnconstants.FrontendParentType, pg, principal)
}

func (h *ReplaceTCPRequestRuleBackendHandlerImpl) Handle(params tcp_request_rule.ReplaceTCPRequestRuleBackendParams, principal interface{}) middleware.Responder {
	g := ReplaceTCPRequestRuleHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *ReplaceTCPRequestRuleFrontendHandlerImpl) Handle(params tcp_request_rule.ReplaceTCPRequestRuleFrontendParams, principal interface{}) middleware.Responder {
	g := ReplaceTCPRequestRuleHandlerImpl(*h)
	pg := tcp_request_rule.ReplaceTCPRequestRuleBackendParams(params)
	return g.Handle(cnconstants.FrontendParentType, pg, principal)
}

func (h *ReplaceAllTCPRequestRuleBackendHandlerImpl) Handle(params tcp_request_rule.ReplaceAllTCPRequestRuleBackendParams, principal interface{}) middleware.Responder {
	g := ReplaceAllTCPRequestRuleHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *ReplaceAllTCPRequestRuleFrontendHandlerImpl) Handle(params tcp_request_rule.ReplaceAllTCPRequestRuleFrontendParams, principal interface{}) middleware.Responder {
	g := ReplaceAllTCPRequestRuleHandlerImpl(*h)
	pg := tcp_request_rule.ReplaceAllTCPRequestRuleBackendParams(params)
	return g.Handle(cnconstants.FrontendParentType, pg, principal)
}
