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
	client_native "github.com/haproxytech/client-native/v5"
	"github.com/haproxytech/client-native/v5/models"

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

// GetAclsHandlerImpl implementation of the GetAclsHandler interface using client-native client
type GetAclsHandlerImpl struct {
	Client client_native.HAProxyClient
}

// ReplaceACLHandlerImpl implementation of the ReplaceACLHandler interface using client-native client
type ReplaceACLHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// Handle executing the request and returning a response
func (h *CreateACLHandlerImpl) Handle(params acl.CreateACLParams, principal interface{}) middleware.Responder {
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
		return acl.NewCreateACLDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return acl.NewCreateACLDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.CreateACL(params.ParentType, params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return acl.NewCreateACLDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return acl.NewCreateACLDefault(int(*e.Code)).WithPayload(e)
			}
			return acl.NewCreateACLCreated().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return acl.NewCreateACLAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return acl.NewCreateACLAccepted().WithPayload(params.Data)
}

// Handle executing the request and returning a response
func (h *DeleteACLHandlerImpl) Handle(params acl.DeleteACLParams, principal interface{}) middleware.Responder {
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
		return acl.NewDeleteACLDefault(int(*e.Code)).WithPayload(e)
	}
	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return acl.NewDeleteACLDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.DeleteACL(params.Index, params.ParentType, params.ParentName, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return acl.NewDeleteACLDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return acl.NewDeleteACLDefault(int(*e.Code)).WithPayload(e)
			}
			return acl.NewDeleteACLNoContent()
		}
		rID := h.ReloadAgent.Reload()
		return acl.NewDeleteACLAccepted().WithReloadID(rID)
	}
	return acl.NewDeleteACLAccepted()
}

// Handle executing the request and returning a response
func (h *GetACLHandlerImpl) Handle(params acl.GetACLParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return acl.NewGetACLDefault(int(*e.Code)).WithPayload(e)
	}

	v, rule, err := configuration.GetACL(params.Index, params.ParentType, params.ParentName, t)
	if err != nil {
		e := misc.HandleError(err)
		return acl.NewGetACLDefault(int(*e.Code)).WithPayload(e)
	}

	return acl.NewGetACLOK().WithPayload(&acl.GetACLOKBody{Version: v, Data: rule})
}

// Handle executing the request and returning a response
func (h *GetAclsHandlerImpl) Handle(params acl.GetAclsParams, principal interface{}) middleware.Responder {
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
		return acl.NewGetACLDefault(int(*e.Code)).WithPayload(e)
	}

	v, rules, err := configuration.GetACLs(params.ParentType, params.ParentName, t, aclName...)
	if err != nil {
		e := misc.HandleContainerGetError(err)
		if *e.Code == misc.ErrHTTPOk {
			return acl.NewGetAclsOK().WithPayload(&acl.GetAclsOKBody{Version: v, Data: models.Acls{}})
		}
		return acl.NewGetAclsDefault(int(*e.Code)).WithPayload(e)
	}
	return acl.NewGetAclsOK().WithPayload(&acl.GetAclsOKBody{Version: v, Data: rules})
}

// Handle executing the request and returning a response
func (h *ReplaceACLHandlerImpl) Handle(params acl.ReplaceACLParams, principal interface{}) middleware.Responder {
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
		return acl.NewReplaceACLDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return acl.NewReplaceACLDefault(int(*e.Code)).WithPayload(e)
	}
	err = configuration.EditACL(params.Index, params.ParentType, params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return acl.NewReplaceACLDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return acl.NewReplaceACLDefault(int(*e.Code)).WithPayload(e)
			}
			return acl.NewReplaceACLOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return acl.NewReplaceACLAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return acl.NewReplaceACLAccepted().WithPayload(params.Data)
}
