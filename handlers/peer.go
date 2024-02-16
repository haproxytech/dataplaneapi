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
	"github.com/haproxytech/dataplaneapi/operations/peer"
)

// CreatePeerHandlerImpl implementation of the CreatePeerHandler interface using client-native client
type CreatePeerHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// DeletePeerHandlerImpl implementation of the DeletePeerHandler interface using client-native client
type DeletePeerHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// GetPeerHandlerImpl implementation of the GetPeerHandler interface using client-native client
type GetPeerHandlerImpl struct {
	Client client_native.HAProxyClient
}

// GetPeersHandlerImpl implementation of the GetPeersHandler interface using client-native client
type GetPeersHandlerImpl struct {
	Client client_native.HAProxyClient
}

// ReplacePeerHandlerImpl implementation of the ReplacePeerHandler interface using client-native client
type ReplacePeerHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// Handle executing the request and returning a response
func (h *CreatePeerHandlerImpl) Handle(params peer.CreatePeerParams, principal interface{}) middleware.Responder {
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
		return peer.NewCreatePeerDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return peer.NewCreatePeerDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.CreatePeerSection(params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return peer.NewCreatePeerDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return peer.NewCreatePeerDefault(int(*e.Code)).WithPayload(e)
			}
			return peer.NewCreatePeerCreated().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return peer.NewCreatePeerAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return peer.NewCreatePeerAccepted().WithPayload(params.Data)
}

// Handle executing the request and returning a response
func (h *DeletePeerHandlerImpl) Handle(params peer.DeletePeerParams, principal interface{}) middleware.Responder {
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
		return peer.NewDeletePeerDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return peer.NewDeletePeerDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.DeletePeerSection(params.Name, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return peer.NewDeletePeerDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return peer.NewDeletePeerDefault(int(*e.Code)).WithPayload(e)
			}
			return peer.NewDeletePeerNoContent()
		}
		rID := h.ReloadAgent.Reload()
		return peer.NewDeletePeerAccepted().WithReloadID(rID)
	}
	return peer.NewDeletePeerAccepted()
}

// Handle executing the request and returning a response
func (h *GetPeerHandlerImpl) Handle(params peer.GetPeerSectionParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return peer.NewGetPeerSectionDefault(int(*e.Code)).WithPayload(e)
	}

	_, p, err := configuration.GetPeerSection(params.Name, t)
	if err != nil {
		e := misc.HandleError(err)
		return peer.NewGetPeerSectionDefault(int(*e.Code)).WithPayload(e)
	}
	return peer.NewGetPeerSectionOK().WithPayload(p)
}

// Handle executing the request and returning a response
func (h *GetPeersHandlerImpl) Handle(params peer.GetPeerSectionsParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return peer.NewGetPeerSectionsDefault(int(*e.Code)).WithPayload(e)
	}

	_, ps, err := configuration.GetPeerSections(t)
	if err != nil {
		e := misc.HandleError(err)
		return peer.NewGetPeerSectionsDefault(int(*e.Code)).WithPayload(e)
	}
	return peer.NewGetPeerSectionsOK().WithPayload(ps)
}
