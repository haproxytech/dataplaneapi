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

package transactions

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	client_native "github.com/haproxytech/client-native/v6"
	"github.com/haproxytech/client-native/v6/models"

	"github.com/haproxytech/dataplaneapi/handlers/middleware"
	"github.com/haproxytech/dataplaneapi/handlers/respond"
	"github.com/haproxytech/dataplaneapi/haproxy"
)

// RegisterRouter registers all SPOE transaction routes onto r using spec-based request validation
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

// HandlerImpl implements ServerInterface for HAProxy SPOE transaction management.
type HandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

func (h *HandlerImpl) GetAllSpoeTransaction(w http.ResponseWriter, r *http.Request, parentName string, params GetAllSpoeTransactionParams) {
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
	ts, err := ss.Transaction.GetTransactions(string(params.Status))
	if err != nil {
		respond.Error(w, err)
		return
	}
	var result models.SpoeTransactions
	for _, t := range *ts {
		m := &models.SpoeTransaction{
			ID:      t.ID,
			Version: t.Version,
			Status:  t.Status,
		}
		result = append(result, m)
	}
	respond.JSON(w, http.StatusOK, result)
}

func (h *HandlerImpl) StartSpoeTransaction(w http.ResponseWriter, r *http.Request, parentName string, params StartSpoeTransactionParams) {
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
	t, err := ss.Transaction.StartTransaction(int64(params.Version))
	if err != nil {
		respond.Error(w, err)
		return
	}
	m := &models.SpoeTransaction{
		ID:      t.ID,
		Version: t.Version,
		Status:  t.Status,
	}
	respond.JSON(w, http.StatusCreated, m)
}

func (h *HandlerImpl) DeleteSpoeTransaction(w http.ResponseWriter, r *http.Request, parentName string, id string) {
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
	if err = ss.Transaction.DeleteTransaction(id); err != nil {
		respond.Error(w, err)
		return
	}
	respond.NoContent(w)
}

func (h *HandlerImpl) GetSpoeTransaction(w http.ResponseWriter, r *http.Request, parentName string, id string) {
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
	t, err := ss.Transaction.GetTransaction(id)
	if err != nil {
		respond.Error(w, err)
		return
	}
	m := &models.SpoeTransaction{
		ID:      t.ID,
		Version: t.Version,
		Status:  t.Status,
	}
	respond.JSON(w, http.StatusOK, m)
}

func (h *HandlerImpl) CommitSpoeTransaction(w http.ResponseWriter, r *http.Request, parentName string, id string, params CommitSpoeTransactionParams) {
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
	t, err := ss.Transaction.CommitTransaction(id)
	if err != nil {
		respond.Error(w, err)
		return
	}
	m := &models.SpoeTransaction{
		ID:      t.ID,
		Version: t.Version,
		Status:  t.Status,
	}
	if params.ForceReload {
		if err = h.ReloadAgent.ForceReload(); err != nil {
			respond.Error(w, err)
			return
		}
		respond.JSON(w, http.StatusOK, m)
		return
	}
	respond.Accepted(w, h.ReloadAgent.Reload(), m)
}
