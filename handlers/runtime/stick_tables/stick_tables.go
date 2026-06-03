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

package stick_tables

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	client_native "github.com/haproxytech/client-native/v6"
	"github.com/haproxytech/client-native/v6/models"

	"github.com/haproxytech/dataplaneapi/handlers/middleware"
	"github.com/haproxytech/dataplaneapi/handlers/respond"
	"github.com/haproxytech/dataplaneapi/misc"
)

// RegisterRouter registers all stick_tables routes onto r using spec-based request validation.
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

// HandlerImpl implements ServerInterface for HAProxy stick tables.
type HandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h *HandlerImpl) GetStickTables(w http.ResponseWriter, r *http.Request) {
	rt, err := h.Client.Runtime()
	if err != nil {
		respond.Error(w, err)
		return
	}
	stkTS, err := rt.ShowTables()
	if err != nil {
		respond.Error(w, err)
		return
	}
	for _, table := range stkTS {
		table.Fields = findTableFields(table.Name, h.Client)
	}
	respond.JSON(w, http.StatusOK, stkTS)
}

func (h *HandlerImpl) GetStickTable(w http.ResponseWriter, r *http.Request, name string) {
	rt, err := h.Client.Runtime()
	if err != nil {
		respond.Error(w, err)
		return
	}
	stkT, err := rt.ShowTable(name)
	if stkT == nil {
		msg := fmt.Sprintf("Stick table %s not found", name)
		c := misc.ErrHTTPNotFound
		respond.JSON(w, int(c), &models.Error{
			Message: &msg,
			Code:    &c,
		})
		return
	}
	stkT.Fields = findTableFields(stkT.Name, h.Client)
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, stkT)
}

func findTableFields(name string, client client_native.HAProxyClient) []*models.StickTableField {
	configuration, err := client.Configuration()
	if err != nil {
		return nil
	}
	_, bck, err := configuration.GetBackend(name, "")
	if err != nil {
		return nil
	}
	if bck.StickTable == nil {
		return nil
	}
	data := strings.Split(bck.StickTable.Store, ",")
	fields := make([]*models.StickTableField, 0)
	for _, d := range data {
		f := &models.StickTableField{}
		spl := strings.Split(d, "(")
		if len(spl) == 1 {
			f.Field = d
			f.Type = "counter"
			fields = append(fields, f)
		} else if len(spl) == 2 {
			p := misc.ParseTimeout(spl[1][:len(spl[1])-1])
			if p != nil {
				f.Field = spl[0]
				f.Period = *p
				f.Type = "rate"
				fields = append(fields, f)
			}
		}
	}
	return fields
}
