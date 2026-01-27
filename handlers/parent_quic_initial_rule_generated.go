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
	"github.com/haproxytech/dataplaneapi/operations/quic_initial_rule"
)

type (
	CreateQUICInitialRuleFrontendHandlerImpl CreateQUICInitialRuleHandlerImpl
	CreateQUICInitialRuleDefaultsHandlerImpl CreateQUICInitialRuleHandlerImpl
)

type (
	GetQUICInitialRuleFrontendHandlerImpl GetQUICInitialRuleHandlerImpl
	GetQUICInitialRuleDefaultsHandlerImpl GetQUICInitialRuleHandlerImpl
)

type (
	GetAllQUICInitialRuleFrontendHandlerImpl GetAllQUICInitialRuleHandlerImpl
	GetAllQUICInitialRuleDefaultsHandlerImpl GetAllQUICInitialRuleHandlerImpl
)

type (
	DeleteQUICInitialRuleFrontendHandlerImpl DeleteQUICInitialRuleHandlerImpl
	DeleteQUICInitialRuleDefaultsHandlerImpl DeleteQUICInitialRuleHandlerImpl
)

type (
	ReplaceQUICInitialRuleFrontendHandlerImpl ReplaceQUICInitialRuleHandlerImpl
	ReplaceQUICInitialRuleDefaultsHandlerImpl ReplaceQUICInitialRuleHandlerImpl
)

type (
	ReplaceAllQUICInitialRuleFrontendHandlerImpl ReplaceAllQUICInitialRuleHandlerImpl
	ReplaceAllQUICInitialRuleDefaultsHandlerImpl ReplaceAllQUICInitialRuleHandlerImpl
)

func (h *CreateQUICInitialRuleFrontendHandlerImpl) Handle(params quic_initial_rule.CreateQUICInitialRuleFrontendParams, principal any) middleware.Responder {
	g := CreateQUICInitialRuleHandlerImpl(*h)
	return g.Handle(cnconstants.FrontendParentType, params, principal)
}

func (h *CreateQUICInitialRuleDefaultsHandlerImpl) Handle(params quic_initial_rule.CreateQUICInitialRuleDefaultsParams, principal any) middleware.Responder {
	g := CreateQUICInitialRuleHandlerImpl(*h)
	pg := quic_initial_rule.CreateQUICInitialRuleFrontendParams(params)
	return g.Handle(cnconstants.DefaultsParentType, pg, principal)
}

func (h *GetQUICInitialRuleFrontendHandlerImpl) Handle(params quic_initial_rule.GetQUICInitialRuleFrontendParams, principal any) middleware.Responder {
	g := GetQUICInitialRuleHandlerImpl(*h)
	return g.Handle(cnconstants.FrontendParentType, params, principal)
}

func (h *GetQUICInitialRuleDefaultsHandlerImpl) Handle(params quic_initial_rule.GetQUICInitialRuleDefaultsParams, principal any) middleware.Responder {
	g := GetQUICInitialRuleHandlerImpl(*h)
	pg := quic_initial_rule.GetQUICInitialRuleFrontendParams(params)
	return g.Handle(cnconstants.DefaultsParentType, pg, principal)
}

func (h *GetAllQUICInitialRuleFrontendHandlerImpl) Handle(params quic_initial_rule.GetAllQUICInitialRuleFrontendParams, principal any) middleware.Responder {
	g := GetAllQUICInitialRuleHandlerImpl(*h)
	return g.Handle(cnconstants.FrontendParentType, params, principal)
}

func (h *GetAllQUICInitialRuleDefaultsHandlerImpl) Handle(params quic_initial_rule.GetAllQUICInitialRuleDefaultsParams, principal any) middleware.Responder {
	g := GetAllQUICInitialRuleHandlerImpl(*h)
	pg := quic_initial_rule.GetAllQUICInitialRuleFrontendParams(params)
	return g.Handle(cnconstants.DefaultsParentType, pg, principal)
}

func (h *DeleteQUICInitialRuleFrontendHandlerImpl) Handle(params quic_initial_rule.DeleteQUICInitialRuleFrontendParams, principal any) middleware.Responder {
	g := DeleteQUICInitialRuleHandlerImpl(*h)
	return g.Handle(cnconstants.FrontendParentType, params, principal)
}

func (h *DeleteQUICInitialRuleDefaultsHandlerImpl) Handle(params quic_initial_rule.DeleteQUICInitialRuleDefaultsParams, principal any) middleware.Responder {
	g := DeleteQUICInitialRuleHandlerImpl(*h)
	pg := quic_initial_rule.DeleteQUICInitialRuleFrontendParams(params)
	return g.Handle(cnconstants.DefaultsParentType, pg, principal)
}

func (h *ReplaceQUICInitialRuleFrontendHandlerImpl) Handle(params quic_initial_rule.ReplaceQUICInitialRuleFrontendParams, principal any) middleware.Responder {
	g := ReplaceQUICInitialRuleHandlerImpl(*h)
	return g.Handle(cnconstants.FrontendParentType, params, principal)
}

func (h *ReplaceQUICInitialRuleDefaultsHandlerImpl) Handle(params quic_initial_rule.ReplaceQUICInitialRuleDefaultsParams, principal any) middleware.Responder {
	g := ReplaceQUICInitialRuleHandlerImpl(*h)
	pg := quic_initial_rule.ReplaceQUICInitialRuleFrontendParams(params)
	return g.Handle(cnconstants.DefaultsParentType, pg, principal)
}

func (h *ReplaceAllQUICInitialRuleFrontendHandlerImpl) Handle(params quic_initial_rule.ReplaceAllQUICInitialRuleFrontendParams, principal any) middleware.Responder {
	g := ReplaceAllQUICInitialRuleHandlerImpl(*h)
	return g.Handle(cnconstants.FrontendParentType, params, principal)
}

func (h *ReplaceAllQUICInitialRuleDefaultsHandlerImpl) Handle(params quic_initial_rule.ReplaceAllQUICInitialRuleDefaultsParams, principal any) middleware.Responder {
	g := ReplaceAllQUICInitialRuleHandlerImpl(*h)
	pg := quic_initial_rule.ReplaceAllQUICInitialRuleFrontendParams(params)
	return g.Handle(cnconstants.DefaultsParentType, pg, principal)
}
