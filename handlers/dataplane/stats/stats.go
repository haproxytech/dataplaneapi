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

package stats

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	client_native "github.com/haproxytech/client-native/v6"
	"github.com/haproxytech/client-native/v6/models"

	"github.com/haproxytech/dataplaneapi/handlers/middleware"
	"github.com/haproxytech/dataplaneapi/handlers/respond"
)

// RegisterRouter registers all stats routes onto r using spec-based request validation
// and a shared error handler.
func RegisterRouter(r chi.Router, client client_native.HAProxyClient) error {
	spec, err := GetSpec()
	if err != nil {
		return err
	}
	HandlerWithOptions(&HandlerImpl{Client: client}, ChiServerOptions{
		BaseRouter:       r,
		Middlewares:      []MiddlewareFunc{middleware.NewValidator(spec)},
		ErrorHandlerFunc: middleware.ErrorHandler,
	})
	return nil
}

// HandlerImpl implements ServerInterface for HAProxy stats.
type HandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h *HandlerImpl) GetStats(w http.ResponseWriter, r *http.Request, params GetStatsParams) {
	if params.Name != "" {
		switch params.Type {
		case "":
			respond.BadRequest(w, "Type required when filtering by name")
			return
		case Server:
			if params.Parent == "" {
				respond.BadRequest(w, "Parent backend required when filtering by server")
				return
			}
		}
	}

	rt, err := h.Client.Runtime()
	if err != nil {
		respond.Error(w, err)
		return
	}

	nStat := rt.GetStats()

	if nStat.Error != "" {
		respond.JSON(w, http.StatusInternalServerError, &nStat)
		return
	}

	retVal := make([]*models.NativeStat, 0, len(nStat.Stats))
	for _, item := range nStat.Stats {
		if params.Name != "" {
			if item.Type == "server" {
				if item.Name == params.Name && item.Type == string(params.Type) && item.BackendName == params.Parent {
					retVal = append(retVal, item)
				}
			} else if item.Name == params.Name && item.Type == string(params.Type) {
				retVal = append(retVal, item)
			}
		} else {
			if params.Type != "" {
				if params.Type == Server && params.Parent != "" {
					if item.Type == string(params.Type) && item.BackendName == params.Parent {
						retVal = append(retVal, item)
					}
				} else {
					if item.Type == string(params.Type) {
						retVal = append(retVal, item)
					}
				}
			} else {
				retVal = append(retVal, item)
			}
		}
	}
	nStat.Stats = retVal

	respond.JSON(w, http.StatusOK, &nStat)
}
