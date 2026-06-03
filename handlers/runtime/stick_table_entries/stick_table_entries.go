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

package stick_table_entries

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-openapi/strfmt"
	client_native "github.com/haproxytech/client-native/v6"
	"github.com/haproxytech/client-native/v6/models"

	"github.com/haproxytech/dataplaneapi/handlers/middleware"
	"github.com/haproxytech/dataplaneapi/handlers/respond"
	"github.com/haproxytech/dataplaneapi/misc"
)

// RegisterRouter registers all stick_table_entries routes onto r using spec-based request validation.
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

// HandlerImpl implements ServerInterface for HAProxy stick table entries.
type HandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h *HandlerImpl) GetStickTableEntries(w http.ResponseWriter, r *http.Request, parentName string, params GetStickTableEntriesParams) {
	filter := make([]string, 0)
	if params.Filter != "" {
		filter = strings.Split(params.Filter, ",")
	}

	rt, err := h.Client.Runtime()
	if err != nil {
		respond.Error(w, err)
		return
	}

	stkEntries, err := rt.GetTableEntries(parentName, filter, params.Key)
	if err != nil {
		respond.Error(w, err)
		return
	}

	// if no entries return empty array
	if len(stkEntries) == 0 {
		respond.JSON(w, http.StatusOK, stkEntries)
		return
	}

	// check for pagination
	offset := params.Offset
	if offset >= len(stkEntries) {
		msg := fmt.Sprintf("Offset %d is larger than the slice size %d", offset, len(stkEntries))
		c := misc.ErrHTTPBadRequest
		respond.JSON(w, int(c), &models.Error{
			Message: &msg,
			Code:    &c,
		})
		return
	}

	if params.Count != nil {
		if offset+*params.Count >= len(stkEntries) {
			stkEntries = stkEntries[offset:]
		} else {
			stkEntries = stkEntries[offset : offset+*params.Count]
		}
	} else {
		stkEntries = stkEntries[offset:]
	}
	respond.JSON(w, http.StatusOK, stkEntries)
}

func (h *HandlerImpl) SetStickTableEntries(w http.ResponseWriter, r *http.Request, parentName string) {
	var data SetStickTableEntriesJSONRequestBody
	if !respond.DecodeJSON(r, w, &data) {
		return
	}
	if err := data.DataType.Validate(strfmt.Default); err != nil {
		respond.Unprocessable(w, err)
		return
	}

	rt, err := h.Client.Runtime()
	if err != nil {
		respond.Error(w, err)
		return
	}

	if err = rt.SetTableEntry(parentName, data.Key, data.DataType); err != nil {
		respond.Error(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
