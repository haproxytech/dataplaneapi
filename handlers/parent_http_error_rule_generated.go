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
	"github.com/haproxytech/dataplaneapi/operations/http_error_rule"
)

type (
	CreateHTTPErrorRuleBackendHandlerImpl  CreateHTTPErrorRuleHandlerImpl
	CreateHTTPErrorRuleFrontendHandlerImpl CreateHTTPErrorRuleHandlerImpl
	CreateHTTPErrorRuleDefaultsHandlerImpl CreateHTTPErrorRuleHandlerImpl
)

type (
	GetHTTPErrorRuleBackendHandlerImpl  GetHTTPErrorRuleHandlerImpl
	GetHTTPErrorRuleFrontendHandlerImpl GetHTTPErrorRuleHandlerImpl
	GetHTTPErrorRuleDefaultsHandlerImpl GetHTTPErrorRuleHandlerImpl
)

type (
	GetAllHTTPErrorRuleBackendHandlerImpl  GetAllHTTPErrorRuleHandlerImpl
	GetAllHTTPErrorRuleFrontendHandlerImpl GetAllHTTPErrorRuleHandlerImpl
	GetAllHTTPErrorRuleDefaultsHandlerImpl GetAllHTTPErrorRuleHandlerImpl
)

type (
	DeleteHTTPErrorRuleBackendHandlerImpl  DeleteHTTPErrorRuleHandlerImpl
	DeleteHTTPErrorRuleFrontendHandlerImpl DeleteHTTPErrorRuleHandlerImpl
	DeleteHTTPErrorRuleDefaultsHandlerImpl DeleteHTTPErrorRuleHandlerImpl
)

type (
	ReplaceHTTPErrorRuleBackendHandlerImpl  ReplaceHTTPErrorRuleHandlerImpl
	ReplaceHTTPErrorRuleFrontendHandlerImpl ReplaceHTTPErrorRuleHandlerImpl
	ReplaceHTTPErrorRuleDefaultsHandlerImpl ReplaceHTTPErrorRuleHandlerImpl
)

type (
	ReplaceAllHTTPErrorRuleBackendHandlerImpl  ReplaceAllHTTPErrorRuleHandlerImpl
	ReplaceAllHTTPErrorRuleFrontendHandlerImpl ReplaceAllHTTPErrorRuleHandlerImpl
	ReplaceAllHTTPErrorRuleDefaultsHandlerImpl ReplaceAllHTTPErrorRuleHandlerImpl
)

func (h *CreateHTTPErrorRuleBackendHandlerImpl) Handle(params http_error_rule.CreateHTTPErrorRuleBackendParams, principal interface{}) middleware.Responder {
	g := CreateHTTPErrorRuleHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *CreateHTTPErrorRuleFrontendHandlerImpl) Handle(params http_error_rule.CreateHTTPErrorRuleFrontendParams, principal interface{}) middleware.Responder {
	g := CreateHTTPErrorRuleHandlerImpl(*h)
	pg := http_error_rule.CreateHTTPErrorRuleBackendParams(params)
	return g.Handle(cnconstants.FrontendParentType, pg, principal)
}

func (h *CreateHTTPErrorRuleDefaultsHandlerImpl) Handle(params http_error_rule.CreateHTTPErrorRuleDefaultsParams, principal interface{}) middleware.Responder {
	g := CreateHTTPErrorRuleHandlerImpl(*h)
	pg := http_error_rule.CreateHTTPErrorRuleBackendParams(params)
	return g.Handle(cnconstants.DefaultsParentType, pg, principal)
}

func (h *GetHTTPErrorRuleBackendHandlerImpl) Handle(params http_error_rule.GetHTTPErrorRuleBackendParams, principal interface{}) middleware.Responder {
	g := GetHTTPErrorRuleHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *GetHTTPErrorRuleFrontendHandlerImpl) Handle(params http_error_rule.GetHTTPErrorRuleFrontendParams, principal interface{}) middleware.Responder {
	g := GetHTTPErrorRuleHandlerImpl(*h)
	pg := http_error_rule.GetHTTPErrorRuleBackendParams(params)
	return g.Handle(cnconstants.FrontendParentType, pg, principal)
}

func (h *GetHTTPErrorRuleDefaultsHandlerImpl) Handle(params http_error_rule.GetHTTPErrorRuleDefaultsParams, principal interface{}) middleware.Responder {
	g := GetHTTPErrorRuleHandlerImpl(*h)
	pg := http_error_rule.GetHTTPErrorRuleBackendParams(params)
	return g.Handle(cnconstants.DefaultsParentType, pg, principal)
}

func (h *GetAllHTTPErrorRuleBackendHandlerImpl) Handle(params http_error_rule.GetAllHTTPErrorRuleBackendParams, principal interface{}) middleware.Responder {
	g := GetAllHTTPErrorRuleHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *GetAllHTTPErrorRuleFrontendHandlerImpl) Handle(params http_error_rule.GetAllHTTPErrorRuleFrontendParams, principal interface{}) middleware.Responder {
	g := GetAllHTTPErrorRuleHandlerImpl(*h)
	pg := http_error_rule.GetAllHTTPErrorRuleBackendParams(params)
	return g.Handle(cnconstants.FrontendParentType, pg, principal)
}

func (h *GetAllHTTPErrorRuleDefaultsHandlerImpl) Handle(params http_error_rule.GetAllHTTPErrorRuleDefaultsParams, principal interface{}) middleware.Responder {
	g := GetAllHTTPErrorRuleHandlerImpl(*h)
	pg := http_error_rule.GetAllHTTPErrorRuleBackendParams(params)
	return g.Handle(cnconstants.DefaultsParentType, pg, principal)
}

func (h *DeleteHTTPErrorRuleBackendHandlerImpl) Handle(params http_error_rule.DeleteHTTPErrorRuleBackendParams, principal interface{}) middleware.Responder {
	g := DeleteHTTPErrorRuleHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *DeleteHTTPErrorRuleFrontendHandlerImpl) Handle(params http_error_rule.DeleteHTTPErrorRuleFrontendParams, principal interface{}) middleware.Responder {
	g := DeleteHTTPErrorRuleHandlerImpl(*h)
	pg := http_error_rule.DeleteHTTPErrorRuleBackendParams(params)
	return g.Handle(cnconstants.FrontendParentType, pg, principal)
}

func (h *DeleteHTTPErrorRuleDefaultsHandlerImpl) Handle(params http_error_rule.DeleteHTTPErrorRuleDefaultsParams, principal interface{}) middleware.Responder {
	g := DeleteHTTPErrorRuleHandlerImpl(*h)
	pg := http_error_rule.DeleteHTTPErrorRuleBackendParams(params)
	return g.Handle(cnconstants.DefaultsParentType, pg, principal)
}

func (h *ReplaceHTTPErrorRuleBackendHandlerImpl) Handle(params http_error_rule.ReplaceHTTPErrorRuleBackendParams, principal interface{}) middleware.Responder {
	g := ReplaceHTTPErrorRuleHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *ReplaceHTTPErrorRuleFrontendHandlerImpl) Handle(params http_error_rule.ReplaceHTTPErrorRuleFrontendParams, principal interface{}) middleware.Responder {
	g := ReplaceHTTPErrorRuleHandlerImpl(*h)
	pg := http_error_rule.ReplaceHTTPErrorRuleBackendParams(params)
	return g.Handle(cnconstants.FrontendParentType, pg, principal)
}

func (h *ReplaceHTTPErrorRuleDefaultsHandlerImpl) Handle(params http_error_rule.ReplaceHTTPErrorRuleDefaultsParams, principal interface{}) middleware.Responder {
	g := ReplaceHTTPErrorRuleHandlerImpl(*h)
	pg := http_error_rule.ReplaceHTTPErrorRuleBackendParams(params)
	return g.Handle(cnconstants.DefaultsParentType, pg, principal)
}

func (h *ReplaceAllHTTPErrorRuleBackendHandlerImpl) Handle(params http_error_rule.ReplaceAllHTTPErrorRuleBackendParams, principal interface{}) middleware.Responder {
	g := ReplaceAllHTTPErrorRuleHandlerImpl(*h)
	return g.Handle(cnconstants.BackendParentType, params, principal)
}

func (h *ReplaceAllHTTPErrorRuleFrontendHandlerImpl) Handle(params http_error_rule.ReplaceAllHTTPErrorRuleFrontendParams, principal interface{}) middleware.Responder {
	g := ReplaceAllHTTPErrorRuleHandlerImpl(*h)
	pg := http_error_rule.ReplaceAllHTTPErrorRuleBackendParams(params)
	return g.Handle(cnconstants.FrontendParentType, pg, principal)
}

func (h *ReplaceAllHTTPErrorRuleDefaultsHandlerImpl) Handle(params http_error_rule.ReplaceAllHTTPErrorRuleDefaultsParams, principal interface{}) middleware.Responder {
	g := ReplaceAllHTTPErrorRuleHandlerImpl(*h)
	pg := http_error_rule.ReplaceAllHTTPErrorRuleBackendParams(params)
	return g.Handle(cnconstants.DefaultsParentType, pg, principal)
}
