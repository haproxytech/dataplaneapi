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

package http_error_rule

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	client_native "github.com/haproxytech/client-native/v6"

	"github.com/haproxytech/dataplaneapi/handlers/middleware"
	"github.com/haproxytech/dataplaneapi/handlers/respond"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/reload_agent"
)

// RegisterRouter registers all HTTP Error Rule routes onto r using spec-based request validation
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

// HandlerImpl implements ServerInterface for HAProxy HTTP Error Rule configuration.
type HandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent reload_agent.IReloadAgent
}

// --- Backend ---

func (h *HandlerImpl) GetAllHTTPErrorRuleBackend(w http.ResponseWriter, r *http.Request, parentName string, params GetAllHTTPErrorRuleBackendParams) {
	h.getAllHTTPErrorRule(w, r, "backend", parentName, params.TransactionId)
}

func (h *HandlerImpl) ReplaceAllHTTPErrorRuleBackend(w http.ResponseWriter, r *http.Request, parentName string, params ReplaceAllHTTPErrorRuleBackendParams) {
	h.replaceAllHTTPErrorRule(w, r, "backend", parentName, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) DeleteHTTPErrorRuleBackend(w http.ResponseWriter, r *http.Request, parentName string, index int, params DeleteHTTPErrorRuleBackendParams) {
	h.deleteHTTPErrorRule(w, r, "backend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) GetHTTPErrorRuleBackend(w http.ResponseWriter, r *http.Request, parentName string, index int, params GetHTTPErrorRuleBackendParams) {
	h.getHTTPErrorRule(w, r, "backend", parentName, index, params.TransactionId)
}

func (h *HandlerImpl) CreateHTTPErrorRuleBackend(w http.ResponseWriter, r *http.Request, parentName string, index int, params CreateHTTPErrorRuleBackendParams) {
	h.createHTTPErrorRule(w, r, "backend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) ReplaceHTTPErrorRuleBackend(w http.ResponseWriter, r *http.Request, parentName string, index int, params ReplaceHTTPErrorRuleBackendParams) {
	h.replaceHTTPErrorRule(w, r, "backend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

// --- Defaults ---

func (h *HandlerImpl) GetAllHTTPErrorRuleDefaults(w http.ResponseWriter, r *http.Request, parentName string, params GetAllHTTPErrorRuleDefaultsParams) {
	h.getAllHTTPErrorRule(w, r, "defaults", parentName, params.TransactionId)
}

func (h *HandlerImpl) ReplaceAllHTTPErrorRuleDefaults(w http.ResponseWriter, r *http.Request, parentName string, params ReplaceAllHTTPErrorRuleDefaultsParams) {
	h.replaceAllHTTPErrorRule(w, r, "defaults", parentName, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) DeleteHTTPErrorRuleDefaults(w http.ResponseWriter, r *http.Request, parentName string, index int, params DeleteHTTPErrorRuleDefaultsParams) {
	h.deleteHTTPErrorRule(w, r, "defaults", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) GetHTTPErrorRuleDefaults(w http.ResponseWriter, r *http.Request, parentName string, index int, params GetHTTPErrorRuleDefaultsParams) {
	h.getHTTPErrorRule(w, r, "defaults", parentName, index, params.TransactionId)
}

func (h *HandlerImpl) CreateHTTPErrorRuleDefaults(w http.ResponseWriter, r *http.Request, parentName string, index int, params CreateHTTPErrorRuleDefaultsParams) {
	h.createHTTPErrorRule(w, r, "defaults", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) ReplaceHTTPErrorRuleDefaults(w http.ResponseWriter, r *http.Request, parentName string, index int, params ReplaceHTTPErrorRuleDefaultsParams) {
	h.replaceHTTPErrorRule(w, r, "defaults", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

// --- Frontend ---

func (h *HandlerImpl) GetAllHTTPErrorRuleFrontend(w http.ResponseWriter, r *http.Request, parentName string, params GetAllHTTPErrorRuleFrontendParams) {
	h.getAllHTTPErrorRule(w, r, "frontend", parentName, params.TransactionId)
}

func (h *HandlerImpl) ReplaceAllHTTPErrorRuleFrontend(w http.ResponseWriter, r *http.Request, parentName string, params ReplaceAllHTTPErrorRuleFrontendParams) {
	h.replaceAllHTTPErrorRule(w, r, "frontend", parentName, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) DeleteHTTPErrorRuleFrontend(w http.ResponseWriter, r *http.Request, parentName string, index int, params DeleteHTTPErrorRuleFrontendParams) {
	h.deleteHTTPErrorRule(w, r, "frontend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) GetHTTPErrorRuleFrontend(w http.ResponseWriter, r *http.Request, parentName string, index int, params GetHTTPErrorRuleFrontendParams) {
	h.getHTTPErrorRule(w, r, "frontend", parentName, index, params.TransactionId)
}

func (h *HandlerImpl) CreateHTTPErrorRuleFrontend(w http.ResponseWriter, r *http.Request, parentName string, index int, params CreateHTTPErrorRuleFrontendParams) {
	h.createHTTPErrorRule(w, r, "frontend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) ReplaceHTTPErrorRuleFrontend(w http.ResponseWriter, r *http.Request, parentName string, index int, params ReplaceHTTPErrorRuleFrontendParams) {
	h.replaceHTTPErrorRule(w, r, "frontend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

// --- Shared implementations ---

func (h *HandlerImpl) getAllHTTPErrorRule(w http.ResponseWriter, r *http.Request, parentType, parentName, txID string) {
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	_, rules, err := cfg.GetHTTPErrorRules(parentType, parentName, txID)
	if err != nil {
		e := misc.HandleContainerGetError(err)
		if *e.Code == misc.ErrHTTPOk {
			respond.JSON(w, http.StatusOK, HttpErrorRules{})
			return
		}
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, rules)
}

func (h *HandlerImpl) getHTTPErrorRule(w http.ResponseWriter, r *http.Request, parentType, parentName string, index int, txID string) {
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	_, rule, err := cfg.GetHTTPErrorRule(int64(index), parentType, parentName, txID)
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, rule)
}

func (h *HandlerImpl) createHTTPErrorRule(w http.ResponseWriter, r *http.Request, parentType, parentName string, index int, txID string, version int64, forceReload bool) {
	if txID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data HttpErrorRule
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.CreateHTTPErrorRule(int64(index), parentType, parentName, &data, txID, version); err != nil {
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

func (h *HandlerImpl) deleteHTTPErrorRule(w http.ResponseWriter, r *http.Request, parentType, parentName string, index int, txID string, version int64, forceReload bool) {
	if txID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.DeleteHTTPErrorRule(int64(index), parentType, parentName, txID, version); err != nil {
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

func (h *HandlerImpl) replaceHTTPErrorRule(w http.ResponseWriter, r *http.Request, parentType, parentName string, index int, txID string, version int64, forceReload bool) {
	if txID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data HttpErrorRule
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.EditHTTPErrorRule(int64(index), parentType, parentName, &data, txID, version); err != nil {
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

func (h *HandlerImpl) replaceAllHTTPErrorRule(w http.ResponseWriter, r *http.Request, parentType, parentName, txID string, version int64, forceReload bool) {
	if txID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data HttpErrorRules
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.ReplaceHTTPErrorRules(parentType, parentName, data, txID, version); err != nil {
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
