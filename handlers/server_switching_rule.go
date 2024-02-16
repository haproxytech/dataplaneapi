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
	"github.com/haproxytech/dataplaneapi/operations/server_switching_rule"
)

// CreateServerSwitchingRuleHandlerImpl implementation of the CreateServerSwitchingRuleHandler interface using client-native client
type CreateServerSwitchingRuleHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// DeleteServerSwitchingRuleHandlerImpl implementation of the DeleteServerSwitchingRuleHandler interface using client-native client
type DeleteServerSwitchingRuleHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// GetServerSwitchingRuleHandlerImpl implementation of the GetServerSwitchingRuleHandler interface using client-native client
type GetServerSwitchingRuleHandlerImpl struct {
	Client client_native.HAProxyClient
}

// GetServerSwitchingRulesHandlerImpl implementation of the GetServerSwitchingRulesHandler interface using client-native client
type GetServerSwitchingRulesHandlerImpl struct {
	Client client_native.HAProxyClient
}

// ReplaceServerSwitchingRuleHandlerImpl implementation of the ReplaceServerSwitchingRuleHandler interface using client-native client
type ReplaceServerSwitchingRuleHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// Handle executing the request and returning a response
func (h *CreateServerSwitchingRuleHandlerImpl) Handle(params server_switching_rule.CreateServerSwitchingRuleParams, principal interface{}) middleware.Responder {
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
		return server_switching_rule.NewCreateServerSwitchingRuleDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return server_switching_rule.NewCreateServerSwitchingRuleDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.CreateServerSwitchingRule(params.Backend, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return server_switching_rule.NewCreateServerSwitchingRuleDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return server_switching_rule.NewCreateServerSwitchingRuleDefault(int(*e.Code)).WithPayload(e)
			}
			return server_switching_rule.NewCreateServerSwitchingRuleCreated().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return server_switching_rule.NewCreateServerSwitchingRuleAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return server_switching_rule.NewCreateServerSwitchingRuleAccepted().WithPayload(params.Data)
}

// Handle executing the request and returning a response
func (h *DeleteServerSwitchingRuleHandlerImpl) Handle(params server_switching_rule.DeleteServerSwitchingRuleParams, principal interface{}) middleware.Responder {
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
		return server_switching_rule.NewDeleteServerSwitchingRuleDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return server_switching_rule.NewDeleteServerSwitchingRuleDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.DeleteServerSwitchingRule(params.Index, params.Backend, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return server_switching_rule.NewDeleteServerSwitchingRuleDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return server_switching_rule.NewDeleteServerSwitchingRuleDefault(int(*e.Code)).WithPayload(e)
			}
			return server_switching_rule.NewDeleteServerSwitchingRuleNoContent()
		}
		rID := h.ReloadAgent.Reload()
		return server_switching_rule.NewDeleteServerSwitchingRuleAccepted().WithReloadID(rID)
	}
	return server_switching_rule.NewDeleteServerSwitchingRuleAccepted()
}

// Handle executing the request and returning a response
func (h *GetServerSwitchingRuleHandlerImpl) Handle(params server_switching_rule.GetServerSwitchingRuleParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return server_switching_rule.NewGetServerSwitchingRuleDefault(int(*e.Code)).WithPayload(e)
	}

	_, rule, err := configuration.GetServerSwitchingRule(params.Index, params.Backend, t)
	if err != nil {
		e := misc.HandleError(err)
		return server_switching_rule.NewGetServerSwitchingRuleDefault(int(*e.Code)).WithPayload(e)
	}
	return server_switching_rule.NewGetServerSwitchingRuleOK().WithPayload(rule)
}

// Handle executing the request and returning a response
func (h *GetServerSwitchingRulesHandlerImpl) Handle(params server_switching_rule.GetServerSwitchingRulesParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return server_switching_rule.NewGetServerSwitchingRulesDefault(int(*e.Code)).WithPayload(e)
	}

	_, rules, err := configuration.GetServerSwitchingRules(params.Backend, t)
	if err != nil {
		e := misc.HandleContainerGetError(err)
		if *e.Code == misc.ErrHTTPOk {
			return server_switching_rule.NewGetServerSwitchingRulesOK().WithPayload(models.ServerSwitchingRules{})
		}
		return server_switching_rule.NewGetServerSwitchingRulesDefault(int(*e.Code)).WithPayload(e)
	}
	return server_switching_rule.NewGetServerSwitchingRulesOK().WithPayload(rules)
}

// Handle executing the request and returning a response
func (h *ReplaceServerSwitchingRuleHandlerImpl) Handle(params server_switching_rule.ReplaceServerSwitchingRuleParams, principal interface{}) middleware.Responder {
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
		return server_switching_rule.NewReplaceServerSwitchingRuleDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return server_switching_rule.NewReplaceServerSwitchingRuleDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.EditServerSwitchingRule(params.Index, params.Backend, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return server_switching_rule.NewReplaceServerSwitchingRuleDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return server_switching_rule.NewReplaceServerSwitchingRuleDefault(int(*e.Code)).WithPayload(e)
			}
			return server_switching_rule.NewReplaceServerSwitchingRuleOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return server_switching_rule.NewReplaceServerSwitchingRuleAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return server_switching_rule.NewReplaceServerSwitchingRuleAccepted().WithPayload(params.Data)
}
