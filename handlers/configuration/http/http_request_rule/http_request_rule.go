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

package http_request_rule

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	client_native "github.com/haproxytech/client-native/v6"

	"github.com/haproxytech/dataplaneapi/handlers/middleware"
	"github.com/haproxytech/dataplaneapi/handlers/respond"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/reload_agent"
)

// RegisterRouter registers all HTTP Request Rule routes onto r using spec-based request validation
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

// HandlerImpl implements ServerInterface for HAProxy HTTP Request Rule configuration.
type HandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent reload_agent.IReloadAgent
}

// --- Backend ---

func (h *HandlerImpl) GetAllHTTPRequestRuleBackend(w http.ResponseWriter, r *http.Request, parentName string, params GetAllHTTPRequestRuleBackendParams) {
	h.getAllHTTPRequestRule(w, r, "backend", parentName, params.TransactionId)
}

func (h *HandlerImpl) ReplaceAllHTTPRequestRuleBackend(w http.ResponseWriter, r *http.Request, parentName string, params ReplaceAllHTTPRequestRuleBackendParams) {
	h.replaceAllHTTPRequestRule(w, r, "backend", parentName, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) DeleteHTTPRequestRuleBackend(w http.ResponseWriter, r *http.Request, parentName string, index int, params DeleteHTTPRequestRuleBackendParams) {
	h.deleteHTTPRequestRule(w, r, "backend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) GetHTTPRequestRuleBackend(w http.ResponseWriter, r *http.Request, parentName string, index int, params GetHTTPRequestRuleBackendParams) {
	h.getHTTPRequestRule(w, r, "backend", parentName, index, params.TransactionId)
}

func (h *HandlerImpl) CreateHTTPRequestRuleBackend(w http.ResponseWriter, r *http.Request, parentName string, index int, params CreateHTTPRequestRuleBackendParams) {
	h.createHTTPRequestRule(w, r, "backend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) ReplaceHTTPRequestRuleBackend(w http.ResponseWriter, r *http.Request, parentName string, index int, params ReplaceHTTPRequestRuleBackendParams) {
	h.replaceHTTPRequestRule(w, r, "backend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

// --- Defaults ---

func (h *HandlerImpl) GetAllHTTPRequestRuleDefaults(w http.ResponseWriter, r *http.Request, parentName string, params GetAllHTTPRequestRuleDefaultsParams) {
	h.getAllHTTPRequestRule(w, r, "defaults", parentName, params.TransactionId)
}

func (h *HandlerImpl) ReplaceAllHTTPRequestRuleDefaults(w http.ResponseWriter, r *http.Request, parentName string, params ReplaceAllHTTPRequestRuleDefaultsParams) {
	h.replaceAllHTTPRequestRule(w, r, "defaults", parentName, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) DeleteHTTPRequestRuleDefaults(w http.ResponseWriter, r *http.Request, parentName string, index int, params DeleteHTTPRequestRuleDefaultsParams) {
	h.deleteHTTPRequestRule(w, r, "defaults", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) GetHTTPRequestRuleDefaults(w http.ResponseWriter, r *http.Request, parentName string, index int, params GetHTTPRequestRuleDefaultsParams) {
	h.getHTTPRequestRule(w, r, "defaults", parentName, index, params.TransactionId)
}

func (h *HandlerImpl) CreateHTTPRequestRuleDefaults(w http.ResponseWriter, r *http.Request, parentName string, index int, params CreateHTTPRequestRuleDefaultsParams) {
	h.createHTTPRequestRule(w, r, "defaults", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) ReplaceHTTPRequestRuleDefaults(w http.ResponseWriter, r *http.Request, parentName string, index int, params ReplaceHTTPRequestRuleDefaultsParams) {
	h.replaceHTTPRequestRule(w, r, "defaults", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

// --- Frontend ---

func (h *HandlerImpl) GetAllHTTPRequestRuleFrontend(w http.ResponseWriter, r *http.Request, parentName string, params GetAllHTTPRequestRuleFrontendParams) {
	h.getAllHTTPRequestRule(w, r, "frontend", parentName, params.TransactionId)
}

func (h *HandlerImpl) ReplaceAllHTTPRequestRuleFrontend(w http.ResponseWriter, r *http.Request, parentName string, params ReplaceAllHTTPRequestRuleFrontendParams) {
	h.replaceAllHTTPRequestRule(w, r, "frontend", parentName, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) DeleteHTTPRequestRuleFrontend(w http.ResponseWriter, r *http.Request, parentName string, index int, params DeleteHTTPRequestRuleFrontendParams) {
	h.deleteHTTPRequestRule(w, r, "frontend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) GetHTTPRequestRuleFrontend(w http.ResponseWriter, r *http.Request, parentName string, index int, params GetHTTPRequestRuleFrontendParams) {
	h.getHTTPRequestRule(w, r, "frontend", parentName, index, params.TransactionId)
}

func (h *HandlerImpl) CreateHTTPRequestRuleFrontend(w http.ResponseWriter, r *http.Request, parentName string, index int, params CreateHTTPRequestRuleFrontendParams) {
	h.createHTTPRequestRule(w, r, "frontend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) ReplaceHTTPRequestRuleFrontend(w http.ResponseWriter, r *http.Request, parentName string, index int, params ReplaceHTTPRequestRuleFrontendParams) {
	h.replaceHTTPRequestRule(w, r, "frontend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

// --- Shared implementations ---

func (h *HandlerImpl) getAllHTTPRequestRule(w http.ResponseWriter, r *http.Request, parentType, parentName, txID string) {
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	_, rules, err := cfg.GetHTTPRequestRules(parentType, parentName, txID)
	if err != nil {
		e := misc.HandleContainerGetError(err)
		if *e.Code == misc.ErrHTTPOk {
			respond.JSON(w, http.StatusOK, HttpRequestRules{})
			return
		}
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, rules)
}

func (h *HandlerImpl) getHTTPRequestRule(w http.ResponseWriter, r *http.Request, parentType, parentName string, index int, txID string) {
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	_, rule, err := cfg.GetHTTPRequestRule(int64(index), parentType, parentName, txID)
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, rule)
}

func (h *HandlerImpl) createHTTPRequestRule(w http.ResponseWriter, r *http.Request, parentType, parentName string, index int, txID string, version int64, forceReload bool) {
	if txID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data HttpRequestRule
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.CreateHTTPRequestRule(int64(index), parentType, parentName, &data, txID, version); err != nil {
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

func (h *HandlerImpl) deleteHTTPRequestRule(w http.ResponseWriter, r *http.Request, parentType, parentName string, index int, txID string, version int64, forceReload bool) {
	if txID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.DeleteHTTPRequestRule(int64(index), parentType, parentName, txID, version); err != nil {
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

func (h *HandlerImpl) replaceHTTPRequestRule(w http.ResponseWriter, r *http.Request, parentType, parentName string, index int, txID string, version int64, forceReload bool) {
	if txID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data HttpRequestRule
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.EditHTTPRequestRule(int64(index), parentType, parentName, &data, txID, version); err != nil {
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

func (h *HandlerImpl) replaceAllHTTPRequestRule(w http.ResponseWriter, r *http.Request, parentType, parentName, txID string, version int64, forceReload bool) {
	if txID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data HttpRequestRules
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.ReplaceHTTPRequestRules(parentType, parentName, data, txID, version); err != nil {
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
