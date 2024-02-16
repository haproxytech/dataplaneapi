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
	"github.com/haproxytech/dataplaneapi/operations/peer_entry"
)

// CreatePeerEntryHandlerImpl implementation of the CreatePeerEntryHandler interface using client-native client
type CreatePeerEntryHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// DeletePeerEntryHandlerImpl implementation of the DeletePeerEntryHandler interface using client-native client
type DeletePeerEntryHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// GetPeerEntryHandlerImpl implementation of the GetPeerEntryHandler interface using client-native client
type GetPeerEntryHandlerImpl struct {
	Client client_native.HAProxyClient
}

// GetPeerEntriesHandlerImpl implementation of the GetPeerEntriesHandler interface using client-native client
type GetPeerEntriesHandlerImpl struct {
	Client client_native.HAProxyClient
}

// ReplacePeerEntryHandlerImpl implementation of the ReplacePeerEntryHandler interface using client-native client
type ReplacePeerEntryHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// Handle executing the request and returning a response
func (h *CreatePeerEntryHandlerImpl) Handle(params peer_entry.CreatePeerEntryParams, principal interface{}) middleware.Responder {
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
		return peer_entry.NewCreatePeerEntryDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return peer_entry.NewCreatePeerEntryDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.CreatePeerEntry(params.PeerSection, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return peer_entry.NewCreatePeerEntryDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return peer_entry.NewCreatePeerEntryDefault(int(*e.Code)).WithPayload(e)
			}
			return peer_entry.NewCreatePeerEntryCreated().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return peer_entry.NewCreatePeerEntryAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return peer_entry.NewCreatePeerEntryAccepted().WithPayload(params.Data)
}

// Handle executing the request and returning a response
func (h *DeletePeerEntryHandlerImpl) Handle(params peer_entry.DeletePeerEntryParams, principal interface{}) middleware.Responder {
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
		return peer_entry.NewDeletePeerEntryDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return peer_entry.NewDeletePeerEntryDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.DeletePeerEntry(params.Name, params.PeerSection, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return peer_entry.NewDeletePeerEntryDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return peer_entry.NewDeletePeerEntryDefault(int(*e.Code)).WithPayload(e)
			}
			return peer_entry.NewDeletePeerEntryNoContent()
		}
		rID := h.ReloadAgent.Reload()
		return peer_entry.NewDeletePeerEntryAccepted().WithReloadID(rID)
	}
	return peer_entry.NewDeletePeerEntryAccepted()
}

// Handle executing the request and returning a response
func (h *GetPeerEntryHandlerImpl) Handle(params peer_entry.GetPeerEntryParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return peer_entry.NewGetPeerEntryDefault(int(*e.Code)).WithPayload(e)
	}

	_, pe, err := configuration.GetPeerEntry(params.Name, params.PeerSection, t)
	if err != nil {
		e := misc.HandleError(err)
		return peer_entry.NewGetPeerEntryDefault(int(*e.Code)).WithPayload(e)
	}
	return peer_entry.NewGetPeerEntryOK().WithPayload(pe)
}

// Handle executing the request and returning a response
func (h *GetPeerEntriesHandlerImpl) Handle(params peer_entry.GetPeerEntriesParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return peer_entry.NewGetPeerEntriesDefault(int(*e.Code)).WithPayload(e)
	}

	_, pes, err := configuration.GetPeerEntries(params.PeerSection, t)
	if err != nil {
		e := misc.HandleError(err)
		return peer_entry.NewGetPeerEntriesDefault(int(*e.Code)).WithPayload(e)
	}
	return peer_entry.NewGetPeerEntriesOK().WithPayload(pes)
}

// Handle executing the request and returning a response
func (h *ReplacePeerEntryHandlerImpl) Handle(params peer_entry.ReplacePeerEntryParams, principal interface{}) middleware.Responder {
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
		return peer_entry.NewReplacePeerEntryDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return peer_entry.NewReplacePeerEntryDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.EditPeerEntry(params.Name, params.PeerSection, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return peer_entry.NewReplacePeerEntryDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return peer_entry.NewReplacePeerEntryDefault(int(*e.Code)).WithPayload(e)
			}
			return peer_entry.NewReplacePeerEntryOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return peer_entry.NewReplacePeerEntryAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return peer_entry.NewReplacePeerEntryAccepted().WithPayload(params.Data)
}
