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
	"github.com/haproxytech/dataplaneapi/operations/dgram_bind"
)

// CreateDgramBindHandlerImpl implementation of the CreateDgramBindHandler interface using client-native client
type CreateDgramBindHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// DeleteDgramBindHandlerImpl implementation of the DeleteDgramBindHandler interface using client-native client
type DeleteDgramBindHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// GetDgramBindHandlerImpl implementation of the GetDgramBindHandler interface using client-native client
type GetDgramBindHandlerImpl struct {
	Client client_native.HAProxyClient
}

// GetDgramBindsHandlerImpl implementation of the GetDgramBindsHandler interface using client-native client
type GetDgramBindsHandlerImpl struct {
	Client client_native.HAProxyClient
}

// ReplaceDgramBindHandlerImpl implementation of the ReplaceDgramBindHandler interface using client-native client
type ReplaceDgramBindHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// Handle executing the request and returning a response
func (h *CreateDgramBindHandlerImpl) Handle(params dgram_bind.CreateDgramBindParams, principal interface{}) middleware.Responder {
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
		return dgram_bind.NewCreateDgramBindDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return dgram_bind.NewCreateDgramBindDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.CreateDgramBind(params.Data.Name, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return dgram_bind.NewCreateDgramBindDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return dgram_bind.NewCreateDgramBindDefault(int(*e.Code)).WithPayload(e)
			}
			return dgram_bind.NewCreateDgramBindCreated().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return dgram_bind.NewCreateDgramBindAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return dgram_bind.NewCreateDgramBindAccepted().WithPayload(params.Data)
}

// Handle executing the request and returning a response
func (h *DeleteDgramBindHandlerImpl) Handle(params dgram_bind.DeleteDgramBindParams, principal interface{}) middleware.Responder {
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
		return dgram_bind.NewDeleteDgramBindDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return dgram_bind.NewDeleteDgramBindDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.DeleteDgramBind(params.Name, params.LogForward, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return dgram_bind.NewDeleteDgramBindDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return dgram_bind.NewDeleteDgramBindDefault(int(*e.Code)).WithPayload(e)
			}
			return dgram_bind.NewDeleteDgramBindNoContent()
		}
		rID := h.ReloadAgent.Reload()
		return dgram_bind.NewDeleteDgramBindAccepted().WithReloadID(rID)
	}
	return dgram_bind.NewDeleteDgramBindAccepted()
}

// Handle executing the request and returning a response
func (h *GetDgramBindHandlerImpl) Handle(params dgram_bind.GetDgramBindParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return dgram_bind.NewGetDgramBindDefault(int(*e.Code)).WithPayload(e)
	}

	v, bck, err := configuration.GetDgramBind(params.Name, params.LogForward, t)
	if err != nil {
		e := misc.HandleError(err)
		return dgram_bind.NewGetDgramBindDefault(int(*e.Code)).WithPayload(e)
	}
	return dgram_bind.NewGetDgramBindOK().WithPayload(&dgram_bind.GetDgramBindOKBody{Version: v, Data: bck})
}

// Handle executing the request and returning a response
func (h *GetDgramBindsHandlerImpl) Handle(params dgram_bind.GetDgramBindsParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return dgram_bind.NewGetDgramBindsDefault(int(*e.Code)).WithPayload(e)
	}

	v, bcks, err := configuration.GetDgramBinds(params.LogForward, t)
	if err != nil {
		e := misc.HandleError(err)
		return dgram_bind.NewGetDgramBindsDefault(int(*e.Code)).WithPayload(e)
	}
	return dgram_bind.NewGetDgramBindsOK().WithPayload(&dgram_bind.GetDgramBindsOKBody{Version: v, Data: bcks})
}

// Handle executing the request and returning a response
func (h *ReplaceDgramBindHandlerImpl) Handle(params dgram_bind.ReplaceDgramBindParams, principal interface{}) middleware.Responder {
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
		return dgram_bind.NewReplaceDgramBindDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return dgram_bind.NewReplaceDgramBindDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.EditDgramBind(params.Name, params.LogForward, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return dgram_bind.NewReplaceDgramBindDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return dgram_bind.NewReplaceDgramBindDefault(int(*e.Code)).WithPayload(e)
			}
			return dgram_bind.NewReplaceDgramBindOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return dgram_bind.NewReplaceDgramBindAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return dgram_bind.NewReplaceDgramBindAccepted().WithPayload(params.Data)
}
