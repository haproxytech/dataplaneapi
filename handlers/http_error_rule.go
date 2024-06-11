// Copyright 2022 HAProxy Technologies
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
	"github.com/haproxytech/dataplaneapi/operations/http_error_rule"
)

// CreateHTTPErrorRuleHandlerImpl implementation of the CreateHTTPErrorRuleHandler interface using client-native client
type CreateHTTPErrorRuleHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// DeleteHTTPErrorRuleHandlerImpl implementation of the DeleteHTTPErrorRuleHandler interface using client-native client
type DeleteHTTPErrorRuleHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// GetHTTPErrorRuleHandlerImpl implementation of the GetHTTPErrorRuleHandler interface using client-native client
type GetHTTPErrorRuleHandlerImpl struct {
	Client client_native.HAProxyClient
}

// GetHTTPErrorRulesHandlerImpl implementation of the GetHTTPErrorRulesHandler interface using client-native client
type GetHTTPErrorRulesHandlerImpl struct {
	Client client_native.HAProxyClient
}

// ReplaceHTTPErrorRuleHandlerImpl implementation of the ReplaceHTTPErrorRuleHandler interface using client-native client
type ReplaceHTTPErrorRuleHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// ReplaceHTTPErrorRulesHandlerImpl implementation of the ReplaceHTTPErrorRulesHandler interface using client-native client
type ReplaceHTTPErrorRulesHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// Handle executing the request and returning a response
func (h *CreateHTTPErrorRuleHandlerImpl) Handle(params http_error_rule.CreateHTTPErrorRuleParams, principal interface{}) middleware.Responder {
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
		return http_error_rule.NewCreateHTTPErrorRuleDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_error_rule.NewCreateHTTPErrorRuleDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.CreateHTTPErrorRule(params.Index, params.ParentType, params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return http_error_rule.NewCreateHTTPErrorRuleDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return http_error_rule.NewCreateHTTPErrorRuleDefault(int(*e.Code)).WithPayload(e)
			}
			return http_error_rule.NewCreateHTTPErrorRuleCreated().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return http_error_rule.NewCreateHTTPErrorRuleAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return http_error_rule.NewCreateHTTPErrorRuleAccepted().WithPayload(params.Data)
}

// Handle executing the request and returning a response
func (h *DeleteHTTPErrorRuleHandlerImpl) Handle(params http_error_rule.DeleteHTTPErrorRuleParams, principal interface{}) middleware.Responder {
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
		return http_error_rule.NewDeleteHTTPErrorRuleDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_error_rule.NewDeleteHTTPErrorRuleDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.DeleteHTTPErrorRule(params.Index, params.ParentType, params.ParentName, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return http_error_rule.NewDeleteHTTPErrorRuleDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return http_error_rule.NewDeleteHTTPErrorRuleDefault(int(*e.Code)).WithPayload(e)
			}
			return http_error_rule.NewDeleteHTTPErrorRuleNoContent()
		}
		rID := h.ReloadAgent.Reload()
		return http_error_rule.NewDeleteHTTPErrorRuleAccepted().WithReloadID(rID)
	}
	return http_error_rule.NewDeleteHTTPErrorRuleAccepted()
}

// Handle executing the request and returning a response
func (h *GetHTTPErrorRuleHandlerImpl) Handle(params http_error_rule.GetHTTPErrorRuleParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_error_rule.NewCreateHTTPErrorRuleDefault(int(*e.Code)).WithPayload(e)
	}

	_, rule, err := configuration.GetHTTPErrorRule(params.Index, params.ParentType, params.ParentName, t)
	if err != nil {
		e := misc.HandleError(err)
		return http_error_rule.NewGetHTTPErrorRuleDefault(int(*e.Code)).WithPayload(e)
	}
	return http_error_rule.NewGetHTTPErrorRuleOK().WithPayload(rule)
}

// Handle executing the request and returning a response
func (h *GetHTTPErrorRulesHandlerImpl) Handle(params http_error_rule.GetHTTPErrorRulesParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_error_rule.NewGetHTTPErrorRuleDefault(int(*e.Code)).WithPayload(e)
	}

	_, rules, err := configuration.GetHTTPErrorRules(params.ParentType, params.ParentName, t)
	if err != nil {
		e := misc.HandleContainerGetError(err)
		if *e.Code == misc.ErrHTTPOk {
			return http_error_rule.NewGetHTTPErrorRulesOK().WithPayload(models.HTTPErrorRules{})
		}
		return http_error_rule.NewGetHTTPErrorRulesDefault(int(*e.Code)).WithPayload(e)
	}
	return http_error_rule.NewGetHTTPErrorRulesOK().WithPayload(rules)
}

// Handle executing the request and returning a response
func (h *ReplaceHTTPErrorRuleHandlerImpl) Handle(params http_error_rule.ReplaceHTTPErrorRuleParams, principal interface{}) middleware.Responder {
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
		return http_error_rule.NewReplaceHTTPErrorRuleDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_error_rule.NewReplaceHTTPErrorRuleDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.EditHTTPErrorRule(params.Index, params.ParentType, params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return http_error_rule.NewReplaceHTTPErrorRuleDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return http_error_rule.NewReplaceHTTPErrorRuleDefault(int(*e.Code)).WithPayload(e)
			}
			return http_error_rule.NewReplaceHTTPErrorRuleOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return http_error_rule.NewReplaceHTTPErrorRuleAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return http_error_rule.NewReplaceHTTPErrorRuleAccepted().WithPayload(params.Data)
}

func (h *ReplaceHTTPErrorRulesHandlerImpl) Handle(params http_error_rule.ReplaceHTTPErrorRulesParams, principal interface{}) middleware.Responder {
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
		return http_error_rule.NewReplaceHTTPErrorRulesDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_error_rule.NewReplaceHTTPErrorRulesDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.ReplaceHTTPErrorRules(params.ParentType, params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return http_error_rule.NewReplaceHTTPErrorRulesDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return http_error_rule.NewReplaceHTTPErrorRulesDefault(int(*e.Code)).WithPayload(e)
			}
			return http_error_rule.NewReplaceHTTPErrorRulesOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return http_error_rule.NewReplaceHTTPErrorRulesAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return http_error_rule.NewReplaceHTTPErrorRulesAccepted().WithPayload(params.Data)
}
