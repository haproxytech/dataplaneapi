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
	"github.com/haproxytech/dataplaneapi/operations/ring"
)

// CreateRingHandlerImpl implementation of the CreateRingHandler interface using client-native client
type CreateRingHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// DeleteRingHandlerImpl implementation of the DeleteRingHandler interface using client-native client
type DeleteRingHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// GetRingHandlerImpl implementation of the GetRingHandler interface using client-native client
type GetRingHandlerImpl struct {
	Client client_native.HAProxyClient
}

// GetRingsHandlerImpl implementation of the GetRingsHandler interface using client-native client
type GetRingsHandlerImpl struct {
	Client client_native.HAProxyClient
}

// ReplaceRingHandlerImpl implementation of the ReplaceRingHandler interface using client-native client
type ReplaceRingHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// Handle executing the request and returning a response
func (h *CreateRingHandlerImpl) Handle(params ring.CreateRingParams, principal interface{}) middleware.Responder {
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
		return ring.NewCreateRingDefault(int(*e.Code)).WithPayload(e)
	}

	err := h.createRing(params, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return ring.NewCreateRingDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return ring.NewCreateRingDefault(int(*e.Code)).WithPayload(e)
			}
			return ring.NewCreateRingCreated().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return ring.NewCreateRingAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return ring.NewCreateRingAccepted().WithPayload(params.Data)
}

func (h *CreateRingHandlerImpl) createRing(params ring.CreateRingParams, t string, v int64) error {
	configuration, err := h.Client.Configuration()
	if err != nil {
		return err
	}
	if params.FullSection != nil && *params.FullSection {
		return configuration.CreateStructuredRing(params.Data, t, v)
	}
	return configuration.CreateRing(params.Data, t, v)
}

// Handle executing the request and returning a response
func (h *DeleteRingHandlerImpl) Handle(params ring.DeleteRingParams, principal interface{}) middleware.Responder {
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
		return ring.NewDeleteRingDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return ring.NewDeleteRingDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.DeleteRing(params.Name, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return ring.NewDeleteRingDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return ring.NewDeleteRingDefault(int(*e.Code)).WithPayload(e)
			}
			return ring.NewDeleteRingNoContent()
		}
		rID := h.ReloadAgent.Reload()
		return ring.NewDeleteRingAccepted().WithReloadID(rID)
	}
	return ring.NewDeleteRingAccepted()
}

// Handle executing the request and returning a response
func (h *GetRingHandlerImpl) Handle(params ring.GetRingParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	_, bck, err := h.getRing(params, t)
	if err != nil {
		e := misc.HandleError(err)
		return ring.NewGetRingDefault(int(*e.Code)).WithPayload(e)
	}
	return ring.NewGetRingOK().WithPayload(bck)
}

func (h *GetRingHandlerImpl) getRing(params ring.GetRingParams, t string) (int64, *models.Ring, error) {
	configuration, err := h.Client.Configuration()
	if err != nil {
		return 0, nil, err
	}
	if params.FullSection != nil && *params.FullSection {
		return configuration.GetStructuredRing(params.Name, t)
	}
	return configuration.GetRing(params.Name, t)
}

// Handle executing the request and returning a response
func (h *GetRingsHandlerImpl) Handle(params ring.GetRingsParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	_, bcks, err := h.getRings(params, t)
	if err != nil {
		e := misc.HandleError(err)
		return ring.NewGetRingsDefault(int(*e.Code)).WithPayload(e)
	}
	return ring.NewGetRingsOK().WithPayload(bcks)
}

func (h *GetRingsHandlerImpl) getRings(params ring.GetRingsParams, t string) (int64, models.Rings, error) {
	configuration, err := h.Client.Configuration()
	if err != nil {
		return 0, nil, err
	}
	if params.FullSection != nil && *params.FullSection {
		return configuration.GetStructuredRings(t)
	}
	return configuration.GetRings(t)
}

// Handle executing the request and returning a response
func (h *ReplaceRingHandlerImpl) Handle(params ring.ReplaceRingParams, principal interface{}) middleware.Responder {
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
		return ring.NewReplaceRingDefault(int(*e.Code)).WithPayload(e)
	}

	err := h.editRing(params, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return ring.NewReplaceRingDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return ring.NewReplaceRingDefault(int(*e.Code)).WithPayload(e)
			}
			return ring.NewReplaceRingOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return ring.NewReplaceRingAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return ring.NewReplaceRingAccepted().WithPayload(params.Data)
}

func (h *ReplaceRingHandlerImpl) editRing(params ring.ReplaceRingParams, t string, v int64) error {
	configuration, err := h.Client.Configuration()
	if err != nil {
		return err
	}
	if params.FullSection != nil && *params.FullSection {
		return configuration.EditStructuredRing(params.Name, params.Data, t, v)
	}
	return configuration.EditRing(params.Name, params.Data, t, v)
}
