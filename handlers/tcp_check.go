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
	"github.com/haproxytech/dataplaneapi/operations/tcp_check"
)

// CreateTCPCheckHandlerImpl implementation of the CreateTCPCheckHandler interface using client-native client
type CreateTCPCheckHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// DeleteTCPCheckHandlerImpl implementation of the DeleteTCPCheckHandler interface using client-native client
type DeleteTCPCheckHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// GetTCPCheckHandlerImpl implementation of the GetTCPCheckHandler interface using client-native client
type GetTCPCheckHandlerImpl struct {
	Client client_native.HAProxyClient
}

// GetAllTCPCheckHandlerImpl implementation of the GetTCPChecksHandler interface using client-native client
type GetAllTCPCheckHandlerImpl struct {
	Client client_native.HAProxyClient
}

// ReplaceTCPCheckHandlerImpl implementation of the ReplaceTCPCheckHandler interface using client-native client
type ReplaceTCPCheckHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// ReplaceAllTCPCheckHandlerImpl implementation of the ReplaceTCPChecksHandler interface using client-native client
type ReplaceAllTCPCheckHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// Handle executing the request and returning a response
func (h *CreateTCPCheckHandlerImpl) Handle(parentType cnconstants.CnParentType, params tcp_check.CreateTCPCheckBackendParams, principal interface{}) middleware.Responder {
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
		return tcp_check.NewCreateTCPCheckBackendDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return tcp_check.NewCreateTCPCheckBackendDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.CreateTCPCheck(params.Index, string(parentType), params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return tcp_check.NewCreateTCPCheckBackendDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return tcp_check.NewCreateTCPCheckBackendDefault(int(*e.Code)).WithPayload(e)
			}
			return tcp_check.NewCreateTCPCheckBackendCreated().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return tcp_check.NewCreateTCPCheckBackendAccepted().WithReloadID(rID).WithPayload(params.Data)
	}

	return tcp_check.NewCreateTCPCheckBackendAccepted().WithPayload(params.Data)
}

// Handle executing the request and returning a response
func (h *DeleteTCPCheckHandlerImpl) Handle(parentType cnconstants.CnParentType, params tcp_check.DeleteTCPCheckBackendParams, principal interface{}) middleware.Responder {
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
		return tcp_check.NewDeleteTCPCheckBackendDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return tcp_check.NewDeleteTCPCheckBackendDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.DeleteTCPCheck(params.Index, string(parentType), params.ParentName, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return tcp_check.NewDeleteTCPCheckBackendDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return tcp_check.NewDeleteTCPCheckBackendDefault(int(*e.Code)).WithPayload(e)
			}
			return tcp_check.NewDeleteTCPCheckBackendNoContent()
		}
		rID := h.ReloadAgent.Reload()
		return tcp_check.NewCreateTCPCheckBackendAccepted().WithReloadID(rID)
	}

	return tcp_check.NewDeleteTCPCheckBackendAccepted()
}

// Handle executing the request and returning a response
func (h *GetTCPCheckHandlerImpl) Handle(parentType cnconstants.CnParentType, params tcp_check.GetTCPCheckBackendParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return tcp_check.NewGetTCPCheckBackendDefault(int(*e.Code)).WithPayload(e)
	}

	_, data, err := configuration.GetTCPCheck(params.Index, string(parentType), params.ParentName, t)
	if err != nil {
		e := misc.HandleError(err)
		return tcp_check.NewGetTCPCheckBackendDefault(int(*e.Code)).WithPayload(e)
	}

	return tcp_check.NewGetTCPCheckBackendOK().WithPayload(data)
}

// Handle executing the request and returning a response
func (h *GetAllTCPCheckHandlerImpl) Handle(parentType cnconstants.CnParentType, params tcp_check.GetAllTCPCheckBackendParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return tcp_check.NewGetAllTCPCheckBackendDefault(int(*e.Code)).WithPayload(e)
	}

	_, data, err := configuration.GetTCPChecks(string(parentType), params.ParentName, t)
	if err != nil {
		e := misc.HandleContainerGetError(err)
		if *e.Code == misc.ErrHTTPOk {
			return tcp_check.NewGetAllTCPCheckBackendOK().WithPayload(models.TCPChecks{})
		}
		return tcp_check.NewGetAllTCPCheckBackendDefault(int(*e.Code)).WithPayload(e)
	}

	return tcp_check.NewGetAllTCPCheckBackendOK().WithPayload(data)
}

// Handle executing the request and returning a response
func (h *ReplaceTCPCheckHandlerImpl) Handle(parentType cnconstants.CnParentType, params tcp_check.ReplaceTCPCheckBackendParams, principal interface{}) middleware.Responder {
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
		return tcp_check.NewReplaceTCPCheckBackendDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return tcp_check.NewReplaceTCPCheckBackendDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.EditTCPCheck(params.Index, string(parentType), params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return tcp_check.NewReplaceTCPCheckBackendDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return tcp_check.NewReplaceTCPCheckBackendDefault(int(*e.Code)).WithPayload(e)
			}
			return tcp_check.NewReplaceTCPCheckBackendOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return tcp_check.NewReplaceTCPCheckBackendAccepted().WithReloadID(rID).WithPayload(params.Data)
	}

	return tcp_check.NewReplaceTCPCheckBackendAccepted().WithPayload(params.Data)
}

// Handle executing the request and returning a response
func (h *ReplaceAllTCPCheckHandlerImpl) Handle(parentType cnconstants.CnParentType, params tcp_check.ReplaceAllTCPCheckBackendParams, principal interface{}) middleware.Responder {
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
		return tcp_check.NewReplaceAllTCPCheckBackendDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return tcp_check.NewReplaceAllTCPCheckBackendDefault(int(*e.Code)).WithPayload(e)
	}
	err = configuration.ReplaceTCPChecks(string(parentType), params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return tcp_check.NewReplaceAllTCPCheckBackendDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return tcp_check.NewReplaceAllTCPCheckBackendDefault(int(*e.Code)).WithPayload(e)
			}
			return tcp_check.NewReplaceAllTCPCheckBackendOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return tcp_check.NewReplaceAllTCPCheckBackendAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return tcp_check.NewReplaceAllTCPCheckBackendAccepted().WithPayload(params.Data)
}
