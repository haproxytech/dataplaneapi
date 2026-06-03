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

package tcp_check

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	client_native "github.com/haproxytech/client-native/v6"

	"github.com/haproxytech/dataplaneapi/handlers/middleware"
	"github.com/haproxytech/dataplaneapi/handlers/respond"
	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/misc"
)

// RegisterRouter registers all TCP check routes onto r using spec-based request validation
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

// HandlerImpl implements ServerInterface for HAProxy TCP check configuration.
type HandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// --- Backend parent ---

func (h *HandlerImpl) GetAllTCPCheckBackend(w http.ResponseWriter, r *http.Request, parentName string, params GetAllTCPCheckBackendParams) {
	h.getAllTCPCheck(w, r, "backend", parentName, params.TransactionId)
}

func (h *HandlerImpl) ReplaceAllTCPCheckBackend(w http.ResponseWriter, r *http.Request, parentName string, params ReplaceAllTCPCheckBackendParams) {
	h.replaceAllTCPCheck(w, r, "backend", parentName, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) DeleteTCPCheckBackend(w http.ResponseWriter, r *http.Request, parentName string, index int, params DeleteTCPCheckBackendParams) {
	h.deleteTCPCheck(w, r, "backend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) GetTCPCheckBackend(w http.ResponseWriter, r *http.Request, parentName string, index int, params GetTCPCheckBackendParams) {
	h.getTCPCheck(w, r, "backend", parentName, index, params.TransactionId)
}

func (h *HandlerImpl) CreateTCPCheckBackend(w http.ResponseWriter, r *http.Request, parentName string, index int, params CreateTCPCheckBackendParams) {
	h.createTCPCheck(w, r, "backend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) ReplaceTCPCheckBackend(w http.ResponseWriter, r *http.Request, parentName string, index int, params ReplaceTCPCheckBackendParams) {
	h.replaceTCPCheck(w, r, "backend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

// --- Defaults parent ---

func (h *HandlerImpl) GetAllTCPCheckDefaults(w http.ResponseWriter, r *http.Request, parentName string, params GetAllTCPCheckDefaultsParams) {
	h.getAllTCPCheck(w, r, "defaults", parentName, params.TransactionId)
}

func (h *HandlerImpl) ReplaceAllTCPCheckDefaults(w http.ResponseWriter, r *http.Request, parentName string, params ReplaceAllTCPCheckDefaultsParams) {
	h.replaceAllTCPCheck(w, r, "defaults", parentName, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) DeleteTCPCheckDefaults(w http.ResponseWriter, r *http.Request, parentName string, index int, params DeleteTCPCheckDefaultsParams) {
	h.deleteTCPCheck(w, r, "defaults", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) GetTCPCheckDefaults(w http.ResponseWriter, r *http.Request, parentName string, index int, params GetTCPCheckDefaultsParams) {
	h.getTCPCheck(w, r, "defaults", parentName, index, params.TransactionId)
}

func (h *HandlerImpl) CreateTCPCheckDefaults(w http.ResponseWriter, r *http.Request, parentName string, index int, params CreateTCPCheckDefaultsParams) {
	h.createTCPCheck(w, r, "defaults", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) ReplaceTCPCheckDefaults(w http.ResponseWriter, r *http.Request, parentName string, index int, params ReplaceTCPCheckDefaultsParams) {
	h.replaceTCPCheck(w, r, "defaults", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

// --- Shared internal methods ---

func (h *HandlerImpl) getAllTCPCheck(w http.ResponseWriter, r *http.Request, parentType, parentName, transactionID string) {
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	_, checks, err := cfg.GetTCPChecks(parentType, parentName, transactionID)
	if err != nil {
		e := misc.HandleContainerGetError(err)
		if *e.Code == misc.ErrHTTPOk {
			respond.JSON(w, http.StatusOK, TcpChecks{})
			return
		}
		respond.JSON(w, int(*e.Code), e)
		return
	}
	respond.JSON(w, http.StatusOK, checks)
}

func (h *HandlerImpl) replaceAllTCPCheck(w http.ResponseWriter, r *http.Request, parentType, parentName, transactionID string, version int64, forceReload bool) {
	if transactionID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data TcpChecks
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.ReplaceTCPChecks(parentType, parentName, data, transactionID, version); err != nil {
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

func (h *HandlerImpl) deleteTCPCheck(w http.ResponseWriter, r *http.Request, parentType, parentName string, index int, transactionID string, version int64, forceReload bool) {
	if transactionID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.DeleteTCPCheck(int64(index), parentType, parentName, transactionID, version); err != nil {
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

func (h *HandlerImpl) getTCPCheck(w http.ResponseWriter, r *http.Request, parentType, parentName string, index int, transactionID string) {
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	_, check, err := cfg.GetTCPCheck(int64(index), parentType, parentName, transactionID)
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, check)
}

func (h *HandlerImpl) createTCPCheck(w http.ResponseWriter, r *http.Request, parentType, parentName string, index int, transactionID string, version int64, forceReload bool) {
	if transactionID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data TcpCheck
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.CreateTCPCheck(int64(index), parentType, parentName, &data, transactionID, version); err != nil {
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

func (h *HandlerImpl) replaceTCPCheck(w http.ResponseWriter, r *http.Request, parentType, parentName string, index int, transactionID string, version int64, forceReload bool) {
	if transactionID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data TcpCheck
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.EditTCPCheck(int64(index), parentType, parentName, &data, transactionID, version); err != nil {
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
