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
	"github.com/haproxytech/dataplaneapi/operations/tcp_response_rule"
)

type (
	CreateTCPResponseRuleBackendHandlerImpl CreateTCPResponseRuleHandlerImpl
)

type (
	GetTCPResponseRuleBackendHandlerImpl GetTCPResponseRuleHandlerImpl
)

type (
	GetAllTCPResponseRuleBackendHandlerImpl GetAllTCPResponseRuleHandlerImpl
)

type (
	DeleteTCPResponseRuleBackendHandlerImpl DeleteTCPResponseRuleHandlerImpl
)

type (
	ReplaceTCPResponseRuleBackendHandlerImpl ReplaceTCPResponseRuleHandlerImpl
)

type (
	ReplaceAllTCPResponseRuleBackendHandlerImpl ReplaceAllTCPResponseRuleHandlerImpl
)

func (h *CreateTCPResponseRuleBackendHandlerImpl) Handle(params tcp_response_rule.CreateTCPResponseRuleBackendParams, principal interface{}) middleware.Responder {
	g := CreateTCPResponseRuleHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *GetTCPResponseRuleBackendHandlerImpl) Handle(params tcp_response_rule.GetTCPResponseRuleBackendParams, principal interface{}) middleware.Responder {
	g := GetTCPResponseRuleHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *GetAllTCPResponseRuleBackendHandlerImpl) Handle(params tcp_response_rule.GetAllTCPResponseRuleBackendParams, principal interface{}) middleware.Responder {
	g := GetAllTCPResponseRuleHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *DeleteTCPResponseRuleBackendHandlerImpl) Handle(params tcp_response_rule.DeleteTCPResponseRuleBackendParams, principal interface{}) middleware.Responder {
	g := DeleteTCPResponseRuleHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *ReplaceTCPResponseRuleBackendHandlerImpl) Handle(params tcp_response_rule.ReplaceTCPResponseRuleBackendParams, principal interface{}) middleware.Responder {
	g := ReplaceTCPResponseRuleHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *ReplaceAllTCPResponseRuleBackendHandlerImpl) Handle(params tcp_response_rule.ReplaceAllTCPResponseRuleBackendParams, principal interface{}) middleware.Responder {
	g := ReplaceAllTCPResponseRuleHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}
