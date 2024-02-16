// Copyright 2023 HAProxy Technologies
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
	"github.com/haproxytech/dataplaneapi/operations/table"
)

// CreateTableHandlerImpl implementation of the CreateTableHandler interface using client-native client
type CreateTableHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// DeleteTableHandlerImpl implementation of the DeleteTableHandler interface using client-native client
type DeleteTableHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// GetTableHandlerImpl implementation of the GetTableHandler interface using client-native client
type GetTableHandlerImpl struct {
	Client client_native.HAProxyClient
}

// GetTablesHandlerImpl implementation of the GetTablesHandler interface using client-native client
type GetTablesHandlerImpl struct {
	Client client_native.HAProxyClient
}

// ReplaceTableHandlerImpl implementation of the ReplaceTableHandler interface using client-native client
type ReplaceTableHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// Handle executing the request and returning a response
func (h *CreateTableHandlerImpl) Handle(params table.CreateTableParams, principal interface{}) middleware.Responder {
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
		return table.NewCreateTableDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return table.NewCreateTableDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.CreateTable(params.PeerSection, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return table.NewCreateTableDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return table.NewCreateTableDefault(int(*e.Code)).WithPayload(e)
			}
			return table.NewCreateTableCreated().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return table.NewCreateTableAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return table.NewCreateTableAccepted().WithPayload(params.Data)
}

// Handle executing the request and returning a response
func (h *DeleteTableHandlerImpl) Handle(params table.DeleteTableParams, principal interface{}) middleware.Responder {
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
		return table.NewDeleteTableDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return table.NewDeleteTableDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.DeleteTable(params.Name, params.PeerSection, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return table.NewDeleteTableDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return table.NewDeleteTableDefault(int(*e.Code)).WithPayload(e)
			}
			return table.NewDeleteTableNoContent()
		}
		rID := h.ReloadAgent.Reload()
		return table.NewDeleteTableAccepted().WithReloadID(rID)
	}
	return table.NewDeleteTableAccepted()
}

// Handle executing the request and returning a response
func (h *GetTableHandlerImpl) Handle(params table.GetTableParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return table.NewGetTableDefault(int(*e.Code)).WithPayload(e)
	}

	_, ta, err := configuration.GetTable(params.Name, params.PeerSection, t)
	if err != nil {
		e := misc.HandleError(err)
		return table.NewGetTableDefault(int(*e.Code)).WithPayload(e)
	}
	return table.NewGetTableOK().WithPayload(ta)
}

// Handle executing the request and returning a response
func (h *GetTablesHandlerImpl) Handle(params table.GetTablesParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return table.NewGetTablesDefault(int(*e.Code)).WithPayload(e)
	}

	_, tables, err := configuration.GetTables(params.PeerSection, t)
	if err != nil {
		e := misc.HandleError(err)
		return table.NewGetTablesDefault(int(*e.Code)).WithPayload(e)
	}
	return table.NewGetTablesOK().WithPayload(tables)
}

// Handle executing the request and returning a response
func (h *ReplaceTableHandlerImpl) Handle(params table.ReplaceTableParams, principal interface{}) middleware.Responder {
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
		return table.NewReplaceTableDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return table.NewReplaceTableDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.EditTable(params.Name, params.PeerSection, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return table.NewReplaceTableDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return table.NewReplaceTableDefault(int(*e.Code)).WithPayload(e)
			}
			return table.NewReplaceTableOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return table.NewReplaceTableAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return table.NewReplaceTableAccepted().WithPayload(params.Data)
}
