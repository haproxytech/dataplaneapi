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

package http_response_rule

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	client_native "github.com/haproxytech/client-native/v6"

	"github.com/haproxytech/dataplaneapi/handlers/middleware"
	"github.com/haproxytech/dataplaneapi/handlers/respond"
	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/misc"
)

// RegisterRouter registers all HTTP Response Rule routes onto r using spec-based request validation
// and a shared error handler.
func RegisterRouter(r chi.Router, client client_native.HAProxyClient, ra haproxy.IReloadAgent) error {
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

// HandlerImpl implements ServerInterface for HAProxy HTTP Response Rule configuration.
type HandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// --- Backend ---

func (h *HandlerImpl) GetAllHTTPResponseRuleBackend(w http.ResponseWriter, r *http.Request, parentName string, params GetAllHTTPResponseRuleBackendParams) {
	h.getAllHTTPResponseRule(w, r, "backend", parentName, params.TransactionId)
}

func (h *HandlerImpl) ReplaceAllHTTPResponseRuleBackend(w http.ResponseWriter, r *http.Request, parentName string, params ReplaceAllHTTPResponseRuleBackendParams) {
	h.replaceAllHTTPResponseRule(w, r, "backend", parentName, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) DeleteHTTPResponseRuleBackend(w http.ResponseWriter, r *http.Request, parentName string, index int, params DeleteHTTPResponseRuleBackendParams) {
	h.deleteHTTPResponseRule(w, r, "backend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) GetHTTPResponseRuleBackend(w http.ResponseWriter, r *http.Request, parentName string, index int, params GetHTTPResponseRuleBackendParams) {
	h.getHTTPResponseRule(w, r, "backend", parentName, index, params.TransactionId)
}

func (h *HandlerImpl) CreateHTTPResponseRuleBackend(w http.ResponseWriter, r *http.Request, parentName string, index int, params CreateHTTPResponseRuleBackendParams) {
	h.createHTTPResponseRule(w, r, "backend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) ReplaceHTTPResponseRuleBackend(w http.ResponseWriter, r *http.Request, parentName string, index int, params ReplaceHTTPResponseRuleBackendParams) {
	h.replaceHTTPResponseRule(w, r, "backend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

// --- Defaults ---

func (h *HandlerImpl) GetAllHTTPResponseRuleDefaults(w http.ResponseWriter, r *http.Request, parentName string, params GetAllHTTPResponseRuleDefaultsParams) {
	h.getAllHTTPResponseRule(w, r, "defaults", parentName, params.TransactionId)
}

func (h *HandlerImpl) ReplaceAllHTTPResponseRuleDefaults(w http.ResponseWriter, r *http.Request, parentName string, params ReplaceAllHTTPResponseRuleDefaultsParams) {
	h.replaceAllHTTPResponseRule(w, r, "defaults", parentName, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) DeleteHTTPResponseRuleDefaults(w http.ResponseWriter, r *http.Request, parentName string, index int, params DeleteHTTPResponseRuleDefaultsParams) {
	h.deleteHTTPResponseRule(w, r, "defaults", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) GetHTTPResponseRuleDefaults(w http.ResponseWriter, r *http.Request, parentName string, index int, params GetHTTPResponseRuleDefaultsParams) {
	h.getHTTPResponseRule(w, r, "defaults", parentName, index, params.TransactionId)
}

func (h *HandlerImpl) CreateHTTPResponseRuleDefaults(w http.ResponseWriter, r *http.Request, parentName string, index int, params CreateHTTPResponseRuleDefaultsParams) {
	h.createHTTPResponseRule(w, r, "defaults", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) ReplaceHTTPResponseRuleDefaults(w http.ResponseWriter, r *http.Request, parentName string, index int, params ReplaceHTTPResponseRuleDefaultsParams) {
	h.replaceHTTPResponseRule(w, r, "defaults", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

// --- Frontend ---

func (h *HandlerImpl) GetAllHTTPResponseRuleFrontend(w http.ResponseWriter, r *http.Request, parentName string, params GetAllHTTPResponseRuleFrontendParams) {
	h.getAllHTTPResponseRule(w, r, "frontend", parentName, params.TransactionId)
}

func (h *HandlerImpl) ReplaceAllHTTPResponseRuleFrontend(w http.ResponseWriter, r *http.Request, parentName string, params ReplaceAllHTTPResponseRuleFrontendParams) {
	h.replaceAllHTTPResponseRule(w, r, "frontend", parentName, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) DeleteHTTPResponseRuleFrontend(w http.ResponseWriter, r *http.Request, parentName string, index int, params DeleteHTTPResponseRuleFrontendParams) {
	h.deleteHTTPResponseRule(w, r, "frontend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) GetHTTPResponseRuleFrontend(w http.ResponseWriter, r *http.Request, parentName string, index int, params GetHTTPResponseRuleFrontendParams) {
	h.getHTTPResponseRule(w, r, "frontend", parentName, index, params.TransactionId)
}

func (h *HandlerImpl) CreateHTTPResponseRuleFrontend(w http.ResponseWriter, r *http.Request, parentName string, index int, params CreateHTTPResponseRuleFrontendParams) {
	h.createHTTPResponseRule(w, r, "frontend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) ReplaceHTTPResponseRuleFrontend(w http.ResponseWriter, r *http.Request, parentName string, index int, params ReplaceHTTPResponseRuleFrontendParams) {
	h.replaceHTTPResponseRule(w, r, "frontend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

// --- Shared implementations ---

func (h *HandlerImpl) getAllHTTPResponseRule(w http.ResponseWriter, r *http.Request, parentType, parentName, txID string) {
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	_, rules, err := cfg.GetHTTPResponseRules(parentType, parentName, txID)
	if err != nil {
		e := misc.HandleContainerGetError(err)
		if *e.Code == misc.ErrHTTPOk {
			respond.JSON(w, http.StatusOK, HttpResponseRules{})
			return
		}
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, rules)
}

func (h *HandlerImpl) getHTTPResponseRule(w http.ResponseWriter, r *http.Request, parentType, parentName string, index int, txID string) {
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	_, rule, err := cfg.GetHTTPResponseRule(int64(index), parentType, parentName, txID)
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, rule)
}

func (h *HandlerImpl) createHTTPResponseRule(w http.ResponseWriter, r *http.Request, parentType, parentName string, index int, txID string, version int64, forceReload bool) {
	if txID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data HttpResponseRule
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.CreateHTTPResponseRule(int64(index), parentType, parentName, &data, txID, version); err != nil {
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

func (h *HandlerImpl) deleteHTTPResponseRule(w http.ResponseWriter, r *http.Request, parentType, parentName string, index int, txID string, version int64, forceReload bool) {
	if txID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.DeleteHTTPResponseRule(int64(index), parentType, parentName, txID, version); err != nil {
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

func (h *HandlerImpl) replaceHTTPResponseRule(w http.ResponseWriter, r *http.Request, parentType, parentName string, index int, txID string, version int64, forceReload bool) {
	if txID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data HttpResponseRule
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.EditHTTPResponseRule(int64(index), parentType, parentName, &data, txID, version); err != nil {
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

func (h *HandlerImpl) replaceAllHTTPResponseRule(w http.ResponseWriter, r *http.Request, parentType, parentName, txID string, version int64, forceReload bool) {
	if txID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data HttpResponseRules
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.ReplaceHTTPResponseRules(parentType, parentName, data, txID, version); err != nil {
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
