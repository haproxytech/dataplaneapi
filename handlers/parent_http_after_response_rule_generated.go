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
	"github.com/haproxytech/dataplaneapi/operations/http_after_response_rule"
)

type (
	CreateHTTPAfterResponseRuleBackendHandlerImpl  CreateHTTPAfterResponseRuleHandlerImpl
	CreateHTTPAfterResponseRuleFrontendHandlerImpl CreateHTTPAfterResponseRuleHandlerImpl
	CreateHTTPAfterResponseRuleDefaultsHandlerImpl CreateHTTPAfterResponseRuleHandlerImpl
)

type (
	GetHTTPAfterResponseRuleBackendHandlerImpl  GetHTTPAfterResponseRuleHandlerImpl
	GetHTTPAfterResponseRuleFrontendHandlerImpl GetHTTPAfterResponseRuleHandlerImpl
	GetHTTPAfterResponseRuleDefaultsHandlerImpl GetHTTPAfterResponseRuleHandlerImpl
)

type (
	GetAllHTTPAfterResponseRuleBackendHandlerImpl  GetAllHTTPAfterResponseRuleHandlerImpl
	GetAllHTTPAfterResponseRuleFrontendHandlerImpl GetAllHTTPAfterResponseRuleHandlerImpl
	GetAllHTTPAfterResponseRuleDefaultsHandlerImpl GetAllHTTPAfterResponseRuleHandlerImpl
)

type (
	DeleteHTTPAfterResponseRuleBackendHandlerImpl  DeleteHTTPAfterResponseRuleHandlerImpl
	DeleteHTTPAfterResponseRuleFrontendHandlerImpl DeleteHTTPAfterResponseRuleHandlerImpl
	DeleteHTTPAfterResponseRuleDefaultsHandlerImpl DeleteHTTPAfterResponseRuleHandlerImpl
)

type (
	ReplaceHTTPAfterResponseRuleBackendHandlerImpl  ReplaceHTTPAfterResponseRuleHandlerImpl
	ReplaceHTTPAfterResponseRuleFrontendHandlerImpl ReplaceHTTPAfterResponseRuleHandlerImpl
	ReplaceHTTPAfterResponseRuleDefaultsHandlerImpl ReplaceHTTPAfterResponseRuleHandlerImpl
)

type (
	ReplaceAllHTTPAfterResponseRuleBackendHandlerImpl  ReplaceAllHTTPAfterResponseRuleHandlerImpl
	ReplaceAllHTTPAfterResponseRuleFrontendHandlerImpl ReplaceAllHTTPAfterResponseRuleHandlerImpl
	ReplaceAllHTTPAfterResponseRuleDefaultsHandlerImpl ReplaceAllHTTPAfterResponseRuleHandlerImpl
)

func (h *CreateHTTPAfterResponseRuleBackendHandlerImpl) Handle(params http_after_response_rule.CreateHTTPAfterResponseRuleBackendParams, principal interface{}) middleware.Responder {
	g := CreateHTTPAfterResponseRuleHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *CreateHTTPAfterResponseRuleFrontendHandlerImpl) Handle(params http_after_response_rule.CreateHTTPAfterResponseRuleFrontendParams, principal interface{}) middleware.Responder {
	g := CreateHTTPAfterResponseRuleHandlerImpl(*h)
	pg := http_after_response_rule.CreateHTTPAfterResponseRuleBackendParams(params)
	return g.Handle(cnconstants.FrontendParentType, pg, principal)
}

func (h *CreateHTTPAfterResponseRuleDefaultsHandlerImpl) Handle(params http_after_response_rule.CreateHTTPAfterResponseRuleDefaultsParams, principal interface{}) middleware.Responder {
	g := CreateHTTPAfterResponseRuleHandlerImpl(*h)
	pg := http_after_response_rule.CreateHTTPAfterResponseRuleBackendParams(params)
	return g.Handle(cnconstants.DefaultsParentType, pg, principal)
}

func (h *GetHTTPAfterResponseRuleBackendHandlerImpl) Handle(params http_after_response_rule.GetHTTPAfterResponseRuleBackendParams, principal interface{}) middleware.Responder {
	g := GetHTTPAfterResponseRuleHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *GetHTTPAfterResponseRuleFrontendHandlerImpl) Handle(params http_after_response_rule.GetHTTPAfterResponseRuleFrontendParams, principal interface{}) middleware.Responder {
	g := GetHTTPAfterResponseRuleHandlerImpl(*h)
	pg := http_after_response_rule.GetHTTPAfterResponseRuleBackendParams(params)
	return g.Handle(cnconstants.FrontendParentType, pg, principal)
}

func (h *GetHTTPAfterResponseRuleDefaultsHandlerImpl) Handle(params http_after_response_rule.GetHTTPAfterResponseRuleDefaultsParams, principal interface{}) middleware.Responder {
	g := GetHTTPAfterResponseRuleHandlerImpl(*h)
	pg := http_after_response_rule.GetHTTPAfterResponseRuleBackendParams(params)
	return g.Handle(cnconstants.DefaultsParentType, pg, principal)
}

func (h *GetAllHTTPAfterResponseRuleBackendHandlerImpl) Handle(params http_after_response_rule.GetAllHTTPAfterResponseRuleBackendParams, principal interface{}) middleware.Responder {
	g := GetAllHTTPAfterResponseRuleHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *GetAllHTTPAfterResponseRuleFrontendHandlerImpl) Handle(params http_after_response_rule.GetAllHTTPAfterResponseRuleFrontendParams, principal interface{}) middleware.Responder {
	g := GetAllHTTPAfterResponseRuleHandlerImpl(*h)
	pg := http_after_response_rule.GetAllHTTPAfterResponseRuleBackendParams(params)
	return g.Handle(cnconstants.FrontendParentType, pg, principal)
}

func (h *GetAllHTTPAfterResponseRuleDefaultsHandlerImpl) Handle(params http_after_response_rule.GetAllHTTPAfterResponseRuleDefaultsParams, principal interface{}) middleware.Responder {
	g := GetAllHTTPAfterResponseRuleHandlerImpl(*h)
	pg := http_after_response_rule.GetAllHTTPAfterResponseRuleBackendParams(params)
	return g.Handle(cnconstants.DefaultsParentType, pg, principal)
}

func (h *DeleteHTTPAfterResponseRuleBackendHandlerImpl) Handle(params http_after_response_rule.DeleteHTTPAfterResponseRuleBackendParams, principal interface{}) middleware.Responder {
	g := DeleteHTTPAfterResponseRuleHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *DeleteHTTPAfterResponseRuleFrontendHandlerImpl) Handle(params http_after_response_rule.DeleteHTTPAfterResponseRuleFrontendParams, principal interface{}) middleware.Responder {
	g := DeleteHTTPAfterResponseRuleHandlerImpl(*h)
	pg := http_after_response_rule.DeleteHTTPAfterResponseRuleBackendParams(params)
	return g.Handle(cnconstants.FrontendParentType, pg, principal)
}

func (h *DeleteHTTPAfterResponseRuleDefaultsHandlerImpl) Handle(params http_after_response_rule.DeleteHTTPAfterResponseRuleDefaultsParams, principal interface{}) middleware.Responder {
	g := DeleteHTTPAfterResponseRuleHandlerImpl(*h)
	pg := http_after_response_rule.DeleteHTTPAfterResponseRuleBackendParams(params)
	return g.Handle(cnconstants.DefaultsParentType, pg, principal)
}

func (h *ReplaceHTTPAfterResponseRuleBackendHandlerImpl) Handle(params http_after_response_rule.ReplaceHTTPAfterResponseRuleBackendParams, principal interface{}) middleware.Responder {
	g := ReplaceHTTPAfterResponseRuleHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *ReplaceHTTPAfterResponseRuleFrontendHandlerImpl) Handle(params http_after_response_rule.ReplaceHTTPAfterResponseRuleFrontendParams, principal interface{}) middleware.Responder {
	g := ReplaceHTTPAfterResponseRuleHandlerImpl(*h)
	pg := http_after_response_rule.ReplaceHTTPAfterResponseRuleBackendParams(params)
	return g.Handle(cnconstants.FrontendParentType, pg, principal)
}

func (h *ReplaceHTTPAfterResponseRuleDefaultsHandlerImpl) Handle(params http_after_response_rule.ReplaceHTTPAfterResponseRuleDefaultsParams, principal interface{}) middleware.Responder {
	g := ReplaceHTTPAfterResponseRuleHandlerImpl(*h)
	pg := http_after_response_rule.ReplaceHTTPAfterResponseRuleBackendParams(params)
	return g.Handle(cnconstants.DefaultsParentType, pg, principal)
}

func (h *ReplaceAllHTTPAfterResponseRuleBackendHandlerImpl) Handle(params http_after_response_rule.ReplaceAllHTTPAfterResponseRuleBackendParams, principal interface{}) middleware.Responder {
	g := ReplaceAllHTTPAfterResponseRuleHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *ReplaceAllHTTPAfterResponseRuleFrontendHandlerImpl) Handle(params http_after_response_rule.ReplaceAllHTTPAfterResponseRuleFrontendParams, principal interface{}) middleware.Responder {
	g := ReplaceAllHTTPAfterResponseRuleHandlerImpl(*h)
	pg := http_after_response_rule.ReplaceAllHTTPAfterResponseRuleBackendParams(params)
	return g.Handle(cnconstants.FrontendParentType, pg, principal)
}

func (h *ReplaceAllHTTPAfterResponseRuleDefaultsHandlerImpl) Handle(params http_after_response_rule.ReplaceAllHTTPAfterResponseRuleDefaultsParams, principal interface{}) middleware.Responder {
	g := ReplaceAllHTTPAfterResponseRuleHandlerImpl(*h)
	pg := http_after_response_rule.ReplaceAllHTTPAfterResponseRuleBackendParams(params)
	return g.Handle(cnconstants.DefaultsParentType, pg, principal)
}
