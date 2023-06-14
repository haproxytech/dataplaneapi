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
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	client_native "github.com/haproxytech/client-native/v5"
	"github.com/haproxytech/client-native/v5/models"

	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/log"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/operations/server"
)

// CreateServerHandlerImpl implementation of the CreateServerHandler interface using client-native client
type CreateServerHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// DeleteServerHandlerImpl implementation of the DeleteServerHandler interface using client-native client
type DeleteServerHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// GetServerHandlerImpl implementation of the GetServerHandler interface using client-native client
type GetServerHandlerImpl struct {
	Client client_native.HAProxyClient
}

// GetServersHandlerImpl implementation of the GetServersHandler interface using client-native client
type GetServersHandlerImpl struct {
	Client client_native.HAProxyClient
}

// ReplaceServerHandlerImpl implementation of the ReplaceServerHandler interface using client-native client
type ReplaceServerHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

func serverTypeParams(backend *string, parentType *string, parentName *string) (pType string, pName string, err error) {
	if backend != nil && *backend != "" {
		return "backend", *backend, nil
	}
	if parentType == nil || *parentType == "" {
		return "", "", fmt.Errorf("parentType empty")
	}
	if parentName == nil || *parentName == "" {
		return "", "", fmt.Errorf("parentName empty")
	}
	return *parentType, *parentName, nil
}

// Handle executing the request and returning a response
func (h *CreateServerHandlerImpl) Handle(params server.CreateServerParams, principal interface{}) middleware.Responder {
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
		return server.NewCreateServerDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return server.NewCreateServerDefault(int(*e.Code)).WithPayload(e)
	}

	pType, pName, err := serverTypeParams(params.Backend, params.ParentType, params.ParentName)
	if err != nil {
		e := misc.HandleError(err)
		return server.NewCreateServerDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.CreateServer(pType, pName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return server.NewCreateServerDefault(int(*e.Code)).WithPayload(e)
	}

	// Try to create the new server dynamically. This is only possible if no `default_server`
	// was defined in the current backend or in the `defaults` section.
	useRuntime := false
	var ras *models.RuntimeAddServer
	_, defaults, err := configuration.GetDefaultsConfiguration(t)
	if err != nil {
		e := misc.HandleError(err)
		return server.NewCreateServerDefault(int(*e.Code)).WithPayload(e)
	}
	_, backend, err := configuration.GetBackend(pName, t)
	if err != nil {
		e := misc.HandleError(err)
		return server.NewCreateServerDefault(int(*e.Code)).WithPayload(e)
	}
	runtime, err := h.Client.Runtime()
	if err == nil && defaults.DefaultServer == nil && backend.DefaultServer == nil {
		// Also make sure the server attributes are supported by the runtime API.
		err = misc.ConvertStruct(params.Data, ras)
		useRuntime = err == nil
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err = h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return server.NewCreateServerDefault(int(*e.Code)).WithPayload(e)
			}
			return server.NewCreateServerCreated().WithPayload(params.Data)
		}
		if useRuntime {
			err = runtime.AddServer(pName, params.Data.Name, SerializeRuntimeAddServer(ras))
			if err == nil {
				// No need to reload.
				return server.NewCreateServerCreated().WithPayload(params.Data)
			}
			log.Warning("failed to add server through runtime:", err)
		}
		rID := h.ReloadAgent.Reload()
		return server.NewCreateServerAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return server.NewCreateServerAccepted().WithPayload(params.Data)
}

// Handle executing the request and returning a response
func (h *DeleteServerHandlerImpl) Handle(params server.DeleteServerParams, principal interface{}) middleware.Responder {
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
		return server.NewDeleteServerDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return server.NewDeleteServerDefault(int(*e.Code)).WithPayload(e)
	}

	pType, pName, err := serverTypeParams(params.Backend, params.ParentType, params.ParentName)
	if err != nil {
		e := misc.HandleError(err)
		return server.NewDeleteServerDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.DeleteServer(params.Name, pType, pName, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return server.NewDeleteServerDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return server.NewDeleteServerDefault(int(*e.Code)).WithPayload(e)
			}
			return server.NewDeleteServerNoContent()
		}
		rID := h.ReloadAgent.Reload()
		return server.NewDeleteServerAccepted().WithReloadID(rID)
	}
	return server.NewDeleteServerAccepted()
}

// Handle executing the request and returning a response
func (h *GetServerHandlerImpl) Handle(params server.GetServerParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return server.NewGetRuntimeServerDefault(int(*e.Code)).WithPayload(e)
	}

	pType, pName, err := serverTypeParams(params.Backend, params.ParentType, params.ParentName)
	if err != nil {
		e := misc.HandleError(err)
		return server.NewGetRuntimeServerDefault(int(*e.Code)).WithPayload(e)
	}

	v, srv, err := configuration.GetServer(params.Name, pType, pName, t)
	if err != nil {
		e := misc.HandleError(err)
		return server.NewGetServerDefault(int(*e.Code)).WithPayload(e)
	}
	return server.NewGetServerOK().WithPayload(&server.GetServerOKBody{Version: v, Data: srv})
}

// Handle executing the request and returning a response
func (h *GetServersHandlerImpl) Handle(params server.GetServersParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return server.NewGetRuntimeServersDefault(int(*e.Code)).WithPayload(e)
	}

	pType, pName, err := serverTypeParams(params.Backend, params.ParentType, params.ParentName)
	if err != nil {
		e := misc.HandleError(err)
		return server.NewGetRuntimeServersDefault(int(*e.Code)).WithPayload(e)
	}

	v, srvs, err := configuration.GetServers(pType, pName, t)
	if err != nil {
		e := misc.HandleContainerGetError(err)
		if *e.Code == misc.ErrHTTPOk {
			return server.NewGetServersOK().WithPayload(&server.GetServersOKBody{Version: v, Data: models.Servers{}})
		}
		return server.NewGetServersDefault(int(*e.Code)).WithPayload(e)
	}
	return server.NewGetServersOK().WithPayload(&server.GetServersOKBody{Version: v, Data: srvs})
}

// Handle executing the request and returning a response
func (h *ReplaceServerHandlerImpl) Handle(params server.ReplaceServerParams, principal interface{}) middleware.Responder {
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
		return server.NewReplaceServerDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return server.NewReplaceServerDefault(int(*e.Code)).WithPayload(e)
	}

	pType, pName, err := serverTypeParams(params.Backend, params.ParentType, params.ParentName)
	if err != nil {
		e := misc.HandleError(err)
		return server.NewReplaceServerDefault(int(*e.Code)).WithPayload(e)
	}

	_, ondisk, err := configuration.GetServer(params.Name, pType, pName, t)
	if err != nil {
		e := misc.HandleError(err)
		return server.NewReplaceServerDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.EditServer(params.Name, pType, pName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return server.NewReplaceServerDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		reload := changeThroughRuntimeAPI(*params.Data, *ondisk, pType, "", h.Client)
		if reload {
			if *params.ForceReload {
				err := h.ReloadAgent.ForceReload()
				if err != nil {
					e := misc.HandleError(err)
					return server.NewReplaceServerDefault(int(*e.Code)).WithPayload(e)
				}
				return server.NewReplaceServerOK().WithPayload(params.Data)
			}
			rID := h.ReloadAgent.Reload()
			return server.NewReplaceServerAccepted().WithReloadID(rID).WithPayload(params.Data)
		}
		return server.NewReplaceServerOK().WithPayload(params.Data)
	}
	return server.NewReplaceServerAccepted().WithPayload(params.Data)
}
