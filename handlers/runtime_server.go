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
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	client_native "github.com/haproxytech/client-native/v2"
	"github.com/haproxytech/client-native/v2/models"

	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/operations/server"
)

// GetRuntimeServerHandlerImpl implementation of the GetRuntimeServerHandler interface using client-native client
type GetRuntimeServerHandlerImpl struct {
	Client *client_native.HAProxyClient
}

// GetRuntimeServersHandlerImpl implementation of the GetRuntimeServersHandler interface using client-native client
type GetRuntimeServersHandlerImpl struct {
	Client *client_native.HAProxyClient
}

// ReplaceRuntimeServerHandlerImpl implementation of the ReplaceRuntimeServerHandler interface using client-native client
type ReplaceRuntimeServerHandlerImpl struct {
	Client *client_native.HAProxyClient
}

// Handle executing the request and returning a response
func (h *GetRuntimeServerHandlerImpl) Handle(params server.GetRuntimeServerParams, principal interface{}) middleware.Responder {
	rs, err := h.Client.Runtime.GetServerState(params.Backend, params.Name)
	if err != nil {
		e := misc.HandleError(err)
		return server.NewGetRuntimeServerDefault(int(*e.Code)).WithPayload(e)
	}

	if rs == nil {
		code := int64(404)
		msg := fmt.Sprintf("Runtime server %s not found in backend %s", params.Name, params.Backend)
		return server.NewGetRuntimeServerNotFound().WithPayload(&models.Error{Code: &code, Message: &msg})
	}

	return server.NewGetRuntimeServerOK().WithPayload(rs)
}

// Handle executing the request and returning a response
func (h *GetRuntimeServersHandlerImpl) Handle(params server.GetRuntimeServersParams, principal interface{}) middleware.Responder {
	rs, err := h.Client.Runtime.GetServersState(params.Backend)
	if err != nil {
		e := misc.HandleContainerGetError(err)
		if *e.Code == misc.ErrHTTPOk {
			return server.NewGetRuntimeServersOK().WithPayload(models.RuntimeServers{})
		}
		return server.NewGetRuntimeServersDefault(int(*e.Code)).WithPayload(e)
	}

	return server.NewGetRuntimeServersOK().WithPayload(rs)
}

// Handle executing the request and returning a response
func (h *ReplaceRuntimeServerHandlerImpl) Handle(params server.ReplaceRuntimeServerParams, principal interface{}) middleware.Responder {
	rs, err := h.Client.Runtime.GetServerState(params.Backend, params.Name)
	if err != nil {
		e := misc.HandleError(err)
		return server.NewReplaceRuntimeServerDefault(int(*e.Code)).WithPayload(e)
	}

	if rs == nil {
		code := int64(404)
		msg := fmt.Sprintf("Runtime server %s not found in backend %s", params.Name, params.Backend)
		return server.NewReplaceRuntimeServerNotFound().WithPayload(&models.Error{Code: &code, Message: &msg})
	}

	// change operational state
	if params.Data.OperationalState != "" && rs.OperationalState != params.Data.OperationalState {
		err = h.Client.Runtime.SetServerHealth(params.Backend, params.Name, params.Data.OperationalState)
		if err != nil {
			e := misc.HandleError(err)
			return server.NewReplaceRuntimeServerDefault(int(*e.Code)).WithPayload(e)
		}
	}

	// change admin state
	if params.Data.AdminState != "" && rs.AdminState != params.Data.AdminState {
		err = h.Client.Runtime.SetServerState(params.Backend, params.Name, params.Data.AdminState)
		if err != nil {
			e := misc.HandleError(err)

			// try to revert operational state and fall silently
			//nolint:errcheck
			h.Client.Runtime.SetServerHealth(params.Backend, params.Name, rs.OperationalState)
			return server.NewReplaceRuntimeServerDefault(int(*e.Code)).WithPayload(e)
		}
	}

	rs, err = h.Client.Runtime.GetServerState(params.Backend, params.Name)
	if err != nil {
		e := misc.HandleError(err)
		return server.NewReplaceRuntimeServerDefault(int(*e.Code)).WithPayload(e)
	}

	return server.NewReplaceRuntimeServerOK().WithPayload(rs)
}
