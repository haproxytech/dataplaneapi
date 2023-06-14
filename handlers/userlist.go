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
	"github.com/haproxytech/dataplaneapi/operations/userlist"
)

// CreateUserListHandlerImpl implementation of the CreateUserlistHandler interface using client-native client
type CreateUserListHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// DeleteUserListHandlerImpl implementation of the DeleteUserListHandler interface using client-native client
type DeleteUserListHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// GetUserListHandlerImpl implementation of the GetUserListHandler interface using client-native client
type GetUserListHandlerImpl struct {
	Client client_native.HAProxyClient
}

// GetUserListsHandlerImpl implementation of the GetUserListsHandler interface using client-native client
type GetUserListsHandlerImpl struct {
	Client client_native.HAProxyClient
}

// Handle executing the request and returning a response
func (h *CreateUserListHandlerImpl) Handle(params userlist.CreateUserlistParams, principal interface{}) middleware.Responder {
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
		return userlist.NewCreateUserlistDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return userlist.NewCreateUserlistDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.CreateUserList(params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return userlist.NewCreateUserlistDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return userlist.NewCreateUserlistDefault(int(*e.Code)).WithPayload(e)
			}
			return userlist.NewCreateUserlistCreated().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return userlist.NewCreateUserlistAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return userlist.NewCreateUserlistAccepted().WithPayload(params.Data)
}

// Handle executing the request and returning a response
func (h *DeleteUserListHandlerImpl) Handle(params userlist.DeleteUserlistParams, principal interface{}) middleware.Responder {
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
		return userlist.NewDeleteUserlistDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return userlist.NewDeleteUserlistDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.DeleteUserList(params.Name, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return userlist.NewDeleteUserlistDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return userlist.NewDeleteUserlistDefault(int(*e.Code)).WithPayload(e)
			}
			return userlist.NewDeleteUserlistNoContent()
		}
		rID := h.ReloadAgent.Reload()
		return userlist.NewDeleteUserlistAccepted().WithReloadID(rID)
	}
	return userlist.NewDeleteUserlistAccepted()
}

// Handle executing the request and returning a response
func (h *GetUserListHandlerImpl) Handle(params userlist.GetUserlistParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return userlist.NewGetUserlistDefault(int(*e.Code)).WithPayload(e)
	}

	v, u, err := configuration.GetUserList(params.Name, t)
	if err != nil {
		e := misc.HandleError(err)
		return userlist.NewGetUserlistDefault(int(*e.Code)).WithPayload(e)
	}
	return userlist.NewGetUserlistOK().WithPayload(&userlist.GetUserlistOKBody{Version: v, Data: u})
}

// Handle executing the request and returning a response
func (h *GetUserListsHandlerImpl) Handle(params userlist.GetUserlistsParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return userlist.NewGetUserlistsDefault(int(*e.Code)).WithPayload(e)
	}

	v, userlists, err := configuration.GetUserLists(t)
	if err != nil {
		e := misc.HandleError(err)
		return userlist.NewGetUserlistsDefault(int(*e.Code)).WithPayload(e)
	}
	return userlist.NewGetUserlistsOK().WithPayload(&userlist.GetUserlistsOKBody{Version: v, Data: userlists})
}
