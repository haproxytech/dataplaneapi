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
//

package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	client_native "github.com/haproxytech/client-native/v6"
	"github.com/haproxytech/client-native/v6/models"

	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/operations/spoe_transactions"
)

// SpoeTransactionsStartSpoeTransactionHandlerImpl implementation of the SpoeTransactionsStartSpoeTransactionHandler interface
type SpoeTransactionsStartSpoeTransactionHandlerImpl struct {
	Client client_native.HAProxyClient
}

// Handle executing the request and returning a response
func (h *SpoeTransactionsStartSpoeTransactionHandlerImpl) Handle(params spoe_transactions.StartSpoeTransactionParams, principal interface{}) middleware.Responder {
	spoeStorage, err := h.Client.Spoe()
	if err != nil {
		e := misc.HandleError(err)
		return spoe_transactions.NewStartSpoeTransactionDefault(int(*e.Code)).WithPayload(e)
	}

	ss, err := spoeStorage.GetSingleSpoe(params.Spoe)
	if err != nil {
		e := misc.HandleError(err)
		return spoe_transactions.NewStartSpoeTransactionDefault(int(*e.Code)).WithPayload(e)
	}
	t, err := ss.Transaction.StartTransaction(params.Version)
	if err != nil {
		e := misc.HandleError(err)
		return spoe_transactions.NewStartSpoeTransactionDefault(int(*e.Code)).WithPayload(e)
	}
	m := &models.SpoeTransaction{
		ID:      t.ID,
		Version: t.Version,
		Status:  t.Status,
	}
	return spoe_transactions.NewStartSpoeTransactionCreated().WithPayload(m)
}

// SpoeTransactionsDeleteSpoeTransactionHandlerImpl implementation of the DeleteTransactionHandler interface
type SpoeTransactionsDeleteSpoeTransactionHandlerImpl struct {
	Client client_native.HAProxyClient
}

// Handle executing the request and returning a response
func (h *SpoeTransactionsDeleteSpoeTransactionHandlerImpl) Handle(params spoe_transactions.DeleteSpoeTransactionParams, principal interface{}) middleware.Responder {
	spoeStorage, err := h.Client.Spoe()
	if err != nil {
		e := misc.HandleError(err)
		return spoe_transactions.NewStartSpoeTransactionDefault(int(*e.Code)).WithPayload(e)
	}

	ss, err := spoeStorage.GetSingleSpoe(params.Spoe)
	if err != nil {
		e := misc.HandleError(err)
		return spoe_transactions.NewStartSpoeTransactionDefault(int(*e.Code)).WithPayload(e)
	}
	err = ss.Transaction.DeleteTransaction(params.ID)
	if err != nil {
		e := misc.HandleError(err)
		return spoe_transactions.NewDeleteSpoeTransactionDefault(int(*e.Code)).WithPayload(e)
	}
	return spoe_transactions.NewDeleteSpoeTransactionNoContent()
}

// SpoeTransactionsGetSpoeTransactionHandlerImpl implementation of the SpoeTransactionsGetSpoeTransactionHandler interface
type SpoeTransactionsGetSpoeTransactionHandlerImpl struct {
	Client client_native.HAProxyClient
}

// Handle executing the request and returning a response
func (h *SpoeTransactionsGetSpoeTransactionHandlerImpl) Handle(params spoe_transactions.GetSpoeTransactionParams, principal interface{}) middleware.Responder {
	spoeStorage, err := h.Client.Spoe()
	if err != nil {
		e := misc.HandleError(err)
		return spoe_transactions.NewStartSpoeTransactionDefault(int(*e.Code)).WithPayload(e)
	}

	ss, err := spoeStorage.GetSingleSpoe(params.Spoe)
	if err != nil {
		e := misc.HandleError(err)
		return spoe_transactions.NewStartSpoeTransactionDefault(int(*e.Code)).WithPayload(e)
	}
	t, err := ss.Transaction.GetTransaction(params.ID)
	if err != nil {
		e := misc.HandleError(err)
		return spoe_transactions.NewGetSpoeTransactionsDefault(int(*e.Code)).WithPayload(e)
	}
	m := &models.SpoeTransaction{
		ID:      t.ID,
		Version: t.Version,
		Status:  t.Status,
	}
	return spoe_transactions.NewGetSpoeTransactionOK().WithPayload(m)
}

// SpoeTransactionsGetSpoeTransactionsHandlerImpl implementation of the SpoeTransactionsGetSpoeTransactionsHandler interface
type SpoeTransactionsGetSpoeTransactionsHandlerImpl struct {
	Client client_native.HAProxyClient
}

// Handle executing the request and returning a response
func (h *SpoeTransactionsGetSpoeTransactionsHandlerImpl) Handle(params spoe_transactions.GetSpoeTransactionsParams, principal interface{}) middleware.Responder {
	spoeStorage, err := h.Client.Spoe()
	if err != nil {
		e := misc.HandleError(err)
		return spoe_transactions.NewStartSpoeTransactionDefault(int(*e.Code)).WithPayload(e)
	}

	ss, err := spoeStorage.GetSingleSpoe(params.Spoe)
	if err != nil {
		e := misc.HandleError(err)
		return spoe_transactions.NewStartSpoeTransactionDefault(int(*e.Code)).WithPayload(e)
	}
	s := ""
	if params.Status != nil {
		s = *params.Status
	}
	ts, err := ss.Transaction.GetTransactions(s)
	if err != nil {
		e := misc.HandleError(err)
		return spoe_transactions.NewGetSpoeTransactionsDefault(int(*e.Code)).WithPayload(e)
	}
	var ms models.SpoeTransactions
	if *ts != nil && len(*ts) > 0 {
		for _, t := range *ts {
			m := &models.SpoeTransaction{
				Version: t.Version,
				ID:      t.ID,
				Status:  t.Status,
			}
			ms = append(ms, m)
		}
	}
	return spoe_transactions.NewGetSpoeTransactionsOK().WithPayload(ms)
}

// SpoeTransactionsCommitSpoeTransactionHandlerImpl implementation of the SpoeTransactionsCommitSpoeTransactionHandler interface
type SpoeTransactionsCommitSpoeTransactionHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// Handle executing the request and returning a response
func (h *SpoeTransactionsCommitSpoeTransactionHandlerImpl) Handle(params spoe_transactions.CommitSpoeTransactionParams, principal interface{}) middleware.Responder {
	spoeStorage, err := h.Client.Spoe()
	if err != nil {
		e := misc.HandleError(err)
		return spoe_transactions.NewStartSpoeTransactionDefault(int(*e.Code)).WithPayload(e)
	}

	ss, err := spoeStorage.GetSingleSpoe(params.Spoe)
	if err != nil {
		e := misc.HandleError(err)
		return spoe_transactions.NewStartSpoeTransactionDefault(int(*e.Code)).WithPayload(e)
	}
	t, err := ss.Transaction.CommitTransaction(params.ID)
	if err != nil {
		e := misc.HandleError(err)
		return spoe_transactions.NewCommitSpoeTransactionDefault(int(*e.Code)).WithPayload(e)
	}
	m := &models.SpoeTransaction{
		ID:      t.ID,
		Version: t.Version,
		Status:  t.Status,
	}
	if *params.ForceReload {
		err := h.ReloadAgent.ForceReload()
		if err != nil {
			e := misc.HandleError(err)
			return spoe_transactions.NewCommitSpoeTransactionDefault(int(*e.Code)).WithPayload(e)
		}
		return spoe_transactions.NewCommitSpoeTransactionOK().WithPayload(m)
	}
	rID := h.ReloadAgent.Reload()
	return spoe_transactions.NewCommitSpoeTransactionAccepted().WithReloadID(rID).WithPayload(m)
}
