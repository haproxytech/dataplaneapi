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

package global

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	client_native "github.com/haproxytech/client-native/v6"
	client_conf "github.com/haproxytech/client-native/v6/configuration"

	cn "github.com/haproxytech/dataplaneapi/client-native"
	"github.com/haproxytech/dataplaneapi/handlers/middleware"
	"github.com/haproxytech/dataplaneapi/handlers/respond"
	"github.com/haproxytech/dataplaneapi/reload_agent"
)

// RegisterRouter registers all global routes onto r using spec-based request validation
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

// HandlerImpl implements ServerInterface for HAProxy global configuration.
type HandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent reload_agent.IReloadAgent
}

func (h *HandlerImpl) GetGlobal(w http.ResponseWriter, r *http.Request, params GetGlobalParams) {
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	var data *Global
	if params.FullSection {
		_, data, err = cfg.GetStructuredGlobalConfiguration(params.TransactionId)
	} else {
		_, data, err = cfg.GetGlobalConfiguration(params.TransactionId)
	}
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, data)
}

func (h *HandlerImpl) ReplaceGlobal(w http.ResponseWriter, r *http.Request, params ReplaceGlobalParams) {
	if params.TransactionId != "" && params.ForceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data Global
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	if err := client_conf.ValidateGlobalSection(&data); err != nil {
		respond.BadRequest(w, err.Error())
		return
	}

	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if params.FullSection {
		err = cfg.PushStructuredGlobalConfiguration(&data, params.TransactionId, int64(params.Version))
	} else {
		err = cfg.PushGlobalConfiguration(&data, params.TransactionId, int64(params.Version))
	}
	if err != nil {
		respond.Error(w, err)
		return
	}

	if params.TransactionId == "" {
		callbackNeeded, reconfigureFunc, err := cn.ReconfigureRuntime(h.Client)
		if err != nil {
			respond.Error(w, err)
			return
		}
		if params.ForceReload {
			if callbackNeeded {
				err = h.ReloadAgent.ForceReloadWithCallback(reconfigureFunc)
			} else {
				err = h.ReloadAgent.ForceReload()
			}
			if err != nil {
				respond.Error(w, err)
				return
			}
			respond.JSON(w, http.StatusOK, &data)
			return
		}
		var rID string
		if callbackNeeded {
			rID = h.ReloadAgent.ReloadWithCallback(reconfigureFunc)
		} else {
			rID = h.ReloadAgent.Reload()
		}
		respond.Accepted(w, rID, &data)
		return
	}
	respond.JSON(w, http.StatusAccepted, &data)
}
