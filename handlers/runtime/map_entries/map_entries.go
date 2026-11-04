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

package map_entries

import (
	"errors"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	client_native "github.com/haproxytech/client-native/v6"
	"github.com/haproxytech/client-native/v6/models"
	"github.com/haproxytech/client-native/v6/storage"

	config "github.com/haproxytech/dataplaneapi/configuration"
	"github.com/haproxytech/dataplaneapi/handlers/middleware"
	"github.com/haproxytech/dataplaneapi/handlers/respond"
)

// RegisterRouter registers all map_entries routes onto r using spec-based request validation.
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

// HandlerImpl implements ServerInterface for HAProxy runtime map entries.
type HandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h *HandlerImpl) ShowRuntimeMap(w http.ResponseWriter, r *http.Request, parentName string) {
	rt, err := h.Client.Runtime()
	if err != nil {
		respond.Error(w, err)
		return
	}
	m, err := rt.ShowMapEntries(parentName)
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

func (h *HandlerImpl) AddMapEntry(w http.ResponseWriter, r *http.Request, parentName string, params AddMapEntryParams) {
	var data AddMapEntryJSONRequestBody
	if !respond.DecodeBody(r, w, &data) {
		return
	}

	rt, err := h.Client.Runtime()
	if err != nil {
		respond.Error(w, err)
		return
	}

	if err = rt.AddMapEntry(parentName, data.Key, data.Value); err != nil {
		respond.RuntimeError(w, err)
		return
	}

	if params.ForceSync {
		m, err := rt.GetMap(parentName)
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
	respond.JSON(w, http.StatusCreated, &data)
}

func (h *HandlerImpl) ReplaceRuntimeMapEntries(w http.ResponseWriter, r *http.Request, parentName string, params ReplaceRuntimeMapEntriesParams) {
	var data ReplaceRuntimeMapEntriesJSONRequestBody
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	contents, err := serializeMapEntries(data)
	if err != nil {
		respond.BadRequest(w, err.Error())
		return
	}

	rt, err := h.Client.Runtime()
	if err != nil {
		respond.Error(w, err)
		return
	}

	var storageName string
	var st storage.Storage
	if params.ForceSync {
		m, err := rt.GetMap(parentName)
		if err != nil {
			respond.RuntimeError(w, err)
			return
		}
		mapsDir, err := rt.GetMapsDir()
		if err != nil {
			respond.RuntimeError(w, err)
			return
		}
		storageName, err = managedMapStorageName(m, mapsDir)
		if err != nil {
			respond.BadRequest(w, err.Error())
			return
		}
		st, err = h.Client.MapStorage()
		if err != nil {
			respond.Error(w, err)
			return
		}
	}

	if err = rt.ReplaceMapPayload(parentName, data); err != nil {
		respond.RuntimeError(w, err)
		return
	}

	if params.ForceSync {
		if _, err = st.Replace(storageName, contents); err != nil {
			respond.RuntimeError(w, err)
			return
		}
	}

	respond.JSON(w, http.StatusOK, data)
}

func managedMapStorageName(m *models.Map, mapsDir string) (string, error) {
	if m == nil || m.File == "" || mapsDir == "" {
		return "", errors.New("cannot identify managed storage for runtime map")
	}
	relative, err := filepath.Rel(filepath.Clean(mapsDir), filepath.Clean(m.File))
	if err != nil || relative == "." || filepath.IsAbs(relative) ||
		relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) ||
		filepath.Base(relative) != relative {
		return "", errors.New("cannot force sync an unmanaged runtime map")
	}
	return relative, nil
}

func serializeMapEntries(entries models.MapEntries) (string, error) {
	var contents strings.Builder
	for _, entry := range entries {
		if entry == nil {
			return "", errors.New("map entry cannot be null")
		}
		contents.WriteString(entry.Key)
		contents.WriteByte(' ')
		contents.WriteString(entry.Value)
		contents.WriteByte('\n')
	}
	return contents.String(), nil
}

func (h *HandlerImpl) DeleteRuntimeMapEntry(w http.ResponseWriter, r *http.Request, parentName string, id string, params DeleteRuntimeMapEntryParams) {
	rt, err := h.Client.Runtime()
	if err != nil {
		respond.Error(w, err)
		return
	}

	if err = rt.DeleteMapEntry(parentName, id); err != nil {
		respond.RuntimeError(w, err)
		return
	}

	if params.ForceSync {
		m, err := rt.GetMap(parentName)
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

func (h *HandlerImpl) GetRuntimeMapEntry(w http.ResponseWriter, r *http.Request, parentName string, id string) {
	rt, err := h.Client.Runtime()
	if err != nil {
		respond.Error(w, err)
		return
	}
	m, err := rt.GetMapEntry(parentName, id)
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

func (h *HandlerImpl) ReplaceRuntimeMapEntry(w http.ResponseWriter, r *http.Request, parentName string, id string, params ReplaceRuntimeMapEntryParams) {
	var data ReplaceRuntimeMapEntryJSONRequestBody
	if !respond.DecodeJSON(r, w, &data) {
		return
	}

	rt, err := h.Client.Runtime()
	if err != nil {
		respond.Error(w, err)
		return
	}

	if err = rt.SetMapEntry(parentName, id, data.Value); err != nil {
		respond.RuntimeError(w, err)
		return
	}

	e, err := rt.GetMapEntry(parentName, id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if params.ForceSync {
		m, err := rt.GetMap(parentName)
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
	respond.JSON(w, http.StatusOK, e)
}
