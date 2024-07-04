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
	"github.com/haproxytech/dataplaneapi/operations/http_response_rule"
)

// CreateHTTPResponseRuleHandlerImpl implementation of the CreateHTTPResponseRuleHandler interface using client-native client
type CreateHTTPResponseRuleHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// DeleteHTTPResponseRuleHandlerImpl implementation of the DeleteHTTPResponseRuleHandler interface using client-native client
type DeleteHTTPResponseRuleHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// GetHTTPResponseRuleHandlerImpl implementation of the GetHTTPResponseRuleHandler interface using client-native client
type GetHTTPResponseRuleHandlerImpl struct {
	Client client_native.HAProxyClient
}

// GetAllHTTPResponseRuleHandlerImpl implementation of the GetHTTPResponseRulesHandler interface using client-native client
type GetAllHTTPResponseRuleHandlerImpl struct {
	Client client_native.HAProxyClient
}

// ReplaceHTTPResponseRuleHandlerImpl implementation of the ReplaceHTTPResponseRuleHandler interface using client-native client
type ReplaceHTTPResponseRuleHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// ReplaceAllHTTPResponseRuleHandlerImpl implementation of the ReplaceHTTPResponseRulesHandler interface using client-native client
type ReplaceAllHTTPResponseRuleHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// Handle executing the request and returning a response
func (h *CreateHTTPResponseRuleHandlerImpl) Handle(parentType cnconstants.CnParentType, params http_response_rule.CreateHTTPResponseRuleBackendParams, principal interface{}) middleware.Responder {
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
		return http_response_rule.NewCreateHTTPResponseRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_response_rule.NewCreateHTTPResponseRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.CreateHTTPResponseRule(params.Index, string(parentType), params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return http_response_rule.NewCreateHTTPResponseRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return http_response_rule.NewCreateHTTPResponseRuleBackendDefault(int(*e.Code)).WithPayload(e)
			}
			return http_response_rule.NewCreateHTTPResponseRuleBackendCreated().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return http_response_rule.NewCreateHTTPResponseRuleBackendAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return http_response_rule.NewCreateHTTPResponseRuleBackendAccepted().WithPayload(params.Data)
}

// Handle executing the request and returning a response
func (h *DeleteHTTPResponseRuleHandlerImpl) Handle(parentType cnconstants.CnParentType, params http_response_rule.DeleteHTTPResponseRuleBackendParams, principal interface{}) middleware.Responder {
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
		return http_response_rule.NewDeleteHTTPResponseRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_response_rule.NewDeleteHTTPResponseRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.DeleteHTTPResponseRule(params.Index, string(parentType), params.ParentName, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return http_response_rule.NewDeleteHTTPResponseRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return http_response_rule.NewDeleteHTTPResponseRuleBackendDefault(int(*e.Code)).WithPayload(e)
			}
			return http_response_rule.NewDeleteHTTPResponseRuleBackendNoContent()
		}
		rID := h.ReloadAgent.Reload()
		return http_response_rule.NewDeleteHTTPResponseRuleBackendAccepted().WithReloadID(rID)
	}
	return http_response_rule.NewDeleteHTTPResponseRuleBackendAccepted()
}

// Handle executing the request and returning a response
func (h *GetHTTPResponseRuleHandlerImpl) Handle(parentType cnconstants.CnParentType, params http_response_rule.GetHTTPResponseRuleBackendParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_response_rule.NewCreateHTTPResponseRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	_, rule, err := configuration.GetHTTPResponseRule(params.Index, string(parentType), params.ParentName, t)
	if err != nil {
		e := misc.HandleError(err)
		return http_response_rule.NewGetHTTPResponseRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}
	return http_response_rule.NewGetHTTPResponseRuleBackendOK().WithPayload(rule)
}

// Handle executing the request and returning a response
func (h *GetAllHTTPResponseRuleHandlerImpl) Handle(parentType cnconstants.CnParentType, params http_response_rule.GetAllHTTPResponseRuleBackendParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_response_rule.NewGetHTTPResponseRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	_, rules, err := configuration.GetHTTPResponseRules(string(parentType), params.ParentName, t)
	if err != nil {
		e := misc.HandleContainerGetError(err)
		if *e.Code == misc.ErrHTTPOk {
			return http_response_rule.NewGetAllHTTPResponseRuleBackendOK().WithPayload(models.HTTPResponseRules{})
		}
		return http_response_rule.NewGetAllHTTPResponseRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}
	return http_response_rule.NewGetAllHTTPResponseRuleBackendOK().WithPayload(rules)
}

// Handle executing the request and returning a response
func (h *ReplaceHTTPResponseRuleHandlerImpl) Handle(parentType cnconstants.CnParentType, params http_response_rule.ReplaceHTTPResponseRuleBackendParams, principal interface{}) middleware.Responder {
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
		return http_response_rule.NewReplaceHTTPResponseRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_response_rule.NewReplaceHTTPResponseRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.EditHTTPResponseRule(params.Index, string(parentType), params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return http_response_rule.NewReplaceHTTPResponseRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return http_response_rule.NewReplaceHTTPResponseRuleBackendDefault(int(*e.Code)).WithPayload(e)
			}
			return http_response_rule.NewReplaceHTTPResponseRuleBackendOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return http_response_rule.NewReplaceHTTPResponseRuleBackendAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return http_response_rule.NewReplaceHTTPResponseRuleBackendAccepted().WithPayload(params.Data)
}

// Handle executing the request and returning a response
func (h *ReplaceAllHTTPResponseRuleHandlerImpl) Handle(parentType cnconstants.CnParentType, params http_response_rule.ReplaceAllHTTPResponseRuleBackendParams, principal interface{}) middleware.Responder {
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
		return http_response_rule.NewReplaceAllHTTPResponseRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_response_rule.NewReplaceAllHTTPResponseRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}
	err = configuration.ReplaceHTTPResponseRules(string(parentType), params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return http_response_rule.NewReplaceAllHTTPResponseRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return http_response_rule.NewReplaceAllHTTPResponseRuleBackendDefault(int(*e.Code)).WithPayload(e)
			}
			return http_response_rule.NewReplaceAllHTTPResponseRuleBackendOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return http_response_rule.NewReplaceAllHTTPResponseRuleBackendAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return http_response_rule.NewReplaceAllHTTPResponseRuleBackendAccepted().WithPayload(params.Data)
}
