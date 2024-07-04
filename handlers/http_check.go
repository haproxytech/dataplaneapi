// Copyright 2021 HAProxy Technologies
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
	"github.com/haproxytech/dataplaneapi/operations/http_check"
)

// CreateHTTPCheckHandlerImpl implementation of the CreateHTTPCheckHandler interface using client-native client
type CreateHTTPCheckHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// DeleteHTTPCheckHandlerImpl implementation of the DeleteHTTPCheckHandler interface using client-native client
type DeleteHTTPCheckHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// GetHTTPCheckHandlerImpl implementation of the GetHTTPCheckHandler interface using client-native client
type GetHTTPCheckHandlerImpl struct {
	Client client_native.HAProxyClient
}

// GetAllHTTPCheckHandlerImpl implementation of the GetHTTPChecksHandler interface using client-native client
type GetAllHTTPCheckHandlerImpl struct {
	Client client_native.HAProxyClient
}

// ReplaceHTTPCheckHandlerImpl implementation of the ReplaceHTTPCheckHandler interface using client-native client
type ReplaceHTTPCheckHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// ReplaceAllHTTPCheckHandlerImpl implementation of the ReplaceHTTPChecksHandler interface using client-native client
type ReplaceAllHTTPCheckHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// Handle executing the request and returning a response
func (h *CreateHTTPCheckHandlerImpl) Handle(parentType cnconstants.CnParentType, params http_check.CreateHTTPCheckBackendParams, principal interface{}) middleware.Responder {
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
		return http_check.NewCreateHTTPCheckBackendDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_check.NewCreateHTTPCheckBackendDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.CreateHTTPCheck(params.Index, string(parentType), params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return http_check.NewCreateHTTPCheckBackendDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return http_check.NewCreateHTTPCheckBackendDefault(int(*e.Code)).WithPayload(e)
			}
			return http_check.NewCreateHTTPCheckBackendCreated().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return http_check.NewCreateHTTPCheckBackendAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return http_check.NewCreateHTTPCheckBackendAccepted().WithPayload(params.Data)
}

// Handle executing the request and returning a response
func (h *DeleteHTTPCheckHandlerImpl) Handle(parentType cnconstants.CnParentType, params http_check.DeleteHTTPCheckBackendParams, principal interface{}) middleware.Responder {
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
		return http_check.NewDeleteHTTPCheckBackendDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_check.NewDeleteHTTPCheckBackendDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.DeleteHTTPCheck(params.Index, string(parentType), params.ParentName, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return http_check.NewDeleteHTTPCheckBackendDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return http_check.NewDeleteHTTPCheckBackendDefault(int(*e.Code)).WithPayload(e)
			}
			return http_check.NewDeleteHTTPCheckBackendNoContent()
		}
		rID := h.ReloadAgent.Reload()
		return http_check.NewDeleteHTTPCheckBackendAccepted().WithReloadID(rID)
	}
	return http_check.NewDeleteHTTPCheckBackendAccepted()
}

// Handle executing the request and returning a response
func (h *GetHTTPCheckHandlerImpl) Handle(parentType cnconstants.CnParentType, params http_check.GetHTTPCheckBackendParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_check.NewGetHTTPCheckBackendDefault(int(*e.Code)).WithPayload(e)
	}

	_, rule, err := configuration.GetHTTPCheck(params.Index, string(parentType), params.ParentName, t)
	if err != nil {
		e := misc.HandleError(err)
		return http_check.NewGetHTTPCheckBackendDefault(int(*e.Code)).WithPayload(e)
	}
	return http_check.NewGetHTTPCheckBackendOK().WithPayload(rule)
}

// Handle executing the request and returning a response
func (h *GetAllHTTPCheckHandlerImpl) Handle(parentType cnconstants.CnParentType, params http_check.GetAllHTTPCheckBackendParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_check.NewGetAllHTTPCheckBackendDefault(int(*e.Code)).WithPayload(e)
	}

	_, rules, err := configuration.GetHTTPChecks(string(parentType), params.ParentName, t)
	if err != nil {
		e := misc.HandleContainerGetError(err)
		if *e.Code == misc.ErrHTTPOk {
			return http_check.NewGetAllHTTPCheckBackendOK().WithPayload(models.HTTPChecks{})
		}
		return http_check.NewGetAllHTTPCheckBackendDefault(int(*e.Code)).WithPayload(e)
	}
	return http_check.NewGetAllHTTPCheckBackendOK().WithPayload(rules)
}

// Handle executing the request and returning a response
func (h *ReplaceHTTPCheckHandlerImpl) Handle(parentType cnconstants.CnParentType, params http_check.ReplaceHTTPCheckBackendParams, principal interface{}) middleware.Responder {
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
		return http_check.NewReplaceHTTPCheckBackendDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_check.NewReplaceHTTPCheckBackendDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.EditHTTPCheck(params.Index, string(parentType), params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return http_check.NewReplaceHTTPCheckBackendDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return http_check.NewReplaceHTTPCheckBackendDefault(int(*e.Code)).WithPayload(e)
			}
			return http_check.NewReplaceHTTPCheckBackendOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return http_check.NewReplaceHTTPCheckBackendAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return http_check.NewReplaceHTTPCheckBackendAccepted().WithPayload(params.Data)
}

// Handle executing the request and returning a response
func (h *ReplaceAllHTTPCheckHandlerImpl) Handle(parentType cnconstants.CnParentType, params http_check.ReplaceAllHTTPCheckBackendParams, principal interface{}) middleware.Responder {
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
		return http_check.NewReplaceAllHTTPCheckBackendDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_check.NewReplaceAllHTTPCheckBackendDefault(int(*e.Code)).WithPayload(e)
	}
	err = configuration.ReplaceHTTPChecks(string(parentType), params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return http_check.NewReplaceAllHTTPCheckBackendDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return http_check.NewReplaceAllHTTPCheckBackendDefault(int(*e.Code)).WithPayload(e)
			}
			return http_check.NewReplaceAllHTTPCheckBackendOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return http_check.NewReplaceAllHTTPCheckBackendAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return http_check.NewReplaceAllHTTPCheckBackendAccepted().WithPayload(params.Data)
}
