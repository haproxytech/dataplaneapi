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

// GetHTTPRequestRulesHandlerImpl implementation of the GetHTTPRequestRulesHandler interface using client-native client
type GetHTTPRequestRulesHandlerImpl struct {
	Client client_native.HAProxyClient
}

// ReplaceHTTPRequestRuleHandlerImpl implementation of the ReplaceHTTPRequestRuleHandler interface using client-native client
type ReplaceHTTPRequestRuleHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// ReplaceHTTPRequestRulesHandlerImpl implementation of the ReplaceHTTPRequestRulesHandler interface using client-native client
type ReplaceHTTPRequestRulesHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// Handle executing the request and returning a response
func (h *CreateHTTPRequestRuleHandlerImpl) Handle(params http_request_rule.CreateHTTPRequestRuleParams, principal interface{}) middleware.Responder {
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
		return http_request_rule.NewCreateHTTPRequestRuleDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_request_rule.NewCreateHTTPRequestRuleDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.CreateHTTPRequestRule(params.Index, params.ParentType, params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return http_request_rule.NewCreateHTTPRequestRuleDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return http_request_rule.NewCreateHTTPRequestRuleDefault(int(*e.Code)).WithPayload(e)
			}
			return http_request_rule.NewCreateHTTPRequestRuleCreated().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return http_request_rule.NewCreateHTTPRequestRuleAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return http_request_rule.NewCreateHTTPRequestRuleAccepted().WithPayload(params.Data)
}

// Handle executing the request and returning a response
func (h *DeleteHTTPRequestRuleHandlerImpl) Handle(params http_request_rule.DeleteHTTPRequestRuleParams, principal interface{}) middleware.Responder {
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
		return http_request_rule.NewDeleteHTTPRequestRuleDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_request_rule.NewDeleteHTTPRequestRuleDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.DeleteHTTPRequestRule(params.Index, params.ParentType, params.ParentName, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return http_request_rule.NewDeleteHTTPRequestRuleDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return http_request_rule.NewDeleteHTTPRequestRuleDefault(int(*e.Code)).WithPayload(e)
			}
			return http_request_rule.NewDeleteHTTPRequestRuleNoContent()
		}
		rID := h.ReloadAgent.Reload()
		return http_request_rule.NewDeleteHTTPRequestRuleAccepted().WithReloadID(rID)
	}
	return http_request_rule.NewDeleteHTTPRequestRuleAccepted()
}

// Handle executing the request and returning a response
func (h *GetHTTPRequestRuleHandlerImpl) Handle(params http_request_rule.GetHTTPRequestRuleParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_request_rule.NewGetHTTPRequestRuleDefault(int(*e.Code)).WithPayload(e)
	}

	_, rule, err := configuration.GetHTTPRequestRule(params.Index, params.ParentType, params.ParentName, t)
	if err != nil {
		e := misc.HandleError(err)
		return http_request_rule.NewGetHTTPRequestRuleDefault(int(*e.Code)).WithPayload(e)
	}
	return http_request_rule.NewGetHTTPRequestRuleOK().WithPayload(rule)
}

// Handle executing the request and returning a response
func (h *GetHTTPRequestRulesHandlerImpl) Handle(params http_request_rule.GetHTTPRequestRulesParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_request_rule.NewGetHTTPRequestRuleDefault(int(*e.Code)).WithPayload(e)
	}

	_, rules, err := configuration.GetHTTPRequestRules(params.ParentType, params.ParentName, t)
	if err != nil {
		e := misc.HandleContainerGetError(err)
		if *e.Code == misc.ErrHTTPOk {
			return http_request_rule.NewGetHTTPRequestRulesOK().WithPayload(models.HTTPRequestRules{})
		}
		return http_request_rule.NewGetHTTPRequestRulesDefault(int(*e.Code)).WithPayload(e)
	}
	return http_request_rule.NewGetHTTPRequestRulesOK().WithPayload(rules)
}

// Handle executing the request and returning a response
func (h *ReplaceHTTPRequestRuleHandlerImpl) Handle(params http_request_rule.ReplaceHTTPRequestRuleParams, principal interface{}) middleware.Responder {
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
		return http_request_rule.NewReplaceHTTPRequestRuleDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_request_rule.NewReplaceHTTPRequestRuleDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.EditHTTPRequestRule(params.Index, params.ParentType, params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return http_request_rule.NewReplaceHTTPRequestRuleDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return http_request_rule.NewReplaceHTTPRequestRuleDefault(int(*e.Code)).WithPayload(e)
			}
			return http_request_rule.NewReplaceHTTPRequestRuleOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return http_request_rule.NewReplaceHTTPRequestRuleAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return http_request_rule.NewReplaceHTTPRequestRuleAccepted().WithPayload(params.Data)
}

// Handle executing the request and returning a response
func (h *ReplaceHTTPRequestRulesHandlerImpl) Handle(params http_request_rule.ReplaceHTTPRequestRulesParams, principal interface{}) middleware.Responder {
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
		return http_request_rule.NewReplaceHTTPRequestRulesDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_request_rule.NewReplaceHTTPRequestRulesDefault(int(*e.Code)).WithPayload(e)
	}
	err = configuration.ReplaceHTTPRequestRules(params.ParentType, params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return http_request_rule.NewReplaceHTTPRequestRulesDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return http_request_rule.NewReplaceHTTPRequestRulesDefault(int(*e.Code)).WithPayload(e)
			}
			return http_request_rule.NewReplaceHTTPRequestRulesOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return http_request_rule.NewReplaceHTTPRequestRulesAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return http_request_rule.NewReplaceHTTPRequestRulesAccepted().WithPayload(params.Data)
}
