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

package filter_sequence

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

// RegisterRouter registers all filter_sequence routes onto r using spec-based request validation
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

// HandlerImpl implements ServerInterface for HAProxy filter-sequence configuration.
type HandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent reload_agent.IReloadAgent
}

// Backend parent variants

func (h *HandlerImpl) GetAllFilterSequenceBackend(w http.ResponseWriter, r *http.Request, parentName string, params GetAllFilterSequenceBackendParams) {
	h.getAllFilterSequence(w, r, "backend", parentName, params.TransactionId)
}

func (h *HandlerImpl) ReplaceAllFilterSequenceBackend(w http.ResponseWriter, r *http.Request, parentName string, params ReplaceAllFilterSequenceBackendParams) {
	h.replaceAllFilterSequence(w, r, "backend", parentName, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) DeleteFilterSequenceBackend(w http.ResponseWriter, r *http.Request, parentName string, index int, params DeleteFilterSequenceBackendParams) {
	h.deleteFilterSequence(w, r, "backend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) GetFilterSequenceBackend(w http.ResponseWriter, r *http.Request, parentName string, index int, params GetFilterSequenceBackendParams) {
	h.getFilterSequence(w, r, "backend", parentName, index, params.TransactionId)
}

func (h *HandlerImpl) CreateFilterSequenceBackend(w http.ResponseWriter, r *http.Request, parentName string, index int, params CreateFilterSequenceBackendParams) {
	h.createFilterSequence(w, r, "backend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) ReplaceFilterSequenceBackend(w http.ResponseWriter, r *http.Request, parentName string, index int, params ReplaceFilterSequenceBackendParams) {
	h.replaceFilterSequence(w, r, "backend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

// Frontend parent variants

func (h *HandlerImpl) GetAllFilterSequenceFrontend(w http.ResponseWriter, r *http.Request, parentName string, params GetAllFilterSequenceFrontendParams) {
	h.getAllFilterSequence(w, r, "frontend", parentName, params.TransactionId)
}

func (h *HandlerImpl) ReplaceAllFilterSequenceFrontend(w http.ResponseWriter, r *http.Request, parentName string, params ReplaceAllFilterSequenceFrontendParams) {
	h.replaceAllFilterSequence(w, r, "frontend", parentName, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) DeleteFilterSequenceFrontend(w http.ResponseWriter, r *http.Request, parentName string, index int, params DeleteFilterSequenceFrontendParams) {
	h.deleteFilterSequence(w, r, "frontend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) GetFilterSequenceFrontend(w http.ResponseWriter, r *http.Request, parentName string, index int, params GetFilterSequenceFrontendParams) {
	h.getFilterSequence(w, r, "frontend", parentName, index, params.TransactionId)
}

func (h *HandlerImpl) CreateFilterSequenceFrontend(w http.ResponseWriter, r *http.Request, parentName string, index int, params CreateFilterSequenceFrontendParams) {
	h.createFilterSequence(w, r, "frontend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) ReplaceFilterSequenceFrontend(w http.ResponseWriter, r *http.Request, parentName string, index int, params ReplaceFilterSequenceFrontendParams) {
	h.replaceFilterSequence(w, r, "frontend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

// Defaults parent variants

func (h *HandlerImpl) GetAllFilterSequenceDefaults(w http.ResponseWriter, r *http.Request, parentName string, params GetAllFilterSequenceDefaultsParams) {
	h.getAllFilterSequence(w, r, "defaults", parentName, params.TransactionId)
}

func (h *HandlerImpl) ReplaceAllFilterSequenceDefaults(w http.ResponseWriter, r *http.Request, parentName string, params ReplaceAllFilterSequenceDefaultsParams) {
	h.replaceAllFilterSequence(w, r, "defaults", parentName, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) DeleteFilterSequenceDefaults(w http.ResponseWriter, r *http.Request, parentName string, index int, params DeleteFilterSequenceDefaultsParams) {
	h.deleteFilterSequence(w, r, "defaults", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) GetFilterSequenceDefaults(w http.ResponseWriter, r *http.Request, parentName string, index int, params GetFilterSequenceDefaultsParams) {
	h.getFilterSequence(w, r, "defaults", parentName, index, params.TransactionId)
}

func (h *HandlerImpl) CreateFilterSequenceDefaults(w http.ResponseWriter, r *http.Request, parentName string, index int, params CreateFilterSequenceDefaultsParams) {
	h.createFilterSequence(w, r, "defaults", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) ReplaceFilterSequenceDefaults(w http.ResponseWriter, r *http.Request, parentName string, index int, params ReplaceFilterSequenceDefaultsParams) {
	h.replaceFilterSequence(w, r, "defaults", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

// Shared internal methods

func (h *HandlerImpl) getAllFilterSequence(w http.ResponseWriter, r *http.Request, parentType, parentName, txID string) {
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	_, rules, err := cfg.GetFilterSequences(parentType, parentName, txID)
	if err != nil {
		e := misc.HandleContainerGetError(err)
		if *e.Code == misc.ErrHTTPOk {
			respond.JSON(w, http.StatusOK, FilterSequences{})
			return
		}
		respond.JSON(w, int(*e.Code), e)
		return
	}
	respond.JSON(w, http.StatusOK, rules)
}

func (h *HandlerImpl) replaceAllFilterSequence(w http.ResponseWriter, r *http.Request, parentType, parentName, txID string, version int64, forceReload bool) {
	if txID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data FilterSequences
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.ReplaceFilterSequences(parentType, parentName, data, txID, version); err != nil {
		respond.Error(w, err)
		return
	}
	if txID == "" {
		respond.Committed(w, h.ReloadAgent, cn.NeedsReload(version), forceReload, http.StatusOK, data)
		return
	}
	respond.JSON(w, http.StatusAccepted, data)
}

func (h *HandlerImpl) deleteFilterSequence(w http.ResponseWriter, r *http.Request, parentType, parentName string, index int, txID string, version int64, forceReload bool) {
	if txID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.DeleteFilterSequence(int64(index), parentType, parentName, txID, version); err != nil {
		respond.Error(w, err)
		return
	}
	if txID == "" {
		respond.Committed(w, h.ReloadAgent, cn.NeedsReload(version), forceReload, http.StatusNoContent, nil)
		return
	}
	respond.Accepted(w, "", nil)
}

func (h *HandlerImpl) getFilterSequence(w http.ResponseWriter, r *http.Request, parentType, parentName string, index int, txID string) {
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	_, rule, err := cfg.GetFilterSequence(int64(index), parentType, parentName, txID)
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, rule)
}

func (h *HandlerImpl) createFilterSequence(w http.ResponseWriter, r *http.Request, parentType, parentName string, index int, txID string, version int64, forceReload bool) {
	if txID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data FilterSequence
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.CreateFilterSequence(int64(index), parentType, parentName, &data, txID, version); err != nil {
		respond.Error(w, err)
		return
	}
	if txID == "" {
		respond.Committed(w, h.ReloadAgent, cn.NeedsReload(version), forceReload, http.StatusCreated, &data)
		return
	}
	respond.JSON(w, http.StatusAccepted, &data)
}

func (h *HandlerImpl) replaceFilterSequence(w http.ResponseWriter, r *http.Request, parentType, parentName string, index int, txID string, version int64, forceReload bool) {
	if txID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data FilterSequence
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.EditFilterSequence(int64(index), parentType, parentName, &data, txID, version); err != nil {
		respond.Error(w, err)
		return
	}
	if txID == "" {
		respond.Committed(w, h.ReloadAgent, cn.NeedsReload(version), forceReload, http.StatusOK, &data)
		return
	}
	respond.JSON(w, http.StatusAccepted, &data)
}
