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
	"github.com/haproxytech/dataplaneapi/operations/crt_store"
)

type GetCrtStoresHandlerImpl struct {
	Client client_native.HAProxyClient
}

// Get all crt-stores
func (h *GetCrtStoresHandlerImpl) Handle(params crt_store.GetCrtStoresParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return crt_store.NewGetCrtStoresDefault(int(*e.Code)).WithPayload(e)
	}

	_, stores, err := configuration.GetCrtStores(t)
	if err != nil {
		e := misc.HandleError(err)
		return crt_store.NewGetCrtStoresDefault(int(*e.Code)).WithPayload(e)
	}

	return crt_store.NewGetCrtStoresOK().WithPayload(stores)
}

type GetCrtStoreHandlerImpl struct {
	Client client_native.HAProxyClient
}

// Get one crt-store
func (h *GetCrtStoreHandlerImpl) Handle(params crt_store.GetCrtStoreParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return crt_store.NewGetCrtStoreDefault(int(*e.Code)).WithPayload(e)
	}

	_, store, err := configuration.GetCrtStore(params.Name, t)
	if err != nil {
		e := misc.HandleError(err)
		return crt_store.NewGetCrtStoreDefault(int(*e.Code)).WithPayload(e)
	}

	return crt_store.NewGetCrtStoreOK().WithPayload(store)
}

type CreateCrtStoreHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

func (h *CreateCrtStoreHandlerImpl) Handle(params crt_store.CreateCrtStoreParams, principal interface{}) middleware.Responder {
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
		return crt_store.NewCreateCrtStoreDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return crt_store.NewCreateCrtStoreDefault(int(*e.Code)).WithPayload(e)
	}

	if err = configuration.CreateCrtStore(params.Data, t, v); err != nil {
		e := misc.HandleError(err)
		return crt_store.NewCreateCrtStoreDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return crt_store.NewCreateCrtStoreDefault(int(*e.Code)).WithPayload(e)
			}
			return crt_store.NewCreateCrtStoreCreated().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return crt_store.NewCreateCrtStoreAccepted().WithReloadID(rID).WithPayload(params.Data)
	}

	return crt_store.NewCreateCrtStoreAccepted().WithPayload(params.Data)
}

type EditCrtStoreHandler struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

func (h *EditCrtStoreHandler) Handle(params crt_store.EditCrtStoreParams, principal interface{}) middleware.Responder {
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
		return crt_store.NewEditCrtStoreDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return crt_store.NewEditCrtStoreDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.EditCrtStore(params.Name, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return crt_store.NewEditCrtStoreDefault(int(*e.Code)).WithPayload(e)
	}

	return crt_store.NewEditCrtStoreOK().WithPayload(params.Data)
}

type DeleteCrtStoreHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

func (h *DeleteCrtStoreHandlerImpl) Handle(params crt_store.DeleteCrtStoreParams, principal interface{}) middleware.Responder {
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
		return crt_store.NewDeleteCrtStoreDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return crt_store.NewDeleteCrtStoreDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.DeleteCrtStore(params.Name, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return crt_store.NewDeleteCrtStoreDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return crt_store.NewDeleteCrtStoreDefault(int(*e.Code)).WithPayload(e)
			}
			return crt_store.NewDeleteCrtStoreNoContent()
		}
		rID := h.ReloadAgent.Reload()
		return crt_store.NewDeleteCrtStoreAccepted().WithReloadID(rID)
	}

	return crt_store.NewDeleteCrtStoreAccepted()
}
