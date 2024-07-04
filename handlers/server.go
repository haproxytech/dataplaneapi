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
	"errors"

	"github.com/haproxytech/client-native/v6/runtime"

	"github.com/go-openapi/runtime/middleware"
	client_native "github.com/haproxytech/client-native/v6"
	cnconstants "github.com/haproxytech/client-native/v6/configuration/parents"
	"github.com/haproxytech/client-native/v6/models"

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

// GetAllServerHandlerImpl implementation of the GetServersHandler interface using client-native client
type GetAllServerHandlerImpl struct {
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
		return "", "", errors.New("parentType empty")
	}
	if parentName == nil || *parentName == "" {
		return "", "", errors.New("parentName empty")
	}
	return *parentType, *parentName, nil
}

// Handle executing the request and returning a response
func (h *CreateServerHandlerImpl) Handle(parentType cnconstants.CnParentType, params server.CreateServerBackendParams, principal interface{}) middleware.Responder {
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
		return server.NewCreateServerBackendDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return server.NewCreateServerBackendDefault(int(*e.Code)).WithPayload(e)
	}

	pType, pName, err := serverTypeParams(nil, misc.StringP(string(parentType)), &params.ParentName)
	if err != nil {
		e := misc.HandleError(err)
		return server.NewCreateServerBackendDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.CreateServer(pType, pName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return server.NewCreateServerBackendDefault(int(*e.Code)).WithPayload(e)
	}

	// Try to create the new server dynamically. This is only possible if parentType is `backend` and no `default_server`
	// was defined in the current backend or in the `defaults` section.
	useRuntime := false
	ras := &models.RuntimeAddServer{}
	var runtimeClient runtime.Runtime
	if pType == "backend" {
		_, defaults, errRuntime := configuration.GetDefaultsConfiguration(t)
		if errRuntime != nil {
			e := misc.HandleError(errRuntime)
			return server.NewCreateServerBackendDefault(int(*e.Code)).WithPayload(e)
		}
		_, backend, errRuntime := configuration.GetBackend(pName, t)
		if errRuntime != nil {
			e := misc.HandleError(errRuntime)
			return server.NewCreateServerBackendDefault(int(*e.Code)).WithPayload(e)
		}
		runtimeClient, errRuntime = h.Client.Runtime()
		if errRuntime == nil && defaults.DefaultServer == nil && backend.DefaultServer == nil {
			// Also make sure the server attributes are supported by the runtime API.
			errRuntime = misc.ConvertStruct(params.Data, ras)
			useRuntime = errRuntime == nil
		}
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err = h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return server.NewCreateServerBackendDefault(int(*e.Code)).WithPayload(e)
			}
			return server.NewCreateServerBackendCreated().WithPayload(params.Data)
		}
		if useRuntime {
			err = runtimeClient.AddServer(pName, params.Data.Name, SerializeRuntimeAddServer(ras))
			if err == nil {
				// No need to reload.
				log.Debugf("backend %s: server %s added though runtime", pName, params.Data.Name)
				return server.NewCreateServerBackendCreated().WithPayload(params.Data)
			}
			log.Warning("failed to add server through runtime:", err)
		}
		rID := h.ReloadAgent.Reload()
		return server.NewCreateServerBackendAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return server.NewCreateServerBackendAccepted().WithPayload(params.Data)
}

// Handle executing the request and returning a response
func (h *DeleteServerHandlerImpl) Handle(parentType cnconstants.CnParentType, params server.DeleteServerBackendParams, principal interface{}) middleware.Responder {
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
		return server.NewDeleteServerBackendDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return server.NewDeleteServerBackendDefault(int(*e.Code)).WithPayload(e)
	}

	pType, pName, err := serverTypeParams(nil, misc.StringP(string(parentType)), &params.ParentName)
	if err != nil {
		e := misc.HandleError(err)
		return server.NewDeleteServerBackendDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.DeleteServer(params.Name, pType, pName, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return server.NewDeleteServerBackendDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return server.NewDeleteServerBackendDefault(int(*e.Code)).WithPayload(e)
			}
			return server.NewDeleteServerBackendNoContent()
		}
		rID := h.ReloadAgent.Reload()
		return server.NewDeleteServerBackendAccepted().WithReloadID(rID)
	}
	return server.NewDeleteServerBackendAccepted()
}

// Handle executing the request and returning a response
func (h *GetServerHandlerImpl) Handle(parentType cnconstants.CnParentType, params server.GetServerBackendParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return server.NewGetServerBackendDefault(int(*e.Code)).WithPayload(e)
	}

	pType, pName, err := serverTypeParams(nil, misc.StringP(string(parentType)), &params.ParentName)
	if err != nil {
		e := misc.HandleError(err)
		return server.NewGetServerBackendDefault(int(*e.Code)).WithPayload(e)
	}

	_, srv, err := configuration.GetServer(params.Name, pType, pName, t)
	if err != nil {
		e := misc.HandleError(err)
		return server.NewGetServerBackendDefault(int(*e.Code)).WithPayload(e)
	}
	return server.NewGetServerBackendOK().WithPayload(srv)
}

// Handle executing the request and returning a response
func (h *GetAllServerHandlerImpl) Handle(parentType cnconstants.CnParentType, params server.GetAllServerBackendParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return server.NewGetAllServerBackendDefault(int(*e.Code)).WithPayload(e)
	}

	pType, pName, err := serverTypeParams(nil, misc.StringP(string(parentType)), &params.ParentName)
	if err != nil {
		e := misc.HandleError(err)
		return server.NewGetAllServerBackendDefault(int(*e.Code)).WithPayload(e)
	}

	_, srvs, err := configuration.GetServers(pType, pName, t)
	if err != nil {
		e := misc.HandleContainerGetError(err)
		if *e.Code == misc.ErrHTTPOk {
			return server.NewGetAllServerBackendOK().WithPayload(models.Servers{})
		}
		return server.NewGetAllServerBackendDefault(int(*e.Code)).WithPayload(e)
	}
	return server.NewGetAllServerBackendOK().WithPayload(srvs)
}

// Handle executing the request and returning a response
func (h *ReplaceServerHandlerImpl) Handle(parentType cnconstants.CnParentType, params server.ReplaceServerBackendParams, principal interface{}) middleware.Responder {
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
		return server.NewReplaceServerBackendDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return server.NewReplaceServerBackendDefault(int(*e.Code)).WithPayload(e)
	}

	pType, pName, err := serverTypeParams(nil, misc.StringP(string(parentType)), &params.ParentName)
	if err != nil {
		e := misc.HandleError(err)
		return server.NewReplaceServerBackendDefault(int(*e.Code)).WithPayload(e)
	}

	_, ondisk, err := configuration.GetServer(params.Name, pType, pName, t)
	if err != nil {
		e := misc.HandleError(err)
		return server.NewReplaceServerBackendDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.EditServer(params.Name, pType, pName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return server.NewReplaceServerBackendDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		reload := changeThroughRuntimeAPI(*params.Data, *ondisk, pName, h.Client)
		if reload {
			if *params.ForceReload {
				err := h.ReloadAgent.ForceReload()
				if err != nil {
					e := misc.HandleError(err)
					return server.NewReplaceServerBackendDefault(int(*e.Code)).WithPayload(e)
				}
				return server.NewReplaceServerBackendOK().WithPayload(params.Data)
			}
			rID := h.ReloadAgent.Reload()
			return server.NewReplaceServerBackendAccepted().WithReloadID(rID).WithPayload(params.Data)
		}
		return server.NewReplaceServerBackendOK().WithPayload(params.Data)
	}
	return server.NewReplaceServerBackendAccepted().WithPayload(params.Data)
}
