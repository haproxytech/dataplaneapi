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
	"github.com/haproxytech/dataplaneapi/operations/tcp_request_rule"
)

// CreateTCPRequestRuleHandlerImpl implementation of the CreateTCPRequestRuleHandler interface using client-native client
type CreateTCPRequestRuleHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// DeleteTCPRequestRuleHandlerImpl implementation of the DeleteTCPRequestRuleHandler interface using client-native client
type DeleteTCPRequestRuleHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// GetTCPRequestRuleHandlerImpl implementation of the GetTCPRequestRuleHandler interface using client-native client
type GetTCPRequestRuleHandlerImpl struct {
	Client client_native.HAProxyClient
}

// GetAllTCPRequestRuleHandlerImpl implementation of the GetTCPRequestRulesHandler interface using client-native client
type GetAllTCPRequestRuleHandlerImpl struct {
	Client client_native.HAProxyClient
}

// ReplaceTCPRequestRuleHandlerImpl implementation of the ReplaceTCPRequestRuleHandler interface using client-native client
type ReplaceTCPRequestRuleHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// ReplaceAllTCPRequestRuleHandlerImpl implementation of the ReplaceTCPRequestRulesHandler interface using client-native client
type ReplaceAllTCPRequestRuleHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// Handle executing the request and returning a response
func (h *CreateTCPRequestRuleHandlerImpl) Handle(parentType cnconstants.CnParentType, params tcp_request_rule.CreateTCPRequestRuleBackendParams, principal interface{}) middleware.Responder {
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
		return tcp_request_rule.NewCreateTCPRequestRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return tcp_request_rule.NewCreateTCPRequestRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.CreateTCPRequestRule(params.Index, string(parentType), params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return tcp_request_rule.NewCreateTCPRequestRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return tcp_request_rule.NewCreateTCPRequestRuleBackendDefault(int(*e.Code)).WithPayload(e)
			}
			return tcp_request_rule.NewCreateTCPRequestRuleBackendCreated().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return tcp_request_rule.NewCreateTCPRequestRuleBackendAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return tcp_request_rule.NewCreateTCPRequestRuleBackendAccepted().WithPayload(params.Data)
}

// Handle executing the request and returning a response
func (h *DeleteTCPRequestRuleHandlerImpl) Handle(parentType cnconstants.CnParentType, params tcp_request_rule.DeleteTCPRequestRuleBackendParams, principal interface{}) middleware.Responder {
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
		return tcp_request_rule.NewDeleteTCPRequestRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return tcp_request_rule.NewDeleteTCPRequestRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.DeleteTCPRequestRule(params.Index, string(parentType), params.ParentName, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return tcp_request_rule.NewDeleteTCPRequestRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return tcp_request_rule.NewDeleteTCPRequestRuleBackendDefault(int(*e.Code)).WithPayload(e)
			}
			return tcp_request_rule.NewDeleteTCPRequestRuleBackendNoContent()
		}
		rID := h.ReloadAgent.Reload()
		return tcp_request_rule.NewDeleteTCPRequestRuleBackendAccepted().WithReloadID(rID)
	}
	return tcp_request_rule.NewDeleteTCPRequestRuleBackendAccepted()
}

// Handle executing the request and returning a response
func (h *GetTCPRequestRuleHandlerImpl) Handle(parentType cnconstants.CnParentType, params tcp_request_rule.GetTCPRequestRuleBackendParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return tcp_request_rule.NewGetTCPRequestRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	_, rule, err := configuration.GetTCPRequestRule(params.Index, string(parentType), params.ParentName, t)
	if err != nil {
		e := misc.HandleError(err)
		return tcp_request_rule.NewGetTCPRequestRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}
	return tcp_request_rule.NewGetTCPRequestRuleBackendOK().WithPayload(rule)
}

// Handle executing the request and returning a response
func (h *GetAllTCPRequestRuleHandlerImpl) Handle(parentType cnconstants.CnParentType, params tcp_request_rule.GetAllTCPRequestRuleBackendParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return tcp_request_rule.NewGetAllTCPRequestRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	_, rules, err := configuration.GetTCPRequestRules(string(parentType), params.ParentName, t)
	if err != nil {
		e := misc.HandleContainerGetError(err)
		if *e.Code == misc.ErrHTTPOk {
			return tcp_request_rule.NewGetAllTCPRequestRuleBackendOK().WithPayload(models.TCPRequestRules{})
		}
		return tcp_request_rule.NewGetAllTCPRequestRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}
	return tcp_request_rule.NewGetAllTCPRequestRuleBackendOK().WithPayload(rules)
}

// Handle executing the request and returning a response
func (h *ReplaceTCPRequestRuleHandlerImpl) Handle(parentType cnconstants.CnParentType, params tcp_request_rule.ReplaceTCPRequestRuleBackendParams, principal interface{}) middleware.Responder {
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
		return tcp_request_rule.NewReplaceTCPRequestRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return tcp_request_rule.NewReplaceTCPRequestRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.EditTCPRequestRule(params.Index, string(parentType), params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return tcp_request_rule.NewReplaceTCPRequestRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return tcp_request_rule.NewReplaceTCPRequestRuleBackendDefault(int(*e.Code)).WithPayload(e)
			}
			return tcp_request_rule.NewReplaceTCPRequestRuleBackendOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return tcp_request_rule.NewReplaceTCPRequestRuleBackendAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return tcp_request_rule.NewReplaceTCPRequestRuleBackendAccepted().WithPayload(params.Data)
}

// Handle executing the request and returning a response
func (h *ReplaceAllTCPRequestRuleHandlerImpl) Handle(parentType cnconstants.CnParentType, params tcp_request_rule.ReplaceAllTCPRequestRuleBackendParams, principal interface{}) middleware.Responder {
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
		return tcp_request_rule.NewReplaceAllTCPRequestRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return tcp_request_rule.NewReplaceAllTCPRequestRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}
	err = configuration.ReplaceTCPRequestRules(string(parentType), params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return tcp_request_rule.NewReplaceAllTCPRequestRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return tcp_request_rule.NewReplaceAllTCPRequestRuleBackendDefault(int(*e.Code)).WithPayload(e)
			}
			return tcp_request_rule.NewReplaceAllTCPRequestRuleBackendOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return tcp_request_rule.NewReplaceAllTCPRequestRuleBackendAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return tcp_request_rule.NewReplaceAllTCPRequestRuleBackendAccepted().WithPayload(params.Data)
}
