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

package tcp_response_rule

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	client_native "github.com/haproxytech/client-native/v6"

	"github.com/haproxytech/dataplaneapi/handlers/middleware"
	"github.com/haproxytech/dataplaneapi/handlers/respond"
	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/misc"
)

// RegisterRouter registers all TCP response rule routes onto r using spec-based request validation
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

// HandlerImpl implements ServerInterface for HAProxy TCP response rule configuration.
type HandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// --- Backend parent ---

func (h *HandlerImpl) GetAllTCPResponseRuleBackend(w http.ResponseWriter, r *http.Request, parentName string, params GetAllTCPResponseRuleBackendParams) {
	h.getAllTCPResponseRule(w, r, "backend", parentName, params.TransactionId)
}

func (h *HandlerImpl) ReplaceAllTCPResponseRuleBackend(w http.ResponseWriter, r *http.Request, parentName string, params ReplaceAllTCPResponseRuleBackendParams) {
	h.replaceAllTCPResponseRule(w, r, "backend", parentName, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) DeleteTCPResponseRuleBackend(w http.ResponseWriter, r *http.Request, parentName string, index int, params DeleteTCPResponseRuleBackendParams) {
	h.deleteTCPResponseRule(w, r, "backend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) GetTCPResponseRuleBackend(w http.ResponseWriter, r *http.Request, parentName string, index int, params GetTCPResponseRuleBackendParams) {
	h.getTCPResponseRule(w, r, "backend", parentName, index, params.TransactionId)
}

func (h *HandlerImpl) CreateTCPResponseRuleBackend(w http.ResponseWriter, r *http.Request, parentName string, index int, params CreateTCPResponseRuleBackendParams) {
	h.createTCPResponseRule(w, r, "backend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) ReplaceTCPResponseRuleBackend(w http.ResponseWriter, r *http.Request, parentName string, index int, params ReplaceTCPResponseRuleBackendParams) {
	h.replaceTCPResponseRule(w, r, "backend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

// --- Defaults parent ---

func (h *HandlerImpl) GetAllTCPResponseRuleDefaults(w http.ResponseWriter, r *http.Request, parentName string, params GetAllTCPResponseRuleDefaultsParams) {
	h.getAllTCPResponseRule(w, r, "defaults", parentName, params.TransactionId)
}

func (h *HandlerImpl) ReplaceAllTCPResponseRuleDefaults(w http.ResponseWriter, r *http.Request, parentName string, params ReplaceAllTCPResponseRuleDefaultsParams) {
	h.replaceAllTCPResponseRule(w, r, "defaults", parentName, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) DeleteTCPResponseRuleDefaults(w http.ResponseWriter, r *http.Request, parentName string, index int, params DeleteTCPResponseRuleDefaultsParams) {
	h.deleteTCPResponseRule(w, r, "defaults", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) GetTCPResponseRuleDefaults(w http.ResponseWriter, r *http.Request, parentName string, index int, params GetTCPResponseRuleDefaultsParams) {
	h.getTCPResponseRule(w, r, "defaults", parentName, index, params.TransactionId)
}

func (h *HandlerImpl) CreateTCPResponseRuleDefaults(w http.ResponseWriter, r *http.Request, parentName string, index int, params CreateTCPResponseRuleDefaultsParams) {
	h.createTCPResponseRule(w, r, "defaults", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) ReplaceTCPResponseRuleDefaults(w http.ResponseWriter, r *http.Request, parentName string, index int, params ReplaceTCPResponseRuleDefaultsParams) {
	h.replaceTCPResponseRule(w, r, "defaults", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

// --- Shared internal methods ---

func (h *HandlerImpl) getAllTCPResponseRule(w http.ResponseWriter, r *http.Request, parentType, parentName, transactionID string) {
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	_, rules, err := cfg.GetTCPResponseRules(parentType, parentName, transactionID)
	if err != nil {
		e := misc.HandleContainerGetError(err)
		if *e.Code == misc.ErrHTTPOk {
			respond.JSON(w, http.StatusOK, TcpResponseRules{})
			return
		}
		respond.JSON(w, int(*e.Code), e)
		return
	}
	respond.JSON(w, http.StatusOK, rules)
}

func (h *HandlerImpl) replaceAllTCPResponseRule(w http.ResponseWriter, r *http.Request, parentType, parentName, transactionID string, version int64, forceReload bool) {
	if transactionID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data TcpResponseRules
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.ReplaceTCPResponseRules(parentType, parentName, data, transactionID, version); err != nil {
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

func (h *HandlerImpl) deleteTCPResponseRule(w http.ResponseWriter, r *http.Request, parentType, parentName string, index int, transactionID string, version int64, forceReload bool) {
	if transactionID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.DeleteTCPResponseRule(int64(index), parentType, parentName, transactionID, version); err != nil {
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

func (h *HandlerImpl) getTCPResponseRule(w http.ResponseWriter, r *http.Request, parentType, parentName string, index int, transactionID string) {
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	_, rule, err := cfg.GetTCPResponseRule(int64(index), parentType, parentName, transactionID)
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, rule)
}

func (h *HandlerImpl) createTCPResponseRule(w http.ResponseWriter, r *http.Request, parentType, parentName string, index int, transactionID string, version int64, forceReload bool) {
	if transactionID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data TcpResponseRule
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.CreateTCPResponseRule(int64(index), parentType, parentName, &data, transactionID, version); err != nil {
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

func (h *HandlerImpl) replaceTCPResponseRule(w http.ResponseWriter, r *http.Request, parentType, parentName string, index int, transactionID string, version int64, forceReload bool) {
	if transactionID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data TcpResponseRule
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.EditTCPResponseRule(int64(index), parentType, parentName, &data, transactionID, version); err != nil {
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
