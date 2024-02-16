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
	"github.com/haproxytech/dataplaneapi/operations/nameserver"
)

// CreateNameserverHandlerImpl implementation of the CreateNameserverHandler interface using client-native client
type CreateNameserverHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// DeleteNameserverHandlerImpl implementation of the DeleteNameserverHandler interface using client-native client
type DeleteNameserverHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// GetNameserverHandlerImpl implementation of the GetNameserverHandler interface using client-native client
type GetNameserverHandlerImpl struct {
	Client client_native.HAProxyClient
}

// GetNameserversHandlerImpl implementation of the GetNameserversHandler interface using client-native client
type GetNameserversHandlerImpl struct {
	Client client_native.HAProxyClient
}

// ReplaceNameserverHandlerImpl implementation of the ReplaceNameserverHandler interface using client-native client
type ReplaceNameserverHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// Handle executing the request and returning a response
func (h *CreateNameserverHandlerImpl) Handle(params nameserver.CreateNameserverParams, principal interface{}) middleware.Responder {
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
		return nameserver.NewCreateNameserverDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return nameserver.NewCreateNameserverDefault(int(*e.Code)).WithPayload(e)
	}
	err = configuration.CreateNameserver(params.Resolver, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return nameserver.NewCreateNameserverDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return nameserver.NewCreateNameserverDefault(int(*e.Code)).WithPayload(e)
			}
			return nameserver.NewCreateNameserverCreated().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return nameserver.NewCreateNameserverAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return nameserver.NewCreateNameserverAccepted().WithPayload(params.Data)
}

// Handle executing the request and returning a response
func (h *DeleteNameserverHandlerImpl) Handle(params nameserver.DeleteNameserverParams, principal interface{}) middleware.Responder {
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
		return nameserver.NewDeleteNameserverDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return nameserver.NewDeleteNameserverDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.DeleteNameserver(params.Name, params.Resolver, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return nameserver.NewDeleteNameserverDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return nameserver.NewDeleteNameserverDefault(int(*e.Code)).WithPayload(e)
			}
			return nameserver.NewDeleteNameserverNoContent()
		}
		rID := h.ReloadAgent.Reload()
		return nameserver.NewDeleteNameserverAccepted().WithReloadID(rID)
	}
	return nameserver.NewDeleteNameserverAccepted()
}

// Handle executing the request and returning a response
func (h *GetNameserverHandlerImpl) Handle(params nameserver.GetNameserverParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return nameserver.NewGetNameserverDefault(int(*e.Code)).WithPayload(e)
	}

	_, b, err := configuration.GetNameserver(params.Name, params.Resolver, t)
	if err != nil {
		e := misc.HandleError(err)
		return nameserver.NewGetNameserverDefault(int(*e.Code)).WithPayload(e)
	}
	return nameserver.NewGetNameserverOK().WithPayload(b)
}

// Handle executing the request and returning a response
func (h *GetNameserversHandlerImpl) Handle(params nameserver.GetNameserversParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return nameserver.NewGetNameserversDefault(int(*e.Code)).WithPayload(e)
	}

	_, bs, err := configuration.GetNameservers(params.Resolver, t)
	if err != nil {
		e := misc.HandleError(err)
		return nameserver.NewGetNameserversDefault(int(*e.Code)).WithPayload(e)
	}
	return nameserver.NewGetNameserversOK().WithPayload(bs)
}

// Handle executing the request and returning a response
func (h *ReplaceNameserverHandlerImpl) Handle(params nameserver.ReplaceNameserverParams, principal interface{}) middleware.Responder {
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
		return nameserver.NewReplaceNameserverDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return nameserver.NewReplaceNameserverDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.EditNameserver(params.Name, params.Resolver, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return nameserver.NewReplaceNameserverDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return nameserver.NewReplaceNameserverDefault(int(*e.Code)).WithPayload(e)
			}
			return nameserver.NewReplaceNameserverOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return nameserver.NewReplaceNameserverAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return nameserver.NewReplaceNameserverAccepted().WithPayload(params.Data)
}
