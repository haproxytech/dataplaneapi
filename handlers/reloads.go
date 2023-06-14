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
	"github.com/haproxytech/client-native/v5/models"

	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/operations/reloads"
)

// GetReloadHandlerImpl implementation of the GetReloadHandler interface using client-native client
type GetReloadHandlerImpl struct {
	ReloadAgent haproxy.IReloadAgent
}

// GetReloadsHandlerImpl implementation of the GetReloadsHandler interface using client-native client
type GetReloadsHandlerImpl struct {
	ReloadAgent haproxy.IReloadAgent
}

// Handle executing the request and returning a response
func (rh *GetReloadHandlerImpl) Handle(params reloads.GetReloadParams, principal interface{}) middleware.Responder {
	r := rh.ReloadAgent.GetReload(params.ID)
	if r == nil {
		msg := fmt.Sprintf("Reload with ID %s does not exist", params.ID)
		c := misc.ErrHTTPNotFound
		e := &models.Error{
			Code:    &c,
			Message: &msg,
		}
		return reloads.NewGetReloadDefault(404).WithPayload(e)
	}
	return reloads.NewGetReloadOK().WithPayload(r)
}

// Handle executing the request and returning a response
func (rh *GetReloadsHandlerImpl) Handle(params reloads.GetReloadsParams, principal interface{}) middleware.Responder {
	rs := rh.ReloadAgent.GetReloads()
	return reloads.NewGetReloadsOK().WithPayload(rs)
}
