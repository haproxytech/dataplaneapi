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

package groups

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	client_native "github.com/haproxytech/client-native/v6"

	"github.com/haproxytech/dataplaneapi/handlers/middleware"
	"github.com/haproxytech/dataplaneapi/handlers/respond"
)

// RegisterRouter registers all SPOE group routes onto r using spec-based request validation
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

// HandlerImpl implements ServerInterface for HAProxy SPOE group configuration.
type HandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h *HandlerImpl) GetAllSpoeGroup(w http.ResponseWriter, r *http.Request, parentName string, scopeName string, params GetAllSpoeGroupParams) {
	spoeStorage, err := h.Client.Spoe()
	if err != nil {
		respond.Error(w, err)
		return
	}
	ss, err := spoeStorage.GetSingleSpoe(parentName)
	if err != nil {
		respond.Error(w, err)
		return
	}
	_, groups, err := ss.GetGroups(scopeName, params.TransactionId)
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, groups)
}

func (h *HandlerImpl) CreateSpoeGroup(w http.ResponseWriter, r *http.Request, parentName string, scopeName string, params CreateSpoeGroupParams) {
	spoeStorage, err := h.Client.Spoe()
	if err != nil {
		respond.Error(w, err)
		return
	}
	ss, err := spoeStorage.GetSingleSpoe(parentName)
	if err != nil {
		respond.Error(w, err)
		return
	}
	var data SpoeGroup
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	if err = ss.CreateGroup(scopeName, &data, params.TransactionId, int64(params.Version)); err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusCreated, &data)
}

func (h *HandlerImpl) DeleteSpoeGroup(w http.ResponseWriter, r *http.Request, parentName string, scopeName string, name string, params DeleteSpoeGroupParams) {
	spoeStorage, err := h.Client.Spoe()
	if err != nil {
		respond.Error(w, err)
		return
	}
	ss, err := spoeStorage.GetSingleSpoe(parentName)
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = ss.DeleteGroup(scopeName, name, params.TransactionId, int64(params.Version)); err != nil {
		respond.Error(w, err)
		return
	}
	respond.NoContent(w)
}

func (h *HandlerImpl) GetSpoeGroup(w http.ResponseWriter, r *http.Request, parentName string, scopeName string, name string, params GetSpoeGroupParams) {
	spoeStorage, err := h.Client.Spoe()
	if err != nil {
		respond.Error(w, err)
		return
	}
	ss, err := spoeStorage.GetSingleSpoe(parentName)
	if err != nil {
		respond.Error(w, err)
		return
	}
	_, group, err := ss.GetGroup(scopeName, name, params.TransactionId)
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, group)
}

func (h *HandlerImpl) ReplaceSpoeGroup(w http.ResponseWriter, r *http.Request, parentName string, scopeName string, name string, params ReplaceSpoeGroupParams) {
	spoeStorage, err := h.Client.Spoe()
	if err != nil {
		respond.Error(w, err)
		return
	}
	ss, err := spoeStorage.GetSingleSpoe(parentName)
	if err != nil {
		respond.Error(w, err)
		return
	}
	var data SpoeGroup
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	if err = ss.EditGroup(scopeName, &data, name, params.TransactionId, int64(params.Version)); err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, &data)
}
