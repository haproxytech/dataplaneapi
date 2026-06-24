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

package log_target

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	client_native "github.com/haproxytech/client-native/v6"
	"github.com/haproxytech/client-native/v6/models"

	"github.com/haproxytech/dataplaneapi/handlers/middleware"
	"github.com/haproxytech/dataplaneapi/handlers/respond"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/reload_agent"
)

// RegisterRouter registers all log_target routes onto r using spec-based request validation
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

// HandlerImpl implements ServerInterface for HAProxy log_target configuration.
type HandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent reload_agent.IReloadAgent
}

// --- Backend ---

func (h *HandlerImpl) GetAllLogTargetBackend(w http.ResponseWriter, r *http.Request, parentName string, params GetAllLogTargetBackendParams) {
	h.getAllLogTargets(w, r, parentName, "backend", params.TransactionId)
}

func (h *HandlerImpl) ReplaceAllLogTargetBackend(w http.ResponseWriter, r *http.Request, parentName string, params ReplaceAllLogTargetBackendParams) {
	h.replaceAllLogTargets(w, r, parentName, "backend", params.TransactionId, params.Version, params.ForceReload)
}

func (h *HandlerImpl) DeleteLogTargetBackend(w http.ResponseWriter, r *http.Request, parentName string, index int, params DeleteLogTargetBackendParams) {
	h.deleteLogTarget(w, r, parentName, "backend", index, params.TransactionId, params.Version, params.ForceReload)
}

func (h *HandlerImpl) GetLogTargetBackend(w http.ResponseWriter, r *http.Request, parentName string, index int, params GetLogTargetBackendParams) {
	h.getLogTarget(w, r, parentName, "backend", index, params.TransactionId)
}

func (h *HandlerImpl) CreateLogTargetBackend(w http.ResponseWriter, r *http.Request, parentName string, index int, params CreateLogTargetBackendParams) {
	h.createLogTarget(w, r, parentName, "backend", index, params.TransactionId, params.Version, params.ForceReload)
}

func (h *HandlerImpl) ReplaceLogTargetBackend(w http.ResponseWriter, r *http.Request, parentName string, index int, params ReplaceLogTargetBackendParams) {
	h.replaceLogTarget(w, r, parentName, "backend", index, params.TransactionId, params.Version, params.ForceReload)
}

// --- Defaults ---

func (h *HandlerImpl) GetAllLogTargetDefaults(w http.ResponseWriter, r *http.Request, parentName string, params GetAllLogTargetDefaultsParams) {
	h.getAllLogTargets(w, r, parentName, "defaults", params.TransactionId)
}

func (h *HandlerImpl) ReplaceAllLogTargetDefaults(w http.ResponseWriter, r *http.Request, parentName string, params ReplaceAllLogTargetDefaultsParams) {
	h.replaceAllLogTargets(w, r, parentName, "defaults", params.TransactionId, params.Version, params.ForceReload)
}

func (h *HandlerImpl) DeleteLogTargetDefaults(w http.ResponseWriter, r *http.Request, parentName string, index int, params DeleteLogTargetDefaultsParams) {
	h.deleteLogTarget(w, r, parentName, "defaults", index, params.TransactionId, params.Version, params.ForceReload)
}

func (h *HandlerImpl) GetLogTargetDefaults(w http.ResponseWriter, r *http.Request, parentName string, index int, params GetLogTargetDefaultsParams) {
	h.getLogTarget(w, r, parentName, "defaults", index, params.TransactionId)
}

func (h *HandlerImpl) CreateLogTargetDefaults(w http.ResponseWriter, r *http.Request, parentName string, index int, params CreateLogTargetDefaultsParams) {
	h.createLogTarget(w, r, parentName, "defaults", index, params.TransactionId, params.Version, params.ForceReload)
}

func (h *HandlerImpl) ReplaceLogTargetDefaults(w http.ResponseWriter, r *http.Request, parentName string, index int, params ReplaceLogTargetDefaultsParams) {
	h.replaceLogTarget(w, r, parentName, "defaults", index, params.TransactionId, params.Version, params.ForceReload)
}

// --- Frontend ---

func (h *HandlerImpl) GetAllLogTargetFrontend(w http.ResponseWriter, r *http.Request, parentName string, params GetAllLogTargetFrontendParams) {
	h.getAllLogTargets(w, r, parentName, "frontend", params.TransactionId)
}

func (h *HandlerImpl) ReplaceAllLogTargetFrontend(w http.ResponseWriter, r *http.Request, parentName string, params ReplaceAllLogTargetFrontendParams) {
	h.replaceAllLogTargets(w, r, parentName, "frontend", params.TransactionId, params.Version, params.ForceReload)
}

func (h *HandlerImpl) DeleteLogTargetFrontend(w http.ResponseWriter, r *http.Request, parentName string, index int, params DeleteLogTargetFrontendParams) {
	h.deleteLogTarget(w, r, parentName, "frontend", index, params.TransactionId, params.Version, params.ForceReload)
}

func (h *HandlerImpl) GetLogTargetFrontend(w http.ResponseWriter, r *http.Request, parentName string, index int, params GetLogTargetFrontendParams) {
	h.getLogTarget(w, r, parentName, "frontend", index, params.TransactionId)
}

func (h *HandlerImpl) CreateLogTargetFrontend(w http.ResponseWriter, r *http.Request, parentName string, index int, params CreateLogTargetFrontendParams) {
	h.createLogTarget(w, r, parentName, "frontend", index, params.TransactionId, params.Version, params.ForceReload)
}

func (h *HandlerImpl) ReplaceLogTargetFrontend(w http.ResponseWriter, r *http.Request, parentName string, index int, params ReplaceLogTargetFrontendParams) {
	h.replaceLogTarget(w, r, parentName, "frontend", index, params.TransactionId, params.Version, params.ForceReload)
}

// --- LogForward ---

func (h *HandlerImpl) GetAllLogTargetLogForward(w http.ResponseWriter, r *http.Request, parentName string, params GetAllLogTargetLogForwardParams) {
	h.getAllLogTargets(w, r, parentName, "log_forward", params.TransactionId)
}

func (h *HandlerImpl) ReplaceAllLogTargetLogForward(w http.ResponseWriter, r *http.Request, parentName string, params ReplaceAllLogTargetLogForwardParams) {
	h.replaceAllLogTargets(w, r, parentName, "log_forward", params.TransactionId, params.Version, params.ForceReload)
}

func (h *HandlerImpl) DeleteLogTargetLogForward(w http.ResponseWriter, r *http.Request, parentName string, index int, params DeleteLogTargetLogForwardParams) {
	h.deleteLogTarget(w, r, parentName, "log_forward", index, params.TransactionId, params.Version, params.ForceReload)
}

func (h *HandlerImpl) GetLogTargetLogForward(w http.ResponseWriter, r *http.Request, parentName string, index int, params GetLogTargetLogForwardParams) {
	h.getLogTarget(w, r, parentName, "log_forward", index, params.TransactionId)
}

func (h *HandlerImpl) CreateLogTargetLogForward(w http.ResponseWriter, r *http.Request, parentName string, index int, params CreateLogTargetLogForwardParams) {
	h.createLogTarget(w, r, parentName, "log_forward", index, params.TransactionId, params.Version, params.ForceReload)
}

func (h *HandlerImpl) ReplaceLogTargetLogForward(w http.ResponseWriter, r *http.Request, parentName string, index int, params ReplaceLogTargetLogForwardParams) {
	h.replaceLogTarget(w, r, parentName, "log_forward", index, params.TransactionId, params.Version, params.ForceReload)
}

// --- Peer ---

func (h *HandlerImpl) GetAllLogTargetPeer(w http.ResponseWriter, r *http.Request, parentName string, params GetAllLogTargetPeerParams) {
	h.getAllLogTargets(w, r, parentName, "peers", params.TransactionId)
}

func (h *HandlerImpl) ReplaceAllLogTargetPeer(w http.ResponseWriter, r *http.Request, parentName string, params ReplaceAllLogTargetPeerParams) {
	h.replaceAllLogTargets(w, r, parentName, "peers", params.TransactionId, params.Version, params.ForceReload)
}

func (h *HandlerImpl) DeleteLogTargetPeer(w http.ResponseWriter, r *http.Request, parentName string, index int, params DeleteLogTargetPeerParams) {
	h.deleteLogTarget(w, r, parentName, "peers", index, params.TransactionId, params.Version, params.ForceReload)
}

func (h *HandlerImpl) GetLogTargetPeer(w http.ResponseWriter, r *http.Request, parentName string, index int, params GetLogTargetPeerParams) {
	h.getLogTarget(w, r, parentName, "peers", index, params.TransactionId)
}

func (h *HandlerImpl) CreateLogTargetPeer(w http.ResponseWriter, r *http.Request, parentName string, index int, params CreateLogTargetPeerParams) {
	h.createLogTarget(w, r, parentName, "peers", index, params.TransactionId, params.Version, params.ForceReload)
}

func (h *HandlerImpl) ReplaceLogTargetPeer(w http.ResponseWriter, r *http.Request, parentName string, index int, params ReplaceLogTargetPeerParams) {
	h.replaceLogTarget(w, r, parentName, "peers", index, params.TransactionId, params.Version, params.ForceReload)
}

// --- Shared internal implementations ---

func (h *HandlerImpl) getAllLogTargets(w http.ResponseWriter, r *http.Request, parentName, parentType, transactionID string) {
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	_, logTargets, err := cfg.GetLogTargets(parentType, parentName, transactionID)
	if err != nil {
		e := misc.HandleContainerGetError(err)
		if *e.Code == misc.ErrHTTPOk {
			respond.JSON(w, http.StatusOK, models.LogTargets{})
			return
		}
		respond.JSON(w, int(*e.Code), e)
		return
	}
	respond.JSON(w, http.StatusOK, logTargets)
}

func (h *HandlerImpl) replaceAllLogTargets(w http.ResponseWriter, r *http.Request, parentName, parentType, transactionID string, version int, forceReload bool) {
	if transactionID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data LogTargets
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.ReplaceLogTargets(parentType, parentName, data, transactionID, int64(version)); err != nil {
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

func (h *HandlerImpl) deleteLogTarget(w http.ResponseWriter, r *http.Request, parentName, parentType string, index int, transactionID string, version int, forceReload bool) {
	if transactionID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.DeleteLogTarget(int64(index), parentType, parentName, transactionID, int64(version)); err != nil {
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

func (h *HandlerImpl) getLogTarget(w http.ResponseWriter, r *http.Request, parentName, parentType string, index int, transactionID string) {
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	_, logTarget, err := cfg.GetLogTarget(int64(index), parentType, parentName, transactionID)
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, logTarget)
}

func (h *HandlerImpl) createLogTarget(w http.ResponseWriter, r *http.Request, parentName, parentType string, index int, transactionID string, version int, forceReload bool) {
	if transactionID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data LogTarget
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.CreateLogTarget(int64(index), parentType, parentName, &data, transactionID, int64(version)); err != nil {
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

func (h *HandlerImpl) replaceLogTarget(w http.ResponseWriter, r *http.Request, parentName, parentType string, index int, transactionID string, version int, forceReload bool) {
	if transactionID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data LogTarget
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.EditLogTarget(int64(index), parentType, parentName, &data, transactionID, int64(version)); err != nil {
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
