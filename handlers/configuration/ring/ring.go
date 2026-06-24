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

package ring

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	client_native "github.com/haproxytech/client-native/v6"

	"github.com/haproxytech/dataplaneapi/handlers/middleware"
	"github.com/haproxytech/dataplaneapi/handlers/respond"
	"github.com/haproxytech/dataplaneapi/reload_agent"
)

// RegisterRouter registers all ring routes onto r using spec-based request validation
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

// HandlerImpl implements ServerInterface for HAProxy ring configuration.
type HandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent reload_agent.IReloadAgent
}

func (h *HandlerImpl) GetRings(w http.ResponseWriter, r *http.Request, params GetRingsParams) {
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	var rings Rings
	if params.FullSection {
		_, rings, err = cfg.GetStructuredRings(params.TransactionId)
	} else {
		_, rings, err = cfg.GetRings(params.TransactionId)
	}
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, rings)
}

func (h *HandlerImpl) CreateRing(w http.ResponseWriter, r *http.Request, params CreateRingParams) {
	if params.TransactionId != "" && params.ForceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data Ring
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if params.FullSection {
		err = cfg.CreateStructuredRing(&data, params.TransactionId, int64(params.Version))
	} else {
		err = cfg.CreateRing(&data, params.TransactionId, int64(params.Version))
	}
	if err != nil {
		respond.Error(w, err)
		return
	}
	if params.TransactionId == "" {
		if params.ForceReload {
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

func (h *HandlerImpl) DeleteRing(w http.ResponseWriter, r *http.Request, name string, params DeleteRingParams) {
	if params.TransactionId != "" && params.ForceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.DeleteRing(name, params.TransactionId, int64(params.Version)); err != nil {
		respond.Error(w, err)
		return
	}
	if params.TransactionId == "" {
		if params.ForceReload {
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

func (h *HandlerImpl) GetRing(w http.ResponseWriter, r *http.Request, name string, params GetRingParams) {
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	var ring *Ring
	if params.FullSection {
		_, ring, err = cfg.GetStructuredRing(name, params.TransactionId)
	} else {
		_, ring, err = cfg.GetRing(name, params.TransactionId)
	}
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, ring)
}

func (h *HandlerImpl) ReplaceRing(w http.ResponseWriter, r *http.Request, name string, params ReplaceRingParams) {
	if params.TransactionId != "" && params.ForceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data Ring
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if params.FullSection {
		err = cfg.EditStructuredRing(name, &data, params.TransactionId, int64(params.Version))
	} else {
		err = cfg.EditRing(name, &data, params.TransactionId, int64(params.Version))
	}
	if err != nil {
		respond.Error(w, err)
		return
	}
	if params.TransactionId == "" {
		if params.ForceReload {
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
