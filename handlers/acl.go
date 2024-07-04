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
	"github.com/haproxytech/dataplaneapi/operations/acl"
)

// CreateACLHandlerImpl implementation of the CreateACLHandler interface using client-native client
type CreateACLHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// DeleteACLHandlerImpl implementation of the DeleteACLHandler interface using client-native client
type DeleteACLHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// GetACLHandlerImpl implementation of the GetACLHandler interface using client-native client
type GetACLHandlerImpl struct {
	Client client_native.HAProxyClient
}

// GetAllACLHandlerImpl implementation of the GetAclsHandler interface using client-native client
type GetAllACLHandlerImpl struct {
	Client client_native.HAProxyClient
}

// ReplaceACLHandlerImpl implementation of the ReplaceACLHandler interface using client-native client
type ReplaceACLHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// ReplaceAllACLHandlerImpl implementation of the ReplaceAclsHandler interface using client-native client
type ReplaceAllACLHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// Handle executing the request and returning a response
func (h *CreateACLHandlerImpl) Handle(parentType cnconstants.CnParentType, params acl.CreateACLBackendParams, principal interface{}) middleware.Responder {
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
		return acl.NewCreateACLBackendDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return acl.NewCreateACLBackendDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.CreateACL(params.Index, string(parentType), params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return acl.NewCreateACLBackendDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return acl.NewCreateACLBackendDefault(int(*e.Code)).WithPayload(e)
			}
			return acl.NewCreateACLBackendCreated().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return acl.NewCreateACLBackendAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return acl.NewCreateACLBackendAccepted().WithPayload(params.Data)
}

// Handle executing the request and returning a response
func (h *DeleteACLHandlerImpl) Handle(parentType cnconstants.CnParentType, params acl.DeleteACLBackendParams, principal interface{}) middleware.Responder {
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
		return acl.NewDeleteACLBackendDefault(int(*e.Code)).WithPayload(e)
	}
	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return acl.NewDeleteACLBackendDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.DeleteACL(params.Index, string(parentType), params.ParentName, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return acl.NewDeleteACLBackendDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return acl.NewDeleteACLBackendDefault(int(*e.Code)).WithPayload(e)
			}
			return acl.NewDeleteACLBackendNoContent()
		}
		rID := h.ReloadAgent.Reload()
		return acl.NewDeleteACLBackendAccepted().WithReloadID(rID)
	}
	return acl.NewDeleteACLBackendAccepted()
}

// Handle executing the request and returning a response
func (h *GetACLHandlerImpl) Handle(parentType cnconstants.CnParentType, params acl.GetACLBackendParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return acl.NewGetACLBackendDefault(int(*e.Code)).WithPayload(e)
	}

	_, rule, err := configuration.GetACL(params.Index, string(parentType), params.ParentName, t)
	if err != nil {
		e := misc.HandleError(err)
		return acl.NewGetACLBackendDefault(int(*e.Code)).WithPayload(e)
	}

	return acl.NewGetACLBackendOK().WithPayload(rule)
}

// Handle executing the request and returning a response
func (h *GetAllACLHandlerImpl) Handle(parentType cnconstants.CnParentType, params acl.GetAllACLBackendParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	var aclName []string
	if params.ACLName != nil {
		aclName = []string{*params.ACLName}
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return acl.NewGetACLBackendDefault(int(*e.Code)).WithPayload(e)
	}

	_, rules, err := configuration.GetACLs(string(parentType), params.ParentName, t, aclName...)
	if err != nil {
		e := misc.HandleContainerGetError(err)
		if *e.Code == misc.ErrHTTPOk {
			return acl.NewGetAllACLBackendOK().WithPayload(rules)
		}
		return acl.NewGetAllACLBackendDefault(int(*e.Code)).WithPayload(e)
	}
	return acl.NewGetAllACLBackendOK().WithPayload(rules)
}

// Handle executing the request and returning a response
func (h *ReplaceACLHandlerImpl) Handle(parentType cnconstants.CnParentType, params acl.ReplaceACLBackendParams, principal interface{}) middleware.Responder {
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
		return acl.NewReplaceACLBackendDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return acl.NewReplaceACLBackendDefault(int(*e.Code)).WithPayload(e)
	}
	err = configuration.EditACL(params.Index, string(parentType), params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return acl.NewReplaceACLBackendDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return acl.NewReplaceACLBackendDefault(int(*e.Code)).WithPayload(e)
			}
			return acl.NewReplaceACLBackendOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return acl.NewReplaceACLBackendAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return acl.NewReplaceACLBackendAccepted().WithPayload(params.Data)
}

// Handle executing the request and returning a response
func (h *ReplaceAllACLHandlerImpl) Handle(parentType cnconstants.CnParentType, params acl.ReplaceAllACLBackendParams, principal interface{}) middleware.Responder {
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
		return acl.NewReplaceAllACLBackendDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return acl.NewReplaceAllACLBackendDefault(int(*e.Code)).WithPayload(e)
	}
	err = configuration.ReplaceAcls(string(parentType), params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return acl.NewReplaceAllACLBackendDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return acl.NewReplaceAllACLBackendDefault(int(*e.Code)).WithPayload(e)
			}
			return acl.NewReplaceAllACLBackendOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return acl.NewReplaceAllACLBackendAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return acl.NewReplaceAllACLBackendAccepted().WithPayload(params.Data)
}
