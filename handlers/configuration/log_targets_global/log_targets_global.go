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

package log_targets_global

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

// RegisterRouter registers all global log_target routes onto r using spec-based request validation
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

// HandlerImpl implements ServerInterface for HAProxy global log_target configuration.
type HandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent reload_agent.IReloadAgent
}

func (h *HandlerImpl) GetAllLogTargetGlobal(w http.ResponseWriter, r *http.Request, params GetAllLogTargetGlobalParams) {
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	_, logTargets, err := cfg.GetLogTargets("global", "", params.TransactionId)
	if err != nil {
		e := misc.HandleContainerGetError(err)
		if *e.Code == misc.ErrHTTPOk {
			respond.JSON(w, http.StatusOK, models.LogTargets{})
			return
		}
		respond.JSON(w, int(*e.Code), e)
		return
	}
	respond.JSON(w, http.StatusOK, logTargets)
}

func (h *HandlerImpl) ReplaceAllLogTargetGlobal(w http.ResponseWriter, r *http.Request, params ReplaceAllLogTargetGlobalParams) {
	if params.TransactionId != "" && params.ForceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data LogTargets
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.ReplaceLogTargets("global", "", data, params.TransactionId, int64(params.Version)); err != nil {
		respond.Error(w, err)
		return
	}
	if params.TransactionId == "" {
		if params.ForceReload {
			if err = h.ReloadAgent.ForceReload(); err != nil {
				respond.Error(w, err)
				return
			}
			respond.JSON(w, http.StatusOK, data)
			return
		}
		respond.Accepted(w, h.ReloadAgent.Reload(), data)
		return
	}
	respond.JSON(w, http.StatusAccepted, data)
}

func (h *HandlerImpl) DeleteLogTargetGlobal(w http.ResponseWriter, r *http.Request, index int, params DeleteLogTargetGlobalParams) {
	if params.TransactionId != "" && params.ForceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.DeleteLogTarget(int64(index), "global", "", params.TransactionId, int64(params.Version)); err != nil {
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

func (h *HandlerImpl) GetLogTargetGlobal(w http.ResponseWriter, r *http.Request, index int, params GetLogTargetGlobalParams) {
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	_, logTarget, err := cfg.GetLogTarget(int64(index), "global", "", params.TransactionId)
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, logTarget)
}

func (h *HandlerImpl) CreateLogTargetGlobal(w http.ResponseWriter, r *http.Request, index int, params CreateLogTargetGlobalParams) {
	if params.TransactionId != "" && params.ForceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data LogTarget
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.CreateLogTarget(int64(index), "global", "", &data, params.TransactionId, int64(params.Version)); err != nil {
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

func (h *HandlerImpl) ReplaceLogTargetGlobal(w http.ResponseWriter, r *http.Request, index int, params ReplaceLogTargetGlobalParams) {
	if params.TransactionId != "" && params.ForceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data LogTarget
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.EditLogTarget(int64(index), "global", "", &data, params.TransactionId, int64(params.Version)); err != nil {
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
