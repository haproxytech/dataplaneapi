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
	"github.com/haproxytech/client-native"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/operations/information"
	"github.com/haproxytech/models"
)

//GetInformationHandlerImpl implementation of the GetInformationHandler interface using client-native client
type GetInformationHandlerImpl struct {
	Client *client_native.HAProxyClient
}

//Handle executing the request and returning a response
func (h *GetInformationHandlerImpl) Handle(params information.GetHaproxyProcessInfoParams, principal interface{}) middleware.Responder {
	info, err := h.Client.Runtime.GetInfo()
	if err != nil || len(info) == 0 {
		code := misc.ErrHTTPInternalServerError
		msg := err.Error()
		e := &models.Error{
			Code:    &code,
			Message: &msg,
		}
		return information.NewGetHaproxyProcessInfoDefault(int(misc.ErrHTTPInternalServerError)).WithPayload(e)
	}

	data := models.ProcessInfo{}
	data.Haproxy = &info[0]

	return information.NewGetHaproxyProcessInfoOK().WithPayload(&data)
}
