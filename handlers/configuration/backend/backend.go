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

package backend

import (
	"net/http"
	"slices"

	"github.com/go-chi/chi/v5"
	client_native "github.com/haproxytech/client-native/v6"
	"github.com/haproxytech/client-native/v6/models"

	"github.com/haproxytech/dataplaneapi/handlers/middleware"
	"github.com/haproxytech/dataplaneapi/handlers/respond"
	"github.com/haproxytech/dataplaneapi/haproxy"
)

// RegisterRouter registers all backend routes onto r using spec-based request validation
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

// HandlerImpl implements ServerInterface for HAProxy backend configuration.
type HandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

func (h *HandlerImpl) GetBackends(w http.ResponseWriter, r *http.Request, params GetBackendsParams) {
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	var bcks models.Backends
	if params.FullSection {
		_, bcks, err = cfg.GetStructuredBackends(params.TransactionId)
	} else {
		_, bcks, err = cfg.GetBackends(params.TransactionId)
	}
	if err != nil {
		respond.Error(w, err)
		return
	}
	for _, bck := range bcks {
		handleDeprecatedBackendFields(http.MethodGet, bck, nil)
	}
	respond.JSON(w, http.StatusOK, bcks)
}

func (h *HandlerImpl) CreateBackend(w http.ResponseWriter, r *http.Request, params CreateBackendParams) {
	if params.TransactionId != "" && params.ForceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data Backend
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	handleDeprecatedBackendFields(http.MethodPost, &data, nil)

	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if params.FullSection {
		err = cfg.CreateStructuredBackend(&data, params.TransactionId, int64(params.Version))
	} else {
		err = cfg.CreateBackend(&data, params.TransactionId, int64(params.Version))
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

func (h *HandlerImpl) DeleteBackend(w http.ResponseWriter, r *http.Request, name string, params DeleteBackendParams) {
	if params.TransactionId != "" && params.ForceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.DeleteBackend(name, params.TransactionId, int64(params.Version)); err != nil {
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

func (h *HandlerImpl) GetBackend(w http.ResponseWriter, r *http.Request, name string, params GetBackendParams) {
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	var bck *models.Backend
	if params.FullSection {
		_, bck, err = cfg.GetStructuredBackend(name, params.TransactionId)
	} else {
		_, bck, err = cfg.GetBackend(name, params.TransactionId)
	}
	if err != nil {
		respond.Error(w, err)
		return
	}
	handleDeprecatedBackendFields(http.MethodGet, bck, nil)
	respond.JSON(w, http.StatusOK, bck)
}

func (h *HandlerImpl) ReplaceBackend(w http.ResponseWriter, r *http.Request, name string, params ReplaceBackendParams) {
	if params.TransactionId != "" && params.ForceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data Backend
	if !respond.DecodeBody(r, w, &data) {
		return
	}

	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if data.ForcePersist != nil || data.IgnorePersist != nil {
		_, onDisk, confErr := cfg.GetBackend(name, params.TransactionId)
		if confErr != nil {
			respond.Error(w, confErr)
			return
		}
		handleDeprecatedBackendFields(http.MethodPut, &data, onDisk)
	}

	if params.FullSection {
		err = cfg.EditStructuredBackend(name, &data, params.TransactionId, int64(params.Version))
	} else {
		err = cfg.EditBackend(name, &data, params.TransactionId, int64(params.Version))
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

// handleDeprecatedBackendFields adds backward compatibility support for the fields
// force_persist and ignore_persist that are deprecated in favour of force_persist_list
// and ignore_persist_list.
func handleDeprecatedBackendFields(method string, payload *models.Backend, onDisk *models.Backend) {
	if method == http.MethodGet {
		if len(payload.ForcePersistList) > 0 {
			payload.ForcePersist = &models.BackendForcePersist{
				Cond:     payload.ForcePersistList[0].Cond,
				CondTest: payload.ForcePersistList[0].CondTest,
			}
		}
		if len(payload.IgnorePersistList) > 0 {
			payload.IgnorePersist = &models.BackendIgnorePersist{
				Cond:     payload.IgnorePersistList[0].Cond,
				CondTest: payload.IgnorePersistList[0].CondTest,
			}
		}
		return
	}

	if payload.ForcePersist != nil && len(payload.ForcePersistList) == 0 {
		if method == http.MethodPost || (method == http.MethodPut && (onDisk == nil || len(onDisk.ForcePersistList) == 0)) {
			payload.ForcePersistList = []*models.ForcePersist{{
				Cond:     payload.ForcePersist.Cond,
				CondTest: payload.ForcePersist.CondTest,
			}}
		} else {
			found := -1
			for i, item := range onDisk.ForcePersistList {
				if *item.Cond == *payload.ForcePersist.Cond && *item.CondTest == *payload.ForcePersist.CondTest {
					found = i
					break
				}
			}
			switch found {
			case -1:
				payload.ForcePersistList = slices.Insert(onDisk.ForcePersistList, 0, &models.ForcePersist{
					Cond:     payload.ForcePersist.Cond,
					CondTest: payload.ForcePersist.CondTest,
				})
			case 0:
				payload.ForcePersistList = onDisk.ForcePersistList
			default:
				payload.ForcePersistList = slices.Concat(onDisk.ForcePersistList[found:found+1], onDisk.ForcePersistList[:found], onDisk.ForcePersistList[found+1:])
			}
		}
	}

	if payload.IgnorePersist != nil && len(payload.IgnorePersistList) == 0 {
		if method == http.MethodPost || (method == http.MethodPut && (onDisk == nil || len(onDisk.IgnorePersistList) == 0)) {
			payload.IgnorePersistList = []*models.IgnorePersist{{
				Cond:     payload.IgnorePersist.Cond,
				CondTest: payload.IgnorePersist.CondTest,
			}}
		} else {
			found := -1
			for i, item := range onDisk.IgnorePersistList {
				if *item.Cond == *payload.IgnorePersist.Cond && *item.CondTest == *payload.IgnorePersist.CondTest {
					found = i
					break
				}
			}
			switch found {
			case -1:
				payload.IgnorePersistList = slices.Insert(onDisk.IgnorePersistList, 0, &models.IgnorePersist{
					Cond:     payload.IgnorePersist.Cond,
					CondTest: payload.IgnorePersist.CondTest,
				})
			case 0:
				payload.IgnorePersistList = onDisk.IgnorePersistList
			default:
				payload.IgnorePersistList = slices.Concat(onDisk.IgnorePersistList[found:found+1], onDisk.IgnorePersistList[:found], onDisk.IgnorePersistList[found+1:])
			}
		}
	}

	payload.ForcePersist = nil
	payload.IgnorePersist = nil
}
