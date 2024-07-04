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
	"github.com/haproxytech/dataplaneapi/operations/log_target"
)

// CreateLogTargetHandlerImpl implementation of the CreateLogTargetHandler interface using client-native client
type CreateLogTargetHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// DeleteLogTargetHandlerImpl implementation of the DeleteLogTargetHandler interface using client-native client
type DeleteLogTargetHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// GetLogTargetHandlerImpl implementation of the GetLogTargetHandler interface using client-native client
type GetLogTargetHandlerImpl struct {
	Client client_native.HAProxyClient
}

// GetAllLogTargetHandlerImpl implementation of the GetLogTargetsHandler interface using client-native client
type GetAllLogTargetHandlerImpl struct {
	Client client_native.HAProxyClient
}

// ReplaceLogTargetHandlerImpl implementation of the ReplaceLogTargetHandler interface using client-native client
type ReplaceLogTargetHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// ReplaceAllLogTargetHandlerImpl implementation of the ReplaceLogTargetsHandler interface using client-native client
type ReplaceAllLogTargetHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// Handle executing the request and returning a response
func (h *CreateLogTargetHandlerImpl) Handle(parentType cnconstants.CnParentType, params log_target.CreateLogTargetBackendParams, principal interface{}) middleware.Responder {
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
		return log_target.NewCreateLogTargetBackendDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return log_target.NewCreateLogTargetBackendDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.CreateLogTarget(params.Index, string(parentType), params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return log_target.NewCreateLogTargetBackendDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return log_target.NewCreateLogTargetBackendDefault(int(*e.Code)).WithPayload(e)
			}
			return log_target.NewCreateLogTargetBackendCreated().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return log_target.NewCreateLogTargetBackendAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return log_target.NewCreateLogTargetBackendAccepted().WithPayload(params.Data)
}

// Handle executing the request and returning a response
func (h *DeleteLogTargetHandlerImpl) Handle(parentType cnconstants.CnParentType, params log_target.DeleteLogTargetBackendParams, principal interface{}) middleware.Responder {
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
		return log_target.NewDeleteLogTargetBackendDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return log_target.NewDeleteLogTargetBackendDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.DeleteLogTarget(params.Index, string(parentType), params.ParentName, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return log_target.NewDeleteLogTargetBackendDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return log_target.NewDeleteLogTargetBackendDefault(int(*e.Code)).WithPayload(e)
			}
			return log_target.NewDeleteLogTargetBackendNoContent()
		}
		rID := h.ReloadAgent.Reload()
		return log_target.NewDeleteLogTargetBackendAccepted().WithReloadID(rID)
	}
	return log_target.NewDeleteLogTargetBackendAccepted()
}

// Handle executing the request and returning a response
func (h *GetLogTargetHandlerImpl) Handle(parentType cnconstants.CnParentType, params log_target.GetLogTargetBackendParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return log_target.NewGetLogTargetBackendDefault(int(*e.Code)).WithPayload(e)
	}

	_, logTarget, err := configuration.GetLogTarget(params.Index, string(parentType), params.ParentName, t)
	if err != nil {
		e := misc.HandleError(err)
		return log_target.NewGetLogTargetBackendDefault(int(*e.Code)).WithPayload(e)
	}
	return log_target.NewGetLogTargetBackendOK().WithPayload(logTarget)
}

// Handle executing the request and returning a response
func (h *GetAllLogTargetHandlerImpl) Handle(parentType cnconstants.CnParentType, params log_target.GetAllLogTargetBackendParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return log_target.NewGetAllLogTargetBackendDefault(int(*e.Code)).WithPayload(e)
	}

	_, logTargets, err := configuration.GetLogTargets(string(parentType), params.ParentName, t)
	if err != nil {
		e := misc.HandleContainerGetError(err)
		if *e.Code == misc.ErrHTTPOk {
			return log_target.NewGetAllLogTargetBackendOK().WithPayload(models.LogTargets{})
		}
		return log_target.NewGetAllLogTargetBackendDefault(int(*e.Code)).WithPayload(e)
	}
	return log_target.NewGetAllLogTargetBackendOK().WithPayload(logTargets)
}

// Handle executing the request and returning a response
func (h *ReplaceLogTargetHandlerImpl) Handle(parentType cnconstants.CnParentType, params log_target.ReplaceLogTargetBackendParams, principal interface{}) middleware.Responder {
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
		return log_target.NewReplaceLogTargetBackendDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return log_target.NewReplaceLogTargetBackendDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.EditLogTarget(params.Index, string(parentType), params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return log_target.NewReplaceLogTargetBackendDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return log_target.NewReplaceLogTargetBackendDefault(int(*e.Code)).WithPayload(e)
			}
			return log_target.NewReplaceLogTargetBackendOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return log_target.NewReplaceLogTargetBackendAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return log_target.NewReplaceLogTargetBackendAccepted().WithPayload(params.Data)
}

func logTargetParentTypeRequiresParentName(parentType string) bool {
	return (parentType == "frontend" || parentType == "backend" || parentType == "peers" || parentType == "log_forward")
}

// Handle executing the request and returning a response
func (h *ReplaceAllLogTargetHandlerImpl) Handle(parentType cnconstants.CnParentType, params log_target.ReplaceAllLogTargetBackendParams, principal interface{}) middleware.Responder {
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
		return log_target.NewReplaceAllLogTargetBackendDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return log_target.NewReplaceAllLogTargetBackendDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.ReplaceLogTargets(string(parentType), params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return log_target.NewReplaceAllLogTargetBackendDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return log_target.NewReplaceAllLogTargetBackendDefault(int(*e.Code)).WithPayload(e)
			}
			return log_target.NewReplaceAllLogTargetBackendOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return log_target.NewReplaceAllLogTargetBackendAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return log_target.NewReplaceAllLogTargetBackendAccepted().WithPayload(params.Data)
}
