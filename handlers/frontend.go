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
	"github.com/haproxytech/dataplaneapi/operations/frontend"
)

// CreateFrontendHandlerImpl implementation of the CreateFrontendHandler interface using client-native client
type CreateFrontendHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// DeleteFrontendHandlerImpl implementation of the DeleteFrontendHandler interface using client-native client
type DeleteFrontendHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// GetFrontendHandlerImpl implementation of the GetFrontendHandler interface using client-native client
type GetFrontendHandlerImpl struct {
	Client client_native.HAProxyClient
}

// GetFrontendsHandlerImpl implementation of the GetFrontendsHandler interface using client-native client
type GetFrontendsHandlerImpl struct {
	Client client_native.HAProxyClient
}

// ReplaceFrontendHandlerImpl implementation of the ReplaceFrontendHandler interface using client-native client
type ReplaceFrontendHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// Handle executing the request and returning a response
func (h *CreateFrontendHandlerImpl) Handle(params frontend.CreateFrontendParams, principal interface{}) middleware.Responder {
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
		return frontend.NewCreateFrontendDefault(int(*e.Code)).WithPayload(e)
	}

	err := h.createFrontend(params, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return frontend.NewCreateFrontendDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return frontend.NewCreateFrontendDefault(int(*e.Code)).WithPayload(e)
			}
			return frontend.NewCreateFrontendCreated().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return frontend.NewCreateFrontendAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return frontend.NewCreateFrontendAccepted().WithPayload(params.Data)
}

func (h *CreateFrontendHandlerImpl) createFrontend(params frontend.CreateFrontendParams, t string, v int64) error {
	configuration, err := h.Client.Configuration()
	if err != nil {
		return err
	}
	if params.FullSection != nil && *params.FullSection {
		return configuration.CreateStructuredFrontend(params.Data, t, v)
	}
	return configuration.CreateFrontend(params.Data, t, v)
}

// Handle executing the request and returning a response
func (h *DeleteFrontendHandlerImpl) Handle(params frontend.DeleteFrontendParams, principal interface{}) middleware.Responder {
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
		return frontend.NewDeleteFrontendDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return frontend.NewDeleteFrontendDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.DeleteFrontend(params.Name, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return frontend.NewDeleteFrontendDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return frontend.NewDeleteFrontendDefault(int(*e.Code)).WithPayload(e)
			}
			return frontend.NewDeleteFrontendNoContent()
		}
		rID := h.ReloadAgent.Reload()
		return frontend.NewDeleteFrontendAccepted().WithReloadID(rID)
	}
	return frontend.NewDeleteFrontendAccepted()
}

// Handle executing the request and returning a response
func (h *GetFrontendHandlerImpl) Handle(params frontend.GetFrontendParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	_, f, err := h.getFrontend(params, t)
	if err != nil {
		e := misc.HandleError(err)
		return frontend.NewGetFrontendDefault(int(*e.Code)).WithPayload(e)
	}
	return frontend.NewGetFrontendOK().WithPayload(f)
}

func (h *GetFrontendHandlerImpl) getFrontend(params frontend.GetFrontendParams, t string) (int64, *models.Frontend, error) {
	configuration, err := h.Client.Configuration()
	if err != nil {
		return 0, nil, err
	}
	if params.FullSection != nil && *params.FullSection {
		return configuration.GetStructuredFrontend(params.Name, t)
	}
	return configuration.GetFrontend(params.Name, t)
}

// Handle executing the request and returning a response
func (h *GetFrontendsHandlerImpl) Handle(params frontend.GetFrontendsParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	_, fs, err := h.getFrontends(params, t)
	if err != nil {
		e := misc.HandleError(err)
		return frontend.NewGetFrontendsDefault(int(*e.Code)).WithPayload(e)
	}
	return frontend.NewGetFrontendsOK().WithPayload(fs)
}

func (h *GetFrontendsHandlerImpl) getFrontends(params frontend.GetFrontendsParams, t string) (int64, models.Frontends, error) {
	configuration, err := h.Client.Configuration()
	if err != nil {
		return 0, nil, err
	}
	if params.FullSection != nil && *params.FullSection {
		return configuration.GetStructuredFrontends(t)
	}
	return configuration.GetFrontends(t)
}

// Handle executing the request and returning a response
func (h *ReplaceFrontendHandlerImpl) Handle(params frontend.ReplaceFrontendParams, principal interface{}) middleware.Responder {
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
		return frontend.NewReplaceFrontendDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return frontend.NewReplaceFrontendDefault(int(*e.Code)).WithPayload(e)
	}

	_, ondisk, err := configuration.GetFrontend(params.Name, t)
	if err != nil {
		e := misc.HandleError(err)
		return frontend.NewReplaceFrontendDefault(int(*e.Code)).WithPayload(e)
	}

	err = h.editFrontend(params, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return frontend.NewReplaceFrontendDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		reload := changeThroughRuntimeAPI(*params.Data, *ondisk, "", h.Client)
		if reload {
			if *params.ForceReload {
				err := h.ReloadAgent.ForceReload()
				if err != nil {
					e := misc.HandleError(err)
					return frontend.NewReplaceFrontendDefault(int(*e.Code)).WithPayload(e)
				}
				return frontend.NewReplaceFrontendOK().WithPayload(params.Data)
			}
			rID := h.ReloadAgent.Reload()
			return frontend.NewReplaceFrontendAccepted().WithReloadID(rID).WithPayload(params.Data)
		}
		return frontend.NewReplaceFrontendOK().WithPayload(params.Data)
	}
	return frontend.NewReplaceFrontendAccepted().WithPayload(params.Data)
}

func (h *ReplaceFrontendHandlerImpl) editFrontend(params frontend.ReplaceFrontendParams, t string, v int64) error {
	configuration, err := h.Client.Configuration()
	if err != nil {
		return err
	}
	if params.FullSection != nil && *params.FullSection {
		return configuration.EditStructuredFrontend(params.Name, params.Data, t, v)
	}
	return configuration.EditFrontend(params.Name, params.Data, t, v)
}
