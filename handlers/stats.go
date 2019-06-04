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
	"github.com/haproxytech/client-native"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/operations/stats"
	"github.com/haproxytech/models"
)

//GetStatsHandlerImpl implementation of the GetStatsHandler interface using client-native client
type GetStatsHandlerImpl struct {
	Client *client_native.HAProxyClient
}

//Handle executing the request and returning a response
func (h *GetStatsHandlerImpl) Handle(params stats.GetStatsParams, principal interface{}) middleware.Responder {
	if params.Name != nil {
		if params.Type == nil {
			code := misc.ErrHTTPBadRequest
			msg := "Type required when filtering by name"
			e := &models.Error{Code: &code, Message: &msg}
			return stats.NewGetStatsDefault(int(code)).WithPayload(e)
		} else if *params.Type == "server" {
			if params.Parent == nil {
				code := misc.ErrHTTPBadRequest
				msg := "Parent backend required when filtering by server"
				e := &models.Error{Code: &code, Message: &msg}
				return stats.NewGetStatsDefault(int(code)).WithPayload(e)
			}
		}
	}

	s, err := h.Client.Runtime.GetStats()
	if err != nil {
		code := misc.ErrHTTPInternalServerError
		msg := err.Error()
		e := &models.Error{
			Code:    &code,
			Message: &msg,
		}
		return stats.NewGetStatsDefault(int(misc.ErrHTTPInternalServerError)).WithPayload(e)
	}

	nativeStats := models.NativeStats{}
	for _, nStat := range s {
		for _, item := range nStat {
			nativeStatItem := *item
			if params.Name != nil {
				if item.Type == "server" {
					if item.Name == *params.Name && item.Type == *params.Type && item.BackendName == *params.Parent {
						nativeStats = append(nativeStats, &nativeStatItem)
					}
				} else if item.Name == *params.Name && item.Type == *params.Type {
					nativeStats = append(nativeStats, &nativeStatItem)
				}
			} else {
				if params.Type != nil {
					if *params.Type == "server" && params.Parent != nil {
						if item.Type == *params.Type && item.BackendName == *params.Parent {
							nativeStats = append(nativeStats, &nativeStatItem)
						}
					} else {
						if item.Type == *params.Type {
							nativeStats = append(nativeStats, &nativeStatItem)
						}
					}
				} else {
					nativeStats = append(nativeStats, &nativeStatItem)
				}
			}
		}
	}
	return stats.NewGetStatsOK().WithPayload(nativeStats)
}
