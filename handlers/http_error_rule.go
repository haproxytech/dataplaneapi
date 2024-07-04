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
	cnconstants "github.com/haproxytech/client-native/v6/configuration/parents"
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

// GetAllHTTPErrorRuleHandlerImpl implementation of the GetHTTPErrorRulesHandler interface using client-native client
type GetAllHTTPErrorRuleHandlerImpl struct {
	Client client_native.HAProxyClient
}

// ReplaceHTTPErrorRuleHandlerImpl implementation of the ReplaceHTTPErrorRuleHandler interface using client-native client
type ReplaceHTTPErrorRuleHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// ReplaceAllHTTPErrorRuleHandlerImpl implementation of the ReplaceHTTPErrorRulesHandler interface using client-native client
type ReplaceAllHTTPErrorRuleHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// Handle executing the request and returning a response
func (h *CreateHTTPErrorRuleHandlerImpl) Handle(parentType cnconstants.CnParentType, params http_error_rule.CreateHTTPErrorRuleBackendParams, principal interface{}) middleware.Responder {
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
		return http_error_rule.NewCreateHTTPErrorRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_error_rule.NewCreateHTTPErrorRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.CreateHTTPErrorRule(params.Index, string(parentType), params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return http_error_rule.NewCreateHTTPErrorRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return http_error_rule.NewCreateHTTPErrorRuleBackendDefault(int(*e.Code)).WithPayload(e)
			}
			return http_error_rule.NewCreateHTTPErrorRuleBackendCreated().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return http_error_rule.NewCreateHTTPErrorRuleBackendAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return http_error_rule.NewCreateHTTPErrorRuleBackendAccepted().WithPayload(params.Data)
}

// Handle executing the request and returning a response
func (h *DeleteHTTPErrorRuleHandlerImpl) Handle(parentType cnconstants.CnParentType, params http_error_rule.DeleteHTTPErrorRuleBackendParams, principal interface{}) middleware.Responder {
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
		return http_error_rule.NewDeleteHTTPErrorRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_error_rule.NewDeleteHTTPErrorRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.DeleteHTTPErrorRule(params.Index, string(parentType), params.ParentName, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return http_error_rule.NewDeleteHTTPErrorRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return http_error_rule.NewDeleteHTTPErrorRuleBackendDefault(int(*e.Code)).WithPayload(e)
			}
			return http_error_rule.NewDeleteHTTPErrorRuleBackendNoContent()
		}
		rID := h.ReloadAgent.Reload()
		return http_error_rule.NewDeleteHTTPErrorRuleBackendAccepted().WithReloadID(rID)
	}
	return http_error_rule.NewDeleteHTTPErrorRuleBackendAccepted()
}

// Handle executing the request and returning a response
func (h *GetHTTPErrorRuleHandlerImpl) Handle(parentType cnconstants.CnParentType, params http_error_rule.GetHTTPErrorRuleBackendParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_error_rule.NewCreateHTTPErrorRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	_, rule, err := configuration.GetHTTPErrorRule(params.Index, string(parentType), params.ParentName, t)
	if err != nil {
		e := misc.HandleError(err)
		return http_error_rule.NewGetHTTPErrorRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}
	return http_error_rule.NewGetHTTPErrorRuleBackendOK().WithPayload(rule)
}

// Handle executing the request and returning a response
func (h *GetAllHTTPErrorRuleHandlerImpl) Handle(parentType cnconstants.CnParentType, params http_error_rule.GetAllHTTPErrorRuleBackendParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_error_rule.NewGetAllHTTPErrorRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	_, rules, err := configuration.GetHTTPErrorRules(string(parentType), params.ParentName, t)
	if err != nil {
		e := misc.HandleContainerGetError(err)
		if *e.Code == misc.ErrHTTPOk {
			return http_error_rule.NewGetAllHTTPErrorRuleBackendOK().WithPayload(models.HTTPErrorRules{})
		}
		return http_error_rule.NewGetAllHTTPErrorRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}
	return http_error_rule.NewGetAllHTTPErrorRuleBackendOK().WithPayload(rules)
}

// Handle executing the request and returning a response
func (h *ReplaceHTTPErrorRuleHandlerImpl) Handle(parentType cnconstants.CnParentType, params http_error_rule.ReplaceHTTPErrorRuleBackendParams, principal interface{}) middleware.Responder {
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
		return http_error_rule.NewReplaceHTTPErrorRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_error_rule.NewReplaceHTTPErrorRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.EditHTTPErrorRule(params.Index, string(parentType), params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return http_error_rule.NewReplaceHTTPErrorRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return http_error_rule.NewReplaceHTTPErrorRuleBackendDefault(int(*e.Code)).WithPayload(e)
			}
			return http_error_rule.NewReplaceHTTPErrorRuleBackendOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return http_error_rule.NewReplaceHTTPErrorRuleBackendAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return http_error_rule.NewReplaceHTTPErrorRuleBackendAccepted().WithPayload(params.Data)
}

func (h *ReplaceAllHTTPErrorRuleHandlerImpl) Handle(parentType cnconstants.CnParentType, params http_error_rule.ReplaceAllHTTPErrorRuleBackendParams, principal interface{}) middleware.Responder {
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
		return http_error_rule.NewReplaceAllHTTPErrorRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_error_rule.NewReplaceAllHTTPErrorRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.ReplaceHTTPErrorRules(string(parentType), params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return http_error_rule.NewReplaceAllHTTPErrorRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return http_error_rule.NewReplaceAllHTTPErrorRuleBackendDefault(int(*e.Code)).WithPayload(e)
			}
			return http_error_rule.NewReplaceAllHTTPErrorRuleBackendOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return http_error_rule.NewReplaceAllHTTPErrorRuleBackendAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return http_error_rule.NewReplaceAllHTTPErrorRuleBackendAccepted().WithPayload(params.Data)
}
