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

package reloads

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/haproxytech/client-native/v6/models"

	"github.com/haproxytech/dataplaneapi/handlers/middleware"
	"github.com/haproxytech/dataplaneapi/handlers/respond"
	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/misc"
)

// RegisterRouter registers all reload routes onto r using spec-based request validation
// and a shared error handler.
func RegisterRouter(r chi.Router, ra haproxy.IReloadAgent) error {
	spec, err := GetSpec()
	if err != nil {
		return err
	}
	HandlerWithOptions(&HandlerImpl{ReloadAgent: ra}, ChiServerOptions{
		BaseRouter:       r,
		Middlewares:      []MiddlewareFunc{middleware.NewValidator(spec)},
		ErrorHandlerFunc: middleware.ErrorHandler,
	})
	return nil
}

// HandlerImpl implements ServerInterface for HAProxy reload operations.
type HandlerImpl struct {
	ReloadAgent haproxy.IReloadAgent
}

func (h *HandlerImpl) GetReloads(w http.ResponseWriter, r *http.Request) {
	rs := h.ReloadAgent.GetReloads()
	respond.JSON(w, http.StatusOK, rs)
}

func (h *HandlerImpl) GetReload(w http.ResponseWriter, r *http.Request, id string) {
	reload := h.ReloadAgent.GetReload(id)
	if reload == nil {
		c := misc.ErrHTTPNotFound
		msg := fmt.Sprintf("Reload with ID %s does not exist", id)
		respond.JSON(w, http.StatusNotFound, &models.Error{
			Code:    &c,
			Message: &msg,
		})
		return
	}
	respond.JSON(w, http.StatusOK, reload)
}
