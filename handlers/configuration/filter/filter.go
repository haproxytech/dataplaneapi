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

package filter

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	client_native "github.com/haproxytech/client-native/v6"

	"github.com/haproxytech/dataplaneapi/handlers/middleware"
	"github.com/haproxytech/dataplaneapi/handlers/respond"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/reload_agent"
)

// RegisterRouter registers all filter routes onto r using spec-based request validation
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

// HandlerImpl implements ServerInterface for HAProxy filter configuration.
type HandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent reload_agent.IReloadAgent
}

// Backend parent variants

func (h *HandlerImpl) GetAllFilterBackend(w http.ResponseWriter, r *http.Request, parentName string, params GetAllFilterBackendParams) {
	h.getAllFilter(w, r, "backend", parentName, params.TransactionId)
}

func (h *HandlerImpl) ReplaceAllFilterBackend(w http.ResponseWriter, r *http.Request, parentName string, params ReplaceAllFilterBackendParams) {
	h.replaceAllFilter(w, r, "backend", parentName, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) DeleteFilterBackend(w http.ResponseWriter, r *http.Request, parentName string, index int, params DeleteFilterBackendParams) {
	h.deleteFilter(w, r, "backend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) GetFilterBackend(w http.ResponseWriter, r *http.Request, parentName string, index int, params GetFilterBackendParams) {
	h.getFilter(w, r, "backend", parentName, index, params.TransactionId)
}

func (h *HandlerImpl) CreateFilterBackend(w http.ResponseWriter, r *http.Request, parentName string, index int, params CreateFilterBackendParams) {
	h.createFilter(w, r, "backend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) ReplaceFilterBackend(w http.ResponseWriter, r *http.Request, parentName string, index int, params ReplaceFilterBackendParams) {
	h.replaceFilter(w, r, "backend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

// Frontend parent variants

func (h *HandlerImpl) GetAllFilterFrontend(w http.ResponseWriter, r *http.Request, parentName string, params GetAllFilterFrontendParams) {
	h.getAllFilter(w, r, "frontend", parentName, params.TransactionId)
}

func (h *HandlerImpl) ReplaceAllFilterFrontend(w http.ResponseWriter, r *http.Request, parentName string, params ReplaceAllFilterFrontendParams) {
	h.replaceAllFilter(w, r, "frontend", parentName, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) DeleteFilterFrontend(w http.ResponseWriter, r *http.Request, parentName string, index int, params DeleteFilterFrontendParams) {
	h.deleteFilter(w, r, "frontend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) GetFilterFrontend(w http.ResponseWriter, r *http.Request, parentName string, index int, params GetFilterFrontendParams) {
	h.getFilter(w, r, "frontend", parentName, index, params.TransactionId)
}

func (h *HandlerImpl) CreateFilterFrontend(w http.ResponseWriter, r *http.Request, parentName string, index int, params CreateFilterFrontendParams) {
	h.createFilter(w, r, "frontend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) ReplaceFilterFrontend(w http.ResponseWriter, r *http.Request, parentName string, index int, params ReplaceFilterFrontendParams) {
	h.replaceFilter(w, r, "frontend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

// Shared internal methods

func (h *HandlerImpl) getAllFilter(w http.ResponseWriter, r *http.Request, parentType, parentName, txID string) {
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	_, fs, err := cfg.GetFilters(parentType, parentName, txID)
	if err != nil {
		e := misc.HandleContainerGetError(err)
		if *e.Code == misc.ErrHTTPOk {
			respond.JSON(w, http.StatusOK, Filters{})
			return
		}
		respond.JSON(w, int(*e.Code), e)
		return
	}
	respond.JSON(w, http.StatusOK, fs)
}

func (h *HandlerImpl) replaceAllFilter(w http.ResponseWriter, r *http.Request, parentType, parentName, txID string, version int64, forceReload bool) {
	if txID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data Filters
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.ReplaceFilters(parentType, parentName, data, txID, version); err != nil {
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

func (h *HandlerImpl) deleteFilter(w http.ResponseWriter, r *http.Request, parentType, parentName string, index int, txID string, version int64, forceReload bool) {
	if txID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.DeleteFilter(int64(index), parentType, parentName, txID, version); err != nil {
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

func (h *HandlerImpl) getFilter(w http.ResponseWriter, r *http.Request, parentType, parentName string, index int, txID string) {
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	_, f, err := cfg.GetFilter(int64(index), parentType, parentName, txID)
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, f)
}

func (h *HandlerImpl) createFilter(w http.ResponseWriter, r *http.Request, parentType, parentName string, index int, txID string, version int64, forceReload bool) {
	if txID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data Filter
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.CreateFilter(int64(index), parentType, parentName, &data, txID, version); err != nil {
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

func (h *HandlerImpl) replaceFilter(w http.ResponseWriter, r *http.Request, parentType, parentName string, index int, txID string, version int64, forceReload bool) {
	if txID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data Filter
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.EditFilter(int64(index), parentType, parentName, &data, txID, version); err != nil {
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
