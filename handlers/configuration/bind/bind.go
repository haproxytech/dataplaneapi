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

package bind

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	client_native "github.com/haproxytech/client-native/v6"

	"github.com/haproxytech/dataplaneapi/handlers/middleware"
	"github.com/haproxytech/dataplaneapi/handlers/respond"
	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/misc"
)

// RegisterRouter registers all bind routes onto r using spec-based request validation
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

// HandlerImpl implements ServerInterface for HAProxy bind configuration.
type HandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// Frontend parent variants

func (h *HandlerImpl) GetAllBindFrontend(w http.ResponseWriter, r *http.Request, parentName string, params GetAllBindFrontendParams) {
	h.getAllBind(w, r, "frontend", parentName, params.TransactionId)
}

func (h *HandlerImpl) CreateBindFrontend(w http.ResponseWriter, r *http.Request, parentName string, params CreateBindFrontendParams) {
	h.createBind(w, r, "frontend", parentName, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) DeleteBindFrontend(w http.ResponseWriter, r *http.Request, parentName string, name string, params DeleteBindFrontendParams) {
	h.deleteBind(w, r, "frontend", parentName, name, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) GetBindFrontend(w http.ResponseWriter, r *http.Request, parentName string, name string, params GetBindFrontendParams) {
	h.getBind(w, r, "frontend", parentName, name, params.TransactionId)
}

func (h *HandlerImpl) ReplaceBindFrontend(w http.ResponseWriter, r *http.Request, parentName string, name string, params ReplaceBindFrontendParams) {
	h.replaceBind(w, r, "frontend", parentName, name, params.TransactionId, int64(params.Version), params.ForceReload)
}

// LogForward parent variants

func (h *HandlerImpl) GetAllBindLogForward(w http.ResponseWriter, r *http.Request, parentName string, params GetAllBindLogForwardParams) {
	h.getAllBind(w, r, "log_forward", parentName, params.TransactionId)
}

func (h *HandlerImpl) CreateBindLogForward(w http.ResponseWriter, r *http.Request, parentName string, params CreateBindLogForwardParams) {
	h.createBind(w, r, "log_forward", parentName, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) DeleteBindLogForward(w http.ResponseWriter, r *http.Request, parentName string, name string, params DeleteBindLogForwardParams) {
	h.deleteBind(w, r, "log_forward", parentName, name, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) GetBindLogForward(w http.ResponseWriter, r *http.Request, parentName string, name string, params GetBindLogForwardParams) {
	h.getBind(w, r, "log_forward", parentName, name, params.TransactionId)
}

func (h *HandlerImpl) ReplaceBindLogForward(w http.ResponseWriter, r *http.Request, parentName string, name string, params ReplaceBindLogForwardParams) {
	h.replaceBind(w, r, "log_forward", parentName, name, params.TransactionId, int64(params.Version), params.ForceReload)
}

// Peer parent variants

func (h *HandlerImpl) GetAllBindPeer(w http.ResponseWriter, r *http.Request, parentName string, params GetAllBindPeerParams) {
	h.getAllBind(w, r, "peers", parentName, params.TransactionId)
}

func (h *HandlerImpl) CreateBindPeer(w http.ResponseWriter, r *http.Request, parentName string, params CreateBindPeerParams) {
	h.createBind(w, r, "peers", parentName, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) DeleteBindPeer(w http.ResponseWriter, r *http.Request, parentName string, name string, params DeleteBindPeerParams) {
	h.deleteBind(w, r, "peers", parentName, name, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) GetBindPeer(w http.ResponseWriter, r *http.Request, parentName string, name string, params GetBindPeerParams) {
	h.getBind(w, r, "peers", parentName, name, params.TransactionId)
}

func (h *HandlerImpl) ReplaceBindPeer(w http.ResponseWriter, r *http.Request, parentName string, name string, params ReplaceBindPeerParams) {
	h.replaceBind(w, r, "peers", parentName, name, params.TransactionId, int64(params.Version), params.ForceReload)
}

// Shared internal methods

func (h *HandlerImpl) getAllBind(w http.ResponseWriter, r *http.Request, parentType, parentName, txID string) {
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	_, bs, err := cfg.GetBinds(parentType, parentName, txID)
	if err != nil {
		e := misc.HandleContainerGetError(err)
		if *e.Code == misc.ErrHTTPOk {
			respond.JSON(w, http.StatusOK, Binds{})
			return
		}
		respond.JSON(w, int(*e.Code), e)
		return
	}
	respond.JSON(w, http.StatusOK, bs)
}

func (h *HandlerImpl) createBind(w http.ResponseWriter, r *http.Request, parentType, parentName, txID string, version int64, forceReload bool) {
	if txID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data Bind
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.CreateBind(parentType, parentName, &data, txID, version); err != nil {
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

func (h *HandlerImpl) deleteBind(w http.ResponseWriter, r *http.Request, parentType, parentName, name, txID string, version int64, forceReload bool) {
	if txID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.DeleteBind(name, parentType, parentName, txID, version); err != nil {
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

func (h *HandlerImpl) getBind(w http.ResponseWriter, r *http.Request, parentType, parentName, name, txID string) {
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	_, b, err := cfg.GetBind(name, parentType, parentName, txID)
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, b)
}

func (h *HandlerImpl) replaceBind(w http.ResponseWriter, r *http.Request, parentType, parentName, name, txID string, version int64, forceReload bool) {
	if txID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data Bind
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.EditBind(name, parentType, parentName, &data, txID, version); err != nil {
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
