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
	"github.com/haproxytech/dataplaneapi/operations/tcp_response_rule"
)

// CreateTCPResponseRuleHandlerImpl implementation of the CreateTCPResponseRuleHandler interface using client-native client
type CreateTCPResponseRuleHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// DeleteTCPResponseRuleHandlerImpl implementation of the DeleteTCPResponseRuleHandler interface using client-native client
type DeleteTCPResponseRuleHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// GetTCPResponseRuleHandlerImpl implementation of the GetTCPResponseRuleHandler interface using client-native client
type GetTCPResponseRuleHandlerImpl struct {
	Client client_native.HAProxyClient
}

// GetAllTCPResponseRuleHandlerImpl implementation of the GetTCPResponseRulesHandler interface using client-native client
type GetAllTCPResponseRuleHandlerImpl struct {
	Client client_native.HAProxyClient
}

// ReplaceTCPResponseRuleHandlerImpl implementation of the ReplaceTCPResponseRuleHandler interface using client-native client
type ReplaceTCPResponseRuleHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// ReplaceAllTCPResponseRuleHandlerImpl implementation of the ReplaceTCPResponseRulesHandler interface using client-native client
type ReplaceAllTCPResponseRuleHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// Handle executing the request and returning a response
func (h *CreateTCPResponseRuleHandlerImpl) Handle(parentType cnconstants.CnParentType, params tcp_response_rule.CreateTCPResponseRuleBackendParams, principal interface{}) middleware.Responder {
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
		return tcp_response_rule.NewCreateTCPResponseRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return tcp_response_rule.NewCreateTCPResponseRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.CreateTCPResponseRule(params.Index, string(parentType), params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return tcp_response_rule.NewCreateTCPResponseRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return tcp_response_rule.NewCreateTCPResponseRuleBackendDefault(int(*e.Code)).WithPayload(e)
			}
			return tcp_response_rule.NewCreateTCPResponseRuleBackendCreated().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return tcp_response_rule.NewCreateTCPResponseRuleBackendAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return tcp_response_rule.NewCreateTCPResponseRuleBackendAccepted().WithPayload(params.Data)
}

// Handle executing the request and returning a response
func (h *DeleteTCPResponseRuleHandlerImpl) Handle(parentType cnconstants.CnParentType, params tcp_response_rule.DeleteTCPResponseRuleBackendParams, principal interface{}) middleware.Responder {
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
		return tcp_response_rule.NewDeleteTCPResponseRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return tcp_response_rule.NewDeleteTCPResponseRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.DeleteTCPResponseRule(params.Index, string(parentType), params.ParentName, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return tcp_response_rule.NewDeleteTCPResponseRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return tcp_response_rule.NewDeleteTCPResponseRuleBackendDefault(int(*e.Code)).WithPayload(e)
			}
			return tcp_response_rule.NewDeleteTCPResponseRuleBackendNoContent()
		}
		rID := h.ReloadAgent.Reload()
		return tcp_response_rule.NewDeleteTCPResponseRuleBackendAccepted().WithReloadID(rID)
	}
	return tcp_response_rule.NewDeleteTCPResponseRuleBackendAccepted()
}

// Handle executing the request and returning a response
func (h *GetTCPResponseRuleHandlerImpl) Handle(parentType cnconstants.CnParentType, params tcp_response_rule.GetTCPResponseRuleBackendParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return tcp_response_rule.NewGetTCPResponseRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}
	_, rule, err := configuration.GetTCPResponseRule(params.Index, string(parentType), params.ParentName, t)
	if err != nil {
		e := misc.HandleError(err)
		return tcp_response_rule.NewGetTCPResponseRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}
	return tcp_response_rule.NewGetTCPResponseRuleBackendOK().WithPayload(rule)
}

// Handle executing the request and returning a response
func (h *GetAllTCPResponseRuleHandlerImpl) Handle(parentType cnconstants.CnParentType, params tcp_response_rule.GetAllTCPResponseRuleBackendParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return tcp_response_rule.NewGetAllTCPResponseRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	_, rules, err := configuration.GetTCPResponseRules(string(parentType), params.ParentName, t)
	if err != nil {
		e := misc.HandleContainerGetError(err)
		if *e.Code == misc.ErrHTTPOk {
			return tcp_response_rule.NewGetAllTCPResponseRuleBackendOK().WithPayload(models.TCPResponseRules{})
		}
		return tcp_response_rule.NewGetAllTCPResponseRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}
	return tcp_response_rule.NewGetAllTCPResponseRuleBackendOK().WithPayload(rules)
}

// Handle executing the request and returning a response
func (h *ReplaceTCPResponseRuleHandlerImpl) Handle(parentType cnconstants.CnParentType, params tcp_response_rule.ReplaceTCPResponseRuleBackendParams, principal interface{}) middleware.Responder {
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
		return tcp_response_rule.NewReplaceTCPResponseRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return tcp_response_rule.NewReplaceTCPResponseRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.EditTCPResponseRule(params.Index, string(parentType), params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return tcp_response_rule.NewReplaceTCPResponseRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return tcp_response_rule.NewReplaceTCPResponseRuleBackendDefault(int(*e.Code)).WithPayload(e)
			}
			return tcp_response_rule.NewReplaceTCPResponseRuleBackendOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return tcp_response_rule.NewReplaceTCPResponseRuleBackendAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return tcp_response_rule.NewReplaceTCPResponseRuleBackendAccepted().WithPayload(params.Data)
}

func (h *ReplaceAllTCPResponseRuleHandlerImpl) Handle(parentType cnconstants.CnParentType, params tcp_response_rule.ReplaceAllTCPResponseRuleBackendParams, principal interface{}) middleware.Responder {
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
		return tcp_response_rule.NewReplaceAllTCPResponseRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return tcp_response_rule.NewReplaceAllTCPResponseRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}
	err = configuration.ReplaceTCPResponseRules(string(parentType), params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return tcp_response_rule.NewReplaceAllTCPResponseRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return tcp_response_rule.NewReplaceAllTCPResponseRuleBackendDefault(int(*e.Code)).WithPayload(e)
			}
			return tcp_response_rule.NewReplaceAllTCPResponseRuleBackendOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return tcp_response_rule.NewReplaceAllTCPResponseRuleBackendAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return tcp_response_rule.NewReplaceAllTCPResponseRuleBackendAccepted().WithPayload(params.Data)
}
