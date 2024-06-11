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

// GetTCPRequestRulesHandlerImpl implementation of the GetTCPRequestRulesHandler interface using client-native client
type GetTCPRequestRulesHandlerImpl struct {
	Client client_native.HAProxyClient
}

// ReplaceTCPRequestRuleHandlerImpl implementation of the ReplaceTCPRequestRuleHandler interface using client-native client
type ReplaceTCPRequestRuleHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// ReplaceTCPRequestRulesHandlerImpl implementation of the ReplaceTCPRequestRulesHandler interface using client-native client
type ReplaceTCPRequestRulesHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// Handle executing the request and returning a response
func (h *CreateTCPRequestRuleHandlerImpl) Handle(params tcp_request_rule.CreateTCPRequestRuleParams, principal interface{}) middleware.Responder {
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
		return tcp_request_rule.NewCreateTCPRequestRuleDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return tcp_request_rule.NewCreateTCPRequestRuleDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.CreateTCPRequestRule(params.Index, params.ParentType, params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return tcp_request_rule.NewCreateTCPRequestRuleDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return tcp_request_rule.NewCreateTCPRequestRuleDefault(int(*e.Code)).WithPayload(e)
			}
			return tcp_request_rule.NewCreateTCPRequestRuleCreated().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return tcp_request_rule.NewCreateTCPRequestRuleAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return tcp_request_rule.NewCreateTCPRequestRuleAccepted().WithPayload(params.Data)
}

// Handle executing the request and returning a response
func (h *DeleteTCPRequestRuleHandlerImpl) Handle(params tcp_request_rule.DeleteTCPRequestRuleParams, principal interface{}) middleware.Responder {
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
		return tcp_request_rule.NewDeleteTCPRequestRuleDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return tcp_request_rule.NewDeleteTCPRequestRuleDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.DeleteTCPRequestRule(params.Index, params.ParentType, params.ParentName, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return tcp_request_rule.NewDeleteTCPRequestRuleDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return tcp_request_rule.NewDeleteTCPRequestRuleDefault(int(*e.Code)).WithPayload(e)
			}
			return tcp_request_rule.NewDeleteTCPRequestRuleNoContent()
		}
		rID := h.ReloadAgent.Reload()
		return tcp_request_rule.NewDeleteTCPRequestRuleAccepted().WithReloadID(rID)
	}
	return tcp_request_rule.NewDeleteTCPRequestRuleAccepted()
}

// Handle executing the request and returning a response
func (h *GetTCPRequestRuleHandlerImpl) Handle(params tcp_request_rule.GetTCPRequestRuleParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return tcp_request_rule.NewGetTCPRequestRuleDefault(int(*e.Code)).WithPayload(e)
	}

	_, rule, err := configuration.GetTCPRequestRule(params.Index, params.ParentType, params.ParentName, t)
	if err != nil {
		e := misc.HandleError(err)
		return tcp_request_rule.NewGetTCPRequestRuleDefault(int(*e.Code)).WithPayload(e)
	}
	return tcp_request_rule.NewGetTCPRequestRuleOK().WithPayload(rule)
}

// Handle executing the request and returning a response
func (h *GetTCPRequestRulesHandlerImpl) Handle(params tcp_request_rule.GetTCPRequestRulesParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return tcp_request_rule.NewGetTCPRequestRulesDefault(int(*e.Code)).WithPayload(e)
	}

	_, rules, err := configuration.GetTCPRequestRules(params.ParentType, params.ParentName, t)
	if err != nil {
		e := misc.HandleContainerGetError(err)
		if *e.Code == misc.ErrHTTPOk {
			return tcp_request_rule.NewGetTCPRequestRulesOK().WithPayload(models.TCPRequestRules{})
		}
		return tcp_request_rule.NewGetTCPRequestRulesDefault(int(*e.Code)).WithPayload(e)
	}
	return tcp_request_rule.NewGetTCPRequestRulesOK().WithPayload(rules)
}

// Handle executing the request and returning a response
func (h *ReplaceTCPRequestRuleHandlerImpl) Handle(params tcp_request_rule.ReplaceTCPRequestRuleParams, principal interface{}) middleware.Responder {
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
		return tcp_request_rule.NewReplaceTCPRequestRuleDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return tcp_request_rule.NewReplaceTCPRequestRuleDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.EditTCPRequestRule(params.Index, params.ParentType, params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return tcp_request_rule.NewReplaceTCPRequestRuleDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return tcp_request_rule.NewReplaceTCPRequestRuleDefault(int(*e.Code)).WithPayload(e)
			}
			return tcp_request_rule.NewReplaceTCPRequestRuleOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return tcp_request_rule.NewReplaceTCPRequestRuleAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return tcp_request_rule.NewReplaceTCPRequestRuleAccepted().WithPayload(params.Data)
}

// Handle executing the request and returning a response
func (h *ReplaceTCPRequestRulesHandlerImpl) Handle(params tcp_request_rule.ReplaceTCPRequestRulesParams, principal interface{}) middleware.Responder {
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
		return tcp_request_rule.NewReplaceTCPRequestRulesDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return tcp_request_rule.NewReplaceTCPRequestRulesDefault(int(*e.Code)).WithPayload(e)
	}
	err = configuration.ReplaceTCPRequestRules(params.ParentType, params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return tcp_request_rule.NewReplaceTCPRequestRulesDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return tcp_request_rule.NewReplaceTCPRequestRulesDefault(int(*e.Code)).WithPayload(e)
			}
			return tcp_request_rule.NewReplaceTCPRequestRulesOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return tcp_request_rule.NewReplaceTCPRequestRulesAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return tcp_request_rule.NewReplaceTCPRequestRulesAccepted().WithPayload(params.Data)
}
