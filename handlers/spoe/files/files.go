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

package files

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	client_native "github.com/haproxytech/client-native/v6"
	"github.com/haproxytech/client-native/v6/models"

	"github.com/haproxytech/dataplaneapi/handlers/middleware"
	"github.com/haproxytech/dataplaneapi/handlers/respond"
	"github.com/haproxytech/dataplaneapi/misc"
)

// RegisterRouter registers all SPOE file routes onto r using spec-based request validation
// and a shared error handler.
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

// HandlerImpl implements ServerInterface for HAProxy SPOE file management.
type HandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h *HandlerImpl) GetAllSpoeFiles(w http.ResponseWriter, r *http.Request) {
	spoeStorage, err := h.Client.Spoe()
	if err != nil {
		respond.Error(w, err)
		return
	}
	files, err := spoeStorage.GetAll()
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, files)
}

func (h *HandlerImpl) CreateSpoe(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		respond.MultipartError(w, err)
		return
	}
	file, header, err := r.FormFile("file_upload")
	if err != nil {
		respond.BadRequest(w, err.Error())
		return
	}
	defer file.Close()

	spoeStorage, err := h.Client.Spoe()
	if err != nil {
		respond.Error(w, err)
		return
	}
	path, err := spoeStorage.Create(header.Filename, file)
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusCreated, path)
}

func (h *HandlerImpl) DeleteSpoeFile(w http.ResponseWriter, r *http.Request, name string) {
	spoeStorage, err := h.Client.Spoe()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = spoeStorage.Delete(name); err != nil {
		respond.Error(w, err)
		return
	}
	respond.NoContent(w)
}

func (h *HandlerImpl) GetOneSpoeFile(w http.ResponseWriter, r *http.Request, name string) {
	spoeStorage, err := h.Client.Spoe()
	if err != nil {
		respond.Error(w, err)
		return
	}
	path, err := spoeStorage.Get(name)
	if err != nil {
		respond.Error(w, err)
		return
	}
	if path == "" {
		code := misc.ErrHTTPNotFound
		msg := "spoe file " + name + " not found"
		respond.JSON(w, http.StatusNotFound, &models.Error{Code: &code, Message: &msg})
		return
	}
	respond.JSON(w, http.StatusOK, struct {
		Data string `json:"data,omitempty"`
	}{Data: path})
}
