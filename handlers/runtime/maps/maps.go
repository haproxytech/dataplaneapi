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

package maps

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	client_native "github.com/haproxytech/client-native/v6"
	"github.com/haproxytech/client-native/v6/models"

	config "github.com/haproxytech/dataplaneapi/configuration"
	"github.com/haproxytech/dataplaneapi/handlers/middleware"
	"github.com/haproxytech/dataplaneapi/handlers/respond"
)

// RegisterRouter registers all maps routes onto r using spec-based request validation.
func RegisterRouter(r chi.Router, client client_native.HAProxyClient) error {
	spec, err := GetSpec()
	if err != nil {
		return err
	}
	HandlerWithOptions(&HandlerImpl{Client: client}, ChiServerOptions{
		BaseRouter:       r,
		Middlewares:      []MiddlewareFunc{middleware.NewValidator(spec)},
		ErrorHandlerFunc: middleware.ErrorHandler,
	})
	return nil
}

// HandlerImpl implements ServerInterface for HAProxy runtime maps.
type HandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h *HandlerImpl) GetAllRuntimeMapFiles(w http.ResponseWriter, r *http.Request, params GetAllRuntimeMapFilesParams) {
	mapList := []*models.Map{}

	rt, err := h.Client.Runtime()
	if err != nil {
		respond.Error(w, err)
		return
	}

	runtimeMaps, err := rt.ShowMaps()
	if err != nil {
		respond.RuntimeError(w, err)
		return
	}

	mapsDir, err := rt.GetMapsDir()
	if err != nil {
		respond.RuntimeError(w, err)
		return
	}

	for _, m := range runtimeMaps {
		if params.IncludeUnmanaged || strings.HasPrefix(filepath.Dir(m.File), mapsDir) {
			if strings.HasPrefix(filepath.Dir(m.File), mapsDir) {
				m.StorageName = filepath.Base(m.File)
			}
			mapList = append(mapList, m)
		}
	}
	respond.JSON(w, http.StatusOK, mapList)
}

func (h *HandlerImpl) ClearRuntimeMap(w http.ResponseWriter, r *http.Request, name string, params ClearRuntimeMapParams) {
	rt, err := h.Client.Runtime()
	if err != nil {
		respond.Error(w, err)
		return
	}

	if err = rt.ClearMap(name, params.ForceDelete); err != nil {
		respond.RuntimeError(w, err)
		return
	}

	if params.ForceSync {
		m, err := rt.GetMap(name)
		if err != nil {
			respond.RuntimeError(w, err)
			return
		}
		ms := config.NewMapSync()
		if _, err = ms.Sync(m, h.Client); err != nil {
			respond.RuntimeError(w, err)
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *HandlerImpl) GetOneRuntimeMap(w http.ResponseWriter, r *http.Request, name string) {
	rt, err := h.Client.Runtime()
	if err != nil {
		respond.Error(w, err)
		return
	}
	m, err := rt.GetMap(name)
	if err != nil {
		respond.RuntimeError(w, err)
		return
	}
	if m == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	respond.JSON(w, http.StatusOK, m)
}

func (h *HandlerImpl) AddPayloadRuntimeMap(w http.ResponseWriter, r *http.Request, name string, params AddPayloadRuntimeMapParams) {
	var data AddPayloadRuntimeMapJSONRequestBody
	if !respond.DecodeBody(r, w, &data) {
		return
	}

	rt, err := h.Client.Runtime()
	if err != nil {
		respond.Error(w, err)
		return
	}

	if err = rt.AddMapPayloadVersioned(name, data); err != nil {
		respond.RuntimeError(w, err)
		return
	}

	if params.ForceSync {
		m, err := rt.GetMap(name)
		if err != nil {
			respond.RuntimeError(w, err)
			return
		}
		ms := config.NewMapSync()
		if _, err = ms.Sync(m, h.Client); err != nil {
			respond.RuntimeError(w, err)
			return
		}
	}
	respond.JSON(w, http.StatusCreated, data)
}
