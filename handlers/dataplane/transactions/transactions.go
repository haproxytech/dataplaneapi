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
	"strings"
	"sync"

	"github.com/go-chi/chi/v5"
	client_native "github.com/haproxytech/client-native/v6"
	"github.com/haproxytech/client-native/v6/models"

	cn "github.com/haproxytech/dataplaneapi/client-native"
	"github.com/haproxytech/dataplaneapi/handlers/middleware"
	"github.com/haproxytech/dataplaneapi/handlers/respond"
	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/log"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/rate"
)

// RegisterRouter registers all transaction routes onto r using spec-based request validation
// and a shared error handler. When maxOpenTransactions > 0 the StartTransaction endpoint
// is rate-limited to that threshold.
func RegisterRouter(r chi.Router, client client_native.HAProxyClient, ra haproxy.IReloadAgent, maxOpenTransactions int) error {
	spec, err := GetSpec()
	if err != nil {
		return err
	}
	impl := &HandlerImpl{Client: client, ReloadAgent: ra, Mutex: &sync.Mutex{}}
	if maxOpenTransactions > 0 {
		// On any failure the count falls back to 0 ("nothing open"), which
		// disables the max-open-transactions limit for that request — so log it.
		actualCount := func() uint64 {
			cfg, err := client.Configuration()
			if err != nil {
				log.Errorf("transaction limiter: cannot get configuration client: %v", err)
				return 0
			}
			ts, err := cfg.GetTransactions(models.TransactionStatusInProgress)
			if err != nil {
				log.Errorf("transaction limiter: cannot count open transactions: %v", err)
				return 0
			}
			if ts == nil {
				return 0
			}
			return uint64(len(*ts))
		}
		impl.Limiter = rate.NewThresholdLimit(uint64(maxOpenTransactions), actualCount)
	}
	HandlerWithOptions(impl, ChiServerOptions{
		BaseRouter:       r,
		Middlewares:      []MiddlewareFunc{middleware.NewValidator(spec)},
		ErrorHandlerFunc: middleware.ErrorHandler,
	})
	return nil
}

// HandlerImpl implements ServerInterface for HAProxy transaction management.
type HandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
	Mutex       *sync.Mutex
	Limiter     rate.Threshold
}

func (h *HandlerImpl) GetTransactions(w http.ResponseWriter, r *http.Request, params GetTransactionsParams) {
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	ts, err := cfg.GetTransactions(string(params.Status))
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, *ts)
}

func (h *HandlerImpl) StartTransaction(w http.ResponseWriter, r *http.Request, params StartTransactionParams) {
	if h.Limiter != nil {
		if err := h.Limiter.LimitReached(); err != nil {
			respond.Error(w, err)
			return
		}
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	t, err := cfg.StartTransaction(int64(params.Version))
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusCreated, t)
}

func (h *HandlerImpl) DeleteTransaction(w http.ResponseWriter, r *http.Request, id string) {
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.DeleteTransaction(id); err != nil {
		e := misc.HandleError(err)
		if strings.HasSuffix(*e.Message, "does not exist") {
			e.Code = new(int64(http.StatusNotFound))
			respond.JSON(w, http.StatusNotFound, e)
			return
		}
		respond.JSON(w, int(*e.Code), e)
		return
	}
	respond.NoContent(w)
}

func (h *HandlerImpl) GetTransaction(w http.ResponseWriter, r *http.Request, id string) {
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	t, err := cfg.GetTransaction(id)
	if err != nil {
		e := misc.HandleError(err)
		if strings.HasSuffix(*e.Message, "does not exist") {
			e.Code = new(int64(http.StatusNotFound))
			respond.JSON(w, http.StatusNotFound, e)
			return
		}
		respond.JSON(w, int(*e.Code), e)
		return
	}
	respond.JSON(w, http.StatusOK, t)
}

func (h *HandlerImpl) CommitTransaction(w http.ResponseWriter, r *http.Request, id string, params CommitTransactionParams) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()

	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}

	transaction, err := cfg.GetTransaction(id)
	if err != nil {
		e := misc.HandleError(err)
		if strings.HasSuffix(*e.Message, "does not exist") {
			e.Code = new(int64(http.StatusNotFound))
			respond.JSON(w, http.StatusNotFound, e)
			return
		}
		respond.JSON(w, int(*e.Code), e)
		return
	}

	switch transaction.Status {
	case models.TransactionStatusOutdated:
		respond.JSON(w, http.StatusNotAcceptable, misc.OutdatedTransactionError(transaction.ID))
		return
	case models.TransactionStatusFailed:
		respond.JSON(w, http.StatusNotAcceptable, misc.FailedTransactionError(transaction.ID))
		return
	}

	t, err := cfg.CommitTransaction(id)
	if err != nil {
		e := misc.HandleError(err)
		if strings.HasSuffix(*e.Message, "does not exist") {
			e.Code = new(int64(http.StatusNotFound))
			respond.JSON(w, http.StatusNotFound, e)
			return
		}
		respond.JSON(w, int(*e.Code), e)
		return
	}

	// Mark outdated in-progress transactions
	txs, err := cfg.GetTransactions(models.TransactionStatusInProgress)
	if err != nil {
		respond.Error(w, err)
		return
	}
	for _, tx := range *txs {
		if tx.Version <= t.Version {
			_ = cfg.MarkTransactionOutdated(tx.ID)
		}
	}

	callbackNeeded, reconfigureFunc, err := cn.ReconfigureRuntime(h.Client)
	if err != nil {
		respond.Error(w, err)
		return
	}

	if params.ForceReload {
		if callbackNeeded {
			err = h.ReloadAgent.ForceReloadWithCallback(reconfigureFunc)
		} else {
			err = h.ReloadAgent.ForceReload()
		}
		if err != nil {
			respond.Error(w, err)
			return
		}
		respond.JSON(w, http.StatusOK, t)
		return
	}

	var rID string
	if callbackNeeded {
		rID = h.ReloadAgent.ReloadWithCallback(reconfigureFunc)
	} else {
		rID = h.ReloadAgent.Reload()
	}
	respond.Accepted(w, rID, t)
}
