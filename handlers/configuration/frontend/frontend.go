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

package frontend

import (
	"net/http"
	"reflect"

	"github.com/go-chi/chi/v5"
	client_native "github.com/haproxytech/client-native/v6"

	"github.com/haproxytech/dataplaneapi/handlers/middleware"
	"github.com/haproxytech/dataplaneapi/handlers/respond"
	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/log"
)

// RegisterRouter registers all frontend routes onto r using spec-based request validation
// and a shared error handler.
func RegisterRouter(r chi.Router, client client_native.HAProxyClient, ra haproxy.IReloadAgent) error {
	spec, err := GetSpec()
	if err != nil {
		return err
	}
	HandlerWithOptions(&HandlerImpl{Client: client, ReloadAgent: ra}, ChiServerOptions{
		BaseRouter:       r,
		Middlewares:      []MiddlewareFunc{middleware.NewValidator(spec)},
		ErrorHandlerFunc: middleware.ErrorHandler,
	})
	return nil
}

// HandlerImpl implements ServerInterface for HAProxy frontend configuration.
type HandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

func (h *HandlerImpl) GetFrontends(w http.ResponseWriter, r *http.Request, params GetFrontendsParams) {
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	var frontends Frontends
	if params.FullSection {
		_, frontends, err = cfg.GetStructuredFrontends(params.TransactionId)
	} else {
		_, frontends, err = cfg.GetFrontends(params.TransactionId)
	}
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, frontends)
}

func (h *HandlerImpl) CreateFrontend(w http.ResponseWriter, r *http.Request, params CreateFrontendParams) {
	if params.TransactionId != "" && params.ForceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data Frontend
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if params.FullSection {
		err = cfg.CreateStructuredFrontend(&data, params.TransactionId, int64(params.Version))
	} else {
		err = cfg.CreateFrontend(&data, params.TransactionId, int64(params.Version))
	}
	if err != nil {
		respond.Error(w, err)
		return
	}
	if params.TransactionId == "" {
		if params.ForceReload {
			if err = h.ReloadAgent.ForceReload(); err != nil {
				respond.Error(w, err)
				return
			}
			respond.JSON(w, http.StatusCreated, &data)
			return
		}
		respond.Accepted(w, h.ReloadAgent.Reload(), &data)
		return
	}
	respond.JSON(w, http.StatusAccepted, &data)
}

func (h *HandlerImpl) DeleteFrontend(w http.ResponseWriter, r *http.Request, name string, params DeleteFrontendParams) {
	if params.TransactionId != "" && params.ForceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.DeleteFrontend(name, params.TransactionId, int64(params.Version)); err != nil {
		respond.Error(w, err)
		return
	}
	if params.TransactionId == "" {
		if params.ForceReload {
			if err = h.ReloadAgent.ForceReload(); err != nil {
				respond.Error(w, err)
				return
			}
			respond.NoContent(w)
			return
		}
		respond.Accepted(w, h.ReloadAgent.Reload(), nil)
		return
	}
	respond.Accepted(w, "", nil)
}

func (h *HandlerImpl) GetFrontend(w http.ResponseWriter, r *http.Request, name string, params GetFrontendParams) {
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	var fe *Frontend
	if params.FullSection {
		_, fe, err = cfg.GetStructuredFrontend(name, params.TransactionId)
	} else {
		_, fe, err = cfg.GetFrontend(name, params.TransactionId)
	}
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, fe)
}

func (h *HandlerImpl) ReplaceFrontend(w http.ResponseWriter, r *http.Request, name string, params ReplaceFrontendParams) {
	if params.TransactionId != "" && params.ForceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data Frontend
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	_, ondisk, err := cfg.GetFrontend(name, params.TransactionId)
	if err != nil {
		respond.Error(w, err)
		return
	}
	if params.FullSection {
		err = cfg.EditStructuredFrontend(name, &data, params.TransactionId, int64(params.Version))
	} else {
		err = cfg.EditFrontend(name, &data, params.TransactionId, int64(params.Version))
	}
	if err != nil {
		respond.Error(w, err)
		return
	}
	if params.TransactionId == "" {
		needsReload := frontendChangeThroughRuntimeAPI(data, *ondisk, h.Client)
		if needsReload {
			if params.ForceReload {
				if err = h.ReloadAgent.ForceReload(); err != nil {
					respond.Error(w, err)
					return
				}
				respond.JSON(w, http.StatusOK, &data)
				return
			}
			respond.Accepted(w, h.ReloadAgent.Reload(), &data)
			return
		}
		respond.JSON(w, http.StatusOK, &data)
		return
	}
	respond.JSON(w, http.StatusAccepted, &data)
}

// runtimeSupportedFrontendFields lists Frontend fields that can be changed without a reload.
var runtimeSupportedFrontendFields = []string{"Maxconn"}

// frontendChangeThroughRuntimeAPI attempts to apply frontend field changes via the HAProxy runtime API.
// Returns true when a reload is still required (unsupported fields changed or runtime call failed).
func frontendChangeThroughRuntimeAPI(data, ondisk Frontend, client client_native.HAProxyClient) (reload bool) {
	defer func() {
		if r := recover(); r != nil {
			log.Warning("frontendChangeThroughRuntimeAPI panic:", r)
			reload = true
		}
	}()

	rt, err := client.Runtime()
	if err != nil {
		return true
	}

	diff := frontendCompareObjects(data, ondisk)
	if len(diff) == 0 {
		return false
	}
	if !frontendCompareChanged(diff, runtimeSupportedFrontendFields) {
		return true
	}

	for _, field := range diff {
		fieldValue := reflect.ValueOf(data).FieldByName(field)
		if !fieldValue.IsValid() {
			continue
		}
		switch field {
		case "Maxconn":
			maxConn := int(fieldValue.Elem().Int())
			if err = rt.SetFrontendMaxConn(data.Name, maxConn); err != nil {
				return true
			}
		}
	}
	return false
}

func frontendCompareObjects(data, ondisk any) []string {
	diff := []string{}
	dataVal := reflect.ValueOf(data)
	ondiskVal := reflect.ValueOf(ondisk)
	for i := range dataVal.NumField() {
		fName := dataVal.Type().Field(i).Name
		dField := dataVal.FieldByName(fName)
		oField := ondiskVal.FieldByName(fName)
		dKind := dField.Kind()
		oKind := oField.Kind()
		if dKind != oKind {
			diff = append(diff, fName)
			continue
		}
		if dKind == reflect.Ptr {
			dField = dField.Elem()
			oField = oField.Elem()
			dKind = dField.Kind()
			oKind = oField.Kind()
			if dKind != oKind {
				diff = append(diff, fName)
				continue
			}
		}
		switch dKind {
		case reflect.Float32, reflect.Float64:
			if dField.Float() != oField.Float() {
				diff = append(diff, fName)
			}
		case reflect.Bool:
			if dField.Bool() != oField.Bool() {
				diff = append(diff, fName)
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if dField.Int() != oField.Int() {
				diff = append(diff, fName)
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if dField.Uint() != oField.Uint() {
				diff = append(diff, fName)
			}
		case reflect.String:
			if dField.String() != oField.String() {
				diff = append(diff, fName)
			}
		case reflect.Struct:
			diff = append(diff, frontendCompareObjects(dField.Interface(), oField.Interface())...)
		}
	}
	return diff
}

func frontendCompareChanged(changed, changeable []string) bool {
	if len(changed) > len(changeable) {
		return false
	}
	set := make(map[string]bool, len(changed))
	for _, f := range changed {
		set[f] = true
	}
	for _, f := range changeable {
		delete(set, f)
	}
	return len(set) == 0
}
