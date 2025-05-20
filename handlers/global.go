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
	client_conf "github.com/haproxytech/client-native/v6/configuration"
	"github.com/haproxytech/client-native/v6/models"
	cn "github.com/haproxytech/dataplaneapi/client-native"
	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/operations/global"
)

// GetGlobalHandlerImpl implementation of the GetGlobalHandler interface
type GetGlobalHandlerImpl struct {
	Client client_native.HAProxyClient
}

// ReplaceGlobalHandlerImpl implementation of the ReplaceGlobalHandler interface
type ReplaceGlobalHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// Handle executing the request and returning a response
func (h *GetGlobalHandlerImpl) Handle(params global.GetGlobalParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	_, data, err := h.getGlobalConfiguration(params, t)
	if err != nil {
		e := misc.HandleError(err)
		return global.NewGetGlobalDefault(int(*e.Code)).WithPayload(e)
	}
	return global.NewGetGlobalOK().WithPayload(data)
}

func (h *GetGlobalHandlerImpl) getGlobalConfiguration(params global.GetGlobalParams, t string) (int64, *models.Global, error) {
	configuration, err := h.Client.Configuration()
	if err != nil {
		return 0, nil, err
	}
	if params.FullSection != nil && *params.FullSection {
		return configuration.GetStructuredGlobalConfiguration(t)
	}
	return configuration.GetGlobalConfiguration(t)
}

// Handle executing the request and returning a response
func (h *ReplaceGlobalHandlerImpl) Handle(params global.ReplaceGlobalParams, principal interface{}) middleware.Responder {
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
		return global.NewReplaceGlobalDefault(int(*e.Code)).WithPayload(e)
	}
	// validate constraints that can not be validated by the swagger 2.0 spec.
	if err := client_conf.ValidateGlobalSection(params.Data); err != nil {
		code := misc.ErrHTTPBadRequest
		msg := err.Error()
		e := &models.Error{Code: &code, Message: &msg}
		return global.NewReplaceGlobalDefault(int(misc.ErrHTTPBadRequest)).WithPayload(e)
	}

	err := h.pushGlobalConfiguration(params, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return global.NewReplaceGlobalDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		callbackNeeded, reconfigureFunc, err := cn.ReconfigureRuntime(h.Client)
		if err != nil {
			e := misc.HandleError(err)
			return global.NewReplaceGlobalDefault(int(*e.Code)).WithPayload(e)
		}

		if *params.ForceReload {
			if callbackNeeded {
				err = h.ReloadAgent.ForceReloadWithCallback(reconfigureFunc)
			} else {
				err = h.ReloadAgent.ForceReload()
			}

			if err != nil {
				e := misc.HandleError(err)
				return global.NewReplaceGlobalDefault(int(*e.Code)).WithPayload(e)
			}
			return global.NewReplaceGlobalOK().WithPayload(params.Data)
		}

		var rID string
		if callbackNeeded {
			rID = h.ReloadAgent.ReloadWithCallback(reconfigureFunc)
		} else {
			rID = h.ReloadAgent.Reload()
		}
		return global.NewReplaceGlobalAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return global.NewReplaceGlobalAccepted().WithPayload(params.Data)
}

func (h *ReplaceGlobalHandlerImpl) pushGlobalConfiguration(params global.ReplaceGlobalParams, t string, v int64) error {
	configuration, err := h.Client.Configuration()
	if err != nil {
		return err
	}
	if params.FullSection != nil && *params.FullSection {
		return configuration.PushStructuredGlobalConfiguration(params.Data, t, v)
	}
	return configuration.PushGlobalConfiguration(params.Data, t, v)
}
