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

package filter_rule

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	client_native "github.com/haproxytech/client-native/v6"

	"github.com/haproxytech/dataplaneapi/handlers/middleware"
	"github.com/haproxytech/dataplaneapi/handlers/respond"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/reload_agent"
)

// RegisterRouter registers all filter_rule routes onto r using spec-based request validation
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

// HandlerImpl implements ServerInterface for HAProxy filter-rule configuration.
type HandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent reload_agent.IReloadAgent
}

// Backend parent variants

func (h *HandlerImpl) GetAllFilterRuleBackend(w http.ResponseWriter, r *http.Request, parentName string, params GetAllFilterRuleBackendParams) {
	h.getAllFilterRule(w, r, "backend", parentName, params.TransactionId)
}

func (h *HandlerImpl) ReplaceAllFilterRuleBackend(w http.ResponseWriter, r *http.Request, parentName string, params ReplaceAllFilterRuleBackendParams) {
	h.replaceAllFilterRule(w, r, "backend", parentName, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) DeleteFilterRuleBackend(w http.ResponseWriter, r *http.Request, parentName string, index int, params DeleteFilterRuleBackendParams) {
	h.deleteFilterRule(w, r, "backend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) GetFilterRuleBackend(w http.ResponseWriter, r *http.Request, parentName string, index int, params GetFilterRuleBackendParams) {
	h.getFilterRule(w, r, "backend", parentName, index, params.TransactionId)
}

func (h *HandlerImpl) CreateFilterRuleBackend(w http.ResponseWriter, r *http.Request, parentName string, index int, params CreateFilterRuleBackendParams) {
	h.createFilterRule(w, r, "backend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) ReplaceFilterRuleBackend(w http.ResponseWriter, r *http.Request, parentName string, index int, params ReplaceFilterRuleBackendParams) {
	h.replaceFilterRule(w, r, "backend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

// Frontend parent variants

func (h *HandlerImpl) GetAllFilterRuleFrontend(w http.ResponseWriter, r *http.Request, parentName string, params GetAllFilterRuleFrontendParams) {
	h.getAllFilterRule(w, r, "frontend", parentName, params.TransactionId)
}

func (h *HandlerImpl) ReplaceAllFilterRuleFrontend(w http.ResponseWriter, r *http.Request, parentName string, params ReplaceAllFilterRuleFrontendParams) {
	h.replaceAllFilterRule(w, r, "frontend", parentName, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) DeleteFilterRuleFrontend(w http.ResponseWriter, r *http.Request, parentName string, index int, params DeleteFilterRuleFrontendParams) {
	h.deleteFilterRule(w, r, "frontend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) GetFilterRuleFrontend(w http.ResponseWriter, r *http.Request, parentName string, index int, params GetFilterRuleFrontendParams) {
	h.getFilterRule(w, r, "frontend", parentName, index, params.TransactionId)
}

func (h *HandlerImpl) CreateFilterRuleFrontend(w http.ResponseWriter, r *http.Request, parentName string, index int, params CreateFilterRuleFrontendParams) {
	h.createFilterRule(w, r, "frontend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) ReplaceFilterRuleFrontend(w http.ResponseWriter, r *http.Request, parentName string, index int, params ReplaceFilterRuleFrontendParams) {
	h.replaceFilterRule(w, r, "frontend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

// Defaults parent variants

func (h *HandlerImpl) GetAllFilterRuleDefaults(w http.ResponseWriter, r *http.Request, parentName string, params GetAllFilterRuleDefaultsParams) {
	h.getAllFilterRule(w, r, "defaults", parentName, params.TransactionId)
}

func (h *HandlerImpl) ReplaceAllFilterRuleDefaults(w http.ResponseWriter, r *http.Request, parentName string, params ReplaceAllFilterRuleDefaultsParams) {
	h.replaceAllFilterRule(w, r, "defaults", parentName, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) DeleteFilterRuleDefaults(w http.ResponseWriter, r *http.Request, parentName string, index int, params DeleteFilterRuleDefaultsParams) {
	h.deleteFilterRule(w, r, "defaults", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) GetFilterRuleDefaults(w http.ResponseWriter, r *http.Request, parentName string, index int, params GetFilterRuleDefaultsParams) {
	h.getFilterRule(w, r, "defaults", parentName, index, params.TransactionId)
}

func (h *HandlerImpl) CreateFilterRuleDefaults(w http.ResponseWriter, r *http.Request, parentName string, index int, params CreateFilterRuleDefaultsParams) {
	h.createFilterRule(w, r, "defaults", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) ReplaceFilterRuleDefaults(w http.ResponseWriter, r *http.Request, parentName string, index int, params ReplaceFilterRuleDefaultsParams) {
	h.replaceFilterRule(w, r, "defaults", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

// Shared internal methods

func (h *HandlerImpl) getAllFilterRule(w http.ResponseWriter, r *http.Request, parentType, parentName, txID string) {
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	_, rules, err := cfg.GetFilterRules(parentType, parentName, txID)
	if err != nil {
		e := misc.HandleContainerGetError(err)
		if *e.Code == misc.ErrHTTPOk {
			respond.JSON(w, http.StatusOK, FilterRules{})
			return
		}
		respond.JSON(w, int(*e.Code), e)
		return
	}
	respond.JSON(w, http.StatusOK, rules)
}

func (h *HandlerImpl) replaceAllFilterRule(w http.ResponseWriter, r *http.Request, parentType, parentName, txID string, version int64, forceReload bool) {
	if txID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data FilterRules
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.ReplaceFilterRules(parentType, parentName, data, txID, version); err != nil {
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

func (h *HandlerImpl) deleteFilterRule(w http.ResponseWriter, r *http.Request, parentType, parentName string, index int, txID string, version int64, forceReload bool) {
	if txID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.DeleteFilterRule(int64(index), parentType, parentName, txID, version); err != nil {
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

func (h *HandlerImpl) getFilterRule(w http.ResponseWriter, r *http.Request, parentType, parentName string, index int, txID string) {
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	_, rule, err := cfg.GetFilterRule(int64(index), parentType, parentName, txID)
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, rule)
}

func (h *HandlerImpl) createFilterRule(w http.ResponseWriter, r *http.Request, parentType, parentName string, index int, txID string, version int64, forceReload bool) {
	if txID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data FilterRule
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.CreateFilterRule(int64(index), parentType, parentName, &data, txID, version); err != nil {
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

func (h *HandlerImpl) replaceFilterRule(w http.ResponseWriter, r *http.Request, parentType, parentName string, index int, txID string, version int64, forceReload bool) {
	if txID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data FilterRule
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.EditFilterRule(int64(index), parentType, parentName, &data, txID, version); err != nil {
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
