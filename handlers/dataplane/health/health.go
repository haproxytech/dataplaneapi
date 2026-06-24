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

package health

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/haproxytech/client-native/v6/models"

	"github.com/haproxytech/dataplaneapi/handlers/middleware"
	"github.com/haproxytech/dataplaneapi/handlers/respond"
	"github.com/haproxytech/dataplaneapi/reload_agent"
)

// RegisterRouter registers all health routes onto r using spec-based request validation
// and a shared error handler.
func RegisterRouter(r chi.Router, ra reload_agent.IReloadAgent) error {
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

// HandlerImpl implements ServerInterface for HAProxy health checks.
type HandlerImpl struct {
	ReloadAgent reload_agent.IReloadAgent
}

func (h *HandlerImpl) GetHealth(w http.ResponseWriter, r *http.Request) {
	data := models.Health{}
	status, err := h.ReloadAgent.Status()
	if err == nil {
		if status {
			data.Haproxy = models.HealthHaproxyUp
		} else {
			data.Haproxy = models.HealthHaproxyDown
		}
	} else {
		data.Haproxy = models.HealthHaproxyUnknown
	}
	respond.JSON(w, http.StatusOK, &data)
}
