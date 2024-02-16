// Copyright 2022 HAProxy Technologies
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
	"github.com/haproxytech/dataplaneapi/operations/cache"
)

// CreateCacheHandlerImpl implementation of the CreateCacheHandler interface using client-native client
type CreateCacheHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// DeleteCacheHandlerImpl implementation of the DeleteCacheHandler interface using client-native client
type DeleteCacheHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// GetCacheHandlerImpl implementation of the GetCacheHandler interface using client-native client
type GetCacheHandlerImpl struct {
	Client client_native.HAProxyClient
}

// GetCachesHandlerImpl implementation of the GetCachesHandler interface using client-native client
type GetCachesHandlerImpl struct {
	Client client_native.HAProxyClient
}

// ReplaceCacheHandlerImpl implementation of the ReplaceCacheHandler interface using client-native client
type ReplaceCacheHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// Handle executing the request and returning a response
func (h *CreateCacheHandlerImpl) Handle(params cache.CreateCacheParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	if t != "" && *params.ForceReload {
		msg := "Both force_reload and transaction specified, specify only one"
		c := misc.ErrHTTPBadRequest
		e := &models.Error{
			Message: &msg,
			Code:    &c,
		}
		return cache.NewCreateCacheDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return cache.NewCreateCacheDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.CreateCache(params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return cache.NewCreateCacheDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return cache.NewCreateCacheDefault(int(*e.Code)).WithPayload(e)
			}
			return cache.NewCreateCacheCreated().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return cache.NewCreateCacheAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return cache.NewCreateCacheAccepted().WithPayload(params.Data)
}

// Handle executing the request and returning a response
func (h *DeleteCacheHandlerImpl) Handle(params cache.DeleteCacheParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	if t != "" && *params.ForceReload {
		msg := "Both force_reload and transaction specified, specify only one"
		c := misc.ErrHTTPBadRequest
		e := &models.Error{
			Message: &msg,
			Code:    &c,
		}
		return cache.NewDeleteCacheDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return cache.NewDeleteCacheDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.DeleteCache(params.Name, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return cache.NewDeleteCacheDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return cache.NewDeleteCacheDefault(int(*e.Code)).WithPayload(e)
			}
			return cache.NewDeleteCacheNoContent()
		}
		rID := h.ReloadAgent.Reload()
		return cache.NewDeleteCacheAccepted().WithReloadID(rID)
	}
	return cache.NewDeleteCacheAccepted()
}

// Handle executing the request and returning a response
func (h *GetCacheHandlerImpl) Handle(params cache.GetCacheParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return cache.NewGetCacheDefault(int(*e.Code)).WithPayload(e)
	}

	_, r, err := configuration.GetCache(params.Name, t)
	if err != nil {
		e := misc.HandleError(err)
		return cache.NewGetCacheDefault(int(*e.Code)).WithPayload(e)
	}
	return cache.NewGetCacheOK().WithPayload(r)
}

// Handle executing the request and returning a response
func (h *GetCachesHandlerImpl) Handle(params cache.GetCachesParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return cache.NewGetCachesDefault(int(*e.Code)).WithPayload(e)
	}

	_, rs, err := configuration.GetCaches(t)
	if err != nil {
		e := misc.HandleError(err)
		return cache.NewGetCachesDefault(int(*e.Code)).WithPayload(e)
	}
	return cache.NewGetCachesOK().WithPayload(rs)
}

// Handle executing the request and returning a response
func (h *ReplaceCacheHandlerImpl) Handle(params cache.ReplaceCacheParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	if t != "" && *params.ForceReload {
		msg := "Both force_reload and transaction specified, specify only one"
		c := misc.ErrHTTPBadRequest
		e := &models.Error{
			Message: &msg,
			Code:    &c,
		}
		return cache.NewReplaceCacheDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return cache.NewReplaceCacheDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.EditCache(params.Name, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return cache.NewReplaceCacheDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return cache.NewReplaceCacheDefault(int(*e.Code)).WithPayload(e)
			}
			return cache.NewReplaceCacheOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return cache.NewReplaceCacheAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return cache.NewReplaceCacheAccepted().WithPayload(params.Data)
}
