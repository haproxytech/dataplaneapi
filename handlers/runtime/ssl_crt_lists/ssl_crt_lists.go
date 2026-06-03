// Copyright 2025 HAProxy Technologies
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

package ssl_crt_lists

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	client_native "github.com/haproxytech/client-native/v6"

	"github.com/haproxytech/dataplaneapi/handlers/middleware"
	"github.com/haproxytech/dataplaneapi/handlers/respond"
)

// RegisterRouter registers all ssl_crt_lists routes onto r using spec-based request validation.
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

// HandlerImpl implements ServerInterface for HAProxy runtime SSL CRT lists.
type HandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h *HandlerImpl) GetAllCrtLists(w http.ResponseWriter, r *http.Request) {
	rt, err := h.Client.Runtime()
	if err != nil {
		respond.Error(w, err)
		return
	}
	files, err := rt.ShowCrtLists()
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, files)
}

func (h *HandlerImpl) DeleteCrtListEntry(w http.ResponseWriter, r *http.Request, params DeleteCrtListEntryParams) {
	rt, err := h.Client.Runtime()
	if err != nil {
		respond.Error(w, err)
		return
	}

	var lineNumber *int64
	if params.LineNumber != nil {
		lineNumber = new(int64(*params.LineNumber))
	}

	if err = rt.DeleteCrtListEntry(params.Name, params.CertFile, lineNumber); err != nil {
		respond.Error(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *HandlerImpl) GetAllCrtListEntries(w http.ResponseWriter, r *http.Request, params GetAllCrtListEntriesParams) {
	rt, err := h.Client.Runtime()
	if err != nil {
		respond.Error(w, err)
		return
	}
	crtList, err := rt.ShowCrtListEntries(params.Name)
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, crtList)
}

func (h *HandlerImpl) AddCrtListEntry(w http.ResponseWriter, r *http.Request, params AddCrtListEntryParams) {
	var data AddCrtListEntryJSONRequestBody
	if !respond.DecodeBody(r, w, &data) {
		return
	}

	rt, err := h.Client.Runtime()
	if err != nil {
		respond.Error(w, err)
		return
	}

	if err = rt.AddCrtListEntry(params.Name, data); err != nil {
		respond.Error(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
