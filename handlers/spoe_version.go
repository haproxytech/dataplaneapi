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
	"github.com/haproxytech/dataplaneapi/operations/spoe"
)

// SpoeGetSpoeConfigurationVersionHandlerImpl implementation of the SpoeGetSpoeConfigurationVersionHandler interface
type SpoeGetSpoeConfigurationVersionHandlerImpl struct {
	Client client_native.HAProxyClient
}

// Handle executing the request and returning a response
func (h *SpoeGetSpoeConfigurationVersionHandlerImpl) Handle(params spoe.GetSpoeConfigurationVersionParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	spoeStorage, err := h.Client.Spoe()
	if err != nil {
		e := misc.HandleError(err)
		return spoe.NewGetSpoeScopeDefault(int(*e.Code)).WithPayload(e)
	}

	ss, err := spoeStorage.GetSingleSpoe(params.Spoe)
	if err != nil {
		status := misc.GetHTTPStatusFromErr(err)
		return spoe.NewGetSpoeScopeDefault(status).WithPayload(misc.SetError(status, err.Error()))
	}
	v, err := ss.GetConfigurationVersion(t)
	if err != nil {
		e := misc.HandleError(err)
		return spoe.NewGetSpoeConfigurationVersionDefault(int(*e.Code)).WithPayload(e)
	}
	return spoe.NewGetSpoeConfigurationVersionOK().WithPayload(v)
}
