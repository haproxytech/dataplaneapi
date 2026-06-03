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

package specification

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/haproxytech/dataplaneapi/handlers/middleware"
	"github.com/haproxytech/dataplaneapi/handlers/respond"
)

// RegisterRouter registers all specification routes onto r using spec-based request validation
// and a shared error handler. swaggerJSON is the raw OpenAPI v2 spec served at the /v2 endpoint
// for backward compatibility.
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

// HandlerImpl implements ServerInterface for the API specification endpoints.
type HandlerImpl struct {
	// swaggerJSON is the raw OpenAPI v2 / swagger spec served at the legacy endpoint
	// for backward compatibility with existing clients.
	swaggerJSON json.RawMessage
}

func (h *HandlerImpl) GetSpecification(w http.ResponseWriter, r *http.Request) {
	b, err := h.swaggerJSON.MarshalJSON()
	if err != nil {
		respond.Error(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	respond.Write(w, b)
}

func (h *HandlerImpl) GetOpenapiv3Specification(w http.ResponseWriter, r *http.Request) {
	b, err := GetDataplaneSpecJSON()
	if err != nil {
		respond.Error(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	respond.Write(w, b)
}
