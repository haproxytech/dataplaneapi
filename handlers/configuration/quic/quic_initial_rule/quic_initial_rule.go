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

package quic_initial_rule

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	client_native "github.com/haproxytech/client-native/v6"

	"github.com/haproxytech/dataplaneapi/handlers/middleware"
	"github.com/haproxytech/dataplaneapi/handlers/respond"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/reload_agent"
)

// RegisterRouter registers all QUIC initial rule routes onto r using spec-based request validation
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

// HandlerImpl implements ServerInterface for HAProxy QUIC initial rule configuration.
type HandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent reload_agent.IReloadAgent
}

// --- Defaults parent ---

func (h *HandlerImpl) GetAllQUICInitialRuleDefaults(w http.ResponseWriter, r *http.Request, parentName string, params GetAllQUICInitialRuleDefaultsParams) {
	h.getAllQUICInitialRule(w, r, "defaults", parentName, params.TransactionId)
}

func (h *HandlerImpl) ReplaceAllQUICInitialRuleDefaults(w http.ResponseWriter, r *http.Request, parentName string, params ReplaceAllQUICInitialRuleDefaultsParams) {
	h.replaceAllQUICInitialRule(w, r, "defaults", parentName, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) DeleteQUICInitialRuleDefaults(w http.ResponseWriter, r *http.Request, parentName string, index int, params DeleteQUICInitialRuleDefaultsParams) {
	h.deleteQUICInitialRule(w, r, "defaults", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) GetQUICInitialRuleDefaults(w http.ResponseWriter, r *http.Request, parentName string, index int, params GetQUICInitialRuleDefaultsParams) {
	h.getQUICInitialRule(w, r, "defaults", parentName, index, params.TransactionId)
}

func (h *HandlerImpl) CreateQUICInitialRuleDefaults(w http.ResponseWriter, r *http.Request, parentName string, index int, params CreateQUICInitialRuleDefaultsParams) {
	h.createQUICInitialRule(w, r, "defaults", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) ReplaceQUICInitialRuleDefaults(w http.ResponseWriter, r *http.Request, parentName string, index int, params ReplaceQUICInitialRuleDefaultsParams) {
	h.replaceQUICInitialRule(w, r, "defaults", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

// --- Frontend parent ---

func (h *HandlerImpl) GetAllQUICInitialRuleFrontend(w http.ResponseWriter, r *http.Request, parentName string, params GetAllQUICInitialRuleFrontendParams) {
	h.getAllQUICInitialRule(w, r, "frontend", parentName, params.TransactionId)
}

func (h *HandlerImpl) ReplaceAllQUICInitialRuleFrontend(w http.ResponseWriter, r *http.Request, parentName string, params ReplaceAllQUICInitialRuleFrontendParams) {
	h.replaceAllQUICInitialRule(w, r, "frontend", parentName, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) DeleteQUICInitialRuleFrontend(w http.ResponseWriter, r *http.Request, parentName string, index int, params DeleteQUICInitialRuleFrontendParams) {
	h.deleteQUICInitialRule(w, r, "frontend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) GetQUICInitialRuleFrontend(w http.ResponseWriter, r *http.Request, parentName string, index int, params GetQUICInitialRuleFrontendParams) {
	h.getQUICInitialRule(w, r, "frontend", parentName, index, params.TransactionId)
}

func (h *HandlerImpl) CreateQUICInitialRuleFrontend(w http.ResponseWriter, r *http.Request, parentName string, index int, params CreateQUICInitialRuleFrontendParams) {
	h.createQUICInitialRule(w, r, "frontend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) ReplaceQUICInitialRuleFrontend(w http.ResponseWriter, r *http.Request, parentName string, index int, params ReplaceQUICInitialRuleFrontendParams) {
	h.replaceQUICInitialRule(w, r, "frontend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

// --- Shared internal methods ---

func (h *HandlerImpl) getAllQUICInitialRule(w http.ResponseWriter, r *http.Request, parentType, parentName, transactionID string) {
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	_, rules, err := cfg.GetQUICInitialRules(parentType, parentName, transactionID)
	if err != nil {
		e := misc.HandleContainerGetError(err)
		if *e.Code == misc.ErrHTTPOk {
			respond.JSON(w, http.StatusOK, QuicInitialRules{})
			return
		}
		respond.JSON(w, int(*e.Code), e)
		return
	}
	respond.JSON(w, http.StatusOK, rules)
}

func (h *HandlerImpl) replaceAllQUICInitialRule(w http.ResponseWriter, r *http.Request, parentType, parentName, transactionID string, version int64, forceReload bool) {
	if transactionID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data QuicInitialRules
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.ReplaceQUICInitialRules(parentType, parentName, data, transactionID, version); err != nil {
		respond.Error(w, err)
		return
	}
	if transactionID == "" {
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

func (h *HandlerImpl) deleteQUICInitialRule(w http.ResponseWriter, r *http.Request, parentType, parentName string, index int, transactionID string, version int64, forceReload bool) {
	if transactionID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.DeleteQUICInitialRule(int64(index), parentType, parentName, transactionID, version); err != nil {
		respond.Error(w, err)
		return
	}
	if transactionID == "" {
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

func (h *HandlerImpl) getQUICInitialRule(w http.ResponseWriter, r *http.Request, parentType, parentName string, index int, transactionID string) {
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	_, rule, err := cfg.GetQUICInitialRule(int64(index), parentType, parentName, transactionID)
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, rule)
}

func (h *HandlerImpl) createQUICInitialRule(w http.ResponseWriter, r *http.Request, parentType, parentName string, index int, transactionID string, version int64, forceReload bool) {
	if transactionID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data QuicInitialRule
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.CreateQUICInitialRule(int64(index), parentType, parentName, &data, transactionID, version); err != nil {
		respond.Error(w, err)
		return
	}
	if transactionID == "" {
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

func (h *HandlerImpl) replaceQUICInitialRule(w http.ResponseWriter, r *http.Request, parentType, parentName string, index int, transactionID string, version int64, forceReload bool) {
	if transactionID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data QuicInitialRule
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.EditQUICInitialRule(int64(index), parentType, parentName, &data, transactionID, version); err != nil {
		respond.Error(w, err)
		return
	}
	if transactionID == "" {
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
