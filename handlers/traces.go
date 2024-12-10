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
	"github.com/haproxytech/dataplaneapi/operations/traces"
)

type GetTracesHandlerImpl struct {
	Client client_native.HAProxyClient
}

// Get the traces section (it is unique)
func (h *GetTracesHandlerImpl) Handle(params traces.GetTracesParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return traces.NewGetTracesDefault(int(*e.Code)).WithPayload(e)
	}

	var ts *models.Traces

	if params.FullSection != nil && *params.FullSection {
		_, ts, err = configuration.GetStructuredTraces(t)
	} else {
		_, ts, err = configuration.GetTraces(t)
	}

	if err != nil {
		e := misc.HandleError(err)
		return traces.NewGetTracesDefault(int(*e.Code)).WithPayload(e)
	}

	return traces.NewGetTracesOK().WithPayload(ts)
}

type CreateTracesHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

func (h *CreateTracesHandlerImpl) Handle(params traces.CreateTracesParams, principal interface{}) middleware.Responder {
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
		return traces.NewCreateTracesDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return traces.NewCreateTracesDefault(int(*e.Code)).WithPayload(e)
	}

	if err = configuration.CreateTraces(params.Data, t, v); err != nil {
		e := misc.HandleError(err)
		return traces.NewCreateTracesDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return traces.NewCreateTracesDefault(int(*e.Code)).WithPayload(e)
			}
			return traces.NewCreateTracesCreated().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return traces.NewCreateTracesAccepted().WithReloadID(rID).WithPayload(params.Data)
	}

	return traces.NewCreateTracesAccepted().WithPayload(params.Data)
}

type ReplaceTracesHandler struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

func (h *ReplaceTracesHandler) Handle(params traces.ReplaceTracesParams, principal interface{}) middleware.Responder {
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
		return traces.NewReplaceTracesDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return traces.NewReplaceTracesDefault(int(*e.Code)).WithPayload(e)
	}

	if params.FullSection != nil && *params.FullSection {
		err = configuration.PushStructuredTraces(params.Data, t, v)
	} else {
		err = configuration.EditTraces(params.Data, t, v)
	}

	if err != nil {
		e := misc.HandleError(err)
		return traces.NewReplaceTracesDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return traces.NewReplaceTracesDefault(int(*e.Code)).WithPayload(e)
			}
			return traces.NewReplaceTracesOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return traces.NewReplaceTracesAccepted().WithReloadID(rID)
	}

	return traces.NewReplaceTracesOK().WithPayload(params.Data)
}

type DeleteTracesHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

func (h *DeleteTracesHandlerImpl) Handle(params traces.DeleteTracesParams, principal interface{}) middleware.Responder {
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
		return traces.NewDeleteTracesDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return traces.NewDeleteTracesDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.DeleteTraces(t, v)
	if err != nil {
		e := misc.HandleError(err)
		return traces.NewDeleteTracesDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return traces.NewDeleteTracesDefault(int(*e.Code)).WithPayload(e)
			}
			return traces.NewDeleteTracesNoContent()
		}
		rID := h.ReloadAgent.Reload()
		return traces.NewDeleteTracesAccepted().WithReloadID(rID)
	}

	return traces.NewDeleteTracesAccepted()
}

// Trace Entries

type CreateTraceEntryHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

func (h *CreateTraceEntryHandlerImpl) Handle(params traces.CreateTraceEntryParams, principal interface{}) middleware.Responder {
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
		return traces.NewCreateTraceEntryDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return traces.NewCreateTraceEntryDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.CreateTraceEntry(params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return traces.NewCreateTraceEntryDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return traces.NewCreateTraceEntryDefault(int(*e.Code)).WithPayload(e)
			}
			return traces.NewCreateTraceEntryCreated().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return traces.NewCreateTraceEntryAccepted().WithReloadID(rID)
	}

	return traces.NewCreateTraceEntryCreated().WithPayload(params.Data)
}

type DeleteTraceEntryHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

func (h *DeleteTraceEntryHandlerImpl) Handle(params traces.DeleteTraceEntryParams, principal interface{}) middleware.Responder {
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
		return traces.NewDeleteTraceEntryDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return traces.NewDeleteTraceEntryDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.DeleteTraceEntry(params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return traces.NewDeleteTraceEntryDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return traces.NewDeleteTraceEntryDefault(int(*e.Code)).WithPayload(e)
			}
			return traces.NewDeleteTraceEntryNoContent()
		}
		rID := h.ReloadAgent.Reload()
		return traces.NewDeleteTraceEntryAccepted().WithReloadID(rID)
	}

	return traces.NewDeleteTraceEntryAccepted()
}
