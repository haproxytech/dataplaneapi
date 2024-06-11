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

// GetTCPResponseRulesHandlerImpl implementation of the GetTCPResponseRulesHandler interface using client-native client
type GetTCPResponseRulesHandlerImpl struct {
	Client client_native.HAProxyClient
}

// ReplaceTCPResponseRuleHandlerImpl implementation of the ReplaceTCPResponseRuleHandler interface using client-native client
type ReplaceTCPResponseRuleHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// ReplaceTCPResponseRulesHandlerImpl implementation of the ReplaceTCPResponseRulesHandler interface using client-native client
type ReplaceTCPResponseRulesHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// Handle executing the request and returning a response
func (h *CreateTCPResponseRuleHandlerImpl) Handle(params tcp_response_rule.CreateTCPResponseRuleParams, principal interface{}) middleware.Responder {
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
		return tcp_response_rule.NewCreateTCPResponseRuleDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return tcp_response_rule.NewCreateTCPResponseRuleDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.CreateTCPResponseRule(params.Index, params.Backend, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return tcp_response_rule.NewCreateTCPResponseRuleDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return tcp_response_rule.NewCreateTCPResponseRuleDefault(int(*e.Code)).WithPayload(e)
			}
			return tcp_response_rule.NewCreateTCPResponseRuleCreated().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return tcp_response_rule.NewCreateTCPResponseRuleAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return tcp_response_rule.NewCreateTCPResponseRuleAccepted().WithPayload(params.Data)
}

// Handle executing the request and returning a response
func (h *DeleteTCPResponseRuleHandlerImpl) Handle(params tcp_response_rule.DeleteTCPResponseRuleParams, principal interface{}) middleware.Responder {
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
		return tcp_response_rule.NewDeleteTCPResponseRuleDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return tcp_response_rule.NewDeleteTCPResponseRuleDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.DeleteTCPResponseRule(params.Index, params.Backend, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return tcp_response_rule.NewDeleteTCPResponseRuleDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return tcp_response_rule.NewDeleteTCPResponseRuleDefault(int(*e.Code)).WithPayload(e)
			}
			return tcp_response_rule.NewDeleteTCPResponseRuleNoContent()
		}
		rID := h.ReloadAgent.Reload()
		return tcp_response_rule.NewDeleteTCPResponseRuleAccepted().WithReloadID(rID)
	}
	return tcp_response_rule.NewDeleteTCPResponseRuleAccepted()
}

// Handle executing the request and returning a response
func (h *GetTCPResponseRuleHandlerImpl) Handle(params tcp_response_rule.GetTCPResponseRuleParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return tcp_response_rule.NewGetTCPResponseRuleDefault(int(*e.Code)).WithPayload(e)
	}
	_, rule, err := configuration.GetTCPResponseRule(params.Index, params.Backend, t)
	if err != nil {
		e := misc.HandleError(err)
		return tcp_response_rule.NewGetTCPResponseRuleDefault(int(*e.Code)).WithPayload(e)
	}
	return tcp_response_rule.NewGetTCPResponseRuleOK().WithPayload(rule)
}

// Handle executing the request and returning a response
func (h *GetTCPResponseRulesHandlerImpl) Handle(params tcp_response_rule.GetTCPResponseRulesParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return tcp_response_rule.NewGetTCPResponseRulesDefault(int(*e.Code)).WithPayload(e)
	}

	_, rules, err := configuration.GetTCPResponseRules(params.Backend, t)
	if err != nil {
		e := misc.HandleContainerGetError(err)
		if *e.Code == misc.ErrHTTPOk {
			return tcp_response_rule.NewGetTCPResponseRulesOK().WithPayload(models.TCPResponseRules{})
		}
		return tcp_response_rule.NewGetTCPResponseRulesDefault(int(*e.Code)).WithPayload(e)
	}
	return tcp_response_rule.NewGetTCPResponseRulesOK().WithPayload(rules)
}

// Handle executing the request and returning a response
func (h *ReplaceTCPResponseRuleHandlerImpl) Handle(params tcp_response_rule.ReplaceTCPResponseRuleParams, principal interface{}) middleware.Responder {
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
		return tcp_response_rule.NewReplaceTCPResponseRuleDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return tcp_response_rule.NewReplaceTCPResponseRuleDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.EditTCPResponseRule(params.Index, params.Backend, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return tcp_response_rule.NewReplaceTCPResponseRuleDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return tcp_response_rule.NewReplaceTCPResponseRuleDefault(int(*e.Code)).WithPayload(e)
			}
			return tcp_response_rule.NewReplaceTCPResponseRuleOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return tcp_response_rule.NewReplaceTCPResponseRuleAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return tcp_response_rule.NewReplaceTCPResponseRuleAccepted().WithPayload(params.Data)
}

func (h *ReplaceTCPResponseRulesHandlerImpl) Handle(params tcp_response_rule.ReplaceTCPResponseRulesParams, principal interface{}) middleware.Responder {
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
		return tcp_response_rule.NewReplaceTCPResponseRulesDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return tcp_response_rule.NewReplaceTCPResponseRulesDefault(int(*e.Code)).WithPayload(e)
	}
	err = configuration.ReplaceTCPResponseRules(params.Backend, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return tcp_response_rule.NewReplaceTCPResponseRulesDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return tcp_response_rule.NewReplaceTCPResponseRulesDefault(int(*e.Code)).WithPayload(e)
			}
			return tcp_response_rule.NewReplaceTCPResponseRulesOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return tcp_response_rule.NewReplaceTCPResponseRulesAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return tcp_response_rule.NewReplaceTCPResponseRulesAccepted().WithPayload(params.Data)
}
