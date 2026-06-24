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

package http_check

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	client_native "github.com/haproxytech/client-native/v6"

	"github.com/haproxytech/dataplaneapi/handlers/middleware"
	"github.com/haproxytech/dataplaneapi/handlers/respond"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/reload_agent"
)

// RegisterRouter registers all HTTP Check routes onto r using spec-based request validation
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

// HandlerImpl implements ServerInterface for HAProxy HTTP Check configuration.
type HandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent reload_agent.IReloadAgent
}

// --- Backend ---

func (h *HandlerImpl) GetAllHTTPCheckBackend(w http.ResponseWriter, r *http.Request, parentName string, params GetAllHTTPCheckBackendParams) {
	h.getAllHTTPCheck(w, r, "backend", parentName, params.TransactionId)
}

func (h *HandlerImpl) ReplaceAllHTTPCheckBackend(w http.ResponseWriter, r *http.Request, parentName string, params ReplaceAllHTTPCheckBackendParams) {
	h.replaceAllHTTPCheck(w, r, "backend", parentName, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) DeleteHTTPCheckBackend(w http.ResponseWriter, r *http.Request, parentName string, index int, params DeleteHTTPCheckBackendParams) {
	h.deleteHTTPCheck(w, r, "backend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) GetHTTPCheckBackend(w http.ResponseWriter, r *http.Request, parentName string, index int, params GetHTTPCheckBackendParams) {
	h.getHTTPCheck(w, r, "backend", parentName, index, params.TransactionId)
}

func (h *HandlerImpl) CreateHTTPCheckBackend(w http.ResponseWriter, r *http.Request, parentName string, index int, params CreateHTTPCheckBackendParams) {
	h.createHTTPCheck(w, r, "backend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) ReplaceHTTPCheckBackend(w http.ResponseWriter, r *http.Request, parentName string, index int, params ReplaceHTTPCheckBackendParams) {
	h.replaceHTTPCheck(w, r, "backend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

// --- Defaults ---

func (h *HandlerImpl) GetAllHTTPCheckDefaults(w http.ResponseWriter, r *http.Request, parentName string, params GetAllHTTPCheckDefaultsParams) {
	h.getAllHTTPCheck(w, r, "defaults", parentName, params.TransactionId)
}

func (h *HandlerImpl) ReplaceAllHTTPCheckDefaults(w http.ResponseWriter, r *http.Request, parentName string, params ReplaceAllHTTPCheckDefaultsParams) {
	h.replaceAllHTTPCheck(w, r, "defaults", parentName, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) DeleteHTTPCheckDefaults(w http.ResponseWriter, r *http.Request, parentName string, index int, params DeleteHTTPCheckDefaultsParams) {
	h.deleteHTTPCheck(w, r, "defaults", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) GetHTTPCheckDefaults(w http.ResponseWriter, r *http.Request, parentName string, index int, params GetHTTPCheckDefaultsParams) {
	h.getHTTPCheck(w, r, "defaults", parentName, index, params.TransactionId)
}

func (h *HandlerImpl) CreateHTTPCheckDefaults(w http.ResponseWriter, r *http.Request, parentName string, index int, params CreateHTTPCheckDefaultsParams) {
	h.createHTTPCheck(w, r, "defaults", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) ReplaceHTTPCheckDefaults(w http.ResponseWriter, r *http.Request, parentName string, index int, params ReplaceHTTPCheckDefaultsParams) {
	h.replaceHTTPCheck(w, r, "defaults", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

// --- Shared implementations ---

func (h *HandlerImpl) getAllHTTPCheck(w http.ResponseWriter, r *http.Request, parentType, parentName, txID string) {
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	_, checks, err := cfg.GetHTTPChecks(parentType, parentName, txID)
	if err != nil {
		e := misc.HandleContainerGetError(err)
		if *e.Code == misc.ErrHTTPOk {
			respond.JSON(w, http.StatusOK, HttpChecks{})
			return
		}
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, checks)
}

func (h *HandlerImpl) getHTTPCheck(w http.ResponseWriter, r *http.Request, parentType, parentName string, index int, txID string) {
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	_, check, err := cfg.GetHTTPCheck(int64(index), parentType, parentName, txID)
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, check)
}

func (h *HandlerImpl) createHTTPCheck(w http.ResponseWriter, r *http.Request, parentType, parentName string, index int, txID string, version int64, forceReload bool) {
	if txID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data HttpCheck
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.CreateHTTPCheck(int64(index), parentType, parentName, &data, txID, version); err != nil {
		respond.Error(w, err)
		return
	}
	if txID == "" {
		if forceReload {
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

func (h *HandlerImpl) deleteHTTPCheck(w http.ResponseWriter, r *http.Request, parentType, parentName string, index int, txID string, version int64, forceReload bool) {
	if txID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.DeleteHTTPCheck(int64(index), parentType, parentName, txID, version); err != nil {
		respond.Error(w, err)
		return
	}
	if txID == "" {
		if forceReload {
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

func (h *HandlerImpl) replaceHTTPCheck(w http.ResponseWriter, r *http.Request, parentType, parentName string, index int, txID string, version int64, forceReload bool) {
	if txID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data HttpCheck
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.EditHTTPCheck(int64(index), parentType, parentName, &data, txID, version); err != nil {
		respond.Error(w, err)
		return
	}
	if txID == "" {
		if forceReload {
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

func (h *HandlerImpl) replaceAllHTTPCheck(w http.ResponseWriter, r *http.Request, parentType, parentName, txID string, version int64, forceReload bool) {
	if txID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data HttpChecks
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.ReplaceHTTPChecks(parentType, parentName, data, txID, version); err != nil {
		respond.Error(w, err)
		return
	}
	if txID == "" {
		if forceReload {
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
