// Copyright 2024 HAProxy Technologies
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
	"github.com/haproxytech/dataplaneapi/operations/crt_load"
)

type GetCrtLoadsHandlerImpl struct {
	Client client_native.HAProxyClient
}

// Get all loads from crt-stores
func (h *GetCrtLoadsHandlerImpl) Handle(params crt_load.GetCrtLoadsParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return crt_load.NewGetCrtLoadsDefault(int(*e.Code)).WithPayload(e)
	}

	_, stores, err := configuration.GetCrtLoads(params.CrtStore, t)
	if err != nil {
		e := misc.HandleError(err)
		return crt_load.NewGetCrtLoadsDefault(int(*e.Code)).WithPayload(e)
	}

	return crt_load.NewGetCrtLoadsOK().WithPayload(stores)
}

type GetCrtLoadHandlerImpl struct {
	Client client_native.HAProxyClient
}

// Get one load from a crt-store
func (h *GetCrtLoadHandlerImpl) Handle(params crt_load.GetCrtLoadParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return crt_load.NewGetCrtLoadDefault(int(*e.Code)).WithPayload(e)
	}

	_, store, err := configuration.GetCrtLoad(params.Certificate, params.CrtStore, t)
	if err != nil {
		e := misc.HandleError(err)
		return crt_load.NewGetCrtLoadDefault(int(*e.Code)).WithPayload(e)
	}

	return crt_load.NewGetCrtLoadOK().WithPayload(store)
}

type CreateCrtLoadHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

func (h *CreateCrtLoadHandlerImpl) Handle(params crt_load.CreateCrtLoadParams, principal interface{}) middleware.Responder {
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
		return crt_load.NewCreateCrtLoadDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return crt_load.NewCreateCrtLoadDefault(int(*e.Code)).WithPayload(e)
	}

	if err = configuration.CreateCrtLoad(params.CrtStore, params.Data, t, v); err != nil {
		e := misc.HandleError(err)
		return crt_load.NewCreateCrtLoadDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return crt_load.NewCreateCrtLoadDefault(int(*e.Code)).WithPayload(e)
			}
			return crt_load.NewCreateCrtLoadCreated().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return crt_load.NewCreateCrtLoadAccepted().WithReloadID(rID).WithPayload(params.Data)
	}

	return crt_load.NewCreateCrtLoadAccepted().WithPayload(params.Data)
}

type ReplaceCrtLoadHandler struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

func (h *ReplaceCrtLoadHandler) Handle(params crt_load.ReplaceCrtLoadParams, principal interface{}) middleware.Responder {
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
		return crt_load.NewReplaceCrtLoadDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return crt_load.NewReplaceCrtLoadDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.EditCrtLoad(params.Certificate, params.CrtStore, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return crt_load.NewReplaceCrtLoadDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return crt_load.NewReplaceCrtLoadDefault(int(*e.Code)).WithPayload(e)
			}
			return crt_load.NewReplaceCrtLoadOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return crt_load.NewReplaceCrtLoadAccepted().WithReloadID(rID)
	}

	return crt_load.NewReplaceCrtLoadOK().WithPayload(params.Data)
}

type DeleteCrtLoadHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

func (h *DeleteCrtLoadHandlerImpl) Handle(params crt_load.DeleteCrtLoadParams, principal interface{}) middleware.Responder {
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
		return crt_load.NewDeleteCrtLoadDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return crt_load.NewDeleteCrtLoadDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.DeleteCrtLoad(params.Certificate, params.CrtStore, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return crt_load.NewDeleteCrtLoadDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return crt_load.NewDeleteCrtLoadDefault(int(*e.Code)).WithPayload(e)
			}
			return crt_load.NewDeleteCrtLoadNoContent()
		}
		rID := h.ReloadAgent.Reload()
		return crt_load.NewDeleteCrtLoadAccepted().WithReloadID(rID)
	}

	return crt_load.NewDeleteCrtLoadAccepted()
}
