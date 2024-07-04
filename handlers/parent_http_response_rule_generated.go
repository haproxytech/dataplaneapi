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
	"github.com/haproxytech/dataplaneapi/operations/http_response_rule"
)

type (
	CreateHTTPResponseRuleBackendHandlerImpl  CreateHTTPResponseRuleHandlerImpl
	CreateHTTPResponseRuleFrontendHandlerImpl CreateHTTPResponseRuleHandlerImpl
)

type (
	GetHTTPResponseRuleBackendHandlerImpl  GetHTTPResponseRuleHandlerImpl
	GetHTTPResponseRuleFrontendHandlerImpl GetHTTPResponseRuleHandlerImpl
)

type (
	GetAllHTTPResponseRuleBackendHandlerImpl  GetAllHTTPResponseRuleHandlerImpl
	GetAllHTTPResponseRuleFrontendHandlerImpl GetAllHTTPResponseRuleHandlerImpl
)

type (
	DeleteHTTPResponseRuleBackendHandlerImpl  DeleteHTTPResponseRuleHandlerImpl
	DeleteHTTPResponseRuleFrontendHandlerImpl DeleteHTTPResponseRuleHandlerImpl
)

type (
	ReplaceHTTPResponseRuleBackendHandlerImpl  ReplaceHTTPResponseRuleHandlerImpl
	ReplaceHTTPResponseRuleFrontendHandlerImpl ReplaceHTTPResponseRuleHandlerImpl
)

type (
	ReplaceAllHTTPResponseRuleBackendHandlerImpl  ReplaceAllHTTPResponseRuleHandlerImpl
	ReplaceAllHTTPResponseRuleFrontendHandlerImpl ReplaceAllHTTPResponseRuleHandlerImpl
)

func (h *CreateHTTPResponseRuleBackendHandlerImpl) Handle(params http_response_rule.CreateHTTPResponseRuleBackendParams, principal interface{}) middleware.Responder {
	g := CreateHTTPResponseRuleHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *CreateHTTPResponseRuleFrontendHandlerImpl) Handle(params http_response_rule.CreateHTTPResponseRuleFrontendParams, principal interface{}) middleware.Responder {
	g := CreateHTTPResponseRuleHandlerImpl(*h)
	pg := http_response_rule.CreateHTTPResponseRuleBackendParams(params)
	return g.Handle(cnconstants.FrontendParentType, pg, principal)
}

func (h *GetHTTPResponseRuleBackendHandlerImpl) Handle(params http_response_rule.GetHTTPResponseRuleBackendParams, principal interface{}) middleware.Responder {
	g := GetHTTPResponseRuleHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *GetHTTPResponseRuleFrontendHandlerImpl) Handle(params http_response_rule.GetHTTPResponseRuleFrontendParams, principal interface{}) middleware.Responder {
	g := GetHTTPResponseRuleHandlerImpl(*h)
	pg := http_response_rule.GetHTTPResponseRuleBackendParams(params)
	return g.Handle(cnconstants.FrontendParentType, pg, principal)
}

func (h *GetAllHTTPResponseRuleBackendHandlerImpl) Handle(params http_response_rule.GetAllHTTPResponseRuleBackendParams, principal interface{}) middleware.Responder {
	g := GetAllHTTPResponseRuleHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *GetAllHTTPResponseRuleFrontendHandlerImpl) Handle(params http_response_rule.GetAllHTTPResponseRuleFrontendParams, principal interface{}) middleware.Responder {
	g := GetAllHTTPResponseRuleHandlerImpl(*h)
	pg := http_response_rule.GetAllHTTPResponseRuleBackendParams(params)
	return g.Handle(cnconstants.FrontendParentType, pg, principal)
}

func (h *DeleteHTTPResponseRuleBackendHandlerImpl) Handle(params http_response_rule.DeleteHTTPResponseRuleBackendParams, principal interface{}) middleware.Responder {
	g := DeleteHTTPResponseRuleHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *DeleteHTTPResponseRuleFrontendHandlerImpl) Handle(params http_response_rule.DeleteHTTPResponseRuleFrontendParams, principal interface{}) middleware.Responder {
	g := DeleteHTTPResponseRuleHandlerImpl(*h)
	pg := http_response_rule.DeleteHTTPResponseRuleBackendParams(params)
	return g.Handle(cnconstants.FrontendParentType, pg, principal)
}

func (h *ReplaceHTTPResponseRuleBackendHandlerImpl) Handle(params http_response_rule.ReplaceHTTPResponseRuleBackendParams, principal interface{}) middleware.Responder {
	g := ReplaceHTTPResponseRuleHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *ReplaceHTTPResponseRuleFrontendHandlerImpl) Handle(params http_response_rule.ReplaceHTTPResponseRuleFrontendParams, principal interface{}) middleware.Responder {
	g := ReplaceHTTPResponseRuleHandlerImpl(*h)
	pg := http_response_rule.ReplaceHTTPResponseRuleBackendParams(params)
	return g.Handle(cnconstants.FrontendParentType, pg, principal)
}

func (h *ReplaceAllHTTPResponseRuleBackendHandlerImpl) Handle(params http_response_rule.ReplaceAllHTTPResponseRuleBackendParams, principal interface{}) middleware.Responder {
	g := ReplaceAllHTTPResponseRuleHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *ReplaceAllHTTPResponseRuleFrontendHandlerImpl) Handle(params http_response_rule.ReplaceAllHTTPResponseRuleFrontendParams, principal interface{}) middleware.Responder {
	g := ReplaceAllHTTPResponseRuleHandlerImpl(*h)
	pg := http_response_rule.ReplaceAllHTTPResponseRuleBackendParams(params)
	return g.Handle(cnconstants.FrontendParentType, pg, principal)
}
