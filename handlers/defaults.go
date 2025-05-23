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
	"github.com/haproxytech/dataplaneapi/operations/defaults"
)

// GetDefaultsHandlerImpl implementation of the GetDefaultsHandler interface
type GetDefaultsHandlerImpl struct {
	Client client_native.HAProxyClient
}

// ReplaceDefaultsHandlerImpl implementation of the ReplaceDefaultsHandler interface
type ReplaceDefaultsHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// Handle executing the request and returning a response
func (h *GetDefaultsHandlerImpl) Handle(params defaults.GetDefaultsSectionParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	_, data, err := h.getDefaultsConfiguration(params, t)
	if err != nil {
		e := misc.HandleError(err)
		return defaults.NewGetDefaultsSectionDefault(int(*e.Code)).WithPayload(e)
	}
	return defaults.NewGetDefaultsSectionOK().WithPayload(data)
}

func (h *GetDefaultsHandlerImpl) getDefaultsConfiguration(params defaults.GetDefaultsSectionParams, t string) (int64, *models.Defaults, error) {
	configuration, err := h.Client.Configuration()
	if err != nil {
		return 0, nil, err
	}
	if params.FullSection != nil && *params.FullSection {
		return configuration.GetStructuredDefaultsConfiguration(t)
	}
	return configuration.GetDefaultsConfiguration(t)
}

// Handle executing the request and returning a response
func (h *ReplaceDefaultsHandlerImpl) Handle(params defaults.ReplaceDefaultsSectionParams, principal interface{}) middleware.Responder {
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
		return defaults.NewReplaceDefaultsSectionDefault(int(*e.Code)).WithPayload(e)
	}

	err := h.pushDefaultsConfiguration(params, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return defaults.NewReplaceDefaultsSectionDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return defaults.NewReplaceDefaultsSectionDefault(int(*e.Code)).WithPayload(e)
			}
			return defaults.NewReplaceDefaultsSectionOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return defaults.NewReplaceDefaultsSectionAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return defaults.NewReplaceDefaultsSectionAccepted().WithPayload(params.Data)
}

func (h *ReplaceDefaultsHandlerImpl) pushDefaultsConfiguration(params defaults.ReplaceDefaultsSectionParams, t string, v int64) error {
	configuration, err := h.Client.Configuration()
	if err != nil {
		return err
	}
	if params.FullSection != nil && *params.FullSection {
		return configuration.PushStructuredDefaultsConfiguration(params.Data, t, v)
	}
	return configuration.PushDefaultsConfiguration(params.Data, t, v)
}

// GetDefaultsHandlerImpl implementation of the GetDefaultsHandler interface
type GetDefaultsSectionsHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h GetDefaultsSectionsHandlerImpl) Handle(params defaults.GetDefaultsSectionsParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return defaults.NewGetDefaultsSectionsDefault(int(*e.Code)).WithPayload(e)
	}

	_, fs, err := configuration.GetDefaultsSections(t)
	if err != nil {
		e := misc.HandleError(err)
		return defaults.NewGetDefaultsSectionsDefault(int(*e.Code)).WithPayload(e)
	}
	return defaults.NewGetDefaultsSectionsOK().WithPayload(fs)
}

func (h *GetDefaultsSectionsHandlerImpl) getDefaultsSections(params defaults.GetDefaultsSectionParams, t string) (int64, models.DefaultsSections, error) {
	configuration, err := h.Client.Configuration()
	if err != nil {
		return 0, nil, err
	}
	if params.FullSection != nil && *params.FullSection {
		return configuration.GetStructuredDefaultsSections(t)
	}
	return configuration.GetDefaultsSections(t)
}

type GetDefaultsSectionHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h GetDefaultsSectionHandlerImpl) Handle(params defaults.GetDefaultsSectionParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	_, f, err := h.getDefaultsSection(params, t)
	if err != nil {
		e := misc.HandleError(err)
		return defaults.NewGetDefaultsSectionDefault(int(*e.Code)).WithPayload(e)
	}
	return defaults.NewGetDefaultsSectionOK().WithPayload(f)
}

func (h *GetDefaultsSectionHandlerImpl) getDefaultsSection(params defaults.GetDefaultsSectionParams, t string) (int64, *models.Defaults, error) {
	configuration, err := h.Client.Configuration()
	if err != nil {
		return 0, nil, err
	}
	if params.FullSection != nil && *params.FullSection {
		return configuration.GetStructuredDefaultsSection(params.Name, t)
	}
	return configuration.GetDefaultsSection(params.Name, t)
}

type CreateDefaultsSectionHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

func (h CreateDefaultsSectionHandlerImpl) Handle(params defaults.CreateDefaultsSectionParams, principal interface{}) middleware.Responder {
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
		return defaults.NewCreateDefaultsSectionDefault(int(*e.Code)).WithPayload(e)
	}

	err := h.createDefaultsSection(params, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return defaults.NewCreateDefaultsSectionDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return defaults.NewCreateDefaultsSectionDefault(int(*e.Code)).WithPayload(e)
			}
			return defaults.NewCreateDefaultsSectionCreated().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return defaults.NewCreateDefaultsSectionAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return defaults.NewCreateDefaultsSectionAccepted().WithPayload(params.Data)
}

func (h *CreateDefaultsSectionHandlerImpl) createDefaultsSection(params defaults.CreateDefaultsSectionParams, t string, v int64) error {
	configuration, err := h.Client.Configuration()
	if err != nil {
		return err
	}
	if params.FullSection != nil && *params.FullSection {
		return configuration.CreateStructuredDefaultsSection(params.Data, t, v)
	}
	return configuration.CreateDefaultsSection(params.Data, t, v)
}

type AddDefaultsSectionHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

func (h AddDefaultsSectionHandlerImpl) Handle(params defaults.AddDefaultsSectionParams, principal interface{}) middleware.Responder {
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
		return defaults.NewAddDefaultsSectionDefault(int(*e.Code)).WithPayload(e)
	}

	err := h.createDefaultsSection(params, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return defaults.NewAddDefaultsSectionDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return defaults.NewAddDefaultsSectionDefault(int(*e.Code)).WithPayload(e)
			}
			return defaults.NewAddDefaultsSectionCreated().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return defaults.NewAddDefaultsSectionAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return defaults.NewAddDefaultsSectionAccepted().WithPayload(params.Data)
}

func (h *AddDefaultsSectionHandlerImpl) createDefaultsSection(params defaults.AddDefaultsSectionParams, t string, v int64) error {
	configuration, err := h.Client.Configuration()
	if err != nil {
		return err
	}
	if params.FullSection != nil && *params.FullSection {
		return configuration.CreateStructuredDefaultsSection(params.Data, t, v)
	}
	return configuration.CreateDefaultsSection(params.Data, t, v)
}

// ReplaceDefaultsHandlerImpl implementation of the ReplaceDefaultsHandler interface
type ReplaceDefaultsSectionHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

func (h ReplaceDefaultsSectionHandlerImpl) Handle(params defaults.ReplaceDefaultsSectionParams, principal interface{}) middleware.Responder {
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
		return defaults.NewReplaceDefaultsSectionDefault(int(*e.Code)).WithPayload(e)
	}

	err := h.editDefaultsSection(params, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return defaults.NewReplaceDefaultsSectionDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return defaults.NewReplaceDefaultsSectionDefault(int(*e.Code)).WithPayload(e)
			}
			return defaults.NewReplaceDefaultsSectionOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return defaults.NewReplaceDefaultsSectionAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return defaults.NewReplaceDefaultsSectionAccepted().WithPayload(params.Data)
}

func (h *ReplaceDefaultsSectionHandlerImpl) editDefaultsSection(params defaults.ReplaceDefaultsSectionParams, t string, v int64) error {
	configuration, err := h.Client.Configuration()
	if err != nil {
		return err
	}
	if params.FullSection != nil && *params.FullSection {
		return configuration.EditStructuredDefaultsSection(params.Name, params.Data, t, v)
	}
	return configuration.EditDefaultsSection(params.Name, params.Data, t, v)
}

type DeleteDefaultsSectionHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

func (h DeleteDefaultsSectionHandlerImpl) Handle(params defaults.DeleteDefaultsSectionParams, principal interface{}) middleware.Responder {
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
		return defaults.NewDeleteDefaultsSectionDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return defaults.NewDeleteDefaultsSectionDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.DeleteDefaultsSection(params.Name, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return defaults.NewDeleteDefaultsSectionDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return defaults.NewDeleteDefaultsSectionDefault(int(*e.Code)).WithPayload(e)
			}
			return defaults.NewDeleteDefaultsSectionNoContent()
		}
		rID := h.ReloadAgent.Reload()
		return defaults.NewDeleteDefaultsSectionAccepted().WithReloadID(rID)
	}
	return defaults.NewDeleteDefaultsSectionAccepted()
}
