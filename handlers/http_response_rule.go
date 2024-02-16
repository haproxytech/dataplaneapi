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

// GetHTTPResponseRulesHandlerImpl implementation of the GetHTTPResponseRulesHandler interface using client-native client
type GetHTTPResponseRulesHandlerImpl struct {
	Client client_native.HAProxyClient
}

// ReplaceHTTPResponseRuleHandlerImpl implementation of the ReplaceHTTPResponseRuleHandler interface using client-native client
type ReplaceHTTPResponseRuleHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// Handle executing the request and returning a response
func (h *CreateHTTPResponseRuleHandlerImpl) Handle(params http_response_rule.CreateHTTPResponseRuleParams, principal interface{}) middleware.Responder {
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
		return http_response_rule.NewCreateHTTPResponseRuleDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_response_rule.NewCreateHTTPResponseRuleDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.CreateHTTPResponseRule(params.ParentType, params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return http_response_rule.NewCreateHTTPResponseRuleDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return http_response_rule.NewCreateHTTPResponseRuleDefault(int(*e.Code)).WithPayload(e)
			}
			return http_response_rule.NewCreateHTTPResponseRuleCreated().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return http_response_rule.NewCreateHTTPResponseRuleAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return http_response_rule.NewCreateHTTPResponseRuleAccepted().WithPayload(params.Data)
}

// Handle executing the request and returning a response
func (h *DeleteHTTPResponseRuleHandlerImpl) Handle(params http_response_rule.DeleteHTTPResponseRuleParams, principal interface{}) middleware.Responder {
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
		return http_response_rule.NewDeleteHTTPResponseRuleDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_response_rule.NewDeleteHTTPResponseRuleDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.DeleteHTTPResponseRule(params.Index, params.ParentType, params.ParentName, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return http_response_rule.NewDeleteHTTPResponseRuleDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return http_response_rule.NewDeleteHTTPResponseRuleDefault(int(*e.Code)).WithPayload(e)
			}
			return http_response_rule.NewDeleteHTTPResponseRuleNoContent()
		}
		rID := h.ReloadAgent.Reload()
		return http_response_rule.NewDeleteHTTPResponseRuleAccepted().WithReloadID(rID)
	}
	return http_response_rule.NewDeleteHTTPResponseRuleAccepted()
}

// Handle executing the request and returning a response
func (h *GetHTTPResponseRuleHandlerImpl) Handle(params http_response_rule.GetHTTPResponseRuleParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_response_rule.NewCreateHTTPResponseRuleDefault(int(*e.Code)).WithPayload(e)
	}

	_, rule, err := configuration.GetHTTPResponseRule(params.Index, params.ParentType, params.ParentName, t)
	if err != nil {
		e := misc.HandleError(err)
		return http_response_rule.NewGetHTTPResponseRuleDefault(int(*e.Code)).WithPayload(e)
	}
	return http_response_rule.NewGetHTTPResponseRuleOK().WithPayload(rule)
}

// Handle executing the request and returning a response
func (h *GetHTTPResponseRulesHandlerImpl) Handle(params http_response_rule.GetHTTPResponseRulesParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_response_rule.NewGetHTTPResponseRuleDefault(int(*e.Code)).WithPayload(e)
	}

	_, rules, err := configuration.GetHTTPResponseRules(params.ParentType, params.ParentName, t)
	if err != nil {
		e := misc.HandleContainerGetError(err)
		if *e.Code == misc.ErrHTTPOk {
			return http_response_rule.NewGetHTTPResponseRulesOK().WithPayload(models.HTTPResponseRules{})
		}
		return http_response_rule.NewGetHTTPResponseRulesDefault(int(*e.Code)).WithPayload(e)
	}
	return http_response_rule.NewGetHTTPResponseRulesOK().WithPayload(rules)
}

// Handle executing the request and returning a response
func (h *ReplaceHTTPResponseRuleHandlerImpl) Handle(params http_response_rule.ReplaceHTTPResponseRuleParams, principal interface{}) middleware.Responder {
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
		return http_response_rule.NewReplaceHTTPResponseRuleDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_response_rule.NewReplaceHTTPResponseRuleDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.EditHTTPResponseRule(params.Index, params.ParentType, params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return http_response_rule.NewReplaceHTTPResponseRuleDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return http_response_rule.NewReplaceHTTPResponseRuleDefault(int(*e.Code)).WithPayload(e)
			}
			return http_response_rule.NewReplaceHTTPResponseRuleOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return http_response_rule.NewReplaceHTTPResponseRuleAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return http_response_rule.NewReplaceHTTPResponseRuleAccepted().WithPayload(params.Data)
}
