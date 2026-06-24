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

package raw

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	client_native "github.com/haproxytech/client-native/v6"

	cn "github.com/haproxytech/dataplaneapi/client-native"
	"github.com/haproxytech/dataplaneapi/handlers/middleware"
	"github.com/haproxytech/dataplaneapi/handlers/respond"
	"github.com/haproxytech/dataplaneapi/reload_agent"
)

// RegisterRouter registers all raw configuration routes onto r using spec-based request validation
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

// HandlerImpl implements ServerInterface for HAProxy raw configuration.
type HandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent reload_agent.IReloadAgent
}

func (h *HandlerImpl) GetHAProxyConfiguration(w http.ResponseWriter, r *http.Request, params GetHAProxyConfigurationParams) {
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	_, clusterVersion, md5Hash, data, err := cfg.GetRawConfigurationWithClusterData(params.TransactionId, int64(params.Version))
	if err != nil {
		respond.Error(w, err)
		return
	}
	writePlainText(w, http.StatusOK, clusterVersion, md5Hash, data)
}

func (h *HandlerImpl) PostHAProxyConfiguration(w http.ResponseWriter, r *http.Request, params PostHAProxyConfigurationParams) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		respond.BadRequest(w, err.Error())
		return
	}
	data := string(body)

	if len(data) > 0 && !strings.ContainsRune(data, '\n') {
		respond.BadRequest(w, "invalid configuration: no newline character found")
		return
	}

	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.PostRawConfiguration(&data, int64(params.Version), params.SkipVersion, params.OnlyValidate); err != nil {
		respond.Error(w, err)
		return
	}

	_, clusterVersion, md5Hash, result, err := cfg.GetRawConfigurationWithClusterData("", 0)
	if err != nil {
		respond.Error(w, err)
		return
	}

	if params.OnlyValidate {
		writePlainText(w, http.StatusAccepted, clusterVersion, md5Hash, result)
		return
	}

	if params.SkipReload {
		if params.XRuntimeActions != "" {
			if err = executeRuntimeActions(params.XRuntimeActions, h.Client); err != nil {
				respond.Error(w, err)
				return
			}
		}
		writePlainText(w, http.StatusCreated, clusterVersion, md5Hash, result)
		return
	}

	callbackNeeded, reconfigureFunc, err := cn.ReconfigureRuntime(h.Client)
	if err != nil {
		respond.Error(w, err)
		return
	}

	if params.ForceReload {
		if callbackNeeded {
			err = h.ReloadAgent.ForceReloadWithCallback(reconfigureFunc)
		} else {
			err = h.ReloadAgent.ForceReload()
		}
		if err != nil {
			respond.Error(w, err)
			return
		}
		writePlainText(w, http.StatusCreated, clusterVersion, md5Hash, result)
		return
	}

	var rID string
	if callbackNeeded {
		rID = h.ReloadAgent.ReloadWithCallback(reconfigureFunc)
	} else {
		rID = h.ReloadAgent.Reload()
	}
	w.Header().Set(respond.ReloadIDHeader, rID)
	writePlainText(w, http.StatusAccepted, clusterVersion, md5Hash, result)
}

func writePlainText(w http.ResponseWriter, status int, clusterVersion int64, md5Hash, data string) {
	w.Header().Set("Content-Type", "text/plain")
	if clusterVersion != 0 {
		w.Header().Set("Cluster-Version", strconv.FormatInt(clusterVersion, 10))
	}
	if md5Hash != "" {
		w.Header().Set("Configuration-Checksum", md5Hash)
	}
	w.WriteHeader(status)
	respond.Write(w, []byte(data))
}

func executeRuntimeActions(actionsStr string, client client_native.HAProxyClient) error {
	rt, err := client.Runtime()
	if err != nil {
		return err
	}
	for a := range strings.SplitSeq(actionsStr, ";") {
		params := strings.Split(a, " ")
		if len(params) == 0 {
			continue
		}
		action := params[0]
		if action == "" {
			continue
		}
		switch action {
		case "SetFrontendMaxConn":
			if len(params) > 2 {
				maxConn, err := strconv.ParseInt(params[2], 10, 64)
				if err != nil {
					return fmt.Errorf("cannot execute %s: %s", action, err.Error())
				}
				if err = rt.SetFrontendMaxConn(params[1], int(maxConn)); err != nil {
					return fmt.Errorf("cannot execute %s: %s", action, err.Error())
				}
			} else {
				return fmt.Errorf("cannot execute %s: not enough parameters", action)
			}
		case "SetServerWeight":
			if len(params) > 3 {
				if err = rt.SetServerWeight(params[1], params[2], params[3]); err != nil {
					return fmt.Errorf("cannot execute %s: %s", action, err.Error())
				}
			} else {
				return fmt.Errorf("cannot execute %s: not enough parameters", action)
			}
		case "SetServerCheckPort":
			if len(params) > 3 {
				port, err := strconv.ParseInt(params[3], 10, 64)
				if err != nil {
					return fmt.Errorf("cannot execute %s: %s", action, err.Error())
				}
				if err = rt.SetServerCheckPort(params[1], params[2], int(port)); err != nil {
					return fmt.Errorf("cannot execute %s: %s", action, err.Error())
				}
			} else {
				return fmt.Errorf("cannot execute %s: not enough parameters", action)
			}
		case "SetServerAddr":
			if len(params) > 4 {
				port, err := strconv.ParseInt(params[4], 10, 64)
				if err != nil {
					return fmt.Errorf("cannot execute %s: %s", action, err.Error())
				}
				if err = rt.SetServerAddr(params[1], params[2], params[3], int(port)); err != nil {
					return fmt.Errorf("cannot execute %s: %s", action, err.Error())
				}
			} else {
				return fmt.Errorf("cannot execute %s: not enough parameters", action)
			}
		case "SetServerState":
			if len(params) > 3 {
				if err = rt.SetServerState(params[1], params[2], params[3]); err != nil {
					return fmt.Errorf("cannot execute %s: %s", action, err.Error())
				}
			} else {
				return fmt.Errorf("cannot execute %s: not enough parameters", action)
			}
		case "EnableAgentCheck":
			if len(params) > 2 {
				if err = rt.EnableAgentCheck(params[1], params[2]); err != nil {
					return fmt.Errorf("cannot execute %s: %s", action, err.Error())
				}
			} else {
				return fmt.Errorf("cannot execute %s: not enough parameters", action)
			}
		case "DisableAgentCheck":
			if len(params) > 2 {
				if err = rt.DisableAgentCheck(params[1], params[2]); err != nil {
					return fmt.Errorf("cannot execute %s: %s", action, err.Error())
				}
			} else {
				return fmt.Errorf("cannot execute %s: not enough parameters", action)
			}
		case "SetServerAgentAddr":
			if len(params) > 3 {
				if err = rt.SetServerAgentAddr(params[1], params[2], params[3]); err != nil {
					return fmt.Errorf("cannot execute %s: %s", action, err.Error())
				}
			} else {
				return fmt.Errorf("cannot execute %s: not enough parameters", action)
			}
		case "SetServerAgentSend":
			if len(params) > 3 {
				if err = rt.SetServerAgentSend(params[1], params[2], params[3]); err != nil {
					return fmt.Errorf("cannot execute %s: %s", action, err.Error())
				}
			} else {
				return fmt.Errorf("cannot execute %s: not enough parameters", action)
			}
		}
	}
	return nil
}
