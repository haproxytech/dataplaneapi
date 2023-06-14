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
	"github.com/haproxytech/dataplaneapi/operations/user"
)

// CreateUserHandlerImpl implementation of the CreateUserHandler interface using client-native client
type CreateUserHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// DeleteUserHandlerImpl implementation of the DeleteUserHandler interface using client-native client
type DeleteUserHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// GetUserHandlerImpl implementation of the GetUserHandler interface using client-native client
type GetUserHandlerImpl struct {
	Client client_native.HAProxyClient
}

// GetUsersHandlerImpl implementation of the GetUsersHandler interface using client-native client
type GetUsersHandlerImpl struct {
	Client client_native.HAProxyClient
}

// ReplaceUserHandlerImpl implementation of the ReplaceUserHandler interface using client-native client
type ReplaceUserHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// Handle executing the request and returning a response
func (h *CreateUserHandlerImpl) Handle(params user.CreateUserParams, principal interface{}) middleware.Responder {
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
		return user.NewCreateUserDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return user.NewCreateUserDefault(int(*e.Code)).WithPayload(e)
	}

	_, userlist, _ := configuration.GetUserList(params.Userlist, t)
	if userlist == nil {
		return user.NewCreateUserBadRequest()
	}

	err = configuration.CreateUser(params.Userlist, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return user.NewCreateUserDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return user.NewCreateUserDefault(int(*e.Code)).WithPayload(e)
			}
			return user.NewCreateUserCreated().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return user.NewCreateUserAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return user.NewCreateUserAccepted().WithPayload(params.Data)
}

// Handle executing the request and returning a response
func (h *DeleteUserHandlerImpl) Handle(params user.DeleteUserParams, principal interface{}) middleware.Responder {
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
		return user.NewDeleteUserDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return user.NewDeleteUserDefault(int(*e.Code)).WithPayload(e)
	}

	_, userlist, err := configuration.GetUserList(params.Userlist, t)
	if userlist == nil {
		return user.NewDeleteUserNotFound()
	}
	if err != nil {
		return user.NewDeleteUserNotFound()
	}

	_, u, err := configuration.GetUser(params.Username, params.Userlist, t)
	if u == nil {
		return user.NewDeleteUserNotFound()
	}
	if err != nil {
		return user.NewDeleteUserNotFound()
	}

	err = configuration.DeleteUser(params.Username, params.Userlist, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return user.NewDeleteUserDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return user.NewDeleteUserDefault(int(*e.Code)).WithPayload(e)
			}
			return user.NewDeleteUserNoContent()
		}
		rID := h.ReloadAgent.Reload()
		return user.NewDeleteUserAccepted().WithReloadID(rID)
	}
	return user.NewDeleteUserAccepted()
}

// Handle executing the request and returning a response
func (h *GetUserHandlerImpl) Handle(params user.GetUserParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return user.NewGetUserDefault(int(*e.Code)).WithPayload(e)
	}

	_, userlist, err := configuration.GetUserList(params.Userlist, t)
	if userlist == nil {
		return user.NewGetUserNotFound()
	}
	if err != nil {
		return user.NewGetUserNotFound()
	}

	v, u, err := configuration.GetUser(params.Username, params.Userlist, t)
	if u == nil {
		return user.NewGetUserNotFound()
	}
	if err != nil {
		e := misc.HandleError(err)
		return user.NewGetUserDefault(int(*e.Code)).WithPayload(e)
	}

	return user.NewGetUserOK().WithPayload(&user.GetUserOKBody{Version: v, Data: u})
}

// Handle executing the request and returning a response
func (h *GetUsersHandlerImpl) Handle(params user.GetUsersParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return user.NewGetUsersDefault(int(*e.Code)).WithPayload(e)
	}
	_, userlist, err := configuration.GetUserList(params.Userlist, t)
	if userlist == nil {
		return user.NewGetUserNotFound()
	}
	if err != nil {
		return user.NewGetUserNotFound()
	}
	v, users, err := configuration.GetUsers(params.Userlist, t)
	if err != nil {
		e := misc.HandleContainerGetError(err)
		if *e.Code == misc.ErrHTTPOk {
			return user.NewGetUsersOK().WithPayload(&user.GetUsersOKBody{Version: v, Data: models.Users{}})
		}
		return user.NewGetUsersDefault(int(*e.Code)).WithPayload(e)
	}
	return user.NewGetUsersOK().WithPayload(&user.GetUsersOKBody{Version: v, Data: users})
}

// Handle executing the request and returning a response
func (h *ReplaceUserHandlerImpl) Handle(params user.ReplaceUserParams, principal interface{}) middleware.Responder {
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
		return user.NewReplaceUserDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return user.NewReplaceUserDefault(int(*e.Code)).WithPayload(e)
	}
	_, userlist, err := configuration.GetUserList(params.Userlist, t)
	if userlist == nil {
		return user.NewReplaceUserNotFound()
	}
	if err != nil {
		return user.NewReplaceUserNotFound()
	}

	_, u, err := configuration.GetUser(params.Username, params.Userlist, t)
	if u == nil {
		return user.NewReplaceUserNotFound()
	}
	if err != nil {
		return user.NewReplaceUserNotFound()
	}

	err = configuration.EditUser(params.Username, params.Userlist, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return user.NewReplaceUserDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return user.NewReplaceUserDefault(int(*e.Code)).WithPayload(e)
			}
			return user.NewReplaceUserOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return user.NewReplaceUserAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return user.NewReplaceUserAccepted().WithPayload(params.Data)
}
