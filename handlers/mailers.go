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
	"github.com/haproxytech/dataplaneapi/operations/mailers"
)

// CreateMailersSectionHandlerImpl implementation of the CreateMailersSectionHandler interface using client-native client
type CreateMailersSectionHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// DeleteMailersSectionHandlerImpl implementation of the DeleteMailersSectionHandler interface using client-native client
type DeleteMailersSectionHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// GetMailersSectionHandlerImpl implementation of the GetMailersSectionHandler interface using client-native client
type GetMailersSectionHandlerImpl struct {
	Client client_native.HAProxyClient
}

// GetMailersSectionsHandlerImpl implementation of the GetMailersSectionsHandler interface using client-native client
type GetMailersSectionsHandlerImpl struct {
	Client client_native.HAProxyClient
}

// EditMailersSectionHandlerImpl implementation of the EditMailersSectionHandler interface using client-native client
type EditMailersSectionHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// Handle executing the request and returning a response
func (h *CreateMailersSectionHandlerImpl) Handle(params mailers.CreateMailersSectionParams, principal interface{}) middleware.Responder {
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
		return mailers.NewCreateMailersSectionDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return mailers.NewCreateMailersSectionDefault(int(*e.Code)).WithPayload(e)
	}

	if err = configuration.CreateMailersSection(params.Data, t, v); err != nil {
		e := misc.HandleError(err)
		return mailers.NewCreateMailersSectionDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return mailers.NewCreateMailersSectionDefault(int(*e.Code)).WithPayload(e)
			}
			return mailers.NewCreateMailersSectionCreated().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return mailers.NewCreateMailersSectionAccepted().WithReloadID(rID).WithPayload(params.Data)
	}

	return mailers.NewCreateMailersSectionAccepted().WithPayload(params.Data)
}

// Handle executing the request and returning a response
func (h *DeleteMailersSectionHandlerImpl) Handle(params mailers.DeleteMailersSectionParams, principal interface{}) middleware.Responder {
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
		return mailers.NewDeleteMailersSectionDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return mailers.NewDeleteMailersSectionDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.DeleteMailersSection(params.Name, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return mailers.NewDeleteMailersSectionDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return mailers.NewDeleteMailersSectionDefault(int(*e.Code)).WithPayload(e)
			}
			return mailers.NewDeleteMailersSectionNoContent()
		}
		rID := h.ReloadAgent.Reload()
		return mailers.NewDeleteMailersSectionAccepted().WithReloadID(rID)
	}

	return mailers.NewDeleteMailersSectionAccepted()
}

// Handle executing the request and returning a response
func (h *GetMailersSectionHandlerImpl) Handle(params mailers.GetMailersSectionParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return mailers.NewGetMailersSectionDefault(int(*e.Code)).WithPayload(e)
	}

	_, ms, err := configuration.GetMailersSection(params.Name, t)
	if err != nil {
		e := misc.HandleError(err)
		return mailers.NewGetMailersSectionDefault(int(*e.Code)).WithPayload(e)
	}

	return mailers.NewGetMailersSectionOK().WithPayload(ms)
}

// Handle executing the request and returning a response
func (h *GetMailersSectionsHandlerImpl) Handle(params mailers.GetMailersSectionsParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return mailers.NewGetMailersSectionsDefault(int(*e.Code)).WithPayload(e)
	}

	_, ms, err := configuration.GetMailersSections(t)
	if err != nil {
		e := misc.HandleError(err)
		return mailers.NewGetMailersSectionsDefault(int(*e.Code)).WithPayload(e)
	}

	return mailers.NewGetMailersSectionsOK().WithPayload(ms)
}

func (h *EditMailersSectionHandlerImpl) Handle(params mailers.EditMailersSectionParams, principal interface{}) middleware.Responder {
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
		return mailers.NewEditMailersSectionDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return mailers.NewEditMailersSectionDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.EditMailersSection(params.Name, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return mailers.NewEditMailersSectionDefault(int(*e.Code)).WithPayload(e)
	}

	return mailers.NewEditMailersSectionOK().WithPayload(params.Data)
}
