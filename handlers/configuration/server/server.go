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

package server

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

// RegisterRouter registers all server routes onto r using spec-based request validation
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

// HandlerImpl implements ServerInterface for HAProxy server configuration.
type HandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent reload_agent.IReloadAgent
}

// --- Backend ---

func (h *HandlerImpl) GetAllServerBackend(w http.ResponseWriter, r *http.Request, parentName string, params GetAllServerBackendParams) {
	h.getAllServer(w, r, "backend", parentName, params.TransactionId)
}

func (h *HandlerImpl) CreateServerBackend(w http.ResponseWriter, r *http.Request, parentName string, params CreateServerBackendParams) {
	h.createServer(w, r, "backend", parentName, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) DeleteServerBackend(w http.ResponseWriter, r *http.Request, parentName string, name string, params DeleteServerBackendParams) {
	h.deleteServer(w, r, "backend", parentName, name, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) GetServerBackend(w http.ResponseWriter, r *http.Request, parentName string, name string, params GetServerBackendParams) {
	h.getServer(w, r, "backend", parentName, name, params.TransactionId)
}

func (h *HandlerImpl) ReplaceServerBackend(w http.ResponseWriter, r *http.Request, parentName string, name string, params ReplaceServerBackendParams) {
	h.replaceServer(w, r, "backend", parentName, name, params.TransactionId, int64(params.Version), params.ForceReload)
}

// --- Peer ---

func (h *HandlerImpl) GetAllServerPeer(w http.ResponseWriter, r *http.Request, parentName string, params GetAllServerPeerParams) {
	h.getAllServer(w, r, "peers", parentName, params.TransactionId)
}

func (h *HandlerImpl) CreateServerPeer(w http.ResponseWriter, r *http.Request, parentName string, params CreateServerPeerParams) {
	h.createServer(w, r, "peers", parentName, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) DeleteServerPeer(w http.ResponseWriter, r *http.Request, parentName string, name string, params DeleteServerPeerParams) {
	h.deleteServer(w, r, "peers", parentName, name, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) GetServerPeer(w http.ResponseWriter, r *http.Request, parentName string, name string, params GetServerPeerParams) {
	h.getServer(w, r, "peers", parentName, name, params.TransactionId)
}

func (h *HandlerImpl) ReplaceServerPeer(w http.ResponseWriter, r *http.Request, parentName string, name string, params ReplaceServerPeerParams) {
	h.replaceServer(w, r, "peers", parentName, name, params.TransactionId, int64(params.Version), params.ForceReload)
}

// --- Ring ---

func (h *HandlerImpl) GetAllServerRing(w http.ResponseWriter, r *http.Request, parentName string, params GetAllServerRingParams) {
	h.getAllServer(w, r, "ring", parentName, params.TransactionId)
}

func (h *HandlerImpl) CreateServerRing(w http.ResponseWriter, r *http.Request, parentName string, params CreateServerRingParams) {
	h.createServer(w, r, "ring", parentName, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) DeleteServerRing(w http.ResponseWriter, r *http.Request, parentName string, name string, params DeleteServerRingParams) {
	h.deleteServer(w, r, "ring", parentName, name, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) GetServerRing(w http.ResponseWriter, r *http.Request, parentName string, name string, params GetServerRingParams) {
	h.getServer(w, r, "ring", parentName, name, params.TransactionId)
}

func (h *HandlerImpl) ReplaceServerRing(w http.ResponseWriter, r *http.Request, parentName string, name string, params ReplaceServerRingParams) {
	h.replaceServer(w, r, "ring", parentName, name, params.TransactionId, int64(params.Version), params.ForceReload)
}

// --- Shared implementations ---

func (h *HandlerImpl) getAllServer(w http.ResponseWriter, r *http.Request, parentType, parentName, txID string) {
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	_, srvs, err := cfg.GetServers(parentType, parentName, txID)
	if err != nil {
		e := misc.HandleContainerGetError(err)
		if *e.Code == misc.ErrHTTPOk {
			respond.JSON(w, http.StatusOK, Servers{})
			return
		}
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, srvs)
}

func (h *HandlerImpl) getServer(w http.ResponseWriter, r *http.Request, parentType, parentName, name, txID string) {
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	_, srv, err := cfg.GetServer(name, parentType, parentName, txID)
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, srv)
}

func (h *HandlerImpl) createServer(w http.ResponseWriter, r *http.Request, parentType, parentName, txID string, version int64, forceReload bool) {
	if txID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data Server
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.CreateServer(parentType, parentName, &data, txID, version); err != nil {
		respond.Error(w, err)
		return
	}
	if txID == "" {
		respond.Committed(w, h.ReloadAgent, cn.NeedsReload(version), forceReload, http.StatusCreated, &data)
		return
	}
	respond.JSON(w, http.StatusAccepted, &data)
}

func (h *HandlerImpl) deleteServer(w http.ResponseWriter, r *http.Request, parentType, parentName, name, txID string, version int64, forceReload bool) {
	if txID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.DeleteServer(name, parentType, parentName, txID, version); err != nil {
		respond.Error(w, err)
		return
	}
	if txID == "" {
		respond.Committed(w, h.ReloadAgent, cn.NeedsReload(version), forceReload, http.StatusNoContent, nil)
		return
	}
	respond.Accepted(w, "", nil)
}

func (h *HandlerImpl) replaceServer(w http.ResponseWriter, r *http.Request, parentType, parentName, name, txID string, version int64, forceReload bool) {
	if txID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data Server
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.EditServer(name, parentType, parentName, &data, txID, version); err != nil {
		respond.Error(w, err)
		return
	}
	if txID == "" {
		respond.Committed(w, h.ReloadAgent, cn.NeedsReload(version), forceReload, http.StatusOK, &data)
		return
	}
	respond.JSON(w, http.StatusAccepted, &data)
}
