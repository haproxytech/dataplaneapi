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
	client_native "github.com/haproxytech/client-native/v2"
	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/operations/transactions"
	log "github.com/sirupsen/logrus"
)

//StartTransactionHandlerImpl implementation of the StartTransactionHandler interface using client-native client
type StartTransactionHandlerImpl struct {
	Client *client_native.HAProxyClient
}

//DeleteTransactionHandlerImpl implementation of the DeleteTransactionHandler interface using client-native client
type DeleteTransactionHandlerImpl struct {
	Client *client_native.HAProxyClient
}

//GetTransactionHandlerImpl implementation of the GetTransactionHandler interface using client-native client
type GetTransactionHandlerImpl struct {
	Client *client_native.HAProxyClient
}

//GetTransactionsHandlerImpl implementation of the GetTransactionsHandler interface using client-native client
type GetTransactionsHandlerImpl struct {
	Client *client_native.HAProxyClient
}

//CommitTransactionHandlerImpl implementation of the CommitTransactionHandlerImpl interface using client-native client
type CommitTransactionHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

//Handle executing the request and returning a response
func (th *StartTransactionHandlerImpl) Handle(params transactions.StartTransactionParams, principal interface{}) middleware.Responder {
	t, err := th.Client.Configuration.StartTransaction(params.Version)
	if err != nil {
		e := misc.HandleError(err)
		return transactions.NewStartTransactionDefault(int(*e.Code)).WithPayload(e)
	}
	return transactions.NewStartTransactionCreated().WithPayload(t)
}

//Handle executing the request and returning a response
func (th *DeleteTransactionHandlerImpl) Handle(params transactions.DeleteTransactionParams, principal interface{}) middleware.Responder {
	err := th.Client.Configuration.DeleteTransaction(params.ID)
	if err != nil {
		e := misc.HandleError(err)
		return transactions.NewDeleteTransactionDefault(int(*e.Code)).WithPayload(e)
	}
	return transactions.NewDeleteTransactionNoContent()
}

//Handle executing the request and returning a response
func (th *GetTransactionHandlerImpl) Handle(params transactions.GetTransactionParams, principal interface{}) middleware.Responder {
	t, err := th.Client.Configuration.GetTransaction(params.ID)
	if err != nil {
		e := misc.HandleError(err)
		return transactions.NewGetTransactionsDefault(int(*e.Code)).WithPayload(e)
	}
	return transactions.NewGetTransactionOK().WithPayload(t)
}

//Handle executing the request and returning a response
func (th *GetTransactionsHandlerImpl) Handle(params transactions.GetTransactionsParams, principal interface{}) middleware.Responder {
	s := ""
	if params.Status != nil {
		s = *params.Status
	}
	ts, err := th.Client.Configuration.GetTransactions(s)
	if err != nil {
		e := misc.HandleError(err)
		return transactions.NewGetTransactionsDefault(int(*e.Code)).WithPayload(e)
	}
	return transactions.NewGetTransactionsOK().WithPayload(*ts)
}

//Handle executing the request and returning a response
func (th *CommitTransactionHandlerImpl) Handle(params transactions.CommitTransactionParams, principal interface{}) middleware.Responder {
	t, err := th.Client.Configuration.CommitTransaction(params.ID)
	if err != nil {
		e := misc.HandleError(err)
		return transactions.NewCommitTransactionDefault(int(*e.Code)).WithPayload(e)
	}
	if *params.ForceReload {
		err := th.ReloadAgent.ForceReload()
		if err != nil {
			if *params.EnableRevert {
				var errR error
				if errR = th.Client.Configuration.RevertTransaction(t); errR != nil {
					log.Debug("revert fail", errR)
				}
			}
			e := misc.HandleError(err)
			return transactions.NewCommitTransactionDefault(int(*e.Code)).WithPayload(e)
		}
		return transactions.NewCommitTransactionOK().WithPayload(t)
	}
	rID := th.ReloadAgent.Reload()
	return transactions.NewCommitTransactionAccepted().WithReloadID(rID).WithPayload(t)
}
