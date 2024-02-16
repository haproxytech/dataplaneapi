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
	"github.com/haproxytech/dataplaneapi/operations/group"
)

// CreateGroupHandlerImpl implementation of the CreateGroupHandler interface using client-native client
type CreateGroupHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// DeleteGroupHandlerImpl implementation of the DeleteGroupHandler interface using client-native client
type DeleteGroupHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// GetGroupHandlerImpl implementation of the GetGroupHandler interface using client-native client
type GetGroupHandlerImpl struct {
	Client client_native.HAProxyClient
}

// GetGroupsHandlerImpl implementation of the GetGroupsHandler interface using client-native client
type GetGroupsHandlerImpl struct {
	Client client_native.HAProxyClient
}

// ReplaceGroupHandlerImpl implementation of the ReplaceGroupHandler interface using client-native client
type ReplaceGroupHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// Handle executing the request and returning a response
func (h *CreateGroupHandlerImpl) Handle(params group.CreateGroupParams, principal interface{}) middleware.Responder {
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
		return group.NewCreateGroupDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return group.NewCreateGroupDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.CreateGroup(params.Userlist, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return group.NewCreateGroupDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return group.NewCreateGroupDefault(int(*e.Code)).WithPayload(e)
			}
			return group.NewCreateGroupCreated().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return group.NewCreateGroupAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return group.NewCreateGroupAccepted().WithPayload(params.Data)
}

// Handle executing the request and returning a response
func (h *DeleteGroupHandlerImpl) Handle(params group.DeleteGroupParams, principal interface{}) middleware.Responder {
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
		return group.NewDeleteGroupDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return group.NewDeleteGroupDefault(int(*e.Code)).WithPayload(e)
	}

	_, userlist, err := configuration.GetUserList(params.Userlist, t)
	if err != nil {
		e := misc.HandleError(err)
		return group.NewDeleteGroupDefault(int(*e.Code)).WithPayload(e)
	}
	if userlist == nil {
		return group.NewDeleteGroupNotFound()
	}

	err = configuration.DeleteGroup(params.Name, params.Userlist, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return group.NewDeleteGroupDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return group.NewDeleteGroupDefault(int(*e.Code)).WithPayload(e)
			}
			return group.NewDeleteGroupNoContent()
		}
		rID := h.ReloadAgent.Reload()
		return group.NewDeleteGroupAccepted().WithReloadID(rID)
	}
	return group.NewDeleteGroupAccepted()
}

// Handle executing the request and returning a response
func (h *GetGroupHandlerImpl) Handle(params group.GetGroupParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return group.NewGetGroupDefault(int(*e.Code)).WithPayload(e)
	}

	_, g, err := configuration.GetGroup(params.Name, params.Userlist, t)
	if err != nil {
		e := misc.HandleError(err)
		return group.NewGetGroupDefault(int(*e.Code)).WithPayload(e)
	}

	if g == nil {
		return group.NewGetGroupNotFound()
	}

	return group.NewGetGroupOK().WithPayload(g)
}

// Handle executing the request and returning a response
func (h *GetGroupsHandlerImpl) Handle(params group.GetGroupsParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return group.NewGetGroupsDefault(int(*e.Code)).WithPayload(e)
	}

	_, userlist, err := configuration.GetUserList(params.Userlist, t)
	if userlist == nil {
		return group.NewGetGroupNotFound()
	}
	if err != nil {
		return group.NewGetGroupNotFound()
	}

	_, groups, err := configuration.GetGroups(params.Userlist, t)
	if err != nil {
		e := misc.HandleContainerGetError(err)
		if *e.Code == misc.ErrHTTPOk {
			return group.NewGetGroupsOK().WithPayload(models.Groups{})
		}
		return group.NewGetGroupsDefault(int(*e.Code)).WithPayload(e)
	}
	return group.NewGetGroupsOK().WithPayload(groups)
}

// Handle executing the request and returning a response
func (h *ReplaceGroupHandlerImpl) Handle(params group.ReplaceGroupParams, principal interface{}) middleware.Responder {
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
		return group.NewReplaceGroupDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return group.NewReplaceGroupDefault(int(*e.Code)).WithPayload(e)
	}

	_, userlist, err := configuration.GetUserList(params.Userlist, t)
	if userlist == nil {
		return group.NewReplaceGroupNotFound()
	}
	if err != nil {
		return group.NewReplaceGroupNotFound()
	}

	_, g, err := configuration.GetGroup(params.Name, params.Userlist, t)
	if err != nil {
		e := misc.HandleError(err)
		return group.NewGetGroupDefault(int(*e.Code)).WithPayload(e)
	}
	if g == nil {
		return group.NewReplaceGroupNotFound()
	}

	err = configuration.EditGroup(params.Name, params.Userlist, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return group.NewReplaceGroupDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return group.NewReplaceGroupDefault(int(*e.Code)).WithPayload(e)
			}
			return group.NewReplaceGroupOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return group.NewReplaceGroupAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return group.NewReplaceGroupAccepted().WithPayload(params.Data)
}
