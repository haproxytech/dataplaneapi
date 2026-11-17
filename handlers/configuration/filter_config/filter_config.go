// Copyright 2026 HAProxy Technologies
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

package filter_config

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	client_native "github.com/haproxytech/client-native/v6"

	cn "github.com/haproxytech/dataplaneapi/client-native"
	"github.com/haproxytech/dataplaneapi/handlers/middleware"
	"github.com/haproxytech/dataplaneapi/handlers/respond"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/reload_agent"
)

// RegisterRouter registers all filter_config routes onto r using spec-based request validation
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

// HandlerImpl implements ServerInterface for HAProxy filter-config configuration.
type HandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent reload_agent.IReloadAgent
}

// Backend parent variants

func (h *HandlerImpl) GetAllFilterConfigBackend(w http.ResponseWriter, r *http.Request, parentName string, params GetAllFilterConfigBackendParams) {
	h.getAllFilterConfig(w, r, "backend", parentName, params.TransactionId)
}

func (h *HandlerImpl) ReplaceAllFilterConfigBackend(w http.ResponseWriter, r *http.Request, parentName string, params ReplaceAllFilterConfigBackendParams) {
	h.replaceAllFilterConfig(w, r, "backend", parentName, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) DeleteFilterConfigBackend(w http.ResponseWriter, r *http.Request, parentName string, index int, params DeleteFilterConfigBackendParams) {
	h.deleteFilterConfig(w, r, "backend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) GetFilterConfigBackend(w http.ResponseWriter, r *http.Request, parentName string, index int, params GetFilterConfigBackendParams) {
	h.getFilterConfig(w, r, "backend", parentName, index, params.TransactionId)
}

func (h *HandlerImpl) CreateFilterConfigBackend(w http.ResponseWriter, r *http.Request, parentName string, index int, params CreateFilterConfigBackendParams) {
	h.createFilterConfig(w, r, "backend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) ReplaceFilterConfigBackend(w http.ResponseWriter, r *http.Request, parentName string, index int, params ReplaceFilterConfigBackendParams) {
	h.replaceFilterConfig(w, r, "backend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

// Frontend parent variants

func (h *HandlerImpl) GetAllFilterConfigFrontend(w http.ResponseWriter, r *http.Request, parentName string, params GetAllFilterConfigFrontendParams) {
	h.getAllFilterConfig(w, r, "frontend", parentName, params.TransactionId)
}

func (h *HandlerImpl) ReplaceAllFilterConfigFrontend(w http.ResponseWriter, r *http.Request, parentName string, params ReplaceAllFilterConfigFrontendParams) {
	h.replaceAllFilterConfig(w, r, "frontend", parentName, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) DeleteFilterConfigFrontend(w http.ResponseWriter, r *http.Request, parentName string, index int, params DeleteFilterConfigFrontendParams) {
	h.deleteFilterConfig(w, r, "frontend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) GetFilterConfigFrontend(w http.ResponseWriter, r *http.Request, parentName string, index int, params GetFilterConfigFrontendParams) {
	h.getFilterConfig(w, r, "frontend", parentName, index, params.TransactionId)
}

func (h *HandlerImpl) CreateFilterConfigFrontend(w http.ResponseWriter, r *http.Request, parentName string, index int, params CreateFilterConfigFrontendParams) {
	h.createFilterConfig(w, r, "frontend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) ReplaceFilterConfigFrontend(w http.ResponseWriter, r *http.Request, parentName string, index int, params ReplaceFilterConfigFrontendParams) {
	h.replaceFilterConfig(w, r, "frontend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

// Defaults parent variants

func (h *HandlerImpl) GetAllFilterConfigDefaults(w http.ResponseWriter, r *http.Request, parentName string, params GetAllFilterConfigDefaultsParams) {
	h.getAllFilterConfig(w, r, "defaults", parentName, params.TransactionId)
}

func (h *HandlerImpl) ReplaceAllFilterConfigDefaults(w http.ResponseWriter, r *http.Request, parentName string, params ReplaceAllFilterConfigDefaultsParams) {
	h.replaceAllFilterConfig(w, r, "defaults", parentName, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) DeleteFilterConfigDefaults(w http.ResponseWriter, r *http.Request, parentName string, index int, params DeleteFilterConfigDefaultsParams) {
	h.deleteFilterConfig(w, r, "defaults", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) GetFilterConfigDefaults(w http.ResponseWriter, r *http.Request, parentName string, index int, params GetFilterConfigDefaultsParams) {
	h.getFilterConfig(w, r, "defaults", parentName, index, params.TransactionId)
}

func (h *HandlerImpl) CreateFilterConfigDefaults(w http.ResponseWriter, r *http.Request, parentName string, index int, params CreateFilterConfigDefaultsParams) {
	h.createFilterConfig(w, r, "defaults", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) ReplaceFilterConfigDefaults(w http.ResponseWriter, r *http.Request, parentName string, index int, params ReplaceFilterConfigDefaultsParams) {
	h.replaceFilterConfig(w, r, "defaults", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

// Shared internal methods

func (h *HandlerImpl) getAllFilterConfig(w http.ResponseWriter, r *http.Request, parentType, parentName, txID string) {
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	_, configs, err := cfg.GetFilterConfigs(parentType, parentName, txID)
	if err != nil {
		e := misc.HandleContainerGetError(err)
		if *e.Code == misc.ErrHTTPOk {
			respond.JSON(w, http.StatusOK, FilterConfigs{})
			return
		}
		respond.JSON(w, int(*e.Code), e)
		return
	}
	respond.JSON(w, http.StatusOK, configs)
}

func (h *HandlerImpl) replaceAllFilterConfig(w http.ResponseWriter, r *http.Request, parentType, parentName, txID string, version int64, forceReload bool) {
	if txID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data FilterConfigs
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.ReplaceFilterConfigs(parentType, parentName, data, txID, version); err != nil {
		respond.Error(w, err)
		return
	}
	if txID == "" {
		respond.Committed(w, h.ReloadAgent, cn.NeedsReload(version), forceReload, http.StatusOK, data)
		return
	}
	respond.JSON(w, http.StatusAccepted, data)
}

func (h *HandlerImpl) deleteFilterConfig(w http.ResponseWriter, r *http.Request, parentType, parentName string, index int, txID string, version int64, forceReload bool) {
	if txID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.DeleteFilterConfig(int64(index), parentType, parentName, txID, version); err != nil {
		respond.Error(w, err)
		return
	}
	if txID == "" {
		respond.Committed(w, h.ReloadAgent, cn.NeedsReload(version), forceReload, http.StatusNoContent, nil)
		return
	}
	respond.Accepted(w, "", nil)
}

func (h *HandlerImpl) getFilterConfig(w http.ResponseWriter, r *http.Request, parentType, parentName string, index int, txID string) {
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	_, config, err := cfg.GetFilterConfig(int64(index), parentType, parentName, txID)
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, config)
}

func (h *HandlerImpl) createFilterConfig(w http.ResponseWriter, r *http.Request, parentType, parentName string, index int, txID string, version int64, forceReload bool) {
	if txID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data FilterConfig
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.CreateFilterConfig(int64(index), parentType, parentName, &data, txID, version); err != nil {
		respond.Error(w, err)
		return
	}
	if txID == "" {
		respond.Committed(w, h.ReloadAgent, cn.NeedsReload(version), forceReload, http.StatusCreated, &data)
		return
	}
	respond.JSON(w, http.StatusAccepted, &data)
}

func (h *HandlerImpl) replaceFilterConfig(w http.ResponseWriter, r *http.Request, parentType, parentName string, index int, txID string, version int64, forceReload bool) {
	if txID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data FilterConfig
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.EditFilterConfig(int64(index), parentType, parentName, &data, txID, version); err != nil {
		respond.Error(w, err)
		return
	}
	if txID == "" {
		respond.Committed(w, h.ReloadAgent, cn.NeedsReload(version), forceReload, http.StatusOK, &data)
		return
	}
	respond.JSON(w, http.StatusAccepted, &data)
}
