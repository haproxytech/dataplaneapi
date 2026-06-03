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

package ssl_crl_files

import (
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	client_native "github.com/haproxytech/client-native/v6"

	"github.com/haproxytech/dataplaneapi/handlers/middleware"
	"github.com/haproxytech/dataplaneapi/handlers/respond"
)

// RegisterRouter registers all ssl_crl_files routes onto r using spec-based request validation.
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

// HandlerImpl implements ServerInterface for HAProxy runtime SSL CRL files.
type HandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h *HandlerImpl) GetAllCrl(w http.ResponseWriter, r *http.Request) {
	rt, err := h.Client.Runtime()
	if err != nil {
		respond.Error(w, err)
		return
	}
	files, err := rt.ShowCrlFiles()
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, files)
}

func (h *HandlerImpl) CreateCrl(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		respond.BadRequest(w, err.Error())
		return
	}
	file, header, err := r.FormFile("file_upload")
	if err != nil {
		respond.BadRequest(w, err.Error())
		return
	}
	defer file.Close()

	payload, err := io.ReadAll(file)
	if err != nil {
		respond.Error(w, err)
		return
	}

	rt, err := h.Client.Runtime()
	if err != nil {
		respond.Error(w, err)
		return
	}

	if err = rt.NewCrlFile(header.Filename); err != nil {
		respond.Error(w, err)
		return
	}

	if err = rt.SetCrlFile(header.Filename, string(payload)); err != nil {
		respond.Error(w, err)
		return
	}

	if err = rt.CommitCrlFile(header.Filename); err != nil {
		respond.Error(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *HandlerImpl) DeleteCrl(w http.ResponseWriter, r *http.Request, name string) {
	rt, err := h.Client.Runtime()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = rt.DeleteCrlFile(name); err != nil {
		respond.Error(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *HandlerImpl) GetCrl(w http.ResponseWriter, r *http.Request, name string, params GetCrlParams) {
	rt, err := h.Client.Runtime()
	if err != nil {
		respond.Error(w, err)
		return
	}

	var idx *int64
	if params.Index != nil {
		idx = new(int64(*params.Index))
	}

	entries, err := rt.ShowCrlFile(name, idx)
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, *entries)
}

func (h *HandlerImpl) ReplaceCrl(w http.ResponseWriter, r *http.Request, name string) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		respond.BadRequest(w, err.Error())
		return
	}
	file, _, err := r.FormFile("file_upload")
	if err != nil {
		respond.BadRequest(w, err.Error())
		return
	}
	defer file.Close()

	payload, err := io.ReadAll(file)
	if err != nil {
		respond.Error(w, err)
		return
	}

	rt, err := h.Client.Runtime()
	if err != nil {
		respond.Error(w, err)
		return
	}

	if err = rt.SetCrlFile(name, string(payload)); err != nil {
		respond.Error(w, err)
		return
	}

	if err = rt.CommitCrlFile(name); err != nil {
		respond.Error(w, err)
		return
	}

	respond.JSON(w, http.StatusOK, nil)
}
