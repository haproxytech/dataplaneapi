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
	"github.com/haproxytech/dataplaneapi/operations/http_request_rule"
)

type (
	CreateHTTPRequestRuleBackendHandlerImpl  CreateHTTPRequestRuleHandlerImpl
	CreateHTTPRequestRuleFrontendHandlerImpl CreateHTTPRequestRuleHandlerImpl
	CreateHTTPRequestRuleDefaultsHandlerImpl CreateHTTPRequestRuleHandlerImpl
)

type (
	GetHTTPRequestRuleBackendHandlerImpl  GetHTTPRequestRuleHandlerImpl
	GetHTTPRequestRuleFrontendHandlerImpl GetHTTPRequestRuleHandlerImpl
	GetHTTPRequestRuleDefaultsHandlerImpl GetHTTPRequestRuleHandlerImpl
)

type (
	GetAllHTTPRequestRuleBackendHandlerImpl  GetAllHTTPRequestRuleHandlerImpl
	GetAllHTTPRequestRuleFrontendHandlerImpl GetAllHTTPRequestRuleHandlerImpl
	GetAllHTTPRequestRuleDefaultsHandlerImpl GetAllHTTPRequestRuleHandlerImpl
)

type (
	DeleteHTTPRequestRuleBackendHandlerImpl  DeleteHTTPRequestRuleHandlerImpl
	DeleteHTTPRequestRuleFrontendHandlerImpl DeleteHTTPRequestRuleHandlerImpl
	DeleteHTTPRequestRuleDefaultsHandlerImpl DeleteHTTPRequestRuleHandlerImpl
)

type (
	ReplaceHTTPRequestRuleBackendHandlerImpl  ReplaceHTTPRequestRuleHandlerImpl
	ReplaceHTTPRequestRuleFrontendHandlerImpl ReplaceHTTPRequestRuleHandlerImpl
	ReplaceHTTPRequestRuleDefaultsHandlerImpl ReplaceHTTPRequestRuleHandlerImpl
)

type (
	ReplaceAllHTTPRequestRuleBackendHandlerImpl  ReplaceAllHTTPRequestRuleHandlerImpl
	ReplaceAllHTTPRequestRuleFrontendHandlerImpl ReplaceAllHTTPRequestRuleHandlerImpl
	ReplaceAllHTTPRequestRuleDefaultsHandlerImpl ReplaceAllHTTPRequestRuleHandlerImpl
)

func (h *CreateHTTPRequestRuleBackendHandlerImpl) Handle(params http_request_rule.CreateHTTPRequestRuleBackendParams, principal interface{}) middleware.Responder {
	g := CreateHTTPRequestRuleHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *CreateHTTPRequestRuleFrontendHandlerImpl) Handle(params http_request_rule.CreateHTTPRequestRuleFrontendParams, principal interface{}) middleware.Responder {
	g := CreateHTTPRequestRuleHandlerImpl(*h)
	pg := http_request_rule.CreateHTTPRequestRuleBackendParams(params)
	return g.Handle(cnconstants.FrontendParentType, pg, principal)
}

func (h *CreateHTTPRequestRuleDefaultsHandlerImpl) Handle(params http_request_rule.CreateHTTPRequestRuleDefaultsParams, principal interface{}) middleware.Responder {
	g := CreateHTTPRequestRuleHandlerImpl(*h)
	pg := http_request_rule.CreateHTTPRequestRuleBackendParams(params)
	return g.Handle(cnconstants.DefaultsParentType, pg, principal)
}

func (h *GetHTTPRequestRuleBackendHandlerImpl) Handle(params http_request_rule.GetHTTPRequestRuleBackendParams, principal interface{}) middleware.Responder {
	g := GetHTTPRequestRuleHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *GetHTTPRequestRuleFrontendHandlerImpl) Handle(params http_request_rule.GetHTTPRequestRuleFrontendParams, principal interface{}) middleware.Responder {
	g := GetHTTPRequestRuleHandlerImpl(*h)
	pg := http_request_rule.GetHTTPRequestRuleBackendParams(params)
	return g.Handle(cnconstants.FrontendParentType, pg, principal)
}

func (h *GetHTTPRequestRuleDefaultsHandlerImpl) Handle(params http_request_rule.GetHTTPRequestRuleDefaultsParams, principal interface{}) middleware.Responder {
	g := GetHTTPRequestRuleHandlerImpl(*h)
	pg := http_request_rule.GetHTTPRequestRuleBackendParams(params)
	return g.Handle(cnconstants.DefaultsParentType, pg, principal)
}

func (h *GetAllHTTPRequestRuleBackendHandlerImpl) Handle(params http_request_rule.GetAllHTTPRequestRuleBackendParams, principal interface{}) middleware.Responder {
	g := GetAllHTTPRequestRuleHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *GetAllHTTPRequestRuleFrontendHandlerImpl) Handle(params http_request_rule.GetAllHTTPRequestRuleFrontendParams, principal interface{}) middleware.Responder {
	g := GetAllHTTPRequestRuleHandlerImpl(*h)
	pg := http_request_rule.GetAllHTTPRequestRuleBackendParams(params)
	return g.Handle(cnconstants.FrontendParentType, pg, principal)
}

func (h *GetAllHTTPRequestRuleDefaultsHandlerImpl) Handle(params http_request_rule.GetAllHTTPRequestRuleDefaultsParams, principal interface{}) middleware.Responder {
	g := GetAllHTTPRequestRuleHandlerImpl(*h)
	pg := http_request_rule.GetAllHTTPRequestRuleBackendParams(params)
	return g.Handle(cnconstants.DefaultsParentType, pg, principal)
}

func (h *DeleteHTTPRequestRuleBackendHandlerImpl) Handle(params http_request_rule.DeleteHTTPRequestRuleBackendParams, principal interface{}) middleware.Responder {
	g := DeleteHTTPRequestRuleHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *DeleteHTTPRequestRuleFrontendHandlerImpl) Handle(params http_request_rule.DeleteHTTPRequestRuleFrontendParams, principal interface{}) middleware.Responder {
	g := DeleteHTTPRequestRuleHandlerImpl(*h)
	pg := http_request_rule.DeleteHTTPRequestRuleBackendParams(params)
	return g.Handle(cnconstants.FrontendParentType, pg, principal)
}

func (h *DeleteHTTPRequestRuleDefaultsHandlerImpl) Handle(params http_request_rule.DeleteHTTPRequestRuleDefaultsParams, principal interface{}) middleware.Responder {
	g := DeleteHTTPRequestRuleHandlerImpl(*h)
	pg := http_request_rule.DeleteHTTPRequestRuleBackendParams(params)
	return g.Handle(cnconstants.DefaultsParentType, pg, principal)
}

func (h *ReplaceHTTPRequestRuleBackendHandlerImpl) Handle(params http_request_rule.ReplaceHTTPRequestRuleBackendParams, principal interface{}) middleware.Responder {
	g := ReplaceHTTPRequestRuleHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *ReplaceHTTPRequestRuleFrontendHandlerImpl) Handle(params http_request_rule.ReplaceHTTPRequestRuleFrontendParams, principal interface{}) middleware.Responder {
	g := ReplaceHTTPRequestRuleHandlerImpl(*h)
	pg := http_request_rule.ReplaceHTTPRequestRuleBackendParams(params)
	return g.Handle(cnconstants.FrontendParentType, pg, principal)
}

func (h *ReplaceHTTPRequestRuleDefaultsHandlerImpl) Handle(params http_request_rule.ReplaceHTTPRequestRuleDefaultsParams, principal interface{}) middleware.Responder {
	g := ReplaceHTTPRequestRuleHandlerImpl(*h)
	pg := http_request_rule.ReplaceHTTPRequestRuleBackendParams(params)
	return g.Handle(cnconstants.DefaultsParentType, pg, principal)
}

func (h *ReplaceAllHTTPRequestRuleBackendHandlerImpl) Handle(params http_request_rule.ReplaceAllHTTPRequestRuleBackendParams, principal interface{}) middleware.Responder {
	g := ReplaceAllHTTPRequestRuleHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *ReplaceAllHTTPRequestRuleFrontendHandlerImpl) Handle(params http_request_rule.ReplaceAllHTTPRequestRuleFrontendParams, principal interface{}) middleware.Responder {
	g := ReplaceAllHTTPRequestRuleHandlerImpl(*h)
	pg := http_request_rule.ReplaceAllHTTPRequestRuleBackendParams(params)
	return g.Handle(cnconstants.FrontendParentType, pg, principal)
}

func (h *ReplaceAllHTTPRequestRuleDefaultsHandlerImpl) Handle(params http_request_rule.ReplaceAllHTTPRequestRuleDefaultsParams, principal interface{}) middleware.Responder {
	g := ReplaceAllHTTPRequestRuleHandlerImpl(*h)
	pg := http_request_rule.ReplaceAllHTTPRequestRuleBackendParams(params)
	return g.Handle(cnconstants.DefaultsParentType, pg, principal)
}
