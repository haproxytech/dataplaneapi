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

	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/operations/configuration"
)

// GetMapsHandlerImpl implementation of the GetAllRuntimeMapFilesHandler interface using client-native client
type ConfigurationGetConfigurationVersionHandlerImpl struct {
	Client client_native.HAProxyClient
}

// Handle executing the request and returning a response
func (h *ConfigurationGetConfigurationVersionHandlerImpl) Handle(params configuration.GetConfigurationVersionParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	cfg, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return configuration.NewGetConfigurationVersionDefault(int(*e.Code)).WithPayload(e)
	}

	v, err := cfg.GetConfigurationVersion(t)
	if err != nil {
		e := misc.HandleError(err)
		return configuration.NewGetConfigurationVersionDefault(int(*e.Code)).WithPayload(e)
	}
	return configuration.NewGetConfigurationVersionOK().WithPayload(v)
}
