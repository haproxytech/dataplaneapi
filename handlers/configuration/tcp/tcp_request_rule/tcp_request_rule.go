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

package tcp_request_rule

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	client_native "github.com/haproxytech/client-native/v6"

	"github.com/haproxytech/dataplaneapi/handlers/middleware"
	"github.com/haproxytech/dataplaneapi/handlers/respond"
	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/misc"
)

// RegisterRouter registers all TCP request rule routes onto r using spec-based request validation
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

// HandlerImpl implements ServerInterface for HAProxy TCP request rule configuration.
type HandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// --- Backend parent ---

func (h *HandlerImpl) GetAllTCPRequestRuleBackend(w http.ResponseWriter, r *http.Request, parentName string, params GetAllTCPRequestRuleBackendParams) {
	h.getAllTCPRequestRule(w, r, "backend", parentName, params.TransactionId)
}

func (h *HandlerImpl) ReplaceAllTCPRequestRuleBackend(w http.ResponseWriter, r *http.Request, parentName string, params ReplaceAllTCPRequestRuleBackendParams) {
	h.replaceAllTCPRequestRule(w, r, "backend", parentName, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) DeleteTCPRequestRuleBackend(w http.ResponseWriter, r *http.Request, parentName string, index int, params DeleteTCPRequestRuleBackendParams) {
	h.deleteTCPRequestRule(w, r, "backend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) GetTCPRequestRuleBackend(w http.ResponseWriter, r *http.Request, parentName string, index int, params GetTCPRequestRuleBackendParams) {
	h.getTCPRequestRule(w, r, "backend", parentName, index, params.TransactionId)
}

func (h *HandlerImpl) CreateTCPRequestRuleBackend(w http.ResponseWriter, r *http.Request, parentName string, index int, params CreateTCPRequestRuleBackendParams) {
	h.createTCPRequestRule(w, r, "backend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) ReplaceTCPRequestRuleBackend(w http.ResponseWriter, r *http.Request, parentName string, index int, params ReplaceTCPRequestRuleBackendParams) {
	h.replaceTCPRequestRule(w, r, "backend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

// --- Defaults parent ---

func (h *HandlerImpl) GetAllTCPRequestRuleDefaults(w http.ResponseWriter, r *http.Request, parentName string, params GetAllTCPRequestRuleDefaultsParams) {
	h.getAllTCPRequestRule(w, r, "defaults", parentName, params.TransactionId)
}

func (h *HandlerImpl) ReplaceAllTCPRequestRuleDefaults(w http.ResponseWriter, r *http.Request, parentName string, params ReplaceAllTCPRequestRuleDefaultsParams) {
	h.replaceAllTCPRequestRule(w, r, "defaults", parentName, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) DeleteTCPRequestRuleDefaults(w http.ResponseWriter, r *http.Request, parentName string, index int, params DeleteTCPRequestRuleDefaultsParams) {
	h.deleteTCPRequestRule(w, r, "defaults", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) GetTCPRequestRuleDefaults(w http.ResponseWriter, r *http.Request, parentName string, index int, params GetTCPRequestRuleDefaultsParams) {
	h.getTCPRequestRule(w, r, "defaults", parentName, index, params.TransactionId)
}

func (h *HandlerImpl) CreateTCPRequestRuleDefaults(w http.ResponseWriter, r *http.Request, parentName string, index int, params CreateTCPRequestRuleDefaultsParams) {
	h.createTCPRequestRule(w, r, "defaults", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) ReplaceTCPRequestRuleDefaults(w http.ResponseWriter, r *http.Request, parentName string, index int, params ReplaceTCPRequestRuleDefaultsParams) {
	h.replaceTCPRequestRule(w, r, "defaults", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

// --- Frontend parent ---

func (h *HandlerImpl) GetAllTCPRequestRuleFrontend(w http.ResponseWriter, r *http.Request, parentName string, params GetAllTCPRequestRuleFrontendParams) {
	h.getAllTCPRequestRule(w, r, "frontend", parentName, params.TransactionId)
}

func (h *HandlerImpl) ReplaceAllTCPRequestRuleFrontend(w http.ResponseWriter, r *http.Request, parentName string, params ReplaceAllTCPRequestRuleFrontendParams) {
	h.replaceAllTCPRequestRule(w, r, "frontend", parentName, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) DeleteTCPRequestRuleFrontend(w http.ResponseWriter, r *http.Request, parentName string, index int, params DeleteTCPRequestRuleFrontendParams) {
	h.deleteTCPRequestRule(w, r, "frontend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) GetTCPRequestRuleFrontend(w http.ResponseWriter, r *http.Request, parentName string, index int, params GetTCPRequestRuleFrontendParams) {
	h.getTCPRequestRule(w, r, "frontend", parentName, index, params.TransactionId)
}

func (h *HandlerImpl) CreateTCPRequestRuleFrontend(w http.ResponseWriter, r *http.Request, parentName string, index int, params CreateTCPRequestRuleFrontendParams) {
	h.createTCPRequestRule(w, r, "frontend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) ReplaceTCPRequestRuleFrontend(w http.ResponseWriter, r *http.Request, parentName string, index int, params ReplaceTCPRequestRuleFrontendParams) {
	h.replaceTCPRequestRule(w, r, "frontend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

// --- Shared internal methods ---

func (h *HandlerImpl) getAllTCPRequestRule(w http.ResponseWriter, r *http.Request, parentType, parentName, transactionID string) {
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	_, rules, err := cfg.GetTCPRequestRules(parentType, parentName, transactionID)
	if err != nil {
		e := misc.HandleContainerGetError(err)
		if *e.Code == misc.ErrHTTPOk {
			respond.JSON(w, http.StatusOK, TcpRequestRules{})
			return
		}
		respond.JSON(w, int(*e.Code), e)
		return
	}
	respond.JSON(w, http.StatusOK, rules)
}

func (h *HandlerImpl) replaceAllTCPRequestRule(w http.ResponseWriter, r *http.Request, parentType, parentName, transactionID string, version int64, forceReload bool) {
	if transactionID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data TcpRequestRules
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.ReplaceTCPRequestRules(parentType, parentName, data, transactionID, version); err != nil {
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

func (h *HandlerImpl) deleteTCPRequestRule(w http.ResponseWriter, r *http.Request, parentType, parentName string, index int, transactionID string, version int64, forceReload bool) {
	if transactionID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.DeleteTCPRequestRule(int64(index), parentType, parentName, transactionID, version); err != nil {
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

func (h *HandlerImpl) getTCPRequestRule(w http.ResponseWriter, r *http.Request, parentType, parentName string, index int, transactionID string) {
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	_, rule, err := cfg.GetTCPRequestRule(int64(index), parentType, parentName, transactionID)
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, rule)
}

func (h *HandlerImpl) createTCPRequestRule(w http.ResponseWriter, r *http.Request, parentType, parentName string, index int, transactionID string, version int64, forceReload bool) {
	if transactionID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data TcpRequestRule
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.CreateTCPRequestRule(int64(index), parentType, parentName, &data, transactionID, version); err != nil {
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

func (h *HandlerImpl) replaceTCPRequestRule(w http.ResponseWriter, r *http.Request, parentType, parentName string, index int, transactionID string, version int64, forceReload bool) {
	if transactionID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data TcpRequestRule
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.EditTCPRequestRule(int64(index), parentType, parentName, &data, transactionID, version); err != nil {
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
