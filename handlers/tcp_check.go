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

// GetTCPChecksHandlerImpl implementation of the GetTCPChecksHandler interface using client-native client
type GetTCPChecksHandlerImpl struct {
	Client client_native.HAProxyClient
}

// ReplaceTCPCheckHandlerImpl implementation of the ReplaceTCPCheckHandler interface using client-native client
type ReplaceTCPCheckHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// ReplaceTCPChecksHandlerImpl implementation of the ReplaceTCPChecksHandler interface using client-native client
type ReplaceTCPChecksHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// Handle executing the request and returning a response
func (h *CreateTCPCheckHandlerImpl) Handle(params tcp_check.CreateTCPCheckParams, principal interface{}) middleware.Responder {
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
		return tcp_check.NewCreateTCPCheckDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return tcp_check.NewCreateTCPCheckDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.CreateTCPCheck(params.Index, params.ParentType, params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return tcp_check.NewCreateTCPCheckDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return tcp_check.NewCreateTCPCheckDefault(int(*e.Code)).WithPayload(e)
			}
			return tcp_check.NewCreateTCPCheckCreated().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return tcp_check.NewCreateTCPCheckAccepted().WithReloadID(rID).WithPayload(params.Data)
	}

	return tcp_check.NewCreateTCPCheckAccepted().WithPayload(params.Data)
}

// Handle executing the request and returning a response
func (h *DeleteTCPCheckHandlerImpl) Handle(params tcp_check.DeleteTCPCheckParams, principal interface{}) middleware.Responder {
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
		return tcp_check.NewDeleteTCPCheckDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return tcp_check.NewDeleteTCPCheckDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.DeleteTCPCheck(params.Index, params.ParentType, params.ParentName, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return tcp_check.NewDeleteTCPCheckDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return tcp_check.NewDeleteTCPCheckDefault(int(*e.Code)).WithPayload(e)
			}
			return tcp_check.NewDeleteTCPCheckNoContent()
		}
		rID := h.ReloadAgent.Reload()
		return tcp_check.NewCreateTCPCheckAccepted().WithReloadID(rID)
	}

	return tcp_check.NewDeleteTCPCheckAccepted()
}

// Handle executing the request and returning a response
func (h *GetTCPCheckHandlerImpl) Handle(params tcp_check.GetTCPCheckParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return tcp_check.NewGetTCPCheckDefault(int(*e.Code)).WithPayload(e)
	}

	_, data, err := configuration.GetTCPCheck(params.Index, params.ParentType, params.ParentName, t)
	if err != nil {
		e := misc.HandleError(err)
		return tcp_check.NewGetTCPCheckDefault(int(*e.Code)).WithPayload(e)
	}

	return tcp_check.NewGetTCPCheckOK().WithPayload(data)
}

// Handle executing the request and returning a response
func (h *GetTCPChecksHandlerImpl) Handle(params tcp_check.GetTCPChecksParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return tcp_check.NewGetTCPChecksDefault(int(*e.Code)).WithPayload(e)
	}

	_, data, err := configuration.GetTCPChecks(params.ParentType, params.ParentName, t)
	if err != nil {
		e := misc.HandleContainerGetError(err)
		if *e.Code == misc.ErrHTTPOk {
			return tcp_check.NewGetTCPChecksOK().WithPayload(models.TCPChecks{})
		}
		return tcp_check.NewGetTCPChecksDefault(int(*e.Code)).WithPayload(e)
	}

	return tcp_check.NewGetTCPChecksOK().WithPayload(data)
}

// Handle executing the request and returning a response
func (h *ReplaceTCPCheckHandlerImpl) Handle(params tcp_check.ReplaceTCPCheckParams, principal interface{}) middleware.Responder {
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
		return tcp_check.NewReplaceTCPCheckDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return tcp_check.NewReplaceTCPCheckDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.EditTCPCheck(params.Index, params.ParentType, params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return tcp_check.NewReplaceTCPCheckDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return tcp_check.NewReplaceTCPCheckDefault(int(*e.Code)).WithPayload(e)
			}
			return tcp_check.NewReplaceTCPCheckOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return tcp_check.NewReplaceTCPCheckAccepted().WithReloadID(rID).WithPayload(params.Data)
	}

	return tcp_check.NewReplaceTCPCheckAccepted().WithPayload(params.Data)
}

// Handle executing the request and returning a response
func (h *ReplaceTCPChecksHandlerImpl) Handle(params tcp_check.ReplaceTCPChecksParams, principal interface{}) middleware.Responder {
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
		return tcp_check.NewReplaceTCPChecksDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return tcp_check.NewReplaceTCPChecksDefault(int(*e.Code)).WithPayload(e)
	}
	err = configuration.ReplaceTCPChecks(params.ParentType, params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return tcp_check.NewReplaceTCPChecksDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return tcp_check.NewReplaceTCPChecksDefault(int(*e.Code)).WithPayload(e)
			}
			return tcp_check.NewReplaceTCPChecksOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return tcp_check.NewReplaceTCPChecksAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return tcp_check.NewReplaceTCPChecksAccepted().WithPayload(params.Data)
}
