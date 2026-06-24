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

package server

import (
	"net/http"
	"reflect"
	"strconv"

	"github.com/go-chi/chi/v5"
	client_native "github.com/haproxytech/client-native/v6"
	"github.com/haproxytech/client-native/v6/models"
	cn_runtime "github.com/haproxytech/client-native/v6/runtime"

	"github.com/haproxytech/dataplaneapi/handlers/middleware"
	"github.com/haproxytech/dataplaneapi/handlers/respond"
	runtime_servers "github.com/haproxytech/dataplaneapi/handlers/runtime/servers"
	"github.com/haproxytech/dataplaneapi/log"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/reload_agent"
)

// RegisterRouter registers all server routes onto r using spec-based request validation
// and a shared error handler.
func RegisterRouter(r chi.Router, client client_native.HAProxyClient, ra reload_agent.IReloadAgent) error {
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

// HandlerImpl implements ServerInterface for HAProxy server configuration.
type HandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent reload_agent.IReloadAgent
}

// --- Backend ---

func (h *HandlerImpl) GetAllServerBackend(w http.ResponseWriter, r *http.Request, parentName string, params GetAllServerBackendParams) {
	h.getAllServer(w, r, "backend", parentName, params.TransactionId)
}

func (h *HandlerImpl) CreateServerBackend(w http.ResponseWriter, r *http.Request, parentName string, params CreateServerBackendParams) {
	h.createServer(w, r, "backend", parentName, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) DeleteServerBackend(w http.ResponseWriter, r *http.Request, parentName string, name string, params DeleteServerBackendParams) {
	h.deleteServer(w, r, "backend", parentName, name, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) GetServerBackend(w http.ResponseWriter, r *http.Request, parentName string, name string, params GetServerBackendParams) {
	h.getServer(w, r, "backend", parentName, name, params.TransactionId)
}

func (h *HandlerImpl) ReplaceServerBackend(w http.ResponseWriter, r *http.Request, parentName string, name string, params ReplaceServerBackendParams) {
	h.replaceServer(w, r, "backend", parentName, name, params.TransactionId, int64(params.Version), params.ForceReload)
}

// --- Peer ---

func (h *HandlerImpl) GetAllServerPeer(w http.ResponseWriter, r *http.Request, parentName string, params GetAllServerPeerParams) {
	h.getAllServer(w, r, "peers", parentName, params.TransactionId)
}

func (h *HandlerImpl) CreateServerPeer(w http.ResponseWriter, r *http.Request, parentName string, params CreateServerPeerParams) {
	h.createServer(w, r, "peers", parentName, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) DeleteServerPeer(w http.ResponseWriter, r *http.Request, parentName string, name string, params DeleteServerPeerParams) {
	h.deleteServer(w, r, "peers", parentName, name, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) GetServerPeer(w http.ResponseWriter, r *http.Request, parentName string, name string, params GetServerPeerParams) {
	h.getServer(w, r, "peers", parentName, name, params.TransactionId)
}

func (h *HandlerImpl) ReplaceServerPeer(w http.ResponseWriter, r *http.Request, parentName string, name string, params ReplaceServerPeerParams) {
	h.replaceServer(w, r, "peers", parentName, name, params.TransactionId, int64(params.Version), params.ForceReload)
}

// --- Ring ---

func (h *HandlerImpl) GetAllServerRing(w http.ResponseWriter, r *http.Request, parentName string, params GetAllServerRingParams) {
	h.getAllServer(w, r, "ring", parentName, params.TransactionId)
}

func (h *HandlerImpl) CreateServerRing(w http.ResponseWriter, r *http.Request, parentName string, params CreateServerRingParams) {
	h.createServer(w, r, "ring", parentName, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) DeleteServerRing(w http.ResponseWriter, r *http.Request, parentName string, name string, params DeleteServerRingParams) {
	h.deleteServer(w, r, "ring", parentName, name, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) GetServerRing(w http.ResponseWriter, r *http.Request, parentName string, name string, params GetServerRingParams) {
	h.getServer(w, r, "ring", parentName, name, params.TransactionId)
}

func (h *HandlerImpl) ReplaceServerRing(w http.ResponseWriter, r *http.Request, parentName string, name string, params ReplaceServerRingParams) {
	h.replaceServer(w, r, "ring", parentName, name, params.TransactionId, int64(params.Version), params.ForceReload)
}

// --- Shared implementations ---

func (h *HandlerImpl) getAllServer(w http.ResponseWriter, r *http.Request, parentType, parentName, txID string) {
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	_, srvs, err := cfg.GetServers(parentType, parentName, txID)
	if err != nil {
		e := misc.HandleContainerGetError(err)
		if *e.Code == misc.ErrHTTPOk {
			respond.JSON(w, http.StatusOK, Servers{})
			return
		}
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, srvs)
}

func (h *HandlerImpl) getServer(w http.ResponseWriter, r *http.Request, parentType, parentName, name, txID string) {
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	_, srv, err := cfg.GetServer(name, parentType, parentName, txID)
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, srv)
}

func (h *HandlerImpl) createServer(w http.ResponseWriter, r *http.Request, parentType, parentName, txID string, version int64, forceReload bool) {
	if txID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data Server
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.CreateServer(parentType, parentName, &data, txID, version); err != nil {
		respond.Error(w, err)
		return
	}

	// Try to add the server dynamically via the runtime API (backend only, no default_server).
	// On success we skip the reload entirely.
	useRuntime := false
	var ras models.RuntimeAddServer
	var runtimeClient cn_runtime.Runtime
	if parentType == "backend" && txID == "" {
		_, defaults, errRT := cfg.GetDefaultsConfiguration(txID)
		if errRT != nil {
			respond.Error(w, errRT)
			return
		}
		_, backend, errRT := cfg.GetBackend(parentName, txID)
		if errRT != nil {
			respond.Error(w, errRT)
			return
		}
		runtimeClient, errRT = h.Client.Runtime()
		if errRT == nil && defaults.DefaultServer == nil && backend.DefaultServer == nil {
			useRuntime = misc.ConvertStruct(data, &ras) == nil
		}
	}

	if txID == "" {
		if forceReload {
			if err = h.ReloadAgent.ForceReload(); err != nil {
				respond.Error(w, err)
				return
			}
			respond.JSON(w, http.StatusCreated, &data)
			return
		}
		if useRuntime {
			haversion, _ := runtimeClient.GetVersion()
			if err = runtimeClient.AddServer(parentName, data.Name, runtime_servers.SerializeRuntimeAddServer(&ras, &haversion)); err == nil {
				log.Debugf("backend %s: server %s added through runtime", parentName, data.Name)
				respond.JSON(w, http.StatusCreated, &data)
				return
			}
			log.Warning("failed to add server through runtime:", err)
		}
		respond.Accepted(w, h.ReloadAgent.Reload(), &data)
		return
	}
	respond.JSON(w, http.StatusAccepted, &data)
}

func (h *HandlerImpl) deleteServer(w http.ResponseWriter, r *http.Request, parentType, parentName, name, txID string, version int64, forceReload bool) {
	if txID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.DeleteServer(name, parentType, parentName, txID, version); err != nil {
		respond.Error(w, err)
		return
	}
	if txID == "" {
		if forceReload {
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

func (h *HandlerImpl) replaceServer(w http.ResponseWriter, r *http.Request, parentType, parentName, name, txID string, version int64, forceReload bool) {
	if txID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data Server
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	_, ondisk, err := cfg.GetServer(name, parentType, parentName, txID)
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.EditServer(name, parentType, parentName, &data, txID, version); err != nil {
		respond.Error(w, err)
		return
	}
	if txID == "" {
		needsReload := serverChangeThroughRuntimeAPI(data, *ondisk, parentName, h.Client)
		if needsReload {
			if forceReload {
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

// runtimeSupportedServerFields lists Server fields that can be changed without a reload.
var runtimeSupportedServerFields = []string{"Weight", "Address", "Port", "Maintenance", "AgentCheck", "AgentAddr", "AgentSend", "HealthCheckPort"}

// serverChangeThroughRuntimeAPI attempts to apply server field changes via the HAProxy runtime API.
// Returns true when a reload is still required (unsupported fields changed or runtime call failed).
func serverChangeThroughRuntimeAPI(data, ondisk models.Server, parentName string, client client_native.HAProxyClient) (reload bool) {
	defer func() {
		if r := recover(); r != nil {
			log.Warning("serverChangeThroughRuntimeAPI panic:", r)
			reload = true
		}
	}()

	rt, err := client.Runtime()
	if err != nil {
		return true
	}

	diff := serverCompareObjects(data, ondisk)
	if len(diff) == 0 {
		return false
	}
	if !serverCompareChanged(diff, runtimeSupportedServerFields) {
		return true
	}

	addrPortChanged := false
	for _, field := range diff {
		fieldValue := reflect.ValueOf(data).FieldByName(field)
		if !fieldValue.IsValid() {
			continue
		}
		switch field {
		case "Weight":
			weight := strconv.FormatInt(fieldValue.Elem().Int(), 10)
			if err = rt.SetServerWeight(parentName, data.Name, weight); err != nil {
				return true
			}
		case "HealthCheckPort":
			port := int(fieldValue.Elem().Int())
			if err = rt.SetServerCheckPort(parentName, data.Name, port); err != nil {
				return true
			}
		case "Address":
			if !addrPortChanged {
				portVal := int(reflect.ValueOf(data).FieldByName("Port").Elem().Int())
				if err = rt.SetServerAddr(parentName, data.Name, fieldValue.String(), portVal); err != nil {
					return true
				}
				addrPortChanged = true
			}
		case "Port":
			if !addrPortChanged {
				port := int(fieldValue.Elem().Int())
				addr := reflect.ValueOf(data).FieldByName("Address").String()
				if err = rt.SetServerAddr(parentName, data.Name, addr, port); err != nil {
					return true
				}
				addrPortChanged = true
			}
		case "Maintenance":
			state := "ready"
			if fieldValue.String() == "enabled" {
				state = "maint"
			}
			if err = rt.SetServerState(parentName, data.Name, state); err != nil {
				return true
			}
		case "AgentCheck":
			switch fieldValue.String() {
			case "enabled":
				if err = rt.EnableAgentCheck(parentName, data.Name); err != nil {
					return true
				}
			case "disabled":
				if err = rt.DisableAgentCheck(parentName, data.Name); err != nil {
					return true
				}
			}
		case "AgentAddr":
			if err = rt.SetServerAgentAddr(parentName, data.Name, fieldValue.String()); err != nil {
				return true
			}
		case "AgentSend":
			if err = rt.SetServerAgentSend(parentName, data.Name, fieldValue.String()); err != nil {
				return true
			}
		}
	}
	return false
}

func serverCompareObjects(data, ondisk any) []string {
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
			diff = append(diff, serverCompareObjects(dField.Interface(), oField.Interface())...)
		}
	}
	return diff
}

func serverCompareChanged(changed, changeable []string) bool {
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
