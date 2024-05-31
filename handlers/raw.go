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
//

package handlers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	client_native "github.com/haproxytech/client-native/v6"
	"github.com/haproxytech/client-native/v6/models"
	cn "github.com/haproxytech/dataplaneapi/client-native"
	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/operations/configuration"
)

// GetRawConfigurationHandlerImpl implementation of the GetHAProxyConfigurationHandler interface
type GetRawConfigurationHandlerImpl struct {
	Client client_native.HAProxyClient
}

// PostRawConfigurationHandlerImpl implementation of the PostHAProxyConfigurationHandler interface
type PostRawConfigurationHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// Handle executing the request and returning a response
func (h *GetRawConfigurationHandlerImpl) Handle(params configuration.GetHAProxyConfigurationParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	v := int64(0)
	if params.Version != nil {
		v = *params.Version
	}

	cfg, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return configuration.NewGetConfigurationVersionDefault(int(*e.Code)).WithPayload(e)
	}

	_, clusterVersion, md5Hash, data, err := cfg.GetRawConfigurationWithClusterData(t, v)
	if err != nil {
		e := misc.HandleError(err)
		return configuration.NewGetHAProxyConfigurationDefault(int(*e.Code)).WithPayload(e)
	}
	cVersion := ""
	if clusterVersion != 0 {
		cVersion = strconv.FormatInt(clusterVersion, 10)
	}
	return configuration.NewGetHAProxyConfigurationOK().WithPayload(&configuration.GetHAProxyConfigurationOKBody{Data: &data}).WithClusterVersion(cVersion).WithConfigurationChecksum(md5Hash)
}

// Handle executing the request and returning a response
func (h *PostRawConfigurationHandlerImpl) Handle(params configuration.PostHAProxyConfigurationParams, principal interface{}) middleware.Responder {
	v := int64(0)
	if params.Version != nil {
		v = *params.Version
	}
	skipReload := false
	if params.SkipReload != nil {
		skipReload = *params.SkipReload
	}
	skipVersion := false
	if params.SkipVersion != nil {
		skipVersion = *params.SkipVersion
	}
	forceReload := false
	if params.ForceReload != nil {
		forceReload = *params.ForceReload
	}
	onlyValidate := false
	if params.OnlyValidate != nil {
		onlyValidate = *params.OnlyValidate
	}

	// Check for a common error where the user ran `curl -d @haproxy.cfg` without
	// converting the file to json, which removes all the \n from the configuration.
	if len(params.Data) > 0 && !strings.ContainsRune(params.Data, '\n') {
		code := misc.ErrHTTPBadRequest
		msg := "invalid configuration: no newline character found"
		e := &models.Error{Code: &code, Message: &msg}
		return configuration.NewPostHAProxyConfigurationBadRequest().WithPayload(e)
	}

	cfg, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return configuration.NewPostHAProxyConfigurationDefault(int(*e.Code)).WithPayload(e)
	}

	err = cfg.PostRawConfiguration(&params.Data, v, skipVersion, onlyValidate)
	if err != nil {
		e := misc.HandleError(err)
		return configuration.NewPostHAProxyConfigurationDefault(int(*e.Code)).WithPayload(e)
	}

	_, clusterVersion, md5Hash, data, err := cfg.GetRawConfigurationWithClusterData("", 0)
	if err != nil {
		e := misc.HandleError(err)
		return configuration.NewPostHAProxyConfigurationDefault(int(*e.Code)).WithPayload(e)
	}

	cVersion := ""
	if clusterVersion != 0 {
		cVersion = strconv.FormatInt(clusterVersion, 10)
	}
	if onlyValidate {
		// return here without reloading, since config is only validated.
		return configuration.NewPostHAProxyConfigurationAccepted().WithPayload(data).WithClusterVersion(cVersion).WithConfigurationChecksum(md5Hash)
	}
	if skipReload {
		if params.XRuntimeActions != nil {
			if err = executeRuntimeActions(*params.XRuntimeActions, h.Client); err != nil {
				e := misc.HandleError(err)
				return configuration.NewPostHAProxyConfigurationDefault(int(*e.Code)).WithPayload(e)
			}
		}
		return configuration.NewPostHAProxyConfigurationCreated().WithPayload(data).WithClusterVersion(cVersion).WithConfigurationChecksum(md5Hash)
	}
	if forceReload {
		var callbackNeeded bool
		var reconfigureFunc func()
		callbackNeeded, reconfigureFunc, err = cn.ReconfigureRuntime(h.Client)
		if err != nil {
			e := misc.HandleError(err)
			return configuration.NewPostHAProxyConfigurationDefault(int(*e.Code)).WithPayload(e)
		}
		if callbackNeeded {
			err = h.ReloadAgent.ForceReloadWithCallback(reconfigureFunc)
		} else {
			err = h.ReloadAgent.ForceReload()
		}
		if err != nil {
			e := misc.HandleError(err)
			return configuration.NewPostHAProxyConfigurationDefault(int(*e.Code)).WithPayload(e)
		}
		return configuration.NewPostHAProxyConfigurationCreated().WithPayload(data).WithClusterVersion(cVersion).WithConfigurationChecksum(md5Hash)
	}
	callbackNeeded, reconfigureFunc, err := cn.ReconfigureRuntime(h.Client)
	if err != nil {
		e := misc.HandleError(err)
		return configuration.NewPostHAProxyConfigurationDefault(int(*e.Code)).WithPayload(e)
	}

	var rID string
	if callbackNeeded {
		rID = h.ReloadAgent.ReloadWithCallback(reconfigureFunc)
	} else {
		rID = h.ReloadAgent.Reload()
	}

	return configuration.NewPostHAProxyConfigurationAccepted().WithReloadID(rID).WithPayload(data).WithClusterVersion(cVersion).WithConfigurationChecksum(md5Hash)
}

func executeRuntimeActions(actionsStr string, client client_native.HAProxyClient) error {
	runtime, err := client.Runtime()
	if err != nil {
		return err
	}
	actions := strings.Split(actionsStr, ";")
	for _, a := range actions {
		params := strings.Split(a, " ")
		if len(params) == 0 {
			continue
		}
		action := params[0]
		switch action {
		case "SetFrontendMaxConn":
			if len(params) > 2 {
				fName := params[1]
				maxConn, err := strconv.ParseInt(params[2], 10, 64)
				if err != nil {
					return fmt.Errorf("cannot execute %s: %s", action, err.Error())
				}
				if err := runtime.SetFrontendMaxConn(fName, int(maxConn)); err != nil {
					return fmt.Errorf("cannot execute %s: %s", action, err.Error())
				}
			} else {
				return fmt.Errorf("cannot execute %s: not enough parameters", action)
			}
		case "SetServerWeight":
			if len(params) > 3 {
				backend := params[1]
				server := params[2]
				weight := params[3]
				if err := runtime.SetServerWeight(backend, server, weight); err != nil {
					return fmt.Errorf("cannot execute %s: %s", action, err.Error())
				}
			} else {
				return fmt.Errorf("cannot execute %s: not enough parameters", action)
			}
		case "SetServerCheckPort":
			if len(params) > 3 {
				backend := params[1]
				server := params[2]
				port, err := strconv.ParseInt(params[3], 10, 64)
				if err != nil {
					return fmt.Errorf("cannot execute %s: %s", action, err.Error())
				}
				if err := runtime.SetServerCheckPort(backend, server, int(port)); err != nil {
					return fmt.Errorf("cannot execute %s: %s", action, err.Error())
				}
			} else {
				return fmt.Errorf("cannot execute %s: not enough parameters", action)
			}
		case "SetServerAddr":
			if len(params) > 4 {
				backend := params[1]
				server := params[2]
				ip := params[3]
				port, err := strconv.ParseInt(params[4], 10, 64)
				if err != nil {
					return fmt.Errorf("cannot execute %s: %s", action, err.Error())
				}
				if err := runtime.SetServerAddr(backend, server, ip, int(port)); err != nil {
					return fmt.Errorf("cannot execute %s: %s", action, err.Error())
				}
			} else {
				return fmt.Errorf("cannot execute %s: not enough parameters", action)
			}
		case "SetServerState":
			if len(params) > 3 {
				backend := params[1]
				server := params[2]
				state := params[3]
				if err := runtime.SetServerState(backend, server, state); err != nil {
					return fmt.Errorf("cannot execute %s: %s", action, err.Error())
				}
			} else {
				return fmt.Errorf("cannot execute %s: not enough parameters", action)
			}
		case "EnableAgentCheck":
			if len(params) > 2 {
				backend := params[1]
				server := params[2]
				if err := runtime.EnableAgentCheck(backend, server); err != nil {
					return fmt.Errorf("cannot execute %s: %s", action, err.Error())
				}
			} else {
				return fmt.Errorf("cannot execute %s: not enough parameters", action)
			}
		case "DisableAgentCheck":
			if len(params) > 2 {
				backend := params[1]
				server := params[2]
				if err := runtime.DisableAgentCheck(backend, server); err != nil {
					return fmt.Errorf("cannot execute %s: %s", action, err.Error())
				}
			} else {
				return fmt.Errorf("cannot execute %s: not enough parameters", action)
			}
		case "SetServerAgentAddr":
			if len(params) > 3 {
				backend := params[1]
				server := params[2]
				addr := params[3]
				if err := runtime.SetServerAgentAddr(backend, server, addr); err != nil {
					return fmt.Errorf("cannot execute %s: %s", action, err.Error())
				}
			} else {
				return fmt.Errorf("cannot execute %s: not enough parameters", action)
			}
		case "SetServerAgentSend":
			if len(params) > 3 {
				backend := params[1]
				server := params[2]
				send := params[3]
				if err := runtime.SetServerAgentSend(backend, server, send); err != nil {
					return fmt.Errorf("cannot execute %s: %s", action, err.Error())
				}
			} else {
				return fmt.Errorf("cannot execute %s: not enough parameters", action)
			}
		}
	}
	return nil
}
