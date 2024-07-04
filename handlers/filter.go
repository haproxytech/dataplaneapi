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
	"github.com/haproxytech/dataplaneapi/operations/filter"
)

// CreateFilterHandlerImpl implementation of the CreateFilterHandler interface using client-native client
type CreateFilterHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// DeleteFilterHandlerImpl implementation of the DeleteFilterHandler interface using client-native client
type DeleteFilterHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// GetFilterHandlerImpl implementation of the GetFilterHandler interface using client-native client
type GetFilterHandlerImpl struct {
	Client client_native.HAProxyClient
}

// GetAllFilterHandlerImpl implementation of the GetAllFilterHandler interface using client-native client
type GetAllFilterHandlerImpl struct {
	Client client_native.HAProxyClient
}

// ReplaceFilterHandlerImpl implementation of the ReplaceFilterHandler interface using client-native client
type ReplaceFilterHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// ReplaceAllFilterHandlerImpl implementation of the ReplaceAllFilterHandler interface using client-native client
type ReplaceAllFilterHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// Handle executing the request and returning a response
func (h *CreateFilterHandlerImpl) Handle(parentType cnconstants.CnParentType, params filter.CreateFilterBackendParams, principal interface{}) middleware.Responder {
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
		return filter.NewCreateFilterBackendDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return filter.NewCreateFilterBackendDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.CreateFilter(params.Index, string(parentType), params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return filter.NewCreateFilterBackendDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return filter.NewCreateFilterBackendDefault(int(*e.Code)).WithPayload(e)
			}
			return filter.NewCreateFilterBackendCreated().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return filter.NewCreateFilterBackendAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return filter.NewCreateFilterBackendAccepted().WithPayload(params.Data)
}

// Handle executing the request and returning a response
func (h *DeleteFilterHandlerImpl) Handle(parentType cnconstants.CnParentType, params filter.DeleteFilterBackendParams, principal interface{}) middleware.Responder {
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
		return filter.NewDeleteFilterBackendDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return filter.NewDeleteFilterBackendDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.DeleteFilter(params.Index, string(parentType), params.ParentName, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return filter.NewDeleteFilterBackendDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return filter.NewDeleteFilterBackendDefault(int(*e.Code)).WithPayload(e)
			}
			return filter.NewDeleteFilterBackendNoContent()
		}
		rID := h.ReloadAgent.Reload()
		return filter.NewDeleteFilterBackendAccepted().WithReloadID(rID)
	}
	return filter.NewDeleteFilterBackendAccepted()
}

// Handle executing the request and returning a response
func (h *GetFilterHandlerImpl) Handle(parentType cnconstants.CnParentType, params filter.GetFilterBackendParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return filter.NewGetFilterBackendDefault(int(*e.Code)).WithPayload(e)
	}

	_, f, err := configuration.GetFilter(params.Index, string(parentType), params.ParentName, t)
	if err != nil {
		e := misc.HandleError(err)
		return filter.NewGetFilterBackendDefault(int(*e.Code)).WithPayload(e)
	}
	return filter.NewGetFilterBackendOK().WithPayload(f)
}

// Handle executing the request and returning a response
func (h *GetAllFilterHandlerImpl) Handle(parentType cnconstants.CnParentType, params filter.GetAllFilterBackendParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return filter.NewGetAllFilterBackendDefault(int(*e.Code)).WithPayload(e)
	}

	_, fs, err := configuration.GetFilters(string(parentType), params.ParentName, t)
	if err != nil {
		e := misc.HandleContainerGetError(err)
		if *e.Code == misc.ErrHTTPOk {
			return filter.NewGetAllFilterBackendOK().WithPayload(models.Filters{})
		}
		return filter.NewGetAllFilterBackendDefault(int(*e.Code)).WithPayload(e)
	}
	return filter.NewGetAllFilterBackendOK().WithPayload(fs)
}

// Handle executing the request and returning a response
func (h *ReplaceFilterHandlerImpl) Handle(parentType cnconstants.CnParentType, params filter.ReplaceFilterBackendParams, principal interface{}) middleware.Responder {
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
		return filter.NewReplaceFilterBackendDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return filter.NewGetFilterBackendDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.EditFilter(params.Index, string(parentType), params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return filter.NewReplaceFilterBackendDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return filter.NewReplaceFilterBackendDefault(int(*e.Code)).WithPayload(e)
			}
			return filter.NewReplaceFilterBackendOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return filter.NewReplaceFilterBackendAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return filter.NewReplaceFilterBackendAccepted().WithPayload(params.Data)
}

// Handle executing the request and returning a response
func (h *ReplaceAllFilterHandlerImpl) Handle(parentType cnconstants.CnParentType, params filter.ReplaceAllFilterBackendParams, principal interface{}) middleware.Responder {
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
		return filter.NewReplaceAllFilterBackendDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return filter.NewReplaceAllFilterBackendDefault(int(*e.Code)).WithPayload(e)
	}
	err = configuration.ReplaceFilters(string(parentType), params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return filter.NewReplaceAllFilterBackendDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return filter.NewReplaceAllFilterBackendDefault(int(*e.Code)).WithPayload(e)
			}
			return filter.NewReplaceAllFilterBackendOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return filter.NewReplaceAllFilterBackendAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return filter.NewReplaceAllFilterBackendAccepted().WithPayload(params.Data)
}
