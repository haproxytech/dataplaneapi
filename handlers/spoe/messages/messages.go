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

package messages

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	client_native "github.com/haproxytech/client-native/v6"

	"github.com/haproxytech/dataplaneapi/handlers/middleware"
	"github.com/haproxytech/dataplaneapi/handlers/respond"
)

// RegisterRouter registers all SPOE message routes onto r using spec-based request validation
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

// HandlerImpl implements ServerInterface for HAProxy SPOE message configuration.
type HandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h *HandlerImpl) GetAllSpoeMessage(w http.ResponseWriter, r *http.Request, parentName string, scopeName string, params GetAllSpoeMessageParams) {
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
	_, messages, err := ss.GetMessages(scopeName, params.TransactionId)
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, messages)
}

func (h *HandlerImpl) CreateSpoeMessage(w http.ResponseWriter, r *http.Request, parentName string, scopeName string, params CreateSpoeMessageParams) {
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
	var data SpoeMessage
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	if err = ss.CreateMessage(scopeName, &data, params.TransactionId, int64(params.Version)); err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusCreated, &data)
}

func (h *HandlerImpl) DeleteSpoeMessage(w http.ResponseWriter, r *http.Request, parentName string, scopeName string, name string, params DeleteSpoeMessageParams) {
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
	if err = ss.DeleteMessage(scopeName, name, params.TransactionId, int64(params.Version)); err != nil {
		respond.Error(w, err)
		return
	}
	respond.NoContent(w)
}

func (h *HandlerImpl) GetSpoeMessage(w http.ResponseWriter, r *http.Request, parentName string, scopeName string, name string, params GetSpoeMessageParams) {
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
	_, message, err := ss.GetMessage(scopeName, name, params.TransactionId)
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, message)
}

func (h *HandlerImpl) ReplaceSpoeMessage(w http.ResponseWriter, r *http.Request, parentName string, scopeName string, name string, params ReplaceSpoeMessageParams) {
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
	var data SpoeMessage
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	if err = ss.EditMessage(scopeName, &data, name, params.TransactionId, int64(params.Version)); err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, &data)
}
