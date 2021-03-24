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
	"github.com/haproxytech/client-native/v2/models"

	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/operations/defaults"
)

// GetDefaultsHandlerImpl implementation of the GetDefaultsHandler interface
type GetDefaultsHandlerImpl struct {
	Client *client_native.HAProxyClient
}

// ReplaceDefaultsHandlerImpl implementation of the ReplaceDefaultsHandler interface
type ReplaceDefaultsHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// Handle executing the request and returning a response
func (h *GetDefaultsHandlerImpl) Handle(params defaults.GetDefaultsParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	v, data, err := h.Client.Configuration.GetDefaultsConfiguration(t)
	if err != nil {
		e := misc.HandleError(err)
		return defaults.NewGetDefaultsDefault(int(*e.Code)).WithPayload(e).WithConfigurationVersion(v)
	}
	return defaults.NewGetDefaultsOK().WithPayload(&defaults.GetDefaultsOKBody{Version: v, Data: data}).WithConfigurationVersion(v)
}

// Handle executing the request and returning a response
func (h *ReplaceDefaultsHandlerImpl) Handle(params defaults.ReplaceDefaultsParams, principal interface{}) middleware.Responder {
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
		return defaults.NewReplaceDefaultsDefault(int(*e.Code)).WithPayload(e)
	}

	err := h.Client.Configuration.PushDefaultsConfiguration(params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return defaults.NewReplaceDefaultsDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return defaults.NewReplaceDefaultsDefault(int(*e.Code)).WithPayload(e)
			}
			return defaults.NewReplaceDefaultsOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return defaults.NewReplaceDefaultsAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return defaults.NewReplaceDefaultsAccepted().WithPayload(params.Data)
}
