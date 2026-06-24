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

package defaults

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	client_native "github.com/haproxytech/client-native/v6"

	"github.com/haproxytech/dataplaneapi/handlers/middleware"
	"github.com/haproxytech/dataplaneapi/handlers/respond"
	"github.com/haproxytech/dataplaneapi/reload_agent"
)

// RegisterRouter registers all defaults routes onto r using spec-based request validation
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

// HandlerImpl implements ServerInterface for HAProxy defaults section configuration.
type HandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent reload_agent.IReloadAgent
}

func (h *HandlerImpl) GetDefaultsSections(w http.ResponseWriter, r *http.Request, params GetDefaultsSectionsParams) {
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	var sections DefaultsSections
	if params.FullSection {
		_, sections, err = cfg.GetStructuredDefaultsSections(params.TransactionId)
	} else {
		_, sections, err = cfg.GetDefaultsSections(params.TransactionId)
	}
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, sections)
}

func (h *HandlerImpl) AddDefaultsSection(w http.ResponseWriter, r *http.Request, params AddDefaultsSectionParams) {
	h.createDefaultsSection(w, r, params.TransactionId, params.Version, params.ForceReload, params.FullSection)
}

// CreateDefaultsSection handles the deprecated PUT alias of AddDefaultsSection.
func (h *HandlerImpl) CreateDefaultsSection(w http.ResponseWriter, r *http.Request, params CreateDefaultsSectionParams) {
	h.createDefaultsSection(w, r, params.TransactionId, params.Version, params.ForceReload, params.FullSection)
}

func (h *HandlerImpl) createDefaultsSection(w http.ResponseWriter, r *http.Request, transactionID string, version int, forceReload, fullSection bool) {
	if transactionID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data Defaults
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if fullSection {
		err = cfg.CreateStructuredDefaultsSection(&data, transactionID, int64(version))
	} else {
		err = cfg.CreateDefaultsSection(&data, transactionID, int64(version))
	}
	if err != nil {
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

func (h *HandlerImpl) DeleteDefaultsSection(w http.ResponseWriter, r *http.Request, name string, params DeleteDefaultsSectionParams) {
	if params.TransactionId != "" && params.ForceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.DeleteDefaultsSection(name, params.TransactionId, int64(params.Version)); err != nil {
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

func (h *HandlerImpl) GetDefaultsSection(w http.ResponseWriter, r *http.Request, name string, params GetDefaultsSectionParams) {
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	var section *Defaults
	if params.FullSection {
		_, section, err = cfg.GetStructuredDefaultsSection(name, params.TransactionId)
	} else {
		_, section, err = cfg.GetDefaultsSection(name, params.TransactionId)
	}
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, section)
}

func (h *HandlerImpl) ReplaceDefaultsSection(w http.ResponseWriter, r *http.Request, name string, params ReplaceDefaultsSectionParams) {
	if params.TransactionId != "" && params.ForceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data Defaults
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if params.FullSection {
		err = cfg.EditStructuredDefaultsSection(name, &data, params.TransactionId, int64(params.Version))
	} else {
		err = cfg.EditDefaultsSection(name, &data, params.TransactionId, int64(params.Version))
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
