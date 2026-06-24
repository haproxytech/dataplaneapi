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

package groups

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	client_native "github.com/haproxytech/client-native/v6"
	"github.com/haproxytech/client-native/v6/models"

	"github.com/haproxytech/dataplaneapi/handlers/middleware"
	"github.com/haproxytech/dataplaneapi/handlers/respond"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/reload_agent"
)

func notFound(w http.ResponseWriter) {
	code := misc.ErrHTTPNotFound
	msg := "not found"
	respond.JSON(w, http.StatusNotFound, &models.Error{Code: &code, Message: &msg})
}

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

type HandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent reload_agent.IReloadAgent
}

func (h *HandlerImpl) GetGroups(w http.ResponseWriter, r *http.Request, params GetGroupsParams) {
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	_, userlist, err := cfg.GetUserList(params.Userlist, params.TransactionId)
	if userlist == nil || err != nil {
		notFound(w)
		return
	}
	_, groups, err := cfg.GetGroups(params.Userlist, params.TransactionId)
	if err != nil {
		e := misc.HandleContainerGetError(err)
		if *e.Code == misc.ErrHTTPOk {
			respond.JSON(w, http.StatusOK, models.Groups{})
			return
		}
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, groups)
}

func (h *HandlerImpl) CreateGroup(w http.ResponseWriter, r *http.Request, params CreateGroupParams) {
	if params.TransactionId != "" && params.ForceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	var data Group
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	if err = cfg.CreateGroup(params.Userlist, &data, params.TransactionId, int64(params.Version)); err != nil {
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

func (h *HandlerImpl) DeleteGroup(w http.ResponseWriter, r *http.Request, name string, params DeleteGroupParams) {
	if params.TransactionId != "" && params.ForceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	_, userlist, err := cfg.GetUserList(params.Userlist, params.TransactionId)
	if userlist == nil || err != nil {
		notFound(w)
		return
	}
	if err = cfg.DeleteGroup(name, params.Userlist, params.TransactionId, int64(params.Version)); err != nil {
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

func (h *HandlerImpl) GetGroup(w http.ResponseWriter, r *http.Request, name string, params GetGroupParams) {
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	_, g, err := cfg.GetGroup(name, params.Userlist, params.TransactionId)
	if err != nil {
		respond.Error(w, err)
		return
	}
	if g == nil {
		notFound(w)
		return
	}
	respond.JSON(w, http.StatusOK, g)
}

func (h *HandlerImpl) ReplaceGroup(w http.ResponseWriter, r *http.Request, name string, params ReplaceGroupParams) {
	if params.TransactionId != "" && params.ForceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	_, userlist, err := cfg.GetUserList(params.Userlist, params.TransactionId)
	if userlist == nil || err != nil {
		notFound(w)
		return
	}
	_, g, err := cfg.GetGroup(name, params.Userlist, params.TransactionId)
	if g == nil || err != nil {
		notFound(w)
		return
	}
	var data Group
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	if err = cfg.EditGroup(name, params.Userlist, &data, params.TransactionId, int64(params.Version)); err != nil {
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
