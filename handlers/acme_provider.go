// Copyright 2025 HAProxy Technologies
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
	"github.com/haproxytech/dataplaneapi/operations/acme"
)

type GetAcmeProvidersHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h *GetAcmeProvidersHandlerImpl) Handle(params acme.GetAcmeProvidersParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return acme.NewGetAcmeProvidersDefault(int(*e.Code)).WithPayload(e)
	}

	_, providers, err := configuration.GetAcmeProviders(t)
	if err != nil {
		e := misc.HandleError(err)
		return acme.NewGetAcmeProvidersDefault(int(*e.Code)).WithPayload(e)
	}

	return acme.NewGetAcmeProvidersOK().WithPayload(providers)
}

type GetAcmeProviderHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h *GetAcmeProviderHandlerImpl) Handle(params acme.GetAcmeProviderParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return acme.NewGetAcmeProviderDefault(int(*e.Code)).WithPayload(e)
	}

	_, provider, err := configuration.GetAcmeProvider(params.Name, t)
	if err != nil {
		e := misc.HandleError(err)
		return acme.NewGetAcmeProviderDefault(int(*e.Code)).WithPayload(e)
	}

	return acme.NewGetAcmeProviderOK().WithPayload(provider)
}

type CreateAcmeProviderHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

func (h *CreateAcmeProviderHandlerImpl) Handle(params acme.CreateAcmeProviderParams, principal interface{}) middleware.Responder {
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
		return acme.NewCreateAcmeProviderDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return acme.NewCreateAcmeProviderDefault(int(*e.Code)).WithPayload(e)
	}

	if err = configuration.CreateAcmeProvider(params.Data, t, v); err != nil {
		e := misc.HandleError(err)
		return acme.NewCreateAcmeProviderDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return acme.NewCreateAcmeProviderDefault(int(*e.Code)).WithPayload(e)
			}
			return acme.NewCreateAcmeProviderCreated().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return acme.NewCreateAcmeProviderAccepted().WithReloadID(rID).WithPayload(params.Data)
	}

	return acme.NewCreateAcmeProviderAccepted().WithPayload(params.Data)
}

type EditAcmeProviderHandler struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

func (h *EditAcmeProviderHandler) Handle(params acme.EditAcmeProviderParams, principal interface{}) middleware.Responder {
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
		return acme.NewEditAcmeProviderDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return acme.NewEditAcmeProviderDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.EditAcmeProvider(params.Name, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return acme.NewEditAcmeProviderDefault(int(*e.Code)).WithPayload(e)
	}

	return acme.NewEditAcmeProviderOK().WithPayload(params.Data)
}

type DeleteAcmeProviderHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

func (h *DeleteAcmeProviderHandlerImpl) Handle(params acme.DeleteAcmeProviderParams, principal interface{}) middleware.Responder {
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
		return acme.NewDeleteAcmeProviderDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return acme.NewDeleteAcmeProviderDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.DeleteAcmeProvider(params.Name, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return acme.NewDeleteAcmeProviderDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return acme.NewDeleteAcmeProviderDefault(int(*e.Code)).WithPayload(e)
			}
			return acme.NewDeleteAcmeProviderNoContent()
		}
		rID := h.ReloadAgent.Reload()
		return acme.NewDeleteAcmeProviderAccepted().WithReloadID(rID)
	}

	return acme.NewDeleteAcmeProviderAccepted()
}
