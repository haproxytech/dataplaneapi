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
	"github.com/haproxytech/dataplaneapi/operations/acl"
	"github.com/haproxytech/dataplaneapi/operations/backend_switching_rule"
)

// CreateBackendSwitchingRuleHandlerImpl implementation of the CreateBackendSwitchingRuleHandler interface using client-native client
type CreateBackendSwitchingRuleHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// DeleteBackendSwitchingRuleHandlerImpl implementation of the DeleteBackendSwitchingRuleHandler interface using client-native client
type DeleteBackendSwitchingRuleHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// GetBackendSwitchingRuleHandlerImpl implementation of the GetBackendSwitchingRuleHandler interface using client-native client
type GetBackendSwitchingRuleHandlerImpl struct {
	Client client_native.HAProxyClient
}

// GetBackendSwitchingRulesHandlerImpl implementation of the GetBackendSwitchingRulesHandler interface using client-native client
type GetBackendSwitchingRulesHandlerImpl struct {
	Client client_native.HAProxyClient
}

// ReplaceBackendSwitchingRuleHandlerImpl implementation of the ReplaceBackendSwitchingRuleHandler interface using client-native client
type ReplaceBackendSwitchingRuleHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// ReplaceBackendSwitchingRulesHandlerImpl implementation of the ReplaceBackendSwitchingRulesHandler interface using client-native client
type ReplaceBackendSwitchingRulesHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// Handle executing the request and returning a response
func (h *CreateBackendSwitchingRuleHandlerImpl) Handle(params backend_switching_rule.CreateBackendSwitchingRuleParams, principal interface{}) middleware.Responder {
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
		return backend_switching_rule.NewCreateBackendSwitchingRuleDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return backend_switching_rule.NewCreateBackendSwitchingRuleDefault(int(*e.Code)).WithPayload(e)
	}
	err = configuration.CreateBackendSwitchingRule(params.Index, params.Frontend, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return backend_switching_rule.NewCreateBackendSwitchingRuleDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return backend_switching_rule.NewCreateBackendSwitchingRuleDefault(int(*e.Code)).WithPayload(e)
			}
			return backend_switching_rule.NewCreateBackendSwitchingRuleCreated().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return backend_switching_rule.NewCreateBackendSwitchingRuleAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return backend_switching_rule.NewCreateBackendSwitchingRuleAccepted().WithPayload(params.Data)
}

// Handle executing the request and returning a response
func (h *DeleteBackendSwitchingRuleHandlerImpl) Handle(params backend_switching_rule.DeleteBackendSwitchingRuleParams, principal interface{}) middleware.Responder {
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
		return backend_switching_rule.NewDeleteBackendSwitchingRuleDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return backend_switching_rule.NewDeleteBackendSwitchingRuleDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.DeleteBackendSwitchingRule(params.Index, params.Frontend, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return backend_switching_rule.NewDeleteBackendSwitchingRuleDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return backend_switching_rule.NewDeleteBackendSwitchingRuleDefault(int(*e.Code)).WithPayload(e)
			}
			return backend_switching_rule.NewDeleteBackendSwitchingRuleNoContent()
		}
		rID := h.ReloadAgent.Reload()
		return backend_switching_rule.NewDeleteBackendSwitchingRuleAccepted().WithReloadID(rID)
	}
	return backend_switching_rule.NewDeleteBackendSwitchingRuleAccepted()
}

// Handle executing the request and returning a response
func (h *GetBackendSwitchingRuleHandlerImpl) Handle(params backend_switching_rule.GetBackendSwitchingRuleParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return backend_switching_rule.NewCreateBackendSwitchingRuleDefault(int(*e.Code)).WithPayload(e)
	}

	_, bckRule, err := configuration.GetBackendSwitchingRule(params.Index, params.Frontend, t)
	if err != nil {
		e := misc.HandleError(err)
		return backend_switching_rule.NewGetBackendSwitchingRuleDefault(int(*e.Code)).WithPayload(e)
	}
	return backend_switching_rule.NewGetBackendSwitchingRuleOK().WithPayload(bckRule)
}

// Handle executing the request and returning a response
func (h *GetBackendSwitchingRulesHandlerImpl) Handle(params backend_switching_rule.GetBackendSwitchingRulesParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return backend_switching_rule.NewGetBackendSwitchingRulesDefault(int(*e.Code)).WithPayload(e)
	}

	_, bckRules, err := configuration.GetBackendSwitchingRules(params.Frontend, t)
	if err != nil {
		e := misc.HandleContainerGetError(err)
		if *e.Code == misc.ErrHTTPOk {
			return backend_switching_rule.NewGetBackendSwitchingRulesOK().WithPayload(models.BackendSwitchingRules{})
		}
		return backend_switching_rule.NewGetBackendSwitchingRulesDefault(int(*e.Code)).WithPayload(e)
	}
	return backend_switching_rule.NewGetBackendSwitchingRulesOK().WithPayload(bckRules)
}

// Handle executing the request and returning a response
func (h *ReplaceBackendSwitchingRuleHandlerImpl) Handle(params backend_switching_rule.ReplaceBackendSwitchingRuleParams, principal interface{}) middleware.Responder {
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
		return backend_switching_rule.NewReplaceBackendSwitchingRuleDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return backend_switching_rule.NewReplaceBackendSwitchingRuleDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.EditBackendSwitchingRule(params.Index, params.Frontend, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return backend_switching_rule.NewReplaceBackendSwitchingRuleDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return backend_switching_rule.NewReplaceBackendSwitchingRuleDefault(int(*e.Code)).WithPayload(e)
			}
			return backend_switching_rule.NewReplaceBackendSwitchingRuleOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return backend_switching_rule.NewReplaceBackendSwitchingRuleAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return backend_switching_rule.NewReplaceBackendSwitchingRuleAccepted().WithPayload(params.Data)
}

// Handle executing the request and returning a response
func (h *ReplaceBackendSwitchingRulesHandlerImpl) Handle(params backend_switching_rule.ReplaceBackendSwitchingRulesParams, principal interface{}) middleware.Responder {
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
		return backend_switching_rule.NewReplaceBackendSwitchingRulesDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return acl.NewReplaceAclsDefault(int(*e.Code)).WithPayload(e)
	}
	err = configuration.ReplaceBackendSwitchingRules(params.Frontend, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return backend_switching_rule.NewReplaceBackendSwitchingRulesDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return backend_switching_rule.NewReplaceBackendSwitchingRulesDefault(int(*e.Code)).WithPayload(e)
			}
			return backend_switching_rule.NewReplaceBackendSwitchingRulesOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return backend_switching_rule.NewReplaceBackendSwitchingRulesAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return backend_switching_rule.NewReplaceBackendSwitchingRulesAccepted().WithPayload(params.Data)
}
