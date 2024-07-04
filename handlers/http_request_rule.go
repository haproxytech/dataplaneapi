// Copyright 2019 HAProxy Technologies
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	client_native "github.com/haproxytech/client-native/v6"
	cnconstants "github.com/haproxytech/client-native/v6/configuration/parents"
	"github.com/haproxytech/client-native/v6/models"

	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/operations/http_request_rule"
)

// CreateHTTPRequestRuleHandlerImpl implementation of the CreateHTTPRequestRuleHandler interface using client-native client
type CreateHTTPRequestRuleHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// DeleteHTTPRequestRuleHandlerImpl implementation of the DeleteHTTPRequestRuleHandler interface using client-native client
type DeleteHTTPRequestRuleHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// GetHTTPRequestRuleHandlerImpl implementation of the GetHTTPRequestRuleHandler interface using client-native client
type GetHTTPRequestRuleHandlerImpl struct {
	Client client_native.HAProxyClient
}

// GetAllHTTPRequestRuleHandlerImpl implementation of the GetHTTPRequestRulesHandler interface using client-native client
type GetAllHTTPRequestRuleHandlerImpl struct {
	Client client_native.HAProxyClient
}

// ReplaceHTTPRequestRuleHandlerImpl implementation of the ReplaceHTTPRequestRuleHandler interface using client-native client
type ReplaceHTTPRequestRuleHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// ReplaceAllHTTPRequestRuleHandlerImpl implementation of the ReplaceHTTPRequestRulesHandler interface using client-native client
type ReplaceAllHTTPRequestRuleHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// Handle executing the request and returning a response
func (h *CreateHTTPRequestRuleHandlerImpl) Handle(parentType cnconstants.CnParentType, params http_request_rule.CreateHTTPRequestRuleBackendParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	if t != "" && *params.ForceReload {
		msg := "Both force_reload and transaction specified, specify only one"
		c := misc.ErrHTTPBadRequest
		e := &models.Error{
			Message: &msg,
			Code:    &c,
		}
		return http_request_rule.NewCreateHTTPRequestRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_request_rule.NewCreateHTTPRequestRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.CreateHTTPRequestRule(params.Index, string(parentType), params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return http_request_rule.NewCreateHTTPRequestRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return http_request_rule.NewCreateHTTPRequestRuleBackendDefault(int(*e.Code)).WithPayload(e)
			}
			return http_request_rule.NewCreateHTTPRequestRuleBackendCreated().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return http_request_rule.NewCreateHTTPRequestRuleBackendAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return http_request_rule.NewCreateHTTPRequestRuleBackendAccepted().WithPayload(params.Data)
}

// Handle executing the request and returning a response
func (h *DeleteHTTPRequestRuleHandlerImpl) Handle(parentType cnconstants.CnParentType, params http_request_rule.DeleteHTTPRequestRuleBackendParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	if t != "" && *params.ForceReload {
		msg := "Both force_reload and transaction specified, specify only one"
		c := misc.ErrHTTPBadRequest
		e := &models.Error{
			Message: &msg,
			Code:    &c,
		}
		return http_request_rule.NewDeleteHTTPRequestRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_request_rule.NewDeleteHTTPRequestRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.DeleteHTTPRequestRule(params.Index, string(parentType), params.ParentName, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return http_request_rule.NewDeleteHTTPRequestRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return http_request_rule.NewDeleteHTTPRequestRuleBackendDefault(int(*e.Code)).WithPayload(e)
			}
			return http_request_rule.NewDeleteHTTPRequestRuleBackendNoContent()
		}
		rID := h.ReloadAgent.Reload()
		return http_request_rule.NewDeleteHTTPRequestRuleBackendAccepted().WithReloadID(rID)
	}
	return http_request_rule.NewDeleteHTTPRequestRuleBackendAccepted()
}

// Handle executing the request and returning a response
func (h *GetHTTPRequestRuleHandlerImpl) Handle(parentType cnconstants.CnParentType, params http_request_rule.GetHTTPRequestRuleBackendParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_request_rule.NewGetHTTPRequestRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	_, rule, err := configuration.GetHTTPRequestRule(params.Index, string(parentType), params.ParentName, t)
	if err != nil {
		e := misc.HandleError(err)
		return http_request_rule.NewGetHTTPRequestRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}
	return http_request_rule.NewGetHTTPRequestRuleBackendOK().WithPayload(rule)
}

// Handle executing the request and returning a response
func (h *GetAllHTTPRequestRuleHandlerImpl) Handle(parentType cnconstants.CnParentType, params http_request_rule.GetAllHTTPRequestRuleBackendParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_request_rule.NewGetHTTPRequestRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	_, rules, err := configuration.GetHTTPRequestRules(string(parentType), params.ParentName, t)
	if err != nil {
		e := misc.HandleContainerGetError(err)
		if *e.Code == misc.ErrHTTPOk {
			return http_request_rule.NewGetAllHTTPRequestRuleBackendOK().WithPayload(models.HTTPRequestRules{})
		}
		return http_request_rule.NewGetAllHTTPRequestRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}
	return http_request_rule.NewGetAllHTTPRequestRuleBackendOK().WithPayload(rules)
}

// Handle executing the request and returning a response
func (h *ReplaceHTTPRequestRuleHandlerImpl) Handle(parentType cnconstants.CnParentType, params http_request_rule.ReplaceHTTPRequestRuleBackendParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	if t != "" && *params.ForceReload {
		msg := "Both force_reload and transaction specified, specify only one"
		c := misc.ErrHTTPBadRequest
		e := &models.Error{
			Message: &msg,
			Code:    &c,
		}
		return http_request_rule.NewReplaceHTTPRequestRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_request_rule.NewReplaceHTTPRequestRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.EditHTTPRequestRule(params.Index, string(parentType), params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return http_request_rule.NewReplaceHTTPRequestRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return http_request_rule.NewReplaceHTTPRequestRuleBackendDefault(int(*e.Code)).WithPayload(e)
			}
			return http_request_rule.NewReplaceHTTPRequestRuleBackendOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return http_request_rule.NewReplaceHTTPRequestRuleBackendAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return http_request_rule.NewReplaceHTTPRequestRuleBackendAccepted().WithPayload(params.Data)
}

// Handle executing the request and returning a response
func (h *ReplaceAllHTTPRequestRuleHandlerImpl) Handle(parentType cnconstants.CnParentType, params http_request_rule.ReplaceAllHTTPRequestRuleBackendParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	if t != "" && *params.ForceReload {
		msg := "Both force_reload and transaction specified, specify only one"
		c := misc.ErrHTTPBadRequest
		e := &models.Error{
			Message: &msg,
			Code:    &c,
		}
		return http_request_rule.NewReplaceAllHTTPRequestRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_request_rule.NewReplaceAllHTTPRequestRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}
	err = configuration.ReplaceHTTPRequestRules(string(parentType), params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return http_request_rule.NewReplaceAllHTTPRequestRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return http_request_rule.NewReplaceAllHTTPRequestRuleBackendDefault(int(*e.Code)).WithPayload(e)
			}
			return http_request_rule.NewReplaceAllHTTPRequestRuleBackendOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return http_request_rule.NewReplaceAllHTTPRequestRuleBackendAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return http_request_rule.NewReplaceAllHTTPRequestRuleBackendAccepted().WithPayload(params.Data)
}
