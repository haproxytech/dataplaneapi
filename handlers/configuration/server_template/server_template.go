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

package server_template

import (
	"net/http"
	"reflect"

	"github.com/go-chi/chi/v5"
	client_native "github.com/haproxytech/client-native/v6"
	"github.com/haproxytech/client-native/v6/models"

	"github.com/haproxytech/dataplaneapi/handlers/middleware"
	"github.com/haproxytech/dataplaneapi/handlers/respond"
	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/log"
	"github.com/haproxytech/dataplaneapi/misc"
)

// RegisterRouter registers all server template routes onto r using spec-based request validation
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

// HandlerImpl implements ServerInterface for HAProxy server template configuration.
type HandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

func (h *HandlerImpl) GetServerTemplates(w http.ResponseWriter, r *http.Request, parentName string, params GetServerTemplatesParams) {
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	_, templates, err := cfg.GetServerTemplates(parentName, params.TransactionId)
	if err != nil {
		e := misc.HandleContainerGetError(err)
		if *e.Code == misc.ErrHTTPOk {
			respond.JSON(w, http.StatusOK, models.ServerTemplates{})
			return
		}
		respond.JSON(w, int(*e.Code), e)
		return
	}
	respond.JSON(w, http.StatusOK, templates)
}

func (h *HandlerImpl) CreateServerTemplate(w http.ResponseWriter, r *http.Request, parentName string, params CreateServerTemplateParams) {
	if params.TransactionId != "" && params.ForceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data ServerTemplate
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.CreateServerTemplate(parentName, &data, params.TransactionId, int64(params.Version)); err != nil {
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

func (h *HandlerImpl) DeleteServerTemplate(w http.ResponseWriter, r *http.Request, parentName string, prefix string, params DeleteServerTemplateParams) {
	if params.TransactionId != "" && params.ForceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.DeleteServerTemplate(prefix, parentName, params.TransactionId, int64(params.Version)); err != nil {
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

func (h *HandlerImpl) GetServerTemplate(w http.ResponseWriter, r *http.Request, parentName string, prefix string, params GetServerTemplateParams) {
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	_, tmpl, err := cfg.GetServerTemplate(prefix, parentName, params.TransactionId)
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, tmpl)
}

func (h *HandlerImpl) ReplaceServerTemplate(w http.ResponseWriter, r *http.Request, parentName string, prefix string, params ReplaceServerTemplateParams) {
	if params.TransactionId != "" && params.ForceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data ServerTemplate
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	_, ondisk, err := cfg.GetServerTemplate(prefix, parentName, params.TransactionId)
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.EditServerTemplate(prefix, parentName, &data, params.TransactionId, int64(params.Version)); err != nil {
		respond.Error(w, err)
		return
	}
	if params.TransactionId == "" {
		if serverTemplateNeedsReload(data, *ondisk, h.Client) {
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

// serverTemplateNeedsReload reports whether replacing a server template requires a
// reload. Server templates have no runtime-API-editable fields, so any change to an
// existing template requires a reload; an identical body (no diff) does not. This
// mirrors the pre-migration changeThroughRuntimeAPI behavior for the ServerTemplate type.
func serverTemplateNeedsReload(data, ondisk models.ServerTemplate, client client_native.HAProxyClient) (reload bool) {
	// reflect kinds and values are loosely checked as they are bound strictly in
	// schema, but in case of any panic, we log and reload to ensure changes go through.
	defer func() {
		if r := recover(); r != nil {
			log.Warning("serverTemplateNeedsReload panic:", r)
			reload = true
		}
	}()
	if _, err := client.Runtime(); err != nil {
		return true
	}
	return len(serverTemplateCompareObjects(data, ondisk)) != 0
}

// serverTemplateCompareObjects returns the names of fields that differ between data and ondisk.
func serverTemplateCompareObjects(data, ondisk any) []string {
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
			diff = append(diff, serverTemplateCompareObjects(dField.Interface(), oField.Interface())...)
		}
	}
	return diff
}
