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
	client_native "github.com/haproxytech/client-native/v5"
	"github.com/haproxytech/client-native/v5/models"

	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/operations/log_forward"
)

// CreateLogForwardHandlerImpl implementation of the CreateLogForwardHandler interface using client-native client
type CreateLogForwardHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// DeleteLogForwardHandlerImpl implementation of the DeleteLogForwardHandler interface using client-native client
type DeleteLogForwardHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// GetLogForwardHandlerImpl implementation of the GetLogForwardHandler interface using client-native client
type GetLogForwardHandlerImpl struct {
	Client client_native.HAProxyClient
}

// GetLogForwardsHandlerImpl implementation of the GetLogForwardsHandler interface using client-native client
type GetLogForwardsHandlerImpl struct {
	Client client_native.HAProxyClient
}

// ReplaceLogForwardHandlerImpl implementation of the ReplaceLogForwardHandler interface using client-native client
type ReplaceLogForwardHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// Handle executing the request and returning a response
func (h *CreateLogForwardHandlerImpl) Handle(params log_forward.CreateLogForwardParams, principal interface{}) middleware.Responder {
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
		return log_forward.NewCreateLogForwardDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return log_forward.NewCreateLogForwardDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.CreateLogForward(params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return log_forward.NewCreateLogForwardDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return log_forward.NewCreateLogForwardDefault(int(*e.Code)).WithPayload(e)
			}
			return log_forward.NewCreateLogForwardCreated().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return log_forward.NewCreateLogForwardAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return log_forward.NewCreateLogForwardAccepted().WithPayload(params.Data)
}

// Handle executing the request and returning a response
func (h *DeleteLogForwardHandlerImpl) Handle(params log_forward.DeleteLogForwardParams, principal interface{}) middleware.Responder {
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
		return log_forward.NewDeleteLogForwardDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return log_forward.NewDeleteLogForwardDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.DeleteLogForward(params.Name, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return log_forward.NewDeleteLogForwardDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return log_forward.NewDeleteLogForwardDefault(int(*e.Code)).WithPayload(e)
			}
			return log_forward.NewDeleteLogForwardNoContent()
		}
		rID := h.ReloadAgent.Reload()
		return log_forward.NewDeleteLogForwardAccepted().WithReloadID(rID)
	}
	return log_forward.NewDeleteLogForwardAccepted()
}

// Handle executing the request and returning a response
func (h *GetLogForwardHandlerImpl) Handle(params log_forward.GetLogForwardParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return log_forward.NewGetLogForwardDefault(int(*e.Code)).WithPayload(e)
	}

	v, bck, err := configuration.GetLogForward(params.Name, t)
	if err != nil {
		e := misc.HandleError(err)
		return log_forward.NewGetLogForwardDefault(int(*e.Code)).WithPayload(e)
	}
	return log_forward.NewGetLogForwardOK().WithPayload(&log_forward.GetLogForwardOKBody{Version: v, Data: bck})
}

// Handle executing the request and returning a response
func (h *GetLogForwardsHandlerImpl) Handle(params log_forward.GetLogForwardsParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return log_forward.NewGetLogForwardsDefault(int(*e.Code)).WithPayload(e)
	}

	v, bcks, err := configuration.GetLogForwards(t)
	if err != nil {
		e := misc.HandleError(err)
		return log_forward.NewGetLogForwardsDefault(int(*e.Code)).WithPayload(e)
	}
	return log_forward.NewGetLogForwardsOK().WithPayload(&log_forward.GetLogForwardsOKBody{Version: v, Data: bcks})
}

// Handle executing the request and returning a response
func (h *ReplaceLogForwardHandlerImpl) Handle(params log_forward.ReplaceLogForwardParams, principal interface{}) middleware.Responder {
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
		return log_forward.NewReplaceLogForwardDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return log_forward.NewReplaceLogForwardDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.EditLogForward(params.Name, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return log_forward.NewReplaceLogForwardDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return log_forward.NewReplaceLogForwardDefault(int(*e.Code)).WithPayload(e)
			}
			return log_forward.NewReplaceLogForwardOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return log_forward.NewReplaceLogForwardAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return log_forward.NewReplaceLogForwardAccepted().WithPayload(params.Data)
}
