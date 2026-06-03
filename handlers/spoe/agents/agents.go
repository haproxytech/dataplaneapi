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

package agents

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	client_native "github.com/haproxytech/client-native/v6"
	"github.com/haproxytech/client-native/v6/models"

	"github.com/haproxytech/dataplaneapi/handlers/middleware"
	"github.com/haproxytech/dataplaneapi/handlers/respond"
)

// RegisterRouter registers all SPOE agent routes onto r using spec-based request validation
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

// HandlerImpl implements ServerInterface for HAProxy SPOE agent configuration.
type HandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h *HandlerImpl) GetAllSpoeAgent(w http.ResponseWriter, r *http.Request, parentName string, scopeName string, params GetAllSpoeAgentParams) {
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
	_, agents, err := ss.GetAgents(scopeName, params.TransactionId)
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, agents)
}

func (h *HandlerImpl) CreateSpoeAgent(w http.ResponseWriter, r *http.Request, parentName string, scopeName string, params CreateSpoeAgentParams) {
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
	var data SpoeAgent
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	if err = ss.CreateAgent(scopeName, &data, params.TransactionId, int64(params.Version)); err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusCreated, &data)
}

func (h *HandlerImpl) DeleteSpoeAgent(w http.ResponseWriter, r *http.Request, parentName string, scopeName string, name string, params DeleteSpoeAgentParams) {
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
	if err = ss.DeleteAgent(scopeName, name, params.TransactionId, int64(params.Version)); err != nil {
		respond.Error(w, err)
		return
	}
	respond.NoContent(w)
}

func (h *HandlerImpl) GetSpoeAgent(w http.ResponseWriter, r *http.Request, parentName string, scopeName string, name string, params GetSpoeAgentParams) {
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
	_, agent, err := ss.GetAgent(scopeName, name, params.TransactionId)
	if err != nil {
		respond.Error(w, err)
		return
	}
	if agent == nil {
		code := int64(404)
		msg := fmt.Sprintf("SPOE agent %s not found in scope %s", name, scopeName)
		respond.JSON(w, http.StatusNotFound, &models.Error{Code: &code, Message: &msg})
		return
	}
	respond.JSON(w, http.StatusOK, agent)
}

func (h *HandlerImpl) ReplaceSpoeAgent(w http.ResponseWriter, r *http.Request, parentName string, scopeName string, name string, params ReplaceSpoeAgentParams) {
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
	var data SpoeAgent
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	if err = ss.EditAgent(scopeName, &data, params.TransactionId, int64(params.Version)); err != nil {
		respond.Error(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
