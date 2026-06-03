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

package discovery

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/haproxytech/dataplaneapi/handlers/middleware"
	"github.com/haproxytech/dataplaneapi/handlers/respond"
	"github.com/haproxytech/dataplaneapi/misc"
)

// RegisterRouter registers all discovery routes onto r using spec-based request validation
// and a shared error handler. swaggerJSON is the raw OpenAPI v2 spec used to enumerate
// child paths for the discovery endpoints.
func RegisterRouter(r chi.Router, swaggerJSON json.RawMessage) error {
	spec, err := GetSpec()
	if err != nil {
		return err
	}
	HandlerWithOptions(&HandlerImpl{swaggerJSON: swaggerJSON}, ChiServerOptions{
		BaseRouter:       r,
		Middlewares:      []MiddlewareFunc{middleware.NewValidator(spec)},
		ErrorHandlerFunc: middleware.ErrorHandler,
	})
	return nil
}

// HandlerImpl implements ServerInterface for the API discovery endpoints.
type HandlerImpl struct {
	swaggerJSON json.RawMessage
}

func (h *HandlerImpl) GetAPIEndpoints(w http.ResponseWriter, r *http.Request) {
	ends, err := misc.DiscoverChildPaths("", h.swaggerJSON)
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, ends)
}

func (h *HandlerImpl) GetServicesEndpoints(w http.ResponseWriter, r *http.Request) {
	rURI := "/" + strings.SplitN(r.RequestURI[1:], "/", 2)[1]
	ends, err := misc.DiscoverChildPaths(rURI, h.swaggerJSON)
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, ends)
}

func (h *HandlerImpl) GetHaproxyEndpoints(w http.ResponseWriter, r *http.Request) {
	rURI := "/" + strings.SplitN(r.RequestURI[1:], "/", 2)[1]
	ends, err := misc.DiscoverChildPaths(rURI, h.swaggerJSON)
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, ends)
}

func (h *HandlerImpl) GetConfigurationEndpoints(w http.ResponseWriter, r *http.Request) {
	rURI := "/" + strings.SplitN(r.RequestURI[1:], "/", 2)[1]
	ends, err := misc.DiscoverChildPaths(rURI, h.swaggerJSON)
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, ends)
}

func (h *HandlerImpl) GetRuntimeEndpoints(w http.ResponseWriter, r *http.Request) {
	rURI := "/" + strings.SplitN(r.RequestURI[1:], "/", 2)[1]
	ends, err := misc.DiscoverChildPaths(rURI, h.swaggerJSON)
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, ends)
}

func (h *HandlerImpl) GetStatsEndpoints(w http.ResponseWriter, r *http.Request) {
	rURI := "/" + strings.SplitN(r.RequestURI[1:], "/", 2)[1]
	ends, err := misc.DiscoverChildPaths(rURI, h.swaggerJSON)
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, ends)
}

func (h *HandlerImpl) GetSpoeEndpoints(w http.ResponseWriter, r *http.Request) {
	rURI := "/" + strings.SplitN(r.RequestURI[1:], "/", 2)[1]
	ends, err := misc.DiscoverChildPaths(rURI, h.swaggerJSON)
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, ends)
}

func (h *HandlerImpl) GetStorageEndpoints(w http.ResponseWriter, r *http.Request) {
	rURI := "/" + strings.SplitN(r.RequestURI[1:], "/", 2)[1]
	ends, err := misc.DiscoverChildPaths(rURI, h.swaggerJSON)
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, ends)
}
