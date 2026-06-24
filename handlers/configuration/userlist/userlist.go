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

package userlist

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	client_native "github.com/haproxytech/client-native/v6"
	"github.com/haproxytech/client-native/v6/models"

	"github.com/haproxytech/dataplaneapi/handlers/middleware"
	"github.com/haproxytech/dataplaneapi/handlers/respond"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/reload_agent"
)

// RegisterRouter registers all userlist and user routes onto r using spec-based request validation
// and a shared error handler.
func RegisterRouter(r chi.Router, client client_native.HAProxyClient, ra reload_agent.IReloadAgent) error {
	spec, err := GetSpec()
	if err != nil {
		return err
	}
	HandlerWithOptions(&HandlerImpl{Client: client, ReloadAgent: ra}, ChiServerOptions{
		BaseRouter:       r,
		Middlewares:      []MiddlewareFunc{middleware.NewValidator(spec)},
		ErrorHandlerFunc: middleware.ErrorHandler,
	})
	return nil
}

// HandlerImpl implements ServerInterface for HAProxy userlist and user configuration.
type HandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent reload_agent.IReloadAgent
}

func (h *HandlerImpl) GetUserlists(w http.ResponseWriter, r *http.Request, params GetUserlistsParams) {
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	var uls models.Userlists
	if params.FullSection {
		_, uls, err = cfg.GetStructuredUserLists(params.TransactionId)
	} else {
		_, uls, err = cfg.GetUserLists(params.TransactionId)
	}
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, uls)
}

func (h *HandlerImpl) CreateUserlist(w http.ResponseWriter, r *http.Request, params CreateUserlistParams) {
	if params.TransactionId != "" && params.ForceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data Userlist
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if params.FullSection {
		err = cfg.CreateStructuredUserList(&data, params.TransactionId, int64(params.Version))
	} else {
		err = cfg.CreateUserList(&data, params.TransactionId, int64(params.Version))
	}
	if err != nil {
		respond.Error(w, err)
		return
	}
	if params.TransactionId == "" {
		if params.ForceReload {
			if err = h.ReloadAgent.ForceReload(); err != nil {
				respond.Error(w, err)
				return
			}
			respond.JSON(w, http.StatusCreated, &data)
			return
		}
		respond.Accepted(w, h.ReloadAgent.Reload(), &data)
		return
	}
	respond.JSON(w, http.StatusAccepted, &data)
}

func (h *HandlerImpl) DeleteUserlist(w http.ResponseWriter, r *http.Request, name string, params DeleteUserlistParams) {
	if params.TransactionId != "" && params.ForceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.DeleteUserList(name, params.TransactionId, int64(params.Version)); err != nil {
		respond.Error(w, err)
		return
	}
	if params.TransactionId == "" {
		if params.ForceReload {
			if err = h.ReloadAgent.ForceReload(); err != nil {
				respond.Error(w, err)
				return
			}
			respond.NoContent(w)
			return
		}
		respond.Accepted(w, h.ReloadAgent.Reload(), nil)
		return
	}
	respond.Accepted(w, "", nil)
}

func (h *HandlerImpl) GetUserlist(w http.ResponseWriter, r *http.Request, name string, params GetUserlistParams) {
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	var ul *models.Userlist
	if params.FullSection {
		_, ul, err = cfg.GetStructuredUserList(name, params.TransactionId)
	} else {
		_, ul, err = cfg.GetUserList(name, params.TransactionId)
	}
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, ul)
}

func (h *HandlerImpl) GetUsers(w http.ResponseWriter, r *http.Request, params GetUsersParams) {
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	_, ul, err := cfg.GetUserList(params.Userlist, params.TransactionId)
	if ul == nil || err != nil {
		code := misc.ErrHTTPNotFound
		msg := fmt.Sprintf("userlist %s not found", params.Userlist)
		respond.JSON(w, http.StatusNotFound, &models.Error{Code: &code, Message: &msg})
		return
	}
	_, users, err := cfg.GetUsers(params.Userlist, params.TransactionId)
	if err != nil {
		e := misc.HandleContainerGetError(err)
		if *e.Code == misc.ErrHTTPOk {
			respond.JSON(w, http.StatusOK, models.Users{})
			return
		}
		respond.JSON(w, int(*e.Code), e)
		return
	}
	respond.JSON(w, http.StatusOK, users)
}

func (h *HandlerImpl) CreateUser(w http.ResponseWriter, r *http.Request, params CreateUserParams) {
	if params.TransactionId != "" && params.ForceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	_, ul, _ := cfg.GetUserList(params.Userlist, params.TransactionId)
	if ul == nil {
		respond.BadRequest(w, "userlist not found")
		return
	}
	var data User
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	if err = cfg.CreateUser(params.Userlist, &data, params.TransactionId, int64(params.Version)); err != nil {
		respond.Error(w, err)
		return
	}
	if params.TransactionId == "" {
		if params.ForceReload {
			if err = h.ReloadAgent.ForceReload(); err != nil {
				respond.Error(w, err)
				return
			}
			respond.JSON(w, http.StatusCreated, &data)
			return
		}
		respond.Accepted(w, h.ReloadAgent.Reload(), &data)
		return
	}
	respond.JSON(w, http.StatusAccepted, &data)
}

func (h *HandlerImpl) DeleteUser(w http.ResponseWriter, r *http.Request, username string, params DeleteUserParams) {
	if params.TransactionId != "" && params.ForceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	_, ul, err := cfg.GetUserList(params.Userlist, params.TransactionId)
	if err != nil {
		respond.Error(w, err)
		return
	}
	if ul == nil {
		code := misc.ErrHTTPNotFound
		msg := fmt.Sprintf("userlist %s not found", params.Userlist)
		respond.JSON(w, http.StatusNotFound, &models.Error{Code: &code, Message: &msg})
		return
	}
	_, u, err := cfg.GetUser(username, params.Userlist, params.TransactionId)
	if err != nil {
		respond.Error(w, err)
		return
	}
	if u == nil {
		code := misc.ErrHTTPNotFound
		msg := fmt.Sprintf("user %s not found in userlist %s", username, params.Userlist)
		respond.JSON(w, http.StatusNotFound, &models.Error{Code: &code, Message: &msg})
		return
	}
	if err = cfg.DeleteUser(username, params.Userlist, params.TransactionId, int64(params.Version)); err != nil {
		respond.Error(w, err)
		return
	}
	if params.TransactionId == "" {
		if params.ForceReload {
			if err = h.ReloadAgent.ForceReload(); err != nil {
				respond.Error(w, err)
				return
			}
			respond.NoContent(w)
			return
		}
		respond.Accepted(w, h.ReloadAgent.Reload(), nil)
		return
	}
	respond.Accepted(w, "", nil)
}

func (h *HandlerImpl) GetUser(w http.ResponseWriter, r *http.Request, username string, params GetUserParams) {
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	_, ul, err := cfg.GetUserList(params.Userlist, params.TransactionId)
	if err != nil {
		respond.Error(w, err)
		return
	}
	if ul == nil {
		code := misc.ErrHTTPNotFound
		msg := fmt.Sprintf("userlist %s not found", params.Userlist)
		respond.JSON(w, http.StatusNotFound, &models.Error{Code: &code, Message: &msg})
		return
	}
	_, u, err := cfg.GetUser(username, params.Userlist, params.TransactionId)
	if err != nil {
		respond.Error(w, err)
		return
	}
	if u == nil {
		code := misc.ErrHTTPNotFound
		msg := fmt.Sprintf("user %s not found in userlist %s", username, params.Userlist)
		respond.JSON(w, http.StatusNotFound, &models.Error{Code: &code, Message: &msg})
		return
	}
	respond.JSON(w, http.StatusOK, u)
}

func (h *HandlerImpl) ReplaceUser(w http.ResponseWriter, r *http.Request, username string, params ReplaceUserParams) {
	if params.TransactionId != "" && params.ForceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	_, ul, err := cfg.GetUserList(params.Userlist, params.TransactionId)
	if err != nil {
		respond.Error(w, err)
		return
	}
	if ul == nil {
		code := misc.ErrHTTPNotFound
		msg := fmt.Sprintf("userlist %s not found", params.Userlist)
		respond.JSON(w, http.StatusNotFound, &models.Error{Code: &code, Message: &msg})
		return
	}
	_, u, err := cfg.GetUser(username, params.Userlist, params.TransactionId)
	if err != nil {
		respond.Error(w, err)
		return
	}
	if u == nil {
		code := misc.ErrHTTPNotFound
		msg := fmt.Sprintf("user %s not found in userlist %s", username, params.Userlist)
		respond.JSON(w, http.StatusNotFound, &models.Error{Code: &code, Message: &msg})
		return
	}
	var data User
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	if err = cfg.EditUser(username, params.Userlist, &data, params.TransactionId, int64(params.Version)); err != nil {
		respond.Error(w, err)
		return
	}
	if params.TransactionId == "" {
		if params.ForceReload {
			if err = h.ReloadAgent.ForceReload(); err != nil {
				respond.Error(w, err)
				return
			}
			respond.JSON(w, http.StatusOK, &data)
			return
		}
		respond.Accepted(w, h.ReloadAgent.Reload(), &data)
		return
	}
	respond.JSON(w, http.StatusAccepted, &data)
}
