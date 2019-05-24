// Copyright 2019 HAProxy Technologies
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this files except in compliance with the License.
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

	client_native "github.com/haproxytech/client-native"
	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/operations/configuration"
)

//GetRawConfigurationHandlerImpl implementation of the GetHAProxyConfigurationHandler interface
type GetRawConfigurationHandlerImpl struct {
	Client *client_native.HAProxyClient
}

// PostRawConfigurationHandlerImpl implementation of the PostHAProxyConfigurationHandler interface
type PostRawConfigurationHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent *haproxy.ReloadAgent
}

//Handle executing the request and returning a response
func (h *GetRawConfigurationHandlerImpl) Handle(params configuration.GetHAProxyConfigurationParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	v, data, err := h.Client.Configuration.GetRawConfiguration(t)
	if err != nil {
		e := misc.HandleError(err)
		return configuration.NewGetHAProxyConfigurationDefault(int(*e.Code)).WithPayload(e)
	}
	return configuration.NewGetHAProxyConfigurationOK().WithPayload(&configuration.GetHAProxyConfigurationOKBody{Version: v, Data: data})
}

//Handle executing the request and returning a response
func (h *PostRawConfigurationHandlerImpl) Handle(params configuration.PostHAProxyConfigurationParams, principal interface{}) middleware.Responder {
	v := int64(0)
	if params.Version != nil {
		v = *params.Version
	}

	err := h.Client.Configuration.PostRawConfiguration(&params.Data, v)
	if err != nil {
		e := misc.HandleError(err)
		return configuration.NewPostHAProxyConfigurationDefault(int(*e.Code)).WithPayload(e)
	}
	if *params.ForceReload {
		err := h.ReloadAgent.ForceReload()
		if err != nil {
			e := misc.HandleError(err)
			return configuration.NewPostHAProxyConfigurationDefault(int(*e.Code)).WithPayload(e)
		}
		return configuration.NewPostHAProxyConfigurationCreated().WithPayload(params.Data)
	}
	rID := h.ReloadAgent.Reload()
	return configuration.NewPostHAProxyConfigurationAccepted().WithReloadID(rID).WithPayload(params.Data)
}
